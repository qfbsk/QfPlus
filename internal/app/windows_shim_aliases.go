//go:build windows

package app

import "strings"

func windowsSafeShimName(name string) string {
	name = strings.TrimSpace(name)
	replacer := strings.NewReplacer("\\", "_", "/", "_", ":", "_", "*", "_", "?", "_", `"`, "_", "<", "_", ">", "_", "|", "_")
	name = replacer.Replace(name)
	if name == "" {
		return "sdk"
	}
	return name
}

func windowsUniqueStrings(values []string) []string {
	seen := make(map[string]bool)
	var result []string
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

func windowsSDKShimAliases(pluginName string) []string {
	aliases := pluginCommandAliases(pluginName)
	for _, def := range systemSDKDefs {
		if def.Name == pluginName {
			// Place the canonical executable right after the plugin name so
			// filtering by real binaries preserves the expected order.
			aliases = append(aliases[:1], append([]string{def.Exe}, aliases[1:]...)...)
			break
		}
	}
	return uniqueStrings(aliases)
}
