package metrics

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// historyStories is about as many stories as this repository has (I-0155),
// each with four tasks on average.
const historyStories = 500

// historySpan is how long the history runs, from its first story's creation
// to its last event.
const historySpan = 120 * 24 * time.Hour

var (
	historyPackages = []string{
		"metrics", "workitem", "serve", "mcp", "manifest", "harness", "mergequeue", "preview",
		"storygit", "threads", "messages", "forge", "issues", "usage", "prime", "check",
		"release", "plan", "orchestrate", "board", "stream", "archive", "docs", "analyze",
	}
	historyFiles      = []string{"state", "hold", "claims", "events", "render", "load", "parse", "api", "run", "doc"}
	historyDocs       = []string{"workflow.md", "flai-cli.md", "metrics.md", "work-hierarchy.md", "flaiover-dashboard.md", "strategic-agents.md", "overview.md", "planning.md"}
	historyComponents = []string{"Board", "Stats", "Claims", "Story", "Inbox", "MergeQueue", "Threads", "HostPanel", "Settings", "Review"}
	historyProjects   = []manifest.Project{
		{Name: "flai", Path: "flai", Kind: "go", Tags: []string{"cli"}},
		{Name: "flaiover", Path: "flaiover", Kind: "sveltekit", Tags: []string{"dashboard"}},
		{Name: "template", Path: "template", Kind: "template", Tags: []string{"template", "conventions"}},
	}
	historyShared = manifest.Claims{Shared: []string{
		"docs/users/flai.md", "docs/users/flai-reference.md", "design/system/flai-cli.md",
		"design/adrs", "design/issues", "template/CHANGELOG.md",
	}}
)

// claimsHistory is a history of n stories and their tasks, the same on every
// call: stories move from backlog through ready and in progress to review and
// done, some cancelled and some sent back, several in progress at a time, with
// touches that overlap, some naming earlier stories in after; their tasks'
// file touches narrow their stories' folders. It returns every item, the
// Options Compute is given, and the window's start as Compute derives it.
func claimsHistory(n int) ([]*workitem.Item, Options, time.Time) {
	rng := rand.New(rand.NewPCG(4155, 471))
	t0 := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	cutoff := t0.Add(historySpan)
	step := historySpan / time.Duration(n)
	// wait is a random duration of mean d, whole seconds, at least a minute.
	wait := func(d time.Duration) time.Duration {
		return max(time.Duration(rng.ExpFloat64()*float64(d)).Truncate(time.Second), time.Minute)
	}
	chance := func(p float64) bool { return rng.Float64() < p }
	pick := func(xs []string) string { return xs[int(rng.ExpFloat64()*float64(len(xs))/3)%len(xs)] }
	file := func(pkg string) string { return "flai/internal/" + pkg + "/" + pick(historyFiles) + ".go" }
	touches := func() []string {
		if chance(0.02) {
			return nil
		}
		var out []string
		for range 1 + rng.IntN(4) {
			pkg := pick(historyPackages)
			if chance(0.2) {
				out = append(out, "flai/internal/"+pkg)
				continue
			}
			for range 1 + rng.IntN(4) {
				f := file(pkg)
				out = append(out, f, f[:len(f)-3]+"_test.go")
			}
		}
		for range rng.IntN(3) {
			out = append(out, "design/system/"+pick(historyDocs))
		}
		if chance(0.4) {
			out = append(out, "docs/users/flai.md")
		}
		if chance(0.3) {
			out = append(out, "design/adrs")
		}
		if chance(0.2) {
			out = append(out, "design/issues/summary.md")
		}
		if chance(0.25) {
			for range 1 + rng.IntN(2) {
				out = append(out, "flaiover/src/lib/components/"+pick(historyComponents)+".svelte")
			}
		}
		for _, p := range []struct {
			name string
			p    float64
		}{{"flai", 0.03}, {"dashboard", 0.02}, {"template", 0.01}} {
			if chance(p.p) {
				out = append(out, p.name)
			}
		}
		return out
	}
	// narrow is a task's touches inside its story's: a file inside a folder,
	// a file as it is, and now and then a file outside them.
	narrow := func(story []string) []string {
		var out []string
		for range 1 + rng.IntN(3) {
			if len(story) == 0 {
				out = append(out, file(pick(historyPackages)))
				continue
			}
			s := story[rng.IntN(len(story))]
			switch {
			case s == "design/adrs":
				out = append(out, fmt.Sprintf("design/adrs/ADR-%04d.md", rng.IntN(200)))
			case len(s) > len("flai/internal/") && s[:len("flai/internal/")] == "flai/internal/" && s[len(s)-3:] != ".go":
				out = append(out, s+"/"+pick(historyFiles)+".go")
			case s == "flai" || s == "dashboard" || s == "template":
				out = append(out, file(pick(historyPackages)))
			default:
				out = append(out, s)
			}
		}
		if chance(0.15) {
			out = append(out, file(pick(historyPackages)))
		}
		return out
	}
	var stories, tasks []*workitem.Item
	for i := range n {
		s := &workitem.Item{ID: fmt.Sprintf("S-%04d", i+1), Type: workitem.Story, Nature: "feature", Status: workitem.Backlog, Touches: touches()}
		at := t0.Add(time.Duration(i)*step + wait(step))
		s.Created = at.Format(workitem.TimeFormat)
		created := at
		// move appends a transition at d after the last one, unless it falls
		// after the history's end.
		move := func(it *workitem.Item, from *time.Time, to string, d time.Duration) {
			*from = from.Add(d)
			if from.After(cutoff) {
				return
			}
			it.Transitions = append(it.Transitions, workitem.Transition{To: to, At: from.Format(workitem.TimeFormat)})
			it.Status = to
		}
		if i > 0 && chance(0.25) {
			for range 1 + rng.IntN(2) {
				s.After = append(s.After, fmt.Sprintf("S-%04d", i-rng.IntN(min(i, 15))))
			}
		}
		move(s, &at, workitem.Ready, wait(17*step))
		var started, stopped time.Time
		switch {
		case chance(0.03):
			move(s, &at, workitem.Cancelled, wait(2*step))
			stopped = at
		default:
			move(s, &at, workitem.InProgress, wait(2*step))
			if chance(0.05) {
				move(s, &at, workitem.Ready, wait(2*step))
				move(s, &at, workitem.InProgress, wait(step))
			}
			started = at
			if chance(0.03) {
				move(s, &at, workitem.Cancelled, wait(2*step))
				stopped = at
				break
			}
			move(s, &at, workitem.Review, wait(3*step/2))
			stopped = at
			if chance(0.08) {
				move(s, &at, workitem.InProgress, wait(time.Hour))
				move(s, &at, workitem.Review, wait(2*time.Hour))
			}
			move(s, &at, workitem.Done, wait(step))
		}
		s.Archived = s.Status == workitem.Done && at.Before(t0.Add(historySpan/3)) && i%3 == 0
		stories = append(stories, s)
		count := 1 + rng.IntN(7)
		for k := range count {
			tk := &workitem.Item{ID: fmt.Sprintf("T-%04d", len(tasks)+1), Type: workitem.Task, Nature: "feature", Status: workitem.Backlog,
				Parent: s.ID, Touches: narrow(s.Touches), Archived: s.Archived}
			tat := created.Add(wait(6 * time.Hour))
			if !started.IsZero() && chance(0.25) {
				tat = started.Add(wait(step))
			}
			tk.Created = tat.Format(workitem.TimeFormat)
			move(tk, &tat, workitem.Ready, wait(10*time.Minute))
			switch {
			case chance(0.06):
				move(tk, &tat, workitem.Cancelled, wait(step))
			case started.IsZero():
				if !stopped.IsZero() {
					move(tk, &tat, workitem.Cancelled, max(stopped.Sub(tat), 0)+time.Minute)
				}
			default:
				slice := max(stopped.Sub(started)/time.Duration(count), 2*time.Minute)
				begin := started.Add(time.Duration(k) * slice)
				move(tk, &tat, workitem.InProgress, max(begin.Sub(tat), 0)+time.Minute)
				move(tk, &tat, workitem.Done, wait(slice/2))
			}
			tasks = append(tasks, tk)
		}
	}
	opt := Options{Now: cutoff.Add(time.Hour), Since: 30 * 24 * time.Hour, Projects: historyProjects, Shared: historyShared, WIPLimit: 5}
	return append(stories, tasks...), opt, opt.Now.Add(-opt.Since)
}

// storiesOf is the items Compute aggregates by default: the stories.
func storiesOf(all []*workitem.Item) []*workitem.Item {
	var out []*workitem.Item
	for _, it := range all {
		if it.Type == workitem.Story {
			out = append(out, it)
		}
	}
	return out
}

// replayCounts is, over the replay of all up to now, the distinct times of
// the events, the times a Holds is built at, and the times an event changes
// what a hold reads: a story created, or entering or leaving ready, in
// progress, or done; or a task of a story in ready or in progress created, or
// entering or leaving done or cancelled.
func replayCounts(items, all []*workitem.Item, now time.Time) (times, built, changed int) {
	tracked := map[string]bool{}
	for _, it := range items {
		if it.Type == workitem.Story && !it.FirstAt(workitem.Ready).IsZero() {
			tracked[it.ID] = true
		}
	}
	type event struct {
		at time.Time
		it *workitem.Item
		to string
	}
	var events []event
	status := map[string]string{}
	for _, it := range all {
		created, _ := time.Parse(workitem.TimeFormat, it.Created)
		events = append(events, event{at: created, it: it})
		for _, tr := range it.Transitions {
			at, _ := time.Parse(workitem.TimeFormat, tr.At)
			events = append(events, event{at: at, it: it, to: tr.To})
		}
	}
	sort.SliceStable(events, func(a, b int) bool { return events[a].at.Before(events[b].at) })
	reads := map[string]bool{workitem.Ready: true, workitem.InProgress: true, workitem.Done: true}
	closed := func(s string) bool { return s == workitem.Done || s == workitem.Cancelled }
	ready := 0
	now = now.Truncate(time.Second)
	for k := 0; k < len(events) && !events[k].at.After(now); {
		at := events[k].at
		change := false
		for ; k < len(events) && events[k].at.Equal(at); k++ {
			e := events[k]
			from, to := status[e.it.ID], e.to
			if to == "" {
				to = workitem.Backlog
			}
			status[e.it.ID] = to
			if tracked[e.it.ID] {
				if from == workitem.Ready {
					ready--
				}
				if to == workitem.Ready {
					ready++
				}
			}
			switch {
			case e.it.Type == workitem.Story:
				change = change || e.to == "" || reads[from] || reads[to]
			case status[e.it.Parent] == workitem.Ready || status[e.it.Parent] == workitem.InProgress:
				change = change || e.to == "" || closed(from) != closed(to) || (from == workitem.Cancelled) != (to == workitem.Cancelled)
			}
		}
		times++
		if ready > 0 {
			built++
		}
		if change {
			changed++
		}
	}
	return times, built, changed
}

// BenchmarkReplayHolds replays the hold rules over a history at this
// repository's scale (I-0155).
func BenchmarkReplayHolds(b *testing.B) {
	all, opt, start := claimsHistory(historyStories)
	items := storiesOf(all)
	times, built, changed := replayCounts(items, all, opt.Now)
	for b.Loop() {
		replayHolds(items, all, start, opt)
	}
	b.ReportMetric(float64(len(all)), "items/op")
	b.ReportMetric(float64(times), "times/op")
	b.ReportMetric(float64(built), "holds/op")
	b.ReportMetric(float64(changed), "changes/op")
}

// BenchmarkComputeClaims computes the report, the replay of the holds among
// it, over a history at this repository's scale (I-0155).
func BenchmarkComputeClaims(b *testing.B) {
	all, opt, _ := claimsHistory(historyStories)
	for b.Loop() {
		Compute(all, opt)
	}
}

// S-0471: the reference replay gives what replayHolds gives on a history of
// this repository's shape, smaller in short mode.
func TestReplayHoldsMatchesTheReference(t *testing.T) {
	n := historyStories
	if testing.Short() {
		n = 60
	}
	all, opt, start := claimsHistory(n)
	items := storiesOf(all)
	got, want := replayHolds(items, all, start, opt), replayHoldsReference(items, all, start, opt)
	held := 0
	for id, w := range want.seconds {
		g, ok := got.seconds[id]
		if !ok || (g == nil) != (w == nil) || (g != nil && *g != *w) {
			t.Errorf("%s held = %v, want %v", id, deref(g), deref(w))
		}
		if w != nil && *w > 0 {
			held++
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
	if held == 0 {
		t.Errorf("no story held in the history, want some")
	}
}

// replayHoldsReference is replayHolds as it was before S-0471, every Holds
// built anew at each event time: what a faster replay must give.
func replayHoldsReference(items, all []*workitem.Item, start time.Time, opt Options) holdReplay {
	// Held time is whole seconds, as timestamps are: a hold still open runs
	// to now's last whole second.
	now := opt.Now.Truncate(time.Second)
	out := holdReplay{seconds: map[string]*float64{}}
	day0, week0 := firstDay(start), monday(start)
	for d := day0; !d.After(now); d = d.AddDate(0, 0, 1) {
		out.days = append(out.days, 0)
	}
	for m := week0; !m.After(now); m = m.AddDate(0, 0, 7) {
		out.weeks = append(out.weeks, HeldByCode{})
	}
	for _, it := range items {
		if it.Type == workitem.Story && !it.FirstAt(workitem.Ready).IsZero() {
			out.seconds[it.ID] = new(float64)
		}
	}
	if len(out.seconds) == 0 {
		return out
	}
	type event struct {
		at time.Time
		i  int    // the item's copy
		to string // the state entered, "" for its creation
	}
	var copies []workitem.Item
	var events []event
	for _, it := range all {
		if it.Type != workitem.Story && it.Type != workitem.Task {
			continue
		}
		c := *it
		c.Status, c.Archived = workitem.Backlog, false
		i := len(copies)
		copies = append(copies, c)
		created, _ := time.Parse(workitem.TimeFormat, it.Created)
		events = append(events, event{at: created, i: i})
		for _, tr := range it.Transitions {
			at, _ := time.Parse(workitem.TimeFormat, tr.At)
			events = append(events, event{at: at, i: i, to: tr.To})
		}
	}
	sort.SliceStable(events, func(a, b int) bool { return events[a].at.Before(events[b].at) })
	// dayEnd is the moment a day's held stories are counted at: its last
	// second, today's at now.
	dayEnd := func(i int) time.Time {
		end := day0.AddDate(0, 0, i).Add(24*time.Hour - time.Second)
		if end.After(now) {
			return now
		}
		return end
	}
	var snapshot []*workitem.Item
	day := 0
	for k := 0; k < len(events) && !events[k].at.After(now); {
		at := events[k].at
		for ; k < len(events) && events[k].at.Equal(at); k++ {
			e := events[k]
			if e.to == "" {
				snapshot = append(snapshot, &copies[e.i])
			} else {
				copies[e.i].Status = e.to
			}
		}
		// The states hold from at up to next; until is next, not after now.
		until, next := now, now.Add(time.Second)
		if k < len(events) && !events[k].at.After(now) {
			until, next = events[k].at, events[k].at
		}
		byCode := map[string]int{}
		var h *workitem.Holds
		for _, it := range snapshot {
			if out.seconds[it.ID] == nil || it.Status != workitem.Ready {
				continue
			}
			if h == nil {
				h = workitem.NewHolds(snapshot, opt.Projects).WithShared(opt.Shared)
			}
			if hold := h.Of(it); hold != nil {
				*out.seconds[it.ID] += until.Sub(at).Seconds()
				byCode[hold.Code]++
			}
		}
		for code, n := range byCode {
			for i := range out.weeks {
				wk := week0.AddDate(0, 0, 7*i)
				from, to := laterOf(at, wk), earlierOf(until, wk.AddDate(0, 0, 7))
				if to.After(from) {
					out.weeks[i].add(code, float64(n)*to.Sub(from).Seconds())
				}
			}
		}
		for day < len(out.days) && dayEnd(day).Before(at) {
			day++ // counted under earlier states, or before any item was
		}
		for ; day < len(out.days) && dayEnd(day).Before(next); day++ {
			for _, n := range byCode {
				out.days[day] += n
			}
		}
	}
	return out
}
