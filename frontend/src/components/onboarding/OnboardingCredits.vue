<script lang="ts" setup>
import { forkNoteGroups } from '../../app/forkNotes';
import { QFPLUS_REPO, VFOX_REPO, VFOXG_REPO } from '../../app/links';
import { t } from '../../i18n';

const emit = defineEmits<{ (event: 'open-link', url: string): void }>();
</script>

<template>
  <div class="ob-credits">
    <section class="ob-block ob-rise">
      <h2 class="ob-label">{{ t('about.author') }}</h2>
      <p class="ob-author">{{ t('about.author.name') }}</p>
    </section>

    <section class="ob-block ob-rise">
      <h2 class="ob-label">{{ t('onboarding.credits') }}</h2>
      <ul class="ob-lines">
        <li>{{ t('onboarding.credit.vfox') }}</li>
        <li>{{ t('onboarding.credit.gui') }}</li>
        <li>{{ t('about.qwen.body') }}</li>
      </ul>
      <div class="ob-links">
        <button class="ob-link" @click="emit('open-link', QFPLUS_REPO)">
          <span>{{ t('about.link.repo') }}</span>
          <span class="material-symbols-outlined">open_in_new</span>
        </button>
        <button class="ob-link" @click="emit('open-link', VFOX_REPO)">
          <span>{{ t('about.link.vfox_repo') }}</span>
          <span class="material-symbols-outlined">open_in_new</span>
        </button>
        <button class="ob-link" @click="emit('open-link', VFOXG_REPO)">
          <span>{{ t('about.link.vfoxg_repo') }}</span>
          <span class="material-symbols-outlined">open_in_new</span>
        </button>
      </div>
    </section>

    <section class="ob-block ob-rise">
      <h2 class="ob-label">{{ t('about.changes.title') }}</h2>
      <p class="ob-note">{{ t('about.changes.intro') }}</p>
      <div v-for="group in forkNoteGroups" :key="group.id" class="ob-group">
        <h3 class="ob-group-label">{{ t(`about.group.${group.id}`) }}</h3>
        <ul class="ob-lines">
          <li v-for="change in group.changes" :key="change">{{ t(`about.change.${change}`) }}</li>
        </ul>
      </div>
    </section>
  </div>
</template>
