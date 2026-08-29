<script lang="ts" setup>
import OnboardingSheet from './OnboardingSheet.vue';
import { useOnboarding } from '../../composables/useOnboarding';
import { openExternalUrl } from '../../services/runtimeService';

const emit = defineEmits<{ (event: 'dismissed'): void }>();

const { phase, isAnimating, requestEnter } = useOnboarding(() => emit('dismissed'));
</script>

<template>
  <div
    class="ob-root"
    :class="{
      'is-animating': isAnimating,
      'is-ready': phase !== 'loading',
      'is-leaving': phase === 'leaving',
    }"
  >
    <OnboardingSheet :phase="phase" @enter="requestEnter" @open-link="openExternalUrl" />
  </div>
</template>
