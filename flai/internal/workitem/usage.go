package workitem

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// What agents spent on an item (S-0143). A story's usage is measured from
// the logs of the agents that worked it, and each of its tasks' from the
// part of those logs the task was in progress for; an item nobody measured
// is given the sum of its children's. Whenever an item enters done, each
// item above it that was not measured is summed again, up to its epic.
//
// Apart from that, an item carries what strategic agents spent on it
// (S-0225): a planner's activity is charged to the item it planned and to
// every item above it at once, and neither a roll-up nor a measurement
// changes it.

// usageBlock is an item's usage as front matter.
func usageBlock(u *usage.Usage) string {
	var b strings.Builder
	b.WriteString("usage:\n")
	fmt.Fprintf(&b, "  source: %s\n", Scalar(u.Source))
	fmt.Fprintf(&b, "  seconds: %d\n", u.Seconds)
	if u.Estimated {
		b.WriteString("  estimated: true\n")
	}
	modelsBlock(&b, "  ", u.Models)
	StrategicBlock(&b, u.Strategic)
	return b.String()
}

// StrategicBlock writes a usage's strategic key and its entries, nested
// under usage; nothing when there are none.
func StrategicBlock(b *strings.Builder, entries []usage.Strategic) {
	if len(entries) == 0 {
		return
	}
	b.WriteString("  strategic:\n")
	for _, s := range entries {
		fmt.Fprintf(b, "    - kind: %s\n", Scalar(s.Kind))
		fmt.Fprintf(b, "      seconds: %d\n", s.Seconds)
		fmt.Fprintf(b, "      estimated: %t\n", s.Estimated)
		modelsBlock(b, "      ", s.Models)
	}
}

// modelsBlock writes a models key and its list, indented by indent.
func modelsBlock(b *strings.Builder, indent string, models []usage.Model) {
	if len(models) == 0 {
		fmt.Fprintf(b, "%smodels: []\n", indent)
		return
	}
	fmt.Fprintf(b, "%smodels:\n", indent)
	for _, m := range models {
		fmt.Fprintf(b, "%s  - model: %s\n", indent, Scalar(m.Model))
		fmt.Fprintf(b, "%[1]s    input: %[2]d\n%[1]s    output: %[3]d\n%[1]s    cache_read: %[4]d\n%[1]s    cache_write: %[5]d\n", indent, m.Input, m.Output, m.CacheRead, m.CacheWrite)
		fmt.Fprintf(b, "%s    cost: %s\n", indent, strconv.FormatFloat(m.Cost, 'f', -1, 64))
	}
}

// usageErrors are what is wrong with an item's usage.
func usageErrors(u *usage.Usage) []string {
	if u == nil {
		return nil
	}
	var errs []string
	if u.Source != usage.SourceLog && u.Source != usage.SourceSum {
		errs = append(errs, fmt.Sprintf("usage.source %q must be %s or %s", u.Source, usage.SourceLog, usage.SourceSum))
	}
	if u.Seconds < 0 {
		errs = append(errs, "usage.seconds is negative")
	}
	errs = append(errs, modelErrors("usage", u.Models)...)
	return append(errs, StrategicErrors(u.Strategic, "an item")...)
}

// StrategicErrors are what is wrong with a usage's strategic entries, on
// what is charged, such as "an item" or "an issue".
func StrategicErrors(entries []usage.Strategic, on string) []string {
	var errs []string
	kinds := map[string]bool{}
	for i, s := range entries {
		at := fmt.Sprintf("usage.strategic[%d]", i)
		switch {
		case !IsActivityKind(s.Kind):
			errs = append(errs, fmt.Sprintf("%s.kind %q must be one of %s", at, s.Kind, strings.Join(ActivityKinds, ", ")))
		case kinds[s.Kind]:
			errs = append(errs, fmt.Sprintf("%s.kind %s is listed twice", at, s.Kind))
		}
		kinds[s.Kind] = true
		if s.Seconds < 0 {
			errs = append(errs, at+".seconds is negative")
		}
		if !s.Estimated {
			errs = append(errs, at+".estimated must be true: a strategic agent's spending on "+on+" is apportioned")
		}
		errs = append(errs, modelErrors(at, s.Models)...)
	}
	return errs
}

// modelErrors are what is wrong with the models of the usage at path.
func modelErrors(path string, models []usage.Model) []string {
	var errs []string
	seen := map[string]bool{}
	for i, m := range models {
		switch {
		case strings.TrimSpace(m.Model) == "":
			errs = append(errs, fmt.Sprintf("%s.models[%d].model is required", path, i))
		case seen[m.Model]:
			errs = append(errs, fmt.Sprintf("%s.models[%d].model %s is listed twice", path, i, m.Model))
		}
		seen[m.Model] = true
		if m.Input < 0 || m.Output < 0 || m.CacheRead < 0 || m.CacheWrite < 0 || m.Cost < 0 {
			errs = append(errs, fmt.Sprintf("%s.models[%d] has a negative count or cost", path, i))
		}
	}
	return errs
}

// Summed is the sum of the usage of id's children, whatever their state;
// nil when none has any.
func Summed(items []*Item, id string) *usage.Usage {
	var of []*usage.Usage
	for _, c := range Children(items, id) {
		of = append(of, c.Usage)
	}
	return usage.Sum(of...)
}

// RollUp sums again the usage of each item above it that was not measured
// from an agent's log, up to its epic, archived items included, and saves
// those that changed. Each keeps what strategic agents spent on it. It
// returns their IDs.
func (r *Repo) RollUp(it *Item) ([]string, error) {
	if it.Parent == "" {
		return nil, nil
	}
	items, err := r.List(true)
	if err != nil {
		return nil, err
	}
	byID := map[string]*Item{}
	for _, x := range items {
		byID[x.ID] = x
	}
	// the item as the caller has it, which may not be saved yet
	for i, x := range items {
		if x.ID == it.ID {
			items[i] = it
		}
	}
	var changed []string
	seen := map[string]bool{it.ID: true}
	for id := it.Parent; id != "" && !seen[id]; {
		seen[id] = true
		p := byID[id]
		if p == nil {
			break
		}
		if p.Usage == nil || p.Usage.Source != usage.SourceLog {
			if sum := usage.WithStrategic(Summed(items, p.ID), p.Usage); !usage.Same(sum, p.Usage) {
				p.Usage = sum
				if err := r.Save(p); err != nil {
					return changed, err
				}
				changed = append(changed, p.ID)
			}
		}
		id = p.Parent
	}
	return changed, nil
}

// ChargeStrategic adds what the strategic agent kind spent on the item with
// id, the agents' figures of u as apportioned to it, to that item's
// strategic usage and to each item's above it, up to its epic, archived
// items included (S-0225). An item with no usage is given empty agents'
// figures summed from nothing. Each item is read again just before it is
// saved. It returns the IDs of the items it changed, the item's first and
// then upward; a u that spent nothing charges nothing.
func (r *Repo) ChargeStrategic(id, kind string, u *usage.Usage) ([]string, error) {
	if u.Empty() {
		return nil, nil
	}
	if !IsActivityKind(kind) {
		return nil, fmt.Errorf("%s is not a strategic agent's kind: %s", kind, strings.Join(ActivityKinds, ", "))
	}
	var changed []string
	seen := map[string]bool{}
	for id != "" && !seen[id] {
		seen[id] = true
		it, err := r.Get(id)
		if err != nil {
			return changed, err
		}
		if it.Usage == nil {
			it.Usage = &usage.Usage{Source: usage.SourceSum, Models: []usage.Model{}}
		}
		it.Usage.AddStrategic(kind, u)
		if err := r.Save(it); err != nil {
			return changed, err
		}
		changed = append(changed, it.ID)
		id = it.Parent
	}
	return changed, nil
}
