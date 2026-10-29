package serve

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// A run that ends on the harness's usage limit (S-0514, ADR-0214): its log
// says so as data, which the harness's usage reader reads as the run's
// Limit: the error the harness names, its status, and when the limit resets.

// runLimit is the usage limit the log of run ends on, with the harness's
// final text, said, or nil when its log ends otherwise, it has none, or it is
// missing. A log that cannot be read is warned of with warn and read as
// ending on none. The log is read as Claude Code's, the harness that runs a
// strategic agent from its definition (Capabilities.StrategicAgents), as its
// activity is (activity.go).
func runLimit(run *AgentRun, warn func(msg string, args ...any)) (lim *usage.Limit, said string) {
	if run == nil || run.Log == "" {
		return nil, ""
	}
	rec, err := usage.Read(harness.Adapters[harness.ClaudeCode].Reader(), run.Log)
	if err != nil {
		warn("run's log not read for a usage limit", "log", run.Log, "err", err)
		return nil, ""
	}
	if len(rec.Runs) == 0 {
		return nil, ""
	}
	last := rec.Runs[len(rec.Runs)-1]
	if last.Limit == nil {
		return nil, ""
	}
	return last.Limit, last.Result
}

// limitResets is when lim resets, as an RFC 3339 time in UTC, or "" when the
// harness named no time.
func limitResets(lim *usage.Limit) string {
	if lim.Resets.IsZero() {
		return ""
	}
	return lim.Resets.UTC().Format(time.RFC3339)
}

// limitWhy says that a run ended on lim: the error and the status the
// harness named, in parentheses, left out when it named neither, and when
// the limit resets, or that the harness named no reset time. Such as "ended
// on the harness's usage limit (usage_limit_reached, 429); it resets at
// 2026-10-28T08:00:00Z".
func limitWhy(lim *usage.Limit) string {
	var named []string
	if lim.Error != "" {
		named = append(named, lim.Error)
	}
	if lim.Status != 0 {
		named = append(named, strconv.Itoa(lim.Status))
	}
	why := "ended on the harness's usage limit"
	if len(named) > 0 {
		why += " (" + strings.Join(named, ", ") + ")"
	}
	if resets := limitResets(lim); resets != "" {
		return fmt.Sprintf("%s; it resets at %s", why, resets)
	}
	return why + "; it named no reset time"
}
