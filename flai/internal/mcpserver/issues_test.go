package mcpserver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0198: issue_story makes a backlog story from an open issue as flai issue
// story does, and refuses a closed issue or one an open story already links,
// naming it. S-0203: the story is a draft carrying the issue's cost of delay
// inputs, and the issue's Remediation section names it.
func TestIssueStoryMakesAStoryFromAnOpenIssue(t *testing.T) {
	f := setup(t)
	defect, err := issues.New(f.repo, issues.NewOptions{Title: "Fixture was ignored", Class: "defect", Cost: "20m", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	defect.Body += "\n## Impact\n- revenue_per_week: 1200\n"
	if err := defect.Save(); err != nil {
		t.Fatal(err)
	}

	out, failed := f.call(t, "issue_story", map[string]any{"id": "I-1"})
	if failed != "" {
		t.Fatal(failed)
	}
	id, _ := out["id"].(string)
	want := map[string]any{"title": "Fixture was ignored", "nature": "remediation", "issue": "I-0001", "path": "wip/kanban/stories/" + id + "-fixture-was-ignored.md", "draft": true}
	for k, v := range want {
		if out[k] != v {
			t.Errorf("issue_story %s = %v, want %v", k, out[k], v)
		}
	}
	story, err := f.repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if story.Type != workitem.Story || story.Status != workitem.Backlog || story.Nature != "remediation" || story.Owner != "claude" || !issues.Links(story.Body, defect) {
		t.Errorf("the story should be a backlog remediation story the agent owns that links I-0001: %+v\n%s", story, story.Body)
	}
	if !story.Draft {
		t.Error("the story should be a draft")
	}
	cod := story.CostOfDelay
	if cod == nil || cod.Inputs == nil || cod.Inputs.TimeLostPerCycle != "20m" || cod.Inputs.RevenuePerWeek == nil || *cod.Inputs.RevenuePerWeek != 1200 || cod.Inputs.By != "flai" || cod.By != "" {
		t.Errorf("the story should carry the issue's time lost per cycle and its Impact's revenue, set by flai: %+v", cod)
	}
	for _, want := range []string{"time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h", "revenue_per_week 1200 carried over from I-0001's Impact section."} {
		if !strings.Contains(story.Body, want) {
			t.Errorf("the story's Notes should say %q:\n%s", want, story.Body)
		}
	}
	if is, err := issues.Get(f.repo, "I-0001"); err != nil || !strings.Contains(is.Body, "## Remediation\n\nStory "+id+" remediates this issue, created from it at ") {
		t.Fatalf("the issue's Remediation section should name %s: %v %+v", id, err, is)
	}

	if _, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001"}); !strings.Contains(failed, "I-0001 is already linked by open story "+id) {
		t.Errorf("an issue an open story links should be refused, naming it: %q", failed)
	}

	closed, err := issues.New(f.repo, issues.NewOptions{Title: "Go not on PATH", Class: "blocker", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if err := issues.Close(closed, "fixed", t0); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "issue_story", map[string]any{"id": closed.ID}); !strings.Contains(failed, "I-0002 is closed, so no story was made for it") {
		t.Errorf("a closed issue should be refused: %q", failed)
	}
	if list, _ := f.repo.List(false); countStories(list) != 2 {
		t.Errorf("a refusal should make no story: %d stories", countStories(list))
	}
}

// S-0198: story reads the issue from that story's worktree, where the issues
// it recorded are until it is accepted, and names the story there (S-0203);
// the story is made in the main wip.
func TestIssueStoryReadsAStorysWorktree(t *testing.T) {
	f := setup(t)
	wt := *f.repo
	wt.Root = f.repo.WorktreePath(f.story.ID)
	if err := os.MkdirAll(wt.Root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := issues.New(&wt, issues.NewOptions{Title: "Only on the branch", Class: "efficiency", Story: f.story.ID, Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001"}); !strings.Contains(failed, "I-0001 not found") {
		t.Errorf("without story the issue is read from the project, which has none: %q", failed)
	}
	out, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001", "story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["nature"] != "improvement" || !strings.HasPrefix(out["path"].(string), "wip/kanban/stories/") {
		t.Errorf("issue_story with story: %v", out)
	}
	if _, err := os.Stat(filepath.Join(wt.Root, "wip")); err == nil {
		t.Error("the story should be made in the main checkout's wip, not the worktree's")
	}
	if is, err := issues.Get(&wt, "I-0001"); err != nil || !strings.Contains(is.Body, "Story "+out["id"].(string)+" remediates this issue") {
		t.Errorf("the issue in the worktree should name the story: %v %+v", err, is)
	}
}

func countStories(items []*workitem.Item) int {
	n := 0
	for _, it := range items {
		if it.Type == workitem.Story {
			n++
		}
	}
	return n
}
