//go:build windows

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *App) windowsPathShimDir() string {
	return a.getVfoxHomePath("path-shims")
}

func (a *App) writeWindowsSDKShims(pluginName string) ([]string, error) {
	if strings.TrimSpace(pluginName) == "" {
		return nil, fmt.Errorf("plugin name cannot be empty")
	}
	shimDir := a.windowsPathShimDir()
	if shimDir == "" {
		return nil, fmt.Errorf("unable to resolve shim directory")
	}
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		return nil, err
	}
	sdkPath := a.getVfoxHomePath("sdks", pluginName)
	candidates := windowsSDKShimAliases(pluginName)
	aliases := make([]string, 0, len(candidates))
	for _, alias := range candidates {
		if windowsShimTargetExists(sdkPath, alias) {
			aliases = append(aliases, alias)
		}
	}
	if len(aliases) == 0 {
		// Nothing resolved through the shared probe matrix. A shim that reports
		// the miss beats a PATH entry with nothing behind it for layouts this
		// helper cannot see, so fall back to the full candidate list.
		aliases = candidates
	}

	previous, _ := readWindowsPathOverrideAliases(a.getVfoxHomePath("hijacked_paths.json"), pluginName)
	if stale := windowsStaleShimAliases(previous, aliases); len(stale) > 0 {
		if err := a.removeWindowsSDKShims(pluginName, stale); err != nil {
			return nil, err
		}
	}

	for _, alias := range aliases {
		shimName := windowsSafeShimName(alias) + ".cmd"
		shimPath := filepath.Join(shimDir, shimName)
		if err := os.WriteFile(shimPath, []byte(windowsShimScript(pluginName, alias, sdkPath)), 0644); err != nil {
			return nil, err
		}
	}
	return aliases, nil
}

// windowsStaleShimAliases returns the aliases recorded for a plugin on an
// earlier apply that no longer resolve, so switching versions can never leave a
// dead shim shadowing a working command.
func windowsStaleShimAliases(previous []string, keep []string) []string {
	kept := make(map[string]bool, len(keep))
	for _, alias := range keep {
		kept[strings.ToLower(alias)] = true
	}
	var stale []string
	for _, alias := range previous {
		if !kept[strings.ToLower(alias)] {
			stale = append(stale, alias)
		}
	}
	return stale
}

func (a *App) removeWindowsSDKShims(pluginName string, aliases []string) error {
	shimDir := a.windowsPathShimDir()
	if shimDir == "" {
		return nil
	}
	if len(aliases) == 0 {
		aliases = windowsSDKShimAliases(pluginName)
	}
	for _, alias := range aliases {
		_ = os.Remove(filepath.Join(shimDir, windowsSafeShimName(alias)+".cmd"))
	}
	return nil
}
