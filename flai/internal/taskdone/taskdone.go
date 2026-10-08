// Package taskdone closes a task in one call (ADR-0107): it commits the
// story's worktree, tells the other open stories whose claim covers what the
// commit changed, syncs the story's branch, moves the task to done, logs the
// narrative, widens the touches, checks the story, and reads the agent's
// inbox, in that order, stopping at the first step that fails. flai task
// done, the MCP tool task_done, and the host channel's task.done answer its
// Result.
package taskdone

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The steps, in the order Run takes them, as Result.Stopped names them.
const (
	StepCommit  = "commit"
	StepTell    = "tell"
	StepSync    = "sync"
	StepMove    = "move"
	StepLog     = "log"
	StepTouches = "touches"
	StepCheck   = "check"
	StepInbox   = "inbox"
)

// Options say which task to close, in which project, as whom.
type Options struct {
	// Repo is the project, opened from its main checkout or from the story's
	// worktree: work items, narratives, and the cache are the main
	// checkout's either way.
	Repo   *workitem.Repo
	Runner execx.Runner
	Task   string // the task's ID, in any zero padding
	// Message is the commit message; its subject line is the log entry
	// when LogEntry is empty.
	Message  string
	LogEntry string
	// Agent is who closes the task: the move's actor, the narrative's
	// agent, and whose inbox is read; "agent" when empty.
	Agent   string
	Session string
	Now     func() time.Time // time.Now when nil
	Version string           // the running flai's, for the inbox
	Logger  *slog.Logger     // nil discards
}

// Result is what each step did, nil for a step not reached, and the step
// the run stopped at.
type Result struct {
	Task    string               `json:"task"`
	Story   string               `json:"story"`
	Commit  *Commit              `json:"commit"`
	Told    []Told               `json:"told"`
	Sync    *storygit.SyncResult `json:"sync"`
	Move    *Move                `json:"move"`
	Log     *Log                 `json:"log"`
	Touches *Touches             `json:"touches"`
	Check   *Check               `json:"check"`
	Inbox   *inbox.Inbox         `json:"inbox"`
	// Stopped is the step that failed, empty when every step ran.
	Stopped string `json:"stopped"`
	// Error says why that step failed and what to do.
	Error string `json:"error,omitempty"`
}

// Commit is the commit step's commit; Result.Commit is nil when there was
// nothing to commit.
type Commit struct {
	Hash    string   `json:"hash"`
	Subject string   `json:"subject"`
	Paths   []string `json:"paths"` // what it changed, as git names them
}

// Told is a story in progress or in review, other than the task's, whose
// claim covers paths the commit changed, and the conversation the notice of
// them went to it in (S-0333).
type Told struct {
	Story        string   `json:"story" jsonschema:"the other open story"`
	Title        string   `json:"title" jsonschema:"its title"`
	Paths        []string `json:"paths" jsonschema:"the changed paths its claim covers"`
	Conversation string   `json:"conversation" jsonschema:"the conversation, MS-nnnn, the notice is in"`
}

// Move is the task's move to done.
type Move struct {
	ID    string `json:"id"`
	State string `json:"state"`
	// Skipped is set when the task was done already, so was not moved again.
	Skipped bool `json:"skipped"`
	// Followed is each story or epic that moved with the task.
	Followed []workitem.Followed `json:"followed"`
	Warnings []string            `json:"warnings"`
}

// Log is the entry written to the story's narrative.
type Log struct {
	Stream string `json:"stream"`
	Entry  string `json:"entry"`
	At     string `json:"at"`
}

// Touches is what the touches step added to the task's touches and to the
// story's, and the stories in progress the story's claim grew into.
type Touches struct {
	Task     []string               `json:"task"`
	Story    []string               `json:"story"`
	Overlaps []itemedit.Overlapping `json:"overlaps"`
}

// Check is flai check --strict scoped to the story (ADR-0085): Findings are
// inside the story, Notes outside it.
type Check struct {
	Passed   bool            `json:"passed"`
	Findings []check.Finding `json:"findings"`
	Notes    []check.Finding `json:"notes"`
}

// ExitCode is flai task done's exit status for r (ADR-0107): 0 when every
// step ran, 3 when the sync stopped the run, 4 when the check did, else 1.
func (r Result) ExitCode() int {
	switch r.Stopped {
	case "":
		return 0
	case StepSync:
		return 3
	case StepCheck:
		return 4
	}
	return 1
}

// Run closes the task as ADR-0107 orders it. It returns an error only when
// it could not start: no message, or a task, story, or worktree it cannot
// find. Once the commit step has begun, a step that fails, whether refused
// or unable to run, is the Result's Stopped and Error, alongside what the
// steps before it did: they have changed the repository, and the caller
// reports them either way.
func Run(ctx context.Context, o Options) (Result, error) {
	r, err := start(o)
	if err != nil {
		return Result{Task: o.Task}, err
	}
	steps := []struct {
		name string
		do   func() error
	}{
		{StepCommit, r.commit},
		{StepTell, r.tell},
		{StepSync, r.sync},
		{StepMove, r.move},
		{StepLog, r.log},
		{StepTouches, r.touches},
		{StepCheck, r.check},
		{StepInbox, func() error { return r.inbox(ctx) }},
	}
	for _, s := range steps {
		if err := s.do(); err != nil {
			r.res.Stopped, r.res.Error = s.name, err.Error()
			r.logger.Info("task close stopped", "component", "taskdone", "task", r.task.ID, "step", s.name, "err", err)
			return r.res, nil
		}
	}
	return r.res, nil
}

// run is one close: the options with their defaults, what start found, and
// the result so far.
type run struct {
	o      Options
	logger *slog.Logger
	task   *workitem.Item
	story  *workitem.Item
	wt     string         // the story's worktree
	wtRepo *workitem.Repo // the project as checked out in the worktree
	res    Result
}

// start fills in the defaults and finds the task, its story, and the
// story's worktree.
func start(o Options) (*run, error) {
	if o.Repo == nil || o.Runner == nil {
		return nil, errors.New("taskdone needs a project and a runner")
	}
	if strings.TrimSpace(o.Message) == "" {
		return nil, fmt.Errorf("a commit message is needed: flai task done %s -m \"<message>\"", o.Task)
	}
	if o.Agent == "" {
		o.Agent = "agent"
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	r := &run{o: o, logger: o.Logger}
	if r.logger == nil {
		r.logger = slog.New(slog.DiscardHandler)
	}
	task, err := o.Repo.Get(o.Task)
	if err != nil {
		return nil, err
	}
	if task.Type != workitem.Task {
		return nil, fmt.Errorf("%s is a %s; flai task done closes a task (T-nnnn)", task.ID, task.Type)
	}
	if task.Parent == "" {
		return nil, fmt.Errorf("%s has no story; give it one with flai edit %s --parent S-nnnn", task.ID, task.ID)
	}
	story, err := o.Repo.Get(task.Parent)
	if err != nil {
		return nil, fmt.Errorf("read %s's story: %w", task.ID, err)
	}
	r.task, r.story = task, story
	r.wt = o.Repo.WorktreePath(story.ID)
	if st, err := os.Stat(r.wt); err != nil || !st.IsDir() {
		return nil, fmt.Errorf("%s has no worktree at %s; open one with flai stream open %s", story.ID, r.rel(r.wt), story.ID)
	}
	if r.wtRepo, err = workitem.Open(r.wt); err != nil {
		return nil, fmt.Errorf("read the project in %s's worktree %s: %w", story.ID, r.rel(r.wt), err)
	}
	// the worktree's items are the main checkout's, whatever the .git file says
	r.wtRepo.MainRoot = mainRoot(o.Repo)
	r.res = Result{Task: task.ID, Story: story.ID}
	return r, nil
}

// again is the command that runs this close again.
func (r *run) again() string { return "flai task done " + r.task.ID }

// git runs git in the story's worktree.
func (r *run) git(args ...string) (string, error) { return r.o.Runner.Run(r.wt, "git", args...) }

// commit stages everything in the worktree and commits it with the message;
// nothing to commit leaves Result.Commit nil. It refuses while a rebase is
// unfinished, as the sync does.
func (r *run) commit() error {
	if storygit.RebaseInProgress(r.o.Runner, r.wt) {
		waiting := ""
		if c := storygit.Conflicts(r.o.Runner, r.wt); len(c) > 0 {
			waiting = " on " + strings.Join(c, ", ")
		}
		return fmt.Errorf("a rebase is unfinished in %s%s: resolve each conflicting path, git add it, and run git rebase --continue there, or undo it with git rebase --abort; then run %s again", r.rel(r.wt), waiting, r.again())
	}
	if _, err := r.git("add", "-A"); err != nil {
		return fmt.Errorf("stage the changes in %s: %w", r.rel(r.wt), err)
	}
	staged, err := r.git("-c", "core.quotePath=false", "diff", "--cached", "--name-only", "--no-renames")
	if err != nil {
		return fmt.Errorf("list the changes staged in %s: %w", r.rel(r.wt), err)
	}
	paths := lines(staged)
	if len(paths) == 0 {
		return nil
	}
	if _, err := r.git("commit", "-q", "-m", r.o.Message); err != nil {
		return fmt.Errorf("commit on %s in %s: %w; fix what git says and run %s again", storygit.Branch(r.story.ID), r.rel(r.wt), err, r.again())
	}
	hash, err := r.git("rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("read the commit just made in %s: %w", r.rel(r.wt), err)
	}
	r.res.Commit = &Commit{Hash: strings.TrimSpace(hash), Subject: subject(r.o.Message), Paths: paths}
	return nil
}

// tell messages each other story in progress or in review whose claim covers
// a path the commit changed, the shared paths included, on the conversation
// open between the two stories or a new one (S-0333). It runs before the
// sync, so that a sync stopped on conflicts, after which a second call has
// nothing to commit, does not lose the notice. It never stops the run: a
// notice that cannot be sent is logged and left out of Result.Told.
func (r *run) tell() error {
	r.res.Told = []Told{}
	c := r.res.Commit
	if c == nil {
		return nil
	}
	items, err := r.o.Repo.List(false)
	if err != nil {
		r.logger.Warn("change notices not sent", "component", "taskdone", "task", r.task.ID, "err", err)
		return nil
	}
	for _, cov := range itemedit.Covering(items, r.o.Repo.Manifest.Projects, r.story.ID, c.Paths) {
		conv, err := messages.Notify(r.o.Repo, messages.SendOptions{
			From: r.story.ID, To: cov.ID, Author: r.o.Agent, Text: r.notice(cov), About: cov.Paths, Now: r.o.Now(),
		})
		if err != nil {
			r.logger.Warn("change notice not sent", "component", "taskdone", "task", r.task.ID, "item", cov.ID, "err", err)
			continue
		}
		r.res.Told = append(r.res.Told, Told{Story: cov.ID, Title: cov.Title, Paths: cov.Paths, Conversation: conv.ID})
	}
	return nil
}

// notice is the message that tells cov's agent of the commit: the task, its
// story, the commit and its subject, the paths cov's claim covers, and what
// to do about them.
func (r *run) notice(cov itemedit.Covered) string {
	c := r.res.Commit
	paths := make([]string, len(cov.Paths))
	for i, p := range cov.Paths {
		paths[i] = code(p)
	}
	hash := abbrev(c.Hash)
	return fmt.Sprintf("%s of %s changed paths %s's claim covers.\n\n%s, %s, committed %s on %s, %s, changing %s. It reaches the main branch when %s is accepted; %s shows it until then. Reply here if it breaks your work, or adjust to it early.",
		r.task.ID, r.story.ID, cov.ID,
		r.task.ID, r.task.Title, hash, storygit.Branch(r.story.ID), code(c.Subject), strings.Join(paths, ", "),
		r.story.ID, code("git show "+hash))
}

// sync is the story's flai stream sync; a refusal or a stop on conflicts
// stops the run with how to continue.
func (r *run) sync() error {
	res, err := storygit.Sync(storygit.SyncOptions{
		Runner: r.o.Runner, Repo: r.o.Repo, Story: r.story, Now: r.o.Now(),
		Files: issues.SyncFiles(r.o.Repo, r.story.ID, r.o.Runner, r.o.Now), Again: r.again(), Log: r.logger,
	})
	r.res.Sync = &res
	if err != nil {
		return fmt.Errorf("sync %s onto the main branch: %w", storygit.Branch(r.story.ID), err)
	}
	switch res.Stopped {
	case "":
		return nil
	case storygit.StopUncommitted:
		return fmt.Errorf("the sync refused %s: it has uncommitted changes in %s; %s", res.Worktree, strings.Join(res.Uncommitted, ", "), res.Continue)
	case storygit.StopRebaseInProgress:
		return fmt.Errorf("the sync refused %s: a rebase waits there on %s; %s, or %s", res.Worktree, orNothing(res.Conflicts), res.Continue, res.Abort)
	}
	return fmt.Errorf("the rebase of %s onto %s stopped on conflicts in %s; %s, or %s", res.Branch, res.Base, orNothing(res.Conflicts), res.Continue, res.Abort)
}

// move moves the task to done under flai move's rules, with what follows it,
// unless it is done already.
func (r *run) move() error {
	task, err := r.o.Repo.Get(r.task.ID)
	if err != nil {
		return fmt.Errorf("read %s to move it: %w", r.task.ID, err)
	}
	m := &Move{ID: task.ID, State: task.Status, Followed: []workitem.Followed{}, Warnings: []string{}}
	if task.Status == workitem.Done {
		m.Skipped = true
		r.res.Move = m
		return nil
	}
	moved, err := r.o.Repo.TransitionAll(task, workitem.Done, r.o.Agent, "", r.o.Now(), false)
	if err != nil {
		return fmt.Errorf("move %s to done: %w", task.ID, err)
	}
	m.State = task.Status
	if moved.Followed != nil {
		m.Followed = append(m.Followed, *moved.Followed)
	}
	m.Warnings = append(m.Warnings, moved.Warnings...)
	for _, w := range moved.Warnings {
		r.logger.Warn("workflow policy warning", "component", "workitem", "item", task.ID, "detail", w)
	}
	r.res.Move = m
	return nil
}

// log appends the log entry, or the message's subject line, to the story's
// narrative and writes the narratives' index again.
func (r *run) log() error {
	entry := strings.TrimSpace(r.o.LogEntry)
	if entry == "" {
		entry = subject(r.o.Message)
	}
	now := r.o.Now()
	n, err := r.o.Repo.LogStream(r.story.ID, entry, workitem.StreamOptions{Agent: r.o.Agent, Session: r.o.Session, Now: now})
	if err != nil {
		return fmt.Errorf("log to %s's narrative: %w", r.story.ID, err)
	}
	items, err := r.o.Repo.List(false)
	if err == nil {
		err = r.o.Repo.WriteIndex(items, now)
	}
	if err != nil {
		return fmt.Errorf("write the narratives' index after logging to %s: %w", r.story.ID, err)
	}
	r.res.Log = &Log{Stream: n.Stream, Entry: entry, At: n.Updated}
	return nil
}

// touches adds the paths the commit changed that the task's touches do not
// cover to the task, and those the story's do not cover to the story, as
// flai touches --add records them; the wip folder is flai's and is left out.
func (r *run) touches() error {
	t := &Touches{Task: []string{}, Story: []string{}, Overlaps: []itemedit.Overlapping{}}
	r.res.Touches = t
	if r.res.Commit == nil {
		return nil
	}
	w, err := itemedit.WidenStory(itemedit.WidenOptions{
		Repo: r.o.Repo, Story: r.story.ID, Task: r.task.ID, Paths: r.res.Commit.Paths, By: r.o.Agent, Now: r.o.Now,
	})
	t.Task, t.Story, t.Overlaps = w.Task, w.Story, w.Overlaps
	if err != nil {
		return err
	}
	if w.Untold != nil {
		// advisory, as flai touches has it: the touches stand
		r.logger.Warn("overlap notices not sent", "component", "touches", "item", r.story.ID, "err", w.Untold)
	}
	return nil
}

// check runs flai check --strict on the worktree, scoped to the story as the
// close-out scopes it (ADR-0085): a finding inside the story stops the run.
func (r *run) check() error {
	res, err := check.RunScoped(r.wtRepo, r.o.Runner, r.o.Now(), r.story.ID)
	if err != nil {
		return fmt.Errorf("run flai check in %s: %w", r.rel(r.wt), err)
	}
	c := &Check{Passed: res.OK(true), Findings: []check.Finding{}, Notes: []check.Finding{}}
	for _, f := range res.Findings {
		if f.Outside {
			c.Notes = append(c.Notes, f)
		} else {
			c.Findings = append(c.Findings, f)
		}
	}
	r.res.Check = c
	if !c.Passed {
		return fmt.Errorf("flai check --strict fails on %s in %s: fix them, then run %s again, which commits the fix", plural(len(c.Findings), "finding"), r.story.ID, r.again())
	}
	return nil
}

// inbox reads the agent's inbox as the MCP tool inbox does, advancing its
// cursor.
func (r *run) inbox(ctx context.Context) error {
	in, err := inbox.Read(ctx, inbox.Options{Repo: r.o.Repo, Agent: r.o.Agent, Own: r.story.ID, Now: r.o.Now, Runner: r.o.Runner, Version: r.o.Version})
	if err != nil {
		return fmt.Errorf("read %s's inbox: %w", r.o.Agent, err)
	}
	r.res.Inbox = &in
	return nil
}

// rel is p relative to the main checkout, for messages.
func (r *run) rel(p string) string {
	if rel, err := filepath.Rel(mainRoot(r.o.Repo), p); err == nil {
		return rel
	}
	return p
}

// mainRoot is the project's main checkout.
func mainRoot(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

// subject is a commit message's first line.
func subject(message string) string {
	first, _, _ := strings.Cut(strings.TrimSpace(message), "\n")
	return strings.TrimSpace(first)
}

// abbrev is a commit's hash cut to the seven characters git shows.
func abbrev(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

// code is s as a markdown code span, fenced with more backticks than any run
// of them in s, so that a subject or a path reads as git wrote it.
func code(s string) string {
	fence := "`"
	for strings.Contains(s, fence) {
		fence += "`"
	}
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		s = " " + s + " "
	}
	return fence + s + fence
}

// lines is git's output, one path a line, without empty lines.
func lines(out string) []string {
	var paths []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			paths = append(paths, l)
		}
	}
	return paths
}

// orNothing joins paths for a message.
func orNothing(paths []string) string {
	if len(paths) == 0 {
		return "no path"
	}
	return strings.Join(paths, ", ")
}

// plural is n with noun, made plural unless n is one.
func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
