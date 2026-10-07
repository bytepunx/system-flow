package mcpserver

import (
	"context"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// readInbox reads agent's inbox through inbox.Read, as a command does from
// the CLI, at the fixture's clock.
func (f *fixture) readInbox(t *testing.T, agent string) inbox.Inbox {
	t.Helper()
	out, err := inbox.Read(context.Background(), inbox.Options{Repo: f.repo, Agent: agent, Now: func() time.Time { return *f.clock }, Version: "test"})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func changeIDs(changes []Event) map[string]bool {
	ids := map[string]bool{}
	for _, c := range changes {
		ids[c.ID] = true
	}
	return ids
}

// The CLI and the MCP tool read one cursor per agent name: a change the CLI
// reported under the name is not reported again by the tool under it.
func TestAChangeTheCLIReportedIsNotReportedAgainByTheTool(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}

	*f.clock = t0.Add(10 * time.Minute)
	s := f.readyStory(t, "Made ready meanwhile", t0.Add(5*time.Minute))
	cli := f.readInbox(t, "claude")
	if !changeIDs(cli.Changes)[s.ID] {
		t.Fatalf("the CLI reports the move: %+v", cli.Changes)
	}

	*f.clock = t0.Add(11 * time.Minute)
	tool, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if got := changeSummaries(tool, "changes"); len(got) != 0 {
		t.Errorf("the tool reports again what the CLI reported: %v", got)
	}
	if ready := tool["ready"].([]any); len(ready) != 1 || ready[0].(map[string]any)["id"] != s.ID {
		t.Errorf("ready work is state and stays listed: %v", tool["ready"])
	}
}

// A first look under a new name covers the last 24 hours of stories and
// epics: not what is older, and not task transitions.
func TestAFirstLookCoversADayOfStoriesAndEpics(t *testing.T) {
	f := setup(t)
	now := t0.Add(48 * time.Hour)
	*f.clock = now
	old := f.readyStory(t, "Ready two days ago", now.Add(-25*time.Hour))
	recent := f.readyStory(t, "Ready yesterday", now.Add(-23*time.Hour))
	epic, err := f.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic of today", Owner: "alex", Now: now.Add(-2 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.Transition(epic, workitem.Ready, "alex", "", now.Add(-2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	task, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Task of today", Parent: f.story.ID, Owner: "alex", Now: now.Add(-time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.Transition(task, workitem.Ready, "alex", "", now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}

	ids := changeIDs(f.readInbox(t, "newcomer").Changes)
	if !ids[recent.ID] || !ids[epic.ID] {
		t.Errorf("the last day's story and epic are reported: %v", ids)
	}
	if ids[old.ID] || ids[task.ID] || ids[f.task.ID] {
		t.Errorf("an older change or a task transition is reported on a first look: %v", ids)
	}
}
