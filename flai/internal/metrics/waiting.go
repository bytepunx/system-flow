package metrics

import (
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Waiting is how long the agents of the items completed in each ISO week
// waited on threads and in review (S-0205), and the empty wakes of the items
// completed in the window (S-0272), and the longest waits of the window on
// their own (S-0215).
type Waiting struct {
	Weeks      []WaitWeek `json:"weeks"`
	EmptyWakes EmptyWakes `json:"empty_wakes"`
	Longest    []LongWait `json:"longest"`
}

// LongWait is one wait of an item, on a thread or in review, with the
// seconds of it inside the window and who was awaited (ADR-0114).
type LongWait struct {
	Item    string `json:"item"`
	Kind    string `json:"kind"`
	Thread  string `json:"thread,omitempty"`
	Started string `json:"started"`
	// Ended is absent while the wait is open.
	Ended   string  `json:"ended,omitempty"`
	Seconds float64 `json:"seconds"`
	// Awaited is absent while the wait is open, and for a thread resolved
	// with no entry that ended its wait.
	Awaited string `json:"awaited,omitempty"`
}

// The kinds of a LongWait.
const (
	waitOnThread = "thread"
	waitInReview = "review"
)

// longestWaits is how many waits waiting.longest lists.
const longestWaits = 10

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

// wait is a thread's wait, or a review wait without a thread, and who ended
// it: for a thread, the kind of ender and the author of the entry that ended
// it, empty without one; for review, the by of the transition out of it. It
// is open while nothing ended it.
type wait struct {
	span
	by      ender
	thread  string
	awaited string
	open    bool
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
	w := wait{span: span{from, now}, thread: th.ID, open: true}
	if th.Status == "resolved" {
		w.to, _ = time.Parse(workitem.TimeFormat, th.Updated)
		w.open = false
	}
	for _, e := range entries[min(1, len(entries)):] {
		if e.Author == opener || e.Recommendation {
			continue
		}
		w.to, _ = time.Parse(workitem.TimeFormat, e.At)
		w.awaited, w.open = e.Author, false
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
// one, with the thread waits of each item by its canonical ID; sums their
// empty wakes, over the window and per week; and lists the longest waits.
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
	out.Longest = longest(items, waits, start, now)
	return out
}

// longest is the waits of the items not cancelled, on their threads and in
// review, that have a part in the window, the most seconds first, at most
// longestWaits of them (ADR-0114). The window's ends are taken to the whole
// second, as the report prints them, so that the seconds are whole.
func longest(items []*workitem.Item, waits map[string][]wait, start, now time.Time) []LongWait {
	start, now = start.Truncate(time.Second), now.Truncate(time.Second)
	out := []LongWait{}
	add := func(lw LongWait, s span, open bool) {
		if lw.Seconds <= 0 {
			return
		}
		lw.Started = s.from.Format(workitem.TimeFormat)
		if !open {
			lw.Ended = s.to.Format(workitem.TimeFormat)
		}
		out = append(out, lw)
	}
	for _, it := range items {
		if it.Status == workitem.Cancelled {
			continue
		}
		for _, w := range waits[workitem.CanonicalID(it.ID)] {
			in := inWindow(w.span, start, now)
			add(LongWait{Item: it.ID, Kind: waitOnThread, Thread: w.thread, Awaited: w.awaited,
				Seconds: inProgressSeconds(it, []span{in}, now)}, w.span, w.open)
		}
		for _, r := range inReview(it, now) {
			in := inWindow(r.span, start, now)
			add(LongWait{Item: it.ID, Kind: waitInReview, Awaited: r.awaited,
				Seconds: max(0, in.to.Sub(in.from).Seconds())}, r.span, r.open)
		}
	}
	sort.Slice(out, func(i, j int) bool { return longerWait(out[i], out[j]) })
	return out[:min(len(out), longestWaits)]
}

// longerWait reports whether a is listed before b: the more seconds first,
// then the earlier start, then the item and the thread in ID order, a
// thread's wait before a review wait.
func longerWait(a, b LongWait) bool {
	if a.Seconds != b.Seconds {
		return a.Seconds > b.Seconds
	}
	if a.Started != b.Started {
		return a.Started < b.Started
	}
	if ia, ib := workitem.CanonicalID(a.Item), workitem.CanonicalID(b.Item); ia != ib {
		return ia < ib
	}
	if (a.Thread == "") != (b.Thread == "") {
		return b.Thread == ""
	}
	return threads.CanonicalID(a.Thread) < threads.CanonicalID(b.Thread)
}

// inWindow is the part of the span from start to now; it ends before it
// starts when the span has no such part.
func inWindow(s span, start, now time.Time) span {
	return span{laterOf(s.from, start), earlierOf(s.to, now)}
}

// inReview is each interval from a transition to review to the next
// transition, ended by that transition's by, or open to now while review is
// the current state.
func inReview(it *workitem.Item, now time.Time) []wait {
	var out []wait
	for i, tr := range it.Transitions {
		if tr.To != workitem.Review {
			continue
		}
		from, _ := time.Parse(workitem.TimeFormat, tr.At)
		w := wait{span: span{from, now}, open: true}
		if i+1 < len(it.Transitions) {
			next := it.Transitions[i+1]
			w.to, _ = time.Parse(workitem.TimeFormat, next.At)
			w.awaited, w.open = next.By, false
		}
		out = append(out, w)
	}
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
