import { isString } from 'lodash';

export default str => {
  if (!isString(str)) return false;
  try {
    JSON.parse(str);
    return true;
  } catch (err) {
    return false;
  }
};
