import {
  CheckCoreUpdate,
  GetCoreInfo,
  ListCoreVersions,
  SetCoreAutoUpdate,
  SwitchCoreVersion,
} from '../../wailsjs/go/app/App';

export type { CoreInfo, CoreRelease } from './appModels';

export const fetchCoreInfo = () => GetCoreInfo();

export const checkCoreUpdate = () => CheckCoreUpdate();

export const fetchCoreVersions = () => ListCoreVersions();

export const switchCoreVersion = (version: string) => SwitchCoreVersion(version);

export const setCoreAutoUpdate = (enabled: boolean) => SetCoreAutoUpdate(enabled);
