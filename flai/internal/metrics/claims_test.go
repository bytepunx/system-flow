package metrics

import (
	"encoding/json"
	"fmt"
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

// holdRuleCase is a small history, the manifest's projects, and the seconds
// each story in it is held.
type holdRuleCase struct {
	name     string
	items    []*workitem.Item
	projects []manifest.Project
	want     map[string]*float64
}

// S-0205: a story's held time is replayed from the states of the time, each
// hold lasting until the next creation or transition, up to now.
func TestHeldSecondsReplayTheHoldRules(t *testing.T) {
	for _, c := range holdRuleCases() {
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

// I-0155: the replay that judges the holds only when they can change gives
// what the reference, which judges them at every event time, gives.
func TestReplayHoldsMatchesTheReferenceOnTheHoldRules(t *testing.T) {
	for _, c := range holdRuleCases() {
		t.Run(c.name, func(t *testing.T) {
			var items []*workitem.Item
			for _, it := range c.items {
				if it.Type == workitem.Story {
					items = append(items, it)
				}
			}
			opt := Options{Now: now, Since: threeDays, Projects: c.projects}
			start := now.Add(-opt.Since)
			got, want := replayHolds(items, c.items, start, opt), replayHoldsReference(items, c.items, start, opt)
			for id, w := range want.seconds {
				if g, ok := got.seconds[id]; !ok || (g == nil) != (w == nil) || (g != nil && *g != *w) {
					t.Errorf("%s held = %v, want %v", id, deref(g), deref(w))
				}
			}
			if len(got.seconds) != len(want.seconds) {
				t.Errorf("stories replayed = %d, want %d", len(got.seconds), len(want.seconds))
			}
			if !reflect.DeepEqual(got.days, want.days) {
				t.Errorf("days = %v, want %v", got.days, want.days)
			}
			if !reflect.DeepEqual(got.weeks, want.weeks) {
				t.Errorf("weeks = %+v, want %+v", got.weeks, want.weeks)
			}
		})
	}
}

// I-0155: over a history at this repository's scale, smaller in short mode,
// the holds are judged anew no more often than an event changes what a hold
// reads, and less often than at every event time with a story in ready.
func TestReplayHoldsJudgesOnlyWhenAHoldCanChange(t *testing.T) {
	n := historyStories
	if testing.Short() {
		n = 60
	}
	all, opt, start := claimsHistory(n)
	items := storiesOf(all)
	_, built, changed := replayCounts(items, all, opt.Now)
	got := replayHolds(items, all, start, opt).judged
	if got == 0 || got > changed || got >= built {
		t.Errorf("judged %d times, want some, at most the %d that change a hold and fewer than the %d with a story in ready", got, changed, built)
	}
}

// holdRuleCases are the hold rules' small histories (S-0205).
func holdRuleCases() []holdRuleCase {
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
	return []holdRuleCase{
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

// I-0286, ADR-0284: a story's commits under flai's folders, the digest
// entries an acceptance swept in, the issue it was made from, and an analysis
// report, are not files it changed, and its touches under them, wip or an
// issue it closes, are not touches it missed: on both sides the drift leaves
// them out, so a story whose own files match its touches is exact. Without
// the folders, the same story is counted inexact on both sides, as I-0286
// found it in 80 of 115 stories.
func TestClaimsDriftLeavesOutFlaiFoldersOnBothSides(t *testing.T) {
	const day = "2026-08-29T00:00:00Z"
	done := []string{workitem.InProgress, "2026-08-30T00:00:00Z", workitem.Done, "2026-08-31T00:00:00Z"}
	task := claimItem("T-0001", workitem.Task, day, []string{"design/issues/I-0001-a.md", "wip/agents/S-0001.md"})
	task.Parent, task.Status = "S-0001", workitem.Done
	items := []*workitem.Item{
		claimItem("S-0001", workitem.Story, day, []string{"flai/internal/metrics/claims.go", "design/issues/summary.md"}, done...),
		task,
		// only flai's bookkeeping: no entry
		claimItem("S-0002", workitem.Story, day, []string{"flai/x.go"}, done...),
	}
	commits := map[string][]string{
		"S-0001": {"flai/internal/metrics/claims.go", "design/digest/D-0001-a.md", "design/issues/I-0001-a.md",
			"design/issues/summary.md", "design/analysis/2026-08-29-all.md"},
		"S-0002": {"design/digest/D-0002-b.md"},
	}
	opt := Options{Now: now, Since: 10 * 24 * time.Hour, Commits: commits,
		FlaiFolders: []string{"wip", "design/digest", "design/analysis", "design/issues"}}
	rep := Compute(items, opt)
	want := []TouchDrift{{ID: "S-0001", Committed: []string{"flai/internal/metrics/claims.go"}, Outside: []string{}, Unchanged: []string{}}}
	if rep.Claims.Drift == nil || !reflect.DeepEqual(*rep.Claims.Drift, want) {
		t.Errorf("drift = %+v, want %+v", rep.Claims.Drift, want)
	}
	exact := 0
	for _, w := range rep.Claims.Weeks {
		if w.Exact != nil {
			exact += *w.Exact
		}
	}
	if exact != 1 {
		t.Errorf("exact stories over the weeks = %d, want 1: %+v", exact, rep.Claims.Weeks)
	}

	opt.FlaiFolders = nil
	drift := *Compute(items, opt).Claims.Drift
	if len(drift) != 2 || drift[0].OutsideCount != 2 || drift[0].UnchangedCount != 1 {
		t.Errorf("drift without flai's folders = %+v, want S-0001 with the digest and analysis files outside and the wip touch unchanged, and S-0002", drift)
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

// A hold still open at a now with a fraction of a second counts whole seconds,
// as metrics.md's precision says.
func TestHeldSecondsAreWholeWhenNowIsNot(t *testing.T) {
	const old = "2026-08-15T00:00:00Z"
	items := []*workitem.Item{
		claimItem("S-0001", workitem.Story, old, []string{"flai/cmd"}, workitem.InProgress, "2026-08-31T12:00:00Z"),
		claimItem("S-0002", workitem.Story, old, []string{"flai/cmd"}, workitem.Ready, "2026-08-31T12:00:00Z"),
	}
	rep := Compute(items, Options{Now: now.Add(500 * time.Millisecond), Since: threeDays})
	for _, m := range rep.Items {
		if m.ID == "S-0002" && (m.HeldSeconds == nil || *m.HeldSeconds != 24*3600) {
			t.Errorf("S-0002 held = %v, want %v", deref(m.HeldSeconds), 24*3600)
		}
	}
	if w := rep.Claims.Weeks[len(rep.Claims.Weeks)-1].HeldSeconds; w.Overlap != 24*3600 {
		t.Errorf("this week's overlap = %v, want %v", w.Overlap, 24*3600)
	}
}

// weeksString shows weeks with their pointers' values.
func weeksString(ws []ClaimWeek) string {
	data, _ := json.Marshal(ws)
	return string(data)
}

// S-0608: the day sweep counts, on every day, the items stateAt finds in
// progress at the day's end, at windows from a day to a year, the longest
// starting long before every item, so that those days count none.
func TestClaimsDaysSweepAgreesWithStateAt(t *testing.T) {
	for _, set := range [][]*workitem.Item{costItems(), unvaluedItems(), sweepItems()} {
		for _, at := range []time.Time{now, time.Date(2026, 8, 31, 23, 59, 59, 0, time.UTC)} {
			for _, days := range []int{1, 3, 7, 30, 365} {
				var items []*workitem.Item
				for _, it := range set {
					if it.Type == workitem.Story {
						items = append(items, it)
					}
				}
				start := at.Add(-time.Duration(days) * 24 * time.Hour)
				opt := Options{Now: at}
				got := claims(items, set, newTimelines(items), start, opt, replayHolds(items, set, start, opt))
				i := 0
				for d := firstDay(start); !d.After(at); d = d.AddDate(0, 0, 1) {
					end := d.Add(24*time.Hour - time.Second)
					want := 0
					for _, it := range items {
						if stateAt(it, end) == workitem.InProgress {
							want++
						}
					}
					if i >= len(got.Days) {
						t.Fatalf("%d days to %s: %d days, want more", days, at, len(got.Days))
					}
					if g := got.Days[i]; g.Date != d.Format("2006-01-02") || g.InProgress != want {
						t.Errorf("%d days to %s: day %d = %s with %d in progress, want %s with %d", days, at, i, g.Date, g.InProgress, d.Format("2006-01-02"), want)
					}
					i++
				}
				if i != len(got.Days) {
					t.Errorf("%d days to %s: %d days, want %d", days, at, len(got.Days), i)
				}
			}
		}
	}
}

// sameReplay reports where got differs from want: a story's held seconds,
// the stories replayed, the days, or the weeks.
func sameReplay(t *testing.T, name string, got, want holdReplay) {
	t.Helper()
	for id, w := range want.seconds {
		if g, ok := got.seconds[id]; !ok || (g == nil) != (w == nil) || (g != nil && *g != *w) {
			t.Errorf("%s: %s held = %v, want %v", name, id, deref(g), deref(w))
		}
	}
	if len(got.seconds) != len(want.seconds) {
		t.Errorf("%s: stories replayed = %d, want %d", name, len(got.seconds), len(want.seconds))
	}
	if !reflect.DeepEqual(got.days, want.days) {
		t.Errorf("%s: days = %v, want %v", name, got.days, want.days)
	}
	if !reflect.DeepEqual(got.weeks, want.weeks) {
		t.Errorf("%s: weeks = %+v, want %+v", name, got.weeks, want.weeks)
	}
}

// keyedHistory is a history of 60 stories with Options keyed to a root of
// the test's own, so that what it keeps is its own.
func keyedHistory(t *testing.T) ([]*workitem.Item, Options) {
	all, opt, _ := claimsHistory(60)
	opt.Root, opt.ItemsGeneration, opt.CommitsKey = "keyed:"+t.Name(), 1, "tips-1"
	return all, opt
}

// S-0608: the hold timeline kept for a root is answered again while the item
// generation and the manifest's projects and shared paths are the same and it
// answers for the request's now; a moved key, or a now before its last event,
// replays again. Without a root or a generation nothing is kept.
func TestTheHoldReplayIsKeptUntilItsKeysMove(t *testing.T) {
	all, opt := keyedHistory(t)
	now := opt.Now.Truncate(time.Second)
	replay := func(opt Options, now time.Time) (*holdTimeline, bool) {
		t.Helper()
		return keptTimeline(all, now, opt)
	}
	first, replayed := replay(opt, now)
	if !replayed {
		t.Fatal("the first request did not replay")
	}
	for _, c := range []struct {
		name   string
		change func(Options) (Options, time.Time)
		again  bool
	}{
		{"the same keys", func(o Options) (Options, time.Time) { return o, now }, false},
		{"an hour later", func(o Options) (Options, time.Time) { return o, now.Add(time.Hour) }, false},
		{"another window", func(o Options) (Options, time.Time) { o.Since = 365 * 24 * time.Hour; return o, now }, false},
		{"a now before the last event", func(o Options) (Options, time.Time) { return o, now.Add(-30 * 24 * time.Hour) }, true},
		{"a moved generation", func(o Options) (Options, time.Time) { o.ItemsGeneration = 2; return o, now }, true},
		{"other projects", func(o Options) (Options, time.Time) { o.ItemsGeneration, o.Projects = 2, o.Projects[:1]; return o, now }, true},
		{"other shared paths", func(o Options) (Options, time.Time) {
			o.ItemsGeneration, o.Shared = 2, manifest.Claims{}
			return o, now
		}, true},
	} {
		o, at := c.change(opt)
		got, replayed := replay(o, at)
		if replayed != c.again {
			t.Errorf("%s: replayed %v, want %v", c.name, replayed, c.again)
		}
		if !c.again && got != first {
			t.Errorf("%s: answered another timeline than the kept one", c.name)
		}
	}
	for _, o := range []Options{{Now: opt.Now, Projects: opt.Projects, ItemsGeneration: 1}, {Now: opt.Now, Projects: opt.Projects, Root: "keyed:none"}} {
		if _, replayed := replay(o, now); !replayed {
			t.Errorf("root %q generation %d: answered from a keep, want a replay", o.Root, o.ItemsGeneration)
		}
		if _, replayed := replay(o, now); !replayed {
			t.Errorf("root %q generation %d: kept, want nothing kept", o.Root, o.ItemsGeneration)
		}
	}
}

// S-0608: the touches drift kept for a root is answered again while the item
// generation, the commits' key, and the projects are the same, and worked out
// again when any of them moves; the stories among a request's items are the
// kept drift's, as touchDrift gives them.
func TestTheTouchesDriftIsKeptUntilItsKeysMove(t *testing.T) {
	all, opt := keyedHistory(t)
	opt.Commits = map[string][]string{}
	for i, it := range storiesOf(all) {
		if i%3 == 0 && len(it.Touches) > 0 {
			opt.Commits[workitem.CanonicalID(it.ID)] = []string{it.Touches[0], "elsewhere/" + it.ID + ".go"}
		}
	}
	if _, worked := keptTouchDrift(all, opt); !worked {
		t.Fatal("the first request did not work the drift out")
	}
	if _, worked := keptTouchDrift(all, opt); worked {
		t.Error("a second request with the same keys worked the drift out again")
	}
	for _, c := range []struct {
		name   string
		change func(*Options)
	}{
		{"a moved commits' key", func(o *Options) { o.CommitsKey = "tips-2" }},
		{"a moved generation", func(o *Options) { o.ItemsGeneration = 2 }},
		{"other projects", func(o *Options) { o.Projects = o.Projects[:1] }},
		{"other flai folders", func(o *Options) { o.FlaiFolders = []string{"wip", "design/digest"} }},
	} {
		c.change(&opt)
		if _, worked := keptTouchDrift(all, opt); !worked {
			t.Errorf("%s: answered from the keep, want the drift worked out again", c.name)
		}
		if _, worked := keptTouchDrift(all, opt); worked {
			t.Errorf("%s: the request after it worked the drift out again", c.name)
		}
	}
	items := storiesOf(all)
	half := items[:len(items)/2]
	for name, its := range map[string][]*workitem.Item{"every story": items, "half the stories": half, "no story": nil} {
		got, want := driftFor(its, all, opt), touchDrift(its, all, opt)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: kept drift = %+v, want %+v", name, got, want)
		}
	}
	if len(touchDrift(items, all, opt)) == 0 {
		t.Error("no story drifted in the history, want some")
	}
}

// S-0608: a replay kept at one now and laid out at a later one, an hour,
// three days, or a new ISO week on, is what a replay at that now gives, at
// each window T-2662 pins, for every story and for half of them. It is
// replayed again only when an event came between the two nows.
func TestAKeptReplayAtALaterNowIsAFreshReplay(t *testing.T) {
	all, opt := keyedHistory(t)
	stories := storiesOf(all)
	if _, replayed := keptTimeline(all, opt.Now.Truncate(time.Second), opt); !replayed {
		t.Fatal("the first request did not replay")
	}
	replays := 0
	for _, later := range []time.Duration{0, time.Hour, 3 * 24 * time.Hour, 8 * 24 * time.Hour} {
		for _, days := range []int{1, 7, 30, 90, 365} {
			for name, items := range map[string][]*workitem.Item{"every story": stories, "half the stories": stories[:len(stories)/2]} {
				o := opt
				o.Now, o.Since = opt.Now.Add(later), time.Duration(days)*24*time.Hour
				start := o.Now.Add(-o.Since)
				label := fmt.Sprintf("%s later, %d days, %s", later, days, name)
				// The history has events after its first now: a now past one
				// of them replays, and one that is not is answered from the
				// kept replay.
				kept.Lock()
				answers := kept.holds[o.Root].timeline.answers(o.Now.Truncate(time.Second))
				kept.Unlock()
				if _, replayed := keptTimeline(all, o.Now.Truncate(time.Second), o); replayed == answers {
					t.Errorf("%s: replayed %v while the kept replay answers %v", label, replayed, answers)
				}
				if !answers {
					replays++
				}
				want := replayHoldsReference(items, all, start, o)
				sameReplay(t, label+", kept", holdsFor(items, all, start, o), want)
				sameReplay(t, label+", fresh", replayHolds(items, all, start, o), want)
			}
		}
	}
	if replays == 0 || replays > 3 {
		t.Errorf("replayed %d times over four nows, want at least one, for the events after the first now, and fewer than one a now", replays)
	}
}

// S-0608: Compute keyed to a root gives the report it gives with nothing
// kept, on the first request and on the next, which is served from the keep.
func TestComputeKeyedGivesTheReportUnkeyed(t *testing.T) {
	all, opt := keyedHistory(t)
	opt.Commits = map[string][]string{"S-0003": {"flai/internal/metrics/state.go", "elsewhere.go"}}
	plain := opt
	plain.Root, plain.ItemsGeneration, plain.CommitsKey = "", 0, ""
	want, err := json.Marshal(Compute(all, plain))
	if err != nil {
		t.Fatal(err)
	}
	for _, round := range []string{"first", "kept"} {
		got, err := json.Marshal(Compute(all, opt))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Errorf("%s request: the keyed report differs from the unkeyed one", round)
		}
	}
}
