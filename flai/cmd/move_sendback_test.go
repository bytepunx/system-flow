package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/sendback"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// storyInReviewOnTheBaseline is a project on the baseline workflow whose
// story S-0001, with a criterion and a task, is in review, as the operator
// moved it.
func storyInReviewOnTheBaseline(t *testing.T) string {
	t.Helper()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	asOperator(t)
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	mustRun(t, root, "config", "set", "author", "alex")
	mustRun(t, root, "epic", "new", "Epic")
	mustRun(t, root, "story", "new", "Sent back", "--epic", "E-0001")
	file := filepath.Join(root, "wip/kanban/stories/S-0001-sent-back.md")
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, root, "task", "new", "Piece", "--story", "S-0001")
	for _, to := range []string{"ready", "in-progress", "review"} {
		mustRun(t, root, "move", "S-0001", to)
	}
	return root
}

// sentBackThreads are the threads on S-0001 in root.
func sentBackThreads(t *testing.T, root string) []*threads.Thread {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	all, err := threads.For(repo, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	return all
}

// I-0397: the operator's move of a story from review back to in-progress, as
// the dashboard's item.move makes it, opens the send back's thread by flai
// with the reason, which flai serve starts the story's agent on.
func TestTheOperatorsSendBackOpensTheThreadFlaiServeStartsTheAgentOn(t *testing.T) {
	root := storyInReviewOnTheBaseline(t)
	out := mustRun(t, root, "--json", "move", "S-0001", "in-progress", "--reason", "tests missing", "--by", "alex")
	var res struct {
		Status string `json:"status"`
		Thread string `json:"sent_back_thread"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	all := sentBackThreads(t, root)
	if len(all) != 1 {
		t.Fatalf("%d threads on S-0001, want the send back's one", len(all))
	}
	th := all[0]
	if res.Status != workitem.InProgress || res.Thread != th.ID {
		t.Errorf("the move answered %+v, want in-progress and %s", res, th.ID)
	}
	if !th.Open() || th.Opener() != storygit.SentBackAuthor || !strings.HasPrefix(th.Title, storygit.SentBackPrefix) {
		t.Errorf("the thread is not the send back's: %q by %s, status %s", th.Title, th.Opener(), th.Status)
	}
	if by, ok := sendback.ByHandMover(th.Title); !ok || by != "alex" || !strings.Contains(th.Entries()[0].Text, "tests missing") {
		t.Errorf("the thread does not say who sent it back and why: %q: %q", th.Title, th.Entries()[0].Text)
	}
}

// I-0397: the story's own agent, by its name or in its session, and flai
// itself open no thread when they send the story back.
func TestTheAgentsOrFlaisSendBackOpensNoThread(t *testing.T) {
	for _, c := range []struct {
		name, by, story string
	}{
		{"by the agent's name", "agent-S-0001", ""},
		{"in the agent's session", "system-flow", "S-0001"},
		{"by flai", storygit.SentBackAuthor, ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := storyInReviewOnTheBaseline(t)
			t.Setenv("FLAI_STORY", c.story)
			out := mustRun(t, root, "move", "S-0001", "in-progress", "--reason", "more to do", "--by", c.by)
			if all := sentBackThreads(t, root); len(all) != 0 || strings.Contains(out, "TH-") {
				t.Errorf("%d threads on S-0001, want none:\n%s", len(all), out)
			}
		})
	}
}
