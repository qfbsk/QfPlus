import {
  BrowserOpenURL,
  EventsOn,
  WindowSetLightTheme,
} from '../../wailsjs/runtime/runtime';

export type RuntimeEventName =
  | 'core-status-changed'
  | 'core-update-done'
  | 'core-update-error'
  | 'core-update-progress'
  | 'environment-import-progress'
  | 'migration-progress'
  | 'proxy-status-changed'
  | 'sdk-list-changed'
  | 'system-sdks-ready'
  | 'vfox-busy'
  | 'vfox-home-changed'
  | 'vfox-log';

export const onRuntimeEvent = <TPayload>(
  eventName: RuntimeEventName,
  handler: (payload: TPayload) => void
) => EventsOn(eventName, handler);

export const setWindowLightTheme = () => WindowSetLightTheme();

export const openExternalUrl = (url: string) => BrowserOpenURL(url);
