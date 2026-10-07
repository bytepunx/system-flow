package check

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var now = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

func TestGoodFixtureIsClean(t *testing.T) {
	repo, err := workitem.Open("../metrics/testdata/good")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	// Its epic is in backlog with a story in progress, from before an epic
	// followed its stories (S-0200): the one finding, advisory.
	for _, f := range res.Findings {
		if f.Rule != "epic.lags-stories" || f.Path != "wip/kanban/epics/E-001-epic.md" {
			t.Errorf("unexpected finding: %+v", f)
		}
	}
	if len(res.Findings) != 1 || res.Advisory != 1 || !res.OK(true) || res.Items != 10 {
		t.Fatalf("result: %+v", res)
	}
}

func TestBadFixtureFindings(t *testing.T) {
	repo, err := workitem.Open("testdata/bad")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{ // rule -> substring of path
		"story.tasks":         "S-004",            // review with no tasks
		"manifest.key":        "system-flow.yaml", // API responses name the project by it (ADR-0024)
		"story.criteria":      "S-001",
		"narrative.missing":   "S-001",
		"item.parent-missing": "T-001",
		"item.stream":         "T-001",
		"item.done-children":  "E-001",
		"item.archive":        "E-001",
		"item.filename":       "S-003-wrong-name",
		"item.chronology":     "S-003-wrong-name",
		"item.front-matter":   "S-003-wrong-name",
		"item.sequence":       "S-003-wrong-name",
		"item.parent-list":    "E-001",
		"board.order":         "board.md",
		"narrative.stream":    "S-007",
		"narrative.section":   "S-007",
		"narrative.index":     "index.md",
		"doc.front-matter":    "nofm.md",
		"doc.updated":         "docs/users/index.md",
		"adr.id":              "0001-first.md",
		"adr.decision":        "0001-first.md", // no ## Decision to brief it by (ADR-0049)
	}
	got := map[string][]string{}
	for _, f := range res.Findings {
		got[f.Rule] = append(got[f.Rule], f.Path)
		if f.Line < 1 {
			t.Errorf("finding without line: %+v", f)
		}
	}
	for rule, path := range want {
		found := false
		for _, p := range got[rule] {
			if strings.Contains(p, path) {
				found = true
			}
		}
		if !found {
			t.Errorf("expected %s on %s, got %v", rule, path, got[rule])
		}
	}
	// S-001 is in-progress and S-002 is ready: tasks are not required until
	// review (ADR-0021), so neither may raise story.tasks.
	for _, p := range got["story.tasks"] {
		if !strings.Contains(p, "S-004") {
			t.Errorf("story.tasks raised before review on %s", p)
		}
	}
	if res.OK(false) {
		t.Error("bad fixture must fail")
	}
	// line numbers point at the offending key
	for _, f := range res.Findings {
		if f.Rule == "item.parent-missing" && f.Line != 7 {
			t.Errorf("parent-missing should point at the parent line (7), got %d", f.Line)
		}
	}
}

func TestMonorepoIsClean(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	repo, err := workitem.Open("../../..")
	if err != nil {
		t.Skip("monorepo not present")
	}
	res, err := Run(repo, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	// In a close-out (scripts/close-out.sh exports the story), the errors
	// outside the story are main's, logged as scripts/check.sh notes them,
	// and do not fail the story (S-0249).
	story := os.Getenv("CLOSE_OUT_STORY")
	if story != "" {
		changed, err := storygit.StoryChanges(execx.System{}, repo, story)
		if err != nil {
			t.Fatal(err)
		}
		if err := ScopeToStory(res, repo, story, changed); err != nil {
			t.Fatal(err)
		}
	}
	// Warnings (a done story awaiting archive, a WIP breach) are board state,
	// not code defects; the strict check job gates them. Errors fail here.
	for _, f := range res.Findings {
		switch {
		case f.Outside:
			t.Logf("outside %s: %s:%d: %s: %s: %s", story, f.Path, f.Line, f.Level, f.Rule, f.Message)
		case f.Level == Error:
			t.Errorf("%s:%d: %s: %s: %s", f.Path, f.Line, f.Level, f.Rule, f.Message)
		default:
			t.Logf("warning: %s:%d: %s: %s", f.Path, f.Line, f.Rule, f.Message)
		}
	}
}

func TestCacheMustBeIgnored(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", ".flai-cache"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	has := func(rule string) bool {
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range res.Findings {
			if f.Rule == rule {
				return true
			}
		}
		return false
	}
	if !has("layout.gitignore") {
		t.Error("expected layout.gitignore without a .gitignore")
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n"), 0o644)
	if !has("layout.gitignore") {
		t.Error("expected layout.gitignore when .flai-cache is not listed")
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte("bin/\n.flai-cache/\n"), 0o644)
	if has("layout.gitignore") {
		t.Error("ignored cache should pass")
	}
}

func TestTouchesOverlap(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	e := mustItem(t, repo, workitem.Epic, "E", "")
	s1 := mustItem(t, repo, workitem.Story, "One", e.ID)
	s2 := mustItem(t, repo, workitem.Story, "Two", e.ID)
	for _, s := range []*workitem.Item{s1, s2} {
		s.Status = workitem.InProgress
		s.Touches = []string{"flaiover/src/lib"}
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	s2.Touches = []string{"flaiover/src/lib/server/auth.ts"}
	_ = repo.Save(s2)
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want := "S-0001 touches flaiover/src/lib, which S-0002 (in progress) also touches as flaiover/src/lib/server/auth.ts"
	if got := overlapFindings(res); len(got) != 1 || got[0].Message != want || !slices.Equal(got[0].Stories, []string{"S-0001", "S-0002"}) || !strings.HasSuffix(got[0].Path, filepath.Base(s1.Path)) {
		t.Errorf("wip.overlap = %+v, want one on S-0001: %s", got, want)
	}
	s2.Touches = []string{"docs"}
	_ = repo.Save(s2)
	res, _ = Run(repo, now)
	for _, f := range res.Findings {
		if f.Rule == "wip.overlap" {
			t.Errorf("disjoint paths must not overlap: %+v", f)
		}
	}
}

// ADR-0096: an overlap whose narrower entry lies wholly inside a shared path
// is no wip.overlap, one outside every shared path still is, and a pattern
// that is not valid is a manifest.claims error on its own line. The list is
// read as the file has it when check runs, not as the repo was opened.
func TestTouchesOverlapInsideASharedPath(t *testing.T) {
	root := t.TempDir()
	base := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	mf := filepath.Join(root, "system-flow.yaml")
	_ = os.WriteFile(mf, []byte(base), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	e := mustItem(t, repo, workitem.Epic, "E", "")
	s1 := mustItem(t, repo, workitem.Story, "One", e.ID)
	s2 := mustItem(t, repo, workitem.Story, "Two", e.ID)
	s1.Touches = []string{"docs/users/flai.md", "design/adrs/0096-x.md", "flai/cmd/x.go"}
	s2.Touches = []string{"docs/users", "design", "flai/cmd"}
	for _, s := range []*workitem.Item{s1, s2} {
		s.Status = workitem.InProgress
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	overlaps := func() []string {
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, f := range res.Findings {
			if f.Rule == "wip.overlap" {
				got = append(got, f.Message)
			}
		}
		return got
	}
	// ADR-0115: one finding for the pair, naming every overlapping entry
	whole := "S-0001 touches docs/users/flai.md, design/adrs/0096-x.md and flai/cmd/x.go, which S-0002 (in progress) also touches as docs/users, design and flai/cmd"
	if got := overlaps(); len(got) != 1 || got[0] != whole {
		t.Fatalf("with no shared paths, overlaps = %q\nwant only %s", got, whole)
	}

	shared := base + "claims:\n  shared:\n    - docs/users/flai.md\n    - design/adrs\n    - ../outside\n"
	if err := os.WriteFile(mf, []byte(shared), 0o644); err != nil {
		t.Fatal(err)
	}
	want := "S-0001 touches flai/cmd/x.go, which S-0002 (in progress) also touches as flai/cmd"
	if got := overlaps(); len(got) != 1 || got[0] != want {
		t.Errorf("with docs/users/flai.md and design/adrs shared, overlaps = %q\nwant only %s", got, want)
	}

	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var claims []Finding
	for _, f := range res.Findings {
		if f.Rule == "manifest.claims" {
			claims = append(claims, f)
		}
	}
	if len(claims) != 1 || claims[0].Level != Error || claims[0].Path != "system-flow.yaml" || claims[0].Line != 12 || !strings.Contains(claims[0].Message, `"../outside" has a .. segment`) {
		t.Errorf("manifest.claims findings = %+v, want one error on line 12 naming ../outside", claims)
	}
}

// I-0076, ADR-0115 (the S-0223 instance): a story whose folder touch its task
// narrows to one file inside it claims that file, so another story in
// progress on a different file in the folder is no wip.overlap, as it is no
// hold.
func TestOverlapComparesClaimsNotTouches(t *testing.T) {
	repo := overlapProject(t)
	e := mustItem(t, repo, workitem.Epic, "E", "")
	a := mustItem(t, repo, workitem.Story, "A", e.ID)
	b := mustItem(t, repo, workitem.Story, "B", e.ID)
	task := mustItem(t, repo, workitem.Task, "Serve a", a.ID)
	a.Touches = []string{"flai/internal/serve"}
	b.Touches = []string{"flai/internal/serve/b.go"}
	task.Touches = []string{"flai/internal/serve/a.go"}
	for _, it := range []*workitem.Item{a, b, task} {
		it.Status = workitem.InProgress
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	if got := overlapFindings(res); len(got) != 0 {
		t.Errorf("a claim narrowed to flai/internal/serve/a.go overlaps flai/internal/serve/b.go: %+v", got)
	}

	// without the task's narrowing, the folder is claimed whole
	task.Touches = nil
	if err := repo.Save(task); err != nil {
		t.Fatal(err)
	}
	res, err = Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want := "S-0001 touches flai/internal/serve, which S-0002 (in progress) also touches as flai/internal/serve/b.go"
	if got := overlapFindings(res); len(got) != 1 || got[0].Message != want {
		t.Errorf("wip.overlap = %+v, want %s", got, want)
	}
}

// I-0076, ADR-0115 (the S-0262 and T-0877 instance): two stories in progress
// on one path, with a task in progress of one on it too, are one wip.overlap
// naming both stories: the task is part of its story's claim, not compared on
// its own.
func TestOverlapIsOneFindingPerPairOfStories(t *testing.T) {
	repo := overlapProject(t)
	e := mustItem(t, repo, workitem.Epic, "E", "")
	a := mustItem(t, repo, workitem.Story, "A", e.ID)
	b := mustItem(t, repo, workitem.Story, "B", e.ID)
	task := mustItem(t, repo, workitem.Task, "B's task", b.ID)
	for _, it := range []*workitem.Item{a, b, task} {
		it.Status = workitem.InProgress
		it.Touches = []string{"design/system/workflow.md"}
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want := "S-0001 touches design/system/workflow.md, which S-0002 (in progress) also touches as design/system/workflow.md"
	got := overlapFindings(res)
	if len(got) != 1 || got[0].Message != want || !slices.Equal(got[0].Stories, []string{"S-0001", "S-0002"}) || !strings.HasSuffix(got[0].Path, filepath.Base(a.Path)) {
		t.Fatalf("wip.overlap = %+v, want one on S-0001: %s", got, want)
	}
	// the inbox's reading is the same finding
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	if inbox := Overlaps(repo, items); len(inbox) != 1 || inbox[0].Message != want || !slices.Equal(inbox[0].Stories, got[0].Stories) {
		t.Errorf("Overlaps = %+v, want the check's %+v", inbox, got)
	}
}

// overlapProject is an empty project to put stories in progress in.
func overlapProject(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// overlapFindings is the wip.overlap findings of res.
func overlapFindings(res *Result) []Finding {
	var out []Finding
	for _, f := range res.Findings {
		if f.Rule == "wip.overlap" {
			out = append(out, f)
		}
	}
	return out
}

// S-0130, ADR-0046: an after: entry that names no story, the story itself,
// or something other than a story is an error, and so is a cycle, once.
func TestAfterRules(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	e := mustItem(t, repo, workitem.Epic, "E", "")
	s1 := mustItem(t, repo, workitem.Story, "One", e.ID)
	s2 := mustItem(t, repo, workitem.Story, "Two", e.ID)
	s3 := mustItem(t, repo, workitem.Story, "Three", e.ID)
	mustItem(t, repo, workitem.Task, "Task", s1.ID) // T-0001, which a story's after cannot name
	set := func(it *workitem.Item, after ...string) {
		it.After = after
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	findings := func() []string {
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, f := range res.Findings {
			if f.Rule == "story.after" || (f.Rule == "item.front-matter" && strings.Contains(f.Message, "after")) {
				out = append(out, filepath.Base(f.Path)[:6]+": "+f.Message)
			}
		}
		return out
	}

	set(s2, "S-0001")
	set(s3, "S-0001", "S-0002")
	if got := findings(); len(got) != 0 {
		t.Errorf("sound after: reported %v", got)
	}

	set(s1, "S-0009", "S-0001")
	set(e, "S-0002")
	got := strings.Join(findings(), "\n")
	for _, want := range []string{
		"S-0001: after names S-0009, which does not exist",
		"S-0001: after names S-0001 itself",
		"E-0001: after is for stories and tasks, and this is an epic",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}

	set(e)
	set(s1, "T-0001")
	got = strings.Join(findings(), "\n")
	for _, want := range []string{
		`after[0] "T-0001" is not a story ID`,
		"S-0001: after names T-0001, a task; a story's after names stories",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("a task in a story's after: missing %q in\n%s", want, got)
		}
	}

	// S-0001 waits for S-0003, which waits for S-0001 and S-0002, and S-0002
	// waits for S-0001: two cycles through S-0001, reported on it only.
	set(s1, "S-0003")
	got = strings.Join(findings(), "\n")
	if !strings.Contains(got, "S-0001: after forms a cycle, S-0001 waits for S-0003 waits for S-0001") {
		t.Errorf("cycle not reported on S-0001:\n%s", got)
	}
	if n := strings.Count(got, "forms a cycle"); n != 1 {
		t.Errorf("a cycle is reported once, got %d:\n%s", n, got)
	}
}

// S-0176: a task's after: names tasks of its own story. An entry that names
// no task, a story, a task of another story, or the task itself is an error,
// and so is a cycle among a story's tasks, once; a story's after: is not
// confused with its tasks'.
func TestTaskAfterRules(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	s1 := mustItem(t, repo, workitem.Story, "One", "")
	s2 := mustItem(t, repo, workitem.Story, "Two", "")
	t1 := mustItem(t, repo, workitem.Task, "First", s1.ID)
	t2 := mustItem(t, repo, workitem.Task, "Second", s1.ID)
	t3 := mustItem(t, repo, workitem.Task, "Third", s1.ID)
	mustItem(t, repo, workitem.Task, "Elsewhere", s2.ID) // T-0004
	set := func(it *workitem.Item, after ...string) {
		it.After = after
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	findings := func() string {
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, f := range res.Findings {
			if strings.HasSuffix(f.Rule, ".after") || (f.Rule == "item.front-matter" && strings.Contains(f.Message, "after")) {
				out = append(out, filepath.Base(f.Path)[:6]+": "+f.Rule+": "+f.Message)
			}
		}
		return strings.Join(out, "\n")
	}

	// T-0002 and T-0003 wait for T-0001, T-0003 for T-0002 too: a plan of
	// three layers, which is sound; so is a story that waits for another.
	set(t2, "T-0001")
	set(t3, "T-0001", "T-0002")
	set(s2, "S-0001")
	if got := findings(); got != "" {
		t.Errorf("sound after: reported\n%s", got)
	}

	set(t1, "T-0009", "T-0001", "T-0004", "S-0002")
	got := findings()
	for _, want := range []string{
		"T-0001: task.after: after names T-0009, which does not exist; fix it or clear it (flai edit T-0001 --clear-after)",
		"T-0001: task.after: after names T-0001 itself; a task cannot wait for itself",
		"T-0001: task.after: after names T-0004, a task of S-0002; a task waits only for tasks of its own story, S-0001",
		"T-0001: task.after: after names S-0002, a story; a task's after names tasks",
		`T-0001: item.front-matter: after[3] "S-0002" is not a task ID like T-0001`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if strings.Contains(got, "cycle") {
		t.Errorf("a task of another story or the task itself is not a cycle:\n%s", got)
	}

	// T-0001 waits for T-0003, which waits for T-0001 and T-0002, and T-0002
	// waits for T-0001: two cycles through T-0001, reported on it only.
	set(t1, "T-0003")
	got = findings()
	if !strings.Contains(got, "T-0001: task.after: after forms a cycle, T-0001 waits for T-0003 waits for T-0001") {
		t.Errorf("cycle not reported on T-0001:\n%s", got)
	}
	if n := strings.Count(got, "forms a cycle"); n != 1 {
		t.Errorf("a cycle is reported once, got %d:\n%s", n, got)
	}
}

func mustItem(t *testing.T, r *workitem.Repo, typ, title, parent string) *workitem.Item {
	t.Helper()
	it, err := r.Create(workitem.NewOptions{Type: typ, Title: title, Parent: parent, Owner: "t", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	return it
}

func TestThreadRules(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "docs/guide.md"), []byte("---\ntitle: Guide\n---\n\n# Guide\n\n## Install\nx\n"), 0o644)
	fm := "---\nid: %s\ntitle: T\nanchor:\n  path: %s\n%sstatus: %s\nparticipants: [alex]\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\n---\n\n# %s T\n\n## Entries\n\n### 2026-09-01T12:00:00Z alex\nhi\n"
	w := func(name, body string) {
		_ = os.WriteFile(filepath.Join(root, "wip/threads", name), []byte(body), 0o644)
	}
	w("TH-0001-ok.md", fmt.Sprintf(fm, "TH-0001", "docs/guide.md", "  heading: Install\n", "open", "TH-0001"))
	w("TH-0002-gone.md", fmt.Sprintf(fm, "TH-0002", "docs/missing.md", "", "open", "TH-0002"))
	w("TH-0003-heading.md", fmt.Sprintf(fm, "TH-0003", "docs/guide.md", "  heading: Removed\n", "answered", "TH-0003"))
	w("TH-0004-item.md", fmt.Sprintf(fm, "TH-0004", "wip/kanban/stories/S-0009-x.md", "  item: S-0009\n", "open", "TH-0004"))
	w("TH-0005-bad.md", fmt.Sprintf(fm, "TH-0005", "docs/guide.md", "", "pending", "TH-0005"))
	// an item anchor follows the item into the archive; the stale path is not an error
	_ = os.MkdirAll(filepath.Join(root, "wip/archive/kanban/epics"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "wip/archive/kanban/epics/E-0001-old.md"), []byte("---\nid: E-0001\ntype: epic\nnature: feature\ntitle: Old\nstatus: done\nowner: a\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\ntransitions:\n  - to: done\n    at: 2026-09-01T12:00:00Z\n    by: a\ntags: []\n---\n\n# E-0001 Old\n\n## Outcome\nx\n\n## Stories\n\n## Notes\n"), 0o644)
	w("TH-0006-archived.md", fmt.Sprintf(fm, "TH-0006", "wip/kanban/epics/E-0001-old.md", "  item: E-0001\n", "resolved", "TH-0006"))
	w("TH-0007-left-open.md", fmt.Sprintf(fm, "TH-0007", "wip/kanban/epics/E-0001-old.md", "  item: E-0001\n", "answered", "TH-0007"))
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, f := range res.Findings {
		if !strings.Contains(got[f.Rule], filepath.Base(f.Path)) {
			got[f.Rule] += filepath.Base(f.Path) + " "
		}
		// threads have no archive folder: the finding names the command
		if f.Rule == "threads.archived" && f.Message != "TH-0007 is answered but E-0001 is archived; resolve it with flai thread resolve TH-0007" {
			t.Errorf("threads.archived message: %q", f.Message)
		}
	}
	for rule, want := range map[string]string{
		"threads.anchor":       "TH-0002-gone.md TH-0004-item.md",
		"threads.heading":      "TH-0003-heading.md",
		"threads.front-matter": "TH-0005-bad.md",
		"threads.archived":     "TH-0007-left-open.md",
	} {
		if strings.TrimSpace(got[rule]) != want {
			t.Errorf("%s: got %q want %q", rule, got[rule], want)
		}
	}
	if strings.Contains(got["threads.anchor"], "TH-0001") || strings.Contains(got["threads.front-matter"], "TH-0001") || strings.Contains(got["threads.anchor"], "TH-0006") || strings.Contains(got["threads.archived"], "TH-0006") {
		t.Errorf("the good and archived threads must pass: %v", got)
	}
}

// ADR-0047: design and tech files declare topics, and every topic on a
// convention, design, tech, or ADR file is one a story can have.
func TestTopicRules(t *testing.T) {
	repo, err := workitem.Open("testdata/bad")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		if f.Rule == "doc.topics" || f.Rule == "doc.topic" {
			if f.Level != Warning {
				t.Errorf("topic findings are warnings: %+v", f)
			}
			msg, _, _ := strings.Cut(f.Message, " (")
			got = append(got, fmt.Sprintf("%s:%d %s %s", f.Path, f.Line, f.Rule, msg))
		}
	}
	want := []string{
		`design/adrs/0001-first.md:10 doc.topic topic "adr-nope" on heading "Scope" is used by nothing: not all, code, a sub-project's name, tag, or kind in system-flow.yaml, nor a story's or epic's topics`,
		`design/conventions/session-start.md:7 doc.topic topic "conv-nope" in the front matter is used by nothing: not all, code, a sub-project's name, tag, or kind in system-flow.yaml, nor a story's or epic's topics`,
		`design/system/topical.md:4 doc.topic topic "nope" in the front matter is used by nothing: not all, code, a sub-project's name, tag, or kind in system-flow.yaml, nor a story's or epic's topics`,
		`design/system/topical.md:13 doc.topic topic "also-nope" on heading "Part" is used by nothing: not all, code, a sub-project's name, tag, or kind in system-flow.yaml, nor a story's or epic's topics`,
		`design/tech/untopical.md:1 doc.topics no topics; say which stories it is for, such as topics: [all] or a sub-project's name, tag, or kind`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("topic findings:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestVocabularyIsTheManifestsWords(t *testing.T) {
	repo, err := workitem.Open("../../..")
	if err != nil {
		t.Skip("monorepo not present")
	}
	v := Vocabulary(repo, nil)
	for _, w := range []string{"all", "code", "flai", "go", "cli", "flaiover", "sveltekit", "dashboard", "template", "conventions"} {
		if !v[w] {
			t.Errorf("%s is a topic here", w)
		}
	}
	if len(v) != 10 {
		t.Errorf("vocabulary = %v", v)
	}
}

// S-0135: the topics stories and epics declare are topics a document may
// name; the fixture's epic says logging and a story release, which
// topical.md's last heading uses without a finding.
func TestVocabularyHasTheItemsTopics(t *testing.T) {
	repo, err := workitem.Open("testdata/bad")
	if err != nil {
		t.Fatal(err)
	}
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	v := Vocabulary(repo, items)
	for _, w := range []string{"logging", "release", "all", "code"} {
		if !v[w] {
			t.Errorf("%s is a topic of the fixture: %v", w, v)
		}
	}
	if v["nope"] {
		t.Errorf("nope is no one's topic")
	}
}

// S-0179: markdown in wip is linted with the project's markdownlint
// configuration, as warnings; without one, nothing is linted.
func TestMarkdownRules(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	thread := "---\nid: TH-0001\ntitle: T\nanchor:\n  path: docs/guide.md\nstatus: open\nparticipants: [alex]\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\n---\n\n# TH-0001 T\n\n## Entries\n\n### 2026-09-01T12:00:00Z alex\nhi\n\n### 2026-09-01T12:00:00Z alex\nagain\n"
	_ = os.WriteFile(filepath.Join(root, "docs/guide.md"), []byte("---\ntitle: Guide\n---\n\n# Guide\n\n# Guide\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "wip/threads/TH-0001-t.md"), []byte(thread), 0o644)
	run := func() []Finding {
		t.Helper()
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var out []Finding
		for _, f := range res.Findings {
			if strings.HasPrefix(f.Rule, "markdown.") {
				out = append(out, f)
			}
		}
		return out
	}
	if got := run(); len(got) != 0 {
		t.Fatalf("no markdownlint configuration, no lint: %v", got)
	}
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\n"), 0o644)
	got := run()
	if len(got) != 1 {
		t.Fatalf("want the thread's duplicate heading only (docs are not wip): %v", got)
	}
	f := got[0]
	if f.Level != Warning || f.Rule != "markdown.MD024" || f.Line != 19 || filepath.Base(f.Path) != "TH-0001-t.md" || !strings.Contains(f.Message, "MD024/no-duplicate-heading Multiple headings with the same content") {
		t.Errorf("finding: %+v", f)
	}
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: [oops\n"), 0o644)
	if got := run(); len(got) != 1 || got[0].Rule != "markdown.config" {
		t.Errorf("an unreadable configuration is reported once: %v", got)
	}
}

// S-0253, I-0066: a conflict git-added unresolved reached main in a design
// document. Each line that opens or closes a conflict, in the markdown under
// design, docs, and wip and at the root, is an error; the divider is not, as
// alone on a line it is a setext heading underline. The markers are built,
// not written out, so that no line of this file begins with one.
func TestConflictMarkersAreErrors(t *testing.T) {
	ours, base, divide, theirs := "<<<<<<"+"<", "||||||"+"|", "======"+"=", ">>>>>>"+">"
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive/agents", "wip/threads"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	w := func(rel, body string) {
		_ = os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644)
	}
	// I-0066 as it reached main: lines 7 and 11 open and close the conflict,
	// line 9 divides it.
	w("design/system/work-hierarchy.md", "---\ntitle: Work hierarchy\nupdated: 2026-09-01\ntopics: [all]\n---\n\n"+ours+" HEAD\nmine\n"+divide+"\ntheirs\n"+theirs+" f0f5443 (docs: [S-0201] x)\n")
	w("docs/users/merge.md", "---\ntitle: Merge\nupdated: 2026-09-01\n---\n\n"+ours+" HEAD\na\n"+base+" base\nb\n"+divide+"\nc\n"+theirs+" topic\n")
	w("docs/users/guide.md", "---\ntitle: Guide\nupdated: 2026-09-01\n---\n\nGuide\n"+divide+"\n\nNo conflict here.\n")
	w("wip/archive/agents/S-0001.md", "# S-0001\n\n"+ours+"\n"+theirs+"\n")
	w("README.md", "# Project\n\n"+theirs+" main\n")
	w("notes.txt", ours+" HEAD\n")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range res.Findings {
		if f.Rule != "markdown.conflict-marker" {
			continue
		}
		if f.Level != Error || !strings.Contains(f.Message, "merge conflict marker") {
			t.Errorf("finding: %+v", f)
		}
		got = append(got, fmt.Sprintf("%s:%d", filepath.ToSlash(f.Path), f.Line))
	}
	want := []string{
		"README.md:3",
		"design/system/work-hierarchy.md:7",
		"design/system/work-hierarchy.md:11",
		"docs/users/merge.md:6",
		"docs/users/merge.md:8",
		"docs/users/merge.md:12",
		"wip/archive/agents/S-0001.md:3",
		"wip/archive/agents/S-0001.md:4",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("conflict markers:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if res.Errors != len(want) || res.OK(false) {
		t.Errorf("the markers are the fixture's only errors and fail the check without --strict: %+v", res)
	}
}

// S-0181: the listing paths read past front-matter fields this flai does not
// know, and check still reports each, on its line.
func TestUnknownFieldsAreReported(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "design/issues", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "docs/guide.md"), []byte("---\ntitle: Guide\n---\n\n# Guide\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/epics/E-0001-e.md"), []byte("---\nid: E-0001\ntype: epic\nnature: feature\ntitle: E\nstatus: backlog\nowner: a\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\ntransitions: []\ntags: []\nholds: [S-0002]\n---\n\n# E-0001 E\n\n## Outcome\nx\n\n## Stories\n\n## Notes\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "wip/threads/TH-0001-t.md"), []byte("---\nid: TH-0001\ntitle: T\nanchor:\n  path: docs/guide.md\nstatus: open\nparticipants: [alex]\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\npriority: high\n---\n\n# TH-0001 T\n\n## Entries\n\n### 2026-09-01T12:00:00Z alex\nhi\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "design/issues/I-0001-i.md"), []byte("---\nid: I-0001\ntitle: I\nclass: defect\nstatus: open\ncount: 1\nfirst_reported: 2026-09-01T12:00:00Z\nlast_reported: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\nseverity: 2\n---\n\n# I-0001 I\n"), 0o644)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatalf("check refused the tree instead of reporting: %v", err)
	}
	want := map[string]string{
		"item.unknown-field":    `E-0001-e.md:12 field "holds"`,
		"threads.unknown-field": `TH-0001-t.md:10 field "priority"`,
		"issues.unknown-field":  `I-0001-i.md:10 field "severity"`,
	}
	for _, f := range res.Findings {
		w, ok := want[f.Rule]
		if !ok {
			continue
		}
		got := fmt.Sprintf("%s:%d %s", filepath.Base(f.Path), f.Line, f.Message)
		if f.Level != Error || !strings.HasPrefix(got, w) {
			t.Errorf("%s: got %s %q, want an error starting %q", f.Rule, f.Level, got, w)
		}
		delete(want, f.Rule)
	}
	for rule := range want {
		t.Errorf("%s not reported", rule)
	}
}

// S-0271: the narrative of a story in progress or in review has its Current
// state and Next steps written, as the close-out reads them, and its Next
// steps as a list; each finding is on the section's heading and advisory, so
// a narrative just opened fails neither --strict nor flai task done. Stories
// in other states are not asked.
func TestNarrativeState(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	s := mustItem(t, repo, workitem.Story, "One", "")
	if _, err := repo.OpenStream(s, workitem.StreamOptions{Agent: "a", Now: now}); err != nil {
		t.Fatal(err)
	}
	path := repo.NarrativePath(s.ID)
	opened, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	head, _, _ := strings.Cut(string(opened), "## Context")
	narrative := func(current, next string) string {
		return head + "## Context\n\n## Current state\n\n" + current + "\n\n## Next steps\n\n" + next + "\n\n## Decisions\n\n## Open questions\n\n## Log\n"
	}
	set := func(status string) {
		t.Helper()
		s.Status = status
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	// found is the narrative.state findings, each "line: message".
	found := func() string {
		t.Helper()
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, f := range res.Findings {
			if f.Rule != "narrative.state" {
				continue
			}
			if f.Level != Warning || !f.advisory || f.Path != "wip/agents/S-0001.md" {
				t.Errorf("finding: %+v, want an advisory warning on the narrative", f)
			}
			out = append(out, fmt.Sprintf("%d: %s", f.Line, f.Message))
		}
		return strings.Join(out, "\n")
	}
	const (
		current = "S-0001 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0001 --current, or the MCP tool stream_state"
		next    = "S-0001 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0001 --next, or the MCP tool stream_state"
		prose   = "S-0001's narrative's ## Next steps is not a list; write them as one, the very next action first, with flai stream state S-0001 --next, or the MCP tool stream_state"
	)

	set(workitem.InProgress)
	for _, tc := range []struct {
		name, doc string
		want      []string // heading, message, ...: each finding on its heading's line
	}{
		{"as flai stream open leaves it", string(opened), []string{"## Current state", current, "## Next steps", next}},
		{"empty", narrative("", ""), []string{"## Current state", current, "## Next steps", next}},
		{"bare list markers", narrative("-", "1.\n2."), []string{"## Current state", current, "## Next steps", next}},
		{"next steps in prose", narrative("T-0001 is done.", "Finish T-0002, then review."), []string{"## Next steps", prose}},
		{"numbered", narrative("T-0001 is done.", "1. Finish T-0002.\n2. Review."), nil},
		{"bulleted after prose", narrative("T-0001 is done.", "In order:\n\n- Finish T-0002.\n* Review."), nil},
	} {
		if err := os.WriteFile(path, []byte(tc.doc), 0o644); err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(tc.doc, "\n")
		var want []string
		for i := 0; i < len(tc.want); i += 2 {
			want = append(want, fmt.Sprintf("%d: %s", slices.Index(lines, tc.want[i])+1, tc.want[i+1]))
		}
		if got := found(); got != strings.Join(want, "\n") {
			t.Errorf("%s:\n%s\nwant:\n%s", tc.name, got, strings.Join(want, "\n"))
		}
	}

	if err := os.WriteFile(path, opened, 0o644); err != nil {
		t.Fatal(err)
	}
	set(workitem.Review)
	if got := found(); strings.Count(got, "S-0001 is review and") != 2 {
		t.Errorf("a story in review is asked too:\n%s", got)
	}
	for _, status := range []string{workitem.Backlog, workitem.Ready, workitem.Done} {
		set(status)
		if got := found(); got != "" {
			t.Errorf("a story %s is not asked:\n%s", status, got)
		}
	}
}

// S-0243, I-0007: the review column over its limit still warns, but --strict
// passes over it, since only the operator's acceptance clears it; ready and
// in-progress over their limits still fail --strict.
func TestReviewOverItsLimitDoesNotFailStrict(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte("---\ntitle: Board\nupdated: 2026-09-01\nstatus: active\nwip_limits:\n  ready: 1\n  in-progress: 1\n  review: 3\norder: []\n---\n\n# Board\n"), 0o644)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var stories []*workitem.Item
	for i := range 4 {
		stories = append(stories, mustItem(t, repo, workitem.Story, fmt.Sprintf("Story %d", i+1), ""))
	}
	set := func(status ...string) {
		for i, s := range stories {
			s.Status = status[i]
			if err := repo.Save(s); err != nil {
				t.Fatal(err)
			}
		}
	}
	// board runs the board rules alone, so the stories' other findings
	// (no tasks, no narrative) do not decide OK.
	board := func() *Result {
		items, err := repo.List(true)
		if err != nil {
			t.Fatal(err)
		}
		c := &checker{repo: repo, now: now, items: items, byID: map[string]*workitem.Item{}, res: &Result{}}
		c.board()
		return c.res
	}

	set(workitem.Review, workitem.Review, workitem.Review, workitem.Review)
	res := board()
	if len(res.Findings) != 1 || res.Findings[0].Level != Warning || res.Findings[0].Rule != "board.wip-limit" || res.Findings[0].Message != "4 stories in review, limit 3" {
		t.Fatalf("review over its limit should warn: %+v", res.Findings)
	}
	if res.Warnings != 1 || res.Advisory != 1 {
		t.Errorf("warnings %d advisory %d, want 1 and 1", res.Warnings, res.Advisory)
	}
	if !res.OK(true) || !res.OK(false) {
		t.Errorf("review over its limit must not fail --strict: %+v", res)
	}
	full, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	if full.Advisory != 1 {
		t.Errorf("Run: advisory %d, want 1", full.Advisory)
	}

	for _, tc := range []struct {
		name   string
		status []string
		msg    string
	}{
		{"ready", []string{workitem.Ready, workitem.Ready, workitem.Backlog, workitem.Backlog}, "2 stories in ready, limit 1"},
		{"in-progress", []string{workitem.InProgress, workitem.InProgress, workitem.Backlog, workitem.Backlog}, "2 stories in in-progress, limit 1"},
	} {
		set(tc.status...)
		res := board()
		if len(res.Findings) != 1 || res.Findings[0].Rule != "board.wip-limit" || res.Findings[0].Message != tc.msg {
			t.Fatalf("%s over its limit should warn: %+v", tc.name, res.Findings)
		}
		if res.Advisory != 0 || res.OK(true) || !res.OK(false) {
			t.Errorf("%s over its limit must still fail --strict only: %+v", tc.name, res)
		}
	}
}
