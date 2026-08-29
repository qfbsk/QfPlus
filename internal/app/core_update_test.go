package app

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCompareCoreVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.11", "1.0.10", 1},
		{"1.0.9", "1.0.10", -1},
		{"1.2.0", "1.2.0", 0},
		{"v1.3.0", "1.2.9", 1},
		{"1.10.0", "1.9.9", 1},
	}
	for _, c := range cases {
		if got := compareCoreVersions(c.a, c.b); got != c.want {
			t.Fatalf("compareCoreVersions(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCoreAssetNameMatchesGoreleaserTemplate(t *testing.T) {
	name, format, err := coreAssetName("1.0.11")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(name, "vfox_1.0.11_") {
		t.Fatalf("unexpected asset name %q", name)
	}
	if !strings.HasSuffix(name, "."+format) {
		t.Fatalf("asset %q must end with format %q", name, format)
	}
	switch runtime.GOOS {
	case "windows":
		if format != "zip" || !strings.Contains(name, "_windows_") {
			t.Fatalf("windows asset mismatch: %q", name)
		}
	case "darwin":
		if !strings.Contains(name, "_macos_") {
			t.Fatalf("macos asset mismatch: %q", name)
		}
	}
}

func TestPlainTextReleaseNotes(t *testing.T) {
	raw := `<h3>🚀 New Features</h3><ul><li>support &lt;foo&gt; sdk</li><li>second &amp; item</li></ul>`
	got := plainTextReleaseNotes(raw)
	// Real markup is stripped before entities decode, so a literal "<foo>" in
	// the output is correct while a surviving "<h3>" would not be.
	if strings.Contains(got, "<h3>") || strings.Contains(got, "<ul>") || strings.Contains(got, "<li>") {
		t.Fatalf("markup not stripped: %q", got)
	}
	if strings.Contains(got, "&lt;") || strings.Contains(got, "&amp;") {
		t.Fatalf("entities not decoded: %q", got)
	}
	if !strings.Contains(got, "- support <foo> sdk") {
		t.Fatalf("list formatting lost: %q", got)
	}
	if !strings.Contains(got, "second & item") {
		t.Fatalf("entity not decoded: %q", got)
	}
}

func TestParseAtomFeed(t *testing.T) {
	fixture := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <entry>
    <id>tag:github.com,2008:Repository/527644591/v1.0.11</id>
    <updated>2025-03-04T10:00:00Z</updated>
    <link rel="stylesheet" href="https://github.com/version-fox/vfox/releases/tag/v1.0.11"/>
    <title>v1.0.11</title>
    <content type="html">&lt;h3&gt;Bug Fixes&lt;/h3&gt;&lt;ul&gt;&lt;li&gt;fix: shim crash&lt;/li&gt;&lt;/ul&gt;</content>
  </entry>
  <entry>
    <id>tag:github.com,2008:Repository/527644591/v1.0.10</id>
    <updated>2025-01-02T10:00:00Z</updated>
    <link rel="stylesheet" href="https://github.com/version-fox/vfox/releases/tag/v1.0.10"/>
    <title>v1.0.10</title>
    <content type="html">&lt;p&gt;older&lt;/p&gt;</content>
  </entry>
</feed>`
	releases, err := parseCoreAtomFeed([]byte(fixture))
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 2 {
		t.Fatalf("want 2 releases, got %d", len(releases))
	}
	if releases[0].Version != "1.0.11" || releases[1].Version != "1.0.10" {
		t.Fatalf("unexpected order: %+v", releases)
	}
	if releases[0].Date != "2025-03-04" {
		t.Fatalf("unexpected date %q", releases[0].Date)
	}
	if !strings.Contains(releases[0].Notes, "fix: shim crash") {
		t.Fatalf("notes not parsed: %q", releases[0].Notes)
	}
	if releases[0].URL != "https://github.com/version-fox/vfox/releases/tag/v1.0.11" {
		t.Fatalf("unexpected url %q", releases[0].URL)
	}
}

func TestCoreVersionRegexOnVfoxOutput(t *testing.T) {
	match := coreVersionRegex.FindStringSubmatch("vfox version 1.0.11\n")
	if match == nil || match[1] != "1.0.11" {
		t.Fatalf("failed to parse version output, got %#v", match)
	}
}

func TestLocalCoreVersionDirLayout(t *testing.T) {
	dir := localCoreVersionDir("9.9.9")
	if filepath.Base(dir) != "9.9.9" {
		t.Fatalf("version dir must end with the version: %q", dir)
	}
	root := localCoreRoot()
	if !strings.Contains(filepath.ToSlash(root), "/QfPlus/core/") {
		t.Fatalf("unexpected local core root %q", root)
	}
	if _, err := os.Stat(filepath.Dir(root)); err != nil {
		t.Skipf("user config dir unavailable: %v", err)
	}
}

func TestInstallCoreBinaryActivatesVersion(t *testing.T) {
	sandboxCoreStore(t)

	source := filepath.Join(t.TempDir(), coreExecutableName())
	if err := os.WriteFile(source, []byte("fake-engine"), 0755); err != nil {
		t.Fatal(err)
	}
	app := NewApp()
	if err := app.installCoreBinary("9.9.9", source); err != nil {
		t.Fatal(err)
	}

	installed := filepath.Join(localCoreVersionDir("9.9.9"), coreExecutableName())
	if !coreFileExists(installed) {
		t.Fatalf("expected installed binary at %s", installed)
	}
	if markerCoreDir() != localCoreVersionDir("9.9.9") {
		t.Fatalf("marker should point at the new version, got %q", markerCoreDir())
	}
	if got := findCoreFile(coreExecutableName()); filepath.Clean(got) != filepath.Clean(installed) {
		t.Fatalf("engine lookup resolved %q, want %q", got, installed)
	}
	if !containsVersion(localCoreVersions(), "9.9.9") {
		t.Fatalf("localCoreVersions missing 9.9.9: %v", localCoreVersions())
	}

	if err := app.activateBundledCore(); err != nil {
		t.Fatal(err)
	}
	if markerCoreDir() != "" {
		t.Fatal("marker must be cleared after restoring bundled core")
	}
	if !coreFileExists(installed) {
		t.Fatal("restoring the bundled core must keep downloaded versions")
	}
}

func TestInstallCoreBinaryOverwritesExistingVersion(t *testing.T) {
	sandboxCoreStore(t)
	dir := localCoreVersionDir("8.8.8")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(dir, coreExecutableName())
	if err := os.WriteFile(stale, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}

	source := filepath.Join(t.TempDir(), coreExecutableName())
	if err := os.WriteFile(source, []byte("fresh"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := NewApp().installCoreBinary("8.8.8", source); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(stale)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "fresh" {
		t.Fatalf("expected overwrite, got %q", data)
	}
}

func containsVersion(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}

// sandboxCoreStore redirects the user config dir so core marker and version
// store tests never touch the real %AppData%\QfPlus.
func sandboxCoreStore(t *testing.T) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("APPDATA", root)
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("HOME", root)
}
