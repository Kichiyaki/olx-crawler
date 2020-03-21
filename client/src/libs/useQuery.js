import { useLocation } from 'react-router-dom';
import URLSearchParams from '@ungap/url-search-params';

export default function useQuery() {
  return new URLSearchParams(useLocation().search);
}
