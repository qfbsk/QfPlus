package app

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"QfPlus/internal/model"
)

// collectEnvironmentStatus builds a read-only snapshot of every SDK command the
// environment exposes: which executable wins on PATH, what version answers, and
// whether QfPlus owns a shim for it. It must never mutate PATH, the registry, or
// any file — the diagnostic console depends on that guarantee.
func (a *App) collectEnvironmentStatus() (*model.EnvironmentStatusReport, error) {
	report := &model.EnvironmentStatusReport{
		GeneratedAt: time.Now(),
		VfoxHome:    a.getVfoxHome(),
		ShimDir:     a.getVfoxHomePath("path-shims"),
	}

	userPath := envUserPath()
	machinePath := envMachinePath()

	sdks, _ := a.getAllSdks()
	selections, _ := a.readVfoxGlobalSelections()
	custom := a.getNonVfoxSdksMap()

	// The status page only reports on the SDKs this vfox instance actually
	// manages (vfox-added plugins) plus locally-registered custom SDKs. It does
	// NOT scan the whole system-SDK catalog — that list is surfaced on the
	// Installed page. This keeps the status page aligned with what vfox
	// supports here instead of a hardcoded grab-bag.
	pluginNames := map[string]string{} // name -> source
	for _, sdk := range sdks {
		if sdk.Source == "system" {
			continue
		}
		pluginNames[sdk.Name] = sdk.Source
	}
	for name := range custom {
		if _, ok := pluginNames[name]; !ok {
			pluginNames[name] = "custom"
		}
	}

	addedPlugins, _ := a.getAddedPlugins()

	// Probe with the LIVE registry PATH, not the process environment frozen at
	// startup. QfPlus writes new shims into the registry PATH at runtime, so a
	// manual refresh must re-scan that live PATH to reflect freshly added SDKs
	// without forcing the user to restart the app. Other environment variables
	// are preserved so subprocess probes (where / go version) still work.
	probeEnv := livePathEnv(os.Environ(), userPath, machinePath)

	for name, source := range pluginNames {
		isCurrent := selections[name] != ""
		for _, alias := range pluginCommandAliases(name) {
			item := a.probeCommand(name, alias, source, isCurrent, probeEnv, userPath, machinePath, addedPlugins)
			report.Items = append(report.Items, item)
		}
	}

	report.VfoxInPath = a.vfoxShimDirOnPath(userPath, machinePath)
	report.PathDrift = pathContainsDir(machinePath, report.ShimDir) &&
		!pathContainsDir(os.Getenv("PATH"), report.ShimDir)
	return report, nil
}

// getEnvironmentStatus returns the read-only status report in the shape the
// frontend expects.
func (a *App) getEnvironmentStatus() (EnvironmentStatusReport, error) {
	report, err := a.collectEnvironmentStatus()
	if err != nil {
		return EnvironmentStatusReport{}, err
	}
	return *report, nil
}

// probeCommand resolves one alias against the live PATH and classifies it.
func (a *App) probeCommand(pluginName, alias, source string, isCurrent bool, env []string, userPath, machinePath string, addedPlugins []string) model.SdkCommandStatus {
	item := model.SdkCommandStatus{
		SdkName:   pluginName,
		Alias:     alias,
		Source:    source,
		IsCurrent: isCurrent,
	}

		hits := findExecutableCandidates(alias, env)
	if len(hits) == 0 {
		managedBy := a.shimOwningPlugin(alias, addedPlugins)
		facts := sdkCommandFacts{resolved: false, managedBy: managedBy}
		item.State = classifyState(facts)
		if managedBy != "" {
			item.ManagedBy = managedBy
			item.Notes = append(item.Notes, "QfPlus manages a shim for this command, but the target is not currently resolvable")
		} else if alias == pluginName {
			item.Notes = append(item.Notes, aliasCommandHint(pluginName))
		}
		return item
	}

	exePath := hits[0]
	item.Resolved = true
	item.ExePath = exePath
	item.ExeDir = filepath.Dir(exePath)
	item.OnUserPath = pathContainsDir(userPath, item.ExeDir)
	item.OnMachinePath = pathContainsDir(machinePath, item.ExeDir)

	// Windows Store redirector stubs (e.g. C:\Users\...\WindowsApps\python.exe)
	// hang or print Store nudges instead of a version, so skip the probe and
	// classify them as broken immediately. This also keeps the diagnostic
	// console from stalling for seconds per stub.
	if isStoreStub(exePath) {
		facts := sdkCommandFacts{
			resolved:      true,
			managedBy:     a.shimOwningPlugin(alias, addedPlugins),
			onUserPath:    item.OnUserPath,
			onMachinePath: item.OnMachinePath,
			broken:        true,
		}
		item.State = classifyState(facts)
		item.Notes = append(item.Notes, "Windows Store stub detected; install a real interpreter or remove the Store alias from PATH")
		return item
	}

	// Some aliases (e.g. gofmt) ship with a toolchain but expose no version flag.
	// For those we still confirm the executable resolves, then skip the version
	// probe and report the command as working without a version string.
	verArgs := versionArgsFor(alias)
	if len(verArgs) > 0 {
		version := a.tryGetVersionWithEnv(exePath, verArgs, env)
		if !isUsableSystemVersion(version) {
			facts := sdkCommandFacts{
				resolved:      true,
				managedBy:     a.shimOwningPlugin(alias, addedPlugins),
				onUserPath:    item.OnUserPath,
				onMachinePath: item.OnMachinePath,
				broken:        true,
			}
			item.State = classifyState(facts)
			item.Notes = append(item.Notes, "command resolves on PATH but its version probe returned no usable output")
			return item
		}
		item.Version = version
	}

	managedBy := a.shimOwningPlugin(alias, addedPlugins)
	facts := sdkCommandFacts{
		resolved:      true,
		managedBy:     managedBy,
		onUserPath:    item.OnUserPath,
		onMachinePath: item.OnMachinePath,
	}
	item.State = classifyState(facts)
	if managedBy != "" {
		item.ManagedBy = managedBy
	}
	return item
}

// shimOwningPlugin returns the plugin name that owns a QfPlus shim for alias,
// or "" when no shim exists for it. It only considers plugins this vfox
// instance has actually added, so the mapping stays in sync with what is
// managed here rather than the full system-SDK catalog.
func (a *App) shimOwningPlugin(alias string, addedPlugins []string) string {
	shimDir := a.getVfoxHomePath("path-shims")
	if !shimExists(alias, shimDir) {
		return ""
	}
	for _, name := range addedPlugins {
		for _, cand := range pluginCommandAliases(name) {
			if strings.EqualFold(cand, alias) {
				return name
			}
		}
	}
	return ""
}

// pathContainsDir reports whether dir is one of the ;/-separated entries of
// pathValue (case-insensitive on Windows).
func pathContainsDir(pathValue, dir string) bool {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return false
	}
	for _, entry := range filepath.SplitList(pathValue) {
		if strings.EqualFold(strings.TrimSpace(entry), dir) {
			return true
		}
	}
	return false
}

// isStoreStub reports whether exePath is a Windows Store redirector under
// %LOCALAPPDATA%\Microsoft\WindowsApps. Running those stubs usually launches
// the Store or hangs, so callers should skip version probing.
func isStoreStub(exePath string) bool {
	lower := strings.ToLower(exePath)
	return strings.Contains(lower, `\microsoft\windowsapps\`) &&
		(strings.HasSuffix(lower, ".exe") || strings.HasSuffix(lower, ".cmd"))
}

// vfoxShimDirOnPath reports whether the QfPlus shim directory is visible on the
// user or machine PATH.
func (a *App) vfoxShimDirOnPath(userPath, machinePath string) bool {
	shimDir := a.getVfoxHomePath("path-shims")
	if strings.TrimSpace(shimDir) == "" {
		return false
	}
	return pathContainsDir(userPath, shimDir) || pathContainsDir(machinePath, shimDir)
}

// livePathEnv returns a copy of env with the PATH entry replaced by the live
// user + machine PATH a freshly opened console would see. QfPlus adds SDK
// shims to the registry PATH while running, so reusing the process environment
// (captured at launch) would miss them until the app restarts. Callers that
// probe the environment for command resolution must use this live PATH, which
// is what makes a manual refresh reflect newly added SDKs immediately.
func livePathEnv(env []string, userPath, machinePath string) []string {
	live := combinePathScopes(userPath, machinePath)
	if live == "" {
		return env
	}
	const pathKey = "PATH="
	out := make([]string, 0, len(env)+1)
	replaced := false
	for _, kv := range env {
		if len(kv) >= len(pathKey) && strings.EqualFold(kv[:len(pathKey)], pathKey) {
			out = append(out, pathKey+live)
			replaced = true
			continue
		}
		out = append(out, kv)
	}
	if !replaced {
		out = append(out, pathKey+live)
	}
	return out
}

// combinePathScopes joins the user and machine PATH scopes the way a new
// console resolves them. Empty scopes are skipped.
func combinePathScopes(userPath, machinePath string) string {
	parts := make([]string, 0, 2)
	if strings.TrimSpace(userPath) != "" {
		parts = append(parts, userPath)
	}
	if strings.TrimSpace(machinePath) != "" {
		parts = append(parts, machinePath)
	}
	return strings.Join(parts, string(os.PathListSeparator))
}
