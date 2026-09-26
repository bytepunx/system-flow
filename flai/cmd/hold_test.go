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
