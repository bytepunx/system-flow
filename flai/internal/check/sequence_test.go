package check

import (
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0173: a history flai move wrote, back moves included (ADR-0055), is a valid history.
func TestBackMovesAreAValidHistory(t *testing.T) {
	repo := cancelledProject(t)
	moves := []struct{ id, to, reason string }{
		{"S-0001", workitem.InProgress, "sent back"},
		{"S-0001", workitem.Ready, ""},
		{"S-0001", workitem.Backlog, ""},
		{"S-0001", workitem.Ready, ""},
		{"S-0002", workitem.Ready, ""},
		{"S-0002", workitem.Backlog, ""},
		{"S-0002", workitem.Cancelled, "not now"},
		{"S-0002", workitem.Backlog, ""},
	}
	for _, m := range moves {
		it, _ := repo.Get(m.id)
		if _, err := repo.Transition(it, m.to, "alex", m.reason, now); err != nil {
			t.Fatalf("%s to %s: %v", m.id, m.to, err)
		}
	}
	if got := rules(t, repo, "item.sequence"); len(got) != 0 {
		t.Errorf("back moves: %v", got)
	}
}

func TestNothingFollowsDone(t *testing.T) {
	repo := cancelledProject(t)
	task, _ := repo.Get("T-0001")
	for _, to := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
		if _, err := repo.Transition(task, to, "alex", "", now); err != nil {
			t.Fatal(err)
		}
	}
	data, _ := os.ReadFile(task.Path)
	text := string(data)
	at := strings.LastIndex(text, "    by: alex\n") + len("    by: alex\n")
	text = text[:at] + "  - to: in-progress\n    at: " + now.UTC().Format(workitem.TimeFormat) + "\n    by: alex\n" + text[at:]
	_ = os.WriteFile(task.Path, []byte(strings.Replace(text, "status: done", "status: in-progress", 1)), 0o644)
	if got := rules(t, repo, "item.sequence"); len(got) != 1 || !strings.Contains(got[0], "in-progress cannot follow done") {
		t.Errorf("done then in-progress: %v", got)
	}
}
