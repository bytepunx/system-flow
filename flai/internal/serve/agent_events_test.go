package serve

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// agentEvents are the agent events in the project's event log.
func agentEvents(repo *workitem.Repo) []inbox.LoggedEvent {
	var out []inbox.LoggedEvent
	for _, e := range inbox.LoggedEvents(repo) {
		if e.Kind == inbox.Agent {
			out = append(out, e)
		}
	}
	return out
}

// S-0532: the helper writes one agent event per call, on the item or the
// role the run's kind says, by the run's agent, with its role and why.
func TestAppendAgentEventWritesOneEventOnTheRunsItemOrRole(t *testing.T) {
	lab := newReadyEmptyLab(t, 2)
	story := lab.ready("Events")
	loose := strings.Replace(story, "-0", "-", 1)
	var warned []string
	warn := func(msg string, args ...any) { warned = append(warned, fmt.Sprint(append([]any{msg}, args...)...)) }
	for i, tc := range []struct {
		run                      *AgentRun
		to, why                  string
		id, typ, title, by, role string
	}{
		{&AgentRun{Story: loose, Agent: "agent-" + story}, inbox.AgentStarted, "pulled", story, workitem.Story, "Events", "agent-" + story, conventions.RoleStory},
		{&AgentRun{Item: lab.epic.ID, Agent: "planner-" + lab.epic.ID}, inbox.AgentEnded, "exit 0", lab.epic.ID, workitem.Epic, "Epic", "planner-" + lab.epic.ID, conventions.RolePlan},
		{&AgentRun{Agent: workitem.ActivityOrchestrator}, inbox.AgentFailed, "exit 1", workitem.ActivityOrchestrator, inbox.RoleType, workitem.ActivityOrchestrator, workitem.ActivityOrchestrator, conventions.RoleOrchestrate},
		{&AgentRun{Agent: workitem.ActivityAnalyzer, Focus: harness.AllFocus}, inbox.AgentStopped, "stopped by the operator", workitem.ActivityAnalyzer, inbox.RoleType, workitem.ActivityAnalyzer, workitem.ActivityAnalyzer, conventions.RoleAnalyze},
		{&AgentRun{Agent: workitem.ActivityCurator, Role: conventions.RoleCurate}, inbox.AgentStarted, "", workitem.ActivityCurator, inbox.RoleType, workitem.ActivityCurator, workitem.ActivityCurator, conventions.RoleCurate},
		{&AgentRun{Agent: "reviewer-1", Role: "reviewer"}, inbox.AgentEnded, "exit 0", "reviewer", inbox.RoleType, "reviewer", "reviewer-1", "reviewer"},
		// a role a project defines (S-0417): on its item, or on its name
		{&AgentRun{Item: lab.epic.ID, Agent: "auditor-" + lab.epic.ID, Role: "audit", RoleName: "auditor"}, inbox.AgentStarted, "the operator asked for it", lab.epic.ID, workitem.Epic, "Epic", "auditor-" + lab.epic.ID, "audit"},
		{&AgentRun{Agent: "auditor", Role: "audit", RoleName: "auditor"}, inbox.AgentFailed, "exit 1", "auditor", inbox.RoleType, "auditor", "auditor", "audit"},
		{&AgentRun{Agent: "curator-two", Role: conventions.RoleCurate, RoleName: "curator-two"}, inbox.AgentEnded, "exit 0", "curator-two", inbox.RoleType, "curator-two", "curator-two", conventions.RoleCurate},
		// an item that cannot be read is named by its ID, its type from it
		{&AgentRun{Story: "S-9999", Agent: "agent-S-9999"}, inbox.AgentFailed, "exit 2", "S-9999", workitem.Story, "", "agent-S-9999", conventions.RoleStory},
	} {
		appendAgentEvent(lab.root, tc.run, tc.to, tc.why, warn)
		got := agentEvents(lab.repo)
		if len(got) != i+1 {
			t.Fatalf("call %d: one event per call, %d in the log", i+1, len(got))
		}
		e := got[i]
		if e.ID != tc.id || e.Type != tc.typ || e.Title != tc.title || e.To != tc.to || e.By != tc.by || e.Fields[inbox.FieldRole] != tc.role || e.Fields[inbox.FieldWhy] != tc.why || e.At == "" {
			t.Errorf("call %d: got %+v, want %s %s %q to %s by %s, role %s, why %q", i+1, e, tc.id, tc.typ, tc.title, tc.to, tc.by, tc.role, tc.why)
		}
	}
	if len(warned) != 0 {
		t.Errorf("nothing warned: %v", warned)
	}
}

// S-0532, ADR-0211 §4: an event that cannot be appended is warned of, and
// the caller goes on.
func TestAppendAgentEventWarnsWhenTheLogCannotBeWritten(t *testing.T) {
	lab := newReadyEmptyLab(t, 2)
	// a folder where the log is, so it cannot be opened to append
	if err := os.MkdirAll(inbox.EventsPath(lab.repo), 0o700); err != nil {
		t.Fatal(err)
	}
	var warned []string
	appendAgentEvent(lab.root, &AgentRun{Agent: workitem.ActivityOrchestrator}, inbox.AgentStarted, "", func(msg string, _ ...any) { warned = append(warned, msg) })
	if len(warned) != 1 || !strings.Contains(warned[0], "agent event not logged") {
		t.Errorf("one warning: %v", warned)
	}
}

// S-0532: a run whose log ends on Claude Code's usage-limit refusal is one
// failed event, whatever it was judged, whose why names the limit and when
// it resets, then the run's own why when it says more, and its exit.
func TestARunOnTheUsageLimitIsAFailedEventNamingTheLimitAndItsReset(t *testing.T) {
	lab := newReadyEmptyLab(t, 2)
	story := lab.ready("Limited")
	log := limitLog(t, limitStream...)
	one, zero := 1, 0
	for i, tc := range []struct {
		name string
		run  *AgentRun
		want string
	}{
		{"failed", &AgentRun{Story: story, Agent: "agent-" + story, Log: log, Exit: &one, Outcome: OutcomeFailed, Why: "ended (exit 1) with " + story + " in in-progress"},
			limitWhyText + "; ended (exit 1) with " + story + " in in-progress"},
		{"asked", &AgentRun{Story: story, Agent: "agent-" + story, Log: log, Exit: &zero, Outcome: OutcomeAsked, Why: "waiting for an answer to TH-0001: Which?"},
			limitWhyText + "; waiting for an answer to TH-0001: Which? (exit 0)"},
		{"worked", &AgentRun{Story: story, Agent: "agent-" + story, Log: log, Exit: &one, Outcome: OutcomeWorked}, limitWhyText + " (exit 1)"},
		// the curator's run says the limit as its why already, and it is not said twice
		{"curator", &AgentRun{Agent: workitem.ActivityCurator, Role: conventions.RoleCurate, Log: log, Exit: &one, Outcome: OutcomeFailed, Why: limitWhyText},
			limitWhyText + " (exit 1)"},
	} {
		var w warned
		appendAgentEnd(lab.root, tc.run, w.warn)
		got := agentEvents(lab.repo)
		if len(got) != i+1 {
			t.Fatalf("%s: one event per end, %d in the log", tc.name, len(got))
		}
		if e := got[i]; e.To != inbox.AgentFailed || e.Fields[inbox.FieldWhy] != tc.want || len(w) != 0 {
			t.Errorf("%s: event %+v, warned %v; want failed for %q", tc.name, e, w, tc.want)
		}
	}
}

// S-0532: a refusal whose log names no reset time says so.
func TestARunOnAUsageLimitWithNoResetTimeSaysItNamedNone(t *testing.T) {
	lab := newReadyEmptyLab(t, 2)
	story := lab.ready("Limited")
	var stream []string
	for _, line := range limitStream {
		if !strings.Contains(line, `"type":"rate_limit_event"`) {
			stream = append(stream, line)
		}
	}
	one := 1
	run := &AgentRun{Story: story, Agent: "agent-" + story, Log: limitLog(t, stream...), Exit: &one, Outcome: OutcomeFailed}
	var w warned
	appendAgentEnd(lab.root, run, w.warn)
	got := inbox.AgentEventsOf(lab.repo, story, time.Time{})
	want := "ended on the harness's usage limit (usage_limit_reached, 429); it named no reset time (exit 1)"
	if len(got) != 1 || got[0].To != inbox.AgentFailed || got[0].Fields[inbox.FieldWhy] != want || len(w) != 0 {
		t.Errorf("events %+v, warned %v; want one failed for %q", got, w, want)
	}
}

// S-0532: a run on no limit comes to what it was judged, for its own why,
// a failure's with its exit beside it, and a run the operator stopped stays
// stopped though its log ends on the limit.
func TestARunsEndIsItsOutcomeForItsOwnWhy(t *testing.T) {
	worked := limitLog(t, streamCall("s", "m1", runStart, 999), streamFinal("Done.", 2000, 0.5))
	limited := limitLog(t, limitStream...)
	one, zero := 1, 0
	for _, tc := range []struct {
		name    string
		run     AgentRun
		to, why string
	}{
		{"failed, exit named", AgentRun{Log: worked, Exit: &one, Outcome: OutcomeFailed, Why: "ended (exit 1) planning E-0001"}, inbox.AgentFailed, "ended (exit 1) planning E-0001"},
		{"failed, exit not named", AgentRun{Log: worked, Exit: &one, Outcome: OutcomeFailed, Why: "the merge failed"}, inbox.AgentFailed, "the merge failed (exit 1)"},
		{"failed, no why", AgentRun{Exit: &one, Outcome: OutcomeFailed}, inbox.AgentFailed, "ended (exit 1)"},
		{"failed, no exit seen", AgentRun{Outcome: OutcomeFailed}, inbox.AgentFailed, "ended (an exit code nobody saw)"},
		{"could not start", AgentRun{Outcome: OutcomeFailed, Why: "could not be started: no such file"}, inbox.AgentFailed, "could not be started: no such file"},
		{"worked", AgentRun{Log: worked, Exit: &zero, Outcome: OutcomeWorked}, inbox.AgentEnded, "ended with its work done"},
		{"asked", AgentRun{Log: worked, Exit: &zero, Outcome: OutcomeAsked, Why: "waiting for an answer to TH-0001: Which?"}, inbox.AgentEnded, "waiting for an answer to TH-0001: Which?"},
		{"asked, no why", AgentRun{Exit: &zero, Outcome: OutcomeAsked}, inbox.AgentEnded, "ended waiting for an answer"},
		{"stopped", AgentRun{Log: worked, Outcome: OutcomeStopped, Why: "stopped by the operator at 2026-10-28T09:00:00Z"}, inbox.AgentStopped, "stopped by the operator at 2026-10-28T09:00:00Z"},
		{"stopped on a limit", AgentRun{Log: limited, Exit: &one, Outcome: OutcomeStopped, Why: "stopped by the operator at 2026-10-28T09:00:00Z"}, inbox.AgentStopped, "stopped by the operator at 2026-10-28T09:00:00Z"},
		{"stopped, no why", AgentRun{Outcome: OutcomeStopped}, inbox.AgentStopped, "stopped by the operator"},
	} {
		var w warned
		if to, why := endEvent(&tc.run, w.warn); to != tc.to || why != tc.why || len(w) != 0 {
			t.Errorf("%s: %s for %q, warned %v; want %s for %q", tc.name, to, why, w, tc.to, tc.why)
		}
	}
}
