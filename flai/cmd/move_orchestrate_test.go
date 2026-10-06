package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0219: in the orchestrator's shell flai move takes a story to ready only
// when flai promote --candidates lists it and only while ready is under its
// WIP limit; a held story, a draft, even with --yes, and any other item are
// refused with the reason. Outside it a move to ready is as it was.
func TestTheOrchestratorMovesOnlyACandidateToReady(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "orchestrator")
	t.Setenv("FLAI_ROLE", "")
	root := tempProject(t)
	run := func(args ...string) string {
		t.Helper()
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		return out
	}
	fill := func(id string) {
		t.Helper()
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		body := strings.Replace(string(s), "## Goal\n", "## Goal\n\nDo it.\n", 1)
		body = strings.Replace(body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)
		if err := os.WriteFile(file[0], []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("epic", "new", "Epic")
	run("story", "new", "Open", "--epic", "E-0001", "--touches", "flai")
	run("story", "new", "Held", "--epic", "E-0001", "--touches", "flai/cmd")
	run("story", "new", "Drafted", "--epic", "E-0001", "--touches", "wip", "--draft")
	run("story", "new", "Candidate", "--epic", "E-0001", "--touches", "docs")
	run("story", "new", "Next", "--epic", "E-0001", "--touches", "scripts")
	run("story", "new", "Unplanned", "--epic", "E-0001", "--touches", "design")
	for _, id := range []string{"S-0001", "S-0002", "S-0003", "S-0004", "S-0005", "S-0006"} {
		fill(id)
	}
	for _, id := range []string{"S-0001", "S-0002", "S-0003", "S-0004", "S-0005"} {
		run("edit", id, "--cost-of-delay-value", "300", "--forecast-duration", "2h")
	}
	run("move", "S-0001", "ready")
	run("move", "S-0001", "in-progress")
	run("task", "new", "Piece", "--story", "S-0004")

	t.Setenv("FLAI_ROLE", "orchestrate")
	if _, errOut, code := runIn(t, root, "move", "S-0004", "ready"); code == 0 || !strings.Contains(errOut, "rule: the orchestrator moves a story to ready only with orchestration.permissions.promote_to_ready, which is off") {
		t.Errorf("a candidate without promote_to_ready: exit %d %s", code, errOut)
	}
	permitOrchestrator(t, root, "promote_to_ready")
	const notCandidate = " is not a candidate to go to ready (flai promote --candidates): "
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"move", "S-0002", "ready"}, "rule: S-0002" + notCandidate + "held (overlap)"},
		{[]string{"move", "S-0003", "ready"}, "rule: S-0003" + notCandidate + "draft: finalize it first"},
		{[]string{"move", "S-0003", "ready", "--yes"}, "rule: S-0003" + notCandidate + "draft: finalize it first"},
		{[]string{"move", "S-0006", "ready"}, "rule: S-0006" + notCandidate + "no forecast duration; no cost of delay value"},
		{[]string{"move", "T-0001", "ready"}, "rule: T-0001 is a task: the orchestrator moves a story to ready"},
	} {
		if _, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: exit %d %q, want %q", c.args, code, errOut, c.want)
		}
	}
	if out := run("move", "S-0004", "ready"); !strings.Contains(out, "S-0004 → ready") {
		t.Errorf("a candidate: %s", out)
	}
	run("board", "limit", "ready", "1")
	if _, errOut, code := runIn(t, root, "move", "S-0005", "ready"); code == 0 || !strings.Contains(errOut, "rule: the ready column is at its WIP limit (1 of 1): the orchestrator moves no story to ready until one leaves it") {
		t.Errorf("a candidate with ready full: exit %d %s", code, errOut)
	}
	if out := run("show", "S-0005", "--json"); !strings.Contains(out, `"status": "backlog"`) {
		t.Errorf("S-0005 moved: %s", out)
	}

	t.Setenv("FLAI_ROLE", "")
	if out := run("move", "S-0002", "ready"); !strings.Contains(out, "S-0002 → ready") {
		t.Errorf("a held story outside the orchestrator's shell, as before: %s", out)
	}
}

// permitOrchestrator gives the orchestrator the permissions on in the
// project at root's manifest, and no other.
func permitOrchestrator(t *testing.T, root string, on ...string) {
	t.Helper()
	manifest := filepath.Join(root, "system-flow.yaml")
	m, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	block := "orchestration:\n  permissions:\n"
	for _, p := range on {
		block += "    " + p + ": true\n"
	}
	if err := os.WriteFile(manifest, append(m, block...), 0o644); err != nil {
		t.Fatal(err)
	}
}
