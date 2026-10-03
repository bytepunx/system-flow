package workitem

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Planning data (S-0199): a story is a draft until it is finalized, epics and
// stories carry a cost of delay, and a story carries a forecast. Amounts are
// in the project's currency (planning.currency in system-flow.yaml), and each
// block records who set it and when; a cost of delay records it for its
// inputs and for its value apart (ADR-0080). A story that was a draft records
// who finalized it and when (S-0201).

// CostOfDelay is what waiting for an item costs: the inputs it is worked out
// from, the value, or both.
type CostOfDelay struct {
	Inputs *CostInputs `yaml:"inputs" json:"inputs,omitempty"`
	// Value is the cost of delay per week, in the project's currency.
	Value *float64 `yaml:"value" json:"value,omitempty"`
	// By and At are who set the value and when; there are none without one.
	By string `yaml:"by" json:"by,omitempty"`
	At string `yaml:"at" json:"at,omitempty"`
}

// CostInputs are what a cost of delay is worked out from; an absent amount is
// unknown, not zero.
type CostInputs struct {
	// RevenuePerWeek is the revenue the item brings in each week it is done.
	RevenuePerWeek *float64 `yaml:"revenue_per_week" json:"revenue_per_week,omitempty"`
	// PenaltyPerWeek is what each week it is not done costs beyond revenue.
	PenaltyPerWeek *float64 `yaml:"penalty_per_week" json:"penalty_per_week,omitempty"`
	// TimeLostPerCycle is the work lost each cycle it is not done, a Go duration.
	TimeLostPerCycle string `yaml:"time_lost_per_cycle" json:"time_lost_per_cycle,omitempty"`
	// By and At are who last added, changed, or removed an input, and when.
	By string `yaml:"by" json:"by,omitempty"`
	At string `yaml:"at" json:"at,omitempty"`
}

// Forecast is when a story is expected to be delivered, and why.
type Forecast struct {
	// Duration is the work it is expected to take, a Go duration.
	Duration string `yaml:"duration" json:"duration,omitempty"`
	// Delivery is when it is expected to be done, a UTC timestamp.
	Delivery string `yaml:"delivery" json:"delivery,omitempty"`
	// Basis is what the forecast rests on, in one sentence.
	Basis string `yaml:"basis" json:"basis,omitempty"`
	By    string `yaml:"by" json:"by,omitempty"`
	At    string `yaml:"at" json:"at,omitempty"`
}

// Finalized is who cleared a story's draft flag, and when (S-0201).
type Finalized struct {
	By string `yaml:"by" json:"by,omitempty"`
	At string `yaml:"at" json:"at,omitempty"`
}

// IsZero reports whether the inputs give no amount or duration, whoever is
// recorded as setting them.
func (in *CostInputs) IsZero() bool {
	return in == nil || (in.RevenuePerWeek == nil && in.PenaltyPerWeek == nil && in.TimeLostPerCycle == "")
}

// stamped reports whether the inputs record who set them or when.
func (in *CostInputs) stamped() bool {
	return in != nil && (in.By != "" || in.At != "")
}

// IsZero reports whether the cost of delay is absent or empty.
func (c *CostOfDelay) IsZero() bool {
	return c == nil || (c.Inputs.IsZero() && !c.Inputs.stamped() && c.Value == nil && c.By == "" && c.At == "")
}

// Stale reports whether the inputs changed after the value was set: there
// are both, and the inputs' at is later than the value's (ADR-0080).
func (c *CostOfDelay) Stale() bool {
	if c == nil || c.Value == nil || c.Inputs.IsZero() {
		return false
	}
	inputs, err := time.Parse(TimeFormat, c.Inputs.At)
	if err != nil {
		return false
	}
	value, err := time.Parse(TimeFormat, c.At)
	if err != nil {
		return false
	}
	return inputs.After(value)
}

// readOneStamp reads a cost of delay written with one by and at for the
// whole block (ADR-0074) as ADR-0080 stamps it: with no value the stamp is
// the inputs', and with a value it is the value's and the inputs' too, so
// that the value does not read as stale. A block whose inputs are stamped is
// left as it is.
func (c *CostOfDelay) readOneStamp() {
	if c == nil || c.Inputs.IsZero() || c.Inputs.stamped() || (c.By == "" && c.At == "") {
		return
	}
	c.Inputs.By, c.Inputs.At = c.By, c.At
	if c.Value == nil {
		c.By, c.At = "", ""
	}
}

// IsZero reports whether the forecast is absent or empty.
func (f *Forecast) IsZero() bool {
	return f == nil || *f == Forecast{}
}

func (c *CostOfDelay) clone() *CostOfDelay {
	if c == nil {
		return nil
	}
	out := *c
	out.Value = cloneAmount(c.Value)
	if c.Inputs != nil {
		in := *c.Inputs
		in.RevenuePerWeek = cloneAmount(c.Inputs.RevenuePerWeek)
		in.PenaltyPerWeek = cloneAmount(c.Inputs.PenaltyPerWeek)
		out.Inputs = &in
	}
	return &out
}

func (f *Forecast) clone() *Forecast {
	if f == nil {
		return nil
	}
	out := *f
	return &out
}

// IsZero reports whether the finalized block is absent or empty.
func (f *Finalized) IsZero() bool {
	return f == nil || *f == Finalized{}
}

func (f *Finalized) clone() *Finalized {
	if f == nil {
		return nil
	}
	out := *f
	return &out
}

// Finalize clears a story's draft flag and records who cleared it and when
// (S-0201); an item that is not a draft is left as it is.
func (it *Item) Finalize(by string, at time.Time) {
	if !it.Draft {
		return
	}
	it.Draft = false
	it.Finalized = &Finalized{By: orDefault(by, "agent"), At: at.UTC().Format(TimeFormat)}
}

func cloneAmount(v *float64) *float64 {
	if v == nil {
		return nil
	}
	x := *v
	return &x
}

// planningBlock is an item's draft flag, cost of delay, forecast, and
// finalized block as front matter, each only when set. It follows the fields
// an older flai knows, where such a flai writes them back, so that either
// writes the same bytes: finalized, which the flai of S-0199 keeps as a field
// it does not know, comes last.
func planningBlock(it *Item) string {
	var b strings.Builder
	if it.Draft {
		b.WriteString("draft: true\n")
	}
	if c := it.CostOfDelay; !c.IsZero() {
		b.WriteString("cost_of_delay:\n")
		if in := c.Inputs; !in.IsZero() {
			b.WriteString("  inputs:\n")
			writeAmount(&b, "    revenue_per_week", in.RevenuePerWeek)
			writeAmount(&b, "    penalty_per_week", in.PenaltyPerWeek)
			writeString(&b, "    time_lost_per_cycle", in.TimeLostPerCycle)
			writeString(&b, "    by", in.By)
			writeTime(&b, "    at", in.At)
		}
		writeAmount(&b, "  value", c.Value)
		writeString(&b, "  by", c.By)
		writeTime(&b, "  at", c.At)
	}
	if f := it.Forecast; !f.IsZero() {
		b.WriteString("forecast:\n")
		writeString(&b, "  duration", f.Duration)
		writeTime(&b, "  delivery", f.Delivery)
		writeString(&b, "  basis", f.Basis)
		writeString(&b, "  by", f.By)
		writeTime(&b, "  at", f.At)
	}
	if f := it.Finalized; !f.IsZero() {
		b.WriteString("finalized:\n")
		writeString(&b, "  by", f.By)
		writeTime(&b, "  at", f.At)
	}
	return b.String()
}

func writeAmount(b *strings.Builder, key string, v *float64) {
	if v != nil {
		fmt.Fprintf(b, "%s: %s\n", key, strconv.FormatFloat(*v, 'f', -1, 64))
	}
}

func writeString(b *strings.Builder, key, v string) {
	if v != "" {
		fmt.Fprintf(b, "%s: %s\n", key, Scalar(v))
	}
}

// writeTime writes a timestamp unquoted, as created and every at are.
func writeTime(b *strings.Builder, key, v string) {
	if v != "" {
		fmt.Fprintf(b, "%s: %s\n", key, v)
	}
}

// planningErrors are what is wrong with an item's draft flag, cost of delay,
// forecast, and finalized block.
func planningErrors(it *Item) []string {
	var errs []string
	if it.Draft && !Carries(it.Type, "draft") {
		errs = append(errs, "only a story is a draft, and this is "+articled(it.Type))
	}
	if c := it.CostOfDelay; !c.IsZero() {
		if !Carries(it.Type, "cost_of_delay") {
			errs = append(errs, "cost_of_delay is for stories and epics, and this is "+articled(it.Type))
		}
		if c.Inputs.IsZero() && c.Value == nil {
			errs = append(errs, "cost_of_delay has neither inputs nor a value: give one or both")
		}
		if in := c.Inputs; !in.IsZero() {
			errs = append(errs, amountErrors("cost_of_delay.inputs.revenue_per_week", in.RevenuePerWeek)...)
			errs = append(errs, amountErrors("cost_of_delay.inputs.penalty_per_week", in.PenaltyPerWeek)...)
			errs = append(errs, durationErrors("cost_of_delay.inputs.time_lost_per_cycle", in.TimeLostPerCycle)...)
			errs = append(errs, setByErrors("cost_of_delay.inputs", in.By, in.At)...)
		} else if in.stamped() {
			errs = append(errs, "cost_of_delay.inputs.by and at say who set the inputs, and there are none: give an input, or remove them")
		}
		if c.Value != nil {
			errs = append(errs, amountErrors("cost_of_delay.value", c.Value)...)
			errs = append(errs, setByErrors("cost_of_delay", c.By, c.At)...)
		} else if c.By != "" || c.At != "" {
			errs = append(errs, "cost_of_delay.by and at say who set the value, and there is none: give a value, or remove them")
		}
	}
	if f := it.Forecast; !f.IsZero() {
		if !Carries(it.Type, "forecast") {
			errs = append(errs, "forecast is for stories, and this is "+articled(it.Type))
		}
		if f.Duration == "" && f.Delivery == "" {
			errs = append(errs, "forecast has neither a duration nor a delivery: give one or both")
		}
		errs = append(errs, durationErrors("forecast.duration", f.Duration)...)
		if f.Delivery != "" {
			if _, err := time.Parse(TimeFormat, f.Delivery); err != nil {
				errs = append(errs, fmt.Sprintf("forecast.delivery %q is not a UTC timestamp like 2026-09-15T16:10:00Z", f.Delivery))
			}
		}
		if strings.ContainsAny(f.Basis, "\r\n") {
			errs = append(errs, "forecast.basis is more than one line: write it as one sentence")
		}
		errs = append(errs, setByErrors("forecast", f.By, f.At)...)
	}
	if f := it.Finalized; !f.IsZero() {
		if !Carries(it.Type, "finalized") {
			errs = append(errs, "finalized is for stories, and this is "+articled(it.Type))
		}
		if it.Draft {
			errs = append(errs, "finalized says who finalized the story, and it is still a draft: remove draft, or remove finalized if it is a draft again")
		}
		errs = append(errs, setByErrors("finalized", f.By, f.At)...)
	}
	return errs
}

// amountErrors refuses an amount that is not a finite number of zero or more.
func amountErrors(key string, v *float64) []string {
	switch {
	case v == nil:
		return nil
	case math.IsNaN(*v) || math.IsInf(*v, 0):
		return []string{fmt.Sprintf("%s is not a finite amount: write a number of zero or more in the project's currency", key)}
	case *v < 0:
		return []string{fmt.Sprintf("%s %s is negative: write a number of zero or more in the project's currency", key, strconv.FormatFloat(*v, 'f', -1, 64))}
	}
	return nil
}

// durationErrors refuses a duration, when given, that is not a Go duration
// longer than zero.
func durationErrors(key, v string) []string {
	if v == "" {
		return nil
	}
	if d, err := time.ParseDuration(v); err != nil || d <= 0 {
		return []string{fmt.Sprintf("%s %q is not a Go duration longer than zero like 4h or 90m", key, v)}
	}
	return nil
}

// setByErrors refuses a planning block that does not say who set it and when.
func setByErrors(key, by, at string) []string {
	var errs []string
	if strings.TrimSpace(by) == "" {
		errs = append(errs, key+".by is required: who set it")
	}
	if _, err := time.Parse(TimeFormat, at); err != nil {
		errs = append(errs, fmt.Sprintf("%s.at %q is not a UTC timestamp like 2026-09-15T16:10:00Z", key, at))
	}
	return errs
}
