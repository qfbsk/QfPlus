<script lang="ts" setup>
import { inject, ref, type CSSProperties, type Ref } from 'vue';
import { t } from '../../i18n';
import type { AppTab } from '../../app/navigation';
import type { ToastStatus, TaskPhase } from '../../composables/useTaskTerminal';
import { useDownloadQueue } from '../../composables/useDownloadQueue';
import DownloadQueuePanel from './DownloadQueuePanel.vue';
import TaskToast from './TaskToast.vue';

defineProps<{
  currentTab: AppTab;
}>();

const emit = defineEmits<{
  (event: 'show-sdk-display'): void;
  (event: 'switch-tab', tab: AppTab): void;
  (event: 'action-diagnostic'): void;
}>();

interface TaskToastState {
  showTaskToast: Ref<boolean>;
  taskStatus: Ref<ToastStatus>;
  taskPhase: Ref<TaskPhase>;
  taskTitle: Ref<string>;
  lastLogLine: Ref<string>;
  taskProgress: Ref<number>;
  taskDownloadSpeed: Ref<string>;
  isDeterminateDownloadProgress: Ref<boolean>;
  showToastProgress: Ref<boolean>;
  toastProgressStyle: Ref<CSSProperties>;
  busyHintVisible: Ref<boolean>;
}

const toastState = inject<TaskToastState | null>('taskToast', null);
const showTaskToast = toastState?.showTaskToast ?? ref(false);
const taskStatus = toastState?.taskStatus ?? ref<ToastStatus>('running');
const taskPhase = toastState?.taskPhase ?? ref<TaskPhase>('default');
const taskTitle = toastState?.taskTitle ?? ref('');
const lastLogLine = toastState?.lastLogLine ?? ref('');
const taskProgress = toastState?.taskProgress ?? ref(0);
const taskDownloadSpeed = toastState?.taskDownloadSpeed ?? ref('');
const isDeterminate = toastState?.isDeterminateDownloadProgress ?? ref(false);
const showToastProgress = toastState?.showToastProgress ?? ref(false);
const toastProgressStyle = toastState?.toastProgressStyle ?? ref<CSSProperties>({});
const busyHintVisible = toastState?.busyHintVisible ?? ref(false);

const { visible: queueVisible } = useDownloadQueue();
</script>

<template>
  <div class="sidebar">
    <div class="logo">
      <h2>QfPlus</h2>
    </div>

    <nav>
      <button class="nav-btn" :class="{active: currentTab === 'sdk'}" @click="emit('show-sdk-display')">
        <svg class="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
          <path d="M12 3.4 20.6 8 12 12.6 3.4 8z" />
          <path d="M4.6 12.4 12 16.5 19.4 12.4" />
          <path d="M4.6 16.6 12 20.7 19.4 16.6" />
        </svg>
        <span class="nav-label">{{ t('nav.display') }}</span>
      </button>
      <button class="nav-btn" :class="{active: currentTab === 'environment'}" @click="emit('switch-tab', 'environment')">
        <svg class="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
          <path d="M12 2a10 10 0 0 0-7.4 16.7l.7.8" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 22a10 10 0 0 0 7.4-16.7l-.7-.8" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"/>
          <circle cx="12" cy="12" r="3" fill="none" stroke="currentColor" stroke-width="1.7"/>
        </svg>
        <span class="nav-label">{{ t('nav.environment') }}</span>
      </button>
      <button class="nav-btn" :class="{active: currentTab === 'plugin'}" @click="emit('switch-tab', 'plugin')">
        <svg class="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
          <path d="M5.6 8.4h12.8l-1.1 11.2H6.7z" />
          <path d="M9 8.4V6.9a3 3 0 0 1 6 0v1.5" />
        </svg>
        <span class="nav-label">{{ t('nav.market') }}</span>
      </button>
    </nav>

    <nav class="sidebar-actions">
      <button class="nav-btn" :title="t('nav.diagnostic')" @click="emit('action-diagnostic')">
        <svg class="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
          <path d="M5 7l4 4-4 4" />
          <path d="M11 15h7" />
        </svg>
        <span class="nav-label">{{ t('nav.diagnostic') }}</span>
      </button>
    </nav>

    <div class="sidebar-foot">
      <div v-if="queueVisible || showTaskToast" class="sidebar-dock">
        <DownloadQueuePanel />
        <TaskToast
          :visible="showTaskToast"
          :status="taskStatus"
          :phase="taskPhase"
          :title="taskTitle"
          :message="lastLogLine"
          :show-progress="showToastProgress"
          :is-determinate-download-progress="isDeterminate"
          :progress="taskProgress"
          :download-speed="taskDownloadSpeed"
          :progress-style="toastProgressStyle"
          :busy-hint-visible="busyHintVisible"
        />
      </div>

      <nav class="bottom">
        <button class="nav-btn" :class="{active: currentTab === 'settings'}" @click="emit('switch-tab', 'settings')">
          <svg class="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M4 8h5.2" />
            <circle cx="12.4" cy="8" r="2.6" />
            <path d="M15 8h5" />
            <path d="M4 16h2.6" />
            <circle cx="9.8" cy="16" r="2.6" />
            <path d="M12.4 16h6.6" />
          </svg>
          <span class="nav-label">{{ t('nav.settings') }}</span>
        </button>
        <button class="nav-btn" :class="{active: currentTab === 'about'}" @click="emit('switch-tab', 'about')">
          <svg class="nav-icon" viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="12" cy="12" r="8.6" />
            <path d="M12 11.2v5" />
            <path d="M12 8.1h.01" />
          </svg>
          <span class="nav-label">{{ t('nav.about') }}</span>
        </button>
      </nav>
    </div>
  </div>
</template>
