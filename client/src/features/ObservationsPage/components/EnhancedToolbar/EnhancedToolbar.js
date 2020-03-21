import React from 'react';
import { string, func } from 'prop-types';

import { Toolbar, Typography, Tooltip, IconButton } from '@material-ui/core';
import { Add as AddIcon } from '@material-ui/icons';

export default function EnhancedToolbar({ name, onClickAddObservation, t }) {
  return (
    <Toolbar>
      <Typography color="inherit" variant="h2" component="h1">
        {name}
      </Typography>
      <div style={{ flex: '1' }}></div>
      <div>
        <Tooltip title={t('enhancedToolbar.addObservation')}>
          <IconButton
            onClick={onClickAddObservation}
            aria-label="add observation"
          >
            <AddIcon />
          </IconButton>
        </Tooltip>
      </div>
    </Toolbar>
  );
}

EnhancedToolbar.propTypes = {
  name: string.isRequired,
  onClickAddObservation: func.isRequired,
  t: func.isRequired
};
