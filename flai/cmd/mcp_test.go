package cmd

import (
	"strings"
	"testing"
)

// S-0114: flai serve names the agent it started with --agent, because the
// harness's own settings can put a FLAI_AGENT of their own into the MCP
// server it starts (I-0037). The flag wins; without it, the environment.
func TestMCPServesTheAgentNamedByTheFlagOverTheEnvironment(t *testing.T) {
	t.Setenv("FLAI_AGENT", "system-flow")
	if got := mcpAgent("agent-S-0114"); got != "agent-S-0114" {
		t.Errorf("--agent: %q", got)
	}
	if got := mcpAgent(""); got != "system-flow" {
		t.Errorf("FLAI_AGENT: %q", got)
	}
	t.Setenv("FLAI_AGENT", "")
	if got := mcpAgent(""); got != "agent" {
		t.Errorf("neither: %q", got)
	}
	c := newMCPCmd(&app{})
	if f := c.Flags().Lookup("agent"); f == nil || !strings.Contains(f.Usage, "FLAI_AGENT") {
		t.Errorf("flai mcp takes --agent: %+v", f)
	}
}
