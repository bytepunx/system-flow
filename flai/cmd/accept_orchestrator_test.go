package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0221, ADR-0093: with accept_reviews on, the orchestrator accepts a story
// in review whose criteria are ticked, whose diff stays within its touches,
// with no thread open, at the head its verifier passed, and with evidence
// naming a changed file for every criterion. The done transition is its own,
// and its evidence is in the story's Notes in the acceptance commit.
func TestTheOrchestratorAcceptsAStoryItCanVouchFor(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	evidence := "Verdict: pass, tests and lint clean\n- 1: `cli/feature.go`"

	out, errOut, code := runStdin(t, root, evidence+"\n", "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-", "--json")
	if code != 0 {
		t.Fatalf("the orchestrator's acceptance: %d %s\n%s", code, errOut, out)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Status != workitem.Done || !res.Merged || res.By != "orchestrator" || res.Verified != head || res.Evidence == nil || res.Evidence.Verdict != "pass, tests and lint clean" || res.Evidence.Text != evidence {
		t.Errorf("--json must carry the acceptance, the commit verified, and the evidence: %+v", res)
	}

	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	var story *workitem.Item
	for _, it := range items {
		if it.ID == "S-0001" {
			story = it
		}
	}
	if story == nil || !story.Archived {
		t.Fatalf("the accepted story is archived: %+v", story)
	}
	last := story.Transitions[len(story.Transitions)-1]
	if last.To != workitem.Done || last.By != "orchestrator" {
		t.Errorf("the done transition is the orchestrator's: %+v", last)
	}
	section := "## Notes\n\n### Accepted by the orchestrator\n\n- Verified: " + head + "\n- At: 2026-09-19T02:00:00Z\n\n" + evidence + "\n"
	if !strings.Contains(story.Body, section) {
		t.Errorf("the Notes must carry the commit and the evidence:\n%s\nwant\n%s", story.Body, section)
	}
	committed := gitIn(t, root, "show", "HEAD:"+filepath.ToSlash(mustRel(t, root, story.Path)))
	if !strings.Contains(committed, "### Accepted by the orchestrator") || !strings.Contains(gitIn(t, root, "log", "-1", "--format=%s"), "chore: [S-0001] accept and archive") {
		t.Errorf("the evidence must be in the acceptance commit:\n%s", committed)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:cli/feature.go"), "package cli") {
		t.Error("the story branch is merged")
	}
}

// A thread open on the story stops the orchestrator, whatever it verified:
// the refusal names the thread, and nothing is merged.
func TestTheOrchestratorDoesNotAcceptAStoryWithAThreadOpen(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	if _, errOut, code := runIn(t, root, "thread", "new", "--on", "S-0001", "--by", "alex", "Is the flag right?", "Say which."); code != 0 {
		t.Fatal(errOut)
	}
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n"))
	if code == 0 {
		t.Fatal("a story with a thread open must not be accepted by the orchestrator")
	}
	if !strings.Contains(errOut, "S-0001 cannot be accepted yet") || !strings.Contains(errOut, "thread TH-0001 on S-0001 is open, not resolved: Is the flag right?") {
		t.Errorf("the refusal must name the thread:\n%s", errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
}

// With accept_reviews off, flai refuses the orchestrator's acceptance itself,
// naming the permission, through flai accept and flai move alike.
func TestTheOrchestratorAcceptsOnlyUnderAcceptReviews(t *testing.T) {
	root, head := orchestratedStoryInReview(t, false)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	want := "the orchestrator accepts a story only with orchestration.permissions.accept_reviews, which is off: ask the operator with thread_open on S-0001"
	for _, args := range [][]string{
		{"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev},
		{"accept", "S-0001", "--by", "orchestrator", "--verified", head, "--dry-run"},
		{"move", "S-0001", "done", "--by", "orchestrator", "--verified", head, "--evidence", ev},
	} {
		_, errOut, code := runIn(t, root, args...)
		if code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("flai %v must be refused naming accept_reviews: %d\n%s", args, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)
}

// The orchestrator's dry run names each condition that fails: an unticked
// criterion, a file outside the touches, and a verified commit that is no
// longer the head.
func TestTheOrchestratorsDryRunNamesEachBlocker(t *testing.T) {
	root, verified := orchestratedStoryInReview(t, true)
	story := filepath.Join(root, "wip/kanban/stories/S-0001-ship-it.md")
	b, _ := os.ReadFile(story)
	_ = os.WriteFile(story, []byte(strings.Replace(string(b), "- [x] the feature ships\n", "- [x] the feature ships\n- [ ] the guide says so\n", 1)), 0o644)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.MkdirAll(filepath.Join(wt, "other"), 0o755)
	_ = os.WriteFile(filepath.Join(wt, "other", "x.go"), []byte("package other\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] stray")
	head := gitIn(t, root, "rev-parse", "story/S-0001")
	mainBefore := gitIn(t, root, "rev-parse", "main")

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", verified, "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	for code, want := range map[string]string{
		preview.BlockCriterionUnticked: "acceptance criterion 2 of S-0001 is unticked: the guide says so",
		preview.BlockOutsideTouches:    "story/S-0001 changes other/x.go, which is under none of the touches of S-0001",
		preview.BlockUnverified:        "--verified " + verified[:12] + " is not the head of story/S-0001, which is " + head[:12],
	} {
		found := false
		for _, b := range res.OrchestratorBlockers {
			found = found || (b.Code == code && strings.HasPrefix(b.Message, want))
		}
		if !found || !strings.Contains(strings.Join(res.Blockers, "\n"), want) {
			t.Errorf("the dry run must name %s, %q: %+v", code, want, res.OrchestratorBlockers)
		}
	}
	if text, _, _ := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", verified, "--dry-run"); !strings.Contains(text, "blocked: story/S-0001 changes other/x.go") {
		t.Errorf("the dry run's text must name each blocker:\n%s", text)
	}
	assertNothingAccepted(t, root, mainBefore)
}

// The orchestrator's acceptance is refused without evidence, and with
// evidence that names no changed file for a criterion: a criterion it cannot
// check against the diff is never accepted.
func TestTheOrchestratorsAcceptanceNeedsEvidenceForEveryCriterion(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head)
	if code == 0 || !strings.Contains(errOut, "the orchestrator's acceptance of S-0001 needs its evidence: --evidence <file> or -, a Verdict: line and one - <n>: <files> item per criterion (ADR-0093)") {
		t.Errorf("an acceptance without evidence must be refused, saying what to give: %d\n%s", code, errOut)
	}
	unmet := evidenceFile(t, "Verdict: pass\n- 1: `cli/elsewhere.go`\n")
	_, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", unmet)
	if code == 0 || !strings.Contains(errOut, "the evidence for acceptance criterion 1 (the feature ships) names no file the branch changes (it names cli/elsewhere.go)") {
		t.Errorf("evidence naming no changed file must be refused: %d\n%s", code, errOut)
	}
	out, _, _ := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", unmet, "--dry-run", "--json")
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if len(res.OrchestratorBlockers) != 1 || res.OrchestratorBlockers[0].Code != preview.BlockCriterionUnevidenced {
		t.Errorf("the dry run must name the criterion unevidenced: %+v", res.OrchestratorBlockers)
	}
	_, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidenceFile(t, "- 1: `cli/feature.go`\n"))
	if code == 0 || !strings.Contains(errOut, "the evidence has no Verdict line") {
		t.Errorf("evidence without a verdict must be refused: %d\n%s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
}

// Under FLAI_ROLE=orchestrate the orchestrator accepts only as itself, and
// --verified and --evidence belong to its acceptance alone.
func TestAnAcceptanceUnderTheOrchestratorsRoleIsItsOwn(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--verified", head, "--by", "alex")
	if code == 0 || !strings.Contains(errOut, "--verified and --evidence are the orchestrator's acceptance: give them with --by orchestrator, or leave them out") {
		t.Errorf("--verified for another must be refused: %d\n%s", code, errOut)
	}
	t.Setenv("FLAI_ROLE", "orchestrate")
	for args, want := range map[string]string{
		"accept S-0001 --by alex":      "the orchestrator accepts a story only as itself, so give --by orchestrator, not --by alex (ADR-0093)",
		"accept S-0001 --dry-run":      "the orchestrator accepts a story only as itself, so give --by orchestrator, not no --by (ADR-0093)",
		"move S-0001 done --by alex":   "the orchestrator accepts a story only as itself, so give --by orchestrator, not --by alex (ADR-0093)",
		"move S-0001 done --by=olive2": "not --by olive2",
	} {
		_, errOut, code := runIn(t, root, strings.Fields(args)...)
		if code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("flai %s under FLAI_ROLE=orchestrate must be refused with %q: %d\n%s", args, want, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)
}

// The orchestrator's section goes at the end of the Notes, wherever they
// are, and a body without Notes gains them.
func TestTheOrchestratorsSectionGoesAtTheEndOfTheNotes(t *testing.T) {
	sec := "### Accepted by the orchestrator\n\nx\n"
	for body, want := range map[string]string{
		"## Goal\ng\n\n## Notes\n- a note\n":                "## Goal\ng\n\n## Notes\n- a note\n\n" + sec,
		"## Goal\ng\n\n## Notes\n- a note\n\n## Later\nl\n": "## Goal\ng\n\n## Notes\n- a note\n\n" + sec + "\n## Later\nl\n",
		"## Goal\ng\n": "## Goal\ng\n\n## Notes\n\n" + sec,
		"## Notes\n":   "## Notes\n\n" + sec,
	} {
		if got := withNotesSection(body, sec); got != want {
			t.Errorf("withNotesSection(%q):\n got %q\nwant %q", body, got, want)
		}
	}
}

// orchestratedStoryInReview makes a project, with accept_reviews on when
// permit is, whose story S-0001, touching cli and with its one criterion
// ticked, is in review with cli/feature.go committed on its branch, and
// returns the main checkout and the branch's head.
func orchestratedStoryInReview(t *testing.T, permit bool) (root, head string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_ROLE", "")
	root = tempProject(t)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if permit {
		m += "orchestration:\n  permissions:\n    accept_reviews: true\n"
	}
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design", "cli"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Ship it", "--epic", "E-0001", "--touches", "cli"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-ship-it.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] the feature ships\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(wt, "cli", "feature.go"), []byte("package cli\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] ship it")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root, gitIn(t, root, "rev-parse", "story/S-0001")
}

// evidenceFile writes text to a file outside the project and returns its path.
func evidenceFile(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "evidence.md")
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// assertNothingAccepted fails unless main is still at mainBefore and S-0001
// is still in review with its branch.
func assertNothingAccepted(t *testing.T, root, mainBefore string) {
	t.Helper()
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "wip/kanban/stories/S-0001-ship-it.md")); !strings.Contains(string(b), "status: review") {
		t.Error("a refused acceptance must leave the story in review")
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") == "" {
		t.Error("a refused acceptance must keep the story branch")
	}
}

// mustRel is path relative to root.
func mustRel(t *testing.T, root, path string) string {
	t.Helper()
	rel, err := filepath.Rel(root, path)
	if err != nil {
		t.Fatal(err)
	}
	return rel
}
