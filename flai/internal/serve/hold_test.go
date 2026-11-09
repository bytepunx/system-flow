//go:build !windows

package serve

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// readyTouching makes a story that touches paths and moves it to ready.
func (lab *agentLab) readyTouching(title string, touches ...string) string {
	lab.t.Helper()
	id := lab.backlog(title, nil, touches...)
	lab.toReady(id)
	return id
}

// S-0128, ADR-0046: the launcher skips a ready story whose claim overlaps an
// open story's, starts the next that is clear within the limit, says why the
// held one waits, and starts it first once the open story is gone.
func TestAHeldStoryIsSkippedAndStartedFirstOnceClear(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(2)
	lab.hold()
	open := lab.readyTouching("Open", "flai/cmd")
	lab.move(open, workitem.InProgress)
	held := lab.readyTouching("Held", "flai/cmd/serve")
	clear := lab.readyTouching("Clear", "docs/clear")
	behind := lab.readyTouching("Behind", "docs/behind")
	lab.l.look(ctx, false)
	waitFor(t, "the clear story's agent runs", func() bool { return lab.run(clear).live() })
	// the lab is no git repository, so there is no trial judge and the
	// overlap holds as it did before (ADR-0230)
	why := "held (overlap): touches flai/cmd/serve, inside flai/cmd which " + open + " (in progress) touches; starts when " + open + " moves to review, is cancelled, or is sent back"
	if lab.run(held) != nil || lab.run(behind) != nil {
		t.Fatalf("started a held story, or one past the limit: %+v", lab.state().Stories)
	}
	if w := lab.state().Waiting; !strings.Contains(w, held+" "+why) || !strings.Contains(w, "limit leaves no room for "+behind) {
		t.Errorf("waiting: %q", w)
	}
	// agent.status: waiting, with the reason, though it never had an agent
	a := Activity(lab.root, lab.state(), lab.cfg.commandSet())[held]
	if a.State != ActivityWaiting || a.Why != why || a.Hold == nil || a.Hold.Code != workitem.HoldOverlap || a.Run == nil || a.Run.Story != held || a.Run.PID != 0 {
		t.Errorf("activity: %+v", a)
	}
	if n := lab.said(held); n != 1 {
		t.Errorf("logged %d times why %s waits", n, held)
	}

	// the open story is sent back: the held story keeps its place and starts first
	st, _ := lab.repo.Get(open)
	if _, err := lab.repo.Transition(st, workitem.Cancelled, "alex", "not now", lab.now); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "the held story's agent runs", func() bool { return lab.run(held).live() })
	if lab.run(behind) != nil || !strings.Contains(lab.state().Waiting, "limit leaves no room for "+behind) {
		t.Errorf("the story behind it took its place: %+v", lab.state())
	}
	if _, ok := Activity(lab.root, lab.state(), lab.cfg.commandSet())[held]; !ok || Activity(lab.root, lab.state(), lab.cfg.commandSet())[held].Hold != nil {
		t.Errorf("still held once started: %+v", Activity(lab.root, lab.state(), lab.cfg.commandSet())[held])
	}
	for _, id := range []string{held, clear} {
		lab.release(id)
	}
	waitFor(t, "both end", func() bool { return !lab.run(held).live() && !lab.run(clear).live() })
}

// One look never starts two stories that overlap: the first one started
// claims its paths before its agent pulls it. A story with no touches is
// held by any story open, and holds every other while it is open itself.
func TestOneLookStartsNoTwoThatOverlap(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.hold()
	first := lab.readyTouching("First", "flai")
	second := lab.readyTouching("Second", "flai/internal/serve")
	bare := lab.readyTouching("Bare", []string{}...)
	lab.l.look(ctx, false)
	waitFor(t, "the first runs", func() bool { return lab.run(first).live() })
	if lab.run(second) != nil || lab.run(bare) != nil {
		t.Fatalf("started beside an overlapping story: %+v", lab.state().Stories)
	}
	act := Activity(lab.root, lab.state(), lab.cfg.commandSet())
	if a := act[second]; a.State != ActivityWaiting || !strings.Contains(a.Why, "inside flai which "+first+" (its agent started) touches") {
		t.Errorf("second: %+v", a)
	}
	if a := act[bare]; a.Hold == nil || a.Hold.Code != workitem.HoldNoTouches || !strings.Contains(a.Why, "declares no touches") {
		t.Errorf("bare: %+v", a)
	}
	lab.release(first)
	waitFor(t, "the first ends", func() bool { return !lab.run(first).live() })
}

// I-0301: a checked creation writes a task, checks it, and removes it when
// the check refuses it. A look that read the task in that window judged the
// ready story's claim by the task's touches, found no hold, and started an
// agent that story_start then refused as held no-touches. The look reads the
// items under the checked-write lock the creation holds, so it waits for the
// removal and holds the story.
func TestALookWaitsForACheckedCreationAndStartsNoStoryHeldOnceItIsRemoved(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(3)
	lab.hold()
	open := lab.readyTouching("Open", "flai/cmd")
	lab.move(open, workitem.InProgress)
	bare := lab.readyTouching("Bare", []string{}...)
	// the creation's write, under the lock, of a task the check will refuse
	unlock, err := lab.repo.LockCheckedWrites()
	if err != nil {
		t.Fatal(err)
	}
	task, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Refused", Parent: bare, Owner: "planner-" + bare, Touches: []string{"docs/elsewhere.md"}, Now: lab.now})
	if err != nil {
		unlock()
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		lab.l.look(ctx, false)
	}()
	time.Sleep(300 * time.Millisecond)
	// the check refused it: the creation removes it, then lets go
	if err := os.Remove(task.Path); err != nil {
		unlock()
		t.Fatal(err)
	}
	unlock()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("the look did not end once the lock was free")
	}
	if r := lab.run(bare); r != nil {
		lab.release(bare)
		t.Fatalf("started %s on a task the check removed: %+v", bare, r)
	}
	if w := lab.state().Waiting; !strings.Contains(w, bare+" held (no-touches)") {
		t.Errorf("waiting: %q", w)
	}
}

// S-0550, ADR-0096: an overlap wholly inside claims.shared holds nothing, so
// one look starts two ready stories that both touch docs/users/flai.md and a
// file each of their own.
func TestOneLookStartsTwoThatShareOnlyASharedPath(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	manifest := filepath.Join(lab.root, "system-flow.yaml")
	m, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, append(m, "claims:\n  shared:\n    - docs/users/flai.md\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.limit(2)
	lab.hold()
	first := lab.readyTouching("First", "docs/users/flai.md", "flai/cmd/accept.go")
	second := lab.readyTouching("Second", "docs/users/flai.md", "flai/internal/serve/hold.go")
	lab.l.look(ctx, false)
	waitFor(t, "both run", func() bool { return lab.run(first).live() && lab.run(second).live() })
	if w := lab.state().Waiting; strings.Contains(w, first) || strings.Contains(w, second) {
		t.Errorf("waiting names a started story: %q", w)
	}
	for _, id := range []string{first, second} {
		lab.release(id)
	}
	waitFor(t, "both end", func() bool { return !lab.run(first).live() && !lab.run(second).live() })
}

// The operator's word starts a held story's agent, with a warning.
func TestAHeldStoryStartsOnTheOperatorsWord(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.hold()
	open := lab.readyTouching("Open", "docs")
	lab.move(open, workitem.InProgress)
	id := lab.readyTouching("Held", "docs/guide.md")
	lab.l.look(ctx, false)
	if lab.run(id) != nil {
		t.Fatal("the launcher started a held story")
	}
	run, err := Start(ctx, lab.o, Entry{Key: "t", Name: "t", Root: lab.root}, id)
	if err != nil || !run.live() {
		t.Fatalf("start: %+v %v", run, err)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if !strings.Contains(logged, `level=WARN msg="agent started for a held story"`) || !strings.Contains(logged, "held (overlap): touches docs/guide.md, inside docs which "+open) {
		t.Errorf("no warning: %s", logged)
	}
	lab.release(id)
}

// S-0130: the launcher does not start a ready story that names, in after:, a
// story that is not done, says why, and starts it once that story is done.
func TestAStoryWaitsForTheStoryItNamesInAfter(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(2)
	lab.hold()
	first := lab.backlog("First", nil)
	then := lab.backlog("Then", nil)
	st, err := lab.repo.Get(then)
	if err != nil {
		t.Fatal(err)
	}
	st.After = []string{first}
	if err := lab.repo.Save(st); err != nil {
		t.Fatal(err)
	}
	lab.toReady(then)
	lab.l.look(ctx, false)
	why := "held (after): waits for " + first + " (in backlog); starts when " + first + " is done"
	if lab.run(then) != nil {
		t.Fatalf("started a story whose after: story is in backlog: %+v", lab.state().Stories)
	}
	if a := Activity(lab.root, lab.state(), lab.cfg.commandSet())[then]; a.State != ActivityWaiting || a.Why != why || a.Hold == nil || a.Hold.Code != workitem.HoldAfter {
		t.Errorf("activity: %+v", a)
	}

	// the named story is done and archived: the waiting one starts
	f, _ := lab.repo.Get(first)
	f.Status = workitem.Done
	if err := lab.repo.Save(f); err != nil {
		t.Fatal(err)
	}
	items, _ := lab.repo.List(false)
	plan, err := lab.repo.PlanArchive(items, []string{first})
	if err != nil {
		t.Fatal(err)
	}
	if err := lab.repo.Archive(plan); err != nil {
		t.Fatal(err)
	}
	lab.l.look(ctx, false)
	waitFor(t, "the waiting story's agent runs", func() bool { return lab.run(then).live() })
	lab.release(then)
	waitFor(t, "it ends", func() bool { return !lab.run(then).live() })
}

// S-0182, I-0050: whether the holder is in progress or has its agent
// started, the launcher starts no ready story whose claim is a folder holding
// a file the holder claims, nor one whose claim is a file in a folder the
// holder claims, and starts it once the holder is gone.
func TestTheLauncherHoldsAFileAndItsFolderBothWays(t *testing.T) {
	for _, c := range []struct {
		name, holder, held, state, label string
	}{
		{"a folder, while a file in it is in progress", "flai/cmd/prime.go", "flai/cmd", workitem.InProgress, "in progress"},
		{"a file, while its folder is in progress", "flai/cmd", "flai/cmd/prime.go", workitem.InProgress, "in progress"},
		{"a folder, while a file in it has its agent started in ready", "flai/cmd/prime.go", "flai/cmd", workitem.Ready, agentStarted},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab := newAgentLab(t)
			ctx := context.Background()
			lab.limit(3)
			lab.hold()
			holder := lab.readyTouching("Holder", c.holder)
			switch c.state {
			case workitem.Ready:
				lab.l.look(ctx, false)
				waitFor(t, "the holder's agent runs", func() bool { return lab.run(holder).live() })
			default:
				lab.move(holder, c.state)
			}
			held := lab.readyTouching("Held", c.held)
			lab.l.look(ctx, false)
			if r := lab.run(held); r != nil {
				t.Fatalf("started a held story: %+v", r)
			}
			if w := lab.state().Waiting; !strings.Contains(w, held+" held (overlap): touches "+c.held) || !strings.Contains(w, holder+" ("+c.label+")") {
				t.Errorf("waiting: %q", w)
			}
			if a := Activity(lab.root, lab.state(), lab.cfg.commandSet())[held]; a.Hold == nil || a.Hold.Code != workitem.HoldOverlap {
				t.Errorf("activity: %+v", a)
			}
			// the holder is cancelled: the held story is started
			st, _ := lab.repo.Get(holder)
			if _, err := lab.repo.Transition(st, workitem.Cancelled, "alex", "not now", lab.now); err != nil {
				t.Fatal(err)
			}
			lab.l.look(ctx, false)
			waitFor(t, "the held story's agent runs", func() bool { return lab.run(held).live() })
			lab.release(holder)
			lab.release(held)
		})
	}
}

// S-0295, I-0087, ADR-0096: a story in review holds nothing. The launcher
// starts a ready story whose claim overlaps only that of a story in review,
// and says nothing of a hold.
func TestTheLauncherStartsAStoryOverlappingOnlyOneInReview(t *testing.T) {
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.limit(2)
	lab.hold()
	review := lab.readyTouching("Review", "flai/cmd")
	lab.move(review, workitem.InProgress)
	lab.move(review, workitem.Review)
	ready := lab.readyTouching("Ready", "flai/cmd/serve")
	lab.l.look(ctx, false)
	waitFor(t, "the ready story's agent runs", func() bool { return lab.run(ready).live() })
	if w := lab.state().Waiting; strings.Contains(w, ready) {
		t.Errorf("waiting names the started story: %q", w)
	}
	if a := Activity(lab.root, lab.state(), lab.cfg.commandSet())[ready]; a.Hold != nil {
		t.Errorf("held by a story in review: %+v", a)
	}
	lab.release(ready)
	waitFor(t, "it ends", func() bool { return !lab.run(ready).live() })
}

// trialFile is the file the stories of the trial-merge tests change.
const trialFile = "flai/cmd/serve.go"

// gitMain makes the lab a git repository on main, whose one commit holds
// trialFile with ten numbered lines, for the trial merge that judges an
// overlap (ADR-0230).
func (lab *agentLab) gitMain() {
	lab.t.Helper()
	if testing.Short() {
		lab.t.Skip("needs git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		lab.t.Skip("git not installed")
	}
	_ = os.WriteFile(filepath.Join(lab.root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	lab.git(lab.root, "init", "-q", "-b", "main")
	lab.git(lab.root, "config", "user.email", "t@t")
	lab.git(lab.root, "config", "user.name", "t")
	lab.git(lab.root, "config", "commit.gpgsign", "false")
	var lines strings.Builder
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(&lines, "line %d\n", i)
	}
	_ = os.MkdirAll(filepath.Dir(filepath.Join(lab.root, trialFile)), 0o755)
	if err := os.WriteFile(filepath.Join(lab.root, trialFile), []byte(lines.String()), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.git(lab.root, "add", "-A")
	lab.git(lab.root, "commit", "-q", "-m", "init")
}

// storyBranch gives story id its branch from main, checked out in its
// worktree, which it is.
func (lab *agentLab) storyBranch(id string) string {
	lab.t.Helper()
	wt := lab.repo.WorktreePath(id)
	lab.git(lab.root, "worktree", "add", "-q", "-b", "story/"+id, wt)
	return wt
}

// commitLine sets line n, from 1, of trialFile in the checkout dir to text
// and commits it there.
func (lab *agentLab) commitLine(dir string, n int, text string) {
	lab.t.Helper()
	path := filepath.Join(dir, trialFile)
	data, err := os.ReadFile(path)
	if err != nil {
		lab.t.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	lines[n-1] = text
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.git(dir, "add", trialFile)
	lab.git(dir, "commit", "-q", "-m", "change line "+strconv.Itoa(n)+" of "+trialFile)
}

// S-0551, ADR-0230: a ready story whose claim overlaps that of a story in
// progress is not held when the holder's branch merges cleanly at main: it
// waits only for room in the in-progress limit, no one is asked about a hold,
// and the launcher starts it and journals the start as any other.
func TestTheLauncherStartsAnOverlappingStoryOnACleanTrialMerge(t *testing.T) {
	t.Parallel()
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.gitMain()
	lab.limit(1)
	lab.hold()
	holder := lab.readyTouching("Holder", trialFile)
	lab.move(holder, workitem.InProgress)
	lab.commitLine(lab.storyBranch(holder), 1, "the holder's line 1")
	lab.commitLine(lab.root, 10, "main's line 10")
	ready := lab.readyTouching("Ready", "flai/cmd")

	// the limit is full: it waits for room, not on a hold, and no one is asked
	lab.l.look(ctx, false)
	if r := lab.run(ready); r != nil {
		t.Fatalf("started past the limit: %+v", r)
	}
	if w := lab.state().Waiting; strings.Contains(w, "held") || !strings.Contains(w, "the in-progress limit leaves no room for "+ready) {
		t.Errorf("waiting: %q", w)
	}
	if a := Activity(lab.root, lab.state(), lab.cfg.commandSet())[ready]; a.Hold != nil {
		t.Errorf("activity says it is held: %+v", a)
	}
	if asks := lab.asks(ready); len(asks) != 0 {
		t.Errorf("asked about a hold a clean trial clears: %v", asks)
	}

	// room: it starts beside the story it overlaps, journaled as any start
	lab.limit(2)
	lab.l.look(ctx, false)
	waitFor(t, "the overlapping story's agent runs", func() bool { return lab.run(ready).live() })
	if j := lab.entries(); len(j) != 1 || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started stub-agent (command) for "+ready) {
		t.Errorf("journal: %+v", j)
	}
	if w := lab.state().Waiting; strings.Contains(w, ready) {
		t.Errorf("waiting names the started story: %q", w)
	}
	lab.release(ready)
	waitFor(t, "it ends", func() bool { return !lab.run(ready).live() })
}

// S-0551, ADR-0230: a ready story whose claim overlaps that of a story in
// progress is held when the holder's branch conflicts with main, and the
// launcher's journal of why it waits names the path the trial merge
// conflicts on, as the activity does; the holder's agent is asked about it.
func TestTheLauncherHoldsAnOverlappingStoryOnAConflictingTrialMerge(t *testing.T) {
	t.Parallel()
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.gitMain()
	lab.limit(2)
	lab.hold()
	holder := lab.readyTouching("Holder", trialFile)
	lab.move(holder, workitem.InProgress)
	lab.commitLine(lab.storyBranch(holder), 1, "the holder's line 1")
	lab.commitLine(lab.root, 1, "main's line 1")
	held := lab.readyTouching("Held", "flai/cmd")
	lab.l.look(ctx, false)
	if r := lab.run(held); r != nil {
		t.Fatalf("started a story whose holder's trial merge conflicts: %+v", r)
	}
	why := held + " held (overlap): touches flai/cmd, which holds " + trialFile + " that " + holder + " (in progress) touches, and a trial merge of story/" + holder + " at " + held + "'s base conflicts in " + trialFile + "; starts when " + holder
	if w := lab.state().Waiting; !strings.Contains(w, why) {
		t.Errorf("waiting: %q\nwant: %q", w, why)
	}
	if n := lab.said(held); n != 1 {
		t.Errorf("logged %d times why %s waits", n, held)
	}
	if a := Activity(lab.root, lab.state(), lab.cfg.commandSet())[held]; a.Hold == nil || a.Hold.Code != workitem.HoldOverlap || !strings.Contains(a.Why, "conflicts in "+trialFile) {
		t.Errorf("activity: %+v", a)
	}
	if asks := lab.asks(held); asks[holder] != 1 {
		t.Errorf("asks: %v", asks)
	}
	if j := lab.entries(); len(j) != 0 {
		t.Errorf("journal: %+v", j)
	}
}

// S-0551, ADR-0230: a story in progress with no branch yet cannot be
// trial-merged, so a ready story whose claim overlaps its claim is held as
// before, and why it waits says the trial could not run.
func TestTheLauncherHoldsAnOverlappingStoryWhileTheHolderHasNoBranch(t *testing.T) {
	t.Parallel()
	lab := newAgentLab(t)
	ctx := context.Background()
	lab.gitMain()
	lab.limit(2)
	lab.hold()
	holder := lab.readyTouching("Holder", trialFile)
	lab.move(holder, workitem.InProgress)
	held := lab.readyTouching("Held", "flai/cmd")
	lab.l.look(ctx, false)
	if r := lab.run(held); r != nil {
		t.Fatalf("started a story whose holder has no branch: %+v", r)
	}
	why := held + " held (overlap): touches flai/cmd, which holds " + trialFile + " that " + holder + " (in progress) touches; the trial merge could not run: " + holder + " has no branch story/" + holder + " yet"
	if w := lab.state().Waiting; !strings.Contains(w, why) {
		t.Errorf("waiting: %q\nwant: %q", w, why)
	}
}
