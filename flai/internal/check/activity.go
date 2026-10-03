package check

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// costTolerance is how far accrued_cost may stand from the sum of the log's
// costs: each is rounded to four decimals, so their sum may drift by less.
const costTolerance = 0.00005

// activities validates the strategic agents' activity documents that exist
// (ADR-0079). They are not narratives: narratives() reads only S-*.md, so
// they are checked here and nowhere else, save the markdown lint of wip.
func (c *checker) activities() {
	for _, kind := range workitem.ActivityKinds {
		path := c.repo.ActivityPath(kind)
		a, err := workitem.ReadActivity(path)
		switch {
		case errors.Is(err, os.ErrNotExist):
			continue
		case err != nil:
			msg := strings.TrimPrefix(err.Error(), path+": ")
			c.add(Error, "activity.document", path, 1, "%s; %s", msg, c.restore(path))
			continue
		}
		c.activityTotals(a)
	}
}

// activityTotals reports a negative total as an error and totals that
// disagree with the log as a warning.
func (c *checker) activityTotals(a *workitem.Activity) {
	p := a.Path
	negative := false
	for _, f := range []struct {
		key      string
		negative bool
		value    any
	}{
		{"accrued_cost", a.AccruedCost < 0 || math.IsNaN(a.AccruedCost), a.AccruedCost},
		{"accrued_seconds", a.AccruedSeconds < 0, a.AccruedSeconds},
		{"tasks_completed", a.TasksCompleted < 0, a.TasksCompleted},
	} {
		if f.negative {
			negative = true
			c.add(Error, "activity.front-matter", p, keyLine(p, f.key), "%s is %v but must be zero or more; %s", f.key, f.value, c.restore(p))
		}
	}
	if negative {
		return
	}
	var seconds int64
	var cost float64
	lastRun := ""
	for _, e := range a.Entries {
		seconds += e.Seconds
		cost += e.Cost
		if at := e.At.UTC().Format(workitem.TimeFormat); at > lastRun {
			lastRun = at
		}
	}
	if a.TasksCompleted != len(a.Entries) {
		c.add(Warning, "activity.totals", p, keyLine(p, "tasks_completed"), "tasks_completed is %d but the log has %d entries; %s", a.TasksCompleted, len(a.Entries), c.restore(p))
	}
	if a.AccruedSeconds != seconds {
		c.add(Warning, "activity.totals", p, keyLine(p, "accrued_seconds"), "accrued_seconds is %d but the log's seconds sum to %d; %s", a.AccruedSeconds, seconds, c.restore(p))
	}
	if math.Abs(a.AccruedCost-cost) > costTolerance {
		c.add(Warning, "activity.totals", p, keyLine(p, "accrued_cost"), "accrued_cost is %.4f but the log's costs sum to %.4f; %s", a.AccruedCost, cost, c.restore(p))
	}
	if a.LastRun != lastRun {
		c.add(Warning, "activity.totals", p, keyLine(p, "last_run"), "last_run is %q but the newest log entry ended at %q; %s", a.LastRun, lastRun, c.restore(p))
	}
}

// restore says how to repair an activity document: flai writes it, so it is
// restored rather than edited.
func (c *checker) restore(path string) string {
	rel := path
	if r, err := filepath.Rel(c.repo.Root, path); err == nil {
		rel = r
	}
	return "flai writes this document, so restore it from git (git checkout -- " + rel + ") rather than edit it by hand"
}
