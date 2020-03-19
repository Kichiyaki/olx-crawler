import React, { useState, useEffect } from 'react';
import axios from 'axios';
import useSnackbar, { SEVERITY } from '@libs/useSnackbar';
import { SUGGESTIONS } from '@config/api_routes';
import isAPIError from '@utils/isAPIError';
import { getSuggestionsReqParams } from './utils';

import InfiniteScroll from 'react-infinite-scroll-component';
import {
  Grid,
  Container,
  Typography,
  Snackbar,
  Button
} from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import AppLayout, { CONTAINER_ID } from '@common/AppLayout/AppLayout';
import Spinner from '@common/Spinner/Spinner';
import Suggestion from './components/Suggestion/Suggestion';

export default function MainPage() {
  const [suggestions, setSuggestions] = useState(undefined);
  const [selected, setSelected] = useState([]);
  const [statusCode, setStatusCode] = useState(200);
  const [error, setError] = useState('');
  const {
    setSeverity,
    setMessage,
    message,
    alertProps,
    snackbarProps
  } = useSnackbar({
    anchorOrigin: { vertical: 'top', horizontal: 'right' }
  });

  useEffect(() => {
    axios
      .get(SUGGESTIONS.BROWSE + '?' + getSuggestionsReqParams())
      .then(response => {
        setSuggestions(response.data.data);
      })
      .catch(error => {
        if (isAPIError(error)) {
          setError(error.response.data.data.errors[0].message);
        }
        if (error.response) {
          setStatusCode(error.response.status);
        }
      });
  }, []);

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
        setMessage(`Pomyślnie usunięto ${response.data.data.length} sugestie.`);
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
        <Container>
          <Typography gutterBottom align="center" variant="h2" component="h1">
            Sugestie
          </Typography>
          <InfiniteScroll
            hasMore={suggestions.items.length !== suggestions.total}
            next={loadMore}
            dataLength={suggestions.items.length}
            loader={<Spinner mt={2} />}
            style={{ overflow: 'visible' }}
            scrollableTarget={CONTAINER_ID}
          >
            <Grid container spacing={2}>
              {suggestions.items.length === 0 ? (
                <Typography variant="h3" component="h2">
                  Brak sugestii
                </Typography>
              ) : (
                suggestions.items.map(suggestion => (
                  <Grid key={suggestion.id} xs={4} item>
                    <Suggestion
                      onSelect={createSelectHandler(suggestion.id)}
                      selected={selected.some(id => id === suggestion.id)}
                      data={suggestion}
                    />
                  </Grid>
                ))
              )}
            </Grid>
          </InfiniteScroll>
          <Snackbar
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right'
            }}
            ClickAwayListenerProps={{ mouseEvent: false }}
            open={selected.length > 0}
            onClose={handleSnackbarClose}
            message={`Wybrano ${selected.length} sugestie.`}
            action={
              <>
                <Button onClick={handleDelete} color="secondary" size="small">
                  Usuń
                </Button>
                <Button
                  color="secondary"
                  onClick={handleSnackbarClose}
                  size="small"
                >
                  Anuluj
                </Button>
              </>
            }
          />
        </Container>
      )}
      <Snackbar {...snackbarProps}>
        <Alert {...alertProps}>{message}</Alert>
      </Snackbar>
    </AppLayout>
  );
}
