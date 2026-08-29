//go:build windows

package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"QfPlus/internal/model"
)

// envUserPath returns the current user PATH as seen by a new console.
func envUserPath() string {
	return powershellEnvPath("User")
}

// envMachinePath returns the machine (system) PATH as seen by a new console.
func envMachinePath() string {
	return powershellEnvPath("Machine")
}

// powershellEnvPath reads one of the persisted PATH scopes via PowerShell. This
// is the only subprocess the status probe uses, and it is strictly read-only.
func powershellEnvPath(scope string) string {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf("[Environment]::GetEnvironmentVariable('Path', '%s')", scope))
	// Read-only probe: never show a console window.
	hideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(out), "\r\n")
}

// shimExists reports whether QfPlus has written a shim for alias into shimDir.
func shimExists(alias, shimDir string) bool {
	if strings.TrimSpace(shimDir) == "" {
		return false
	}
	for _, ext := range []string{".cmd", ".exe"} {
		if _, err := os.Stat(filepath.Join(shimDir, alias+ext)); err == nil {
			return true
		}
	}
	return false
}

// buildEnvironmentDiagnosticScript assembles a read-only batch script that
// prints the report and the PATH a new console will inherit. It must never
// contain a write operation (see the plan's read-only guarantee).
func buildEnvironmentDiagnosticScript(report *model.EnvironmentStatusReport, userPath, machinePath string) string {
	var b strings.Builder
	b.WriteString("@echo off\n")
	b.WriteString("title QfPlus Environment Check\n")
	b.WriteString("echo QfPlus Environment Diagnostic\n")
	b.WriteString("echo Generated: " + report.GeneratedAt.Format("2006-01-02 15:04:05") + "\n")
	b.WriteString("echo.\n")
	b.WriteString("echo " + diagnosticReadOnlyNotice + "\n")
	b.WriteString("echo.\n")

	for _, item := range report.Items {
		b.WriteString("echo ========================================\n")
		b.WriteString(fmt.Sprintf("echo Command: %s  (plugin %s, %s)\n", item.Alias, item.SdkName, item.Source))
		b.WriteString("echo State: " + item.State + "\n")
		if item.Resolved {
			b.WriteString("echo Resolved: " + item.ExePath + "\n")
			b.WriteString("echo Version: " + item.Version + "\n")
			b.WriteString("echo On user PATH: " + boolStr(item.OnUserPath) + "   On machine PATH: " + boolStr(item.OnMachinePath) + "\n")
		} else {
			b.WriteString("echo Not resolved on PATH.\n")
		}
		for _, note := range item.Notes {
			b.WriteString("echo Note: " + note + "\n")
		}
		b.WriteString("where " + item.Alias + "\n")
		b.WriteString("echo.\n")
	}

	b.WriteString("echo ----------------------------------------\n")
	b.WriteString("echo PATH this console sees (user):\n")
	b.WriteString("echo " + userPath + "\n")
	b.WriteString("echo.\n")
	b.WriteString("echo PATH a NEW console will see (machine):\n")
	b.WriteString("echo " + machinePath + "\n")
	b.WriteString("echo.\n")
	b.WriteString("echo If a command you added is missing, close and reopen your terminal.\n")
	b.WriteString("pause\n")
	return b.String()
}

func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
