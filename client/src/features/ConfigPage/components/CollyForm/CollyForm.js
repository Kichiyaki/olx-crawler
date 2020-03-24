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
  Button
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
    initialValues: config.colly
      ? config.colly
      : {
          limit: 1,
          delay: 5
        },
    onSubmit: async (values, { setSubmitting }) => {
      try {
        const response = await axios.patch(CONFIG.UPDATE.COLLY, values);
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
    validate: values => {
      const errors = {};

      for (let property in values) {
        if (!values[property]) {
          errors[property] = t('collyForm.error.required');
        }
      }

      return errors;
    }
  });

  return (
    <Card>
      <CardContent>
        <Typography variant="h3" gutterBottom component="h2">
          {t('collyForm.title')}
        </Typography>
        <form noValidate onSubmit={handleSubmit}>
          <TextField
            name={`limit`}
            type="number"
            label={t('collyForm.inputLabel.limit')}
            value={values.limit}
            onChange={handleChange}
            onBlur={handleBlur}
            fullWidth
            error={touched.limit && !!errors.limit}
            helperText={touched.limit && errors.limit}
          />
          <TextField
            name={`delay`}
            type="number"
            label={t('collyForm.inputLabel.delay')}
            value={values.delay}
            onChange={handleChange}
            onBlur={handleBlur}
            fullWidth
            error={touched.delay && !!errors.delay}
            helperText={touched.delay && errors.delay}
          />
          <Button type="submit" disabled={isSubmitting}>
            {t('save')}
          </Button>
          <Button type="button" disabled={isSubmitting} onClick={handleReset}>
            {t('reset')}
          </Button>
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
