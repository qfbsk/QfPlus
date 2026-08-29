<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { t } from '../../i18n';
import type { CoreInfo, CoreRelease } from '../../services/appModels';
import {
  fetchCoreInfo,
  fetchCoreVersions,
  setCoreAutoUpdate,
  switchCoreVersion,
} from '../../services/coreService';
import { getErrorMessage, useSettingsNotice } from '../../composables/useSettingsNotice';
import { onRuntimeEvent } from '../../services/runtimeService';

const emit = defineEmits(['notify']);
const { notifyError, notifySuccess } = useSettingsNotice((notice) => {
  emit('notify', notice);
});

const info = ref<CoreInfo | null>(null);
const releases = ref<CoreRelease[]>([]);
const loading = ref(true);
const checking = ref(false);
const switching = ref('');
const savingAutoUpdate = ref(false);
const notesExpanded = ref(false);
let unsubscribe: (() => void) | null = null;

const busy = computed(() => !!switching.value);

const stateLabel = computed(() => {
  if (!info.value?.currentVersion) return 'error';
  if (info.value?.updateAvailable) return 'update';
  return 'latest';
});

const headline = computed(() => {
  const current = info.value?.currentVersion;
  if (!current) return t('settings.core.unavailable');
  return `vfox ${current}`;
});

const subtitle = computed(() => {
  if (!info.value) return '';
  const parts = [info.value.osArch];
  if (info.value.usesLocalCore && info.value.bundledVersion) {
    parts.push(`${t('settings.core.bundled')} ${info.value.bundledVersion}`);
  }
  return parts.join(' · ');
});

const latestRelease = computed(() => releases.value[0] ?? null);

const updateNotes = computed(() => latestRelease.value?.notes || info.value?.releaseNotes || '');

const lastCheckText = computed(() => info.value?.lastCheck?.slice(0, 10) ?? '');

const historyVersions = computed(() => releases.value.slice(0, 8));

const loadInfo = async () => {
  try {
    info.value = await fetchCoreInfo();
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.core.load_error')));
  } finally {
    loading.value = false;
  }
};

// loadVersions runs silently on mount: a cold or offline release feed should not
// toast every time the user opens Settings, and 检查更新 covers the explicit path.
const loadVersions = async () => {
  try {
    releases.value = await fetchCoreVersions();
  } catch {
    releases.value = [];
  }
};

const checkUpdates = async () => {
  if (checking.value) return;
  checking.value = true;
  try {
    releases.value = await fetchCoreVersions();
    info.value = await fetchCoreInfo();
    notesExpanded.value = false;
    notifySuccess(t('settings.core.check_success'));
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.core.check_error')));
  } finally {
    checking.value = false;
  }
};

const doSwitch = async (version: string) => {
  if (busy.value) return;
  switching.value = version;
  try {
    info.value = await switchCoreVersion(version);
    if (info.value?.error) {
      notifyError(info.value.error);
    } else {
      notifySuccess(t('settings.core.switch_success', { version }));
    }
    releases.value = releases.value.map((release) => ({
      ...release,
      isCurrent: release.version === info.value?.currentVersion,
    }));
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.core.switch_error')));
    await loadInfo();
  } finally {
    switching.value = '';
  }
};

const toggleAutoUpdate = async () => {
  if (!info.value || savingAutoUpdate.value) return;
  savingAutoUpdate.value = true;
  try {
    info.value = await setCoreAutoUpdate(!info.value.autoUpdate);
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.core.auto_update_error')));
  } finally {
    savingAutoUpdate.value = false;
  }
};

onMounted(() => {
  loadInfo();
  loadVersions();
  unsubscribe = onRuntimeEvent('core-status-changed', () => loadInfo());
});

onUnmounted(() => unsubscribe?.());
</script>

<template>
  <section class="set-group">
    <h2>{{ t('settings.core.title') }}</h2>
    <div class="set-list">
      <div class="set-row">
        <div class="left">
          <div class="label-line">
            <span class="label">{{ headline }}</span>
            <span v-if="info?.usesLocalCore" class="tag current">{{ t('settings.core.downloaded') }}</span>
          </div>
          <span class="desc">{{ subtitle || t('settings.core.desc') }}</span>
        </div>
        <div class="right">
          <span class="status-pill" :class="stateLabel">
            <span class="status-dot"></span>
            <template v-if="stateLabel === 'update'">
              {{ t('settings.core.update_found', { version: info?.latestVersion ?? '' }) }}
            </template>
            <template v-else-if="stateLabel === 'latest'">{{ t('settings.core.up_to_date') }}</template>
            <template v-else>{{ t('settings.core.unavailable') }}</template>
          </span>
          <button class="btn tonal" :disabled="checking || loading" @click="checkUpdates">
            <div v-if="checking" class="spinner small-spinner"></div>
            <template v-else>{{ t('settings.core.check') }}</template>
          </button>
        </div>
      </div>

      <div v-if="loading" class="core-loading">
        <div class="spinner"></div>
      </div>

      <template v-else>
        <p v-if="lastCheckText || info?.error" class="core-meta">
          <span v-if="info?.error" class="core-error">{{ info.error }}</span>
          <span v-else>{{ t('settings.core.last_check') }} {{ lastCheckText }}</span>
        </p>

        <div v-if="stateLabel === 'update' && info?.latestVersion" class="set-row stack">
          <div class="left">
            <span class="label">{{ t('settings.core.update_to', { version: info!.latestVersion }) }}</span>
            <button class="notes-toggle" @click="notesExpanded = !notesExpanded">
              <span class="material-symbols-outlined notes-chevron">{{ notesExpanded ? 'expand_less' : 'expand_more' }}</span>
              {{ t('settings.core.release_notes') }}
            </button>
          </div>
          <div class="right">
            <button class="btn primary" :disabled="busy" @click="doSwitch(info!.latestVersion)">
              <span v-if="switching === info!.latestVersion" class="spinner small-spinner"></span>
              <template v-else>{{ t('settings.core.update_to', { version: info!.latestVersion }) }}</template>
            </button>
          </div>
          <pre v-if="notesExpanded" class="core-notes">{{ updateNotes || t('settings.core.no_notes') }}</pre>
        </div>

        <div v-if="historyVersions.length" class="set-row stack">
          <span class="block-label">{{ t('settings.core.versions') }}</span>
          <div class="core-version-rows">
            <div
              v-for="release in historyVersions"
              :key="release.version"
              class="core-version-row"
              :class="{ current: release.isCurrent }"
            >
              <div class="core-version-main" :title="release.notes">
                <span class="core-version-name">v{{ release.version }}</span>
                <span v-if="release.isCurrent" class="tag current">{{ t('settings.core.current') }}</span>
                <span v-else-if="release.downloaded" class="tag local">{{ t('settings.core.downloaded') }}</span>
                <span v-if="!release.isCurrent && release.version === info?.bundledVersion" class="tag bundled">
                  {{ t('settings.core.bundled') }}
                </span>
                <span v-if="!release.isCurrent && switching === release.version" class="tag working">
                  {{ t('settings.core.switching') }}
                </span>
              </div>
              <span class="core-version-date">{{ release.date }}</span>
              <button
                v-if="!release.isCurrent"
                class="btn text small"
                :disabled="busy"
                @click="doSwitch(release.version)"
              >{{ t('settings.core.switch') }}</button>
            </div>
          </div>

          <button
            v-if="info?.usesLocalCore"
            class="btn outlined small restore-bundled"
            :disabled="busy"
            @click="doSwitch('bundled')"
          >
            <span v-if="switching === 'bundled'" class="spinner small-spinner"></span>
            <template v-else>{{ t('settings.core.restore_bundled') }}</template>
          </button>
        </div>

        <div v-else class="core-hint">
          {{ checking ? t('settings.core.checking') : t('settings.core.versions_hint') }}
        </div>

        <div class="set-row">
          <div class="left">
            <span class="label">{{ t('settings.core.auto_update') }}</span>
            <span class="desc">{{ t('settings.core.auto_update.desc') }}</span>
          </div>
          <label class="switch">
            <input
              type="checkbox"
              :checked="info?.autoUpdate ?? false"
              :disabled="savingAutoUpdate"
              :aria-label="t('settings.core.auto_update')"
              @change="toggleAutoUpdate"
            >
            <span class="slider"></span>
          </label>
        </div>
      </template>
    </div>
  </section>
</template>

<style scoped>
.core-loading {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.core-meta {
  margin: -6px 0 0;
  font-size: var(--md-typescale-body-small);
  line-height: var(--md-typescale-body-small-line-height);
  color: var(--ink-3);
}

.core-error {
  color: var(--danger);
}

.core-hint {
  padding: 16px 0;
  font-size: var(--md-typescale-body-small);
  line-height: var(--md-typescale-body-small-line-height);
  color: var(--ink-3);
}

.notes-toggle {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  align-self: flex-start;
  padding: 0;
  background: transparent;
  color: var(--link);
  font-size: var(--md-typescale-body-small);
  cursor: pointer;
}

.notes-chevron {
  font-size: 16px;
}

.core-notes {
  width: 100%;
  max-height: 200px;
  margin: 0;
  padding: 12px 14px;
  overflow: auto;
  border-radius: var(--md-shape-small);
  background: var(--code-bg);
  border: 1px solid var(--hairline);
  color: var(--ink-2);
  font-family: var(--font-mono);
  font-size: var(--md-typescale-body-small);
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.core-version-rows {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.core-version-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto 72px;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  border-radius: var(--md-shape-small);
  font-size: var(--md-typescale-body-medium);
}

.core-version-row:hover {
  background: var(--hover);
}

.core-version-row.current {
  background: var(--bubble);
}

.core-version-main {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.core-version-name {
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 500;
  font-variant-numeric: tabular-nums;
  color: var(--ink);
}

.core-version-date {
  font-size: var(--md-typescale-label-medium);
  color: var(--ink-3);
  font-variant-numeric: tabular-nums;
}

.restore-bundled {
  align-self: flex-start;
}
</style>
