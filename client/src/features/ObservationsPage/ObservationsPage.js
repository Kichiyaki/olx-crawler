import React from 'react';
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import useSnackbar, { SEVERITY } from '@libs/useSnackbar';
import useList from '@libs/useList';
import useRequest from '@libs/useRequest';
import { OBSERVATIONS_PAGE } from '@config/namespaces';
import { OBSERVATIONS } from '@config/api_routes';
import { DEFAULT_LIMIT, DEFAULT_ORDER, HEAD_CELLS } from './constants';
import isAPIError from '@utils/isAPIError';
import { getRequestURL } from './utils';

import { Container, Snackbar, Button } from '@material-ui/core';
import { Alert } from '@material-ui/lab';
import ErrorPage from '@features/ErrorPage/ErrorPage';
import AppLayout from '@common/AppLayout/AppLayout';
import Spinner from '@common/Spinner/Spinner';
import Table from '@common/Table/Table';
import ObservationForm from './components/ObservationForm/ObservationForm';
import ListItem from './components/ListItem/ListItem';

function ObservationsPage() {
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
          >
            {observations.items.map(item => (
              <ListItem
                selected={selected.some(id => id === item.id)}
                onSelect={createSelectHandler(item.id)}
                item={item}
                key={item.id}
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
      <ObservationForm open={true} />
    </AppLayout>
  );
}

export default ObservationsPage;
