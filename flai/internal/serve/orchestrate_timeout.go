package serve

import (
	"fmt"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// An orchestrator run that holds past orchestration.run_timeout is stopped
// (ADR-0209 §4), so that a prompt that holds the run anyway still ends.
//
// The run is marked before the signal with the why of its stop, in Why, and
// not with Stopped, which is the operator's stop and which the launcher
// reads as theirs (launcher.ended). orchestrateEnded, whoever settles the
// run, finds the mark on the run recorded and records the run stopped for
// that why, its activity summary and its agent event's why included. A stop
// is no failure: orchestrateRetry does not hold the next start, and the
// reasons queued while the run went start the next run at the first look
// once it has ended, after their coalesce window.

// overduePrefix begins the why of a run stopped for orchestration.run_timeout.
const overduePrefix = "stopped: ran past orchestration.run_timeout"

// orchestrateOverdue is the why, and the activity, of a run stopped because it
// held past orchestration.run_timeout, timeout.
func orchestrateOverdue(timeout time.Duration) string {
	return fmt.Sprintf("%s (%s)", overduePrefix, timeoutText(timeout))
}

// stopOverdue stops run, the orchestrator's run going, once now is past its
// start and orchestration.run_timeout of o, as halt stops it, marked with
// orchestrateOverdue, and journals and logs the stop. A run_timeout that is
// not valid is warned of once and taken as manifest.DefaultRunTimeout. A run
// being stopped already, by the operator or by an earlier look, and a run
// whose start cannot be read, are left alone.
func (r *orchestrator) stopOverdue(o manifest.Orchestration, run *AgentRun, now time.Time) {
	if run.Stopped != "" || strings.HasPrefix(run.Why, overduePrefix) {
		return
	}
	started, err := time.Parse(time.RFC3339, run.Started)
	if err != nil {
		return
	}
	timeout, err := o.RunTimeoutDuration()
	if err != nil {
		r.warnOnce("orchestration run timeout not valid", err)
		timeout = manifest.DefaultRunTimeout
	}
	if !now.After(started.Add(timeout)) {
		return
	}
	why := orchestrateOverdue(timeout)
	h := r.halt(run, func(marked *AgentRun, _ string) { marked.Why = why })
	r.journalStop(h.at, fmt.Sprintf("stopped the orchestrator %s, since it ran past orchestration.run_timeout (%s): %s", run.Agent, timeoutText(timeout), h.what))
	r.starter.log("orchestrator stopped past its run timeout", "pid", run.PID, "run_timeout", timeoutText(timeout), "signalled", h.signalled, "killed", h.killed)
}

// overdueWhy is the why stopOverdue marked run with, read from recorded, the
// orchestrator's run in serve/agents.json; "" when recorded is another run
// or stopOverdue did not stop it.
func overdueWhy(recorded, run *AgentRun) string {
	if recorded.same(run) && strings.HasPrefix(recorded.Why, overduePrefix) {
		return recorded.Why
	}
	return ""
}

// timeoutText is a duration as a manifest writes it: 20m, not 20m0s.
func timeoutText(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}
