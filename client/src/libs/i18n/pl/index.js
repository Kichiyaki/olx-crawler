import * as namespaces from '@config/namespaces';
import common from './common';
import configPage from './config-page';
import mainPage from './main-page';
import observationsPage from './observations-page';

export default {
  common,
  [namespaces.MAIN_PAGE]: mainPage,
  [namespaces.OBSERVATIONS_PAGE]: observationsPage,
  [namespaces.CONFIG_PAGE]: configPage
};
