import type { TranslationDict } from './keys';
import { enEnvironment } from './en/environment';
import { enInfo } from './en/info';
import { enSdk } from './en/sdk';
import { enSettings } from './en/settings';
import { enShell } from './en/shell';

export const en: TranslationDict = {
  ...enShell,
  ...enSettings,
  ...enSdk,
  ...enEnvironment,
  ...enInfo,
};
