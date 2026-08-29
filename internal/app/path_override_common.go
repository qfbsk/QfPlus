package app

import (
	"errors"
	"fmt"
	"strings"

	"QfPlus/internal/parser"
)

// errNoActiveSdkVersion is returned when SDK PATH is enabled for a vfox plugin
// that has no version in global use. The frontend matches on this text to show
// localized guidance, so keep both sides in sync with taskLogParser.ts.
var errNoActiveSdkVersion = errors.New("no version is in use")

// ensureVfoxSdkJunction points sdks/<name> at the version in global use. Every
// managed SDK PATH entry resolves through that link, so a missing link leaves
// a PATH entry that can never find the executable.
func (a *App) ensureVfoxSdkJunction(name string) error {
	out, err := a.runVfoxCommandWithLock("current", name)
	version := parser.CurrentSdkVersion(name, out)
	if version == "" {
		// vfox exits non-zero when nothing is in global use, but genuine failures
		// such as a missing binary or a timeout must still surface unchanged.
		if err != nil && !strings.Contains(strings.ToLower(out), "no current") {
			return err
		}
		return fmt.Errorf("%w: %s", errNoActiveSdkVersion, name)
	}
	runtimeRoot, err := a.resolveVersionRuntimeRoot(name, version)
	if err != nil {
		return err
	}
	return a.ensureJunction(a.getVfoxHomePath("sdks", name), runtimeRoot)
}

func (a *App) checkPluginPathOverride(pluginName string) bool {
	return a.checkPluginWin11CompatMode(pluginName)
}

func (a *App) checkAnyPathOverride() bool {
	return a.checkWin11CompatMode()
}
