package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// researchProject is a git project with one code component (cli, tagged
// cli/v1.0.0), a bare remote named origin, and a story of the given nature in
// review whose branch holds a finding in design/ and, when withCode, a change
// to the component. It returns the project root and the remote's path.
func researchProject(t *testing.T, nature string, withCode bool) (root, remote string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root = t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design", "cli"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	remote = filepath.Join(t.TempDir(), "origin.git")
	gitIn(t, filepath.Dir(remote), "init", "-q", "--bare", "-b", "main", remote)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "remote", "add", "origin", remote)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	gitIn(t, root, "tag", "-a", "cli/v1.0.0", "-m", "cli 1.0.0")
	gitIn(t, root, "push", "-q", "origin", "main", "cli/v1.0.0")
	for _, args := range [][]string{
		{"epic", "new", "Epic"},
		{"story", "new", "What we found", "--epic", "E-0001", "--nature", nature},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-what-we-found.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] found\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Look", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(wt, "design", "finding.md"), []byte("# Finding\n"), 0o644)
	if withCode {
		_ = os.WriteFile(filepath.Join(wt, "cli", "probe.go"), []byte("package cli\n"), 0o644)
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] the finding")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root, remote
}

// assertResearchAccepted checks what ADR-0025 promises of an accepted
// research story: merged, archived, committed, pushed, and nothing released.
func assertResearchAccepted(t *testing.T, root, remote string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, "design", "finding.md")); err != nil {
		t.Errorf("the finding is on main: %v", err)
	}
	if m, _ := filepath.Glob(filepath.Join(root, "wip/archive/kanban/stories/S-0001-*.md")); len(m) != 1 {
		t.Errorf("the story is archived: %v", m)
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "worktrees", "S-0001")); err == nil {
		t.Error("the story worktree is removed")
	}
	if got := strings.TrimSpace(gitIn(t, root, "tag", "--list")); got != "cli/v1.0.0" {
		t.Errorf("no tag is created for research, got %q", got)
	}
	subject := strings.TrimSpace(gitIn(t, root, "log", "-1", "--format=%s"))
	if subject != "chore: [S-0001] accept and archive" {
		t.Errorf("the acceptance commit names no release, got %q", subject)
	}
	if status := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); status != "" {
		t.Errorf("everything is committed, got %q", status)
	}
	local := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD"))
	pushed := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main"))
	if local != pushed {
		t.Errorf("the acceptance commit is pushed although nothing was released: local %s, remote %s", local, pushed)
	}
	if got := strings.TrimSpace(gitIn(t, remote, "tag", "--list")); got != "cli/v1.0.0" {
		t.Errorf("no tag reaches the remote, got %q", got)
	}
}

// S-0053, ADR-0025: a research story is accepted and pushed with no release.
func TestAcceptResearchEndToEnd(t *testing.T) {
	root, remote := researchProject(t, "research", false)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the preview accepts research: %d %s", code, errOut)
	}
	var pre struct {
		Blockers []string `json:"blockers"`
		Plan     struct {
			Level      string `json:"level"`
			Skipped    string `json:"skipped"`
			Steps      []any  `json:"steps"`
			Unreleased []any  `json:"unreleased"`
		} `json:"plan"`
	}
	if err := json.Unmarshal([]byte(out), &pre); err != nil {
		t.Fatalf("preview json: %v\n%s", err, out)
	}
	if len(pre.Blockers) != 0 || pre.Plan.Level != "none" || !strings.Contains(pre.Plan.Skipped, "is research") || len(pre.Plan.Steps) != 0 || len(pre.Plan.Unreleased) != 0 {
		t.Errorf("preview: %+v", pre)
	}

	out, errOut, code = runIn(t, root, "accept", "S-0001")
	if code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	if !strings.Contains(out, "no release (S-0001 is research") {
		t.Errorf("the plan says research releases nothing and why:\n%s", out)
	}
	assertResearchAccepted(t, root, remote)
}

// The board's way in: flai move <story> done is the same acceptance. This
// story also changed a component, which lands unreleased and is said so.
func TestMoveResearchToDoneNamesCodeLandingUnreleased(t *testing.T) {
	root, remote := researchProject(t, "research", true)
	out, errOut, code := runIn(t, root, "move", "S-0001", "done")
	if code != 0 {
		t.Fatalf("move to done: %d %s", code, errOut)
	}
	if !strings.Contains(out, "cli") || !strings.Contains(out, "lands on main without a release") {
		t.Errorf("the plan names the component landing unreleased:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "cli", "probe.go")); err != nil {
		t.Errorf("the component change is on main: %v", err)
	}
	assertResearchAccepted(t, root, remote)
}

// An experiment stays on its branch: refused before anything is merged, in
// the preview as a blocker and in the real run as an error.
func TestAcceptRefusesAnExperimentBeforeMerging(t *testing.T) {
	root, remote := researchProject(t, "experiment", false)
	before := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD"))

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the preview succeeds and carries the blocker: %d %s", code, errOut)
	}
	var pre struct {
		Blockers []string `json:"blockers"`
	}
	_ = json.Unmarshal([]byte(out), &pre)
	if len(pre.Blockers) != 1 || !strings.Contains(pre.Blockers[0], "S-0001 is an experiment") {
		t.Errorf("blockers: %q", pre.Blockers)
	}

	_, errOut, code = runIn(t, root, "accept", "S-0001")
	if code == 0 || !strings.Contains(errOut, "stays on its branch") {
		t.Errorf("accept must refuse: %d %s", code, errOut)
	}
	if after := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")); after != before {
		t.Errorf("nothing is merged: HEAD moved from %s to %s", before, after)
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "worktrees", "S-0001")); err != nil {
		t.Errorf("the worktree is still there: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "design", "finding.md")); err == nil {
		t.Error("the experiment's files must not be on main")
	}
	if got := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")); got != before {
		t.Errorf("nothing is pushed: remote main %s", got)
	}
}

// Nothing pinned that the push follows an acceptance that cuts no release:
// a feature story that touched no component.
func TestAcceptPushesWhenThePlanIsSkipped(t *testing.T) {
	root, remote := researchProject(t, "feature", false)
	out, errOut, code := runIn(t, root, "accept", "S-0001")
	if code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	if !strings.Contains(out, "no release (no component touched") {
		t.Errorf("plan:\n%s", out)
	}
	assertResearchAccepted(t, root, remote)
}
