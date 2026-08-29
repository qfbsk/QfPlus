<script lang="ts" setup>
import { computed, inject, onMounted, onUnmounted, ref } from 'vue';
import { t } from '../../i18n';
import type { AppTab } from '../../app/navigation';
import type { ProxyQuickStatus } from '../../services/appModels';
import {
  fetchProxyQuickStatus,
  setProxyQuickEnabled,
  testProxyQuickDelay,
} from '../../services/proxyService';
import { onRuntimeEvent } from '../../services/runtimeService';

const switchTab = inject<((tab: AppTab) => void) | null>('switchTab', null);

const status = ref<ProxyQuickStatus | null>(null);
const busy = ref(false);
const testing = ref(false);
let unsubscribe: (() => void) | null = null;

const isOn = computed(() => !!status.value?.enabled && !!status.value?.running);
const visible = computed(() => !!status.value?.hasConfig);
const exitName = computed(() => status.value?.exitNode || status.value?.exitGroup || '');

const delayText = computed(() => {
  const delay = status.value?.delay ?? 0;
  if (!isOn.value) return '';
  if (testing.value && delay === 0) return t('home.proxy.testing');
  if (delay < 0) return t('home.proxy.timeout');
  if (delay === 0) return '—';
  return `${delay}ms`;
});

const delayClass = computed(() => {
  const delay = status.value?.delay ?? 0;
  if (delay < 0) return 'fail';
  if (delay > 0 && delay < 300) return 'good';
  if (delay >= 300 && delay < 800) return 'mid';
  if (delay >= 800) return 'slow';
  return 'idle';
});

const refresh = async (withDelay = false) => {
  try {
    status.value = withDelay ? await testProxyQuickDelay() : await fetchProxyQuickStatus();
  } catch {
    status.value = null;
  }
};

const toggle = async () => {
  if (busy.value || !status.value) return;
  busy.value = true;
  try {
    status.value = await setProxyQuickEnabled(!status.value.enabled);
  } finally {
    busy.value = false;
  }
  measure();
};

const measure = async () => {
  if (testing.value || !isOn.value) return;
  testing.value = true;
  try {
    status.value = await testProxyQuickDelay();
  } finally {
    testing.value = false;
  }
};

const openSettings = () => switchTab?.('settings');

const load = async () => {
  await refresh();
  measure();
};

onMounted(() => {
  load();
  unsubscribe = onRuntimeEvent('proxy-status-changed', load);
});

onUnmounted(() => unsubscribe?.());
</script>

<template>
  <div v-if="visible" class="proxy-quick" :class="{ on: isOn }">
    <button
      class="quick-hotspot"
      :title="isOn ? `${status?.exitGroup} → ${exitName}` : t('home.proxy.off')"
      @click="openSettings"
    >
      <span class="material-symbols-outlined quick-icon">
        {{ isOn ? 'cloud_done' : 'cloud_off' }}
      </span>
      <span class="quick-label">
        <template v-if="isOn">{{ exitName || t('home.proxy.on') }}</template>
        <template v-else>{{ t('home.proxy.off') }}</template>
      </span>
      <span
        v-if="isOn"
        class="quick-delay"
        :class="delayClass"
        :title="t('home.proxy.test_delay')"
        @click.stop="measure"
      >{{ delayText }}</span>
    </button>

    <label class="switch small quick-switch" :title="t('home.proxy.toggle')">
      <input
        type="checkbox"
        :checked="isOn"
        :disabled="busy"
        :aria-label="t('home.proxy.toggle')"
        @change="toggle"
      >
      <span class="slider"></span>
    </label>
  </div>
</template>

<style scoped>
.proxy-quick {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  max-width: 420px;
  padding: 3px 6px 3px 3px;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border);
  background: var(--bg);
  transition: border-color var(--t), background-color var(--t);
}

.proxy-quick.on {
  border-color: var(--border-strong);
  background: var(--bubble);
}

.quick-hotspot {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  padding: 4px 8px;
  border: none;
  background: transparent;
  color: var(--ink-2);
  font-size: var(--md-typescale-body-medium);
  cursor: pointer;
  border-radius: var(--md-shape-full);
}

.proxy-quick.on .quick-hotspot {
  color: var(--ink);
}

.quick-hotspot:hover {
  background: var(--hover);
  color: var(--ink);
}

.quick-icon {
  font-size: 18px;
  line-height: 1;
}

.quick-label {
  min-width: 0;
  max-width: 220px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.quick-delay {
  flex-shrink: 0;
  padding: 1px 8px;
  border-radius: var(--md-shape-full);
  font-size: var(--md-typescale-label-medium);
  font-variant-numeric: tabular-nums;
  background: var(--bg);
  color: var(--ink-2);
  cursor: pointer;
}

.quick-delay.good {
  color: var(--ink);
  font-weight: 600;
}

.quick-delay.slow,
.quick-delay.fail {
  color: var(--ink-3);
}

.quick-delay.idle {
  color: var(--ink-3);
}
</style>
