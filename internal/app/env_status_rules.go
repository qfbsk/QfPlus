package app

// sdkCommandFacts captures the inputs classifyState needs. Every field is a
// plain boolean/string so the function stays pure and table-testable.
type sdkCommandFacts struct {
	resolved      bool
	managedBy     string // non-empty when QfPlus owns a shim for this command
	onUserPath    bool
	onMachinePath bool
	broken        bool
}

// classifyState maps the observed facts of one command to its presentation
// state. Pure: no IO, no clock. State values mirror the plan:
// ok | managed | unmanaged | broken | missing.
func classifyState(f sdkCommandFacts) string {
	switch {
	case !f.resolved:
		if f.managedBy != "" {
			return "managed"
		}
		return "missing"
	case f.broken:
		return "broken"
	case f.managedBy != "":
		return "managed"
	case f.onUserPath || f.onMachinePath:
		return "ok"
	default:
		return "unmanaged"
	}
}
