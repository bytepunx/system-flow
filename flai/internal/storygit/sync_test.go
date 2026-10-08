package storygit

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/messages"
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

// syncClock dates what Sync writes.
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
	return SyncOptions{Runner: execx.System{}, Repo: repo, Story: story, Now: syncClock, Files: Files{Generated: gen}}
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

// mergedFiles merges docs/m.md, at a stop, into its three versions in turn,
// read from the index of the worktree wt.
func mergedFiles(t *testing.T, wt string) MergedFiles {
	t.Helper()
	return MergedFiles{
		Match: func(p string) bool { return p == "docs/m.md" },
		Merge: func(p string) error {
			base, ours, theirs, err := Stages(execx.System{}, wt, p)
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(wt, filepath.FromSlash(p)), []byte(base+ours+theirs), 0o644)
		},
	}
}

// ADR-0126: a stop on a file flai merges and a generated file is continued:
// the file is merged from the stop's three versions, and then every generated
// file is written again, the one that did not conflict too, from the merge.
func TestSyncMergesTheFilesItMergesAndRegeneratesFromThem(t *testing.T) {
	repo := syncProject(t)
	commitFiles(t, repo.MainRoot, "docs: base", "docs/m.md")
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] ours", "docs/m.md")
	commitFiles(t, repo.MainRoot, "docs: main's", "docs/m.md")
	gen := GeneratedFile{Path: "docs/gen.md", Regenerate: func() error {
		m, err := os.ReadFile(filepath.Join(wt, "docs", "m.md"))
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(wt, "docs", "gen.md"), []byte("from "+string(m)), 0o644)
	}}
	o := syncOptions(repo, story, gen)
	o.Files.Merged = mergedFiles(t, wt)

	res, err := Sync(o)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || strings.Join(res.Merged, ",") != "docs/m.md" || strings.Join(res.Regenerated, ",") != "docs/gen.md" {
		t.Fatalf("result: %+v", res)
	}
	merged := "docs: base\ndocs: main's\nfeat: [S-0001] ours\n"
	if got, _ := (execx.System{}).Run(wt, "git", "show", "HEAD:docs/m.md"); got+"\n" != merged {
		t.Errorf("docs/m.md: %q", got)
	}
	if got, _ := (execx.System{}).Run(wt, "git", "show", "HEAD:docs/gen.md"); got+"\n" != "from "+merged {
		t.Errorf("docs/gen.md is not written from the merge: %q", got)
	}
	if st, _ := Uncommitted(execx.System{}, wt); len(st) != 0 {
		t.Errorf("the worktree is not clean: %v", st)
	}
	data, _ := json.Marshal(res)
	if !strings.Contains(string(data), `"merged":["docs/m.md"]`) || !strings.Contains(string(data), `"folded":[]`) {
		t.Errorf("json: %s", data)
	}
}

// A merge that fails leaves the stop to the agent, as a stop on any other
// path is.
func TestSyncLeavesAStopWhoseMergeFailsToTheAgent(t *testing.T) {
	repo := syncProject(t)
	commitFiles(t, repo.MainRoot, "docs: base", "docs/m.md")
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] ours", "docs/m.md")
	commitFiles(t, repo.MainRoot, "docs: main's", "docs/m.md")
	o := syncOptions(repo, story)
	o.Files.Merged = MergedFiles{Match: func(p string) bool { return p == "docs/m.md" }, Merge: func(string) error { return errors.New("not merged") }}

	res, err := Sync(o)
	if err != nil {
		t.Fatal(err)
	}
	if res.Synced || res.Stopped != StopConflicts || strings.Join(res.Conflicts, ",") != "docs/m.md" || len(res.Merged) != 0 {
		t.Fatalf("result: %+v", res)
	}
	if !RebaseInProgress(execx.System{}, wt) {
		t.Error("the rebase is not left waiting")
	}
}

// ADR-0126: after a clean rebase the fold runs with the base branch; what it
// folds is committed on the branch with the generated files written again,
// a deleted file included.
func TestSyncCommitsTheFoldAfterACleanRebase(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] dup", "docs/dup.md")
	commitFiles(t, repo.MainRoot, "docs: main's keep", "docs/keep.md")
	gen := GeneratedFile{Path: "docs/gen.md", Regenerate: func() error {
		return os.WriteFile(filepath.Join(wt, "docs", "gen.md"), []byte("after the fold\n"), 0o644)
	}}
	o := syncOptions(repo, story, gen)
	var gotBase string
	o.Files.Fold = func(base string) ([]Folded, []string, error) {
		gotBase = base
		if err := os.Remove(filepath.Join(wt, "docs", "dup.md")); err != nil {
			return nil, nil, err
		}
		if err := os.WriteFile(filepath.Join(wt, "docs", "keep.md"), []byte("kept both\n"), 0o644); err != nil {
			return nil, nil, err
		}
		return []Folded{{From: "I-0002", Into: "I-0001"}, {From: "I-0003", Into: "I-0001"}}, []string{"docs/keep.md", "docs/dup.md"}, nil
	}

	res, err := Sync(o)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || gotBase != "main" || len(res.Folded) != 2 || res.Folded[0] != (Folded{From: "I-0002", Into: "I-0001"}) {
		t.Fatalf("result: %+v, base %q", res, gotBase)
	}
	r := execx.System{}
	if got, _ := r.Run(wt, "git", "log", "-1", "--format=%s"); got != "docs: [S-0001] fold I-0002 into I-0001, I-0003 into I-0001" {
		t.Errorf("the fold's commit: %q", got)
	}
	if got, _ := r.Run(wt, "git", "show", "--name-status", "--format=", "HEAD"); got != "D\tdocs/dup.md\nA\tdocs/gen.md\nM\tdocs/keep.md" {
		t.Errorf("the fold's commit changed:\n%s", got)
	}
	if st, _ := Uncommitted(r, wt); len(st) != 0 {
		t.Errorf("the worktree is not clean: %v", st)
	}
	data, _ := json.Marshal(res)
	if !strings.Contains(string(data), `"folded":[{"from":"I-0002","into":"I-0001"},{"from":"I-0003","into":"I-0001"}]`) {
		t.Errorf("json: %s", data)
	}
}

// A fold that fails is logged and the sync stands, with nothing committed;
// one that folds nothing commits nothing.
func TestSyncStandsWhenTheFoldFailsOrFoldsNothing(t *testing.T) {
	repo := syncProject(t)
	story, wt := openStory(t, repo, "S-0001", "docs")
	commitFiles(t, wt, "feat: [S-0001] guide", "docs/guide.md")
	r := execx.System{}
	head, _ := r.Run(wt, "git", "rev-parse", "HEAD")
	for _, fold := range []FoldFunc{
		func(string) ([]Folded, []string, error) { return nil, nil, errors.New("unreadable") },
		func(string) ([]Folded, []string, error) { return []Folded{}, []string{}, nil },
	} {
		o := syncOptions(repo, story)
		o.Files.Fold = fold
		res, err := Sync(o)
		if err != nil || !res.Synced || len(res.Folded) != 0 {
			t.Fatalf("result: %+v %v", res, err)
		}
		if now, _ := r.Run(wt, "git", "rev-parse", "HEAD"); now != head {
			t.Errorf("a commit was made: %s, was %s", now, head)
		}
	}
}

// conflictingPair is a project with S-0001 and S-0002 in progress, whose
// branches change docs/guide.md each its own way. It returns S-0001 and the
// two worktrees.
func conflictingPair(t *testing.T) (repo *workitem.Repo, story *workitem.Item, one, two string) {
	t.Helper()
	repo = syncProject(t)
	story, one = openStory(t, repo, "S-0001", "docs")
	_, two = openStory(t, repo, "S-0002", "docs")
	commitFiles(t, one, "feat: [S-0001] guide", "docs/guide.md")
	commitFiles(t, two, "feat: [S-0002] guide", "docs/guide.md")
	return repo, story, one, two
}

// syncAt syncs story minutes after syncClock, so that each sync's entries
// carry a time of their own, and returns its one trial-merged branch.
func syncAt(t *testing.T, repo *workitem.Repo, story *workitem.Item, minutes int) BranchCheck {
	t.Helper()
	o := syncOptions(repo, story)
	o.Now = syncClock.Add(time.Duration(minutes) * time.Minute)
	res, err := Sync(o)
	if err != nil {
		t.Fatal(err)
	}
	if res.TrialMergeSkipped != "" {
		t.Skip(res.TrialMergeSkipped)
	}
	if len(res.Branches) != 1 {
		t.Fatalf("branches: %+v", res.Branches)
	}
	return res.Branches[0]
}

// conversation reads the conversation id.
func conversation(t *testing.T, repo *workitem.Repo, id string) *messages.Conversation {
	t.Helper()
	c, err := messages.Get(repo, id)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// noThreads fails the test when the project has any thread.
func noThreads(t *testing.T, repo *workitem.Repo) {
	t.Helper()
	if all, err := threads.List(repo); err != nil || len(all) != 0 {
		t.Errorf("threads: %v %+v", err, all)
	}
}

// ADR-0121: another open story's branch that changes the same path
// conflicts in the trial merge, which writes nothing to either worktree and
// tells the other story in the pair's conversation, about the path, which
// then awaits it. A sync that finds the same paths adds nothing, new paths
// add a message, and the sync that finds the two merging cleanly closes it.
func TestSyncTrialMergeConflictIsAMessageBetweenThePair(t *testing.T) {
	repo, story, one, two := conflictingPair(t)

	b := syncAt(t, repo, story, 1)
	if b.Story != "S-0002" || b.Branch != "story/S-0002" || b.Status != workitem.InProgress || b.Clean || strings.Join(b.Conflicts, ",") != "docs/guide.md" || b.Conversation != "MS-0001" || b.Thread != "" {
		t.Errorf("pair: %+v", b)
	}
	for _, wt := range []string{one, two} {
		if dirty, _ := Uncommitted(execx.System{}, wt); len(dirty) != 0 {
			t.Errorf("%s is not clean: %v", wt, dirty)
		}
	}
	noThreads(t, repo)
	c := conversation(t, repo, "MS-0001")
	e := c.Entries()
	if c.From != "S-0001" || c.To != "S-0002" || strings.Join(c.About, ",") != "docs/guide.md" || c.Status != messages.StatusOpen || c.Awaiting() != "S-0002" {
		t.Errorf("conversation: %+v", c)
	}
	if len(e) != 1 || e[0].Author != ConflictAuthor || e[0].Story != "S-0001" || e[0].Text != ConflictText("S-0002", "S-0001", []string{"docs/guide.md"}) {
		t.Fatalf("entries: %+v", e)
	}
	for _, want := range []string{"story/S-0001 and story/S-0002 conflict when merged.", "- `docs/guide.md`", "`flai message escalate`", "`message_escalate`"} {
		if !strings.Contains(e[0].Text, want) {
			t.Errorf("message lacks %q:\n%s", want, e[0].Text)
		}
	}

	// the same paths add nothing
	if b := syncAt(t, repo, story, 2); b.Conversation != "MS-0001" || len(conversation(t, repo, "MS-0001").Entries()) != 1 {
		t.Errorf("same paths: %+v %+v", b, conversation(t, repo, "MS-0001").Entries())
	}

	// other paths add a message to the same conversation
	commitFiles(t, one, "feat: [S-0001] more", "docs/more.md")
	commitFiles(t, two, "feat: [S-0002] more", "docs/more.md")
	if b := syncAt(t, repo, story, 3); b.Conversation != "MS-0001" || strings.Join(b.Conflicts, ",") != "docs/guide.md,docs/more.md" {
		t.Errorf("new paths: %+v", b)
	}
	c = conversation(t, repo, "MS-0001")
	if e := c.Entries(); len(e) != 2 || !strings.Contains(e[1].Text, "- `docs/guide.md`\n- `docs/more.md`") || strings.Join(c.About, ",") != "docs/guide.md,docs/more.md" {
		t.Fatalf("new paths: %+v %+v", c, e)
	}

	// S-0002 takes S-0001's lines, so the two merge cleanly
	write(t, two, "docs/guide.md", "feat: [S-0001] guide\n")
	write(t, two, "docs/more.md", "feat: [S-0001] more\n")
	git(t, two, "commit", "-q", "-am", "feat: [S-0002] take S-0001's lines")
	if b := syncAt(t, repo, story, 4); !b.Clean || b.Conversation != "" || b.Thread != "" {
		t.Errorf("clean: %+v", b)
	}
	c = conversation(t, repo, "MS-0001")
	if closed, why := c.Closed(repo); !closed || why != "story/S-0001 and story/S-0002 merge cleanly at the sync of S-0001" {
		t.Errorf("not closed on a clean merge: %v %q", closed, why)
	}
	noThreads(t, repo)
}

// A clean merge leaves open a conversation of the pair that tells of no
// conflict: one the agents started, or one whose newest message by flai is
// another notice.
func TestSyncCleanMergeLeavesAConversationWithoutAConflictOpen(t *testing.T) {
	repo := syncProject(t)
	story, _ := openStory(t, repo, "S-0001", "docs")
	openStory(t, repo, "S-0002", "design")
	agents, err := messages.Send(repo, messages.SendOptions{From: "S-0002", To: "S-0001", Author: "agent-S-0002", Text: "Which of us changes the guide?", Now: syncClock})
	if err != nil {
		t.Fatal(err)
	}
	if b := syncAt(t, repo, story, 1); !b.Clean || b.Conversation != "" {
		t.Errorf("pair: %+v", b)
	}
	if c := conversation(t, repo, agents.ID); c.Status != messages.StatusOpen || len(c.Entries()) != 1 {
		t.Errorf("the agents' conversation was touched: %+v %+v", c, c.Entries())
	}

	// a conflict message followed by another of flai's notices
	text := ConflictText("S-0001", "S-0002", []string{"docs/guide.md"})
	for i, m := range []string{text, "S-0001 renamed a path S-0002's claim covers."} {
		if _, err := messages.Notify(repo, messages.SendOptions{From: "S-0001", To: "S-0002", Author: ConflictAuthor, Text: m, Now: syncClock.Add(time.Duration(2+i) * time.Minute)}); err != nil {
			t.Fatal(err)
		}
	}
	syncAt(t, repo, story, 4)
	if c := conversation(t, repo, agents.ID); c.Status != messages.StatusOpen {
		t.Errorf("closed though flai's newest message tells of no conflict: %+v", c.Entries())
	}
}

// A pair's conversation that tells of a conflict is closed by the next sync
// once the other story is sent back out of in progress, though it has no
// branch to trial-merge then.
func TestSyncClosesAConflictConversationWhenTheOtherStoryIsNoLongerOpen(t *testing.T) {
	repo, story, _, _ := conflictingPair(t)
	if b := syncAt(t, repo, story, 1); b.Conversation != "MS-0001" {
		t.Fatalf("pair: %+v", b)
	}
	other, err := repo.Get("S-0002")
	if err != nil {
		t.Fatal(err)
	}
	other.Status = workitem.Ready
	if err := repo.Save(other); err != nil {
		t.Fatal(err)
	}

	o := syncOptions(repo, story)
	o.Now = syncClock.Add(2 * time.Minute)
	res, err := Sync(o)
	if err != nil || len(res.Branches) != 0 {
		t.Fatalf("sync: %v %+v", err, res)
	}
	c := conversation(t, repo, "MS-0001")
	if closed, why := c.Closed(repo); c.Status != messages.StatusClosed || !closed || why != "S-0002 is ready, no longer open, at the sync of S-0001" {
		t.Errorf("not closed: %v %q", closed, why)
	}
	noThreads(t, repo)
}

// ADR-0121: a conflict thread an older flai opened gets no entry while the
// pair conflicts, which is told in a conversation instead, and is resolved
// once the two merge cleanly, naming it on the branch; no thread is opened.
func TestSyncResolvesAnOldConflictThreadAndOpensNone(t *testing.T) {
	repo, story, _, two := conflictingPair(t)
	th, err := threads.New(repo, threads.NewOptions{Title: ConflictTitle("S-0001", "S-0002"), On: "S-0001", Author: ConflictAuthor, Text: "An old conflict thread.", Now: syncClock})
	if err != nil {
		t.Fatal(err)
	}

	if b := syncAt(t, repo, story, 1); b.Conversation != "MS-0001" || b.Thread != "" {
		t.Errorf("conflict: %+v", b)
	}
	if th, err = threads.Get(repo, th.ID); err != nil || !th.Open() || len(th.Entries()) != 1 {
		t.Fatalf("the old thread changed: %v %+v", err, th.Entries())
	}

	write(t, two, "docs/guide.md", "feat: [S-0001] guide\n")
	git(t, two, "commit", "-q", "-am", "feat: [S-0002] take S-0001's guide")
	if b := syncAt(t, repo, story, 2); !b.Clean || b.Thread != th.ID || b.Conversation != "" {
		t.Errorf("clean: %+v", b)
	}
	if th, err = threads.Get(repo, th.ID); err != nil || th.Open() {
		t.Errorf("the old thread is still open: %v", err)
	}
	if all, err := threads.List(repo); err != nil || len(all) != 1 {
		t.Errorf("threads: %v %+v", err, all)
	}
	if c := conversation(t, repo, "MS-0001"); c.Status != messages.StatusClosed {
		t.Errorf("the conversation is still open: %+v", c.Entries())
	}
}
