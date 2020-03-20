import { OBSERVATIONS } from '@config/api_routes';
import prepareParams from '@utils/prepareParams';

export const getRequestURL = obj => {
  return (
    OBSERVATIONS.BROWSE +
    '?' +
    prepareParams(obj, {}, ['limit', 'offset', 'order', 'filter'])
  );
};
