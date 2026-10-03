package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0200: accepting an epic's last open story accepts the epic with it, in
// the one acceptance commit, and archives what of it is left on the board: a
// cancelled story does not hold the epic open.
func TestAcceptingTheLastOpenStoryAcceptsItsEpic(t *testing.T) {
	root, _ := researchProject(t, "remediation", true)
	for _, args := range [][]string{
		{"story", "new", "Dropped", "--epic", "E-0001"},
		{"move", "S-0002", "cancelled", "--reason", "not needed", "--yes"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.Review {
		t.Fatalf("E-0001 followed S-0001 to review: %s", got)
	}

	dry, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run")
	if code != 0 || !strings.Contains(dry, "would also move E-0001 Epic from review to done, following S-0001, and archive it") {
		t.Fatalf("the dry run names the epic's move: %d %s %s", code, dry, errOut)
	}
	js, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	var pre preview.Acceptance
	if err := json.Unmarshal([]byte(js), &pre); code != 0 || err != nil {
		t.Fatalf("dry run --json: %d %v %s %s", code, err, js, errOut)
	}
	if want := (workitem.Followed{ID: "E-0001", Type: workitem.Epic, Title: "Epic", From: workitem.Review, To: workitem.Done, Story: "S-0001"}); pre.Epic == nil || *pre.Epic != want || len(pre.Blockers) != 0 {
		t.Errorf("dry run --json epic: %+v, blockers %v", pre.Epic, pre.Blockers)
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.Review {
		t.Errorf("the dry run changes nothing: E-0001 is %s", got)
	}

	before := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD"))
	out, errOut, code := runIn(t, root, "accept", "S-0001")
	if code != 0 || !strings.Contains(out, "E-0001 followed it from review to done, archived") {
		t.Fatalf("accept: %d %s %s", code, out, errOut)
	}
	for _, glob := range []string{"epics/E-0001-*.md", "stories/S-0001-*.md", "stories/S-0002-*.md"} {
		if m, _ := filepath.Glob(filepath.Join(root, "wip/archive/kanban", glob)); len(m) != 1 {
			t.Errorf("%s is archived: %v", glob, m)
		}
		if m, _ := filepath.Glob(filepath.Join(root, "wip/kanban", glob)); len(m) != 0 {
			t.Errorf("%s is off the board: %v", glob, m)
		}
	}
	m, _ := filepath.Glob(filepath.Join(root, "wip/archive/kanban/epics/E-0001-*.md"))
	if len(m) == 1 {
		data, _ := os.ReadFile(m[0])
		if !strings.Contains(string(data), "status: done") || !strings.Contains(string(data), "follows S-0001, which moved to done") {
			t.Errorf("the epic is done, with a note naming its story:\n%s", data)
		}
	}
	if subject := strings.TrimSpace(gitIn(t, root, "log", "-1", "--format=%s")); subject != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the acceptance commit names the epic: %q", subject)
	}
	if n := strings.Fields(gitIn(t, root, "rev-list", "--count", before+"..HEAD")); len(n) != 1 || n[0] != "2" {
		t.Errorf("the merged story commit and one acceptance commit: %v", n)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); !strings.Contains(files, "wip/archive/kanban/epics/E-0001-epic.md") {
		t.Errorf("the epic's archive is in the acceptance commit:\n%s", files)
	}
	if status := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); status != "" {
		t.Errorf("everything is committed: %q", status)
	}
}

// S-0200: a story that is not its epic's last open one is accepted alone.
func TestAcceptingAStoryWithOthersOpenLeavesItsEpic(t *testing.T) {
	root, _ := pendingProject(t) // S-0001 and S-0002 in review, S-0003 in backlog
	js, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	var pre preview.Acceptance
	if err := json.Unmarshal([]byte(js), &pre); code != 0 || err != nil || pre.Epic != nil {
		t.Fatalf("dry run --json has no epic: %d %v %s %s", code, err, js, errOut)
	}
	for _, id := range []string{"S-0001", "S-0002"} {
		out, errOut, code := runIn(t, root, "accept", id)
		if code != 0 || strings.Contains(out, "E-0001") {
			t.Fatalf("accept %s: %d %s %s", id, code, out, errOut)
		}
		if subject := strings.TrimSpace(gitIn(t, root, "log", "-1", "--format=%s")); subject != "chore: ["+id+"] accept and archive" {
			t.Errorf("the acceptance commit of %s names no epic: %q", id, subject)
		}
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.Review {
		t.Errorf("E-0001 stays where it was, on the board: %s", got)
	}
}

// moveProject is an epic in backlog with two stories in backlog, each with
// a ticked criterion and a task, in a project without git.
func moveProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	for _, args := range [][]string{
		{"epic", "new", "Epic"},
		{"story", "new", "First", "--epic", "E-0001"},
		{"story", "new", "Second", "--epic", "E-0001"},
		{"task", "new", "One", "--story", "S-0001"},
		{"task", "new", "Two", "--story", "S-0002"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	for _, name := range []string{"S-0001-first.md", "S-0002-second.md"} {
		file := filepath.Join(root, "wip/kanban/stories", name)
		s, _ := os.ReadFile(file)
		_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] ok\n", 1)), 0o644)
	}
	return root
}

// S-0200: flai move names the move a story's epic made with it, in text and
// in --json, and says nothing of an epic that stayed.
func TestMoveReportsTheEpicFollowing(t *testing.T) {
	root := moveProject(t)
	out, errOut, code := runIn(t, root, "move", "S-0001", "ready")
	if code != 0 || out != "S-0001 → ready\n  E-0001 → ready, following S-0001\n" {
		t.Fatalf("move ready: %d %q %s", code, out, errOut)
	}
	js, errOut, code := runIn(t, root, "move", "S-0001", "in-progress", "--json")
	var moved struct {
		Followed *workitem.Followed `json:"followed"`
	}
	if err := json.Unmarshal([]byte(js), &moved); code != 0 || err != nil {
		t.Fatalf("move in-progress --json: %d %v %s %s", code, err, js, errOut)
	}
	if want := (workitem.Followed{ID: "E-0001", Type: workitem.Epic, Title: "Epic", From: workitem.Ready, To: workitem.InProgress, Story: "S-0001"}); moved.Followed == nil || *moved.Followed != want {
		t.Errorf("--json followed: %+v", moved.Followed)
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.InProgress {
		t.Errorf("E-0001 is saved in progress: %s", got)
	}
	js, _, code = runIn(t, root, "move", "S-0002", "ready", "--json")
	if code != 0 || strings.Contains(js, `"followed"`) {
		t.Errorf("an epic that stays is not in --json: %d %s", code, js)
	}
	out, _, code = runIn(t, root, "move", "S-0002", "in-progress")
	if code != 0 || out != "S-0002 → in-progress\n" {
		t.Errorf("an epic that stays is not in the text: %d %q", code, out)
	}
}

// S-0200: a cancelled story whose siblings are all in review takes its epic
// to review, and the cancellation says so, previewed and done.
func TestCancellingAStoryReportsTheEpicFollowing(t *testing.T) {
	root := moveProject(t)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"},
		{"move", "S-0002", "ready"}, {"move", "S-0002", "in-progress"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"}, {"move", "T-0002", "done"},
		{"move", "S-0002", "review"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.InProgress {
		t.Fatalf("S-0001 holds E-0001 in progress: %s", got)
	}
	dry, errOut, code := runIn(t, root, "move", "S-0001", "cancelled", "--reason", "dropped", "--dry-run")
	if code != 0 || !strings.Contains(dry, "would also move E-0001 Epic from in-progress to review, following S-0001") {
		t.Fatalf("cancel dry run: %d %s %s", code, dry, errOut)
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.InProgress {
		t.Errorf("the dry run changes nothing: %s", got)
	}
	js, errOut, code := runIn(t, root, "move", "S-0001", "cancelled", "--reason", "dropped", "--yes", "--json")
	var res preview.Cancellation
	if err := json.Unmarshal([]byte(js), &res); code != 0 || err != nil {
		t.Fatalf("cancel --json: %d %v %s %s", code, err, js, errOut)
	}
	if want := (workitem.Followed{ID: "E-0001", Type: workitem.Epic, Title: "Epic", From: workitem.InProgress, To: workitem.Review, Story: "S-0001"}); res.Followed == nil || *res.Followed != want {
		t.Errorf("cancel --json followed: %+v", res.Followed)
	}
	if got := statusOf(t, root, "epics/E-0001-*.md"); got != workitem.Review {
		t.Errorf("E-0001 follows to review, no further: %s", got)
	}
	// back from cancelled, S-0001 holds the epic in progress again; cancelled
	// once more, it lets the epic go, and the text says so each time
	out, errOut, code := runIn(t, root, "move", "S-0001", "backlog")
	if code != 0 || out != "S-0001 → backlog\n  E-0001 → in-progress, following S-0001\n" {
		t.Errorf("back from cancelled: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "move", "S-0001", "cancelled", "--reason", "dropped again", "--yes")
	if code != 0 || !strings.Contains(out, "S-0001 → cancelled\n  E-0001 → review, following S-0001\n") {
		t.Errorf("cancel: %d %q %s", code, out, errOut)
	}
}
