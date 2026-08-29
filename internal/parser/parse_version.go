package parser

import (
	"strconv"
	"strings"
)

func NormalizeSdkVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	return version
}

func SameSdkVersion(leftVersion string, rightVersion string) bool {
	return strings.EqualFold(NormalizeSdkVersion(leftVersion), NormalizeSdkVersion(rightVersion))
}

// MajorVersion returns the leading numeric segment of a version, e.g. "20.11.1" -> "20".
func MajorVersion(version string) string {
	version = NormalizeSdkVersion(version)
	if version == "" {
		return ""
	}
	parts := strings.Split(version, ".")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return ""
	}
	return parts[0]
}

// CompareVersions compares two normalized dot-separated versions.
// Returns >0 if left is newer, <0 if right is newer, 0 if equal.
func CompareVersions(left string, right string) int {
	left = NormalizeSdkVersion(left)
	right = NormalizeSdkVersion(right)
	if left == "" && right == "" {
		return 0
	}
	if left == "" {
		return -1
	}
	if right == "" {
		return 1
	}
	lp := strings.Split(left, ".")
	rp := strings.Split(right, ".")
	maxLen := len(lp)
	if len(rp) > maxLen {
		maxLen = len(rp)
	}
	for i := 0; i < maxLen; i++ {
		var ln, rn int
		if i < len(lp) {
			ln, _ = strconv.Atoi(lp[i])
		}
		if i < len(rp) {
			rn, _ = strconv.Atoi(rp[i])
		}
		if ln != rn {
			if ln > rn {
				return 1
			}
			return -1
		}
	}
	return 0
}

func trimSdkCurrentVersionMarker(line string) (string, bool) {
	for _, marker := range []string{"<-- current", "<— current", "<- current"} {
		if strings.HasSuffix(line, marker) {
			return strings.TrimSpace(strings.TrimSuffix(line, marker)), true
		}
	}
	return line, false
}
