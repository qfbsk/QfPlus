import { ref, watch } from 'vue';
import { en } from './en';
import type { SupportedLang } from './keys';
import { zh } from './zh';

const langStorageKey = 'qfplus-lang';

const getInitialLang = (): SupportedLang => {
  const saved = localStorage.getItem(langStorageKey) ?? localStorage.getItem('vfox-lang');
  if (saved === 'en' || saved === 'zh') {
    return saved;
  }
  return 'en';
};

export const currentLang = ref<SupportedLang>(getInitialLang());

watch(currentLang, (newLang) => {
  localStorage.setItem(langStorageKey, newLang);
});

const messages = { en, zh };

export const t = (key: string, params?: Record<string, string | number | boolean>): string => {
  const dict = messages[currentLang.value] || messages.en;
  let value = dict[key] || key;
  if (params) {
    for (const [name, replacement] of Object.entries(params)) {
      value = value.replaceAll(`{${name}}`, String(replacement));
    }
  }
  return value;
};
