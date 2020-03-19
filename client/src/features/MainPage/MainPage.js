import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { CONFIG, SUGGESTIONS } from '@config/api_routes';
import isAPIError from '@utils/isAPIError';
import { getSuggestionsReqParams } from './utils';

import InfiniteScroll from 'react-infinite-scroll-component';
import { Grid, Container, Typography } from '@material-ui/core';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import AppLayout, { CONTAINER_ID } from '@common/AppLayout/AppLayout';
import Spinner from '@common/Spinner/Spinner';
import Suggestion from './components/Suggestion/Suggestion';

export default function MainPage() {
  const [config, setConfig] = useState(undefined);
  const [suggestions, setSuggestions] = useState(undefined);
  const [statusCode, setStatusCode] = useState(200);
  const [error, setError] = useState('');

  useEffect(() => {
    Promise.all([
      axios.get(CONFIG.READ),
      axios.get(SUGGESTIONS.BROWSE + '?' + getSuggestionsReqParams())
    ])
      .then(responses => {
        setConfig(responses[0].data.data);
        setSuggestions(responses[1].data.data);
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

  if (error || statusCode !== 200) {
    return <ErrorPage statusCode={statusCode} error={error} />;
  }

  const loading = !suggestions || !config;

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
              {suggestions.items.map(suggestion => (
                <Grid key={suggestion.id} xs={4} item>
                  <Suggestion data={suggestion} />
                </Grid>
              ))}
            </Grid>
          </InfiniteScroll>
        </Container>
      )}
    </AppLayout>
  );
}
