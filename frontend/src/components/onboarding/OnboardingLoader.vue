<script lang="ts" setup>
import { t } from '../../i18n';
import type { OnboardingPhase } from '../../composables/useOnboarding';

defineProps<{
  phase: OnboardingPhase;
}>();
</script>

<template>
  <header class="ob-head">
    <span class="ob-mark" aria-hidden="true">
      <svg viewBox="0 0 24 24">
        <path d="M12 3.4 20.6 8 12 12.6 3.4 8z" />
        <path d="M4.6 12.4 12 16.5 19.4 12.4" />
        <path d="M4.6 16.6 12 20.7 19.4 16.6" />
      </svg>
    </span>
    <span class="ob-rule" aria-hidden="true"></span>
    <div class="ob-loading" :class="{ 'is-done': phase !== 'loading' }">
      <span v-if="phase === 'loading'" class="spinner ob-spinner"></span>
      <span class="ob-loading-text">{{ t('onboarding.loading') }}</span>
    </div>
  </header>
</template>
