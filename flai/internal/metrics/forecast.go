package metrics

import (
	"math"
	"sort"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Forecasts is how far the planner's forecasts and the human's estimates were
// from what happened, over the items completed in the window (S-0205).
type Forecasts struct {
	Forecast ErrorStats `json:"forecast"`
	Delivery ErrorStats `json:"delivery"`
	Estimate ErrorStats `json:"estimate"`
}

// ErrorStats spreads one kind of error over the items that have it, in all,
// per nature, and per story agent model.
type ErrorStats struct {
	ErrorSpread
	ByNature map[string]ErrorSpread `json:"by_nature"`
	ByModel  map[string]ErrorSpread `json:"by_model"`
}

// ErrorSpread is the count of a set of errors and the percentiles of their
// absolute values in seconds, absent for an empty set.
type ErrorSpread struct {
	Count int      `json:"count"`
	P50   *float64 `json:"p50_seconds,omitempty"`
	P85   *float64 `json:"p85_seconds,omitempty"`
}

// noModel is the model of an item whose agent names none.
const noModel = "(none)"

// deriveForecast sets the item's forecast and its errors, and the error of its
// estimate in seconds, from the cycle time and estimate already derived.
func deriveForecast(m *ItemMetrics, it *workitem.Item, completed time.Time) {
	m.Model = noModel
	if it.Agent != nil && it.Agent.Model != "" {
		m.Model = it.Agent.Model
	}
	if m.CycleTime != nil && m.Estimate != nil {
		e := *m.CycleTime - *m.Estimate
		m.EstErrorSeconds = &e
	}
	f := it.Forecast
	if f == nil {
		return
	}
	if d, err := time.ParseDuration(f.Duration); err == nil && d > 0 {
		m.Forecast = secs(d)
		if m.CycleTime != nil {
			e := *m.CycleTime - d.Seconds()
			m.ForecastError = &e
		}
	}
	if delivery, err := time.Parse(workitem.TimeFormat, f.Delivery); err == nil && !completed.IsZero() {
		m.DeliveryError = secs(completed.Sub(delivery))
	}
}

// forecasts spreads each error over the items done in the window that have it.
func forecasts(items []*workitem.Item, per map[string]ItemMetrics, inWindow func(*workitem.Item) bool) Forecasts {
	var forecast, delivery, estimate errorSet
	for _, it := range items {
		if it.Status != workitem.Done || !inWindow(it) {
			continue
		}
		m := per[it.ID]
		forecast.add(m.ForecastError, it.Nature, m.Model)
		delivery.add(m.DeliveryError, it.Nature, m.Model)
		estimate.add(m.EstErrorSeconds, it.Nature, m.Model)
	}
	return Forecasts{Forecast: forecast.stats(), Delivery: delivery.stats(), Estimate: estimate.stats()}
}

// errorSet gathers absolute errors, in all and by nature and model.
type errorSet struct {
	all           []float64
	nature, model map[string][]float64
}

func (s *errorSet) add(e *float64, nature, model string) {
	if e == nil {
		return
	}
	if s.nature == nil {
		s.nature, s.model = map[string][]float64{}, map[string][]float64{}
	}
	v := math.Abs(*e)
	s.all = append(s.all, v)
	s.nature[nature] = append(s.nature[nature], v)
	s.model[model] = append(s.model[model], v)
}

func (s errorSet) stats() ErrorStats {
	out := ErrorStats{ErrorSpread: spread(s.all), ByNature: map[string]ErrorSpread{}, ByModel: map[string]ErrorSpread{}}
	for k, v := range s.nature {
		out.ByNature[k] = spread(v)
	}
	for k, v := range s.model {
		out.ByModel[k] = spread(v)
	}
	return out
}

func spread(v []float64) ErrorSpread {
	if len(v) == 0 {
		return ErrorSpread{}
	}
	sorted := append([]float64{}, v...)
	sort.Float64s(sorted)
	p50, p85 := percentile(sorted, 50), percentile(sorted, 85)
	return ErrorSpread{Count: len(v), P50: &p50, P85: &p85}
}
