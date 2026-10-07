package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// I-0073: an acceptance resolves the threads still open or answered on what
// it archives (the story, its tasks, and the epic that followed it), in the
// acceptance commit, so flai check does not warn threads.archived after it.
// A thread on another item is left as it is, and the dry run names the
// threads and changes nothing.
func TestAcceptResolvesTheThreadsOnWhatItArchives(t *testing.T) {
	root := storyInReviewWithThreads(t)
	status, head := gitIn(t, root, "status", "--porcelain"), gitIn(t, root, "rev-parse", "HEAD")
	paths := threadPaths(t, root)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the dry run: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if got := strings.Join(res.ResolvedThreads, " "); got != "TH-0001 TH-0002 TH-0003" || len(res.Blockers) != 0 {
		t.Errorf("the dry run names the threads it would resolve, not as blockers: %q %v", got, res.Blockers)
	}
	if out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run"); !strings.Contains(out, "would resolve TH-0001, TH-0002, TH-0003, open or answered on what it archives\n") {
		t.Errorf("the dry run's text names them:\n%s", out)
	}
	if gitIn(t, root, "status", "--porcelain") != status || gitIn(t, root, "rev-parse", "HEAD") != head {
		t.Error("the dry run changes nothing")
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 {
		t.Fatalf("the acceptance: %d %s\n%s", code, errOut, out)
	}
	if !strings.Contains(out, "resolved TH-0001, TH-0002, TH-0003, open or answered on what it archived\n") {
		t.Errorf("the output names the threads resolved:\n%s", out)
	}
	assertResolvedByAcceptance(t, root, "TH-0001", "TH-0002", "TH-0003")
	if th := thread(t, root, "TH-0004"); th.Status != "open" {
		t.Errorf("a thread on another item is left open: %s", th.Status)
	}
	committed := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	for _, id := range []string{"TH-0001", "TH-0002", "TH-0003"} {
		if !strings.Contains(committed, paths[id]) {
			t.Errorf("%s is in the acceptance commit:\n%s", id, committed)
		}
	}
	if strings.Contains(committed, paths["TH-0004"]) {
		t.Errorf("a thread the acceptance did not resolve is not in its commit:\n%s", committed)
	}
	if s := gitIn(t, root, "status", "--porcelain"); s != "" {
		t.Errorf("everything is committed: %q", s)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "threads.archived") {
		t.Errorf("flai check finds no thread open on an archived item:\n%s", out)
	}
}

// An acceptance whose commit failed is finished by running it again, and the
// finish resolves a thread left open on what it archived, as an older flai's
// acceptance left them, into the commit it makes.
func TestFinishingAnAcceptanceResolvesTheThreadsLeftOpen(t *testing.T) {
	root := storyInReviewWithThreads(t)
	paths := threadPaths(t, root)
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	_ = os.MkdirAll(filepath.Dir(hook), 0o755)
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 {
		t.Fatal("the acceptance whose commit failed exits non-zero")
	}
	// TH-0002 as the acceptance found it, open on the archived task
	gitIn(t, root, "checkout", "HEAD", "--", paths["TH-0002"])
	if th := thread(t, root, "TH-0002"); th.Status != "open" {
		t.Fatalf("TH-0002 is open again: %s", th.Status)
	}
	if out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run"); !strings.Contains(out, "would resolve TH-0002, open or answered on what it archives\n") {
		t.Errorf("the dry run names the thread the finish would resolve:\n%s", out)
	}

	_ = os.Remove(hook)
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 || !strings.Contains(out, "completed the acceptance of S-0001") || !strings.Contains(out, "resolved TH-0002, open or answered on what it archived\n") {
		t.Fatalf("run again, the acceptance resolves the thread and commits: %d %s %s", code, out, errOut)
	}
	assertResolvedByAcceptance(t, root, "TH-0001", "TH-0002", "TH-0003")
	committed := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, paths["TH-0002"]) {
		t.Errorf("TH-0002 is in the acceptance commit:\n%s", committed)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "threads.archived") {
		t.Errorf("flai check finds no thread open on an archived item:\n%s", out)
	}
}

// storyInReviewWithThreads is orchestratedStoryInReview's project, whose
// S-0001 is its epic E-0001's only story, with S-0002 under no epic and these
// threads, committed: TH-0001 alex's, answered on S-0001, and claude's open
// ones, TH-0002 on its task T-0001, TH-0003 on E-0001, and TH-0004 on S-0002.
// The clock is fixed, so an entry by alex after alex's would join it.
func storyInReviewWithThreads(t *testing.T) string {
	t.Helper()
	root, _ := orchestratedStoryInReview(t, false)
	for _, args := range [][]string{
		{"story", "new", "Other"},
		{"thread", "new", "--on", "S-0001", "--by", "alex", "Is the flag right?", "Say which."},
		{"thread", "new", "--on", "T-0001", "--by", "claude", "Is the task done?", "Say so."},
		{"thread", "new", "--on", "E-0001", "--by", "claude", "Is the epic whole?", "Say so."},
		{"thread", "new", "--on", "S-0002", "--by", "claude", "When is it next?", "Say when."},
		{"thread", "reply", "TH-0001", "--by", "claude", "Use the flag."},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	if th := thread(t, root, "TH-0001"); th.Status != "answered" {
		t.Fatalf("TH-0001 is answered: %s", th.Status)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: threads")
	return root
}

// assertResolvedByAcceptance fails unless each thread is resolved by alex
// with S-0001's acceptance as the reason.
func assertResolvedByAcceptance(t *testing.T, root string, ids ...string) {
	t.Helper()
	for _, id := range ids {
		th := thread(t, root, id)
		es := th.Entries()
		last := es[len(es)-1]
		if th.Status != "resolved" || last.Author != "alex" || last.Text != "Resolved: S-0001 was accepted" {
			t.Errorf("%s is resolved by the acceptance: status %s, last entry %+v", id, th.Status, last)
		}
	}
}

// threadPaths are the threads' files relative to root, by ID.
func threadPaths(t *testing.T, root string) map[string]string {
	t.Helper()
	paths := map[string]string{}
	for _, id := range []string{"TH-0001", "TH-0002", "TH-0003", "TH-0004"} {
		paths[id] = filepath.ToSlash(mustRel(t, root, thread(t, root, id).Path))
	}
	return paths
}

// thread reads thread id of the project at root.
func thread(t *testing.T, root, id string) *threads.Thread {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	th, err := threads.Get(repo, id)
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// ADR-0120: an acceptance closes the conversations still open of the story it
// accepts, each with the entry "Closed: S-0001 was accepted" by whoever
// accepts, in the acceptance commit. A conversation between two other
// stories, and one closed already, are left as they are, and the dry run
// names the conversations it would close and changes nothing.
func TestAcceptClosesTheConversationsOfTheStory(t *testing.T) {
	root, paths := storyInReviewWithConversations(t)
	status, head := gitIn(t, root, "status", "--porcelain"), gitIn(t, root, "rev-parse", "HEAD")

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the dry run: %d %s", code, errOut)
	}
	var res struct {
		Blockers []string `json:"blockers"`
		Closed   []string `json:"closed_messages"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if got := strings.Join(res.Closed, " "); got != "MS-0001 MS-0003" || len(res.Blockers) != 0 {
		t.Errorf("the dry run names the conversations it would close, not as blockers: %q %v", got, res.Blockers)
	}
	out, _, _ = runIn(t, root, "accept", "S-0001", "--dry-run")
	if !strings.Contains(out, "would close MS-0001, open with S-0002\nwould close MS-0003, open with S-0003\n") {
		t.Errorf("the dry run's text names each with the story it is open with:\n%s", out)
	}
	if gitIn(t, root, "status", "--porcelain") != status || gitIn(t, root, "rev-parse", "HEAD") != head {
		t.Error("the dry run changes nothing")
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("the acceptance: %d %s\n%s", code, errOut, out)
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || strings.Join(res.Closed, " ") != "MS-0001 MS-0003" {
		t.Errorf("--json lists the conversations closed in closed_messages: %v %s", err, out)
	}
	assertClosedAs(t, root, "Closed: S-0001 was accepted", "MS-0001", "MS-0003")
	for _, id := range []string{"MS-0002", "MS-0004"} {
		if now, was := strings.TrimSpace(conversation(t, root, id).Marshal()), gitIn(t, root, "show", head+":"+paths[id]); now != was {
			t.Errorf("%s, not open on S-0001, is left as it was:\n%s", id, now)
		}
	}
	committed := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	for id, in := range map[string]bool{"MS-0001": true, "MS-0003": true, "MS-0002": false, "MS-0004": false} {
		if strings.Contains(committed, paths[id]) != in {
			t.Errorf("%s in the acceptance commit is %v:\n%s", id, !in, committed)
		}
	}
	if s := gitIn(t, root, "status", "--porcelain"); s != "" {
		t.Errorf("everything is committed: %q", s)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "messages.") {
		t.Errorf("flai check finds nothing in the conversations:\n%s", out)
	}
}

// An acceptance whose commit failed is finished by running it again, and the
// finish closes a conversation left open on the story it accepted, as an older
// flai's acceptance left them, into the commit it makes.
func TestFinishingAnAcceptanceClosesTheConversationsLeftOpen(t *testing.T) {
	root, paths := storyInReviewWithConversations(t)
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	_ = os.MkdirAll(filepath.Dir(hook), 0o755)
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code == 0 {
		t.Fatal("the acceptance whose commit failed exits non-zero")
	}
	// MS-0001 as the acceptance found it, open with S-0002
	gitIn(t, root, "checkout", "HEAD", "--", paths["MS-0001"])
	if c := conversation(t, root, "MS-0001"); c.Status != messages.StatusOpen {
		t.Fatalf("MS-0001 is open again: %s", c.Status)
	}
	if out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run"); !strings.Contains(out, "would close MS-0001, open with S-0002\n") || strings.Contains(out, "MS-0003") {
		t.Errorf("the dry run names the conversation the finish would close, and only it:\n%s", out)
	}

	_ = os.Remove(hook)
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 || !strings.Contains(out, "completed the acceptance of S-0001") || !strings.Contains(out, "closed MS-0001, open with S-0002\n") {
		t.Fatalf("run again, the acceptance closes the conversation and commits: %d %s %s", code, out, errOut)
	}
	assertClosedAs(t, root, "Closed: S-0001 was accepted", "MS-0001", "MS-0003")
	if committed := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); !strings.Contains(committed, paths["MS-0001"]) {
		t.Errorf("MS-0001 is in the acceptance commit:\n%s", committed)
	}
}

// ADR-0120: messages are kept apart from threads. With a conversation and a
// thread in the project, flai thread list, with --all or not, lists the
// thread only, and the inbox's threads, awaiting anyone, hold no conversation.
func TestMessagesStayApartFromThreads(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "alex")
	root := tempProject(t)
	for _, args := range [][]string{
		{"story", "new", "Slice"},
		{"story", "new", "Other"},
		{"thread", "new", "--on", "S-0001", "--by", "claude", "Is the flag right?", "Say which."},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	writeConversation(t, root, "MS-0001", "S-0001", "S-0002", messages.StatusOpen)
	for _, args := range [][]string{{"thread", "list"}, {"thread", "list", "--all"}, {"thread", "list", "--json"}, {"thread", "list", "--all", "--json"}} {
		out, errOut, code := runIn(t, root, args...)
		if code != 0 || !strings.Contains(out, "TH-0001") || strings.Contains(out, "MS-0001") {
			t.Errorf("flai %v lists the thread and no conversation: %d %s\n%s", args, code, errOut, out)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, agent := range []string{"alex", "claude", "agent-S-0001", "agent-S-0002"} {
		box, err := inbox.Read(context.Background(), inbox.Options{Repo: repo, Agent: agent, All: true})
		if err != nil {
			t.Fatal(err)
		}
		if len(box.Threads) != 1 || box.Threads[0].ID != "TH-0001" {
			t.Errorf("%s's inbox holds the thread and no conversation: %+v", agent, box.Threads)
		}
	}
}

// storyInReviewWithConversations is orchestratedStoryInReview's project, whose
// S-0001 is in review, with S-0002 and S-0003 and these conversations,
// committed: MS-0001 open from S-0001 to S-0002, MS-0002 open between S-0002
// and S-0003, MS-0003 open from S-0003 to S-0001, and MS-0004 closed between
// S-0001 and S-0002. It returns the root and the conversations' files
// relative to it, by ID.
func storyInReviewWithConversations(t *testing.T) (string, map[string]string) {
	t.Helper()
	root, _ := orchestratedStoryInReview(t, false)
	for _, args := range [][]string{{"story", "new", "Other"}, {"story", "new", "Third"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	paths := map[string]string{
		"MS-0001": writeConversation(t, root, "MS-0001", "S-0001", "S-0002", messages.StatusOpen),
		"MS-0002": writeConversation(t, root, "MS-0002", "S-0002", "S-0003", messages.StatusOpen),
		"MS-0003": writeConversation(t, root, "MS-0003", "S-0003", "S-0001", messages.StatusOpen),
		"MS-0004": writeConversation(t, root, "MS-0004", "S-0001", "S-0002", messages.StatusClosed),
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: messages")
	return root, paths
}

// writeConversation writes conversation id from story from to story to, with
// one message from from's agent, as flai message send writes one, and the
// status given. It returns the file relative to root.
func writeConversation(t *testing.T, root, id, from, to, status string) string {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	at := "2026-09-15T20:00:00Z"
	c := &messages.Conversation{
		ID: id, Title: "Who changes the parser", From: from, To: to, Status: status,
		Participants: []string{"agent-" + from}, Created: at, Updated: at,
		Path: filepath.Join(messages.Dir(repo), id+"-who-changes-the-parser.md"),
		Body: fmt.Sprintf("\n# %s Who changes the parser\n\nBetween %s and %s.\n\n## Entries\n\n### %s agent-%s %s\nWho changes the parser?\n", id, from, to, at, from, from),
	}
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(mustRel(t, root, c.Path))
}

// assertClosedAs fails unless each conversation is closed by alex with the
// entry text given.
func assertClosedAs(t *testing.T, root, text string, ids ...string) {
	t.Helper()
	for _, id := range ids {
		c := conversation(t, root, id)
		es := c.Entries()
		last := es[len(es)-1]
		if c.Status != messages.StatusClosed || last.Author != "alex" || last.Story != "" || last.Text != text {
			t.Errorf("%s is closed with %q: status %s, last entry %+v", id, text, c.Status, last)
		}
	}
}

// conversation reads conversation id of the project at root.
func conversation(t *testing.T, root, id string) *messages.Conversation {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	c, err := messages.Get(repo, id)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
