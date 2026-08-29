<script lang="ts" setup>
import { t } from '../../i18n';
import { VFOX_LICENSE, VFOX_REPO, VFOXG_REPO } from '../../app/links';

const emit = defineEmits<{ (event: 'open-link', url: string): void }>();
</script>

<template>
  <section class="about-section">
    <h3 class="section-heading">{{ t('about.engine.title') }}</h3>
    <p class="about-body">{{ t('about.engine.body') }}</p>
    <p class="about-body">{{ t('about.engine.ownership') }}</p>
    <p class="about-note">{{ t('about.engine.license_note') }}</p>
    <div class="about-links">
      <button class="about-link" @click="emit('open-link', VFOX_REPO)">
        <span>{{ t('about.link.vfox_repo') }}</span>
        <span class="material-symbols-outlined">open_in_new</span>
      </button>
      <button class="about-link" @click="emit('open-link', VFOX_LICENSE)">
        <span>{{ t('about.link.apache') }}</span>
        <span class="material-symbols-outlined">open_in_new</span>
      </button>
    </div>
  </section>

  <section class="about-section">
    <h3 class="section-heading">{{ t('about.gui.title') }}</h3>
    <p class="about-body">{{ t('about.gui.body') }}</p>
    <div class="about-links">
      <button class="about-link" @click="emit('open-link', VFOXG_REPO)">
        <span>{{ t('about.link.vfoxg_repo') }}</span>
        <span class="material-symbols-outlined">open_in_new</span>
      </button>
    </div>
  </section>

  <section class="about-section">
    <h3 class="section-heading">{{ t('about.qwen.title') }}</h3>
    <p class="about-body">{{ t('about.qwen.body') }}</p>
  </section>
</template>
