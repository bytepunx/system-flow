package verify

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// fixProc fakes fix commands by their first argument: upper rewrites each
// file argument, relative to dir, in upper case; any other exits with the
// status, and fails to start with the error, that exits and errs give it.
type fixProc struct {
	exits map[string]int
	errs  map[string]error
	ran   []string
	dirs  []string
}

func (f *fixProc) Run(_ context.Context, dir string, argv, _ []string, stdout, _ io.Writer) (int, error) {
	f.ran = append(f.ran, strings.Join(argv, " "))
	f.dirs = append(f.dirs, dir)
	if err := f.errs[argv[0]]; err != nil {
		return 0, err
	}
	if argv[0] == "upper" {
		for _, a := range argv[1:] {
			p := filepath.Join(dir, filepath.FromSlash(a))
			data, err := os.ReadFile(p)
			if err != nil {
				return 0, err
			}
			if err := os.WriteFile(p, bytes.ToUpper(data), 0o644); err != nil {
				return 0, err
			}
		}
		return 0, nil
	}
	_, _ = io.WriteString(stdout, "lint: 3 problems\n")
	return f.exits[argv[0]], nil
}

// tree writes the files, path to content, into a temporary checkout.
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func read(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0367: a tier's fix command runs in its dir over the existing files its
// paths select, relative to that dir, and Fix answers the ones it changed;
// a path no tier selects, a path that is not there, and a tier without a fix
// command run nothing.
func TestFixRewritesTheFilesTheTierSelects(t *testing.T) {
	root := tree(t, map[string]string{"docs/a.md": "hello\n", "docs/b.md": "SAME\n", "docs/drafts/c.md": "draft\n", "main.go": "package main\n"})
	tiers := []Tier{
		{Name: "unit", Command: []string{"go", "test", "{packages}"}, Paths: []string{"**"}},
		{Name: "md", Command: []string{"lint", "{files}"}, Dir: "docs", Paths: []string{"**/*.md", "!docs/drafts/**"}, FixCommand: []string{"upper", "{files}"}},
	}
	proc := &fixProc{}
	got := Fix(context.Background(), root, tiers, []string{"docs/b.md", "main.go", "docs/a.md", "docs/gone.md", "docs/drafts/c.md", "./docs/a.md"}, proc, nil)
	if want := []string{"docs/a.md"}; !reflect.DeepEqual(got, want) {
		t.Errorf("changed = %v, want %v", got, want)
	}
	if want := []string{"upper a.md b.md"}; !reflect.DeepEqual(proc.ran, want) {
		t.Errorf("ran %q, want %q: only the md tier's fix, over the existing files it selects", proc.ran, want)
	}
	if want := filepath.Join(root, "docs"); len(proc.dirs) != 1 || proc.dirs[0] != want {
		t.Errorf("ran in %v, want %s", proc.dirs, want)
	}
	if s := read(t, root, "docs/a.md"); s != "HELLO\n" {
		t.Errorf("docs/a.md = %q, want it rewritten", s)
	}
	if s := read(t, root, "docs/drafts/c.md"); s != "draft\n" {
		t.Errorf("docs/drafts/c.md = %q, want it untouched", s)
	}

	proc = &fixProc{}
	if got := Fix(context.Background(), root, tiers, []string{"main.go", "docs/gone.md"}, proc, nil); got != nil || proc.ran != nil {
		t.Errorf("nothing a fix tier selects: changed %v, ran %q; want neither", got, proc.ran)
	}
}

// S-0367: a fix command that exits non-zero or cannot start is logged at
// warn with its tier and exit code, and the tiers after it still run and
// count.
func TestFixLogsAFailingCommandAndGoesOn(t *testing.T) {
	root := tree(t, map[string]string{"a.md": "hello\n"})
	tiers := []Tier{
		{Name: "strict", Paths: []string{"**"}, FixCommand: []string{"lint", "--fix", "{files}"}},
		{Name: "absent", Paths: []string{"**"}, FixCommand: []string{"nope", "{files}"}},
		{Name: "md", Paths: []string{"*.md"}, FixCommand: []string{"upper", "{files}"}},
	}
	proc := &fixProc{exits: map[string]int{"lint": 2}, errs: map[string]error{"nope": errors.New("executable file not found")}}
	var logs bytes.Buffer
	got := Fix(context.Background(), root, tiers, []string{"a.md"}, proc, slog.New(slog.NewTextHandler(&logs, nil)))
	if want := []string{"a.md"}; !reflect.DeepEqual(got, want) {
		t.Errorf("changed = %v, want %v", got, want)
	}
	if len(proc.ran) != 3 {
		t.Errorf("ran %q, want every tier's fix", proc.ran)
	}
	out := logs.String()
	for _, want := range []string{"level=WARN", "component=verify", "tier=strict", "exit_code=2", `output="lint: 3 problems"`, "tier=absent", `err="executable file not found"`} {
		if !strings.Contains(out, want) {
			t.Errorf("log lacks %s:\n%s", want, out)
		}
	}
	if strings.Contains(out, "tier=md") {
		t.Errorf("a fix that passed is logged:\n%s", out)
	}
}

// S-0367: a tier whose placeholder comes to nothing, as {packages} does for
// files no Go package holds, is skipped; a ctx already ended runs nothing.
func TestFixSkipsWhatItCannotFill(t *testing.T) {
	root := tree(t, map[string]string{"a.md": "hello\n"})
	tiers := []Tier{{Name: "go", Paths: []string{"**"}, FixCommand: []string{"upper", "{packages}"}}}
	proc := &fixProc{}
	if got := Fix(context.Background(), root, tiers, []string{"a.md"}, proc, nil); got != nil || proc.ran != nil {
		t.Errorf("an empty placeholder: changed %v, ran %q; want neither", got, proc.ran)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tiers[0].FixCommand = []string{"upper", "{files}"}
	if got := Fix(ctx, root, tiers, []string{"a.md"}, proc, nil); got != nil || proc.ran != nil {
		t.Errorf("an ended ctx: changed %v, ran %q; want neither", got, proc.ran)
	}
}

// S-0367: CheckoutFix reads the fix command from the checkout's manifest and
// refuses tiers that are not valid; a command that cannot start changes
// nothing and is no error.
func TestCheckoutFix(t *testing.T) {
	root := checkout(t, "tests:\n  - name: md\n    command: [lint]\n    paths: [\"**/*.md\"]\n    fix_command: [flai-no-such-fixer, \"{files}\"]\n", "a.md")
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"flai-no-such-fixer", "{files}"}; !reflect.DeepEqual(tiers[0].FixCommand, want) {
		t.Errorf("fix command = %q, want %q", tiers[0].FixCommand, want)
	}
	got, err := CheckoutFix(context.Background(), root, []string{"a.md"}, nil)
	if err != nil || got != nil {
		t.Errorf("CheckoutFix = %v, %v; want nothing changed and no error", got, err)
	}

	_, err = CheckoutFix(context.Background(), checkout(t, "tests:\n  - name: md\n    command: [lint]\n    paths: [\"**\"]\n    fix_command: [fix, \"{story}\"]\n"), nil, nil)
	if err == nil || !strings.Contains(err.Error(), "tests[0].fix_command") {
		t.Errorf("a bad fix command: %v, want it refused naming tests[0].fix_command", err)
	}
}

// S-0367: Unfixed runs the command of each tier that carries a fix command
// over the existing files it selects, and answers the first failure with its
// findings; a tier without a fix command, and a path no such tier selects,
// are not run.
func TestUnfixedRunsTheCheckOfTheTiersThatFix(t *testing.T) {
	root := tree(t, map[string]string{"docs/a.md": "a", "flai/b.go": "b"})
	tiers := []Tier{
		{Name: "gofmt", Command: []string{"gofmt", "{files}"}, Paths: []string{"flai/**"}},
		{Name: "markdown", Command: []string{"mdlint", "{files}"}, Paths: []string{"**/*.md"}, FixCommand: []string{"upper", "{files}"}},
	}
	proc := &fixProc{exits: map[string]int{"mdlint": 1}}

	res, err := Unfixed(context.Background(), root, tiers, []string{"docs/a.md", "flai/b.go", "docs/gone.md"}, RunOptions{Proc: proc})
	if err != nil {
		t.Fatal(err)
	}
	if res.Passed || len(res.Tiers) != 1 || res.Tiers[0].Name != "markdown" || res.Tiers[0].State != Failed || len(res.Tiers[0].Findings) == 0 {
		t.Fatalf("result %+v", res)
	}
	if !reflect.DeepEqual(proc.ran, []string{"mdlint docs/a.md"}) {
		t.Errorf("ran %v; want only the markdown tier's command over the file that exists", proc.ran)
	}

	proc = &fixProc{}
	res, err = Unfixed(context.Background(), root, tiers, []string{"flai/b.go"}, RunOptions{Proc: proc})
	if err != nil || !res.Passed || len(proc.ran) != 0 {
		t.Errorf("with no tier that fixes selected: %+v, %v, ran %v", res, err, proc.ran)
	}
}

// remainsProc fakes fix commands for FixRemaining: "mdfix" prints a
// markdownlint finding and exits with exit; "gone" cannot start.
type remainsProc struct {
	exit int
	ran  []string
}

func (f *remainsProc) Run(_ context.Context, _ string, argv, _ []string, stdout, _ io.Writer) (int, error) {
	f.ran = append(f.ran, strings.Join(argv, " "))
	if argv[0] == "gone" {
		return 0, errors.New("executable file not found")
	}
	if f.exit != 0 {
		_, _ = io.WriteString(stdout, "lint-md: markdownlint-cli2 0.20.0\ndocs/a.md:3 MD022/blanks-around-headings Headings should be surrounded by blank lines\n")
	}
	return f.exit, nil
}

// S-0609: FixRemaining runs each tier's fix once, and its exit says what
// remains: 0 passes, non-zero fails with the tier's findings, with no second
// command run to check.
func TestFixRemainingAnswersWhatTheFixExitSays(t *testing.T) {
	root := tree(t, map[string]string{"docs/a.md": "a\n", "main.go": "package main\n"})
	tiers := []Tier{
		{Name: "unit", Command: []string{"go", "test"}, Paths: []string{"**/*.go"}},
		{Name: "markdown", Command: []string{"mdlint", "{files}"}, Paths: []string{"**/*.md"}, FixCommand: []string{"mdfix", "--fix", "{files}"}},
	}

	proc := &remainsProc{}
	o := FixRemaining(context.Background(), root, tiers, []string{"docs/a.md", "main.go"}, proc, nil)
	if !o.Known() || !o.Remains.Passed || len(o.Remains.Tiers) != 1 || o.Remains.Tiers[0].State != Passed {
		t.Errorf("exit 0: %+v, want known and passed", o)
	}
	if want := []string{"mdfix --fix docs/a.md"}; !reflect.DeepEqual(proc.ran, want) {
		t.Errorf("ran %q, want %q: the fix once and no check", proc.ran, want)
	}
	if want := []string{"docs/a.md", "main.go"}; !reflect.DeepEqual(o.Remains.Paths, want) {
		t.Errorf("paths = %v, want %v", o.Remains.Paths, want)
	}

	proc = &remainsProc{exit: 1}
	o = FixRemaining(context.Background(), root, tiers, []string{"docs/a.md"}, proc, nil)
	if !o.Known() || o.Remains.Passed {
		t.Fatalf("exit 1: %+v, want known and failed", o)
	}
	tr := o.Remains.Tiers[0]
	if tr.Name != "markdown" || tr.State != Failed || tr.ExitCode == nil || *tr.ExitCode != 1 || len(tr.Findings) != 1 {
		t.Fatalf("tier %+v", tr)
	}
	if msg := tr.Findings[0].Message; !strings.Contains(msg, "docs/a.md:3 MD022") || !strings.Contains(msg, "exit status 1") {
		t.Errorf("finding %q, want the fix's remaining lint and its exit", msg)
	}
	if len(proc.ran) != 1 {
		t.Errorf("ran %q, want one process", proc.ran)
	}
	// the same findings Unfixed gives from the check
	check, err := Unfixed(context.Background(), root, tiers, []string{"docs/a.md"}, RunOptions{Proc: &remainsProc{exit: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(check.Tiers[0].Findings[0].Message, tr.Findings[0].Message) {
		t.Errorf("fix finding %q differs from the check's %q", tr.Findings[0].Message, check.Tiers[0].Findings[0].Message)
	}
}

// S-0609: a fix command that cannot start, a fix a ctx already ended does not
// run, and a structured format whose output gives no finding of its own are
// unknown, never a pass; a tier no path selects leaves nothing.
func TestFixRemainingUnknownIsNotAPass(t *testing.T) {
	root := tree(t, map[string]string{"a.md": "a\n"})
	tiers := []Tier{{Name: "markdown", Paths: []string{"*.md"}, FixCommand: []string{"gone", "{files}"}}}
	var logs bytes.Buffer
	o := FixRemaining(context.Background(), root, tiers, []string{"a.md"}, &remainsProc{}, slog.New(slog.NewTextHandler(&logs, nil)))
	if o.Known() || o.Remains.Passed || !reflect.DeepEqual(o.Unknown, []string{"markdown"}) {
		t.Errorf("cannot start: %+v, want unknown and not passed", o)
	}
	if !strings.Contains(logs.String(), "level=WARN") {
		t.Errorf("a fix that cannot start is not logged:\n%s", logs.String())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	proc := &remainsProc{}
	tiers[0].FixCommand = []string{"mdfix", "{files}"}
	o = FixRemaining(ctx, root, tiers, []string{"a.md"}, proc, nil)
	if o.Known() || o.Remains.Passed || proc.ran != nil {
		t.Errorf("ended ctx: %+v, ran %q; want unknown and nothing run", o, proc.ran)
	}

	tiers[0].Format = FormatGolangciJSON
	o = FixRemaining(context.Background(), root, tiers, []string{"a.md"}, &remainsProc{exit: 1}, nil)
	if o.Known() || o.Remains.Passed {
		t.Errorf("structured format with no finding of its own: %+v, want unknown", o)
	}

	proc = &remainsProc{exit: 1}
	o = FixRemaining(context.Background(), root, tiers, []string{"b.go"}, proc, nil)
	if !o.Known() || !o.Remains.Passed || len(o.Remains.Tiers) != 0 || proc.ran != nil {
		t.Errorf("nothing selected: %+v, ran %q; want a pass with nothing run", o, proc.ran)
	}
}

// S-0609: CheckoutFixRemaining reads the tiers from the checkout's manifest.
func TestCheckoutFixRemaining(t *testing.T) {
	root := checkout(t, "tests:\n  - name: md\n    command: [lint]\n    paths: [\"**/*.md\"]\n    fix_command: [flai-no-such-fixer, \"{files}\"]\n", "a.md")
	o, err := CheckoutFixRemaining(context.Background(), root, []string{"a.md"}, nil)
	if err != nil || o.Known() {
		t.Errorf("CheckoutFixRemaining = %+v, %v; want a fixer that cannot start unknown, no error", o, err)
	}
	if _, err := CheckoutFixRemaining(context.Background(), checkout(t, "tests:\n  - name: md\n    command: [lint]\n    paths: [\"**\"]\n    fix_command: [fix, \"{story}\"]\n"), nil, nil); err == nil {
		t.Error("a bad fix command is not refused")
	}
}
