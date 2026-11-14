package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// resolveS0001 is what flai accept S-0001 --by alex --resolve --json says at
// at, with its exit code.
func resolveS0001(t *testing.T, root string, at time.Time) (string, int) {
	t.Helper()
	out, _, code := runInAt(t, root, at, "accept", "S-0001", "--by", "alex", "--resolve", "--json")
	return out, code
}

// S-0788, ADR-0302: flai accept --resolve on a story in review that its
// acceptance's blockers hold accepts nothing: it moves the story back to
// in-progress as flai, with a reason naming who asked, the branch, and the
// blocker, and opens one send-back thread on it by flai, titled as flai serve
// finds it, listing each blocker and telling the agent how to resolve them;
// --json prints the send back. A second --resolve after the story is back in
// review and held again writes on the same open thread, and one on a story
// no longer in review is refused.
func TestAcceptResolveSendsBackAStoryItsBlockersHold(t *testing.T) {
	t.Parallel()
	root, head := orchestratedStoryInReview(t, false)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	if err := os.WriteFile(filepath.Join(wt, "cli", "extra.go"), []byte("package cli\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mainBefore := gitIn(t, root, "rev-parse", "main")
	first := time.Date(2026, 11, 14, 9, 0, 0, 0, time.UTC)

	out, code := resolveS0001(t, root, first)
	if code != 0 {
		t.Fatalf("flai accept --resolve: %d\n%s", code, out)
	}
	if !strings.Contains(out, `"conflicts": []`) {
		t.Errorf("--json must give conflicts as an empty list, as a send back does:\n%s", out)
	}
	var sb sentBack
	if err := json.Unmarshal([]byte(out), &sb); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	blocker := "the worktree .flai-cache/worktrees/S-0001 has uncommitted changes (cli/extra.go); commit them on story/S-0001, or discard them, before accepting"
	if sb.ID != "S-0001" || !sb.SentBack || sb.Branch != "story/S-0001" || sb.Base != "main" || len(sb.Blockers) != 1 || sb.Blockers[0] != blocker || !strings.HasPrefix(sb.Thread, "TH-") {
		t.Errorf("--json must say the story was sent back, with its blockers and its thread: %+v", sb)
	}
	if want := "S-0001 was sent back to in-progress: alex asked its agent to rebase story/S-0001 onto main and resolve 1 blocker, on " + sb.Thread + " (ADR-0302)"; sb.Reason != want {
		t.Errorf("reason\n got %q\nwant %q", sb.Reason, want)
	}

	// nothing is accepted: main and the branch are where they were
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}
	if got := gitIn(t, root, "rev-parse", "story/S-0001"); got != head {
		t.Errorf("story/S-0001 moved from %s to %s", head, got)
	}

	// the story is in progress, moved by flai with who asked and the blocker
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	last := story.Transitions[len(story.Transitions)-1]
	if story.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != storygit.SentBackAuthor {
		t.Errorf("the story is sent back by flai: %s, last %+v", story.Status, last)
	}
	if want := "moved to in-progress: alex asked its agent to rebase story/S-0001 onto main and resolve what holds its acceptance, 1 blocker: " + blocker + "; its agent syncs it with flai stream sync S-0001"; !strings.Contains(story.Body, want) {
		t.Errorf("the reason must name who asked, the branch, and the blocker:\n%s", story.Body)
	}

	// one thread, by flai, with the prefix flai serve starts the agent on
	all, err := threads.For(repo, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	title := storygit.SentBackPrefix + "story/S-0001 cannot be accepted until its agent rebases it onto main and resolves 1 blocker, as alex asked"
	if len(all) != 1 || !all[0].Open() || all[0].Opener() != storygit.SentBackAuthor || all[0].Title != title || all[0].ID != sb.Thread {
		t.Fatalf("threads: %+v", all)
	}
	text := all[0].Entries()[0].Text
	for _, want := range []string{"alex asked S-0001's agent to rebase `story/S-0001` onto `main`", "- " + blocker + "\n", "This thread is for S-0001's agent", "run `flai stream sync S-0001`", "`flai verify S-0001 --job`", "move S-0001 to review again", "flai serve resolves this thread"} {
		if !strings.Contains(text, want) {
			t.Errorf("the thread must say %q:\n%s", want, text)
		}
	}

	// a story no longer in review is refused, and nothing changes
	out, code = resolveS0001(t, root, first.Add(time.Minute))
	var refused struct {
		Refused resolveRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &refused); err != nil || code != exitDocRefused || refused.Refused.ID != "S-0001" || !strings.Contains(refused.Refused.Reason, "S-0001 is in in-progress, not in review") {
		t.Errorf("a story in progress must be refused with exit %d: %d %v\n%s", exitDocRefused, code, err, out)
	}

	// back in review and held again, it writes on the open thread
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] the extra")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	if err := os.WriteFile(filepath.Join(wt, "cli", "more.go"), []byte("package cli\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, code = resolveS0001(t, root, first.Add(time.Hour)); code != 0 {
		t.Fatalf("the second --resolve: %d\n%s", code, out)
	}
	if again, _ := threads.For(repo, "S-0001"); len(again) != 1 || len(again[0].Entries()) != 2 || again[0].Entries()[1].Author != storygit.SentBackAuthor || !strings.Contains(again[0].Entries()[1].Text, "cli/more.go") {
		t.Errorf("the second send back writes on the open thread: %+v", again)
	}
}

// S-0788, ADR-0302: flai accept --resolve on a story nothing holds is
// refused, saying to accept it, and changes nothing; beside an acceptance's
// flags, --yes, or as the orchestrator it is refused before anything is read.
func TestAcceptResolveRefusesAStoryNothingHoldsAndTheFlagsOfAnAcceptance(t *testing.T) {
	t.Parallel()
	root, _ := orchestratedStoryInReview(t, false)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	out, code := resolveS0001(t, root, time.Date(2026, 11, 14, 9, 0, 0, 0, time.UTC))
	var refused struct {
		Refused resolveRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &refused); err != nil || code != exitDocRefused {
		t.Fatalf("a story nothing holds must be refused with exit %d: %d %v\n%s", exitDocRefused, code, err, out)
	}
	if want := "nothing holds S-0001's acceptance, so its agent has nothing to resolve: accept it with flai accept S-0001"; refused.Refused.ID != "S-0001" || refused.Refused.Reason != want {
		t.Errorf("refusal %+v, want reason %q", refused.Refused, want)
	}

	for _, args := range [][]string{
		{"--dry-run"}, {"--batch", "20261114-090000"}, {"--yes"}, {"--without-kluster"}, {"--trailer", "X: y"},
	} {
		_, errOut, code := runIn(t, root, append([]string{"accept", "S-0001", "--by", "alex", "--resolve"}, args...)...)
		if code == 0 || !strings.Contains(errOut, "--resolve accepts nothing") {
			t.Errorf("--resolve %v must be refused: %d %s", args, code, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "orchestrator", "--resolve"); code == 0 || !strings.Contains(errOut, "not the orchestrator's") {
		t.Errorf("--resolve by the orchestrator must be refused: %d %s", code, errOut)
	}

	assertNothingAccepted(t, root, mainBefore)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if all, _ := threads.For(repo, "S-0001"); len(all) != 0 {
		t.Errorf("a refused --resolve opens no thread: %+v", all)
	}
}
