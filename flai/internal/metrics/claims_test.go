package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// claimItem is an item created at the time given, with the touches given,
// moved as the pairs of state and time say.
func claimItem(id, typ, created string, touches []string, moves ...string) *workitem.Item {
	it := &workitem.Item{ID: id, Type: typ, Status: workitem.Backlog, Created: created, Touches: touches}
	for i := 0; i < len(moves); i += 2 {
		it.Transitions = append(it.Transitions, workitem.Transition{To: moves[i], At: moves[i+1]})
		it.Status = moves[i]
	}
	return it
}

// S-0205: a story's held time is replayed from the states of the time, each
// hold lasting until the next creation or transition, up to now.
func TestHeldSecondsReplayTheHoldRules(t *testing.T) {
	const day = "2026-08-29T00:00:00Z"
	s, task := workitem.Story, workitem.Task
	rd, ip, rv, done := workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done
	archived := func(it *workitem.Item) *workitem.Item {
		it.Archived = true
		return it
	}
	after := func(it *workitem.Item, ids ...string) *workitem.Item {
		it.After = ids
		return it
	}
	cases := []struct {
		name     string
		items    []*workitem.Item
		projects []manifest.Project
		want     map[string]*float64
	}{
		{
			name: "overlap while the other is in progress, not in review, archived since",
			items: []*workitem.Item{
				archived(claimItem("S-0001", s, day, []string{"flai"}, ip, "2026-08-31T00:00:00Z", rv, "2026-08-31T06:00:00Z", done, "2026-08-31T08:00:00Z")),
				claimItem("S-0002", s, day, []string{"flai/internal"}, rd, "2026-08-30T12:00:00Z", ip, "2026-08-31T10:00:00Z", done, "2026-08-31T12:00:00Z"),
			},
			want: map[string]*float64{"S-0001": nil, "S-0002": val(6 * 3600)},
		},
		{
			name: "after until the story it names is done, archived since",
			items: []*workitem.Item{
				archived(claimItem("S-0001", s, day, []string{"a"}, ip, "2026-08-30T06:00:00Z", done, "2026-08-31T06:00:00Z")),
				after(claimItem("S-0002", s, day, []string{"b"}, rd, "2026-08-29T12:00:00Z", ip, "2026-08-31T09:00:00Z"), "S-1"),
			},
			want: map[string]*float64{"S-0002": val(42 * 3600)},
		},
		{
			name: "no touches while another is in progress",
			items: []*workitem.Item{
				claimItem("S-0001", s, day, []string{"a"}, ip, "2026-08-30T06:00:00Z", rv, "2026-08-30T18:00:00Z", done, "2026-08-31T00:00:00Z"),
				claimItem("S-0002", s, day, nil, rd, "2026-08-30T03:00:00Z", ip, "2026-08-31T03:00:00Z"),
			},
			want: map[string]*float64{"S-0002": val(12 * 3600)},
		},
		{
			name: "never held",
			items: []*workitem.Item{
				claimItem("S-0001", s, day, []string{"a"}, rd, "2026-08-30T00:00:00Z", ip, "2026-08-30T06:00:00Z", done, "2026-08-30T07:00:00Z"),
			},
			want: map[string]*float64{"S-0001": val(0)},
		},
		{
			name: "a task's touches while the task is open",
			items: []*workitem.Item{
				claimItem("S-0001", s, day, []string{"a"}, ip, "2026-08-30T00:00:00Z", done, "2026-08-31T00:00:00Z"),
				func() *workitem.Item {
					it := claimItem("T-0001", task, day, []string{"b"}, ip, "2026-08-30T01:00:00Z", done, "2026-08-30T10:00:00Z")
					it.Parent = "S-0001"
					return it
				}(),
				claimItem("S-0002", s, day, []string{"b"}, rd, "2026-08-30T00:00:00Z", ip, "2026-08-31T00:00:00Z"),
			},
			want: map[string]*float64{"S-0001": nil, "S-0002": val(10 * 3600)},
		},
		{
			name: "still in ready at now",
			items: []*workitem.Item{
				claimItem("S-0001", s, day, []string{"a"}, ip, "2026-09-01T00:00:00Z"),
				claimItem("S-0002", s, day, []string{"a"}, rd, "2026-08-31T12:00:00Z"),
			},
			want: map[string]*float64{"S-0002": val(12 * 3600)},
		},
		{
			name:     "a project's tag read as its path",
			projects: []manifest.Project{{Name: "flai", Path: "flai/", Tags: []string{"go"}}},
			items: []*workitem.Item{
				claimItem("S-0001", s, day, []string{"flai/internal/metrics"}, ip, "2026-08-30T00:00:00Z", done, "2026-08-30T05:00:00Z"),
				claimItem("S-0002", s, day, []string{"go"}, rd, "2026-08-29T12:00:00Z", ip, "2026-08-31T00:00:00Z"),
			},
			want: map[string]*float64{"S-0002": val(5 * 3600)},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rep := Compute(c.items, Options{Now: now, Projects: c.projects})
			got := map[string]*float64{}
			for _, m := range rep.Items {
				got[m.ID] = m.HeldSeconds
			}
			for id, want := range c.want {
				if g := got[id]; (g == nil) != (want == nil) || (g != nil && *g != *want) {
					t.Errorf("%s held = %v, want %v", id, deref(g), deref(want))
				}
			}
		})
	}
}

// S-0205: an item that is not a story has no held time.
func TestHeldSecondsAreAbsentForTasks(t *testing.T) {
	tk := claimItem("T-0001", workitem.Task, "2026-08-29T00:00:00Z", nil, workitem.Ready, "2026-08-30T00:00:00Z")
	rep := Compute([]*workitem.Item{tk}, Options{Now: now, Type: workitem.Task})
	if len(rep.Items) != 1 || rep.Items[0].HeldSeconds != nil {
		t.Errorf("task items = %+v, want one without held time", rep.Items)
	}
	// S-0214: only stories are held, so a report of tasks holds none.
	for _, d := range rep.Claims.Days {
		if d.Held != 0 {
			t.Errorf("task days = %+v, want none held", rep.Claims.Days)
		}
	}
	for _, w := range rep.Claims.Weeks {
		if w.HeldSeconds != (HeldByCode{}) {
			t.Errorf("task weeks = %+v, want no time held", rep.Claims.Weeks)
		}
	}
}

// S-0205: the stories in progress at the end of each day of the window, with
// the limit when the board has one.
func TestClaimsDaysCountInProgressAgainstTheLimit(t *testing.T) {
	const day = "2026-08-29T00:00:00Z"
	ip := workitem.InProgress
	items := []*workitem.Item{
		claimItem("S-0001", workitem.Story, day, nil, ip, "2026-08-29T10:00:00Z", workitem.Done, "2026-08-30T08:00:00Z"),
		// in review at the end of 30 August, back in progress after
		claimItem("S-0002", workitem.Story, day, nil, ip, "2026-08-30T06:00:00Z", workitem.Review, "2026-08-30T20:00:00Z", ip, "2026-08-31T02:00:00Z"),
		claimItem("S-0003", workitem.Story, "2026-08-31T00:00:00Z", nil, ip, "2026-09-01T01:00:00Z"),
		claimItem("E-0001", workitem.Epic, day, nil, ip, "2026-08-29T00:00:00Z"),
		claimItem("T-0001", workitem.Task, day, nil, ip, "2026-08-29T00:00:00Z"),
	}
	rep := Compute(items, Options{Now: now, Since: threeDays, WIPLimit: 3})
	want := []ClaimDay{{"2026-08-29", 1, 0}, {"2026-08-30", 0, 0}, {"2026-08-31", 1, 0}, {"2026-09-01", 2, 0}}
	if !reflect.DeepEqual(rep.Claims.Days, want) {
		t.Errorf("days = %+v, want %+v", rep.Claims.Days, want)
	}
	data, _ := json.Marshal(rep.Claims)
	if !strings.Contains(string(data), `"limit":3`) {
		t.Errorf("claims = %s, want the limit", data)
	}
	data, _ = json.Marshal(Compute(items, Options{Now: now, Since: threeDays}).Claims)
	if strings.Contains(string(data), `"limit"`) {
		t.Errorf("claims = %s, want no limit without one", data)
	}
}

// S-0205: each story's committed files against its touches and its tasks'
// not cancelled, read as a claim reads them.
func TestClaimsDriftComparesCommitsWithTouches(t *testing.T) {
	const day = "2026-08-29T00:00:00Z"
	task := func(id, status string, touches ...string) *workitem.Item {
		it := claimItem(id, workitem.Task, day, touches)
		it.Parent, it.Status = "S-0001", status
		return it
	}
	items := []*workitem.Item{
		claimItem("S-0003", workitem.Story, day, []string{"flai"}),
		claimItem("S-0001", workitem.Story, day, []string{"flai/internal/metrics/", "./design/system/metrics.md"}),
		claimItem("S-0002", workitem.Story, day, []string{"flai"}),
		task("T-0001", workitem.Done, "dashboard"),
		task("T-0002", workitem.Cancelled, "scripts"),
		task("T-0003", workitem.Backlog, "docs/users"),
	}
	commits := map[string][]string{
		"S-0001": {"flai/internal/metrics/claims.go", "scripts/flai.sh", "flaiover/src/app.ts", "flai/internal/metrics/claims.go",
			"design/system/metrics.md", "flai/internal/metricsx/a.go", "README.md"},
		"S-0003": {"flai/main.go"},
		"T-0001": {"flai/x.go"},
	}
	opt := Options{Now: now, Commits: commits, Projects: []manifest.Project{{Name: "flaiover", Path: "flaiover", Tags: []string{"dashboard"}}}}
	rep := Compute(items, opt)
	want := []TouchDrift{
		{
			ID:             "S-0001",
			Committed:      []string{"README.md", "design/system/metrics.md", "flai/internal/metrics/claims.go", "flai/internal/metricsx/a.go", "flaiover/src/app.ts", "scripts/flai.sh"},
			Outside:        []string{"README.md", "flai/internal/metricsx/a.go", "scripts/flai.sh"},
			Unchanged:      []string{"docs/users"},
			OutsideCount:   3,
			UnchangedCount: 1,
		},
		{ID: "S-0003", Committed: []string{"flai/main.go"}, Outside: []string{}, Unchanged: []string{}},
	}
	if rep.Claims.Drift == nil || !reflect.DeepEqual(*rep.Claims.Drift, want) {
		t.Errorf("drift = %+v, want %+v", rep.Claims.Drift, want)
	}
	data, _ := json.Marshal(rep.Claims)
	if !strings.Contains(string(data), `"outside":[],"unchanged":[]`) {
		t.Errorf("claims = %s, want empty lists, not null", data)
	}

	opt.Commits = map[string][]string{}
	if data, _ := json.Marshal(Compute(items, opt).Claims); !strings.Contains(string(data), `"drift":[]`) {
		t.Errorf("claims = %s, want an empty drift with commits read and none naming a story", data)
	}
	opt.Commits = nil
	if data, _ := json.Marshal(Compute(items, opt).Claims); strings.Contains(string(data), `"drift"`) {
		t.Errorf("claims = %s, want no drift when git could not be read", data)
	}
}

// S-0214: the stories held at each day's end, the time held per ISO week by
// the reason each hold is named by, and per week the stories of the drift
// completed in it and how many of them had exact touches (ADR-0113).
func TestClaimsWeeksAndHeldDays(t *testing.T) {
	const old = "2026-08-15T00:00:00Z"
	s := workitem.Story
	rd, ip, rv, done := workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done
	s2 := claimItem("S-0002", s, old, []string{"a/x"}, rd, "2026-08-28T12:00:00Z", ip, "2026-09-01T06:00:00Z")
	s2.After = []string{"S-0003"}
	items := []*workitem.Item{
		claimItem("S-0001", s, old, []string{"a"}, ip, "2026-08-28T00:00:00Z", rv, "2026-09-01T00:30:00Z", done, "2026-09-01T02:00:00Z"),
		// held from 28 August 12:00 by S-0003 in after and by S-0001's
		// overlap: after; from 30 August 12:00, S-0003 done, by the overlap
		// alone, until S-0001 moves to review
		s2,
		claimItem("S-0003", s, old, []string{"b"}, rd, "2026-08-21T00:00:00Z", ip, "2026-08-29T00:00:00Z", done, "2026-08-30T12:00:00Z"),
		// no touches, held from 31 August 10:00 to now by S-0005
		claimItem("S-0004", s, old, nil, rd, "2026-08-31T10:00:00Z"),
		claimItem("S-0005", s, old, []string{"c"}, ip, "2026-08-31T09:00:00Z"),
		claimItem("S-0009", s, old, []string{"h"}, ip, "2026-08-30T00:00:00Z", workitem.Cancelled, "2026-08-30T01:00:00Z"),
		claimItem("S-0010", s, old, []string{"d"}, ip, "2026-08-31T18:00:00Z", done, "2026-08-31T20:00:00Z"),
		claimItem("S-0011", s, old, []string{"e", "f"}, ip, "2026-08-31T19:00:00Z", done, "2026-08-31T21:00:00Z"),
		// done before the window's start, in its first week
		claimItem("S-0012", s, old, []string{"g"}, ip, "2026-08-19T00:00:00Z", done, "2026-08-20T00:00:00Z"),
	}
	commits := map[string][]string{
		"S-0001": {"a/f.go"},
		"S-0002": {"a/x/y.go"},         // not done
		"S-0003": {"b/x.go", "c/y.go"}, // c/y.go outside
		"S-0009": {"z/1.go"},           // cancelled
		"S-0010": {"d/x.go"},           // exact
		"S-0011": {"e/1.go"},           // f unchanged
		"S-0012": {"g/1.go"},           // before the window
	}
	opt := Options{Now: now, Since: 10 * 24 * time.Hour, Commits: commits}
	rep := Compute(items, opt)

	held := map[string]*float64{}
	for _, m := range rep.Items {
		held[m.ID] = m.HeldSeconds
	}
	for id, want := range map[string]*float64{"S-0001": nil, "S-0002": val((48 + 12 + 24.5) * 3600), "S-0003": val(0), "S-0004": val(26 * 3600)} {
		if g := held[id]; (g == nil) != (want == nil) || (g != nil && *g != *want) {
			t.Errorf("%s held = %v, want %v", id, deref(g), deref(want))
		}
	}

	wantDays := []ClaimDay{
		{"2026-08-22", 0, 0}, {"2026-08-23", 0, 0}, {"2026-08-24", 0, 0}, {"2026-08-25", 0, 0}, {"2026-08-26", 0, 0}, {"2026-08-27", 0, 0},
		{"2026-08-28", 1, 1}, {"2026-08-29", 2, 1}, {"2026-08-30", 1, 1}, {"2026-08-31", 2, 2}, {"2026-09-01", 2, 1},
	}
	if !reflect.DeepEqual(rep.Claims.Days, wantDays) {
		t.Errorf("days = %+v, want %+v", rep.Claims.Days, wantDays)
	}

	n := func(v int) *int { return &v }
	wantWeeks := []ClaimWeek{
		// nothing held, and no story of the drift completed in the window
		{Week: "2026-W34", Start: "2026-08-17", Stories: n(0), Exact: n(0)},
		// the whole first week counts: S-0002 from 28 August 12:00
		{Week: "2026-W35", Start: "2026-08-24", HeldSeconds: HeldByCode{Overlap: 12 * 3600, After: 48 * 3600},
			Stories: n(1), Exact: n(0), ExactShare: val(0)},
		{Week: "2026-W36", Start: "2026-08-31", HeldSeconds: HeldByCode{Overlap: 24.5 * 3600, NoTouches: 26 * 3600},
			Stories: n(3), Exact: n(2), ExactShare: val(2.0 / 3)},
	}
	if !reflect.DeepEqual(rep.Claims.Weeks, wantWeeks) {
		t.Errorf("weeks = %s, want %s", weeksString(rep.Claims.Weeks), weeksString(wantWeeks))
	}
	data, _ := json.Marshal(rep.Claims.Weeks[0])
	if want := `{"week":"2026-W34","start":"2026-08-17","held_seconds":{"overlap":0,"after":0,"no-touches":0},"stories":0,"exact":0}`; string(data) != want {
		t.Errorf("week = %s, want %s", data, want)
	}

	opt.Commits = nil
	rep = Compute(items, opt)
	data, _ = json.Marshal(rep.Claims.Weeks[2])
	if want := `{"week":"2026-W36","start":"2026-08-31","held_seconds":{"overlap":88200,"after":0,"no-touches":93600}}`; string(data) != want {
		t.Errorf("week without git = %s, want %s", data, want)
	}
}

// ADR-0096: an overlap wholly inside the manifest's shared paths holds
// nothing, in the replay as on the board.
func TestHeldSecondsLeaveSharedPathsOut(t *testing.T) {
	const old = "2026-08-15T00:00:00Z"
	items := []*workitem.Item{
		claimItem("S-0001", workitem.Story, old, []string{"docs/users/flai.md"}, workitem.InProgress, "2026-08-31T12:00:00Z"),
		claimItem("S-0002", workitem.Story, old, []string{"docs/users/flai.md", "flai/cmd"}, workitem.Ready, "2026-08-31T12:00:00Z"),
	}
	for _, c := range []struct {
		name   string
		shared []string
		want   float64
	}{
		{"not shared", nil, 24 * 3600},
		{"shared", []string{"docs/users/flai.md"}, 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			rep := Compute(items, Options{Now: now, Since: threeDays, Shared: manifest.Claims{Shared: c.shared}})
			var got *float64
			for _, m := range rep.Items {
				if m.ID == "S-0002" {
					got = m.HeldSeconds
				}
			}
			if got == nil || *got != c.want {
				t.Errorf("S-0002 held = %v, want %v", deref(got), c.want)
			}
			if w := rep.Claims.Weeks[len(rep.Claims.Weeks)-1].HeldSeconds; w.Overlap != c.want {
				t.Errorf("this week's overlap = %v, want %v", w.Overlap, c.want)
			}
			wantHeld := 0
			if c.want > 0 {
				wantHeld = 1
			}
			if d := rep.Claims.Days[len(rep.Claims.Days)-1]; d.Held != wantHeld {
				t.Errorf("today held = %d, want %d", d.Held, wantHeld)
			}
		})
	}
}

// weeksString shows weeks with their pointers' values.
func weeksString(ws []ClaimWeek) string {
	data, _ := json.Marshal(ws)
	return string(data)
}
