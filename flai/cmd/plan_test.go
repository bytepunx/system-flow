//go:build !windows

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// S-0208: flai plan starts the planner for an epic or a story while the plan
// action is on, records the run by item where flai serve tracks agents, and
// says why when it will not; plan.run does the same for a dashboard.
func TestPlanStartsThePlannerForAnItem(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	inProcess(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	runIn(t, root, "task", "new", "Piece", "--story", "S-0001")
	if _, errOut, code := runIn(t, root, "plan", "E-1"); code == 0 || !strings.Contains(errOut, "rule: the plan host action is off for this project: flai serve enable plan") {
		t.Errorf("action off: %d %s", code, errOut)
	}
	out, _, code := runIn(t, root, "hostapi", "plan.run", `{"id":"S-0001","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a30"}`)
	if code == 0 || !strings.Contains(out, `the host action \"plan\" is not enabled`) || !strings.Contains(out, "flai serve enable plan") {
		t.Errorf("plan.run, action off: %d %s", code, out)
	}
	if out, errOut, code := runIn(t, root, "serve", "enable", "plan"); code != 0 || !strings.Contains(out, "plan enabled for t\n  it lets a dashboard start the planner") {
		t.Fatalf("enable: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "plan", "T-1"); code == 0 || !strings.Contains(errOut, "rule: T-0001 is a task; the planner plans an epic or a story") {
		t.Errorf("a task: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "plan", "E-1"); code == 0 || !strings.Contains(errOut, "rule: the planner's agent names no harness") {
		t.Errorf("nothing to start it with: %d %s", code, errOut)
	}
	runIn(t, root, "serve", "agent", "set", "--", "true")
	out, errOut, code := runIn(t, root, "plan", "E-1")
	if code != 0 || !strings.HasPrefix(out, "started true (command) to plan E-0001 as planner-E-0001 (pid ") || !strings.Contains(out, "-planner-") {
		t.Fatalf("plan: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "flai serve is not running") {
		t.Errorf("not told that no flai serve settles the run: %s", errOut)
	}
	out, _, code = runIn(t, root, "hostapi", "plan.run", `{"id":"S-0001","request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a31"}`)
	var started struct {
		Data struct {
			Item, Agent, Log, Session string
			PID                       int
		}
	}
	if err := json.Unmarshal([]byte(out), &started); code != 0 || err != nil || started.Data.Item != "S-0001" || started.Data.Agent != "planner-S-0001" || started.Data.PID == 0 || started.Data.Log == "" || started.Data.Session == "" {
		t.Fatalf("plan.run: %d %s", code, out)
	}
	var all map[string]serve.AgentState
	data, _ := os.ReadFile(filepath.Join(string(serve.DirFor(cfg)), "agents.json"))
	if err := json.Unmarshal(data, &all); err != nil {
		t.Fatal(err)
	}
	for _, st := range all {
		if st.Plans["E-0001"] == nil || st.Plans["S-0001"] == nil || len(st.Stories) != 0 {
			t.Errorf("state = %+v, want both runs by item, and no story's", st)
		}
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"outcome": "disabled"`, "started true (command) to plan E-0001 as planner-E-0001", "started true to plan S-0001 as planner-S-0001"} {
		if !strings.Contains(js, want) {
			t.Errorf("the journal has no %q: %s", want, js)
		}
	}
}

// S-0208: flai mcp's plan starts the planner as flai plan does, under the
// same plan host action, and journals the agent that asked.
func TestTheMCPServerStartsThePlannerUnderThePlanAction(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_STARTED_BY", "") // the operator's own agent, whoever runs the tests
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Epic")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 3, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	if _, err := a.mcpPlan(ctx, root, "E-0001", "agent-ops"); err == nil || !strings.Contains(err.Error(), "the plan host action is off for this project: flai serve enable plan") {
		t.Errorf("action off: %v", err)
	}
	runIn(t, root, "serve", "enable", "plan")
	runIn(t, root, "serve", "agent", "set", "--", "true")
	if _, err := a.mcpPlan(ctx, root, "T-0001", "agent-ops"); err == nil || !strings.Contains(err.Error(), "T-0001 is a task") {
		t.Errorf("a task: %v", err)
	}
	got, err := a.mcpPlan(ctx, root, "S-0001", "agent-ops")
	if err != nil || got.Item != "S-0001" || got.Agent != "planner-S-0001" || got.PID == 0 || got.Command != "true" || got.Log == "" || got.Session == "" {
		t.Fatalf("a story: %+v %v", got, err)
	}
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	if _, err := a.mcpPlan(ctx, root, "E-0001", "agent-S-0001"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start the planner") {
		t.Errorf("an agent flai serve started: %v", err)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	for _, want := range []string{`"method": "mcp.plan"`, `"by": "agent-ops"`, `"outcome": "disabled"`, "agent-ops asked to plan S-0001: started true as planner-S-0001", "agent-ops asked to plan T-0001: T-0001 is a task"} {
		if !strings.Contains(js, want) {
			t.Errorf("journal lacks %q: %s", want, js)
		}
	}
}
