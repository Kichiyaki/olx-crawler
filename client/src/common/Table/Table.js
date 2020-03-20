import React from 'react';
import { bool, object } from 'prop-types';
import { pick, isNil, isArray, isPlainObject } from 'lodash';

import { Table, TableContainer, TableBody, Paper } from '@material-ui/core';
import TableToolbar from './TableToolbar';
import TableHead from './TableHead';
import TablePagination from './TablePagination';
import EmptyBody from './EmptyBody';

export default function MyTable({
  children,
  showPagination,
  showToolbar,
  CustomToolbarCmp,
  tableProps,
  ...rest
}) {
  const toolbarProps = pick(rest, ['name']);
  return (
    <Paper>
      {showToolbar && !isNil(CustomToolbarCmp) ? (
        <CustomToolbarCmp {...toolbarProps} />
      ) : showToolbar ? (
        <TableToolbar {...toolbarProps} />
      ) : null}
      <TableContainer component={Paper}>
        <Table {...tableProps}>
          <TableHead
            {...pick(rest, [
              'orderBy',
              'order',
              'onRequestSort',
              'headCells',
              'sortable'
            ])}
          />
          <TableBody>
            {!isNil(children) &&
            ((isArray(children) && children.length > 0) ||
              isPlainObject(children)) ? (
              children
            ) : (
              <EmptyBody />
            )}
          </TableBody>
        </Table>
      </TableContainer>
      {showPagination && (
        <TablePagination
          {...pick(rest, [
            'rowsPerPageOptions',
            'count',
            'page',
            'rowsPerPage',
            'rowsPerPageOptions',
            'onRequestNewPage',
            'onRequestChangeRowsPerPage'
          ])}
        />
      )}
    </Paper>
  );
}

MyTable.defaultProps = {
  showPagination: false,
  showToolbar: false,
  tableProps: {}
};

MyTable.propTypes = {
  showPagination: bool.isRequired,
  showToolbar: bool.isRequired,
  tableProps: object.isRequired
};
