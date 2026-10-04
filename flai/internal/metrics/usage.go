package metrics

import (
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Usage metrics (S-0143, ADR-0051): what agents spent on items, from each
// item's usage, grouped by the model that spent it. design/system/metrics.md
// defines them.

// ModelSpend is what one model spent on an item, or on a set of items.
type ModelSpend struct {
	Model  string  `json:"model"`
	Tokens int64   `json:"tokens"`
	Cost   float64 `json:"cost"`
	// TokensPerMinute is tokens over the minutes agents worked on the item,
	// or on the items the model worked on; absent when they worked none.
	// TokensPerHour is sixty times that, for what was written to flai 1.25
	// (ADR-0053).
	TokensPerMinute *float64 `json:"tokens_per_minute,omitempty"`
	TokensPerHour   *float64 `json:"tokens_per_hour,omitempty"`
	// Items counts the items the model worked on, in a total.
	Items int `json:"items,omitempty"`
}

// ItemUsage is what agents spent on one item, and apart from that what
// strategic agents spent on it.
type ItemUsage struct {
	Source          string       `json:"source"`
	Tokens          int64        `json:"tokens"`
	Cost            float64      `json:"cost"`
	Seconds         int64        `json:"seconds"`
	TokensPerMinute *float64     `json:"tokens_per_minute,omitempty"`
	TokensPerHour   *float64     `json:"tokens_per_hour,omitempty"`
	Estimated       bool         `json:"estimated,omitempty"`
	Models          []ModelSpend `json:"models"`
	// Strategic is what each kind of strategic agent spent on it, never
	// added to the agents' figures above (ADR-0083).
	Strategic []ItemStrategic `json:"strategic"`
}

// ItemStrategic is what one kind of strategic agent spent on an item.
type ItemStrategic struct {
	Kind      string  `json:"kind"`
	Tokens    int64   `json:"tokens"`
	Cost      float64 `json:"cost"`
	Seconds   int64   `json:"seconds"`
	Estimated bool    `json:"estimated"`
}

// StrategicSpend is what strategic agents spent on a set of items: those
// that carry any of it, and their tokens, cost, and seconds (ADR-0083).
type StrategicSpend struct {
	Items   int     `json:"items"`
	Tokens  int64   `json:"tokens"`
	Cost    float64 `json:"cost"`
	Seconds int64   `json:"seconds"`
}

// add adds what strategic agents spent on one item, if anything.
func (s *StrategicSpend) add(u *usage.Usage) {
	if u == nil || len(u.Strategic) == 0 {
		return
	}
	s.Items++
	s.Tokens += u.StrategicTokens()
	s.Cost += u.StrategicCost()
	s.Seconds += u.StrategicSeconds()
}

// rounded is s with its cost rounded to four decimals.
func (s StrategicSpend) rounded() StrategicSpend {
	s.Cost = round4(s.Cost)
	return s
}

// StrategicKindSpend is what one kind of strategic agent spent on the items
// it worked on.
type StrategicKindSpend struct {
	Kind string `json:"kind"`
	StrategicSpend
}

// StrategicTotals is what strategic agents spent on the items done in the
// window, in all and per kind.
type StrategicTotals struct {
	StrategicSpend
	Estimated bool `json:"estimated"`
	// Kinds are the same per kind, in the order planner, orchestrator,
	// analyzer, only those that spent.
	Kinds []StrategicKindSpend `json:"kinds"`
}

// SpendPoint is one item done, in order of completion, with what had been
// done and spent by then.
type SpendPoint struct {
	At string `json:"at"`
	ID string `json:"id"`
	// Done counts the items done by then, this one included; Tokens and Cost
	// are what they spent between them.
	Done   int     `json:"done"`
	Tokens int64   `json:"tokens"`
	Cost   float64 `json:"cost"`
}

// UsageReport is what agents spent on the items of the report's type.
type UsageReport struct {
	// Items counts the items completed in the window that carry usage.
	Items     int     `json:"items"`
	Tokens    int64   `json:"tokens"`
	Cost      float64 `json:"cost"`
	Seconds   int64   `json:"seconds"`
	Estimated bool    `json:"estimated,omitempty"`
	// Models are the same per model, in order of name.
	Models []ModelSpend `json:"models"`
	// Done is completion against time and cost: every item done in the
	// window that carries usage, oldest first, cumulatively; ByModel is the
	// same for the items each model worked on, counting that model's spend.
	Done    []SpendPoint            `json:"done"`
	ByModel map[string][]SpendPoint `json:"by_model"`
	// Bucket is what Spend's series are laid out in: hour, day, or week.
	// Spend is what was spent on epics, on stories, and on tasks, in all
	// and over time, whatever type the report is about (S-0163).
	Bucket string                `json:"bucket"`
	Spend  map[string]*TypeSpend `json:"spend"`
	// Strategic is what strategic agents spent on the items of the report's
	// type done in the window, apart from every figure above (ADR-0083).
	Strategic StrategicTotals `json:"strategic"`
	// CostPerAgentHour is the project's mean cost of an hour of agent work,
	// over every story measured from its logs; absent before the first.
	CostPerAgentHour *float64 `json:"cost_per_agent_hour,omitempty"`
}

func perHour(tokens, seconds int64) *float64 {
	return over(float64(tokens), float64(seconds)/3600)
}

func perMinute(tokens, seconds int64) *float64 {
	return over(float64(tokens), float64(seconds)/60)
}

// itemUsage is what an item's usage comes to; nil when neither agents nor
// strategic agents spent anything on it. An item that carries only what
// strategic agents spent has the agents' figures at zero.
func itemUsage(u *usage.Usage) *ItemUsage {
	if u.Nothing() {
		return nil
	}
	out := &ItemUsage{Source: u.Source, Tokens: u.Tokens(), Cost: u.Cost(), Seconds: u.Seconds, Estimated: u.Estimated, Models: []ModelSpend{}, Strategic: []ItemStrategic{}}
	out.TokensPerMinute, out.TokensPerHour = perMinute(out.Tokens, u.Seconds), perHour(out.Tokens, u.Seconds)
	for _, m := range u.Models {
		out.Models = append(out.Models, ModelSpend{Model: m.Model, Tokens: m.Tokens(), Cost: m.Cost,
			TokensPerMinute: perMinute(m.Tokens(), u.Seconds), TokensPerHour: perHour(m.Tokens(), u.Seconds)})
	}
	for _, s := range u.Strategic {
		out.Strategic = append(out.Strategic, ItemStrategic{Kind: s.Kind, Tokens: s.Tokens(), Cost: round4(s.Cost()), Seconds: s.Seconds, Estimated: s.Estimated})
	}
	return out
}

// agentsSpent says agents spent something on the item; spentAny says agents
// or strategic agents did; strategicSpent says strategic agents did.
func agentsSpent(u *usage.Usage) bool    { return !u.Empty() }
func spentAny(u *usage.Usage) bool       { return !u.Nothing() }
func strategicSpent(u *usage.Usage) bool { return u != nil && len(u.Strategic) > 0 }

// strategicTotals sums what strategic agents spent on the items done in the
// window, in all and per kind.
func strategicTotals(items []*workitem.Item, start, now time.Time) StrategicTotals {
	var all StrategicSpend
	kinds := map[string]*StrategicSpend{}
	var names []string
	estimated := false
	for _, c := range spent(items, start, now, strategicSpent) {
		u := c.it.Usage
		all.add(u)
		for _, s := range u.Strategic {
			k := kinds[s.Kind]
			if k == nil {
				k = &StrategicSpend{}
				kinds[s.Kind] = k
				names = append(names, s.Kind)
			}
			k.Items++
			k.Tokens += s.Tokens()
			k.Cost += s.Cost()
			k.Seconds += s.Seconds
			estimated = estimated || s.Estimated
		}
	}
	out := StrategicTotals{StrategicSpend: all.rounded(), Estimated: estimated, Kinds: []StrategicKindSpend{}}
	rank := func(kind string) int {
		if i := slices.Index(usage.StrategicKinds, kind); i >= 0 {
			return i
		}
		return len(usage.StrategicKinds)
	}
	slices.SortFunc(names, func(a, b string) int {
		if d := rank(a) - rank(b); d != 0 {
			return d
		}
		return strings.Compare(a, b)
	})
	for _, n := range names {
		out.Kinds = append(out.Kinds, StrategicKindSpend{Kind: n, StrategicSpend: kinds[n].rounded()})
	}
	return out
}

// CostPerAgentHour is the project's mean cost of an hour of agent work: the
// agents' cost over their hours, summed over every story, whatever its status
// and archived ones included, whose usage was measured from its logs and
// took agent time; what strategic agents spent is left out. It is rounded to
// four decimals, and nil when no story has been measured (ADR-0083).
func CostPerAgentHour(items []*workitem.Item) *float64 {
	var cost float64
	var seconds int64
	for _, it := range items {
		u := it.Usage
		if it.Type != workitem.Story || u == nil || u.Source != usage.SourceLog || u.Seconds <= 0 {
			continue
		}
		cost += u.Cost()
		seconds += u.Seconds
	}
	r := over(cost, float64(seconds)/3600)
	if r != nil {
		*r = round4(*r)
	}
	return r
}

// ExpectedCost is what an item is expected to cost before agents work it:
// the duration of its forecast, or else of its estimate, priced at the
// project's cost per agent hour. It is always an estimate.
type ExpectedCost struct {
	Cost float64 `json:"cost"`
	// From is what the duration came from: "forecast" or "estimate".
	From      string `json:"from"`
	Estimated bool   `json:"estimated"`
}

// ExpectedCostOf is the item's expected cost at rate, CostPerAgentHour's:
// its forecast duration, or else its estimate, in hours, times rate, rounded
// to four decimals. It is nil when rate is, or the item has neither duration.
func ExpectedCostOf(it *workitem.Item, rate *float64) *ExpectedCost {
	if rate == nil {
		return nil
	}
	hours := func(v string) float64 {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			return 0
		}
		return d.Hours()
	}
	from, h := "", 0.0
	if it.Forecast != nil {
		from, h = "forecast", hours(it.Forecast.Duration)
	}
	if h == 0 && it.Estimate != "" {
		from, h = "estimate", hours(it.Estimate)
	}
	if h == 0 {
		return nil
	}
	return &ExpectedCost{Cost: round4(h * *rate), From: from, Estimated: true}
}

// spendReport totals the usage of the items done in the window, and lays
// them out against time and cost. An item on which only strategic agents
// spent is left out.
func spendReport(items []*workitem.Item, start, now time.Time) UsageReport {
	rep := UsageReport{Models: []ModelSpend{}, Done: []SpendPoint{}, ByModel: map[string][]SpendPoint{}}
	list := spent(items, start, now, agentsSpent)
	models := map[string]*ModelSpend{}
	seconds := map[string]int64{}
	for _, d := range list {
		u := d.it.Usage
		rep.Items++
		rep.Tokens += u.Tokens()
		rep.Cost += u.Cost()
		rep.Seconds += u.Seconds
		rep.Estimated = rep.Estimated || u.Estimated
		at := d.at.Format(workitem.TimeFormat)
		rep.Done = append(rep.Done, SpendPoint{At: at, ID: d.it.ID, Done: rep.Items, Tokens: rep.Tokens, Cost: rep.Cost})
		for _, m := range u.Models {
			ms := models[m.Model]
			if ms == nil {
				ms = &ModelSpend{Model: m.Model}
				models[m.Model] = ms
			}
			ms.Items++
			ms.Tokens += m.Tokens()
			ms.Cost += m.Cost
			seconds[m.Model] += u.Seconds
			rep.ByModel[m.Model] = append(rep.ByModel[m.Model], SpendPoint{At: at, ID: d.it.ID, Done: ms.Items, Tokens: ms.Tokens, Cost: ms.Cost})
		}
	}
	for name, ms := range models {
		ms.TokensPerMinute, ms.TokensPerHour = perMinute(ms.Tokens, seconds[name]), perHour(ms.Tokens, seconds[name])
		rep.Models = append(rep.Models, *ms)
	}
	slices.SortFunc(rep.Models, func(a, b ModelSpend) int { return strings.Compare(a.Model, b.Model) })
	return rep
}
