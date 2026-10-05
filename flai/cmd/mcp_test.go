package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// S-0257: flai mcp gives permission_prompt the operator's auto-approve host
// action for the project, read afresh at every request: enabling it in a
// shell turns it on and disabling it off without restarting the server, and
// a configuration flai cannot read enables nothing, so permission_prompt asks.
func TestMCPAutoApproveReadsTheShellsSettingAtEveryRequest(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfgPath)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	main := mainRootOf(repo)
	a := &app{out: &bytes.Buffer{}, errOut: &bytes.Buffer{}}
	if a.mcpAutoApprove(main) {
		t.Error("auto-approve is off by default")
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "auto-approve"); code != 0 || !strings.Contains(out, "auto-approve enabled for t\n  it has flai's MCP tool permission_prompt allow") || strings.Contains(out, "journal") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if !a.mcpAutoApprove(main) {
		t.Error("enabled in a shell, the same server now allows without asking")
	}
	if a.mcpAutoApprove(filepath.Join(t.TempDir(), "other")) {
		t.Error("enabled for this project only, another project still asks")
	}
	if _, errOut, code := runIn(t, root, "serve", "disable", "auto-approve"); code != 0 {
		t.Fatalf("disable: %s", errOut)
	}
	if a.mcpAutoApprove(main) {
		t.Error("disabled in a shell, the same server asks again")
	}
	if _, errOut, code := runIn(t, root, "serve", "enable", "auto-approve"); code != 0 {
		t.Fatalf("enable again: %s", errOut)
	}
	if err := os.WriteFile(cfgPath, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if a.mcpAutoApprove(main) {
		t.Error("a configuration flai cannot read enables nothing")
	}
}
