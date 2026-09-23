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

func openStreamWithQuestion(t *testing.T, r *Repo, question string) *Item {
	t.Helper()
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
	n.Body = strings.Replace(n.Body, "## Open questions\n", "## Open questions\n- "+question+"\n", 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	return s
}

// Answering moves the question out of Open questions and into Decisions
// with the answer beside it (design/system/agent-narrative.md: "Answered
// questions move to Decisions", S-0090).
func TestAnswerOpenQuestion(t *testing.T) {
	r := newProject(t)
	openStreamWithQuestion(t, r, "Which port should it use?")

	n, err := r.AnswerOpenQuestion("S-0001", "Which port should it use?", "Nine.", "alex", t0)
	if err != nil {
		t.Fatal(err)
	}
	if OpenQuestions(n.Body) != nil {
		t.Errorf("question still open: %q", OpenQuestions(n.Body))
	}
	if !strings.Contains(n.Body, "Which port should it use?") || !strings.Contains(n.Body, "Nine.") ||
		!strings.Contains(n.Body, "answered by alex") {
		t.Errorf("decision not recorded:\n%s", n.Body)
	}
	decisions := n.Body[strings.Index(n.Body, "## Decisions"):strings.Index(n.Body, "## Open questions")]
	if !strings.Contains(decisions, "Which port should it use?") {
		t.Errorf("answer landed outside Decisions:\n%s", n.Body)
	}

	// once it is gone, answering it again is refused, not silently a no-op
	if _, err := r.AnswerOpenQuestion("S-0001", "Which port should it use?", "again", "alex", t0); err == nil {
		t.Error("answering a question twice must be refused")
	}
}

// The two halves meet (S-0089, S-0090): a question blocks review until it is
// answered, and answering it, not deleting it by hand, is what lets the story
// through.
func TestAnsweringAnOpenQuestionUnblocksReview(t *testing.T) {
	r := newProject(t)
	s := openStreamWithQuestion(t, r, "Which port should it use?")
	items, _ := r.List(false)
	if _, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items}); err == nil {
		t.Fatal("review with an unanswered question must be refused")
	}
	if _, err := r.AnswerOpenQuestion("S-0001", "Which port should it use?", "Nine.", "alex", t0); err != nil {
		t.Fatal(err)
	}
	s, _ = r.Get("S-0001")
	items, _ = r.List(false)
	if _, err := r.Move(s, Review, MoveOptions{Now: t0, Items: items}); err != nil {
		t.Errorf("review once the question is answered: %v", err)
	}
}

func TestAnswerOpenQuestionRefusesWhatDoesNotMatch(t *testing.T) {
	r := newProject(t)
	openStreamWithQuestion(t, r, "Which port should it use?")
	if _, err := r.AnswerOpenQuestion("S-0001", "A different question", "x", "alex", t0); err == nil ||
		!strings.Contains(err.Error(), "no open question matching") {
		t.Errorf("unmatched question: %v", err)
	}
	if _, err := r.AnswerOpenQuestion("S-9999", "x", "y", "alex", t0); err == nil {
		t.Error("no such story must be refused")
	}
}

// A multi-line question (a continuation, as OpenQuestions joins it) is
// removed as one bullet, continuation lines included.
func TestAnswerOpenQuestionRemovesAContinuedBullet(t *testing.T) {
	r := newProject(t)
	openStreamWithQuestion(t, r, "First line")
	n, _ := ReadNarrative(r.NarrativePath("S-0001"))
	n.Body = strings.Replace(n.Body, "- First line\n", "- First line\n  more of it\n- Second?\n", 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	n, err := r.AnswerOpenQuestion("S-0001", "First line more of it", "done", "alex", t0)
	if err != nil {
		t.Fatal(err)
	}
	if got := OpenQuestions(n.Body); len(got) != 1 || got[0] != "Second?" {
		t.Errorf("the other question should remain: %q", got)
	}
}
