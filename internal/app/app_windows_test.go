//go:build windows

package app

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestWindowsSafeShimName(t *testing.T) {
	if got := windowsSafeShimName(` a:b/c*d?e"f<g>h|i `); got != "a_b_c_d_e_f_g_h_i" {
		t.Fatalf("windowsSafeShimName() = %q", got)
	}
	if got := windowsSafeShimName("   "); got != "sdk" {
		t.Fatalf("empty shim name = %q, want sdk", got)
	}
}

func TestWindowsSDKShimAliases(t *testing.T) {
	aliases := windowsSDKShimAliases("python")
	want := []string{"python", "python3", "pip", "pip3"}
	for _, alias := range want {
		if !containsStringFold(aliases, alias) {
			t.Fatalf("missing alias %q in %v", alias, aliases)
		}
	}

	seen := map[string]bool{}
	for _, alias := range aliases {
		key := strings.ToLower(alias)
		if seen[key] {
			t.Fatalf("duplicate alias %q in %v", alias, aliases)
		}
		seen[key] = true
	}
}

func TestWindowsShimTargetExists(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		filepath.Join(root, "bin", "go.exe"),
		filepath.Join(root, "gofmt.cmd"),
		filepath.Join(root, "bin", "notepad"),
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if !windowsShimTargetExists(root, "go") {
		t.Fatal("go should resolve through the bin directory")
	}
	if !windowsShimTargetExists(root, "gofmt") {
		t.Fatal("gofmt should resolve from the SDK root")
	}
	if windowsShimTargetExists(root, "golang") {
		t.Fatal("golang has no executable here and must not get a shim")
	}
	if windowsShimTargetExists(root, "notepad") {
		t.Fatal("a directory without an extension is not an executable")
	}
	if windowsShimTargetExists("", "go") {
		t.Fatal("an unresolved SDK root must not match")
	}
}

func TestWindowsStaleShimAliases(t *testing.T) {
	previous := []string{"golang", "go", "gofmt"}
	got := windowsStaleShimAliases(previous, []string{"GO", "gofmt"})
	if !reflect.DeepEqual(got, []string{"golang"}) {
		t.Fatalf("stale aliases = %v, want [golang]", got)
	}
	if got := windowsStaleShimAliases(previous, previous); got != nil {
		t.Fatalf("nothing should be stale, got %v", got)
	}
}

func TestWindowsShimScriptUsesAllPlaceholders(t *testing.T) {
	script := windowsShimScript("python", "pip", `C:\SDK Root`)
	for _, want := range []string{
		`set "SDK_ROOT=C:\SDK Root"`,
		`%SDK_ROOT%\Scripts\pip.exe`,
		`%SDK_ROOT%\bin\pip.cmd`,
		`where "%ALIAS_NAME%"`,
		`if /I not "%%~fI"=="%~f0"`,
		`QfPlus: pip for python is not available`,
		`no fallback pip was found on PATH`,
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("shim script missing %q:\n%s", want, script)
		}
	}
	if strings.Contains(script, "%!") || strings.Contains(script, "(MISSING)") {
		t.Fatalf("shim script has fmt placeholder leak:\n%s", script)
	}
}

func TestWriteWindowsSDKShimsOnlyCoversRealCommands(t *testing.T) {
	home := t.TempDir()
	app := NewApp()
	app.setVfoxHome(home)

	sdkBin := filepath.Join(home, "sdks", "golang", "bin")
	if err := os.MkdirAll(sdkBin, 0755); err != nil {
		t.Fatal(err)
	}
	for _, exe := range []string{"go.exe", "gofmt.exe"} {
		if err := os.WriteFile(filepath.Join(sdkBin, exe), []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// What an earlier build left behind: the plugin name itself as a shim, plus
	// the alias record that has to name it for it to be cleaned up.
	shimDir := filepath.Join(home, "path-shims")
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shimDir, "golang.cmd"), []byte("@echo off\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "hijacked_paths.json"),
		[]byte(`{"golang":{"Version":2,"Aliases":["golang","go","gofmt"]}}`), 0644); err != nil {
		t.Fatal(err)
	}

	aliases, err := app.writeWindowsSDKShims("golang")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(aliases, []string{"go", "gofmt"}) {
		t.Fatalf("aliases = %v, want [go gofmt]", aliases)
	}
	if _, err := os.Stat(filepath.Join(shimDir, "golang.cmd")); !os.IsNotExist(err) {
		t.Fatalf("dead golang.cmd survived the rewrite: %v", err)
	}
	script, err := os.ReadFile(filepath.Join(shimDir, "go.cmd"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(script), `%SDK_ROOT%\bin\go.exe`) {
		t.Fatalf("go.cmd does not forward through bin:\n%s", script)
	}
}

func TestReadWindowsPathOverrideAliases(t *testing.T) {
	tmp := t.TempDir()
	hijackFile := filepath.Join(tmp, "hijacked_paths.json")
	data := `{
		"python": {
			"Version": 2,
			"Aliases": ["python", "python3", "pip", "pip3"]
		},
		"java": {
			"Version": 2,
			"Aliases": [
				{"value": ["java", "javac"], "Count": 2}
			]
		}
	}`
	if err := os.WriteFile(hijackFile, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		want []string
	}{
		{name: "python", want: []string{"python", "python3", "pip", "pip3"}},
		{name: "java", want: []string{"java", "javac"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := readWindowsPathOverrideAliases(hijackFile, tt.name)
			if !ok {
				t.Fatal("expected aliases to be read")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("aliases = %v, want %v", got, tt.want)
			}
		})
	}
}

func containsStringFold(values []string, want string) bool {
	for _, value := range values {
		if strings.EqualFold(value, want) {
			return true
		}
	}
	return false
}
