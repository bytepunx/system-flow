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
// naming it; the issue's file is not changed.
func TestIssueStoryMakesAStoryFromAnOpenIssue(t *testing.T) {
	f := setup(t)
	defect, err := issues.New(f.repo, issues.NewOptions{Title: "Fixture was ignored", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	was, _ := os.ReadFile(defect.Path)

	out, failed := f.call(t, "issue_story", map[string]any{"id": "I-1"})
	if failed != "" {
		t.Fatal(failed)
	}
	id, _ := out["id"].(string)
	want := map[string]any{"title": "Fixture was ignored", "nature": "remediation", "issue": "I-0001", "path": "wip/kanban/stories/" + id + "-fixture-was-ignored.md"}
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
	if now, _ := os.ReadFile(defect.Path); string(now) != string(was) {
		t.Error("issue_story should not change the issue's file")
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
// it recorded are until it is accepted; the story is made in the main wip.
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
