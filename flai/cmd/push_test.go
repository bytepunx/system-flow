package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// S-0063: an acceptance made with no credential at hand is pushed later, from
// the host, by flai push --pending; flai board says when there is one.
func TestPushPending(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	// accept as the dashboard container does when it holds nothing
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--no-push")
	if code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	if !strings.Contains(out, "cli/v1.1.0") {
		t.Fatalf("the acceptance tags a release: %s", out)
	}
	remoteHead := func() string { return strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")) }
	before := remoteHead()

	board, _, _ := runIn(t, root, "board")
	if !strings.Contains(board, "accepted, not pushed: S-0001") || !strings.Contains(board, "tags cli/v1.1.0") || !strings.Contains(board, "flai push --pending") {
		t.Errorf("flai board says what is pending and what to run:\n%s", board)
	}
	var view struct {
		Unpushed *struct {
			Commits     int      `json:"commits"`
			Acceptances []string `json:"acceptances"`
			Tags        []string `json:"tags"`
			Upstream    string   `json:"upstream"`
		} `json:"unpushed"`
	}
	js, _, _ := runIn(t, root, "board", "--json")
	if err := json.Unmarshal([]byte(js), &view); err != nil || view.Unpushed == nil || view.Unpushed.Upstream != "origin/main" || len(view.Unpushed.Acceptances) != 1 || view.Unpushed.Acceptances[0] != "S-0001" || len(view.Unpushed.Tags) != 1 {
		t.Errorf("board --json: %v %+v", err, view.Unpushed)
	}

	if _, errOut, code := runIn(t, root, "push"); code == 0 || !strings.Contains(errOut, "--pending") {
		t.Errorf("flai push alone does nothing and says why: %d %s", code, errOut)
	}
	dry, _, _ := runIn(t, root, "push", "--pending", "--dry-run")
	if !strings.Contains(dry, "would push S-0001 to origin: main, tags cli/v1.1.0") || remoteHead() != before {
		t.Errorf("dry run: %s (remote moved: %v)", dry, remoteHead() != before)
	}
	out, errOut, code = runIn(t, root, "push", "--pending")
	if code != 0 || !strings.Contains(out, "pushed S-0001 to origin: main, tags cli/v1.1.0") {
		t.Fatalf("push: %d %s %s", code, out, errOut)
	}
	if got := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")); remoteHead() != got {
		t.Errorf("the remote has the acceptance: %s vs %s", remoteHead(), got)
	}
	if tags := gitIn(t, remote, "tag", "--list"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the release tag is on the remote: %s", tags)
	}
	again, _, _ := runIn(t, root, "push", "--pending")
	if !strings.Contains(again, "nothing pending") {
		t.Errorf("a second run: %s", again)
	}
	board, _, _ = runIn(t, root, "board", "--json")
	if strings.Contains(board, "unpushed") {
		t.Errorf("the board stops saying so at once:\n%s", board)
	}
}

func TestPushPendingLeavesOrdinaryCommitsAndDivergenceAlone(t *testing.T) {
	root, remote := researchProject(t, "feature", false)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: wip that is nobody's acceptance")
	out, _, code := runIn(t, root, "push", "--pending")
	if code != 0 || !strings.Contains(out, "none of them an acceptance") {
		t.Errorf("ordinary commits are yours to push: %d %s", code, out)
	}
	if head := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")); head == strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")) {
		t.Error("nothing should have been pushed")
	}

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--no-push", "--yes"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	// someone pushes to the remote from another clone, and this one fetches
	other := filepath.Join(t.TempDir(), "other")
	gitIn(t, filepath.Dir(other), "clone", "-q", remote, other)
	gitIn(t, other, "config", "user.email", "o@o")
	gitIn(t, other, "config", "user.name", "o")
	gitIn(t, other, "commit", "-q", "--allow-empty", "-m", "docs: from elsewhere")
	gitIn(t, other, "push", "-q", "origin", "main")
	gitIn(t, root, "fetch", "-q", "origin")
	_, errOut, code := runIn(t, root, "push", "--pending")
	if code == 0 || !strings.Contains(errOut, "have diverged") || !strings.Contains(errOut, "never forces") {
		t.Errorf("diverged: refuse, say to fetch and merge, never force: %d %s", code, errOut)
	}
}
