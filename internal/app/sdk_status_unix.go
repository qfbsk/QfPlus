//go:build !windows

package app

import (
	"os"

	"QfPlus/internal/model"
)

// envUserPath returns the process PATH; on Unix there is no separate machine
// scope, so the user and machine PATH are the same view.
func envUserPath() string {
	return os.Getenv("PATH")
}

// envMachinePath is empty on Unix: there is no machine PATH concept.
func envMachinePath() string {
	return ""
}

// shimExists is false on Unix: QfPlus does not generate shim files there.
func shimExists(alias, shimDir string) bool {
	return false
}

// buildEnvironmentDiagnosticScript returns "" on Unix because no terminal window
// is launched; OpenEnvironmentDiagnostic reports errNoVisibleConsole instead.
func buildEnvironmentDiagnosticScript(report *model.EnvironmentStatusReport, userPath, machinePath string) string {
	return ""
}
