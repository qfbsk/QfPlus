package app

import (
	"reflect"
	"testing"
)

func TestPluginCommandAliasesUsesCanonicalExecutable(t *testing.T) {
	cases := []struct {
		plugin   string
		expected []string
	}{
		{"golang", []string{"go", "gofmt"}},
		{"nodejs", []string{"node", "npm", "npx", "corepack"}},
		{"python", []string{"python", "python3", "pip", "pip3"}},
		{"java", []string{"java", "javac", "jar", "jshell", "jlink", "jpackage", "keytool"}},
		{"rust", []string{"rustc", "cargo", "rustdoc", "rustup"}},
		{"ruby", []string{"ruby", "gem", "bundle", "bundler", "irb", "rake"}},
		{"php", []string{"php", "composer"}},
		{"perl", []string{"perl", "cpan"}},
	}

	for _, c := range cases {
		t.Run(c.plugin, func(t *testing.T) {
			got := pluginCommandAliases(c.plugin)
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("pluginCommandAliases(%q) = %v, want %v", c.plugin, got, c.expected)
			}
		})
	}
}

func TestVersionArgsForSpecialCases(t *testing.T) {
	cases := []struct {
		alias    string
		expected []string
	}{
		{"go", []string{"version"}},
		{"gofmt", nil},
		{"node", []string{"--version"}},
		{"java", []string{"-version"}},
		{"javac", []string{"-version"}},
		{"jshell", []string{"-version"}},
		{"keytool", []string{"-version"}},
		{"jar", []string{"--version"}},
		{"jlink", []string{"--version"}},
		{"jpackage", []string{"--version"}},
	}

	for _, c := range cases {
		t.Run(c.alias, func(t *testing.T) {
			got := versionArgsFor(c.alias)
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("versionArgsFor(%q) = %v, want %v", c.alias, got, c.expected)
			}
		})
	}
}
