package app

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func (a *App) proxyDir() string {
	return dataPath("proxy")
}

func (a *App) proxyConfigFile() string {
	return filepath.Join(a.proxyDir(), "config.yaml")
}

// subscriptionUserAgents lists client identifiers in probe order. Many
// subscription providers only deliver real nodes to Clash Meta family
// clients (older Clash Premium cannot handle vless/reality), so the Meta
// UA must be tried first.
var subscriptionUserAgents = []string{
	"clash.meta",
	"clash-verge/v2.4.2",
	"ClashforWindows/0.20.39",
}

func (a *App) fetchSubscription(url string) ([]byte, error) {
	var best []byte
	bestScore := -1
	var lastErr error
	for _, ua := range subscriptionUserAgents {
		data, err := a.fetchSubscriptionWithUA(url, ua)
		if err != nil {
			lastErr = err
			continue
		}
		score := scoreSubscription(data)
		if score > bestScore {
			best, bestScore = data, score
		}
		if bestScore >= 10 {
			break
		}
	}
	if best == nil {
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, fmt.Errorf("subscription server returned no usable content")
	}
	if bestScore <= 0 {
		return nil, fmt.Errorf("subscription did not provide usable nodes; the provider may restrict which clients can use it")
	}
	return best, nil
}

func (a *App) fetchSubscriptionWithUA(url, userAgent string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid subscription URL: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/yaml, application/yaml, */*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("subscription server returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to read subscription body: %w", err)
	}
	return data, nil
}

// scoreSubscription counts usable (non-loopback) proxies so the best
// UA response can be picked. Placeholder nodes pointing at 127.0.0.1
// ("your client is not supported" style) do not count.
func scoreSubscription(data []byte) int {
	config, err := parseSubscriptionConfig(data)
	if err != nil {
		return -1
	}
	proxies, _ := config["proxies"].([]any)
	_, hasProvider := config["proxy-providers"]
	score := 0
	for _, p := range proxies {
		pm, ok := p.(map[string]any)
		if !ok {
			continue
		}
		server, _ := pm["server"].(string)
		if isLoopbackHost(server) {
			continue
		}
		score++
	}
	if score == 0 && hasProvider {
		score = 1
	}
	return score
}

func parseSubscriptionConfig(data []byte) (map[string]any, error) {
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err == nil && config != nil {
		return config, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("not valid YAML")
	}
	if err := yaml.Unmarshal(decoded, &config); err != nil || config == nil {
		return nil, fmt.Errorf("not valid YAML after base64 decode")
	}
	return config, nil
}

func subscriptionHasGroup(data []byte, group string) bool {
	config, err := parseSubscriptionConfig(data)
	if err != nil {
		return false
	}
	rawGroups, _ := config["proxy-groups"].([]any)
	for _, g := range rawGroups {
		if gm, ok := g.(map[string]any); ok {
			if name, ok := gm["name"].(string); ok && name == group {
				return true
			}
		}
	}
	return false
}

func isLoopbackHost(host string) bool {
	switch strings.ToLower(strings.TrimSpace(host)) {
	case "localhost", "127.0.0.1", "0.0.0.0", "::1", "[::1]":
		return true
	}
	return false
}

func (a *App) patchSubscriptionYAML(raw []byte, mixedPort, apiPort int, secret, selectedGroup string) ([]byte, error) {
	var config map[string]any
	if err := yaml.Unmarshal(raw, &config); err != nil {
		decoded, b64err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(raw)))
		if b64err != nil {
			return nil, fmt.Errorf("subscription is not valid YAML: %w", err)
		}
		if err := yaml.Unmarshal(decoded, &config); err != nil {
			return nil, fmt.Errorf("subscription is not valid YAML after base64 decode: %w", err)
		}
	}

	if config == nil {
		return nil, fmt.Errorf("subscription config is empty")
	}

	proxies, ok := config["proxies"].([]any)
	if !ok || len(proxies) == 0 {
		if _, hasProvider := config["proxy-providers"]; !hasProvider {
			return nil, fmt.Errorf("subscription has no proxies defined")
		}
	}

	removed := map[string]bool{}
	var keptProxies []any
	for _, p := range proxies {
		pm, isMap := p.(map[string]any)
		if isMap {
			server, _ := pm["server"].(string)
			if isLoopbackHost(server) {
				if name, _ := pm["name"].(string); name != "" {
					removed[name] = true
				}
				continue
			}
		}
		keptProxies = append(keptProxies, p)
	}
	if len(keptProxies) > 0 {
		config["proxies"] = keptProxies
	}
	if len(removed) > 0 {
		if groups, ok := config["proxy-groups"].([]any); ok {
			for _, g := range groups {
				gm, ok := g.(map[string]any)
				if !ok {
					continue
				}
				members, ok := gm["proxies"].([]any)
				if !ok {
					continue
				}
				var kept []any
				for _, m := range members {
					if name, isStr := m.(string); isStr && removed[name] {
						continue
					}
					kept = append(kept, m)
				}
				gm["proxies"] = kept
			}
		}
	}

	config["mixed-port"] = mixedPort
	config["allow-lan"] = false
	config["external-controller"] = fmt.Sprintf("127.0.0.1:%d", apiPort)
	config["secret"] = secret
	config["mode"] = "rule"
	config["log-level"] = "warning"
	delete(config, "external-ui")
	delete(config, "external-ui-url")
	delete(config, "tun")
	delete(config, "dns")
	delete(config, "listeners")

	groupName := selectedGroup
	if groupName == "" {
		groups, hasGroups := config["proxy-groups"].([]any)
		if hasGroups && len(groups) > 0 {
			if first, ok := groups[0].(map[string]any); ok {
				if name, ok := first["name"].(string); ok {
					groupName = name
				}
			}
		}
		if groupName == "" {
			groupName = "PROXY"
			var names []any
			if proxies != nil {
				for _, p := range proxies {
					if pm, ok := p.(map[string]any); ok {
						if name, ok := pm["name"].(string); ok {
							names = append(names, name)
						}
					}
				}
			}
			config["proxy-groups"] = []any{
				map[string]any{
					"name":    "PROXY",
					"type":    "select",
					"proxies": names,
				},
			}
		}
	}

	config["rules"] = []any{
		fmt.Sprintf("MATCH,%s", groupName),
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize patched config: %w", err)
	}
	return out, nil
}

func (a *App) importSubscription(url string) error {
	raw, err := a.fetchSubscription(url)
	if err != nil {
		return err
	}

	config, err := a.readAppConfig()
	if err != nil {
		return err
	}

	if config.Proxy.SelectedGroup != "" && !subscriptionHasGroup(raw, config.Proxy.SelectedGroup) {
		config.Proxy.SelectedGroup = ""
		config.Proxy.SelectedNode = ""
	}

	mixedPort := config.Proxy.MixedPort
	if mixedPort == 0 {
		mixedPort = 17890
	}
	apiPort := config.Proxy.APIPort
	if apiPort == 0 {
		apiPort = 19090
	}

	patched, err := a.patchSubscriptionYAML(raw, mixedPort, apiPort, config.Proxy.APISecret, config.Proxy.SelectedGroup)
	if err != nil {
		return err
	}

	dir := a.proxyDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create proxy config dir: %w", err)
	}
	if err := os.WriteFile(a.proxyConfigFile(), patched, 0644); err != nil {
		return fmt.Errorf("failed to write proxy config: %w", err)
	}

	config.Proxy.SubscriptionURL = url
	if config.Proxy.SelectedGroup == "" {
		config.Proxy.SelectedGroup = detectFirstGroupName(patched)
	}
	return a.saveAppConfig(config)
}

func detectFirstGroupName(patched []byte) string {
	var config map[string]any
	if err := yaml.Unmarshal(patched, &config); err != nil {
		return ""
	}
	groups, ok := config["proxy-groups"].([]any)
	if !ok || len(groups) == 0 {
		return ""
	}
	if first, ok := groups[0].(map[string]any); ok {
		if name, ok := first["name"].(string); ok {
			return name
		}
	}
	return ""
}

// updateRulesGroup rewrites the rules section of the generated mihomo
// config so that all traffic is routed through the given group.
func (a *App) updateRulesGroup(group string) error {
	data, err := os.ReadFile(a.proxyConfigFile())
	if err != nil {
		return err
	}
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	config["rules"] = []any{fmt.Sprintf("MATCH,%s", group)}
	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(a.proxyConfigFile(), out, 0644)
}
