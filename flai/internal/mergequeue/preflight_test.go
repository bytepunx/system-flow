package mergequeue

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// markdownTier is a markdown tier that fails on a line holding BAD or
// FIXABLE, naming it as markdownlint-cli2 names a broken link fragment, and
// whose fix makes FIXABLE fixed and leaves BAD. The fix writes each file
// through a temporary one rather than with sed -i, whose argument BSD sed on
// macOS reads as a backup suffix, so that it runs the same on every host.
const markdownTier = `tests:
  - name: markdown
    command: [sh, -c, 's=0; for f in "$@"; do if grep -q -e BAD -e FIXABLE "$f"; then echo "$f:1 error MD051/link-fragments Link fragments should be valid"; s=1; fi; done; exit $s', sh, "{files}"]
    paths: ["**/*.md"]
    fix_command: [sh, -c, 'for f in "$@"; do sed s/FIXABLE/fixed/ "$f" > "$f.tmp" && mv "$f.tmp" "$f"; done', sh, "{files}"]
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
	g.script(status, "-c", "core.quotePath=false", "status", "--porcelain", "--no-renames", "--untracked-files=all")
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

// preflight runs Preflight on repo with g at t1, by a flai git cannot
// place, failing the test on an error or a report that it is behind.
func preflight(t *testing.T, g *mainGit, repo *workitem.Repo) []verify.FileFinding {
	t.Helper()
	got, err := Preflight(context.Background(), g, repo, t1, buildinfo.Info{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Behind != nil {
		t.Fatalf("behind %+v, want nothing for a flai git cannot place", got.Behind)
	}
	return got.Findings
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

// S-0549: an error of flai check on an uncommitted file outside the shared
// set, the operator's, which a landing leaves as it is (ADR-0178, ADR-0248)
// and the batch never sees, holds nothing; one on a committed file, or with
// no file, still holds.
func TestThePreflightLeavesTheChecksErrorsOnFilesALandingLeaves(t *testing.T) {
	res := &check.Result{Findings: []check.Finding{
		{Level: check.Error, Rule: "doc.front-matter", Path: "docs/plan.md", Line: 1, Message: "no front matter"},
		{Level: check.Error, Rule: "doc.front-matter", Path: "notes/2026-10-19.md", Line: 1, Message: "no front matter"},
		{Level: check.Error, Rule: "layout.folder", Path: ".", Line: 1, Message: "no design folder"},
	}}
	left := []string{"docs/plan.md", "notes/"}
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

// ADR-0234, ADR-0248: uncommitted files of the shared set the fixing tiers
// pass in place hold nothing, and need no scratch worktree; nor do files
// outside it, the operator's, which a landing leaves, whatever they hold.
func TestThePreflightLintPassesInPlace(t *testing.T) {
	repo := preflightProject(t, markdownTier)
	writeFile(t, repo.Root, "wip/notes/n.md", "# A note\n")
	writeFile(t, repo.Root, "design/notes/n.md", "# A design note\n")
	writeFile(t, repo.Root, "notes.md", "BAD\n")
	g := newMainGit(repo.Root, "?? wip/notes/n.md\n?? design/notes/n.md\n?? notes.md\n")
	if got := preflight(t, g, repo); len(got) != 0 {
		t.Errorf("findings %+v, want none", got)
	}
	if len(g.ran) != 0 {
		t.Errorf("ran %v, want no scratch worktree for a lint that passes in place", g.ran)
	}
}

// ADR-0234, ADR-0248: a lint finding in place in a file of the shared set,
// under wip or design, is confirmed in a scratch worktree of main's head,
// with the file copied in and the fix run there: one the fix leaves holds
// the queue, by the main checkout's path, file, and rule; one the fix clears
// holds nothing. The main checkout's file is never changed, and the scratch
// worktree is removed after.
func TestThePreflightConfirmsALintFindingInAScratchWorktree(t *testing.T) {
	for _, tc := range []struct {
		name, path, note string
		hold             bool
	}{
		{name: "the fix leaves it", path: "wip/notes/n.md", note: "# A note\n\n[gone](#BAD)\n", hold: true},
		{name: "the fix leaves it under design", path: "design/notes/n.md", note: "# A note\n\n[gone](#BAD)\n", hold: true},
		{name: "the fix clears it", path: "wip/notes/n.md", note: "# A note\n\n[gone](#FIXABLE)\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := preflightProject(t, markdownTier)
			writeFile(t, repo.Root, tc.path, tc.note)
			g := newMainGit(repo.Root, "?? "+tc.path+"\n")
			got := preflight(t, g, repo)
			if tc.hold {
				if len(got) != 1 || got[0].Path != tc.path || got[0].Rule != "MD051" || got[0].Tier != "markdown" || got[0].Checkout != verify.CheckoutMain || got[0].Line != 1 {
					t.Errorf("findings %+v, want %s:1 MD051 of the markdown tier, in the main checkout", got, tc.path)
				}
			} else if len(got) != 0 {
				t.Errorf("findings %+v, want none: the landing's fix clears it", got)
			}
			if data, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(tc.path))); err != nil || string(data) != tc.note {
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
	_, err := Preflight(context.Background(), scriptedGit{}, repo, t1, buildinfo.Info{}, nil)
	if err == nil || !strings.Contains(err.Error(), "run the merge queue's preflight") {
		t.Errorf("err %v, want one naming the preflight", err)
	}
}

// flaiProjectManifest is a manifest's flai sub-project at flai/.
const flaiProjectManifest = "projects:\n  - name: flai\n    path: flai\n    kind: go\n"

// fakeFlai is the main checkout's flai.command in behindProject: it records
// its arguments in its folder's runs, prints check.out and exits with
// check.status for check, and prints version.out for version.
type fakeFlai struct {
	dir string
}

// set writes name in the fake's folder as content.
func (f fakeFlai) set(t *testing.T, name, content string) {
	t.Helper()
	writeFile(t, f.dir, name, content)
}

// runs are the argument lists the fake ran with, none when it never ran.
func (f fakeFlai) runs(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.dir, "runs"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// flaiGit is mainGit that runs every program but git for real, as the
// main checkout's flai.command runs.
type flaiGit struct {
	*mainGit
}

func (g flaiGit) Run(dir, name string, args ...string) (string, error) {
	if name != "git" {
		return execx.System{}.Run(dir, name, args...)
	}
	return g.mainGit.Run(dir, name, args...)
}

func (g flaiGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

// behindOld is the commit the running flai is built from in behindProject,
// one commit under flai/ behind main's head, c0ffee.
const behindOld = "a1d0000"

// behindProject is a project, as main's history holds one in which the
// running flai is behind, and git scripted to say so: its flai sub-project,
// flai/, has one commit at main's head, c0ffee, after behindOld, and git
// places no other commit and holds no flai/v* tag. Its layout's docs folder
// is missing, an error flai check reports (layout.folder). With command,
// its manifest names scripts/fake-flai, a fakeFlai, as flai.command.
func behindProject(t *testing.T, command bool) (*workitem.Repo, flaiGit, fakeFlai) {
	t.Helper()
	extra := flaiProjectManifest
	if command {
		extra += "flai:\n  command: scripts/fake-flai\n"
	}
	repo := preflightProject(t, extra)
	fake := fakeFlai{dir: t.TempDir()}
	script := fmt.Sprintf("#!/bin/sh\necho \"$*\" >> '%[1]s/runs'\ncase \"$1\" in\ncheck) cat '%[1]s/check.out'; exit \"$(cat '%[1]s/check.status')\" ;;\nversion) cat '%[1]s/version.out' ;;\nesac\n", fake.dir)
	writeFile(t, repo.Root, "scripts/fake-flai", script)
	if err := os.Chmod(filepath.Join(repo.Root, "scripts", "fake-flai"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, repo.Root, "flai/main.go", "package main\n")
	if err := os.Remove(filepath.Join(repo.Root, "docs")); err != nil {
		t.Fatal(err)
	}
	g := flaiGit{newMainGit(repo.Root, "")}
	g.script("c0ffee", "rev-parse", "--verify", "c0ffee^{commit}")
	g.script("", "tag", "--merged", "c0ffee", "--list", "flai/v*")
	for _, c := range []string{behindOld, "c0ffee"} {
		g.script("", "cat-file", "-e", c+"^{commit}")
	}
	g.script("1", "rev-list", "--count", behindOld+"..c0ffee", "--", "flai")
	g.script("0", "rev-list", "--count", "c0ffee..c0ffee", "--", "flai")
	return repo, g, fake
}

// cleanCheck is flai check --json's result with no finding, after a build
// line, as scripts/flai.sh prints it.
const cleanCheck = "building bin/flai\n{\n  \"findings\": [],\n  \"errors\": 0,\n  \"warnings\": 0,\n  \"advisory\": 0,\n  \"items\": 0\n}\n"

// confirmedCheck is flai check --json's result with the docs folder's
// layout.folder error, as main's own flai reports it.
const confirmedCheck = "{\n  \"findings\": [\n    {\n      \"level\": \"error\",\n      \"rule\": \"layout.folder\",\n      \"path\": \"docs\",\n      \"line\": 1,\n      \"message\": \"confirmed by main's flai\"\n    }\n  ],\n  \"errors\": 1,\n  \"warnings\": 0,\n  \"advisory\": 0,\n  \"items\": 0\n}\n"

// versionJSON is flai version --json naming version built from commit.
func versionJSON(version, commit string) string {
	return fmt.Sprintf("{\n  \"version\": %q,\n  \"commit\": %q,\n  \"date\": \"\",\n  \"go\": \"go1.26\",\n  \"task_done_rules\": 1\n}\n", version, commit)
}

// I-0382, ADR-0311: an error the running flai's check reports on main, as
// flai 1.99.0 reported roles.file on role files in a form only a later flai
// reads, does not hold the queue when the running flai is behind main's
// flai sub-project and the flai main builds, run through flai.command, finds
// main clean: the preflight checks main again with it, and the queue is not
// held.
func TestThePreflightChecksMainAgainWithMainsFlaiWhenTheRunningFlaiIsBehind(t *testing.T) {
	repo, g, fake := behindProject(t, true)
	fake.set(t, "check.out", cleanCheck)
	fake.set(t, "check.status", "0")
	fake.set(t, "version.out", "building bin/flai\n"+versionJSON("1.100.0", "c0ffee"))
	running := buildinfo.Info{Version: "1.99.0", Commit: behindOld}
	var logs strings.Builder
	d := Deps{Runner: g, Now: func() time.Time { return t1 }, Log: slog.New(slog.NewJSONHandler(&logs, nil)),
		Preflight: func(ctx context.Context, repo *workitem.Repo) (PreflightReport, error) {
			return Preflight(ctx, g, repo, t1, running, slog.New(slog.NewJSONHandler(&logs, nil)))
		}}
	clean, err := preflightLook(context.Background(), d, repo, true)
	if err != nil {
		t.Fatal(err)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if !clean || st.Preflight != nil {
		t.Fatalf("clean %v, hold %+v; want main clean by main's flai and the queue not held", clean, st.Preflight)
	}
	if got, want := fake.runs(t), []string{"check --json", "version --json"}; !slices.Equal(got, want) {
		t.Errorf("flai.command ran %v, want %v", got, want)
	}
	for _, want := range []string{"checked main again", `"findings_before":1`, `"findings_after":0`} {
		if !strings.Contains(logs.String(), want) {
			t.Errorf("logs lack %q: %s", want, logs.String())
		}
	}
}

// ADR-0311: where the running flai is behind main's flai sub-project, the
// findings of a second check that counts replace the running flai's, with
// no cause; where it cannot settle them, the running flai's stand with the
// cause beside them, saying why. A running flai that is not behind checks
// once, and flai.command never runs.
func TestThePreflightSaysWhyTheRunningFlaisFindingsStand(t *testing.T) {
	for _, tc := range []struct {
		name    string
		command bool
		// atHead builds the running flai from main's head, which is not behind.
		atHead                 bool
		check, status, version string
		// versionAt is the commit the fake's version names.
		versionAt string
		// confirmed says the findings are the fake's.
		confirmed bool
		// why is what the cause says, "" for none.
		why  string
		runs []string
	}{
		{name: "main's flai confirms the finding", command: true, check: confirmedCheck, status: "1", version: "1.100.0", versionAt: "c0ffee", confirmed: true,
			runs: []string{"check --json", "version --json"}},
		{name: "the command gives no result", command: true, check: "building bin/flai\ngo: build failed\n", status: "1", version: "1.100.0", versionAt: "c0ffee",
			why: "scripts/fake-flai check --json, the manifest's flai.command, gave no check result", runs: []string{"check --json"}},
		{name: "the command's flai is behind too", command: true, check: cleanCheck, status: "0", version: "1.99.0", versionAt: behindOld,
			why: "flai 1.99.0, which scripts/fake-flai runs, is behind main too", runs: []string{"check --json", "version --json"}},
		{name: "the command does not say its version", command: true, check: cleanCheck, status: "0",
			why: "scripts/fake-flai version --json, the manifest's flai.command, did not say its version", runs: []string{"check --json", "version --json"}},
		{name: "no flai.command", why: "the manifest names no flai.command"},
		{name: "the running flai is built from main's head", command: true, atHead: true, check: cleanCheck, status: "0", version: "1.100.0", versionAt: "c0ffee"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, g, fake := behindProject(t, tc.command)
			fake.set(t, "check.out", tc.check)
			fake.set(t, "check.status", tc.status)
			if tc.version != "" {
				fake.set(t, "version.out", versionJSON(tc.version, tc.versionAt))
			}
			running := buildinfo.Info{Version: "1.99.0", Commit: behindOld}
			if tc.atHead {
				running = buildinfo.Info{Version: "1.100.0", Commit: "c0ffee"}
			}
			got, err := Preflight(context.Background(), g, repo, t1, running, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Findings) != 1 || got.Findings[0].Rule != "layout.folder" || got.Findings[0].Tier != PreflightCheckTier {
				t.Fatalf("findings %+v, want the docs folder's layout.folder", got.Findings)
			}
			if confirmed := got.Findings[0].Message == "confirmed by main's flai"; confirmed != tc.confirmed {
				t.Errorf("finding %+v; want it from main's flai %v", got.Findings[0], tc.confirmed)
			}
			switch {
			case tc.why == "" && got.Behind != nil:
				t.Errorf("behind %+v, want none", got.Behind)
			case tc.why != "" && got.Behind == nil:
				t.Errorf("no cause, want one saying %q", tc.why)
			case tc.why != "":
				want := release.Behind{Running: "1.99.0", Commit: behindOld, By: "commit", From: behindOld, Main: "c0ffee", Path: "flai", Commits: 1}
				b := got.Behind.Behind
				b.Message = ""
				if b != want || !strings.Contains(got.Behind.Why, tc.why) || !strings.Contains(got.Behind.Message, "flai release --pending") {
					t.Errorf("behind %+v, want %+v, why saying %q, and the remedy", got.Behind, want, tc.why)
				}
			}
			if runs := fake.runs(t); !slices.Equal(runs, tc.runs) {
				t.Errorf("flai.command ran %v, want %v", runs, tc.runs)
			}
		})
	}
}

// ADR-0311: the hold records the cause the preflight reports beside its
// findings, refreshes it on each look as it does the findings, keeping when
// the hold began and its thread, and drops it once a flai built from main
// confirms the findings; the preflight_held event keeps its fields.
func TestThePreflightHoldRecordsWhyTheRunningFlaiMayBeWrong(t *testing.T) {
	repo := project(t, "")
	g := scriptedGit{}
	g.script("main", "rev-parse", "--abbrev-ref", "HEAD")
	g.script("c0ffee", "rev-parse", "--verify", "--quiet", "refs/heads/main^{commit}")
	role := verify.FileFinding{Tier: PreflightCheckTier, Checkout: verify.CheckoutMain, Path: "design/workflows/roles/story.md", Line: 1, Level: "error", Rule: "roles.file", Message: "not a role file"}
	behind := &FlaiBehind{Behind: release.Behind{Running: "1.99.0", Commit: "abc1234", By: "commit", From: "abc1234", Main: "c0ffee", Path: "flai", Commits: 2, Message: "flai 1.99.0 lacks 2 commits"},
		Why: "the manifest names no flai.command"}
	report, at := PreflightReport{Findings: []verify.FileFinding{role}, Behind: behind}, t1
	d := Deps{Runner: g, Now: func() time.Time { return at }, Log: slog.New(slog.DiscardHandler),
		Preflight: func(context.Context, *workitem.Repo) (PreflightReport, error) { return report, nil }}
	look := func() *PreflightHold {
		t.Helper()
		if clean, err := preflightLook(context.Background(), d, repo, true); err != nil || clean {
			t.Fatalf("clean %v, %v; want the queue held", clean, err)
		}
		st, err := ReadState(repo.CacheDir())
		if err != nil {
			t.Fatal(err)
		}
		return st.Preflight
	}
	if got, want := look(), (&PreflightHold{Since: t1, Main: "c0ffee", Findings: []verify.FileFinding{role}, FlaiBehind: behind}); !reflect.DeepEqual(got, want) {
		t.Fatalf("hold %+v, want %+v", got, want)
	}
	if set, err := SetPreflightThread(repo.CacheDir(), "TH-0042"); err != nil || !set {
		t.Fatalf("thread set %v, %v", set, err)
	}
	later := *behind
	later.Why = "flai 1.99.0, which scripts/flai.sh runs, is behind main too"
	report.Behind, at = &later, t1.Add(time.Minute)
	if got := look(); !reflect.DeepEqual(got.FlaiBehind, &later) || !got.Since.Equal(t1) || got.Thread != "TH-0042" {
		t.Errorf("hold %+v; want the new cause, since %s, on TH-0042", got, t1)
	}
	report.Behind, at = nil, t1.Add(2*time.Minute)
	if got := look(); got.FlaiBehind != nil || got.Thread != "TH-0042" {
		t.Errorf("hold %+v; want no cause once main's flai confirms the findings", got)
	}
	held := Event{At: t1, Event: EventPreflightHeld, Stories: []string{}, Files: []string{role.Path}, Rules: []string{"roles.file"}}
	if got := preflightEventsIn(t, repo.CacheDir()); !reflect.DeepEqual(got, []Event{held}) {
		t.Errorf("events %+v, want %+v alone", got, held)
	}
}
