import React from 'react';
import { isNil } from 'lodash';
import { string, oneOfType, number } from 'prop-types';

import { makeStyles } from '@material-ui/core/styles';
import { Typography } from '@material-ui/core';
import AppLayout from '@common/AppLayout/AppLayout';
import ResponsiveImage from '@common/ResponsiveImage/ResponsiveImage';
import errorIcon from './error.svg';

export const STATUS_CODES = {
  400: 'Bad Request',
  404: 'This page could not be found',
  405: 'Method Not Allowed',
  500: 'Internal Server Error'
};

const DEFAULT_MESSAGE = 'An unexpected error has occurred';

const useStyles = makeStyles(theme => ({
  image: {
    maxHeight: '25vh'
  },
  container: {
    '& > *:not(:last-child)': {
      marginBottom: theme.spacing(2)
    }
  },
  appLayout: {
    display: 'flex',
    justifyContent: 'center',
    flexDirection: 'column',
    textAlign: 'center'
  }
}));

export default function ErrorPage({ error, statusCode }) {
  const classes = useStyles();
  return (
    <AppLayout className={classes.appLayout}>
      <ResponsiveImage className={classes.image} src={errorIcon} alt="error" />
      <Typography variant="h2" component="h1">
        {statusCode}
      </Typography>
      <Typography variant="h5" component="h4">
        {!error
          ? isNil(STATUS_CODES[statusCode])
            ? DEFAULT_MESSAGE
            : STATUS_CODES[statusCode]
          : error}
      </Typography>
    </AppLayout>
  );
}

ErrorPage.defaultProps = {
  statusCode: 404
};

ErrorPage.propTypes = {
  statusCode: oneOfType([number, string]).isRequired,
  error: string
};
