import { computed, onMounted, ref, watch } from 'vue';
import { t } from '../i18n';
import {
  PendingMigrationPlan,
  shouldAskDownloadPathMigration,
} from './downloadPathMigration';
import {
  DownloadPathInfo,
  fetchDownloadPathInfo,
  planEnvironmentMigration,
  resetDownloadPath,
  resetDownloadPathWithMigration,
  saveDownloadPath,
  saveDownloadPathWithMigration,
  selectDownloadPath,
} from '../services/settingsService';
import {
  SettingsNotice,
  getErrorMessage,
  useSettingsNotice,
} from './useSettingsNotice';

export const useDownloadPathSettings = (notify: (notice: SettingsNotice) => void) => {
  const downloadPathInfo = ref<DownloadPathInfo | null>(null);
  const downloadPathInput = ref('');
  const loadingDownloadPath = ref(true);
  const savingDownloadPath = ref(false);
  const selectingDownloadPath = ref(false);
  const resettingDownloadPath = ref(false);
  const pendingPlan = ref<PendingMigrationPlan | null>(null);
  const { notifyError, notifySuccess } = useSettingsNotice(notify);

  const syncDownloadPathInfo = (info: DownloadPathInfo) => {
    downloadPathInfo.value = info;
    downloadPathInput.value = info.path;
  };

  const shouldAskMigration = (targetPath: string) => {
    return shouldAskDownloadPathMigration(downloadPathInfo.value, targetPath);
  };

  const isDownloadPathBusy = computed(() => savingDownloadPath.value || resettingDownloadPath.value);

  const loadDownloadPath = async () => {
    loadingDownloadPath.value = true;
    try {
      syncDownloadPathInfo(await fetchDownloadPathInfo());
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.load_error')));
    } finally {
      loadingDownloadPath.value = false;
    }
  };

  const openPlanModal = async (targetPath: string, mode: 'save' | 'reset') => {
    savingDownloadPath.value = true;
    try {
      const plan = await planEnvironmentMigration(targetPath);
      if (plan.blockingReason) {
        notifyError(plan.blockingReason);
        return;
      }
      pendingPlan.value = { mode, targetPath, plan };
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.error')));
      // Fall back to a direct migration so existing data is never stranded.
      if (mode === 'save') {
        await runSaveWithMigration(targetPath);
      } else {
        await runResetWithMigration();
      }
    } finally {
      savingDownloadPath.value = false;
    }
  };

  const requestSaveDownloadPath = async () => {
    const targetPath = downloadPathInput.value.trim();
    if (shouldAskMigration(targetPath)) {
      await openPlanModal(targetPath, 'save');
      return;
    }

    pendingPlan.value = null;
    savingDownloadPath.value = true;
    try {
      syncDownloadPathInfo(await saveDownloadPath(targetPath));
      notifySuccess(t('settings.download.path.success'));
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.error')));
    } finally {
      savingDownloadPath.value = false;
    }
  };

  const chooseDownloadPath = async () => {
    selectingDownloadPath.value = true;
    try {
      const selected = await selectDownloadPath();
      if (selected && selected.trim()) {
        downloadPathInput.value = selected;
      }
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.select_error')));
    } finally {
      selectingDownloadPath.value = false;
    }
  };

  const requestResetDownloadPath = async () => {
    const targetPath = downloadPathInfo.value?.defaultPath || '';
    if (shouldAskMigration(targetPath)) {
      await openPlanModal(targetPath, 'reset');
      return;
    }

    pendingPlan.value = null;
    resettingDownloadPath.value = true;
    try {
      syncDownloadPathInfo(await resetDownloadPath());
      notifySuccess(t('settings.download.path.reset_success'));
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.error')));
    } finally {
      resettingDownloadPath.value = false;
    }
  };

  const runSaveWithMigration = async (targetPath: string) => {
    savingDownloadPath.value = true;
    try {
      syncDownloadPathInfo(await saveDownloadPathWithMigration(targetPath));
      pendingPlan.value = null;
      notifySuccess(t('settings.download.path.migrate_success'));
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.error')));
    } finally {
      savingDownloadPath.value = false;
    }
  };

  const runResetWithMigration = async () => {
    resettingDownloadPath.value = true;
    try {
      syncDownloadPathInfo(await resetDownloadPathWithMigration());
      pendingPlan.value = null;
      notifySuccess(t('settings.download.path.migrate_success'));
    } catch (err) {
      notifyError(getErrorMessage(err, t('settings.download.path.error')));
    } finally {
      resettingDownloadPath.value = false;
    }
  };

  const confirmPendingMigration = async () => {
    const plan = pendingPlan.value;
    if (!plan || isDownloadPathBusy.value) return;
    if (plan.mode === 'save') {
      await runSaveWithMigration(plan.targetPath);
    } else {
      await runResetWithMigration();
    }
  };

  const cancelPendingMigration = () => {
    if (isDownloadPathBusy.value) return;
    pendingPlan.value = null;
  };

  watch(downloadPathInput, (value) => {
    if (pendingPlan.value?.mode === 'save' && value.trim() !== pendingPlan.value.targetPath) {
      pendingPlan.value = null;
    }
  });

  onMounted(() => {
    loadDownloadPath();
  });

  return {
    downloadPathInfo,
    downloadPathInput,
    loadingDownloadPath,
    savingDownloadPath,
    selectingDownloadPath,
    resettingDownloadPath,
    pendingPlan,
    isDownloadPathBusy,
    chooseDownloadPath,
    requestSaveDownloadPath,
    requestResetDownloadPath,
    confirmPendingMigration,
    cancelPendingMigration,
  };
};
