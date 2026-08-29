package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	stdruntime "runtime"
)

func (a *App) getCoreDir() string {
	suffix := filepath.Join(coreOSName(), coreArchName())
	candidates := vfoxCoreDirCandidates()

	for _, candidate := range candidates {
		coreDir := filepath.Join(candidate, suffix)
		if coreFileExists(filepath.Join(coreDir, getVfoxExeName())) {
			return coreDir
		}
	}

	baseDir, _ := filepath.Abs("core")
	return filepath.Join(baseDir, suffix)
}

func coreFileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func vfoxCoreDirCandidates() []string {
	var candidates []string
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "core"),
			filepath.Join(exeDir, "..", "Resources", "core"),
			filepath.Join(exeDir, "..", "..", "core"),
		)
	}
	if stdruntime.GOOS == "linux" {
		candidates = append(candidates, "/usr/lib/qfplus/core")
	}
	if abs, err := filepath.Abs("core"); err == nil {
		candidates = append(candidates, abs)
	}
	return candidates
}

// localCoreRoot is the user-writable store for engine versions downloaded by
// the built-in updater. Installing here keeps the shipped core intact, so a
// Program Files installation never needs elevation to switch versions.
func localCoreRoot() string {
	return dataPath("core", coreOSName(), coreArchName())
}

func localCoreVersionDir(version string) string {
	return filepath.Join(localCoreRoot(), version)
}

func coreActiveMarker() string {
	return dataPath("core-active")
}

// markerCoreDir returns the explicitly activated core directory, if any.
func markerCoreDir() string {
	data, err := os.ReadFile(coreActiveMarker())
	if err != nil {
		return ""
	}
	dir := strings.TrimSpace(string(data))
	if dir == "" || !coreFileExists(filepath.Join(dir, getVfoxExeName())) {
		return ""
	}
	return dir
}

// setMarkerCoreDir points the engine lookup at dir; an empty dir clears it.
func setMarkerCoreDir(dir string) error {
	marker := coreActiveMarker()
	if dir == "" {
		if err := os.Remove(marker); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(marker), 0755); err != nil {
		return err
	}
	return os.WriteFile(marker, []byte(dir), 0644)
}

func coreFileCandidates(name string) []string {
	var dirs []string
	if dir := markerCoreDir(); dir != "" {
		dirs = append(dirs, dir)
	}
	suffix := filepath.Join(coreOSName(), coreArchName())
	for _, base := range vfoxCoreDirCandidates() {
		dirs = append(dirs, filepath.Join(base, suffix))
	}
	paths := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		paths = append(paths, filepath.Join(dir, name))
	}
	return paths
}

func findCoreFile(name string) string {
	for _, path := range coreFileCandidates(name) {
		if coreFileExists(path) {
			return path
		}
	}
	return ""
}

func (a *App) getVfoxExecutable() (string, error) {
	if path := findCoreFile(getVfoxExeName()); path != "" {
		return path, nil
	}
	exePath := filepath.Join(a.getCoreDir(), getVfoxExeName())
	return "", fmt.Errorf("vfox core executable not found at %s; install or bundle core/%s/%s", exePath, coreOSName(), coreArchName())
}

func coreOSName() string {
	osName := stdruntime.GOOS
	if osName == "darwin" {
		return "macos"
	}
	return osName
}

func coreArchName() string {
	switch stdruntime.GOARCH {
	case "amd64":
		return "x86_64"
	case "386":
		return "x86"
	default:
		return stdruntime.GOARCH
	}
}
