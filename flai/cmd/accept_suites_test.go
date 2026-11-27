package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// acceptSuitesTiers are the tiers of the projects below: go-test of suite
// behavior and integration of suite integration, which select cli/**, the
// latter marked queue when queue is set, and lint of suite static, which
// selects docs/** alone.
func acceptSuitesTiers(queue bool) string {
	integration := "  - name: integration\n    suite: integration\n    command: [\"true\"]\n    paths: [\"cli/**\"]\n"
	if queue {
		integration += "    queue: true\n"
	}
	return "tests:\n" +
		"  - name: go-test\n    suite: behavior\n    command: [\"true\"]\n    paths: [\"cli/**\"]\n" +
		integration +
		"  - name: lint\n    suite: static\n    command: [\"true\"]\n    paths: [\"docs/**\"]\n"
}

// commitAcceptSuites commits on the main checkout at root the baseline
// workflow file with the verify of its acceptance replaced by suites, and
// returns main's new head.
func commitAcceptSuites(t *testing.T, root, suites string) string {
	t.Helper()
	base, err := os.ReadFile(filepath.Join("..", "internal", "workflow", "baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	const acceptVerify = "effects: [accept], rules: [children-closed, criteria-ticked], verify: all}"
	if !strings.Contains(string(base), acceptVerify) {
		t.Fatalf("the baseline workflow file's acceptance no longer gives %q", acceptVerify)
	}
	file := strings.Replace(string(base), acceptVerify, strings.Replace(acceptVerify, "verify: all", "verify: "+suites, 1), 1)
	return commitOnMain(t, root, "design/workflows/board.md", file, "chore: the acceptance needs "+suites)
}

// acceptSuitesProject is orchestratedStoryInReview, permit as given, with
// acceptSuitesTiers(queue) and a Go module committed on main, and story/S-0001
// rebased onto it; it answers the root and the branch's head.
func acceptSuitesProject(t *testing.T, permit, queue bool) (root, head string) {
	t.Helper()
	root, _ = orchestratedStoryInReview(t, permit)
	m, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	commitOnMain(t, root, "system-flow.yaml", string(m)+acceptSuitesTiers(queue), "chore: the tiers and their suites")
	commitOnMain(t, root, "go.mod", "module example.com/fixture\n", "chore: the Go module")
	return root, rebaseS0001(t, root)
}

// rebaseS0001 rebases story/S-0001 onto main in its worktree and returns its
// new head.
func rebaseS0001(t *testing.T, root string) string {
	t.Helper()
	gitIn(t, filepath.Join(root, ".flai-cache", "worktrees", "S-0001"), "rebase", "-q", "main")
	return strings.TrimSpace(gitIn(t, root, "rev-parse", "story/S-0001"))
}

// storeS0001SuitesVerify stores a passing flai verify report for S-0001 at
// head, run over main at base, of the suites given (nil for every suite): the
// branch changed cli/feature.go, and go-test and integration passed.
func storeS0001SuitesVerify(t *testing.T, root, head, base string, suites []string) {
	t.Helper()
	rep := verify.Report{Story: "S-0001", Commit: head, Base: "main", BaseCommit: base, MergeBase: base, Passed: true, Paths: []string{"cli/feature.go"}, Suites: suites, Steps: []verify.Step{
		{Name: "check", State: verify.Passed},
		{Name: "go-test", Tier: true, State: verify.Passed},
		{Name: "integration", Tier: true, State: verify.Passed},
	}}
	if err := verify.SaveReport(openProject(t, root), rep); err != nil {
		t.Fatal(err)
	}
}

// S-0808, ADR-0312 §3: verify_behind_main judges only the tiers of the suites
// the verify of the workflow's acceptance lists. While it lists static alone,
// a main that changed the story's Go package, selecting go-test and
// integration, which the verify ran, holds nothing; once it lists behavior,
// the refusal sends back go-test alone, not integration, and --dry-run names
// go-test alone among the stale tiers.
func TestVerifyBehindMainSendsBackOnlyTiersOfTheSuitesTheAcceptanceNeeds(t *testing.T) {
	t.Parallel()
	root, _ := acceptSuitesProject(t, false, false)
	commitAcceptSuites(t, root, "[static]")
	base := commitOnMain(t, root, "cli/conditions.go", "package cli\n\nconst Condition = \"the old words\"\n", "feat: a condition")
	head := rebaseS0001(t, root)
	storeS0001SuitesVerify(t, root, head, base, nil)

	commitOnMain(t, root, "cli/conditions.go", "package cli\n\nconst Condition = \"the new words\"\n", "feat: reword the condition")
	if res := dryRunAccept(t, root); res.VerifyBehindMainBlocker != nil || len(res.StaleTiers) > 0 || len(res.Blockers) > 0 {
		t.Errorf("a main whose overlapping paths select only tiers of suites the acceptance does not need must not hold it: %+v %q %q", res.VerifyBehindMainBlocker, res.StaleTiers, res.Blockers)
	}

	mainBefore := commitAcceptSuites(t, root, "[behavior]")
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, "send back go-test, which that verify ran: cli/conditions.go") || strings.Contains(errOut, "integration") {
		t.Fatalf("an acceptance that needs behavior must be refused sending back go-test alone: %d\n%s", code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
	res := dryRunAccept(t, root)
	if b := res.VerifyBehindMainBlocker; b == nil || b.Code != preview.BlockVerifyBehindMain || !slices.Equal(res.StaleTiers, []string{"go-test"}) || !slices.Contains(res.Blockers, b.Message) {
		t.Errorf("--dry-run must list verify_behind_main naming go-test alone: %+v %q %q", res.VerifyBehindMainBlocker, res.StaleTiers, res.Blockers)
	}
}

// S-0808, ADR-0312 §3: while the merge queue batches, queue_tiers judges only
// the queue tiers of the suites the acceptance needs: a story whose close-out
// left integration to the queue is accepted by hand while the acceptance
// needs behavior alone, and held while it needs integration too.
func TestQueueTiersIgnoresAQueueTierOfASuiteTheAcceptanceDoesNotNeed(t *testing.T) {
	t.Parallel()
	root, _ := acceptSuitesProject(t, true, true)
	mergeQueueOn(t, root)
	// held, so the hand acceptance lands rather than enqueues (S-0552)
	writeQueueState(t, root, mergequeue.State{Held: true})
	commitAcceptSuites(t, root, "[behavior, integration]")
	head := rebaseS0001(t, root)
	storeS0001Report(t, root, head, verify.Queued)

	if res := dryRunAccept(t, root); res.QueueTiersBlocker == nil || !slices.Equal(res.QueueTiers, []string{"integration"}) {
		t.Errorf("an acceptance that needs integration must be held on its queue tier: %+v %q %q", res.QueueTiersBlocker, res.QueueTiers, res.Blockers)
	}

	commitAcceptSuites(t, root, "[behavior]")
	head = rebaseS0001(t, root)
	storeS0001Report(t, root, head, verify.Queued)
	if res := dryRunAccept(t, root); res.QueueTiersBlocker != nil || len(res.QueueTiers) > 0 || res.QueueTiersWhy != "" || len(res.Blockers) > 0 {
		t.Errorf("a queue tier of a suite the acceptance does not need must not hold it: %+v %q %q %q", res.QueueTiersBlocker, res.QueueTiers, res.QueueTiersWhy, res.Blockers)
	}
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("with the tiers the acceptance needs passed at the head the story is accepted: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done {
		t.Errorf("the story must be accepted and archived: %s %v", story.Status, story.Archived)
	}
}

// S-0808, ADR-0312 §3: a story whose last flai verify ran suites that lack one
// the acceptance needs is refused, before anything is merged, with words
// naming the suite and flai verify S-0001 --job, and --dry-run shows it as the
// blocker verify_suites; on the baseline the same report holds nothing, and
// once a verify ran the suite the story is accepted.
func TestAnAcceptanceIsRefusedWhileTheLastVerifyDidNotRunASuiteItNeeds(t *testing.T) {
	t.Parallel()
	root, head := acceptSuitesProject(t, false, false)
	base := strings.TrimSpace(gitIn(t, root, "rev-parse", "main"))
	storeS0001SuitesVerify(t, root, head, base, []string{"behavior"})
	if res := dryRunAccept(t, root); res.VerifySuitesBlocker != nil || len(res.SuitesNotRun) > 0 || len(res.Blockers) > 0 {
		t.Errorf("on the baseline a verify of some suites must not hold the acceptance: %+v %q %q", res.VerifySuitesBlocker, res.SuitesNotRun, res.Blockers)
	}

	base = commitAcceptSuites(t, root, "[behavior, integration]")
	head = rebaseS0001(t, root)
	storeS0001SuitesVerify(t, root, head, base, []string{"behavior"})
	want := "S-0001 cannot be accepted yet: S-0001's last flai verify did not run suite integration, which its acceptance needs"
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, want) || !strings.Contains(errOut, "flai verify S-0001 --job") {
		t.Fatalf("an acceptance whose last verify did not run integration must be refused with %q: %d\n%s", want, code, errOut)
	}
	assertNothingAccepted(t, root, base)
	res := dryRunAccept(t, root)
	if b := res.VerifySuitesBlocker; b == nil || b.Code != preview.BlockVerifySuites || !slices.Equal(res.SuitesNotRun, []string{"integration"}) || !slices.Equal(res.Blockers, []string{b.Message}) {
		t.Errorf("--dry-run must list the blocker verify_suites naming integration: %+v %q %q", res.VerifySuitesBlocker, res.SuitesNotRun, res.Blockers)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run")
	if code != 0 || !strings.Contains(out, "blocked: S-0001's last flai verify did not run suite integration") {
		t.Errorf("flai accept --dry-run must show the refusal: %d %s\n%s", code, errOut, out)
	}

	storeS0001SuitesVerify(t, root, head, base, []string{"behavior", "integration"})
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("with a verify of the suites the acceptance needs the story is accepted: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done {
		t.Errorf("the story must be accepted and archived: %s %v", story.Status, story.Archived)
	}
}
