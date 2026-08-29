package model

type AppConfig struct {
	// SchemaVersion is stamped on every save. A config imported from the old
	// vfoxG data directory carries 0 until the first write.
	SchemaVersion int         `json:"schemaVersion"`
	VfoxHome      string      `json:"vfoxHome"`
	Proxy         ProxyConfig `json:"proxy,omitempty"`
	Core          CoreConfig  `json:"core,omitempty"`
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

type DownloadPathInfo struct {
	Path              string `json:"path"`
	DefaultPath       string `json:"defaultPath"`
	IsDefault         bool   `json:"isDefault"`
	HasMigratableData bool   `json:"hasMigratableData"`
}
