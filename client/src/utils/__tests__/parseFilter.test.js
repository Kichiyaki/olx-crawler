import parseFilter from '../parseFilter';

describe('parseFilter', () => {
  const obj = {
    createdAtGT: 'asd',
    createdAtLTE: 'asdff',
    countGTE: 123,
    countLT: 299,
    testMATCH: 'asdd',
    testNEQ: '123asd',
    hello: 'asdd',
    anotherObjPropertyID: '123'
  };

  test('without customNameMapping', () => {
    const result = parseFilter(obj);
    expect(result['created_at__gt']).toBe(obj.createdAtGT);
    expect(result['created_at__lte']).toBe(obj.createdAtLTE);
    expect(result['count__gte']).toBe(obj.countGTE);
    expect(result['count__lt']).toBe(obj.countLT);
    expect(result['test__match']).toBe(obj.testMATCH);
    expect(result['test__neq']).toBe(obj.testNEQ);
    expect(result['hello']).toBe(obj.hello);
    expect(result['another_obj_property_id']).toBe(obj.anotherObjPropertyID);
  });

  test('with customNameMapping', () => {
    const result = parseFilter(obj, name => {
      return name === 'createdAtGT' ? 'other_name' : '';
    });
    expect(result['other_name']).toBe(obj.createdAtGT);
    expect(result['created_at__lte']).toBe(obj.createdAtLTE);
    expect(result['count__gte']).toBe(obj.countGTE);
    expect(result['count__lt']).toBe(obj.countLT);
    expect(result['test__match']).toBe(obj.testMATCH);
    expect(result['test__neq']).toBe(obj.testNEQ);
    expect(result['hello']).toBe(obj.hello);
    expect(result['another_obj_property_id']).toBe(obj.anotherObjPropertyID);
  });
});
