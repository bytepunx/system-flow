package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0063: an acceptance made with no credential at hand is pushed later, from
// the host, by flai push --pending; flai board says when there is one.
func TestPushPending(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	// accept as the dashboard container does when it holds nothing; S-0087:
	// acceptance alone tags nothing, so there is a pending acceptance commit
	// but no tag until the operator publishes.
	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	remoteHead := func() string { return strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")) }
	before := remoteHead()

	board, _, _ := runIn(t, root, "board")
	if !strings.Contains(board, "accepted, not pushed: S-0001") || !strings.Contains(board, "flai push --pending") {
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
	if err := json.Unmarshal([]byte(js), &view); err != nil || view.Unpushed == nil || view.Unpushed.Upstream != "origin/main" || len(view.Unpushed.Acceptances) != 1 || view.Unpushed.Acceptances[0] != "S-0001" || len(view.Unpushed.Tags) != 0 {
		t.Errorf("board --json: %v %+v", err, view.Unpushed)
	}

	if _, errOut, code := runIn(t, root, "push"); code == 0 || !strings.Contains(errOut, "--pending") {
		t.Errorf("flai push alone does nothing and says why: %d %s", code, errOut)
	}
	dry, _, _ := runIn(t, root, "push", "--pending", "--dry-run")
	if !strings.Contains(dry, "would push S-0001 to origin: main") || remoteHead() != before {
		t.Errorf("dry run: %s (remote moved: %v)", dry, remoteHead() != before)
	}
	out, errOut, code := runIn(t, root, "push", "--pending")
	if code != 0 || !strings.Contains(out, "pushed S-0001 to origin: main") {
		t.Fatalf("push: %d %s %s", code, out, errOut)
	}
	if got := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")); remoteHead() != got {
		t.Errorf("the remote has the acceptance: %s vs %s", remoteHead(), got)
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

// S-0094: flai accept computes no release (S-0087), so nothing tags it until
// something does; this is the one place it happens, before the push itself,
// not a separate step left for someone to remember.
func TestPushPendingTagsWhatAcceptLeftUnreleased(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	if tags := gitIn(t, root, "tag", "--list"); strings.Contains(tags, "v1.1.0") {
		t.Fatalf("accept alone must tag nothing: %s", tags)
	}

	dry, _, _ := runIn(t, root, "push", "--pending", "--dry-run")
	if !strings.Contains(dry, "cli") || !strings.Contains(dry, "1.1.0") || !strings.Contains(dry, "would also tag") {
		t.Errorf("dry run previews the release it would tag first: %s", dry)
	}
	if tags := gitIn(t, root, "tag", "--list"); strings.Contains(tags, "v1.1.0") {
		t.Errorf("a dry run must tag nothing: %s", tags)
	}

	out, errOut, code := runIn(t, root, "push", "--pending", "--json")
	if code != 0 {
		t.Fatalf("push: %s", errOut)
	}
	var got struct {
		Pushed bool     `json:"pushed"`
		Tags   []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || !got.Pushed || len(got.Tags) != 1 || got.Tags[0] != "cli/v1.1.0" {
		t.Fatalf("push --json: %v %s", err, out)
	}
	if tags := gitIn(t, root, "tag", "--list"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the tag exists locally: %s", tags)
	}
	if remoteTags := gitIn(t, remote, "tag", "--list"); !strings.Contains(remoteTags, "cli/v1.1.0") {
		t.Errorf("and it reached the remote, with the push, not after it: %s", remoteTags)
	}
	if remoteHead := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")); remoteHead != strings.TrimSpace(gitIn(t, root, "rev-parse", "main")) {
		t.Error("the version-bump commit reached the remote too")
	}
	if again, _, _ := runIn(t, root, "push", "--pending"); !strings.Contains(again, "nothing pending") {
		t.Errorf("a second run retags nothing: %s", again)
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

	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--yes"); code != 0 {
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
	if code != exitPushDiverged || !strings.Contains(errOut, "have diverged") || !strings.Contains(errOut, "never forces") {
		t.Errorf("diverged: refuse with its own exit code, say to fetch and merge, never force: %d %s", code, errOut)
	}
}

// I-0026: more than three tags in one push create no tag events on GitHub, so
// the push is made in parts, and every tag still arrives.
func TestPushPendingManyTags(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	for _, tag := range []string{"x/v1", "x/v2", "x/v3", "x/v4", "x/v5", "x/v6"} {
		gitIn(t, root, "tag", "-a", tag, "-m", tag)
	}
	out, errOut, code := runIn(t, root, "push", "--pending")
	if code != 0 || !strings.Contains(out, "pushed S-0001 to origin") {
		t.Fatalf("push: %d %s %s", code, out, errOut)
	}
	there := gitIn(t, remote, "tag", "--list")
	for _, tag := range []string{"x/v1", "x/v6"} {
		if !strings.Contains(there, tag) {
			t.Errorf("%s did not arrive: %s", tag, there)
		}
	}
	if again, _, _ := runIn(t, root, "push", "--pending"); !strings.Contains(again, "nothing pending") {
		t.Errorf("after the push: %s", again)
	}
}

// S-0078: --publish also publishes a template component whose version the
// pushed commits moved, after the push, as flai template push --tag does.
func TestPushPendingPublishesAMovedTemplate(t *testing.T) {
	root, remote := researchProject(t, "feature", true)
	bare := filepath.Join(t.TempDir(), "tpl.git")
	gitIn(t, filepath.Dir(bare), "init", "-q", "--bare", "-b", "main", bare)
	tpl := filepath.Join(root, "tpl")
	_ = filepath.WalkDir(miniTemplate, func(path string, d os.DirEntry, _ error) error {
		rel, _ := filepath.Rel(miniTemplate, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(tpl, rel), 0o755)
		}
		b, _ := os.ReadFile(path)
		return os.WriteFile(filepath.Join(tpl, rel), b, 0o644)
	})
	manifest := filepath.Join(tpl, "template.yaml")
	y, _ := os.ReadFile(manifest)
	_ = os.WriteFile(manifest, []byte(string(y)+"publish:\n  repo: "+bare+"\n  ref: main\n"), 0o644)
	m, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(string(m)+"  - name: tpl\n    path: tpl\n    kind: template\n"), 0o644)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: a template component")
	gitIn(t, root, "push", "-q", "origin", "main")

	// nothing of the template moved: an acceptance alone publishes nothing
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--yes"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	dry, _, _ := runIn(t, root, "push", "--pending", "--publish", "--dry-run")
	if strings.Contains(dry, "publish") {
		t.Errorf("no template version moved: %s", dry)
	}
	// its version moves in a commit that is ahead too
	y, _ = os.ReadFile(manifest)
	_ = os.WriteFile(manifest, []byte(strings.Replace(string(y), "version: 9.9.9", "version: 9.9.10", 1)), 0o644)
	gitIn(t, root, "commit", "-q", "-am", "chore: template 9.9.10")
	dry, _, _ = runIn(t, root, "push", "--pending", "--publish", "--dry-run")
	if !strings.Contains(dry, "then publish tpl") || strings.Contains(gitIn(t, bare, "tag", "--list"), "v9.9.10") {
		t.Errorf("the dry run says what would be published and publishes nothing: %s", dry)
	}
	if plain, _, _ := runIn(t, root, "push", "--pending", "--dry-run"); strings.Contains(plain, "publish") {
		t.Errorf("without --publish nothing is said of templates: %s", plain)
	}
	out, errOut, code := runIn(t, root, "push", "--pending", "--publish", "--json")
	if code != 0 || !strings.Contains(out, `"pushed": true`) || !strings.Contains(out, "v9.9.10") {
		t.Fatalf("push --publish: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(gitIn(t, bare, "tag", "--list"), "v9.9.10") {
		t.Error("the template's remote has the new version's tag")
	}
	if strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")) != strings.TrimSpace(gitIn(t, root, "rev-parse", "main")) {
		t.Error("and the project's remote has main")
	}
	if again, _, _ := runIn(t, root, "push", "--pending", "--publish"); !strings.Contains(again, "nothing pending") {
		t.Errorf("twice is once: %s", again)
	}
}
