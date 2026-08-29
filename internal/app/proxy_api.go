package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func (a *App) mihomoBaseURL() string {
	config, _ := a.readAppConfig()
	port := config.Proxy.APIPort
	if port == 0 {
		port = 19090
	}
	return fmt.Sprintf("http://127.0.0.1:%d", port)
}

func (a *App) mihomoRequest(method, path string, body io.Reader) (*http.Response, error) {
	config, _ := a.readAppConfig()
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(method, a.mihomoBaseURL()+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+config.Proxy.APISecret)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func (a *App) getProxyGroups() ([]ProxyGroup, error) {
	resp, err := a.mihomoRequest("GET", "/proxies", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mihomo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mihomo returned status %d", resp.StatusCode)
	}

	var result struct {
		Proxies map[string]struct {
			Type string   `json:"type"`
			Now  string   `json:"now"`
			All  []string `json:"all"`
		} `json:"proxies"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse proxies response: %w", err)
	}

	groupTypes := map[string]bool{
		"Selector": true, "URLTest": true, "Fallback": true, "LoadBalance": true,
	}
	skipNames := map[string]bool{"GLOBAL": true, "DIRECT": true, "REJECT": true}
	builtinNames := map[string]bool{"GLOBAL": true, "DIRECT": true, "REJECT": true, "PASS": true, "COMPATIBLE": true}

	orderIndex := map[string]int{}
	for i, name := range a.proxyGroupOrder() {
		orderIndex[name] = i
	}

	var groups []ProxyGroup
	for name, proxy := range result.Proxies {
		if !groupTypes[proxy.Type] || skipNames[name] {
			continue
		}
		hasRealMember := false
		for _, memberName := range proxy.All {
			if !builtinNames[memberName] {
				hasRealMember = true
				break
			}
		}
		if !hasRealMember {
			continue
		}
		group := ProxyGroup{
			Name: name,
			Type: proxy.Type,
			Now:  proxy.Now,
		}
		for _, memberName := range proxy.All {
			if member, ok := result.Proxies[memberName]; ok {
				group.Nodes = append(group.Nodes, ProxyNode{
					Name: memberName,
					Type: member.Type,
				})
			} else {
				group.Nodes = append(group.Nodes, ProxyNode{
					Name: memberName,
					Type: "Unknown",
				})
			}
		}
		groups = append(groups, group)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		oi, oki := orderIndex[groups[i].Name]
		oj, okj := orderIndex[groups[j].Name]
		if oki && okj {
			return oi < oj
		}
		if oki != okj {
			return oki
		}
		return groups[i].Name < groups[j].Name
	})
	return groups, nil
}

// proxyGroupOrder returns group names in the order they are declared in
// the generated mihomo config, used to keep UI ordering stable.
func (a *App) proxyGroupOrder() []string {
	data, err := os.ReadFile(a.proxyConfigFile())
	if err != nil {
		return nil
	}
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil
	}
	rawGroups, _ := config["proxy-groups"].([]any)
	var order []string
	for _, g := range rawGroups {
		if gm, ok := g.(map[string]any); ok {
			if name, ok := gm["name"].(string); ok {
				order = append(order, name)
			}
		}
	}
	return order
}

func (a *App) selectProxyNode(group, node string) error {
	payload, err := json.Marshal(map[string]string{"name": node})
	if err != nil {
		return err
	}
	path := "/proxies/" + url.PathEscape(group)
	resp, err := a.mihomoRequest("PUT", path, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("failed to select node: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("mihomo returned status %d when selecting node", resp.StatusCode)
	}
	return nil
}

// reloadMihomoConfig asks the running mihomo process to reload the
// generated config file so rule/group changes apply without a restart.
func (a *App) reloadMihomoConfig() error {
	payload, err := json.Marshal(map[string]string{"path": filepath.ToSlash(a.proxyConfigFile())})
	if err != nil {
		return err
	}
	resp, err := a.mihomoRequest("PATCH", "/configs", strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("failed to reach mihomo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("mihomo failed to reload config (status %d): %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (a *App) testGroupDelay(group string) (map[string]int, error) {
	path := fmt.Sprintf("/group/%s/delay?url=%s&timeout=5000",
		url.PathEscape(group),
		url.QueryEscape("https://www.gstatic.com/generate_204"))
	resp, err := a.mihomoRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to test delay: %w", err)
	}
	defer resp.Body.Close()

	var delays map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&delays); err != nil {
		return nil, fmt.Errorf("failed to parse delay response: %w", err)
	}
	return delays, nil
}

func (a *App) testNodeDelay(name string) (int, error) {
	path := fmt.Sprintf("/proxies/%s/delay?url=%s&timeout=5000",
		url.PathEscape(name),
		url.QueryEscape("https://www.gstatic.com/generate_204"))
	resp, err := a.mihomoRequest("GET", path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to test delay: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Delay int `json:"delay"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return -1, nil
	}
	return result.Delay, nil
}

// matchRuleGroup returns the proxy group the generated MATCH rule sends all
// traffic to, which is the entry point of the exit chain.
func (a *App) matchRuleGroup() string {
	data, err := os.ReadFile(a.proxyConfigFile())
	if err != nil {
		return ""
	}
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return ""
	}
	rules, _ := config["rules"].([]any)
	for _, raw := range rules {
		rule, ok := raw.(string)
		if !ok {
			continue
		}
		parts := strings.Split(rule, ",")
		if strings.EqualFold(strings.TrimSpace(parts[0]), "MATCH") && len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}

// resolveExitNode walks selector/url-test groups until it reaches the concrete
// proxy that actually carries traffic.
func (a *App) resolveExitNode(startGroup string) (string, string, error) {
	resp, err := a.mihomoRequest("GET", "/proxies", nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to mihomo: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Proxies map[string]struct {
			Type string `json:"type"`
			Now  string `json:"now"`
		} `json:"proxies"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to parse proxies response: %w", err)
	}

	groupTypes := map[string]bool{
		"Selector": true, "URLTest": true, "Fallback": true, "LoadBalance": true, "Policy": true,
	}
	if startGroup == "" {
		startGroup = a.matchRuleGroup()
	}
	if startGroup == "" {
		return "", "", fmt.Errorf("未找到出口分组")
	}

	name := startGroup
	seen := map[string]bool{}
	for !seen[name] && len(seen) < 12 {
		seen[name] = true
		proxy, ok := result.Proxies[name]
		if !ok {
			break
		}
		if !groupTypes[proxy.Type] {
			return startGroup, name, nil
		}
		if proxy.Now == "" {
			break
		}
		name = proxy.Now
	}
	return startGroup, name, nil
}
