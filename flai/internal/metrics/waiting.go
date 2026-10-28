package metrics

import (
	"maps"
	"slices"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Waiting is how long the agents of the items completed in each ISO week
// waited on threads and in review (S-0205), and in review on checks and
// reviewers (ADR-0168), and the empty wakes of the items
// completed in the window (S-0272), and the longest waits of the window on
// their own (S-0215).
type Waiting struct {
	Weeks      []WaitWeek `json:"weeks"`
	EmptyWakes EmptyWakes `json:"empty_wakes"`
	Longest    []LongWait `json:"longest"`
}

// LongWait is one wait of an item, on a thread, in review, or on its pull
// request's checks or reviewers, with the seconds of it inside the window and
// who was awaited (ADR-0114, ADR-0168).
type LongWait struct {
	Item    string `json:"item"`
	Kind    string `json:"kind"`
	Thread  string `json:"thread,omitempty"`
	Started string `json:"started"`
	// Ended is absent while the wait is open.
	Ended   string  `json:"ended,omitempty"`
	Seconds float64 `json:"seconds"`
	// Awaited is absent while the wait is open, for a thread resolved with
	// no entry that ended its wait, for a checks wait whose checks had
	// finished when its head was first read, and for a wait a push ended.
	Awaited string `json:"awaited,omitempty"`
}

// The kinds of a LongWait.
const (
	waitOnThread    = "thread"
	waitInReview    = "review"
	waitOnChecks    = "checks"
	waitOnReviewers = "reviewers"
)

// waitKinds is the order of the kinds of LongWait among waits otherwise tied.
var waitKinds = map[string]int{waitOnThread: 0, waitInReview: 1, waitOnChecks: 2, waitOnReviewers: 3}

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
	// Checks and Reviewers are the parts of Review waited on the pull
	// requests' required checks and on their reviewers, an item without kept
	// forge state counting 0 (ADR-0168).
	Checks    WaitTotal `json:"checks"`
	Reviewers WaitTotal `json:"reviewers"`
}

// WaitTotal is the seconds waited over a week's items and their mean, absent
// for a week without items.
type WaitTotal struct {
	Total float64  `json:"total_seconds"`
	Mean  *float64 `json:"mean_seconds,omitempty"`
}

// ThreadWaitTotal is the seconds waited on threads over a week's items, and
// the waits among them the orchestrator ended, those the operator's
// confirmation of a recommendation ended (S-0220), and those flai serve's
// confirmation by timeout ended (ADR-0206).
type ThreadWaitTotal struct {
	WaitTotal
	Orchestrator WaitCount `json:"orchestrator"`
	Confirmed    WaitCount `json:"confirmed"`
	Timeout      WaitCount `json:"timeout"`
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
	// byConfirmation is an entry confirming a recommendation (ADR-0090), the
	// operator's.
	byConfirmation
	// byTimeout is flai serve's confirmation of a recommendation the operator
	// left unanswered past confirm_after (ADR-0206).
	byTimeout
)

// wait is a thread's wait, or a review, checks, or reviewers wait without a
// thread, and who ended it: for a thread, the kind of ender and the author of
// the entry that ended it, empty without one; for review, the by of the
// transition out of it; for checks, the check that finished last; for
// reviewers, the reviewer whose review ended it. It is open while nothing
// ended it.
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
// orchestrator's, by a timeout when it is flai serve's confirmation in the
// operator's absence, and by a confirmation when it confirms a recommendation
// otherwise. It is false when the thread's start cannot be read.
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
		case e.ConfirmsByTimeout():
			w.by = byTimeout
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
	return within(spans, inProgress(it, now))
}

// within is the seconds of the union of the spans that fall in the
// intervals, which do not overlap.
func within(spans, intervals []span) float64 {
	var total float64
	for _, w := range union(spans) {
		for _, p := range intervals {
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
// one, with the thread waits and the pull request waits of each item by its
// canonical ID; sums their empty wakes, over the window and per week; and
// lists the longest waits.
func waiting(items []*workitem.Item, per map[string]ItemMetrics, waits map[string][]wait, pulls map[string]forgeWait, start, now time.Time) Waiting {
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
		b.Checks.Total += orZero(m.WaitChecks)
		b.Reviewers.Total += orZero(m.WaitReviewers)
		ws := waits[workitem.CanonicalID(it.ID)]
		b.Threads.Orchestrator.add(it, ws, byOrchestrator, now)
		b.Threads.Confirmed.add(it, ws, byConfirmation, now)
		b.Threads.Timeout.add(it, ws, byTimeout, now)
	}
	for i := range out.Weeks {
		b := &out.Weeks[i]
		if b.Items > 0 {
			n := float64(b.Items)
			t, r, c, v := b.Threads.Total/n, b.Review.Total/n, b.Checks.Total/n, b.Reviewers.Total/n
			b.Threads.Mean, b.Review.Mean, b.Checks.Mean, b.Reviewers.Mean = &t, &r, &c, &v
		}
		out.EmptyWakes.Count += b.EmptyWakes
	}
	out.EmptyWakes.Mean = over(float64(out.EmptyWakes.Count), float64(measured))
	out.Longest = longest(items, waits, pulls, start, now)
	return out
}

// longest is the waits of the items not cancelled, on their threads, in
// review, and on their pull requests' checks and reviewers, that have a part
// in the window, the most seconds first, at most longestWaits of them
// (ADR-0114, ADR-0168). The window's ends are taken to the whole second, as
// the report prints them, so that the seconds are whole.
func longest(items []*workitem.Item, waits map[string][]wait, pulls map[string]forgeWait, start, now time.Time) []LongWait {
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
		reviews := inReview(it, now)
		for _, r := range reviews {
			in := inWindow(r.span, start, now)
			add(LongWait{Item: it.ID, Kind: waitInReview, Awaited: r.awaited,
				Seconds: max(0, in.to.Sub(in.from).Seconds())}, r.span, r.open)
		}
		pull, ok := pulls[workitem.CanonicalID(it.ID)]
		if !ok {
			continue
		}
		review := spansOf(reviews)
		for _, k := range []struct {
			kind  string
			waits []wait
		}{{waitOnChecks, pull.checks}, {waitOnReviewers, pull.reviewers}} {
			for _, w := range k.waits {
				add(LongWait{Item: it.ID, Kind: k.kind, Awaited: w.awaited,
					Seconds: within([]span{inWindow(w.span, start, now)}, review)}, w.span, w.open)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return longerWait(out[i], out[j]) })
	return out[:min(len(out), longestWaits)]
}

// longerWait reports whether a is listed before b: the more seconds first,
// then the earlier start, then the item in ID order, then a thread's wait,
// a review wait, a checks wait, and a reviewers wait in that order, then the
// thread in ID order.
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
	if ka, kb := waitKinds[a.Kind], waitKinds[b.Kind]; ka != kb {
		return ka < kb
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

// spansOf is the spans of the waits.
func spansOf(waits []wait) []span {
	out := make([]span, len(waits))
	for i, w := range waits {
		out[i] = w.span
	}
	return out
}

// forgeWait is a story's checks waits and reviewers waits, replayed from the
// events of its kept forge state (ADR-0168).
type forgeWait struct{ checks, reviewers []wait }

// forgeWaits is the waits replayed from each story's kept forge events, by its
// canonical ID; a story is in it exactly when it is in events.
func forgeWaits(events map[string][]forge.StatusEvent, now time.Time) map[string]forgeWait {
	out := make(map[string]forgeWait, len(events))
	for id, evs := range events {
		out[workitem.CanonicalID(id)] = replayForge(evs, now)
	}
	return out
}

// deriveWaitForge sets the seconds of the item's checks waits and reviewers
// waits that fall in its review intervals, when it has kept forge state.
func deriveWaitForge(m *ItemMetrics, it *workitem.Item, pulls map[string]forgeWait, now time.Time) {
	pull, ok := pulls[workitem.CanonicalID(it.ID)]
	if !ok {
		return
	}
	review := spansOf(inReview(it, now))
	checks, reviewers := within(spansOf(pull.checks), review), within(spansOf(pull.reviewers), review)
	m.WaitChecks, m.WaitReviewers = &checks, &reviewers
}

// The waits a story's pull request can be in while its forge events replay.
const (
	onNothing = iota
	onChecks
	onReviewers
)

// replayForge replays a story's forge events read by read, the events of one
// read sharing their time, each check's state and each reviewer's latest
// review carried from read to read. Each push, an event of kind head, starts
// a checks wait, ending the open wait of the head before it with nobody
// awaited. The checks wait ends at the first read after which no required
// check is pending, awaiting the check that finished last, or nobody when
// none was pending at the push's own read. A reviewers wait starts then and
// ends at the first later review that approves or requests changes, or as it
// starts while an approval stands: a request for changes stands until its
// reviewer reviews again, which is what the wait after the fix's push waits for. A wait still open runs to now, and a
// read whose time cannot be parsed is passed over.
func replayForge(events []forge.StatusEvent, now time.Time) forgeWait {
	var out forgeWait
	checks := map[string]string{}  // each required check's state
	reviews := map[string]string{} // each reviewer's latest review's state
	phase, cur, last := onNothing, wait{}, ""
	keep := func() {
		if phase == onChecks {
			out.checks = append(out.checks, cur)
		} else {
			out.reviewers = append(out.reviewers, cur)
		}
		phase = onNothing
	}
	end := func(at time.Time, awaited string) {
		cur.to, cur.awaited, cur.open = at, awaited, false
		keep()
	}
	for i := 0; i < len(events); {
		j := i + 1
		for j < len(events) && events[j].At == events[i].At {
			j++
		}
		read := events[i:j]
		i = j
		at, err := time.Parse(time.RFC3339, read[0].At)
		if err != nil {
			continue
		}
		pushed := slices.ContainsFunc(read, func(e forge.StatusEvent) bool { return e.Kind == forge.EventHead })
		if pushed {
			if phase != onNothing {
				end(at, "")
			}
			phase, cur, last = onChecks, wait{span: span{at, now}, open: true}, ""
		}
		for _, e := range read {
			switch e.Kind {
			case forge.EventCheck:
				if checks[e.Name] == string(forge.CheckPending) && e.To != string(forge.CheckPending) {
					last = e.Name
				}
				checks[e.Name] = e.To
			case forge.EventReview:
				reviews[e.Name] = e.To
				if phase == onReviewers && decides(e.To) {
					end(at, e.Name)
				}
			}
		}
		if phase != onChecks || pending(checks) {
			continue
		}
		if pushed {
			last = ""
		}
		end(at, last)
		phase, cur = onReviewers, wait{span: span{at, now}, open: true}
		if who := standing(reviews); who != "" {
			end(at, who)
		}
	}
	if phase != onNothing {
		keep()
	}
	return out
}

// pending reports whether some check is pending.
func pending(checks map[string]string) bool {
	for _, state := range checks {
		if state == string(forge.CheckPending) {
			return true
		}
	}
	return false
}

// decides reports whether a review in state approves or requests changes.
func decides(state string) bool {
	return state == string(forge.ReviewApproved) || state == string(forge.ReviewChangesRequested)
}

// standing is the first reviewer, in login order, whose latest review
// approves, empty when none does.
func standing(reviews map[string]string) string {
	for _, who := range slices.Sorted(maps.Keys(reviews)) {
		if reviews[who] == string(forge.ReviewApproved) {
			return who
		}
	}
	return ""
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
