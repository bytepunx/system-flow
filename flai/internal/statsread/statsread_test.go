package statsread

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
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
	} {
		if out, err := r.Run(repo.Root, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	// the items lock lives in .flai-cache, which the commit must not carry
	if err := os.WriteFile(filepath.Join(repo.Root, ".git", "info", "exclude"), []byte(".flai-cache/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "feat: [S-0001] the story"}} {
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

// logCounter is the real runner, counting the git log processes it starts.
type logCounter struct {
	execx.System
	logs int
}

func (c *logCounter) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && slices.Contains(args, "log") {
		c.logs++
	}
	return c.System.Run(dir, name, args...)
}

// I-0184: a second Read with nothing changed in git starts no git log.
func TestReadStartsNoLogWhileNothingChanged(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	repo := project(t)
	r := &logCounter{}
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "system-flow.yaml"},
		{"commit", "-q", "-m", "feat: [S-0001] the story"},
	} {
		if out, err := r.Run(repo.Root, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	for i := range 2 {
		_, opt, err := Read(r, repo, slog.New(slog.DiscardHandler))
		if err != nil {
			t.Fatal(err)
		}
		if got := opt.Commits["S-0001"]; len(got) != 1 || got[0] != "system-flow.yaml" {
			t.Errorf("read %d: commits = %v, want the manifest", i+1, opt.Commits)
		}
	}
	if r.logs != 1 {
		t.Errorf("two reads with nothing changed started git log %d times, want once", r.logs)
	}
}

// freshRootEnv names the project TestReadInANewProcess reads.
const freshRootEnv = "FLAI_STATSREAD_FRESH_ROOT"

// I-0209: a Read in a new process, with nothing changed in git since one
// before it read the log, starts no git log. The new process is this test
// binary again, running TestReadInANewProcess alone.
func TestReadStartsNoLogInANewProcessWhileNothingChanged(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	repo := project(t)
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "system-flow.yaml"},
		{"commit", "-q", "-m", "feat: [S-0001] the story"},
	} {
		if out, err := (execx.System{}).Run(repo.Root, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	if _, _, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler)); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestReadInANewProcess$", "-test.count=1", "-test.v")
	cmd.Env = append(os.Environ(), freshRootEnv+"="+repo.Root)
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "--- PASS: TestReadInANewProcess") {
		t.Errorf("the new process did not pass (%v):\n%s", err, out)
	}
}

// TestReadInANewProcess is the new process of
// TestReadStartsNoLogInANewProcessWhileNothingChanged: it reads the project
// that names, with nothing kept in memory, and wants the story's commits
// answered without a git log.
func TestReadInANewProcess(t *testing.T) {
	root := os.Getenv(freshRootEnv)
	if root == "" {
		t.Skip("the new process of TestReadStartsNoLogInANewProcessWhileNothingChanged")
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	r := &logCounter{}
	_, opt, err := Read(r, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	if got := opt.Commits["S-0001"]; len(got) != 1 || got[0] != "system-flow.yaml" {
		t.Errorf("commits = %v, want the manifest", opt.Commits)
	}
	if r.logs != 0 {
		t.Errorf("a new process with nothing changed started git log %d times, want none", r.logs)
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

// S-0607: the reads run at once, and when two fail the error is still the
// one of the read that comes first, the threads before the board, whichever
// ends first; the board's error is the answer once the threads read.
func TestReadReturnsTheFirstReadsErrorWhenTwoFail(t *testing.T) {
	repo := project(t)
	thread := filepath.Join(threads.Dir(repo), "TH-0009-bad.md")
	if err := os.WriteFile(thread, []byte("---\nid: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo.KanbanDir(), workitem.BoardFile), []byte("---\nwip_limits: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := range 20 {
		_, _, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
		if err == nil || !strings.Contains(err.Error(), "cannot read the threads") || strings.Contains(err.Error(), "in-progress limit") {
			t.Fatalf("read %d with an unreadable thread and board: %v", i, err)
		}
	}
	if err := os.Remove(thread); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler)); err == nil || !strings.Contains(err.Error(), "cannot read the board's in-progress limit") {
		t.Errorf("an unreadable board alone: %v", err)
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

// digestEntryDoc is a digest entry with the ID and status given, dated
// 2026-09-30, that lists TH-0001 in asked_again.
func digestEntryDoc(id, status string) string {
	return "---\nid: " + id + "\ntitle: Left it is.\nkind: answer\nsources: [TH-0001]\ntopics: []\npaths: []\n" +
		"date: 2026-09-30\nstatus: " + status + "\nabsorbed_by:\nasked_again: [TH-0001]\ncited: []\n---\n\n# " + id + " Left it is.\n\nLeft.\n"
}

// S-0434: Read passes the manifest's owner and the digest's entries; a
// missing digest folder is no entries and no warning, and a digest that
// cannot be read sets DigestUnreadable and is warned of.
func TestReadLoadsTheOwnerAndTheDigestsEntries(t *testing.T) {
	repo := reopen(t, project(t), "owner: alex\n")
	var logged bytes.Buffer
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil || opt.Owner != "alex" || opt.Entries != nil || opt.DigestUnreadable {
		t.Errorf("no digest folder: %v, owner %q, entries %v, unreadable %v", err, opt.Owner, opt.Entries, opt.DigestUnreadable)
	}
	for _, ev := range warnings(t, &logged) {
		if strings.Contains(fmt.Sprint(ev["msg"]), "digest") {
			t.Errorf("a missing digest is warned of: %v", ev)
		}
	}

	dir := filepath.Join(repo.Root, filepath.FromSlash(digest.Dir(repo.Manifest)))
	writeFile(t, filepath.Join(dir, "README.md"), "# Digest\n")
	writeFile(t, filepath.Join(dir, "D-0001-left-it-is.md"), digestEntryDoc("D-0001", "active"))
	writeFile(t, filepath.Join(dir, "D-0002-left-again.md"), digestEntryDoc("D-0002", "compacted"))
	_, opt, err = Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err != nil || opt.DigestUnreadable || len(opt.Entries) != 2 || opt.Entries[0].Status != "active" || opt.Entries[1].Status != "compacted" ||
		opt.Entries[0].Date != "2026-09-30" || len(opt.Entries[0].AskedAgain) != 1 || opt.Entries[0].AskedAgain[0] != "TH-0001" {
		t.Errorf("two entries: %v, unreadable %v, %+v", err, opt.DigestUnreadable, opt.Entries)
	}

	writeFile(t, filepath.Join(dir, "D-0003-bad.md"), "---\nid: [\n---\n")
	logged.Reset()
	_, opt, err = Read(execx.System{}, repo, slog.New(slog.NewJSONHandler(&logged, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if !opt.DigestUnreadable || opt.Entries != nil {
		t.Errorf("a bad entry: unreadable %v, %+v", opt.DigestUnreadable, opt.Entries)
	}
	warned := false
	for _, ev := range warnings(t, &logged) {
		warned = warned || (ev["level"] == "WARN" && ev["component"] == "stats" && ev["msg"] == "cannot read the digest" && ev["err"] != nil && ev["detail"] != nil)
	}
	if !warned {
		t.Errorf("no warning of the bad digest: %s", logged.String())
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

// S-0546: a missing serve journal leaves the close-outs out with no warning;
// one with a line that is not an entry leaves them out and is warned of; and
// a readable one yields the figures of the project's close-out end entries,
// its main checkout's, other actions and projects ignored.
func TestCloseOutsReadsTheServeJournal(t *testing.T) {
	t.Parallel()
	repo := project(t)
	path := filepath.Join(t.TempDir(), "serve", "journal.jsonl")
	var logged bytes.Buffer
	if j := CloseOuts(path, repo, slog.New(slog.NewJSONHandler(&logged, nil))); j != nil || logged.Len() != 0 {
		t.Errorf("a missing journal: %+v, logged %s", j, logged.String())
	}

	root := filepath.Dir(repo.CacheDir())
	lines := []string{
		`{"at":"2026-10-01T10:00:00Z","action":"close-out","method":"verify.job","project":"t","root":"` + root + `","by":"agent-S-0001","outcome":"started","detail":"S-0001 job S-0001-20261001T100000Z"}`,
		`{"at":"2026-10-01T10:05:00Z","action":"close-out","method":"verify.job","project":"t","root":"` + root + `","by":"agent-S-0001","outcome":"done","detail":"S-0001 job","duration_ms":300000,"verdict":"passed"}`,
		`{"at":"2026-10-01T10:06:00Z","action":"close-out","method":"verify.job","project":"t","root":"` + root + `","by":"agent-S-0001","outcome":"done","detail":"from the record: S-0001's close-out passed at abc1234, run 2026-10-01T10:05:00Z","verdict":"passed"}`,
		`{"at":"2026-10-01T10:07:00Z","action":"close-out","method":"verify.job","project":"u","root":"/elsewhere","by":"agent-S-0009","outcome":"failed","duration_ms":1000,"verdict":"failed"}`,
		`{"at":"2026-10-01T10:08:00Z","action":"checks","method":"verify.run","project":"t","root":"` + root + `","by":"alex","outcome":"done"}`,
		``,
	}
	writeFile(t, path, strings.Join(lines, "\n"))
	j := CloseOuts(path, repo, slog.New(slog.DiscardHandler))
	if j == nil || j.Root != root || len(j.Entries) != 5 {
		t.Fatalf("a journal: %+v", j)
	}
	_, opt, err := Read(execx.System{}, repo, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	opt.CloseOuts, opt.Now, opt.Since = j, time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC), 7*24*time.Hour
	rep := metrics.Compute(nil, opt)
	five := 300.0
	want := &metrics.CloseOuts{Weeks: []metrics.CloseOutWeek{{Week: "2026-W39", Start: "2026-09-21"}, {Week: "2026-W40", Start: "2026-09-28", Passed: 1, Recorded: 1, P50: &five, P85: &five}}}
	if !reflect.DeepEqual(rep.CloseOuts, want) {
		got, _ := json.Marshal(rep.CloseOuts)
		t.Errorf("the figures: %s", got)
	}

	writeFile(t, path, lines[1]+"\n{not json\n")
	logged.Reset()
	if j := CloseOuts(path, repo, slog.New(slog.NewJSONHandler(&logged, nil))); j != nil {
		t.Errorf("a bad journal: %+v", j)
	}
	warned := false
	for _, ev := range warnings(t, &logged) {
		warned = warned || (ev["level"] == "WARN" && ev["component"] == "stats" && ev["detail"] != nil && strings.Contains(fmt.Sprint(ev["err"]), "line 2"))
	}
	if !warned {
		t.Errorf("no warning of the bad journal naming its line: %s", logged.String())
	}
}
