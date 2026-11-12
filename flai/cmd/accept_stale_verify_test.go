package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// commitOnMain writes content to rel in the main checkout at root and
// commits that path alone, leaving the checkout's other changes as they are,
// and returns main's new head.
func commitOnMain(t *testing.T, root, rel, content, msg string) string {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "--", rel)
	gitIn(t, root, "commit", "-q", "-m", msg, "--", rel)
	return gitIn(t, root, "rev-parse", "main")
}

// storeS0001PassingVerify stores a passing flai verify report for S-0001 at
// head, run over main at base, as a close-out leaves it: the branch changed
// cli/feature.go, and its go-test tier passed.
func storeS0001PassingVerify(t *testing.T, root, head, base string) {
	t.Helper()
	rep := verify.Report{Story: "S-0001", Commit: head, Base: "main", BaseCommit: base, MergeBase: base, Passed: true, Paths: []string{"cli/feature.go"}, Steps: []verify.Step{
		{Name: "check", State: verify.Passed},
		{Name: "go-test", Tier: true, State: verify.Passed},
	}}
	if err := verify.SaveReport(openProject(t, root), rep); err != nil {
		t.Fatal(err)
	}
}

// dryRunAccept is what flai accept S-0001 --by alex --dry-run --json says.
func dryRunAccept(t *testing.T, root string) accepted {
	t.Helper()
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("flai accept --dry-run: %d %s\n%s", code, errOut, out)
	}
	var res accepted
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return res
}

// I-0229, ADR-0291: a story whose close-out passed over an older main, while
// main has since changed a file its go-test tier selects, as S-0535 reworded
// a condition S-0416's guard test copied, is not accepted: flai accept
// refuses it with verify_behind_main, naming the tier and the path, before
// anything is merged, so main is not left red. A main that brought in only
// wip paths does not hold it, and once a close-out passes over the new main
// it lands.
func TestAnAcceptanceIsRefusedWhileMainOutranTheStorysVerify(t *testing.T) {
	t.Parallel()
	root, _ := orchestratedStoryInReview(t, false)
	m, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	commitOnMain(t, root, "system-flow.yaml", string(m)+"tests:\n  - name: go-test\n    command: [\"true\"]\n    paths: [\"cli/**\"]\n", "chore: the go-test tier")
	base := commitOnMain(t, root, "cli/conditions.go", "package cli\n\nconst Condition = \"the old words\"\n", "feat: a condition")
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	gitIn(t, wt, "rebase", "-q", "main")
	head := gitIn(t, root, "rev-parse", "story/S-0001")
	storeS0001PassingVerify(t, root, head, base)

	commitOnMain(t, root, "wip/agents/.gitkeep", "another story's narrative\n", "chore: wip only")
	if res := dryRunAccept(t, root); res.VerifyBehindMainBlocker != nil || len(res.StaleTiers) > 0 || len(res.Blockers) > 0 {
		t.Errorf("a main that brought in only wip paths must not hold the acceptance: %+v %q %q", res.VerifyBehindMainBlocker, res.StaleTiers, res.Blockers)
	}

	mainBefore := commitOnMain(t, root, "cli/conditions.go", "package cli\n\nconst Condition = \"the new words\"\n", "feat: reword the condition")
	want := "S-0001 cannot be accepted yet: S-0001's verify does not cover the main branch it would land on: its last flai verify, at " + head[:7] + ", ran over main at " + base[:7] + ", and the 1 path main brought in since send back go-test, which that verify ran: cli/conditions.go; ask for its close-out again"
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, want) || !strings.Contains(errOut, "flai verify S-0001 --job") || !strings.Contains(errOut, "(ADR-0291)") {
		t.Fatalf("an acceptance over a main its verify did not cover must be refused with %q: %d\n%s", want, code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
	if got := gitIn(t, root, "rev-parse", "story/S-0001"); got != head {
		t.Errorf("a refused acceptance must leave story/S-0001 at %s, not %s", head, got)
	}
	res := dryRunAccept(t, root)
	if b := res.VerifyBehindMainBlocker; b == nil || b.Code != preview.BlockVerifyBehindMain || !slices.Equal(res.StaleTiers, []string{"go-test"}) || !slices.Contains(res.Blockers, b.Message) {
		t.Errorf("--dry-run must list the blocker verify_behind_main naming go-test: %+v %q %q", res.VerifyBehindMainBlocker, res.StaleTiers, res.Blockers)
	}

	// the close-out again: synced onto the new main and passed over it
	gitIn(t, wt, "rebase", "-q", "main")
	head = gitIn(t, root, "rev-parse", "story/S-0001")
	storeS0001PassingVerify(t, root, head, mainBefore)
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("with a verify over the new main the story is accepted: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done {
		t.Errorf("the story must be accepted and archived: %s %v", story.Status, story.Archived)
	}
}
