package metrics

import (
	"math"
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// CostOfDelay is what waiting for the items cost, by the day and by the ISO
// week (S-0205), and what waiting for the ready stories will cost under
// three orders (S-0213).
type CostOfDelay struct {
	Days  []CostDay  `json:"days"`
	Weeks []CostWeek `json:"weeks"`
	Order CostOrder  `json:"order"`
}

// CostDay is the cost of delay outstanding per column at the end of a day, the
// number of items in each column then without a value, and what waiting cost
// during the day.
type CostDay struct {
	Date         string             `json:"date"`
	Outstanding  map[string]float64 `json:"outstanding"`
	WithoutValue map[string]int     `json:"without_value"`
	Incurred     float64            `json:"incurred"`
}

// CostWeek is what waiting cost during one ISO week.
type CostWeek struct {
	Week     string  `json:"week"`
	Start    string  `json:"start"`
	Incurred float64 `json:"incurred"`
}

// secondsPerWeek is what a cost of delay, given per week, is spread over.
const secondsPerWeek = 7 * 24 * 60 * 60

// costColumns are the columns the cost of delay outstanding is laid out in.
var costColumns = []string{workitem.Backlog, workitem.Ready, workitem.InProgress, workitem.Review}

// deriveCostOfDelay sets the item's cost of delay and what it incurred over
// the time in backlog and ready already derived.
func deriveCostOfDelay(m *ItemMetrics, it *workitem.Item) {
	v := costValue(it)
	if v == nil {
		return
	}
	value := *v
	cost := round2(value * (m.InState[workitem.Backlog] + m.InState[workitem.Ready]) / secondsPerWeek)
	m.CostOfDelay, m.CostIncurred = &value, &cost
}

// costOfDelay lays out the cost of delay of the items with a value, by the day
// from the one that holds the window's start to today, and by the ISO week
// from the one that holds the window's start to this one. Each day also counts
// the items without a value, of a type that carries one, in each column. The
// days are swept once, each item's state moved on as the days pass its
// transitions, and a day before every item is answered without reading them
// (S-0608).
func costOfDelay(items []*workitem.Item, tl timelines, start, now time.Time) CostOfDelay {
	var valued, unvalued []*workitem.Item
	for _, it := range items {
		switch {
		case costValue(it) != nil:
			valued = append(valued, it)
		case workitem.Carries(it.Type, "cost_of_delay"):
			unvalued = append(unvalued, it)
		}
	}
	waits := make([]waitedItem, len(valued))
	for i, it := range valued {
		waits[i] = waitedItem{value: *costValue(it), spans: waitSpans(it, tl.of(it), now)}
	}
	worth, without := newDayStates(valued, tl), newDayStates(unvalued, tl)
	out := CostOfDelay{Days: []CostDay{}, Weeks: []CostWeek{}}
	for d := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC); !d.After(now); d = d.AddDate(0, 0, 1) {
		end, next := d.Add(24*time.Hour-time.Second), d.AddDate(0, 0, 1)
		p := CostDay{Date: d.Format("2006-01-02"), Outstanding: map[string]float64{}, WithoutValue: map[string]int{}}
		for _, col := range costColumns {
			p.Outstanding[col], p.WithoutValue[col] = 0, 0
		}
		if !worth.before(next) {
			for i := range valued {
				if st := worth.at(i, end); slices.Contains(costColumns, st) {
					p.Outstanding[st] += waits[i].value
				}
			}
			for col, v := range p.Outstanding {
				p.Outstanding[col] = round2(v)
			}
			p.Incurred = incurred(waits, d, next, now)
		}
		if !without.before(next) {
			for i := range unvalued {
				if st := without.at(i, end); slices.Contains(costColumns, st) {
					p.WithoutValue[st]++
				}
			}
		}
		out.Days = append(out.Days, p)
	}
	for m := monday(start); !m.After(now); m = m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		out.Weeks = append(out.Weeks, CostWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02"),
			Incurred: incurred(waits, m, m.AddDate(0, 0, 7), now)})
	}
	return out
}

// waitedItem is an item's cost of delay per week and its spans in backlog or
// ready.
type waitedItem struct {
	value float64
	spans []span
}

// incurred is what waiting for the items cost from a up to b, not after now.
func incurred(items []waitedItem, a, b, now time.Time) float64 {
	var total float64
	for _, it := range items {
		total += it.value * waited(it.spans, a, b, now)
	}
	return round2(total / secondsPerWeek)
}

// waitSpans is the item's spans in backlog or ready, in the order of its
// transitions: backlog from created, and the last state up to now while it
// is open. A span may run backwards where the transitions are out of order;
// waited counts nothing of it.
func waitSpans(it *workitem.Item, tl *timeline, now time.Time) []span {
	var out []span
	add := func(state string, from, to time.Time) {
		if state == workitem.Backlog || state == workitem.Ready {
			out = append(out, span{from, to})
		}
	}
	state, at := workitem.Backlog, tl.created
	for i, next := range tl.at {
		add(state, at, next)
		state, at = tl.to[i], next
	}
	if !it.Closed() {
		add(state, at, now)
	}
	return out
}

// waited is the seconds the spans cover from a up to b, not after now, each
// span's added in turn.
func waited(spans []span, a, b, now time.Time) float64 {
	if b.After(now) {
		b = now
	}
	var total float64
	for _, s := range spans {
		from, to := s.from, s.to
		if from.Before(a) {
			from = a
		}
		if to.After(b) {
			to = b
		}
		if to.After(from) {
			total += to.Sub(from).Seconds()
		}
	}
	return total
}

// costValue is the item's cost of delay per week, nil without one.
func costValue(it *workitem.Item) *float64 {
	if it.CostOfDelay == nil {
		return nil
	}
	return it.CostOfDelay.Value
}

// round2 rounds an amount to two decimals.
func round2(x float64) float64 {
	return math.Round(x*100) / 100
}

// round4 rounds an amount in US dollars to four decimals, as activity costs
// are written.
func round4(x float64) float64 {
	return math.Round(x*1e4) / 1e4
}
