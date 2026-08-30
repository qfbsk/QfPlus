package app

import "strings"

// pluginCommandAliases returns every command a plugin is expected to provide,
// cross-platform. The primary executable comes from systemSDKDefs when the
// plugin is listed there, otherwise the plugin name itself is used. This keeps
// shim generation and the environment status probe aligned, especially for
// plugins whose package name differs from the real command (e.g. golang -> go).
func pluginCommandAliases(pluginName string) []string {
	primary := pluginName
	for _, def := range systemSDKDefs {
		if def.Name == pluginName {
			primary = def.Exe
			break
		}
	}
	aliases := []string{primary}
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
// An empty slice means the command has no reliable version flag (e.g. gofmt);
// callers should skip the version probe but still treat the command as resolvable.
func versionArgsFor(alias string) []string {
	switch strings.ToLower(alias) {
	case "java", "javac", "jshell", "keytool":
		// Standard JDK tools accept -version.
		return []string{"-version"}
	case "jar", "jlink", "jpackage":
		// These JDK tools reject -version and only accept --version.
		return []string{"--version"}
	case "go":
		return []string{"version"}
	case "gofmt":
		// gofmt has no version flag; its version is the same as go.
		return nil
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
