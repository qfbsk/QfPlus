//go:build !windows

package app

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestUnixManagedPathBlockQuotesAndPreservesPath(t *testing.T) {
	block := unixManagedPathBlock(vfoxSDKPathMarkerPrefix+"py'thon", []string{
		"/opt/python/bin",
		"/Applications/SDK Tools/bin",
		"",
	})

	if !strings.Contains(block, "# >>> "+vfoxSDKPathMarkerPrefix+"py'thon >>>") {
		t.Fatalf("missing start marker: %s", block)
	}
	if !strings.Contains(block, "'/opt/python/bin'") {
		t.Fatalf("missing quoted unix path: %s", block)
	}
	if !strings.Contains(block, "'/Applications/SDK Tools/bin'") {
		t.Fatalf("path with spaces must be quoted: %s", block)
	}
	if !strings.Contains(block, `:"$PATH"`) {
		t.Fatalf("block must preserve existing PATH at the end: %s", block)
	}
}

func TestUnixRemoveManagedBlockFromString(t *testing.T) {
	original := strings.Join([]string{
		"export KEEP=1",
		unixManagedPathBlock(vfoxPathMarkerLabel, []string{"/opt/vfox"}),
		"export AFTER=1",
	}, "\n")

	got := unixRemoveManagedBlockFromString(original, vfoxPathMarkerLabel)
	if strings.Contains(got, vfoxPathMarkerLabel) || strings.Contains(got, "/opt/vfox") {
		t.Fatalf("managed block was not removed: %s", got)
	}
	if !strings.Contains(got, "export KEEP=1") || !strings.Contains(got, "export AFTER=1") {
		t.Fatalf("unrelated profile lines should be preserved: %s", got)
	}
}

func TestUnixRemoveManagedBlockLeavesBrokenBlockUnchanged(t *testing.T) {
	data := "# >>> " + vfoxPathMarkerLabel + " >>>\nexport PATH='/opt/vfox':\"$PATH\"\n"
	got := unixRemoveManagedBlockFromString(data, vfoxPathMarkerLabel)
	if got != data {
		t.Fatalf("unterminated block should be left untouched: got %q want %q", got, data)
	}
}

func TestUnixManagedBlockLabelsIncludesLegacyName(t *testing.T) {
	labels := unixManagedBlockLabels(vfoxPathMarkerLabel)
	if len(labels) != 2 || labels[0] != vfoxPathMarkerLabel || labels[1] != legacyDataDirName+" PATH" {
		t.Fatalf("unexpected labels %v", labels)
	}
	// Labels the app does not own must not gain a phantom alias.
	if extra := unixManagedBlockLabels("Homebrew PATH"); len(extra) != 1 {
		t.Fatalf("unrelated label should stay single: %v", extra)
	}
}

func TestUnixRemoveAllManagedBlocksStripsLegacyBlock(t *testing.T) {
	profile := strings.Join([]string{
		"export KEEP=1",
		unixManagedPathBlock(legacyDataDirName+" PATH", []string{"/opt/vfox"}),
		"export AFTER=1",
	}, "\n")

	got := unixRemoveAllManagedBlocks(profile, vfoxPathMarkerLabel)
	if strings.Contains(got, "/opt/vfox") {
		t.Fatalf("block written before the rename must still be removable: %s", got)
	}
	if !strings.Contains(got, "export KEEP=1") || !strings.Contains(got, "export AFTER=1") {
		t.Fatalf("unrelated profile lines should be preserved: %s", got)
	}
}

func TestUnixSDKMarkerPrefixesCoverLegacyName(t *testing.T) {
	prefixes := unixSDKMarkerPrefixes()
	want := []string{"# >>> " + vfoxSDKPathMarkerPrefix, "# >>> " + legacyDataDirName + " SDK PATH "}
	if len(prefixes) != len(want) || prefixes[0] != want[0] || prefixes[1] != want[1] {
		t.Fatalf("unexpected prefixes %v", prefixes)
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote(filepath.Join("/tmp", "Bob's SDK", "bin"))
	if got != "'/tmp/Bob'\\''s SDK/bin'" {
		t.Fatalf("unexpected quoted path: %q", got)
	}
}
