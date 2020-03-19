import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { CONFIG, SUGGESTIONS } from '@config/api_routes';
import prepareParams from '@utils/prepareParams';
import { SUGGESTIONS_REQ_PARAMS } from './constants';

import InfiniteScroll from 'react-infinite-scroll-component';
import { Grid, Container } from '@material-ui/core';
import AppLayout from '@common/AppLayout/AppLayout';
import Suggestion from './components/Suggestion/Suggestion';

export default function MainPage() {
  const [config, setConfig] = useState(undefined);
  const [suggestions, setSuggestions] = useState(undefined);

  useEffect(() => {
    Promise.all([
      axios.get(CONFIG.READ),
      axios.get(SUGGESTIONS.BROWSE + '?' + SUGGESTIONS_REQ_PARAMS)
    ]).then(responses => {
      setConfig(responses[0].data.data);
      setSuggestions(responses[1].data.data);
    });
  }, []);
  const loading = !suggestions || !config;
  console.log(loading);

  const loadMore = async () => {
    const response = await axios.get(
      SUGGESTIONS.BROWSE +
        '?' +
        prepareParams({ limit: 10, offset: suggestions.items.length }, {}, [
          'limit',
          'offset'
        ])
    );
    setSuggestions({
      items: [
        ...suggestions.items,
        ...response.data.data.items.filter(
          item => !suggestions.items.some(otherItem => otherItem.id === item.id)
        )
      ],
      total: response.data.data.total
    });
  };

  return (
    <AppLayout>
      <Container>
        <Grid container spacing={2}>
          <Grid xs={8} item>
            {!loading && (
              <InfiniteScroll
                hasMore={suggestions.items.length !== suggestions.total}
                next={loadMore}
                dataLength={suggestions.items.length}
                loader={
                  <div className="loader" key={0}>
                    Loading ...
                  </div>
                }
                style={{ overflow: 'visible' }}
                scrollableTarget="main-content"
              >
                {suggestions.items.map(suggestion => (
                  <Suggestion key={suggestion.id} data={suggestion} />
                ))}
              </InfiniteScroll>
            )}
          </Grid>
          <Grid xs={4} item>
            Wersja/Aktywne obserwacje/Ustawienia
          </Grid>
        </Grid>
      </Container>
    </AppLayout>
  );
}
