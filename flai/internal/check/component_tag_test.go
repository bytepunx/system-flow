package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// I-0024: a story whose touches reach two components while no tag of its own
// or its epic's names one of them cannot be released; check says which tag
// to add before acceptance finds out.
func TestComponentTagRule(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\n  - name: web\n    path: web\n    kind: sveltekit\n    tags: [dashboard]\n"), 0o644)
	for _, d := range []string{"design/adrs", "design/system", "design/tech", "design/conventions", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "cli", "web"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	untagged := mustItem(t, repo, workitem.Epic, "Untagged", "")
	tagged := mustItem(t, repo, workitem.Epic, "Tagged", "")
	tagged.Tags = []string{"dashboard"}
	_ = repo.Save(tagged)
	story := func(title, parent string, tags, touches []string, nature string) *workitem.Item {
		t.Helper()
		s := mustItem(t, repo, workitem.Story, title, parent)
		s.Tags, s.Touches, s.Nature = tags, touches, nature
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	warned := story("Both sides", untagged.ID, []string{"ux"}, []string{"cli/cmd", "web/src"}, "feature")
	own := story("Tagged itself", untagged.ID, []string{"cli"}, []string{"cli/cmd", "web/src"}, "feature")
	byEpic := story("Tagged by its epic", tagged.ID, nil, []string{"cli/cmd", "web/src"}, "remediation")
	one := story("One side", untagged.ID, nil, []string{"cli/cmd", "docs/users"}, "feature")
	research := story("A finding", untagged.ID, nil, []string{"cli", "web"}, "research")
	byTask := story("Through a task", untagged.ID, nil, []string{"cli/cmd"}, "improvement")
	task := mustItem(t, repo, workitem.Task, "The other side", byTask.ID)
	task.Touches = []string{"dashboard"} // a component's tag reads as its path
	_ = repo.Save(task)

	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, f := range res.Findings {
		if f.Rule == "story.component-tag" {
			got[filepath.Base(f.Path)[:6]] = f.Message
			if f.Level != Warning {
				t.Errorf("a warning, not %s: %+v", f.Level, f)
			}
		}
	}
	msg := got[warned.ID]
	if !strings.Contains(msg, "touches cli and web") || !strings.Contains(msg, "flai edit "+warned.ID+" --tag <cli|web>,ux") {
		t.Errorf("the untagged story is warned with the tag to add, keeping its own: %q", msg)
	}
	if _, ok := got[byTask.ID]; !ok {
		t.Errorf("an open task's touches count toward the story's: %v", got)
	}
	for _, s := range []*workitem.Item{own, byEpic, one, research} {
		if m, ok := got[s.ID]; ok {
			t.Errorf("%s %s is not warned: %q", s.ID, s.Title, m)
		}
	}

	warned.Status = workitem.Done
	_ = repo.Save(warned)
	res, _ = Run(repo, now)
	for _, f := range res.Findings {
		if f.Rule == "story.component-tag" && strings.Contains(f.Path, warned.ID) {
			t.Errorf("a closed story is not warned: %+v", f)
		}
	}
}
