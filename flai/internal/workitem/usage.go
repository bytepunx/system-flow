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

// usageBlock is an item's usage as front matter.
func usageBlock(u *usage.Usage) string {
	var b strings.Builder
	b.WriteString("usage:\n")
	fmt.Fprintf(&b, "  source: %s\n", Scalar(u.Source))
	fmt.Fprintf(&b, "  seconds: %d\n", u.Seconds)
	if u.Estimated {
		b.WriteString("  estimated: true\n")
	}
	if len(u.Models) == 0 {
		b.WriteString("  models: []\n")
		return b.String()
	}
	b.WriteString("  models:\n")
	for _, m := range u.Models {
		fmt.Fprintf(&b, "    - model: %s\n", Scalar(m.Model))
		fmt.Fprintf(&b, "      input: %d\n      output: %d\n      cache_read: %d\n      cache_write: %d\n", m.Input, m.Output, m.CacheRead, m.CacheWrite)
		fmt.Fprintf(&b, "      cost: %s\n", strconv.FormatFloat(m.Cost, 'f', -1, 64))
	}
	return b.String()
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
	seen := map[string]bool{}
	for i, m := range u.Models {
		switch {
		case strings.TrimSpace(m.Model) == "":
			errs = append(errs, fmt.Sprintf("usage.models[%d].model is required", i))
		case seen[m.Model]:
			errs = append(errs, fmt.Sprintf("usage.models[%d].model %s is listed twice", i, m.Model))
		}
		seen[m.Model] = true
		if m.Input < 0 || m.Output < 0 || m.CacheRead < 0 || m.CacheWrite < 0 || m.Cost < 0 {
			errs = append(errs, fmt.Sprintf("usage.models[%d] has a negative count or cost", i))
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
// those that changed. It returns their IDs.
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
			if sum := Summed(items, p.ID); !usage.Same(sum, p.Usage) {
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
