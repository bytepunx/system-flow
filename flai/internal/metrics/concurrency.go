package metrics

import (
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Concurrency is how long, over the window, exactly n items of the report's
// type were in progress at once, for each n from 0 (S-0551, ADR-0231).
type Concurrency struct {
	WindowSeconds int `json:"window_seconds"`
	// Limit is the board's in-progress limit today, 0 without one.
	Limit int `json:"limit,omitempty"`
	// Seconds and Share run from 0 through the larger of Limit and the most
	// in progress at once in the window.
	Seconds []int     `json:"seconds"`
	Share   []float64 `json:"share"`
	// MeanBusy is the mean count over the seconds with one or more in
	// progress, nil when none was.
	MeanBusy *float64 `json:"mean_busy,omitempty"`
}

// ReadyWait is the queue time of the items summary.queue_time counts at p50
// and p85, under the name the board uses (S-0551, ADR-0231). The percentiles
// are nil when it counts none.
type ReadyWait struct {
	Count int      `json:"count"`
	P50   *float64 `json:"p50_seconds,omitempty"`
	P85   *float64 `json:"p85_seconds,omitempty"`
}

// concurrency counts, over the window from start to now in whole seconds, the
// seconds with exactly n of the items in in-progress, from each item's
// in-progress intervals cut to the window. An interval runs from its
// transition up to, not including, the next.
func concurrency(items []*workitem.Item, start, now time.Time, limit int) Concurrency {
	from, to := start.Truncate(time.Second), now.Truncate(time.Second)
	window := int(to.Unix() - from.Unix())
	delta := map[int64]int{}
	for _, it := range items {
		for _, s := range inProgress(it, now) {
			s = inWindow(s, from, to)
			if !s.to.After(s.from) {
				continue
			}
			delta[s.from.Unix()]++
			delta[s.to.Unix()]--
		}
	}
	at := make([]int64, 0, len(delta))
	for t := range delta {
		at = append(at, t)
	}
	sort.Slice(at, func(i, j int) bool { return at[i] < at[j] })
	seconds := make([]int, limit+1)
	count, prev := 0, from.Unix()
	for _, t := range at {
		seconds[count] += int(t - prev)
		count, prev = count+delta[t], t
		for len(seconds) <= count {
			seconds = append(seconds, 0)
		}
	}
	seconds[count] += int(to.Unix() - prev)
	out := Concurrency{WindowSeconds: window, Limit: limit, Seconds: seconds, Share: make([]float64, len(seconds))}
	busy, weighted := 0, 0
	for n, s := range seconds {
		if window > 0 {
			out.Share[n] = float64(s) / float64(window)
		}
		if n > 0 {
			busy += s
			weighted += n * s
		}
	}
	if busy > 0 {
		m := float64(weighted) / float64(busy)
		out.MeanBusy = &m
	}
	return out
}

// readyWait is the queue time distribution's count, p50, and p85.
func readyWait(queue Distribution) ReadyWait {
	out := ReadyWait{Count: queue.Count}
	if queue.Count > 0 {
		p50, p85 := queue.P50, queue.P85
		out.P50, out.P85 = &p50, &p85
	}
	return out
}
