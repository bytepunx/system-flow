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

// planned is a project with two stories, S-0001 and S-0002, each with one
// task, T-0001 and T-0002: what the planner finds when it drafts a story's
// tasks.
func planned(t *testing.T, now time.Time) *workitem.Repo {
	t.Helper()
	repo := project(t)
	for _, o := range []workitem.NewOptions{
		{Type: workitem.Story, Title: "Guard the planner", Owner: "alex", Now: now},
		{Type: workitem.Story, Title: "Another story", Owner: "alex", Now: now},
		{Type: workitem.Task, Title: "Write the prompt", Parent: "S-0001", Owner: "alex", Now: now},
		{Type: workitem.Task, Title: "Elsewhere", Parent: "S-0002", Owner: "alex", Now: now},
	} {
		if _, err := repo.Create(o); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

// plannedTask is a task of S-0001 as the planner writes it: its work and
// what done means, a nature, tags, touches, and after, its sibling T-0001.
func plannedTask(now time.Time) workitem.NewOptions {
	return workitem.NewOptions{
		Type: workitem.Task, Title: "Test the guard", Nature: "improvement", Parent: "S-0001", Owner: "alex", Now: now,
		Tags: []string{"planner"}, Touches: []string{"flai/internal/guard"}, After: []string{"T-0001"},
		Body: "## Work\n\nTest that the guard lets the planner run flai task new.\n\n## Done when\n\n- The test fails without the guard's rule.\n\n## Notes\n\nNone.\n",
	}
}

// refusedLeavingNothing fails t unless err is a refusal whose every finding
// is rule's on the new task T-0003 and says what, no task but T-0001 and
// T-0002 is left, and S-0001 is as it was.
func refusedLeavingNothing(t *testing.T, repo *workitem.Repo, err error, rule, what string, storyWas []byte) {
	t.Helper()
	refused, ok := docedit.IsRefused(err)
	if !ok {
		t.Fatalf("want a refusal, got %v", err)
	}
	if len(refused.Findings) == 0 {
		t.Errorf("the refusal names no finding: %+v", refused)
	}
	for _, f := range refused.Findings {
		if f.Rule != rule || f.Path != "wip/kanban/tasks/T-0003-test-the-guard.md" || !strings.Contains(f.Message, what) {
			t.Errorf("a finding other than %s's %q on the new task: %+v", rule, what, f)
		}
	}
	m, _ := filepath.Glob(filepath.Join(repo.Root, "wip/kanban/tasks", "*"))
	for _, p := range m {
		if b := filepath.Base(p); !strings.HasPrefix(b, "T-0001-") && !strings.HasPrefix(b, "T-0002-") {
			t.Errorf("nothing is created: %s", b)
		}
	}
	if now := storyFile(t, repo); string(now) != string(storyWas) {
		t.Errorf("the story is left as it was:\n%s", now)
	}
}

// storyFile is S-0001's file as it is now.
func storyFile(t *testing.T, repo *workitem.Repo) []byte {
	t.Helper()
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(story.Path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// S-0255: a task the planner writes with its work, what done means, a
// nature, tags, touches, and after its sibling passes the check and is kept
// as written, and its story lists it.
func TestAWellFormedTaskIsKept(t *testing.T) {
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	repo := planned(t, now)
	res, err := Create(repo, nil, Options{New: plannedTask(now)})
	if err != nil {
		t.Fatalf("the task is kept: %v", err)
	}
	it, err := repo.Get("T-0003")
	if err != nil {
		t.Fatal(err)
	}
	if res.Path != "wip/kanban/tasks/T-0003-test-the-guard.md" || it.Parent != "S-0001" || it.Nature != "improvement" ||
		strings.Join(it.Tags, ",") != "planner" || strings.Join(it.Touches, ",") != "flai/internal/guard" || strings.Join(it.After, ",") != "T-0001" {
		t.Errorf("kept as written: %s %+v", res.Path, it)
	}
	if !strings.Contains(it.Body, "## Done when\n\n- The test fails without the guard's rule.") {
		t.Errorf("the body is the planner's:\n%s", it.Body)
	}
	if story := storyFile(t, repo); !strings.Contains(string(story), "T-0003") {
		t.Errorf("the story lists the task:\n%s", story)
	}
}

// S-0255: a task whose body the markdown lint rejects, here an indented
// list, is refused with MD007 and leaves nothing behind.
func TestATaskWithAnIndentedListIsRefusedWithMD007(t *testing.T) {
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	repo := planned(t, now)
	storyWas := storyFile(t, repo)
	task := plannedTask(now)
	task.Body = strings.ReplaceAll(task.Body, "\n- ", "\n - ")
	_, err := Create(repo, nil, Options{New: task})
	refusedLeavingNothing(t, repo, err, "markdown.MD007", "[Expected: 0; Actual: 1]", storyWas)
}

// S-0255: a task after a task that does not exist is refused by the check's
// task.after and leaves nothing behind.
func TestATaskAfterATaskThatDoesNotExistIsRefused(t *testing.T) {
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	repo := planned(t, now)
	storyWas := storyFile(t, repo)
	task := plannedTask(now)
	task.After = []string{"T-0009"}
	_, err := Create(repo, nil, Options{New: task})
	refusedLeavingNothing(t, repo, err, "task.after", "after names T-0009, which does not exist", storyWas)
}

// S-0255: a task after a task of another story is refused by the check's
// task.after and leaves nothing behind.
func TestATaskAfterATaskOfAnotherStoryIsRefused(t *testing.T) {
	now := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	repo := planned(t, now)
	storyWas := storyFile(t, repo)
	task := plannedTask(now)
	task.After = []string{"T-0002"}
	_, err := Create(repo, nil, Options{New: task})
	refusedLeavingNothing(t, repo, err, "task.after", "after names T-0002, a task of S-0002", storyWas)
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
