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
	"sync"
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
	want := [][]string{
		{"CLOSE_OUT_STORY=S-0001", "FLAI_ROLE=verify", "FLAI_VERIFY_PASSED="},
		{"CLOSE_OUT_STORY=S-0001", "FLAI_ROLE=verify", "FLAI_VERIFY_PASSED=a"},
	}
	if !reflect.DeepEqual(proc.envs, want) || !reflect.DeepEqual(env, []string{"CLOSE_OUT_STORY=S-0001"}) {
		t.Errorf("environments %q, asked %q", proc.envs, env)
	}
	proc = &fakeProc{steps: map[string]step{}}
	if _, err := Run(context.Background(), "/repo", []Selected{sel("a", "", "a")}, RunOptions{Proc: proc, FS: fstest.MapFS{}}); err != nil ||
		!reflect.DeepEqual(proc.envs, [][]string{{"FLAI_ROLE=verify", "FLAI_VERIFY_PASSED="}}) {
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

// liveProc runs each command by its first argument with a function of its
// own, at once from several goroutines, and keeps the order commands
// started and ended in, and which saw their context end.
type liveProc struct {
	mu      sync.Mutex
	fns     map[string]func(ctx context.Context, stdout io.Writer) int
	events  []string
	stopped []string
}

func (p *liveProc) Run(ctx context.Context, _ string, argv, _ []string, stdout, _ io.Writer) (int, error) {
	name := argv[0]
	p.note("start " + name)
	exit := 0
	if fn := p.fns[name]; fn != nil {
		exit = fn(ctx, stdout)
	}
	if ctx.Err() != nil {
		p.mu.Lock()
		p.stopped = append(p.stopped, name)
		p.mu.Unlock()
	}
	p.note("end " + name)
	return exit, nil
}

func (p *liveProc) note(e string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, e)
}

// at is where event e came in the events, -1 when it did not.
func (p *liveProc) at(e string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return slices.Index(p.events, e)
}

// psel is a selected tier marked Parallel.
func psel(name string) Selected {
	s := sel(name, "", name)
	s.Tier.Parallel = true
	return s
}

// resultStates are the tiers' states, in order.
func resultStates(res Result) []State {
	out := make([]State, 0, len(res.Tiers))
	for _, tr := range res.Tiers {
		out = append(out, tr.State)
	}
	return out
}

// awaitClose answers whether ch closed within a few seconds.
func awaitClose(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	case <-time.After(5 * time.Second):
		return false
	}
}

// S-0610, ADR-0266: two tiers marked parallel run at once, each starting
// only once the unmarked tier before them has ended, and the unmarked tier
// after them waits for both; a covered tier between them stays skipped.
func TestRunStartsTheParallelTiersAtOnceAfterTheUnmarkedOnesBeforeThem(t *testing.T) {
	var started sync.WaitGroup
	started.Add(2)
	both := make(chan struct{})
	go func() { started.Wait(); close(both) }()
	meet := func(context.Context, io.Writer) int {
		started.Done()
		if !awaitClose(both) {
			return 1 // the other never started beside it
		}
		return 0
	}
	proc := &liveProc{fns: map[string]func(context.Context, io.Writer) int{"a": meet, "b": meet}}
	covered := psel("c")
	covered.CoveredBy = "b"
	res, err := Run(context.Background(), "/repo", []Selected{sel("vet", "", "vet"), psel("a"), covered, psel("b"), sel("smoke", "", "smoke")},
		RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if err != nil || !res.Passed || !reflect.DeepEqual(resultStates(res), []State{Passed, Passed, Skipped, Passed, Passed}) {
		t.Fatalf("run %+v, %v; events %q", res, err, proc.events)
	}
	if proc.at("end vet") > proc.at("start a") || proc.at("end vet") > proc.at("start b") {
		t.Errorf("a parallel tier started before the unmarked tier before it ended: %q", proc.events)
	}
	if proc.at("start smoke") < proc.at("end a") || proc.at("start smoke") < proc.at("end b") {
		t.Errorf("the unmarked tier after the parallel ones started before both ended: %q", proc.events)
	}
	if proc.at("start c") != -1 {
		t.Errorf("a covered tier ran: %q", proc.events)
	}
}

// S-0610, ADR-0266: a tier that fails stops nothing before it. The tier
// before it that was running finishes, and when it fails too it is the
// failure the run names, with its findings, as a serial run would.
func TestRunLetsARunningEarlierTierFinishWhenALaterOneFails(t *testing.T) {
	for _, earlierFails := range []bool{false, true} {
		bEnded := make(chan struct{})
		proc := &liveProc{fns: map[string]func(context.Context, io.Writer) int{
			"a": func(_ context.Context, w io.Writer) int {
				awaitClose(bEnded)
				time.Sleep(20 * time.Millisecond) // past the run seeing b fail
				if earlierFails {
					_, _ = io.WriteString(w, "a broke\n")
					return 1
				}
				return 0
			},
			"b": func(_ context.Context, w io.Writer) int {
				defer close(bEnded)
				_, _ = io.WriteString(w, "b broke\n")
				return 1
			},
		}}
		res, err := Run(context.Background(), "/repo", []Selected{psel("a"), psel("b")}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
		if err != nil || res.Passed {
			t.Fatalf("earlier fails %v: run %+v, %v", earlierFails, res, err)
		}
		if slices.Contains(proc.stopped, "a") {
			t.Errorf("earlier fails %v: the earlier tier was stopped", earlierFails)
		}
		want, named := []State{Passed, Failed}, "b broke"
		if earlierFails {
			want, named = []State{Failed, NotReached}, "a broke"
		}
		if got := resultStates(res); !reflect.DeepEqual(got, want) {
			t.Errorf("earlier fails %v: states %v, want %v", earlierFails, got, want)
		}
		var findings []Finding
		for _, tr := range res.Tiers {
			findings = append(findings, tr.Findings...)
		}
		if len(findings) != 1 || !strings.Contains(findings[0].Message, named) {
			t.Errorf("earlier fails %v: findings %+v, want one naming %q", earlierFails, findings, named)
		}
	}
}

// S-0610, ADR-0266: a failure stops the tier after it that is running,
// through its context, lists every tier after it not reached, the one that
// had already passed too, starts none after it, and keeps the log of none
// of them.
func TestRunStopsTheTiersAfterAFailureAndListsThemNotReached(t *testing.T) {
	bEnded, cStarted := make(chan struct{}), make(chan struct{})
	proc := &liveProc{fns: map[string]func(context.Context, io.Writer) int{
		"a": func(context.Context, io.Writer) int {
			if !awaitClose(bEnded) || !awaitClose(cStarted) {
				return 2
			}
			time.Sleep(20 * time.Millisecond) // past the run seeing b pass
			return 1
		},
		"b": func(context.Context, io.Writer) int { close(bEnded); return 0 },
		"c": func(ctx context.Context, _ io.Writer) int {
			close(cStarted)
			select {
			case <-ctx.Done():
				return -1
			case <-time.After(5 * time.Second):
				return 0
			}
		},
	}}
	root := t.TempDir()
	res, err := Run(context.Background(), root, []Selected{psel("a"), psel("b"), psel("c"), sel("d", "", "d")},
		RunOptions{Proc: proc, FS: fstest.MapFS{}, LogDir: "logs"})
	if err != nil || res.Passed || !reflect.DeepEqual(resultStates(res), []State{Failed, NotReached, NotReached, NotReached}) {
		t.Fatalf("run %+v, %v; events %q", res, err, proc.events)
	}
	if res.Tiers[0].ExitCode == nil || *res.Tiers[0].ExitCode != 1 {
		t.Errorf("a's exit %v, want 1: it did not see b end and c start", res.Tiers[0].ExitCode)
	}
	if !slices.Contains(proc.stopped, "c") {
		t.Errorf("the running tier after the failure was not stopped: %q", proc.stopped)
	}
	if proc.at("start d") != -1 {
		t.Errorf("a tier after the failure started: %q", proc.events)
	}
	for _, tr := range res.Tiers[1:] {
		if tr.Log != "" || tr.ExitCode != nil || tr.Duration != "" {
			t.Errorf("a tier not reached carries what it did: %+v", tr)
		}
	}
	logs, _ := os.ReadDir(filepath.Join(root, "logs"))
	if len(logs) != 1 || logs[0].Name() != "a.log" {
		t.Errorf("logs %v, want only a.log", logs)
	}
}

// S-0610, ADR-0266: over the same outcomes, a run with every tier marked
// parallel answers what the serial run answers: the verdict, the states,
// the tier it names, and the findings.
func TestRunAtOnceAnswersWhatTheSerialRunAnswers(t *testing.T) {
	steps := map[string]step{
		"a": {},
		"b": {stdout: "first\n", exit: 1},
		"c": {stdout: "second\n", exit: 1},
		"d": {},
	}
	fixed := time.Unix(0, 0)
	run := func(parallel bool) Result {
		var mu sync.Mutex
		proc := &liveProc{fns: map[string]func(context.Context, io.Writer) int{}}
		for name, s := range steps {
			proc.fns[name] = func(_ context.Context, w io.Writer) int {
				mu.Lock()
				defer mu.Unlock()
				_, _ = io.WriteString(w, s.stdout)
				return s.exit
			}
		}
		var tiers []Selected
		for _, name := range []string{"a", "b", "c", "d"} {
			s := sel(name, "", name)
			s.Tier.Parallel = parallel
			tiers = append(tiers, s)
		}
		res, err := Run(context.Background(), "/repo", tiers, RunOptions{Proc: proc, FS: fstest.MapFS{}, Now: func() time.Time { return fixed }})
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	serial, atOnce := run(false), run(true)
	if !reflect.DeepEqual(serial, atOnce) {
		t.Errorf("at once\n%+v\nserial\n%+v", atOnce, serial)
	}
	if serial.Passed || !reflect.DeepEqual(resultStates(serial), []State{Passed, Failed, NotReached, NotReached}) {
		t.Errorf("serial %+v", serial)
	}
}

// S-0610, ADR-0266: when the caller's context ends, every tier running is
// stopped; the first in list order is the failure, and Run returns the
// context's error.
func TestRunStopsEveryRunningTierWhenTheContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var started sync.WaitGroup
	started.Add(2)
	go func() { started.Wait(); cancel() }()
	block := func(ctx context.Context, _ io.Writer) int {
		started.Done()
		select {
		case <-ctx.Done():
			return -1
		case <-time.After(5 * time.Second):
			return 0
		}
	}
	proc := &liveProc{fns: map[string]func(context.Context, io.Writer) int{"a": block, "b": block}}
	res, err := Run(ctx, "/repo", []Selected{psel("a"), psel("b")}, RunOptions{Proc: proc, FS: fstest.MapFS{}})
	if !errors.Is(err, context.Canceled) || res.Passed || !reflect.DeepEqual(resultStates(res), []State{Failed, NotReached}) {
		t.Fatalf("run %+v, %v", res, err)
	}
	if f := res.Tiers[0].Findings; len(f) != 1 || !strings.Contains(f[0].Message, "stopped: context canceled") {
		t.Errorf("findings %+v", f)
	}
	if slices.Sort(proc.stopped); !reflect.DeepEqual(proc.stopped, []string{"a", "b"}) {
		t.Errorf("stopped %q, want both", proc.stopped)
	}
}

// S-0610, ADR-0266: a parallel tier waits for the unmarked tiers before it,
// an unmarked one for every tier before it, and a covered tier neither
// waits nor is waited for.
func TestWaitsFollowTheParallelMarks(t *testing.T) {
	covered := psel("c")
	covered.CoveredBy = "x"
	got := waits([]Selected{sel("u", "", "u"), psel("p"), covered, psel("q"), sel("v", "", "v"), psel("r")})
	want := [][]int{nil, {0}, nil, {0}, {0, 1, 3}, {0, 4}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("waits %v, want %v", got, want)
	}
}

// envProc keeps the environment each command, by its first argument, ran
// with, from several goroutines at once, and fails the commands fail names.
type envProc struct {
	mu   sync.Mutex
	envs map[string][]string
	fail map[string]bool
}

func (p *envProc) Run(_ context.Context, _ string, argv, env []string, _, _ io.Writer) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.envs == nil {
		p.envs = map[string][]string{}
	}
	p.envs[argv[0]] = env
	if p.fail[argv[0]] {
		return 1, nil
	}
	return 0, nil
}

// passed is the FLAI_VERIFY_PASSED the command named name ran with, and
// whether it ran.
func (p *envProc) passed(name string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	env, ok := p.envs[name]
	if !ok {
		return "", false
	}
	got := ""
	for _, kv := range env {
		if v, found := strings.CutPrefix(kv, PassedEnv+"="); found {
			got = v
		}
	}
	return got, true
}

// S-0610, ADR-0266: each tier is told, in FLAI_VERIFY_PASSED, what passed
// before the run, from the PassedEnv entry the caller sets, then the tiers
// it waited for, in list order, a tier that ran its whole command as
// <name>:all: never a parallel tier beside it, whatever ended first.
func TestRunTellsEachTierTheTiersItWaitedForThatPassed(t *testing.T) {
	whole := Selected{Tier: Tier{Name: "markdown", Command: []string{"lint", "{files}"}, AllCommand: []string{"lint-all"}}, Argv: []string{"lint-all"}}
	paths := Selected{Tier: Tier{Name: "lint", Command: []string{"lint", "{files}"}}, Argv: []string{"lint", "a.md"}}
	covered := sel("unit", "", "unit")
	covered.CoveredBy = "integration"
	selected := []Selected{paths, whole, covered, psel("flaiover"), psel("integration"), sel("smoke", "", "smoke")}
	proc := &envProc{}
	env := []string{"CLOSE_OUT_STORY=S-0001", PassedEnv + "=step:check gofmt"}
	res, err := Run(context.Background(), "/repo", selected, RunOptions{Proc: proc, FS: fstest.MapFS{}, Env: env})
	if err != nil || !res.Passed {
		t.Fatalf("result %+v, %v", res, err)
	}
	want := map[string]string{
		"lint":        "step:check gofmt",
		"lint-all":    "step:check gofmt lint",
		"flaiover":    "step:check gofmt lint markdown:all",
		"integration": "step:check gofmt lint markdown:all",
		"smoke":       "step:check gofmt lint markdown:all flaiover integration",
	}
	for name, w := range want {
		if got, ok := proc.passed(name); !ok || got != w {
			t.Errorf("%s was told %q (ran %v), want %q", name, got, ok, w)
		}
	}
	if _, ran := proc.passed("unit"); ran {
		t.Error("the covered tier ran")
	}
}
