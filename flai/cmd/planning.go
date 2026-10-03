package cmd

import (
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planningLines are an item's draft flag, who finalized it (S-0201), cost of
// delay, and forecast (S-0199) as flai show and flai edit --show print them,
// one line each and only when set; amounts are in currency.
func planningLines(draft bool, fin *workitem.Finalized, c *workitem.CostOfDelay, f *workitem.Forecast, currency string) []string {
	var out []string
	if draft {
		out = append(out, "draft: yes, finalized before it is ready")
	}
	if !fin.IsZero() {
		out = append(out, "finalized by "+orDefault(fin.By, "?")+" at "+orDefault(fin.At, "?"))
	}
	if !c.IsZero() {
		amount := func(v float64) string {
			return strconv.FormatFloat(v, 'f', -1, 64) + " " + currency + "/week"
		}
		// the inputs and the value each say who set them (ADR-0079)
		var parts []string
		if in := c.Inputs; !in.IsZero() {
			var inputs []string
			if in.RevenuePerWeek != nil {
				inputs = append(inputs, "revenue "+amount(*in.RevenuePerWeek))
			}
			if in.PenaltyPerWeek != nil {
				inputs = append(inputs, "penalty "+amount(*in.PenaltyPerWeek))
			}
			if in.TimeLostPerCycle != "" {
				inputs = append(inputs, "time lost "+in.TimeLostPerCycle+"/cycle")
			}
			parts = append(parts, "inputs "+strings.Join(inputs, ", ")+" · "+setBy(in.By, in.At))
		}
		if c.Value != nil {
			parts = append(parts, "value "+amount(*c.Value)+" · "+setBy(c.By, c.At))
		}
		if c.Stale() {
			parts = append(parts, "stale: the inputs changed after the value")
		}
		out = append(out, "cost of delay: "+strings.Join(parts, "; "))
	}
	if !f.IsZero() {
		var parts []string
		if f.Duration != "" {
			parts = append(parts, "duration "+f.Duration)
		}
		if f.Delivery != "" {
			parts = append(parts, "delivery "+f.Delivery)
		}
		if f.Basis != "" {
			parts = append(parts, "basis: "+f.Basis)
		}
		out = append(out, "forecast: "+strings.Join(append(parts, setBy(f.By, f.At)), " · "))
	}
	return out
}

func setBy(by, at string) string {
	return "set by " + orDefault(by, "?") + " at " + orDefault(at, "?")
}
