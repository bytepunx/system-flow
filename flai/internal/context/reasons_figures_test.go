package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// storedReport is a last flai verify result as flai verify stores it, which
// the context package reads as the file it is.
const storedReport = `{
  "story": "S-0001",
  "commit": "0123456789abcdef",
  "base": "main",
  "ran_at": "2026-10-29T10:00:00Z",
  "duration_ms": 61000,
  "duration": "1m1s",
  "passed": false,
  "stopped_at": "tier go",
  "steps": [
    {"name": "check", "state": "passed", "duration_ms": 900, "duration": "0.9s"},
    {"name": "go", "tier": true, "state": "failed", "duration_ms": 60000, "duration": "1m0s",
     "findings": [{"name": "TestX", "path": "flai/x_test.go", "line": 12, "message": "want 1, got 2"}]}
  ]
}
`

// figureItem is the reason item of the pack at path, failing the test when
// there is none.
func figureItem(t *testing.T, p *Pack, path string) Item {
	t.Helper()
	it, ok := reasonItems(p)[path]
	if !ok {
		t.Fatalf("the pack loads no %s: %+v", path, p.Items)
	}
	return it
}

// ADR-0217 §3, S-0537 criterion 1: a review loads the story whole and its
// last flai verify result as flai verify stores it, with the story's ID; a
// review of what is not a story, or of a story that does not exist, is
// refused.
func TestAReviewReasonLoadsTheStoryAndItsLastVerifyResult(t *testing.T) {
	repo, epic, story := reasonsRepo(t)
	writeRepo(t, repo, ".flai-cache/verify/"+story.ID+".json", storedReport)
	p := mustReasons(t, repo, "review:S-1")
	raw, _ := os.ReadFile(story.Path)
	if len(p.Items) != 2 || p.Items[0].Path != fixturePath(t, repo, story.Path) || p.Items[0].Text != string(raw) || p.Items[0].Reason != "review:"+story.ID {
		t.Errorf("items %+v; want the story whole, then its result", p.Items)
	}
	got := figureItem(t, p, ".flai-cache/verify/"+story.ID+".json")
	if got.Text != storedReport || got.ID != story.ID || got.Title != "last flai verify result of "+story.ID || got.Reason != "review:"+story.ID || got.Size != len(storedReport) {
		t.Errorf("verify item %+v", got)
	}
	if !strings.Contains(p.Body(), "\n.flai-cache/verify/"+story.ID+".json\n") {
		t.Errorf("body lacks the result under its path:\n%s", p.Body())
	}
	for _, reason := range []string{"review:" + epic.ID, "review:S-0099"} {
		if _, err := ForReasons(repo, []string{reason}, ""); err == nil {
			t.Errorf("%s loaded; want it refused", reason)
		}
	}
}

// ADR-0217 §3: a story in review that flai verify never ran on loads an item
// that says it has no result and how to run it.
func TestAReviewReasonSaysAStoryHasNoVerifyResult(t *testing.T) {
	repo, _, story := reasonsRepo(t)
	p := mustReasons(t, repo, "review:"+story.ID)
	got := figureItem(t, p, ".flai-cache/verify/"+story.ID+".json")
	for _, want := range []string{story.ID + " has no stored flai verify result", "flai verify " + story.ID + ",", "flai verify " + story.ID + " --last"} {
		if !strings.Contains(got.Text, want) {
			t.Errorf("no-result item lacks %q: %q", want, got.Text)
		}
	}
	if got.ID != story.ID || got.Size != len(got.Text) || len(p.Items) != 2 {
		t.Errorf("no-result item %+v; items %d", got, len(p.Items))
	}
}

// ADR-0217 §3, S-0537 criterion 1: a release loads the release policy's
// figures as release_evaluate gives them, in its JSON.
func TestAReleaseReasonLoadsTheReleaseFigures(t *testing.T) {
	repo, _, _ := reasonsRepo(t)
	manifest, err := os.ReadFile(filepath.Join(repo.Root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	writeRepo(t, repo, "system-flow.yaml", string(manifest)+"orchestration:\n  release:\n    policy: threshold\n    count: 2\n")
	repo, err = workitem.Open(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	p := mustReasons(t, repo, "release")
	got := figureItem(t, p, FigureRelease)
	var ev release.Evaluation
	if err := json.Unmarshal([]byte(got.Text), &ev); err != nil {
		t.Fatalf("the release figures are not release_evaluate's JSON: %v\n%s", err, got.Text)
	}
	want, err := release.EvaluateRepo(execx.System{}, repo.Root, repo.Manifest, repo)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ev, want) || ev.Policy != "threshold" || ev.CountThreshold == nil || *ev.CountThreshold != 2 || ev.Met {
		t.Errorf("figures %+v; want %+v", ev, want)
	}
	if got.Reason != "release" || got.Size != len(got.Text) || len(p.Items) != 1 {
		t.Errorf("release item %+v; items %d", got, len(p.Items))
	}
}

// ADR-0217 §3, S-0537 criterion 1: ready loads the board, its columns with
// their limits, each ready story's hold, and the pull order, the backlog as
// a count, and the promotion candidates with why each other backlog story
// is not one, as flai's commands print them.
func TestAReadyReasonLoadsTheBoardAndThePromotionCandidates(t *testing.T) {
	repo, epic, story := reasonsRepo(t)
	now := time.Date(2026, 10, 29, 9, 0, 0, 0, time.UTC)
	body := "## Goal\n\nWork.\n\n## Acceptance criteria\n\n- [ ] it works\n"
	open, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Open story", Parent: epic.ID, Owner: "alex", Touches: []string{"flai/x.go"}, Now: now, Body: body})
	if err != nil {
		t.Fatal(err)
	}
	held, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Held story", Parent: epic.ID, Owner: "alex", Touches: []string{"flai/x.go"}, Now: now, Body: body})
	if err != nil {
		t.Fatal(err)
	}
	for _, mv := range []struct {
		it *workitem.Item
		to string
	}{{story, workitem.Ready}, {held, workitem.Ready}, {open, workitem.Ready}, {open, workitem.InProgress}} {
		if _, err := repo.Transition(mv.it, mv.to, "alex", "", now); err != nil {
			t.Fatal(err)
		}
	}
	p := mustReasons(t, repo, "ready")
	board := figureItem(t, p, FigureBoard)
	for _, want := range []string{
		"backlog: stories 1, each under promote_candidates\n",
		"ready (2/5)\n",
		"  " + story.ID + "  Roles prime  feature  ",
		"  " + held.ID + "  Held story  feature  ",
		" HELD\n      touches flai/x.go\n      ",
		"in-progress (1/2)\n  " + open.ID + "  Open story  ",
		"review (0/3)\n  -\n",
	} {
		if !strings.Contains(board.Text, want) {
			t.Errorf("board lacks %q:\n%s", want, board.Text)
		}
	}
	if strings.Contains(board.Text, "Second story") {
		t.Errorf("board lists the backlog's stories:\n%s", board.Text)
	}
	candidates := figureItem(t, p, FigureCandidates)
	for _, want := range []string{"no candidates by ", "not candidates:\n  S-0002  Second story\n    - no goal\n"} {
		if !strings.Contains(candidates.Text, want) {
			t.Errorf("candidates lack %q:\n%s", want, candidates.Text)
		}
	}
	if len(p.Items) != 2 || p.Items[0].Path != FigureBoard || board.Reason != "ready" || candidates.Size != len(candidates.Text) {
		t.Errorf("items %+v; want the board, then the candidates", p.Items)
	}
}

// candidatesText prints each candidate's rank, ID, figure, and title, and
// each other story's reasons, as flai promote --candidates does.
func TestCandidatesText(t *testing.T) {
	c := workitem.Candidates{Policy: "cod",
		Candidates: []workitem.Ranked{{Position: 1, ID: "S-0007", Title: "Worth most", Text: "500 EUR/week"}, {Position: 2, ID: "S-0003", Title: "Worth less", Text: "80 EUR/week"}},
		Others:     []workitem.Refused{{ID: "S-0009", Title: "Unready", Reasons: []string{"no goal", "no forecast duration"}}},
	}
	want := "candidates by cod:\n  1  S-0007  500 EUR/week  Worth most\n  2  S-0003  80 EUR/week  Worth less\nnot candidates:\n  S-0009  Unready\n    - no goal\n    - no forecast duration\n"
	if got := candidatesText(c); got != want {
		t.Errorf("candidates text\n%s\nwant\n%s", got, want)
	}
	if got := candidatesText(workitem.Candidates{Policy: "cod"}); got != "no backlog stories\n" {
		t.Errorf("no backlog: %q", got)
	}
}

// ADR-0104, ADR-0217: a pack of every figure reason parts within the limit,
// and its parts put back together are the pack.
func TestAFigurePackParts(t *testing.T) {
	repo, _, story := reasonsRepo(t)
	writeRepo(t, repo, ".flai-cache/verify/"+story.ID+".json", storedReport)
	p := mustReasons(t, repo, "review:"+story.ID, "release", "ready")
	if len(p.Items) != 5 {
		t.Fatalf("items %+v; want the story, its result, the release figures, the board, and the candidates", p.Items)
	}
	limit := encoded(p) / 3
	parts := p.Parts(limit)
	if len(parts) < 2 {
		t.Fatalf("the pack is in %d part; the test wants more", len(parts))
	}
	for _, q := range parts {
		if n := size(t, q); n > limit {
			t.Errorf("part %d encodes to %d bytes, over %d", q.PartNumber, n, limit)
		}
	}
	convs, items := rejoin(t, parts)
	if !reflect.DeepEqual(convs, p.Conventions) || !reflect.DeepEqual(items, p.Items) {
		t.Errorf("the parts' conventions and items are not the pack's")
	}
}
