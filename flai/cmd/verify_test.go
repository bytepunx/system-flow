package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// verifyGit answers the git commands flai verify runs for a story whose
// branch changed docs/README.md, has nothing uncommitted, and has no rebase
// left, and refuses every other, as git refuses a repository it is not run
// in. The branch contains main but for commits, the commits of main it
// lacks, which changed paths since the branch left it; none when they are "".
type verifyGit struct {
	commits, paths string
}

func (g verifyGit) Run(_, _ string, args ...string) (string, error) {
	switch strings.Join(args, " ") {
	case "rev-parse --abbrev-ref HEAD":
		return "main", nil
	case "rev-parse HEAD":
		return "0123abc4567\n", nil
	case "diff --name-only --diff-filter=d main...HEAD":
		return "docs/README.md\n", nil
	case "status --porcelain --untracked-files=all":
		return "", nil
	case "diff --no-renames --name-only HEAD...main":
		return g.paths, nil
	case "rev-list --abbrev-commit HEAD..main":
		return g.commits, nil
	}
	return "", errors.New("fatal: not a git repository")
}

func (g verifyGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (verifyGit) LookPath(name string) (string, error) { return name, nil }

// verifyTier is a manifest tier, unit, that runs script with sh for the
// files under docs/ the branch changed.
func verifyTier(script string) string {
	return "tests:\n  - name: unit\n    command: [sh, -c, '" + script + "', unit, \"{files}\"]\n    paths: [\"docs/**\"]\n"
}

// verifyFixture is scopeFixture, whose S-004 is in progress with its
// narrative written and whose S-002's branch was never merged, a note
// outside S-004, with S-004's worktree under .flai-cache/worktrees: a copy
// of the project whose manifest declares the unit tier running script. It
// answers the project's main checkout and the worktree.
func verifyFixture(t *testing.T, script string) (root, worktree string) {
	t.Helper()
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root = scopeFixture(t)
	worktree = filepath.Join(root, ".flai-cache", "worktrees", "S-004")
	if err := os.CopyFS(worktree, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(worktree, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"system-flow.yaml": string(manifest) + verifyTier(script),
		// a linked worktree names its main checkout's git directory
		".git": "gitdir: " + filepath.Join(root, ".git", "worktrees", "S-004") + "\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(worktree, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root, worktree
}

// verifyReport runs flai verify --json in dir and decodes its answer.
func verifyReport(t *testing.T, dir string, wantCode int, args ...string) verify.Report {
	t.Helper()
	out, errOut, code := runVerify(t, dir, verifyGit{}, append([]string{"verify", "--json"}, args...)...)
	if code != wantCode {
		t.Fatalf("flai verify %v: exit %d, want %d\n%s%s", args, code, wantCode, out, errOut)
	}
	var rep verify.Report
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatalf("flai verify %v: not JSON: %v\n%s", args, err, out)
	}
	return rep
}

// stepStates are the report's steps as name:state, a tier's name after
// "tier ".
func stepStates(rep verify.Report) string {
	var out []string
	for _, s := range rep.Steps {
		name := s.Name
		if s.Tier {
			name = "tier " + name
		}
		out = append(out, name+":"+string(s.State))
	}
	return strings.Join(out, ",")
}

// runVerify runs the CLI in dir with git answered by git and the clock
// fixed, and logs a refusal's message on stderr as the command boundary
// does.
func runVerify(t *testing.T, dir string, git execx.Runner, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: dir, runner: git, clock: func() time.Time { return time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC) }}
	root := newRootCmdWith(a)
	root.SetArgs(args)
	code := 0
	if err := root.Execute(); err != nil {
		var ee *exitError
		switch {
		case errors.As(err, &ee):
			code = ee.code
			if ee.msg != "" {
				a.fail(err)
			}
		default:
			a.fail(err)
			code = 1
		}
	}
	return out.String(), errOut.String(), code
}

// openProject opens the project at root.
func openProject(t *testing.T, root string) *workitem.Repo {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// lastLine is the text's last line.
func lastLine(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	return lines[len(lines)-1]
}

func TestVerifyPassesEveryStepAndSaysSoLast(t *testing.T) {
	root, worktree := verifyFixture(t, "test -f \"$1\"")

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-4")
	if code != 0 {
		t.Fatalf("exit %d\n%s%s", code, out, errOut)
	}
	for _, want := range []string{"passed rebase (", "passed sync (", "passed narrative (", "passed check (", "passed tier unit ("} {
		if !strings.Contains(out, "\n"+want) && !strings.HasPrefix(out, want) {
			t.Errorf("text lacks a line %q:\n%s", want, out)
		}
	}
	if got := lastLine(out); got != "verify: S-004 passed every step" {
		t.Errorf("last line %q:\n%s", got, out)
	}

	// from the story's worktree, as data
	rep := verifyReport(t, worktree, 0, "S-004")
	want := "rebase:passed,sync:passed,narrative:passed,check:passed,tier unit:passed"
	if got := stepStates(rep); got != want || !rep.Passed || rep.StoppedAt != "" || rep.Story != "S-004" || rep.Commit != "0123abc4567" || rep.Base != "main" {
		t.Errorf("report %+v, steps %s, want %s", rep, got, want)
	}
	if strings.Join(rep.Paths, ",") != "docs/README.md" || rep.Steps[4].Duration == "" || rep.Steps[4].ExitCode == nil || *rep.Steps[4].ExitCode != 0 {
		t.Errorf("report %+v, unit %+v", rep, rep.Steps[4])
	}
}

func TestVerifyStopsAtAFailingStepWithItsFindingsAndExitsOne(t *testing.T) {
	root, _ := verifyFixture(t, "echo \"docs/README.md is broken\" >&2; exit 3")

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-004")
	if code != 1 {
		t.Fatalf("a failing tier: exit %d, want 1\n%s%s", code, out, errOut)
	}
	if !strings.Contains(out, "\nfailed tier unit (") || !strings.Contains(out, "\n    ") || !strings.Contains(out, "docs/README.md is broken") {
		t.Errorf("text should give the failing tier and its finding indented under it:\n%s", out)
	}
	if got := lastLine(out); got != "verify: S-004 stopped at unit (exit 3)" {
		t.Errorf("last line %q:\n%s", got, out)
	}
	rep := verifyReport(t, root, 1, "S-004")
	unit := rep.Steps[len(rep.Steps)-1]
	if rep.Passed || rep.StoppedAt != "unit" || unit.State != verify.Failed || len(unit.Findings) == 0 || unit.ExitCode == nil || *unit.ExitCode != 3 {
		t.Errorf("report %+v, unit %+v", rep, unit)
	}

	// a step before the tiers: the narrative's Current state is empty
	narrative := filepath.Join(root, "wip", "agents", "S-004.md")
	data, err := os.ReadFile(narrative)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(narrative, []byte(strings.Replace(string(data), "## Current state\ns\n", "## Current state\n- \n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runVerify(t, root, verifyGit{}, "verify", "S-004")
	if code != 1 {
		t.Fatalf("an empty narrative: exit %d, want 1\n%s%s", code, out, errOut)
	}
	for _, want := range []string{"\nfailed narrative (", "\n    wip/agents/S-004.md:", "## Current state is empty", "\nnot-reached check\n", "\nnot-reached tier unit\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("text lacks %q:\n%s", want, out)
		}
	}
	if got := lastLine(out); got != "verify: S-004 stopped at narrative (failed)" {
		t.Errorf("last line %q:\n%s", got, out)
	}
}

// S-0342: a tier that another selected tier covers is skipped, its line
// naming the tier that covers it, in the text, in --json, and in --last.
func TestVerifySkipsATierThatASelectedTierCovers(t *testing.T) {
	root, worktree := verifyFixture(t, "exit 0")
	manifest := filepath.Join(worktree, "system-flow.yaml")
	data, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	whole := "  - name: whole\n    command: [sh, -c, 'exit 0']\n    paths: [\"docs/**\"]\n    covers: [unit]\n"
	if err := os.WriteFile(manifest, append(data, whole...), 0o644); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-004")
	if code != 0 || !strings.Contains(out, "\nskipped tier unit: covered by whole\npassed tier whole (") || lastLine(out) != "verify: S-004 passed every step" {
		t.Errorf("exit %d, want 0 and unit skipped, covered by whole\n%s%s", code, out, errOut)
	}
	for _, args := range [][]string{{"S-004"}, {"S-004", "--last"}} {
		rep := verifyReport(t, root, 0, args...)
		want := "rebase:passed,sync:passed,narrative:passed,check:passed,tier unit:skipped,tier whole:passed"
		if got := stepStates(rep); got != want || !rep.Passed || rep.Steps[4].CoveredBy != "whole" || rep.Steps[4].Duration != "" || rep.Steps[5].CoveredBy != "" {
			t.Errorf("verify %v: steps %s, want %s; unit %+v, whole %+v", args, got, want, rep.Steps[4], rep.Steps[5])
		}
	}
}

// S-0249: the check's findings outside the story are notes that do not fail
// it, given after the steps and in --json.
func TestVerifyGivesTheNotesOutsideTheStory(t *testing.T) {
	root, _ := verifyFixture(t, "exit 0")
	note := "    wip/archive/kanban/stories/S-002-two.md:6: warning: story.unaccepted: S-002 is done but its branch story/S-002 was never merged; merge or delete it\n"

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-004")
	if code != 0 || !strings.Contains(out, "\noutside S-004, notes that do not fail it:\n") || !strings.Contains(out, note) {
		t.Fatalf("exit %d, text should give the note %q:\n%s%s", code, note, out, errOut)
	}
	if strings.Index(out, note) > strings.Index(out, "verify: S-004 passed every step") {
		t.Errorf("the notes should come before the last line:\n%s", out)
	}

	rep := verifyReport(t, root, 0, "S-004")
	found := false
	for _, n := range rep.Notes {
		found = found || n.Rule == "story.unaccepted" && n.Level == "warning" && n.Path == "wip/archive/kanban/stories/S-002-two.md"
	}
	if !found {
		t.Errorf("notes %+v lack S-002's branch", rep.Notes)
	}
}

// --record-issues records the notes in design/issues of the story's
// worktree through flai check's recorder, once per story and notes.
func TestVerifyRecordIssuesRecordsTheNotesInTheWorktree(t *testing.T) {
	root, worktree := verifyFixture(t, "exit 0")
	finding := "`wip/archive/kanban/stories/S-002-two.md`: S-002 is done but its branch story/S-002 was never merged; merge or delete it"

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-004", "--record-issues")
	if code != 0 {
		t.Fatalf("exit %d\n%s%s", code, out, errOut)
	}
	var issue string
	for _, name := range issueFiles(t, worktree) {
		if strings.Contains(name, "story-unaccepted") {
			issue = filepath.Join(worktree, "design", "issues", name)
		}
	}
	if issue == "" {
		t.Fatalf("no story.unaccepted issue in the worktree: %v\n%s", issueFiles(t, worktree), out)
	}
	doc, err := os.ReadFile(issue)
	if err != nil || !strings.Contains(string(doc), finding+"\n") || !strings.Contains(string(doc), "Story: S-0004.\n") {
		t.Errorf("issue %v:\n%s", err, doc)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("the main checkout should have no issue written, got %v", names)
	}
	if !strings.Contains(out, "\nrecorded story.unaccepted outside S-004 in I-") || !strings.Contains(out, " (opened)\n") ||
		lastLine(out) != "verify: S-004 passed every step" {
		t.Errorf("text should say what it recorded, then the outcome:\n%s", out)
	}

	out, errOut, code = runVerify(t, root, verifyGit{}, "verify", "S-004", "--record-issues", "--json")
	var res struct {
		Passed   bool `json:"passed"`
		Recorded []struct {
			Rule, Issue, Outcome string
		} `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("again: exit %d %v\n%s%s", code, err, out, errOut)
	}
	already := false
	for _, r := range res.Recorded {
		already = already || r.Rule == "story.unaccepted" && r.Outcome == "already recorded"
	}
	if !res.Passed || !already {
		t.Errorf("an identical run should record nothing again: %s", out)
	}
	if again, _ := os.ReadFile(issue); string(again) != string(doc) {
		t.Errorf("an identical run changed the issue:\n%s", again)
	}
}

// --last prints the stored report and runs nothing, or says there is none,
// and exits 0 whatever the report holds.
func TestVerifyLastPrintsTheStoredReport(t *testing.T) {
	root, _ := verifyFixture(t, "exit 3")

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-4", "--last")
	if code != 0 || out != "verify: S-004 has no stored result; run flai verify S-004\n" {
		t.Errorf("none stored: exit %d %q %s", code, out, errOut)
	}
	out, errOut, code = runVerify(t, root, verifyGit{}, "verify", "S-004", "--last", "--json")
	if code != 0 || out != "null\n" {
		t.Errorf("none stored, as data: exit %d %q %s", code, out, errOut)
	}

	ran := verifyReport(t, root, 1, "S-004")
	// what --last reads is what was stored, not a new run: no git answers
	out, errOut, code = runVerify(t, root, gitScript{}, "verify", "S-004", "--last")
	if code != 0 || !strings.HasPrefix(out, "last verified ") || lastLine(out) != "verify: S-004 stopped at unit (exit 3)" {
		t.Errorf("stored: exit %d\n%s%s", code, out, errOut)
	}
	out, errOut, code = runVerify(t, root, gitScript{}, "verify", "S-004", "--last", "--json")
	var last verify.Report
	if err := json.Unmarshal([]byte(out), &last); err != nil || code != 0 {
		t.Fatalf("stored, as data: exit %d %v\n%s%s", code, err, out, errOut)
	}
	if last.Passed || last.StoppedAt != "unit" || !last.RanAt.Equal(ran.RanAt) || stepStates(last) != stepStates(ran) {
		t.Errorf("stored %+v, ran %+v", last, ran)
	}

	_, errOut, code = runVerify(t, root, verifyGit{}, "verify", "S-004", "--last", "--record-issues")
	if code != exitVerifyUnusable || !strings.Contains(errOut, "drop --record-issues") {
		t.Errorf("--last with --record-issues: exit %d %s", code, errOut)
	}
}

// roleTierEnv, set to 1, has the test binary run as roleTier's command.
const roleTierEnv = "FLAI_TEST_ROLE_TIER"

// roleTier is a tier script that passes only under the verify role and when
// a story's move to ready, made in process as TH-0260's go-test tier made
// it, is not refused: it runs TestRoleTierMovesAStoryToReady in this test
// binary.
func roleTier(t *testing.T) string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return `test "$FLAI_ROLE" = verify || { echo "FLAI_ROLE is $FLAI_ROLE, not verify"; exit 1; }; ` +
		roleTierEnv + `=1 exec "` + exe + `" -test.run=^TestRoleTierMovesAStoryToReady$`
}

// TestRoleTierMovesAStoryToReady is roleTier's command, not a test of its
// own: in a tier it makes a project and moves its story to ready in process,
// under the role the tier gave it, and fails when flai refuses the move.
func TestRoleTierMovesAStoryToReady(t *testing.T) {
	if os.Getenv(roleTierEnv) != "1" {
		t.Skip("roleTier's command; it runs only in a tier")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Moved", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-moved.md")
	story, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(story), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code != 0 {
		t.Fatalf("flai move S-0001 ready under FLAI_ROLE=%s: exit %d %s", os.Getenv("FLAI_ROLE"), code, errOut)
	}
}

// S-0311: flai verify run as the orchestrator runs its tiers under the
// verify role, so a tier whose tests make flai's writes is not refused as
// the orchestrator (TH-0260).
func TestVerifyRunsItsTiersUnderTheVerifyRoleForTheOrchestrator(t *testing.T) {
	root, _ := verifyFixture(t, roleTier(t))
	t.Setenv("FLAI_ROLE", "orchestrate")

	out, errOut, code := runVerify(t, root, verifyGit{}, "verify", "S-004")
	if code != 0 || lastLine(out) != "verify: S-004 passed every step" {
		t.Errorf("exit %d\n%s%s", code, out, errOut)
	}
}

func TestVerifyRefusesAStoryWithNoWorktree(t *testing.T) {
	root, worktree := verifyFixture(t, "exit 0")
	if err := os.RemoveAll(worktree); err != nil {
		t.Fatal(err)
	}
	cases := map[string]struct {
		args []string
		want string
	}{
		"no worktree":      {[]string{"verify", "S-004"}, "it has no worktree at " + worktree + "; open it with flai stream open S-004"},
		"not a story":      {[]string{"verify", "T-003"}, "it is a task; name a story"},
		"no such story":    {[]string{"verify", "S-099"}, "name a story of this project"},
		"no story named":   {[]string{"verify"}, "name the story to verify"},
		"a cap of nothing": {[]string{"verify", "S-004", "--max", "0"}, "--max 0 is not a count of findings"},
		"sync-only last":   {[]string{"verify", "S-004", "--sync-only", "--last"}, "--sync-only runs the rebase and sync steps and --last runs nothing; drop one of them"},
		"sync-only record": {[]string{"verify", "S-004", "--sync-only", "--record-issues"}, "--sync-only runs no check, so it has no notes to record; drop --record-issues"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, errOut, code := runVerify(t, root, verifyGit{}, tc.args...)
			if code != exitVerifyUnusable || out != "" || !strings.Contains(errOut, tc.want) {
				t.Errorf("exit %d, want %d\nout %q\nerr %q, want %q", code, exitVerifyUnusable, out, errOut, tc.want)
			}
		})
	}
	if _, err := os.Stat(verify.ReportPath(openProject(t, root), "S-004")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a refused run should store no report: %v", err)
	}
}

// I-0119, ADR-0135: the close-out's last check, flai verify --sync-only,
// passes over commits flai made on main's wip/ while the close-out ran, and
// says so under the sync step; a commit outside wip/ still fails it, naming
// flai stream sync. Neither touches the stored report.
func TestVerifySyncOnlyPassesOverWipCommitsAndStoresNoReport(t *testing.T) {
	root, _ := verifyFixture(t, "exit 3")
	report := verify.ReportPath(openProject(t, root), "S-004")
	replan := verifyGit{commits: "a1a1a1a", paths: "wip/kanban/board.md"}
	note := "\npassed sync (0s)\n    passed over 1 commit of main that changes only wip/ paths the branch does not change: a1a1a1a\n"

	out, errOut, code := runVerify(t, root, replan, "verify", "S-004", "--sync-only")
	if code != 0 || !strings.HasPrefix(out, "passed rebase (") || !strings.Contains(out, note) || lastLine(out) != "verify: S-004 passed every step" {
		t.Errorf("a branch behind main by a replan: exit %d, want 0 and the note %q\n%s%s", code, note, out, errOut)
	}
	if strings.Contains(out, "narrative") || strings.Contains(out, "tier unit") {
		t.Errorf("--sync-only ran more than the rebase and sync steps:\n%s", out)
	}
	if _, err := os.Stat(report); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("--sync-only stored a report where there was none: %v", err)
	}

	// a full run, whose unit tier fails, stores the report --sync-only keeps
	if out, errOut, code = runVerify(t, root, replan, "verify", "S-004"); code != 1 || !strings.Contains(out, note) {
		t.Fatalf("a full run: exit %d, want 1 and the note %q\n%s%s", code, note, out, errOut)
	}
	stored, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}

	out, errOut, code = runVerify(t, root, replan, "--json", "verify", "S-004", "--sync-only")
	var rep verify.Report
	if err := json.Unmarshal([]byte(out), &rep); err != nil || code != 0 {
		t.Fatalf("as data: exit %d %v\n%s%s", code, err, out, errOut)
	}
	if !rep.Passed || stepStates(rep) != "rebase:passed,sync:passed" || !strings.HasSuffix(rep.Steps[1].Note, ": a1a1a1a") {
		t.Errorf("report %+v", rep)
	}

	release := verifyGit{commits: "b2b2b2b\na1a1a1a", paths: "flai/go.mod\nwip/kanban/board.md"}
	out, errOut, code = runVerify(t, root, release, "verify", "S-004", "--sync-only")
	if code != 1 || !strings.Contains(out, "\nfailed sync (") || !strings.Contains(out, "\n    ") || !strings.Contains(out, "outside wip/: flai/go.mod; run flai stream sync S-004") ||
		strings.Contains(out, "passed over") || lastLine(out) != "verify: S-004 stopped at sync (failed)" {
		t.Errorf("a branch behind main by a release: exit %d, want 1 naming flai stream sync\n%s%s", code, out, errOut)
	}
	if again, err := os.ReadFile(report); err != nil || !bytes.Equal(again, stored) {
		t.Errorf("--sync-only changed the stored report: %v\n%s\nwas\n%s", err, again, stored)
	}
}

// A step's note is printed indented under its line, on one line, as a
// failing step's findings are; --last prints it with the same text.
func TestVerifyTextPrintsAStepsNoteUnderIt(t *testing.T) {
	rep := verify.Report{Story: "S-004", Passed: true, Steps: []verify.Step{
		{Name: verify.StepRebase, State: verify.Passed, Duration: "0s"},
		{Name: verify.StepSync, State: verify.Passed, Duration: "1ms", Note: "passed over 2 commits of main that change only wip/\npaths the branch does not change: b2b2b2b, a1a1a1a"},
		{Name: verify.StepNarrative, State: verify.Passed, Duration: "0s"},
	}}
	want := "passed rebase (0s)\n" +
		"passed sync (1ms)\n" +
		"    passed over 2 commits of main that change only wip/ paths the branch does not change: b2b2b2b, a1a1a1a\n" +
		"passed narrative (0s)\n" +
		"verify: S-004 passed every step\n"
	if got := verifyText(rep, nil, false, verify.DefaultMax); got != want {
		t.Errorf("text\n%s\nwant\n%s", got, want)
	}
}

// resumeGit is verifyGit with the main branch's commit answered as well, so
// a story's run can resume from its last record.
type resumeGit struct{ verifyGit }

func (g resumeGit) Run(dir, name string, args ...string) (string, error) {
	if strings.Join(args, " ") == "rev-parse --verify --quiet main^{commit}" {
		return "89abcdef0123\n", nil
	}
	return g.verifyGit.Run(dir, name, args...)
}

func (g resumeGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

// resumeFixture is verifyFixture with two tiers for docs/: unit, which
// passes, then lint, which fails with exit 4 until the file marker exists.
func resumeFixture(t *testing.T) (root, marker string) {
	t.Helper()
	root, worktree := verifyFixture(t, "exit 0")
	marker = filepath.Join(t.TempDir(), "lint-passes")
	path := filepath.Join(worktree, "system-flow.yaml")
	manifest, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lint := "  - name: lint\n    command: [sh, -c, 'test -f \"" + marker + "\" || { echo \"docs/README.md does not lint\" >&2; exit 4; }', lint, \"{files}\"]\n    paths: [\"docs/**\"]\n"
	if err := os.WriteFile(path, []byte(string(manifest)+lint), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, marker
}

// resumeReport runs flai verify --json in dir with resumeGit and decodes its
// answer, and answers it raw too.
func resumeReport(t *testing.T, dir string, wantCode int, args ...string) (verify.Report, string) {
	t.Helper()
	out, errOut, code := runVerify(t, dir, resumeGit{}, append([]string{"verify", "--json"}, args...)...)
	if code != wantCode {
		t.Fatalf("flai verify %v: exit %d, want %d\n%s%s", args, code, wantCode, out, errOut)
	}
	var rep verify.Report
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatalf("flai verify %v: not JSON: %v\n%s", args, err, out)
	}
	return rep, out
}

// S-0341: a run at the same head, base, and tiers as the story's last
// resumes at the tier that failed, prints the tiers it reused with the time
// of the run they come from, and says so on its last line.
func TestVerifyResumesAtTheTierThatFailedAndSaysWhatItReused(t *testing.T) {
	root, marker := resumeFixture(t)

	out, errOut, code := runVerify(t, root, resumeGit{}, "verify", "S-004")
	if code != 1 || lastLine(out) != "verify: S-004 stopped at lint (exit 4)" {
		t.Fatalf("first run: exit %d\n%s%s", code, out, errOut)
	}
	// the first run began an hour before the clock runVerify gives
	project := openProject(t, root)
	first, ok, err := verify.LastReport(project, "S-004")
	if err != nil || !ok {
		t.Fatalf("no stored report: %v", err)
	}
	first.RanAt = time.Date(2026, 9, 17, 0, 0, 0, 0, time.UTC)
	if err := verify.SaveReport(project, first); err != nil {
		t.Fatal(err)
	}

	// lint fails again: unit is reused, and the stop says so
	out, errOut, code = runVerify(t, root, resumeGit{}, "verify", "S-004")
	if code != 1 {
		t.Fatalf("second run: exit %d\n%s%s", code, out, errOut)
	}
	for _, want := range []string{"\npassed check (", "\nreused tier unit (from 2026-09-17T00:00:00Z)\n", "\nfailed tier lint (", "docs/README.md does not lint"} {
		if !strings.Contains(out, want) {
			t.Errorf("second run's text lacks %q:\n%s", want, out)
		}
	}
	if got := lastLine(out); got != "verify: S-004 stopped at lint (exit 4) (reused 1 tier from 2026-09-17T00:00:00Z)" {
		t.Errorf("second run's last line %q:\n%s", got, out)
	}

	// lint passes: unit is still reused from the first run
	if err := os.WriteFile(marker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runVerify(t, root, resumeGit{}, "verify", "S-004")
	if code != 0 || !strings.Contains(out, "\nreused tier unit (from 2026-09-17T00:00:00Z)\n") || !strings.Contains(out, "\npassed tier lint (") {
		t.Fatalf("third run: exit %d\n%s%s", code, out, errOut)
	}
	if got := lastLine(out); got != "verify: S-004 passed every step (reused 1 tier from 2026-09-17T00:00:00Z)" {
		t.Errorf("third run's last line %q:\n%s", got, out)
	}

	// every tier passed: both are reused, and --json carries it
	out, errOut, code = runVerify(t, root, resumeGit{}, "verify", "S-004")
	if got := lastLine(out); code != 0 || got != "verify: S-004 passed every step (reused 2 tiers from 2026-09-17T00:00:00Z)" {
		t.Errorf("fourth run: exit %d, last line %q\n%s%s", code, got, out, errOut)
	}
	rep, raw := resumeReport(t, root, 0, "S-004")
	if want := "rebase:passed,sync:passed,narrative:passed,check:passed,tier unit:reused,tier lint:reused"; stepStates(rep) != want {
		t.Errorf("steps %s, want %s", stepStates(rep), want)
	}
	if !strings.Contains(raw, `"state": "reused"`) || !strings.Contains(raw, `"reused_from": "2026-09-17T00:00:00Z"`) || !strings.Contains(raw, `"reused_from": "2026-09-17T01:00:00Z"`) {
		t.Errorf("--json should carry the reused states and when each tier ran:\n%s", raw)
	}
	if rep.BaseCommit != "89abcdef0123" {
		t.Errorf("base commit %q", rep.BaseCommit)
	}
}

// --fresh reaches the story's run: it runs every selected tier where a run
// without it reuses them.
func TestVerifyFreshRunsEveryTier(t *testing.T) {
	root, marker := resumeFixture(t)
	if err := os.WriteFile(marker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if rep, _ := resumeReport(t, root, 0, "S-004"); stepStates(rep) != "rebase:passed,sync:passed,narrative:passed,check:passed,tier unit:passed,tier lint:passed" {
		t.Fatalf("first run: steps %s", stepStates(rep))
	}

	// lint would fail now, but a run without --fresh does not run it
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	if rep, _ := resumeReport(t, root, 0, "S-004"); stepStates(rep) != "rebase:passed,sync:passed,narrative:passed,check:passed,tier unit:reused,tier lint:reused" {
		t.Errorf("without --fresh: steps %s", stepStates(rep))
	}
	out, errOut, code := runVerify(t, root, resumeGit{}, "verify", "S-004", "--fresh")
	if code != 1 || !strings.Contains(out, "\npassed tier unit (") || strings.Contains(out, "reused") {
		t.Errorf("--fresh: exit %d, want 1 with every tier run\n%s%s", code, out, errOut)
	}
	if got := lastLine(out); got != "verify: S-004 stopped at lint (exit 4)" {
		t.Errorf("--fresh: last line %q:\n%s", got, out)
	}
}

// taggedGit is verifyGit in a project whose history holds the flai release
// tags given.
type taggedGit struct {
	verifyGit
	tags string
}

func (g taggedGit) Run(dir, name string, args ...string) (string, error) {
	if strings.Join(args, " ") == "tag --merged HEAD --list flai/v*" {
		return g.tags, nil
	}
	return g.verifyGit.Run(dir, name, args...)
}

func (g taggedGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

// I-0141, ADR-0158: flai verify refuses tiers its flai does not take, in the
// story's worktree, naming both versions and the upgrade when the history
// holds a newer flai release, and the manifest's fix when it does not.
func TestVerifySaysTheFlaiIsOlderThanTheProjectWhenTheTiersRefuse(t *testing.T) {
	root, worktree := verifyFixture(t, "exit 0")
	path := filepath.Join(worktree, "system-flow.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(data), `"{files}"]`, `"{dependents}"]`, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	was := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = was })
	// the log line quotes the tier's name, so the problem is matched without it
	const problem = "argument {dependents} is not a placeholder flai fills"
	// each case its own running version, as FlaiOutdated keeps an answer per
	// root and version for a minute
	for _, c := range []struct {
		version string
		want    []string
		not     string
	}{
		{"1.40.1", []string{"verify S-004: ", "tests[0].command", problem, "this flai is 1.40.1, older than flai 1.49.0 in the project's history", "upgrade it with flai host upgrade"}, "flai manifest set tests="},
		{"1.49.0", []string{"verify S-004: ", "tests[0].command", problem, "flai manifest set tests="}, "upgrade it"},
	} {
		buildinfo.Version = c.version
		out, errOut, code := runVerify(t, root, taggedGit{tags: "flai/v1.40.1\nflai/v1.49.0\n"}, "verify", "S-004")
		if code != exitVerifyUnusable || out != "" {
			t.Errorf("flai %s: exit %d, want %d\nout %q\nerr %q", c.version, code, exitVerifyUnusable, out, errOut)
		}
		for _, want := range c.want {
			if !strings.Contains(errOut, want) {
				t.Errorf("flai %s: %q, want it to say %q", c.version, errOut, want)
			}
		}
		if strings.Contains(errOut, c.not) {
			t.Errorf("flai %s: %q, want it not to say %q", c.version, errOut, c.not)
		}
	}
}
