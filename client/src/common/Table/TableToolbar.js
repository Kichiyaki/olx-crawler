import React from 'react';
import { string } from 'prop-types';

import { Toolbar, Typography } from '@material-ui/core';

export default function TableToolbar({ name }) {
  return (
    <Toolbar>
      <Typography color="inherit" variant="h2" component="h1">
        {name}
      </Typography>
    </Toolbar>
  );
}

TableToolbar.propTypes = {
  name: string.isRequired
};
