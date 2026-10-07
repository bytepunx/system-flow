package mcpserver

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// streamFixture is setup with the story's narrative holding a known body in
// every section, as session s9.
func streamFixture(t *testing.T) (*fixture, string) {
	t.Helper()
	t.Setenv("FLAI_SESSION", "s9")
	f := setup(t)
	n, err := workitem.ReadNarrative(f.repo.NarrativePath(f.story.ID))
	if err != nil {
		t.Fatal(err)
	}
	n.Body = fmt.Sprintf("\n# %s Story\n\n## Context\n\nWhy.\n\n## Current state\n\nOld state.\n\n## Next steps\n\n1. Old step.\n\n## Decisions\n\n- One.\n\n## Open questions\n\n- Which?\n\n## Log\n\n### 2026-09-18T17:00:00Z\n\nStream opened.\n", f.story.ID)
	n.Agent, n.Session = "codex", "s1"
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	return f, n.Body
}

// section is the text under a narrative's ## name, or fails the test.
func section(t *testing.T, body, name string) string {
	t.Helper()
	s, ok := workitem.NarrativeSection(body, name)
	if !ok {
		t.Fatalf("no ## %s in:\n%s", name, body)
	}
	return s.Text
}

// S-0271: stream_state replaces Current state and Next steps as the
// connected agent and leaves every other section as it was, appending
// nothing to the log; the answer names the narrative by its path from the
// project's root.
func TestStreamStateWritesCurrentStateAndNextSteps(t *testing.T) {
	f, before := streamFixture(t)
	f.lintAsThisProjectDoes(t)
	out, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID, "current": "T-1 is done.", "next": "1. Do T-2."})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["story"] != f.story.ID || out["path"] != "wip/agents/"+f.story.ID+".md" || out["changed"] != true || out["updated"] != "2026-09-18T17:01:00Z" {
		t.Errorf("the answer: %v", out)
	}
	if fmt.Sprint(out["written"]) != "[Current state Next steps]" {
		t.Errorf("written: %v", out["written"])
	}
	n, err := workitem.ReadNarrative(f.repo.NarrativePath(f.story.ID))
	if err != nil {
		t.Fatal(err)
	}
	if got := section(t, n.Body, workitem.CurrentState); got != "T-1 is done." {
		t.Errorf("current state: %q", got)
	}
	if got := section(t, n.Body, workitem.NextSteps); got != "1. Do T-2." {
		t.Errorf("next steps: %q", got)
	}
	for _, name := range []string{"Context", "Decisions", "Open questions", "Log"} {
		if got, want := section(t, n.Body, name), section(t, before, name); got != want {
			t.Errorf("## %s changed: %q, was %q", name, got, want)
		}
	}
	if n.Agent != "claude" || n.Session != "s9" || n.Updated != "2026-09-18T17:01:00Z" {
		t.Errorf("the narrative records the agent, session, and time: %s %s %s", n.Agent, n.Session, n.Updated)
	}
	if index, err := os.ReadFile(f.repo.NarrativePath("index")); err != nil || !strings.Contains(string(index), f.story.ID) {
		t.Errorf("the narratives' index is written again: %v\n%s", err, index)
	}

	// the same text again writes nothing
	out, failed = f.call(t, "stream_state", map[string]any{"story": f.story.ID, "current": "T-1 is done.", "next": "1. Do T-2."})
	if failed != "" || out["changed"] != false || out["updated"] != "2026-09-18T17:01:00Z" {
		t.Errorf("the same text again: %v %s", out, failed)
	}
}

// Each text alone writes its own section and leaves the other as it is.
func TestStreamStateWritesEachTextAlone(t *testing.T) {
	f, _ := streamFixture(t)
	path := f.repo.NarrativePath(f.story.ID)
	out, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID, "current": "Reviewing."})
	if failed != "" || fmt.Sprint(out["written"]) != "[Current state]" {
		t.Fatalf("current alone: %v %s", out, failed)
	}
	n, _ := workitem.ReadNarrative(path)
	if section(t, n.Body, workitem.CurrentState) != "Reviewing." || section(t, n.Body, workitem.NextSteps) != "1. Old step." {
		t.Errorf("current alone:\n%s", n.Body)
	}
	out, failed = f.call(t, "stream_state", map[string]any{"story": f.story.ID, "current": "  ", "next": "1. Answer the review."})
	if failed != "" || fmt.Sprint(out["written"]) != "[Next steps]" {
		t.Fatalf("next alone: %v %s", out, failed)
	}
	n, _ = workitem.ReadNarrative(path)
	if section(t, n.Body, workitem.CurrentState) != "Reviewing." || section(t, n.Body, workitem.NextSteps) != "1. Answer the review." {
		t.Errorf("next alone:\n%s", n.Body)
	}
}

// A story not in progress or in review, one with no narrative, text the
// lint rejects, and no text at all are refused with what to do, and nothing
// is written.
func TestStreamStateRefusals(t *testing.T) {
	f, _ := streamFixture(t)
	backlog, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Later", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	bare, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Bare", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(bare.Path)
	if err := os.WriteFile(bare.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if bare, err = f.repo.Get(bare.ID); err != nil {
		t.Fatal(err)
	}
	for _, st := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := f.repo.Transition(bare, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	f.lintAsThisProjectDoes(t)
	path := f.repo.NarrativePath(f.story.ID)
	was, _ := os.ReadFile(path)

	if _, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID}); !strings.Contains(failed, "give current, next, or both") {
		t.Errorf("neither text: %q", failed)
	}
	if _, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID, "current": " ", "next": "\n"}); !strings.Contains(failed, "give current, next, or both") {
		t.Errorf("blank texts: %q", failed)
	}
	_, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID, "next": "1. Do it.\n\n**Looks like a heading**"})
	if !strings.Contains(failed, "wip/agents/"+f.story.ID+".md would fail the project's markdown lint") || !strings.Contains(failed, "MD036") || !strings.Contains(failed, "call stream_state again") {
		t.Errorf("text the lint rejects: %q", failed)
	}
	if _, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID, "current": "Fine.\n\n## Decisions\n\n- Sneaked in."}); !strings.Contains(failed, "use ### or deeper") {
		t.Errorf("a second-level heading in the text: %q", failed)
	}
	if got, _ := os.ReadFile(path); string(got) != string(was) {
		t.Errorf("a refused write changed the narrative:\n%s", got)
	}

	if _, failed := f.call(t, "stream_state", map[string]any{"story": backlog.ID, "current": "Started."}); !strings.Contains(failed, backlog.ID+" is backlog") || !strings.Contains(failed, "move it to in-progress first") {
		t.Errorf("a story in backlog: %q", failed)
	}

	if _, failed := f.call(t, "stream_state", map[string]any{"story": bare.ID, "current": "Started."}); !strings.Contains(failed, "has no narrative") || !strings.Contains(failed, "flai stream open "+bare.ID) {
		t.Errorf("a story with no narrative: %q", failed)
	}
	if _, err := os.Stat(f.repo.NarrativePath(bare.ID)); !os.IsNotExist(err) {
		t.Errorf("no narrative is made: %v", err)
	}
}
