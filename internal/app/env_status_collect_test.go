package app

import (
	"os"
	"strings"
	"testing"
)

func TestLivePathEnv(t *testing.T) {
	base := []string{"PATH=stale", "SystemRoot=C:\\Windows", "OTHER=1"}
	got := livePathEnv(base, "C:\\user", "C:\\machine")

	var pathVal string
	found := false
	for _, kv := range got {
		if strings.HasPrefix(strings.ToUpper(kv), "PATH=") {
			pathVal = kv[len("PATH="):]
			found = true
		}
	}
	if !found {
		t.Fatalf("PATH entry missing from result: %v", got)
	}
	want := "C:\\user" + string(os.PathListSeparator) + "C:\\machine"
	if pathVal != want {
		t.Errorf("live PATH = %q, want %q", pathVal, want)
	}
	// Other variables must be preserved.
	if !containsKV(got, "SystemRoot=C:\\Windows") || !containsKV(got, "OTHER=1") {
		t.Errorf("livePathEnv dropped unrelated env vars: %v", got)
	}
	// Original slice must not be mutated.
	if base[0] != "PATH=stale" {
		t.Errorf("livePathEnv mutated the input slice: %v", base)
	}
}

func TestLivePathEnvEmptyScopes(t *testing.T) {
	base := []string{"PATH=stale"}
	got := livePathEnv(base, "", "")
	for _, kv := range got {
		if strings.HasPrefix(strings.ToUpper(kv), "PATH=") && kv != "PATH=stale" {
			t.Errorf("empty scopes should leave PATH unchanged, got %q", kv)
		}
	}
}

func TestCombinePathScopes(t *testing.T) {
	sep := string(os.PathListSeparator)
	if got := combinePathScopes("A", "B"); got != "A"+sep+"B" {
		t.Errorf("combinePathScopes(A,B) = %q, want %q", got, "A"+sep+"B")
	}
	if got := combinePathScopes("", "B"); got != "B" {
		t.Errorf("combinePathScopes('',B) = %q, want B", got)
	}
	if got := combinePathScopes("A", ""); got != "A" {
		t.Errorf("combinePathScopes(A,'') = %q, want A", got)
	}
}

func containsKV(env []string, kv string) bool {
	for _, e := range env {
		if e == kv {
			return true
		}
	}
	return false
}
