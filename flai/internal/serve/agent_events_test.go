package serve

import (
	"fmt"
	"os"
	"strings"
	"testing"

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
