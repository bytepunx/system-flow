package storygit

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// sharedRepo is a repository at dir whose manifest lays wip and design out
// under their own names and shares docs/users/flai.md and docs/adrs.
func sharedRepo(dir string) *workitem.Repo {
	return &workitem.Repo{Root: dir, MainRoot: dir, Manifest: manifest.Manifest{
		Layout: map[string]string{"wip": "wip/", "design": "design"},
		Claims: manifest.Claims{Shared: []string{"docs/users/flai.md", "docs/adrs"}},
	}}
}

// ADR-0248: the main checkout's uncommitted files are split by the shared
// set: a file under wip or design, an untracked one in a new folder under
// design among them, at a shared path, or under a shared folder, is shared;
// one anywhere else, the operator's, is the rest. Each file of an untracked
// folder is named, in git's order: the tracked files changed, then the
// untracked.
func TestDirtyMainSplitsTheFilesByTheSharedSet(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "base", "cli/a.go", "docs/users/flai.md", "docs/guide.md")
	write(t, dir, "cli/a.go", "the operator's edit\n")
	write(t, dir, "docs/users/flai.md", "flai's guide\n")
	write(t, dir, "docs/guide.md", "the operator's guide\n")
	writeFile(t, dir, "wip/agents/S-0001.md")
	writeFile(t, dir, "design/analysis/2026-10-19-all.md")
	writeFile(t, dir, "design/issues/I-0150-x.md")
	writeFile(t, dir, "docs/adrs/0001-x.md")
	writeFile(t, dir, "notes.md")

	got, err := DirtyMain(r, sharedRepo(dir))
	if err != nil {
		t.Fatal(err)
	}
	want := Dirty{
		Shared: []string{"docs/users/flai.md", "design/analysis/2026-10-19-all.md", "design/issues/I-0150-x.md", "docs/adrs/0001-x.md", "wip/agents/S-0001.md"},
		Rest:   []string{"cli/a.go", "docs/guide.md", "notes.md"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DirtyMain = %+v, want %+v", got, want)
	}
}

// InSharedSet reads the layout's folders with or without a trailing slash,
// a shared path as the file itself and everything below it, a shared glob
// as the files it matches, and nothing else; a manifest with no layout
// shares nothing.
func TestInSharedSetReadsTheLayoutAndTheSharedPatterns(t *testing.T) {
	m := sharedRepo("").Manifest
	m.Claims.Shared = append(m.Claims.Shared, "template/*.md")
	for p, want := range map[string]bool{
		"wip/kanban/stories/S-0001-x.md": true,
		"design/issues/summary.md":       true,
		"docs/users/flai.md":             true,
		"docs/adrs/0001-x.md":            true,
		"template/CHANGELOG.md":          true,
		"template/root/README.md":        false,
		"docs/users/flai-reference.md":   false,
		"designs/x.md":                   false,
		"wip.md":                         false,
		"cli/a.go":                       false,
	} {
		if got := InSharedSet(m, p); got != want {
			t.Errorf("InSharedSet(%s) = %v, want %v", p, got, want)
		}
	}
	if InSharedSet(manifest.Manifest{}, "wip/x.md") {
		t.Error("a manifest with no layout shares nothing")
	}
}

// S-0402: of the main checkout's uncommitted files, only those the change
// also changes are named: a dirty file it changes, an untracked file it adds,
// and not a dirty or untracked file it leaves alone, as the analyzer's report
// and issues on TH-0462.
func TestChangedAmongNamesOnlyTheDirtyFilesTheChangeChanges(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "base", "cli/a.go", "cli/b.go")
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	write(t, dir, "cli/a.go", "changed\n")
	commitFiles(t, dir, "feat: [S-0001] change a, add c", "cli/c.go")
	git(t, dir, "checkout", "-q", "main")

	write(t, dir, "cli/a.go", "the operator's edit\n")
	write(t, dir, "cli/b.go", "left alone by the change\n")
	writeFile(t, dir, "cli/c.go")
	writeFile(t, dir, "design/analysis/2026-10-19-all.md")
	writeFile(t, dir, "design/issues/I-0150-x.md")

	d, err := DirtyMain(r, sharedRepo(dir))
	if err != nil {
		t.Fatal(err)
	}
	got, err := ChangedAmong(r, dir, "main", "story/S-0001", d.Rest)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"cli/a.go", "cli/c.go"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ChangedAmong = %v, want %v", got, want)
	}
}

// A clean checkout, or one whose dirty files the change leaves alone, has
// nothing overwritten; and a branch made before main moved on is read from
// the commit the two share, so main's own later changes are not the
// branch's.
func TestChangedAmongIsEmptyWhenTheChangeLeavesTheDirtyFilesAlone(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, dir, "feat: [S-0001] add x", "cli/x.go")
	git(t, dir, "checkout", "-q", "main")
	commitFiles(t, dir, "main moves on", "cli/m.go")

	overwritten := func() []string {
		t.Helper()
		d, err := DirtyMain(r, sharedRepo(dir))
		if err != nil {
			t.Fatal(err)
		}
		got, err := ChangedAmong(r, dir, "main", "story/S-0001", d.Rest)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	if got := overwritten(); len(got) != 0 {
		t.Fatalf("clean checkout: overwritten = %v; want none", got)
	}
	write(t, dir, "cli/m.go", "the operator's edit of main's own file\n")
	writeFile(t, dir, "notes.md")
	if got := overwritten(); len(got) != 0 {
		t.Fatalf("dirty elsewhere: overwritten = %v; want none", got)
	}
}

// processCounter is the real runner, counting the processes it starts.
type processCounter struct {
	execx.System
	calls []string
}

func (c *processCounter) Run(dir, name string, args ...string) (string, error) {
	c.calls = append(c.calls, name+" "+strings.Join(args, " "))
	return c.System.Run(dir, name, args...)
}

// S-0698: ReadHeads reads every local branch's head and the branch checked
// out with one git process, each head the commit rev-parse's ^{commit}
// names; with HEAD detached no branch is checked out, and Main says why. Integration: it
// runs real git.
func TestReadHeadsReadsEveryBranchAndTheCheckedOutOneInOneProcess(t *testing.T) {
	dir := gitRepo(t)
	git(t, dir, "branch", Branch("S-0001"))
	git(t, dir, "checkout", "-q", "-b", Branch("S-0002"))
	write(t, dir, "b.md", "two\n")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "two")
	git(t, dir, "checkout", "-q", "main")

	r := &processCounter{}
	heads, err := ReadHeads(r, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.calls) != 1 || !strings.HasPrefix(r.calls[0], "git for-each-ref ") {
		t.Errorf("ReadHeads ran %q, want one git for-each-ref", r.calls)
	}
	want := map[string]string{}
	for _, b := range []string{"main", Branch("S-0001"), Branch("S-0002")} {
		head, err := (execx.System{}).Run(dir, "git", "rev-parse", "--verify", "--quiet", b+"^{commit}")
		if err != nil {
			t.Fatal(head, err)
		}
		want[b] = head
	}
	if !reflect.DeepEqual(heads.Commits, want) {
		t.Errorf("heads %v, want rev-parse's %v", heads.Commits, want)
	}
	if branch, head, err := heads.Main(); err != nil || branch != "main" || head != want["main"] {
		t.Errorf("Main = %q %q %v, want main at %s", branch, head, err, want["main"])
	}

	git(t, dir, "checkout", "-q", "--detach", "main")
	heads, err = ReadHeads(execx.System{}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if heads.Checked != "" || heads.Commits["main"] != want["main"] {
		t.Errorf("detached: checked %q, main at %q, want none checked and main at %s", heads.Checked, heads.Commits["main"], want["main"])
	}
	if _, _, err := heads.Main(); err == nil || !strings.Contains(err.Error(), "detached HEAD") {
		t.Errorf("detached: Main error %v, want one saying HEAD is detached", err)
	}
}
