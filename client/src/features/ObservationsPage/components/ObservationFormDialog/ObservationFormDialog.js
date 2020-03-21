import React from 'react';
import { useFormik } from 'formik';
import { format, isAfter } from 'date-fns';

import { makeStyles } from '@material-ui/core/styles';
import {
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  MenuItem,
  Typography,
  FormControlLabel,
  Checkbox
} from '@material-ui/core';

const useStyles = makeStyles(theme => ({
  dialogContent: {
    '& > *:not(:last-child)': {
      marginBottom: theme.spacing(1)
    }
  }
}));

function ObservationFormDialog({ observation, onClose, open, onSubmit, t }) {
  const classes = useStyles();
  const {
    handleBlur,
    handleChange,
    handleSubmit,
    isSubmitting,
    values,
    touched,
    errors,
    setFieldValue
  } = useFormik({
    initialValues: {
      name: observation && observation.name ? observation.name : '',
      url: observation && observation.url ? observation.url : '',
      keywords: observation && observation.keywords ? observation.keywords : [],
      last_check_at: format(
        observation && observation.last_check_at
          ? new Date(observation.last_check_at)
          : new Date(),
        'yyyy-MM-dd HH:mm'
      )
        .split(' ')
        .join('T'),
      deleted_keywords: [],
      started:
        observation && typeof observation.started === 'boolean'
          ? observation.started
          : false
    },
    onSubmit: async (values, { setSubmitting }) => {
      const success = await onSubmit({
        ...values,
        keywords: values.keywords.filter(
          keyword => keyword.for && keyword.value && keyword.type
        ),
        last_check_at: new Date(values.last_check_at)
      });
      if (success) {
        onClose();
      } else {
        setSubmitting(false);
      }
    },
    validate: values => {
      const errors = {};
      if (!values.name) {
        errors.name = t('observationFormDialog.error.required');
      }
      if (!values.url) {
        errors.url = t('observationFormDialog.error.required');
      }
      if (isAfter(new Date(values.last_check_at), new Date())) {
        errors.last_check_at = t('observationFormDialog.error.invalidDate');
      }
      return errors;
    }
  });

  const addKeyword = () => {
    const obj = { for: '', value: '', type: '', group: '' };
    setFieldValue('keywords', [...values.keywords, obj]);
  };

  const deleteKeyword = (keyword, index) => {
    if (keyword.id) {
      setFieldValue(
        'keywords',
        values.keywords.filter(_keyword => _keyword.id !== keyword.id)
      );
      setFieldValue('deleted_keywords', [
        ...values.deleted_keywords,
        keyword.id
      ]);
    } else {
      setFieldValue(
        'keywords',
        values.keywords.filter((_, _index) => _index !== index)
      );
    }
  };

  return (
    <Dialog open={open} onClose={isSubmitting ? undefined : onClose}>
      <DialogTitle>
        {observation
          ? t('observationFormDialog.title.edit')
          : t('observationFormDialog.title.add')}
      </DialogTitle>
      <DialogContent className={classes.dialogContent}>
        <TextField
          name="name"
          label={t('observationFormDialog.inputLabel.name')}
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.name}
          error={touched.name && !!errors.name}
          helperText={touched.name && errors.name}
          fullWidth
        />
        <TextField
          name="url"
          label={t('observationFormDialog.inputLabel.url')}
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.url}
          error={touched.url && !!errors.url}
          helperText={touched.url && errors.url}
          fullWidth
        />
        <TextField
          type="datetime-local"
          name="last_check_at"
          label={t('observationFormDialog.inputLabel.last_check_at')}
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.last_check_at}
          error={touched.last_check_at && !!errors.last_check_at}
          helperText={touched.last_check_at && errors.last_check_at}
          fullWidth
          inputProps={{ max: format(new Date(), 'yyyy-MM-dd') }}
        />
        <FormControlLabel
          control={
            <Checkbox
              checked={values.started}
              onChange={handleChange}
              name="started"
            />
          }
          label={t('observationFormDialog.inputLabel.started')}
        />
        {values.keywords.map((keyword, index) => {
          return (
            <div key={index}>
              <Typography component="p">
                {t('observationFormDialog.keyword_title', { index: index + 1 })}{' '}
                <Button onClick={() => deleteKeyword(keyword, index)}>
                  {t('delete')}
                </Button>
              </Typography>
              <TextField
                select
                name={`keywords[${index}].type`}
                label={t('observationFormDialog.inputLabel.keyword_type', {
                  index: index + 1
                })}
                onChange={handleChange}
                onBlur={handleBlur}
                value={keyword.type}
                fullWidth
              >
                <MenuItem value="">
                  {t('observationFormDialog.select')}
                </MenuItem>
                <MenuItem value="required">
                  {t('observationFormDialog.required')}
                </MenuItem>
                <MenuItem value="one_of">
                  {t('observationFormDialog.one_of')}
                </MenuItem>
                <MenuItem value="excluded">
                  {t('observationFormDialog.excluded')}
                </MenuItem>
              </TextField>
              <TextField
                select
                name={`keywords[${index}].for`}
                label={t('observationFormDialog.inputLabel.keyword_for', {
                  index: index + 1
                })}
                onChange={handleChange}
                onBlur={handleBlur}
                value={keyword.for}
                fullWidth
              >
                <MenuItem value="">
                  {t('observationFormDialog.select')}
                </MenuItem>
                <MenuItem value="title">
                  {t('observationFormDialog.for_title')}
                </MenuItem>
                <MenuItem value="description">
                  {t('observationFormDialog.description')}
                </MenuItem>
              </TextField>
              <TextField
                name={`keywords[${index}].value`}
                label={t('observationFormDialog.inputLabel.keyword_value', {
                  index: index + 1
                })}
                onChange={handleChange}
                onBlur={handleBlur}
                value={keyword.value}
                fullWidth
              />
              {keyword.type === 'one_of' ? (
                <TextField
                  name={`keywords[${index}].group`}
                  label={t('observationFormDialog.inputLabel.keyword_group', {
                    index: index + 1
                  })}
                  onChange={handleChange}
                  onBlur={handleBlur}
                  value={keyword.group}
                  fullWidth
                />
              ) : null}
            </div>
          );
        })}
        <Button type="button" fullWidth onClick={() => addKeyword()}>
          {t('observationFormDialog.button.addKeyword')}
        </Button>
      </DialogContent>
      <DialogActions>
        <Button disabled={isSubmitting} onClick={handleSubmit} color="primary">
          {observation
            ? t('observationFormDialog.button.save')
            : t('observationFormDialog.button.add')}
        </Button>
        <Button
          disabled={isSubmitting}
          onClick={onClose}
          color="primary"
          autoFocus
        >
          {t('cancel')}
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export default ObservationFormDialog;
