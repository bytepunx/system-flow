package itemedit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0598: an edit runs the rules a work item write can affect, before and
// after, in place of the whole check. These tests hold it to what ADR-0197
// refuses: a new finding of a kept rule inside the item's story refuses the
// edit, and every file is put back.

var ruleAt = time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)

// snapshot is what the files at paths hold, "" for one that is not there.
func snapshot(paths ...string) map[string]string {
	out := map[string]string{}
	for _, p := range paths {
		data, _ := os.ReadFile(p)
		out[p] = string(data)
	}
	return out
}

// sameFiles fails the test on each file that no longer holds what was holds.
func sameFiles(t *testing.T, was map[string]string) {
	t.Helper()
	for p, data := range was {
		now, _ := os.ReadFile(p)
		if string(now) != data {
			t.Errorf("%s was not put back:\n%s", p, now)
		}
	}
}

// refusedOn fails the test unless err refuses with a finding of rule.
func refusedOn(t *testing.T, err error, rule string) {
	t.Helper()
	r, ok := docedit.IsRefused(err)
	if !ok {
		t.Fatalf("want a refusal on %s, got %v", rule, err)
	}
	for _, f := range r.Findings {
		if f.Rule == rule {
			return
		}
	}
	t.Errorf("the refusal names no %s: %+v", rule, r.Findings)
}

// betweenItemWrite has another writer run write between an edit's two runs
// of the item-write rules (I-0182).
func betweenItemWrite(t *testing.T, write func(repo *workitem.Repo)) {
	t.Helper()
	calls := 0
	runCheck = func(repo *workitem.Repo, now time.Time, w check.Write) (*check.Result, error) {
		if calls++; calls == 2 {
			write(repo)
		}
		return check.RunItemWrite(repo, now, w)
	}
	t.Cleanup(func() { runCheck = check.RunItemWrite })
}

// An edit that brings a story.after naming no story is refused, and the
// story and the index are left as they were.
func TestAnEditBringingAKeptRulesFindingIsRefused(t *testing.T) {
	repo, a, _, _, _ := raceRepo(t)
	was := snapshot(a.Path, filepath.Join(repo.AgentsDir(), "index.md"), repo.NarrativePath(a.ID))
	_, err := Apply(repo, nil, a.ID, Change{After: &[]string{"S-0099"}}, Options{By: "alex", NoCommit: true, Now: ruleAt})
	refusedOn(t, err, "story.after")
	sameFiles(t, was)
}

// A retitle that brings a finding is refused: the old file is back, the new
// one is gone, and the narrative and the file whose link it rewrote are as
// they were.
func TestARetitleBringingAKeptRulesFindingIsRefused(t *testing.T) {
	repo, a, task, _, _ := raceRepo(t)
	link := "\nSee [A](../stories/" + filepath.Base(a.Path) + ").\n"
	data, _ := os.ReadFile(task.Path)
	if err := os.WriteFile(task.Path, append(data, []byte(link)...), 0o644); err != nil {
		t.Fatal(err)
	}
	epic := mustGet(t, repo, a.Parent)
	was := snapshot(a.Path, task.Path, epic.Path, repo.NarrativePath(a.ID), filepath.Join(repo.AgentsDir(), "index.md"))
	title := "Story A renamed"
	_, err := Apply(repo, nil, a.ID, Change{Title: &title, After: &[]string{"S-0099"}}, Options{By: "alex", NoCommit: true, Now: ruleAt})
	refusedOn(t, err, "story.after")
	sameFiles(t, was)
	renamed := filepath.Join(filepath.Dir(a.Path), workitem.FileName(a.ID, title))
	if _, err := os.Stat(renamed); !os.IsNotExist(err) {
		t.Errorf("the renamed file %s is still there: %v", renamed, err)
	}
}

// A reparent that brings a finding is refused: a task moved to story B that
// still waits for a task of story A brings a task.after inside B, and the
// task and both stories are as they were.
func TestAReparentBringingAKeptRulesFindingIsRefused(t *testing.T) {
	repo, a, task, b, _ := raceRepo(t)
	other, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Another task of A", Parent: a.ID, Owner: "alex", Now: ruleAt})
	if err != nil {
		t.Fatal(err)
	}
	a = mustGet(t, repo, a.ID)
	was := snapshot(task.Path, a.Path, b.Path)
	_, err = Apply(repo, nil, task.ID, Change{Parent: &b.ID, After: &[]string{other.ID}}, Options{By: "alex", NoCommit: true, Now: ruleAt})
	refusedOn(t, err, "task.after")
	sameFiles(t, was)
}

// A wip.overlap an edit brings is advisory: the edit stands, with the
// overlap among its warnings (ADR-0019, I-0169).
func TestAnEditsOverlapIsAWarning(t *testing.T) {
	repo, a, _, b, _ := raceRepo(t)
	touches := []string{"flai/internal/x.go"}
	if _, err := Apply(repo, nil, b.ID, Change{Touches: &touches}, Options{By: "alex", NoCommit: true, Now: ruleAt}); err != nil {
		t.Fatal(err)
	}
	res, err := Apply(repo, nil, a.ID, Change{Touches: &touches}, Options{By: "alex", NoCommit: true, Now: ruleAt})
	if err != nil {
		t.Fatalf("an overlap refused the edit: %v", err)
	}
	warned := false
	for _, f := range res.Warnings {
		warned = warned || f.Rule == "wip.overlap"
	}
	if !warned {
		t.Errorf("no wip.overlap warning: %+v", res.Warnings)
	}
}

// A finding of a rule an item write cannot affect, here a design document
// without front matter, refuses nothing and warns of nothing, whether it was
// there before the edit or another writer brought it between the two runs.
func TestASkippedRulesFindingRefusesNothing(t *testing.T) {
	repo, a, _, _, _ := raceRepo(t)
	design := filepath.Join(repo.Root, "design")
	if err := os.MkdirAll(design, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(design, "before.md"), []byte("# No front matter\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	betweenItemWrite(t, func(*workitem.Repo) {
		if err := os.WriteFile(filepath.Join(design, "between.md"), []byte("# No front matter either\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	})
	res, err := Apply(repo, nil, a.ID, Change{Title: str("Story A, renamed")}, Options{By: "alex", NoCommit: true, Now: ruleAt})
	if err != nil {
		t.Fatalf("a finding of a skipped rule refused the edit: %v", err)
	}
	if whole, err := check.Run(repo, ruleAt); err != nil || !hasRule(whole.Findings, "doc.front-matter") {
		t.Fatalf("the whole check has no doc.front-matter, so the test reproduces nothing: %v", err)
	}
	for _, f := range res.Warnings {
		if strings.HasPrefix(f.Path, "design") {
			t.Errorf("a skipped rule's finding is a warning: %+v", f)
		}
	}
}

// A markdown finding already on a file of the story that a retitle relinks
// is known before the edit, so the retitle stands and the link follows.
func TestARetitleStandsOverAFindingAlreadyOnAFileItRelinks(t *testing.T) {
	repo, a, task, _, _ := raceRepo(t)
	link := "\nSee [A](../stories/" + filepath.Base(a.Path) + ").\n"
	data, _ := os.ReadFile(task.Path)
	if err := os.WriteFile(task.Path, append(data, []byte(link+"\n\n\nMore.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if len(findings(t, repo, "markdown.MD012", task.Path)) == 0 {
		t.Fatal("the task has no markdown.MD012: the test reproduces nothing")
	}
	title := "Story A renamed"
	if _, err := Apply(repo, nil, a.ID, Change{Title: &title}, Options{By: "alex", NoCommit: true, Now: ruleAt}); err != nil {
		t.Fatalf("a finding already on a relinked file refused the retitle: %v", err)
	}
	got, _ := os.ReadFile(task.Path)
	if want := workitem.FileName(a.ID, title); !strings.Contains(string(got), want) {
		t.Errorf("the task's link was not rewritten to %s:\n%s", want, got)
	}
}

func hasRule(fs []check.Finding, rule string) bool {
	for _, f := range fs {
		if f.Rule == rule {
			return true
		}
	}
	return false
}
