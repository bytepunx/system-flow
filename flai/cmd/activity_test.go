package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
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

	got, err := a.mcpActivity(ctx, root, workitem.ActivityOrchestrator, "Ordered the board", []string{"S-0001"}, "orchestrator")
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.At != "2026-10-03T10:01:30Z" || e.Seconds != 90 || e.Cost != half || !e.Estimated || strings.Join(e.Items, ",") != "S-0001" {
		t.Errorf("first entry: %+v", e)
	}
	now = start.Add(3 * time.Minute) // past the run's last call, at 10:02:00
	got, err = a.mcpActivity(ctx, root, workitem.ActivityOrchestrator, "Planned nothing new", nil, "orchestrator")
	if err != nil {
		t.Fatal(err)
	}
	if e := got.Entry; e.Seconds != 31 || e.Cost != half {
		t.Errorf("second entry, from the first's end to the second after the run's last call: %+v", e)
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
