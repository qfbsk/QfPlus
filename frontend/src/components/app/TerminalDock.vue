<script lang="ts" setup>
import { ref } from 'vue';
import { t } from '../../i18n';
import type { TerminalLogEntry } from '../../composables/useTaskTerminal';

defineProps<{
  expanded: boolean;
  logs: TerminalLogEntry[];
}>();

const emit = defineEmits(['clear', 'toggle']);
const terminalBody = ref<HTMLElement | null>(null);

const scrollToBottom = () => {
  if (!terminalBody.value) return;
  terminalBody.value.scrollTop = terminalBody.value.scrollHeight;
};

defineExpose({
  scrollToBottom,
});
</script>

<template>
  <section class="terminal-dock" :class="{ collapsed: !expanded }" :aria-label="t('terminal.aria')">
    <!-- Collapsed: a single thin status line, like a browser status bar -->
    <button
      v-if="!expanded"
      type="button"
      class="terminal-collapsed-bar"
      :aria-label="t('terminal.expand')"
      :title="t('terminal.expand')"
      @click="emit('toggle')"
    >
      <span class="terminal-collapsed-label">
        <span class="material-symbols-outlined">terminal</span>
        {{ t('terminal.title') }}
        <span v-if="logs.length" class="terminal-count">{{ logs.length }}</span>
      </span>
      <span class="terminal-chevron-wrap" aria-hidden="true">
        <svg class="terminal-chevron" viewBox="0 0 24 24">
          <path d="M7 10 L12 15 L17 10" />
        </svg>
      </span>
    </button>

    <!-- Expanded: the full terminal panel -->
    <template v-else>
      <div class="terminal-header">
        <div class="terminal-title">
          <span class="material-symbols-outlined">terminal</span>
          {{ t('terminal.title') }}
          <span class="terminal-count">{{ logs.length }}</span>
        </div>
        <div class="terminal-actions">
          <button class="terminal-icon-btn" :title="t('terminal.clear')" @click="emit('clear')">
            <span class="material-symbols-outlined">delete_sweep</span>
          </button>
          <button
            class="terminal-icon-btn terminal-toggle"
            :title="t('terminal.collapse')"
            :aria-label="t('terminal.collapse')"
            @click="emit('toggle')"
          >
            <svg class="terminal-chevron up" viewBox="0 0 24 24">
              <path d="M7 10 L12 15 L17 10" />
            </svg>
          </button>
        </div>
      </div>
      <div ref="terminalBody" class="terminal-body">
        <div v-if="logs.length === 0" class="terminal-empty">{{ t('terminal.empty') }}</div>
        <div v-for="entry in logs" :key="entry.id" class="terminal-line" :class="entry.level">
          <span class="terminal-time">{{ entry.time }}</span>
          <span class="terminal-prompt">$</span>
          <span class="terminal-text">{{ entry.text }}</span>
        </div>
      </div>
    </template>
  </section>
</template>
