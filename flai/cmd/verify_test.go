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

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// verifyGit answers the git commands flai verify runs for a story whose
// branch changed docs/README.md, contains main, has nothing uncommitted, and
// has no rebase left, and refuses every other, as git refuses a repository
// it is not run in.
type verifyGit struct{}

func (verifyGit) Run(_, _ string, args ...string) (string, error) {
	switch strings.Join(args, " ") {
	case "rev-parse --abbrev-ref HEAD":
		return "main", nil
	case "rev-parse HEAD":
		return "0123abc4567\n", nil
	case "diff --name-only --diff-filter=d main...HEAD":
		return "docs/README.md\n", nil
	case "status --porcelain --untracked-files=all":
		return "", nil
	case "rev-list --count HEAD..main":
		return "0", nil
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
