<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { t } from '../../i18n';
import type { GitHubSource, GitHubSourceSettings } from '../../services/appModels';
import {
  fetchGitHubSourceSettings,
  saveGitHubSourceSettings,
} from '../../services/githubSourceService';
import { getErrorMessage, useSettingsNotice } from '../../composables/useSettingsNotice';

const emit = defineEmits(['notify']);
const { notifyError, notifySuccess } = useSettingsNotice((notice) => {
  emit('notify', notice);
});

const settings = ref<GitHubSourceSettings | null>(null);
const loading = ref(true);
const saving = ref(false);

const customNameInput = ref('');
const customUrlInput = ref('');
const addingCustom = ref(false);

const selectedID = computed({
  get: () => settings.value?.selectedId ?? '',
  set: (value: string) => {
    if (settings.value) {
      settings.value.selectedId = value;
    }
  },
});

const toggleEnabled = () => {
  if (!settings.value) return;
  settings.value.enabled = !settings.value.enabled;
};

const selectedSource = computed<GitHubSource | null>(() => {
  if (!settings.value) return null;
  return settings.value.sources.find((s) => s.id === selectedID.value) ?? null;
});

const canAddCustom = computed(() => {
  const name = customNameInput.value.trim();
  const url = customUrlInput.value.trim();
  return name !== '' && url !== '' && (url.startsWith('http://') || url.startsWith('https://'));
});

const loadSettings = async () => {
  try {
    settings.value = await fetchGitHubSourceSettings();
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.github_source.load_error')));
  } finally {
    loading.value = false;
  }
};

const doSave = async () => {
  if (!settings.value || saving.value) return;
  saving.value = true;
  try {
    settings.value = await saveGitHubSourceSettings(settings.value);
    notifySuccess(t('settings.github_source.save_success'));
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.github_source.save_error')));
    await loadSettings();
  } finally {
    saving.value = false;
  }
};

const generateCustomSourceID = () => {
  const ts = Date.now();
  const rand = Math.floor(Math.random() * 0xffffff)
    .toString(16)
    .padStart(6, '0');
  return `custom-${ts}-${rand}`;
};

const addCustomSource = () => {
  if (!canAddCustom.value || !settings.value) return;
  const id = generateCustomSourceID();
  const newSource: GitHubSource = {
    id,
    name: customNameInput.value.trim(),
    url: customUrlInput.value.trim(),
    isPreset: false,
  };
  settings.value.customSources = [...(settings.value.customSources ?? []), newSource];
  settings.value.sources = [...settings.value.sources, newSource];
  settings.value.selectedId = id;
  customNameInput.value = '';
  customUrlInput.value = '';
  addingCustom.value = false;
};

const removeCustomSource = (id: string) => {
  if (!settings.value) return;
  settings.value.customSources = settings.value.customSources.filter((s) => s.id !== id);
  settings.value.sources = settings.value.sources.filter((s) => s.id !== id);
  if (settings.value.selectedId === id) {
    settings.value.selectedId = settings.value.presetSources[0]?.id ?? '';
  }
};

onMounted(loadSettings);
</script>

<template>
  <section class="set-group">
    <h2>{{ t('settings.github_source.title') }}</h2>
    <div class="set-list">
      <div class="set-row">
        <div class="left">
          <span class="label">{{ t('settings.github_source.name') }}</span>
          <span class="desc">{{ t('settings.github_source.desc') }}</span>
        </div>
        <label class="switch">
          <input
            type="checkbox"
            :checked="settings?.enabled ?? false"
            :disabled="loading || saving"
            :aria-label="t('settings.github_source.name')"
            @change="toggleEnabled"
          >
          <span class="slider"></span>
        </label>
      </div>

      <div v-if="loading" class="github-source-loading">
        <div class="spinner"></div>
      </div>

      <template v-else-if="settings">
        <div class="set-row stack">
          <div class="left">
            <span class="label">{{ t('settings.github_source.source') }}</span>
            <span class="desc">{{ t('settings.github_source.source.desc') }}</span>
          </div>
          <select
            v-model="selectedID"
            class="source-select"
            :disabled="!settings.enabled || saving"
          >
            <optgroup :label="t('settings.github_source.presets')">
              <option
                v-for="source in settings.presetSources"
                :key="source.id"
                :value="source.id"
              >
                {{ source.name }}
              </option>
            </optgroup>
            <optgroup v-if="settings.customSources.length" :label="t('settings.github_source.custom')">
              <option
                v-for="source in settings.customSources"
                :key="source.id"
                :value="source.id"
              >
                {{ source.name }}
              </option>
            </optgroup>
          </select>
        </div>

        <div v-if="selectedSource" class="source-detail">
          <span class="source-url">{{ selectedSource.url }}</span>
          <span v-if="selectedSource.isPreset" class="source-tag preset">
            {{ t('settings.github_source.preset_tag') }}
          </span>
          <span v-else class="source-tag custom">
            {{ t('settings.github_source.custom_tag') }}
          </span>
        </div>

        <div class="set-row stack">
          <div class="left">
            <span class="label">{{ t('settings.github_source.custom_add') }}</span>
            <span class="desc">{{ t('settings.github_source.custom_add.desc') }}</span>
          </div>
          <div class="custom-source-form">
            <input
              v-model="customNameInput"
              class="path-input"
              :placeholder="t('settings.github_source.custom_name_placeholder')"
              :disabled="saving"
            >
            <input
              v-model="customUrlInput"
              class="path-input"
              :placeholder="t('settings.github_source.custom_url_placeholder')"
              :disabled="saving"
            >
            <button
              class="btn primary"
              :disabled="!canAddCustom || saving"
              @click="addCustomSource"
            >
              {{ t('settings.github_source.custom_add_button') }}
            </button>
          </div>
        </div>

        <div v-if="settings.customSources.length" class="set-row stack">
          <span class="block-label">{{ t('settings.github_source.saved_custom') }}</span>
          <div class="custom-source-list">
            <div
              v-for="source in settings.customSources"
              :key="source.id"
              class="custom-source-item"
            >
              <div class="custom-source-info">
                <span class="custom-source-name">{{ source.name }}</span>
                <span class="custom-source-url">{{ source.url }}</span>
              </div>
              <button
                class="btn text small"
                :disabled="saving"
                @click="removeCustomSource(source.id)"
              >
                {{ t('settings.github_source.remove') }}
              </button>
            </div>
          </div>
        </div>

        <div class="set-row">
          <div class="left">
            <span class="desc">{{ t('settings.github_source.save_hint') }}</span>
          </div>
          <button
            class="btn primary"
            :disabled="saving"
            @click="doSave"
          >
            <div v-if="saving" class="spinner small-spinner"></div>
            <template v-else>{{ t('settings.github_source.save') }}</template>
          </button>
        </div>
      </template>
    </div>
  </section>
</template>

<style scoped>
.github-source-loading {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.source-select {
  width: 100%;
  max-width: 360px;
  padding: 8px 12px;
  border-radius: var(--md-shape-small);
  border: 1px solid var(--border-strong);
  background: var(--bg);
  color: var(--ink);
  font-size: var(--md-typescale-body-medium);
}

.source-detail {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: var(--md-shape-small);
  background: var(--code-bg);
  border: 1px solid var(--hairline);
  margin-top: -4px;
}

.source-url {
  flex: 1 1 auto;
  min-width: 0;
  font-family: var(--font-mono);
  font-size: var(--md-typescale-body-small);
  color: var(--ink-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.source-tag {
  flex: 0 0 auto;
  padding: 2px 8px;
  border-radius: var(--radius-pill);
  font-size: var(--md-typescale-label-medium);
}

.source-tag.preset {
  background: var(--bubble);
  color: var(--ink-2);
}

.source-tag.custom {
  background: var(--primary-container);
  color: var(--on-primary-container);
}

.custom-source-form {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  width: 100%;
}

.custom-source-form .path-input {
  flex: 1 1 200px;
}

.custom-source-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.custom-source-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: var(--md-shape-small);
  border: 1px solid var(--border-strong);
  background: var(--bg);
}

.custom-source-info {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-width: 0;
  gap: 2px;
}

.custom-source-name {
  font-weight: 500;
  color: var(--ink);
}

.custom-source-url {
  font-family: var(--font-mono);
  font-size: var(--md-typescale-body-small);
  color: var(--ink-3);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
