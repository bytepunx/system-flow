package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// streamOpenMoved is the moved key of flai stream open --json's output.
func streamOpenMoved(t *testing.T, out string) bool {
	t.Helper()
	var got struct {
		Moved *bool `json:"moved"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || got.Moved == nil {
		t.Fatalf("no moved key in the JSON: %v\n%s", err, out)
	}
	return *got.Moved
}

// I-0097, ADR-0153: flai stream open on a ready story starts it as flai
// story start does: it is in progress, moved by FLAI_AGENT, and the output
// says so before the narrative; --json gives moved true.
func TestStreamOpenStartsAReadyStory(t *testing.T) {
	root := storyStartProject(t)
	out, errOut, code := runIn(t, root, "stream", "open", "S-0001")
	if code != 0 {
		t.Fatalf("stream open: %d\n%s\n%s", code, out, errOut)
	}
	if !strings.HasPrefix(out, "S-0001 → in-progress\n") || !strings.Contains(out, "opened wip/agents/S-0001.md\n") {
		t.Errorf("output should say the move, then the narrative:\n%s", out)
	}
	it := storyStartGet(t, root, "S-0001")
	if last := it.Transitions[len(it.Transitions)-1]; it.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != "agent-a" {
		t.Errorf("S-0001 is %s, last moved to %s by %s", it.Status, last.To, last.By)
	}

	out, errOut, code = runIn(t, root, "stream", "open", "S-0002", "--json")
	if code != 0 {
		t.Fatalf("stream open --json: %d\n%s\n%s", code, out, errOut)
	}
	if !streamOpenMoved(t, out) {
		t.Errorf("moved should be true: %s", out)
	}
	if it := storyStartGet(t, root, "S-0002"); it.Status != workitem.InProgress {
		t.Errorf("S-0002 is %s", it.Status)
	}
}

// ADR-0153: a ready story the board holds is refused as flai story start
// refuses it, with exit 4 and the hold's reason, and nothing is written: no
// move, no narrative, no worktree.
func TestStreamOpenRefusesAHeldReadyStory(t *testing.T) {
	root := storyStartProject(t)
	before := storyStartGet(t, root, "S-0004")
	out, errOut, code := runIn(t, root, "stream", "open", "S-0004")
	if code != exitDocRefused {
		t.Fatalf("want exit %d, got %d\n%s\n%s", exitDocRefused, code, out, errOut)
	}
	if reason := fatalErr(t, errOut); !strings.Contains(reason, "starts when S-0001 is done") || !strings.Contains(reason, "nothing was changed") {
		t.Errorf("reason: %q\n%s", reason, errOut)
	}
	storyStartUnchanged(t, root, before)
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "worktrees", "S-0004")); err == nil {
		t.Error("S-0004 has a worktree")
	}
}

// ADR-0153: a story in any state but ready is opened as before, unmoved:
// one in progress already and one in backlog; --json gives moved false.
func TestStreamOpenLeavesAStoryThatIsNotReadyUnmoved(t *testing.T) {
	root := storyStartProject(t)
	if _, errOut, code := runIn(t, root, "move", "S-0002", workitem.InProgress, "--by", "alex"); code != 0 {
		t.Fatalf("move S-0002: %s", errOut)
	}
	for _, id := range []string{"S-0002", "S-0003"} {
		before := storyStartGet(t, root, id)
		out, errOut, code := runIn(t, root, "stream", "open", id, "--json")
		if code != 0 {
			t.Fatalf("stream open %s --json: %d\n%s\n%s", id, code, out, errOut)
		}
		if streamOpenMoved(t, out) {
			t.Errorf("%s: moved should be false: %s", id, out)
		}
		after := storyStartGet(t, root, id)
		if after.Status != before.Status || len(after.Transitions) != len(before.Transitions) {
			t.Errorf("%s changed: %s with %d transitions, was %s with %d", id, after.Status, len(after.Transitions), before.Status, len(before.Transitions))
		}
		if _, err := os.Stat(filepath.Join(root, "wip", "agents", id+".md")); err != nil {
			t.Errorf("%s has no narrative: %v", id, err)
		}
	}
}

// A story sent back to ready keeps its narrative: stream open moves it to in
// progress again and says to take it up where it is, rather than to run
// stream open again.
func TestStreamOpenOfAStorySentBackSaysToTakeItUp(t *testing.T) {
	root := storyStartProject(t)
	if _, errOut, code := runIn(t, root, "stream", "open", "S-0001"); code != 0 {
		t.Fatalf("stream open: %s", errOut)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", workitem.Ready, "--reason", "sent back"); code != 0 {
		t.Fatalf("move S-0001 back: %s", errOut)
	}
	_, errOut, code := runIn(t, root, "stream", "open", "S-0001")
	if code == 0 || !strings.Contains(errOut, "take it up in") {
		t.Errorf("want a refusal saying to take it up, got %d\n%s", code, errOut)
	}
	if it := storyStartGet(t, root, "S-0001"); it.Status != workitem.InProgress {
		t.Errorf("S-0001 is %s", it.Status)
	}
}
