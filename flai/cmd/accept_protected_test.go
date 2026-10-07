package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0106: a story whose branch changes a path Claude Code protects is
// accepted by its operator only. The orchestrator under accept_reviews and an
// agent on its own name are refused before anything is merged, each refusal
// naming the files.
func TestOnlyTheOperatorAcceptsAStoryThatChangesProtectedPaths(t *testing.T) {
	root, head := protectedStoryInReview(t, "")
	mainBefore := gitIn(t, root, "rev-parse", "main")
	files := ".claude/settings.json, .mcp.json"

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", evidenceFile(t, "Verdict: pass\n- 1: `cli/feature.go`\n"))
	if code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet") || !strings.Contains(errOut, "story/S-0001 changes paths Claude Code protects, "+files+", so only the operator (alex) accepts S-0001, whatever the orchestrator judges (ADR-0106)") {
		t.Errorf("the orchestrator's acceptance must be refused naming the files: %d\n%s", code, errOut)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if len(res.OrchestratorBlockers) != 1 || res.OrchestratorBlockers[0].Code != preview.BlockProtectedPaths {
		t.Errorf("the protected files must be the orchestrator's one blocker: %+v", res.OrchestratorBlockers)
	}

	want := "only the operator (alex) accepts S-0001, not agent-S-0001: its branch changes paths Claude Code protects, " + files + "; the operator reads them in the diff and accepts it (ADR-0106)"
	for _, args := range [][]string{
		{"accept", "S-0001", "--by", "agent-S-0001"},
		{"move", "S-0001", "done", "--by", "agent-S-0001"},
	} {
		_, errOut, code := runIn(t, root, args...)
		if code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("flai %v by an agent on its own name must be refused naming the files: %d\n%s", args, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)

	text, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run")
	if code != 0 || !strings.Contains(text, "only the operator (alex) accepts S-0001: its branch changes paths Claude Code protects: "+files+"\n") || strings.Contains(text, "blocked:") {
		t.Errorf("the operator's dry run must list the files, unblocked: %d %s\n%s", code, errOut, text)
	}
	out, _, _ = runIn(t, root, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	res = preview.Acceptance{}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if !slices.Equal(res.Protected, []string{".claude/settings.json", ".mcp.json"}) || res.OperatorOnly != "only the operator (alex) accepts S-0001: its branch changes paths Claude Code protects" || len(res.Blockers) != 0 {
		t.Errorf("--dry-run --json must carry the files and the sentence: %s", out)
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("the operator's acceptance: %d %s\n%s", code, errOut, out)
	}
	res = preview.Acceptance{}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Status != workitem.Done || !res.Merged || len(res.Protected) != 2 {
		t.Errorf("the operator accepts it, the files still shown: %+v", res)
	}
	if !strings.Contains(gitIn(t, root, "show", "main:.mcp.json"), "mcpServers") {
		t.Error("the story branch is merged")
	}
}

// The operator is the story's owner or the project's owner; with neither
// named, anyone but the orchestrator.
func TestTheOperatorIsTheStorysOwnerOrTheProjects(t *testing.T) {
	for _, c := range []struct {
		name, projectOwner, storyOwner, by string
	}{
		{"the project's owner", "olive", "alex", "olive"},
		{"anyone, with no owner named", "", "", "agent-S-0001"},
	} {
		t.Run(c.name, func(t *testing.T) {
			root, _ := protectedStoryInReview(t, c.projectOwner)
			setOwner(t, root, c.storyOwner)
			if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", c.by); code != 0 {
				t.Errorf("%s accepts it: %d\n%s", c.by, code, errOut)
			}
		})
	}
}

// protectedStoryInReview is orchestratedStoryInReview whose story, owned by
// alex in a project owned by projectOwner when one is given, also commits
// .claude/settings.json and .mcp.json on its branch, both under its touches;
// it returns the main checkout and the branch's head.
func protectedStoryInReview(t *testing.T, projectOwner string) (root, head string) {
	t.Helper()
	root, _ = orchestratedStoryInReview(t, true)
	if projectOwner != "" {
		m := filepath.Join(root, "system-flow.yaml")
		b, _ := os.ReadFile(m)
		_ = os.WriteFile(m, append(b, []byte("owner: "+projectOwner+"\n")...), 0o644)
		gitIn(t, root, "commit", "-q", "-m", "chore: owner", "system-flow.yaml")
	}
	setOwner(t, root, "alex")
	if _, errOut, code := runIn(t, root, "touches", "S-0001", "--add", ".claude", ".mcp.json"); code != 0 {
		t.Fatal(errOut)
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.MkdirAll(filepath.Join(wt, ".claude"), 0o755)
	_ = os.WriteFile(filepath.Join(wt, ".claude", "settings.json"), []byte("{}\n"), 0o644)
	_ = os.WriteFile(filepath.Join(wt, ".mcp.json"), []byte("{\"mcpServers\": {}}\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] configure the agent")
	return root, gitIn(t, root, "rev-parse", "story/S-0001")
}

// setOwner makes owner the owner of S-0001, none when it is empty.
func setOwner(t *testing.T, root, owner string) {
	t.Helper()
	story := filepath.Join(root, "wip/kanban/stories/S-0001-ship-it.md")
	b, err := os.ReadFile(story)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(story, ownerLine.ReplaceAll(b, []byte("owner: "+owner)), 0o644)
}

// ownerLine is an item's owner line in its front matter.
var ownerLine = regexp.MustCompile(`(?m)^owner: .*$`)
