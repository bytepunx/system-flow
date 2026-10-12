package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/verify"
)

// badCheck is a tier command that fails, naming path:line, on each of its
// files that holds BAD.
const badCheck = `[sh, -c, 'rc=0; for f; do if grep -q BAD "$f"; then grep -n BAD "$f" | sed "s|^|$f:|"; rc=1; fi; done; exit $rc', sh, "{files}"]`

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
func TestAcceptFixesTheMainCheckoutsMarkdownBeforeItCommits(t *testing.T) {
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
// review with main unchanged.
func TestAcceptRefusesABranchWhoseLintTheFixLeaves(t *testing.T) {
	root := lintProject(t, "  - name: md\n    command: "+badCheck+"\n    paths: [\"docs/**\"]\n    fix_command: [\"true\"]\n")
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
// refuses the acceptance before anything changes.
func TestAcceptRefusesTheMainCheckoutsLintBeforeAnythingChanges(t *testing.T) {
	root := lintProject(t, "  - name: notes\n    command: "+badCheck+"\n    paths: [\"wip/notes/**\"]\n    fix_command: [\"true\"]\n")
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

func TestUnfixedErrorNamesEachFinding(t *testing.T) {
	res := verify.Result{Tiers: []verify.TierResult{
		{Name: "go", State: verify.Passed},
		{Name: "markdown", State: verify.Failed, Omitted: 2, Findings: []verify.Finding{
			{Name: "MD009", Path: "wip/a.md", Line: 3, Message: "trailing spaces"},
			{Path: "wip/b.md", Message: "no line"},
			{Name: "markdown", Message: "wip/c.md:1 MD041\nexit status 1"},
		}},
		{Name: "later", State: verify.NotReached},
	}}
	err := unfixedError("S-0009", "the main checkout's uncommitted files", "in the main checkout", res)
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
