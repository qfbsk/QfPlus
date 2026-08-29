package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const migrationMarkerFile = ".migration-v1"

// legacyStateFiles hold app state that must be present before the first screen
// renders. They are copied, never moved: the old directory can still contain a
// live mihomo.pid or a locked cache.db, and a half-failed move is the one way to
// destroy the only irreplaceable file in there.
var legacyStateFiles = []string{
	"config.json",
	"core-releases.json",
	filepath.Join("proxy", "config.yaml"),
}

func legacyMigrationComplete() bool {
	_, err := os.Stat(filepath.Join(dataRoot(), migrationMarkerFile))
	return err == nil
}

func markLegacyMigrationComplete() {
	if err := os.MkdirAll(dataRoot(), 0755); err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dataRoot(), migrationMarkerFile), []byte(time.Now().Format(time.RFC3339)), 0644)
}

// ensureLegacyAppConfig imports the state written by vfoxG builds. It must run
// before the first config read, otherwise saving a fresh config would shadow the
// only copy of the user's proxy subscription.
func (a *App) ensureLegacyAppConfig() error {
	if legacyMigrationComplete() {
		return nil
	}
	oldRoot := legacyDataRoot()
	info, err := os.Stat(oldRoot)
	if err != nil || !info.IsDir() {
		markLegacyMigrationComplete()
		return nil
	}

	var failures []string
	copied := 0
	for _, name := range legacyStateFiles {
		src := filepath.Join(oldRoot, name)
		dst := filepath.Join(dataRoot(), name)
		if _, err := os.Stat(dst); err == nil {
			continue
		}
		srcInfo, err := os.Stat(src)
		if err != nil || srcInfo.IsDir() {
			continue
		}
		if err := copyFileNoOverwrite(src, dst, srcInfo.Mode().Perm()); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", name, err))
			continue
		}
		copied++
	}

	if len(failures) > 0 {
		return fmt.Errorf("could not import %s: %s", legacyDataDirName, strings.Join(failures, "; "))
	}
	if copied > 0 {
		a.emitEvent("vfox-log", "[INFO] Imported settings from "+oldRoot)
	}
	markLegacyMigrationComplete()
	return nil
}

// migrateLegacyVfoxHome relocates the SDK data directory when it still sits
// inside the old app directory. Called from a goroutine: it can copy gigabytes.
func (a *App) migrateLegacyVfoxHome() {
	config, err := a.readAppConfig()
	if err != nil {
		return
	}
	source := strings.TrimSpace(config.VfoxHome)
	if source == "" {
		source = filepath.Join(legacyDataRoot(), "vfox-home")
	}
	if !isPathWithin(source, legacyDataRoot()) {
		return
	}
	hasData, err := hasMigratableVfoxHomeData(source)
	if err != nil || !hasData {
		return
	}

	target, err := normalizeDownloadPath(filepath.Join(dataRoot(), "vfox-home"))
	if err != nil {
		return
	}
	a.emitEvent("vfox-log", "[INFO] Moving vfox SDK data out of "+legacyDataDirName)
	if err := a.migrateVfoxHomeData(source, target); err != nil {
		a.reportMigrationFailure(err)
		return
	}

	config, err = a.readAppConfig()
	if err != nil {
		return
	}
	config.VfoxHome = target
	if err := a.saveAppConfig(config); err != nil {
		a.reportMigrationFailure(err)
		return
	}
	a.setVfoxHome(target)
	a.emitEvent("vfox-home-changed")
	a.emitEvent("sdk-list-changed")
}

// reportMigrationFailure re-emits after a delay because events fired while the
// webview is still loading reach no subscriber.
func (a *App) reportMigrationFailure(err error) {
	message := "[APP ERROR] Data migration failed: " + err.Error()
	a.emitEvent("vfox-log", message)
	time.AfterFunc(2500*time.Millisecond, func() {
		a.emitEvent("vfox-log", message)
	})
}
