//go:build !windows

package app

// openEnvironmentDiagnostic streams the read-only diagnostic report into the
// in-app terminal. On non-Windows platforms there was never a visible console
// window, but the output now goes through the same terminal channel as vfox
// tasks so the UX is consistent.
func (a *App) openEnvironmentDiagnostic() error {
	return a.streamDiagnosticToTerminal()
}
