package metrics

import (
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// TimerDefinition is one of the workflow's own timers as the report lists it
// (ADR-0317 §2).
type TimerDefinition struct {
	Name   string   `json:"name"`
	Title  string   `json:"title"`
	Types  []string `json:"types"`
	Counts string   `json:"counts"`
}

// reportedApart are the timers the report gives under their own field
// names, lead_time_seconds and the others, and so not under timers
// (ADR-0305 §5 and §6, ADR-0317 §1).
var reportedApart = []string{"lead", "cycle", "queue", "review"}

// ownTimers are flow's timers in file order, but for those reported apart.
func ownTimers(flow workflow.Workflow) []workflow.Timer {
	var out []workflow.Timer
	for _, t := range flow.Timers {
		if !slices.Contains(reportedApart, t.Name) {
			out = append(out, t)
		}
	}
	return out
}

// timerDefinitions are flow's own timers as the report lists them, nil with
// none.
func timerDefinitions(flow workflow.Workflow) []TimerDefinition {
	var out []TimerDefinition
	for _, t := range ownTimers(flow) {
		out = append(out, TimerDefinition{Name: t.Name, Title: t.Title, Types: slices.Clone(t.Types), Counts: t.Counts})
	}
	return out
}

// itemTimers is the value in seconds of each of flow's own timers that takes
// the history's type and has a value for it, nil with none (ADR-0317 §3).
func itemTimers(flow workflow.Workflow, h workflow.History, now time.Time) map[string]float64 {
	var out map[string]float64
	for _, t := range ownTimers(flow) {
		v := timed(flow, t.Name, h, now)
		if v == nil {
			continue
		}
		if out == nil {
			out = map[string]float64{}
		}
		out[t.Name] = *v
	}
	return out
}

// timerDistributions is the distribution of each of flow's own timers that
// takes typ over the items completed in the window its counts takes, those
// accepted or, for closed, cancelled too, and that have a value for it; nil
// when no own timer takes typ (ADR-0317 §4).
func timerDistributions(flow workflow.Workflow, typ string, items []*workitem.Item, per map[string]ItemMetrics, inWindow func(*workitem.Item) bool) map[string]Distribution {
	var out map[string]Distribution
	for _, t := range ownTimers(flow) {
		if !slices.Contains(t.Types, typ) {
			continue
		}
		var values []float64
		for _, it := range items {
			if !inWindow(it) || !counted(flow, t, it.Status) {
				continue
			}
			if v, ok := per[it.ID].Timers[t.Name]; ok {
				values = append(values, v)
			}
		}
		if out == nil {
			out = map[string]Distribution{}
		}
		out[t.Name] = distribution(values)
	}
	return out
}

// counted reports whether the timer's aggregates take an item closed in
// state: one of the accepted kind, or with counts closed of the cancelled
// kind too.
func counted(flow workflow.Workflow, t workflow.Timer, state string) bool {
	switch flow.Kind(state) {
	case workflow.KindAccepted:
		return true
	case workflow.KindCancelled:
		return t.Counts == workflow.CountsClosed
	}
	return false
}
