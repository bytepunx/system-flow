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
