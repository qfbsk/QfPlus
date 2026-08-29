<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import AboutAcknowledgments from './AboutAcknowledgments.vue';
import AboutChanges from './AboutChanges.vue';
import { QFPLUS_REPO } from '../../app/links';
import { t } from '../../i18n';
import type { CoreInfo } from '../../services/appModels';
import { fetchCoreInfo } from '../../services/coreService';
import { openExternalUrl } from '../../services/runtimeService';

const core = ref<CoreInfo | null>(null);

onMounted(async () => {
  try {
    core.value = await fetchCoreInfo();
  } catch {
    core.value = null;
  }
});
</script>

<template>
  <div class="view-container about-view">
    <h2>{{ t('about.title') }}</h2>
    <p class="about-tagline">{{ t('about.tagline') }}</p>

    <div class="about-facts">
      <div class="about-fact">
        <span class="about-fact-label">{{ t('about.app') }}</span>
        <span class="about-fact-value">QfPlus</span>
      </div>
      <div class="about-fact">
        <span class="about-fact-label">{{ t('about.author') }}</span>
        <span class="about-fact-value">{{ t('about.author.name') }}</span>
      </div>
      <div class="about-fact">
        <span class="about-fact-label">{{ t('about.engine') }}</span>
        <span class="about-fact-value">{{ core?.currentVersion || '—' }}</span>
      </div>
      <div class="about-fact">
        <span class="about-fact-label">{{ t('about.engine.bundled') }}</span>
        <span class="about-fact-value">{{ core?.bundledVersion || '—' }}</span>
      </div>
      <div class="about-fact">
        <span class="about-fact-label">{{ t('about.platform') }}</span>
        <span class="about-fact-value">{{ core?.osArch || '—' }}</span>
      </div>
      <div class="about-fact">
        <span class="about-fact-label">{{ t('about.license') }}</span>
        <span class="about-fact-value">Apache License 2.0</span>
      </div>
    </div>

    <AboutAcknowledgments @open-link="openExternalUrl" />

    <AboutChanges />

    <section class="about-section">
      <h3 class="section-heading">{{ t('about.author.title') }}</h3>
      <p class="about-body">{{ t('about.author.body') }}</p>
      <p class="about-author-name">{{ t('about.author.name') }}</p>
      <div class="about-links">
        <button class="about-link" @click="openExternalUrl(QFPLUS_REPO)">
          <span>{{ t('about.link.repo') }}</span>
          <span class="material-symbols-outlined">open_in_new</span>
        </button>
      </div>
    </section>
  </div>
</template>
