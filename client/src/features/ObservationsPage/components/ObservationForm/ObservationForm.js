import React, { Fragment } from 'react';
import { useFormik } from 'formik';
import { format } from 'date-fns';

import {
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  MenuItem,
  Typography
} from '@material-ui/core';

function ObservationForm({ observation, onClose, open }) {
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
      one_of: observation ? observation.one_of : [],
      excluded: observation ? observation.excluded : [],
      last_check_at: format(
        observation ? new Date(observation.last_check_at) : new Date(),
        'yyyy-MM-dd'
      )
    },
    onSubmit: (_, { setSubmitting }) => {
      setSubmitting(false);
    }
  });

  const addOneOfOrExcluded = f => {
    const obj = { for: '', value: '' };
    setFieldValue(f, [...values[f], obj]);
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
        />
        {values.one_of.map((oneOf, index) => {
          return (
            <Fragment key={index}>
              <Typography>
                Słowo #{index + 1} <Button>Usuń</Button>
              </Typography>
              <TextField
                select
                name={`one_of[${index}].for`}
                label={`${index + 1}. Dla`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={oneOf.for}
                fullWidth
              >
                <MenuItem value="">Wybierz</MenuItem>
                <MenuItem value="title">Tytuł</MenuItem>
                <MenuItem value="description">Opis</MenuItem>
              </TextField>
              <TextField
                name={`one_of[${index}].value`}
                label={`${index + 1}. Słowo`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={oneOf.value}
                fullWidth
              />
            </Fragment>
          );
        })}
        <Button
          type="button"
          fullWidth
          onClick={() => addOneOfOrExcluded('one_of')}
        >
          Dodaj mile widziane słowo
        </Button>
        {values.excluded.map((excluded, index) => {
          return (
            <Fragment key={index}>
              <Typography>Słowo #{index + 1}</Typography>
              <TextField
                select
                name={`excluded[${index}].for`}
                label={`${index + 1}. Dla`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={excluded.for}
                fullWidth
              >
                <MenuItem value="">Wybierz</MenuItem>
                <MenuItem value="title">Tytuł</MenuItem>
                <MenuItem value="description">Opis</MenuItem>
              </TextField>
              <TextField
                name={`excluded[${index}].value`}
                label={`${index + 1}. Słowo`}
                onChange={handleChange}
                onBlur={handleBlur}
                value={excluded.value}
                fullWidth
              />
            </Fragment>
          );
        })}
        <Button
          type="button"
          fullWidth
          onClick={() => addOneOfOrExcluded('excluded')}
        >
          Wyklucz słowo
        </Button>
      </DialogContent>
      <DialogActions>
        <Button disabled={isSubmitting} onClick={handleSubmit} color="primary">
          Utwórz
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

export default ObservationForm;
