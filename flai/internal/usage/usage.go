// Package usage is what agents spent on a work item: tokens and cost per
// model, and the seconds they worked (S-0143), and apart from that what
// strategic agents spent on it (S-0225). The types here are what a work
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
	// Strategic is what strategic agents spent on the item, one entry per
	// kind in StrategicKinds order (S-0225). It is kept apart from the
	// agents' figures above: Tokens, Cost, Add, and Sum leave it out.
	Strategic []Strategic `yaml:"strategic,omitempty" json:"strategic,omitempty"`
}

// StrategicKinds are the kinds of strategic agent whose spending an item
// may carry, in the order it is written: workitem.ActivityKinds, which this
// package cannot import.
var StrategicKinds = []string{"planner", "orchestrator", "analyzer"}

// Strategic is what one kind of strategic agent spent on an item: the part
// of its activities' cost apportioned to the item, so always estimated.
type Strategic struct {
	Kind      string  `yaml:"kind" json:"kind"`
	Seconds   int64   `yaml:"seconds" json:"seconds"`
	Estimated bool    `yaml:"estimated" json:"estimated"`
	Models    []Model `yaml:"models" json:"models"`
}

// Tokens is every token the kind's models read or wrote.
func (s Strategic) Tokens() int64 {
	var n int64
	for _, m := range s.Models {
		n += m.Tokens()
	}
	return n
}

// Cost is what the kind's models cost, in US dollars.
func (s Strategic) Cost() float64 {
	var c float64
	for _, m := range s.Models {
		c += m.Cost
	}
	return c
}

// Empty says the agents spent nothing; what strategic agents spent does not
// count.
func (u *Usage) Empty() bool { return u == nil || (len(u.Models) == 0 && u.Seconds == 0) }

// Nothing says neither agents nor strategic agents spent anything.
func (u *Usage) Nothing() bool { return u.Empty() && (u == nil || len(u.Strategic) == 0) }

// Clone copies u and everything it refers to.
func (u *Usage) Clone() *Usage {
	if u == nil {
		return nil
	}
	c := *u
	c.Models = slices.Clone(u.Models)
	if u.Strategic != nil {
		c.Strategic = make([]Strategic, len(u.Strategic))
		for i, s := range u.Strategic {
			s.Models = slices.Clone(s.Models)
			c.Strategic[i] = s
		}
	}
	return &c
}

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

// StrategicTokens is every token strategic agents' models read or wrote.
func (u *Usage) StrategicTokens() int64 {
	if u == nil {
		return 0
	}
	var n int64
	for _, s := range u.Strategic {
		n += s.Tokens()
	}
	return n
}

// StrategicCost is what strategic agents' models cost, in US dollars.
func (u *Usage) StrategicCost() float64 {
	if u == nil {
		return 0
	}
	var c float64
	for _, s := range u.Strategic {
		c += s.Cost()
	}
	return c
}

// StrategicSeconds is how long strategic agents worked on the item.
func (u *Usage) StrategicSeconds() int64 {
	if u == nil {
		return 0
	}
	var n int64
	for _, s := range u.Strategic {
		n += s.Seconds
	}
	return n
}

// AddStrategic adds the agents' figures of o, its seconds and models, to
// what the strategic agent kind spent on u, as an estimate. Its own
// Strategic is left out.
func (u *Usage) AddStrategic(kind string, o *Usage) {
	if o.Empty() {
		return
	}
	i := slices.IndexFunc(u.Strategic, func(s Strategic) bool { return s.Kind == kind })
	if i < 0 {
		u.Strategic = append(u.Strategic, Strategic{Kind: kind})
		i = len(u.Strategic) - 1
	}
	s := &u.Strategic[i]
	s.Seconds += o.Seconds
	s.Estimated = true
	for _, m := range o.Models {
		s.Models = addModel(s.Models, m)
	}
	u.tidyStrategic()
}

// Add adds o to u, model by model; o's Strategic is left out.
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

func (u *Usage) addModel(m Model) { u.Models = addModel(u.Models, m) }

// addModel adds m to the model of its name in ms, or appends it.
func addModel(ms []Model, m Model) []Model {
	for i := range ms {
		if ms[i].Model == m.Model {
			ms[i].add(m)
			return ms
		}
	}
	return append(ms, m)
}

// Tidy sorts the models by name, drops those that spent nothing, and rounds
// each cost to a hundredth of a cent, so that the same usage is always
// written the same way; and does the same to each strategic entry, kept in
// StrategicKinds order.
func (u *Usage) Tidy() {
	u.Models = tidyModels(u.Models)
	u.tidyStrategic()
}

func tidyModels(ms []Model) []Model {
	ms = slices.DeleteFunc(ms, func(m Model) bool { return m.Tokens() == 0 && m.Cost == 0 })
	for i := range ms {
		ms[i].Cost = math.Round(ms[i].Cost*1e4) / 1e4
	}
	slices.SortFunc(ms, func(a, b Model) int { return strings.Compare(a.Model, b.Model) })
	if ms == nil {
		ms = []Model{}
	}
	return ms
}

// tidyStrategic tidies each strategic entry's models and orders the entries
// by kind, any kind StrategicKinds does not name last, by name.
func (u *Usage) tidyStrategic() {
	if len(u.Strategic) == 0 {
		return
	}
	for i := range u.Strategic {
		u.Strategic[i].Models = tidyModels(u.Strategic[i].Models)
	}
	rank := func(kind string) int {
		if i := slices.Index(StrategicKinds, kind); i >= 0 {
			return i
		}
		return len(StrategicKinds)
	}
	slices.SortStableFunc(u.Strategic, func(a, b Strategic) int {
		if d := rank(a.Kind) - rank(b.Kind); d != 0 {
			return d
		}
		return strings.Compare(a.Kind, b.Kind)
	})
}

// WithStrategic is u, the agents' figures of an item measured or summed
// again, carrying what strategic agents spent on it from old, its usage
// until now: a copy of u when old carries any, and when u is nil, empty
// agents' figures summed from nothing. Without old's Strategic it is u.
func WithStrategic(u, old *Usage) *Usage {
	if old == nil || len(old.Strategic) == 0 {
		return u
	}
	out := u.Clone()
	if out == nil {
		out = &Usage{Source: SourceSum, Models: []Model{}}
	}
	out.Strategic = old.Clone().Strategic
	return out
}

// Sum is the usage of items summed, with SourceSum and without Strategic;
// nil when none of them has any.
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

// Same says a and b say the same, what strategic agents spent included.
func Same(a, b *Usage) bool {
	if a.Nothing() || b.Nothing() {
		return a.Nothing() && b.Nothing()
	}
	if !slices.EqualFunc(a.Strategic, b.Strategic, func(x, y Strategic) bool {
		return x.Kind == y.Kind && x.Seconds == y.Seconds && x.Estimated == y.Estimated && slices.Equal(x.Models, y.Models)
	}) {
		return false
	}
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
