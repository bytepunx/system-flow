package metrics

import (
	"math"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Queue is, per ISO week of the window, the stories the merge queue landed
// and sent back with their queue wait, the batches it built with their size,
// and its bisections, from its event log (S-0404, ADR-0177).
type Queue struct {
	Weeks []QueueWeek `json:"weeks"`
}

// QueueWeek is the merge queue's events counted in one ISO week. The waits
// are those of its landed and sent_back lines, in whole seconds.
type QueueWeek struct {
	Week          string  `json:"week"`
	Start         string  `json:"start"`
	Landed        int     `json:"landed"`
	SentBack      int     `json:"sent_back"`
	WaitMean      float64 `json:"wait_mean_seconds"`
	WaitMax       int     `json:"wait_max_seconds"`
	Batches       int     `json:"batches"`
	BatchSizeMean float64 `json:"batch_size_mean"`
	BatchSizeMax  int     `json:"batch_size_max"`
	Bisections    int     `json:"bisections"`
}

// queueWait is the wait of one landed or sent_back line: its story's
// canonical ID, when it ended, and its seconds.
type queueWait struct {
	story   string
	ended   time.Time
	seconds float64
}

// mergeQueue counts, in the ISO weeks from the one that holds start to the
// one that holds now, the events of the merge queue's log, each in the week
// that holds it and only from start to now, and returns beside it the sum of
// each story's queue waits that end in the window, by its canonical ID.
func mergeQueue(events []mergequeue.Event, start, now time.Time) (*Queue, map[string]float64) {
	out := &Queue{Weeks: []QueueWeek{}}
	at := map[string]int{}
	for i, m := 0, monday(start); !m.After(now); i, m = i+1, m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		at[weekKey(y, w)] = i
		out.Weeks = append(out.Weeks, QueueWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02")})
	}
	week := func(t time.Time) int {
		t = t.UTC()
		if t.IsZero() || t.Before(start) || t.After(now) {
			return -1
		}
		y, w := t.ISOWeek()
		if i, ok := at[weekKey(y, w)]; ok {
			return i
		}
		return -1
	}
	waits := map[int][]float64{}
	sizes := map[int][]int{}
	perStory := map[string]float64{}
	for _, wt := range queueWaits(events) {
		i := week(wt.ended)
		if i < 0 {
			continue
		}
		waits[i] = append(waits[i], wt.seconds)
		perStory[wt.story] += wt.seconds
	}
	for _, e := range events {
		i := week(e.At)
		if i < 0 {
			continue
		}
		wk := &out.Weeks[i]
		switch e.Event {
		case mergequeue.EventLanded:
			wk.Landed++
		case mergequeue.EventSentBack:
			wk.SentBack++
		case mergequeue.EventBatchBuilt:
			wk.Batches++
			sizes[i] = append(sizes[i], len(e.Stories))
		case mergequeue.EventBisected:
			wk.Bisections++
		}
	}
	for i := range out.Weeks {
		wk := &out.Weeks[i]
		if w := waits[i]; len(w) > 0 {
			wk.WaitMean = mean(w)
			for _, s := range w {
				wk.WaitMax = max(wk.WaitMax, int(s))
			}
		}
		if s := sizes[i]; len(s) > 0 {
			total := 0
			for _, n := range s {
				total += n
				wk.BatchSizeMax = max(wk.BatchSizeMax, n)
			}
			wk.BatchSizeMean = float64(total) / float64(len(s))
		}
	}
	return out, perStory
}

// queueWaits is the wait of every landed and sent_back line of events that
// has its story's enqueued line before it: from the latest such line to it,
// the lines taken in the order of their times, and those of a time in the
// order of the log.
func queueWaits(events []mergequeue.Event) []queueWait {
	sorted := append([]mergequeue.Event{}, events...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].At.Before(sorted[j].At) })
	enqueued := map[string]time.Time{}
	var out []queueWait
	for _, e := range sorted {
		if len(e.Stories) == 0 {
			continue
		}
		story := workitem.CanonicalID(e.Stories[0])
		switch e.Event {
		case mergequeue.EventEnqueued:
			enqueued[story] = e.At
		case mergequeue.EventLanded, mergequeue.EventSentBack:
			from, ok := enqueued[story]
			if !ok {
				continue
			}
			out = append(out, queueWait{story: story, ended: e.At, seconds: math.Round(e.At.Sub(from).Seconds())})
		}
	}
	return out
}

// deriveQueueWait sets the sum of the item's queue waits that end in the
// window, when it has any.
func deriveQueueWait(m *ItemMetrics, waits map[string]float64) {
	if s, ok := waits[workitem.CanonicalID(m.ID)]; ok {
		m.QueueWait = &s
	}
}
