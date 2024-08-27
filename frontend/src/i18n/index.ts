import { createI18n } from 'vue-i18n'
import elEnLocale from 'element-plus/es/locale/lang/en'
import elZhLocale from 'element-plus/es/locale/lang/zh-cn'
import zhLang from './lang/zh';
import enLang from './lang/en';

import { getSystemKey, setSystemKey } from '../stores/config'

export function getLang() {
  let currentLang = getSystemKey('lang')
  if (!currentLang) {
    try {
      const supported = ["en", "zh-cn"]
      const browserLang = (navigator.language || (navigator as any).browserLanguage).toLowerCase()
      if (supported.includes(browserLang)) {
        currentLang = browserLang
      } else {
        currentLang = "en";
      }
      setLang(currentLang)
    } catch (e) {
      console.log(e);
    }
  }
  return currentLang
}


const messages = {
  en: {
    ...enLang,
    ...elEnLocale
  },
  'zh-cn': {
    ...zhLang,
    ...elZhLocale
  }
}

export const i18n = createI18n({
  globalInjection: true,
  legacy: false,
  locale: getLang(),
  messages: messages
})
export function setLang(lang: string) {
  //currentLang = lang
  setSystemKey('lang', lang)
}
export function changeLang() {
  const lang = getLang()
  const setlang = lang == 'en' ? 'zh-cn' : 'en'
  setLang(setlang)
  return setlang
}
export function t(textkey: string, rep?: any) {
  return i18n.global.t(textkey, rep)
}
