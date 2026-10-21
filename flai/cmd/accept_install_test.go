package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/preview"
)

// installProc is a verify.Proc that records each run and answers the exit
// status or start error its command's program names, printing out first.
type installProc struct {
	ran   [][]string
	dirs  []string
	exits map[string]int
	fails map[string]error
	out   string
}

func (f *installProc) Run(_ context.Context, dir string, argv, _ []string, stdout, _ io.Writer) (int, error) {
	f.ran = append(f.ran, argv)
	f.dirs = append(f.dirs, dir)
	if err := f.fails[argv[0]]; err != nil {
		return 0, err
	}
	_, _ = io.WriteString(stdout, f.out)
	return f.exits[argv[0]], nil
}

// installProjects are three sub-projects, two of which install.
func installProjects() []manifest.Project {
	return []manifest.Project{
		{Name: "web", Path: "web", Install: &manifest.Install{Command: []string{"pnpm", "install"}, Paths: []string{"package.json", "pnpm-lock.yaml"}}},
		{Name: "cli", Path: "cli"},
		{Name: "api", Path: "svc/api", Install: &manifest.Install{Command: []string{"go", "mod", "download"}, Paths: []string{"go.sum"}}},
	}
}

// ADR-0185: each sub-project whose install paths a changed path matches runs
// its command once, in manifest order, in its folder under the main
// checkout, with its output on the writer given; the others run nothing.
func TestRunInstallsRunsEachMatchingProjectOnceInOrder(t *testing.T) {
	proc := &installProc{out: "added 1 package\n"}
	var out bytes.Buffer
	changed := []string{"svc/api/go.sum", "web/pnpm-lock.yaml", "web/package.json", "cli/go.sum", "docs/x.md"}

	runs := runInstalls(context.Background(), proc, installProjects(), "/main", changed, &out)

	want := []installRun{
		{Project: "web", Command: []string{"pnpm", "install"}, OK: true, folder: filepath.Join("/main", "web")},
		{Project: "api", Command: []string{"go", "mod", "download"}, OK: true, folder: filepath.Join("/main", "svc", "api")},
	}
	if len(runs) != len(want) {
		t.Fatalf("runs %+v, want %+v", runs, want)
	}
	for i := range want {
		if r := runs[i]; r.Project != want[i].Project || !slices.Equal(r.Command, want[i].Command) || r.OK != want[i].OK || r.folder != want[i].folder || r.Error != "" {
			t.Errorf("run %d: %+v, want %+v", i, r, want[i])
		}
	}
	if !slices.Equal(proc.dirs, []string{want[0].folder, want[1].folder}) {
		t.Errorf("ran in %v", proc.dirs)
	}
	if got := out.String(); got != "added 1 package\nadded 1 package\n" {
		t.Errorf("the commands' output was not passed on: %q", got)
	}
	if runs := runInstalls(context.Background(), proc, installProjects(), "/main", []string{"web/src/app.ts", "cli/go.sum"}, io.Discard); len(runs) != 0 {
		t.Errorf("a change to no install path ran %+v", runs)
	}
}

// ADR-0185: a command that exits non-zero or cannot start is recorded with
// why, and the last lines of what it printed; the runs after it still run.
func TestRunInstallsRecordsAFailureAndGoesOn(t *testing.T) {
	proc := &installProc{exits: map[string]int{"pnpm": 3}, fails: map[string]error{"go": errors.New("exec: \"go\": executable file not found in $PATH")}}
	proc.out = strings.Repeat("progress\n", 30) + "\nERR_PNPM_FETCH_404 left-pad\n"

	runs := runInstalls(context.Background(), proc, installProjects(), "/main", []string{"web/package.json", "svc/api/go.sum"}, io.Discard)

	if len(runs) != 2 {
		t.Fatalf("runs %+v", runs)
	}
	web, api := runs[0], runs[1]
	if web.OK || web.Error != "exit status 3" || web.exit != 3 {
		t.Errorf("the failing run: %+v", web)
	}
	if lines := strings.Split(web.output, "\n"); len(lines) != installTailLines || lines[len(lines)-1] != "ERR_PNPM_FETCH_404 left-pad" {
		t.Errorf("the failing run kept %q, want its last %d lines that are not blank", web.output, installTailLines)
	}
	if api.OK || api.Error != `cannot start: exec: "go": executable file not found in $PATH` || api.exit != -1 {
		t.Errorf("the run that could not start: %+v", api)
	}
}

// ADR-0185: the text output names each run, and for one that failed the
// reason and the command to run by hand, quoted for a shell, and where.
func TestPrintAcceptNamesEachInstall(t *testing.T) {
	var out bytes.Buffer
	a := &app{out: &out}
	acc := &accepted{Acceptance: &preview.Acceptance{ID: "S-0001"}, Installs: []installRun{
		{Project: "web", Command: []string{"pnpm", "install"}, OK: true, folder: "/main/web"},
		{Project: "site", Command: []string{"sh", "-c", "it's done"}, Error: "exit status 3", folder: "/main/site"},
	}}
	if err := a.printAccept(acc); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"installed web: pnpm install\n",
		`install for site failed: exit status 3; run sh -c 'it'\''s done' in /main/site` + "\n",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the output lacks %q:\n%s", want, out.String())
		}
	}
}

// installRuns is what flai accept --json says of the installs it ran.
type installRuns struct {
	Installs []struct {
		Project string   `json:"project"`
		Command []string `json:"command"`
		OK      bool     `json:"ok"`
		Error   string   `json:"error"`
	} `json:"installs"`
}

// installProject is syncProject with the sub-project web, whose install
// command is command on web/lock.yaml, and S-0001 in review, whose branch
// writes content to rel. Integration: it runs real git.
func installProject(t *testing.T, command, rel, content string) string {
	t.Helper()
	root := syncProject(t)
	m, err := os.ReadFile(filepath.Join(root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	m = append(m, "projects:\n  - name: web\n    path: web\n    kind: node\n    install:\n      command: "+command+"\n      paths: [lock.yaml]\n"...)
	if err := os.WriteFile(filepath.Join(root, manifest.File), m, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "web", "lock.yaml"), []byte("left-pad: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\nweb/installed.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: the web sub-project")
	wt := openSyncStory(t, root, 1, "web")
	commitIn(t, wt, rel, content)
	issueStoryInReview(t, root, "S-0001", "T-0001")
	return root
}

// copyLock is an install command that copies the lock file it installs
// from to installed.txt, so a test sees which lock it ran on.
const copyLock = "[sh, -c, 'cp lock.yaml installed.txt']"

// I-0179: accepting a story that changed a sub-project's lock file left the
// main checkout's install as it was, since nothing installed after the
// merge. The acceptance now runs the sub-project's install command in the
// main checkout, on the lock file it merged.
func TestAcceptInstallsASubProjectWhoseLockTheStoryChanged(t *testing.T) {
	root := installProject(t, copyLock, "web/lock.yaml", "left-pad: 2\n")

	out, errOut, code := runIn(t, root, "accept", "S-0001")
	if code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	got, err := os.ReadFile(filepath.Join(root, "web", "installed.txt"))
	if err != nil || string(got) != "left-pad: 2\n" {
		t.Errorf("the main checkout's install is %q (%v), want the story's lock file", got, err)
	}
	if want := "installed web: sh -c 'cp lock.yaml installed.txt'\n"; !strings.Contains(out, want) {
		t.Errorf("the output lacks %q:\n%s", want, out)
	}
}

// ADR-0185: a story that changes none of a sub-project's install paths runs
// no install.
func TestAcceptInstallsNothingForAStoryThatChangesNoInstallPath(t *testing.T) {
	root := installProject(t, copyLock, "web/readme.md", "web\n")

	out, errOut, code := runIn(t, root, "--json", "accept", "S-0001")
	if code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "web", "installed.txt")); !os.IsNotExist(err) {
		t.Errorf("an install ran: %v", err)
	}
	var got installRuns
	if err := json.Unmarshal([]byte(out), &got); err != nil || len(got.Installs) != 0 {
		t.Errorf("installs %+v (%v):\n%s", got.Installs, err, out)
	}
}

// ADR-0185: an install that fails leaves the acceptance done and committed
// with exit 0, and --json reports it with its command and why, while what it
// printed goes to standard error, out of the JSON.
func TestAcceptReportsAFailedInstallAndStillAccepts(t *testing.T) {
	cmd := "[sh, -c, 'echo boom; exit 3']"
	root := installProject(t, cmd, "web/lock.yaml", "left-pad: 2\n")

	out, errOut, code := runIn(t, root, "--json", "accept", "S-0001")
	if code != 0 {
		t.Fatalf("accept failed with the install: %d %s", code, errOut)
	}
	var got installRuns
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("accept --json: %v\n%s", err, out)
	}
	want := []string{"sh", "-c", "echo boom; exit 3"}
	if len(got.Installs) != 1 || got.Installs[0].Project != "web" || !slices.Equal(got.Installs[0].Command, want) || got.Installs[0].OK || got.Installs[0].Error != "exit status 3" {
		t.Errorf("installs %+v", got.Installs)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); !strings.Contains(subject, "[S-0001] accept") {
		t.Errorf("the acceptance is not committed: %q", subject)
	}
	if left := gitIn(t, root, "status", "--porcelain"); left != "" {
		t.Errorf("the acceptance left uncommitted: %s", left)
	}
	for _, field := range []string{"boom", `"project":"web"`, `"exit":3`, `"err":"exit status 3"`} {
		if !strings.Contains(errOut, field) {
			t.Errorf("standard error lacks %s:\n%s", field, errOut)
		}
	}
}

// ADR-0185: --dry-run runs no install.
func TestAcceptDryRunInstallsNothing(t *testing.T) {
	root := installProject(t, copyLock, "web/lock.yaml", "left-pad: 2\n")

	if _, errOut, code := runIn(t, root, "accept", "--dry-run", "S-0001"); code != 0 {
		t.Fatalf("accept --dry-run: %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "web", "installed.txt")); !os.IsNotExist(err) {
		t.Errorf("an install ran: %v", err)
	}
}
