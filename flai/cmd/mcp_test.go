package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/serve"
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

// S-0534: flai mcp reads the host's budget state from budget.json in flai
// serve's folder beside FLAI_CONFIG's configuration, afresh at every inbox,
// without writing it. A missing serve folder or file is no budget, and so is
// one it cannot read, which it logs at warn once.
func TestMCPBudgetsReadsTheHostsStateBesideTheConfiguration(t *testing.T) {
	cfgDir := t.TempDir()
	t.Setenv("FLAI_CONFIG", filepath.Join(cfgDir, "config.json"))
	errOut := &bytes.Buffer{}
	a := &app{out: &bytes.Buffer{}, errOut: errOut}
	read := a.mcpBudgets()
	if got := read(); len(got) != 0 {
		t.Errorf("no serve folder, no budget: %+v", got)
	}
	dir := filepath.Join(cfgDir, "serve")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := read(); len(got) != 0 {
		t.Errorf("no budget.json, no budget: %+v", got)
	}
	east := time.FixedZone("east", 2*60*60)
	state := map[string]serve.Budget{
		"anthropic":   {Provider: "anthropic", Limit: "usage_limit_reached", Status: 429, Since: time.Date(2026, 10, 28, 12, 0, 0, 0, east), Until: time.Date(2026, 10, 28, 15, 0, 0, 0, east), Backoff: "5m0s"},
		"claude-code": {Provider: "claude-code", Since: time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC), Until: time.Date(2026, 10, 28, 9, 5, 0, 0, time.UTC), Backoff: "5m0s", Probe: &serve.BudgetProbe{Root: "/p", Work: "S-1", Started: "2026-10-28T09:05:00Z"}},
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, "budget.json")
	if err := os.WriteFile(file, data, 0o600); err != nil {
		t.Fatal(err)
	}
	want := []inbox.Budget{
		{Provider: "anthropic", Limit: "usage_limit_reached", Status: 429, Since: "2026-10-28T10:00:00Z", Until: "2026-10-28T13:00:00Z"},
		{Provider: "claude-code", Since: "2026-10-28T09:00:00Z", Until: "2026-10-28T09:05:00Z", Probing: true},
	}
	got := read()
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("the state set, each provider's, UTC:\n got %+v\nwant %+v", got, want)
	}
	after, err := os.ReadFile(file)
	if err != nil || !bytes.Equal(after, data) {
		t.Errorf("reading wrote budget.json: %s %v", after, err)
	}
	if err := os.WriteFile(file, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if got := read(); len(got) != 0 {
			t.Errorf("an unreadable file is no budget: %+v", got)
		}
	}
	if n := strings.Count(errOut.String(), "cannot read the host's budget state"); n != 1 {
		t.Errorf("an unreadable file is logged once, not %d times: %s", n, errOut)
	}
}
