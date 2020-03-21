import React from 'react';
import { useTranslation } from 'react-i18next';
import { COMMON } from '@config/namespaces';

import { TableRow, TableCell } from '@material-ui/core';

export default function EmptyBody() {
  const { t } = useTranslation(COMMON);
  return (
    <TableRow>
      <TableCell colSpan="1000" align="center" component="th">
        {t('emptyTableBody')}
      </TableCell>
    </TableRow>
  );
}
