package workitem

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Ordering policies (S-0217). A policy computes an order for a set of
// stories from their planning data: cost of delay, weighted shortest job
// first, shortest first, or first in first out. A story without the figure
// its policy needs goes after those with it, and ties keep the order the
// stories came in, which is their current board order.

// The ordering policies, named as orchestration.policy names them.
const (
	PolicyCOD        = manifest.OrderCOD
	PolicyWSJF       = manifest.OrderWSJF
	PolicyThroughput = manifest.OrderThroughput
	PolicyFIFO       = manifest.OrderFIFO
)

// OrderPolicies are the ordering policies, in the order they are offered.
var OrderPolicies = manifest.OrderPolicies

// Ranked is a story's place in an order a policy computed, and the figure
// it was placed by.
type Ranked struct {
	// Position counts from 1.
	Position int    `json:"position"`
	ID       string `json:"id"`
	Title    string `json:"title"`
	// Figure is what the story was ordered by: the cost of delay value per
	// week for cod, that value per hour of forecast duration for wsjf, the
	// forecast duration in hours for throughput, and created in Unix seconds
	// for fifo. There is none when the story lacks it.
	Figure *float64 `json:"figure,omitempty"`
	// Unit names what Figure counts.
	Unit string `json:"unit"`
	// Text is the figure as a person reads it.
	Text string `json:"text,omitempty"`
	// Missing names the figure the story lacks, when it does.
	Missing string `json:"missing,omitempty"`
}

// OrderByPolicy orders stories, given in their current order, by a policy.
func OrderByPolicy(items []*Item, policy string) ([]Ranked, error) {
	var figure func(*Item) (*float64, string, string)
	var unit string
	higherFirst := false
	switch policy {
	case PolicyCOD:
		figure, unit, higherFirst = codFigure, "per week", true
	case PolicyWSJF:
		figure, unit, higherFirst = wsjfFigure, "per hour", true
	case PolicyThroughput:
		figure, unit = throughputFigure, "hours"
	case PolicyFIFO:
		figure, unit = fifoFigure, "unix seconds"
	default:
		return nil, fmt.Errorf("unknown order policy %q: use one of %s", policy, strings.Join(OrderPolicies, ", "))
	}
	out := make([]Ranked, len(items))
	for i, it := range items {
		f, text, missing := figure(it)
		out[i] = Ranked{ID: it.ID, Title: it.Title, Figure: f, Unit: unit, Text: text, Missing: missing}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].Figure, out[j].Figure
		switch {
		case a == nil || b == nil:
			return a != nil && b == nil
		case higherFirst:
			return *a > *b
		default:
			return *a < *b
		}
	})
	for i := range out {
		out[i].Position = i + 1
	}
	return out, nil
}

func codValue(it *Item) *float64 {
	if it.CostOfDelay == nil || it.CostOfDelay.Value == nil {
		return nil
	}
	return it.CostOfDelay.Value
}

// forecastHours is the forecast duration in hours, when there is one longer
// than zero.
func forecastHours(it *Item) (float64, bool) {
	if it.Forecast == nil || it.Forecast.Duration == "" {
		return 0, false
	}
	d, err := time.ParseDuration(it.Forecast.Duration)
	if err != nil || d <= 0 {
		return 0, false
	}
	return d.Hours(), true
}

func codFigure(it *Item) (*float64, string, string) {
	v := codValue(it)
	if v == nil {
		return nil, "", "cost of delay value"
	}
	x := *v
	return &x, formatFigure(x), ""
}

func wsjfFigure(it *Item) (*float64, string, string) {
	v := codValue(it)
	h, ok := forecastHours(it)
	switch {
	case v == nil && !ok:
		return nil, "", "cost of delay value and forecast duration"
	case v == nil:
		return nil, "", "cost of delay value"
	case !ok:
		return nil, "", "forecast duration"
	}
	x := *v / h
	return &x, formatFigure(x) + "/h", ""
}

func throughputFigure(it *Item) (*float64, string, string) {
	h, ok := forecastHours(it)
	if !ok {
		return nil, "", "forecast duration"
	}
	return &h, it.Forecast.Duration, ""
}

func fifoFigure(it *Item) (*float64, string, string) {
	t, err := time.Parse(TimeFormat, it.Created)
	if err != nil {
		return nil, "", "created"
	}
	x := float64(t.Unix())
	return &x, it.Created, ""
}

// formatFigure writes a figure to at most two decimal places.
func formatFigure(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// ApplyReadyOrder makes the ready column's pull order ids, through Place, as
// dragging each story to the bottom in turn would.
func (b *Board) ApplyReadyOrder(items []*Item, ids []string) error {
	for _, id := range ids {
		if err := b.Place(items, id, Placement{Bottom: true}); err != nil {
			return err
		}
	}
	return nil
}
