<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { t } from '../../i18n';
import { useEnvironmentStatus } from '../../composables/useEnvironmentStatus';

const emit = defineEmits(['refreshed']);

const { loading, report, error, load, openDiagnostic } = useEnvironmentStatus();

onMounted(async () => {
  await load();
  emit('refreshed');
});

const handleRefresh = async () => {
  await load();
  emit('refreshed');
};

const handleDiagnostic = async () => {
  try {
    await openDiagnostic();
  } catch (err) {
    // Errors are surfaced through the composable's error ref.
  }
};

const stateClass = (state?: string) => {
  switch (state) {
    case 'ok': return 'state-ok';
    case 'managed': return 'state-managed';
    case 'broken': return 'state-broken';
    case 'missing': return 'state-missing';
    case 'unmanaged': return 'state-unmanaged';
    default: return 'state-unknown';
  }
};

const stateLabel = (state?: string) => {
  switch (state) {
    case 'ok': return t('environment.state.ok');
    case 'managed': return t('environment.state.managed');
    case 'broken': return t('environment.state.broken');
    case 'missing': return t('environment.state.missing');
    case 'unmanaged': return t('environment.state.unmanaged');
    default: return t('environment.state.unknown');
  }
};

const tooltip = (state?: string) => {
  switch (state) {
    case 'ok': return t('environment.state.tooltip.ok');
    case 'managed': return t('environment.state.tooltip.managed');
    case 'broken': return t('environment.state.tooltip.broken');
    case 'missing': return t('environment.state.tooltip.missing');
    case 'unmanaged': return t('environment.state.tooltip.unmanaged');
    default: return t('environment.state.tooltip.unknown');
  }
};

// Only non-ideal states get the circle + exclamation affordance.
const needsDetail = (state?: string) =>
  state === 'broken' || state === 'missing' || state === 'unmanaged';

const expanded = ref<Set<number>>(new Set());
const toggleDetail = (index: number) => {
  const next = new Set(expanded.value);
  if (next.has(index)) next.delete(index);
  else next.add(index);
  expanded.value = next;
};

const pathLabel = (item: { onUserPath: boolean; onMachinePath: boolean }) => {
  if (item.onUserPath && item.onMachinePath) return t('environment.detail.path_both');
  if (item.onUserPath) return t('environment.detail.path_user');
  if (item.onMachinePath) return t('environment.detail.path_machine');
  return t('environment.detail.path_none');
};
</script>

<template>
  <section class="environment-status-card">
    <div class="environment-card-header">
      <h3>{{ t('environment.status.title') }}</h3>
      <div class="environment-card-actions">
        <button class="btn text small" :disabled="loading" @click="handleRefresh">
          <span class="material-symbols-outlined">refresh</span>
          {{ t('environment.status.refresh') }}
        </button>
        <button class="btn tonal small" :disabled="loading" @click="handleDiagnostic">
          <span class="material-symbols-outlined">terminal</span>
          {{ t('environment.status.detect') }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="environment-loading">
      <div class="spinner"></div>
      <span>{{ t('environment.status.loading') }}</span>
    </div>

    <div v-else-if="error" class="environment-error">{{ error }}</div>

    <div v-else-if="report" class="environment-status-body">
      <div class="environment-status-summary">
        <div class="environment-status-row">
          <span class="environment-status-key">{{ t('environment.status.vfox_home') }}</span>
          <span class="environment-status-value">{{ report.vfoxHome }}</span>
        </div>
        <div class="environment-status-row">
          <span class="environment-status-key">{{ t('environment.status.shim_dir') }}</span>
          <span class="environment-status-value">{{ report.shimDir }}</span>
        </div>
        <div class="environment-status-row">
          <span class="environment-status-key">{{ t('environment.status.vfox_in_path') }}</span>
          <span class="environment-status-value">{{ report.vfoxInPath ? t('common.yes') : t('common.no') }}</span>
        </div>
        <div v-if="report.pathDrift" class="environment-status-row warning">
          <span class="environment-status-key">{{ t('environment.status.path_drift') }}</span>
          <span class="environment-status-value">{{ t('common.yes') }}</span>
        </div>
      </div>

      <ul class="environment-status-items">
        <li
          v-for="(item, index) in report.items"
          :key="index"
          class="environment-status-item"
          :class="stateClass(item.state)"
        >
          <div class="environment-status-row-main">
            <span class="sdk-state" :class="stateClass(item.state)" :title="tooltip(item.state)"></span>
            <span class="environment-status-alias">{{ item.alias }}</span>
            <span class="environment-status-source">{{ item.sdkName }}</span>
            <span class="environment-status-version">{{ item.version || '-' }}</span>
            <span class="environment-status-state">{{ stateLabel(item.state) }}</span>
            <button
              v-if="needsDetail(item.state)"
              class="status-detail-toggle"
              type="button"
              :title="t('environment.detail.toggle')"
              :aria-expanded="expanded.has(index)"
              @click="toggleDetail(index)"
            >
              <svg class="detail-icon" viewBox="0 0 24 24" aria-hidden="true">
                <circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2" />
                <path d="M12 7v6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
                <circle cx="12" cy="16.6" r="1.1" fill="currentColor" />
              </svg>
            </button>
          </div>

          <div v-if="expanded.has(index)" class="environment-status-detail">
            <div class="detail-line" v-if="item.notes && item.notes.length">
              <span class="detail-key">{{ t('environment.detail.reason') }}</span>
              <span class="detail-val">
                <span v-for="(note, n) in item.notes" :key="n" class="detail-note">{{ note }}</span>
              </span>
            </div>
            <div class="detail-line">
              <span class="detail-key">{{ t('environment.detail.executable') }}</span>
              <span class="detail-val">{{ item.exePath || '-' }}</span>
            </div>
            <div class="detail-line">
              <span class="detail-key">{{ t('environment.detail.version') }}</span>
              <span class="detail-val">{{ item.version || '-' }}</span>
            </div>
            <div class="detail-line">
              <span class="detail-key">{{ t('environment.detail.on_path') }}</span>
              <span class="detail-val">{{ pathLabel(item) }}</span>
            </div>
            <div class="detail-line">
              <span class="detail-key">{{ t('environment.detail.managed_by') }}</span>
              <span class="detail-val">{{ item.managedBy || '-' }}</span>
            </div>
            <div class="detail-line">
              <span class="detail-key">{{ t('environment.detail.source') }}</span>
              <span class="detail-val">{{ item.source || '-' }}</span>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <div v-else class="environment-empty">{{ t('environment.status.empty') }}</div>
  </section>
</template>

<style scoped>
.environment-status-card {
  width: 100%;
}

.environment-status-items {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.environment-status-item {
  border-radius: var(--md-shape-small);
  padding: 8px 10px;
}

.environment-status-item:hover {
  background: var(--hover);
}

.environment-status-row-main {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.environment-status-alias {
  font-weight: 600;
  color: var(--ink);
}

.environment-status-source {
  color: var(--ink-3);
  font-size: var(--md-typescale-body-small);
}

.environment-status-version {
  margin-left: auto;
  color: var(--ink-2);
  font-variant-numeric: tabular-nums;
}

.environment-status-state {
  min-width: 56px;
  text-align: right;
  color: var(--ink-2);
  font-size: var(--md-typescale-body-small);
}

.status-detail-toggle {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 2px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-2);
}

.status-detail-toggle:hover {
  color: var(--ink);
}

.detail-icon {
  width: 18px;
  height: 18px;
}

.environment-status-item.state-broken .detail-icon {
  color: var(--danger);
}

.environment-status-item.state-missing .detail-icon {
  color: #d97706;
}

.environment-status-item.state-unmanaged .detail-icon {
  color: var(--info, #2563eb);
}

.environment-status-detail {
  margin-top: 8px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--md-shape-small);
  background: var(--bg-side, var(--bg));
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-line {
  display: flex;
  gap: 10px;
  font-size: var(--md-typescale-body-small);
  line-height: 1.5;
}

.detail-key {
  flex-shrink: 0;
  width: 76px;
  color: var(--ink-3);
}

.detail-val {
  flex: 1;
  min-width: 0;
  color: var(--ink);
  word-break: break-all;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
</style>
