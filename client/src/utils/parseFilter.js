export const GT_OPERATOR = 'gt';
export const LT_OPERATOR = 'lt';
export const GTE_OPERATOR = 'gte';
export const LTE_OPERATOR = 'lte';
export const MATCH_OPERATOR = 'match';
export const NEQ_OPERATOR = 'neq';
const validOperators = [
  GT_OPERATOR,
  LT_OPERATOR,
  GTE_OPERATOR,
  LTE_OPERATOR,
  MATCH_OPERATOR,
  NEQ_OPERATOR
];

const findOperator = (name = '') => {
  return validOperators.find(operator => {
    return name.slice(-1 * operator.length).toLowerCase() === operator;
  });
};

export default (obj = {}, customNameMapping = () => {}) => {
  const newObj = {};
  for (let i in obj) {
    let name = customNameMapping(i);
    if (!name) {
      name = i;
      const operator = findOperator(i);
      if (operator) {
        name = name.substring(0, name.length - operator.length);
      }
      name = name
        .replace(/([a-z])([A-Z])/g, '$1 $2')
        .split(' ')
        .filter(val => val !== '')
        .map(val => val.toLowerCase())
        .join('_');
      if (operator) {
        name += `__${operator}`;
      }
    }
    newObj[name] = obj[i];
  }
  return newObj;
};
