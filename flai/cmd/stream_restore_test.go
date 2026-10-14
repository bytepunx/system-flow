package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// undoneStory is I-0124's sequence: S-0001 adds a line to docs/guide.md and
// reverts it, main makes the same change since, and flai stream sync replays
// the revert over main's change, cleanly, and records the path. It returns
// the project's root and S-0001's worktree. Integration: it runs real git.
func undoneStory(t *testing.T) (root, wt string) {
	t.Helper()
	root = syncProject(t)
	wt = openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/guide.md", "line one\nline two\n")
	commitIn(t, wt, "docs/guide.md", "line one\n")
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("line one\nline two\n"), 0o644)
	gitIn(t, root, "add", "docs/guide.md")
	gitIn(t, root, "commit", "-q", "-m", "feat: [S-0002] line two")
	if out, errOut, code := runIn(t, wt, "stream", "sync", "S-0001"); code != 0 {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	if got, _ := os.ReadFile(filepath.Join(wt, "docs", "guide.md")); string(got) != "line one\n" {
		t.Fatalf("the revert was not replayed over main's change: %q", got)
	}
	if got := recordedPaths(t, root); !reflect.DeepEqual(got, []string{"docs/guide.md"}) {
		t.Fatalf("the sync recorded %v, not docs/guide.md", got)
	}
	return root, wt
}

// recordedPaths is the paths S-0001's record holds that its branch still
// undoes.
func recordedPaths(t *testing.T, root string) []string {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	undone, err := storygit.RecordedUndone(execx.System{}, repo, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	for _, u := range undone {
		paths = append(paths, u.Path)
	}
	return paths
}

// restoreRefusalJSON is a refused restore's --json.
type restoreRefusalJSON struct {
	Refused restoreRefused `json:"refused"`
}

// ADR-0155, I-0124: flai stream restore commits main's version of a path the
// rebased branch undoes on the story's branch, so the branch holds main's
// lines again and the path leaves the record; a path the record no longer
// holds is then refused.
func TestStreamRestoreCommitsMainsVersionOfAnUndonePath(t *testing.T) {
	root, wt := undoneStory(t)

	out, errOut, code := runIn(t, wt, "--json", "stream", "restore", "S-0001", "--path", "./docs/guide.md", "--trailer", "Co-Authored-By: t <t@t>")
	if code != 0 {
		t.Fatalf("restore: %d %s %s", code, out, errOut)
	}
	var res restored
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	head := gitIn(t, wt, "rev-parse", "HEAD")
	want := restored{Story: "S-0001", Branch: "story/S-0001", Commit: head, Paths: []string{"docs/guide.md"}}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("json: %+v, want %+v", res, want)
	}
	if got, _ := os.ReadFile(filepath.Join(wt, "docs", "guide.md")); string(got) != "line one\nline two\n" {
		t.Errorf("docs/guide.md is %q, not main's", got)
	}
	if got := gitIn(t, wt, "log", "-1", "--format=%B"); got != "fix: [S-0001] restore main's docs/guide.md, which the rebased branch undid (ADR-0155)\n\nCo-Authored-By: t <t@t>" {
		t.Errorf("message: %q", got)
	}
	if got := gitIn(t, wt, "show", "--name-only", "--format=", "HEAD"); got != "docs/guide.md" {
		t.Errorf("the commit changed %q", got)
	}
	if got := gitIn(t, root, "diff", "main", "story/S-0001", "--", "docs/guide.md"); got != "" {
		t.Errorf("the branch still differs from main on docs/guide.md:\n%s", got)
	}
	if st := gitIn(t, wt, "status", "--porcelain"); st != "" {
		t.Errorf("the worktree is not clean: %s", st)
	}
	if got := recordedPaths(t, root); len(got) != 0 {
		t.Errorf("the record still holds %v", got)
	}
	repo, _ := workitem.Open(root)
	if _, err := os.Stat(storygit.UndoneRecordPath(repo.CacheDir(), "S-0001")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the emptied record is still there: %v", err)
	}

	out, errOut, code = runIn(t, wt, "stream", "restore", "S-0001", "--path", "docs/guide.md")
	if code != exitDocRefused || !strings.Contains(errOut, "flai stream restore S-0001 changed nothing: its record of the paths story/S-0001 undoes of the main branch holds none, so docs/guide.md is not one to restore") {
		t.Errorf("a second restore: %d %s %s", code, out, errOut)
	}
	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("a refused restore committed: %s, was %s", got, head)
	}
}

// The text output names the branch, the commit, and each path restored.
func TestStreamRestorePrintsWhatItRestored(t *testing.T) {
	_, wt := undoneStory(t)
	out, errOut, code := runIn(t, wt, "stream", "restore", "S-0001", "--path", "docs/guide.md")
	if code != 0 {
		t.Fatalf("restore: %d %s %s", code, out, errOut)
	}
	want := "restored main's version on story/S-0001 in " + short(gitIn(t, wt, "rev-parse", "HEAD")) + ":\n  docs/guide.md\n"
	if out != want {
		t.Errorf("output %q, want %q", out, want)
	}
}

// A worktree with uncommitted changes is refused, naming them, and nothing
// changes: no commit, the changes as they were, the record as it was.
func TestStreamRestoreRefusesUncommittedChangesChangingNothing(t *testing.T) {
	root, wt := undoneStory(t)
	head := gitIn(t, wt, "rev-parse", "HEAD")
	_ = os.WriteFile(filepath.Join(wt, "docs", "guide.md"), []byte("half done\n"), 0o644)
	_ = os.WriteFile(filepath.Join(wt, "notes.txt"), []byte("untracked\n"), 0o644)

	out, errOut, code := runIn(t, wt, "stream", "restore", "S-0001", "--path", "docs/guide.md")
	if code != exitDocRefused || !strings.Contains(errOut, "flai stream restore S-0001 changed nothing: its worktree .flai-cache/worktrees/S-0001 has uncommitted changes in docs/guide.md, notes.txt; commit them on story/S-0001, or discard them, then restore") {
		t.Errorf("restore: %d %s %s", code, out, errOut)
	}
	out, _, code = runIn(t, wt, "--json", "stream", "restore", "S-0001", "--path", "docs/guide.md")
	var res restoreRefusalJSON
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if code != exitDocRefused || res.Refused.Story != "S-0001" || !reflect.DeepEqual(res.Refused.Uncommitted, []string{"docs/guide.md", "notes.txt"}) || !strings.Contains(res.Refused.Reason, "uncommitted changes") {
		t.Errorf("json: %d %+v", code, res)
	}

	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("HEAD moved: %s, was %s", got, head)
	}
	for rel, want := range map[string]string{"docs/guide.md": "half done\n", "notes.txt": "untracked\n"} {
		if got, _ := os.ReadFile(filepath.Join(wt, rel)); string(got) != want {
			t.Errorf("%s is %q", rel, got)
		}
	}
	if got := recordedPaths(t, root); !reflect.DeepEqual(got, []string{"docs/guide.md"}) {
		t.Errorf("the record holds %v", got)
	}
}

// A path the record does not hold is refused, naming those it holds, and
// nothing changes, not even the paths given that it holds.
func TestStreamRestoreRefusesAPathTheRecordDoesNotHold(t *testing.T) {
	root, wt := undoneStory(t)
	head := gitIn(t, wt, "rev-parse", "HEAD")

	out, errOut, code := runIn(t, wt, "stream", "restore", "S-0001", "--path", "docs/guide.md", "--path", "docs/other.md")
	if code != exitDocRefused || !strings.Contains(errOut, "flai stream restore S-0001 changed nothing: its record of the paths story/S-0001 undoes of the main branch does not hold docs/other.md; it holds docs/guide.md") {
		t.Errorf("restore: %d %s %s", code, out, errOut)
	}
	out, _, code = runIn(t, wt, "--json", "stream", "restore", "S-0001", "--path", "docs/other.md")
	var res restoreRefusalJSON
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if code != exitDocRefused || !reflect.DeepEqual(res.Refused.NotHeld, []string{"docs/other.md"}) || !reflect.DeepEqual(res.Refused.Held, []string{"docs/guide.md"}) {
		t.Errorf("json: %d %+v", code, res)
	}

	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("HEAD moved: %s, was %s", got, head)
	}
	if got, _ := os.ReadFile(filepath.Join(wt, "docs", "guide.md")); string(got) != "line one\n" {
		t.Errorf("docs/guide.md was restored: %q", got)
	}
	if st := gitIn(t, wt, "status", "--porcelain"); st != "" {
		t.Errorf("the worktree is not clean: %s", st)
	}
	if got := recordedPaths(t, root); !reflect.DeepEqual(got, []string{"docs/guide.md"}) {
		t.Errorf("the record holds %v", got)
	}
}

// A story not in progress or in review, and an item that is not a story, are
// refused before git is asked anything.
func TestStreamRestoreRefusesAStoryNotInProgressOrReview(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Story", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}

	out, errOut, code := runIn(t, root, "--json", "stream", "restore", "S-0001", "--path", "docs/guide.md")
	var res restoreRefusalJSON
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if code != exitDocRefused || res.Refused.Status != "backlog" ||
		!strings.Contains(errOut, "flai stream restore S-0001 changed nothing: it is in backlog, and a path is restored only on the branch of a story in progress or in review") {
		t.Errorf("backlog: %d %+v %s", code, res, errOut)
	}
	if _, errOut, code := runIn(t, root, "stream", "restore", "E-0001", "--path", "docs/guide.md"); code != exitDocRefused || !strings.Contains(errOut, "flai stream restore E-0001 changed nothing: it is not a story; name the story (S-nnnn)") {
		t.Errorf("epic: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "stream", "restore", "S-0001"); code == 0 || !strings.Contains(errOut, "required flag(s)") || !strings.Contains(errOut, "path") {
		t.Errorf("no --path: %d %s", code, errOut)
	}
}
