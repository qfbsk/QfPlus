package app

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// PlanEnvironmentImport reads a JSON environment document and returns a plan.
func (a *App) planEnvironmentImportFromFile(fallbackAllowed bool) (EnvironmentImportPlan, error) {
	var empty EnvironmentImportPlan
	if a.ctx == nil {
		return empty, fmt.Errorf("application context is not ready")
	}

	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Import QfPlus environment",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files (*.json)", Pattern: "*.json"},
		},
	})
	if err != nil {
		return empty, err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return empty, nil
	}

	rows, doc, warnings, err := parseEnvironmentDocument(path)
	if err != nil {
		return empty, err
	}

	plan, err := a.planEnvironmentImport(rows, fallbackAllowed, doc.VfoxHome)
	if err != nil {
		return empty, err
	}
	plan.Warnings = append(plan.Warnings, warnings...)
	return plan, nil
}

// ApplyEnvironmentImport runs the prepared import plan through the queue.
func (a *App) applyEnvironmentImport(plan EnvironmentImportPlan) (EnvironmentImportResult, error) {
	if a.ctx == nil {
		return EnvironmentImportResult{}, fmt.Errorf("application context is not ready")
	}
	return a.runImportQueue(plan), nil
}

// getEnvironmentInventory returns the local state needed to plan imports offline.
func (a *App) getEnvironmentInventory() (EnvironmentInventory, error) {
	inv := EnvironmentInventory{
		AddedPlugins:  []string{},
		InstalledSdks: []SdkInfo{},
		CustomSdksMap: map[string][]SdkInfo{},
	}

	plugins, err := a.getAddedPlugins()
	if err == nil {
		inv.AddedPlugins = plugins
	}

	if sdks, err := a.getInstalledSdks(); err == nil {
		inv.InstalledSdks = sdks
	}

	inv.CustomSdksMap = a.getNonVfoxSdksMap()
	return inv, nil
}
