package verify

import (
	"context"
	"errors"
	"io"
	"reflect"
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
}

func (f *fakeProc) Run(_ context.Context, dir string, argv []string, stdout, stderr io.Writer) (int, error) {
	f.ran = append(f.ran, strings.Join(argv, " "))
	f.dirs = append(f.dirs, dir)
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

func (p *cancellingProc) Run(context.Context, string, []string, io.Writer, io.Writer) (int, error) {
	p.cancel()
	return -1, nil
}
