package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
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

// heldStoryInReview is orchestratedStoryInReview's project with S-0001 held
// by an uncommitted file in its worktree, which the acceptance's preview
// lists as a blocker, and the story's file as it is then.
func heldStoryInReview(t *testing.T) (root, head string, story []byte) {
	t.Helper()
	root, head = orchestratedStoryInReview(t, false)
	if err := os.WriteFile(filepath.Join(root, ".flai-cache", "worktrees", "S-0001", "cli", "extra.go"), []byte("package cli\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, head, storyFile(t, root)
}

// storyFile is S-0001's file in the main checkout.
func storyFile(t *testing.T, root string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, "wip/kanban/stories/S-0001-ship-it.md"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// assertResolveChangedNothing fails unless S-0001's file is still before,
// its branch still at head, main still at mainBefore, and no thread is on
// it: flai accept --resolve moves nothing and opens no thread (S-0804).
func assertResolveChangedNothing(t *testing.T, root, head, mainBefore string, before []byte) {
	t.Helper()
	assertNothingAccepted(t, root, mainBefore)
	if got := gitIn(t, root, "rev-parse", "story/S-0001"); got != head {
		t.Errorf("story/S-0001 moved from %s to %s", head, got)
	}
	if after := storyFile(t, root); string(after) != string(before) {
		t.Errorf("--resolve changed S-0001's file:\n%s", after)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if all, err := threads.For(repo, "S-0001"); err != nil || len(all) != 0 {
		t.Errorf("--resolve opens no thread: %+v %v", all, err)
	}
}

// takeResolve stands in for flai serve, as takeCloseOut does: it waits for
// the request of the fixture's S-0001 in the project's cache dir, takes it,
// and answers with what answer makes of it. The request it read arrives on
// the channel, which is closed when none came within five seconds.
func takeResolve(t *testing.T, cacheDir string, answer func(closeout.Request) closeout.Answer) <-chan closeout.Request {
	const story = "S-0001"
	t.Helper()
	got := make(chan closeout.Request, 1)
	dir := closeout.Dir(cacheDir)
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			taken := closeout.File(dir, story, closeout.TakenSuffix)
			if os.Rename(closeout.File(dir, story, closeout.RequestSuffix), taken) == nil {
				var req closeout.Request
				data, err := os.ReadFile(taken)
				if err == nil {
					err = json.Unmarshal(data, &req)
				}
				if err != nil {
					t.Errorf("the request taken: %v", err)
				}
				got <- req
				if err := closeout.WriteAnswer(dir, answer(req)); err != nil {
					t.Errorf("the answer: %v", err)
				}
				return
			}
			time.Sleep(2 * time.Millisecond)
		}
		close(got)
	}()
	return got
}

// refusedOf is the refusal flai accept --resolve --json printed in out.
func refusedOf(t *testing.T, out string) resolveRefusal {
	t.Helper()
	var r struct {
		Refused resolveRefusal `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return r.Refused
}

// S-0804 criteria 1 and 8: flai accept --resolve on a story in review that
// its acceptance's blockers hold asks flai serve for a rebase and resolve job
// as the operator and answers at once with the job, saying what it does and
// where its output is; it moves nothing and opens no thread. A request while
// a job of the story runs answers that job, a rebase and resolve or a
// close-out, and starts nothing.
func TestAcceptResolveAsksFlaiServeForTheJobAndAnswersAtOnce(t *testing.T) {
	t.Parallel()
	root, head, before := heldStoryInReview(t)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	cache := openProject(t, root).CacheDir()
	log := closeout.LogPath(cache, "S-0001")
	const id = "S-0001-20261116T090000Z"

	got := takeResolve(t, cache, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Job: id, Story: r.Story, Started: true, Log: log}
	})
	out, code := resolveS0001(t, root, time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC))
	if code != 0 {
		t.Fatalf("flai accept --resolve: %d\n%s", code, out)
	}
	if req := <-got; req.Story != "S-0001" || req.Kind != closeout.KindResolve || req.By != "alex" || req.Message != "" || req.All || req.Fresh {
		t.Errorf("the request flai serve read: %+v", req)
	}
	var ans resolveAnswer
	if err := json.Unmarshal([]byte(out), &ans); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	does := "flai serve rebases story/S-0001 onto the main branch and closes it out at the rebased head, and S-0001 stays in review, ready to accept, unless the rebase stops on conflicts or a step of the close-out fails, when the job sends it back to its agent"
	want := resolveAnswer{ID: "S-0001", Job: id, Started: true, Log: log, Reason: "rebase and resolve job " + id + " started for S-0001: " + does + "; its output is in " + log}
	if ans != want {
		t.Errorf("--json\n got %+v\nwant %+v", ans, want)
	}
	for _, key := range []string{`"id"`, `"job"`, `"started"`, `"running"`, `"log"`, `"reason"`} {
		if !strings.Contains(out, key) {
			t.Errorf("--json must carry %s:\n%s", key, out)
		}
	}
	assertResolveChangedNothing(t, root, head, mainBefore, before)

	// a rebase and resolve job of the story runs: its id, and nothing starts
	dir := closeout.Dir(cache)
	if err := closeout.WriteJob(dir, closeout.Job{ID: id, Story: "S-0001", Kind: closeout.KindResolve, By: "alex", Started: time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	takeResolve(t, cache, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Job: id, Story: r.Story, Log: log}
	})
	text, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--resolve")
	if want := "rebase and resolve job " + id + " is already running for S-0001, so nothing more starts: " + does + "; its output is in " + log + "\n"; code != 0 || text != want {
		t.Errorf("a job running: exit %d\n got %q\nwant %q\n%s", code, text, want, errOut)
	}

	// a close-out job of the story runs: no rebase and resolve starts
	const closeOut = "S-0001-20261116T091500Z"
	if err := closeout.WriteJob(dir, closeout.Job{ID: closeOut, Story: "S-0001", By: "agent-S-0001", Started: time.Date(2026, 11, 16, 9, 15, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	takeResolve(t, cache, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Job: closeOut, Story: r.Story, Log: log}
	})
	out, code = resolveS0001(t, root, time.Date(2026, 11, 16, 9, 16, 0, 0, time.UTC))
	ans = resolveAnswer{}
	if err := json.Unmarshal([]byte(out), &ans); err != nil || code != 0 || ans.Job != closeOut || ans.Started || !ans.Running ||
		ans.Reason != "close-out job "+closeOut+" is running for S-0001, so no rebase and resolve starts: ask again once it ends; its output is in "+log {
		t.Errorf("a close-out running: exit %d %v\n%s", code, err, out)
	}
	assertResolveChangedNothing(t, root, head, mainBefore, before)
}

// S-0804, ADR-0302: flai accept --resolve on a story nothing holds is
// refused, saying to accept it, as one on a task and one not in review are;
// beside an acceptance's flags, --yes, or as the orchestrator it is refused
// before anything is read. None asks flai serve for anything, and none
// changes the story or opens a thread.
func TestAcceptResolveRefusesAStoryNothingHoldsAndTheFlagsOfAnAcceptance(t *testing.T) {
	t.Parallel()
	root, head := orchestratedStoryInReview(t, false)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	before := storyFile(t, root)
	cache := openProject(t, root).CacheDir()

	out, code := resolveS0001(t, root, time.Date(2026, 11, 14, 9, 0, 0, 0, time.UTC))
	if want := (resolveRefusal{ID: "S-0001", Reason: "nothing holds S-0001's acceptance, so there is nothing to rebase and resolve: accept it with flai accept S-0001"}); code != exitDocRefused || refusedOf(t, out) != want {
		t.Errorf("a story nothing holds must be refused with exit %d: %d\n%s", exitDocRefused, code, out)
	}
	out, _, code = runIn(t, root, "accept", "T-0001", "--by", "alex", "--resolve", "--json")
	if r := refusedOf(t, out); code != exitDocRefused || r.ID != "T-0001" || !strings.HasPrefix(r.Reason, "T-0001 is not a story: ") {
		t.Errorf("a task must be refused with exit %d: %d\n%s", exitDocRefused, code, out)
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
	if _, err := os.Stat(closeout.File(closeout.Dir(cache), "S-0001", closeout.RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("a refused --resolve asked flai serve: %v", err)
	}
	assertResolveChangedNothing(t, root, head, mainBefore, before)

	// a story not in review
	if _, errOut, code := runIn(t, root, "move", "S-0001", "in-progress", "--reason", "more to do"); code != 0 {
		t.Fatal(errOut)
	}
	out, code = resolveS0001(t, root, time.Date(2026, 11, 14, 9, 5, 0, 0, time.UTC))
	if r := refusedOf(t, out); code != exitDocRefused || r.ID != "S-0001" || !strings.HasPrefix(r.Reason, "S-0001 is in in-progress, not in review: ") {
		t.Errorf("a story in progress must be refused with exit %d: %d\n%s", exitDocRefused, code, out)
	}
}

// S-0804: a request flai serve refuses is refused with exit 4 and its
// reason, naming flai serve, and so is an answer with no job, as a flai serve
// older than the resolve job gives; neither changes the story or opens a
// thread.
func TestAcceptResolveRefusesWhatFlaiServeRefuses(t *testing.T) {
	t.Parallel()
	root, head, before := heldStoryInReview(t)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	cache := openProject(t, root).CacheDir()

	why := "S-0001 is in in-progress, not in review: rebase and resolve runs only for a story in review, so nothing runs"
	takeResolve(t, cache, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Story: r.Story, Reason: why}
	})
	out, code := resolveS0001(t, root, time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC))
	if want := (resolveRefusal{ID: "S-0001", Reason: "flai serve runs no rebase and resolve for S-0001: " + why}); code != exitDocRefused || refusedOf(t, out) != want {
		t.Errorf("a refusal: exit %d\n%s", code, out)
	}

	takeResolve(t, cache, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Story: r.Story, Recorded: true, Verdict: closeout.VerdictPassed, Commit: head}
	})
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--resolve")
	if code != exitDocRefused || !strings.Contains(errOut, "as a flai serve older than S-0804 does, so nothing runs: upgrade flai") {
		t.Errorf("an answer with no job: exit %d\n%s", code, errOut)
	}
	assertResolveChangedNothing(t, root, head, mainBefore, before)
}

// S-0804: with no flai serve to take the request, flai accept --resolve
// withdraws it and is refused with exit 4, saying to start flai serve; the
// story stays as it was. It shortens resolveJobWait, so it does not run in
// parallel.
func TestAcceptResolveWithoutFlaiServeIsRefused(t *testing.T) {
	root, head, before := heldStoryInReview(t)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	was := resolveJobWait
	resolveJobWait = closeout.Wait{Take: 50 * time.Millisecond, Poll: 5 * time.Millisecond}
	t.Cleanup(func() { resolveJobWait = was })

	out, code := resolveS0001(t, root, time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC))
	want := resolveRefusal{ID: "S-0001", Reason: "flai serve did not take the request for S-0001's rebase and resolve, so nothing runs and S-0001 stays in review: start flai serve for this project with flai host start serve, and ask again"}
	if code != exitDocRefused || refusedOf(t, out) != want {
		t.Errorf("no flai serve: exit %d\n%s", code, out)
	}
	cache := openProject(t, root).CacheDir()
	if _, err := os.Stat(closeout.File(closeout.Dir(cache), "S-0001", closeout.RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
	assertResolveChangedNothing(t, root, head, mainBefore, before)
}
