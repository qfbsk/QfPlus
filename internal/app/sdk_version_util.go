package app

import "strings"

func isUnknownSdkVersion(version string) bool {
	version = strings.TrimSpace(version)
	return version == "" || strings.EqualFold(version, "unknown") || strings.EqualFold(version, "(unknown)")
}
