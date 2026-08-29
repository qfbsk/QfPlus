//go:build windows

package app

import (
	"strings"
	"testing"
	"time"

	"QfPlus/internal/model"
)

func TestDiagnosticScriptReadOnly(t *testing.T) {
	report := &model.EnvironmentStatusReport{
		GeneratedAt: time.Now(),
		ShimDir:     `C:\Users\me\QfPlus\vfox-home\path-shims`,
		Items: []model.SdkCommandStatus{
			{Alias: "node", SdkName: "nodejs", Source: "vfox", State: "ok", Resolved: true, ExePath: `C:\nvm\node.exe`, ExeDir: `C:\nvm`, Version: "20.11.1"},
			{Alias: "go", SdkName: "golang", Source: "vfox", State: "missing", Resolved: false},
		},
	}
	script := buildEnvironmentDiagnosticScript(report, `C:\nvm`, `C:\Program Files\Go\bin`)

	forbidden := []string{
		"setx", "reg add", "reg delete", "mklink",
		"New-Item", "Remove-Item", "Move-Item", "Copy-Item",
		"SetEnvironmentVariable('Path',", ">nul", "2>nul",
	}
	for _, token := range forbidden {
		if strings.Contains(script, token) {
			t.Errorf("diagnostic script must not contain %q", token)
		}
	}

	for _, want := range []string{"node", "go", "where", diagnosticReadOnlyNotice} {
		if !strings.Contains(script, want) {
			t.Errorf("diagnostic script should mention %q", want)
		}
	}
}
