package mcpserver

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// sectionsThreads is a threads block holding one thread, as flai mirrors
// threads into a narrative's ## Open questions.
const sectionsThreads = "<!-- threads:start -->\n- [TH-0007](../threads/TH-0007-x.md) Which way? (awaiting operator)\n<!-- threads:end -->\n"

// I-0277: a story's agent in a container has wip/ read-only (ADR-0253), so
// it can write its narrative's Context, Decisions, and Open questions only
// through stream_state, which the host's server answers over HTTP
// (flai/cmd/mcp.go) and which writes the main checkout. This test calls the
// same tool in-process, as the other stream_state tests do: the transport
// is not what I-0277 is about. The agent's session writes all three and
// reads them back from the narrative, its threads block unchanged.
func TestStreamStateWritesContextDecisionsAndQuestions(t *testing.T) {
	t.Parallel()
	f, _ := streamFixture(t)
	f.lintAsThisProjectDoes(t)
	path := f.repo.NarrativePath(f.story.ID)
	n, err := workitem.ReadNarrative(path)
	if err != nil {
		t.Fatal(err)
	}
	n.Body = strings.Replace(n.Body, "- Which?\n", "- Which?\n\n"+sectionsThreads, 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	before := n.Body

	out, failed := f.call(t, "stream_state", map[string]any{
		"story":     f.story.ID,
		"context":   "Why, from the container.",
		"decisions": []string{"Widen stream_state, not the mount.", "Append, never replace."},
		"questions": []string{"Should the CLI take them too?"},
	})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["changed"] != true || fmt.Sprint(out["written"]) != "[Context Decisions Open questions]" {
		t.Errorf("the answer: %v", out)
	}
	n, err = workitem.ReadNarrative(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := section(t, n.Body, workitem.ContextSection); got != "Why, from the container." {
		t.Errorf("context: %q", got)
	}
	if got := section(t, n.Body, workitem.DecisionsSection); got != "- One.\n- Widen stream_state, not the mount.\n- Append, never replace." {
		t.Errorf("decisions: %q", got)
	}
	questions := section(t, n.Body, workitem.OpenQuestionsSection)
	if want := "- Which?\n- 2026-09-18: Should the CLI take them too?\n\n" + strings.TrimSuffix(sectionsThreads, "\n"); questions != want {
		t.Errorf("open questions:\n%q\nwant\n%q", questions, want)
	}
	if !strings.Contains(n.Body, sectionsThreads) {
		t.Errorf("the threads block changed:\n%s", n.Body)
	}
	if got := workitem.OpenQuestions(n.Body); fmt.Sprint(got) != "[Which? 2026-09-18: Should the CLI take them too?]" {
		t.Errorf("OpenQuestions: %q", got)
	}
	for _, name := range []string{workitem.CurrentState, workitem.NextSteps, "Log"} {
		if got, want := section(t, n.Body, name), section(t, before, name); got != want {
			t.Errorf("## %s changed: %q, was %q", name, got, want)
		}
	}
	if n.Agent != "claude" || n.Session != "s9" {
		t.Errorf("the narrative records the agent and session: %s %s", n.Agent, n.Session)
	}

	// a refused write, none of the five, names them all and writes nothing
	was, _ := os.ReadFile(path)
	if _, failed := f.call(t, "stream_state", map[string]any{"story": f.story.ID, "decisions": []string{" "}, "questions": []string{}}); !strings.Contains(failed, "context, decisions, or questions") || !strings.Contains(failed, "give current, next") {
		t.Errorf("no text: %q", failed)
	}
	if got, _ := os.ReadFile(path); string(got) != string(was) {
		t.Errorf("a refused write changed the narrative:\n%s", got)
	}
}
