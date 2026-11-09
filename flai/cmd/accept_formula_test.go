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

// formulaStoryInReview is a project, accept_reviews on, whose S-0001 is in
// review with two ticked criteria and their two tasks, T-0001 done and T-0002
// done when stepDone, else in progress, and a branch that changes nothing
// beyond main; with formula, S-0001 is that formula's instance. head is the
// branch's head, which is main's.
func formulaStoryInReview(t *testing.T, formula string, stepDone bool) (root, head string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  permissions:\n    accept_reviews: true\n"
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Ship it", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-ship-it.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] run the check\n- [x] publish the report\n", 1)), 0o644)
	if formula != "" {
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		story, err := repo.Get("S-0001")
		if err != nil {
			t.Fatal(err)
		}
		story.Formula = formula
		if err := repo.Save(story); err != nil {
			t.Fatal(err)
		}
	}
	steps := [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Run the check", "--story", "S-0001"}, {"task", "new", "Publish the report", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"},
	}
	if stepDone {
		steps = append(steps, []string{"move", "T-0002", "done"})
	}
	for _, args := range append(steps, []string{"move", "S-0001", "review"}) {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	head = gitIn(t, root, "rev-parse", "story/S-0001")
	if main := gitIn(t, root, "rev-parse", "main"); head != main {
		t.Fatalf("the branch must change nothing beyond main: story/S-0001 %s, main %s", head, main)
	}
	return root, head
}

// S-0422 criterion 4: an instance of a formula whose steps changed no file
// has a branch with nothing beyond main, and flai accept accepts it as any
// story: the preview names no blocker, the merge is a fast-forward of
// nothing, and the story is moved to done, archived, and committed.
func TestFlaiAcceptsAnInstanceWhoseBranchChangesNothing(t *testing.T) {
	t.Parallel()
	root, _ := formulaStoryInReview(t, "weekly-check", true)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the dry run: %d %s\n%s", code, errOut, out)
	}
	var dry preview.Acceptance
	if err := json.Unmarshal([]byte(out), &dry); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if len(dry.Blockers) != 0 || dry.Branch != "story/S-0001" {
		t.Errorf("an empty branch blocks nothing: branch %q, blockers %q", dry.Branch, dry.Blockers)
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("the acceptance of an instance whose branch changes nothing: %d %s\n%s", code, errOut, out)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Status != workitem.Done || !res.Merged || res.Archived == 0 {
		t.Errorf("the instance is merged, done, and archived: %+v", res)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if !story.Archived || story.Status != workitem.Done || story.Formula != "weekly-check" {
		t.Errorf("the accepted instance is archived done, its formula kept: %+v", story)
	}
	if got := gitIn(t, root, "rev-parse", "main^"); got != mainBefore {
		t.Errorf("main gains the acceptance commit alone: its parent is %s, not %s", got, mainBefore)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); !strings.HasPrefix(subject, "chore: [S-0001] accept and archive") {
		t.Errorf("the acceptance commit: %q", subject)
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("the merged branch is deleted")
	}
}

// S-0422 criterion 4: the orchestrator accepts an instance whose branch
// changes nothing on the done transition of each step's task, "- <n>: T-nnnn
// done", in place of changed files, and writes that evidence in its Notes.
func TestTheOrchestratorAcceptsAnInstanceOnItsStepsTasks(t *testing.T) {
	t.Parallel()
	root, head := formulaStoryInReview(t, "weekly-check", true)
	evidence := "Verdict: pass, every step exited 0\n- 1: T-0001 done\n- 2: `T-0002` done"

	out, errOut, code := runStdin(t, root, evidence+"\n", "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-", "--json")
	if code != 0 {
		t.Fatalf("the orchestrator's acceptance of an instance: %d %s\n%s", code, errOut, out)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Status != workitem.Done || !res.Merged || res.By != "orchestrator" || res.Verified != head || len(res.OrchestratorBlockers) != 0 {
		t.Errorf("the orchestrator accepts the instance: %+v", res)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if !story.Archived || !strings.Contains(story.Body, "### Accepted by the orchestrator\n\n- Verified: "+head+"\n- At: 2026-09-19T02:00:00Z\n\n"+evidence+"\n") {
		t.Errorf("the evidence is in the archived instance's Notes:\n%s", story.Body)
	}
}

// S-0422 criterion 4: a step whose task is not done, or evidence that names
// no task of the instance, refuses the orchestrator's acceptance, naming the
// task, and the dry run lists it as criterion_unevidenced; a criterion's
// changed file is not an instance's evidence.
func TestTheOrchestratorRefusesAnInstanceWhoseStepIsNotDone(t *testing.T) {
	t.Parallel()
	root, head := formulaStoryInReview(t, "weekly-check", false)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	for _, c := range []struct{ evidence, want string }{
		{"Verdict: pass\n- 1: T-0001 done\n- 2: T-0002 done\n", "the evidence for acceptance criterion 2 (publish the report) names T-0002, which is in-progress, not done: its step has not passed"},
		{"Verdict: pass\n- 1: T-0001 done\n- 2: T-0042 done\n", "the evidence for acceptance criterion 2 (publish the report) names T-0042, which is not a task of S-0001"},
		{"Verdict: pass\n- 1: T-0001 done\n- 2: `report.md`\n", "the evidence for acceptance criterion 2 (publish the report) names no task's done transition (it says \"`report.md`\"): S-0001 is an instance of formula weekly-check, so write \"- 2: T-nnnn done\" naming the task of its step"},
		{"Verdict: pass\n- 1: T-0001 done\n", "the evidence has no item for acceptance criterion 2 (publish the report): S-0001 is an instance of formula weekly-check, so add \"- 2: T-nnnn done\" naming the task of its step"},
	} {
		evidence, want := c.evidence, c.want
		out, _, code := runStdin(t, root, evidence, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-", "--dry-run", "--json")
		if code != 0 {
			t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, out)
		}
		var res preview.Acceptance
		if err := json.Unmarshal([]byte(out), &res); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
		var unevidenced []string
		for _, b := range res.OrchestratorBlockers {
			if b.Code == preview.BlockCriterionUnevidenced {
				unevidenced = append(unevidenced, b.Message)
			}
		}
		if len(unevidenced) != 1 || unevidenced[0] != want {
			t.Errorf("the dry run must name criterion 2 alone as unevidenced:\n%q\nwant %q", unevidenced, want)
		}
		_, errOut, code := runStdin(t, root, evidence, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-")
		// the refusal is logged as JSON, its quotes escaped
		quoted, _ := json.Marshal(want)
		if code == 0 || !strings.Contains(errOut, string(quoted[1:len(quoted)-1])) {
			t.Errorf("the acceptance must be refused naming the step: %d\n%s", code, errOut)
		}
	}
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head)
	if code == 0 || !strings.Contains(errOut, "a Verdict: line and one - <n>: T-nnnn done item per criterion") {
		t.Errorf("an instance's acceptance without evidence names its form: %d\n%s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
}

// An ordinary story's evidence is checked as before: its criteria are met by
// changed files, so a task's done transition is no evidence for it, and an
// empty branch has no file to name.
func TestAnOrdinaryStorysEvidenceIsStillItsChangedFiles(t *testing.T) {
	t.Parallel()
	root, head := formulaStoryInReview(t, "", true)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runStdin(t, root, "Verdict: pass\n- 1: T-0001 done\n- 2: T-0002 done\n", "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-")
	if code == 0 || !strings.Contains(errOut, "the evidence for acceptance criterion 1 (run the check) names no file the branch changes (it names T-0001 done)") {
		t.Errorf("an ordinary story's evidence names changed files: %d\n%s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
}

// stoppedInstanceInReview is a git project, accept_reviews on, with the
// publish-like formula and its instance S-0001 in review, stopped at its
// evaluate step with exit 2 by flai formula stop (ADR-0271 §2): T-0001, the
// check, and T-0002, the evaluation, done, and T-0003 and T-0004 cancelled by
// the stop, every criterion ticked: the check's by hand, the rest by the
// stop. head is the head of story/S-0001, which
// is main's.
func stoppedInstanceInReview(t *testing.T) (root, head string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  permissions:\n    accept_reviews: true\n"
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design/workflows/formulas", "scripts"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, "scripts/test.sh"), []byte("#!/bin/sh\n"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design/workflows/formulas/publish-like.md"), []byte(publishLike), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{
		{"formula", "run", "publish-like", "--ready"},
		{"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"criteria", "tick", "S-0001", "1"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"},
	} {
		if out, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %d %s %s", args, code, out, errOut)
		}
	}
	if out, errOut, code := runStdin(t, root, "not met: 0 accepted stories not yet released\n", "formula", "stop", "S-0001", "--step", "evaluate", "--exit", "2"); code != 0 {
		t.Fatalf("formula stop: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatalf("the stopped instance goes to review: %d %s %s", code, out, errOut)
	}
	head = gitIn(t, root, "rev-parse", "story/S-0001")
	return root, head
}

// dropBranch removes story/S-0001 and its worktree, as I-0283's third
// instance found story/S-0585 gone.
func dropBranch(t *testing.T, root string) {
	t.Helper()
	if wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001"); dirExists(wt) {
		gitIn(t, root, "worktree", "remove", "--force", wt)
	}
	gitIn(t, root, "branch", "-D", "story/S-0001")
}

// stoppedEvidence is the orchestrator's evidence for the stopped instance:
// the steps that ran by their done tasks, the skipped ones by their tasks the
// stop cancelled (ADR-0271 §4).
const stoppedEvidence = "Verdict: pass, step evaluate stopped the run with exit 2\n- 1: T-0001 done\n- 2: T-0002 done\n- 3: T-0003 skipped\n- 4: `T-0004` skipped"

// acceptStopped runs the orchestrator's acceptance of S-0001 at verified
// with evidence, and answers its result and exit code.
func acceptStopped(t *testing.T, root, verified, evidence string, dry bool) (preview.Acceptance, int, string) {
	t.Helper()
	args := []string{"accept", "S-0001", "--by", "orchestrator", "--verified", verified, "--evidence", "-", "--json"}
	if dry {
		args = append(args, "--dry-run")
	}
	out, errOut, code := runStdin(t, root, evidence+"\n", args...)
	var res preview.Acceptance
	if code == 0 {
		if err := json.Unmarshal([]byte(out), &res); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	return res, code, errOut
}

// S-0621 criterion 1 (I-0283's second instance): an instance stopped cleanly
// at its evaluate step is accepted by the orchestrator with no operator's
// word, its skipped steps evidenced by their tasks the stop cancelled.
func TestTheOrchestratorAcceptsAnInstanceAStopEnded(t *testing.T) {
	t.Parallel()
	root, head := stoppedInstanceInReview(t)
	res, code, errOut := acceptStopped(t, root, head, stoppedEvidence, false)
	if code != 0 {
		t.Fatalf("the orchestrator's acceptance of a stopped instance: %d %s", code, errOut)
	}
	if res.Status != workitem.Done || res.By != "orchestrator" || res.Verified != head || len(res.OrchestratorBlockers) != 0 {
		t.Errorf("the orchestrator accepts the stopped instance: %+v", res)
	}
}

// S-0621 criterion 1 (I-0283's third instance): an instance whose branch is
// gone is accepted by the orchestrator on the commit its close-out passed
// at, on main, with nothing to merge; a commit not on main is refused.
func TestTheOrchestratorAcceptsAStoppedInstanceWhoseBranchIsGone(t *testing.T) {
	t.Parallel()
	root, head := stoppedInstanceInReview(t)
	dropBranch(t, root)

	stray := gitIn(t, root, "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", "not on main")
	res, code, errOut := acceptStopped(t, root, stray, stoppedEvidence, true)
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	if b := res.OrchestratorBlockers; len(b) != 1 || b[0].Code != preview.BlockUnverified || !strings.Contains(b[0].Message, "is not on main, and S-0001, a formula's instance, has no branch story/S-0001") {
		t.Errorf("a verified commit not on main blocks a branchless instance: %+v", b)
	}

	res, code, errOut = acceptStopped(t, root, head, stoppedEvidence, false)
	if code != 0 {
		t.Fatalf("the orchestrator's acceptance of a branchless instance: %d %s", code, errOut)
	}
	if res.Status != workitem.Done || res.Merged || res.Verified != head || len(res.OrchestratorBlockers) != 0 {
		t.Errorf("the orchestrator accepts the branchless instance, merging nothing: %+v", res)
	}
}

// S-0621: evidence that names a skipped step's task the stop did not cancel
// is refused, naming why; a story that is no instance and has no branch
// stays the operator's (ADR-0271 §4).
func TestTheOrchestratorRefusesSkippedEvidenceNoStopWrote(t *testing.T) {
	t.Parallel()
	root, head := stoppedInstanceInReview(t)
	res, code, errOut := acceptStopped(t, root, head, "Verdict: pass\n- 1: T-0001 skipped\n- 2: T-0002 done\n- 3: T-0003 skipped\n- 4: T-0004 skipped", true)
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	if b := res.OrchestratorBlockers; len(b) != 1 || b[0].Code != preview.BlockCriterionUnevidenced || !strings.Contains(b[0].Message, "names T-0001 skipped, but T-0001 is done, not cancelled") {
		t.Errorf("a done task is no skipped step's evidence: %+v", b)
	}
}

// S-0621 (ADR-0271 §4): a story that is no formula's instance and has no
// branch is still the operator's to accept.
func TestABranchlessStoryThatIsNoInstanceStaysTheOperators(t *testing.T) {
	t.Parallel()
	root, head := formulaStoryInReview(t, "", true)
	dropBranch(t, root)
	res, code, errOut := acceptStopped(t, root, head, "Verdict: pass\n- 1: `a.go`\n- 2: `a.go`", true)
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	found := false
	for _, b := range res.OrchestratorBlockers {
		found = found || (b.Code == preview.BlockUnverified && strings.Contains(b.Message, "has no branch story/S-0001 to verify here, so no commit its verifier passed can be its head; the operator accepts it"))
	}
	if !found {
		t.Errorf("a branchless story that is no instance is the operator's: %+v", res.OrchestratorBlockers)
	}
}
