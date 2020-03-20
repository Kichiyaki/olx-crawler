import React from 'react';
import { object, string as propString, arrayOf, func, bool } from 'prop-types';

import {
  TableHead,
  TableRow,
  TableSortLabel,
  TableCell
} from '@material-ui/core';

const MyTableHead = ({
  order,
  orderBy,
  onRequestSort,
  headCells,
  sortable
}) => {
  const createSortHandler = property => () => {
    if (property === orderBy) {
      onRequestSort(`${property} ${order === 'asc' ? 'desc' : 'asc'}`);
    } else {
      onRequestSort(property + ' asc');
    }
  };

  return (
    <TableHead>
      <TableRow>
        {sortable
          ? headCells.map(headCell => (
              <TableCell
                key={headCell.id}
                align={headCell.numeric ? 'right' : 'left'}
                padding={headCell.disablePadding ? 'none' : 'default'}
                sortDirection={orderBy === headCell.id ? order : false}
              >
                {headCell.sortable ? (
                  <TableSortLabel
                    active={orderBy === headCell.id}
                    direction={order}
                    onClick={createSortHandler(headCell.id)}
                  >
                    {headCell.label}
                  </TableSortLabel>
                ) : (
                  headCell.label
                )}
              </TableCell>
            ))
          : headCells.map(headCell => (
              <TableCell key={headCell.id}>{headCell.label}</TableCell>
            ))}
      </TableRow>
    </TableHead>
  );
};

MyTableHead.propTypes = {
  order: propString.isRequired,
  orderBy: propString,
  onRequestSort: func.isRequired,
  headCells: arrayOf(object).isRequired,
  sortable: bool.isRequired
};

MyTableHead.defaultProps = {
  order: 'asc',
  orderBy: '',
  onRequestSort: () => {},
  headCells: [],
  sortable: true
};

export default MyTableHead;
