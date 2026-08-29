package model

import "time"

// EnvironmentGenerator identifies the application that produced the document.
type EnvironmentGenerator struct {
	App     string `json:"app"`
	Version string `json:"version"`
}

// EnvironmentPlatform describes the target platform the document was captured on.
type EnvironmentPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// EnvironmentSdkSource is one of vfox / system / custom.
type EnvironmentSdkSource string

const (
	EnvironmentSdkSourceVfox   EnvironmentSdkSource = "vfox"
	EnvironmentSdkSourceSystem EnvironmentSdkSource = "system"
	EnvironmentSdkSourceCustom EnvironmentSdkSource = "custom"
)

// EnvironmentSdkEntry is a single row in the exported environment document.
type EnvironmentSdkEntry struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	PluginAdded bool     `json:"pluginAdded"`
	Current     string   `json:"current"`
	Versions    []string `json:"versions"`
	Path        string   `json:"path"`
	Version     string   `json:"version"`
}

// EnvironmentDocument is the JSON-serializable snapshot shared with other machines.
// schemaVersion is 1; kind is always "qfplus.environment".
type EnvironmentDocument struct {
	SchemaVersion int                   `json:"schemaVersion"`
	Kind          string                `json:"kind"`
	GeneratedAt   time.Time             `json:"generatedAt"`
	Generator     EnvironmentGenerator  `json:"generator"`
	Platform      EnvironmentPlatform   `json:"platform"`
	VfoxHome      string                `json:"vfoxHome"`
	Sdks          []EnvironmentSdkEntry `json:"sdks"`
	Warnings      []string              `json:"warnings"`
}

// EnvironmentImportResolution describes how a requested version can be satisfied.
type EnvironmentImportResolution string

const (
	EnvironmentResolutionExact            EnvironmentImportResolution = "exact"
	EnvironmentResolutionAlreadyInstalled EnvironmentImportResolution = "already_installed"
	EnvironmentResolutionUnavailable        EnvironmentImportResolution = "unavailable"
	EnvironmentResolutionFallback         EnvironmentImportResolution = "fallback"
	EnvironmentResolutionPathMissing      EnvironmentImportResolution = "path_missing"
	EnvironmentResolutionInvalidName        EnvironmentImportResolution = "invalid_name"
	EnvironmentResolutionNotExported        EnvironmentImportResolution = "not_exported"
)

// EnvironmentImportItem is one line in the human-readable import plan.
type EnvironmentImportItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Source          string `json:"source"`
	Path            string `json:"path"`
	Resolution      string `json:"resolution"`
	FallbackVersion string `json:"fallbackVersion"`
	Action          string `json:"action"`
	SkipReason      string `json:"skipReason"`
	SkipMessage     string `json:"skipMessage"`
	Current         bool   `json:"current"`
}

// EnvironmentImportPlan is what the UI renders before applying.
type EnvironmentImportPlan struct {
	SchemaVersion   int                     `json:"schemaVersion"`
	GeneratedAt     time.Time               `json:"generatedAt"`
	SourceVfoxHome  string                  `json:"sourceVfoxHome"`
	Items           []EnvironmentImportItem `json:"items"`
	FallbackAllowed bool                    `json:"fallbackAllowed"`
	Warnings        []string                `json:"warnings"`
}

// EnvironmentImportResult reports what happened after a plan was applied.
type EnvironmentImportResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Failed   int      `json:"failed"`
	Warnings []string `json:"warnings"`
}

// EnvironmentImportProgress is emitted while the queue is running.
type EnvironmentImportProgress struct {
	Stage   string `json:"stage"`
	Index   int    `json:"index"`
	Total   int    `json:"total"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Phase   string `json:"phase"`
	Status  string `json:"status"`
	Percent int    `json:"percent"`
	Message string `json:"message"`
}

// EnvironmentInventory helps the import planner know what already exists locally.
type EnvironmentInventory struct {
	AddedPlugins  []string             `json:"addedPlugins"`
	InstalledSdks []SdkInfo            `json:"installedSdks"`
	CustomSdksMap map[string][]SdkInfo `json:"customSdksMap"`
}
