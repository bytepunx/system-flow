package storygit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// commitFiles writes each file and commits them with the subject, leaving
// out the cache folders Committed keeps its answer in, as a project's
// .gitignore does.
func commitFiles(t *testing.T, dir, subject string, files ...string) {
	t.Helper()
	exclude, err := (execx.System{}).Run(dir, "git", "rev-parse", "--path-format=absolute", "--git-path", "info/exclude")
	if err != nil {
		t.Fatal(exclude, err)
	}
	if err := os.MkdirAll(filepath.Dir(exclude), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(exclude, []byte(".flai-cache\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(f)), 0o755); err != nil {
			t.Fatal(err)
		}
		write(t, dir, f, subject+"\n")
	}
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", subject)
}

// S-0205: a story's commits are those on the main branch and the story
// branches whose subject names it in brackets, merges left out, and their
// files under wip left out.
func TestCommittedReadsEachStorysFilesFromTheMainAndStoryBranches(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	repo := &workitem.Repo{Root: dir, MainRoot: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}

	commitFiles(t, dir, "feat: [S-0001] one", "cli/a.go", "wip/kanban/stories/S-0001.md")
	commitFiles(t, dir, "fix: [S-1] and [S-0002] both", "cli/b.go")
	commitFiles(t, dir, "docs: no story", "README.md")
	commitFiles(t, dir, "chore: [T-0003] a task, [S-0002 unclosed", "cli/t.go")
	git(t, dir, "mv", "cli/a.go", "cli/z.go")
	git(t, dir, "commit", "-q", "-m", "refactor: [S-0007] rename")

	// a story branch not merged, and one merged with a merge commit that
	// names a story of its own
	git(t, dir, "checkout", "-q", "-b", "story/S-0003")
	commitFiles(t, dir, "feat: [S-0003] on its branch", "web/d.go")
	git(t, dir, "checkout", "-q", "-b", "story/S-0004", "main")
	commitFiles(t, dir, "feat: [S-0004] merged", "web/e.go")
	git(t, dir, "checkout", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "merge: [S-0005] the branch", "story/S-0004")
	git(t, dir, "branch", "-q", "-D", "story/S-0004")
	// a branch that is neither
	git(t, dir, "checkout", "-q", "-b", "other/x")
	commitFiles(t, dir, "feat: [S-0006] elsewhere", "web/f.go")
	git(t, dir, "checkout", "-q", "main")

	got, err := Committed(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		"S-0001": {"cli/a.go", "cli/b.go"},
		"S-0002": {"cli/b.go"},
		"S-0003": {"web/d.go"},
		"S-0004": {"web/e.go"},
		"S-0007": {"cli/a.go", "cli/z.go"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("committed = %v, want %v", got, want)
	}
}

// A project in a folder of a larger repository reads only its own files,
// named from its root.
func TestCommittedKeepsToTheProjectsFolder(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "feat: [S-0001] both", "proj/cli/a.go", "proj/wip/x.md", "outside.go")
	sub := filepath.Join(dir, "proj")
	repo := &workitem.Repo{Root: sub, MainRoot: sub, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip/"}}}
	got, err := Committed(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string][]string{"S-0001": {"cli/a.go"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("committed = %v, want %v", got, want)
	}
}

// logCounter is the real runner, counting the git log processes it starts.
type logCounter struct {
	execx.System
	logs int
}

func (c *logCounter) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && slices.Contains(args, "log") {
		c.logs++
	}
	return c.System.Run(dir, name, args...)
}

// committedRepo is a repository with a commit of story S-0001 on main and a
// cache folder, and the workitem.Repo that reads it.
func committedRepo(t *testing.T) (string, *workitem.Repo) {
	t.Helper()
	dir := gitRepo(t)
	commitFiles(t, dir, "feat: [S-0001] one", "cli/a.go")
	repo := &workitem.Repo{Root: dir, MainRoot: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}
	if err := os.MkdirAll(repo.CacheDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir, repo
}

// committedOnce calls Committed and reports whether it read the log.
func committedOnce(t *testing.T, r *logCounter, repo *workitem.Repo) (map[string][]string, bool) {
	t.Helper()
	before := r.logs
	got, err := Committed(r, repo)
	if err != nil {
		t.Fatal(err)
	}
	return got, r.logs > before
}

// I-0184: while no branch tip moved, Committed answers what it read before
// and starts no git log.
func TestCommittedReadsTheLogOnceWhileNoBranchMoves(t *testing.T) {
	_, repo := committedRepo(t)
	r := &logCounter{}
	first, read := committedOnce(t, r, repo)
	if !read {
		t.Fatal("the first call read no log")
	}
	again, read := committedOnce(t, r, repo)
	if read {
		t.Error("a second call with no branch moved read the log again")
	}
	if want := map[string][]string{"S-0001": {"cli/a.go"}}; !reflect.DeepEqual(first, want) || !reflect.DeepEqual(again, want) {
		t.Errorf("committed = %v then %v, want %v both times", first, again, want)
	}
}

// A commit on main or on a story branch, and a new story branch, are seen by
// the next call, which reads the log again; the call after it, with nothing
// moved, does not.
func TestCommittedReadsTheLogAgainWhenABranchMoves(t *testing.T) {
	dir, repo := committedRepo(t)
	git(t, dir, "branch", "story/S-0002")
	r := &logCounter{}
	committedOnce(t, r, repo)
	seen := func(after, id string, want []string) {
		t.Helper()
		got, read := committedOnce(t, r, repo)
		if !read || !reflect.DeepEqual(got[id], want) {
			t.Errorf("after %s: read the log %v, %s = %v, want %v", after, read, id, got[id], want)
		}
		if _, read := committedOnce(t, r, repo); read {
			t.Errorf("after %s: the call after the one that saw it read the log again", after)
		}
	}

	commitFiles(t, dir, "feat: [S-0001] two", "cli/b.go")
	seen("a commit on main", "S-0001", []string{"cli/a.go", "cli/b.go"})

	git(t, dir, "checkout", "-q", "story/S-0002")
	commitFiles(t, dir, "feat: [S-0002] on its branch", "web/c.go")
	git(t, dir, "checkout", "-q", "main")
	seen("a commit on a story branch", "S-0002", []string{"web/c.go"})

	git(t, dir, "checkout", "-q", "-b", "story/S-0003")
	commitFiles(t, dir, "feat: [S-0003] a new branch", "web/d.go")
	git(t, dir, "checkout", "-q", "main")
	seen("a new story branch", "S-0003", []string{"web/d.go"})
}

// A caller who changes Committed's answer does not change the next one, which
// is answered from what was kept.
func TestCommittedAnswersACopy(t *testing.T) {
	_, repo := committedRepo(t)
	r := &logCounter{}
	got, _ := committedOnce(t, r, repo)
	got["S-0001"][0] = "changed.go"
	got["S-0009"] = []string{"added.go"}
	again, read := committedOnce(t, r, repo)
	if read {
		t.Error("the second call read the log again rather than answer what was kept")
	}
	if want := map[string][]string{"S-0001": {"cli/a.go"}}; !reflect.DeepEqual(again, want) {
		t.Errorf("committed after the caller changed its answer = %v, want %v", again, want)
	}
}

// Each repository root answers its own commits.
func TestCommittedKeepsEachRootApart(t *testing.T) {
	_, one := committedRepo(t)
	dir := gitRepo(t)
	commitFiles(t, dir, "feat: [S-0002] other", "web/b.go")
	two := &workitem.Repo{Root: dir, MainRoot: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}
	r := &logCounter{}
	for range 2 {
		if got, _ := committedOnce(t, r, one); !reflect.DeepEqual(got, map[string][]string{"S-0001": {"cli/a.go"}}) {
			t.Errorf("first root = %v", got)
		}
		if got, _ := committedOnce(t, r, two); !reflect.DeepEqual(got, map[string][]string{"S-0002": {"web/b.go"}}) {
			t.Errorf("second root = %v", got)
		}
	}
	if r.logs != 2 {
		t.Errorf("read the log %d times for two roots called twice each, want once per root", r.logs)
	}
}

// forgetKept drops what Committed keeps in memory for repo, as a new process
// starts with none.
func forgetKept(repo *workitem.Repo) {
	keptCommitted.Lock()
	defer keptCommitted.Unlock()
	delete(keptCommitted.roots, repo.MainRoot)
}

// freshOnce calls Committed as a new process's first call, with nothing kept
// in memory, and reports whether it read the log.
func freshOnce(t *testing.T, r *logCounter, repo *workitem.Repo) (map[string][]string, bool) {
	t.Helper()
	forgetKept(repo)
	return committedOnce(t, r, repo)
}

// I-0209: a new process with no branch moved answers from the file the last
// read was kept in, and starts no git log.
func TestCommittedAnswersANewProcessFromTheKeptFile(t *testing.T) {
	_, repo := committedRepo(t)
	r := &logCounter{}
	first, read := committedOnce(t, r, repo)
	if !read {
		t.Fatal("the first call read no log")
	}
	if _, err := os.Stat(CommittedPath(repo.CacheDir())); err != nil {
		t.Fatalf("the first call kept no file: %v", err)
	}
	again, read := freshOnce(t, r, repo)
	if read {
		t.Error("a new process with no branch moved read the log rather than the kept file")
	}
	if want := map[string][]string{"S-0001": {"cli/a.go"}}; !reflect.DeepEqual(first, want) || !reflect.DeepEqual(again, want) {
		t.Errorf("committed = %v then %v, want %v both times", first, again, want)
	}
	again["S-0001"][0] = "changed.go"
	if got, _ := committedOnce(t, r, repo); !reflect.DeepEqual(got, map[string][]string{"S-0001": {"cli/a.go"}}) {
		t.Errorf("committed after the caller changed the answer read from the file = %v", got)
	}
}

// After a branch tip moves, a new process reads the log again and keeps what
// it read, so that the process after it reads none.
func TestCommittedRewritesTheKeptFileWhenABranchMoves(t *testing.T) {
	dir, repo := committedRepo(t)
	r := &logCounter{}
	committedOnce(t, r, repo)
	commitFiles(t, dir, "feat: [S-0001] two", "cli/b.go")
	want := map[string][]string{"S-0001": {"cli/a.go", "cli/b.go"}}
	if got, read := freshOnce(t, r, repo); !read || !reflect.DeepEqual(got, want) {
		t.Errorf("after a commit on main, a new process read the log %v and answered %v, want a read and %v", read, got, want)
	}
	if got, read := freshOnce(t, r, repo); read || !reflect.DeepEqual(got, want) {
		t.Errorf("the next new process read the log %v and answered %v, want no read and %v", read, got, want)
	}
}

// A kept file that is damaged, or was read in another root or at another
// key, is no answer: a new process reads the log, answers what it read, and
// keeps it in place of the file.
func TestCommittedReadsTheLogOverAFileThatIsNoAnswer(t *testing.T) {
	_, repo := committedRepo(t)
	r := &logCounter{}
	committedOnce(t, r, repo)
	path := CommittedPath(repo.CacheDir())
	good, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var kept keptCommittedFile
	if err := json.Unmarshal(good, &kept); err != nil {
		t.Fatal(err)
	}
	changed := func(change func(*keptCommittedFile)) string {
		k := kept
		k.Committed = map[string][]string{"S-0009": {"wrong.go"}}
		change(&k)
		data, err := json.Marshal(k)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	want := map[string][]string{"S-0001": {"cli/a.go"}}
	for name, content := range map[string]string{
		"garbage":       "\x00not json{",
		"another shape": `["S-0001"]`,
		"no answer":     changed(func(k *keptCommittedFile) { k.Committed = nil }),
		"another key":   changed(func(k *keptCommittedFile) { k.Key += "x" }),
		"another root":  changed(func(k *keptCommittedFile) { k.Root += "x" }),
	} {
		write(t, repo.CacheDir(), filepath.Base(path), content)
		if got, read := freshOnce(t, r, repo); !read || !reflect.DeepEqual(got, want) {
			t.Errorf("%s: a new process read the log %v and answered %v, want a read and %v", name, read, got, want)
		}
		if got, read := freshOnce(t, r, repo); read || !reflect.DeepEqual(got, want) {
			t.Errorf("%s: the next new process read the log %v and answered %v, want the rewritten file's %v", name, read, got, want)
		}
	}
}

// A cache folder that cannot be written costs each new process a read of
// the log, never the answer.
func TestCommittedAnswersWhenTheCacheCannotBeWritten(t *testing.T) {
	dir, repo := committedRepo(t)
	if err := os.Remove(repo.CacheDir()); err != nil {
		t.Fatal(err)
	}
	write(t, dir, ".flai-cache", "a file where the cache folder goes\n")
	r := &logCounter{}
	want := map[string][]string{"S-0001": {"cli/a.go"}}
	for i := range 2 {
		got, read := freshOnce(t, r, repo)
		if !read || !reflect.DeepEqual(got, want) {
			t.Errorf("new process %d: read the log %v and answered %v, want a read and %v", i+1, read, got, want)
		}
	}
}

// A project with no cache folder gains none, since flai check reports a
// .flai-cache that no .gitignore ignores (layout.gitignore): each new process
// reads the log, as before I-0209.
func TestCommittedMakesNoCacheFolder(t *testing.T) {
	_, repo := committedRepo(t)
	if err := os.Remove(repo.CacheDir()); err != nil {
		t.Fatal(err)
	}
	r := &logCounter{}
	want := map[string][]string{"S-0001": {"cli/a.go"}}
	for i := range 2 {
		if got, read := freshOnce(t, r, repo); !read || !reflect.DeepEqual(got, want) {
			t.Errorf("new process %d: read the log %v and answered %v, want a read and %v", i+1, read, got, want)
		}
	}
	if _, err := os.Stat(repo.CacheDir()); !os.IsNotExist(err) {
		t.Errorf("Committed made the cache folder %s (%v), want none", repo.CacheDir(), err)
	}
}

func TestCommittedFailsWithoutARepository(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	dir := t.TempDir()
	repo := &workitem.Repo{Root: dir, Manifest: manifest.Manifest{Layout: map[string]string{"wip": "wip"}}}
	got, err := Committed(execx.System{}, repo)
	if err == nil || got != nil || !strings.Contains(err.Error(), "cannot read the stories' commits in "+dir) {
		t.Errorf("not a repository: %v %v", got, err)
	}
}

func TestSubjectStories(t *testing.T) {
	for subject, want := range map[string][]string{
		"feat: [S-12] x":                {"S-0012"},
		"[S-0001][S-002] [S-1]":         {"S-0001", "S-0002"},
		"[E-0001] [T-0002] [s-3] [S-x]": nil,
		"S-0004 without brackets [":     nil,
	} {
		if got := subjectStories(subject); !reflect.DeepEqual(got, want) {
			t.Errorf("%q: %v, want %v", subject, got, want)
		}
	}
}
