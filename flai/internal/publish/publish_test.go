package publish

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// remote creates a bare repository with one commit on main and returns its path.
func remote(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	r := execx.System{}
	bare := filepath.Join(t.TempDir(), "tpl.git")
	if _, err := r.Run("", "git", "init", "-q", "--bare", "-b", "main", bare); err != nil {
		t.Fatal(err)
	}
	seed := t.TempDir()
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}} {
		_, _ = r.Run(seed, "git", args...)
	}
	_ = os.WriteFile(filepath.Join(seed, "README.md"), []byte("old\n"), 0o644)
	_ = os.WriteFile(filepath.Join(seed, "stale.txt"), []byte("gone after push\n"), 0o644)
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "seed"}, {"push", "-q", bare, "main"}} {
		if out, err := r.Run(seed, "git", args...); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	return bare
}

// local copies the mini fixture into a temp dir with a publish section.
func local(t *testing.T, bare string) string {
	t.Helper()
	src := "../template/testdata/mini"
	dst := t.TempDir()
	_ = filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(path)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	y, _ := os.ReadFile(filepath.Join(dst, "template.yaml"))
	_ = os.WriteFile(filepath.Join(dst, "template.yaml"), []byte(string(y)+"publish:\n  repo: "+bare+"\n  ref: main\n"), 0o644)
	return dst
}

func TestPushLifecycle(t *testing.T) {
	bare := remote(t)
	dir := local(t, bare)
	r := execx.System{}
	cache := t.TempDir()
	opt := Options{Dir: dir, Remote: bare, Ref: "main", CacheDir: cache, Source: dir}

	dry := opt
	dry.DryRun = true
	res, err := Run(r, dry)
	if err != nil {
		t.Fatal(err)
	}
	if res.Pushed || res.Commit != "" || len(res.Files) < 3 || !strings.HasPrefix(res.Message, "template 9.9.9") {
		t.Errorf("dry run: %+v", res)
	}
	if _, err := os.Stat(filepath.Join(cache, "publish")); err == nil {
		if entries, _ := os.ReadDir(filepath.Join(cache, "publish")); len(entries) != 0 {
			t.Error("dry run left a working clone behind")
		}
	}
	res, err = Run(r, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Pushed || res.Commit == "" {
		t.Fatalf("push: %+v", res)
	}
	// verify the remote: stale file gone, template files present, message
	check := t.TempDir()
	if _, err := r.Run("", "git", "clone", "-q", bare, check); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(check, "stale.txt")); err == nil {
		t.Error("stale file should be removed by replace")
	}
	if _, err := os.Stat(filepath.Join(check, "root", "README.md.tmpl")); err != nil {
		t.Error("template file missing on remote")
	}
	msg, _ := r.Run(check, "git", "log", "-1", "--format=%s")
	if strings.TrimSpace(msg) != "template 9.9.9" {
		t.Errorf("message: %q", msg)
	}
	// nothing to push the second time
	res, err = Run(r, opt)
	if err != nil || !res.Nothing || res.Pushed {
		t.Errorf("second push: %+v %v", res, err)
	}
	// tag
	tagged := opt
	tagged.Tag = true
	res, err = Run(r, tagged)
	if err != nil || res.Tag != "v9.9.9" || !res.Pushed {
		t.Fatalf("tag push: %+v %v", res, err)
	}
	tags, _ := r.Run("", "git", "ls-remote", "--tags", bare)
	if !strings.Contains(tags, "refs/tags/v9.9.9") {
		t.Errorf("tag not on remote: %s", tags)
	}
	if _, err := Run(r, tagged); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("existing tag must refuse: %v", err)
	}
	// new branch is created from the default branch
	branch := opt
	branch.Ref = "next"
	res, err = Run(r, branch)
	if err != nil || !res.Created {
		t.Errorf("branch creation: %+v %v", res, err)
	}
	// non-fast-forward: someone else pushed; the git error surfaces verbatim
	other := t.TempDir()
	_, _ = r.Run("", "git", "clone", "-q", bare, other)
	_ = os.WriteFile(filepath.Join(other, "theirs.txt"), []byte("x\n"), 0o644)
	for _, args := range [][]string{{"config", "user.email", "o@o"}, {"config", "user.name", "o"}, {"add", "-A"}, {"commit", "-q", "-m", "theirs"}, {"push", "-q", "origin", "main"}} {
		_, _ = r.Run(other, "git", args...)
	}
	_ = os.WriteFile(filepath.Join(dir, "root", "extra.txt"), []byte("mine\n"), 0o644)
	res, err = Run(r, opt) // clone is fresh, so this is a fast-forward on top of theirs; theirs.txt is removed by replace
	if err != nil || !res.Pushed {
		t.Fatalf("push after external commit: %+v %v", res, err)
	}
	// missing remote: verbatim git error, no panic
	bad := opt
	bad.Remote = filepath.Join(t.TempDir(), "nope.git")
	if _, err := Run(r, bad); err == nil || !strings.Contains(err.Error(), "git clone") {
		t.Errorf("missing remote: %v", err)
	}
	if _, err := Run(r, Options{Dir: dir, Ref: "main", CacheDir: cache}); err == nil || !strings.Contains(err.Error(), "--remote") {
		t.Errorf("no remote configured: %v", err)
	}
}
