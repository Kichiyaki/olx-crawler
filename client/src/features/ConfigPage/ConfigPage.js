import React from 'react';
import { useTranslation } from 'react-i18next';
import useSnackbar from '@libs/useSnackbar';
import useConfig from '@libs/config/useConfig';
import { CONFIG_PAGE } from '@config/namespaces';

import { Container, Grid, Snackbar } from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import AppLayout from '@common/AppLayout/AppLayout';
import ProxyForm from './components/ProxyForm/ProxyForm';
import CollyForm from './components/CollyForm/CollyForm';
import DiscordNotificationsForm from './components/DiscordNotificationsForm/DiscordNotificationsForm';
import AppSettingsForm from './components/AppSettingsForm/AppSettingsForm';

function ConfigPage() {
  const configBag = useConfig();
  const { message, alertProps, snackbarProps, ...snackbarBag } = useSnackbar({
    anchorOrigin: { vertical: 'top', horizontal: 'right' }
  });
  const { t } = useTranslation(CONFIG_PAGE);
  const props = {
    ...configBag,
    ...snackbarBag,
    t
  };
  return (
    <AppLayout>
      <Container>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <ProxyForm {...props} />
          </Grid>
          <Grid item xs={12} md={6}>
            <CollyForm {...props} />
          </Grid>
          <Grid item xs={12} md={6}>
            <DiscordNotificationsForm {...props} />
          </Grid>
          <Grid item xs={12} md={6}>
            <AppSettingsForm {...props} />
          </Grid>
        </Grid>
      </Container>
      <Snackbar {...snackbarProps}>
        <Alert {...alertProps}>{message}</Alert>
      </Snackbar>
    </AppLayout>
  );
}

export default ConfigPage;
