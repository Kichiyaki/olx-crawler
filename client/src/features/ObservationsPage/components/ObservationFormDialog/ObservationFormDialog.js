import React from 'react';
import { useFormik } from 'formik';
import { format, isAfter } from 'date-fns';

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

function ObservationFormDialog({ observation, onClose, open, onSubmit }) {
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
      name: observation ? observation.name : '',
      url: observation ? observation.url : '',
      keywords: observation ? observation.keywords : [],
      last_check_at: format(
        observation ? new Date(observation.last_check_at) : new Date(),
        'yyyy-MM-dd'
      ),
      deleted: [],
      started:
        observation && typeof observation.started == 'bool'
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
        errors.name = 'Wymagane';
      }
      if (!values.url) {
        errors.url = 'Wymagane';
      }
      if (isAfter(new Date(values.last_check_at), new Date())) {
        errors.last_check_at = 'Niepoprawna data';
      }
      return errors;
    }
  });

  const addKeyword = () => {
    const obj = { for: '', value: '', type: '' };
    setFieldValue('keywords', [...values.keywords, obj]);
  };

  const deleteKeyword = (keyword, index) => {
    if (keyword.id) {
      setFieldValue(
        'keywords',
        values.keywords.filter(_keyword => _keyword.id !== keyword.id)
      );
      setFieldValue('deleted', [...values.deleted, keyword.id]);
    } else {
      setFieldValue(
        'keywords',
        values.keywords.filter((_, _index) => _index !== index)
      );
    }
  };

  return (
    <Dialog open={open} onClose={isSubmitting ? undefined : onClose}>
      <DialogTitle>{observation ? 'Edytowanie' : 'Tworzenie'}</DialogTitle>
      <DialogContent>
        <TextField
          name="name"
          label="Nazwa"
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.name}
          error={touched.name && !!errors.name}
          helperText={touched.name && errors.name}
          fullWidth
        />
        <TextField
          name="url"
          label="URL"
          onChange={handleChange}
          onBlur={handleBlur}
          value={values.url}
          error={touched.url && !!errors.url}
          helperText={touched.url && errors.url}
          fullWidth
        />
        <TextField
          type="date"
          name="last_check_at"
          label="Od kiedy ma zacząć sprawdzać"
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
          label="Włączony"
        />
        {values.keywords.map((keyword, index) => {
          return (
            <div key={index}>
              <Typography component="p">
                Słowo #{index + 1}{' '}
                <Button onClick={() => deleteKeyword(keyword, index)}>
                  Usuń
                </Button>
              </Typography>
              <TextField
                select
                name={`keywords[${index}].type`}
                label={`${index + 1}. Typ`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={keyword.type}
                fullWidth
              >
                <MenuItem value="">Wybierz</MenuItem>
                <MenuItem value="one_of">Jedno z</MenuItem>
                <MenuItem value="excluded">Wykluczone</MenuItem>
              </TextField>
              <TextField
                select
                name={`keywords[${index}].for`}
                label={`${index + 1}. Dla`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={keyword.for}
                fullWidth
              >
                <MenuItem value="">Wybierz</MenuItem>
                <MenuItem value="title">Tytuł</MenuItem>
                <MenuItem value="description">Opis</MenuItem>
              </TextField>
              <TextField
                name={`keywords[${index}].value`}
                label={`${index + 1}. Słowo`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={keyword.value}
                fullWidth
              />
            </div>
          );
        })}
        <Button type="button" fullWidth onClick={() => addKeyword()}>
          Dodaj słowo kluczowe
        </Button>
      </DialogContent>
      <DialogActions>
        <Button disabled={isSubmitting} onClick={handleSubmit} color="primary">
          {observation ? 'Zapisz' : 'Dodaj'}
        </Button>
        <Button
          disabled={isSubmitting}
          onClick={onClose}
          color="primary"
          autoFocus
        >
          Anuluj
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export default ObservationFormDialog;
