export const DEFAULT_LIMIT = 25;
export const DEFAULT_ORDER = 'id desc';
export const HEAD_CELLS = [
  {
    id: 'select',
    disablePadding: false,
    label: 'table.head.select',
    sortable: false
  },
  {
    id: 'id',
    disablePadding: false,
    label: 'ID',
    sortable: true
  },
  {
    id: 'name',
    disablePadding: false,
    label: 'table.head.name',
    sortable: true
  },
  {
    id: 'started',
    disablePadding: false,
    label: 'table.head.started',
    sortable: true
  },
  {
    id: 'last_check_at',
    disablePadding: false,
    label: 'table.head.lastCheckAt',
    sortable: true
  },
  {
    id: 'action',
    disablePadding: false,
    label: 'table.head.action',
    sortable: false
  }
];
