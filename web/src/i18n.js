import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import translations_de_DE from './translations/de-DE.json';
import translations_en_US from './translations/en-US.json';

export const fallbackLanguage = 'en-US';

export const languages = {
  'de-DE': {
    displayName: 'Deutsch',
    translation: translations_de_DE,
  },
  'en-US': {
    displayName: 'English',
    translation: translations_en_US,
  },
}

const i18nConfig = {
  resources: languages,
  fallbackLng: fallbackLanguage,
  debug: true,
  interpolation: {
    escapeValue: false,
  },
}

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init(i18nConfig);

export default i18n;
