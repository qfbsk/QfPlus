import type { PluginInfo } from '../services/appModels';
import { ref, type Ref } from 'vue';
import { t } from '../i18n';
import {
  addCustomSdkReference,
  addVfoxPlugin,
  detectSystemSdkVersion,
  fetchCachedSystemSdks,
  fetchNonVfoxSdks,
  hijackSystemSdkPath,
  installPluginVersion,
  refreshSystemSdkCache,
  removePlugin,
  useCustomSdkReference,
} from '../services/pluginMarketService';
import { useDownloadQueue } from './useDownloadQueue';
import { extractDownloadSpeed, extractProgressPercent } from './taskLogParser';
import { onRuntimeEvent } from '../services/runtimeService';

type UsePluginMarketActionsOptions = {
  selectedPlugin: Ref<PluginInfo | null>;
  installedVersions: Ref<Set<string>>;
  runTask: (title: string, task: () => Promise<void>) => Promise<boolean>;
  notifyError: (message: string) => void;
  notifyTaskError: (err: unknown, fallback: string) => void;
  formatError: (err: unknown, fallback: string) => string;
  fetchPlugins: () => Promise<void>;
  fetchPluginVersions: (pluginName: string) => Promise<void>;
};

export const usePluginMarketActions = (options: UsePluginMarketActionsOptions) => {
  const queue = useDownloadQueue();
  const addingPlugin = ref<string | null>(null);
  const removingPlugin = ref<string | null>(null);
  const installingVersion = ref<string | null>(null);
  const confirmRemove = ref<string | null>(null);
  const confirmRemoveCustomSdks = ref<Array<{ path: string; version: string }>>([]);

  const linkSystemSdkAfterAdd = async (pluginName: string) => {
    try {
      const systemSdks = await fetchCachedSystemSdks();
      const matchingSdk = systemSdks?.find(sdk => sdk.name === pluginName);
      if (!matchingSdk?.path) return;

      const version = await detectSystemSdkVersion(pluginName, matchingSdk.path);
      await addCustomSdkReference(pluginName, matchingSdk.path, version || 'unknown');
      await useCustomSdkReference(pluginName, matchingSdk.path);
      await hijackSystemSdkPath(pluginName, matchingSdk.path);
    } catch (err) {
      options.notifyError(options.formatError(err, t('market.add_auto_link_error', { name: pluginName })));
    }
  };

  const addPlugin = async (pluginName: string) => {
    addingPlugin.value = pluginName;
    try {
      await queue.enqueue({
        kind: 'plugin',
        title: pluginName,
        unitsTotal: 1,
        run: async (task) => {
          const off = onRuntimeEvent<string>('vfox-log', (log) => {
            const percent = extractProgressPercent(log);
            if (percent !== null) {
              task.unitPercent = percent;
              task.unitsDone = 1;
            }
            const speed = extractDownloadSpeed(log);
            if (speed) task.speed = speed;
          });
          try {
            await addVfoxPlugin(pluginName);
            await options.fetchPlugins();
            await linkSystemSdkAfterAdd(pluginName);
            refreshSystemSdkCache();

            if (options.selectedPlugin.value?.name === pluginName) {
              await options.fetchPluginVersions(pluginName);
            }
          } finally {
            off();
          }
        },
      });
    } catch (err) {
      options.notifyTaskError(err, t('market.add_error', { name: pluginName }));
    } finally {
      if (addingPlugin.value === pluginName) {
        addingPlugin.value = null;
      }
    }
  };

  const requestRemovePlugin = async (pluginName: string) => {
    try {
      const nonVfoxMap = await fetchNonVfoxSdks();
      const customSdks = nonVfoxMap[pluginName] || [];
      confirmRemoveCustomSdks.value = customSdks.map((sdk: any) => ({
        path: sdk.path || sdk.Path || '',
        version: sdk.versions?.[0]?.version || sdk.version || sdk.Version || 'unknown',
      }));
    } catch (err) {
      confirmRemoveCustomSdks.value = [];
      options.notifyError(options.formatError(err, t('market.custom_refs_error', { name: pluginName })));
    }
    confirmRemove.value = pluginName;
  };

  const executeRemovePlugin = async (chosenPath: string | null) => {
    if (!confirmRemove.value) return;
    const pluginName = confirmRemove.value;
    confirmRemove.value = null;

    removingPlugin.value = pluginName;
    try {
      await options.runTask(t('task.plugin.remove', { name: pluginName }), async () => {
        await removePlugin(pluginName, chosenPath || '');
        await options.fetchPlugins();
        refreshSystemSdkCache();
      });
    } catch (err) {
      options.notifyTaskError(err, t('market.remove_error', { name: pluginName }));
    } finally {
      removingPlugin.value = null;
    }
  };

  const installVersion = async (pluginName: string, version: string) => {
    installingVersion.value = version;
    try {
      await queue.enqueue({
        kind: 'plugin',
        title: `${pluginName}@${version}`,
        unitsTotal: 1,
        run: async (task) => {
          const off = onRuntimeEvent<string>('vfox-log', (log) => {
            const percent = extractProgressPercent(log);
            if (percent !== null) {
              task.unitPercent = percent;
              task.unitsDone = 1;
            }
            const speed = extractDownloadSpeed(log);
            if (speed) task.speed = speed;
          });
          try {
            await installPluginVersion(pluginName, version);
            await options.fetchPlugins();
            if (options.selectedPlugin.value?.name === pluginName) {
              await options.fetchPluginVersions(pluginName);
            }
            options.installedVersions.value.add(version);
          } finally {
            off();
          }
        },
      });
    } catch (err) {
      options.notifyTaskError(err, t('market.install_error', { name: pluginName, version }));
    } finally {
      if (installingVersion.value === version) {
        installingVersion.value = null;
      }
    }
  };

  return {
    addingPlugin,
    removingPlugin,
    installingVersion,
    confirmRemove,
    confirmRemoveCustomSdks,
    addPlugin,
    requestRemovePlugin,
    executeRemovePlugin,
    installVersion,
  };
};
