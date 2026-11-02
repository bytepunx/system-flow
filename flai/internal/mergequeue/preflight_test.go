package mergequeue

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// markdownTier is a markdown tier that fails on a line holding BAD or
// FIXABLE, naming it as markdownlint-cli2 names a broken link fragment, and
// whose fix makes FIXABLE fixed and leaves BAD.
const markdownTier = `tests:
  - name: markdown
    command: [sh, -c, 's=0; for f in "$@"; do if grep -q -e BAD -e FIXABLE "$f"; then echo "$f:1 error MD051/link-fragments Link fragments should be valid"; s=1; fi; done; exit $s', sh, "{files}"]
    paths: ["**/*.md"]
    fix_command: [sh, -c, 'sed -i s/FIXABLE/fixed/ "$@"', sh, "{files}"]
`

// mainGit is git for the preflight's tests, with no repository: the main
// checkout at root on main at c0ffee whose uncommitted files git status
// gives as status, and a worktree add that checks out a copy of root's
// manifest alone, as main's head holds it, and is removed by worktree
// remove. It keeps the worktree commands it ran.
type mainGit struct {
	scriptedGit
	root  string
	added map[string]bool
	ran   []string
}

// newMainGit is mainGit for the main checkout at root, with status as its
// uncommitted files.
func newMainGit(root, status string) *mainGit {
	g := &mainGit{scriptedGit: scriptedGit{}, root: root, added: map[string]bool{}}
	g.script("main", "rev-parse", "--abbrev-ref", "HEAD")
	g.script("c0ffee", "rev-parse", "--verify", "--quiet", "refs/heads/main^{commit}")
	g.script(status, "status", "--porcelain", "--untracked-files=all")
	g.script(status, "status", "--porcelain")
	return g
}

func (g *mainGit) Run(dir, name string, args ...string) (string, error) {
	switch {
	case name == "git" && len(args) == 6 && args[0] == "worktree" && args[1] == "add":
		path := args[4]
		g.ran = append(g.ran, "worktree add "+path+" "+args[5])
		data, err := os.ReadFile(filepath.Join(g.root, "system-flow.yaml"))
		if err != nil {
			return "", err
		}
		if err := os.MkdirAll(path, 0o755); err != nil {
			return "", err
		}
		g.added[path] = true
		return "", os.WriteFile(filepath.Join(path, "system-flow.yaml"), data, 0o644)
	case name == "git" && len(args) == 4 && args[0] == "worktree" && args[1] == "remove":
		g.ran = append(g.ran, "worktree remove "+args[3])
		delete(g.added, args[3])
		return "", os.RemoveAll(args[3])
	case name == "git" && len(args) == 2 && args[0] == "worktree" && args[1] == "prune":
		return "", nil
	case name == "git" && len(args) == 2 && args[0] == "rev-parse" && args[1] == "--show-toplevel" && g.added[dir]:
		return filepath.EvalSymlinks(dir)
	}
	return g.scriptedGit.Run(dir, name, args...)
}

func (g *mainGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

// preflightProject is a project, not a git repository, with its layout's
// folders and a .gitignore of its cache, so that flai check finds no error
// in it, and extra added to its manifest.
func preflightProject(t *testing.T, extra string) *workitem.Repo {
	t.Helper()
	repo := project(t, "")
	writeFile(t, repo.Root, "system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"+extra)
	writeFile(t, repo.Root, ".gitignore", ".flai-cache/\n")
	for _, d := range []string{"design", "docs"} {
		if err := os.MkdirAll(filepath.Join(repo.Root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// preflight runs Preflight on repo with g at t1, failing the test on an
// error.
func preflight(t *testing.T, g *mainGit, repo *workitem.Repo) []verify.FileFinding {
	t.Helper()
	got, err := Preflight(context.Background(), g, repo, t1, nil)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// ADR-0234: only an error of flai check holds the queue; a warning, of which
// main always has some (D-0338), holds nothing.
func TestThePreflightKeepsOnlyFlaiChecksErrors(t *testing.T) {
	res := &check.Result{Findings: []check.Finding{
		{Level: check.Warning, Rule: "item.archive", Path: filepath.Join("wip", "kanban", "stories", "S-0001-a.md"), Line: 1, Message: "archive it"},
		{Level: check.Error, Rule: "manifest.roles", Path: "system-flow.yaml", Line: 4, Message: "a role is not valid"},
		{Level: check.Warning, Rule: "issues.no-story", Path: "design/issues/I-0001-a.md", Message: "no story"},
	}}
	want := []verify.FileFinding{{Tier: PreflightCheckTier, Checkout: verify.CheckoutMain, Path: "system-flow.yaml", Line: 4, Level: check.Error, Rule: "manifest.roles", Message: "a role is not valid"}}
	if got := checkFindings(res, nil); !reflect.DeepEqual(got, want) {
		t.Errorf("findings %+v, want %+v", got, want)
	}
}

// S-0549: an error of flai check on an uncommitted file outside wip/, which
// a landing leaves as it is (ADR-0178) and the batch never sees, holds
// nothing; one on a committed file, or with no file, still holds.
func TestThePreflightLeavesTheChecksErrorsOnFilesALandingLeaves(t *testing.T) {
	res := &check.Result{Findings: []check.Finding{
		{Level: check.Error, Rule: "issue.front-matter", Path: "design/issues/I-0150-x.md", Line: 1, Message: "no front matter"},
		{Level: check.Error, Rule: "doc.front-matter", Path: "design/analysis/2026-10-19-all.md", Line: 1, Message: "no front matter"},
		{Level: check.Error, Rule: "layout.folder", Path: ".", Line: 1, Message: "no design folder"},
	}}
	left := []string{"design/issues/I-0150-x.md", "design/analysis/"}
	want := []verify.FileFinding{{Tier: PreflightCheckTier, Checkout: verify.CheckoutMain, Path: ".", Line: 1, Level: check.Error, Rule: "layout.folder", Message: "no design folder"}}
	if got := checkFindings(res, left); !reflect.DeepEqual(got, want) {
		t.Errorf("findings %+v, want %+v", got, want)
	}
}

// ADR-0234: the preflight runs flai check on the main checkout and holds on
// its error, here a layout folder that is missing, and finds nothing once it
// is there.
func TestThePreflightHoldsOnAnErrorOfFlaiCheck(t *testing.T) {
	repo := preflightProject(t, "")
	if err := os.Remove(filepath.Join(repo.Root, "docs")); err != nil {
		t.Fatal(err)
	}
	g := newMainGit(repo.Root, "")
	got := preflight(t, g, repo)
	if len(got) != 1 || got[0].Tier != PreflightCheckTier || got[0].Rule != "layout.folder" || got[0].Checkout != verify.CheckoutMain {
		t.Fatalf("findings %+v, want flai check's layout.folder in the main checkout", got)
	}
	if err := os.MkdirAll(filepath.Join(repo.Root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := preflight(t, g, repo); len(got) != 0 {
		t.Errorf("findings %+v once the folder is there, want none", got)
	}
}

// ADR-0234: uncommitted wip files the fixing tiers pass in place hold
// nothing, and need no scratch worktree; nor do files outside wip, which a
// landing leaves, whatever they hold.
func TestThePreflightLintPassesInPlace(t *testing.T) {
	repo := preflightProject(t, markdownTier)
	writeFile(t, repo.Root, "wip/notes/n.md", "# A note\n")
	writeFile(t, repo.Root, "notes.md", "BAD\n")
	g := newMainGit(repo.Root, "?? wip/notes/n.md\n?? notes.md\n")
	if got := preflight(t, g, repo); len(got) != 0 {
		t.Errorf("findings %+v, want none", got)
	}
	if len(g.ran) != 0 {
		t.Errorf("ran %v, want no scratch worktree for a lint that passes in place", g.ran)
	}
}

// ADR-0234: a lint finding in place is confirmed in a scratch worktree of
// main's head, with the file copied in and the fix run there: one the fix
// leaves holds the queue, by the main checkout's path, file, and rule; one
// the fix clears holds nothing. The main checkout's file is never changed,
// and the scratch worktree is removed after.
func TestThePreflightConfirmsALintFindingInAScratchWorktree(t *testing.T) {
	for _, tc := range []struct {
		name, note string
		hold       bool
	}{
		{name: "the fix leaves it", note: "# A note\n\n[gone](#BAD)\n", hold: true},
		{name: "the fix clears it", note: "# A note\n\n[gone](#FIXABLE)\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := preflightProject(t, markdownTier)
			writeFile(t, repo.Root, "wip/notes/n.md", tc.note)
			g := newMainGit(repo.Root, "?? wip/notes/n.md\n")
			got := preflight(t, g, repo)
			if tc.hold {
				if len(got) != 1 || got[0].Path != "wip/notes/n.md" || got[0].Rule != "MD051" || got[0].Tier != "markdown" || got[0].Checkout != verify.CheckoutMain || got[0].Line != 1 {
					t.Errorf("findings %+v, want wip/notes/n.md:1 MD051 of the markdown tier, in the main checkout", got)
				}
			} else if len(got) != 0 {
				t.Errorf("findings %+v, want none: the landing's fix clears it", got)
			}
			if data, err := os.ReadFile(filepath.Join(repo.Root, "wip", "notes", "n.md")); err != nil || string(data) != tc.note {
				t.Errorf("the main checkout's note is %q, %v; want it as it was", data, err)
			}
			scratch := filepath.Join(repo.CacheDir(), "worktrees", PreflightWorktreeName)
			want := []string{"worktree add " + scratch + " c0ffee", "worktree remove " + scratch}
			if !slices.Equal(g.ran, want) {
				t.Errorf("ran %v, want %v", g.ran, want)
			}
			if _, err := os.Stat(scratch); !os.IsNotExist(err) {
				t.Errorf("the scratch worktree %s is still there: %v", scratch, err)
			}
		})
	}
}

// A preflight that cannot read the main checkout's uncommitted files is an
// error saying what it could not do.
func TestAPreflightThatCannotRunIsAnError(t *testing.T) {
	repo := preflightProject(t, "")
	_, err := Preflight(context.Background(), scriptedGit{}, repo, t1, nil)
	if err == nil || !strings.Contains(err.Error(), "run the merge queue's preflight") {
		t.Errorf("err %v, want one naming the preflight", err)
	}
}
