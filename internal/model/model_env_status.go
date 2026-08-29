package model

import "time"

// SdkCommandStatus describes one SDK command as seen by the live environment:
// which executable wins on PATH, what version answers, and whether QfPlus owns
// a shim for it. It is a read-only observation — never a directive.
type SdkCommandStatus struct {
	SdkName      string   `json:"sdkName"`
	Alias        string   `json:"alias"`
	Source       string   `json:"source"`
	ManagedBy    string   `json:"managedBy"`
	Resolved     bool     `json:"resolved"`
	ExePath      string   `json:"exePath"`
	ExeDir       string   `json:"exeDir"`
	OnUserPath   bool     `json:"onUserPath"`
	OnMachinePath bool    `json:"onMachinePath"`
	Version      string   `json:"version"`
	IsCurrent    bool     `json:"isCurrent"`
	State        string   `json:"state"`
	Notes        []string `json:"notes"`
}

// EnvironmentStatusReport is the full, read-only snapshot consumed by the
// environment status card and the visible diagnostic console.
type EnvironmentStatusReport struct {
	GeneratedAt time.Time           `json:"generatedAt"`
	VfoxHome    string              `json:"vfoxHome"`
	VfoxInPath  bool                `json:"vfoxInPath"`
	ShimDir     string              `json:"shimDir"`
	PathDrift   bool                `json:"pathDrift"`
	Items       []SdkCommandStatus `json:"items"`
	Warnings    []string            `json:"warnings"`
}
