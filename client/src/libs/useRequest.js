import { useState, useEffect } from 'react';
import axios from 'axios';
import isAPIError from '@utils/isAPIError';

export default function useRequest(url) {
  const [data, setData] = useState(undefined);
  const [statusCode, setStatusCode] = useState(200);
  const [error, setError] = useState('');

  const refresh = () => {
    axios
      .get(url)
      .then(response => {
        setData(response.data.data);
      })
      .catch(err => {
        if (isAPIError(err)) {
          setError(err.response.data.errors[0].message);
        }
        setStatusCode(err.response ? err.response.status : 500);
      });
  };

  useEffect(refresh, [url]);

  return {
    data,
    statusCode,
    error,
    setData,
    setStatusCode,
    setError,
    refresh
  };
}
