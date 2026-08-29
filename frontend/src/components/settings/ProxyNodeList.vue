<script lang="ts" setup>
import { computed } from 'vue';
import { t } from '../../i18n';
import type { ProxyGroup, ProxyNode } from '../../services/appModels';

const props = defineProps<{
  group: ProxyGroup;
  selectable: boolean;
  busy: boolean;
  testingDelay: boolean;
}>();

const emit = defineEmits<{
  select: [node: string];
  test: [];
}>();

const groupTypeNames: Record<string, string> = {
  Selector: 'settings.proxy.type.selector',
  URLTest: 'settings.proxy.type.urltest',
  Fallback: 'settings.proxy.type.fallback',
  LoadBalance: 'settings.proxy.type.loadbalance',
};

const groupTypeLabel = computed(() => {
  const key = groupTypeNames[props.group.type];
  return key ? t(key) : props.group.type;
});

const isNestedGroup = (node: ProxyNode) => node.type in groupTypeNames;

const delayClass = (delay: number) => {
  if (delay < 0) return 'delay-fail';
  if (delay === 0) return 'delay-idle';
  if (delay < 300) return 'delay-good';
  if (delay < 800) return 'delay-mid';
  return 'delay-slow';
};

const onRowClick = (node: ProxyNode) => {
  if (!props.selectable || props.busy) return;
  emit('select', node.name);
};
</script>

<template>
  <div class="node-panel">
    <div class="node-panel-header">
      <div class="node-panel-title">
        <span class="node-panel-name">{{ group.name }}</span>
        <span class="node-panel-type">{{ groupTypeLabel }}</span>
        <span class="node-panel-count">{{ group.nodes.length }} {{ t('settings.proxy.nodes_unit') }}</span>
      </div>
      <button
        class="btn outlined small"
        :disabled="testingDelay || busy"
        @click="emit('test')"
      >
        <div v-if="testingDelay" class="spinner small-spinner"></div>
        <template v-else>{{ t('settings.proxy.test_all') }}</template>
      </button>
    </div>

    <p v-if="!selectable" class="auto-hint">
      <span class="material-symbols-outlined auto-hint-icon">speed</span>
      {{ t('settings.proxy.auto_group_hint') }}
    </p>

    <div class="node-rows">
      <div
        v-for="node in group.nodes"
        :key="node.name"
        class="node-row"
        :class="{
          active: node.name === group.now,
          clickable: selectable && !busy,
          nested: isNestedGroup(node),
        }"
        @click="onRowClick(node)"
      >
        <span class="node-indicator">
          <span v-if="node.name === group.now" class="material-symbols-outlined">check_circle</span>
          <span v-else-if="isNestedGroup(node)" class="material-symbols-outlined">folder_open</span>
          <span v-else class="radio-dot"></span>
        </span>
        <span class="node-name" :title="node.name">{{ node.name }}</span>
        <span class="node-type">{{ node.type }}</span>
        <span class="node-delay" :class="delayClass(node.delay)">
          <template v-if="node.delay > 0">{{ node.delay }}ms</template>
          <template v-else-if="node.delay < 0">{{ t('settings.proxy.timeout') }}</template>
          <template v-else>—</template>
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.node-panel {
  display: flex;
  flex-direction: column;
  gap: 10px;
  border: 1px solid var(--hairline);
  border-radius: var(--radius-card);
  background: var(--bg-side);
  padding: 14px;
}

.node-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.node-panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.node-panel-name {
  font-size: var(--md-typescale-body-medium);
  font-weight: 600;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.node-panel-type {
  flex-shrink: 0;
  font-size: var(--md-typescale-label-small);
  font-weight: 500;
  color: var(--ink-2);
  background: var(--bubble);
  padding: 2px 8px;
  border-radius: var(--md-shape-full);
}

.node-panel-count {
  flex-shrink: 0;
  font-size: var(--md-typescale-label-medium);
  color: var(--ink-3);
}

.auto-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
  font-size: var(--md-typescale-body-small);
  line-height: var(--md-typescale-body-small-line-height);
  color: var(--ink-3);
}

.auto-hint-icon {
  font-size: 15px;
}

.node-rows {
  display: flex;
  flex-direction: column;
  gap: 2px;
  max-height: 340px;
  overflow-y: auto;
  padding-right: 2px;
}

.node-row {
  display: grid;
  grid-template-columns: 22px minmax(0, 1fr) 84px 88px;
  align-items: center;
  gap: 10px;
  padding: 7px 10px;
  border-radius: var(--md-shape-small);
  border: 1px solid transparent;
  transition: background-color var(--t);
}

.node-row.clickable {
  cursor: pointer;
}

.node-row.clickable:hover {
  background: var(--hover);
}

.node-row.active {
  background: var(--active);
}

.node-row.active .node-name {
  font-weight: 600;
  color: var(--ink);
}

.node-indicator {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-3);
}

.node-indicator .material-symbols-outlined {
  font-size: 17px;
  color: var(--ink);
}

.radio-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  border: 1.5px solid var(--border-strong);
}

.node-name {
  font-size: var(--md-typescale-body-medium);
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.node-row.nested .node-name {
  font-weight: 600;
}

.node-type {
  font-size: var(--md-typescale-label-medium);
  color: var(--ink-3);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-align: right;
}

/* Latency reads as shape + weight, not colour. */
.node-delay {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 5px;
  font-size: var(--md-typescale-label-medium);
  font-variant-numeric: tabular-nums;
  text-align: right;
  white-space: nowrap;
  color: var(--ink-2);
}

.node-delay::before {
  font-size: 9px;
  line-height: 1;
}

.delay-good { color: var(--ink); font-weight: 600; }
.delay-good::before { content: "\25CF"; }

.delay-mid { color: var(--ink-2); }
.delay-mid::before { content: "\25CF"; }

.delay-slow { color: var(--ink-3); }
.delay-slow::before { content: "\25D1"; }

.delay-fail { color: var(--ink-3); }
.delay-fail::before { content: "\2715"; }

.delay-idle { color: var(--ink-3); }
.delay-idle::before { content: "\25CB"; }
</style>
