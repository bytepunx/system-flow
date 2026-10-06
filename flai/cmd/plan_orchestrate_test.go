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
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0218, S-0219: of the agents flai serve starts, the orchestrator may have
// the planner plan an epic flai plan --candidates lists, while the project
// gives it plan_backlog_epics, and nothing else, through flai mcp's plan and
// flai plan alike; its run records the orchestrator as what started it, and
// the journal names it.
func TestTheOrchestratorAsksForThePlannerOnACandidateEpic(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Empty")      // E-0001: a candidate
	runIn(t, root, "epic", "new", "Dropped")    // E-0002: cancelled
	runIn(t, root, "epic", "new", "Open story") // E-0003: a story not done
	runIn(t, root, "epic", "new", "Waiting")    // E-0004: left out, its planner asked
	runIn(t, root, "epic", "new", "Also empty") // E-0005: a candidate
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0003")
	runIn(t, root, "move", "E-0002", "cancelled", "--reason", "dropped")
	runIn(t, root, "thread", "new", "--on", "E-0004", "--by", "planner-E-0004", "Which outcome?", "Say which.")
	runIn(t, root, "serve", "enable", "plan")
	runIn(t, root, "serve", "agent", "set", "--", "true")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	agents := filepath.Join(string(serve.DirFor(cfg)), "agents.json")
	data, _ := json.Marshal(map[string]serve.AgentState{mainRootOf(repo): {Plans: map[string]*serve.AgentRun{
		"E-0004": {Item: "E-0004", Agent: "planner-E-0004", Started: "2026-09-15T19:00:00Z", Ended: "2026-09-15T20:00:00Z", Outcome: serve.OutcomeAsked, Thread: "TH-0001"},
	}}})
	if err := os.MkdirAll(filepath.Dir(agents), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agents, data, 0o644); err != nil {
		t.Fatal(err)
	}
	trigger := func(item string) string {
		t.Helper()
		var all map[string]serve.AgentState
		data, _ := os.ReadFile(agents)
		if err := json.Unmarshal(data, &all); err != nil {
			t.Fatal(err)
		}
		if run := all[mainRootOf(repo)].Plans[item]; run != nil {
			return run.Trigger
		}
		return ""
	}

	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	t.Setenv("FLAI_ROLE", "orchestrate")
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC) }}
	ctx := context.Background()
	if _, err := a.mcpPlan(ctx, root, "E-0001", "agent-o"); err == nil || !strings.Contains(err.Error(), "only with orchestration.permissions.plan_backlog_epics, which is off") {
		t.Errorf("permission off: %v", err)
	}
	manifest := filepath.Join(root, "system-flow.yaml")
	m, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, append(m, "orchestration:\n  permissions:\n    plan_backlog_epics: true\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	const alone = "the orchestrator asks for the planner on an epic flai plan --candidates lists alone, and "
	for item, want := range map[string]string{
		"S-0001": alone + "S-0001 is not an epic",
		"T-0001": alone + "T-0001 is not an epic",
		"E-0002": alone + "E-0002 is not one: it is cancelled",
		"E-0003": alone + "E-0003 is not one: it is neither in the backlog with no stories nor open with every story done or cancelled and one done",
		"E-0004": alone + "E-0004 is left out: its planner asked on TH-0001, which awaits the operator: Which outcome?",
		"E-0009": "E-0009 not found",
	} {
		if _, err := a.mcpPlan(ctx, root, item, "agent-o"); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v, want %q", item, err, want)
		}
	}
	got, err := a.mcpPlan(ctx, root, "E-0001", "agent-o")
	if err != nil || got.Item != "E-0001" || got.Agent != "planner-E-0001" || got.PID == 0 {
		t.Fatalf("a candidate: %+v %v", got, err)
	}
	if tr := trigger("E-0001"); tr != serve.TriggerOrchestrator {
		t.Errorf("the run's trigger = %q, want orchestrator", tr)
	}

	// flai plan in the orchestrator's shell is held alike
	if _, stderr, code := runIn(t, root, "plan", "E-0003"); code == 0 || !strings.Contains(stderr, "rule: "+alone+"E-0003 is not one") {
		t.Errorf("flai plan on an epic with an open story: exit %d %s", code, stderr)
	}
	if out, stderr, code := runIn(t, root, "plan", "E-0005"); code != 0 || !strings.Contains(out, "to plan E-0005 as planner-E-0005") {
		t.Fatalf("flai plan on a candidate: exit %d %s %s", code, out, stderr)
	}
	if tr := trigger("E-0005"); tr != serve.TriggerOrchestrator {
		t.Errorf("flai plan's run trigger = %q, want orchestrator", tr)
	}

	// another agent flai serve started is refused, and the operator's own
	// agent plans a story as asked, as before
	t.Setenv("FLAI_ROLE", "plan")
	if _, err := a.mcpPlan(ctx, root, "E-0003", "planner-S-0001"); err == nil || !strings.Contains(err.Error(), "an agent flai serve started does not start the planner") {
		t.Errorf("another agent flai serve started: %v", err)
	}
	t.Setenv("FLAI_ROLE", "")
	t.Setenv("FLAI_STARTED_BY", "")
	if got, err := a.mcpPlan(ctx, root, "S-0001", "agent-ops"); err != nil || got.Item != "S-0001" {
		t.Fatalf("the operator's agent on a story: %+v %v", got, err)
	}
	if tr := trigger("S-0001"); tr != "asked" {
		t.Errorf("the operator's agent's run trigger = %q, want asked", tr)
	}
	js, _, _ := runIn(t, root, "serve", "journal", "--json")
	if !strings.Contains(js, `"by": "orchestrator (agent-o)"`) || !strings.Contains(js, "orchestrator (agent-o) asked to plan E-0001: started true as planner-E-0001") {
		t.Errorf("the journal does not name the orchestrator: %s", js)
	}
}
