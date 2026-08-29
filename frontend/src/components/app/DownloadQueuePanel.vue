<script lang="ts" setup>
import { t } from '../../i18n';
import { useDownloadQueue, type DownloadTask } from '../../composables/useDownloadQueue';

const { tasks, visible, activeCount, dismiss, clearFinished } = useDownloadQueue();

const kindLabel = (task: DownloadTask) => t('download_queue.plugin_title');

const statusLabel = (status: DownloadTask['status']) => {
  switch (status) {
    case 'queued':
      return t('download_queue.queued');
    case 'running':
      return t('download_queue.running');
    case 'done':
      return t('download_queue.done');
    case 'error':
      return t('download_queue.error');
    default:
      return status;
  }
};

const phaseLabel = (phase: string) => {
  if (!phase) return '';
  const key = `environment.queue.phase.${phase}`;
  const label = t(key);
  return label !== key ? label : phase;
};
</script>

<template>
  <transition name="dq-fade">
    <div v-if="visible" class="download-queue" role="region" aria-label="download queue">
      <header class="dq-header">
        <div class="dq-header-title">
          <span class="dq-dot" :class="{ active: activeCount > 0 }"></span>
          <span>{{ t('download_queue.title') }}</span>
          <span class="dq-count">{{ activeCount }} / {{ tasks.length }}</span>
        </div>
        <button
          v-if="tasks.some((t) => t.status === 'done' || t.status === 'error')"
          class="dq-clear"
          type="button"
          @click="clearFinished"
        >
          {{ t('download_queue.clear') }}
        </button>
      </header>

      <ul class="dq-list">
        <li
          v-for="task in tasks"
          :key="task.id"
          class="dq-task"
          :class="[task.status]"
        >
          <div class="dq-task-head">
            <span class="dq-kind">{{ kindLabel(task) }}</span>
            <span class="dq-title" :title="task.title">{{ task.title }}</span>
            <span class="dq-status">{{ statusLabel(task.status) }}</span>
          </div>

          <div class="dq-sub">
            <span v-if="task.subtitle" class="dq-subtitle">{{ task.subtitle }}</span>
            <span v-else-if="phaseLabel(task.phase)" class="dq-phase">{{ phaseLabel(task.phase) }}</span>
          </div>

          <div class="dq-bar">
            <div
              class="dq-fill"
              :class="{ indeterminate: task.unitPercent === null }"
              :style="task.unitPercent !== null ? { width: task.unitPercent + '%' } : undefined"
            ></div>
          </div>

          <div class="dq-meta">
            <span class="dq-files">
              {{ t('download_queue.files', { done: task.unitsDone, total: task.unitsTotal }) }}
            </span>
            <span v-if="task.speed" class="dq-speed">{{ t('download_queue.speed', { speed: task.speed }) }}</span>
            <span v-if="task.status === 'error'" class="dq-err">{{ task.error }}</span>
            <button
              v-if="task.status === 'done' || task.status === 'error'"
              class="dq-dismiss"
              type="button"
              :title="t('download_queue.dismiss')"
              @click="dismiss(task.id)"
            >
              ×
            </button>
          </div>
        </li>
      </ul>
    </div>
  </transition>
</template>

<style scoped>
.download-queue {
  width: 100%;
  max-height: 46vh;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  border: 1px solid #e5e5e5;
  border-radius: var(--radius-card);
  font-size: 13px;
  color: #1a1a1a;
  overflow: hidden;
}

.dq-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
}

.dq-header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.dq-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #999999;
}

.dq-dot.active {
  background: #1a1a1a;
  animation: dq-pulse 1.4s ease-in-out infinite;
}

.dq-count {
  color: #666666;
  font-weight: 500;
}

.dq-clear {
  border: none;
  background: transparent;
  color: #1a1a1a;
  cursor: pointer;
  font-size: 12px;
  padding: 2px 4px;
}

.dq-list {
  list-style: none;
  margin: 0;
  padding: 6px;
  overflow-y: auto;
}

.dq-task {
  padding: 8px 10px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  margin-bottom: 6px;
  background: #ffffff;
}

.dq-task:last-child {
  margin-bottom: 0;
}

.dq-task.error {
  border-color: #e0e0e0;
  background: #fafafa;
}

.dq-task.done {
  border-color: #e0e0e0;
  background: #fafafa;
}

.dq-task-head {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dq-kind {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 999px;
  background: #f0f0f0;
  color: #333333;
  white-space: nowrap;
}

.dq-title {
  flex: 1;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dq-status {
  font-size: 11px;
  color: #666666;
  white-space: nowrap;
}

.dq-task.error .dq-status {
  color: #1a1a1a;
}

.dq-task.done .dq-status {
  color: #1a1a1a;
}

.dq-sub {
  min-height: 16px;
  margin: 4px 0;
  color: #555555;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dq-bar {
  height: 6px;
  border-radius: 999px;
  background: #e5e5e5;
  overflow: hidden;
}

.dq-fill {
  height: 100%;
  background: #1a1a1a;
  border-radius: 999px;
  transition: width 0.25s ease;
}

.dq-task.done .dq-fill,
.dq-task.error .dq-fill {
  background: #1a1a1a;
}

.dq-fill.indeterminate {
  width: 40%;
  animation: dq-indeterminate 1.2s ease-in-out infinite;
}

.dq-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 5px;
  font-size: 11px;
  color: #666666;
}

.dq-speed {
  color: #333333;
  font-variant-numeric: tabular-nums;
}

.dq-err {
  color: #333333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.dq-dismiss {
  margin-left: auto;
  border: none;
  background: transparent;
  color: #999999;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 0 2px;
}

.dq-dismiss:hover {
  color: #333333;
}

@keyframes dq-indeterminate {
  0% {
    margin-left: -40%;
  }
  100% {
    margin-left: 100%;
  }
}

@keyframes dq-pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

.dq-fade-enter-active,
.dq-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.dq-fade-enter-from,
.dq-fade-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
