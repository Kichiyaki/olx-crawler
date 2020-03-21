import React, { useState } from 'react';
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import useSnackbar, { SEVERITY } from '@libs/useSnackbar';
import useList from '@libs/useList';
import useRequest from '@libs/useRequest';
import { OBSERVATIONS_PAGE } from '@config/namespaces';
import { OBSERVATIONS, KEYWORDS } from '@config/api_routes';
import { DEFAULT_LIMIT, DEFAULT_ORDER, HEAD_CELLS } from './constants';
import isAPIError from '@utils/isAPIError';
import { getRequestURL } from './utils';

import { Container, Snackbar, Button } from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import AppLayout from '@common/AppLayout/AppLayout';
import Spinner from '@common/Spinner/Spinner';
import Table from '@common/Table/Table';
import EnhancedToolbar from './components/EnhancedToolbar/EnhancedToolbar';
import ObservationFormDialog from './components/ObservationFormDialog/ObservationFormDialog';
import ListItem from './components/ListItem/ListItem';

function ObservationsPage() {
  const [dialogID, setDialogID] = useState(0);
  const [editedObservationID, setEditedObservationID] = useState(0);
  const { t } = useTranslation(OBSERVATIONS_PAGE);
  const {
    order,
    limit,
    offset,
    filter,
    selected,
    changeLimit,
    changeOffset,
    changeOrder,
    setSelected
  } = useList({
    limit: DEFAULT_LIMIT,
    order: DEFAULT_ORDER
  });
  const {
    setSeverity,
    setMessage,
    message,
    alertProps,
    snackbarProps
  } = useSnackbar({
    anchorOrigin: { vertical: 'top', horizontal: 'right' }
  });
  const { data: observations, statusCode, error, refresh } = useRequest(
    getRequestURL({ order, limit, offset, filter })
  );
  const loading = !observations;

  const handleNewPage = page => {
    changeOffset(page * limit);
  };

  const handleDelete = async () => {
    try {
      const response = await axios.delete(
        OBSERVATIONS.DELETE + '?id=' + selected.join(',')
      );
      if (Array.isArray(response.data.data)) {
        setSeverity(SEVERITY.SUCCESS);
        setMessage(t('deleted', { count: response.data.data.length }));
        setSelected([]);
        refresh();
      }
    } catch (error) {
      if (isAPIError(error)) {
        setSeverity(SEVERITY.ERROR);
        setMessage(error.response.data.errors[0].message);
      }
    }
  };

  const handleSnackbarClose = () => {
    setSelected([]);
  };

  const createSelectHandler = id => () => {
    if (selected.some(otherID => otherID === id)) {
      setSelected(selected.filter(otherID => otherID !== id));
    } else {
      setSelected([...selected, id]);
    }
  };

  const handleCreate = input => {
    return axios
      .post(OBSERVATIONS.STORE, input)
      .then(() => {
        setSeverity(SEVERITY.SUCCESS);
        setMessage(t('created'));
        refresh();
        return true;
      })
      .catch(err => {
        if (isAPIError(err)) {
          setSeverity(SEVERITY.ERROR);
          setMessage(err.response.data.errors[0].message);
        }
        return false;
      });
  };

  const handleEdit = async input => {
    try {
      await axios.patch(OBSERVATIONS.UPDATE + `/${editedObservationID}`, input);
      if (input.deleted_keywords.length > 0) {
        await axios.delete(
          KEYWORDS.DELETE + '?id=' + input.deleted_keywords.join(',')
        );
      }
      setSeverity(SEVERITY.SUCCESS);
      setMessage(t('updated'));
      refresh();
      return true;
    } catch (error) {
      if (isAPIError(error)) {
        setSeverity(SEVERITY.ERROR);
        setMessage(error.response.data.errors[0].message);
      }
      return false;
    }
  };

  const createEditObservationHandler = id => () => {
    setEditedObservationID(id);
    setDialogID(2);
  };

  const handleDialogClose = () => {
    setDialogID(0);
  };

  const renderDialog = () => {
    switch (dialogID) {
      case 1:
        return (
          <ObservationFormDialog
            onClose={handleDialogClose}
            onSubmit={handleCreate}
            open={true}
            t={t}
          />
        );
      case 2:
        return (
          <ObservationFormDialog
            observation={observations.items.find(
              item => item.id === editedObservationID
            )}
            onClose={handleDialogClose}
            onSubmit={handleEdit}
            open={true}
            t={t}
          />
        );
      default:
        return null;
    }
  };

  if (error || statusCode !== 200) {
    return <ErrorPage statusCode={statusCode} error={error} />;
  }

  return (
    <AppLayout>
      {loading && <Spinner />}
      {!loading && (
        <Container>
          <Table
            name={t('table.title')}
            showToolbar
            orderBy={order.split(' ')[0]}
            order={order.split(' ')[1]}
            headCells={HEAD_CELLS.map(cell => ({
              ...cell,
              label: t(cell.label)
            }))}
            showPagination
            onRequestSort={changeOrder}
            count={observations.total}
            onRequestChangeRowsPerPage={changeLimit}
            onRequestNewPage={handleNewPage}
            rowsPerPage={limit}
            page={Math.floor(offset / limit)}
            CustomToolbarCmp={props => (
              <EnhancedToolbar
                {...props}
                onClickAddObservation={() => setDialogID(1)}
                t={t}
              />
            )}
          >
            {observations.items.map(item => (
              <ListItem
                selected={selected.some(id => id === item.id)}
                onSelect={createSelectHandler(item.id)}
                item={item}
                key={item.id}
                onEdit={createEditObservationHandler(item.id)}
                t={t}
              />
            ))}
          </Table>
          <Snackbar
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right'
            }}
            ClickAwayListenerProps={{ mouseEvent: false }}
            open={selected.length > 0}
            onClose={handleSnackbarClose}
            message={t('selected', { count: selected.length })}
            action={
              <>
                <Button onClick={handleDelete} color="secondary" size="small">
                  {t('delete')}
                </Button>
                <Button
                  color="secondary"
                  onClick={handleSnackbarClose}
                  size="small"
                >
                  {t('cancel')}
                </Button>
              </>
            }
          />
        </Container>
      )}
      <Snackbar {...snackbarProps}>
        <Alert {...alertProps}>{message}</Alert>
      </Snackbar>
      {renderDialog()}
    </AppLayout>
  );
}

export default ObservationsPage;
