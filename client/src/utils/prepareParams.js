import {
  isString,
  defaultsDeep as fillWithDefaultValues,
  pick,
  isNumber,
  isPlainObject
} from 'lodash';
import URLSearchParams from '@ungap/url-search-params';
import isJSON from './isJSON';
import parseFilter from './parseFilter';

export default (obj = {}, defaults = {}, allowed = []) => {
  let { order, limit, offset, filter } = fillWithDefaultValues(
    pick(obj, allowed),
    defaults
  );
  const params = new URLSearchParams();
  if (limit && isNumber(parseInt(limit))) {
    params.set('limit', limit);
  }
  if (offset && isNumber(parseInt(offset))) {
    params.set('offset', offset);
  }
  if (order && isString(order)) {
    params.set('order', order);
  }
  if (filter && isJSON(filter)) {
    filter = JSON.parse(filter);
  }
  if (filter && isPlainObject(filter)) {
    const parsedFilterObj = parseFilter(filter);
    for (let property in parsedFilterObj) {
      params.set(property, parsedFilterObj[property]);
    }
  }

  return params;
};
