package mcpserver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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
	t.Parallel()
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
	t.Parallel()
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

// otherTask makes a story with a task in repo and returns the task's file.
func otherTask(t *testing.T, repo *workitem.Repo) string {
	t.Helper()
	other, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Other", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Half-written", Parent: other.ID, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	return task.Path
}

// halfWrite overwrites the item file at path with PLACEHOLDER, which does not
// parse, as a file holds for the moment it is half-written (I-0344).
func halfWrite(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("PLACEHOLDER"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// I-0344: a task file of another story that does not parse once failed every
// list of the project, and inbox with it. The inbox answers past it with the
// other items' changes and the ready stories.
func TestInboxAnswersPastAnUnreadableTaskFile(t *testing.T) {
	t.Parallel()
	f := setup(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	*f.clock = t0.Add(10 * time.Minute)
	s := f.readyStory(t, "Made ready meanwhile", t0.Add(5*time.Minute))
	halfWrite(t, otherTask(t, f.repo))

	*f.clock = t0.Add(11 * time.Minute)
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatalf("inbox with an unreadable task file: %s", failed)
	}
	reported := false
	for _, c := range out["changes"].([]any) {
		reported = reported || c.(map[string]any)["id"] == s.ID
	}
	if !reported {
		t.Errorf("the other story's move is not reported: %v", out["changes"])
	}
	if ready := out["ready"].([]any); len(ready) != 1 || ready[0].(map[string]any)["id"] != s.ID {
		t.Errorf("ready: %v", out["ready"])
	}
}

// I-0344: wait_for_events failed three times on a task file that held only
// PLACEHOLDER. A held wait that wakes on such a file answers with the path
// that changed, and the next wait answers with what others did meanwhile.
func TestWaitForEventsAnswersPastAnUnreadableTaskFile(t *testing.T) {
	t.Parallel()
	f := setupWith(t, holdLong)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	// the task is written before the wait, so that the wait wakes on its
	// file's half-write alone
	path := otherTask(t, f.repo)
	done := f.callHeld(t, "wait_for_events", map[string]any{"timeout_seconds": heldSeconds})
	halfWrite(t, path)
	rel, err := filepath.Rel(f.repo.Root, path)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case a := <-done:
		if a.failed != "" {
			t.Fatalf("wait_for_events with an unreadable task file: %s", a.failed)
		}
		changed, _ := a.out["changed"].([]any)
		named := false
		for _, p := range changed {
			named = named || p == filepath.ToSlash(rel)
		}
		if a.out["timed_out"] == true || !named {
			t.Errorf("the wait does not name the file that changed (%s): %v", rel, a.out)
		}
	case <-time.After(heldWait + 10*time.Second):
		t.Fatal("wait_for_events never answered")
	}

	s := f.readyStory(t, "Made ready meanwhile", t0.Add(2*time.Minute))
	out, failed := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 1})
	if failed != "" {
		t.Fatalf("the next wait: %s", failed)
	}
	events, _ := out["events"].([]any)
	reported := false
	for _, e := range events {
		reported = reported || e.(map[string]any)["id"] == s.ID
	}
	if !reported {
		t.Errorf("the next wait does not report the other story's move: %v", out)
	}
}

// budgetState is a fixed host budget state for the server's Budgets, which
// a test sets and clears through the pointer.
func budgetState(state *[]inbox.Budget) BudgetRead {
	return func() []inbox.Budget { return *state }
}

// anthropicOut is a budget state set on anthropic, held until 22:00 UTC.
var anthropicOut = []inbox.Budget{{Provider: "anthropic", Limit: "usage_limit_reached", Status: 429, Since: "2026-09-18T17:00:00Z", Until: "2026-09-18T22:00:00Z"}}

// budgetWrong says what is wrong with the budget an inbox answered, against
// the anthropicOut state; "" when it is that state.
func budgetWrong(got any) string {
	list, ok := got.([]any)
	if !ok || len(list) != 1 {
		return "one entry"
	}
	b := list[0].(map[string]any)
	if b["provider"] != "anthropic" || b["limit"] != "usage_limit_reached" || b["until"] != "2026-09-18T22:00:00Z" || b["since"] != "2026-09-18T17:00:00Z" || b["status"] != float64(429) {
		return "the provider, the limit, the status, since, and until"
	}
	if _, ok := b["probing"]; ok {
		return "no probing while no probe runs"
	}
	return ""
}

// S-0534: over stdio, inbox carries the host's budget state while it is set,
// with the provider, the limit, and until, and leaves budget out once it
// clears; the instructions tell an agent to hold rather than retry.
func TestInboxCarriesTheHostsBudgetWhileItIsSet(t *testing.T) {
	t.Parallel()
	state := anthropicOut
	f := setupWith(t, func(o *Options) { o.Budgets = budgetState(&state) })
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if wrong := budgetWrong(out["budget"]); wrong != "" {
		t.Errorf("budget, %s: %v", wrong, out["budget"])
	}
	state = nil
	if out, failed = f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	if _, ok := out["budget"]; ok {
		t.Errorf("a cleared state leaves budget out: %v", out["budget"])
	}
	if in := f.cs.InitializeResult().Instructions; !strings.Contains(in, budgetInstructions) || !strings.Contains(budgetInstructions, "When inbox carries budget") || !strings.Contains(budgetInstructions, "do not retry") {
		t.Errorf("the instructions name budget and say to hold rather than retry: %s", in)
	}
}

// A server with no budget reader, as a test or a host-less one has it,
// answers no budget.
func TestInboxWithoutABudgetReaderCarriesNone(t *testing.T) {
	t.Parallel()
	f := setup(t)
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if _, ok := out["budget"]; ok {
		t.Errorf("no reader, no budget: %v", out["budget"])
	}
}
