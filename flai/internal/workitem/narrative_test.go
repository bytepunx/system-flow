package workitem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
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

// stateBody is a narrative with something in every section.
const stateBody = "\n# S-0001 S\n\n## Context\n\nWhy.\n\n## Current state\n\nOld state.\n\n## Next steps\n\n1. Old step.\n\n## Decisions\n\n- One.\n\n## Open questions\n\n- Which?\n\n## Log\n\n### 2026-09-15T20:00:00Z\nStream opened.\n"

// inProgressStream is S-0001, in progress, with stateBody for its narrative,
// in a project whose markdown lint is the template's.
func inProgressStream(t *testing.T) *Repo {
	t.Helper()
	r := lintProject(t)
	s, err := r.Create(NewOptions{Type: Story, Title: "S", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	s.Status = InProgress
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	n, err := r.OpenStream(s, StreamOptions{Agent: "a", Session: "s1", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	n.Body = stateBody
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	return r
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0271: both sections are replaced, every other byte of the body stays as
// it was, nothing is logged, and the index shows the new stamp and agent.
func TestSetStreamStateReplacesBothSections(t *testing.T) {
	r := inProgressStream(t)
	later := t0.Add(time.Hour)
	res, err := r.SetStreamState("S-0001", "New state.\n\nMore of it.\n", "1. Step one.\n2. Step two.", StreamOptions{Agent: "b", Now: later})
	if err != nil {
		t.Fatal(err)
	}
	stamp := later.Format(TimeFormat)
	if !res.Changed || strings.Join(res.Written, "|") != "Current state|Next steps" || res.Updated != stamp ||
		res.Stream != "S-0001" || res.Path != r.NarrativePath("S-0001") {
		t.Errorf("result: %+v", res)
	}
	n, err := ReadNarrative(r.NarrativePath("S-0001"))
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Replace(stateBody, "## Current state\n\nOld state.\n\n## Next steps\n\n1. Old step.\n",
		"## Current state\n\nNew state.\n\nMore of it.\n\n## Next steps\n\n1. Step one.\n2. Step two.\n", 1)
	if n.Body != want {
		t.Errorf("body:\n%s\nwant:\n%s", n.Body, want)
	}
	if n.Updated != stamp || n.Agent != "b" || n.Session != "s1" || n.Title != "S" {
		t.Errorf("front matter: %+v", n)
	}
	if index := readFile(t, filepath.Join(r.AgentsDir(), "index.md")); !strings.Contains(index, "| b | "+stamp+" |") {
		t.Errorf("index does not show the write:\n%s", index)
	}
}

// Text left out leaves its section as it was; giving neither is refused; the
// same text again writes nothing.
func TestSetStreamStateLeavesASectionLeftOut(t *testing.T) {
	r := inProgressStream(t)
	path := r.NarrativePath("S-0001")
	res, err := r.SetStreamState("S-0001", "", "- Only the steps.", StreamOptions{Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res.Written, "|") != "Next steps" {
		t.Errorf("written: %q", res.Written)
	}
	body := readFile(t, path)
	if s, _ := NarrativeSection(body, CurrentState); s.Text != "Old state." {
		t.Errorf("current state left out was changed: %q", s.Text)
	}
	if s, _ := NarrativeSection(body, NextSteps); s.Text != "- Only the steps." {
		t.Errorf("next steps: %q", s.Text)
	}

	if _, err := r.SetStreamState("S-0001", "Only the state.", "  \n", StreamOptions{Now: t0}); err != nil {
		t.Fatal(err)
	}
	body = readFile(t, path)
	if s, _ := NarrativeSection(body, NextSteps); s.Text != "- Only the steps." {
		t.Errorf("next steps left out was changed: %q", s.Text)
	}
	if s, _ := NarrativeSection(body, CurrentState); s.Text != "Only the state." {
		t.Errorf("current state: %q", s.Text)
	}

	if _, err := r.SetStreamState("S-0001", " ", "", StreamOptions{Now: t0}); err == nil || !strings.Contains(err.Error(), "give the current state, the next steps, or both") {
		t.Errorf("neither text: %v", err)
	}
	res, err = r.SetStreamState("S-0001", "Only the state.", "", StreamOptions{Now: t0.Add(time.Hour)})
	if err != nil || res.Changed || res.Updated != t0.Format(TimeFormat) {
		t.Errorf("the same text again: %+v, %v", res, err)
	}
	if readFile(t, path) != body {
		t.Error("nothing to change must write nothing")
	}
}

// A story not in progress or in review, or with no narrative, is refused
// with its state and what to do, and its narrative is left as it was.
func TestSetStreamStateRefusesAStoryItMayNotWrite(t *testing.T) {
	r := inProgressStream(t)
	path := r.NarrativePath("S-0001")
	was := readFile(t, path)
	var refused *StreamRefusedError
	for status, todo := range map[string]string{Ready: "`flai move S-0001 in-progress`", Done: "narrative is finished", Backlog: "`flai move S-0001 in-progress`"} {
		s, _ := r.Get("S-0001")
		s.Status = status
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
		_, err := r.SetStreamState("S-0001", "New.", "1. Next.", StreamOptions{Now: t0})
		if !errors.As(err, &refused) || refused.Status != status || !strings.Contains(err.Error(), "S-0001 is "+status) || !strings.Contains(err.Error(), todo) {
			t.Errorf("%s: %v", status, err)
		}
		if readFile(t, path) != was {
			t.Errorf("%s: the narrative was changed", status)
		}
	}
	s, _ := r.Get("S-0001")
	s.Status = Review
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	if _, err := r.SetStreamState("S-0001", "In review.", "", StreamOptions{Now: t0}); err != nil {
		t.Errorf("a story in review: %v", err)
	}

	other, err := r.Create(NewOptions{Type: Story, Title: "Other", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	other.Status = InProgress
	if err := r.Save(other); err != nil {
		t.Fatal(err)
	}
	_, err = r.SetStreamState("S-0002", "New.", "", StreamOptions{Now: t0})
	if !errors.As(err, &refused) || refused.Status != InProgress || !strings.Contains(err.Error(), "no narrative") || !strings.Contains(err.Error(), "`flai stream open S-0002`") {
		t.Errorf("no narrative: %v", err)
	}
	if _, err := os.Stat(r.NarrativePath("S-0002")); !os.IsNotExist(err) {
		t.Errorf("no narrative must be made: %v", err)
	}
	if _, err := r.SetStreamState("S-9999", "New.", "", StreamOptions{Now: t0}); !errors.Is(err, ErrNotFound) {
		t.Errorf("no such story: %v", err)
	}
}

// Text the lint rejects, or that holds a heading that would end the section,
// is refused and the narrative is left as it was.
func TestSetStreamStateRefusesTextTheLintRejects(t *testing.T) {
	r := inProgressStream(t)
	path := r.NarrativePath("S-0001")
	was := readFile(t, path)
	var le *mdlint.Error
	if _, err := r.SetStreamState("S-0001", "Fine.", "1. Do it.\n\n**Looks like a heading**", StreamOptions{Now: t0}); !errors.As(err, &le) || le.Findings[0].Rule != "MD036" {
		t.Errorf("text the lint rejects: %v", err)
	}
	if _, err := r.SetStreamState("S-0001", "Fine.\n\n## Decisions\n\n- Sneaked in.", "", StreamOptions{Now: t0}); err == nil || !strings.Contains(err.Error(), "use ### or deeper") {
		t.Errorf("a second-level heading in the text: %v", err)
	}
	if readFile(t, path) != was {
		t.Error("a refused write changed the narrative")
	}
}

// A section the narrative lacks is inserted where the template puts it.
func TestSetStreamStateInsertsAMissingSection(t *testing.T) {
	r := inProgressStream(t)
	n, _ := ReadNarrative(r.NarrativePath("S-0001"))
	n.Body = strings.Replace(n.Body, "## Next steps\n\n1. Old step.\n\n", "", 1)
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.SetStreamState("S-0001", "", "1. Back again.", StreamOptions{Now: t0}); err != nil {
		t.Fatal(err)
	}
	n, _ = ReadNarrative(r.NarrativePath("S-0001"))
	if want := strings.Replace(stateBody, "1. Old step.", "1. Back again.", 1); n.Body != want {
		t.Errorf("body:\n%s\nwant:\n%s", n.Body, want)
	}

	n.Body = "\n# S-0001 S\n\n## Context\n\nWhy.\n"
	if err := n.Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.SetStreamState("S-0001", "Now.", "1. Then.", StreamOptions{Now: t0}); err != nil {
		t.Fatal(err)
	}
	n, _ = ReadNarrative(r.NarrativePath("S-0001"))
	if want := "\n# S-0001 S\n\n## Context\n\nWhy.\n\n## Current state\n\nNow.\n\n## Next steps\n\n1. Then.\n"; n.Body != want {
		t.Errorf("appended:\n%q\nwant:\n%q", n.Body, want)
	}
}

func TestNarrativeSection(t *testing.T) {
	text := "---\nstream: S-0001\n---\n\n## Current state\n\n\nSaid.\n\n## Next steps\n1.\n\n## Log\n"
	s, ok := NarrativeSection(text, CurrentState)
	if !ok || s.Line != 5 || s.Text != "Said." || s.Name != CurrentState || !s.Written() {
		t.Errorf("current state: %+v, %v", s, ok)
	}
	s, ok = NarrativeSection(text, NextSteps)
	if !ok || s.Line != 10 || s.Text != "1." || s.Written() {
		t.Errorf("the template's lone 1. is not written: %+v, %v", s, ok)
	}
	if s, ok := NarrativeSection(text, "Log"); !ok || s.Text != "" || s.Written() {
		t.Errorf("an empty last section: %+v, %v", s, ok)
	}
	if _, ok := NarrativeSection(text, "Decisions"); ok {
		t.Error("a missing section must not be found")
	}
	for text, want := range map[string]bool{"-": false, "  * ": false, "12.": false, "- [ ] x": true, "x": true, "\n\n2. y": true} {
		if got := (Section{Text: text}).Written(); got != want {
			t.Errorf("Written(%q) = %v, want %v", text, got, want)
		}
	}
}
