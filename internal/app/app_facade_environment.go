package app

// GetEnvironmentStatus returns the read-only environment status report.
func (a *App) GetEnvironmentStatus() (EnvironmentStatusReport, error) {
	return a.getEnvironmentStatus()
}

// OpenEnvironmentDiagnostic opens a visible console so the user can verify PATH state.
func (a *App) OpenEnvironmentDiagnostic() error {
	return a.openEnvironmentDiagnostic()
}

// ExportEnvironmentJson writes the current environment document to a .json file.
func (a *App) ExportEnvironmentJson() (string, error) {
	return a.exportEnvironmentJson()
}

// PreviewEnvironmentDocument returns the current environment document.
func (a *App) PreviewEnvironmentDocument() (EnvironmentDocument, error) {
	return a.previewEnvironmentDocument()
}

// PlanEnvironmentImport reads an environment JSON file and builds an import plan.
func (a *App) PlanEnvironmentImport(fallbackAllowed bool) (EnvironmentImportPlan, error) {
	return a.planEnvironmentImportFromFile(fallbackAllowed)
}

// ApplyEnvironmentImport applies a prepared import plan through the batch queue.
func (a *App) ApplyEnvironmentImport(plan EnvironmentImportPlan) (EnvironmentImportResult, error) {
	return a.applyEnvironmentImport(plan)
}

// GetEnvironmentInventory returns local state useful for import planning.
func (a *App) GetEnvironmentInventory() (EnvironmentInventory, error) {
	return a.getEnvironmentInventory()
}
