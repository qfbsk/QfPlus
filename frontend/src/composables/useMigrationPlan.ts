import { t } from '../i18n';
import type { MigrationItem, MigrationPlan } from '../services/appModels';

// MigrationItemKind is a plain string field in the generated model (wails does
// not emit TS types for Go const string blocks), so we alias it locally.
export type MigrationItemKind = string;

// Human-readable byte formatting: bytes → KB / MB / GB.
export const formatBytes = (bytes: number): string => {
  if (!bytes || bytes < 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let value = bytes;
  let unit = 0;
  while (value >= 1024 && unit < units.length - 1) {
    value /= 1024;
    unit += 1;
  }
  const rounded = unit === 0 ? value : Math.round(value * 10) / 10;
  return `${rounded} ${units[unit]}`;
};

// Localized label for a migration item kind.
export const migrationKindLabel = (kind: MigrationItemKind | string): string => {
  const key = `settings.download.path.plan_kind.${(kind || 'other').toString()}`;
  const label = t(key);
  return label === key ? (kind || 'other').toString() : label;
};

// Splits a plan into the shell-rendered groups used by the preview modal.
export const splitPlanItems = (plan: MigrationPlan | null | undefined) => {
  const movable = (plan?.movableItems || []).filter((item) => item.willMove);
  const excluded = (plan?.excludedItems || []).filter((item) => !item.willMove);
  return { movable, excluded };
};

export type { MigrationItem };
