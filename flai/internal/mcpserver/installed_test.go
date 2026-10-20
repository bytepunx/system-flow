package mcpserver

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/taskdone"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeFlai is an installed flai written as a POSIX sh script: flai version
// --json prints version and exits versionExit; flai task done records its
// arguments, working directory, FLAI_AGENT, and FLAI_SESSION beside it,
// logs a line on standard error, prints result, and exits resultExit.
type fakeFlai struct {
	path string
	dir  string // where it records
}

func newFakeFlai(t *testing.T, version string, versionExit int, result string, resultExit int) fakeFlai {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake installed flai is a POSIX sh script")
	}
	dir := t.TempDir()
	writeIn(t, dir, "version.json", version)
	writeIn(t, dir, "result.json", result)
	script := "#!/bin/sh\nrec='" + dir + "'\n" +
		"if [ \"$1\" = version ]; then\n" +
		"  cat \"$rec/version.json\"\n" +
		"  echo 'level=WARN msg=\"version failed\"' >&2\n" +
		"  exit " + strconv.Itoa(versionExit) + "\n" +
		"fi\n" +
		"printf '%s\\n' \"$@\" > \"$rec/args\"\n" +
		"pwd -P > \"$rec/cwd\"\n" +
		"printf '%s\\n' \"$FLAI_AGENT\" > \"$rec/agent\"\n" +
		"printf '%s\\n' \"$FLAI_SESSION\" > \"$rec/session\"\n" +
		"echo 'level=INFO msg=\"task close stopped\"' >&2\n" +
		"cat \"$rec/result.json\"\n" +
		"exit " + strconv.Itoa(resultExit) + "\n"
	path := filepath.Join(dir, "flai")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return fakeFlai{path: path, dir: dir}
}

// recorded is what the fake's task done wrote to name, "" when it did not run.
func (f fakeFlai) recorded(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.dir, name))
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(data))
}

const (
	closedResult  = `{"task":"T-003","story":"S-004","commit":{"hash":"abc","subject":"docs: [S-004] T-003 the readme","paths":["README.md"]},"left":[],"move":{"id":"T-003","state":"done","skipped":false,"followed":[],"warnings":[]},"stopped":""}`
	stoppedResult = `{"task":"T-003","story":"S-004","commit":null,"left":[],"stopped":"sync","error":"the rebase stopped on conflicts in README.md"}`
)

// outdatedProject is I-0147's case: S-004's worktree with T-003's change in
// it, on a main branch tagged flai/v1.10.0, served by flai 1.9.3 as
// agent-S-004 in session s1, with the installed flai at installed. With
// rules above 0, the tagged commit carries that task-done rules number
// (ADR-0180); with 0 it carries none.
func outdatedProject(t *testing.T, installed string, rules int) (*fixture, string) {
	t.Helper()
	repo, wt := closeProject(t)
	if rules > 0 {
		writeIn(t, repo.Root, buildinfo.TaskDoneRulesPath, strconv.Itoa(rules)+"\n")
		gitIn(t, repo.Root, "add", buildinfo.TaskDoneRulesPath)
		gitIn(t, repo.Root, "commit", "-q", "-m", "flai 1.10.0")
	}
	gitIn(t, repo.Root, "tag", "flai/v1.10.0")
	writeIn(t, wt, "README.md", "# good\n\nChanged on the story's branch.\n")
	cs := connectServer(t, Options{Repo: repo, Runner: execx.System{}, Agent: "agent-S-004", Version: "1.9.3", Installed: installed, Now: func() time.Time { return closeClock }})
	return &fixture{repo: repo, cs: cs}, wt
}

// I-0147, ADR-0160: a server a minor version behind the newest release,
// whose installed flai has caught up, closes the task with the installed
// flai: run in the main checkout as the server's agent and session, with the
// message and the log, and answered with closed_by naming it.
func TestTaskDoneClosesWithTheInstalledFlaiOnceItHasCaughtUp(t *testing.T) {
	fake := newFakeFlai(t, `{"version":"1.10.0","commit":"c","date":"d","go":"go"}`, 0, closedResult, 0)
	f, _ := outdatedProject(t, fake.path, 0)

	out, failed := f.call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme", "log": "the readme"})
	if failed != "" {
		t.Fatalf("a tool error: %s", failed)
	}
	if out["closed_by"] != "flai 1.10.0 at "+fake.path || out["task"] != "T-003" || out["stopped"] != "" || field(out, "move.state") != workitem.Done || field(out, "commit.subject") != "docs: [S-004] T-003 the readme" {
		t.Errorf("answered %v", out)
	}
	if got, want := fake.recorded(t, "args"), "task\ndone\nT-003\n--json\n-m\ndocs: [S-004] T-003 the readme\n--log\nthe readme"; got != want {
		t.Errorf("the installed flai ran with\n%s\nwant\n%s", got, want)
	}
	root, err := filepath.EvalSymlinks(f.repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	if got := fake.recorded(t, "cwd"); got != root {
		t.Errorf("the installed flai ran in %s, want the main checkout %s", got, root)
	}
	if agent, session := fake.recorded(t, "agent"), fake.recorded(t, "session"); agent != "agent-S-004" || session != "s1" {
		t.Errorf("the installed flai ran as %q in session %q", agent, session)
	}
}

// ADR-0160: a close the installed flai stops at a step exits non-zero and
// prints its result, which is the answer, not a tool error.
func TestTaskDoneAnswersTheInstalledFlaisStop(t *testing.T) {
	fake := newFakeFlai(t, `{"version":"1.10.1"}`, 0, stoppedResult, 3)
	f, _ := outdatedProject(t, fake.path, 0)

	out, failed := f.call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
	if failed != "" {
		t.Fatalf("a stop is a tool error: %s", failed)
	}
	if out["stopped"] != taskdone.StepSync || out["error"] != "the rebase stopped on conflicts in README.md" || out["closed_by"] != "flai 1.10.1 at "+fake.path {
		t.Errorf("answered %v", out)
	}
}

// ADR-0160: an installed flai behind too, or one whose version cannot be
// read, does not close the task: task_done refuses as ADR-0142 does, adds
// the installed flai's version or why it could not be read and to call
// again once the host is upgraded, and changes nothing.
func TestTaskDoneRefusesWhenTheInstalledFlaiHasNotCaughtUp(t *testing.T) {
	for _, tc := range []struct {
		name, version string
		exit          int
		want          []string
	}{
		{"behind too", `{"version":"1.9.4"}`, 0, []string{"is 1.9.4, behind too"}},
		{"version fails", `{}`, 1, []string{"could not be read", "exit status 1", "version failed"}},
		{"version not JSON", "flai 1.10.0", 0, []string{"could not be read", "printed no version", "flai 1.10.0"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := newFakeFlai(t, tc.version, tc.exit, closedResult, 0)
			f, wt := outdatedProject(t, fake.path, 0)
			head := gitIn(t, wt, "rev-parse", "HEAD")

			out, failed := f.call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
			for _, want := range append([]string{"this flai is 1.9.3, a minor version behind flai 1.10.0", "the installed flai at " + fake.path, "call task_done again once the host is upgraded"}, tc.want...) {
				if !strings.Contains(failed, want) {
					t.Errorf("answered %v, failed %q, want %q in it", out, failed, want)
				}
			}
			if args := fake.recorded(t, "args"); args != "" {
				t.Errorf("the installed flai's task done ran: %s", args)
			}
			if st := t003(t, f.repo); st != workitem.InProgress {
				t.Errorf("T-003 is %s", st)
			}
			if gitIn(t, wt, "rev-parse", "HEAD") != head || gitIn(t, wt, "diff", "--name-only") != "README.md" {
				t.Error("the worktree changed")
			}
		})
	}
}

// I-0157, ADR-0180: a server a minor version behind a release that kept its
// task-done rules number is not refused, so it closes the task itself and
// the installed flai, behind too, is not asked.
func TestTaskDoneClosesAMinorVersionBehindAReleaseThatKeptTheRules(t *testing.T) {
	fake := newFakeFlai(t, `{"version":"1.9.3"}`, 0, closedResult, 0)
	f, _ := outdatedProject(t, fake.path, buildinfo.TaskDoneRules)

	out, failed := f.call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
	if failed != "" {
		t.Fatalf("a tool error: %s", failed)
	}
	if out["closed_by"] != nil || out["stopped"] != "" || field(out, "move.state") != workitem.Done || field(out, "commit.subject") != "docs: [S-004] T-003 the readme" {
		t.Errorf("answered %v", out)
	}
	if args := fake.recorded(t, "args"); args != "" {
		t.Errorf("the installed flai's task done ran: %s", args)
	}
}

// I-0157, ADR-0180: under a release that raised the task-done rules number,
// which refuses this server's flai, an installed flai that carries that
// number or a higher one has caught up and closes the task though it is a
// minor version behind; one that prints no number keeps ADR-0142's minor
// rule, so it closes within a patch of the release.
func TestTaskDoneClosesWithAnInstalledFlaiThatCarriesTheNewestRules(t *testing.T) {
	newest := buildinfo.TaskDoneRules + 1
	for _, tc := range []struct {
		name, version, closedBy string
	}{
		{"the newest rules", `{"version":"1.9.4","task_done_rules":` + strconv.Itoa(newest) + `}`, "flai 1.9.4"},
		{"higher rules", `{"version":"1.9.4","task_done_rules":` + strconv.Itoa(newest+1) + `}`, "flai 1.9.4"},
		{"no rules, a patch ahead", `{"version":"1.10.1"}`, "flai 1.10.1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := newFakeFlai(t, tc.version, 0, closedResult, 0)
			f, _ := outdatedProject(t, fake.path, newest)

			out, failed := f.call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
			if failed != "" {
				t.Fatalf("a tool error: %s", failed)
			}
			if out["closed_by"] != tc.closedBy+" at "+fake.path || field(out, "move.state") != workitem.Done {
				t.Errorf("answered %v", out)
			}
			if got := fake.recorded(t, "args"); !strings.HasPrefix(got, "task\ndone\nT-003\n--json") {
				t.Errorf("the installed flai ran with %q", got)
			}
		})
	}
}

// ADR-0180: under a release that raised the task-done rules number, an
// installed flai below that number is behind too, even at the release's own
// version, and one that prints no number is behind a minor version; the
// refusal names the server's rules and the installed flai's version and
// number, and nothing changes.
func TestTaskDoneRefusesAnInstalledFlaiBelowTheNewestRules(t *testing.T) {
	newest := buildinfo.TaskDoneRules + 1
	for _, tc := range []struct {
		name, version, says string
	}{
		{"lower rules", `{"version":"1.10.0","task_done_rules":` + strconv.Itoa(buildinfo.TaskDoneRules) + `}`, "is 1.10.0 with task-done rules " + strconv.Itoa(buildinfo.TaskDoneRules) + ", behind too"},
		{"no rules, a minor behind", `{"version":"1.9.4"}`, "is 1.9.4, behind too"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := newFakeFlai(t, tc.version, 0, closedResult, 0)
			f, wt := outdatedProject(t, fake.path, newest)
			head := gitIn(t, wt, "rev-parse", "HEAD")

			out, failed := f.call(t, "task_done", map[string]any{"task": "T-003", "message": "docs: [S-004] T-003 the readme"})
			raised := "this flai is 1.9.3, older than flai 1.10.0 in the project's history, which raised task done's rules from " + strconv.Itoa(buildinfo.TaskDoneRules) + " to " + strconv.Itoa(newest)
			for _, want := range []string{raised, "the installed flai at " + fake.path + " " + tc.says, "call task_done again once the host is upgraded", "ADR-0180"} {
				if !strings.Contains(failed, want) {
					t.Errorf("answered %v, failed %q, want %q in it", out, failed, want)
				}
			}
			if args := fake.recorded(t, "args"); args != "" {
				t.Errorf("the installed flai's task done ran: %s", args)
			}
			if st := t003(t, f.repo); st != workitem.InProgress {
				t.Errorf("T-003 is %s", st)
			}
			if gitIn(t, wt, "rev-parse", "HEAD") != head || gitIn(t, wt, "diff", "--name-only") != "README.md" {
				t.Error("the worktree changed")
			}
		})
	}
}

// ADR-0160: output of the installed flai's task done that is not a result is
// the tool's error, naming the command and quoting what it said.
func TestInstalledFlaiOutputThatIsNotAResultIsAnError(t *testing.T) {
	fake := newFakeFlai(t, `{"version":"1.10.0"}`, 0, "Error: T-003 not found", 1)
	root := t.TempDir()
	f := installedFlai{path: fake.path, dir: root, env: append(os.Environ(), "FLAI_AGENT=a", "FLAI_SESSION=s")}

	res, err := f.taskDone(context.Background(), TaskDoneIn{Task: "T-003"})
	if err == nil {
		t.Fatalf("answered %+v", res)
	}
	for _, want := range []string{"the installed flai printed no result closing T-003", fake.path + " task done T-003 --json: exit status 1", "printed \"Error: T-003 not found\"", "task close stopped"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error lacks %q: %v", want, err)
		}
	}
	if got := fake.recorded(t, "args"); got != "task\ndone\nT-003\n--json" {
		t.Errorf("ran with %q: no message or log was given", got)
	}
}
