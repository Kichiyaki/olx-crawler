import React from 'react';
import axios from 'axios';
import { func, object } from 'prop-types';
import { useFormik } from 'formik';
import isURL from 'validator/lib/isURL';
import { SEVERITY } from '@libs/useSnackbar';
import { CONFIG } from '@config/api_routes';
import isAPIError from '@utils/isAPIError';

import {
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  InputAdornment,
  IconButton
} from '@material-ui/core';
import { Delete as DeleteIcon } from '@material-ui/icons';

function ProxyForm({ update, config, setMessage, setSeverity, t }) {
  const {
    handleBlur,
    handleChange,
    handleSubmit,
    isSubmitting,
    values,
    touched,
    errors,
    setFieldValue,
    handleReset
  } = useFormik({
    initialValues: {
      proxy: config.proxy || []
    },
    onSubmit: async (values, { setSubmitting }) => {
      try {
        const response = await axios.patch(CONFIG.UPDATE.PROXY, values.proxy);
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
    validate: ({ proxy }) => {
      const errors = {};
      proxy.forEach((proxyaddr, index) => {
        if (
          !isURL(proxyaddr, {
            protocols: ['http', 'https', 'socks4', 'socks5']
          })
        ) {
          errors[`proxy[${index}]`] = t('proxyForm.error.invalidProxyAddress');
        }
      });
      return errors;
    }
  });

  const createDeleteHandler = index => () => {
    setFieldValue(
      'proxy',
      values.proxy.filter((_, i) => i !== index)
    );
  };

  const addProxy = () => setFieldValue('proxy', [...values.proxy, '']);

  return (
    <Card>
      <CardContent>
        <Typography variant="h3" gutterBottom component="h2">
          {t('proxyForm.title')}
        </Typography>
        <form noValidate onSubmit={handleSubmit}>
          {values.proxy.map((proxyaddr, index) => {
            return (
              <TextField
                key={`proxy[${index}]`}
                name={`proxy[${index}]`}
                value={proxyaddr}
                onChange={handleChange}
                onBlur={handleBlur}
                fullWidth
                InputProps={{
                  endAdornment: (
                    <InputAdornment disablePointerEvents={false} position="end">
                      <IconButton onClick={createDeleteHandler(index)}>
                        <DeleteIcon />
                      </IconButton>
                    </InputAdornment>
                  )
                }}
                error={
                  touched.proxy &&
                  touched.proxy[index] &&
                  !!errors[`proxy[${index}]`]
                }
                helperText={
                  touched.proxy &&
                  touched.proxy[index] &&
                  errors[`proxy[${index}]`]
                }
              />
            );
          })}
          <Button type="submit" disabled={isSubmitting}>
            {t('save')}
          </Button>
          <Button type="button" disabled={isSubmitting} onClick={addProxy}>
            {t('proxyForm.button.addProxy')}
          </Button>
          <Button type="button" disabled={isSubmitting} onClick={handleReset}>
            {t('reset')}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}

ProxyForm.propTypes = {
  update: func.isRequired,
  setMessage: func.isRequired,
  setSeverity: func.isRequired,
  t: func.isRequired,
  config: object.isRequired
};

export default ProxyForm;
