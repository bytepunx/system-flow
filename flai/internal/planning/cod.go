package planning

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// week is the period a cost of delay value is counted over.
const week = 7 * 24 * time.Hour

// Cost is an epic's or story's cost of delay per week and what it rests on.
type Cost struct {
	ID string `json:"id"`
	// Value is the cost of delay per week in Currency, to two decimals.
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
	// From is "inputs" when the value is worked out from the item's own
	// inputs, "epic" when it is the story's share of its epic's.
	From string `json:"from"`
	// Basis is what the value rests on, in one sentence.
	Basis  string      `json:"basis"`
	Inputs *CostInputs `json:"inputs,omitempty"`
	Epic   *EpicShare  `json:"epic,omitempty"`
}

// CostInputs are the inputs a value was worked out from, and the settings
// that turn time lost into an amount; absent ones are left out.
type CostInputs struct {
	RevenuePerWeek   *float64 `json:"revenue_per_week,omitempty"`
	PenaltyPerWeek   *float64 `json:"penalty_per_week,omitempty"`
	TimeLostPerCycle string   `json:"time_lost_per_cycle,omitempty"`
	// HourRate, Cycle, and CyclesPerWeek are given with time lost only.
	HourRate      *float64 `json:"hour_rate,omitempty"`
	Cycle         string   `json:"cycle,omitempty"`
	CyclesPerWeek float64  `json:"cycles_per_week,omitempty"`
}

// EpicShare is how a story's value was apportioned from its epic's.
type EpicShare struct {
	ID string `json:"id"`
	// Value is the epic's cost of delay per week, to two decimals.
	Value float64 `json:"value"`
	// From is "inputs" when the epic's value is worked out from its inputs,
	// "value" when it is the value recorded on it.
	From string `json:"from"`
	// Share is the fraction of the epic's value that is the story's.
	Share float64 `json:"share"`
	// Stories are the epic's open stories without inputs of their own that
	// the value is apportioned over, by duration.
	Stories []SharedStory `json:"stories"`
}

// SharedStory is a story the epic's value is apportioned over.
type SharedStory struct {
	ID              string `json:"id"`
	Duration        string `json:"duration"`
	DurationSeconds int64  `json:"duration_seconds"`
	// FromForecast says the duration is the story's own forecast.duration
	// rather than one worked out from history.
	FromForecast bool `json:"from_forecast"`
}

// CostOfDelay works out the cost of delay per week of epic or story id: from
// its own inputs, or, for a story without any, as its share of its epic's
// value apportioned by duration over the epic's open stories without inputs.
// items includes archived ones, which the durations are forecast from.
func CostOfDelay(items []*workitem.Item, plan manifest.Planning, id string) (Cost, error) {
	id = workitem.CanonicalID(id)
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	target := byID[id]
	fail := func(format string, args ...any) (Cost, error) {
		return Cost{}, fmt.Errorf("cannot work out the cost of delay of %s: %s", id, fmt.Sprintf(format, args...))
	}
	switch {
	case target == nil:
		return fail("there is no such item; check the ID with flai board")
	case target.Type != workitem.Epic && target.Type != workitem.Story:
		return fail("it is a %s, and only an epic or a story has a cost of delay", target.Type)
	case target.Archived:
		return fail("it is archived; work out the cost of delay of an item that is still open")
	case target.Closed():
		return fail("it is %s; work out the cost of delay of an item that is still open", target.Status)
	}
	currency := plan.CurrencyCode()
	res := Cost{ID: id, Currency: currency}
	if hasInputs(target) {
		value, in, basis, err := fromInputs(target, plan)
		if err != nil {
			return fail("%v", err)
		}
		res.Value, res.From, res.Inputs = value, "inputs", in
		res.Basis = fmt.Sprintf("%s: %.2f %s a week.", basis, value, currency)
		return res, nil
	}
	const give = "flai edit --revenue-per-week, --penalty-per-week, or --time-lost-per-cycle"
	if target.Type == workitem.Epic {
		return fail("%s has no cost of delay inputs; the operator gives them with %s", id, give)
	}
	epic := byID[target.Parent]
	if epic == nil || epic.Type != workitem.Epic {
		return fail("%s has no cost of delay inputs and no epic to take a share of; the operator gives its inputs with %s", id, give)
	}
	share := &EpicShare{ID: epic.ID}
	switch {
	case hasInputs(epic):
		value, _, _, err := fromInputs(epic, plan)
		if err != nil {
			return fail("%v", err)
		}
		share.Value, share.From = value, "inputs"
	case epic.CostOfDelay != nil && epic.CostOfDelay.Value != nil:
		share.Value, share.From = round2(*epic.CostOfDelay.Value), "value"
	default:
		return fail("%s has no cost of delay inputs, and its epic %s has neither inputs nor a value; the operator gives inputs on the story or the epic with %s", id, epic.ID, give)
	}
	fallback, err := plan.FallbackDuration()
	if err != nil {
		return fail("%v", err)
	}
	f := newForecaster(items, fallback)
	var own, total time.Duration
	for _, it := range sharers(items, epic.ID) {
		d, fromForecast := f.durationOf(it)
		share.Stories = append(share.Stories, SharedStory{ID: it.ID, Duration: FormatDuration(d), DurationSeconds: int64(d / time.Second), FromForecast: fromForecast})
		total += d
		if it.ID == id {
			own = d
		}
	}
	share.Share = float64(own) / float64(total)
	res.Value, res.From, res.Epic = round2(share.Value*share.Share), "epic", share
	res.Basis = fmt.Sprintf("%s's share of %s's %.2f %s a week, %s of %s forecast over its %d open %s without inputs: %.2f %s a week.",
		id, epic.ID, share.Value, currency, FormatDuration(own), FormatDuration(total), len(share.Stories), plural(len(share.Stories), "story", "stories"), res.Value, currency)
	return res, nil
}

// hasInputs reports whether an item's cost of delay has at least one input.
func hasInputs(it *workitem.Item) bool {
	return it.CostOfDelay != nil && !it.CostOfDelay.Inputs.IsZero()
}

// fromInputs works out an item's value per week from its inputs: revenue
// and penalty per week, plus the hours lost per cycle at the hour rate times
// the cycles in a week. It gives the inputs it used and a basis to end with
// the value.
func fromInputs(it *workitem.Item, plan manifest.Planning) (float64, *CostInputs, string, error) {
	src := it.CostOfDelay.Inputs
	currency := plan.CurrencyCode()
	in := &CostInputs{RevenuePerWeek: src.RevenuePerWeek, PenaltyPerWeek: src.PenaltyPerWeek, TimeLostPerCycle: src.TimeLostPerCycle}
	var value float64
	var parts []string
	if v := src.RevenuePerWeek; v != nil {
		value += *v
		parts = append(parts, fmt.Sprintf("%.2f %s of revenue a week", *v, currency))
	}
	if v := src.PenaltyPerWeek; v != nil {
		value += *v
		parts = append(parts, fmt.Sprintf("%.2f %s of penalty a week", *v, currency))
	}
	if src.TimeLostPerCycle != "" {
		lost, err := time.ParseDuration(src.TimeLostPerCycle)
		if err != nil || lost <= 0 {
			return 0, nil, "", fmt.Errorf("%s's time_lost_per_cycle %q is not a Go duration longer than zero; the operator corrects it with flai edit %s --time-lost-per-cycle", it.ID, src.TimeLostPerCycle, it.ID)
		}
		if plan.HourRate == nil {
			return 0, nil, "", fmt.Errorf("%s's time lost per cycle needs planning.hour_rate in system-flow.yaml, and it is unset; the operator sets it to what an hour of work costs in %s", it.ID, currency)
		}
		cycle, err := plan.CycleDuration()
		if err != nil {
			return 0, nil, "", err
		}
		perWeek := float64(week) / float64(cycle)
		value += lost.Hours() * *plan.HourRate * perWeek
		in.HourRate, in.Cycle, in.CyclesPerWeek = plan.HourRate, FormatDuration(cycle), perWeek
		parts = append(parts, fmt.Sprintf("%s of time lost per %s cycle at %s %s an hour, %.2f cycles a week",
			FormatDuration(lost), FormatDuration(cycle), strconv.FormatFloat(*plan.HourRate, 'f', -1, 64), currency, perWeek))
	}
	return round2(value), in, strings.Join(parts, " plus "), nil
}

// sharers are an epic's open stories without inputs of their own, by ID.
func sharers(items []*workitem.Item, epic string) []*workitem.Item {
	var out []*workitem.Item
	for _, it := range items {
		if it.Type == workitem.Story && it.Parent == epic && !it.Archived && !it.Closed() && !hasInputs(it) {
			out = append(out, it)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// round2 rounds an amount to two decimals.
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
