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

// ItemUsage is what agents spent on one item.
type ItemUsage struct {
	Source          string       `json:"source"`
	Tokens          int64        `json:"tokens"`
	Cost            float64      `json:"cost"`
	Seconds         int64        `json:"seconds"`
	TokensPerMinute *float64     `json:"tokens_per_minute,omitempty"`
	TokensPerHour   *float64     `json:"tokens_per_hour,omitempty"`
	Estimated       bool         `json:"estimated,omitempty"`
	Models          []ModelSpend `json:"models"`
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
}

func perHour(tokens, seconds int64) *float64 {
	return over(float64(tokens), float64(seconds)/3600)
}

func perMinute(tokens, seconds int64) *float64 {
	return over(float64(tokens), float64(seconds)/60)
}

// itemUsage is what an item's usage comes to; nil when it has none.
func itemUsage(u *usage.Usage) *ItemUsage {
	if u.Empty() {
		return nil
	}
	out := &ItemUsage{Source: u.Source, Tokens: u.Tokens(), Cost: u.Cost(), Seconds: u.Seconds, Estimated: u.Estimated, Models: []ModelSpend{}}
	out.TokensPerMinute, out.TokensPerHour = perMinute(out.Tokens, u.Seconds), perHour(out.Tokens, u.Seconds)
	for _, m := range u.Models {
		out.Models = append(out.Models, ModelSpend{Model: m.Model, Tokens: m.Tokens(), Cost: m.Cost,
			TokensPerMinute: perMinute(m.Tokens(), u.Seconds), TokensPerHour: perHour(m.Tokens(), u.Seconds)})
	}
	return out
}

// spendReport totals the usage of the items done in the window, and lays
// them out against time and cost.
func spendReport(items []*workitem.Item, start, now time.Time) UsageReport {
	rep := UsageReport{Models: []ModelSpend{}, Done: []SpendPoint{}, ByModel: map[string][]SpendPoint{}}
	list := spent(items, start, now)
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
