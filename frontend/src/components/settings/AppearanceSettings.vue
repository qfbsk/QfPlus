<script lang="ts" setup>
import { currentLang, t } from '../../i18n';
import type { SupportedLang } from '../../i18n/keys';

defineProps<{
  showTerminal: boolean;
}>();

const emit = defineEmits(['update:showTerminal']);

const setLang = (lang: SupportedLang) => {
  currentLang.value = lang;
};
</script>

<template>
  <section class="set-group">
    <h2>{{ t('settings.appearance') }}</h2>

    <div class="set-row">
      <div class="left">
        <div class="label">{{ t('settings.language') }}</div>
        <div class="desc">{{ t('settings.language.desc') }}</div>
      </div>
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
    </div>

    <div class="set-row">
      <div class="left">
        <div class="label">{{ t('settings.terminal') }}</div>
        <div class="desc">{{ t('settings.terminal.desc') }}</div>
      </div>
      <label class="switch">
        <input
          type="checkbox"
          :checked="showTerminal"
          :aria-label="t('settings.terminal.show')"
          @change="emit('update:showTerminal', ($event.target as HTMLInputElement).checked)"
        >
        <span class="slider"></span>
      </label>
    </div>
  </section>
</template>
