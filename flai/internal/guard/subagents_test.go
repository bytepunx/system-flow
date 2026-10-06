package guard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// started is a SubagentStart in session s1 for the sub-agent id.
func started(id, typ string) Event {
	return Event{HookEventName: HookSubagentStart, SessionID: "s1", AgentID: id, AgentType: typ}
}

// S-0285: Record adds a started sub-agent to its session's record and
// removes a stopped one, and removes the record once none runs.
func TestRecordKeepsTheSessionsRunningSubAgents(t *testing.T) {
	dir := filepath.Join(t.TempDir(), RecordDir)
	at := time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC)
	steps := []struct {
		e    Event
		want []SubAgent
	}{
		{started("a1", "general-purpose"), []SubAgent{{ID: "a1", Type: "general-purpose", Started: at}}},
		{started("a2", "verifier"), []SubAgent{{ID: "a1", Type: "general-purpose", Started: at}, {ID: "a2", Type: "verifier", Started: at}}},
		{Event{HookEventName: HookSubagentStop, SessionID: "s1", AgentID: "a1"}, []SubAgent{{ID: "a2", Type: "verifier", Started: at}}},
		{Event{HookEventName: HookSubagentStop, SessionID: "s1", AgentID: "a2"}, nil},
	}
	for _, s := range steps {
		if err := Record(dir, s.e, at); err != nil {
			t.Fatal(err)
		}
		got, err := Running(dir, "s1")
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got, s.want) {
			t.Errorf("after %s of %s: %+v, want %+v", s.e.HookEventName, s.e.AgentID, got, s.want)
		}
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("left in %s: %v", dir, entries)
	}
}

// S-0285: a session ID that is not a safe file name, an event with no agent,
// and an event of another kind are ignored.
func TestRecordIgnoresWhatItCannotFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), RecordDir)
	for _, e := range []Event{
		{HookEventName: HookSubagentStart, SessionID: "../s1", AgentID: "a1"},
		{HookEventName: HookSubagentStart, SessionID: "s 1", AgentID: "a1"},
		{HookEventName: HookSubagentStart, SessionID: "", AgentID: "a1"},
		{HookEventName: HookSubagentStart, SessionID: "s1"},
		{HookEventName: "PreToolUse", SessionID: "s1", AgentID: "a1"},
	} {
		if err := Record(dir, e, time.Now()); err != nil {
			t.Errorf("%+v: %v", e, err)
		}
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("recorded: %v", err)
	}
	if got, err := Running(dir, "../s1"); got != nil || err != nil {
		t.Errorf("unsafe session read: %v, %v", got, err)
	}
}

// S-0285: a layer's sub-agents start at once, each in a hook process of its
// own; every start lands in the record.
func TestRecordKeepsEveryStartOfALayer(t *testing.T) {
	dir := filepath.Join(t.TempDir(), RecordDir)
	const n = 16
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := range n {
		wg.Go(func() { errs <- Record(dir, started(fmt.Sprintf("a%02d", i), "general-purpose"), time.Now()) })
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	got, err := Running(dir, "s1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != n {
		t.Errorf("%d of %d starts recorded: %+v", len(got), n, got)
	}
}

// S-0285: a lock left by a guard that died holding it is taken once it is
// old.
func TestRecordTakesALockLeftBehind(t *testing.T) {
	dir := filepath.Join(t.TempDir(), RecordDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	lock := filepath.Join(dir, "s1.json.lock")
	if err := os.WriteFile(lock, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Minute)
	if err := os.Chtimes(lock, old, old); err != nil {
		t.Fatal(err)
	}
	if err := Record(dir, started("a1", "verifier"), time.Now()); err != nil {
		t.Fatal(err)
	}
	if got, _ := Running(dir, "s1"); len(got) != 1 {
		t.Errorf("running %+v", got)
	}
}

// storyGuard is a story's agent's guard for S-0001, with running sub-agents
// and an answer for whether a thread on the story is open.
func storyGuard(running []SubAgent, runErr error, open bool, openErr error) Guard {
	gd := g
	gd.Story = "S-0001"
	gd.Running = func() ([]SubAgent, error) { return running, runErr }
	gd.ThreadOpen = func() (bool, error) { return open, openErr }
	return gd
}

// waitCall is a wait_for_events call, the session's own when agent is "".
func waitCall(agent string) Event {
	e := Event{HookEventName: "PreToolUse", SessionID: "s1", ToolName: MCPPrefix + "wait_for_events"}
	if agent != "" {
		e.AgentID, e.AgentType = "a9", agent
	}
	return e
}

// S-0285: a story's agent's own wait_for_events is refused while a
// sub-agent of its session runs and no thread on its story is open, naming
// each sub-agent and how to wait instead.
func TestTheStorysAgentDoesNotWaitOnItsSubAgentsWithWaitForEvents(t *testing.T) {
	running := []SubAgent{{ID: "ac73f87145abc4315", Type: "general-purpose"}, {ID: "b2"}}
	why := storyGuard(running, nil, false, nil).Check(waitCall(""))
	want := "the story's agent cannot call wait_for_events while its sub-agents run (general-purpose ac73f87145abc4315, unnamed b2) and no thread on S-0001 awaits the designer: " + waitsOnSubAgents
	if why != want {
		t.Errorf("got  %q\nwant %q", why, want)
	}
	for _, s := range []string{"run_in_background set to false", "in one message", "only for a thread awaiting the designer"} {
		if !strings.Contains(why, s) {
			t.Errorf("the refusal does not say %q", s)
		}
	}
}

// S-0285: the wait passes with no sub-agent running, with a thread on the
// story open, when either cannot be read, for another tool, and outside a
// story's agent's session; a sub-agent's wait is refused as before.
func TestTheStorysAgentWaitsForTheDesigner(t *testing.T) {
	one := []SubAgent{{ID: "a1", Type: "verifier"}}
	fail := errors.New("unreadable")
	planner, orchestrating, outside := storyGuard(one, nil, false, nil), storyGuard(one, nil, false, nil), storyGuard(one, nil, false, nil)
	planner.Role, orchestrating.Role, outside.Story = RolePlan, RoleOrchestrate, ""
	unknown := g
	unknown.Story = "S-0001"
	for what, c := range map[string]struct {
		gd Guard
		e  Event
	}{
		"no sub-agent running":          {storyGuard(nil, nil, false, nil), waitCall("")},
		"a thread on the story open":    {storyGuard(one, nil, true, nil), waitCall("")},
		"the record unreadable":         {storyGuard(one, fail, false, nil), waitCall("")},
		"the threads unreadable":        {storyGuard(one, nil, false, fail), waitCall("")},
		"another tool":                  {storyGuard(one, nil, false, nil), Event{ToolName: MCPPrefix + "inbox"}},
		"the planner":                   {planner, waitCall("")},
		"the orchestrator":              {orchestrating, waitCall("")},
		"outside a story's agent":       {outside, waitCall("")},
		"no record of sub-agents given": {unknown, waitCall("")},
	} {
		if why := c.gd.Check(c.e); why != "" {
			t.Errorf("%s: refused: %s", what, why)
		}
	}
	why := storyGuard(nil, nil, true, nil).Check(waitCall("general-purpose"))
	if !strings.Contains(why, "a sub-agent (general-purpose) cannot call wait_for_events") {
		t.Errorf("a sub-agent's wait: %q", why)
	}
}
