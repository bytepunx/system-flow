package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0206, ADR-0079: flai mcp's activity_log measures a strategic agent's
// activity from its run's log in the serve folder and writes it in the
// project's activity document, returning the entry and the totals.
func TestTheMCPServerLogsAStrategicAgentsActivity(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	start := time.Date(2026, 10, 3, 10, 0, 0, 0, time.UTC)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: start})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story", Parent: epic.ID, Owner: "alex", Now: start}); err != nil {
		t.Fatal(err)
	}
	call := func(id string, at time.Time, read int) string {
		return fmt.Sprintf(`{"type":"assistant","timestamp":%q,"session_id":"s","message":{"id":%q,"model":"claude-opus-5-5","content":[],"usage":{"input_tokens":1,"output_tokens":1,"cache_read_input_tokens":%d,"cache_creation_input_tokens":0}}}`, at.Format(time.RFC3339Nano), id, read)
	}
	// calls weigh 1000, 1000, and 2000 of the session's 4000
	lines := []string{
		call("m1", start, 999), call("m2", start.Add(time.Minute), 999), call("m3", start.Add(2*time.Minute), 1999),
		`{"type":"result","subtype":"success","session_id":"s","modelUsage":{"claude-opus-5-5":{"inputTokens":2,"outputTokens":300,"cacheReadInputTokens":4000,"cacheCreationInputTokens":0,"costUSD":1.2345}}}`,
	}
	logs := filepath.Join(filepath.Dir(cfg), "serve", "agents")
	if err := os.MkdirAll(logs, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(logs, "t-orchestrator-20261003T100000Z.log"), []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	now := start.Add(90 * time.Second)
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, clock: func() time.Time { return now }}
	ctx := context.Background()
	half := math.Round(1.2345*2000/4000*1e4) / 1e4

	got, err := a.mcpActivity(ctx, root, workitem.ActivityOrchestrator, "Ordered the board", []string{"S-0001", "TH-0001"}, "orchestrator")
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.At != "2026-10-03T10:01:30Z" || e.Seconds != 90 || e.Cost != half || !e.Estimated || strings.Join(e.Items, ",") != "S-0001,TH-0001" {
		t.Errorf("first entry: %+v", e)
	}
	// ADR-0095: the orchestrator's cost is charged to the work item it named
	if strings.Join(got.ChargedTo, ",") != "S-0001" || got.Charge != "charged to S-0001 and the items above it" {
		t.Errorf("first charge: %v, %q, want S-0001's and the items above it", got.ChargedTo, got.Charge)
	}
	now = start.Add(3 * time.Minute) // past the run's last call, at 10:02:00
	got, err = a.mcpActivity(ctx, root, workitem.ActivityOrchestrator, "Planned nothing new", nil, "orchestrator")
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.Seconds != 31 || e.Cost != half {
		t.Errorf("second entry, from the first's end to the second after the run's last call: %+v", e)
	}
	if got.ChargedTo != nil || got.Charge != "charged to no item: left in the orchestrator's project strategic total" {
		t.Errorf("second charge, naming no item: %v, %q", got.ChargedTo, got.Charge)
	}
	want := mcpserver.ActivityTotals{Kind: "orchestrator", AccruedCost: math.Round(2*half*1e4) / 1e4, AccruedSeconds: 121, TasksCompleted: 2, LastRun: "2026-10-03T10:03:00Z"}
	if got.Activity != want {
		t.Errorf("totals: %+v, want %+v", got.Activity, want)
	}
	doc, err := workitem.ReadActivity(filepath.Join(root, "wip", "agents", "orchestrator.md"))
	if err != nil {
		t.Fatalf("the document in the main checkout: %v", err)
	}
	if len(doc.Entries) != 2 || doc.TasksCompleted != 2 || doc.Entries[0].Summary != "Ordered the board" {
		t.Errorf("document: %+v", doc)
	}
	if _, err := a.mcpActivity(ctx, root, "story", "Did it", nil, "orchestrator"); err == nil || !strings.Contains(err.Error(), "planner, orchestrator, analyzer") {
		t.Errorf("a kind that is no strategic agent: %v", err)
	}
}

// ADR-0083, ADR-0095, S-0227: activity_log says what an activity's cost was
// charged to: the item a planner planned, the items an orchestrator named,
// the issues an analyzer named, or the kind's project strategic total.
func TestActivityLogSaysWhatTheCostWasChargedTo(t *testing.T) {
	logged := func(kind string, cost float64, planned string, shared ...string) *serve.Logged {
		return &serve.Logged{Entry: workitem.ActivityEntry{Cost: cost}, Activity: &workitem.Activity{Kind: kind}, Planned: planned, Shared: shared}
	}
	for _, c := range []struct {
		name   string
		logged *serve.Logged
		failed error
		to     string
		line   string
	}{
		{"planned", logged("planner", 1, "S-0001"), nil, "S-0001", "charged to S-0001, the item its run planned, and the items above it"},
		{"one named", logged("orchestrator", 1, "", "T-0001"), nil, "T-0001", "charged to T-0001 and the items above it"},
		{"named", logged("orchestrator", 1, "", "T-0001", "S-0002"), nil, "T-0001,S-0002", "charged evenly to T-0001, S-0002 and the items above them"},
		{"none named", logged("orchestrator", 1, ""), nil, "", "charged to no item: left in the orchestrator's project strategic total"},
		{"analyzer", logged("analyzer", 1, ""), nil, "", "charged to no item: left in the analyzer's project strategic total"},
		{"one issue", logged("analyzer", 1, "", "I-0001"), nil, "I-0001", "charged to I-0001, the issue it named"},
		{"issues", logged("analyzer", 1, "", "I-0001", "I-0002"), nil, "I-0001,I-0002", "charged evenly to I-0001, I-0002, the issues it named"},
		{"issue failed", logged("analyzer", 1, "", "I-0001"), errors.New("I-0002: locked"), "I-0001", "charged to I-0001, the issue it named; I-0002: locked; what was not charged is left in the project strategic total"},
		{"nothing spent", logged("orchestrator", 0, ""), nil, "", ""},
		{"failed", logged("orchestrator", 1, "", "T-0001"), errors.New("S-0002: locked"), "T-0001", "charged to T-0001 and the items above it; S-0002: locked; what was not charged is left in the project strategic total"},
		{"failed, nothing spent", logged("planner", 0, ""), errors.New("locked"), "", "locked; what was not charged is left in the project strategic total"},
	} {
		to, line := chargedTo(c.logged, c.failed)
		if strings.Join(to, ",") != c.to || line != c.line {
			t.Errorf("%s: %v, %q; want %s, %q", c.name, to, line, c.to, c.line)
		}
	}
}
