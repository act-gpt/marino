import i18n from 'i18next'
import Backend from 'i18next-http-backend'
import LanguageDetector from 'i18next-browser-languagedetector'
import { initReactI18next } from 'react-i18next'
//import { createRequire } from "module";
//const require = createRequire(import.meta.url);
import data from '../package.json'

export default i18n
  // load translation using http -> see /public/locales
  // learn more: https://github.com/i18next/i18next-http-backend
  .use(Backend)
  // detect user language
  // learn more: https://github.com/i18next/i18next-browser-languageDetector
  .use(LanguageDetector)
  // pass the i18n instance to react-i18next.
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    backend: {
      loadPath: `${process.env.REACT_APP_CDN || ''}/locales/{{lng}}/{{ns}}.json`,
      crossDomain: true,
      withCredentials: true,
      queryStringParams: { v: data.version }
    },
    preload: ['en'],
    fallbackLng: {
      'en-US': ['en'],
      'zh-CN': ['zh-Hans', 'en'],
      default: ['en']
    },
    load: 'currentOnly',
    debug: process.env.NODE_ENV !== 'production',
    detection: {
      lookupQuerystring: 'lang',
      order: ['querystring', 'localStorage', 'navigator', 'cookie']
    },
    detectLngQS: 'lang',
    interpolation: {
      escapeValue: false // not needed for react as it escapes by default
    },
    react: {
      wait: true
    }
  })
