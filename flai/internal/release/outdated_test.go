package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
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
	if out.NewestRules != 0 {
		t.Errorf("a tag whose commit has no task-done rules file carries rules %d", out.NewestRules)
	}
	for _, running := range []string{"1.10.0", "v1.10.1", "dev"} {
		if got := FlaiOutdated(r, root, running); got != nil {
			t.Errorf("%s: %+v", running, got)
		}
	}
}

// ADR-0180: the newest tag's task-done rules number is the rules file at the
// commit it tags, not at HEAD, and a file that holds no number carries none.
func TestFlaiOutdatedReadsTheTagsTaskDoneRules(t *testing.T) {
	root, r := gitRepo(t)
	commit := func(rules string) {
		t.Helper()
		p := filepath.Join(root, buildinfo.TaskDoneRulesPath)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(rules), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "rules " + rules}} {
			if out, err := r.Run(root, "git", args...); err != nil {
				t.Fatalf("git %v: %v %s", args, err, out)
			}
		}
	}
	tag := func(name string) {
		t.Helper()
		if out, err := r.Run(root, "git", "tag", "-a", name, "-m", name); err != nil {
			t.Fatalf("git tag %s: %v %s", name, err, out)
		}
	}
	commit("3\n")
	tag("flai/v1.12.0")
	commit("4\n")
	if out := FlaiOutdated(r, root, "1.11.0"); out == nil || out.Newest != "1.12.0" || out.NewestRules != 3 {
		t.Errorf("1.11.0 under 1.12.0 with rules 3: %+v", out)
	}
	commit("none\n")
	tag("flai/v1.13.0")
	if out := FlaiOutdated(r, root, "1.11.1"); out == nil || out.Newest != "1.13.0" || out.NewestRules != 0 {
		t.Errorf("1.11.1 under 1.13.0 with no number: %+v", out)
	}
}
