package model

type ProxyNode struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Delay int    `json:"delay"`
}

type ProxyGroup struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Now   string      `json:"now"`
	Nodes []ProxyNode `json:"nodes"`
}

type ProxyStatus struct {
	Enabled       bool   `json:"enabled"`
	Running       bool   `json:"running"`
	HasConfig     bool   `json:"hasConfig"`
	MixedPort     int    `json:"mixedPort"`
	Subscription  string `json:"subscriptionUrl"`
	SelectedGroup string `json:"selectedGroup"`
	SelectedNode  string `json:"selectedNode"`
	Error         string `json:"error,omitempty"`
}

// ProxyQuickStatus is the compact payload the home page toggle renders.
type ProxyQuickStatus struct {
	Enabled   bool   `json:"enabled"`
	Running   bool   `json:"running"`
	HasConfig bool   `json:"hasConfig"`
	ExitGroup string `json:"exitGroup"`
	ExitNode  string `json:"exitNode"`
	Delay     int    `json:"delay"`
	Error     string `json:"error,omitempty"`
}
