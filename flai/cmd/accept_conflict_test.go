package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// I-0066: a rebase conflict during a story's sync was git-added unresolved,
// and acceptance fast-forwarded main onto the markers. Acceptance must refuse
// such a branch, naming each marker's path and line and merging nothing, and
// accept it once the conflict is resolved and committed (S-0253). The refusal
// says to send the story back for its agent to resolve them, not to resolve
// them by hand (ADR-0248).
func TestAcceptRefusesABranchCarryingConflictMarkers(t *testing.T) {
	t.Parallel()
	root, wt, doc := conflictedStoryInReview(t)
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-conflicted.md")
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 {
		t.Fatal("acceptance must refuse a branch carrying conflict markers")
	}
	for _, want := range []string{
		"story/S-0001 carries merge conflict markers at design/x.md:4, design/x.md:8;",
		"send S-0001 back to in-progress so that its agent resolves each conflict in .flai-cache/worktrees/S-0001, keeping what both sides meant, removes the markers, commits, and moves it to review again (ADR-0248)",
	} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal must say %q:\n%s", want, errOut)
		}
	}
	if strings.Contains(errOut, "accept again") {
		t.Errorf("the refusal must not ask the operator to resolve the markers:\n%s", errOut)
	}
	if strings.Contains(errOut, "design/x.md:6") {
		t.Errorf("the divider is not a marker:\n%s", errOut)
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}
	if after, _ := os.ReadFile(storyFile); !strings.Contains(string(after), "status: review") {
		t.Error("a refused acceptance must leave the story in review")
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") == "" {
		t.Error("a refused acceptance must keep the story branch")
	}
	if _, err := os.Stat(wt); err != nil {
		t.Errorf("a refused acceptance must keep the worktree: %v", err)
	}

	// resolved and committed, the branch is accepted
	_ = os.WriteFile(doc, []byte("\n# X\n\nours and theirs\nend\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] resolve x")
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("a resolved branch must be accepted: %d %s", code, errOut)
	}
	if got := gitIn(t, root, "show", "main:design/x.md"); !strings.Contains(got, "ours and theirs") {
		t.Errorf("main must carry the resolution:\n%s", got)
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("an accepted story's branch is removed")
	}
}

// S-0276: the acceptance preview, flai accept --dry-run, which the
// dashboard's Accept dialog shows, names the same markers as a blocker before
// Accept is pressed, and names none once they are resolved.
func TestAcceptPreviewReportsConflictMarkers(t *testing.T) {
	t.Parallel()
	root, wt, doc := conflictedStoryInReview(t)
	mainBefore := gitIn(t, root, "rev-parse", "main")

	js, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res struct {
		Blockers []string `json:"blockers"`
	}
	if err := json.Unmarshal([]byte(js), &res); err != nil {
		t.Fatalf("%v: %s", err, js)
	}
	want := "story/S-0001 carries merge conflict markers at design/x.md:4, design/x.md:8; send S-0001 back to in-progress so that its agent resolves each conflict in .flai-cache/worktrees/S-0001"
	if len(res.Blockers) != 1 || !strings.HasPrefix(res.Blockers[0], want) {
		t.Errorf("the preview must name each marker as the one blocker, starting %q: %q", want, res.Blockers)
	}
	if text, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run"); !strings.Contains(text, "blocked: "+want) {
		t.Errorf("the dry run's text must name the blocker:\n%s", text)
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("a dry run moved main from %s to %s", mainBefore, got)
	}

	// resolved and committed, the preview names no blocker
	_ = os.WriteFile(doc, []byte("\n# X\n\nours and theirs\nend\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] resolve x")
	js, errOut, code = runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("%d %s", code, errOut)
	}
	res.Blockers = nil
	if err := json.Unmarshal([]byte(js), &res); err != nil {
		t.Fatalf("%v: %s", err, js)
	}
	if len(res.Blockers) != 0 {
		t.Errorf("a clean branch has no blocker: %q", res.Blockers)
	}
}

// S-0337: an acceptance whose rebase of the story's branch stops on a
// conflict is refused, and the stop is a line of the conflict log, met at
// acceptance; the rebase itself is aborted by the send back (ADR-0248), so
// none waits in the worktree.
func TestAcceptLogsTheConflictItsRebaseStopsOn(t *testing.T) {
	t.Parallel()
	root, wt, doc := conflictedStoryInReview(t)
	_ = os.WriteFile(doc, []byte("\n# X\n\nthe story's\nend\n"), 0o644)
	gitIn(t, wt, "commit", "-q", "-am", "docs: [S-0001] resolve x")
	_ = os.WriteFile(filepath.Join(root, "design", "x.md"), []byte("\n# X\n\nmain's\nend\n"), 0o644)
	gitIn(t, root, "add", "design/x.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: main's x")

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 || !strings.Contains(errOut, "stopped on conflicts in design/x.md") {
		t.Fatalf("acceptance must refuse the stopped rebase: %d %s", code, errOut)
	}
	got, err := storygit.ReadConflicts(filepath.Join(root, ".flai-cache"))
	if err != nil {
		t.Fatal(err)
	}
	want := []storygit.Conflict{{At: time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC), Story: "S-0001", Where: storygit.WhereAccept, Paths: []string{"design/x.md"}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("conflict log: %+v, want %+v", got, want)
	}
	if storygit.RebaseInProgress(execx.System{}, wt) {
		t.Error("the rebase still waits in the worktree")
	}
}

// ADR-0248: the acceptance by hand, without --by orchestrator, of a story
// whose rebase stops on a conflict the merge steps cannot settle leaves no
// rebase for the operator: it aborts the rebase, leaving the worktree clean
// at the branch's head and main where it was, moves the story back to
// in-progress as flai with the path in the reason, opens one thread on it by
// flai naming the path, the branch, and the acceptance, and says the story
// was sent back without telling the operator to resolve anything; with --json
// it is the refusal sent_back, exiting 4. Accepted again once the story is
// back in review, it writes on the same thread.
func TestAnAcceptanceByHandSendsBackAStoryWhoseRebaseItCannotSettle(t *testing.T) {
	t.Parallel()
	root, head := orchestratedStoryInReview(t, false)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(root, "cli", "feature.go"), []byte("package cli\n\n// main's\n"), 0o644)
	gitIn(t, root, "add", "cli/feature.go")
	gitIn(t, root, "commit", "-q", "-m", "feat: main's feature")
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code == 0 {
		t.Fatalf("the send back is not an acceptance: %d %s", code, errOut)
	}
	for _, want := range []string{"S-0001 was sent back to in-progress", "stopped on conflicts in cli/feature.go", "which the acceptance of S-0001 by alex cannot settle", "its agent resolves them on TH-"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal must say %q:\n%s", want, errOut)
		}
	}
	for _, no := range []string{"git rebase --continue", "resolve them and", "orchestrator"} {
		if strings.Contains(errOut, no) {
			t.Errorf("the refusal must not say %q:\n%s", no, errOut)
		}
	}

	// the worktree is clean at the branch's head, and main has not moved
	if storygit.RebaseInProgress(execx.System{}, wt) {
		t.Error("the rebase still waits in the worktree")
	}
	if st := gitIn(t, wt, "status", "--porcelain"); st != "" {
		t.Errorf("the worktree is not clean:\n%s", st)
	}
	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("the worktree is at %s, want the branch's head %s", got, head)
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}

	// the story is in progress, moved by flai with the path in the reason
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	last := story.Transitions[len(story.Transitions)-1]
	if story.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != "flai" {
		t.Errorf("the story is sent back by flai: %s, last %+v", story.Status, last)
	}
	if want := "moved to in-progress: the rebase of story/S-0001 onto main stopped on conflicts in cli/feature.go, which the acceptance of S-0001 by alex cannot settle, so it aborted the rebase; its agent settles it"; !strings.Contains(story.Body, want) {
		t.Errorf("the reason must name the path and the acceptance:\n%s", story.Body)
	}

	// one thread, by flai, naming the path, the branch, and the acceptance,
	// addressed to the agent
	all, err := threads.For(repo, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || !all[0].Open() || all[0].Opener() != "flai" ||
		all[0].Title != "Acceptance sent back: story/S-0001 stopped rebasing onto main on conflicts in cli/feature.go" {
		t.Fatalf("threads: %+v", all)
	}
	text := all[0].Entries()[0].Text
	for _, want := range []string{"The acceptance of S-0001 by alex rebased `story/S-0001` onto `main`", "- `cli/feature.go`", "This thread is for S-0001's agent", "run `flai stream sync S-0001`", "move S-0001 to review again"} {
		if !strings.Contains(text, want) {
			t.Errorf("the thread must say %q:\n%s", want, text)
		}
	}

	// back in review, the dashboard's acceptance, --json, sends it back again
	// as the refusal sent_back, writing on the same thread; later, since an
	// entry at the same time by the same author joins the last
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	later := time.Date(2026, 9, 15, 21, 5, 0, 0, time.UTC)
	out, errOut, code := runInAt(t, root, later, "accept", "S-0001", "--by", "alex", "--json")
	if code != exitDocRefused {
		t.Fatalf("the send back exits %d: %d %s\n%s", exitDocRefused, code, errOut, out)
	}
	var res struct {
		Refused sentBack `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	r := res.Refused
	if r.ID != "S-0001" || !r.SentBack || r.Branch != "story/S-0001" || r.Base != "main" || !reflect.DeepEqual(r.Conflicts, []string{"cli/feature.go"}) || r.Thread != all[0].ID {
		t.Errorf("--json must say the story was sent back, with the paths and the thread: %+v", r)
	}
	if !strings.Contains(r.Reason, "which the acceptance of S-0001 by alex cannot settle") || !strings.Contains(r.Reason, "its agent resolves them on "+all[0].ID) {
		t.Errorf("the reason names the acceptance and the thread: %s", r.Reason)
	}
	if again, _ := threads.For(repo, "S-0001"); len(again) != 1 || len(again[0].Entries()) != 2 || again[0].Entries()[1].Author != "flai" {
		t.Errorf("the second send back writes on the open thread: %+v", again)
	}
	if story, _ = repo.Get("S-0001"); story.Status != workitem.InProgress {
		t.Errorf("the story is sent back again: %s", story.Status)
	}

	// each stop is a line of the conflict log
	got, _ := storygit.ReadConflicts(filepath.Join(root, ".flai-cache"))
	if len(got) != 2 || got[0].Where != storygit.WhereAccept || got[1].Where != storygit.WhereAccept || !reflect.DeepEqual(got[1].Paths, []string{"cli/feature.go"}) {
		t.Errorf("conflict log: %+v", got)
	}
}

// ADR-0189, ADR-0248: the orchestrator's acceptance whose rebase stops on a
// conflict the merge steps cannot settle logs it, aborts the rebase, leaving
// the worktree clean at the branch's head and main where it was, moves the
// story back to in-progress as flai with the path in the reason, opens one
// thread on it naming the path and the branch, and says so, under --json as
// a refusal that exits 4, as the acceptance by hand does.
func TestTheOrchestratorSendsBackAStoryWhoseRebaseItCannotSettle(t *testing.T) {
	t.Parallel()
	root, head := orchestratedStoryInReview(t, true)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(root, "cli", "feature.go"), []byte("package cli\n\n// main's\n"), 0o644)
	gitIn(t, root, "add", "cli/feature.go")
	gitIn(t, root, "commit", "-q", "-m", "feat: main's feature")
	mainBefore := gitIn(t, root, "rev-parse", "main")

	evidence := "Verdict: pass\n- 1: `cli/feature.go`\n"
	out, errOut, code := runStdin(t, root, evidence, "accept", "S-0001", "--by", "orchestrator", "--verified", head, "--evidence", "-", "--json")
	if code != exitDocRefused {
		t.Fatalf("the send back exits %d: %d %s\n%s", exitDocRefused, code, errOut, out)
	}
	var res struct {
		Refused sentBack `json:"refused"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	r := res.Refused
	if r.ID != "S-0001" || !r.SentBack || r.Branch != "story/S-0001" || r.Base != "main" || !reflect.DeepEqual(r.Conflicts, []string{"cli/feature.go"}) || r.Thread == "" {
		t.Errorf("--json must say the story was sent back, with the paths and the thread: %+v", r)
	}
	for _, want := range []string{"S-0001 was sent back to in-progress", "conflicts in cli/feature.go", "its agent resolves them on " + r.Thread} {
		if !strings.Contains(r.Reason, want) || !strings.Contains(errOut, want) {
			t.Errorf("the refusal must say %q:\nreason %s\nstderr %s", want, r.Reason, errOut)
		}
	}

	// the worktree is clean at the branch's head, and main has not moved
	if storygit.RebaseInProgress(execx.System{}, wt) {
		t.Error("the rebase still waits in the worktree")
	}
	if st := gitIn(t, wt, "status", "--porcelain"); st != "" {
		t.Errorf("the worktree is not clean:\n%s", st)
	}
	if got := gitIn(t, wt, "rev-parse", "HEAD"); got != head {
		t.Errorf("the worktree is at %s, want the branch's head %s", got, head)
	}
	if got := gitIn(t, wt, "symbolic-ref", "--short", "HEAD"); got != "story/S-0001" {
		t.Errorf("the worktree has %s checked out", got)
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}

	// the story is in progress, moved by flai with the path in the reason
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	last := story.Transitions[len(story.Transitions)-1]
	if story.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != "flai" {
		t.Errorf("the story is sent back by flai: %s, last %+v", story.Status, last)
	}
	if want := "moved to in-progress: the rebase of story/S-0001 onto main stopped on conflicts in cli/feature.go, which the acceptance of S-0001 by orchestrator cannot settle"; !strings.Contains(story.Body, want) {
		t.Errorf("the reason must name the path:\n%s", story.Body)
	}

	// one thread, by flai, naming the path and the branch
	all, err := threads.For(repo, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].ID != r.Thread || !all[0].Open() || all[0].Opener() != "flai" ||
		all[0].Title != "Acceptance sent back: story/S-0001 stopped rebasing onto main on conflicts in cli/feature.go" {
		t.Fatalf("threads: %+v", all)
	}
	text := all[0].Entries()[0].Text
	for _, want := range []string{"- `cli/feature.go`", "`story/S-0001`", "run `flai stream sync S-0001`", "move S-0001 to review again"} {
		if !strings.Contains(text, want) {
			t.Errorf("the thread must say %q:\n%s", want, text)
		}
	}

	// the stop is a line of the conflict log, as the operator's is
	if got, _ := storygit.ReadConflicts(filepath.Join(root, ".flai-cache")); len(got) != 1 || got[0].Where != storygit.WhereAccept || !reflect.DeepEqual(got[0].Paths, []string{"cli/feature.go"}) {
		t.Errorf("conflict log: %+v", got)
	}
}

// I-0315, ADR-0276: an acceptance whose rebase stopped on a design document
// outside claims.shared, where the story and a story accepted before it each
// added a section at one spot, sent the story back. It now merges the file by
// sections and lands the story, with main's section and then the story's on
// main, and opens no thread.
func testAcceptMergesBySectionsADesignDocumentOutsideTheSharedPaths(t *testing.T) {
	root, wt := sectionStories(t)
	issueStoryInReview(t, root, "S-0001", "T-0001")

	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 {
		t.Fatalf("the acceptance stopped on %s: %d %s", causesDoc, code, errOut)
	}
	if strings.Contains(errOut, "sent back") {
		t.Errorf("the story was sent back:\n%s", errOut)
	}
	hasBothCauses(t, "main's "+causesDoc, gitIn(t, root, "show", "main:"+causesDoc))
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("an accepted story's branch is removed")
	}
	if _, err := os.Stat(wt); err == nil {
		t.Error("an accepted story's worktree is removed")
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if all, err := threads.For(repo, "S-0001"); err != nil || len(all) != 0 {
		t.Errorf("an acceptance that lands opens no thread: %v %+v", err, all)
	}
}

// conflictedStoryInReview makes a project with story S-0001 in review whose
// branch commits design/x.md with an unresolved conflict, its markers at
// lines 4 and 8, and returns the main checkout, the story's worktree, and
// the file.
func conflictedStoryInReview(t *testing.T) (root, wt, doc string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root = tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Conflicted", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-conflicted.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt = filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	doc = filepath.Join(wt, "design", "x.md")
	// The markers are built so that no line of this file is one. The file
	// opens with a blank line, which a reader that trims git's output would
	// lose and so misnumber every marker.
	conflicted := "\n# X\n\n" +
		strings.Repeat("<", 7) + " HEAD\nours\n" +
		strings.Repeat("=", 7) + "\ntheirs\n" +
		strings.Repeat(">", 7) + " f0f5443 (docs: [S-0201] keep both)\nend\n"
	_ = os.WriteFile(doc, []byte(conflicted), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] x")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root, wt, doc
}

// ADR-0248: a rebase the acceptance's sync found already in progress in the
// story's worktree, which it did not start, is sent back too, with a thread
// that tells the agent flai left the rebase as it was and how to finish or
// abort it; nothing in it, nor in the words the sync stop gives the operator
// otherwise, asks the operator to resolve a conflict.
func TestASendBackOverARebaseFoundInProgressLeavesItToTheAgent(t *testing.T) {
	t.Parallel()
	stop := &syncStopped{Story: "S-0001", Base: "main", Worktree: ".flai-cache/worktrees/S-0001", RebaseInProgress: true, Found: true, Conflicts: []string{"cli/x.go"}}
	text := sentBackText("S-0001", "story/S-0001", "the acceptance of S-0001 by alex", stop)
	for _, want := range []string{
		"The acceptance of S-0001 by alex found a rebase of `story/S-0001` already in progress in `.flai-cache/worktrees/S-0001`, which it did not start and left as it was",
		"- `cli/x.go`",
		"This thread is for S-0001's agent, which settles the conflicts; the operator does not.",
		"or run `git rebase --abort` and then `flai stream sync S-0001`",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("the thread text lacks %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "aborted") {
		t.Errorf("the thread text says flai aborted a rebase it did not start:\n%s", text)
	}
	stop.Found = false
	if text := sentBackText("S-0001", "story/S-0001", "the acceptance of S-0001 by alex", stop); !strings.Contains(text, "flai aborted the rebase, so `story/S-0001` and `main` are as they were") {
		t.Errorf("the thread text of a stop the acceptance made does not say the rebase was aborted:\n%s", text)
	}
}
