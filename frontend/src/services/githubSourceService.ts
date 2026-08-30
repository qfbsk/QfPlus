import {
  GetGitHubSourceSettings,
  SaveGitHubSourceSettings,
} from '../../wailsjs/go/app/App';
import type { GitHubSourceSettings } from './appModels';
export type { GitHubSource, GitHubSourceSettings } from './appModels';

export const fetchGitHubSourceSettings = () => GetGitHubSourceSettings();

export const saveGitHubSourceSettings = (settings: GitHubSourceSettings) =>
  SaveGitHubSourceSettings(settings);
