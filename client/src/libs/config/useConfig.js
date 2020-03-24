import { useContext } from 'react';
import ctx from './config.context';

export default () => {
  return useContext(ctx);
};
