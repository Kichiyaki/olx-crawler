import prepareParams from '@utils/prepareParams';

export const getSuggestionsReqParams = (offset = 0) =>
  prepareParams(
    {
      limit: 12,
      offset,
      order: 'id desc'
    },
    {},
    ['limit', 'offset', 'order']
  );
