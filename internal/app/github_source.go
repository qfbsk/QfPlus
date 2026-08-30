package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// presetGitHubSources are well-known GitHub mirrors / accelerators popular in
// mainland China. The URL may contain a "{url}" placeholder, act as a prefix,
// or replace the github.com hostname, depending on the mirror's convention.
var presetGitHubSources = []GitHubSource{
	{ID: "github", Name: "GitHub 官方", URL: "https://github.com", IsPreset: true},
	{ID: "ghproxy", Name: "GHProxy", URL: "https://mirror.ghproxy.com/{url}", IsPreset: true},
	{ID: "ghps", Name: "GHPS", URL: "https://ghps.cc/https://github.com", IsPreset: true},
	{ID: "kkgithub", Name: "KK GitHub", URL: "https://kkgithub.com", IsPreset: true},
	{ID: "gitclone", Name: "GitClone", URL: "https://gitclone.com/github.com", IsPreset: true},
	{ID: "moeyy", Name: "Moeyy", URL: "https://github.moeyy.xyz/{url}", IsPreset: true},
}

// githubSourceByID returns a source from the merged preset + custom list.
func (a *App) githubSourceByID(id string) (GitHubSource, bool) {
	for _, s := range presetGitHubSources {
		if s.ID == id {
			return s, true
		}
	}
	config, _ := a.readAppConfig()
	for _, s := range config.GitHubSource.CustomSources {
		if s.ID == id {
			return s, true
		}
	}
	return GitHubSource{}, false
}

// selectedGitHubSource returns the currently active source, or the official
// GitHub endpoint when none is selected or the feature is disabled.
func (a *App) selectedGitHubSource() GitHubSource {
	config, err := a.readAppConfig()
	if err != nil || !config.GitHubSource.Enabled {
		return presetGitHubSources[0]
	}
	id := config.GitHubSource.SelectedID
	if id == "" {
		return presetGitHubSources[0]
	}
	if s, ok := a.githubSourceByID(id); ok {
		return s
	}
	return presetGitHubSources[0]
}

// transformGitHubURL rewrites a github.com URL through the active mirror.
// The source URL may use one of three forms:
//   - template:  "https://mirror.ghproxy.com/{url}" -> prefix the full URL
//   - prefix:    "https://ghps.cc/https://github.com" -> concatenate
//   - hostname:  "https://kkgithub.com" -> replace github.com hostname
func transformGitHubURL(sourceURL, originalURL string) string {
	if sourceURL == "" || originalURL == "" {
		return originalURL
	}

	// Template form.
	if strings.Contains(sourceURL, "{url}") {
		return strings.ReplaceAll(sourceURL, "{url}", originalURL)
	}

	// Prefix form: source URL starts with http and the original URL is absolute.
	if strings.HasPrefix(sourceURL, "http") && strings.HasPrefix(originalURL, "http") {
		return strings.TrimSuffix(sourceURL, "/") + "/" + strings.TrimPrefix(originalURL, "/")
	}

	// Hostname replacement form.
	u, err := url.Parse(originalURL)
	if err != nil {
		return originalURL
	}
	host := strings.TrimPrefix(strings.TrimPrefix(sourceURL, "https://"), "http://")
	host = strings.Split(host, "/")[0]
	if host == "" {
		return originalURL
	}
	u.Scheme = "https"
	u.Host = host
	return u.String()
}

// applyGitHubSource rewrites a github.com URL through the currently selected
// GitHub source when the feature is enabled.
func (a *App) applyGitHubSource(originalURL string) string {
	if !strings.Contains(originalURL, "github.com") {
		return originalURL
	}
	source := a.selectedGitHubSource()
	if source.ID == "github" || source.URL == "" {
		return originalURL
	}
	return transformGitHubURL(source.URL, originalURL)
}

// githubSourceEnvValue returns the URL of the active source, suitable for
// passing into child processes such as vfox plugins that may honor a mirror
// environment variable.
func (a *App) githubSourceEnvValue() string {
	source := a.selectedGitHubSource()
	if source.ID == "github" {
		return ""
	}
	return source.URL
}

// vfoxPluginsRegistryURL is the canonical raw GitHub URL for the vfox plugin
// registry. Many GitHub mirrors can accelerate this endpoint as well.
const vfoxPluginsRegistryURL = "https://raw.githubusercontent.com/version-fox/vfox-plugins/main/plugins"

// vfoxRegistryMirrorURL maps the selected GitHub source to a vfox plugin
// registry mirror address. It returns "" when the official GitHub endpoint is
// selected or when no safe mapping is known for the custom source.
func vfoxRegistryMirrorURL(source GitHubSource) string {
	switch source.ID {
	case "github":
		return ""
	case "ghproxy":
		return transformGitHubURL(source.URL, vfoxPluginsRegistryURL)
	case "ghps":
		return transformGitHubURL(source.URL, vfoxPluginsRegistryURL)
	case "kkgithub":
		return transformGitHubURL(source.URL, vfoxPluginsRegistryURL)
	case "gitclone":
		return transformGitHubURL(source.URL, vfoxPluginsRegistryURL)
	case "moeyy":
		return transformGitHubURL(source.URL, vfoxPluginsRegistryURL)
	}
	// For custom sources, make a best-effort guess if the user pasted a
	// GitHub-mirror-style URL. Otherwise leave vfox at its default.
	if !source.IsPreset && strings.Contains(source.URL, "github.com") {
		return transformGitHubURL(source.URL, vfoxPluginsRegistryURL)
	}
	return ""
}

// syncVfoxGitHubSourceConfig pushes the selected GitHub mirror into vfox's
// plugin registry configuration. This accelerates plugin discovery and plugin
// manifest downloads, which are the parts of vfox that QfPlus can actually
// steer. SDK package downloads are controlled by individual vfox plugins and
// may still need the bundled mihomo proxy.
func (a *App) syncVfoxGitHubSourceConfig() error {
	source := a.selectedGitHubSource()
	registryURL := vfoxRegistryMirrorURL(source)

	if registryURL == "" {
		// Restore vfox's default registry.
		_, _ = a.runVfoxCommand("config", "--unset", "registry.address")
		return nil
	}

	_, err := a.runVfoxCommand("config", "registry.address", registryURL)
	return err
}

func generateSourceID() string {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("custom-%d", time.Now().Unix())
	}
	return "custom-" + hex.EncodeToString(buf)
}
