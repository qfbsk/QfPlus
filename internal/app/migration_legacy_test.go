package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sandboxDataRoots points both the legacy vfoxG and the current QfPlus data
// roots at one temporary directory so tests never touch the real AppData.
func sandboxDataRoots(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("APPDATA", root)
	t.Setenv("AppData", root)
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("HOME", root)
	t.Setenv("USERPROFILE", root)
	return root
}

func writeStateFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureLegacyAppConfigImportsState(t *testing.T) {
	root := sandboxDataRoots(t)
	old := filepath.Join(root, legacyDataDirName)
	writeStateFile(t, filepath.Join(old, "config.json"), `{"vfoxHome":"D:\\sdk-data"}`)
	writeStateFile(t, filepath.Join(old, "core-releases.json"), `[]`)
	writeStateFile(t, filepath.Join(old, "proxy", "config.yaml"), "mixed-port: 17890\n")
	writeStateFile(t, filepath.Join(old, "proxy", "mihomo.pid"), "4242")

	if err := NewApp().ensureLegacyAppConfig(); err != nil {
		t.Fatal(err)
	}

	config, err := NewApp().readAppConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.VfoxHome != `D:\sdk-data` {
		t.Fatalf("imported config lost VfoxHome: %#v", config)
	}
	if !legacyMigrationComplete() {
		t.Fatal("marker must be written once the import succeeds")
	}
	for _, name := range []string{"core-releases.json", filepath.Join("proxy", "config.yaml")} {
		if _, err := os.Stat(dataPath(name)); err != nil {
			t.Fatalf("expected import of %s: %v", name, err)
		}
	}
	// Runtime state belongs to whichever process is alive; importing a stale pid
	// file would make the new build kill an unrelated process.
	if _, err := os.Stat(filepath.Join(dataRoot(), "proxy", "mihomo.pid")); !os.IsNotExist(err) {
		t.Fatal("mihomo.pid must not be imported")
	}
	if _, err := os.Stat(filepath.Join(old, "config.json")); err != nil {
		t.Fatal("the legacy directory must be left intact")
	}
}

func TestEnsureLegacyAppConfigKeepsNewerFiles(t *testing.T) {
	sandboxDataRoots(t)
	writeStateFile(t, filepath.Join(legacyDataRoot(), "config.json"), `{"vfoxHome":"legacy"}`)
	writeStateFile(t, appConfigFile(), `{"vfoxHome":"current"}`)

	if err := NewApp().ensureLegacyAppConfig(); err != nil {
		t.Fatal(err)
	}
	config, err := NewApp().readAppConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.VfoxHome != "current" {
		t.Fatalf("import must never overwrite an existing config, got %q", config.VfoxHome)
	}
}

func TestEnsureLegacyAppConfigRunsOnce(t *testing.T) {
	sandboxDataRoots(t)
	app := NewApp()
	if err := app.ensureLegacyAppConfig(); err != nil {
		t.Fatal(err)
	}
	writeStateFile(t, filepath.Join(legacyDataRoot(), "config.json"), `{"vfoxHome":"late"}`)
	if err := app.ensureLegacyAppConfig(); err != nil {
		t.Fatal(err)
	}
	config, err := app.readAppConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.VfoxHome == "late" {
		t.Fatal("a completed migration must not re-import on the next launch")
	}
}

func TestMigrateLegacyVfoxHomeRelocatesDefaultData(t *testing.T) {
	sandboxDataRoots(t)
	source := filepath.Join(legacyDataRoot(), "vfox-home")
	prepareMigratableVfoxHome(t, source)

	app := NewApp()
	config, err := app.readAppConfig()
	if err != nil {
		t.Fatal(err)
	}
	config.VfoxHome = source
	if err := app.saveAppConfig(config); err != nil {
		t.Fatal(err)
	}

	app.migrateLegacyVfoxHome()

	var migrated AppConfig
	if err := readJSON(t, appConfigFile(), &migrated); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dataRoot(), "vfox-home")
	if !samePath(migrated.VfoxHome, want) {
		t.Fatalf("VfoxHome = %q, want %q", migrated.VfoxHome, want)
	}
	if _, err := os.Stat(filepath.Join(want, "cache", "python", "version.txt")); err != nil {
		t.Fatalf("SDK data not relocated: %v", err)
	}
	if !strings.HasPrefix(filepath.Clean(app.getVfoxHome()), filepath.Clean(dataRoot())) {
		t.Fatalf("running instance still points at the legacy home: %q", app.getVfoxHome())
	}
}

func TestMigrateLegacyVfoxHomeLeavesCustomPathAlone(t *testing.T) {
	root := sandboxDataRoots(t)
	custom := filepath.Join(root, "D", "qfpacket", "vFox")
	prepareMigratableVfoxHome(t, custom)

	app := NewApp()
	config, err := app.readAppConfig()
	if err != nil {
		t.Fatal(err)
	}
	config.VfoxHome = custom
	if err := app.saveAppConfig(config); err != nil {
		t.Fatal(err)
	}

	app.migrateLegacyVfoxHome()

	var unchanged AppConfig
	if err := readJSON(t, appConfigFile(), &unchanged); err != nil {
		t.Fatal(err)
	}
	if !samePath(unchanged.VfoxHome, custom) {
		t.Fatalf("a home outside the app directory must not move, got %q", unchanged.VfoxHome)
	}
}

func readJSON(t *testing.T, path string, target interface{}) error {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
