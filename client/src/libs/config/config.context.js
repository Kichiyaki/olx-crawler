import { createContext } from 'react';

const ctx = createContext({
  config: {},
  reload: () => {},
  update: () => {}
});

ctx.displayName = 'Config';

export default ctx;
