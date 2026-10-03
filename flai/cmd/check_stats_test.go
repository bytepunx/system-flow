package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func TestCheckCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	out, errOut, code := runIn(t, "../internal/metrics/testdata/good", "check", "--strict")
	if code != 0 || !strings.Contains(out, "10 items checked, 0 errors, 1 warnings (1 that --strict passes over") || !strings.Contains(out, "epic.lags-stories") {
		t.Fatalf("good: %d %s %s", code, out, errOut)
	}
	out, _, code = runIn(t, ".", "check", "../internal/check/testdata/bad")
	if code != 1 || !strings.Contains(out, "wip/kanban/tasks/T-001-t.md:7: error: item.parent-missing: parent S-009 does not exist") {
		t.Fatalf("bad: %d\n%s", code, out)
	}
	out, _, code = runIn(t, "../internal/check/testdata/bad", "check", "--json")
	if code != 1 {
		t.Fatal("json check should still exit 1")
	}
	var res struct {
		Errors   int `json:"errors"`
		Findings []struct {
			Rule string `json:"rule"`
		} `json:"findings"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || res.Errors == 0 || len(res.Findings) == 0 {
		t.Fatalf("json: %v %s", err, out)
	}
}

// S-0243: the summary says how many warnings --strict passes over, and only
// when there are some.
func TestCheckSummaryNamesTheAdvisoryWarnings(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 4 {
		s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: fmt.Sprintf("Story %d", i+1), Owner: "t", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		s.Status = workitem.Review
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	out, _, _ := runIn(t, root, "check")
	if !strings.Contains(out, "warning: board.wip-limit: 4 stories in review, limit 3") ||
		!strings.Contains(out, " warnings (1 that --strict passes over: only the operator clears them, by accepting or by moving an epic)\n") {
		t.Errorf("summary should name the advisory warning:\n%s", out)
	}
}

func TestStatsCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dir := "../internal/metrics/testdata/good"
	a := func(args ...string) (string, string, int) {
		return runInAt(t, dir, time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC), args...)
	}
	out, errOut, code := a("stats")
	if code != 0 {
		t.Fatalf("stats: %s", errOut)
	}
	for _, want := range []string{"completed 2 · cancelled 1 (33%)", "throughput 0.5/week", "WIP now 1", "cycle time   p50 1d · p85 1d2h · max 1d2h", "flow efficiency 88%", "S-004  in-progress  1d2h   Four"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// what the two stories done in the window spent, per story, per minute
	// of agent work, and per dollar, then per model (S-0163)
	for _, want := range []string{"usage: 4.1M tokens · $2.02 (estimated in part) · 30m0s of agent work, over 2 done",
		"per story 2.0M tokens, $1.01, 15m of agent work · per agent minute 136.7K tokens · per dollar 2.0M tokens",
		"claude-haiku-4-5  100.0K tokens · $0.02 · 10.0K tokens/min (1)", "claude-opus-5-5  4.0M tokens · $2.00 · 133.3K tokens/min (2)"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	out, _, _ = a("stats", "--by", "nature", "--since", "6w")
	if !strings.Contains(out, "[feature]") || !strings.Contains(out, "[research]") {
		t.Errorf("grouped:\n%s", out)
	}
	out, _, code = a("stats", "--json", "--type", "task")
	var rep struct {
		Type    string `json:"type"`
		Summary struct {
			Completed int `json:"completed"`
		} `json:"summary"`
		CFD []any `json:"cfd"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &rep) != nil || rep.Type != "task" || len(rep.CFD) == 0 {
		t.Fatalf("json: %d %s", code, out)
	}
	if _, errOut, code := a("stats", "--since", "soon"); code == 0 || !strings.Contains(errOut, "--since") {
		t.Errorf("bad window: %s", errOut)
	}
	out, _, code = a("stats", "--json", "--bucket", "week")
	var spent struct {
		Usage struct {
			Bucket string `json:"bucket"`
			Spend  map[string]struct {
				Items   int `json:"items"`
				Buckets []struct {
					At    string `json:"at"`
					Items int    `json:"items"`
				} `json:"buckets"`
			} `json:"spend"`
		} `json:"usage"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &spent) != nil {
		t.Fatalf("json by the week: %d %s", code, out)
	}
	u := spent.Usage
	// S-002 was done on 12 August, a Wednesday; now is in the week of 31 August
	if st := u.Spend["story"]; u.Bucket != "week" || st.Items != 2 || len(st.Buckets) != 5 || st.Buckets[0].At != "2026-08-03T00:00:00Z" || st.Buckets[1].Items != 1 {
		t.Errorf("stories by the week: %+v", u)
	}
	if ta := u.Spend["task"]; ta.Items != 1 || len(u.Spend["epic"].Buckets) != 0 {
		t.Errorf("tasks and epics by the week: %+v", u)
	}
	if _, errOut, code := a("stats", "--bucket", "minute"); code == 0 || !strings.Contains(errOut, "hour, day, or week") {
		t.Errorf("a bucket that is not one: %s", errOut)
	}
	if _, errOut, code := a("stats", "--bucket", "hour", "--since", "90d"); code == 0 || !strings.Contains(errOut, "31 days or less") {
		t.Errorf("an hour over 90 days: %s", errOut)
	}
}

// ADR-0079: flai stats prints each strategic agent's totals from its
// activity document, --json carries them and the window's log entries under
// strategic, and an unreadable document stops stats as an unreadable item
// does.
func TestStatsReportsTheStrategicAgents(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)
	a := func(args ...string) (string, string, int) { return runInAt(t, root, at, args...) }
	out, errOut, code := a("stats")
	if code != 0 || strings.Contains(out, "strategic") {
		t.Fatalf("no documents: %d %s %s", code, out, errOut)
	}
	write := func(doc *workitem.Activity) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, "wip", "agents", doc.Kind+".md"), []byte(doc.Marshal()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(&workitem.Activity{Kind: workitem.ActivityPlanner, AccruedCost: 0.5213, AccruedSeconds: 823, TasksCompleted: 2, LastRun: "2026-10-03T18:00:00Z",
		Entries: []workitem.ActivityEntry{
			{At: time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC), Summary: "Long ago", Seconds: 100, Cost: 0.1},
			{At: time.Date(2026, 10, 3, 18, 0, 0, 0, time.UTC), Summary: "Drafted five stories", Items: []string{"E-0016", "S-0230"}, Seconds: 723, Cost: 0.4213, Estimated: true},
		}})
	write(&workitem.Activity{Kind: workitem.ActivityAnalyzer, AccruedCost: 0.05, AccruedSeconds: 60, TasksCompleted: 1, LastRun: "2026-10-01T08:00:00Z",
		Entries: []workitem.ActivityEntry{{At: time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC), Summary: "Read the issues", Seconds: 60, Cost: 0.05}}})
	out, errOut, code = a("stats")
	want := "\nstrategic agents (all time):\n  planner: 2 activities, 0.5213 USD, 823 s, last 2026-10-03T18:00:00Z\n  analyzer: 1 activity, 0.0500 USD, 60 s, last 2026-10-01T08:00:00Z\n"
	if code != 0 || !strings.Contains(out, want) || strings.Contains(out, "orchestrator") {
		t.Errorf("text: %d %s\n%s", code, errOut, out)
	}
	out, errOut, code = a("stats", "--json")
	var rep struct {
		Strategic []struct {
			Kind       string  `json:"kind"`
			Cost       float64 `json:"cost"`
			Seconds    int64   `json:"seconds"`
			Activities int     `json:"activities"`
			Log        []struct {
				At        string   `json:"at"`
				Estimated bool     `json:"estimated"`
				Items     []string `json:"items"`
			} `json:"log"`
		} `json:"strategic"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &rep) != nil {
		t.Fatalf("json: %d %s %s", code, errOut, out)
	}
	if s := rep.Strategic; len(s) != 2 || s[0].Kind != "planner" || s[0].Seconds != 823 || len(s[0].Log) != 1 || s[0].Log[0].At != "2026-10-03T18:00:00Z" ||
		!s[0].Log[0].Estimated || len(s[0].Log[0].Items) != 2 || s[1].Kind != "analyzer" || s[1].Activities != 1 || len(s[1].Log) != 1 {
		t.Errorf("strategic = %+v", s)
	}
	if err := os.WriteFile(filepath.Join(root, "wip", "agents", "orchestrator.md"), []byte("---\nkind: orchestrator\nmood: busy\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := a("stats"); code == 0 || !strings.Contains(errOut, "orchestrator.md") || !strings.Contains(errOut, "activity documents") {
		t.Errorf("an unreadable document: %d %s", code, errOut)
	}
}
