import type { TranslationDict } from './keys';
import { zhEnvironment } from './zh/environment';
import { zhInfo } from './zh/info';
import { zhSdk } from './zh/sdk';
import { zhSettings } from './zh/settings';
import { zhShell } from './zh/shell';

export const zh: TranslationDict = {
  ...zhShell,
  ...zhSettings,
  ...zhSdk,
  ...zhEnvironment,
  ...zhInfo,
};
