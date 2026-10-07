package metrics

import (
	"fmt"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/planning"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orders the ready column's cost of delay is projected under: the
// board's pull order, and S-0217's cod and wsjf policies.
const (
	OrderCurrent = "current"
	OrderCOD     = workitem.PolicyCOD
	OrderWSJF    = workitem.PolicyWSJF
)

// CostOrder is the ready column's cost of delay projected from At until each
// story is pulled, under the pull order, by cost of delay, and by WSJF
// (S-0213).
type CostOrder struct {
	At      string        `json:"at"`
	Horizon string        `json:"horizon"`
	Series  []OrderSeries `json:"series"`
	// Saving is current's total less the lower of cod's and wsjf's, and
	// Cheaper the order with that lower total.
	Saving  float64 `json:"saving"`
	Cheaper string  `json:"cheaper"`
	// LeftOut are the ready stories without a value or a forecast duration,
	// in pull order.
	LeftOut []string `json:"left_out"`
}

// OrderSeries is one order's cumulative projected cost against each placed
// story's pull.
type OrderSeries struct {
	By     string       `json:"by"`
	Total  float64      `json:"total"`
	Points []OrderPoint `json:"points"`
}

// OrderPoint is a story's projected pull and the cost of the stories pulled
// up to and including it; the first point of a series has no story.
type OrderPoint struct {
	At       string  `json:"at"`
	ID       string  `json:"id,omitempty"`
	Incurred float64 `json:"incurred"`
}

// costOrder projects the cost of delay of the ready stories in all that have
// a value and a forecast duration, played out on the board's lanes from the
// report's now under each order.
func costOrder(all []*workitem.Item, opt Options) CostOrder {
	at := opt.Now.Truncate(time.Second)
	byID := map[string]*workitem.Item{}
	for _, it := range all {
		if !it.Archived {
			byID[it.ID] = it
		}
	}
	var placed []*workitem.Item
	out := CostOrder{At: at.Format(workitem.TimeFormat), LeftOut: []string{}}
	for _, id := range workitem.PullSequence(opt.Order, all, workitem.Ready) {
		if it := byID[id]; costValue(it) != nil && hasForecastDuration(it) {
			placed = append(placed, it)
		} else {
			out.LeftOut = append(out.LeftOut, id)
		}
	}
	fallback := opt.Fallback
	if fallback <= 0 {
		fallback = manifest.DefaultDuration
	}
	horizon := at
	for _, by := range []string{OrderCurrent, OrderCOD, OrderWSJF} {
		seq := orderedBy(placed, by)
		s, last := projectSeries(by, seq, planning.PullTimes(all, opt.Order, opt.WIPLimit, fallback, seq, at), at)
		if last.After(horizon) {
			horizon = last
		}
		out.Series = append(out.Series, s)
	}
	out.Horizon = horizon.Format(workitem.TimeFormat)
	current, cod, wsjf := out.Series[0].Total, out.Series[1].Total, out.Series[2].Total
	out.Cheaper, out.Saving = OrderCOD, round2(current-cod)
	if wsjf < cod {
		out.Cheaper, out.Saving = OrderWSJF, round2(current-wsjf)
	}
	return out
}

// orderedBy is the placed stories, given in pull order, in the order by
// names: as given for current, else as its policy orders them.
func orderedBy(placed []*workitem.Item, by string) []*workitem.Item {
	if by == OrderCurrent {
		return placed
	}
	ranked, err := workitem.OrderByPolicy(placed, by)
	if err != nil {
		// by is one of flai's own policy names, so this is a bug in flai
		panic(fmt.Sprintf("cannot order the ready stories by %s: %v", by, err))
	}
	byID := map[string]*workitem.Item{}
	for _, it := range placed {
		byID[it.ID] = it
	}
	out := make([]*workitem.Item, len(ranked))
	for i, r := range ranked {
		out[i] = byID[r.ID]
	}
	return out
}

// projectSeries lays out the stories, pulled at pulls, in order of pull, ties
// in their sequence, with the cumulative cost of their wait from at to each
// pull rounded up to the second. It gives the series and its last pull.
func projectSeries(by string, seq []*workitem.Item, pulls []time.Time, at time.Time) (OrderSeries, time.Time) {
	type pulled struct {
		it *workitem.Item
		at time.Time
	}
	ps := make([]pulled, len(seq))
	for i, it := range seq {
		ps[i] = pulled{it, ceilSecond(pulls[i])}
	}
	sort.SliceStable(ps, func(i, j int) bool { return ps[i].at.Before(ps[j].at) })
	s := OrderSeries{By: by, Points: []OrderPoint{{At: at.Format(workitem.TimeFormat)}}}
	last := at
	var sum float64
	for _, p := range ps {
		sum += *costValue(p.it) * p.at.Sub(at).Seconds() / secondsPerWeek
		s.Points = append(s.Points, OrderPoint{At: p.at.Format(workitem.TimeFormat), ID: p.it.ID, Incurred: round2(sum)})
		last = p.at
	}
	s.Total = s.Points[len(s.Points)-1].Incurred
	return s, last
}

// hasForecastDuration says the item's forecast.duration is a positive Go
// duration.
func hasForecastDuration(it *workitem.Item) bool {
	if it.Forecast == nil {
		return false
	}
	d, err := time.ParseDuration(it.Forecast.Duration)
	return err == nil && d > 0
}

// ceilSecond rounds t up to the second.
func ceilSecond(t time.Time) time.Time {
	s := t.Truncate(time.Second)
	if s.Before(t) {
		s = s.Add(time.Second)
	}
	return s
}
