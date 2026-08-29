package app

import (
	"fmt"
	"os"
)

// GetProxyStatus returns the current proxy state.
func (a *App) GetProxyStatus() (ProxyStatus, error) {
	config, err := a.readAppConfig()
	if err != nil {
		return ProxyStatus{}, err
	}
	_, fileErr := os.Stat(a.proxyConfigFile())
	return ProxyStatus{
		Enabled:       config.Proxy.Enabled,
		Running:       a.isMihomoRunning(),
		HasConfig:     fileErr == nil,
		MixedPort:     config.Proxy.MixedPort,
		Subscription:  config.Proxy.SubscriptionURL,
		SelectedGroup: config.Proxy.SelectedGroup,
		SelectedNode:  config.Proxy.SelectedNode,
	}, nil
}

// SetProxyEnabled starts or stops the proxy.
func (a *App) SetProxyEnabled(enabled bool) (ProxyStatus, error) {
	config, err := a.readAppConfig()
	if err != nil {
		return ProxyStatus{}, err
	}

	config.Proxy.Enabled = enabled
	if err := a.saveAppConfig(config); err != nil {
		return ProxyStatus{}, err
	}

	if enabled {
		if _, err := a.ensureAPISecret(); err != nil {
			return a.getProxyStatusWithError(err)
		}
		if err := a.startMihomo(); err != nil {
			return a.getProxyStatusWithError(err)
		}
		if config.Proxy.SelectedGroup != "" && config.Proxy.SelectedNode != "" {
			_ = a.selectProxyNode(config.Proxy.SelectedGroup, config.Proxy.SelectedNode)
		}
	} else {
		a.stopMihomo()
	}

	a.emitEvent("proxy-status-changed")
	return a.GetProxyStatus()
}

// ImportProxySubscription fetches and applies a subscription link.
func (a *App) ImportProxySubscription(url string) (ProxyStatus, error) {
	if url == "" {
		return ProxyStatus{}, fmt.Errorf("subscription URL cannot be empty")
	}

	wasRunning := a.isMihomoRunning()
	if wasRunning {
		a.stopMihomo()
	}

	if err := a.importSubscription(url); err != nil {
		return ProxyStatus{}, err
	}

	config, _ := a.readAppConfig()
	if config.Proxy.Enabled {
		if err := a.startMihomo(); err != nil {
			return a.getProxyStatusWithError(err)
		}
		if config.Proxy.SelectedGroup != "" && config.Proxy.SelectedNode != "" {
			_ = a.selectProxyNode(config.Proxy.SelectedGroup, config.Proxy.SelectedNode)
		}
	}

	a.emitEvent("proxy-status-changed")
	return a.GetProxyStatus()
}

// GetProxyGroups returns all proxy groups with their nodes.
func (a *App) GetProxyGroups() ([]ProxyGroup, error) {
	if !a.isMihomoRunning() {
		return nil, fmt.Errorf("proxy is not running")
	}
	return a.getProxyGroups()
}

// SelectProxyNode switches the active node in a group.
func (a *App) SelectProxyNode(group string, node string) (ProxyStatus, error) {
	if err := a.selectProxyNode(group, node); err != nil {
		return ProxyStatus{}, err
	}

	config, err := a.readAppConfig()
	if err == nil {
		config.Proxy.SelectedGroup = group
		config.Proxy.SelectedNode = node
		_ = a.saveAppConfig(config)
	}

	a.emitEvent("proxy-status-changed")
	return a.GetProxyStatus()
}

// SetProxyGroup switches which proxy group carries all QfPlus traffic by
// rewriting the MATCH rule and hot-reloading mihomo.
func (a *App) SetProxyGroup(group string) (ProxyStatus, error) {
	if group == "" {
		return ProxyStatus{}, fmt.Errorf("group name cannot be empty")
	}
	if !a.isMihomoRunning() {
		return ProxyStatus{}, fmt.Errorf("proxy is not running")
	}

	if err := a.updateRulesGroup(group); err != nil {
		return ProxyStatus{}, err
	}
	if err := a.reloadMihomoConfig(); err != nil {
		a.stopMihomo()
		if startErr := a.startMihomo(); startErr != nil {
			return a.getProxyStatusWithError(startErr)
		}
	}

	config, err := a.readAppConfig()
	if err == nil {
		config.Proxy.SelectedGroup = group
		config.Proxy.SelectedNode = ""
		_ = a.saveAppConfig(config)
	}

	a.emitEvent("proxy-status-changed")
	return a.GetProxyStatus()
}

// TestProxyGroupDelay tests latency for all nodes in a group.
func (a *App) TestProxyGroupDelay(group string) ([]ProxyGroup, error) {
	delays, err := a.testGroupDelay(group)
	if err != nil {
		return nil, err
	}

	groups, err := a.getProxyGroups()
	if err != nil {
		return nil, err
	}
	groupTypes := map[string]bool{
		"Selector": true, "URLTest": true, "Fallback": true, "LoadBalance": true,
	}
	builtinTypes := map[string]bool{
		"Direct": true, "Reject": true, "Compatible": true, "Pass": true, "Unknown": true,
	}
	for i := range groups {
		for j := range groups[i].Nodes {
			node := &groups[i].Nodes[j]
			if d, ok := delays[node.Name]; ok {
				node.Delay = d
				continue
			}
			if !groupTypes[node.Type] && !builtinTypes[node.Type] {
				node.Delay = -1
			}
		}
	}
	return groups, nil
}

// TestProxyNodeDelay tests latency for a single node.
func (a *App) TestProxyNodeDelay(name string) (ProxyNode, error) {
	delay, err := a.testNodeDelay(name)
	if err != nil {
		return ProxyNode{Name: name, Delay: -1}, err
	}
	return ProxyNode{Name: name, Delay: delay}, nil
}

// GetProxyQuickStatus returns the compact state used by the home page toggle.
func (a *App) GetProxyQuickStatus() ProxyQuickStatus {
	return a.quickStatus(false)
}

// TestProxyQuickDelay resolves the live exit node and measures its latency.
func (a *App) TestProxyQuickDelay() ProxyQuickStatus {
	return a.quickStatus(true)
}

// SetProxyQuickEnabled flips the proxy on or off from the home page shortcut.
func (a *App) SetProxyQuickEnabled(enabled bool) ProxyQuickStatus {
	if _, err := a.SetProxyEnabled(enabled); err != nil {
		status := a.quickStatus(false)
		status.Error = err.Error()
		return status
	}
	return a.quickStatus(enabled)
}

func (a *App) quickStatus(withDelay bool) ProxyQuickStatus {
	config, _ := a.readAppConfig()
	_, fileErr := os.Stat(a.proxyConfigFile())
	status := ProxyQuickStatus{
		Enabled:   config.Proxy.Enabled,
		Running:   a.isMihomoRunning(),
		HasConfig: fileErr == nil,
	}
	if !status.Running {
		return status
	}

	group, node, err := a.resolveExitNode(config.Proxy.SelectedGroup)
	if err != nil {
		status.Error = err.Error()
		return status
	}
	status.ExitGroup = group
	status.ExitNode = node
	if !withDelay {
		return status
	}
	if delay, err := a.testNodeDelay(node); err == nil {
		status.Delay = delay
	} else {
		status.Delay = -1
	}
	return status
}

func (a *App) getProxyStatusWithError(err error) (ProxyStatus, error) {
	status, _ := a.GetProxyStatus()
	status.Error = err.Error()
	return status, nil
}

func (a *App) restoreProxyState() {
	config, err := a.readAppConfig()
	if err != nil || !config.Proxy.Enabled {
		return
	}
	if _, err := os.Stat(a.proxyConfigFile()); err != nil {
		return
	}
	if err := a.startMihomo(); err != nil {
		a.emitEvent("vfox-log", "[PROXY] Failed to restore proxy: "+err.Error())
		return
	}
	if config.Proxy.SelectedGroup != "" && config.Proxy.SelectedNode != "" {
		_ = a.selectProxyNode(config.Proxy.SelectedGroup, config.Proxy.SelectedNode)
	}
	a.emitEvent("proxy-status-changed")
}
