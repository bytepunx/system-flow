package metrics

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Spend over time (S-0163, ADR-0053): what agents spent, placed at the
// moment each item entered done, in buckets, for every item type at once.
// design/system/metrics.md defines it.

// The buckets a series is laid out in.
const (
	BucketHour = "hour"
	BucketDay  = "day"
	BucketWeek = "week"
	// DefaultBucket is the bucket when none is asked for.
	DefaultBucket = BucketDay
	// MaxHourWindow is the longest window laid out by the hour.
	MaxHourWindow = 31 * 24 * time.Hour
)

// CheckBucket says whether a bucket is one, and whether the window is short
// enough for it; an empty bucket is the default.
func CheckBucket(bucket string, window time.Duration) error {
	switch bucket {
	case "", BucketDay, BucketWeek:
		return nil
	case BucketHour:
		if window > MaxHourWindow {
			return fmt.Errorf("--bucket hour needs a window of 31 days or less; use --bucket day, or a shorter --since")
		}
		return nil
	}
	return fmt.Errorf("--bucket expects hour, day, or week, got %q", bucket)
}

// Spend is what was spent on a set of items, and what that comes to per
// item, per minute of agent work, and per dollar. A value whose divisor is
// zero is absent, and so are minutes per item when no item took agent time.
type Spend struct {
	Items           int      `json:"items"`
	Tokens          int64    `json:"tokens"`
	Cost            float64  `json:"cost"`
	Seconds         int64    `json:"seconds"`
	Estimated       bool     `json:"estimated,omitempty"`
	TokensPerItem   *float64 `json:"tokens_per_item,omitempty"`
	CostPerItem     *float64 `json:"cost_per_item,omitempty"`
	MinutesPerItem  *float64 `json:"minutes_per_item,omitempty"`
	TokensPerMinute *float64 `json:"tokens_per_minute,omitempty"`
	TokensPerDollar *float64 `json:"tokens_per_dollar,omitempty"`
}

// ModelShare is what one model spent on the items it worked on: its own
// tokens and cost, and those items' agent time.
type ModelShare struct {
	Model string `json:"model"`
	Spend
}

// Bucket is what was spent on the items done in one bucket of time.
type Bucket struct {
	// At is the bucket's start.
	At string `json:"at"`
	Spend
	// MeanTokens and MeanCost are the running means per bucket, from the
	// first bucket of the series to this one.
	MeanTokens float64 `json:"mean_tokens"`
	MeanCost   float64 `json:"mean_cost"`
	// Models are the same per model, in order of name.
	Models []ModelShare `json:"models,omitempty"`
}

// TypeSpend is what was spent on the items of one type done in the window,
// in all and over time.
type TypeSpend struct {
	Spend
	Models  []ModelShare `json:"models"`
	Buckets []Bucket     `json:"buckets"`
}

// tally sums a set of items, or one model's share of them.
type tally struct {
	items     int
	tokens    int64
	cost      float64
	seconds   int64
	estimated bool
}

func (t *tally) add(tokens int64, cost float64, seconds int64, estimated bool) {
	t.items++
	t.tokens += tokens
	t.cost += cost
	t.seconds += seconds
	t.estimated = t.estimated || estimated
}

func over(n, d float64) *float64 {
	if d <= 0 {
		return nil
	}
	r := n / d
	return &r
}

// minutesPerItem is the mean agent time an item took, in minutes (S-0169);
// none when no item took any.
func (t *tally) minutesPerItem() *float64 {
	if t.seconds <= 0 {
		return nil
	}
	return over(float64(t.seconds)/60, float64(t.items))
}

func (t *tally) spend() Spend {
	return Spend{
		Items: t.items, Tokens: t.tokens, Cost: t.cost, Seconds: t.seconds, Estimated: t.estimated,
		TokensPerItem:   over(float64(t.tokens), float64(t.items)),
		CostPerItem:     over(t.cost, float64(t.items)),
		MinutesPerItem:  t.minutesPerItem(),
		TokensPerMinute: over(float64(t.tokens), float64(t.seconds)/60),
		TokensPerDollar: over(float64(t.tokens), t.cost),
	}
}

// group is a set of items and each model's share of them.
type group struct {
	all    tally
	models map[string]*tally
}

func (g *group) add(it *workitem.Item) {
	u := it.Usage
	g.all.add(u.Tokens(), u.Cost(), u.Seconds, u.Estimated)
	if g.models == nil {
		g.models = map[string]*tally{}
	}
	for _, m := range u.Models {
		t := g.models[m.Model]
		if t == nil {
			t = &tally{}
			g.models[m.Model] = t
		}
		t.add(m.Tokens(), m.Cost, u.Seconds, u.Estimated)
	}
}

func (g *group) shares() []ModelShare {
	var out []ModelShare
	for name, t := range g.models {
		out = append(out, ModelShare{Model: name, Spend: t.spend()})
	}
	slices.SortFunc(out, func(a, b ModelShare) int { return strings.Compare(a.Model, b.Model) })
	return out
}

// bucketStart is the start of the bucket that holds t.
func bucketStart(t time.Time, bucket string) time.Time {
	t = t.UTC()
	switch bucket {
	case BucketHour:
		return t.Truncate(time.Hour)
	case BucketWeek:
		day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return day.AddDate(0, 0, -((int(day.Weekday()) + 6) % 7))
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// nextBucket is the start of the bucket after the one that starts at t.
func nextBucket(t time.Time, bucket string) time.Time {
	switch bucket {
	case BucketHour:
		return t.Add(time.Hour)
	case BucketWeek:
		return t.AddDate(0, 0, 7)
	}
	return t.AddDate(0, 0, 1)
}

// completion is an item done in the window that carries usage.
type completion struct {
	at time.Time
	it *workitem.Item
}

// spent lists the items done in the window that carry usage, oldest first,
// then by ID.
func spent(items []*workitem.Item, start, now time.Time) []completion {
	var list []completion
	for _, it := range items {
		at := it.FirstAt(workitem.Done)
		if at.IsZero() || at.Before(start) || at.After(now) || it.Usage.Empty() {
			continue
		}
		list = append(list, completion{at, it})
	}
	sort.SliceStable(list, func(i, j int) bool {
		if !list[i].at.Equal(list[j].at) {
			return list[i].at.Before(list[j].at)
		}
		return workitem.CanonicalID(list[i].it.ID) < workitem.CanonicalID(list[j].it.ID)
	})
	return list
}

// spendOverTime lays out what was spent on each type's items, whatever type
// the report is about.
func spendOverTime(all []*workitem.Item, start, now time.Time, bucket string) map[string]*TypeSpend {
	out := map[string]*TypeSpend{}
	for _, typ := range []string{workitem.Epic, workitem.Story, workitem.Task} {
		var items []*workitem.Item
		for _, it := range all {
			if it.Type == typ {
				items = append(items, it)
			}
		}
		out[typ] = typeSpend(spent(items, start, now), now, bucket)
	}
	return out
}

func typeSpend(list []completion, now time.Time, bucket string) *TypeSpend {
	ts := &TypeSpend{Models: []ModelShare{}, Buckets: []Bucket{}}
	if len(list) == 0 {
		return ts
	}
	var window group
	buckets := map[time.Time]*group{}
	for _, c := range list {
		window.add(c.it)
		at := bucketStart(c.at, bucket)
		g := buckets[at]
		if g == nil {
			g = &group{}
			buckets[at] = g
		}
		g.add(c.it)
	}
	ts.Spend = window.all.spend()
	ts.Models = window.shares()
	var tokens, cost float64
	n := 0
	last := bucketStart(now, bucket)
	for at := bucketStart(list[0].at, bucket); !at.After(last); at = nextBucket(at, bucket) {
		g := buckets[at]
		if g == nil {
			g = &group{}
		}
		n++
		tokens += float64(g.all.tokens)
		cost += g.all.cost
		ts.Buckets = append(ts.Buckets, Bucket{
			At: at.Format(workitem.TimeFormat), Spend: g.all.spend(),
			MeanTokens: tokens / float64(n), MeanCost: cost / float64(n), Models: g.shares(),
		})
	}
	return ts
}
