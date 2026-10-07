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
// the items without a value, of a type that carries one, in each column.
func costOfDelay(items []*workitem.Item, start, now time.Time) CostOfDelay {
	var valued, unvalued []*workitem.Item
	for _, it := range items {
		switch {
		case costValue(it) != nil:
			valued = append(valued, it)
		case workitem.Carries(it.Type, "cost_of_delay"):
			unvalued = append(unvalued, it)
		}
	}
	out := CostOfDelay{Days: []CostDay{}, Weeks: []CostWeek{}}
	for d := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC); !d.After(now); d = d.AddDate(0, 0, 1) {
		end := d.Add(24*time.Hour - time.Second)
		p := CostDay{Date: d.Format("2006-01-02"), Outstanding: map[string]float64{}, WithoutValue: map[string]int{}}
		for _, col := range costColumns {
			p.Outstanding[col], p.WithoutValue[col] = 0, 0
		}
		for _, it := range valued {
			if st := stateAt(it, end); slices.Contains(costColumns, st) {
				p.Outstanding[st] += *costValue(it)
			}
		}
		for _, it := range unvalued {
			if st := stateAt(it, end); slices.Contains(costColumns, st) {
				p.WithoutValue[st]++
			}
		}
		for col, v := range p.Outstanding {
			p.Outstanding[col] = round2(v)
		}
		p.Incurred = incurred(valued, d, d.AddDate(0, 0, 1), now)
		out.Days = append(out.Days, p)
	}
	for m := monday(start); !m.After(now); m = m.AddDate(0, 0, 7) {
		y, w := m.ISOWeek()
		out.Weeks = append(out.Weeks, CostWeek{Week: weekKey(y, w), Start: m.Format("2006-01-02"),
			Incurred: incurred(valued, m, m.AddDate(0, 0, 7), now)})
	}
	return out
}

// incurred is what waiting for the items cost from a up to b, not after now.
func incurred(items []*workitem.Item, a, b, now time.Time) float64 {
	var total float64
	for _, it := range items {
		total += *costValue(it) * waited(it, a, b, now)
	}
	return round2(total / secondsPerWeek)
}

// waited is the seconds the item spent in backlog or ready from a up to b,
// not after now: backlog from created, and the last state up to now while it
// is open.
func waited(it *workitem.Item, a, b, now time.Time) float64 {
	if b.After(now) {
		b = now
	}
	var total float64
	add := func(state string, from, to time.Time) {
		if state != workitem.Backlog && state != workitem.Ready {
			return
		}
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
	created, _ := time.Parse(workitem.TimeFormat, it.Created)
	state, at := workitem.Backlog, created
	for _, tr := range it.Transitions {
		next, _ := time.Parse(workitem.TimeFormat, tr.At)
		add(state, at, next)
		state, at = tr.To, next
	}
	if !it.Closed() {
		add(state, at, now)
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
