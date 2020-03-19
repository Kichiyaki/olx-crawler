import prepareParams from '@utils/prepareParams';

export const SUGGESTIONS_REQ_PARAMS = prepareParams(
  {
    limit: 10,
    offset: 0
  },
  {},
  ['limit', 'offset']
);
