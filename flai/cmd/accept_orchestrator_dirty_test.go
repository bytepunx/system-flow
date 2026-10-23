package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0189, I-0150: the orchestrator accepts on a main checkout holding
// uncommitted files outside wip that the story does not change, and leaves
// them as they were, uncommitted and out of the acceptance commit, where the
// operator's acceptance by hand still refuses them unless --yes.
func TestTheOrchestratorAcceptsOnADirtyMainAndLeavesTheFilesAsTheyAre(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	dirty := dirtyMain(t, root)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 || !strings.Contains(errOut, "uncommitted changes outside wip") || !strings.Contains(errOut, "--yes") {
		t.Fatalf("the operator's acceptance by hand on a dirty checkout must refuse, naming --yes: %d %s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)

	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	dry, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev, "--dry-run")
	if code != 0 || !strings.Contains(dry, "the orchestrator's acceptance leaves them as they are") || strings.Contains(dry, "blocked:") {
		t.Errorf("the orchestrator's dry run says the files are left as they are, with nothing blocked: %d %s\n%s", code, errOut, dry)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev)
	if code != 0 {
		t.Fatalf("the orchestrator's acceptance on a dirty checkout: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done {
		t.Errorf("S-0001 must be accepted and archived: %s, archived %v", story.Status, story.Archived)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:cli/feature.go"), "package cli") {
		t.Error("the story's change must be on main")
	}
	if got := porcelainOutsideWip(t, root); !slices.Equal(got, dirty) {
		t.Errorf("uncommitted after the acceptance %v, want the operator's files %v as they were", got, dirty)
	}
	for _, p := range strings.Split(gitIn(t, root, "show", "--name-only", "--format=", "HEAD"), "\n") {
		if isLeft(p, dirty) {
			t.Errorf("the acceptance commit holds the operator's file %s", p)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(root, "design/.gitkeep")); string(b) != "the operator's edit\n" {
		t.Errorf("the operator's edit was changed: %q", b)
	}
}

// ADR-0189: a story whose change would overwrite a file the main checkout
// has not committed is refused the orchestrator before anything changes, and
// its dry run lists the blocker dirty_main, naming the file and the command
// that sets it aside.
func TestTheOrchestratorIsRefusedAStoryWhoseChangeWouldOverwriteADirtyFile(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	dirtyMain(t, root)
	if err := os.WriteFile(filepath.Join(root, "cli/feature.go"), []byte("the operator's own\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mainBefore := gitIn(t, root, "rev-parse", "main")
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	want := "S-0001's change would overwrite cli/feature.go, which the main checkout has not committed, so the orchestrator does not accept it: commit them, or set them aside with git stash push -- cli/feature.go, and it is accepted once they are clean (ADR-0189)"

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev, "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	found := false
	for _, b := range res.OrchestratorBlockers {
		found = found || (b.Code == preview.BlockDirtyMain && b.Message == want)
	}
	if !found || !slices.Contains(res.Blockers, want) {
		t.Errorf("the dry run must list %s, %q: %+v", preview.BlockDirtyMain, want, res.OrchestratorBlockers)
	}
	if len(res.OrchestratorBlockers) != 1 {
		t.Errorf("the files the story does not change block nothing: %+v", res.OrchestratorBlockers)
	}

	_, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev)
	if code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet: "+want) {
		t.Errorf("the orchestrator's acceptance must be refused naming the file: %d\n%s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
	if b, _ := os.ReadFile(filepath.Join(root, "cli/feature.go")); string(b) != "the operator's own\n" {
		t.Errorf("the operator's file was changed: %q", b)
	}
}

// failCommits makes every git commit in the main checkout at root fail, as
// one does while another process holds the index lock beyond flai's wait,
// and returns the hook to remove.
func failCommits(t *testing.T, root string) string {
	t.Helper()
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	if err := os.MkdirAll(filepath.Dir(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hook, []byte("#!/bin/sh\necho 'pre-commit refuses the commit' >&2\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return hook
}

// ADR-0189, I-0150 (S-0338): the orchestrator finishes an acceptance it made
// itself that stopped before its commit, on a main checkout that has gained
// uncommitted files since, which it leaves as they are.
func TestTheOrchestratorFinishesItsOwnAcceptanceThatStoppedBeforeItsCommit(t *testing.T) {
	root, head := orchestratedStoryInReview(t, true)
	hook := failCommits(t, root)
	ev := evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n")
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", ev); code == 0 || !strings.Contains(errOut, "pre-commit refuses the commit") {
		t.Fatalf("the orchestrator's acceptance stops at its commit: %d %s", code, errOut)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Transitions[len(story.Transitions)-1].By != "orchestrator" {
		t.Fatalf("S-0001 is done and archived by the orchestrator: %+v", story.Transitions)
	}
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	dirty := dirtyMain(t, root)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--dry-run")
	if code != 0 || !strings.Contains(out, "would commit the acceptance of S-0001, done and archived but not committed") {
		t.Errorf("the orchestrator's dry run says it would commit: %d %s\n%s", code, errOut, out)
	}
	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "orchestrator")
	if code != 0 || !strings.Contains(out, "completed the acceptance of S-0001: done and archived already, committed") {
		t.Fatalf("the orchestrator finishes its own acceptance: %d %s\n%s", code, errOut, out)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the usual subject: %q", subject)
	}
	if got := porcelainOutsideWip(t, root); !slices.Equal(got, dirty) {
		t.Errorf("uncommitted after the acceptance %v, want the operator's files %v as they were", got, dirty)
	}
	if st := workitem.PorcelainPaths(gitIn(t, root, "status", "--porcelain")); !slices.Equal(st, dirty) {
		t.Errorf("everything but the operator's files is committed: %v", st)
	}
}

// ADR-0189: an acceptance that stopped before its commit, made by anyone
// but the orchestrator, stays the operator's to finish, and the refusal of
// the orchestrator names who made it.
func TestTheOrchestratorDoesNotFinishAnAcceptanceSomeoneElseMade(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, true)
	hook := failCommits(t, root)
	if _, _, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 {
		t.Fatal("the operator's acceptance stops at its commit")
	}
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	head := gitIn(t, root, "rev-parse", "HEAD")
	want := "S-0001 is done and archived but its acceptance was never committed, and alex moved it to done; the orchestrator finishes only an acceptance it made itself, so the operator completes this acceptance with flai accept S-0001 (ADR-0189)"
	for _, args := range [][]string{
		{"accept", "S-0001", "--by", "orchestrator", "--dry-run"},
		{"accept", "S-0001", "--by", "orchestrator"},
	} {
		if _, errOut, code := runIn(t, root, args...); code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("flai %v must be refused naming alex: %d\n%s", args, code, errOut)
		}
	}
	if got := gitIn(t, root, "rev-parse", "HEAD"); got != head {
		t.Errorf("a refused acceptance commits nothing: HEAD moved from %s to %s", head, got)
	}
}
