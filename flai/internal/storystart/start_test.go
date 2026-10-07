package storystart

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)

// clock is when the tests start a story, an hour after the fixture was made.
var clock = t0.Add(time.Hour)

// noGit is a host without git: no branch or worktree is opened.
type noGit struct{ execx.System }

func (noGit) LookPath(string) (string, error) { return "", errors.New("not found") }

func (noGit) Run(string, string, ...string) (string, error) {
	return "", errors.New("git is not installed")
}

const convention = `---
title: Git
updated: 2026-08-01
audience: agent
order: 10
status: active
---

# Git

## Rules
- Commit at landing.

<!-- system-flow:end-of-baseline -->

## Project additions
`

// fixture is a project with one epic in backlog and its conventions; with
// conventions false it has no conventions folder, so no pack can be built.
type fixture struct {
	repo *workitem.Repo
	epic *workitem.Item
}

func newFixture(t *testing.T, conventions bool) *fixture {
	t.Helper()
	root := t.TempDir()
	write(t, root, "system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	for _, d := range []string{"design/system", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, root, "design/system/plan.md", "---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\ntext\n")
	if conventions {
		write(t, root, "design/conventions/README.md", "# Agent conventions\n\n| Order | File | Governs |\n|-------|------|---------|\n| 10 | [git.md](git.md) | Commits |\n")
		write(t, root, "design/conventions/git.md", convention)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	return &fixture{repo: repo, epic: epic}
}

// story makes a story of the epic with acceptance criteria, moved by the
// operator through states, and returns it as it is then.
func (f *fixture) story(t *testing.T, title string, opt workitem.NewOptions, states ...string) *workitem.Item {
	t.Helper()
	opt.Type, opt.Title, opt.Parent, opt.Owner, opt.Now = workitem.Story, title, f.epic.ID, "alex", t0
	if opt.Touches == nil {
		opt.Touches = []string{"docs/" + strings.ToLower(title)}
	}
	s, err := f.repo.Create(opt)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Dir(s.Path), filepath.Base(s.Path), strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1))
	for _, st := range states {
		s = f.get(t, s.ID)
		if _, err := f.repo.Transition(s, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	return f.get(t, s.ID)
}

func (f *fixture) get(t *testing.T, id string) *workitem.Item {
	t.Helper()
	it, err := f.repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	return it
}

func (f *fixture) options(id string, runner execx.Runner) Options {
	return Options{Repo: f.repo, Runner: runner, Story: id, Agent: "agent-" + id, Session: "s1", Host: "h1", Budget: "40KB", Now: func() time.Time { return clock }, Version: "test"}
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// A ready story starts in one call: it is in progress, moved by the agent,
// its narrative names the agent, the stream index lists it, the pack is
// built within its budget, and the inbox is the agent's, without its own
// move among the changes. Without git no branch or worktree is opened.
func TestStartStartsAReadyStory(t *testing.T) {
	f := newFixture(t, true)
	s := f.story(t, "One", workitem.NewOptions{}, workitem.Ready)
	res, err := Start(context.Background(), f.options(s.ID, noGit{}))
	if err != nil {
		t.Fatal(err)
	}
	got := f.get(t, s.ID)
	last := got.Transitions[len(got.Transitions)-1]
	if got.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != "agent-"+s.ID || last.At != clock.Format(workitem.TimeFormat) {
		t.Errorf("story %s, last transition %+v", got.Status, last)
	}
	if res.Story.ID != s.ID || res.Story.Status != workitem.InProgress || res.Story.Title != "One" || res.Story.Warnings == nil {
		t.Errorf("result story: %+v", res.Story)
	}
	n, err := workitem.ReadNarrative(f.repo.NarrativePath(s.ID))
	if err != nil {
		t.Fatalf("no narrative: %v", err)
	}
	if n.Agent != "agent-"+s.ID || n.Session != "s1" || n.Host != "h1" {
		t.Errorf("narrative agent %q, session %q, host %q", n.Agent, n.Session, n.Host)
	}
	index, err := os.ReadFile(filepath.Join(f.repo.Root, "wip/agents/index.md"))
	if err != nil || !strings.Contains(string(index), s.ID) {
		t.Errorf("the stream index does not list %s: %v\n%s", s.ID, err, index)
	}
	if res.Worktree != "" || res.Branch != "" || res.From != "" {
		t.Errorf("without git: worktree %q, branch %q, from %q", res.Worktree, res.Branch, res.From)
	}
	want, _ := ctxpack.ParseSize("40KB")
	if res.Pack == nil || res.Pack.Story != s.ID || len(res.Pack.Conventions) == 0 || res.Pack.Budget != want || res.Pack.Size.Bytes == 0 || res.Pack.Size.Bytes > want {
		t.Fatalf("pack: %+v", res.Pack)
	}
	if res.Inbox == nil || res.Inbox.Agent != "agent-"+s.ID {
		t.Fatalf("inbox: %+v", res.Inbox)
	}
	for _, c := range res.Inbox.Changes {
		if c.ID == s.ID && c.To == workitem.InProgress {
			t.Errorf("the agent's own move is reported to it: %s", c.Summary)
		}
	}
}

// The epic follows its first started story to in-progress, as flai move
// makes it, and the result says so; a second story's start moves it no more.
func TestStartFollowsTheEpicOnItsFirstStory(t *testing.T) {
	f := newFixture(t, true)
	first := f.story(t, "First", workitem.NewOptions{}, workitem.Ready)
	second := f.story(t, "Second", workitem.NewOptions{}, workitem.Ready)
	if e := f.get(t, f.epic.ID); e.Status != workitem.Ready {
		t.Fatalf("the epic should have followed its stories to ready: %s", e.Status)
	}
	res, err := Start(context.Background(), f.options(first.ID, noGit{}))
	if err != nil {
		t.Fatal(err)
	}
	if res.Followed == nil || res.Followed.ID != f.epic.ID || res.Followed.From != workitem.Ready || res.Followed.To != workitem.InProgress || res.Followed.Story != first.ID {
		t.Errorf("followed: %+v", res.Followed)
	}
	if e := f.get(t, f.epic.ID); e.Status != workitem.InProgress || e.Transitions[len(e.Transitions)-1].By != "agent-"+first.ID {
		t.Errorf("epic %s, transitions %+v", e.Status, e.Transitions)
	}
	res, err = Start(context.Background(), f.options(second.ID, noGit{}))
	if err != nil {
		t.Fatal(err)
	}
	if res.Followed != nil {
		t.Errorf("the epic is in progress already, yet followed: %+v", res.Followed)
	}
}

// refused asserts that err is a *Refused naming id in status, with a hold
// of code when code is not empty, and that the story was left as it was:
// the same transitions, no worktree. It returns the refusal's reason.
func refused(t *testing.T, f *fixture, before *workitem.Item, err error, code string) string {
	t.Helper()
	var r *Refused
	if !errors.As(err, &r) {
		t.Fatalf("want a refusal, got %v", err)
	}
	if r.Story != before.ID || r.Status != before.Status || !strings.Contains(r.Error(), before.ID) {
		t.Errorf("refusal: %+v", r)
	}
	switch {
	case code == "" && r.Hold != nil:
		t.Errorf("refused for its state, yet with a hold: %+v", r.Hold)
	case code != "" && (r.Hold == nil || r.Hold.Code != code || !strings.Contains(r.Reason, r.Hold.Reason)):
		t.Errorf("want a hold %s in the reason: %+v", code, r)
	}
	after := f.get(t, before.ID)
	if after.Status != before.Status || len(after.Transitions) != len(before.Transitions) {
		t.Errorf("%s changed: %s, %d transitions, was %s, %d", before.ID, after.Status, len(after.Transitions), before.Status, len(before.Transitions))
	}
	if exists(f.repo.WorktreePath(before.ID)) {
		t.Errorf("%s has a worktree", before.ID)
	}
	return r.Reason
}

// A story that is not ready is refused before anything changes: one in
// backlog, and one in progress, which another agent has started already.
func TestStartRefusesAStoryThatIsNotReady(t *testing.T) {
	f := newFixture(t, true)
	backlog := f.story(t, "Backlog", workitem.NewOptions{})
	_, err := Start(context.Background(), f.options(backlog.ID, noGit{}))
	if reason := refused(t, f, backlog, err, ""); !strings.Contains(reason, "not ready") {
		t.Errorf("reason: %s", reason)
	}
	if exists(f.repo.NarrativePath(backlog.ID)) {
		t.Errorf("%s has a narrative", backlog.ID)
	}

	started := f.story(t, "Started", workitem.NewOptions{}, workitem.Ready, workitem.InProgress)
	_, err = Start(context.Background(), f.options(started.ID, noGit{}))
	reason := refused(t, f, started, err, "")
	if !strings.Contains(reason, "in progress already") || !strings.Contains(reason, "flai stream open "+started.ID) {
		t.Errorf("reason: %s", reason)
	}
}

// A ready story whose claim overlaps a story in progress is refused with
// the hold's reason, which names the other story and what clears it.
func TestStartRefusesAStoryHeldByAnOverlap(t *testing.T) {
	f := newFixture(t, true)
	theirs := f.story(t, "Theirs", workitem.NewOptions{Touches: []string{"docs/shared"}}, workitem.Ready, workitem.InProgress)
	mine := f.story(t, "Mine", workitem.NewOptions{Touches: []string{"docs/shared/page.md"}}, workitem.Ready)
	_, err := Start(context.Background(), f.options(mine.ID, noGit{}))
	reason := refused(t, f, mine, err, workitem.HoldOverlap)
	if !strings.Contains(reason, theirs.ID+" moves to review") {
		t.Errorf("reason does not say what clears it: %s", reason)
	}
	if exists(f.repo.NarrativePath(mine.ID)) {
		t.Errorf("%s has a narrative", mine.ID)
	}
}

// A ready story that names in after: a story not done is refused with the
// hold's reason.
func TestStartRefusesAStoryHeldByAfter(t *testing.T) {
	f := newFixture(t, true)
	first := f.story(t, "First", workitem.NewOptions{})
	next := f.story(t, "Next", workitem.NewOptions{After: []string{first.ID}}, workitem.Ready)
	_, err := Start(context.Background(), f.options(next.ID, noGit{}))
	reason := refused(t, f, next, err, workitem.HoldAfter)
	if !strings.Contains(reason, "starts when "+first.ID+" is done") {
		t.Errorf("reason: %s", reason)
	}
	if exists(f.repo.NarrativePath(next.ID)) {
		t.Errorf("%s has a narrative", next.ID)
	}
}

// A step that fails after the move leaves the story in progress and says
// which step failed and the command that finishes it, with what the steps
// before it did.
func TestStartSaysWhichStepFailedAfterTheMove(t *testing.T) {
	f := newFixture(t, false)
	s := f.story(t, "One", workitem.NewOptions{}, workitem.Ready)
	res, err := Start(context.Background(), f.options(s.ID, noGit{}))
	var step *StepError
	if !errors.As(err, &step) {
		t.Fatalf("want a step error, got %v", err)
	}
	if step.Step != StepPrime || step.Finish != "flai prime --story "+s.ID || !strings.Contains(err.Error(), s.ID+" is in progress") {
		t.Errorf("step error: %v", err)
	}
	if got := f.get(t, s.ID); got.Status != workitem.InProgress {
		t.Errorf("the story is %s", got.Status)
	}
	if res.Story.Status != workitem.InProgress || !exists(f.repo.NarrativePath(s.ID)) || res.Pack != nil || res.Inbox != nil {
		t.Errorf("result: %+v", res)
	}
}

// With git, the start checks the story's new branch out from the main branch
// in its worktree. Integration: it runs real git.
func TestStartOpensTheBranchAndWorktree(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	f := newFixture(t, true)
	s := f.story(t, "One", workitem.NewOptions{}, workitem.Ready)
	root := f.repo.Root
	write(t, root, ".gitignore", ".flai-cache/\n")
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"commit", "-q", "-m", "project"}} {
		if out, err := (execx.System{}).Run(root, "git", args...); err != nil {
			t.Fatal(out, err)
		}
	}
	res, err := Start(context.Background(), f.options(s.ID, execx.System{}))
	if err != nil {
		t.Fatal(err)
	}
	if res.Branch != storygit.Branch(s.ID) || res.Worktree != f.repo.WorktreePath(s.ID) || res.From != storygit.FromMain {
		t.Errorf("branch %q, worktree %q, from %q", res.Branch, res.Worktree, res.From)
	}
	head, err := (execx.System{}).Run(res.Worktree, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil || strings.TrimSpace(head) != storygit.Branch(s.ID) {
		t.Errorf("the worktree is on %q: %v", head, err)
	}
	if f.get(t, s.ID).Status != workitem.InProgress || res.Pack == nil || res.Inbox == nil {
		t.Errorf("result: %+v", res)
	}
}
