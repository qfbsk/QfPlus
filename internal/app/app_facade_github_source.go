package app

import (
	"fmt"
	"strings"
)

// GetGitHubSourceSettings returns the current GitHub source configuration,
// merged with built-in presets.
func (a *App) GetGitHubSourceSettings() (GitHubSourceSettings, error) {
	config, err := a.readAppConfig()
	if err != nil {
		return GitHubSourceSettings{}, err
	}

	custom := make([]GitHubSource, len(config.GitHubSource.CustomSources))
	copy(custom, config.GitHubSource.CustomSources)

	sources := make([]GitHubSource, len(presetGitHubSources)+len(custom))
	copy(sources, presetGitHubSources)
	copy(sources[len(presetGitHubSources):], custom)

	selectedID := config.GitHubSource.SelectedID
	if selectedID == "" {
		selectedID = presetGitHubSources[0].ID
	}

	return GitHubSourceSettings{
		Enabled:       config.GitHubSource.Enabled,
		SelectedID:    selectedID,
		Sources:       sources,
		PresetSources: append([]GitHubSource(nil), presetGitHubSources...),
		CustomSources: custom,
	}, nil
}

// SaveGitHubSourceSettings saves the GitHub source selection and any custom
// sources the user wants to keep. Custom sources without a name or URL are
// dropped.
func (a *App) SaveGitHubSourceSettings(settings GitHubSourceSettings) (GitHubSourceSettings, error) {
	config, err := a.readAppConfig()
	if err != nil {
		return GitHubSourceSettings{}, err
	}

	var custom []GitHubSource
	seen := map[string]bool{}
	for _, s := range settings.CustomSources {
		name := strings.TrimSpace(s.Name)
		url := strings.TrimSpace(s.URL)
		if name == "" || url == "" {
			continue
		}
		id := strings.TrimSpace(s.ID)
		if id == "" || seen[id] {
			id = generateSourceID()
		}
		seen[id] = true
		custom = append(custom, GitHubSource{
			ID:       id,
			Name:     name,
			URL:      url,
			IsPreset: false,
		})
	}

	selectedID := strings.TrimSpace(settings.SelectedID)
	if selectedID == "" {
		selectedID = presetGitHubSources[0].ID
	}

	config.GitHubSource = GitHubSourceConfig{
		Enabled:       settings.Enabled,
		SelectedID:    selectedID,
		CustomSources: custom,
	}

	if err := a.saveAppConfig(config); err != nil {
		return GitHubSourceSettings{}, fmt.Errorf("failed to save GitHub source settings: %w", err)
	}

	// Push the mirror choice into vfox's own plugin registry config so that
	// plugin discovery/manifest downloads can benefit from it. Errors here are
	// non-fatal: vfox will simply fall back to its default registry.
	_ = a.syncVfoxGitHubSourceConfig()

	return a.GetGitHubSourceSettings()
}
