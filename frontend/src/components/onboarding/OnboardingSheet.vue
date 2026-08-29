<script lang="ts" setup>
import OnboardingCredits from './OnboardingCredits.vue';
import OnboardingIntro from './OnboardingIntro.vue';
import OnboardingLoader from './OnboardingLoader.vue';
import { currentLang, t } from '../../i18n';
import type { SupportedLang } from '../../i18n/keys';
import type { OnboardingPhase } from '../../composables/useOnboarding';

defineProps<{
  phase: OnboardingPhase;
}>();

const emit = defineEmits<{
  (event: 'enter'): void;
  (event: 'open-link', url: string): void;
}>();

const setLang = (lang: SupportedLang) => {
  currentLang.value = lang;
};
</script>

<template>
  <div class="ob-sheet" role="dialog" aria-modal="true" :aria-label="t('onboarding.title')">
    <OnboardingLoader :phase="phase" />

    <div v-if="phase !== 'loading'" class="ob-body">
      <OnboardingIntro class="ob-rise" />
      <OnboardingCredits @open-link="emit('open-link', $event)" />
    </div>

    <footer v-if="phase !== 'loading'" class="ob-footer ob-rise">
      <div class="seg" role="radiogroup" :aria-label="t('settings.language')">
        <button
          type="button"
          :class="{ on: currentLang === 'en' }"
          role="radio"
          :aria-checked="currentLang === 'en'"
          @click="setLang('en')"
        >{{ t('settings.language.en') }}</button>
        <button
          type="button"
          :class="{ on: currentLang === 'zh' }"
          role="radio"
          :aria-checked="currentLang === 'zh'"
          @click="setLang('zh')"
        >{{ t('settings.language.zh') }}</button>
      </div>
      <button type="button" class="btn filled ob-enter" @click="emit('enter')">
        {{ t('onboarding.enter') }}
        <span class="material-symbols-outlined">arrow_forward</span>
      </button>
    </footer>
  </div>
</template>
