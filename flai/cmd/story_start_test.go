package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// storyStartConvention is the one convention the start tests' project has,
// so that a context pack can be built.
const storyStartConvention = `---
title: Git
updated: 2026-08-01
audience: agent
order: 10
status: active
---

# Git

## Rules
- Commit at landing.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// storyStartProject is a project outside git, with conventions, holding four
// stories, each with criteria and touches of its own: S-0001 and S-0002
// ready, S-0003 in backlog, and S-0004 ready but held, as it names S-0001 in
// after:. FLAI_AGENT is agent-a.
func storyStartProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "agent-a")
	t.Setenv("FLAI_SESSION", "s1")
	t.Setenv("LOG_FORMAT", "json")
	root := tempProject(t)
	writeIn(t, root, "design/conventions/README.md", "# Agent conventions\n\n| Order | File | Governs |\n|-------|------|---------|\n| 10 | [git.md](git.md) | Commits |\n")
	writeIn(t, root, "design/conventions/git.md", storyStartConvention)
	writeIn(t, root, "design/system/plan.md", "---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\ntext\n")
	for _, s := range []struct{ title, touches, after, to string }{
		{"One", "docs/one", "", workitem.Ready},
		{"Two", "docs/two", "", workitem.Ready},
		{"Three", "docs/three", "", ""},
		{"Four", "docs/four", "S-0001", workitem.Ready},
	} {
		args := []string{"story", "new", s.title, "--touches", s.touches}
		if s.after != "" {
			args = append(args, "--after", s.after)
		}
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("story new %s: %d %s %s", s.title, code, out, errOut)
		}
		id := strings.Fields(out)[0]
		p := filepath.Join(root, "wip", "kanban", "stories", id+"-"+strings.ToLower(s.title)+".md")
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(strings.Replace(string(data), "- [ ]\n", "- [ ] it works\n", 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		if s.to != "" {
			if _, errOut, code := runIn(t, root, "move", id, s.to, "--by", "alex"); code != 0 {
				t.Fatalf("move %s %s: %s", id, s.to, errOut)
			}
		}
	}
	return root
}

func storyStartGet(t *testing.T, root, id string) *workitem.Item {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	return it
}

// fatalErr is the err field of the fatal event the command boundary logged
// on stderr, "" when there is none.
func fatalErr(t *testing.T, errOut string) string {
	t.Helper()
	sc := bufio.NewScanner(strings.NewReader(errOut))
	for sc.Scan() {
		var ev map[string]any
		if json.Unmarshal(sc.Bytes(), &ev) == nil && ev["level"] == "FATAL" {
			s, _ := ev["err"].(string)
			return s
		}
	}
	return ""
}

// S-0274: a ready story starts in one call: it is in progress, moved by
// FLAI_AGENT, and the output is the move, the narrative, the branch (none
// outside git), the pack, and the inbox, in that order.
func TestStoryStartStartsAReadyStoryAsText(t *testing.T) {
	root := storyStartProject(t)
	out, errOut, code := runIn(t, root, "story", "start", "S-0001")
	if code != 0 {
		t.Fatalf("story start: %d\n%s\n%s", code, out, errOut)
	}
	at := 0
	for _, want := range []string{
		"S-0001 → in-progress\n",
		"narrative wip/agents/S-0001.md\n",
		"no branch or worktree: the project is not a git work tree",
		"S-0001 context pack\n",
		"story: S-0001 One\n",
		"inbox: 0 threads awaiting agent-a\n",
		"ready to pull, in pull order",
		"since you last looked\n",
	} {
		i := strings.Index(out[at:], want)
		if i < 0 {
			t.Fatalf("output lacks %q after byte %d:\n%s", want, at, out)
		}
		at += i + len(want)
	}
	it := storyStartGet(t, root, "S-0001")
	if last := it.Transitions[len(it.Transitions)-1]; it.Status != workitem.InProgress || last.By != "agent-a" {
		t.Errorf("S-0001 is %s, last moved by %s", it.Status, last.By)
	}
	n, err := workitem.ReadNarrative(filepath.Join(root, "wip", "agents", "S-0001.md"))
	if err != nil || n.Agent != "agent-a" || n.Session != "s1" {
		t.Errorf("narrative: %+v, %v", n, err)
	}
}

// S-0274: --json prints the result as one object with every key.
func TestStoryStartAnswersJSON(t *testing.T) {
	root := storyStartProject(t)
	out, errOut, code := runIn(t, root, "story", "start", "S-0002", "--json", "--budget", "40KB")
	if code != 0 {
		t.Fatalf("story start --json: %d\n%s\n%s", code, out, errOut)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not one JSON object: %v\n%s", err, out)
	}
	for _, key := range []string{"story", "followed", "worktree", "branch", "from", "pack", "inbox"} {
		if _, ok := got[key]; !ok {
			t.Errorf("no %q key: %s", key, out)
		}
	}
	if len(got) != 7 {
		t.Errorf("want 7 keys, got %d: %s", len(got), out)
	}
	var story struct{ ID, Status string }
	var pack struct {
		Story  string `json:"story"`
		Budget int    `json:"budget"`
	}
	var inbox struct {
		Agent string `json:"agent"`
	}
	if json.Unmarshal(got["story"], &story) != nil || story.ID != "S-0002" || story.Status != workitem.InProgress {
		t.Errorf("story: %s", got["story"])
	}
	if json.Unmarshal(got["pack"], &pack) != nil || pack.Story != "S-0002" || pack.Budget != 40*1024 {
		t.Errorf("pack: story %q, budget %d", pack.Story, pack.Budget)
	}
	if json.Unmarshal(got["inbox"], &inbox) != nil || inbox.Agent != "agent-a" {
		t.Errorf("inbox: %s", got["inbox"])
	}
	if it := storyStartGet(t, root, "S-0002"); it.Status != workitem.InProgress {
		t.Errorf("S-0002 is %s", it.Status)
	}
}

// S-0274: a story in backlog is refused with exit 4, its state in the
// reason, and nothing changes: no move, no narrative.
func TestStoryStartRefusesABacklogStory(t *testing.T) {
	root := storyStartProject(t)
	before := storyStartGet(t, root, "S-0003")
	out, errOut, code := runIn(t, root, "story", "start", "S-0003")
	if code != exitDocRefused {
		t.Fatalf("want exit %d, got %d\n%s\n%s", exitDocRefused, code, out, errOut)
	}
	if reason := fatalErr(t, errOut); !strings.Contains(reason, "S-0003 is backlog, not ready") {
		t.Errorf("reason: %q\n%s", reason, errOut)
	}
	storyStartUnchanged(t, root, before)
}

// S-0274: a ready story the board holds is refused with exit 4 and its
// hold's reason, as text and as --json, and nothing changes.
func TestStoryStartRefusesAHeldStory(t *testing.T) {
	root := storyStartProject(t)
	before := storyStartGet(t, root, "S-0004")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	hold := repo.Holds(items).Of(before)
	if hold == nil || hold.Code != workitem.HoldAfter {
		t.Fatalf("the fixture's S-0004 should be held by after: %+v", hold)
	}

	out, errOut, code := runIn(t, root, "story", "start", "S-0004")
	if code != exitDocRefused {
		t.Fatalf("want exit %d, got %d\n%s\n%s", exitDocRefused, code, out, errOut)
	}
	if reason := fatalErr(t, errOut); !strings.Contains(reason, hold.Reason) {
		t.Errorf("reason %q lacks the hold's %q", reason, hold.Reason)
	}
	storyStartUnchanged(t, root, before)

	out, errOut, code = runIn(t, root, "story", "start", "S-0004", "--json")
	if code != exitDocRefused {
		t.Fatalf("--json: want exit %d, got %d\n%s\n%s", exitDocRefused, code, out, errOut)
	}
	var got struct {
		Refused struct {
			Story  string         `json:"story"`
			Status string         `json:"status"`
			Hold   *workitem.Hold `json:"hold"`
			Reason string         `json:"reason"`
		} `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	r := got.Refused
	if r.Story != "S-0004" || r.Status != workitem.Ready || r.Hold == nil || !reflect.DeepEqual(*r.Hold, *hold) || !strings.Contains(r.Reason, hold.Reason) {
		t.Errorf("refused: %+v", r)
	}
	storyStartUnchanged(t, root, before)
}

// storyStartUnchanged fails when the story moved or has a narrative.
func storyStartUnchanged(t *testing.T, root string, before *workitem.Item) {
	t.Helper()
	after := storyStartGet(t, root, before.ID)
	if after.Status != before.Status || len(after.Transitions) != len(before.Transitions) {
		t.Errorf("%s changed: %s with %d transitions, was %s with %d", before.ID, after.Status, len(after.Transitions), before.Status, len(before.Transitions))
	}
	if _, err := os.Stat(filepath.Join(root, "wip", "agents", before.ID+".md")); err == nil {
		t.Errorf("%s has a narrative", before.ID)
	}
}

// S-0274: the help says what the command does and the four calls it
// replaces at the start of a story.
func TestStoryStartHelpNamesWhatItReplaces(t *testing.T) {
	out, _, code := runCLI(t, "story", "start", "--help")
	if code != 0 {
		t.Fatalf("help: %d", code)
	}
	for _, want := range []string{"flai move S-nnnn in-progress", "flai stream open S-nnnn", "flai prime --story S-nnnn", "the MCP tool inbox", "--budget"} {
		if !strings.Contains(out, want) {
			t.Errorf("help lacks %q:\n%s", want, out)
		}
	}
}

// storyStartThreads opens, as alex, 32 threads on other items, S-0002 and
// a design document, each with a last entry of 1,500 bytes, and one on
// S-0001, which it returns: S-0399's instance (I-0167).
func storyStartThreads(t *testing.T, root string) string {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)
	long := strings.Repeat("Which library do we take here? ", 50)[:1500]
	for i := range 32 {
		on := "S-0002"
		if i%2 == 1 {
			on = "design/system/plan.md"
		}
		if _, err := threads.New(repo, threads.NewOptions{Title: fmt.Sprintf("Elsewhere %d", i), On: on, Author: "alex", Text: long, Now: at}); err != nil {
			t.Fatal(err)
		}
	}
	th, err := threads.New(repo, threads.NewOptions{Title: "On the story", On: "S-0001", Author: "alex", Text: long, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	return th.ID
}

// ADR-0179: --json lists the story's threads alone, counts the others in
// threads_omitted, and is at most 40,000 bytes, as story_start answers.
func TestStoryStartJSONListsOnlyTheStorysThreads(t *testing.T) {
	root := storyStartProject(t)
	mine := storyStartThreads(t, root)
	out, errOut, code := runIn(t, root, "story", "start", "S-0001", "--json")
	if code != 0 {
		t.Fatalf("story start --json: %d\n%s\n%s", code, out, errOut)
	}
	if n := len(out); n > ctxpack.PartLimit {
		t.Errorf("the answer is %d bytes, over %d", n, ctxpack.PartLimit)
	}
	var got struct {
		Inbox struct {
			AwaitingYou    int                   `json:"awaiting_you"`
			Threads        []struct{ ID string } `json:"threads"`
			ThreadsOmitted int                   `json:"threads_omitted"`
		} `json:"inbox"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	in := got.Inbox
	if len(in.Threads) != 1 || in.Threads[0].ID != mine || in.AwaitingYou != 1 || in.ThreadsOmitted != 32 {
		t.Errorf("inbox: %d awaiting, %d omitted, threads %+v; want %s alone and 32 omitted", in.AwaitingYou, in.ThreadsOmitted, in.Threads, mine)
	}
}

// ADR-0179: the text form lists the story's threads and says how many are
// on other items, and that flai thread list lists them.
func TestStoryStartSaysHowManyThreadsAreElsewhere(t *testing.T) {
	root := storyStartProject(t)
	mine := storyStartThreads(t, root)
	out, errOut, code := runIn(t, root, "story", "start", "S-0001")
	if code != 0 {
		t.Fatalf("story start: %d\n%s\n%s", code, out, errOut)
	}
	for _, want := range []string{"inbox: 1 thread awaiting agent-a\n", "  " + mine + " On the story", "inbox: 32 threads on other items: flai thread list lists them\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Elsewhere") {
		t.Errorf("output lists a thread of another item:\n%s", out)
	}
}
