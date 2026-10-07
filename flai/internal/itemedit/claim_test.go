package itemedit

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/messages"
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

// ADR-0096: a gained path whose overlap with another story in progress lies
// wholly inside a shared path tells nobody, the narrower entry deciding: the
// gained folder docs/users meets docs/users/flai.md, which is shared. One
// outside every shared path is still told. The list is read as the manifest
// has it at the write, not as the repo was opened.
func TestGrownIsSilentInsideASharedPath(t *testing.T) {
	for _, c := range []struct {
		name   string
		shared string
		want   []string
	}{
		{"none shared", "", []string{"docs/users", "flai/cmd/x.go"}},
		{"docs/users/flai.md shared", "claims:\n  shared:\n    - docs/users/flai.md\n", []string{"flai/cmd/x.go"}},
		{"docs shared", "claims:\n  shared: [docs]\n", []string{"flai/cmd/x.go"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			at := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
			repo, story := claimRepo(t, at)
			mine := story("Mine", workitem.InProgress, "design")
			theirs := story("Theirs", workitem.InProgress, "docs/users/flai.md", "flai/cmd")
			mf := filepath.Join(repo.Root, "system-flow.yaml")
			data, _ := os.ReadFile(mf)
			if err := os.WriteFile(mf, append(data, c.shared...), 0o644); err != nil {
				t.Fatal(err)
			}
			w := WatchClaim(repo, mine.ID)
			if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Reach", Parent: mine.ID, Owner: "alex", Touches: []string{"docs/users", "flai/cmd/x.go"}, Now: at}); err != nil {
				t.Fatal(err)
			}
			got, err := w.Grown("claude", at)
			if err != nil {
				t.Fatal(err)
			}
			if want := []Overlapping{{Story: theirs.ID, Title: "Theirs", Paths: c.want}}; !reflect.DeepEqual(got, want) {
				t.Errorf("reached:\n got %+v\nwant %+v", got, want)
			}
		})
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

// ADR-0121: a claim that grows into another story in progress messages that
// story in the pair's conversation, from the story whose claim grew, by the
// writer, about the paths gained, and the conversation awaits the other. A
// second growth adds to the same conversation and its about, and so does one
// of the other story's, the conversation then awaiting the first.
func TestGrownMessagesTheOtherStoryInThePairsConversation(t *testing.T) {
	at := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "design")
	theirs := story("Theirs", workitem.InProgress, "cli")
	// each growth a minute later: an entry with the same heading as the last
	// joins it
	grow := func(s *workitem.Item, by string, touches ...string) {
		t.Helper()
		at = at.Add(time.Minute)
		w := WatchClaim(repo, s.ID)
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Grow", Parent: s.ID, Owner: "alex", Touches: touches, Now: at}); err != nil {
			t.Fatal(err)
		}
		got, err := w.Grown(by, at)
		if err != nil || len(got) != 1 {
			t.Fatalf("grown: %+v %v", got, err)
		}
	}
	only := func() *messages.Conversation {
		t.Helper()
		all, err := messages.List(repo)
		if err != nil || len(all) != 1 {
			t.Fatalf("one conversation: %+v %v", all, err)
		}
		return all[0]
	}

	grow(mine, "claude", "design", "flai/cmd/x.go")
	c := only()
	if c.From != mine.ID || c.To != theirs.ID || c.Awaiting() != theirs.ID || !reflect.DeepEqual(c.About, []string{"flai/cmd/x.go"}) {
		t.Errorf("from %s to %s awaiting %s about %v", c.From, c.To, c.Awaiting(), c.About)
	}
	es := c.Entries()
	if len(es) != 1 || es[0].Author != "claude" || es[0].Story != mine.ID {
		t.Fatalf("one entry by claude for %s: %+v", mine.ID, es)
	}
	for _, want := range []string{"`flai/cmd/x.go`", "first", "`flai message escalate`", "`message_escalate`"} {
		if !strings.Contains(es[0].Text, want) {
			t.Errorf("the message does not say %s:\n%s", want, es[0].Text)
		}
	}
	if len(Overlaps(repo)) != 2 {
		t.Errorf("the overlapped notices stay: %+v", Overlaps(repo))
	}

	grow(mine, "claude", "flai/internal/y.go")
	c = only()
	if es := c.Entries(); len(es) != 2 || !strings.Contains(es[1].Text, "`flai/internal/y.go`") || strings.Contains(es[1].Text, "flai/cmd/x.go") {
		t.Errorf("a second entry about the paths gained only: %+v", es)
	}
	if !reflect.DeepEqual(c.About, []string{"flai/cmd/x.go", "flai/internal/y.go"}) || c.Awaiting() != theirs.ID {
		t.Errorf("about %v awaiting %s", c.About, c.Awaiting())
	}

	grow(theirs, "codex", "design/a.md")
	c = only()
	if es := c.Entries(); len(es) != 3 || es[2].Author != "codex" || es[2].Story != theirs.ID || c.Awaiting() != mine.ID {
		t.Errorf("the other story's growth adds to the same conversation, awaiting %s: %+v awaiting %s", mine.ID, es, c.Awaiting())
	}
}

// S-0295: a growth whose overlap lies wholly inside a shared path of the
// manifest messages nobody.
func TestGrownMessagesNobodyInsideASharedPath(t *testing.T) {
	at := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "design")
	story("Theirs", workitem.InProgress, "docs/users/flai.md")
	mf := filepath.Join(repo.Root, "system-flow.yaml")
	data, _ := os.ReadFile(mf)
	if err := os.WriteFile(mf, append(data, "claims:\n  shared: [docs]\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	w := WatchClaim(repo, mine.ID)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Shared", Parent: mine.ID, Owner: "alex", Touches: []string{"design", "docs/users"}, Now: at}); err != nil {
		t.Fatal(err)
	}
	if got, err := w.Grown("claude", at); got != nil || err != nil {
		t.Errorf("reached %+v %v", got, err)
	}
	if all, err := messages.List(repo); len(all) != 0 || err != nil {
		t.Errorf("messaged: %+v %v", all, err)
	}
}

// ADR-0121: a message that cannot be written is Grown's error, and the
// stories reached, the overlapped notices, and the write stand.
func TestGrownKeepsTheWriteWhenTheMessageFails(t *testing.T) {
	at := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "design")
	theirs := story("Theirs", workitem.InProgress, "cli")
	// a file where the messages folder goes: no conversation can be saved
	if err := os.WriteFile(messages.Dir(repo), []byte("not a folder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	w := WatchClaim(repo, mine.ID)
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Reach", Parent: mine.ID, Owner: "alex", Touches: []string{"design", "flai/cmd/x.go"}, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	got, err := w.Grown("claude", at)
	if err == nil || !strings.Contains(err.Error(), theirs.ID) {
		t.Errorf("the failed message is the error, naming %s: %v", theirs.ID, err)
	}
	if want := []Overlapping{{Story: theirs.ID, Title: "Theirs", Paths: []string{"flai/cmd/x.go"}}}; !reflect.DeepEqual(got, want) {
		t.Errorf("reached:\n got %+v\nwant %+v", got, want)
	}
	if len(Overlaps(repo)) != 2 {
		t.Errorf("the overlapped notices stand: %+v", Overlaps(repo))
	}
	if it, err := repo.Get(task.ID); err != nil || !reflect.DeepEqual(it.Touches, []string{"design", "flai/cmd/x.go"}) {
		t.Errorf("the write stands: %+v %v", it, err)
	}
}
