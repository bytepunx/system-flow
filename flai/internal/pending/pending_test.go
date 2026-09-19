package pending

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// lab is a clone with a bare remote named origin and main pushed to it.
func lab(t *testing.T) (root, remote string, git func(dir string, args ...string) string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	git = func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return string(out)
	}
	base := t.TempDir()
	remote = filepath.Join(base, "origin.git")
	root = filepath.Join(base, "clone")
	git(base, "init", "-q", "--bare", "-b", "main", remote)
	git(base, "init", "-q", "-b", "main", root)
	git(root, "commit", "-q", "--allow-empty", "-m", "init")
	git(root, "remote", "add", "origin", remote)
	git(root, "push", "-q", "-u", "origin", "main")
	return root, remote, git
}

func TestDetect(t *testing.T) {
	r := execx.System{}
	t.Run("nothing ahead", func(t *testing.T) {
		root, _, _ := lab(t)
		if u := Detect(r, root); u != nil {
			t.Errorf("got %+v", u)
		}
	})
	t.Run("ahead without an acceptance is not pending", func(t *testing.T) {
		root, _, git := lab(t)
		git(root, "commit", "-q", "--allow-empty", "-m", "docs: a note")
		u := Detect(r, root)
		if u == nil || u.Commits != 1 || len(u.Acceptances) != 0 || u.Pending() {
			t.Errorf("got %+v", u)
		}
	})
	t.Run("an acceptance with its tags", func(t *testing.T) {
		root, _, git := lab(t)
		git(root, "tag", "-a", "cli/v1.0.0", "-m", "already on the remote")
		git(root, "push", "-q", "origin", "cli/v1.0.0")
		git(root, "commit", "-q", "--allow-empty", "-m", "feat: [S-0007] the change")
		git(root, "commit", "-q", "--allow-empty", "-m", "chore: [S-0007] accept and archive; release cli 1.1.0")
		git(root, "tag", "-a", "cli/v1.1.0", "-m", "cli 1.1.0")
		git(root, "commit", "-q", "--allow-empty", "-m", "chore: [S-0008] accept and archive")
		u := Detect(r, root)
		if u == nil || !u.Pending() {
			t.Fatalf("got %+v", u)
		}
		if u.Branch != "main" || u.Upstream != "origin/main" || u.Remote != "origin" || u.Commits != 3 || u.Behind != 0 {
			t.Errorf("got %+v", u)
		}
		if want := []string{"S-0008", "S-0007"}; !reflect.DeepEqual(u.Acceptances, want) {
			t.Errorf("acceptances, newest first: %v", u.Acceptances)
		}
		if want := []string{"cli/v1.1.0"}; !reflect.DeepEqual(u.Tags, want) {
			t.Errorf("only the tags on unpushed commits: %v", u.Tags)
		}
		if want := []string{"main", "cli/v1.1.0"}; !reflect.DeepEqual(u.Refs(), want) {
			t.Errorf("refs: %v", u.Refs())
		}
		// after a push from this clone it is no longer true, at once
		git(root, "push", "-q", "origin", "main", "cli/v1.1.0")
		if u := Detect(r, root); u != nil {
			t.Errorf("after the push: %+v", u)
		}
	})
	t.Run("diverged", func(t *testing.T) {
		root, remote, git := lab(t)
		other := filepath.Join(t.TempDir(), "other")
		git(filepath.Dir(other), "clone", "-q", remote, other)
		git(other, "commit", "-q", "--allow-empty", "-m", "docs: from elsewhere")
		git(other, "push", "-q", "origin", "main")
		git(root, "commit", "-q", "--allow-empty", "-m", "chore: [S-0009] accept and archive")
		if u := Detect(r, root); u == nil || u.Behind != 0 {
			t.Errorf("a push from another clone is not seen until someone fetches here: %+v", u)
		}
		git(root, "fetch", "-q", "origin")
		if u := Detect(r, root); u == nil || u.Behind != 1 || u.Commits != 1 {
			t.Errorf("after a fetch: %+v", u)
		}
	})
	t.Run("pushed without -u: origin's branch of the same name stands in", func(t *testing.T) {
		root, _, git := lab(t)
		git(root, "branch", "--unset-upstream")
		git(root, "commit", "-q", "--allow-empty", "-m", "chore: [S-0011] accept and archive")
		if u := Detect(r, root); u == nil || u.Upstream != "origin/main" || !u.Pending() {
			t.Errorf("got %+v", u)
		}
	})
	t.Run("no upstream, no remote, no repository", func(t *testing.T) {
		root, _, git := lab(t)
		git(root, "checkout", "-q", "-b", "side")
		git(root, "commit", "-q", "--allow-empty", "-m", "chore: [S-0010] accept and archive")
		if u := Detect(r, root); u != nil {
			t.Errorf("a branch with no upstream: %+v", u)
		}
		if u := Detect(r, t.TempDir()); u != nil {
			t.Errorf("not a repository: %+v", u)
		}
		if u := Detect(nil, root); u != nil {
			t.Errorf("no runner: %+v", u)
		}
	})
}
