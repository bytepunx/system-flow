package itemedit

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// claimRepo is a project with one sub-project, flai, tagged cli, and story,
// which creates a story with touches and moves it to status, ready or in
// progress.
func claimRepo(t *testing.T, at time.Time) (*workitem.Repo, func(title, status string, touches ...string) *workitem.Item) {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: flai\n    path: flai\n    tags: [cli]\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story := func(title, status string, touches ...string) *workitem.Item {
		t.Helper()
		s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Owner: "alex", Touches: touches, Now: at})
		if err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(s.Path)
		_ = os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
		s, _ = repo.Get(s.ID)
		for _, st := range []string{workitem.Ready, workitem.InProgress} {
			if st == workitem.InProgress && status == workitem.Ready {
				break
			}
			if _, err := repo.Transition(s, st, "alex", "", at); err != nil {
				t.Fatal(err)
			}
		}
		return s
	}
	return repo, story
}

// Which stories a grown claim reaches (I-0059): only stories in progress,
// only through the paths the claim gained, never one whose claim is empty,
// and nothing from a story that has not started.
func TestGrownReachesOnlyStoriesInProgressThroughWhatTheClaimGained(t *testing.T) {
	at := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "design")
	theirs := story("Theirs", workitem.InProgress, "cli")
	story("Empty", workitem.InProgress)
	story("Waiting", workitem.Ready, "docs")
	later := story("Later", workitem.Ready, "docs/later")

	w := WatchClaim(repo, mine.ID)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Reach", Parent: mine.ID, Owner: "alex", Touches: []string{"design", "flai/cmd/x.go", "docs/a.md"}, Now: at}); err != nil {
		t.Fatal(err)
	}
	got, err := w.Grown("claude", at)
	if err != nil {
		t.Fatal(err)
	}
	if want := []Overlapping{{Story: theirs.ID, Title: "Theirs", Paths: []string{"flai/cmd/x.go"}}}; !reflect.DeepEqual(got, want) {
		t.Errorf("reached:\n got %+v\nwant %+v", got, want)
	}
	notices := Overlaps(repo)
	if len(notices) != 2 || notices[0].ID != mine.ID || notices[1].ID != theirs.ID || notices[0].Cause() != theirs.ID || notices[1].Cause() != mine.ID || notices[0].By != "claude" {
		t.Errorf("one notice for each story: %+v", notices)
	}

	// a story that has not started reaches nobody: the hold judges it
	w = WatchClaim(repo, later.ID)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Early", Parent: later.ID, Owner: "alex", Touches: []string{"flai"}, Now: at}); err != nil {
		t.Fatal(err)
	}
	if got, err := w.Grown("claude", at); err != nil || got != nil {
		t.Errorf("not started: %+v %v", got, err)
	}
	if _, err := WatchClaim(repo, "T-0099").Grown("claude", at); err == nil {
		t.Error("a watch on no item says so once the write is done")
	}
	if len(Overlaps(repo)) != 2 {
		t.Errorf("nothing more told: %+v", Overlaps(repo))
	}
}

// ADR-0096: a task that names a file inside its story's folder touch narrows
// the claim to it. That gains the claim nothing, so it reaches no story, even
// one in progress that touches the same file; a task that names a path outside
// the folder does reach one.
func TestGrownIgnoresANarrowedClaim(t *testing.T) {
	at := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "docs")
	story("Same file", workitem.InProgress, "docs/a.md")
	elsewhere := story("Elsewhere", workitem.InProgress, "flai/cmd/y.go")
	task := func(title string, touches ...string) []Overlapping {
		t.Helper()
		w := WatchClaim(repo, mine.ID)
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: title, Parent: mine.ID, Owner: "alex", Touches: touches, Now: at}); err != nil {
			t.Fatal(err)
		}
		got, err := w.Grown("claude", at)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}

	if got := task("Narrow", "docs/a.md"); got != nil {
		t.Errorf("narrowing docs to docs/a.md reached %+v", got)
	}
	if len(Overlaps(repo)) != 0 {
		t.Errorf("narrowing told: %+v", Overlaps(repo))
	}
	got := task("Reach", "docs/b.md", "flai/cmd/y.go")
	if want := []Overlapping{{Story: elsewhere.ID, Title: "Elsewhere", Paths: []string{"flai/cmd/y.go"}}}; !reflect.DeepEqual(got, want) {
		t.Errorf("reached:\n got %+v\nwant %+v", got, want)
	}
}
