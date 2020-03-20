import React, { useEffect, useState } from 'react';
import axios from 'axios';
import i18n from '@libs/i18n/i18n';
import { CONFIG } from '@config/api_routes';
import isAPIError from '@utils/isAPIError';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import Spinner from '@common/Spinner/Spinner';
import Context from './config.context';

export default function Provider({ children }) {
  const [config, setConfig] = useState(undefined);
  const [err, setErr] = useState('');
  const [statusCode, setStatusCode] = useState(200);

  useEffect(() => {
    reload();
  }, []);

  const reload = () => {
    axios
      .get(CONFIG.READ)
      .then(response => {
        i18n.changeLanguage(response.data.data.lang, () => {
          setConfig(response.data.data);
        });
      })
      .catch(err => {
        if (isAPIError(err)) {
          setErr(err.response.data.errors[0].message);
        }
        setStatusCode(err.response ? err.response.status : 500);
      });
  };

  const update = changes => {
    if (changes.lang !== config.lang) {
      i18n.changeLanguage(changes.lang);
    }
    setConfig({
      ...config,
      ...changes
    });
  };

  const loading = !config;

  if (statusCode !== 200 || err) {
    return <ErrorPage error={err} statusCode={statusCode} />;
  }

  if (loading) {
    return <Spinner height="100vh" />;
  }

  return (
    <Context.Provider value={{ config, reload, update }}>
      {children}
    </Context.Provider>
  );
}
