import {
  GetProxyStatus,
  SetProxyEnabled,
  ImportProxySubscription,
  GetProxyGroups,
  GetProxyQuickStatus,
  SelectProxyNode,
  SetProxyGroup,
  SetProxyQuickEnabled,
  TestProxyGroupDelay,
  TestProxyNodeDelay,
  TestProxyQuickDelay,
} from '../../wailsjs/go/app/App';

export type { ProxyStatus, ProxyGroup, ProxyNode, ProxyQuickStatus } from './appModels';

export const fetchProxyStatus = () => GetProxyStatus();

export const setProxyEnabled = (enabled: boolean) => SetProxyEnabled(enabled);

export const importProxySubscription = (url: string) => ImportProxySubscription(url);

export const fetchProxyGroups = () => GetProxyGroups();

export const selectProxyNode = (group: string, node: string) => SelectProxyNode(group, node);

export const setProxyGroup = (group: string) => SetProxyGroup(group);

export const testProxyGroupDelay = (group: string) => TestProxyGroupDelay(group);

export const testProxyNodeDelay = (name: string) => TestProxyNodeDelay(name);

export const fetchProxyQuickStatus = () => GetProxyQuickStatus();

export const testProxyQuickDelay = () => TestProxyQuickDelay();

export const setProxyQuickEnabled = (enabled: boolean) => SetProxyQuickEnabled(enabled);
