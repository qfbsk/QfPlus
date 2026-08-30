package model

type AppConfig struct {
	// SchemaVersion is stamped on every save. A config imported from the old
	// vfoxG data directory carries 0 until the first write.
	SchemaVersion int                `json:"schemaVersion"`
	VfoxHome      string             `json:"vfoxHome"`
	Proxy         ProxyConfig        `json:"proxy,omitempty"`
	Core          CoreConfig         `json:"core,omitempty"`
	GitHubSource  GitHubSourceConfig `json:"githubSource,omitempty"`
}

// CoreConfig persists the vfox engine update policy. LatestVersion is only a
// cache of the newest release seen on the network.
type CoreConfig struct {
	AutoUpdate   bool   `json:"autoUpdate,omitempty"`
	LatestKnown  string `json:"latestKnown,omitempty"`
	LastCheck    string `json:"lastCheck,omitempty"`
	CurrentKnown string `json:"currentKnown,omitempty"`
}

type ProxyConfig struct {
	Enabled         bool   `json:"enabled"`
	SubscriptionURL string `json:"subscriptionUrl,omitempty"`
	MixedPort       int    `json:"mixedPort,omitempty"`
	APIPort         int    `json:"apiPort,omitempty"`
	APISecret       string `json:"apiSecret,omitempty"`
	SelectedGroup   string `json:"selectedGroup,omitempty"`
	SelectedNode    string `json:"selectedNode,omitempty"`
}

// GitHubSourceConfig persists the user's GitHub mirror / accelerator settings.
// It is independent from the bundled proxy and can be used together with it.
type GitHubSourceConfig struct {
	Enabled       bool           `json:"enabled"`
	SelectedID    string         `json:"selectedId,omitempty"`
	CustomSources []GitHubSource `json:"customSources,omitempty"`
}

// GitHubSource describes a single GitHub mirror or accelerator endpoint.
type GitHubSource struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	IsPreset bool   `json:"isPreset"`
}

// GitHubSourceSettings is the DTO returned to the frontend. It contains both
// built-in presets and the user's saved custom sources.
type GitHubSourceSettings struct {
	Enabled        bool           `json:"enabled"`
	SelectedID     string         `json:"selectedId"`
	Sources        []GitHubSource `json:"sources"`
	PresetSources  []GitHubSource `json:"presetSources"`
	CustomSources  []GitHubSource `json:"customSources"`
}

type DownloadPathInfo struct {
	Path              string `json:"path"`
	DefaultPath       string `json:"defaultPath"`
	IsDefault         bool   `json:"isDefault"`
	HasMigratableData bool   `json:"hasMigratableData"`
}
