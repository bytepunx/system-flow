package serve

import (
	"slices"
	"time"
)

// MCP is how flai serve has each served project's HTTP MCP server kept
// (S-0096, ADR-0034): one `flai mcp http` per project, as ADR-0030 has it.
// Since S-0106 flai serve starts none itself: flai host does, and flai
// serve tells it which projects it serves. Package cmd implements it with
// the host's client.
type MCP interface {
	// Keep says every project whose server is to be kept; one left out is
	// stopped, if the host started it.
	Keep(roots []string) error
}

// mcpEvery is how often the host is told again, changed or not, when
// Options.MCPEvery is zero.
const mcpEvery = 15 * time.Second

// mcpTeller tells the host which projects are served: at once when that
// changes, and again every so often, so that a host that lost it hears it.
type mcpTeller struct {
	o      Options
	told   []string
	at     time.Time
	warned string
}

func (m *mcpTeller) tell(roots []string) {
	if m.o.MCP == nil {
		return
	}
	every := m.o.MCPEvery
	if every <= 0 {
		every = mcpEvery
	}
	slices.Sort(roots)
	now := m.o.Now()
	if m.told != nil && slices.Equal(roots, m.told) && now.Sub(m.at) < every {
		return
	}
	if err := m.o.MCP.Keep(roots); err != nil {
		// the same failure every look (the host restarting) is said once
		if err.Error() != m.warned {
			m.warned = err.Error()
			m.o.Logger.Warn("mcp servers not kept", "component", "serve", "err", m.warned)
		}
		return
	}
	if m.warned != "" || !slices.Equal(roots, m.told) {
		m.o.Logger.Info("mcp servers kept", "component", "serve", "projects", len(roots))
	}
	m.told, m.at, m.warned = roots, now, ""
}
