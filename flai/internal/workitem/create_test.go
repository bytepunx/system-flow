package workitem

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// S-0199: a story is made a draft when it is created; nothing else is.
func TestCreateDraft(t *testing.T) {
	r := newProject(t)
	s, err := r.Create(NewOptions{Type: Story, Title: "Drafted", Draft: true, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(s.Path)
	if !s.Draft || !strings.Contains(string(data), "\ndraft: true\n") {
		t.Errorf("the story is a draft:\n%s", data)
	}
	if plain, err := r.Create(NewOptions{Type: Story, Title: "Plain", Now: t0}); err != nil || plain.Draft {
		t.Errorf("a story is not a draft unless asked: %v", err)
	}
	for _, typ := range []string{Epic, Task} {
		if _, err := r.Create(NewOptions{Type: typ, Title: "No", Parent: map[string]string{Task: s.ID}[typ], Draft: true, Now: t0}); err == nil || !strings.Contains(err.Error(), "only a story is a draft") {
			t.Errorf("a draft %s must be refused: %v", typ, err)
		}
	}
}

// S-0203: a story is created with a cost of delay, as one made from an issue
// is; a task is not, and a block without who set it is refused.
func TestCreateCostOfDelay(t *testing.T) {
	r := newProject(t)
	cod := &CostOfDelay{Inputs: &CostInputs{TimeLostPerCycle: "2h", By: "flai", At: t0.Format(TimeFormat)}}
	s, err := r.Create(NewOptions{Type: Story, Title: "Costed", CostOfDelay: cod, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	cod.Inputs.TimeLostPerCycle = "9h"
	got, err := r.Get(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if c := got.CostOfDelay; c.IsZero() || c.Inputs.TimeLostPerCycle != "2h" || c.Inputs.By != "flai" || c.Inputs.At != t0.Format(TimeFormat) || c.By != "" {
		t.Errorf("the story carries the cost of delay it was given: %+v", c)
	}
	if _, err := r.Create(NewOptions{Type: Task, Title: "No", Parent: s.ID, CostOfDelay: cod, Now: t0}); err == nil || !strings.Contains(err.Error(), "not a task") {
		t.Errorf("a task with a cost of delay must be refused: %v", err)
	}
	unsigned := &CostOfDelay{Inputs: &CostInputs{TimeLostPerCycle: "2h", At: t0.Format(TimeFormat)}}
	if _, err := r.Create(NewOptions{Type: Story, Title: "Unsigned", CostOfDelay: unsigned, Now: t0}); err == nil || !strings.Contains(err.Error(), "cost_of_delay.inputs.by is required") {
		t.Errorf("a cost of delay without who set it must be refused: %v", err)
	}
}

// I-0071: a touch is a path in the repository or a component name; one that
// starts with a dot is a path like any other, and one that leaves the
// repository, or would split or read as a flag, is refused.
func TestCleanTouches(t *testing.T) {
	got, err := CleanTouches([]string{" .claude/agents/planner.md ", ".github/workflows/", ".dockerignore", "flai/cmd", "", ".github/workflows", "a..b/c", "flai"})
	want := []string{".claude/agents/planner.md", ".github/workflows", ".dockerignore", "flai/cmd", "a..b/c", "flai"}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("CleanTouches = %q, %v; want %q", got, err, want)
	}
	if got, err := CleanTouches([]string{" ", "/"}); err != nil || got != nil {
		t.Errorf("nothing given is nil: %q, %v", got, err)
	}
	for _, c := range []struct{ touch, says string }{
		{"../x", "climbs out of the repository"},
		{"a/../b", "climbs out of the repository"},
		{"..", "climbs out of the repository"},
		{"a/..", "climbs out of the repository"},
		{"/etc", "is absolute"},
		{"a,b", "holds a comma"},
		{"-x", "starts with a dash"},
		{"a\x07b", "holds a control character"},
		{"a\nb", "holds a control character"},
	} {
		_, err := CleanTouches([]string{"flai", c.touch})
		if err == nil || !strings.Contains(err.Error(), c.says) || !strings.Contains(err.Error(), "give a repository path or component name") {
			t.Errorf("touch %q: want an error saying %q and what a touch may be, got %v", c.touch, c.says, err)
		}
	}
}

// I-0071: a story or a task is created keeping a touch that starts with a
// dot, and refused with one that leaves the repository, writing nothing.
func TestCreateTouches(t *testing.T) {
	r := newProject(t)
	s, err := r.Create(NewOptions{Type: Story, Title: "Dotted", Touches: []string{".claude/agents/", "flai/cmd"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := r.Get(s.ID); err != nil || !slices.Equal(got.Touches, []string{".claude/agents", "flai/cmd"}) {
		t.Errorf("the story keeps its touches: %v %v", got, err)
	}
	if _, err := r.Create(NewOptions{Type: Task, Title: "Outside", Parent: s.ID, Touches: []string{"../elsewhere"}, Now: t0}); err == nil || !strings.Contains(err.Error(), `touch "../elsewhere"`) {
		t.Errorf("a touch outside the repository must be refused: %v", err)
	}
	if items, err := r.List(false); err != nil || len(items) != 1 {
		t.Errorf("nothing is written for a refused touch: %d items, %v", len(items), err)
	}
}

// S-0468: a task is created with a file, a path absent from the tree, or a
// folder inside a shared pattern, and refused a folder of the tree holding
// files, writing nothing, with the folder and its file count named; a story
// keeps its folder touch (ADR-0188).
func TestCreateRefusesATaskAFolderTouch(t *testing.T) {
	r := newProject(t)
	r.Manifest.Projects = []manifest.Project{{Name: "flai", Path: "flai", Tags: []string{"cli"}}}
	writeTree(t, r.Root, "flai/cmd/a.go", "flai/cmd/b.go", "design/adrs/0001-a.md")
	if err := os.WriteFile(filepath.Join(r.Root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nclaims:\n  shared: [design/adrs]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := r.Create(NewOptions{Type: Story, Title: "Story", Touches: []string{"flai/cmd"}, Now: t0})
	if err != nil {
		t.Fatalf("a story keeps its folder touch: %v", err)
	}
	for _, c := range []struct {
		touches        []string
		touch, message string
	}{
		{[]string{"flai/cmd", "flai/cmd/a.go"}, "flai/cmd", "T-0001's touch flai/cmd is a folder holding 2 files: a task takes no folder touch"},
		{[]string{"flai/cmd/a.go", "cli"}, "cli", "T-0001's touch cli is the folder flai, holding 2 files: a task takes no folder touch"},
	} {
		_, err := r.Create(NewOptions{Type: Task, Title: "Folder", Parent: s.ID, Touches: c.touches, Now: t0})
		var fe *FolderTouchError
		if !errors.As(err, &fe) || fe.Touch != c.touch || !strings.HasPrefix(err.Error(), c.message) {
			t.Errorf("%v: a task's folder touch must be refused naming %s: %v", c.touches, c.touch, err)
		}
	}
	if items, err := r.List(false); err != nil || len(items) != 1 {
		t.Errorf("nothing is written for a refused task: %d items, %v", len(items), err)
	}
	task, err := r.Create(NewOptions{Type: Task, Title: "Files", Parent: s.ID, Touches: []string{"flai/cmd/a.go", "flai/cmd/new.go", "flai/internal/new", "design/adrs"}, Now: t0})
	if err != nil {
		t.Fatalf("a file, a path to be made, and a shared folder are taken: %v", err)
	}
	if !slices.Equal(task.Touches, []string{"flai/cmd/a.go", "flai/cmd/new.go", "flai/internal/new", "design/adrs"}) {
		t.Errorf("the task keeps its touches: %v", task.Touches)
	}
}

// I-0231: S-0536's title made a file name of 240 bytes, which Claude Code's
// Edit could not write, since its temporary file adds ".tmp.<pid>.<hex>" and
// passes the 255-byte limit; the slug is now cut at 80 bytes on a word.
func TestCreateCutsALongTitlesFileName(t *testing.T) {
	r := newProject(t)
	title := "The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted"
	s, err := r.Create(NewOptions{Type: Story, Title: title, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(s.Path)
	if max := len(s.ID) + 1 + 80 + len(".md"); len(base) > max {
		t.Errorf("the file name is %d bytes, more than %d: %s", len(base), max, base)
	}
	if tmp := base + ".tmp.12345.0123456789abcdef"; len(tmp) > 255 {
		t.Errorf("a temporary file beside it is %d bytes, more than 255: %s", len(tmp), tmp)
	}
	if !strings.HasPrefix(base, s.ID+"-the-orchestrator-runs-") || strings.HasSuffix(strings.TrimSuffix(base, ".md"), "-") {
		t.Errorf("the slug is cut on a word: %s", base)
	}
	if got, err := r.Get(s.ID); err != nil || got.Title != title {
		t.Errorf("the item keeps its whole title: %v", err)
	}
}
