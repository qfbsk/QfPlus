package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"QfPlus/internal/parser"
)

// planEnvironmentImport builds a human-readable plan from parsed rows.
func (a *App) planEnvironmentImport(rows []EnvironmentImportRow, fallbackAllowed bool, sourceVfoxHome string) (EnvironmentImportPlan, error) {
	plan := EnvironmentImportPlan{
		SchemaVersion:   1,
		GeneratedAt:     time.Now(),
		SourceVfoxHome:  sourceVfoxHome,
		Items:           make([]EnvironmentImportItem, 0, len(rows)),
		FallbackAllowed: fallbackAllowed,
	}

	installed, err := a.getInstalledSdks()
	if err != nil {
		plan.Warnings = append(plan.Warnings, "Unable to read installed SDKs: "+err.Error())
	}
	installedMap := make(map[string]bool)
	for _, sdk := range installed {
		for _, v := range sdk.Versions {
			installedMap[strings.ToLower(sdk.Name)+"@"+strings.ToLower(v.Version)] = true
		}
	}

	addedPlugins, err := a.getAddedPlugins()
	if err != nil {
		plan.Warnings = append(plan.Warnings, "Unable to read added plugins: "+err.Error())
	}
	addedMap := make(map[string]bool)
	for _, p := range addedPlugins {
		addedMap[strings.ToLower(strings.TrimSpace(p))] = true
	}

	customMap := a.getNonVfoxSdksMap()
	seen := make(map[string]bool)

	for _, row := range rows {
		item := a.buildImportItem(row, installedMap, addedMap, customMap, fallbackAllowed)
		key := strings.ToLower(item.Name) + "@" + strings.ToLower(item.Version) + "@" + item.Source + "@" + item.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		plan.Items = append(plan.Items, item)
	}

	sort.Slice(plan.Items, func(i, j int) bool {
		if plan.Items[i].Source != plan.Items[j].Source {
			return plan.Items[i].Source < plan.Items[j].Source
		}
		return strings.ToLower(plan.Items[i].Name) < strings.ToLower(plan.Items[j].Name)
	})

	return plan, nil
}

func (a *App) buildImportItem(row EnvironmentImportRow, installedMap map[string]bool, addedMap map[string]bool, customMap map[string][]SdkInfo, fallbackAllowed bool) EnvironmentImportItem {
	name := strings.TrimSpace(row.Name)
	version := parser.NormalizeSdkVersion(row.Version)
	item := EnvironmentImportItem{
		ID:      fmt.Sprintf("%s-%s-%s-%d", row.Kind, strings.ToLower(name), strings.ToLower(version), time.Now().UnixNano()),
		Name:    name,
		Version: version,
		Source:  row.Kind,
		Path:    row.Path,
		Current: row.Current,
	}

	switch row.Kind {
	case "vfox":
		if name == "" || isUnknownSdkVersion(version) {
			item.Resolution = string(EnvironmentResolutionInvalidName)
			item.Action = "skip"
			item.SkipReason = "invalid_name"
			return item
		}

		key := strings.ToLower(name) + "@" + strings.ToLower(version)
		if installedMap[key] {
			item.Resolution = string(EnvironmentResolutionAlreadyInstalled)
			item.Action = "skip"
			item.SkipReason = "already_installed"
			item.SkipMessage = fmt.Sprintf("%s@%s is already installed.", name, version)
			return item
		}

		// If the exact version is already available locally, install it.
		if addedMap[strings.ToLower(name)] {
			item.Resolution = string(EnvironmentResolutionExact)
			item.Action = "install"
			return item
		}

		// Try fallback to the same major version if allowed.
		if fallbackAllowed {
			if fallback, err := a.resolveFallbackVersion(name, version); err == nil && fallback != "" {
				item.Resolution = string(EnvironmentResolutionFallback)
				item.FallbackVersion = fallback
				item.Action = "install"
				item.SkipMessage = fmt.Sprintf("%s@%s unavailable, using fallback %s", name, version, fallback)
				return item
			}
		}

		item.Resolution = string(EnvironmentResolutionUnavailable)
		item.Action = "skip"
		item.SkipReason = "version_unavailable"
		item.SkipMessage = fmt.Sprintf("%s@%s is not available", name, version)
		return item

	case "custom":
		if name == "" || row.Path == "" {
			item.Resolution = string(EnvironmentResolutionInvalidName)
			item.Action = "skip"
			item.SkipReason = "invalid_name"
			return item
		}
		if err := validateSDKExecutablePath(row.Path); err != nil {
			item.Resolution = string(EnvironmentResolutionPathMissing)
			item.Action = "skip"
			item.SkipReason = "path_missing"
			item.SkipMessage = err.Error()
			return item
		}
		for _, existing := range customMap[name] {
			if samePath(existing.Path, row.Path) {
				item.Resolution = string(EnvironmentResolutionAlreadyInstalled)
				item.Action = "skip"
				item.SkipReason = "already_installed"
				item.SkipMessage = fmt.Sprintf("%s at %s is already registered.", name, row.Path)
				return item
			}
		}
		item.Resolution = string(EnvironmentResolutionExact)
		item.Action = "import_custom"
		return item

	default:
		item.Resolution = string(EnvironmentResolutionNotExported)
		item.Action = "skip"
		item.SkipReason = "not_exported"
		return item
	}
}

// resolveFallbackVersion tries to find the highest available version sharing the same major.
func (a *App) resolveFallbackVersion(name string, version string) (string, error) {
	major := parser.MajorVersion(version)
	if major == "" {
		return "", fmt.Errorf("cannot determine major version of %q", version)
	}
	versions, err := a.searchVersions(name)
	if err != nil {
		return "", err
	}
	var best string
	for _, v := range versions {
		if strings.HasPrefix(v, major) && parser.CompareVersions(v, best) > 0 {
			best = v
		}
	}
	if best == "" {
		return "", fmt.Errorf("no %s.x version available for %s", major, name)
	}
	return best, nil
}
