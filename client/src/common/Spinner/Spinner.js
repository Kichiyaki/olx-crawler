import React from 'react';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { Box } from '@material-ui/core';

export default function Spinner({ spinnerProps, ...rest }) {
  return (
    <Box {...rest}>
      <ScaleLoader {...spinnerProps} />
    </Box>
  );
}

Spinner.defaultProps = {
  spinnerProps: {},
  display: 'flex',
  height: '100%',
  alignItems: 'center',
  justifyContent: 'center'
};
