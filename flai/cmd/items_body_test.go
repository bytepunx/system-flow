package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// bodyProject is a git project with one epic, committed, and a clean tree.
func bodyProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "olive@example.invalid")
	gitIn(t, root, "config", "user.name", "Olive")
	if _, errOut, code := runIn(t, root, "epic", "new", "Workbench"); code != 0 {
		t.Fatal(errOut)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	return root
}

const storyBody = "## Goal\nThe operator creates work from the board.\n\n## Acceptance criteria\n- [ ] A form exists\n- [ ] flai does the rest\n\n## Tasks\n\n## Notes\nRaised on the board.\n"

// S-0059: the author's markdown becomes the item's body, and everything else
// is what flai story new would have written.
func TestStoryNewWithABodyFromStandardInput(t *testing.T) {
	root := bodyProject(t)
	out, errOut, code := runStdin(t, root, storyBody, "story", "new", "Created from the board", "--epic", "E-0001", "--nature", "improvement", "--owner", "olive", "--tag", "dashboard", "--body-stdin", "--autocommit", "--trailer", "Created-with: flaiover", "--json")
	if code != 0 {
		t.Fatalf("create: %d %s", code, errOut)
	}
	var res struct {
		Item struct {
			ID, Title, Nature, Parent, Owner, Status string
		} `json:"item"`
		Path      string `json:"path"`
		Committed bool   `json:"committed"`
		Commit    string `json:"commit"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if res.Item.ID != "S-0001" || res.Item.Nature != "improvement" || res.Item.Parent != "E-0001" || res.Item.Owner != "olive" || res.Item.Status != "backlog" {
		t.Errorf("item: %+v", res.Item)
	}
	if res.Path != "wip/kanban/stories/S-0001-created-from-the-board.md" || !res.Committed || res.Commit == "" {
		t.Errorf("result: %+v", res)
	}
	file, _ := os.ReadFile(filepath.Join(root, res.Path))
	text := string(file)
	if !strings.HasPrefix(text, "---\nid: S-0001\n") || !strings.Contains(text, "tags: [dashboard]") {
		t.Errorf("front matter is flai's:\n%s", text)
	}
	if !strings.Contains(text, "\n# S-0001 Created from the board\n\n## Goal\nThe operator creates work from the board.\n") || !strings.Contains(text, "- [ ] flai does the rest\n") || !strings.HasSuffix(text, "Raised on the board.\n") {
		t.Errorf("the heading is flai's and the body is the author's:\n%s", text)
	}
	epic, _ := filepath.Glob(filepath.Join(root, "wip/kanban/epics/E-0001-*.md"))
	if data, _ := os.ReadFile(epic[0]); !strings.Contains(string(data), "- S-0001 Created from the board") {
		t.Errorf("the story is linked into its epic:\n%s", data)
	}
	// one commit, holding the story and its epic and nothing else, authored by the git identity, with the trailer
	show := gitIn(t, root, "show", "--stat", "--format=%an|%s|%b", "HEAD")
	if !strings.HasPrefix(show, "Olive|chore: [S-0001] create story: Created from the board|Created-with: flaiover") {
		t.Errorf("commit: %s", show)
	}
	if !strings.Contains(show, "S-0001-created-from-the-board.md") || !strings.Contains(show, "E-0001-workbench.md") || !strings.Contains(show, "2 files changed") {
		t.Errorf("the commit holds the two paths only: %s", show)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("nothing left uncommitted: %q", st)
	}
	// it can go to ready with no tasks (ADR-0021)
	if _, errOut, code := runIn(t, root, "move", "S-0001", "ready"); code != 0 {
		t.Errorf("a story written this way is ready: %s", errOut)
	}
}

func TestEpicNewWithABodyAndAutocommitOff(t *testing.T) {
	root := bodyProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	m, _ := os.ReadFile(manifest)
	_ = os.WriteFile(manifest, append(m, []byte("dashboard:\n  autocommit: false\n")...), 0o644)
	gitIn(t, root, "commit", "-q", "-am", "autocommit off")
	head := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD"))
	out, errOut, code := runStdin(t, root, "## Outcome\nA second epic.\n\n## Stories\n\n## Notes\n", "epic", "new", "Second", "--nature", "research", "--body-stdin", "--autocommit")
	if code != 0 || !strings.HasPrefix(out, "E-0002 Second") || strings.Contains(out, "committed") {
		t.Fatalf("create: %d %s %s", code, out, errOut)
	}
	if got := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")); got != head {
		t.Error("dashboard.autocommit: false means no commit")
	}
	file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/epics/E-0002-*.md"))
	if data, _ := os.ReadFile(file[0]); !strings.Contains(string(data), "# E-0002 Second\n\n## Outcome\nA second epic.\n") || !strings.Contains(string(data), "nature: research") {
		t.Errorf("epic:\n%s", data)
	}
}

// What flai check would report with the new item in place refuses it, and
// nothing is left behind: no file, no link in the parent, no commit.
func TestItemNewIsRefusedByTheCheckAndLeavesNothingBehind(t *testing.T) {
	root := bodyProject(t)
	epic, _ := filepath.Glob(filepath.Join(root, "wip/kanban/epics/E-0001-*.md"))
	epicWas, _ := os.ReadFile(epic[0])
	head := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD"))
	// a story without its Notes section is a finding (item.heading) the item introduces
	body := "## Goal\nNo notes section.\n\n## Acceptance criteria\n- [ ] x\n\n## Tasks\n"
	out, errOut, code := runStdin(t, root, body, "story", "new", "Refused", "--epic", "E-0001", "--body-stdin", "--autocommit", "--json")
	if code != 4 {
		t.Fatalf("exit 4, as a refused save: %d\n%s\n%s", code, out, errOut)
	}
	var res struct {
		Refused struct {
			Reason   string `json:"reason"`
			Findings []struct {
				Rule, Path, Message string
			} `json:"findings"`
		} `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || len(res.Refused.Findings) != 1 || res.Refused.Findings[0].Rule != "item.heading" || !strings.Contains(res.Refused.Reason, "nothing was created") {
		t.Fatalf("the findings come back: %v\n%s", err, out)
	}
	if m, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories/*.md")); len(m) != 0 {
		t.Errorf("no story file is left: %v", m)
	}
	if now, _ := os.ReadFile(epic[0]); string(now) != string(epicWas) {
		t.Errorf("the epic is as it was:\n%s", now)
	}
	if got := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")); got != head {
		t.Error("nothing was committed")
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("a clean tree: %q", st)
	}
	// the ID was not spent: the next creation gets it
	out, _, code = runStdin(t, root, storyBody, "story", "new", "Accepted", "--epic", "E-0001", "--body-stdin")
	if code != 0 || !strings.HasPrefix(out, "S-0001 Accepted") {
		t.Errorf("next: %d %s", code, out)
	}
}

func TestItemNewPrintBodyAndItsRefusals(t *testing.T) {
	root := bodyProject(t)
	out, errOut, code := runIn(t, root, "story", "new", "--print-body")
	if code != 0 {
		t.Fatal(errOut)
	}
	if strings.Contains(out, "---") || strings.Contains(out, "# ") && strings.HasPrefix(out, "# ") || !strings.Contains(out, "## Goal") || !strings.Contains(out, "## Acceptance criteria") || !strings.Contains(out, "## Notes") {
		t.Errorf("the template's sections, no front matter, no heading:\n%s", out)
	}
	out, _, _ = runIn(t, root, "epic", "new", "--print-body")
	if !strings.Contains(out, "## Stories") {
		t.Errorf("an epic's body:\n%s", out)
	}
	if m, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories/*.md")); len(m) != 0 {
		t.Errorf("--print-body creates nothing: %v", m)
	}
	if _, errOut, code := runStdin(t, root, "  \n", "story", "new", "Empty", "--epic", "E-0001", "--body-stdin"); code == 0 || !strings.Contains(errOut, "standard input is empty") {
		t.Errorf("an empty body is refused: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "No parent"); code == 0 || !strings.Contains(errOut, "epic") {
		t.Errorf("a story still needs its epic: %d %s", code, errOut)
	}
}
