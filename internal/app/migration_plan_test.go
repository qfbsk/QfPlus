package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPlanEnvironmentMigrationEmptyHome(t *testing.T) {
	a := NewApp()
	home := t.TempDir()
	a.setVfoxHome(home)

	plan, err := a.planEnvironmentMigration(filepath.Join(filepath.Dir(home), "newhome"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.BlockingReason != "" {
		t.Fatalf("expected no blocking reason for empty home, got %q", plan.BlockingReason)
	}
	if plan.TotalCount != 0 {
		t.Fatalf("expected 0 total count, got %d", plan.TotalCount)
	}
	if plan.TotalSizeBytes != 0 {
		t.Fatalf("expected 0 total size, got %d", plan.TotalSizeBytes)
	}
	if len(plan.MovableItems) != 0 {
		t.Fatalf("expected no movable items, got %d", len(plan.MovableItems))
	}
}

func TestPlanEnvironmentMigrationTargetWithinSourceBlocking(t *testing.T) {
	a := NewApp()
	home := t.TempDir()
	a.setVfoxHome(home)

	plan, err := a.planEnvironmentMigration(filepath.Join(home, "sub"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.BlockingReason == "" {
		t.Fatalf("expected blocking reason when target is inside source")
	}
}

func TestPlanEnvironmentMigrationExcludesExternalSdks(t *testing.T) {
	a := NewApp()
	home := t.TempDir()
	a.setVfoxHome(home)

	// A real, movable entry under the home.
	if err := os.MkdirAll(filepath.Join(home, "cache", "nodejs"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "cache", "nodejs", "v20.zip"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	// A third-party SDK recorded only by external path — must be listed, not moved.
	registry := `{"python":[{"name":"python","source":"system","path":"C:/python/python.exe","versions":[{"version":"3.12"}]}]}`
	if err := os.WriteFile(filepath.Join(home, "gui-non-vfox-sdks.json"), []byte(registry), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := a.planEnvironmentMigration(filepath.Join(t.TempDir(), "newhome"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.ExcludedItems) != 1 {
		t.Fatalf("expected 1 excluded item, got %d", len(plan.ExcludedItems))
	}
	excluded := plan.ExcludedItems[0]
	if excluded.Name != "python" {
		t.Fatalf("expected excluded name python, got %q", excluded.Name)
	}
	if excluded.WillMove {
		t.Fatalf("external SDK must not be marked willMove")
	}
	if plan.TotalCount == 0 {
		t.Fatalf("expected movable cache entry to be counted")
	}
}

func TestScanTreeSizeDoesNotFollowJunctions(t *testing.T) {
	root := t.TempDir()
	// A deep directory tree that a junction might point at.
	deep := filepath.Join(root, "deep")
	if err := os.MkdirAll(filepath.Join(deep, "a", "b", "c"), 0755); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(filepath.Join(deep, "a", "b", "c", fmt.Sprintf("f%d", i)), []byte("xxxxxxxx"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	linkPath := filepath.Join(root, "link")
	if err := os.Symlink(deep, linkPath); err != nil {
		t.Skipf("symlink not supported in this environment: %v", err)
	}

	count, _, err := scanTreeSize(linkPath)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("junction/link must count as a single item, got %d", count)
	}
}

func TestScanTreeSizeRegularTree(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "top.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "inner.txt"), []byte("world!!"), 0644); err != nil {
		t.Fatal(err)
	}

	count, size, err := scanTreeSize(root)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 { // root + top.txt + sub + sub/inner.txt
		t.Fatalf("expected 4 items, got %d", count)
	}
	if size != int64(len("hello")+len("world!!")) {
		t.Fatalf("expected size %d, got %d", len("hello")+len("world!!"), size)
	}
}
