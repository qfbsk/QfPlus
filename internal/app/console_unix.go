//go:build !windows

package app

// openEnvironmentDiagnostic is unsupported off Windows.
func (a *App) openEnvironmentDiagnostic() error {
	return errNoVisibleConsole
}

// openVisibleConsole is unsupported off Windows. Callers should surface
// errNoVisibleConsole to the user and fall back to the in-app report.
func openVisibleConsole(script string) error {
	return errNoVisibleConsole
}
