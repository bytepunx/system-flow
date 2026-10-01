package release

import (
	"strings"
	"testing"
)

// S-0181: the running flai is older when a flai/v* tag reachable from HEAD
// is newer; a tag on a branch HEAD does not include, a dev build, and a flai
// at or above the newest say nothing.
func TestFlaiOutdated(t *testing.T) {
	root, r := gitRepo(t)
	git := func(args ...string) {
		t.Helper()
		if out, err := r.Run(root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("tag", "-a", "flai/v1.9.0", "-m", "flai 1.9.0")
	git("tag", "-a", "flai/v1.10.0", "-m", "flai 1.10.0")
	git("checkout", "-q", "-b", "later")
	git("commit", "-q", "--allow-empty", "-m", "later")
	git("tag", "-a", "flai/v1.11.0", "-m", "flai 1.11.0")
	git("checkout", "-q", "main")

	out := FlaiOutdated(r, root, "1.9.3")
	if out == nil || out.Running != "1.9.3" || out.Newest != "1.10.0" || !strings.HasPrefix(out.Upgrade, "flai host upgrade") ||
		!strings.Contains(out.Message, "this flai is 1.9.3, older than flai 1.10.0") || !strings.Contains(out.Message, "flai host upgrade") {
		t.Fatalf("1.9.3 under 1.10.0: %+v", out)
	}
	for _, running := range []string{"1.10.0", "v1.10.1", "dev"} {
		if got := FlaiOutdated(r, root, running); got != nil {
			t.Errorf("%s: %+v", running, got)
		}
	}
}
