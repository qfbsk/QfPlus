//go:build windows

package app

import (
	"os"
	"os/exec"
	"path/filepath"
)

// openEnvironmentDiagnostic builds the read-only diagnostic script and opens
// a visible console on Windows.
func (a *App) openEnvironmentDiagnostic() error {
	report, err := a.collectEnvironmentStatus()
	if err != nil {
		return err
	}
	script := a.buildDiagnosticScript(report)
	return openVisibleConsole(script)
}

// openVisibleConsole writes the diagnostic script to a fixed temp file (so a
// re-run overwrites rather than littering) and launches it in a real, visible
// cmd window. The window is the feature itself, so the command is started
// without hideWindow.
func openVisibleConsole(script string) error {
	scriptPath := filepath.Join(os.TempDir(), "QfPlus-env-check.cmd")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return err
	}
	cmd := exec.Command("cmd.exe", "/c", "start", "QfPlus Environment Check", scriptPath)
	if err := cmd.Start(); err != nil {
		return err
	}
	// `start` returns almost immediately; only reap the process handle.
	go func() { _ = cmd.Wait() }()
	return nil
}
