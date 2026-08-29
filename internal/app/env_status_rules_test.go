package app

import "testing"

func TestClassifyState(t *testing.T) {
	cases := []struct {
		name string
		f    sdkCommandFacts
		want string
	}{
		{"unresolved managed", sdkCommandFacts{resolved: false, managedBy: "nodejs"}, "managed"},
		{"unresolved missing", sdkCommandFacts{resolved: false}, "missing"},
		{"broken", sdkCommandFacts{resolved: true, broken: true}, "broken"},
		{"resolved managed", sdkCommandFacts{resolved: true, managedBy: "nodejs", onUserPath: true}, "managed"},
		{"resolved ok on user", sdkCommandFacts{resolved: true, onUserPath: true}, "ok"},
		{"resolved ok on machine", sdkCommandFacts{resolved: true, onMachinePath: true}, "ok"},
		{"resolved unmanaged", sdkCommandFacts{resolved: true}, "unmanaged"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := classifyState(c.f); got != c.want {
				t.Errorf("classifyState(%+v) = %q, want %q", c.f, got, c.want)
			}
		})
	}
}
