package itemnew

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// lintCases holds S-0231's original file and this project's markdownlint
// configuration, beside the MD007 tests that agree with markdownlint-cli2.
const lintCases = "../mdlint/testdata/cases"

// project is an empty project linted with this project's configuration.
func project(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	cfg, err := os.ReadFile(filepath.Join(lintCases, ".markdownlint.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), cfg, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// s0231Body is what S-0231's author wrote below its heading: a goal whose
// list is indented one space.
func s0231Body(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(lintCases, "license-story.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, body, ok := strings.Cut(string(data), "# S-0231 Author LICENSE.md file\n")
	if !ok {
		t.Fatal("license-story.md has no S-0231 heading")
	}
	return body
}

// S-0240: a story made with S-0231's body, as the dashboard's new form makes
// it, is refused with markdownlint's MD007 on each item of the indented list,
// as flai check would report it, and nothing is left behind.
func TestABodyWithAnIndentedListIsRefusedWithMD007(t *testing.T) {
	repo := project(t)
	now := time.Date(2026, 10, 2, 12, 24, 10, 0, time.UTC)
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Licensing", Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	parentWas, _ := os.ReadFile(epic.Path)
	story := func(body string) workitem.NewOptions {
		return workitem.NewOptions{Type: workitem.Story, Title: "Author LICENSE.md file", Nature: "feature", Parent: epic.ID, Owner: "alex", Now: now, Body: body}
	}

	body := s0231Body(t)
	_, err = Create(repo, nil, Options{New: story(body)})
	refused, ok := docedit.IsRefused(err)
	if !ok {
		t.Fatalf("want a refusal, got %v", err)
	}
	var lines []int
	for _, f := range refused.Findings {
		if f.Rule != "markdown.MD007" || f.Path != "wip/kanban/stories/S-0001-author-license-md-file.md" || !strings.Contains(f.Message, "[Expected: 0; Actual: 1]") {
			t.Errorf("a finding other than MD007 on the new story: %+v", f)
		}
		lines = append(lines, f.Line)
	}
	if m, _ := filepath.Glob(filepath.Join(repo.Root, "wip/kanban/stories", "*")); len(m) != 0 {
		t.Errorf("nothing is created: %v", m)
	}
	if now, _ := os.ReadFile(epic.Path); string(now) != string(parentWas) {
		t.Errorf("the parent is left as it was:\n%s", now)
	}

	// The same story with its list unindented is kept, and its items sit on
	// the lines the refusal named.
	res, err := Create(repo, nil, Options{New: story(strings.ReplaceAll(body, "\n - ", "\n- "))})
	if err != nil {
		t.Fatalf("the body with its list unindented is kept: %v", err)
	}
	kept, _ := os.ReadFile(res.Item.Path)
	var want []int
	for i, l := range strings.Split(string(kept), "\n") {
		if strings.HasPrefix(l, "- ") && !strings.HasPrefix(l, "- [ ]") {
			want = append(want, i+1)
		}
	}
	if len(want) != 4 || len(lines) != len(want) {
		t.Fatalf("MD007 on lines %v, want one on each of the goal's four items, lines %v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("MD007 on lines %v, want %v", lines, want)
		}
	}
}

// gitLog is a git that is always inside a work tree and records what it is
// asked to do.
type gitLog struct{ calls [][]string }

func (g *gitLog) Run(_, name string, args ...string) (string, error) {
	g.calls = append(g.calls, append([]string{name}, args...))
	return "", nil
}

func (g *gitLog) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *gitLog) LookPath(name string) (string, error) { return name, nil }

// committed is what the recorded git commit was given after "--".
func (g *gitLog) committed() []string {
	for _, c := range g.calls {
		if len(c) > 1 && c[1] == "commit" {
			for i, a := range c {
				if a == "--" {
					return c[i+1:]
				}
			}
		}
	}
	return nil
}

// S-0203: what Then writes with the new item is committed with it; a path
// outside the repository is not handed to git.
func TestThenIsCommittedWithTheItem(t *testing.T) {
	repo := project(t)
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	issue := filepath.Join(repo.Root, "design", "issues", "I-0001.md")
	g := &gitLog{}
	res, err := Create(repo, g, Options{
		New:        workitem.NewOptions{Type: workitem.Story, Title: "From an issue", Now: now},
		Autocommit: true,
		Then: func(it *workitem.Item) ([]string, error) {
			if err := os.MkdirAll(filepath.Dir(issue), 0o755); err != nil {
				return nil, err
			}
			return []string{issue, "design/issues/summary.md", filepath.Join(t.TempDir(), "elsewhere.md")}, os.WriteFile(issue, []byte(it.ID+"\n"), 0o644)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{res.Path, "design/issues/I-0001.md", "design/issues/summary.md"}
	if got := g.committed(); strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("committed %v, want %v", got, want)
	}
	if !res.Committed {
		t.Errorf("the item is committed: %+v", res)
	}
}

// S-0203: when Then fails, the new item is removed, its parent restored, and
// the error returned.
func TestThenFailingLeavesNothing(t *testing.T) {
	repo := project(t)
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Issues", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	parentWas, _ := os.ReadFile(epic.Path)
	boom := errors.New("the issue file is gone")
	g := &gitLog{}
	_, err = Create(repo, g, Options{
		New:        workitem.NewOptions{Type: workitem.Story, Title: "From an issue", Parent: epic.ID, Now: now},
		Autocommit: true,
		Then:       func(*workitem.Item) ([]string, error) { return nil, boom },
	})
	if !errors.Is(err, boom) {
		t.Fatalf("want the hook's error, got %v", err)
	}
	if m, _ := filepath.Glob(filepath.Join(repo.Root, "wip/kanban/stories", "*")); len(m) != 0 {
		t.Errorf("nothing is created: %v", m)
	}
	if now, _ := os.ReadFile(epic.Path); string(now) != string(parentWas) {
		t.Errorf("the parent is left as it was:\n%s", now)
	}
	if len(g.calls) != 0 {
		t.Errorf("nothing is committed: %v", g.calls)
	}
}
