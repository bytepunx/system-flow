package template

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// cachedClone resolves a git source whose clone is already in the cache.
func cachedClone(t *testing.T, ref string) Source {
	t.Helper()
	s, err := Resolve("https://example.com/t.git", ref, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.Dir, ManifestFile), []byte("version: 1.0.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestEnsureFetchesACachedBranchClone(t *testing.T) {
	for _, ref := range []string{"main", ""} {
		f := &fakeGit{out: map[string]string{"rev-parse": "main"}}
		if err := cachedClone(t, ref).Ensure(f, false); err != nil {
			t.Fatalf("ref %q: %v", ref, err)
		}
		if !f.ran("fetch") || !f.ran("reset") || f.ran("clone") {
			t.Fatalf("ref %q: want fetch and reset, no clone; ran %v", ref, f.calls)
		}
	}
}

func TestEnsureReusesACachedTagOrCommitClone(t *testing.T) {
	f := &fakeGit{out: map[string]string{"rev-parse": "HEAD"}}
	if err := cachedClone(t, "v1.0.18").Ensure(f, false); err != nil {
		t.Fatal(err)
	}
	if f.ran("fetch") || f.ran("reset") || f.ran("clone") {
		t.Fatalf("a detached clone is not fetched again; ran %v", f.calls)
	}
}

func TestEnsureKeepsABranchCloneWhenTheFetchFails(t *testing.T) {
	f := &fakeGit{
		out:  map[string]string{"rev-parse": "main"},
		fail: map[string]error{"fetch": errors.New("could not resolve host")},
	}
	s := cachedClone(t, "main")
	err := s.Ensure(f, false)
	if !errors.Is(err, ErrStale) {
		t.Fatalf("want ErrStale, got %v", err)
	}
	if f.ran("reset") || !s.Cached() {
		t.Fatalf("the clone should be kept as it was; ran %v", f.calls)
	}
}

// recorder runs real git and records the subcommands it ran.
type recorder struct {
	execx.System
	subs []string
}

func (r *recorder) Run(dir, name string, args ...string) (string, error) {
	if len(args) > 0 {
		r.subs = append(r.subs, args[0])
	}
	return r.System.Run(dir, name, args...)
}

func TestEnsureSeesNewCommitsOnABranchButNotOnATagWithRealGit(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: needs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	run := func(dir string, args ...string) {
		t.Helper()
		if out, err := (execx.System{}).Run(dir, "git", args...); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	write := func(dir, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, RootDir, "a.txt"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	read := func(s Source) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(s.Dir, RootDir, "a.txt"))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	remote := filepath.Join(t.TempDir(), "t.git")
	run("", "init", "-q", "--bare", "-b", "main", remote)
	work := t.TempDir()
	run(work, "init", "-q", "-b", "main")
	run(work, "config", "user.email", "t@t")
	run(work, "config", "user.name", "t")
	_ = os.WriteFile(filepath.Join(work, ManifestFile), []byte("version: 1.0.0\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(work, RootDir), 0o755)
	write(work, "one\n")
	run(work, "add", "-A")
	run(work, "commit", "-q", "-m", "one")
	run(work, "tag", "v1.0.0")
	run(work, "remote", "add", "origin", "file://"+remote)
	run(work, "push", "-q", "origin", "main", "v1.0.0")

	cache := t.TempDir()
	branch, _ := Resolve("file://"+remote, "main", cache)
	deflt, _ := Resolve("file://"+remote, "", cache)
	tag, _ := Resolve("file://"+remote, "v1.0.0", cache)
	for _, s := range []Source{branch, deflt, tag} {
		if err := s.Ensure(execx.System{}, false); err != nil {
			t.Fatalf("clone %q: %v", s.Ref, err)
		}
	}

	write(work, "two\n")
	run(work, "commit", "-q", "-am", "two")
	run(work, "push", "-q", "origin", "main")

	for _, s := range []Source{branch, deflt} {
		if err := s.Ensure(execx.System{}, false); err != nil {
			t.Fatalf("update %q: %v", s.Ref, err)
		}
		if got := read(s); got != "two\n" {
			t.Fatalf("clone of %q reads %q after a push, want two", s.Ref, got)
		}
	}
	rec := &recorder{}
	if err := tag.Ensure(rec, false); err != nil {
		t.Fatal(err)
	}
	if got := read(tag); got != "one\n" {
		t.Fatalf("clone of the tag reads %q, want one", got)
	}
	for _, sub := range rec.subs {
		if sub == "fetch" || sub == "clone" {
			t.Fatalf("the tag clone was fetched again: ran %v", rec.subs)
		}
	}

	run(work, "tag", "-a", "-m", "release", "v1.0.10")
	run(work, "push", "-q", "origin", "v1.0.10")
	rm, err := ListRemote(execx.System{}, "file://"+remote)
	if err != nil {
		t.Fatal(err)
	}
	if latest, _ := rm.Latest(); rm.DefaultBranch != "main" || latest.Name != "v1.0.10" || len(rm.Tags) != 2 {
		t.Fatalf("real ls-remote: %+v", rm)
	}

	if err := os.RemoveAll(remote); err != nil {
		t.Fatal(err)
	}
	if err := branch.Ensure(execx.System{}, false); !errors.Is(err, ErrStale) {
		t.Fatalf("with the remote gone want ErrStale, got %v", err)
	}
	if got := read(branch); got != "two\n" {
		t.Fatalf("the stale clone should be kept, reads %q", got)
	}
}
