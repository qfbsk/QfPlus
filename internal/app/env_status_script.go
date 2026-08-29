package app

import (
	"errors"
	"fmt"
	"strings"
)

// errNoVisibleConsole is returned when a platform cannot launch a real terminal
// window. Only Windows opens one; every other platform returns this from
// OpenEnvironmentDiagnostic.
var errNoVisibleConsole = errors.New("a visible diagnostic console is only available on Windows")

// diagnosticReadOnlyNotice is printed at the top of every diagnostic script so
// the user can see, in the window itself, that nothing was changed.
const diagnosticReadOnlyNotice = "This window is READ-ONLY. QfPlus did not modify PATH or any file."

// buildDiagnosticScript generates a batch script that prints every SDK alias,
// its resolved path, version, and PATH membership without making any changes.
func (a *App) buildDiagnosticScript(report *EnvironmentStatusReport) string {
	var b strings.Builder
	writeLine := func(format string, args ...interface{}) {
		if len(args) == 0 {
			b.WriteString(format)
		} else {
			b.WriteString(fmt.Sprintf(format, args...))
		}
		b.WriteString("\r\n")
	}

	writeLine("@echo off")
	writeLine("title QfPlus Environment Check")
	writeLine("echo %s", diagnosticReadOnlyNotice)
	writeLine("echo.")
	writeLine("echo VfoxHome: %s", report.VfoxHome)
	writeLine("echo ShimDir:  %s", report.ShimDir)
	writeLine("echo VfoxInPath: %v", report.VfoxInPath)
	writeLine("echo.")
	writeLine("echo Current terminal PATH:")
	writeLine("echo %s", "%PATH%")
	writeLine("echo.")
	writeLine("echo Resolved SDK commands:")
	for _, item := range report.Items {
		state := item.State
		if state == "" {
			state = "unknown"
		}
		writeLine("echo %s (%s): %s", item.Alias, item.SdkName, state)
		if item.ExePath != "" {
			writeLine("echo   Path: %s", item.ExePath)
			writeLine("echo   Dir:  %s", item.ExeDir)
		}
		if item.Version != "" {
			writeLine("echo   Version: %s", item.Version)
		}
		writeLine("echo   UserPATH: %v  MachinePATH: %v", item.OnUserPath, item.OnMachinePath)
		writeLine("echo.")
	}
	writeLine("echo.")
	writeLine("echo Open a NEW terminal window if the PATH above does not match what you expect.")
	writeLine("pause")
	return b.String()
}
