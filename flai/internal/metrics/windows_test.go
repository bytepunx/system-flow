package metrics

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var updateWindows = flag.Bool("update-windows", false, "rewrite testdata/windows.golden.json from what Compute gives")

// windowCases are the windows the report is pinned at (S-0608).
var windowCases = []struct {
	name string
	days int
}{
	{"1d", 1}, {"7d", 7}, {"30d", 30}, {"90d", 90}, {"365d", 365},
}

// windowsGolden is where the reports at windowCases are pinned.
var windowsGolden = filepath.Join("testdata", "windows.golden.json")

// windowsItems are testdata/good's items with touches, after, costs of
// delay, and forecasts set in memory, and stories and a task added, so that
// the hold replay holds on overlap, after, and no touches, the ready column
// has a cost order, and the touch drift has stories inside, outside, and
// short of their touches. The 90 and 365-day windows start before the first
// item was created.
func windowsItems(t *testing.T) []*workitem.Item {
	t.Helper()
	repo, err := workitem.Open("testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	value := func(v float64) *workitem.CostOfDelay { return &workitem.CostOfDelay{Value: &v} }
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	set := func(id string, touches []string, cost *workitem.CostOfDelay) {
		it := byID[id]
		if it == nil {
			t.Fatalf("testdata/good holds no %s", id)
		}
		it.Touches, it.CostOfDelay = touches, cost
	}
	set("S-001", []string{"flai/internal/a.go"}, value(15))
	set("S-002", []string{"flai/internal/b.go", "docs"}, nil)
	set("S-004", []string{"flai/internal/a.go", "docs/a.md"}, value(20))
	set("S-005", []string{"dashboard"}, value(5))

	s6 := claimItem("S-006", workitem.Story, "2026-08-24T09:00:00Z", []string{"flai/internal/a.go"},
		workitem.Ready, "2026-08-24T10:00:00Z", workitem.InProgress, "2026-08-26T10:00:00Z",
		workitem.Review, "2026-08-28T10:00:00Z", workitem.Done, "2026-08-29T10:00:00Z")
	s6.CostOfDelay = value(30)
	t6 := claimItem("T-006", workitem.Task, "2026-08-24T09:30:00Z", []string{"flai/internal/a_test.go"},
		workitem.Ready, "2026-08-24T10:00:00Z", workitem.InProgress, "2026-08-26T10:00:00Z",
		workitem.Done, "2026-08-28T09:00:00Z")
	t6.Parent = "S-006"
	s7 := claimItem("S-007", workitem.Story, "2026-08-27T09:00:00Z", nil,
		workitem.Ready, "2026-08-27T10:00:00Z")
	s7.CostOfDelay, s7.Forecast = value(50), &workitem.Forecast{Duration: "4h"}
	s8 := claimItem("S-008", workitem.Story, "2026-08-28T09:00:00Z", []string{"flaiover/src/x.ts"},
		workitem.Ready, "2026-08-28T10:00:00Z")
	s8.After, s8.CostOfDelay, s8.Forecast = []string{"S-004"}, value(10), &workitem.Forecast{Duration: "2h"}
	s9 := claimItem("S-009", workitem.Story, "2026-08-30T09:00:00Z", nil)
	for _, it := range []*workitem.Item{s6, t6, s7, s8, s9} {
		it.Nature, it.Parent = "feature", cmp.Or(it.Parent, "E-001")
	}
	return append(items, s6, t6, s7, s8, s9)
}

// windowsOptions are the options every window shares: commits that fall
// inside, outside, and short of the stories' touches, a project whose tag a
// touch names, the board's limit and pull order, and testdata/good's threads.
// Commits are keyed by the canonical ID, as storygit gives them.
func windowsOptions(t *testing.T, days int) Options {
	t.Helper()
	repo, err := workitem.Open("testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	list, err := threads.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	return Options{
		Now:      now,
		Since:    time.Duration(days) * 24 * time.Hour,
		Projects: []manifest.Project{{Name: "flaiover", Path: "flaiover", Tags: []string{"dashboard"}}},
		WIPLimit: 1,
		Commits: map[string][]string{
			"S-0001": {"flai/internal/a.go", "README.md"},
			"S-0002": {"flai/internal/b.go"},
			"S-0004": {"flai/internal/a.go"},
			"S-0005": {"flaiover/src/y.ts"},
			"S-0006": {"flai/internal/a.go", "flai/internal/a_test.go"},
		},
		Order:   []string{"S-008", "S-007"},
		Threads: list,
	}
}

// S-0608: the report over testdata/good, enriched as windowsItems says, stays
// what it was before S-0608 changed how the replay is kept and the days are
// swept, at windows of 1, 7, 30, 90, and 365 days. Run go test -update-windows
// to rewrite the golden file once a change to the report is intended, and
// review its diff. It holds one report a line, compact.
func TestReportAtEveryWindowIsTheGolden(t *testing.T) {
	var golden bytes.Buffer
	golden.WriteString("{\n")
	got := map[string][]byte{}
	for i, c := range windowCases {
		rep := Compute(windowsItems(t), windowsOptions(t, c.days))
		checkWindowFilled(t, c.name, rep)
		a, err := json.Marshal(rep)
		if err != nil {
			t.Fatal(err)
		}
		b, err := json.Marshal(Compute(windowsItems(t), windowsOptions(t, c.days)))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(a, b) {
			t.Errorf("%s: two computes of the same items differ at %s", c.name, firstDifference(a, b))
		}
		got[c.name] = a
		name, _ := json.Marshal(c.name)
		golden.Write(name)
		golden.WriteString(": ")
		golden.Write(a)
		if i < len(windowCases)-1 {
			golden.WriteString(",")
		}
		golden.WriteString("\n")
	}
	golden.WriteString("}\n")
	data, err := os.ReadFile(windowsGolden)
	if *updateWindows || errors.Is(err, fs.ErrNotExist) {
		if err := os.WriteFile(windowsGolden, golden.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		if !*updateWindows {
			t.Fatalf("%s was missing and is written now; review it and run the test again", windowsGolden)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	var want map[string]json.RawMessage
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatalf("%s: %v; run go test -update-windows to write it again", windowsGolden, err)
	}
	if len(want) != len(windowCases) {
		t.Errorf("%s holds %d windows, the test %d; run go test -update-windows and review the diff", windowsGolden, len(want), len(windowCases))
	}
	for _, c := range windowCases {
		if w, ok := want[c.name]; !ok {
			t.Errorf("%s holds no window %s", windowsGolden, c.name)
		} else if !bytes.Equal(got[c.name], w) {
			t.Errorf("%s: the report differs from %s at %s; if the change is intended, run go test -update-windows and review the diff", c.name, windowsGolden, firstDifference(got[c.name], w))
		}
	}
}

// checkWindowFilled fails the test when a section the golden is there to pin
// is empty at a window, so that the golden cannot pass by pinning nothing.
func checkWindowFilled(t *testing.T, name string, rep *Report) {
	t.Helper()
	c, cod := rep.Claims, rep.CostOfDelay
	if len(c.Days) == 0 || len(c.Weeks) == 0 || c.Drift == nil || len(*c.Drift) == 0 {
		t.Errorf("%s: claims has %d days, %d weeks, and drift %v", name, len(c.Days), len(c.Weeks), c.Drift)
	}
	if len(cod.Days) == 0 || len(cod.Weeks) == 0 || len(cod.Order.Series) == 0 || len(cod.Order.Series[0].Points) < 2 || cod.Order.Series[0].Total == 0 {
		t.Errorf("%s: cost of delay has %d days, %d weeks, and order %+v", name, len(cod.Days), len(cod.Weeks), cod.Order)
	}
	if name != "30d" {
		return
	}
	var held HeldByCode
	for _, w := range c.Weeks {
		held.Overlap += w.HeldSeconds.Overlap
		held.After += w.HeldSeconds.After
		held.NoTouches += w.HeldSeconds.NoTouches
	}
	if held.Overlap == 0 || held.After == 0 || held.NoTouches == 0 {
		t.Errorf("%s: the replay held %+v; every reason should hold some time", name, held)
	}
}

// firstDifference shows where got first differs from want, with some of
// each around it.
func firstDifference(got, want []byte) string {
	i := 0
	for i < len(got) && i < len(want) && got[i] == want[i] {
		i++
	}
	around := func(b []byte) string {
		lo, hi := max(0, i-80), min(len(b), i+80)
		return string(b[lo:hi])
	}
	return "byte " + strconv.Itoa(i) + "\n got: …" + around(got) + "…\nwant: …" + around(want) + "…"
}
