import React from 'react';
import axios from 'axios';
import { func, object } from 'prop-types';
import { useFormik } from 'formik';
import { SEVERITY } from '@libs/useSnackbar';
import { CONFIG } from '@config/api_routes';
import isAPIError from '@utils/isAPIError';

import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Checkbox,
  FormControlLabel
} from '@material-ui/core';

function CollyForm({ update, config, setMessage, setSeverity, t }) {
  const {
    handleBlur,
    handleChange,
    handleSubmit,
    isSubmitting,
    values,
    touched,
    errors,
    handleReset
  } = useFormik({
    initialValues: config.discord_notifications
      ? config.discord_notifications
      : {
          channel_id: '',
          token: '',
          enabled: false
        },
    onSubmit: async (values, { setSubmitting }) => {
      try {
        const response = await axios.patch(
          CONFIG.UPDATE.DISCORD_NOTIFICATIONS,
          values
        );
        update(response.data.data);
        setSeverity(SEVERITY.SUCCESS);
        setMessage(t('success'));
      } catch (error) {
        setSeverity(SEVERITY.ERROR);
        if (isAPIError(error)) {
          setMessage(error.response.data.errors[0].message);
        } else {
          setMessage(t('error'));
        }
      }
      setSubmitting(false);
    },
    validate: ({ token, channel_id }) => {
      const errors = {};
      if (token && !token.match(/[MN][A-Za-z\d]{23}\.[\w-]{6}\.[\w-]{27}/g)) {
        errors.token = t('discordNotificationsForm.error.invalidToken');
      }
      if (channel_id && isNaN(parseInt(channel_id))) {
        errors.channel_id = t(
          'discordNotificationsForm.error.invalidChannelID'
        );
      }
      return errors;
    }
  });

  return (
    <Card>
      <CardContent>
        <Typography variant="h3" gutterBottom component="h2">
          {t('discordNotificationsForm.title')}
        </Typography>
        <form noValidate onSubmit={handleSubmit}>
          <TextField
            name={'token'}
            label={t('discordNotificationsForm.inputLabel.token')}
            value={values.token}
            onChange={handleChange}
            onBlur={handleBlur}
            fullWidth
            error={touched.token && !!errors.token}
            helperText={touched.token && errors.token}
          />
          <TextField
            name={'channel_id'}
            label={t('discordNotificationsForm.inputLabel.channelID')}
            value={values.channel_id}
            onChange={handleChange}
            onBlur={handleBlur}
            fullWidth
            error={touched.channel_id && !!errors.channel_id}
            helperText={touched.channel_id && errors.channel_id}
          />
          <FormControlLabel
            control={
              <Checkbox
                checked={values.enabled}
                onChange={handleChange}
                name="enabled"
              />
            }
            label={t('discordNotificationsForm.inputLabel.enabled')}
          />
          <div>
            <Button type="submit" disabled={isSubmitting}>
              {t('save')}
            </Button>
            <Button type="button" disabled={isSubmitting} onClick={handleReset}>
              {t('reset')}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

CollyForm.propTypes = {
  update: func.isRequired,
  setMessage: func.isRequired,
  setSeverity: func.isRequired,
  t: func.isRequired,
  config: object.isRequired
};

export default CollyForm;
