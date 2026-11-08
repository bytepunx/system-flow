package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// badCheck is a tier command that fails, naming path:line, on each of its
// files that holds BAD.
const badCheck = `[sh, -c, 'rc=0; for f; do if grep -q BAD "$f"; then grep -n BAD "$f" | sed "s|^|$f:|"; rc=1; fi; done; exit $rc', sh, "{files}"]`

// mdCheck is a tier command that fails, printing a line as markdownlint-cli2
// does, on each of its files that holds BAD.
const mdCheck = `[sh, -c, 'rc=0; for f; do n=$(grep -n BAD "$f" | head -n 1 | cut -d: -f1); if [ -n "$n" ]; then echo "$f:$n:3 error MD009/no-trailing-spaces Trailing spaces [Expected: 0; Actual: 1]"; rc=1; fi; done; exit $rc', sh, "{files}"]`

// lintProject is syncProject with tiers appended to its manifest and
// committed. Integration: it runs real git.
func lintProject(t *testing.T, tiers string) string {
	t.Helper()
	root := syncProject(t)
	m, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), append(m, "tests:\n"+tiers...), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: lint tiers")
	return root
}

// storyStatus is the status line of story id's file under wip/kanban.
func storyStatus(t *testing.T, root, id string) string {
	t.Helper()
	matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
	if len(matches) != 1 {
		t.Fatalf("story file for %s: %v", id, matches)
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "status: ") {
			return strings.TrimPrefix(line, "status: ")
		}
	}
	return ""
}

// S-0367 criterion 2: the markdown flai accept commits from the main
// checkout, what was there before and what the acceptance writes, is fixed
// before it is committed.
func testAcceptFixesTheMainCheckoutsMarkdownBeforeItCommits(t *testing.T) {
	root := lintProject(t,
		"  - name: notes\n    command: [\"true\"]\n    paths: [\"wip/notes/**\"]\n"+
			"    fix_command: [sh, -c, 'for f; do printf \"good\\n\" > \"$f\"; done', sh, \"{files}\"]\n"+
			"  - name: archive\n    command: [\"true\"]\n    paths: [\"wip/archive/**\"]\n"+
			"    fix_command: [sh, -c, 'for f; do sed \"s/x] done/x] done (fixed)/\" \"$f\" > \"$f.t\" && mv \"$f.t\" \"$f\"; done', sh, \"{files}\"]\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/new.md", "new\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	if err := os.MkdirAll(filepath.Join(root, "wip/notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/notes/n.md"), []byte("bad\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	if got := gitIn(t, root, "show", "HEAD:wip/notes/n.md"); got != "good" {
		t.Errorf("the acceptance committed wip/notes/n.md as %q", got)
	}
	if got := gitIn(t, root, "grep", "-l", "done (fixed)", "HEAD", "--", "wip/archive"); !strings.Contains(got, "S-0001") {
		t.Errorf("the archived story was not fixed before the commit: %q", got)
	}
}

// TH-0413: a story branch whose markdown the fix after the rebase leaves
// failing is refused before the merge, with the file and line, and stays in
// review with main unchanged. The fix fixes nothing and its exit says what it
// left, as markdownlint-cli2 --fix does (S-0609).
func testAcceptRefusesABranchWhoseLintTheFixLeaves(t *testing.T) {
	root := lintProject(t, "  - name: md\n    command: "+badCheck+"\n    paths: [\"docs/**\"]\n    fix_command: "+badCheck+"\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/bad.md", "fine\nBAD line\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	before := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 {
		t.Fatal("accept passed a branch whose lint the fix left")
	}
	for _, want := range []string{"docs/bad.md:2:BAD line", "story/S-0001", "md tier", "flai accept S-0001 again"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal lacks %q:\n%s", want, errOut)
		}
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != before {
		t.Errorf("main moved from %s to %s", before, got)
	}
	if got := storyStatus(t, root, "S-0001"); got != "review" {
		t.Errorf("S-0001 is %q, not in review", got)
	}
	if _, err := os.Stat(wt); err != nil {
		t.Errorf("the worktree was removed: %v", err)
	}
}

// TH-0413: an uncommitted main-checkout file whose lint the fix leaves
// refuses the acceptance before anything changes. The fix cannot start, so
// its exit says nothing and the check runs after it (S-0609, ADR-0140).
func testAcceptRefusesTheMainCheckoutsLintBeforeAnythingChanges(t *testing.T) {
	root := lintProject(t, "  - name: notes\n    command: "+badCheck+"\n    paths: [\"wip/notes/**\"]\n    fix_command: [\"./no-such-fix\"]\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/new.md", "new\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	if err := os.MkdirAll(filepath.Join(root, "wip/notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/notes/n.md"), []byte("BAD\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before, branch := gitIn(t, root, "rev-parse", "main"), gitIn(t, root, "rev-parse", "story/S-0001")

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 {
		t.Fatal("accept passed an uncommitted file whose lint the fix left")
	}
	for _, want := range []string{"wip/notes/n.md:1:BAD", "main checkout", "flai accept S-0001 again"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal lacks %q:\n%s", want, errOut)
		}
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != before {
		t.Errorf("main moved from %s to %s", before, got)
	}
	if got := gitIn(t, root, "rev-parse", "story/S-0001"); got != branch {
		t.Errorf("story/S-0001 moved from %s to %s: it was rebased", branch, got)
	}
	if got := storyStatus(t, root, "S-0001"); got != "review" {
		t.Errorf("S-0001 is %q, not in review", got)
	}
}

// acceptRefusal runs flai accept id --json in root and answers its exit code
// and the refusal it printed.
func acceptRefusal(t *testing.T, root, id string) (int, unfixedRefusal) {
	t.Helper()
	out, errOut, code := runIn(t, root, "--json", "accept", id)
	var got struct {
		Refused unfixedRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("accept --json printed no refusal: %v\n%s\n%s", err, out, errOut)
	}
	return code, got.Refused
}

// S-0378: with --json, a branch whose markdown the fix leaves failing is
// refused with exit 4 and each finding as data, naming the story's branch as
// where to correct it.
func testAcceptJSONRefusesABranchsLintWithEachFinding(t *testing.T) {
	root := lintProject(t, "  - name: markdown\n    command: "+mdCheck+"\n    paths: [\"docs/**\"]\n    fix_command: "+mdCheck+"\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/bad.md", "fine\nBAD line\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")

	code, r := acceptRefusal(t, root, "S-0001")
	if code != exitDocRefused {
		t.Errorf("exit %d, want %d", code, exitDocRefused)
	}
	want := verify.FileFinding{Tier: "markdown", Checkout: "story", Story: "S-0001", Path: "docs/bad.md", Line: 2, Column: 3, Level: "error",
		Rule: "MD009", Alias: "no-trailing-spaces", Message: "Trailing spaces [Expected: 0; Actual: 1]"}
	if r.ID != "S-0001" || len(r.Findings) != 1 || r.Findings[0] != want {
		t.Errorf("the refusal's data: %+v, want the finding %+v", r, want)
	}
	if strings.HasPrefix(r.Reason, "rule:") || !strings.Contains(r.Reason, "docs/bad.md:2:3 error MD009") || !strings.Contains(r.Reason, "flai accept S-0001 again") {
		t.Errorf("the reason is not the text refusal: %q", r.Reason)
	}
	if got := storyStatus(t, root, "S-0001"); got != "review" {
		t.Errorf("S-0001 is %q, not in review", got)
	}
}

// S-0378: with --json, an uncommitted main-checkout file whose markdown the
// fix leaves failing is refused with exit 4, naming the main checkout.
func testAcceptJSONRefusesTheMainCheckoutsLintWithEachFinding(t *testing.T) {
	root := lintProject(t, "  - name: markdown\n    command: "+mdCheck+"\n    paths: [\"wip/notes/**\"]\n    fix_command: "+mdCheck+"\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/new.md", "new\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	if err := os.MkdirAll(filepath.Join(root, "wip/notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/notes/n.md"), []byte("BAD\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, r := acceptRefusal(t, root, "S-0001")
	if code != exitDocRefused {
		t.Errorf("exit %d, want %d", code, exitDocRefused)
	}
	if len(r.Findings) != 1 || r.Findings[0].Checkout != "main" || r.Findings[0].Story != "" ||
		r.Findings[0].Path != "wip/notes/n.md" || r.Findings[0].Line != 1 || r.Findings[0].Rule != "MD009" {
		t.Errorf("the refusal's findings: %+v", r.Findings)
	}
	// without --json the refusal is the text it was, with exit 1
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 1 || !strings.Contains(errOut, "wip/notes/n.md:1:3 error MD009") {
		t.Errorf("without --json: exit %d\n%s", code, errOut)
	}
}

// TestAcceptLintsEachSetOnce runs, as the agent tester, the tests of how
// often an acceptance lints each set (S-0609). Integration: they run real git.
func TestAcceptLintsEachSetOnce(t *testing.T) {
	t.Setenv("FLAI_AGENT", "tester")
	together(t,
		testAcceptLintsTheUncommittedSetAndTheBranchsSetOnceEach,
		testAcceptChecksABranchWhoseFixCannotSay,
	)
}

// countingTier is a tier over docs and wip/notes whose fix and check each
// log a line naming their files to log, and whose fix makes fixme fixed and
// exits 0, so nothing remains.
func countingTier(log string) string {
	return "  - name: md\n" +
		"    command: [sh, -c, 'echo \"check $*\" >> " + log + "', sh, \"{files}\"]\n" +
		"    paths: [\"docs/**\", \"wip/notes/**\"]\n" +
		"    fix_command: [sh, -c, 'echo \"fix $*\" >> " + log + "; for f; do sed s/fixme/fixed/ \"$f\" > \"$f.t\" && mv \"$f.t\" \"$f\"; done', sh, \"{files}\"]\n"
}

// runsNaming are the lines of the log at path that begin with kind and name
// file.
func runsNaming(t *testing.T, path, kind, file string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	n := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, kind+" ") && strings.Contains(line, file) {
			n++
		}
	}
	return n
}

// S-0609 criterion 1: an acceptance with markdown in the main checkout's
// uncommitted set and in the branch's set fixes each once, and checks
// neither after the fix, whose exit says nothing remains. The acceptance's
// own writes fall under no path the tier selects, so the fixes before its
// commits run nothing; before S-0609 the uncommitted file was fixed twice and
// checked once, and the branch's file fixed and checked.
func testAcceptLintsTheUncommittedSetAndTheBranchsSetOnceEach(t *testing.T) {
	log := filepath.Join(t.TempDir(), "lint.log")
	root := lintProject(t, countingTier(log))
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/new.md", "fixme\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	if err := os.MkdirAll(filepath.Join(root, "wip/notes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/notes/n.md"), []byte("fixme\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	for _, c := range []struct {
		kind, file string
		want       int
	}{{"fix", "wip/notes/n.md", 1}, {"check", "wip/notes/n.md", 0}, {"fix", "docs/new.md", 1}, {"check", "docs/new.md", 0}} {
		if got := runsNaming(t, log, c.kind, c.file); got != c.want {
			data, _ := os.ReadFile(log)
			t.Errorf("%s runs naming %s: %d, want %d; the log:\n%s", c.kind, c.file, got, c.want, data)
		}
	}
	if got := gitIn(t, root, "show", "HEAD:wip/notes/n.md"); got != "fixed" {
		t.Errorf("the acceptance committed wip/notes/n.md as %q", got)
	}
	if got := gitIn(t, root, "show", "main:docs/new.md"); got != "fixed" {
		t.Errorf("main has docs/new.md as %q", got)
	}
}

// S-0609, ADR-0140: a branch whose fix cannot start says nothing of what it
// left, so the check runs after it and refuses the acceptance before the
// merge.
func testAcceptChecksABranchWhoseFixCannotSay(t *testing.T) {
	root := lintProject(t, "  - name: md\n    command: "+badCheck+"\n    paths: [\"docs/**\"]\n    fix_command: [\"./no-such-fix\"]\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/bad.md", "fine\nBAD line\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	before := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 {
		t.Fatal("accept passed a branch whose fix could not start and whose check fails")
	}
	for _, want := range []string{"docs/bad.md:2:BAD line", "story/S-0001", "md tier"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal lacks %q:\n%s", want, errOut)
		}
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != before {
		t.Errorf("main moved from %s to %s", before, got)
	}
}

func TestUnfixedErrorNamesEachFinding(t *testing.T) {
	t.Parallel()
	res := verify.Result{Tiers: []verify.TierResult{
		{Name: "go", State: verify.Passed},
		{Name: "markdown", State: verify.Failed, Omitted: 2, Findings: []verify.Finding{
			{Name: "MD009", Path: "wip/a.md", Line: 3, Message: "trailing spaces"},
			{Path: "wip/b.md", Message: "no line"},
			{Name: "markdown", Message: "wip/c.md:1 MD041\nexit status 1"},
		}},
		{Name: "later", State: verify.NotReached},
	}}
	at := unfixedAt{root: "/repo", checkout: verify.CheckoutMain, where: "the main checkout's uncommitted files", fixIn: "in the main checkout"}
	err := unfixedError("S-0009", at, res)
	want := "rule: S-0009 cannot be accepted: the main checkout's uncommitted files hold lint the fix of the markdown tier could not fix:\n" +
		"  wip/a.md:3: trailing spaces (markdown, MD009)\n" +
		"  wip/b.md: no line (markdown)\n" +
		"  markdown: wip/c.md:1 MD041\n    exit status 1\n" +
		"  and 2 findings more\n" +
		"correct them in the main checkout and run flai accept S-0009 again"
	if err.Error() != want {
		t.Errorf("got:\n%s\nwant:\n%s", err, want)
	}
}
