package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// scopeFixture copies the good fixture, whose S-004 is in progress, and
// gives the archived S-002 a branch that was never merged (I-0057).
func scopeFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	refs := filepath.Join(root, ".git", "refs", "heads", "story")
	if err := os.MkdirAll(refs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refs, "S-002"), []byte("0123456789abcdef0123456789abcdef01234567\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// S-0249: --story S-nnnn passes over the findings outside the story, says
// how many in the summary, and marks them in --json.
func TestCheckStoryPassesOverFindingsOutsideIt(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := scopeFixture(t)
	out, errOut, code := runIn(t, root, "check", "--strict")
	if code != 1 || !strings.Contains(out, "warning: story.unaccepted: S-002 is done but its branch story/S-002 was never merged; merge or delete it\n") {
		t.Fatalf("unscoped --strict should fail on S-002's branch: %d\n%s%s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "check", "--strict", "--story", "S-4")
	if code != 0 {
		t.Fatalf("scoped to S-004, --strict should pass: %d\n%s%s", code, out, errOut)
	}
	for _, want := range []string{
		"warning: story.unaccepted: S-002 is done but its branch story/S-002 was never merged; merge or delete it (outside S-004)\n",
		"10 items checked, 0 errors, 2 warnings; 2 outside S-004, notes the run passes over\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}

	out, _, code = runIn(t, root, "check", "--strict", "--story", "S-004", "--json")
	var res struct {
		Warnings int `json:"warnings"`
		Advisory int `json:"advisory"`
		Outside  int `json:"outside"`
		Findings []struct {
			Rule    string `json:"rule"`
			Outside bool   `json:"outside"`
		} `json:"findings"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("json: %d %v %s", code, err, out)
	}
	if res.Outside != 2 || res.Warnings != 2 || res.Advisory != 0 || len(res.Findings) != 2 {
		t.Fatalf("json counts: %+v", res)
	}
	for _, f := range res.Findings {
		if !f.Outside {
			t.Errorf("%s should carry outside: %s", f.Rule, out)
		}
	}
	// Unscoped, --json has no outside at all.
	out, _, _ = runIn(t, root, "check", "--json")
	if strings.Contains(out, `"outside"`) {
		t.Errorf("unscoped json should not mention outside:\n%s", out)
	}

	if _, errOut, code := runIn(t, root, "check", "--story", "T-003"); code == 0 || !strings.Contains(errOut, "scope the check to T-003: it is a task; name a story") {
		t.Errorf("a task is not a scope: %d %s", code, errOut)
	}
}

// gitScript answers the git commands flai stream diff and git status run,
// for a story branch that changes design/system/overview.md, or, with
// porcelain, a worktree that has it uncommitted; statusErr fails git status.
type gitScript struct {
	changed, porcelain string
	statusErr          error
}

func (g gitScript) Run(_, _ string, args ...string) (string, error) {
	cmd := strings.Join(args, " ")
	switch {
	case cmd == "rev-parse --is-inside-work-tree":
		return "true", nil
	case cmd == "rev-parse --verify --quiet refs/heads/story/S-004":
		return "abc", nil
	case cmd == "rev-parse --abbrev-ref HEAD":
		return "main", nil
	case strings.HasPrefix(cmd, "merge-base "), strings.HasPrefix(cmd, "rev-parse --short "):
		return "base", nil
	case strings.HasPrefix(cmd, "rev-list --count "):
		return "1", nil
	case strings.HasPrefix(cmd, "diff --name-status "):
		if g.changed == "" {
			return "", nil
		}
		return "M\x00" + g.changed + "\x00", nil
	case strings.HasPrefix(cmd, "diff --numstat "):
		if g.changed == "" {
			return "", nil
		}
		return "1\t1\t" + g.changed + "\x00", nil
	case strings.HasPrefix(cmd, "diff --no-color "):
		return "@@ -1 +1 @@\n", nil
	case cmd == "status --porcelain":
		return g.porcelain, g.statusErr
	}
	return "", errors.New("unscripted: git " + cmd)
}

func (g gitScript) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (gitScript) LookPath(name string) (string, error) { return name, nil }

// S-0249: a finding on a path the story's branch changes, committed or
// uncommitted in its worktree, is the story's and still fails --strict.
func TestCheckStoryFailsOnAPathItChanges(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	for _, tc := range []struct {
		name string
		git  gitScript
		code int
	}{
		{"committed on its branch", gitScript{changed: "design/system/overview.md"}, 1},
		{"uncommitted in its worktree", gitScript{porcelain: " M design/system/overview.md\n"}, 1},
		{"changed elsewhere", gitScript{changed: "flai/main.go"}, 0},
		{"a worktree that does not read", gitScript{statusErr: errors.New("fatal: not a git repository")}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := scopeFixture(t)
			overview := filepath.Join(root, "design/system/overview.md")
			data, _ := os.ReadFile(overview)
			if err := os.WriteFile(overview, []byte(strings.Replace(string(data), "title: Overview\n", "", 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(root, ".flai-cache", "worktrees", "S-004"), 0o755); err != nil {
				t.Fatal(err)
			}
			out, errOut, code := runWithApp(t, &app{cwd: root, runner: tc.git}, "check", "--strict", "--story", "S-004")
			if code != tc.code {
				t.Fatalf("exit %d, want %d:\n%s%s", code, tc.code, out, errOut)
			}
			line := "design/system/overview.md:1: warning: doc.title: front matter needs a title"
			if tc.git.statusErr != nil {
				if !strings.Contains(errOut, "read what S-004 has uncommitted in") || strings.Contains(out, "items checked") {
					t.Errorf("an unreadable worktree should stop the run:\n%s%s", out, errOut)
				}
				return
			}
			if tc.code == 1 && !strings.Contains(out, line+"\n") {
				t.Errorf("doc.title should be the story's:\n%s", out)
			}
			if tc.code == 0 && !strings.Contains(out, line+" (outside S-004)\n") {
				t.Errorf("doc.title should be outside the story:\n%s", out)
			}
		})
	}
}
