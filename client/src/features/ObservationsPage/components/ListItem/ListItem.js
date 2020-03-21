import React from 'react';
import formatDate from 'date-fns/format';
import { DATE_FORMAT } from '@config/application';

import {
  TableRow,
  TableCell,
  Checkbox,
  Tooltip,
  IconButton
} from '@material-ui/core';
import { Edit as EditIcon } from '@material-ui/icons';

export default function TableItem({ item, t, selected, onSelect, onEdit }) {
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
      <TableCell>
        <Tooltip title="Edytuj obserwację">
          <IconButton
            onClick={onEdit}
            aria-label={'edit observation ' + item.name}
          >
            <EditIcon />
          </IconButton>
        </Tooltip>
      </TableCell>
    </TableRow>
  );
}
