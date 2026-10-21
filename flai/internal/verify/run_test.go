package verify

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

// step is what a fake command does.
type step struct {
	stdout, stderr string
	exit           int
	err            error
	took           time.Duration
}

// fakeProc runs commands by their first argument, advancing its clock by
// each one's time.
type fakeProc struct {
	steps map[string]step
	now   time.Time
	ran   []string
	dirs  []string
	envs  [][]string
}

func (f *fakeProc) Run(_ context.Context, dir string, argv, env []string, stdout, stderr io.Writer) (int, error) {
	f.ran = append(f.ran, strings.Join(argv, " "))
	f.dirs = append(f.dirs, dir)
	f.envs = append(f.envs, env)
	s := f.steps[argv[0]]
	_, _ = io.WriteString(stdout, s.stdout)
	_, _ = io.WriteString(stderr, s.stderr)
	f.now = f.now.Add(s.took)
	return s.exit, s.err
}

func (f *fakeProc) clock() time.Time { return f.now }

func sel(name, dir string, argv ...string) Selected {
	return Selected{Tier: Tier{Name: name, Dir: dir}, Argv: argv}
}

func TestRunStopsAtTheFirstFailingTier(t *testing.T) {
	proc := &fakeProc{steps: map[string]step{
		"vet":  {took: 1500 * time.Millisecond},
		"test": {stdout: "ok\n", stderr: "boom\n", exit: 1, took: 42 * time.Millisecond},
		"lint": {took: time.Second},
	}}
	res, err := Run(context.Background(), "/repo", []Selected{
		sel("vet", "flai", "vet", "./..."), sel("test", "flai", "test"), sel("lint", "", "lint"),
	}, RunOptions{Proc: proc, Now: proc.clock, FS: fstest.MapFS{}})
	if err != nil || res.Passed {
		t.Fatalf("run %+v, %v", res, err)
	}
	if !reflect.DeepEqual(proc.ran, []string{"vet ./...", "test"}) || proc.dirs[0] != "/repo/flai" {
		t.Errorf("ran %q in %q", proc.ran, proc.dirs)
	}
	zero, one := 0, 1
	want := []TierResult{
		{Name: "vet", Command: []string{"vet", "./..."}, Dir: "flai", State: Passed, ExitCode: &zero, DurationMS: 1500, Duration: "1.5s"},
		{Name: "test", Command: []string{"test"}, Dir: "flai", State: Failed, ExitCode: &one, DurationMS: 42, Duration: "42ms",
			Findings: []Finding{{Name: "test", Message: "ok\nboom\nexit status 1"}}},
		{Name: "lint", Command: []string{"lint"}, State: NotReached},
	}
	if !reflect.DeepEqual(res.Tiers, want) {
		t.Errorf("tiers\n%+v\nwant\n%+v", res.Tiers, want)
	}
}

// S-0457: with a log folder, each tier that runs writes its whole output,
// stdout and stderr, to its own file there, which its result names, below
// the run's root or at an absolute folder; a tier not reached writes none,
// and with no folder nothing is written.
func TestRunWritesEachTierThatRunsItsWholeOutputToALog(t *testing.T) {
	steps := map[string]step{
		"vet":  {stdout: "vetted\n"},
		"test": {stdout: "--- FAIL: TestX\n", stderr: "boom\n", exit: 1},
	}
	tiers := []Selected{sel("vet", "", "vet"), sel("go test", "", "test"), sel("lint", "", "lint")}
	root := t.TempDir()
	res, _ := Run(context.Background(), root, tiers, RunOptions{Proc: &fakeProc{steps: steps}, FS: fstest.MapFS{}, LogDir: "logs"})
	for _, c := range []struct{ tier, log, holds string }{
		{"vet", filepath.Join("logs", "vet.log"), "vetted\n"},
		{"go test", filepath.Join("logs", "go-test.log"), "--- FAIL: TestX\nboom\n"},
	} {
		tr := res.Tiers[slices.IndexFunc(res.Tiers, func(tr TierResult) bool { return tr.Name == c.tier })]
		got, err := os.ReadFile(filepath.Join(root, tr.Log))
		if tr.Log != c.log || err != nil || string(got) != c.holds {
			t.Errorf("%s: log %q holds %q, %v; want %q holding %q", c.tier, tr.Log, got, err, c.log, c.holds)
		}
	}
	if lint := res.Tiers[2]; lint.State != NotReached || lint.Log != "" {
		t.Errorf("a tier not reached: %+v", lint)
	}
	if _, err := os.Stat(filepath.Join(root, "logs", "lint.log")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a tier not reached wrote a log: %v", err)
	}
	if step := tierStep(res.Tiers[1]); step.Log != res.Tiers[1].Log {
		t.Errorf("its step names %q", step.Log)
	}

	abs := filepath.Join(t.TempDir(), "batch")
	res, _ = Run(context.Background(), "/repo", tiers[:1], RunOptions{Proc: &fakeProc{steps: steps}, FS: fstest.MapFS{}, LogDir: abs})
	if got, err := os.ReadFile(filepath.Join(abs, "vet.log")); res.Tiers[0].Log != filepath.Join(abs, "vet.log") || err != nil || string(got) != "vetted\n" {
		t.Errorf("an absolute folder: log %q holds %q, %v", res.Tiers[0].Log, got, err)
	}

	root = t.TempDir()
	res, _ = Run(context.Background(), root, tiers, RunOptions{Proc: &fakeProc{steps: steps}, FS: fstest.MapFS{}})
	if entries, err := os.ReadDir(root); err != nil || len(entries) != 0 || res.Tiers[0].Log != "" || res.Tiers[1].Log != "" {
		t.Errorf("with no folder: wrote %v, %v, tiers %+v", entries, err, res.Tiers)
	}
}

// S-0342: a covered tier is skipped in its place, naming the tier that
// covers it, without its command started, before and after a failure, and
// never fails the run.
func TestRunSkipsACoveredTierWithoutStartingIt(t *testing.T) {
	covered := sel("short", "flai", "short")
	covered.CoveredBy = "full"
	proc := &fakeProc{steps: map[string]step{}}
	res, err := Run(context.Background(), "/repo", []Selected{sel("vet", "", "vet"), covered, sel("full", "flai", "full")}, RunOptions{Proc: proc, Now: proc.clock, FS: fstest.MapFS{}})
	if err != nil || !res.Passed || !reflect.DeepEqual(proc.ran, []string{"vet", "full"}) {
		t.Fatalf("run %+v, %v, ran %q", res, err, proc.ran)
	}
	zero := 0
	want := []TierResult{
		{Name: "vet", Command: []string{"vet"}, State: Passed, ExitCode: &zero, Duration: "0s"},
		{Name: "short", Command: []string{"short"}, Dir: "flai", State: Skipped, CoveredBy: "full"},
		{Name: "full", Command: []string{"full"}, Dir: "flai", State: Passed, ExitCode: &zero, Duration: "0s"},
	}
	if !reflect.DeepEqual(res.Tiers, want) {
		t.Errorf("tiers\n%+v\nwant\n%+v", res.Tiers, want)
	}

	proc = &fakeProc{steps: map[string]step{"vet": {exit: 1}}}
	res, _ = Run(context.Background(), "/repo", []Selected{sel("vet", "", "vet"), covered, sel("full", "flai", "full")}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	var states []State
	for _, tr := range res.Tiers {
		states = append(states, tr.State)
	}
	if res.Passed || !reflect.DeepEqual(states, []State{Failed, Skipped, NotReached}) || res.Tiers[1].CoveredBy != "full" {
		t.Errorf("after a failure: %+v", res)
	}
}

// S-0342: with integration covering go-test, flai test --all and a story's
// close-out for a Go file report go-test skipped and run only integration;
// for paths alone go-test runs.
func TestIntegrationCoveringGoTestSkipsIt(t *testing.T) {
	tiers := []Tier{
		{Name: "go-test", Dir: "flai", Paths: []string{"flai"}, Command: []string{"go-short", "{packages}"}},
		{Name: "integration", AllOnly: true, Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go-full", "{packages}"}, Covers: []string{"go-test"}},
	}
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/main.go": {}}
	for _, tc := range []struct {
		name     string
		selected []Selected
		ran      []string
		skipped  bool
	}{
		{"all", Select(fsys, tiers, nil, true), []string{"go-full ./..."}, true},
		{"a story's close-out", SelectStory(fsys, tiers, []string{"flai/main.go"}), []string{"go-full ./..."}, true},
		{"paths alone", Select(fsys, tiers, []string{"flai/main.go"}, false), []string{"go-short ."}, false},
	} {
		proc := &fakeProc{steps: map[string]step{}}
		res, err := Run(context.Background(), "/repo", tc.selected, RunOptions{Proc: proc, FS: fsys})
		if err != nil || !res.Passed || !reflect.DeepEqual(proc.ran, tc.ran) {
			t.Errorf("%s: ran %q, want %q: %+v, %v", tc.name, proc.ran, tc.ran, res, err)
			continue
		}
		if got := res.Tiers[0]; got.Name != "go-test" || (got.State == Skipped) != tc.skipped || tc.skipped && got.CoveredBy != "integration" {
			t.Errorf("%s: go-test %+v", tc.name, got)
		}
	}
}

func TestRunPassesWhenEveryTierDoesOrNoneIsSelected(t *testing.T) {
	proc := &fakeProc{steps: map[string]step{}}
	res, err := Run(context.Background(), "/repo", []Selected{sel("a", "", "a"), sel("b", "", "b")}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if err != nil || !res.Passed || len(res.Tiers) != 2 || res.Tiers[1].State != Passed {
		t.Errorf("run %+v, %v", res, err)
	}
	res, err = Run(context.Background(), "/repo", nil, RunOptions{Proc: proc})
	if err != nil || !res.Passed || res.Tiers == nil || len(res.Tiers) != 0 {
		t.Errorf("nothing selected: %+v, %v", res, err)
	}
}

func TestRunFailsATierThatCannotStartOrGofmtListingFiles(t *testing.T) {
	proc := &fakeProc{steps: map[string]step{
		"missing": {err: errors.New(`exec: "missing": executable file not found in $PATH`)},
		"gofmt":   {stdout: "a.go\n"},
	}}
	res, _ := Run(context.Background(), "/repo", []Selected{sel("m", "", "missing")}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if tr := res.Tiers[0]; res.Passed || tr.State != Failed || tr.ExitCode != nil || len(tr.Findings) != 1 || !strings.Contains(tr.Findings[0].Message, "not found") {
		t.Errorf("a tier that cannot start: %+v", res)
	}
	res, _ = Run(context.Background(), "/repo", []Selected{{Tier: Tier{Name: "m"}}}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if res.Passed || len(proc.ran) != 1 || !strings.Contains(res.Tiers[0].Findings[0].Message, "no command") {
		t.Errorf("a tier with no command ran or passed: %+v", res)
	}
	gofmt := Selected{Tier: Tier{Name: "gofmt", Format: FormatGofmtList}, Argv: []string{"gofmt", "-l", "."}}
	res, _ = Run(context.Background(), "/repo", []Selected{gofmt}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	want := []Finding{{Name: "gofmt", Path: "a.go", Message: "not gofmt-formatted"}}
	if res.Passed || !reflect.DeepEqual(res.Tiers[0].Findings, want) {
		t.Errorf("gofmt listing a file: %+v", res)
	}
}

func TestRunCapsFindingsAcrossTheRunAtFiveByDefault(t *testing.T) {
	proc := &fakeProc{steps: map[string]step{"gofmt": {stdout: "a.go\nb.go\nc.go\nd.go\ne.go\nf.go\ng.go\n"}}}
	gofmt := Selected{Tier: Tier{Name: "gofmt", Format: FormatGofmtList}, Argv: []string{"gofmt"}}
	res, _ := Run(context.Background(), "/repo", []Selected{gofmt}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if tr := res.Tiers[0]; len(tr.Findings) != DefaultMax || tr.Omitted != 2 || tr.Findings[4].Path != "e.go" {
		t.Errorf("findings %+v, omitted %d", tr.Findings, tr.Omitted)
	}
}

func TestTestReadsThePathsItIsGivenOrTheChangedOnes(t *testing.T) {
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/a.go": {}, "flai/sub/b.go": {}, "docs/x.md": {}}
	tiers := []Tier{
		{Name: "vet", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "vet", "{packages}"}},
		{Name: "md", Paths: []string{"**/*.md"}, Command: []string{"mdcheck", "{files}"}},
	}
	git := &fakeGit{out: map[string]string{
		"diff --name-only --diff-filter=d main...HEAD": "docs/x.md",
		"status --porcelain --untracked-files=all":     "",
	}}
	opts := Options{Root: "/repo", Tiers: tiers, Base: "main", Git: git}
	opts.FS = fsys

	opts.Proc = &fakeProc{}
	res, err := Test(context.Background(), opts)
	if err != nil || !res.Passed || !reflect.DeepEqual(res.Paths, []string{"docs/x.md"}) || len(res.Tiers) != 1 ||
		!reflect.DeepEqual(res.Tiers[0].Command, []string{"mdcheck", "docs/x.md"}) {
		t.Errorf("changed paths: %+v, %v", res, err)
	}

	opts.Args = []string{"flai/sub"}
	proc := &fakeProc{}
	opts.Proc = proc
	res, err = Test(context.Background(), opts)
	if err != nil || !reflect.DeepEqual(res.Paths, []string{"flai/sub/b.go"}) || !reflect.DeepEqual(proc.ran, []string{"go vet ./sub"}) {
		t.Errorf("args: %+v, %v, ran %q", res, err, proc.ran)
	}

	opts.All = true
	proc = &fakeProc{}
	opts.Proc = proc
	res, err = Test(context.Background(), opts)
	if err != nil || res.Paths != nil || !reflect.DeepEqual(proc.ran, []string{"go vet ./...", "mdcheck ."}) {
		t.Errorf("all: %+v, %v, ran %q", res, err, proc.ran)
	}

	opts.All = false
	opts.Args = []string{"nowhere"}
	if _, err := Test(context.Background(), opts); err == nil {
		t.Error("a path that is not in the checkout gave no error")
	}
}

func TestRunReportsTheContextEndingDuringATier(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	proc := &cancellingProc{cancel: cancel}
	res, err := Run(ctx, "/repo", []Selected{sel("a", "", "a"), sel("b", "", "b")}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if !errors.Is(err, context.Canceled) || res.Passed || res.Tiers[0].State != Failed || res.Tiers[1].State != NotReached {
		t.Errorf("run %+v, %v", res, err)
	}
	if f := res.Tiers[0].Findings; len(f) != 1 || !strings.Contains(f[0].Message, "stopped: context canceled") {
		t.Errorf("findings %+v", f)
	}
}

// cancellingProc ends the run's context while its command runs, as a
// cancelled MCP call does, and answers as a killed process would.
type cancellingProc struct{ cancel context.CancelFunc }

func (p *cancellingProc) Run(context.Context, string, []string, []string, io.Writer, io.Writer) (int, error) {
	p.cancel()
	return -1, nil
}

func TestRunGivesEveryTierTheEnvironmentItIsAsked(t *testing.T) {
	proc := &fakeProc{steps: map[string]step{}}
	env := []string{"CLOSE_OUT_STORY=S-0001"}
	if _, err := Run(context.Background(), "/repo", []Selected{sel("a", "", "a"), sel("b", "", "b")}, RunOptions{Proc: proc, FS: fstest.MapFS{}, Env: env}); err != nil {
		t.Fatal(err)
	}
	want := []string{"CLOSE_OUT_STORY=S-0001", "FLAI_ROLE=verify"}
	if !reflect.DeepEqual(proc.envs, [][]string{want, want}) || !reflect.DeepEqual(env, []string{"CLOSE_OUT_STORY=S-0001"}) {
		t.Errorf("environments %q, asked %q", proc.envs, env)
	}
	proc = &fakeProc{steps: map[string]step{}}
	if _, err := Run(context.Background(), "/repo", []Selected{sel("a", "", "a")}, RunOptions{Proc: proc, FS: fstest.MapFS{}}); err != nil ||
		!reflect.DeepEqual(proc.envs, [][]string{{"FLAI_ROLE=verify"}}) {
		t.Errorf("with no environment asked a tier got %q, %v", proc.envs, err)
	}
}

// role is the FLAI_ROLE a command run with env, added to its own, sees: the
// last one, as exec takes it.
func role(env []string) string {
	got := ""
	for _, kv := range env {
		if v, ok := strings.CutPrefix(kv, "FLAI_ROLE="); ok {
			got = v
		}
	}
	return got
}

func TestEveryTierRunsAsTheVerifyRoleWhateverRoleAskedIt(t *testing.T) {
	for _, env := range [][]string{nil, {"CLOSE_OUT_STORY=S-0001"}, {"FLAI_ROLE=orchestrate", "CLOSE_OUT_STORY=S-0001"}} {
		proc := &fakeProc{steps: map[string]step{}}
		if _, err := Run(context.Background(), "/repo", []Selected{sel("a", "", "a"), sel("b", "", "b")}, RunOptions{Proc: proc, FS: fstest.MapFS{}, Env: env}); err != nil {
			t.Fatal(err)
		}
		for i, got := range proc.envs {
			if role(got) != "verify" || (len(env) > 0 && !slices.Contains(got, "CLOSE_OUT_STORY=S-0001")) {
				t.Errorf("asked %q, %s ran with %q", env, proc.ran[i], got)
			}
		}
	}

	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/a.go": {}}
	tiers := []Tier{{Name: "vet", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "vet", "{packages}"}}}
	proc := &fakeProc{}
	opts := Options{Root: "/repo", Tiers: tiers, All: true}
	opts.FS, opts.Proc, opts.Env = fsys, proc, []string{"FLAI_ROLE=orchestrate"}
	if _, err := Test(context.Background(), opts); err != nil || len(proc.envs) != 1 || role(proc.envs[0]) != "verify" {
		t.Errorf("flai test ran with %q, %v", proc.envs, err)
	}
}
