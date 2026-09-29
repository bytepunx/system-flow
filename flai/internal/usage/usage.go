// Package usage is what agents spent on a work item: tokens and cost per
// model, and the seconds they worked (S-0143). The types here are what a work
// item's front matter carries; log.go measures them from an agent's log.
package usage

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Sources of an item's usage.
const (
	// SourceLog is usage measured from the logs of the agents that worked it.
	SourceLog = "log"
	// SourceSum is usage summed from the item's children.
	SourceSum = "sum"
)

// Model is what one model spent.
type Model struct {
	Model      string `yaml:"model" json:"model"`
	Input      int64  `yaml:"input" json:"input"`
	Output     int64  `yaml:"output" json:"output"`
	CacheRead  int64  `yaml:"cache_read" json:"cache_read"`
	CacheWrite int64  `yaml:"cache_write" json:"cache_write"`
	// Cost is in US dollars.
	Cost float64 `yaml:"cost" json:"cost"`
}

// Tokens is every token the model read or wrote.
func (m Model) Tokens() int64 { return m.Input + m.Output + m.CacheRead + m.CacheWrite }

func (m *Model) add(o Model) {
	m.Input += o.Input
	m.Output += o.Output
	m.CacheRead += o.CacheRead
	m.CacheWrite += o.CacheWrite
	m.Cost += o.Cost
}

// scaled is m with every count multiplied by share, rounded to whole tokens.
func (m Model) scaled(share float64) Model {
	r := func(n int64) int64 { return int64(math.Round(float64(n) * share)) }
	return Model{Model: m.Model, Input: r(m.Input), Output: r(m.Output), CacheRead: r(m.CacheRead), CacheWrite: r(m.CacheWrite), Cost: m.Cost * share}
}

// Usage is what agents spent on one item.
type Usage struct {
	// Source says where it came from: SourceLog or SourceSum.
	Source string `yaml:"source" json:"source"`
	// Seconds is how long agents worked on it.
	Seconds int64 `yaml:"seconds" json:"seconds"`
	// Estimated says some of the cost was not reported by the harness but
	// estimated or apportioned.
	Estimated bool `yaml:"estimated,omitempty" json:"estimated,omitempty"`
	// Models are what each model spent, in order of name.
	Models []Model `yaml:"models" json:"models"`
}

// Empty says nothing was spent.
func (u *Usage) Empty() bool { return u == nil || (len(u.Models) == 0 && u.Seconds == 0) }

// Tokens is every token every model read or wrote.
func (u *Usage) Tokens() int64 {
	if u == nil {
		return 0
	}
	var n int64
	for _, m := range u.Models {
		n += m.Tokens()
	}
	return n
}

// Cost is what every model cost, in US dollars.
func (u *Usage) Cost() float64 {
	if u == nil {
		return 0
	}
	var c float64
	for _, m := range u.Models {
		c += m.Cost
	}
	return c
}

// Add adds o to u, model by model.
func (u *Usage) Add(o *Usage) {
	if o == nil {
		return
	}
	u.Seconds += o.Seconds
	u.Estimated = u.Estimated || o.Estimated
	for _, m := range o.Models {
		u.addModel(m)
	}
}

func (u *Usage) addModel(m Model) {
	for i := range u.Models {
		if u.Models[i].Model == m.Model {
			u.Models[i].add(m)
			return
		}
	}
	u.Models = append(u.Models, m)
}

// Tidy sorts the models by name, drops those that spent nothing, and rounds
// each cost to a hundredth of a cent, so that the same usage is always
// written the same way.
func (u *Usage) Tidy() {
	u.Models = slices.DeleteFunc(u.Models, func(m Model) bool { return m.Tokens() == 0 && m.Cost == 0 })
	for i := range u.Models {
		u.Models[i].Cost = math.Round(u.Models[i].Cost*1e4) / 1e4
	}
	slices.SortFunc(u.Models, func(a, b Model) int { return strings.Compare(a.Model, b.Model) })
	if u.Models == nil {
		u.Models = []Model{}
	}
}

// Sum is the usage of items summed, with SourceSum; nil when none of them
// has any.
func Sum(of ...*Usage) *Usage {
	var out *Usage
	for _, u := range of {
		if u.Empty() {
			continue
		}
		if out == nil {
			out = &Usage{Source: SourceSum}
		}
		out.Add(u)
	}
	if out != nil {
		out.Tidy()
	}
	return out
}

// Same says a and b say the same.
func Same(a, b *Usage) bool {
	if a.Empty() || b.Empty() {
		return a.Empty() && b.Empty()
	}
	return a.Source == b.Source && a.Seconds == b.Seconds && a.Estimated == b.Estimated && slices.Equal(a.Models, b.Models)
}

// Summary says in a line what was spent: tokens, cost, time, and where the
// numbers came from.
func (u *Usage) Summary() string {
	cost := fmt.Sprintf("$%.2f", u.Cost())
	if u.Estimated {
		cost += " (estimated)"
	}
	from := "measured from its agents' logs"
	if u.Source == SourceSum {
		from = "summed from its children"
	}
	return fmt.Sprintf("%s tokens · %s · %s of agent work · %s", Count(u.Tokens()), cost, (time.Duration(u.Seconds) * time.Second).String(), from)
}

// String is a model's line: its tokens by kind, and its cost.
func (m Model) String() string {
	return fmt.Sprintf("%s  input %s · output %s · cache read %s · cache write %s · $%.4f", m.Model, Count(m.Input), Count(m.Output), Count(m.CacheRead), Count(m.CacheWrite), m.Cost)
}

// Count is a count of tokens made short: 950, 12.3K, 20.1M.
func Count(n int64) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1e9)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 10_000:
		return fmt.Sprintf("%.1fK", float64(n)/1e3)
	}
	return strconv.FormatInt(n, 10)
}
