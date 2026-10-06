package metrics

import (
	"time"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// StrategicAgent is a strategic agent's activity document as flai stats
// reports it: its totals as written, all time, what of them the items carry,
// what the issues carry that no story does (S-0227), and the project
// strategic total (ADR-0095), and its log entries in the window (ADR-0079).
type StrategicAgent struct {
	Kind       string           `json:"kind"`
	Cost       float64          `json:"cost"`
	Seconds    int64            `json:"seconds"`
	Activities int              `json:"activities"`
	LastRun    string           `json:"last_run"`
	Items      StrategicAmount  `json:"items"`
	Issues     StrategicAmount  `json:"issues"`
	Project    StrategicAmount  `json:"project"`
	Log        []StrategicEntry `json:"log"`
}

// StrategicIssue is what strategic agents spent on one issue (S-0227): its
// entries per kind, the story made from it, "" for none, and whether its
// usage is counted in the totals, which it is until such a story carries it.
type StrategicIssue struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Status    string          `json:"status"`
	Story     string          `json:"story"`
	Counted   bool            `json:"counted"`
	Strategic []ItemStrategic `json:"strategic"`
}

// StrategicAmount is a cost and a time a strategic agent spent.
type StrategicAmount struct {
	Cost    float64 `json:"cost"`
	Seconds int64   `json:"seconds"`
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
	// CostPerItem is the mean agents' cost of those on which agents spent,
	// what strategic agents spent left out, and
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
		if !it.Usage.Empty() {
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

// strategic reports each activity document, in the order given, with what
// of its totals the items carry, what the issues carry that no story does,
// onIssues, the rest as the project strategic total, and the entries that
// ended from start to now.
func strategic(docs []*workitem.Activity, items []*workitem.Item, onIssues map[string]StrategicAmount, start, now time.Time) []StrategicAgent {
	carried := onItems(items)
	out := make([]StrategicAgent, 0, len(docs))
	for _, d := range docs {
		s := StrategicAgent{Kind: d.Kind, Cost: d.AccruedCost, Seconds: d.AccruedSeconds, Activities: d.TasksCompleted, LastRun: d.LastRun, Log: []StrategicEntry{}}
		on, is := carried[d.Kind], onIssues[d.Kind]
		s.Items = StrategicAmount{Cost: round4(on.Cost), Seconds: on.Seconds}
		s.Issues = StrategicAmount{Cost: round4(is.Cost), Seconds: is.Seconds}
		s.Project = StrategicAmount{Cost: max(0, round4(d.AccruedCost-s.Items.Cost-s.Issues.Cost)), Seconds: max(0, d.AccruedSeconds-on.Seconds-is.Seconds)}
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

// onItems is what the items carry of each strategic agent kind's spending:
// the sum over the items at the top of the hierarchy, those with no parent
// among them, which carry every charge made below them (ADR-0095).
func onItems(items []*workitem.Item) map[string]StrategicAmount {
	ids := map[string]bool{}
	for _, it := range items {
		ids[workitem.CanonicalID(it.ID)] = true
	}
	out := map[string]StrategicAmount{}
	for _, it := range items {
		if it.Usage == nil || (it.Parent != "" && ids[workitem.CanonicalID(it.Parent)]) {
			continue
		}
		for _, s := range it.Usage.Strategic {
			a := out[s.Kind]
			a.Cost += s.Cost()
			a.Seconds += s.Seconds
			out[s.Kind] = a
		}
	}
	return out
}

// strategicIssues lists, in the order given, the issues that carry strategic
// usage, and sums per kind what those carry from which no story among the
// items was made (S-0227). A story made from an issue carries the issue's
// entries, so it counts them among the items instead; one the items do not
// hold is no such story.
func strategicIssues(list []*issues.Issue, items []*workitem.Item) ([]StrategicIssue, map[string]StrategicAmount) {
	stories := map[string]bool{}
	for _, it := range items {
		if it.Type == workitem.Story {
			stories[workitem.CanonicalID(it.ID)] = true
		}
	}
	out, counted := []StrategicIssue{}, map[string]StrategicAmount{}
	for _, is := range list {
		if is.Usage == nil || len(is.Usage.Strategic) == 0 {
			continue
		}
		s := StrategicIssue{ID: is.ID, Title: is.Title, Status: is.Status, Strategic: []ItemStrategic{}}
		for _, id := range issues.StoriesMade(is) {
			if stories[id] {
				s.Story = id
				break
			}
		}
		s.Counted = s.Story == ""
		for _, e := range is.Usage.Strategic {
			s.Strategic = append(s.Strategic, ItemStrategic{Kind: e.Kind, Tokens: e.Tokens(), Cost: round4(e.Cost()), Seconds: e.Seconds, Estimated: e.Estimated})
			if s.Counted {
				a := counted[e.Kind]
				a.Cost += e.Cost()
				a.Seconds += e.Seconds
				counted[e.Kind] = a
			}
		}
		out = append(out, s)
	}
	return out, counted
}
