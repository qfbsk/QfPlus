<script lang="ts" setup>
import { computed } from 'vue';
import { t } from '../../i18n';
import type { PendingMigrationPlan } from '../../composables/downloadPathMigration';
import { formatBytes, migrationKindLabel, splitPlanItems } from '../../composables/useMigrationPlan';

const props = defineProps<{
  pendingPlan: PendingMigrationPlan | null;
  isDownloadPathBusy: boolean;
}>();

const emit = defineEmits(['confirm', 'cancel']);

const movable = computed(() => splitPlanItems(props.pendingPlan?.plan).movable);
const excluded = computed(() => splitPlanItems(props.pendingPlan?.plan).excluded);
const totalSize = computed(() =>
  formatBytes(props.pendingPlan?.plan?.totalSizeBytes ?? 0)
);
const totalCount = computed(() => props.pendingPlan?.plan?.totalCount ?? 0);
</script>

<template>
  <Teleport to="body">
    <Transition name="modal" appear>
      <div
        v-if="pendingPlan"
        class="modal-overlay migration-plan-overlay"
        role="dialog"
        aria-modal="true"
        @click="emit('cancel')"
      >
        <div class="modal-content migration-plan-modal" @click.stop>
          <div class="migration-plan-heading">
            <div class="migration-plan-icon">
              <span class="material-symbols-outlined">sync_alt</span>
            </div>
            <div>
              <h2 class="modal-title">{{ t('settings.download.path.plan_title') }}</h2>
              <p class="modal-message">{{ t('settings.download.path.plan_desc') }}</p>
            </div>
          </div>

          <div class="migration-plan-paths">
            <div class="migration-plan-path">
              <span class="migration-plan-path-label">{{ t('settings.download.path.plan_source') }}</span>
              <code>{{ pendingPlan.plan.sourcePath }}</code>
            </div>
            <div class="migration-plan-path">
              <span class="migration-plan-path-label">{{ t('settings.download.path.plan_target') }}</span>
              <code>{{ pendingPlan.plan.targetPath }}</code>
            </div>
          </div>

          <div v-if="pendingPlan.plan.blockingReason" class="migration-plan-blocking">
            {{ pendingPlan.plan.blockingReason }}
          </div>

          <div v-else class="migration-plan-body">
            <p class="migration-plan-summary">
              {{ t('settings.download.path.plan_summary', { count: totalCount, size: totalSize }) }}
            </p>

            <h3 class="migration-plan-group-title">
              {{ t('settings.download.path.plan_will_move') }}
            </h3>
            <ul class="migration-plan-list">
              <li v-for="item in movable" :key="item.name" class="migration-plan-row">
                <span class="migration-plan-name">{{ item.name }}</span>
                <span class="migration-plan-kind">{{ migrationKindLabel(item.kind) }}</span>
                <span class="migration-plan-meta">
                  {{ t('settings.download.path.plan_count', { count: item.count }) }}
                  · {{ formatBytes(item.sizeBytes) }}
                </span>
              </li>
              <li v-if="movable.length === 0" class="migration-plan-empty">
                {{ t('settings.download.path.plan_nothing') }}
              </li>
            </ul>

            <h3 class="migration-plan-group-title">
              {{ t('settings.download.path.plan_will_not_move') }}
            </h3>
            <ul class="migration-plan-list">
              <li v-for="item in excluded" :key="item.name" class="migration-plan-row excluded">
                <span class="migration-plan-name">{{ item.name }}</span>
                <span class="migration-plan-kind">{{ migrationKindLabel(item.kind) }}</span>
                <span class="migration-plan-meta">{{ item.reason }}</span>
              </li>
              <li v-if="excluded.length === 0" class="migration-plan-empty">
                {{ t('settings.download.path.plan_no_excluded') }}
              </li>
            </ul>

            <p class="migration-plan-note">{{ t('settings.download.path.plan_note') }}</p>
          </div>

          <div class="modal-actions">
            <button class="btn tonal" :disabled="isDownloadPathBusy" @click="emit('cancel')">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn primary"
              :disabled="isDownloadPathBusy || !!pendingPlan?.plan.blockingReason"
              @click="emit('confirm')"
            >
              <div v-if="isDownloadPathBusy" class="spinner small-spinner"></div>
              <template v-else>{{ t('settings.download.path.plan_confirm') }}</template>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
