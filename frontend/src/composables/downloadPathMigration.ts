import type { DownloadPathInfo, MigrationPlan } from '../services/settingsService';

export type PendingMigrationAction = {
  type: 'save' | 'reset';
  targetPath: string;
};

export type PendingMigrationPlan = {
  mode: 'save' | 'reset';
  targetPath: string;
  plan: MigrationPlan;
};

export const shouldAskDownloadPathMigration = (
  downloadPathInfo: DownloadPathInfo | null,
  targetPath: string
) => {
  const currentPath = downloadPathInfo?.path?.trim();
  const normalizedTargetPath = targetPath.trim();

  return Boolean(
    downloadPathInfo?.hasMigratableData &&
      currentPath &&
      normalizedTargetPath &&
      currentPath !== normalizedTargetPath
  );
};
