package serve

import (
	"context"
	"time"
)

// MCP starts and stops the HTTP MCP server each served project's agents
// reach (S-0096, ADR-0034): one `flai mcp http` per project, as ADR-0030
// has it, but kept running by flai serve rather than started by hand.
// Package cmd implements it, since the server and its state files are cmd's.
type MCP interface {
	// Ensure makes sure the project at root has one running. It starts one
	// unless one already runs, and says its PID and whether this call
	// started it: one already running (flai mcp start) is not flai serve's.
	Ensure(root string) (pid int, started bool, err error)
	// Stop ends one flai serve started and waits, briefly, for it to go.
	Stop(pid int)
}

// mcpEvery is how often a served project's MCP server is looked at, and
// started again when it has stopped, when Options.MCPEvery is zero.
const mcpEvery = 15 * time.Second

// superviseMCP keeps a project's MCP server running for as long as ctx lasts,
// then stops the one it started, if it is still the one serving. A server
// someone else started is used and left alone, and one of ours that stopped
// is started again on the next look. The child also watches this process
// (flai mcp http --exit-with), so it does not outlive a flai serve that was
// killed and never got here.
func superviseMCP(ctx context.Context, o Options, e Entry) {
	every := o.MCPEvery
	if every <= 0 {
		every = mcpEvery
	}
	owned, warned := 0, ""
	for {
		pid, started, err := o.MCP.Ensure(e.Root)
		switch {
		case err != nil:
			// the same failure every look (a port taken) is said once
			if err.Error() != warned {
				warned = err.Error()
				o.Logger.Warn("mcp server not started", "component", "serve", "root", e.Root, "err", warned)
			}
		case started:
			owned, warned = pid, ""
			o.Logger.Info("mcp server started", "component", "serve", "root", e.Root, "pid", pid)
		case pid != owned:
			owned = 0 // another process serves it now: not ours to stop
		}
		select {
		case <-ctx.Done():
			if owned != 0 && Alive(owned) {
				o.MCP.Stop(owned)
				o.Logger.Info("mcp server stopped", "component", "serve", "root", e.Root, "pid", owned)
			}
			return
		case <-time.After(every):
		}
	}
}
