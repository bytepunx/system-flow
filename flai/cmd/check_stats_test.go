package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/usage"
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
		!strings.Contains(out, " warnings (1 that --strict passes over: review over its limit and an epic behind its stories, which only the operator clears, or a narrative not yet written, which flai stream state writes)\n") {
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
		"claude-haiku-4-5  100.0K tokens · $0.02 · 10.0K tokens/min (1)", "claude-opus-5-5  4.0M tokens · $2.00 · 133.3K tokens/min (2)",
		// $2.02 over the half hour of S-001 and S-002, measured from their logs (ADR-0083)
		"  cost per agent hour $4.04, over every story measured from its logs\n"} {
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

// ADR-0083: flai stats prints what strategic agents spent on the items done
// in the window on a line of its own, after the agents' figures, and does so
// when agents spent nothing on them.
func TestStatsPrintsStrategicUsageApart(t *testing.T) {
	var out bytes.Buffer
	a := &app{out: &out}
	rate := 4.04
	u := metrics.UsageReport{CostPerAgentHour: &rate, Strategic: metrics.StrategicTotals{
		StrategicSpend: metrics.StrategicSpend{Items: 2, Tokens: 12_500, Cost: 0.4213, Seconds: 723}, Estimated: true,
		Kinds: []metrics.StrategicKindSpend{
			{Kind: "planner", StrategicSpend: metrics.StrategicSpend{Items: 2, Tokens: 12_000, Cost: 0.4113, Seconds: 700}},
			{Kind: "analyzer", StrategicSpend: metrics.StrategicSpend{Items: 1, Tokens: 500, Cost: 0.01, Seconds: 23}},
		}}}
	printUsage(a, workitem.Story, u)
	want := "  strategic usage, apart: 12.5K tokens · $0.42 (estimated) · 12m3s of strategic agent work, over 2 done (planner $0.41 · analyzer $0.01)\n" +
		"  cost per agent hour $4.04, over every story measured from its logs\n"
	if out.String() != want {
		t.Errorf("printed:\n%s\nwant:\n%s", out.String(), want)
	}
	out.Reset()
	printUsage(a, workitem.Story, metrics.UsageReport{})
	if out.Len() != 0 {
		t.Errorf("nothing spent printed %q", out.String())
	}
}

// S-0272, ADR-0105: the waiting section ends with the empty wakes, their
// mean left out when no item completed carries agents' usage.
func TestStatsPrintsTheEmptyWakes(t *testing.T) {
	var out bytes.Buffer
	a := &app{out: &out}
	mean := 1.5
	week := metrics.WaitWeek{Items: 8, Threads: metrics.ThreadWaitTotal{WaitTotal: metrics.WaitTotal{Total: 3600}}, Review: metrics.WaitTotal{Total: 7200}}
	printWaiting(a, metrics.Waiting{EmptyWakes: metrics.EmptyWakes{Count: 12, Mean: &mean}, Weeks: []metrics.WaitWeek{week}})
	if want := "  empty wakes  total 12 · mean 1.5 per item with usage\n"; !strings.HasSuffix(out.String(), want) {
		t.Errorf("printed:\n%s\nwant it to end with:\n%s", out.String(), want)
	}
	out.Reset()
	printWaiting(a, metrics.Waiting{Weeks: []metrics.WaitWeek{week}})
	if want := "  in review    total 2h · mean 15m\n  empty wakes  total 0\n"; !strings.HasSuffix(out.String(), want) {
		t.Errorf("printed:\n%s\nwant it to end with:\n%s", out.String(), want)
	}
	out.Reset()
	printWaiting(a, metrics.Waiting{EmptyWakes: metrics.EmptyWakes{Count: 3}, Weeks: []metrics.WaitWeek{{}}})
	if out.Len() != 0 {
		t.Errorf("nothing completed printed %q", out.String())
	}
}

// S-0213: the cost of delay section adds the items without a value now, per
// column with any, and what the ready stories will cost until pulled under
// each order, each only when there is something to show.
func TestStatsPrintsTheItemsWithoutAValueAndThePullOrdersCost(t *testing.T) {
	var out bytes.Buffer
	a := &app{out: &out}
	value := 10.0
	cols := func(ready, review float64) map[string]float64 {
		return map[string]float64{workitem.Backlog: 0, workitem.Ready: ready, workitem.InProgress: 0, workitem.Review: review}
	}
	counts := func(backlog, ready int) map[string]int {
		return map[string]int{workitem.Backlog: backlog, workitem.Ready: ready, workitem.InProgress: 0, workitem.Review: 0}
	}
	series := func(by string, total float64) metrics.OrderSeries {
		return metrics.OrderSeries{By: by, Total: total, Points: []metrics.OrderPoint{{At: "2026-10-07T12:00:00Z"}, {At: "2026-10-08T12:00:00Z", ID: "S-0002", Incurred: total}}}
	}
	report := func(without map[string]int, o metrics.CostOrder) *metrics.Report {
		return &metrics.Report{Items: []metrics.ItemMetrics{{CostOfDelay: &value}}, CostOfDelay: metrics.CostOfDelay{
			Days:  []metrics.CostDay{{Date: "2026-10-07", Outstanding: cols(30, 10), WithoutValue: without, Incurred: 1.5}},
			Order: o,
		}}
	}
	head := "\ncost of delay (per week of waiting):\n  outstanding now  backlog 0.00 · ready 30.00 · in-progress 0.00 · review 10.00\n  incurred in the window 1.50\n"
	for _, c := range []struct {
		name string
		rep  *metrics.Report
		want string
	}{
		{"both", report(counts(2, 1), metrics.CostOrder{Series: []metrics.OrderSeries{series(metrics.OrderCurrent, 12.5), series(metrics.OrderCOD, 9), series(metrics.OrderWSJF, 8)},
			Saving: 4.5, Cheaper: metrics.OrderWSJF, LeftOut: []string{"S-0003"}}),
			head + "  without a value now  backlog 2 · ready 1\n" +
				"  ready until pulled  pull order 12.50 · by cost of delay 9.00 · by WSJF 8.00; by WSJF saves 4.50; 1 left out, without a value or a forecast duration\n"},
		{"the pull order the cheapest", report(counts(0, 0), metrics.CostOrder{Series: []metrics.OrderSeries{series(metrics.OrderCurrent, 8), series(metrics.OrderCOD, 9), series(metrics.OrderWSJF, 8.5)},
			Saving: -0.5, Cheaper: metrics.OrderWSJF, LeftOut: []string{}}),
			head + "  ready until pulled  pull order 8.00 · by cost of delay 9.00 · by WSJF 8.50; the pull order is already the cheapest\n"},
		{"nothing to show", report(counts(0, 0), metrics.CostOrder{Series: []metrics.OrderSeries{
			{By: metrics.OrderCurrent, Points: []metrics.OrderPoint{{}}}, {By: metrics.OrderCOD, Points: []metrics.OrderPoint{{}}}, {By: metrics.OrderWSJF, Points: []metrics.OrderPoint{{}}},
		}, Cheaper: metrics.OrderCOD, LeftOut: []string{}}),
			head},
	} {
		out.Reset()
		printCostOfDelay(a, c.rep)
		if out.String() != c.want {
			t.Errorf("%s printed:\n%s\nwant:\n%s", c.name, out.String(), c.want)
		}
	}
}

// S-0293: the story agents' turns print with each class's share of the
// window's, then the days that have any, then each story; nothing prints
// when there are none.
func TestStatsPrintsTheStoryAgentsTurns(t *testing.T) {
	var out bytes.Buffer
	a := &app{out: &out}
	counts := func(ceremony, tests, wakes, edits, work int) metrics.TurnCounts {
		return metrics.TurnCounts{Turns: ceremony + tests + wakes + edits + work, Ceremony: ceremony, TestRuns: tests, EmptyWakes: wakes, HandEdits: edits, Work: work}
	}
	turns := metrics.Turns{
		Classes: usage.TurnClasses,
		Total:   counts(4, 2, 1, 1, 12),
		Days: []metrics.TurnsDay{
			{Day: "2026-10-04", TurnCounts: counts(3, 2, 0, 1, 4)},
			{Day: "2026-10-05", TurnCounts: counts(0, 0, 0, 0, 0)},
			{Day: "2026-10-06", TurnCounts: counts(1, 0, 1, 0, 8)},
		},
		Stories: []metrics.StoryTurns{
			{ID: "S-0270", Title: "First", Status: workitem.Done, TurnCounts: counts(3, 2, 0, 1, 4)},
			{ID: "S-0293", Title: "Second", Status: workitem.InProgress, TurnCounts: counts(1, 0, 1, 0, 8)},
		},
	}
	printTurns(a, turns)
	want := "\nstory agents' turns, over 2 stories in the window:\n" +
		"  total 20 · ceremony 4 (20%) · test runs 2 (10%) · empty wakes 1 (5%) · hand edits 1 (5%) · work 12 (60%)\n" +
		"  by day:\n" +
		"    2026-10-04  10 · ceremony 3 · test runs 2 · empty wakes 0 · hand edits 1 · work 4\n" +
		"    2026-10-06  10 · ceremony 1 · test runs 0 · empty wakes 1 · hand edits 0 · work 8\n" +
		"  by story:\n" +
		"    S-0270  10 · ceremony 3 · test runs 2 · empty wakes 0 · hand edits 1 · work 4  First\n" +
		"    S-0293  10 · ceremony 1 · test runs 0 · empty wakes 1 · hand edits 0 · work 8  Second\n"
	if out.String() != want {
		t.Errorf("printed:\n%s\nwant:\n%s", out.String(), want)
	}
	out.Reset()
	printTurns(a, metrics.Turns{Classes: usage.TurnClasses, Days: []metrics.TurnsDay{{Day: "2026-10-04"}}})
	if out.Len() != 0 {
		t.Errorf("no turns printed %q", out.String())
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
	// a story the planner planned carries what its run spent on it, and the
	// rest of the planner's totals is its project total (ADR-0095)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	planned, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Planned", Owner: "t", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	planned.Usage = &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}, Strategic: []usage.Strategic{
		{Kind: "planner", Seconds: 723, Estimated: true, Models: []usage.Model{{Model: "claude-opus-5-5", Output: 1000, Cost: 0.4213}}},
	}}
	if err := repo.Save(planned); err != nil {
		t.Fatal(err)
	}
	write(&workitem.Activity{Kind: workitem.ActivityPlanner, AccruedCost: 0.5213, AccruedSeconds: 823, TasksCompleted: 2, LastRun: "2026-10-03T18:00:00Z",
		Entries: []workitem.ActivityEntry{
			{At: time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC), Summary: "Long ago", Seconds: 100, Cost: 0.1},
			{At: time.Date(2026, 10, 3, 18, 0, 0, 0, time.UTC), Summary: "Drafted five stories", Items: []string{"E-0016", "S-0230"}, Seconds: 723, Cost: 0.4213, Estimated: true},
		}})
	write(&workitem.Activity{Kind: workitem.ActivityAnalyzer, AccruedCost: 0.05, AccruedSeconds: 60, TasksCompleted: 1, LastRun: "2026-10-01T08:00:00Z",
		Entries: []workitem.ActivityEntry{{At: time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC), Summary: "Read the issues", Seconds: 60, Cost: 0.05}}})
	out, errOut, code = a("stats")
	want := "\nstrategic agents (all time):\n" +
		"  planner: 2 activities, 0.5213 USD, 823 s (on items 0.4213 USD, 723 s; on issues 0.0000 USD, 0 s; project 0.1000 USD, 100 s), last 2026-10-03T18:00:00Z\n" +
		"  analyzer: 1 activity, 0.0500 USD, 60 s (on items 0.0000 USD, 0 s; on issues 0.0000 USD, 0 s; project 0.0500 USD, 60 s), last 2026-10-01T08:00:00Z\n"
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
			Items      struct {
				Cost    float64 `json:"cost"`
				Seconds int64   `json:"seconds"`
			} `json:"items"`
			Project struct {
				Cost    float64 `json:"cost"`
				Seconds int64   `json:"seconds"`
			} `json:"project"`
			Log []struct {
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
	if s := rep.Strategic; len(s) == 2 && (s[0].Items.Cost != 0.4213 || s[0].Items.Seconds != 723 || s[0].Project.Cost != 0.1 || s[0].Project.Seconds != 100 ||
		s[1].Items.Cost != 0 || s[1].Project.Cost != 0.05 || s[1].Project.Seconds != 60) {
		t.Errorf("items and project = %+v", s)
	}
	if err := os.WriteFile(filepath.Join(root, "wip", "agents", "orchestrator.md"), []byte("---\nkind: orchestrator\nmood: busy\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := a("stats"); code == 0 || !strings.Contains(errOut, "orchestrator.md") || !strings.Contains(errOut, "activity documents") {
		t.Errorf("an unreadable document: %d %s", code, errOut)
	}
}

// S-0227: flai stats reads the issues and reports the analyzer's usage on
// each, counted in the kind's issues amount until flai issue story makes a
// story from it, and then once, as the story's, under items; a project with
// no design/issues folder reports none.
func TestStatsReportsTheIssuesStrategicUsageOnce(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_STORY", "")
	root := tempProject(t)
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)
	a := func(args ...string) (string, string, int) { return runInAt(t, root, at, args...) }
	if out, errOut, code := a("stats", "--json"); code != 0 || !strings.Contains(out, `"strategic_issues": []`) {
		t.Fatalf("no issues folder: %d %s %s", code, errOut, out)
	}
	if _, errOut, code := a("issue", "new", "Review waits a day", "--class", "efficiency"); code != 0 {
		t.Fatalf("issue new: %s", errOut)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := issues.ChargeStrategic(repo, "I-0001", workitem.ActivityAnalyzer,
		&usage.Usage{Seconds: 30, Models: []usage.Model{{Model: "claude-opus-5-5", Output: 100, Cost: 0.03}}}); err != nil {
		t.Fatal(err)
	}
	doc := &workitem.Activity{Kind: workitem.ActivityAnalyzer, AccruedCost: 0.05, AccruedSeconds: 60, TasksCompleted: 1, LastRun: "2026-10-01T08:00:00Z"}
	if err := os.WriteFile(filepath.Join(root, "wip", "agents", "analyzer.md"), []byte(doc.Marshal()), 0o644); err != nil {
		t.Fatal(err)
	}
	type report struct {
		Strategic []struct {
			Kind                   string
			Items, Issues, Project metrics.StrategicAmount
		} `json:"strategic"`
		StrategicIssues []metrics.StrategicIssue `json:"strategic_issues"`
	}
	read := func() report {
		t.Helper()
		out, errOut, code := a("stats", "--json")
		var rep report
		if code != 0 || json.Unmarshal([]byte(out), &rep) != nil || len(rep.Strategic) != 1 || len(rep.StrategicIssues) != 1 {
			t.Fatalf("json: %d %s %s", code, errOut, out)
		}
		return rep
	}
	rep := read()
	if s := rep.Strategic[0]; s.Items != (metrics.StrategicAmount{}) || s.Issues != (metrics.StrategicAmount{Cost: 0.03, Seconds: 30}) || s.Project != (metrics.StrategicAmount{Cost: 0.02, Seconds: 30}) {
		t.Errorf("before the story: %+v", s)
	}
	if is := rep.StrategicIssues[0]; is.ID != "I-0001" || is.Status != "open" || is.Story != "" || !is.Counted || len(is.Strategic) != 1 || is.Strategic[0].Cost != 0.03 {
		t.Errorf("the issue before the story: %+v", is)
	}
	out, _, _ := a("stats")
	want := "\nstrategic usage on issues (all time):\n  I-0001 open: analyzer 0.0300 USD, 30 s; counted on the issue  Review waits a day\n"
	if !strings.Contains(out, want) || !strings.Contains(out, "(on items 0.0000 USD, 0 s; on issues 0.0300 USD, 30 s; project 0.0200 USD, 30 s)") {
		t.Errorf("text before the story:\n%s", out)
	}

	if _, errOut, code := a("issue", "story", "I-0001"); code != 0 {
		t.Fatalf("issue story: %s", errOut)
	}
	rep = read()
	if s := rep.Strategic[0]; s.Items != (metrics.StrategicAmount{Cost: 0.03, Seconds: 30}) || s.Issues != (metrics.StrategicAmount{}) || s.Project != (metrics.StrategicAmount{Cost: 0.02, Seconds: 30}) {
		t.Errorf("after the story: %+v", s)
	}
	if is := rep.StrategicIssues[0]; is.Story != "S-0001" || is.Counted || len(is.Strategic) != 1 {
		t.Errorf("the issue after the story: %+v", is)
	}
	if out, _, _ := a("stats"); !strings.Contains(out, "  I-0001 open: analyzer 0.0300 USD, 30 s; counted on S-0001  Review waits a day\n") {
		t.Errorf("text after the story:\n%s", out)
	}
}

// S-0205: flai stats carries forecasts, cost of delay, waiting, claims, and
// the strategic agents' use by the day, and prints them. Outside git it
// leaves the touches drift out, warns on stderr, and still succeeds.
func TestStatsReportsPlanningWaitingAndClaims(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("LOG_FORMAT", "json")
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story := func(title string, value float64, created string, moves ...string) *workitem.Item {
		t.Helper()
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Owner: "t", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		it.Created = created
		it.Transitions = []workitem.Transition{{To: workitem.Backlog, At: created, By: "t"}}
		for i := 0; i < len(moves); i += 2 {
			it.Transitions = append(it.Transitions, workitem.Transition{To: moves[i], At: moves[i+1], By: "t"})
		}
		it.Status = it.Transitions[len(it.Transitions)-1].To
		it.CostOfDelay = &workitem.CostOfDelay{Value: &value, By: "designer", At: created}
		return it
	}
	// done after 48 hours in progress, forecast at 24, delivery forecast 10
	// hours early, estimated at 50 hours; one in ready, valued at 50
	done := story("Done", 100, "2026-09-28T09:00:00Z", workitem.Ready, "2026-09-28T10:00:00Z", workitem.InProgress, "2026-09-29T10:00:00Z",
		workitem.Review, "2026-10-01T06:00:00Z", workitem.Done, "2026-10-01T10:00:00Z")
	done.Forecast = &workitem.Forecast{Duration: "24h", Delivery: "2026-10-01T00:00:00Z"}
	done.Estimate = "50h"
	// its agents woke to nothing four times (S-0272)
	done.Usage = &usage.Usage{Source: usage.SourceLog, Seconds: 3600, EmptyWakes: 4, Models: []usage.Model{{Model: "claude-opus-5-5", Input: 10, Cost: 1}}}
	// and its agent took ten turns on two days of the window, and some
	// before it (S-0293)
	done.Usage.Turns = []usage.TurnDay{
		{Day: "2026-08-01", Work: 50},
		{Day: "2026-09-29", Ceremony: 2, Work: 6},
		{Day: "2026-10-01", TestRuns: 1, HandEdits: 1},
	}
	waiting := story("Waiting", 50, "2026-09-30T09:00:00Z", workitem.Ready, "2026-09-30T12:00:00Z")
	for _, it := range []*workitem.Item{done, waiting} {
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	// its agent asked at noon on its first day and was answered two hours on
	th, err := threads.New(repo, threads.NewOptions{Title: "Which way", On: "S-0001", Author: "agent", Text: "Left?", Now: time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := threads.Reply(repo, th.ID, "designer", "Left.", time.Date(2026, 9, 29, 14, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AppendActivity(workitem.ActivityPlanner, workitem.ActivityEntry{At: time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC), Summary: "Planned.", Seconds: 90, Cost: 0.25}); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 10, 3, 20, 0, 0, 0, time.UTC)

	out, errOut, code := runInAt(t, root, at, "stats", "--json")
	if code != 0 {
		t.Fatalf("stats --json: %s", errOut)
	}
	var rep map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"forecasts", "cost_of_delay", "waiting", "claims", "strategic_days", "turns"} {
		if len(rep[key]) == 0 || string(rep[key]) == "null" {
			t.Errorf("--json lacks %s", key)
		}
	}
	// the window's 31 days, from the 3rd of September to the 3rd of
	// October, and the story's turns on them (S-0293)
	var turns metrics.Turns
	if err := json.Unmarshal(rep["turns"], &turns); err != nil {
		t.Fatal(err)
	}
	if n := len(turns.Days); n != 31 || turns.Days[0].Day != "2026-09-03" || turns.Days[n-1].Day != "2026-10-03" {
		t.Errorf("turns.days = %d from %+v", n, turns.Days)
	}
	if want := (metrics.TurnCounts{Turns: 10, Ceremony: 2, TestRuns: 1, HandEdits: 1, Work: 6}); turns.Total != want {
		t.Errorf("turns.total = %+v, want %+v", turns.Total, want)
	}
	if len(turns.Stories) != 1 || turns.Stories[0].ID != done.ID || turns.Stories[0].Status != workitem.Done || turns.Stories[0].Turns != 10 {
		t.Errorf("turns.stories = %+v, want %s with 10", turns.Stories, done.ID)
	}
	var wakes struct {
		Waiting struct {
			EmptyWakes struct {
				Count int      `json:"count"`
				Mean  *float64 `json:"mean"`
			} `json:"empty_wakes"`
			Weeks []struct {
				EmptyWakes int `json:"empty_wakes"`
			} `json:"weeks"`
		} `json:"waiting"`
		Items []struct {
			ID    string `json:"id"`
			Usage struct {
				EmptyWakes *int `json:"empty_wakes"`
			} `json:"usage"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &wakes); err != nil {
		t.Fatal(err)
	}
	if e := wakes.Waiting.EmptyWakes; e.Count != 4 || e.Mean == nil || *e.Mean != 4 {
		t.Errorf("waiting.empty_wakes = %+v, want 4 over 1 item with usage", e)
	}
	weekWakes := 0
	for _, w := range wakes.Waiting.Weeks {
		weekWakes += w.EmptyWakes
	}
	if weekWakes != 4 {
		t.Errorf("weeks' empty wakes add up to %d, want 4", weekWakes)
	}
	for _, it := range wakes.Items {
		if it.ID == done.ID && (it.Usage.EmptyWakes == nil || *it.Usage.EmptyWakes != 4) {
			t.Errorf("%s usage.empty_wakes = %v, want 4", it.ID, it.Usage.EmptyWakes)
		}
	}
	var claims map[string]json.RawMessage
	if err := json.Unmarshal(rep["claims"], &claims); err != nil || claims["drift"] != nil || string(claims["limit"]) != "2" {
		t.Errorf("claims outside git: %v %s", err, rep["claims"])
	}
	// Waiting was held for want of touches the 18 hours Done was in progress
	// beside it; outside git no week counts exact touches (S-0214)
	var weeks []map[string]json.RawMessage
	if err := json.Unmarshal(claims["weeks"], &weeks); err != nil || len(weeks) == 0 {
		t.Fatalf("claims.weeks: %v %s", err, claims["weeks"])
	}
	var heldTotal metrics.HeldByCode
	for _, w := range weeks {
		var h metrics.HeldByCode
		if err := json.Unmarshal(w["held_seconds"], &h); err != nil {
			t.Fatalf("claims.weeks[].held_seconds: %v %s", err, w["held_seconds"])
		}
		heldTotal.Overlap += h.Overlap
		heldTotal.After += h.After
		heldTotal.NoTouches += h.NoTouches
		for _, key := range []string{"stories", "exact", "exact_share"} {
			if w[key] != nil {
				t.Errorf("claims.weeks[].%s outside git: %s", key, w[key])
			}
		}
	}
	if heldTotal != (metrics.HeldByCode{NoTouches: 18 * 3600}) {
		t.Errorf("claims.weeks held_seconds add up to %+v, want 18h of no-touches", heldTotal)
	}
	warned := false
	for _, line := range strings.Split(errOut, "\n") {
		var ev map[string]any
		if json.Unmarshal([]byte(line), &ev) == nil && ev["level"] == "WARN" && ev["component"] == "stats" && ev["err"] != nil && ev["detail"] != nil {
			warned = true
		}
	}
	if !warned {
		t.Errorf("no warning that the commits could not be read:\n%s", errOut)
	}

	out, _, code = runInAt(t, root, at, "stats")
	for _, want := range []string{
		"\nforecast and estimate error (absolute):\n  forecast     p50 1d · p85 1d (n=1)\n  delivery     p50 10h · p85 10h (n=1)\n  estimate     p50 2h · p85 2h (n=1)\n",
		// the story in ready has a value but no forecast duration (S-0213)
		"\ncost of delay (per week of waiting):\n  outstanding now  backlog 0.00 · ready 50.00 · in-progress 0.00 · review 0.00\n  incurred in the window 39.57\n" +
			"  ready until pulled  none placed; 1 left out, without a value or a forecast duration\n",
		"\nwaiting, over 1 completed:\n  on threads   total 2h · mean 2h\n  in review    total 4h · mean 4h\n  empty wakes  total 4 · mean 4.0 per item with usage\n",
		// the turns follow the waiting and come before the claims (S-0293)
		"per item with usage\n\nstory agents' turns, over 1 story in the window:\n" +
			"  total 10 · ceremony 2 (20%) · test runs 1 (10%) · empty wakes 0 (0%) · hand edits 1 (10%) · work 6 (60%)\n" +
			"  by day:\n" +
			"    2026-09-29  8 · ceremony 2 · test runs 0 · empty wakes 0 · hand edits 0 · work 6\n" +
			"    2026-10-01  2 · ceremony 0 · test runs 1 · empty wakes 0 · hand edits 1 · work 0\n" +
			"  by story:\n" +
			"    S-0001  10 · ceremony 2 · test runs 1 · empty wakes 0 · hand edits 1 · work 6  Done\n\nclaims:\n",
		"\nclaims:\n  held in ready  total 0m · mean 0m, over 1 completed\n  held by reason  overlap 0m · after 0m · no-touches 18h, every story over the window's weeks\n  in progress now 0 of a limit of 2\n",
		"\nstrategic agents in the window:\n  $0.25 · 1m, beside 1 completed (usage $1.00 per item)\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if code != 0 || strings.Contains(out, "touches drift") || strings.Contains(out, "exact touches") {
		t.Errorf("text outside git: %d\n%s", code, out)
	}
	// the dashboard's read answers the same and warns the same
	if _, rerr := sameAnswer(t, root, hostapi.Host{}, "stats.get", `{}`, "stats"); rerr != nil {
		t.Errorf("stats.get outside git: %+v", rerr)
	}
}

// S-0214: the claims print the time held by reason over the window's weeks
// and how many stories with commits completed in them had exact touches; a
// week git could not be read for, or with none completed, adds nothing.
func TestStatsPrintsHeldByReasonAndExactTouches(t *testing.T) {
	n := func(v int) *int { return &v }
	drift := []metrics.TouchDrift{{ID: "S-0001"}, {ID: "S-0002", OutsideCount: 1}, {ID: "S-0003"}}
	rep := &metrics.Report{WindowStart: "2026-09-01T00:00:00Z", Claims: metrics.Claims{
		Drift: &drift,
		Weeks: []metrics.ClaimWeek{
			{Week: "2026-W36", HeldSeconds: metrics.HeldByCode{Overlap: 3 * 3600, NoTouches: 1800}, Stories: n(2), Exact: n(1)},
			{Week: "2026-W37", HeldSeconds: metrics.HeldByCode{After: 26 * 3600}, Stories: n(0), Exact: n(0)},
			{Week: "2026-W38", Stories: n(1), Exact: n(1)},
		},
	}}
	var out bytes.Buffer
	printClaims(&app{out: &out}, rep)
	want := "\nclaims:\n" +
		"  held by reason  overlap 3h · after 1d2h · no-touches 30m, every story over the window's weeks\n" +
		"  touches drift  3 stories with commits · 1 files outside their touches · 0 touches unchanged\n" +
		"  exact touches  2 of 3 completed with commits (67%)\n"
	if out.String() != want {
		t.Errorf("claims:\n%q\nwant\n%q", out.String(), want)
	}

	// outside git and with nothing held, neither line
	rep.Claims = metrics.Claims{Weeks: []metrics.ClaimWeek{{Week: "2026-W36"}}}
	out.Reset()
	printClaims(&app{out: &out}, rep)
	if out.Len() != 0 {
		t.Errorf("claims with nothing held, outside git:\n%s", out.String())
	}
}
