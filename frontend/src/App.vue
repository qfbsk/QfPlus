<script lang="ts" setup>
import { nextTick, provide, readonly, ref } from 'vue';
import { t } from './i18n';
import SdkManager from './components/sdk/SdkManager.vue';
import EnvironmentView from './components/environment/EnvironmentView.vue';
import PluginMarket from './components/plugin/PluginMarket.vue';
import Settings from './components/settings/Settings.vue';
import AboutView from './components/about/AboutView.vue';
import AppSidebar from './components/app/AppSidebar.vue';
import MigrationProgressOverlay from './components/app/MigrationProgressOverlay.vue';
import OnboardingOverlay from './components/onboarding/OnboardingOverlay.vue';
import TerminalDock from './components/app/TerminalDock.vue';
import { getNavTransition, type AppTab, type NavTransition } from './app/navigation';
import { shouldShowOnboarding } from './composables/useOnboarding';
import { useTerminalDock } from './composables/useTerminalDock';
import { useMigrationProgress } from './composables/useMigrationProgress';
import { useTaskTerminal } from './composables/useTaskTerminal';
import { useEnvironmentStatus } from './composables/useEnvironmentStatus';

type SdkSidebarAction = { id: number; type: 'display' };

const terminalDock = ref<InstanceType<typeof TerminalDock> | null>(null);
const scrollTerminalToBottom = () => {
  terminalDock.value?.scrollToBottom();
};

const { showTerminal } = useTerminalDock({
  onShowTerminal: () => nextTick(scrollTerminalToBottom),
});
provide('showTerminal', showTerminal);

const {
  showTaskToast,
  terminalTaskRunning,
  busyHintVisible,
  taskTitle,
  taskStatus,
  lastLogLine,
  taskProgress,
  taskDownloadSpeed,
  taskPhase,
  terminalLogs,
  isDeterminateDownloadProgress,
  showToastProgress,
  toastProgressStyle,
  handleStartTask,
  handleNotify,
  runTerminalTask,
  clearTerminalLogs,
} = useTaskTerminal({ showTerminal, scrollTerminalToBottom });
provide('terminalTaskRunning', readonly(terminalTaskRunning));
provide('runTerminalTask', runTerminalTask);

// Share the live task-toast state with the sidebar dock so the toast is rendered
// inside the sidebar instead of as a floating window.
provide('taskToast', {
  showTaskToast,
  taskStatus,
  taskPhase,
  taskTitle,
  lastLogLine,
  taskProgress,
  taskDownloadSpeed,
  isDeterminateDownloadProgress,
  showToastProgress,
  toastProgressStyle,
  busyHintVisible,
});

const envStatus = useEnvironmentStatus();

const handleActionDiagnostic = async () => {
  // Expand the terminal so the diagnostic output is visible immediately.
  showTerminal.value = true;
  try {
    await envStatus.openDiagnostic();
  } catch (err) {
    handleNotify({
      type: 'error',
      title: t('common.error'),
      message: err instanceof Error ? err.message : String(err),
    });
  }
};

const {
  migrationVisible,
  migrationProgress,
  formatRemainingTime,
} = useMigrationProgress();

const isOnboardingVisible = ref(shouldShowOnboarding());

const currentTab = ref<AppTab>('sdk');
const navTransition = ref<NavTransition>('slide-up');
const sdkSidebarAction = ref<SdkSidebarAction | null>(null);
let sdkSidebarActionId = 0;

const switchTab = (tab: AppTab) => {
  if (tab === currentTab.value) return;
  navTransition.value = getNavTransition(currentTab.value, tab);
  currentTab.value = tab;
};
provide('switchTab', switchTab);

const triggerSdkSidebarAction = (type: SdkSidebarAction['type']) => {
  if (currentTab.value !== 'sdk') {
    switchTab('sdk');
  }
  sdkSidebarAction.value = { id: ++sdkSidebarActionId, type };
};

const clearSdkSidebarAction = (id: number) => {
  if (sdkSidebarAction.value?.id === id) {
    sdkSidebarAction.value = null;
  }
};
</script>

<template>
  <div id="layout">
    <AppSidebar
      :current-tab="currentTab"
      @show-sdk-display="triggerSdkSidebarAction('display')"
      @switch-tab="switchTab"
      @action-diagnostic="handleActionDiagnostic"
    />
    <div class="main-shell">
      <div class="main-content">
        <Transition :name="navTransition" mode="out-in">
          <SdkManager
            v-if="currentTab === 'sdk'"
            key="sdk"
            :sidebar-action="sdkSidebarAction"
            @sidebar-action-done="clearSdkSidebarAction"
            @start-task="handleStartTask"
            @notify="handleNotify"
            @open-plugin-market="switchTab('plugin')"
            @open-environment="switchTab('environment')"
          />
          <EnvironmentView v-else-if="currentTab === 'environment'" key="environment" @notify="handleNotify" @start-task="handleStartTask" />
          <PluginMarket v-else-if="currentTab === 'plugin'" key="plugin" @start-task="handleStartTask" @notify="handleNotify" />
          <Settings v-else-if="currentTab === 'settings'" key="settings" @notify="handleNotify" />
          <AboutView v-else-if="currentTab === 'about'" key="about" />
        </Transition>
      </div>

      <TerminalDock
        ref="terminalDock"
        :expanded="showTerminal"
        :logs="terminalLogs"
        @clear="clearTerminalLogs"
        @toggle="showTerminal = !showTerminal"
      />
    </div>

    <MigrationProgressOverlay
      :visible="migrationVisible"
      :progress="migrationProgress"
      :format-remaining-time="formatRemainingTime"
    />

    <OnboardingOverlay v-if="isOnboardingVisible" @dismissed="isOnboardingVisible = false" />
  </div>
</template>
