package metrics

import (
	"math"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var now = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func load(t *testing.T) *Report {
	t.Helper()
	repo, err := workitem.Open("testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	return Compute(items, Options{Now: now})
}

func TestFixtureNumbers(t *testing.T) {
	rep := load(t)
	s := rep.Summary
	if s.Completed != 2 || s.Cancelled != 1 || !near(s.CancellationRate, 1.0/3) || s.WIP != 1 {
		t.Errorf("counts: %+v", s)
	}
	if !near(s.ThroughputWeek, 2/(30.0/7)) {
		t.Errorf("throughput %v", s.ThroughputWeek)
	}
	if s.CycleTime.P50 != 86400 || s.CycleTime.P85 != 93600 || s.CycleTime.Max != 93600 || s.CycleTime.Mean != 90000 || s.CycleTime.Count != 2 {
		t.Errorf("cycle: %+v", s.CycleTime)
	}
	if s.LeadTime.P50 != 174600 || s.LeadTime.P85 != 183600 || s.LeadTime.Mean != 179100 {
		t.Errorf("lead: %+v", s.LeadTime)
	}
	if s.QueueTime.P50 != 86400 || s.QueueTime.Max != 86400 {
		t.Errorf("queue: %+v", s.QueueTime)
	}
	if !near(s.FlowEfficiency, (72000.0/93600+1)/2) {
		t.Errorf("flow efficiency %v", s.FlowEfficiency)
	}
	total := 358200.0
	for st, want := range map[string]float64{"backlog": 5400, "ready": 172800, "in-progress": 122400, "review": 57600} {
		if !near(s.InStateShare[st], want/total) {
			t.Errorf("share %s = %v want %v", st, s.InStateShare[st], want/total)
		}
	}
	var s1 ItemMetrics
	for _, m := range rep.Items {
		if m.ID == "S-001" {
			s1 = m
		}
	}
	if s1.Blocked != 21600 || s1.EstError == nil || !near(*s1.EstError, 0.3) || *s1.CycleTime != 93600 || *s1.LeadTime != 183600 || *s1.QueueTime != 86400 {
		t.Errorf("S-001: %+v", s1)
	}
	if s1.InState["in-progress"] != 86400 || s1.InState["review"] != 7200 || s1.InState["backlog"] != 3600 {
		t.Errorf("S-001 in-state: %v", s1.InState)
	}
	if len(rep.Aging) != 1 || rep.Aging[0].ID != "S-004" || rep.Aging[0].Age != 93600 || rep.Aging[0].OverP85 {
		t.Errorf("aging: %+v", rep.Aging)
	}
	doneTotal := 0
	for _, w := range rep.Throughput {
		doneTotal += w.Done
	}
	if doneTotal != 2 {
		t.Errorf("throughput buckets: %+v", rep.Throughput)
	}
	all := rep.Burnup["all"]
	if len(all) == 0 || all[0].Date != "2026-07-01" || all[len(all)-1].Date != "2026-09-01" {
		t.Fatalf("burnup range: %d points", len(all))
	}
	last := all[len(all)-1]
	if last.Scope != 4 || last.Done != 3 {
		t.Errorf("burnup last: %+v", last)
	}
	if e := rep.Burnup["E-001"]; len(e) != len(all) || e[len(e)-1].Done != 3 {
		t.Errorf("per-epic burnup: %+v", e[len(e)-1])
	}
	c := rep.CFD[len(rep.CFD)-1].Counts
	if c["in-progress"] != 1 || c["done"] != 3 || c["cancelled"] != 1 || c["backlog"] != 0 {
		t.Errorf("cfd last: %v", c)
	}
	// on 2026-08-02 S-001 was in progress and S-005 done
	for _, p := range rep.CFD {
		if p.Date == "2026-08-02" && (p.Counts["in-progress"] != 1 || p.Counts["done"] != 1) {
			t.Errorf("cfd 08-02: %v", p.Counts)
		}
	}
}

func TestGroupingAndTypes(t *testing.T) {
	repo, _ := workitem.Open("testdata/good")
	items, _ := repo.List(true)
	rep := Compute(items, Options{Now: now, By: "nature"})
	if len(rep.Groups) != 3 || rep.Groups[0].Group != "feature" || rep.Groups[0].Completed != 1 || rep.Groups[1].Group != "improvement" || rep.Groups[2].Group != "research" || rep.Groups[2].Cancelled != 1 {
		for _, g := range rep.Groups {
			t.Logf("%+v", g)
		}
		t.Error("grouping by nature wrong")
	}
	tasks := Compute(items, Options{Now: now, Type: workitem.Task, Since: 90 * 24 * time.Hour})
	if tasks.Summary.Completed != 3 || tasks.Summary.WIP != 1 {
		t.Errorf("tasks: %+v", tasks.Summary)
	}
}

func TestHuman(t *testing.T) {
	for in, want := range map[float64]string{30: "0m", 2700: "45m", 7200: "2h", 9000: "2h30m", 93600: "1d2h", 172800: "2d"} {
		if got := Human(in); got != want {
			t.Errorf("Human(%v) = %s want %s", in, got, want)
		}
	}
}
