package model

// CoreInfo describes the bundled vfox engine plus the cached update verdict.
type CoreInfo struct {
	CurrentVersion  string `json:"currentVersion"`
	BundledVersion  string `json:"bundledVersion"`
	UsesLocalCore   bool   `json:"usesLocalCore"`
	ExecutablePath  string `json:"executablePath"`
	OsArch          string `json:"osArch"`
	LatestVersion   string `json:"latestVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
	ReleaseNotes    string `json:"releaseNotes"`
	ReleaseURL      string `json:"releaseUrl"`
	AutoUpdate      bool   `json:"autoUpdate"`
	LastCheck       string `json:"lastCheck"`
	Error           string `json:"error,omitempty"`
}

// CoreRelease is one entry from the upstream release feed.
type CoreRelease struct {
	Version    string `json:"version"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	Notes      string `json:"notes"`
	URL        string `json:"url"`
	IsCurrent  bool   `json:"isCurrent"`
	Downloaded bool   `json:"downloaded"`
}
