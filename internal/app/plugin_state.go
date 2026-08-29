package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *App) getAddedPlugins() ([]string, error) {
	vfoxHome := a.getVfoxHome()
	if strings.TrimSpace(vfoxHome) == "" {
		return []string{}, fmt.Errorf("unable to resolve vfox home directory")
	}
	pluginDir := filepath.Join(vfoxHome, "plugin")

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var plugins []string
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}
	return plugins, nil
}

// loadAddedPluginSet returns a lowercase set of plugin names currently added.
func (a *App) loadAddedPluginSet(result *SdkEnvironmentImportResult) map[string]bool {
	addedPlugins := make(map[string]bool)
	plugins, err := a.getAddedPlugins()
	if err != nil {
		if result != nil {
			result.Warnings = append(result.Warnings, "Unable to list added vfox plugins: "+err.Error())
		}
		return addedPlugins
	}
	for _, plugin := range plugins {
		addedPlugins[strings.ToLower(strings.TrimSpace(plugin))] = true
	}
	return addedPlugins
}

// applyIsAddedStatus marks market plugins that already have vfox-managed SDKs.
func (a *App) applyIsAddedStatus(plugins []PluginInfo) []PluginInfo {
	installedSdks, _ := a.getInstalledSdks()
	addedMap := make(map[string]bool)
	for _, sdk := range installedSdks {
		if sdk.Source == "vfox" {
			addedMap[sdk.Name] = true
		}
	}
	for i := range plugins {
		plugins[i].IsAdded = addedMap[plugins[i].Name]
	}
	return plugins
}
