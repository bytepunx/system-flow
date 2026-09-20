package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0070: an epic with a story in review (one task) and a story in backlog.
func cancelledProject(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic, _ := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Now: now})
	for _, title := range []string{"Reviewed", "Waiting"} {
		s, _ := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Now: now})
		data, _ := os.ReadFile(s.Path)
		_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Task of " + title, Parent: s.ID, Now: now}); err != nil {
			t.Fatal(err)
		}
	}
	s, _ := repo.Get("S-0001")
	for _, to := range []string{workitem.Ready, workitem.InProgress, workitem.Review} {
		if _, err := repo.Transition(s, to, "alex", "", now); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

func rules(t *testing.T, repo *workitem.Repo, rule string) []string {
	t.Helper()
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, f := range res.Findings {
		if f.Rule == rule {
			out = append(out, f.Message)
		}
	}
	return out
}

func TestAnOpenItemUnderACancelledParentIsReported(t *testing.T) {
	repo := cancelledProject(t)
	if got := rules(t, repo, "item.parent-cancelled"); len(got) != 0 {
		t.Fatalf("clean project: %v", got)
	}
	// As an older flai did it: the epic alone, and one story alone.
	items, _ := repo.List(false)
	for _, id := range []string{"E-0001", "S-0002"} {
		it, _ := repo.Get(id)
		if _, err := repo.Move(it, workitem.Cancelled, workitem.MoveOptions{By: "alex", Reason: "old", Now: now, Items: items}); err != nil {
			t.Fatal(err)
		}
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	got := strings.Join(rules(t, repo, "item.parent-cancelled"), "\n")
	for _, want := range []string{"S-0001 is review under E-0001, which is cancelled", "T-0002 is backlog under S-0002, which is cancelled", "flai move T-0002 cancelled"} {
		if !strings.Contains(got, want) {
			t.Errorf("findings do not say %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "S-0002 is") {
		t.Errorf("a cancelled story under a cancelled epic is fine:\n%s", got)
	}
}

func TestACascadeOutOfReviewIsAValidHistory(t *testing.T) {
	repo := cancelledProject(t)
	epic, _ := repo.Get("E-0001")
	if _, err := repo.TransitionAll(epic, workitem.Cancelled, "alex", "a different route", now); err != nil {
		t.Fatal(err)
	}
	if got := rules(t, repo, "item.sequence"); len(got) != 0 {
		t.Errorf("review to cancelled through the epic: %v", got)
	}
	if got := rules(t, repo, "item.parent-cancelled"); len(got) != 0 {
		t.Errorf("after the cascade: %v", got)
	}
	// The same history without the parent's cancellation is still wrong.
	s, _ := repo.Get("S-0001")
	data, _ := os.ReadFile(s.Path)
	_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "parent: E-0001", "parent: E-0002", 1)), 0o644)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Other", Now: now}); err != nil {
		t.Fatal(err)
	}
	if got := rules(t, repo, "item.sequence"); len(got) != 1 || !strings.Contains(got[0], "cancelled cannot follow review") {
		t.Errorf("review to cancelled with no cancelled parent: %v", got)
	}
}
