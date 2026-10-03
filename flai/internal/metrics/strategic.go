package metrics

import (
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// StrategicAgent is a strategic agent's activity document as flai stats
// reports it: its totals as written, all time, and its log entries in the
// window (ADR-0079).
type StrategicAgent struct {
	Kind       string           `json:"kind"`
	Cost       float64          `json:"cost"`
	Seconds    int64            `json:"seconds"`
	Activities int              `json:"activities"`
	LastRun    string           `json:"last_run"`
	Log        []StrategicEntry `json:"log"`
}

// StrategicEntry is one activity logged in the window.
type StrategicEntry struct {
	At        string   `json:"at"`
	Seconds   int64    `json:"seconds"`
	Cost      float64  `json:"cost"`
	Estimated bool     `json:"estimated,omitempty"`
	Items     []string `json:"items"`
}

// StrategicDay is what the strategic agents spent on one day, beside what was
// delivered that day (S-0205).
type StrategicDay struct {
	Date string `json:"date"`
	// Agents is the use of each kind with an entry that ended that day.
	Agents  map[string]StrategicUse `json:"agents"`
	Cost    float64                 `json:"cost"`
	Seconds int64                   `json:"seconds"`
	// Completed counts the items of the report's type done that day;
	// CostPerItem is the mean usage cost of those carrying usage, and
	// CycleTime the mean cycle time of those with one.
	Completed   int      `json:"completed"`
	CostPerItem *float64 `json:"cost_per_item,omitempty"`
	CycleTime   *float64 `json:"cycle_time_seconds,omitempty"`
}

// StrategicUse is what one strategic agent spent on a day.
type StrategicUse struct {
	Cost      float64 `json:"cost"`
	Seconds   int64   `json:"seconds"`
	Estimated bool    `json:"estimated,omitempty"`
}

// strategicDays lays out, by the day from the one that holds the window's
// start to today, the activities that ended each day, not after now, and the
// items done each day with their usage cost and cycle time.
func strategicDays(docs []*workitem.Activity, items []*workitem.Item, per map[string]ItemMetrics, start, now time.Time) []StrategicDay {
	out := []StrategicDay{}
	at := map[string]int{}
	for d := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC); !d.After(now); d = d.AddDate(0, 0, 1) {
		at[d.Format("2006-01-02")] = len(out)
		out = append(out, StrategicDay{Date: d.Format("2006-01-02"), Agents: map[string]StrategicUse{}})
	}
	for _, doc := range docs {
		for _, e := range doc.Entries {
			i, ok := at[e.At.UTC().Format("2006-01-02")]
			if !ok || e.At.After(now) {
				continue
			}
			p := &out[i]
			u := p.Agents[doc.Kind]
			u.Cost += e.Cost
			u.Seconds += e.Seconds
			u.Estimated = u.Estimated || e.Estimated
			p.Agents[doc.Kind] = u
			p.Cost += e.Cost
			p.Seconds += e.Seconds
		}
	}
	type delivery struct {
		cost, cycle float64
		used, timed int
	}
	done := make([]delivery, len(out))
	for _, it := range items {
		c := it.CompletedAt()
		i, ok := at[c.Format("2006-01-02")]
		if it.Status != workitem.Done || !ok || c.After(now) {
			continue
		}
		out[i].Completed++
		m := per[it.ID]
		if m.Usage != nil {
			done[i].cost += m.Usage.Cost
			done[i].used++
		}
		if m.CycleTime != nil {
			done[i].cycle += *m.CycleTime
			done[i].timed++
		}
	}
	for i := range out {
		p := &out[i]
		for kind, u := range p.Agents {
			u.Cost = round4(u.Cost)
			p.Agents[kind] = u
		}
		p.Cost = round4(p.Cost)
		if c := over(done[i].cost, float64(done[i].used)); c != nil {
			*c = round4(*c)
			p.CostPerItem = c
		}
		p.CycleTime = over(done[i].cycle, float64(done[i].timed))
	}
	return out
}

// strategic reports each activity document, in the order given, with the
// entries that ended from start to now.
func strategic(docs []*workitem.Activity, start, now time.Time) []StrategicAgent {
	out := make([]StrategicAgent, 0, len(docs))
	for _, d := range docs {
		s := StrategicAgent{Kind: d.Kind, Cost: d.AccruedCost, Seconds: d.AccruedSeconds, Activities: d.TasksCompleted, LastRun: d.LastRun, Log: []StrategicEntry{}}
		for _, e := range d.Entries {
			if e.At.Before(start) || e.At.After(now) {
				continue
			}
			items := e.Items
			if items == nil {
				items = []string{}
			}
			s.Log = append(s.Log, StrategicEntry{At: e.At.UTC().Format(workitem.TimeFormat), Seconds: e.Seconds, Cost: e.Cost, Estimated: e.Estimated, Items: items})
		}
		out = append(out, s)
	}
	return out
}
