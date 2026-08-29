package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetCoreInfo returns the local engine state plus any cached update verdict.
// It never touches the network.
func (a *App) GetCoreInfo() (CoreInfo, error) {
	config, _ := a.readAppConfig()
	info := CoreInfo{
		OsArch:     coreOSName() + "/" + coreArchName(),
		AutoUpdate: config.Core.AutoUpdate,
		LastCheck:  config.Core.LastCheck,
	}
	info.ExecutablePath = findCoreFile(coreExecutableName())

	current, err := a.getCurrentCoreVersion()
	if err != nil {
		info.Error = err.Error()
		info.LatestVersion = config.Core.LatestKnown
		return info, nil
	}
	info.CurrentVersion = current
	info.UsesLocalCore = markerCoreDir() != ""
	if info.UsesLocalCore {
		info.BundledVersion = a.coreVersionForDir(a.bundledCoreDir())
	} else {
		info.BundledVersion = current
	}

	if cached := a.readCoreReleasesCache(); len(cached) > 0 {
		info.LatestVersion = cached[0].Version
		info.ReleaseNotes = cached[0].Notes
		info.ReleaseURL = cached[0].URL
	}
	if config.Core.LatestKnown != "" {
		info.LatestVersion = config.Core.LatestKnown
	}
	info.UpdateAvailable = info.LatestVersion != "" &&
		compareCoreVersions(info.LatestVersion, current) > 0

	if config.Core.CurrentKnown != current {
		config.Core.CurrentKnown = current
		_ = a.saveAppConfig(config)
	}
	return info, nil
}

// CheckCoreUpdate fetches the upstream release feed and refreshes the cache.
func (a *App) CheckCoreUpdate() (CoreInfo, error) {
	releases, err := a.fetchCoreReleases()
	if err != nil {
		info, _ := a.GetCoreInfo()
		info.Error = err.Error()
		return info, err
	}
	now := time.Now().Format(time.RFC3339)
	config, _ := a.readAppConfig()
	config.Core.LastCheck = now
	config.Core.LatestKnown = releases[0].Version
	_ = a.saveAppConfig(config)
	_ = a.writeCoreReleasesCache(releases)

	info, _ := a.GetCoreInfo()
	return info, nil
}

// ListCoreVersions returns upstream releases with a flag marking which one is
// active and which ones are already downloaded locally.
func (a *App) ListCoreVersions() ([]CoreRelease, error) {
	releases, err := a.fetchCoreReleases()
	if err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	config, _ := a.readAppConfig()
	config.Core.LastCheck = now
	config.Core.LatestKnown = releases[0].Version
	_ = a.saveAppConfig(config)
	_ = a.writeCoreReleasesCache(releases)

	current, _ := a.getCurrentCoreVersion()
	bundledVersion := a.coreVersionForDir(a.bundledCoreDir())
	downloaded := map[string]bool{}
	for _, version := range localCoreVersions() {
		downloaded[version] = true
	}
	for i := range releases {
		releases[i].IsCurrent = releases[i].Version == current
		releases[i].Downloaded = downloaded[releases[i].Version] || releases[i].Version == bundledVersion
	}
	return releases, nil
}

// SwitchCoreVersion activates an engine version, downloading it on demand.
// Passing "bundled" restores the version shipped inside the installer.
func (a *App) SwitchCoreVersion(version string) (CoreInfo, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(version), "v")
	if trimmed == "" {
		return CoreInfo{}, fmt.Errorf("版本号不能为空")
	}
	if strings.EqualFold(trimmed, "bundled") {
		if err := a.activateBundledCore(); err != nil {
			return CoreInfo{}, err
		}
		a.emitEvent("core-status-changed")
		return a.GetCoreInfo()
	}

	if !coreVersionRegex.MatchString(trimmed) {
		return CoreInfo{}, fmt.Errorf("版本号格式不正确：%s", version)
	}

	localDir := localCoreVersionDir(trimmed)
	if !coreFileExists(filepath.Join(localDir, coreExecutableName())) {
		if trimmed == a.coreVersionForDir(a.bundledCoreDir()) {
			if err := a.activateBundledCore(); err != nil {
				return CoreInfo{}, err
			}
			a.emitEvent("core-status-changed")
			return a.GetCoreInfo()
		}
		if err := a.updateCore(trimmed); err != nil {
			return CoreInfo{}, err
		}
		a.emitEvent("core-status-changed")
		return a.GetCoreInfo()
	}

	if err := setMarkerCoreDir(localDir); err != nil {
		return CoreInfo{}, err
	}
	config, _ := a.readAppConfig()
	config.Core.CurrentKnown = trimmed
	_ = a.saveAppConfig(config)

	a.emitEvent("core-status-changed")
	info, err := a.GetCoreInfo()
	if err != nil {
		return info, err
	}
	if info.CurrentVersion != trimmed {
		return info, fmt.Errorf("内核切换后版本校验失败：期望 %s，实际 %s", trimmed, info.CurrentVersion)
	}
	return info, nil
}

// SetCoreAutoUpdate toggles the silent once-a-day engine upgrade.
func (a *App) SetCoreAutoUpdate(enabled bool) (CoreInfo, error) {
	config, err := a.readAppConfig()
	if err != nil {
		return CoreInfo{}, err
	}
	config.Core.AutoUpdate = enabled
	if err := a.saveAppConfig(config); err != nil {
		return CoreInfo{}, err
	}
	if enabled {
		go a.autoUpdateCore()
	}
	return a.GetCoreInfo()
}

func (a *App) coreReleasesCacheFile() string {
	return dataPath("core-releases.json")
}

func (a *App) writeCoreReleasesCache(releases []CoreRelease) error {
	return a.writeJSONFile(a.coreReleasesCacheFile(), releases)
}

func (a *App) readCoreReleasesCache() []CoreRelease {
	data, err := os.ReadFile(a.coreReleasesCacheFile())
	if err != nil {
		return nil
	}
	var releases []CoreRelease
	if err := json.Unmarshal(data, &releases); err != nil {
		return nil
	}
	return releases
}
