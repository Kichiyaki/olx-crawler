import '@testing-library/jest-dom/extend-expect';
import 'react-act';
import { advanceTo, clear } from 'jest-date-mock';

beforeEach(() => {
  advanceTo(0);
});

afterEach(() => {
  clear();
});

jest.setTimeout(10000);
