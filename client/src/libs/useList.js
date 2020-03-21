import { useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { isPlainObject } from 'lodash';
import useQuery from './useQuery';
import isJSON from '@utils/isJSON';

export default ({
  offset: defaultOffset = 0,
  limit: defaultLimit = 25,
  filter: defaultFilter = {},
  order: defaultOrder = 'id asc'
} = {}) => {
  const query = useQuery();
  const offset = parseInt(query.get('offset')) || defaultOffset;
  const limit = parseInt(query.get('limit')) || defaultLimit;
  const filter = isJSON(query.get('filter'))
    ? JSON.parse(query.get('filter'))
    : defaultFilter;
  const order = query.get('order') || defaultOrder;
  const { pathname } = useLocation();
  const history = useHistory();
  const [selected, setSelected] = useState([]);

  const changeOffset = offset => {
    offset = parseInt(offset);
    if (isNaN(offset) || offset < 0) return;
    query.set('offset', offset);
    history.push(pathname + '?' + query.toString());
  };

  const changeLimit = e => {
    const limit = parseInt(e.target.value);
    if (isNaN(limit) || limit < 0) return;
    query.set('limit', limit);
    history.push(pathname + '?' + query.toString());
  };

  const changeFilter = filter => {
    if (isPlainObject(filter)) {
      filter = JSON.stringify(filter);
    }
    if (isJSON(filter)) {
      query.set('filter', filter);
      history.push(pathname + '?' + query.toString());
    }
  };

  const changeOrder = order => {
    if (
      typeof order !== 'string' ||
      !order ||
      (!order.toLowerCase().includes('asc') &&
        !order.toLowerCase().includes('desc'))
    )
      return;
    query.set('order', order);
    history.push(pathname + '?' + query.toString());
  };

  return {
    offset,
    limit,
    filter,
    order,
    selected,
    changeOffset,
    changeLimit,
    changeFilter,
    changeOrder,
    setSelected
  };
};
