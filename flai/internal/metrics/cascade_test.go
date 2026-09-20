package metrics

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0070: items cancelled with their parent count as cancelled at the time
// of the cascade, by their own transition, with no change to the definitions.
func TestCascadedItemsCountAsCancelledAtTheCascade(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	t0 := time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)
	epic, _ := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Now: t0})
	for _, title := range []string{"Started", "Waiting"} {
		s, _ := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Now: t0})
		data, _ := os.ReadFile(s.Path)
		_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
	}
	s, _ := repo.Get("S-0001")
	for i, to := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.Transition(s, to, "agent", "", t0.Add(time.Duration(i+1)*time.Hour)); err != nil {
			t.Fatal(err)
		}
	}
	cascade := t0.Add(48 * time.Hour)
	epic, _ = repo.Get(epic.ID)
	if _, err := repo.TransitionAll(epic, workitem.Cancelled, "alex", "a different route", cascade); err != nil {
		t.Fatal(err)
	}
	items, _ := repo.List(false)
	rep := Compute(items, Options{Now: cascade.Add(time.Hour)})
	if rep.Summary.Cancelled != 2 || rep.Summary.Completed != 0 || rep.Summary.CancellationRate != 1 {
		t.Errorf("summary: %+v", rep.Summary)
	}
	for _, m := range rep.Items {
		if m.Status != workitem.Cancelled || m.Completed != cascade.Format(workitem.TimeFormat) {
			t.Errorf("%s: status %s completed %s", m.ID, m.Status, m.Completed)
		}
		if m.Age != nil {
			t.Errorf("%s still ages after the cascade", m.ID)
		}
	}
}
