<script lang="ts" setup>
import { t } from '../../i18n';
import { forkNoteGroups } from '../../app/forkNotes';
</script>

<template>
  <section class="about-section">
    <h3 class="section-heading">{{ t('about.changes.title') }}</h3>
    <p class="about-body">{{ t('about.changes.intro') }}</p>
    <ul class="about-groups">
      <li v-for="group in forkNoteGroups" :key="group.id">
        <h4 class="about-group-label">{{ t(`about.group.${group.id}`) }}</h4>
        <ul class="about-changes">
          <li v-for="change in group.changes" :key="change" class="about-change">
            {{ t(`about.change.${change}`) }}
          </li>
        </ul>
      </li>
    </ul>
  </section>
</template>
