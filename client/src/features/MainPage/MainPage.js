import React, { useState } from 'react';
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import useRequest from '@libs/useRequest';
import useSnackbar, { SEVERITY } from '@libs/useSnackbar';
import { SUGGESTIONS } from '@config/api_routes';
import { MAIN_PAGE } from '@config/namespaces';
import isAPIError from '@utils/isAPIError';
import { getSuggestionsReqParams } from './utils';

import InfiniteScroll from 'react-infinite-scroll-component';
import {
  Grid,
  Typography,
  Snackbar,
  Button,
  Container
} from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import AppLayout, { CONTAINER_ID } from '@common/AppLayout/AppLayout';
import Spinner from '@common/Spinner/Spinner';
import Suggestion from './components/Suggestion/Suggestion';

export default function MainPage() {
  const [selected, setSelected] = useState([]);
  const {
    setSeverity,
    setMessage,
    message,
    alertProps,
    snackbarProps
  } = useSnackbar({
    anchorOrigin: { vertical: 'top', horizontal: 'right' }
  });
  const { t } = useTranslation(MAIN_PAGE);
  const {
    data: suggestions,
    setData: setSuggestions,
    statusCode,
    error
  } = useRequest(SUGGESTIONS.BROWSE + '?' + getSuggestionsReqParams());

  const loadMore = async () => {
    try {
      const response = await axios.get(
        SUGGESTIONS.BROWSE +
          '?' +
          getSuggestionsReqParams(suggestions.items.length)
      );
      setSuggestions({
        items: [
          ...suggestions.items,
          ...response.data.data.items.filter(
            item =>
              !suggestions.items.some(otherItem => otherItem.id === item.id)
          )
        ],
        total: response.data.data.total
      });
    } catch (error) {}
  };

  const createSelectHandler = id => () => {
    if (selected.some(otherID => otherID === id)) {
      setSelected(selected.filter(otherID => otherID !== id));
    } else {
      setSelected([...selected, id]);
    }
  };

  const handleDelete = async () => {
    try {
      const response = await axios.delete(
        SUGGESTIONS.DELETE + '?id=' + selected.join(',')
      );
      if (Array.isArray(response.data.data)) {
        setSuggestions({
          total: suggestions.total - response.data.data.length,
          items: suggestions.items.filter(
            item =>
              !response.data.data.some(otherItem => otherItem.id === item.id)
          )
        });
        setSelected([]);
        setSeverity(SEVERITY.SUCCESS);
        setMessage(t('deleted', { count: response.data.data.length }));
      }
    } catch (error) {
      if (isAPIError(error)) {
        setSeverity(SEVERITY.ERROR);
        setMessage(error.response.data.errors[0].message);
      }
    }
  };

  const handleSnackbarClose = () => {
    setSelected([]);
  };

  if (error || statusCode !== 200) {
    return <ErrorPage statusCode={statusCode} error={error} />;
  }

  const loading = !suggestions;

  return (
    <AppLayout>
      {loading && <Spinner />}
      {!loading && (
        <div>
          <Typography gutterBottom align="center" variant="h2" component="h1">
            {t('title')}
          </Typography>
          <InfiniteScroll
            hasMore={suggestions.items.length !== suggestions.total}
            next={loadMore}
            dataLength={suggestions.items.length}
            loader={<Spinner mt={2} />}
            style={{ overflow: 'visible' }}
            scrollableTarget={CONTAINER_ID}
          >
            <Container>
              <Grid container spacing={2}>
                {suggestions.items.length === 0 ? (
                  <Typography variant="h3" component="h2">
                    {t('emptyArray')}
                  </Typography>
                ) : (
                  suggestions.items.map(suggestion => (
                    <Grid key={suggestion.id} xs={12} sm={6} lg={4} item>
                      <Suggestion
                        t={t}
                        onSelect={createSelectHandler(suggestion.id)}
                        selected={selected.some(id => id === suggestion.id)}
                        data={suggestion}
                      />
                    </Grid>
                  ))
                )}
              </Grid>
            </Container>
          </InfiniteScroll>
          <Snackbar
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right'
            }}
            ClickAwayListenerProps={{ mouseEvent: false }}
            open={selected.length > 0}
            onClose={handleSnackbarClose}
            message={t('selected', { count: selected.length })}
            action={
              <>
                <Button onClick={handleDelete} color="secondary" size="small">
                  {t('delete')}
                </Button>
                <Button
                  color="secondary"
                  onClick={handleSnackbarClose}
                  size="small"
                >
                  {t('cancel')}
                </Button>
              </>
            }
          />
        </div>
      )}
      <Snackbar {...snackbarProps}>
        <Alert {...alertProps}>{message}</Alert>
      </Snackbar>
    </AppLayout>
  );
}
