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
func TestAcceptJSONRefusesABranchsLintWithEachFinding(t *testing.T) {
	root := lintProject(t, "  - name: markdown\n    command: "+mdCheck+"\n    paths: [\"docs/**\"]\n    fix_command: [\"true\"]\n")
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/bad.md", "fine\nBAD line\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")

	code, r := acceptRefusal(t, root, "S-0001")
	if code != exitDocRefused {
		t.Errorf("exit %d, want %d", code, exitDocRefused)
	}
	want := unfixedFinding{Tier: "markdown", Checkout: "story", Story: "S-0001", Path: "docs/bad.md", Line: 2, Column: 3, Level: "error",
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
func TestAcceptJSONRefusesTheMainCheckoutsLintWithEachFinding(t *testing.T) {
	root := lintProject(t, "  - name: markdown\n    command: "+mdCheck+"\n    paths: [\"wip/notes/**\"]\n    fix_command: [\"true\"]\n")
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
	at := unfixedAt{root: "/repo", checkout: checkoutMain, where: "the main checkout's uncommitted files", fixIn: "in the main checkout"}
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

// S-0378: each markdownlint-cli2 line of a plain tier's output is a finding
// with its path, line, column, level, rule, and alias when it gives them,
// and the checkout to correct it in; other lines are none.
func TestUnfixedErrorGivesEachMarkdownlintLineAsData(t *testing.T) {
	out := "lint-md: markdownlint-cli2 0.20.0\n" +
		"Summary: 3 error(s)\n" +
		"docs/a.md:3 error MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]\n" +
		"docs/a.md:7:1 error MD004/ul-style Unordered list style [Expected: asterisk; Actual: dash]\n" +
		"/wt/sub/b.md:1 MD041/first-line-heading/first-line-h1 First line in a file should be a top-level heading\n" +
		"c.md:2:5 warning MD013 Line length\n" +
		"exit status 1"
	res := verify.Result{Tiers: []verify.TierResult{
		{Name: "markdown", Dir: "sub", State: verify.Failed, Findings: []verify.Finding{{Name: "markdown", Message: out}}},
		{Name: "go", State: verify.Failed, Findings: []verify.Finding{{Name: "errcheck", Path: "x.go", Line: 4, Message: "unchecked "}}},
	}}
	at := unfixedAt{root: "/wt", checkout: checkoutStory, story: "S-0009", where: "the files story/S-0009 changes", fixIn: "in the worktree"}
	got := unfixedError("S-0009", at, res).Findings
	story := unfixedFinding{Checkout: "story", Story: "S-0009"}
	want := []unfixedFinding{
		{Path: "sub/docs/a.md", Line: 3, Level: "error", Rule: "MD012", Alias: "no-multiple-blanks", Message: "Multiple consecutive blank lines [Expected: 1; Actual: 2]"},
		{Path: "sub/docs/a.md", Line: 7, Column: 1, Level: "error", Rule: "MD004", Alias: "ul-style", Message: "Unordered list style [Expected: asterisk; Actual: dash]"},
		{Path: "sub/b.md", Line: 1, Rule: "MD041", Alias: "first-line-heading", Message: "First line in a file should be a top-level heading"},
		{Path: "sub/c.md", Line: 2, Column: 5, Level: "warning", Rule: "MD013", Message: "Line length"},
		{Tier: "go", Path: "x.go", Line: 4, Rule: "errcheck", Message: "unchecked"},
	}
	for i := range want {
		if want[i].Tier == "" {
			want[i].Tier = "markdown"
		}
		want[i].Checkout, want[i].Story = story.Checkout, story.Story
	}
	if len(got) != len(want) {
		t.Fatalf("got %d findings, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("finding %d:\n got %+v\nwant %+v", i, got[i], want[i])
		}
	}
}
