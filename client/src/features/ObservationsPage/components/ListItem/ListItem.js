import React from 'react';
import formatDate from 'date-fns/format';
import { DATE_FORMAT } from '@config/application';

import { TableRow, TableCell, Checkbox } from '@material-ui/core';

export default function TableItem({ item, t, selected, onSelect }) {
  return (
    <TableRow>
      <TableCell>
        <Checkbox checked={selected} onChange={onSelect} />
      </TableCell>
      <TableCell component="th">{item.id}</TableCell>
      <TableCell>{item.name}</TableCell>
      <TableCell>{t('table.started_' + item.started)}</TableCell>
      <TableCell>
        {formatDate(new Date(item.last_check_at), DATE_FORMAT)}
      </TableCell>
      <TableCell>Edit</TableCell>
    </TableRow>
  );
}
