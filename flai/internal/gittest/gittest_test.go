package gittest

import (
	"os/exec"
	"strings"
	"testing"
)

func committer(t *testing.T) (string, error) {
	t.Helper()
	dir := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", dir).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, out)
	}
	c := exec.Command("git", "var", "GIT_COMMITTER_IDENT")
	c.Dir = dir
	out, err := c.Output()
	return string(out), err
}

func TestGitHasOnlyTheIdentityTheTestGivesIt(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Run("none", func(t *testing.T) {
		NoIdentity(t)
		if got, err := committer(t); err == nil {
			t.Errorf("git worked out an identity: %s", got)
		}
	})
	t.Run("given", func(t *testing.T) {
		Identity(t, "Dana Designer", "dana@example.com")
		got, err := committer(t)
		if err != nil || !strings.HasPrefix(got, "Dana Designer <dana@example.com>") {
			t.Errorf("identity: %q %v", got, err)
		}
	})
}
