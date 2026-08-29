<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { t } from '../../i18n';
import type { ProxyGroup, ProxyStatus } from '../../services/appModels';
import {
  fetchProxyStatus,
  setProxyEnabled,
  importProxySubscription,
  fetchProxyGroups,
  selectProxyNode,
  setProxyGroup,
  testProxyGroupDelay,
} from '../../services/proxyService';
import { getErrorMessage, useSettingsNotice } from '../../composables/useSettingsNotice';
import ProxyNodeList from './ProxyNodeList.vue';

const emit = defineEmits(['notify']);
const { notifyError, notifySuccess } = useSettingsNotice((notice) => {
  emit('notify', notice);
});

const proxyStatus = ref<ProxyStatus | null>(null);
const groups = ref<ProxyGroup[]>([]);
const subscriptionInput = ref('');
const loading = ref(true);
const toggling = ref(false);
const importing = ref(false);
const testingDelay = ref(false);
const selectingNode = ref(false);
const switchingGroup = ref('');

const statusState = computed(() => {
  if (proxyStatus.value?.enabled && proxyStatus.value?.running) return 'on';
  if (!proxyStatus.value?.enabled) return 'off';
  return 'error';
});

const activeGroup = computed<ProxyGroup | null>(() => {
  const selected = proxyStatus.value?.selectedGroup;
  if (selected) {
    const found = groups.value.find((g) => g.name === selected);
    if (found) return found;
  }
  return groups.value[0] ?? null;
});

const exitNode = computed(() => {
  if (!activeGroup.value) return '';
  return activeGroup.value.now || proxyStatus.value?.selectedNode || '';
});

const isSelectorGroup = (group: ProxyGroup) => group.type === 'Selector';

const loadStatus = async () => {
  try {
    proxyStatus.value = await fetchProxyStatus();
    if (proxyStatus.value?.subscriptionUrl) {
      subscriptionInput.value = proxyStatus.value.subscriptionUrl;
    }
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.proxy.load_error')));
  } finally {
    loading.value = false;
  }
};

const loadGroups = async () => {
  if (!proxyStatus.value?.running) {
    groups.value = [];
    return;
  }
  try {
    groups.value = await fetchProxyGroups();
  } catch {
    groups.value = [];
  }
};

const toggleProxy = async () => {
  if (!proxyStatus.value) return;
  toggling.value = true;
  try {
    const enabled = !proxyStatus.value.enabled;
    proxyStatus.value = await setProxyEnabled(enabled);
    if (enabled) {
      await loadGroups();
      notifySuccess(t('settings.proxy.enabled_success'));
    } else {
      groups.value = [];
      notifySuccess(t('settings.proxy.disabled_success'));
    }
    if (proxyStatus.value?.error) {
      notifyError(proxyStatus.value.error);
    }
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.proxy.toggle_error')));
    await loadStatus();
  } finally {
    toggling.value = false;
  }
};

const doImport = async () => {
  const url = subscriptionInput.value.trim();
  if (!url) return;
  importing.value = true;
  try {
    proxyStatus.value = await importProxySubscription(url);
    await loadGroups();
    notifySuccess(t('settings.proxy.import_success'));
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.proxy.import_error')));
  } finally {
    importing.value = false;
  }
};

const switchGroup = async (group: ProxyGroup) => {
  if (switchingGroup.value || !proxyStatus.value) return;
  if (group.name === proxyStatus.value.selectedGroup) return;
  switchingGroup.value = group.name;
  try {
    proxyStatus.value = await setProxyGroup(group.name);
    await loadGroups();
    notifySuccess(t('settings.proxy.group_switch_success'));
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.proxy.group_switch_error')));
    await loadStatus();
    await loadGroups();
  } finally {
    switchingGroup.value = '';
  }
};

const doSelectNode = async (node: string) => {
  const group = activeGroup.value;
  if (!group || !isSelectorGroup(group)) return;
  if (node === group.now) return;
  selectingNode.value = true;
  try {
    proxyStatus.value = await selectProxyNode(group.name, node);
    await loadGroups();
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.proxy.select_error')));
  } finally {
    selectingNode.value = false;
  }
};

const doTestDelays = async () => {
  const group = activeGroup.value;
  if (!group) return;
  testingDelay.value = true;
  try {
    groups.value = await testProxyGroupDelay(group.name);
  } catch (err) {
    notifyError(getErrorMessage(err, t('settings.proxy.delay_error')));
  } finally {
    testingDelay.value = false;
  }
};

onMounted(() => {
  loadStatus().then(loadGroups);
});
</script>

<template>
  <section class="set-group">
    <h2>{{ t('settings.proxy.title') }}</h2>
    <div class="set-list">
      <div class="set-row">
        <div class="left">
          <span class="label">{{ t('settings.proxy.name') }}</span>
          <span class="desc">{{ t('settings.proxy.desc') }}</span>
        </div>
        <label class="switch">
          <input
            type="checkbox"
            :checked="proxyStatus?.enabled ?? false"
            :disabled="toggling || loading"
            :aria-label="t('settings.proxy.name')"
            @change="toggleProxy"
          >
          <span class="slider"></span>
        </label>
      </div>

      <div v-if="loading" class="proxy-loading">
        <div class="spinner"></div>
      </div>

      <template v-else>
        <div class="set-row">
          <div class="left">
            <span class="label">{{ t('settings.proxy.current_exit') }}</span>
            <span v-if="statusState === 'on' && proxyStatus?.selectedGroup" class="desc current-exit" :title="`${proxyStatus.selectedGroup}${exitNode ? ' → ' + exitNode : ''}`">
              {{ proxyStatus.selectedGroup }}<template v-if="exitNode"> → {{ exitNode }}</template>
            </span>
            <span v-else class="desc">—</span>
          </div>
          <div class="right">
            <span v-if="proxyStatus?.error" class="proxy-error">{{ proxyStatus.error }}</span>
            <span class="status-pill" :class="statusState">
              <span class="status-dot"></span>
              <template v-if="statusState === 'on'">{{ t('settings.proxy.status.on') }}</template>
              <template v-else-if="statusState === 'off'">{{ t('settings.proxy.status.off') }}</template>
              <template v-else>{{ t('settings.proxy.status.error') }}</template>
            </span>
          </div>
        </div>

        <div class="set-row stack">
          <div class="left">
            <span class="label">{{ t('settings.proxy.subscription') }}</span>
            <span class="desc">{{ t('settings.proxy.subscription.desc') }}</span>
          </div>
          <div class="path-input-row">
            <input
              v-model="subscriptionInput"
              class="path-input"
              :placeholder="t('settings.proxy.subscription_placeholder')"
              :disabled="importing"
              @keyup.enter="doImport"
            >
            <button
              class="btn primary"
              :disabled="importing || !subscriptionInput.trim()"
              @click="doImport"
            >
              <div v-if="importing" class="spinner small-spinner"></div>
              <template v-else>{{ t('settings.proxy.import') }}</template>
            </button>
          </div>
        </div>

        <template v-if="statusState === 'on' && groups.length > 0">
          <div class="set-row stack">
            <span class="block-label">{{ t('settings.proxy.exit_group') }}</span>
            <div class="group-chips">
              <button
                v-for="group in groups"
                :key="group.name"
                class="group-chip"
                :class="{ active: activeGroup?.name === group.name }"
                :disabled="!!switchingGroup"
                :title="group.name"
                @click="switchGroup(group)"
              >
                <span v-if="switchingGroup === group.name" class="spinner small-spinner"></span>
                <template v-else>{{ group.name }}</template>
              </button>
            </div>
          </div>

          <ProxyNodeList
            v-if="activeGroup"
            :group="activeGroup"
            :selectable="isSelectorGroup(activeGroup)"
            :busy="selectingNode || !!switchingGroup"
            :testing-delay="testingDelay"
            @select="doSelectNode"
            @test="doTestDelays"
          />
        </template>

        <p v-else-if="proxyStatus?.enabled && !proxyStatus?.hasConfig" class="proxy-hint">
          {{ t('settings.proxy.no_config') }}
        </p>
      </template>
    </div>
  </section>
</template>

<style scoped>
.proxy-loading {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.current-exit {
  overflow-wrap: anywhere;
}

.proxy-error {
  font-size: var(--md-typescale-body-small);
  color: var(--danger);
}

.group-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.group-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 32px;
  padding: 4px 14px;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-strong);
  background: var(--bg);
  color: var(--ink-2);
  font-size: var(--md-typescale-body-medium);
  cursor: pointer;
  transition: background-color var(--t), color var(--t), border-color var(--t);
}

.group-chip:hover:not(:disabled):not(.active) {
  background: var(--hover);
  color: var(--ink);
}

.group-chip.active {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--on-ink);
  font-weight: 600;
}

.group-chip:disabled {
  opacity: 0.55;
  cursor: default;
}

.proxy-hint {
  margin: 0;
  padding: 16px 0;
  font-size: var(--md-typescale-body-medium);
  line-height: var(--md-typescale-body-medium-line-height);
  color: var(--ink-3);
}
</style>
