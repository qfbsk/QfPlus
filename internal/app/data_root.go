package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	dataDirName       = "QfPlus"
	legacyDataDirName = "vfoxG"
)

// userConfigBase resolves the OS per-user config directory, falling back to
// ~/.config when the platform variable is missing.
func userConfigBase() (string, error) {
	base, err := os.UserConfigDir()
	if err == nil && strings.TrimSpace(base) != "" {
		return base, nil
	}
	if home, homeErr := os.UserHomeDir(); homeErr == nil && strings.TrimSpace(home) != "" {
		return filepath.Join(home, ".config"), nil
	}
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("unable to resolve user config directory")
}

func dataRoot() string {
	base, err := userConfigBase()
	if err != nil {
		return filepath.Join(os.TempDir(), dataDirName)
	}
	return filepath.Join(base, dataDirName)
}

// legacyDataRoot is the pre-rename vfoxG directory. Migration reads it; nothing
// ever writes to it and it is never deleted.
func legacyDataRoot() string {
	base, err := userConfigBase()
	if err != nil {
		return filepath.Join(os.TempDir(), legacyDataDirName)
	}
	return filepath.Join(base, legacyDataDirName)
}

func dataPath(elem ...string) string {
	return filepath.Join(append([]string{dataRoot()}, elem...)...)
}
