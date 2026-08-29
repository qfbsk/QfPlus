<script lang="ts" setup>
import { t } from '../../i18n';
import type { DownloadPathInfo } from '../../services/settingsService';

defineProps<{
  downloadPathInfo: DownloadPathInfo | null;
  loadingDownloadPath: boolean;
  savingDownloadPath: boolean;
  selectingDownloadPath: boolean;
  resettingDownloadPath: boolean;
}>();

const downloadPathInput = defineModel<string>('downloadPathInput', { required: true });
const emit = defineEmits(['save', 'choose', 'reset']);
</script>

<template>
  <section class="set-group">
    <h2>{{ t('settings.system') }}</h2>
    <div class="set-list">
      <div class="set-row stack">
        <div class="left">
          <div class="label-line">
            <span class="label">{{ t('settings.download.path') }}</span>
            <span class="status-pill" :class="{ off: downloadPathInfo?.isDefault }">
              {{ downloadPathInfo?.isDefault ? t('settings.download.path.default_state') : t('settings.download.path.custom_state') }}
            </span>
          </div>
          <span class="desc">{{ t('settings.download.path.desc') }}</span>
        </div>

        <div class="path-input-row">
          <input
            v-model="downloadPathInput"
            class="path-input"
            :placeholder="downloadPathInfo?.defaultPath || t('settings.download.path.placeholder')"
            :disabled="loadingDownloadPath || savingDownloadPath || resettingDownloadPath"
            @keyup.enter="emit('save')"
          >
          <button
            class="path-icon-btn"
            :title="t('settings.download.path.browse')"
            :disabled="selectingDownloadPath || savingDownloadPath || resettingDownloadPath"
            @click="emit('choose')"
          >
            <span v-if="!selectingDownloadPath" class="material-symbols-outlined">folder_open</span>
            <div v-else class="spinner small-spinner"></div>
          </button>
          <button
            class="btn primary"
            :disabled="savingDownloadPath || resettingDownloadPath || !downloadPathInput.trim()"
            @click="emit('save')"
          >
            <div v-if="savingDownloadPath" class="spinner small-spinner"></div>
            <template v-else>{{ t('settings.download.path.save') }}</template>
          </button>
          <button
            class="btn outlined"
            :disabled="resettingDownloadPath || savingDownloadPath || downloadPathInfo?.isDefault"
            @click="emit('reset')"
          >
            <div v-if="resettingDownloadPath" class="spinner small-spinner"></div>
            <template v-else>{{ t('settings.download.path.reset') }}</template>
          </button>
        </div>

        <div class="path-meta">
          <div class="path-meta-row">
            <span>{{ t('settings.download.path.current') }}</span>
            <code>{{ downloadPathInfo?.path || '-' }}</code>
          </div>
          <div class="path-meta-row">
            <span>{{ t('settings.download.path.default') }}</span>
            <code>{{ downloadPathInfo?.defaultPath || '-' }}</code>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
