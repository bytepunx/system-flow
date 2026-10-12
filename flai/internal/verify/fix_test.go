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
