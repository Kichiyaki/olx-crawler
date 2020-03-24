import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import pl from './pl';
import { DEFAULT_LANGUAGE } from '@config/application';

i18n
  // pass the i18n instance to react-i18next.
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    fallbackLng: DEFAULT_LANGUAGE,
    lng: DEFAULT_LANGUAGE,
    debug: true,
    load: 'languageOnly',
    resources: {
      pl
    },
    defaultNS: 'common',

    interpolation: {
      escapeValue: false // not needed for react as it escapes by default
    }
  });

export default i18n;
