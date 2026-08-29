import {
  GetDownloadPathInfo,
  PlanEnvironmentMigration,
  ResetDownloadPath,
  ResetDownloadPathWithMigration,
  SelectDownloadPath,
  SetDownloadPath,
  SetDownloadPathWithMigration,
} from '../../wailsjs/go/app/App';
export type { DownloadPathInfo, MigrationPlan } from './appModels';

export const fetchDownloadPathInfo = () => GetDownloadPathInfo();

export const saveDownloadPath = (targetPath: string) => SetDownloadPath(targetPath);

export const saveDownloadPathWithMigration = (targetPath: string) =>
  SetDownloadPathWithMigration(targetPath, true);

export const resetDownloadPath = () => ResetDownloadPath();

export const resetDownloadPathWithMigration = () => ResetDownloadPathWithMigration(true);

export const selectDownloadPath = () => SelectDownloadPath();

export const planEnvironmentMigration = (targetPath: string) =>
  PlanEnvironmentMigration(targetPath);
