//go:build windows

package app

// openEnvironmentDiagnostic builds the read-only diagnostic report and streams
// it into the in-app terminal instead of opening a separate visible console
// window. The report is read-only: it never mutates PATH, the registry, or any
// file.
func (a *App) openEnvironmentDiagnostic() error {
	return a.streamDiagnosticToTerminal()
}
