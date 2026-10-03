package metrics

import (
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Waiting is how long the agents of the items completed in each ISO week
// waited on threads and in review (S-0205).
type Waiting struct {
	Weeks []WaitWeek `json:"weeks"`
}

// WaitWeek is the waiting of the items completed in one ISO week.
type WaitWeek struct {
	Week    string    `json:"week"`
	Start   string    `json:"start"`
	Items   int       `json:"items"`
	Threads WaitTotal `json:"threads"`
	Review  WaitTotal `json:"review"`
}

// WaitTotal is the seconds waited over a week's items and their mean, absent
// for a week without items.
type WaitTotal struct {
	Total float64  `json:"total_seconds"`
	Mean  *float64 `json:"mean_seconds,omitempty"`
}

// span is the time from one moment up to another.
type span struct{ from, to time.Time }

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
func threadWaits(all []*workitem.Item, ths []*threads.Thread, now time.Time) map[string][]span {
	parent := map[string]string{}
	for _, it := range all {
		if it.Type == workitem.Task && it.Parent != "" {
			parent[workitem.CanonicalID(it.ID)] = workitem.CanonicalID(it.Parent)
		}
	}
	out := map[string][]span{}
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
// without entries, to the first later entry by another author; without one,
// to updated once resolved and to now while open. It is false when the
// thread's start cannot be read.
func threadWait(th *threads.Thread, now time.Time) (span, bool) {
	entries := th.Entries()
	at, opener := th.Created, ""
	if len(entries) > 0 {
		at, opener = entries[0].At, entries[0].Author
	}
	from, err := time.Parse(workitem.TimeFormat, at)
	if err != nil {
		return span{}, false
	}
	to := now
	if th.Status == "resolved" {
		to, _ = time.Parse(workitem.TimeFormat, th.Updated)
	}
	for _, e := range entries[min(1, len(entries)):] {
		if e.Author != opener {
			to, _ = time.Parse(workitem.TimeFormat, e.At)
			break
		}
	}
	return span{from, to}, true
}

// waitInProgress is the seconds of the union of the waits that fall in the
// item's in-progress intervals, nil without waits.
func waitInProgress(it *workitem.Item, waits []span, now time.Time) *float64 {
	if waits == nil {
		return nil
	}
	var total float64
	for _, w := range union(waits) {
		for _, p := range inProgress(it, now) {
			from, to := laterOf(w.from, p.from), earlierOf(w.to, p.to)
			if to.After(from) {
				total += to.Sub(from).Seconds()
			}
		}
	}
	return &total
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
// one.
func waiting(items []*workitem.Item, per map[string]ItemMetrics, start, now time.Time) Waiting {
	out := Waiting{Weeks: []WaitWeek{}}
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
		m := per[it.ID]
		b.Threads.Total += orZero(m.WaitThreads)
		b.Review.Total += orZero(m.WaitReview)
	}
	for i := range out.Weeks {
		b := &out.Weeks[i]
		if b.Items > 0 {
			t, r := b.Threads.Total/float64(b.Items), b.Review.Total/float64(b.Items)
			b.Threads.Mean, b.Review.Mean = &t, &r
		}
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
