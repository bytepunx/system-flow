package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// heldProject has S-0001 in progress touching flai/cmd, and S-0002 ready
// touching cli, which the manifest names as flai's tag: S-0002 is held.
func heldProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	data, _ := os.ReadFile(manifest)
	_ = os.WriteFile(manifest, append(data, "projects:\n  - name: flai\n    path: flai\n    kind: go\n    tags: [cli]\n"...), 0o644)
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	for _, s := range []struct{ title, touches string }{{"Open", "flai/cmd"}, {"Waits", "cli"}} {
		if _, errOut, code := runIn(t, root, "story", "new", s.title, "--epic", "E-0001", "--touches", s.touches); code != 0 {
			t.Fatal(errOut)
		}
	}
	for _, id := range []string{"S-0001", "S-0002"} {
		matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		body, _ := os.ReadFile(matches[0])
		_ = os.WriteFile(matches[0], []byte(strings.Replace(string(body), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] done\n", 1)), 0o644)
		if _, errOut, code := runIn(t, root, "move", id, "ready"); code != 0 {
			t.Fatal(errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "in-progress"); code != 0 {
		t.Fatal(errOut)
	}
	return root
}

const heldWhy = "held (overlap): touches flai, which holds flai/cmd that S-0001 (in progress) touches; starts when S-0001 is accepted, cancelled, or sent back"

// S-0128: moving a held story to in-progress warns, and moves it.
func TestMovingAHeldStoryWarns(t *testing.T) {
	root := heldProject(t)
	out, errOut, code := runIn(t, root, "move", "S-0002", "in-progress")
	if code != 0 || !strings.Contains(out, "S-0002 → in-progress") {
		t.Fatalf("move: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "workflow policy warning") || !strings.Contains(errOut, "S-0002 is "+heldWhy) {
		t.Errorf("no warning: %s", errOut)
	}
}

// S-0128: flai board marks a held story and says why, as text and as JSON.
func TestBoardMarksAHeldStory(t *testing.T) {
	root := heldProject(t)
	out, errOut, code := runIn(t, root, "board")
	if code != 0 {
		t.Fatal(errOut)
	}
	if !strings.Contains(out, "HELD\n") || !strings.Contains(out, "         "+heldWhy+"\n") {
		t.Errorf("board text:\n%s", out)
	}
	out, errOut, code = runIn(t, root, "board", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var view struct {
		Columns map[string][]struct {
			ID   string `json:"id"`
			Held *struct {
				Code   string `json:"code"`
				Reason string `json:"reason"`
			} `json:"held"`
		} `json:"columns"`
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatal(err)
	}
	ready := view.Columns["ready"]
	if len(ready) != 1 || ready[0].Held == nil || ready[0].Held.Code != "overlap" || ready[0].Held.Reason != heldWhy {
		t.Errorf("ready: %+v", ready)
	}
	if ip := view.Columns["in-progress"]; len(ip) != 1 || ip[0].Held != nil {
		t.Errorf("in progress: %+v", ip)
	}
}

// S-0130: flai edit --after makes a ready story wait for another, flai board
// says so, flai check refuses a story that does not exist and a cycle, and
// --clear-after lets it go.
func TestEditAfterHoldsAStoryUntilTheOtherIsDone(t *testing.T) {
	root := heldProject(t)
	if _, errOut, code := runIn(t, root, "story", "new", "Later", "--epic", "E-0001", "--touches", "docs"); code != 0 {
		t.Fatal(errOut)
	}
	matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", "S-0003-*.md"))
	body, _ := os.ReadFile(matches[0])
	_ = os.WriteFile(matches[0], []byte(strings.Replace(string(body), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] done\n", 1)), 0o644)
	if _, errOut, code := runIn(t, root, "move", "S-0003", "ready"); code != 0 {
		t.Fatal(errOut)
	}

	out, errOut, code := runIn(t, root, "edit", "S-0003", "--after", "S-2")
	if code != 0 || !strings.Contains(out, "S-0003: changed after") {
		t.Fatalf("edit: %d %s %s", code, out, errOut)
	}
	if item, _ := os.ReadFile(matches[0]); !strings.Contains(string(item), "\nafter: [S-0002]\n") {
		t.Errorf("the item:\n%s", item)
	}
	if out, _, _ := runIn(t, root, "edit", "S-0003", "--show"); !strings.Contains(out, "  after: S-0002\n") {
		t.Errorf("show:\n%s", out)
	}
	if out, _, _ := runIn(t, root, "board"); !strings.Contains(out, "held (after): waits for S-0002 (ready); starts when S-0002 is done\n") {
		t.Errorf("board:\n%s", out)
	}

	for _, c := range []struct {
		args []string
		code int
		says string
	}{
		{[]string{"edit", "S-0003", "--after", "S-0009"}, 4, "after names S-0009, which does not exist"},
		{[]string{"edit", "S-0002", "--after", "S-0003"}, 4, "after forms a cycle, S-0002 waits for S-0003 waits for S-0002"},
		{[]string{"edit", "S-0003", "--after", "S-3"}, 1, "a story cannot wait for itself"},
		{[]string{"edit", "S-0003", "--after", "E-0001"}, 1, "name stories, such as S-0001"},
		{[]string{"edit", "S-0003", "--after", "S-0001", "--clear-after"}, 1, "contradict"},
	} {
		out, errOut, code := runIn(t, root, c.args...)
		if code != c.code || !strings.Contains(out+errOut, c.says) {
			t.Errorf("%v: %d\n%s%s", c.args, code, out, errOut)
		}
	}

	if out, errOut, code := runIn(t, root, "edit", "S-0003", "--clear-after"); code != 0 || !strings.Contains(out, "changed after") {
		t.Fatalf("clear: %d %s %s", code, out, errOut)
	}
	if item, _ := os.ReadFile(matches[0]); strings.Contains(string(item), "after:") {
		t.Errorf("after: stays:\n%s", item)
	}
	if out, _, _ := runIn(t, root, "board"); strings.Contains(out, "held (after)") {
		t.Errorf("still held:\n%s", out)
	}
}
