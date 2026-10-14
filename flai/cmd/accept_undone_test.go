package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// I-0124, ADR-0155: a story branch changed docs/guide.md and reverted it, main
// then made the same change, and the sync replayed the revert cleanly over
// it. The acceptance preview lists the path as data and as a blocker, the
// acceptance refuses it, merging nothing, and accepts the branch once main's
// version is committed back on it.
func TestAcceptRefusesABranchThatUndoesMainsChange(t *testing.T) {
	root, wt := undoingStoryInReview(t, true)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	acc := dryRunUndone(t, root, "accept", "S-0001", "--dry-run", "--json")
	if len(acc.UndoesMain) != 1 || acc.UndoesMain[0].Path != "docs/guide.md" || strings.Join(acc.UndoesMain[0].Lines, "|") != "line two" || acc.UndoesMain[0].OldBase == "" || acc.UndoesMain[0].KeptReason != "" {
		t.Errorf("the preview must list docs/guide.md and main's line two under undoes_main: %+v", acc.UndoesMain)
	}
	want := "story/S-0001 undoes the main branch's change to docs/guide.md, removing 1 line main added since the branch left it; restore main's version with flai stream restore S-0001 --path docs/guide.md, keep the removal with flai accept S-0001 --keep-removal docs/guide.md --reason \"<why>\""
	if len(acc.Blockers) != 1 || !strings.HasPrefix(acc.Blockers[0], want) {
		t.Errorf("the preview's one blocker must start %q: %q", want, acc.Blockers)
	}
	if _, rerr := sameAnswer(t, root, hostapi.Host{}, "accept.preview", `{"id":"S-0001"}`, "accept", "S-0001", "--dry-run"); rerr != nil {
		t.Errorf("accept.preview: %+v", rerr)
	}

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet: story/S-0001 undoes the main branch's change to docs/guide.md") {
		t.Errorf("acceptance must refuse the undone path: %d %s", code, errOut)
	}
	out, _, code := runIn(t, root, "accept", "S-0001", "--json")
	var refused struct {
		Refused undoneRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &refused); err != nil || code != exitDocRefused {
		t.Fatalf("--json must exit %d with the refusal as data: %d %v %s", exitDocRefused, code, err, out)
	}
	if r := refused.Refused; r.ID != "S-0001" || !strings.Contains(r.Reason, want) || len(r.UndoesMain) != 1 || r.UndoesMain[0].Path != "docs/guide.md" || strings.Join(r.UndoesMain[0].Lines, "|") != "line two" {
		t.Errorf("refusal: %+v", r)
	}
	unchanged(t, root, wt, mainBefore)

	// main's version committed back on the branch, as flai stream restore
	// does, the path leaves the record and the branch is accepted
	_ = os.WriteFile(filepath.Join(wt, "docs", "guide.md"), []byte("line one\nline two\n"), 0o644)
	gitIn(t, wt, "commit", "-q", "-am", "fix: [S-0001] restore main's line two")
	if acc := dryRunUndone(t, root, "accept", "S-0001", "--dry-run", "--json"); len(acc.UndoesMain) != 0 || len(acc.Blockers) != 0 {
		t.Errorf("a restored path is no blocker: %+v", acc)
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("a restored branch must be accepted: %d %s", code, errOut)
	}
	if got := gitIn(t, root, "show", "main:docs/guide.md"); got != "line one\nline two" {
		t.Errorf("main must keep its change: %q", got)
	}
}

// ADR-0155: the operator keeps the branch's removal with --keep-removal and a
// reason; the preview shows the reason and no blocker, and the acceptance
// goes through with the reason under ### Removals kept in the story's Notes.
func TestAcceptKeepsARemovalTheOperatorGivesAReasonFor(t *testing.T) {
	root, _ := undoingStoryInReview(t, true)
	keep := []string{"--keep-removal", "docs/guide.md", "--reason", "the story drops line two on purpose"}

	acc := dryRunUndone(t, root, append([]string{"accept", "S-0001", "--dry-run", "--json"}, keep...)...)
	if len(acc.Blockers) != 0 || len(acc.UndoesMain) != 1 || acc.UndoesMain[0].KeptReason != "the story drops line two on purpose" {
		t.Errorf("a kept path is listed with its reason and blocks nothing: %+v", acc)
	}
	if text, _, _ := runIn(t, root, append([]string{"accept", "S-0001", "--dry-run"}, keep...)...); !strings.Contains(text, "would keep the removal of the main branch's lines at docs/guide.md: the story drops line two on purpose") {
		t.Errorf("the dry run's text must name the kept path:\n%s", text)
	}
	out, errOut, code := runIn(t, root, append([]string{"accept", "S-0001"}, keep...)...)
	if code != 0 {
		t.Fatalf("a kept removal must be accepted: %d %s", code, errOut)
	}
	if !strings.Contains(out, "kept the removal of the main branch's lines at docs/guide.md") {
		t.Errorf("the acceptance must say what it kept:\n%s", out)
	}
	if got := gitIn(t, root, "show", "main:docs/guide.md"); got != "line one" {
		t.Errorf("main must take the kept removal: %q", got)
	}
	story, err := filepath.Glob(filepath.Join(root, "wip", "archive", "*", "*", "S-0001-*.md"))
	if err != nil || len(story) == 0 {
		story, _ = filepath.Glob(filepath.Join(root, "wip", "archive", "*", "S-0001-*.md"))
	}
	if len(story) != 1 {
		t.Fatalf("the archived story: %v", story)
	}
	body, _ := os.ReadFile(story[0])
	if want := "### Removals kept\n\n- `docs/guide.md`, kept by "; !strings.Contains(string(body), want) || !strings.Contains(string(body), "at 2026-09-15T21:00:00Z: the story drops line two on purpose\n") {
		t.Errorf("the story's Notes must record the kept removal:\n%s", body)
	}
	if got := gitIn(t, root, "status", "--porcelain"); got != "" {
		t.Errorf("the Notes must be in the acceptance commit: %s", got)
	}
}

// ADR-0155: a keep the operator may not give, or that names no recorded
// path, is refused before anything changes.
func TestAcceptRefusesAKeepItCannotHonour(t *testing.T) {
	root, wt := undoingStoryInReview(t, true)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--by", "orchestrator", "--keep-removal", "docs/guide.md", "--reason", "mine"}, "the orchestrator keeps no removal of the main branch's lines"},
		{[]string{"--keep-removal", "docs/other.md", "--reason", "why"}, "--keep-removal docs/other.md names a path story/S-0001 does not undo of the main branch's changes: its record holds only docs/guide.md"},
		{[]string{"--keep-removal", "docs/guide.md"}, "--keep-removal docs/guide.md needs a --reason"},
		{[]string{"--keep-removal", "docs/guide.md", "--reason", " "}, "the --reason for --keep-removal docs/guide.md is empty"},
		{[]string{"--keep-removal", "docs/guide.md", "--reason", "a", "--reason", "b"}, "2 reasons for 1 path"},
		{[]string{"--reason", "why"}, "--reason says why a --keep-removal is kept"},
	} {
		for _, dry := range []bool{true, false} {
			args := append([]string{"accept", "S-0001"}, c.args...)
			if dry {
				args = append(args, "--dry-run")
			}
			if _, errOut, code := runIn(t, root, args...); code == 0 || !strings.Contains(errOut, c.want) {
				t.Errorf("flai %v must be refused with %q: %d %s", args, c.want, code, errOut)
			}
		}
	}
	unchanged(t, root, wt, mainBefore)
}

// ADR-0155: a path only the acceptance's own rebase finds, main having made
// its change after the story's last sync, is refused as a recorded one is,
// with the path as data, merging nothing; the rebase records it, so the
// preview lists it after.
func TestAcceptRefusesAPathOnlyItsRebaseFinds(t *testing.T) {
	root, wt := undoingStoryInReview(t, false)
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("line one\nline two\n"), 0o644)
	gitIn(t, root, "commit", "-q", "-am", "feat: [S-0002] line two")
	mainBefore := gitIn(t, root, "rev-parse", "main")
	if acc := dryRunUndone(t, root, "accept", "S-0001", "--dry-run", "--json"); len(acc.UndoesMain) != 0 || len(acc.Blockers) != 0 {
		t.Fatalf("before the rebase nothing is recorded: %+v", acc)
	}

	out, _, code := runIn(t, root, "accept", "S-0001", "--json")
	var refused struct {
		Refused undoneRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &refused); err != nil || code != exitDocRefused {
		t.Fatalf("--json must exit %d with the refusal as data: %d %v %s", exitDocRefused, code, err, out)
	}
	if r := refused.Refused; r.ID != "S-0001" || !strings.Contains(r.Reason, "story/S-0001 undoes the main branch's change to docs/guide.md") || len(r.UndoesMain) != 1 || strings.Join(r.UndoesMain[0].Lines, "|") != "line two" {
		t.Errorf("refusal: %+v", r)
	}
	unchanged(t, root, wt, mainBefore)
	if acc := dryRunUndone(t, root, "accept", "S-0001", "--dry-run", "--json"); len(acc.UndoesMain) != 1 || len(acc.Blockers) != 1 {
		t.Errorf("the rebase's finding is recorded for the preview: %+v", acc)
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--keep-removal", "docs/guide.md", "--reason", "meant"); code != 0 {
		t.Errorf("the operator keeps what the rebase found: %d %s", code, errOut)
	}
}

// undoneAcceptance is the part of flai accept --dry-run --json these tests
// read.
type undoneAcceptance struct {
	Blockers   []string `json:"blockers"`
	UndoesMain []struct {
		Path       string   `json:"path"`
		Lines      []string `json:"lines"`
		OldBase    string   `json:"old_base"`
		KeptReason string   `json:"kept_reason"`
	} `json:"undoes_main"`
}

// dryRunUndone runs args, a dry run with --json, and reads its answer.
func dryRunUndone(t *testing.T, root string, args ...string) undoneAcceptance {
	t.Helper()
	js, errOut, code := runIn(t, root, args...)
	if code != 0 {
		t.Fatalf("flai %v: %d %s", args, code, errOut)
	}
	var acc undoneAcceptance
	if err := json.Unmarshal([]byte(js), &acc); err != nil {
		t.Fatalf("%v: %s", err, js)
	}
	return acc
}

// unchanged fails unless main is still at mainBefore with its line two, and
// S-0001 is in review with its branch and worktree.
func unchanged(t *testing.T, root, wt, mainBefore string) {
	t.Helper()
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}
	if got := gitIn(t, root, "show", "main:docs/guide.md"); got != "line one\nline two" {
		t.Errorf("main must keep its change: %q", got)
	}
	if after, _ := os.ReadFile(filepath.Join(root, "wip/kanban/stories/S-0001-undoing.md")); !strings.Contains(string(after), "status: review") {
		t.Error("a refused acceptance must leave the story in review")
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") == "" {
		t.Error("a refused acceptance must keep the story branch")
	}
	if _, err := os.Stat(wt); err != nil {
		t.Errorf("a refused acceptance must keep the worktree: %v", err)
	}
}

// undoingStoryInReview makes a project with story S-0001 in review whose
// branch added line two to docs/guide.md and reverted it, as S-0232 did in
// I-0124, and returns the main checkout and the story's worktree. With
// mainFirst, main makes the same change and the story is synced before
// review, which replays the revert over main's change and records the path;
// without it, main is left for the caller.
func undoingStoryInReview(t *testing.T, mainFirst bool) (root, wt string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.MkdirAll(filepath.Join(root, "docs"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("line one\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Undoing", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-undoing.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt = filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	guide := filepath.Join(wt, "docs", "guide.md")
	_ = os.WriteFile(guide, []byte("line one\nline two\n"), 0o644)
	gitIn(t, wt, "commit", "-q", "-am", "feat: [S-0001] line two")
	_ = os.WriteFile(guide, []byte("line one\n"), 0o644)
	gitIn(t, wt, "commit", "-q", "-am", "revert: [S-0001] line two, as main of the time")
	if mainFirst {
		_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("line one\nline two\n"), 0o644)
		gitIn(t, root, "commit", "-q", "-am", "feat: [S-0002] line two")
		if _, errOut, code := runIn(t, root, "stream", "sync", "S-0001"); code != 0 {
			t.Fatalf("flai stream sync: %s", errOut)
		}
		if got, _ := os.ReadFile(guide); string(got) != "line one\n" {
			t.Fatalf("the sync did not replay the revert over main's change: %q", got)
		}
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root, wt
}
