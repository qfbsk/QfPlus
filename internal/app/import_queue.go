package app

import (
	"fmt"
	"strings"

	"QfPlus/internal/parser"
)

// runImportQueue applies every installable item in the plan. It acquires the
// vfox task lock once for the whole batch, emits per-item progress, and
// releases the lock even if individual items fail. The batch itself never aborts
// early because a single item failed.
func (a *App) runImportQueue(plan EnvironmentImportPlan) EnvironmentImportResult {
	result := EnvironmentImportResult{Warnings: plan.Warnings}

	releaseTask, err := a.tryStartVfoxTask()
	if err != nil {
		// Should not happen at runtime because the UI guards against concurrent
		// vfox work, but if it does, report it cleanly without the raw lock text.
		result.Failed = len(plan.Items)
		result.Warnings = append(result.Warnings, "Could not start import batch because another task is running.")
		return result
	}
	defer releaseTask()

	a.emitImportProgress(EnvironmentImportProgress{
		Stage:   "installing",
		Index:   0,
		Total:   len(plan.Items),
		Name:    "",
		Version: "",
		Phase:   "",
		Status:  "running",
		Percent: 0,
		Message: "Preparing import batch",
	})

	addedPlugins := a.loadAddedPluginSet(&SdkEnvironmentImportResult{})
	for i, item := range plan.Items {
		progress := EnvironmentImportProgress{
			Stage:   "installing",
			Index:   i + 1,
			Total:   len(plan.Items),
			Name:    item.Name,
			Version: item.Version,
		}

		if item.Action == "skip" {
			progress.Phase = "skip"
			progress.Status = "warning"
			progress.Message = item.SkipMessage
			if progress.Message == "" {
				progress.Message = item.SkipReason
			}
			a.emitImportProgress(progress)
			result.Skipped++
			continue
		}

		a.runImportQueueItem(item, addedPlugins, &progress, &result)
	}

	a.emitImportProgress(EnvironmentImportProgress{
		Stage:   "done",
		Index:   len(plan.Items),
		Total:   len(plan.Items),
		Status:  "ok",
		Percent: 100,
		Message: fmt.Sprintf("Imported %d, skipped %d, failed %d", result.Imported, result.Skipped, result.Failed),
	})

	return result
}

func (a *App) runImportQueueItem(item EnvironmentImportItem, addedPlugins map[string]bool, progress *EnvironmentImportProgress, result *EnvironmentImportResult) {
	name := item.Name
	version := parser.NormalizeSdkVersion(item.Version)
	if item.Resolution == string(EnvironmentResolutionFallback) && item.FallbackVersion != "" {
		version = parser.NormalizeSdkVersion(item.FallbackVersion)
	}

	switch item.Source {
	case "vfox":
		pluginKey := strings.ToLower(name)
		if !addedPlugins[pluginKey] {
			progress.Phase = "add_plugin"
			progress.Status = "running"
			progress.Message = "Adding plugin"
			a.emitImportProgress(*progress)
			if err := a.runVfoxWithProgress([]string{"add", name}); err != nil {
				progress.Status = "failed"
				progress.Message = fmt.Sprintf("Failed to add plugin: %v", err)
				a.emitImportProgress(*progress)
				result.Failed++
				result.Warnings = append(result.Warnings, fmt.Sprintf("%s: plugin add failed: %v", name, err))
				return
			}
			addedPlugins[pluginKey] = true
		}

		progress.Phase = "download"
		progress.Status = "running"
		progress.Message = "Downloading"
		a.emitImportProgress(*progress)
		if err := a.installVersionUnlocked(name, version); err != nil {
			progress.Status = "failed"
			progress.Message = fmt.Sprintf("Failed to install: %v", err)
			a.emitImportProgress(*progress)
			result.Failed++
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s@%s: install failed: %v", name, version, err))
			return
		}

		if item.Current {
			progress.Phase = "activate"
			progress.Status = "running"
			progress.Message = "Activating"
			a.emitImportProgress(*progress)
			if _, err := a.useVersionUnlocked(name, version); err != nil {
				progress.Status = "warning"
				progress.Message = fmt.Sprintf("Installed but failed to activate: %v", err)
				a.emitImportProgress(*progress)
				result.Warnings = append(result.Warnings, fmt.Sprintf("%s@%s: installed but activation failed: %v", name, version, err))
				// Still count as imported; activation failure is a warning.
			} else {
				progress.Status = "ok"
				progress.Message = "Installed and activated"
				a.emitImportProgress(*progress)
			}
		} else {
			progress.Status = "ok"
			progress.Message = "Installed"
			a.emitImportProgress(*progress)
		}
		result.Imported++

	case "custom":
		progress.Phase = "register"
		progress.Status = "running"
		progress.Message = "Registering custom SDK"
		a.emitImportProgress(*progress)

		customMap := a.getNonVfoxSdksMap()
		customVersion := version
		if isUnknownSdkVersion(customVersion) {
			customVersion = a.detectSdkPathVersion(name, item.Path)
		}
		if customVersion == "" {
			customVersion = "unknown"
		}
		customMap[name] = append(customMap[name], SdkInfo{
			Name:     name,
			Source:   "system",
			Path:     item.Path,
			Versions: []SdkVersion{{Version: customVersion}},
		})
		if err := a.saveNonVfoxSdksMap(customMap); err != nil {
			progress.Status = "failed"
			progress.Message = fmt.Sprintf("Failed to save custom SDK: %v", err)
			a.emitImportProgress(*progress)
			result.Failed++
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s: custom SDK save failed: %v", name, err))
			return
		}
		a.emitEvent("sdk-list-changed")
		progress.Status = "ok"
		progress.Message = "Registered custom SDK"
		a.emitImportProgress(*progress)
		result.Imported++

	default:
		progress.Phase = "skip"
		progress.Status = "warning"
		progress.Message = "Unknown source"
		a.emitImportProgress(*progress)
		result.Skipped++
	}
}

func (a *App) emitImportProgress(progress EnvironmentImportProgress) {
	a.emitEvent("environment-import-progress", progress)
}
