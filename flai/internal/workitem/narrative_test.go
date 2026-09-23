package workitem

import (
	"strings"
	"testing"
)

func TestOpenQuestions(t *testing.T) {
	body := "## Open questions\n- First?\n  continued here\n* Second?\n<!-- threads:start -->\n- TH-0001 generated\n<!-- threads:end -->\n- Third?\n\n## Log\n"
	got := OpenQuestions(body)
	if strings.Join(got, "|") != "First? continued here|Second?|Third?" {
		t.Errorf("questions: %q", got)
	}
	if OpenQuestions("## Log\n") != nil {
		t.Error("nothing to find must find nothing")
	}
}

// A story's narrative that still has a hand-written open question refuses
// to go to review (S-0089): the operator's answer, not the agent's guess,
// is what the story needs before it can be reviewed.
func TestMoveToReviewRefusesAnOpenQuestion(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] works\n", 1)
	_ = r.Save(s)
	mustMove(t, r, s, Ready, "")
	mustMove(t, r, s, InProgress, "")
	mustCreate(t, r, Task, "T", "S-0001")
	s, _ = r.Get("S-0001")

	if _, err := r.OpenStream(s, StreamOptions{Agent: "test", Now: t0}); err != nil {
		t.Fatal(err)
	}
	n, err := ReadNarrative(r.NarrativePath("S-0001"))
	if err != nil {
		t.Fatal(err)
	}
	n.Body = strings.Replace(n.Body, "## Open questions\n", "## Open questions\n- Which port should it use?\n", 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}

	items, _ := r.List(false)
	if _, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items}); err == nil ||
		!strings.Contains(err.Error(), "Which port should it use?") {
		t.Errorf("review with an open question must be refused, saying which: %v", err)
	}

	// removing it (by hand, as the convention has it) lets the move through
	n, _ = ReadNarrative(r.NarrativePath("S-0001"))
	n.Body = strings.Replace(n.Body, "- Which port should it use?\n", "", 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	items, _ = r.List(false)
	if _, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items}); err != nil {
		t.Errorf("review once the question is gone: %v", err)
	}
}

// A story with no stream open yet has no narrative to check, so nothing
// here refuses it: the "at least one task" rule still applies on its own.
func TestMoveToReviewWithNoNarrativeIsUnaffected(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] works\n", 1)
	_ = r.Save(s)
	mustMove(t, r, s, Ready, "")
	mustMove(t, r, s, InProgress, "")
	mustCreate(t, r, Task, "T", "S-0001")
	s, _ = r.Get("S-0001")
	items, _ := r.List(false)
	if _, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items}); err != nil {
		t.Errorf("review with no narrative: %v", err)
	}
}
