package app

import "strings"

// pluginCommandAliases returns every command a plugin is expected to provide,
// cross-platform. The plugin name itself is always the first entry. This is the
// single source of truth shared by both Windows shim generation and the
// environment status probe, so the two can never drift apart.
func pluginCommandAliases(pluginName string) []string {
	aliases := []string{pluginName}
	switch strings.ToLower(pluginName) {
	case "python":
		aliases = append(aliases, "python3", "pip", "pip3")
	case "nodejs":
		aliases = append(aliases, "npm", "npx", "corepack")
	case "java":
		aliases = append(aliases, "javac", "jar", "jshell", "jlink", "jpackage", "keytool")
	case "golang":
		aliases = append(aliases, "gofmt")
	case "rust":
		aliases = append(aliases, "cargo", "rustdoc", "rustup")
	case "ruby":
		aliases = append(aliases, "gem", "bundle", "bundler", "irb", "rake")
	case "php":
		aliases = append(aliases, "composer")
	case "perl":
		aliases = append(aliases, "cpan")
	}
	return uniqueStrings(aliases)
}

// uniqueStrings de-duplicates, case-insensitively, preserving first-seen order.
func uniqueStrings(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, value)
	}
	return result
}

// versionArgsFor returns the version-probe arguments for a command alias.
func versionArgsFor(alias string) []string {
	switch strings.ToLower(alias) {
	case "java", "javac", "jar", "jshell", "jlink", "jpackage", "keytool":
		return []string{"-version"}
	default:
		return []string{"--version"}
	}
}

// aliasCommandHint explains, for a plugin whose name is not itself a command,
// which command the user should actually type.
func aliasCommandHint(pluginName string) string {
	switch strings.ToLower(pluginName) {
	case "golang":
		return "plugin name is not a command: type `go` or `gofmt`"
	case "rust":
		return "plugin name is not a command: type `cargo`, `rustc` or `rustup`"
	default:
		return "no executable found on PATH for this command"
	}
}
