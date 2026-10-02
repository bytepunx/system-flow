package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/release"
)

// ADR-0067: a clone whose remote branch has commits it lacks offers no plan
// and refuses to publish, changing nothing, naming the fetch and the merge
// or rebase; once fetched and merged, the publish goes through and pushes.
func TestReleasePendingFromACloneBehindItsRemoteBranch(t *testing.T) {
	defer func(d time.Duration) { release.RemoteTTL = d }(release.RemoteTTL)
	release.RemoteTTL = 0
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	// someone else pushed: the remote has a commit this clone lacks
	other := filepath.Join(t.TempDir(), "other")
	gitIn(t, filepath.Dir(other), "clone", "-q", remote, other)
	gitIn(t, other, "commit", "-q", "--allow-empty", "-m", "docs: from elsewhere")
	gitIn(t, other, "push", "-q", "origin", "main")
	theirs := gitIn(t, remote, "rev-parse", "main")
	head := gitIn(t, root, "rev-parse", "HEAD")

	dry, _, code := runIn(t, root, "release", "--pending", "--dry-run")
	if code != 0 || !strings.Contains(dry, "nothing can publish") || !strings.Contains(dry, "origin/main has commits this clone lacks") ||
		!strings.Contains(dry, "git fetch origin && git merge origin/main") || strings.Contains(dry, "S-0001") {
		t.Fatalf("dry run from a clone behind its remote branch: %d %s", code, dry)
	}
	js, _, _ := runIn(t, root, "release", "--pending", "--dry-run", "--json")
	var pub preview.Publishing
	if err := json.Unmarshal([]byte(js), &pub); err != nil || len(pub.Plans) != 0 || !pub.Remote.Lagging() || pub.Remote.Branch == nil ||
		*pub.Remote.Branch != (release.BranchLag{Upstream: "origin/main", Head: theirs}) || pub.Remote.Fix != "git fetch origin && git merge origin/main" {
		t.Errorf("dry run --json from a clone behind its remote branch: %v %s", err, js)
	}

	_, errOut, code := runIn(t, root, "release", "--pending")
	if code != exitPushDiverged || !strings.Contains(errOut, "refusing to publish") || !strings.Contains(errOut, "git fetch origin && git merge origin/main") ||
		!strings.Contains(errOut, "rebase onto origin/main") {
		t.Fatalf("publish from a clone behind its remote branch: %d %s", code, errOut)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.0.1") {
		t.Errorf("a refused publish tags nothing: %s", tags)
	}
	if now := gitIn(t, root, "rev-parse", "HEAD"); now != head {
		t.Errorf("a refused publish commits nothing: %s, was %s", now, head)
	}
	if status := gitIn(t, root, "status", "--porcelain"); strings.TrimSpace(status) != "" {
		t.Errorf("a refused publish leaves the checkout as it was: %s", status)
	}

	// fetched but not merged: still refused, with only the merge left
	gitIn(t, root, "fetch", "-q", "origin")
	if _, errOut, code := runIn(t, root, "release", "--pending"); code != exitPushDiverged || !strings.Contains(errOut, "fetched here but not merged") || !strings.Contains(errOut, "git merge origin/main") {
		t.Fatalf("publish with the remote branch fetched, not merged: %d %s", code, errOut)
	}

	gitIn(t, root, "merge", "-q", "--no-edit", "origin/main")
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.0.1") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("publish once merged: %d %s %s", code, out, errOut)
	}
	if remoteTags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(remoteTags, "cli/v1.0.1") {
		t.Errorf("the tag reaches the remote: %s", remoteTags)
	}
	if local, now := gitIn(t, root, "rev-parse", "HEAD"), gitIn(t, remote, "rev-parse", "main"); local != now {
		t.Errorf("main reaches the remote too: %s vs %s", local, now)
	}
}
