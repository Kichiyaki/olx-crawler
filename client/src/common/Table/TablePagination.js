import React from 'react';
import { number, arrayOf, func, bool } from 'prop-types';
import { useTranslation } from 'react-i18next';
import { COMMON } from '@config/namespaces';

import { TablePagination } from '@material-ui/core';

export default function MyTablePagination({
  count,
  rowsPerPageOptions,
  rowsPerPage,
  page,
  onRequestNewPage,
  onRequestChangeRowsPerPage,
  fetching
}) {
  const createOnChangePageHandler = (_, newPage = 0) => {
    onRequestNewPage(newPage);
  };
  const { t } = useTranslation(COMMON);

  return (
    <TablePagination
      rowsPerPageOptions={rowsPerPageOptions}
      component="div"
      count={count}
      SelectProps={{
        native: process.env.NODE_ENV === 'test',
        disabled: fetching
      }}
      rowsPerPage={rowsPerPage}
      page={page}
      onChangePage={createOnChangePageHandler}
      onChangeRowsPerPage={onRequestChangeRowsPerPage}
      labelRowsPerPage={t('rowsPerPage')}
      labelDisplayedRows={obj => t('pagination', obj)}
    />
  );
}

MyTablePagination.propTypes = {
  total: number.isRequired,
  page: number.isRequired,
  rowsPerPage: number.isRequired,
  rowsPerPageOptions: arrayOf(number).isRequired,
  onRequestNewPage: func.isRequired,
  onRequestChangeRowsPerPage: func.isRequired,
  fetching: bool.isRequired
};

MyTablePagination.defaultProps = {
  total: 0,
  page: 1,
  fetching: false,
  rowsPerPageOptions: [10, 25, 50],
  rowsPerPage: 10,
  onRequestNewPage: () => {},
  onRequestChangeRowsPerPage: () => {}
};
