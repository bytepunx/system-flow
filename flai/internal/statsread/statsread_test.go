package statsread

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// project writes a project with a sub-project, a board whose in-progress
// limit is 3, a story, and a thread on it, and opens it.
func project(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	board := "---\ntitle: Board\nstatus: active\nwip_limits:\n  in-progress: 3\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "One", Owner: "t", Now: now}); err != nil {
		t.Fatal(err)
	}
	if _, err := threads.New(repo, threads.NewOptions{Title: "Which way", On: "S-0001", Author: "agent", Text: "Left or right?", Now: now}); err != nil {
		t.Fatal(err)
	}
	return repo
}

// warnings are the events log took, decoded.
func warnings(t *testing.T, b *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, l := range strings.Split(strings.TrimSpace(b.String()), "\n") {
		if l == "" {
			continue
		}
		var ev map[string]any
		if err := json.Unmarshal([]byte(l), &ev); err != nil {
			t.Fatal(err)
		}
		out = append(out, ev)
	}
	return out
}

func TestReadWithoutGitLeavesTheCommitsOutAndWarns(t *testing.T) {
	repo := project(t)
	var logged bytes.Buffer
	items, opt, err := Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || len(opt.Threads) != 1 || opt.WIPLimit != 3 || len(opt.Projects) != 1 || opt.Projects[0].Name != "cli" {
		t.Errorf("read %d items, options %+v", len(items), opt)
	}
	if opt.Commits != nil {
		t.Errorf("commits outside a repository: %v", opt.Commits)
	}
	evs := warnings(t, &logged)
	if len(evs) != 1 || evs[0]["level"] != "WARN" || evs[0]["component"] != "stats" || evs[0]["err"] == nil || evs[0]["detail"] == nil {
		t.Errorf("want one warning with err and detail, got %v", evs)
	}
}

func TestReadInARepositoryReadsTheCommits(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	repo := project(t)
	r := execx.System{}
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"},
		{"add", "-A"}, {"commit", "-q", "-m", "feat: [S-0001] the story"},
	} {
		if out, err := r.Run(repo.Root, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	var logged bytes.Buffer
	_, opt, err := Read(r, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if got := opt.Commits["S-0001"]; len(got) != 1 || got[0] != "system-flow.yaml" {
		t.Errorf("commits = %v, want the manifest alone, wip left out", opt.Commits)
	}
	if logged.Len() != 0 {
		t.Errorf("warned in a repository: %s", logged.String())
	}
}

// reopen rewrites the project's manifest with a planning section and its
// board with a pull order, and opens it again.
func reopen(t *testing.T, repo *workitem.Repo, planning string) *workitem.Repo {
	t.Helper()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + planning
	if err := os.WriteFile(filepath.Join(repo.Root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	board := "---\ntitle: Board\nstatus: active\nwip_limits:\n  in-progress: 2\norder:\n  - S-0001\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(repo.Root, "wip/kanban/board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	again, err := workitem.Open(repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	return again
}

// S-0213: the board's pull order and planning.default_duration, with the
// limit, are what the ready column's cost of delay is projected from.
func TestReadPassesTheBoardsOrderAndTheFallbackDuration(t *testing.T) {
	repo := reopen(t, project(t), "planning:\n  default_duration: 90m\n")
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if opt.WIPLimit != 2 || len(opt.Order) != 1 || opt.Order[0] != "S-0001" || opt.Fallback != 90*time.Minute {
		t.Errorf("limit %d, order %v, fallback %v; want 2, [S-0001], 1h30m", opt.WIPLimit, opt.Order, opt.Fallback)
	}
}

// S-0213: a planning.default_duration that is not a duration stops the read,
// as it stops flai forecast, saying what to write.
func TestReadStopsOnABadFallbackDuration(t *testing.T) {
	repo := reopen(t, project(t), "planning:\n  default_duration: soon\n")
	_, _, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err == nil || !strings.Contains(err.Error(), "cannot project the ready column's cost of delay") || !strings.Contains(err.Error(), "planning.default_duration") {
		t.Errorf("a bad default duration: %v", err)
	}
}

func TestReadStopsOnAnUnreadableThread(t *testing.T) {
	repo := project(t)
	if err := os.WriteFile(filepath.Join(threads.Dir(repo), "TH-0009-bad.md"), []byte("---\nid: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler)); err == nil || !strings.Contains(err.Error(), "cannot read the threads") {
		t.Errorf("an unreadable thread: %v", err)
	}
}

// writeFile writes data to the path under the project's root, making its
// folder.
func writeFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

// conversationDoc is a conversation between S-0001 and S-0002 with one entry.
const conversationDoc = "---\nid: MS-0001\ntitle: Hello\nfrom: S-0001\nto: S-0002\nstatus: open\nparticipants: [agent]\n" +
	"created: 2026-10-01T09:00:00Z\nupdated: 2026-10-01T09:00:00Z\n---\n\n### 2026-10-01T09:00:00Z agent S-0001\n\nHello.\n"

// S-0337: Read loads the conversations under wip/messages and the conflict
// log, leaving out with a warning a conversation it cannot read.
func TestReadLoadsTheConversationsAndTheConflictLog(t *testing.T) {
	repo := project(t)
	writeFile(t, filepath.Join(messages.Dir(repo), "MS-0001-hello.md"), conversationDoc)
	writeFile(t, filepath.Join(messages.Dir(repo), "MS-0002-bad.md"), "---\nid: [\n---\n")
	if err := storygit.AppendConflict(repo.CacheDir(), storygit.Conflict{At: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC), Story: "S-0001", Where: storygit.WhereSync, Paths: []string{"a"}}); err != nil {
		t.Fatal(err)
	}
	var logged bytes.Buffer
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if len(opt.Conversations) != 1 || opt.Conversations[0].ID != "MS-0001" || len(opt.Conversations[0].Entries()) != 1 {
		t.Errorf("conversations = %+v, want MS-0001 alone", opt.Conversations)
	}
	if opt.ConflictsUnreadable || len(opt.Conflicts) != 1 || opt.Conflicts[0].Where != storygit.WhereSync {
		t.Errorf("conflicts = %+v unreadable %v, want the sync line", opt.Conflicts, opt.ConflictsUnreadable)
	}
	var skipped map[string]any
	for _, ev := range warnings(t, &logged) {
		if ev["msg"] == "cannot read some conversations" {
			skipped = ev
		}
	}
	if skipped == nil || skipped["level"] != "WARN" || skipped["count"] != float64(1) || skipped["err"] == nil || skipped["detail"] == nil {
		t.Errorf("want a warning naming the conversation left out, got %s", logged.String())
	}
}

// S-0337: a missing conflict log means no conflicts and is no warning; one
// that cannot be read sets ConflictsUnreadable and is warned of.
func TestReadWarnsOfAConflictLogItCannotRead(t *testing.T) {
	repo := project(t)
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err != nil || opt.ConflictsUnreadable || opt.Conflicts != nil {
		t.Errorf("a missing log: %v, unreadable %v, %v", err, opt.ConflictsUnreadable, opt.Conflicts)
	}
	writeFile(t, storygit.ConflictsPath(repo.CacheDir()), "{not json\n")
	var logged bytes.Buffer
	_, opt, err = Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if !opt.ConflictsUnreadable || opt.Conflicts != nil {
		t.Errorf("a bad log: unreadable %v, %v", opt.ConflictsUnreadable, opt.Conflicts)
	}
	warned := false
	for _, ev := range warnings(t, &logged) {
		warned = warned || (ev["level"] == "WARN" && ev["component"] == "stats" && ev["err"] != nil && ev["detail"] != nil && strings.Contains(ev["err"].(string), "conflict log"))
	}
	if !warned {
		t.Errorf("no warning of the bad log: %s", logged.String())
	}
}

// S-0404: a missing merge queue event log leaves the queue out with no
// warning; an empty one is a log with no events; one that cannot be read
// leaves the queue out and is warned of; and a readable one is passed on.
func TestReadLoadsTheMergeQueuesEventLog(t *testing.T) {
	repo := project(t)
	var logged bytes.Buffer
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil || opt.QueueLogged || opt.QueueEvents != nil {
		t.Errorf("a missing log: %v, logged %v, %v", err, opt.QueueLogged, opt.QueueEvents)
	}
	for _, ev := range warnings(t, &logged) {
		if strings.Contains(fmt.Sprint(ev["msg"]), "merge queue") {
			t.Errorf("a missing log is warned of: %v", ev)
		}
	}

	writeFile(t, mergequeue.EventsPath(repo.CacheDir()), "")
	if _, opt, err = Read(execx.System{}, repo, slog.New(slog.DiscardHandler)); err != nil || !opt.QueueLogged || len(opt.QueueEvents) != 0 {
		t.Errorf("an empty log: %v, logged %v, %v", err, opt.QueueLogged, opt.QueueEvents)
	}

	if err := mergequeue.AppendEvents(repo.CacheDir(), []mergequeue.Event{{At: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC), Event: mergequeue.EventEnqueued, Stories: []string{"S-0001"}, By: "alex"}}); err != nil {
		t.Fatal(err)
	}
	if _, opt, err = Read(execx.System{}, repo, slog.New(slog.DiscardHandler)); err != nil || !opt.QueueLogged || len(opt.QueueEvents) != 1 || opt.QueueEvents[0].Event != mergequeue.EventEnqueued {
		t.Errorf("a log of one line: %v, logged %v, %+v", err, opt.QueueLogged, opt.QueueEvents)
	}

	writeFile(t, mergequeue.EventsPath(repo.CacheDir()), "{not json\n")
	logged.Reset()
	_, opt, err = Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if opt.QueueLogged || opt.QueueEvents != nil {
		t.Errorf("a bad log: logged %v, %v", opt.QueueLogged, opt.QueueEvents)
	}
	warned := false
	for _, ev := range warnings(t, &logged) {
		warned = warned || (ev["level"] == "WARN" && ev["component"] == "stats" && ev["err"] != nil && ev["detail"] != nil && strings.Contains(ev["err"].(string), "merge queue's event log"))
	}
	if !warned {
		t.Errorf("no warning of the bad log: %s", logged.String())
	}
}

// keepStatus keeps a pull request status for S-0001 under the project's
// cache, with a push and a check.
func keepStatus(t *testing.T, repo *workitem.Repo) {
	t.Helper()
	st := forge.Status{Story: "S-0001", PullRequest: 1, HeadSHA: "a", ReadAt: "2026-10-01T10:00:00Z", Events: []forge.StatusEvent{
		{At: "2026-10-01T10:00:00Z", Kind: forge.EventHead, To: "a", Commit: "a"},
		{At: "2026-10-01T10:00:00Z", Kind: forge.EventCheck, Name: "build", To: "pending", Commit: "a"},
	}}
	if err := forge.WriteStatus(repo.CacheDir(), st); err != nil {
		t.Fatal(err)
	}
}

// S-0396, ADR-0168: in local mode Read loads no kept forge state, even with
// a status file left from an earlier mode.
func TestReadInLocalModeLoadsNoForgeState(t *testing.T) {
	repo := project(t)
	keepStatus(t, repo)
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if opt.Forge != nil {
		t.Errorf("local mode loaded forge state: %+v", opt.Forge)
	}
}

// S-0396, ADR-0168: under pull-request acceptance Read loads each story's
// kept forge events by its ID, a story with none kept left out; a status it
// cannot read is left out with a warning naming it, and the rest still load.
func TestReadUnderPullRequestAcceptanceLoadsTheForgeEvents(t *testing.T) {
	repo := reopen(t, project(t), "acceptance:\n  mode: pull-request\n")
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if opt.Forge == nil || len(opt.Forge) != 0 {
		t.Errorf("no status kept: %+v, want an empty map", opt.Forge)
	}
	keepStatus(t, repo)
	writeFile(t, forge.StatusPath(repo.CacheDir(), "S-0002"), "{not json\n")
	var logged bytes.Buffer
	_, opt, err = Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if evs, ok := opt.Forge["S-0001"]; len(opt.Forge) != 1 || !ok || len(evs) != 2 || evs[1].Name != "build" {
		t.Errorf("forge = %+v, want S-0001's two events alone", opt.Forge)
	}
	var skipped map[string]any
	for _, ev := range warnings(t, &logged) {
		if ev["msg"] == "cannot read some kept pull request statuses" {
			skipped = ev
		}
	}
	if skipped == nil || skipped["level"] != "WARN" || skipped["component"] != "stats" || skipped["count"] != float64(1) ||
		!strings.Contains(fmt.Sprint(skipped["paths"]), "S-0002.json") || skipped["err"] == nil || skipped["detail"] == nil {
		t.Errorf("want a warning naming the status left out, got %s", logged.String())
	}
}
