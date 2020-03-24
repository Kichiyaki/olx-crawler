import React, { useState } from 'react';
import axios from 'axios';
import { func, object } from 'prop-types';
import { SEVERITY } from '@libs/useSnackbar';
import { CONFIG } from '@config/api_routes';
import { DEFAULT_LANGUAGE, AVAILABLE_LANGUAGES } from '@config/application';
import isAPIError from '@utils/isAPIError';

import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Checkbox,
  FormControlLabel,
  MenuItem
} from '@material-ui/core';

function AppSettingsForm({ update, config, setMessage, setSeverity, t }) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [lang, setLang] = useState(config.lang || DEFAULT_LANGUAGE);
  const [debug, setDebug] = useState(config.debug || false);

  const handleSelectLanguageForm = async e => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      const response = await axios.patch(CONFIG.UPDATE.LANG + '/' + lang);
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
    setIsSubmitting(false);
  };

  const handleDebugForm = async e => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      const response = await axios.patch(CONFIG.UPDATE.DEBUG + '/' + debug);
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
    setIsSubmitting(false);
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h3" gutterBottom component="h2">
          {t('appSettingsForm.title')}
        </Typography>
        <Typography variant="h4" gutterBottom component="h3">
          {t('appSettingsForm.subtitle.1')}
        </Typography>
        <form noValidate onSubmit={handleSelectLanguageForm}>
          <TextField
            onChange={e => setLang(e.target.value)}
            select
            name={'lang'}
            label={t('appSettingsForm.inputLabel.lang')}
            value={lang}
          >
            {AVAILABLE_LANGUAGES.map(lang => {
              return (
                <MenuItem key={lang} value={lang}>
                  {t('common:languages.' + lang)}
                </MenuItem>
              );
            })}
          </TextField>
          <Button type="submit" disabled={isSubmitting}>
            {t('save')}
          </Button>
        </form>
        <Typography variant="h4" gutterBottom component="h3">
          {t('appSettingsForm.subtitle.2')}
        </Typography>
        <form noValidate onSubmit={handleDebugForm}>
          <FormControlLabel
            control={
              <Checkbox
                checked={debug}
                name="debug"
                onChange={e => setDebug(e.target.checked)}
              />
            }
            label={t('appSettingsForm.inputLabel.debug')}
          />
          <Button type="submit" disabled={isSubmitting}>
            {t('save')}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}

AppSettingsForm.propTypes = {
  update: func.isRequired,
  setMessage: func.isRequired,
  setSeverity: func.isRequired,
  t: func.isRequired,
  config: object.isRequired
};

export default AppSettingsForm;
