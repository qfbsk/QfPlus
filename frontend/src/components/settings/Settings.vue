<script lang="ts" setup>
import { computed, inject, Ref } from 'vue';
import { t } from '../../i18n';
import { useDownloadPathSettings } from '../../composables/useDownloadPathSettings';
import AppearanceSettings from './AppearanceSettings.vue';
import CoreSettings from './CoreSettings.vue';
import DownloadPathSettings from './DownloadPathSettings.vue';
import GitHubSourceSettings from './GitHubSourceSettings.vue';
import MigrationPlanModal from './MigrationPlanModal.vue';
import ProxySettings from './ProxySettings.vue';

const showTerminal = inject<Ref<boolean>>('showTerminal');
const emit = defineEmits(['notify']);

const terminalVisible = computed({
  get: () => showTerminal?.value ?? false,
  set: (value: boolean) => {
    if (showTerminal) {
      showTerminal.value = value;
    }
  },
});

const downloadPathSettings = useDownloadPathSettings((notice) => {
  emit('notify', notice);
});
</script>

<template>
  <div class="settings-view view-container">
    <h2>{{ t('settings.title') }}</h2>

    <AppearanceSettings v-model:show-terminal="terminalVisible" />
    <ProxySettings @notify="emit('notify', $event)" />
    <GitHubSourceSettings @notify="emit('notify', $event)" />
    <DownloadPathSettings
      v-model:download-path-input="downloadPathSettings.downloadPathInput.value"
      :download-path-info="downloadPathSettings.downloadPathInfo.value"
      :loading-download-path="downloadPathSettings.loadingDownloadPath.value"
      :saving-download-path="downloadPathSettings.savingDownloadPath.value"
      :selecting-download-path="downloadPathSettings.selectingDownloadPath.value"
      :resetting-download-path="downloadPathSettings.resettingDownloadPath.value"
      @save="downloadPathSettings.requestSaveDownloadPath"
      @choose="downloadPathSettings.chooseDownloadPath"
      @reset="downloadPathSettings.requestResetDownloadPath"
    />
    <CoreSettings @notify="emit('notify', $event)" />

    <MigrationPlanModal
      :pending-plan="downloadPathSettings.pendingPlan.value"
      :is-download-path-busy="downloadPathSettings.isDownloadPathBusy.value"
      @confirm="downloadPathSettings.confirmPendingMigration"
      @cancel="downloadPathSettings.cancelPendingMigration"
    />
  </div>
</template>

<style scoped>
.settings-view {
  display: flex;
  flex-direction: column;
}
</style>
