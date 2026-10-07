package storygit

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Integration: these run real git.
func gitRepo(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	git(t, dir, "init", "-q", "-b", "main")
	git(t, dir, "config", "user.email", "t@t")
	git(t, dir, "config", "user.name", "t")
	write(t, dir, "a.md", "one\n")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "init")
	return dir
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, rel), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	if out, err := (execx.System{}).Run(dir, "git", args...); err != nil {
		t.Fatal(out, err)
	}
}

func TestUncommittedRebaseInProgressAndConflicts(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	if got, err := Uncommitted(r, dir); err != nil || len(got) != 0 {
		t.Fatalf("clean: %v %v", got, err)
	}
	if RebaseInProgress(r, dir) || len(Conflicts(r, dir)) != 0 {
		t.Fatal("clean repository reports a rebase")
	}

	write(t, dir, "a.md", "changed\n")
	write(t, dir, "new.md", "new\n")
	if got, err := Uncommitted(r, dir); err != nil || strings.Join(got, ",") != "a.md,new.md" {
		t.Fatalf("uncommitted: %v %v", got, err)
	}
	git(t, dir, "checkout", "-q", "--", "a.md")
	_ = os.Remove(filepath.Join(dir, "new.md"))

	// a branch and main change the same line, so the rebase stops
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	write(t, dir, "a.md", "story\n")
	git(t, dir, "commit", "-q", "-am", "story")
	git(t, dir, "checkout", "-q", "main")
	write(t, dir, "a.md", "main\n")
	git(t, dir, "commit", "-q", "-am", "main")
	git(t, dir, "checkout", "-q", "story/S-0001")
	if _, err := r.Run(dir, "git", "rebase", "main"); err == nil {
		t.Fatal("rebase did not stop")
	}
	if !RebaseInProgress(r, dir) {
		t.Error("stopped rebase not reported")
	}
	if got := Conflicts(r, dir); strings.Join(got, ",") != "a.md" {
		t.Errorf("conflicts: %v", got)
	}
	git(t, dir, "rebase", "--abort")
	if RebaseInProgress(r, dir) {
		t.Error("aborted rebase still reported")
	}
}

// ContinueRebase continues a resolved stop with no editor, stopping again on
// a later commit's conflict and finishing after the last.
func TestContinueRebase(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	// no editor can run here: one that did would fail the continue
	git(t, dir, "config", "core.editor", "false")
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	write(t, dir, "a.md", "story\n")
	git(t, dir, "commit", "-q", "-am", "story one")
	write(t, dir, "a.md", "story again\n")
	git(t, dir, "commit", "-q", "-am", "story two")
	git(t, dir, "checkout", "-q", "main")
	write(t, dir, "a.md", "main\n")
	git(t, dir, "commit", "-q", "-am", "main")
	git(t, dir, "checkout", "-q", "story/S-0001")
	if _, err := r.Run(dir, "git", "rebase", "main"); err == nil {
		t.Fatal("rebase did not stop")
	}

	// the first commit resolved, the second stops on the same line
	write(t, dir, "a.md", "main and story\n")
	git(t, dir, "add", "a.md")
	if err := ContinueRebase(r, dir); err == nil {
		t.Fatal("the second commit's conflict did not stop the rebase")
	}
	if !RebaseInProgress(r, dir) || strings.Join(Conflicts(r, dir), ",") != "a.md" {
		t.Fatalf("not stopped on a.md: %v", Conflicts(r, dir))
	}
	write(t, dir, "a.md", "main and story again\n")
	git(t, dir, "add", "a.md")
	if err := ContinueRebase(r, dir); err != nil {
		t.Fatalf("continue: %v", err)
	}
	if RebaseInProgress(r, dir) {
		t.Fatal("the rebase is still in progress")
	}
	out, err := r.Run(dir, "git", "log", "--format=%s", "main..HEAD")
	if err != nil || out != "story two\nstory one" {
		t.Errorf("commits after the rebase: %q %v", out, err)
	}
	if got, _ := os.ReadFile(filepath.Join(dir, "a.md")); string(got) != "main and story again\n" {
		t.Errorf("a.md: %q", got)
	}
}

// syncClock dates the conflict threads Sync writes.
var syncClock = time.Date(2026, 10, 6, 9, 0, 0, 0, time.UTC)

// syncProject is a project in a git repository whose main branch has
// docs/guide.md, for stories to sync. Integration: it runs real git.
func syncProject(t *testing.T) *workitem.Repo {
	t.Helper()
	dir := gitRepo(t)
	write(t, dir, "system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	write(t, dir, ".gitignore", ".flai-cache/\n")
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "docs"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, dir, "docs/guide.md", "line one\n")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "project")
	repo, err := workitem.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// openStory saves story id in progress with touches and checks its branch
// out from main in its worktree, which it returns.
func openStory(t *testing.T, repo *workitem.Repo, id string, touches ...string) (*workitem.Item, string) {
	t.Helper()
	it := &workitem.Item{ID: id, Type: workitem.Story, Nature: "feature", Title: "Story " + id, Status: workitem.InProgress, Touches: touches}
	if err := repo.Save(it); err != nil {
		t.Fatal(err)
	}
	wt := repo.WorktreePath(id)
	git(t, repo.MainRoot, "worktree", "add", "-q", "-b", Branch(id), wt, "main")
	return it, wt
}

// syncOptions are the options to sync story with, its generated files gen.
func syncOptions(repo *workitem.Repo, story *workitem.Item, gen ...GeneratedFile) SyncOptions {
	return SyncOptions{Runner: execx.System{}, Repo: repo, Story: story, Now: syncClock, Generated: gen}
}

// A clean sync rebases the branch onto main and lists what it changed
// outside the story's claim; with no other open story there is nothing to
// trial-merge.
func TestSyncRebasesOntoMainAndListsWhatIsOutsideTheClaim(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] guide", "docs/guide.md", "README.md")
	commitFiles(t, repo.MainRoot, "docs: main moves on", "design/note.md")

	res, err := Sync(syncOptions(repo, story))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || res.Stopped != "" || res.Base != "main" || res.Branch != "story/S-0001" || res.Worktree != ".flai-cache/worktrees/S-0001" {
		t.Fatalf("result: %+v", res)
	}
	if strings.Join(res.Outside, ",") != "README.md" || len(res.Branches) != 0 || res.TrialMergeSkipped != "" || len(res.Conflicts) != 0 || len(res.Regenerated) != 0 {
		t.Errorf("checks: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(wt, "design", "note.md")); err != nil {
		t.Errorf("main's commit is not under the branch: %v", err)
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"synced":true`, `"uncommitted":[]`, `"conflicts":[]`, `"branches":[]`, `"outside_touches":["README.md"]`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("json lacks %s: %s", want, data)
		}
	}
}

// ADR-0069: uncommitted changes refuse the sync, touching nothing, and name
// each path and the command to run again.
func TestSyncRefusesUncommittedChanges(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, repo.MainRoot, "docs: main moves on", "design/note.md")
	head, _ := (execx.System{}).Run(wt, "git", "rev-parse", "HEAD")
	write(t, wt, "docs/guide.md", "half done\n")

	o := syncOptions(repo, story)
	o.Again = "flai task done T-0001"
	res, err := Sync(o)
	if err != nil {
		t.Fatal(err)
	}
	if res.Synced || res.Stopped != StopUncommitted || res.RebaseInProgress() || strings.Join(res.Uncommitted, ",") != "docs/guide.md" || res.Abort != "" {
		t.Fatalf("result: %+v", res)
	}
	if want := "commit them on story/S-0001 (or stash them), then run flai task done T-0001 again"; res.Continue != want {
		t.Errorf("continue: %q, want %q", res.Continue, want)
	}
	if now, _ := (execx.System{}).Run(wt, "git", "rev-parse", "HEAD"); now != head {
		t.Errorf("HEAD moved: %s, was %s", now, head)
	}
}

// A rebase stopped on a conflict waits in the worktree with each path and
// how to continue or abort; a sync while it waits is refused the same way,
// touching nothing.
func TestSyncStopsOnConflictsAndRefusesARebaseInProgress(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] branch's guide", "docs/guide.md")
	commitFiles(t, repo.MainRoot, "docs: main's guide", "docs/guide.md")

	res, err := Sync(syncOptions(repo, story))
	if err != nil {
		t.Fatal(err)
	}
	if res.Synced || res.Stopped != StopConflicts || !res.RebaseInProgress() || strings.Join(res.Conflicts, ",") != "docs/guide.md" || len(res.Uncommitted) != 0 {
		t.Fatalf("result: %+v", res)
	}
	if want := "in .flai-cache/worktrees/S-0001, resolve each conflicting path, git add it, and run git rebase --continue; then run flai stream sync S-0001 again"; res.Continue != want {
		t.Errorf("continue: %q", res.Continue)
	}
	if want := "in .flai-cache/worktrees/S-0001, run git rebase --abort, which puts story/S-0001 back as it was before the sync"; res.Abort != want {
		t.Errorf("abort: %q", res.Abort)
	}
	if !RebaseInProgress(execx.System{}, wt) {
		t.Fatal("the rebase is not left waiting")
	}

	res, err = Sync(syncOptions(repo, story))
	if err != nil {
		t.Fatal(err)
	}
	if res.Synced || res.Stopped != StopRebaseInProgress || strings.Join(res.Conflicts, ",") != "docs/guide.md" || res.Abort == "" {
		t.Fatalf("in progress: %+v", res)
	}
	if !RebaseInProgress(execx.System{}, wt) || strings.Join(Conflicts(execx.System{}, wt), ",") != "docs/guide.md" {
		t.Error("the refused sync touched the rebase")
	}
}

// ADR-0098: a stop on generated files alone is resolved by writing them
// again and continued, for each commit that stops on them.
func TestSyncRegeneratesGeneratedFilesItStopsOnAlone(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] one", "docs/gen.md")
	commitFiles(t, wt, "feat: [S-0001] two", "docs/gen.md", "docs/two.md")
	commitFiles(t, repo.MainRoot, "docs: main's gen", "docs/gen.md")
	gen := GeneratedFile{Path: "docs/gen.md", Regenerate: func() error {
		return os.WriteFile(filepath.Join(wt, "docs", "gen.md"), []byte("regenerated\n"), 0o644)
	}}

	res, err := Sync(syncOptions(repo, story, gen))
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || strings.Join(res.Regenerated, ",") != "docs/gen.md,docs/gen.md" {
		t.Fatalf("result: %+v", res)
	}
	if got, _ := os.ReadFile(filepath.Join(wt, "docs", "gen.md")); string(got) != "regenerated\n" {
		t.Errorf("docs/gen.md: %q", got)
	}
}

// Another open story's branch that changes the same path conflicts in the
// trial merge, which writes nothing to either worktree and opens one thread
// for the pair; once the two merge cleanly the next sync resolves it.
func TestSyncTrialMergeConflictOpensAThreadForThePair(t *testing.T) {
	repo := syncProject(t)
	story, one := openStory(t, repo, "S-0001", "docs")
	_, two := openStory(t, repo, "S-0002", "docs")
	commitFiles(t, one, "feat: [S-0001] guide", "docs/guide.md")
	commitFiles(t, two, "feat: [S-0002] guide", "docs/guide.md")

	res, err := Sync(syncOptions(repo, story))
	if err != nil {
		t.Fatal(err)
	}
	if res.TrialMergeSkipped != "" {
		t.Skip(res.TrialMergeSkipped)
	}
	if len(res.Branches) != 1 {
		t.Fatalf("branches: %+v", res.Branches)
	}
	b := res.Branches[0]
	if b.Story != "S-0002" || b.Branch != "story/S-0002" || b.Status != workitem.InProgress || b.Clean || strings.Join(b.Conflicts, ",") != "docs/guide.md" || b.Thread != "TH-0001" {
		t.Errorf("pair: %+v", b)
	}
	for _, wt := range []string{one, two} {
		if dirty, _ := Uncommitted(execx.System{}, wt); len(dirty) != 0 {
			t.Errorf("%s is not clean: %v", wt, dirty)
		}
	}
	th, err := threads.Get(repo, "TH-0001")
	if err != nil {
		t.Fatal(err)
	}
	if th.Title != ConflictTitle("S-0002", "S-0001") || !th.Open() || th.Opener() != ConflictAuthor || th.Entries()[0].Text != ConflictText("S-0001", "S-0002", []string{"docs/guide.md"}) {
		t.Errorf("thread: %+v %+v", th, th.Entries())
	}

	// S-0002 takes S-0001's line, so the two merge cleanly
	write(t, two, "docs/guide.md", "feat: [S-0001] guide\n")
	git(t, two, "commit", "-q", "-am", "feat: [S-0002] take S-0001's guide")
	if res, err = Sync(syncOptions(repo, story)); err != nil || len(res.Branches) != 1 || !res.Branches[0].Clean || res.Branches[0].Thread != "" {
		t.Fatalf("clean: %v %+v", err, res.Branches)
	}
	if th, err = threads.Get(repo, "TH-0001"); err != nil || th.Open() {
		t.Errorf("the pair's thread is still open: %v", err)
	}
}
