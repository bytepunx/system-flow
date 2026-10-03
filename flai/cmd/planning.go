package cmd

import (
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// planningLines are an item's draft flag, cost of delay, and forecast
// (S-0199) as flai show and flai edit --show print them, one line each and
// only when set; amounts are in currency.
func planningLines(draft bool, c *workitem.CostOfDelay, f *workitem.Forecast, currency string) []string {
	var out []string
	if draft {
		out = append(out, "draft: yes, finalized before it is ready")
	}
	if !c.IsZero() {
		amount := func(v float64) string {
			return strconv.FormatFloat(v, 'f', -1, 64) + " " + currency + "/week"
		}
		var parts []string
		if c.Value != nil {
			parts = append(parts, "value "+amount(*c.Value))
		}
		if in := c.Inputs; in != nil {
			if in.RevenuePerWeek != nil {
				parts = append(parts, "revenue "+amount(*in.RevenuePerWeek))
			}
			if in.PenaltyPerWeek != nil {
				parts = append(parts, "penalty "+amount(*in.PenaltyPerWeek))
			}
			if in.TimeLostPerCycle != "" {
				parts = append(parts, "time lost "+in.TimeLostPerCycle+"/cycle")
			}
		}
		out = append(out, "cost of delay: "+strings.Join(append(parts, setBy(c.By, c.At)), " · "))
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
