package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// moveRemoteDuringPush makes the remote's main move after the publish has
// checked it and before its push of the branch lands, as when someone pushes
// in between (S-0242): another clone's commit waits on the remote's branch
// moved, and a pre-push hook in the publishing clone, which git runs once it
// has the remote's refs and before it sends anything, sets main to it when
// the push carries a branch. It returns the hook's path, to remove it.
func moveRemoteDuringPush(t *testing.T, root, remote string) string {
	t.Helper()
	other := filepath.Join(t.TempDir(), "other")
	gitIn(t, filepath.Dir(other), "clone", "-q", remote, other)
	gitIn(t, other, "commit", "-q", "--allow-empty", "-m", "docs: from elsewhere")
	gitIn(t, other, "push", "-q", "origin", "HEAD:refs/heads/moved")
	hook := gitIn(t, root, "rev-parse", "--path-format=absolute", "--git-path", "hooks/pre-push")
	script := "#!/bin/sh\ngrep -q ' refs/heads/' || exit 0\ngit --git-dir='" + remote + "' update-ref refs/heads/main refs/heads/moved\n"
	if err := os.MkdirAll(filepath.Dir(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hook, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return hook
}

// S-0242, ADR-0067, TH-0070: a publish whose push is refused because the
// remote's branch moved after the check deletes the release tags it made,
// leaves none on the remote, and says to fetch and rebase or merge; once
// merged, publishing again tags again and pushes.
func TestReleasePendingWhenTheRemoteMovesDuringThePush(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	hook := moveRemoteDuringPush(t, root, remote)

	_, errOut, code := runIn(t, root, "release", "--pending")
	if code != exitPushDiverged {
		t.Fatalf("a push refused because the remote moved exits 3: %d %s", code, errOut)
	}
	for _, s := range []string{"the push was refused because origin/main moved", "deleted the release tags cli/v1.0.1 here", "git fetch origin", "rebase onto origin/main or merge it", "flai release --pending again, which tags again"} {
		if !strings.Contains(errOut, s) {
			t.Errorf("message lacks %q: %s", s, errOut)
		}
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/v1.0.1"); tags != "" {
		t.Errorf("the tag the push did not send is deleted here: %q", tags)
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/v1.0.1"); tags != "" {
		t.Errorf("the refused push leaves no tag on the remote: %q", tags)
	}

	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "fetch", "-q", "origin")
	gitIn(t, root, "merge", "-q", "--no-edit", "origin/main")
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.0.1") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("publish once merged: %d %s %s", code, out, errOut)
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/v1.0.1"); tags != "cli/v1.0.1" {
		t.Errorf("tagged again, the tag reaches the remote: %q", tags)
	}
	if local, now := gitIn(t, root, "rev-parse", "HEAD"), gitIn(t, remote, "rev-parse", "main"); local != now {
		t.Errorf("main reaches the remote too: %s vs %s", local, now)
	}
}

// S-0242: tags an earlier batch already pushed are kept and named, with the
// advice to merge rather than rebase; a tag of the user's own that did not
// go is not a release tag and is not deleted.
func TestReleasePendingKeepsTagsAlreadyPushedWhenTheRemoteMoves(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	// with the release tag, four tags ride the publish: three go first on
	// their own (I-0026), and the last goes with main
	for _, tag := range []string{"mine-1", "mine-2", "mine-3"} {
		gitIn(t, root, "tag", tag)
	}
	moveRemoteDuringPush(t, root, remote)

	_, errOut, code := runIn(t, root, "release", "--pending")
	if code != exitPushDiverged {
		t.Fatalf("a push refused because the remote moved exits 3: %d %s", code, errOut)
	}
	for _, s := range []string{"origin/main moved", "Tags cli/v1.0.1, mine-1, mine-2 already reached origin and are kept", "merge origin/main rather than rebase onto it"} {
		if !strings.Contains(errOut, s) {
			t.Errorf("message lacks %q: %s", s, errOut)
		}
	}
	if strings.Contains(errOut, "deleted") {
		t.Errorf("no release tag was left to delete: %s", errOut)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/v1.0.1", "mine-3"); tags != "cli/v1.0.1\nmine-3" {
		t.Errorf("the pushed release tag and the user's own are kept here: %q", tags)
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/v1.0.1", "mine-3"); tags != "cli/v1.0.1" {
		t.Errorf("only the earlier batch reached the remote: %q", tags)
	}
}

// S-0242: with --json, the refusal is the result with remote_moved and the
// tags deleted and kept, and exit 3.
func TestReleasePendingJSONWhenTheRemoteMovesDuringThePush(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	moveRemoteDuringPush(t, root, remote)

	out, errOut, code := runIn(t, root, "release", "--pending", "--json")
	if code != exitPushDiverged {
		t.Fatalf("a push refused because the remote moved exits 3: %d %s %s", code, out, errOut)
	}
	var result struct {
		Pushed      bool     `json:"pushed"`
		RemoteMoved bool     `json:"remote_moved"`
		DeletedTags []string `json:"deleted_tags"`
		KeptTags    []string `json:"kept_tags"`
		PushError   string   `json:"push_error"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("--json: %v %s", err, out)
	}
	if result.Pushed || !result.RemoteMoved || strings.Join(result.DeletedTags, " ") != "cli/v1.0.1" || result.KeptTags == nil || len(result.KeptTags) != 0 || result.PushError == "" {
		t.Errorf("--json: %+v", result)
	}
}
