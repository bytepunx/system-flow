package inbox

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0532: an agent event appended to the log is reported once through
// CatchUp, after the cursor, sorted by time with the other kinds, with its
// fields and in words.
func TestAnAgentEventIsReportedOnceSortedWithTheOthers(t *testing.T) {
	f := setup(t)
	f.catchUp(t, "claude") // the first look sets the cursor
	f.clock = t0.Add(10 * time.Minute)
	f.appendAt(t, AgentEvent(f.story.ID, workitem.Story, "Story", AgentFailed, "agent-"+f.story.ID, conventions.RoleStory, "exit 1"), t0.Add(6*time.Minute))
	f.appendAt(t, PlannedEvent(f.story.ID, workitem.Story, "Story", "planner-"+f.story.ID, []string{"T-0001"}), t0.Add(5*time.Minute))
	f.appendAt(t, AgentEvent("orchestrator", RoleType, "orchestrator", AgentStarted, "orchestrator", conventions.RoleOrchestrate, ""), t0.Add(7*time.Minute))
	if _, err := f.repo.Transition(f.storyIn(t, "Next", []string{"docs"}), workitem.Ready, "alex", "", t0.Add(8*time.Minute)); err != nil {
		t.Fatal(err)
	}

	events, omitted := f.catchUp(t, "claude")
	if omitted != 0 {
		t.Fatalf("nothing omitted: %d", omitted)
	}
	var kinds []string
	for _, e := range events {
		kinds = append(kinds, e.Kind+" "+e.To)
	}
	if got, want := strings.Join(kinds, "; "), "planned T-0001; agent failed; agent started; moved ready"; got != want {
		t.Fatalf("sorted by time with the other kinds:\n got %s\nwant %s", got, want)
	}
	failed := events[1]
	if failed.ID != f.story.ID || failed.By != "agent-"+f.story.ID || failed.Fields[FieldWhy] != "exit 1" || failed.Fields[FieldRole] != conventions.RoleStory {
		t.Errorf("the failed event and its fields: %+v", failed)
	}
	if want := f.story.ID + " Story's agent agent-" + f.story.ID + " failed: exit 1" + letPass; failed.Summary != want {
		t.Errorf("summary:\n got %s\nwant %s", failed.Summary, want)
	}
	if again, omitted := f.catchUp(t, "claude"); len(again) != 0 || omitted != 0 {
		t.Errorf("told once: %+v, %d omitted", again, omitted)
	}
}

// S-0532: agent events count under MaxEvents with the rest, the newest kept,
// and one that falls out is counted in omitted and not told again.
func TestAgentEventsCountUnderTheCap(t *testing.T) {
	f := setup(t)
	f.catchUp(t, "claude")
	f.clock = t0.Add(time.Hour)
	for i := range MaxEvents - 2 {
		s := f.storyIn(t, fmt.Sprintf("Story %d", i), []string{"docs"})
		if _, err := f.repo.Transition(s, workitem.Ready, "alex", "", t0.Add(5*time.Minute+time.Duration(i)*time.Second)); err != nil {
			t.Fatal(err)
		}
	}
	// two after the cursor and older than every move, two newer
	for i, at := range []time.Duration{2 * time.Minute, 3 * time.Minute, 50 * time.Minute, 51 * time.Minute} {
		f.appendAt(t, AgentEvent(f.story.ID, workitem.Story, "Story", AgentEnded, fmt.Sprintf("agent-%d", i), conventions.RoleStory, ""), t0.Add(at))
	}
	events, omitted := f.catchUp(t, "claude")
	if len(events) != MaxEvents || omitted != 2 {
		t.Fatalf("%d moves and 4 agent events, capped at %d: %d listed, %d omitted", MaxEvents-2, MaxEvents, len(events), omitted)
	}
	agents := 0
	for _, e := range events {
		if e.Kind == Agent {
			agents++
		}
	}
	if agents != 2 || events[len(events)-1].By != "agent-3" || events[len(events)-2].By != "agent-2" {
		t.Errorf("the two newer agent events kept, last: %d kept, %+v, %+v", agents, events[len(events)-2], events[len(events)-1])
	}
	if again, omitted := f.catchUp(t, "claude"); len(again) != 0 || omitted != 0 {
		t.Errorf("what fell out is not told later: %+v, %d omitted", again, omitted)
	}
}

// S-0532: an agent is not told of its own runs' events, as it is not of its
// own moves; other agents are, and other kinds by it still reach it.
func TestAnAgentIsNotToldOfItsOwnRuns(t *testing.T) {
	f := setup(t)
	own := "agent-" + f.story.ID
	f.catchUp(t, own)
	f.catchUp(t, "orchestrator")
	f.clock = t0.Add(10 * time.Minute)
	f.appendAt(t, AgentEvent(f.story.ID, workitem.Story, "Story", AgentStarted, own, conventions.RoleStory, ""), t0.Add(2*time.Minute))
	f.appendAt(t, AgentEvent("S-0099", workitem.Story, "Other", AgentEnded, "agent-S-0099", conventions.RoleStory, ""), t0.Add(3*time.Minute))
	f.appendAt(t, ReleaseDueEvent(f.story.ID, "Story", own, "", "met"), t0.Add(4*time.Minute))

	events, _ := f.catchUp(t, own)
	if len(events) != 2 || events[0].By != "agent-S-0099" || events[1].Kind != ReleaseDue {
		t.Errorf("the other agent's run and its own release-due, not its own run: %+v", events)
	}
	if events, _ := f.catchUp(t, "orchestrator"); len(events) != 3 {
		t.Errorf("another reader is told of all three: %+v", events)
	}
}

// S-0532: the summary says each to, with the item or the role, and the why.
func TestDescribeSaysEachAgentEvent(t *testing.T) {
	for _, tc := range []struct {
		e    LoggedEvent
		want string
	}{
		{AgentEvent("S-0532", workitem.Story, "Events", AgentStarted, "agent-S-0532", conventions.RoleStory, "pulled from ready"),
			"S-0532 Events's agent agent-S-0532 started: pulled from ready"},
		{AgentEvent("S-0532", workitem.Story, "Events", AgentEnded, "agent-S-0532", conventions.RoleStory, "worked: left in review"),
			"S-0532 Events's agent agent-S-0532 ended: worked: left in review"},
		{AgentEvent("S-0532", workitem.Story, "Events", AgentFailed, "agent-S-0532", conventions.RoleStory, "exit 1 with S-0532 in in-progress"),
			"S-0532 Events's agent agent-S-0532 failed: exit 1 with S-0532 in in-progress"},
		{AgentEvent("S-0532", workitem.Story, "Events", AgentStopped, "agent-S-0532", conventions.RoleStory, "stopped by the operator"),
			"S-0532 Events's agent agent-S-0532 was stopped: stopped by the operator"},
		{AgentEvent("E-0030", workitem.Epic, "Agents", AgentEnded, "planner-E-0030", conventions.RolePlan, "exit 0"),
			"E-0030 Agents's planner planner-E-0030 ended: exit 0"},
		{AgentEvent("orchestrator", RoleType, "orchestrator", AgentFailed, "orchestrator", conventions.RoleOrchestrate, "the usage limit"),
			"the orchestrator's run failed: the usage limit"},
		{AgentEvent("curator", RoleType, "curator", AgentStarted, "curator", conventions.RoleCurate, ""),
			"the curator's run started"},
		{AgentEvent("reviewer", RoleType, "reviewer", AgentStopped, "reviewer-1", "reviewer", ""),
			"the reviewer's run reviewer-1 was stopped"},
	} {
		if got, want := describeLogged(tc.e.Change(), tc.e.Fields), tc.want+letPass; got != want {
			t.Errorf("describeLogged(%s %s):\n got %s\nwant %s", tc.e.ID, tc.e.To, got, want)
		}
	}
	// without its fields, from the change alone
	c := AgentEvent("S-0532", workitem.Story, "Events", AgentFailed, "agent-S-0532", conventions.RoleStory, "exit 1").Change()
	if got, want := Describe(c), "S-0532 Events's agent agent-S-0532 failed"+letPass; got != want {
		t.Errorf("Describe:\n got %s\nwant %s", got, want)
	}
}

// S-0532: AgentEventsOf lists one item's agent events from since, in any zero
// padding, oldest first, and nothing else of the log.
func TestAgentEventsOfListsOneItemsSince(t *testing.T) {
	f := setup(t)
	id := f.story.ID
	loose := strings.Replace(id, "-0", "-", 1)
	f.appendAt(t, AgentEvent(id, workitem.Story, "Story", AgentStarted, "a", conventions.RoleStory, ""), t0)
	f.appendAt(t, AgentEvent(id, workitem.Story, "Story", AgentFailed, "c", conventions.RoleStory, ""), t0.Add(3*time.Minute))
	f.appendAt(t, AgentEvent(loose, workitem.Story, "Story", AgentStarted, "b", conventions.RoleStory, ""), t0.Add(2*time.Minute))
	f.appendAt(t, AgentEvent("S-0099", workitem.Story, "Other", AgentStarted, "x", conventions.RoleStory, ""), t0.Add(2*time.Minute))
	f.appendAt(t, PlannedEvent(id, workitem.Story, "Story", "planner", []string{"T-0001"}), t0.Add(2*time.Minute))

	var by []string
	for _, e := range AgentEventsOf(f.repo, loose, t0.Add(2*time.Minute)) {
		by = append(by, e.By)
	}
	if got := strings.Join(by, ","); got != "b,c" {
		t.Errorf("the story's agent events from since, oldest first: %s", got)
	}
	if got := AgentEventsOf(f.repo, "orchestrator", t0); len(got) != 0 {
		t.Errorf("none on a role with none: %+v", got)
	}
}
