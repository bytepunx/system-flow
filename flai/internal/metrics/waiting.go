package metrics

import (
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Waiting is how long the agents of the items completed in each ISO week
// waited on threads and in review (S-0205), and the empty wakes of the items
// completed in the window (S-0272).
type Waiting struct {
	Weeks      []WaitWeek `json:"weeks"`
	EmptyWakes EmptyWakes `json:"empty_wakes"`
}

// EmptyWakes is the sum of the empty wakes of the items completed in the
// window, and its mean over those of them that carry agents' usage, absent
// when none does (ADR-0105).
type EmptyWakes struct {
	Count int      `json:"count"`
	Mean  *float64 `json:"mean,omitempty"`
}

// WaitWeek is the waiting of the items completed in one ISO week.
type WaitWeek struct {
	Week  string `json:"week"`
	Start string `json:"start"`
	Items int    `json:"items"`
	// EmptyWakes is the sum of the items' empty wakes, 0 when none.
	EmptyWakes int             `json:"empty_wakes"`
	Threads    ThreadWaitTotal `json:"threads"`
	Review     WaitTotal       `json:"review"`
}

// WaitTotal is the seconds waited over a week's items and their mean, absent
// for a week without items.
type WaitTotal struct {
	Total float64  `json:"total_seconds"`
	Mean  *float64 `json:"mean_seconds,omitempty"`
}

// ThreadWaitTotal is the seconds waited on threads over a week's items, and
// the waits among them the orchestrator ended and those the operator's
// confirmation of a recommendation ended (S-0220).
type ThreadWaitTotal struct {
	WaitTotal
	Orchestrator WaitCount `json:"orchestrator"`
	Confirmed    WaitCount `json:"confirmed"`
}

// WaitCount is how many of a week's thread waits one kind of entry ended,
// and the seconds of their part in the items' thread waits.
type WaitCount struct {
	Count int     `json:"count"`
	Total float64 `json:"total_seconds"`
}

// span is the time from one moment up to another.
type span struct{ from, to time.Time }

// ender is who ended a thread's wait (S-0220); the zero value is anyone but
// the orchestrator without confirming a recommendation, or nobody yet.
type ender int

const (
	// byOrchestrator is an entry of the orchestrator's.
	byOrchestrator ender = iota + 1
	// byConfirmation is an entry confirming a recommendation (ADR-0090).
	byConfirmation
)

// wait is a thread's wait and who ended it.
type wait struct {
	span
	by ender
}

// deriveWaitReview sets the seconds the item spent in review, from the time
// in state already derived, when it was ever there.
func deriveWaitReview(m *ItemMetrics, it *workitem.Item) {
	if it.FirstAt(workitem.Review).IsZero() {
		return
	}
	r := m.InState[workitem.Review]
	m.WaitReview = &r
}

// threadWaits is the waits of the threads anchored to each item, or to one of
// its tasks, by the item's canonical ID.
func threadWaits(all []*workitem.Item, ths []*threads.Thread, now time.Time) map[string][]wait {
	parent := map[string]string{}
	for _, it := range all {
		if it.Type == workitem.Task && it.Parent != "" {
			parent[workitem.CanonicalID(it.ID)] = workitem.CanonicalID(it.Parent)
		}
	}
	out := map[string][]wait{}
	for _, th := range ths {
		if th.Anchor.Item == "" {
			continue
		}
		w, ok := threadWait(th, now)
		if !ok {
			continue
		}
		id := workitem.CanonicalID(th.Anchor.Item)
		out[id] = append(out[id], w)
		if p := parent[id]; p != "" {
			out[p] = append(out[p], w)
		}
	}
	return out
}

// threadWait is the time from the thread's first entry, or its creation
// without entries, to the first later entry by another author that is not a
// recommendation; without one, to updated once resolved and to now while
// open. It is ended by the orchestrator when that entry is the
// orchestrator's, and by a confirmation when it confirms a recommendation.
// It is false when the thread's start cannot be read.
func threadWait(th *threads.Thread, now time.Time) (wait, bool) {
	entries := th.Entries()
	at, opener := th.Created, ""
	if len(entries) > 0 {
		at, opener = entries[0].At, entries[0].Author
	}
	from, err := time.Parse(workitem.TimeFormat, at)
	if err != nil {
		return wait{}, false
	}
	w := wait{span: span{from, now}}
	if th.Status == "resolved" {
		w.to, _ = time.Parse(workitem.TimeFormat, th.Updated)
	}
	for _, e := range entries[min(1, len(entries)):] {
		if e.Author == opener || e.Recommendation {
			continue
		}
		w.to, _ = time.Parse(workitem.TimeFormat, e.At)
		switch {
		case workitem.IsOrchestrator(e.Author):
			w.by = byOrchestrator
		case e.Confirms():
			w.by = byConfirmation
		}
		break
	}
	return w, true
}

// deriveWaitThreads sets the seconds the item's agent waited on its threads
// while it was in progress, and the part of them the orchestrator ended,
// both absent without threads.
func deriveWaitThreads(m *ItemMetrics, it *workitem.Item, waits []wait, now time.Time) {
	if waits == nil {
		return
	}
	all := make([]span, len(waits))
	for i, w := range waits {
		all[i] = w.span
	}
	total := inProgressSeconds(it, all, now)
	orchestrator := inProgressSeconds(it, endedBy(waits, byOrchestrator), now)
	m.WaitThreads, m.WaitThreadsOrchestrator = &total, &orchestrator
}

// endedBy is the spans of the waits by ended.
func endedBy(waits []wait, by ender) []span {
	var out []span
	for _, w := range waits {
		if w.by == by {
			out = append(out, w.span)
		}
	}
	return out
}

// inProgressSeconds is the seconds of the union of the spans that fall in
// the item's in-progress intervals.
func inProgressSeconds(it *workitem.Item, spans []span, now time.Time) float64 {
	var total float64
	for _, w := range union(spans) {
		for _, p := range inProgress(it, now) {
			from, to := laterOf(w.from, p.from), earlierOf(w.to, p.to)
			if to.After(from) {
				total += to.Sub(from).Seconds()
			}
		}
	}
	return total
}

// add counts the item's waits by ended that have a part in its in-progress
// intervals, and adds the seconds of the union of those waits that fall in
// them.
func (c *WaitCount) add(it *workitem.Item, waits []wait, by ender, now time.Time) {
	spans := endedBy(waits, by)
	for _, s := range spans {
		if inProgressSeconds(it, []span{s}, now) > 0 {
			c.Count++
		}
	}
	c.Total += inProgressSeconds(it, spans, now)
}

// union merges overlapping spans into disjoint ones, in order.
func union(spans []span) []span {
	sorted := append([]span{}, spans...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].from.Before(sorted[j].from) })
	var out []span
	for _, s := range sorted {
		if n := len(out); n > 0 && !s.from.After(out[n-1].to) {
			out[n-1].to = laterOf(out[n-1].to, s.to)
			continue
		}
		out = append(out, s)
	}
	return out
}

// inProgress is each interval from a transition to in-progress to the next
// transition, or to now while it is the current state.
func inProgress(it *workitem.Item, now time.Time) []span {
	var out []span
	for i, tr := range it.Transitions {
		if tr.To != workitem.InProgress {
			continue
		}
		from, _ := time.Parse(workitem.TimeFormat, tr.At)
		to := now
		if i+1 < len(it.Transitions) {
			to, _ = time.Parse(workitem.TimeFormat, it.Transitions[i+1].At)
		}
		out = append(out, span{from, to})
	}
	return out
}

// waiting lays out the waits of the items done in the window by the ISO week
// they were completed in, from the one that holds the window's start to this
// one, with the thread waits of each item by its canonical ID; and sums their
// empty wakes, over the window and per week.
func waiting(items []*workitem.Item, per map[string]ItemMetrics, waits map[string][]wait, start, now time.Time) Waiting {
	out := Waiting{Weeks: []WaitWeek{}}
	measured := 0 // the items that carry agents' usage
	at := map[string]int{}
	for m := monday(start); !m.After(now); m = m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		at[weekKey(y, w)] = len(out.Weeks)
		out.Weeks = append(out.Weeks, WaitWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02")})
	}
	for _, it := range items {
		done := it.CompletedAt()
		if it.Status != workitem.Done || done.Before(start) || done.After(now) {
			continue
		}
		y, w := done.ISOWeek()
		b := &out.Weeks[at[weekKey(y, w)]]
		b.Items++
		if it.Usage != nil {
			b.EmptyWakes += it.Usage.EmptyWakes
		}
		if agentsSpent(it.Usage) {
			measured++
		}
		m := per[it.ID]
		b.Threads.Total += orZero(m.WaitThreads)
		b.Review.Total += orZero(m.WaitReview)
		ws := waits[workitem.CanonicalID(it.ID)]
		b.Threads.Orchestrator.add(it, ws, byOrchestrator, now)
		b.Threads.Confirmed.add(it, ws, byConfirmation, now)
	}
	for i := range out.Weeks {
		b := &out.Weeks[i]
		if b.Items > 0 {
			t, r := b.Threads.Total/float64(b.Items), b.Review.Total/float64(b.Items)
			b.Threads.Mean, b.Review.Mean = &t, &r
		}
		out.EmptyWakes.Count += b.EmptyWakes
	}
	out.EmptyWakes.Mean = over(float64(out.EmptyWakes.Count), float64(measured))
	return out
}

func orZero(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func laterOf(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func earlierOf(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
