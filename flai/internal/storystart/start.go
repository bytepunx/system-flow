// Package storystart starts a ready story in one call for flai story start, story_start, and story.start.
package storystart

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The steps after the move, in the order Start takes them, as StepError
// names them.
const (
	StepStream = "stream"
	StepPrime  = "prime"
	StepInbox  = "inbox"
)

// Options say which story to start, in which project, as whom.
type Options struct {
	Repo   *workitem.Repo
	Runner execx.Runner
	Story  string // the story's ID, in any zero padding
	// Agent is who starts the story: the move's actor, the narrative's
	// agent, and whose inbox is read; "agent" when empty.
	Agent   string
	Session string
	// Host is the host the narrative records; workitem.ThisHost when empty.
	Host string
	// Budget is the prime pack's size, such as 80KB; empty takes the
	// manifest's prime.budget, else flai's default.
	Budget  string
	Now     func() time.Time // time.Now when nil
	Version string           // the running flai's, for the inbox
	// RelativePaths reports whether to link a new worktree with relative
	// paths, as storygit.OpenOptions has it; nil links an ordinary one.
	RelativePaths func() bool
	Log           *slog.Logger // nil discards
}

// Result is what a start did. After a failed step it holds what the steps
// before it did.
type Result struct {
	Story    Moved              `json:"story"`
	Followed *workitem.Followed `json:"followed"`
	Worktree string             `json:"worktree"`
	Branch   string             `json:"branch"`
	// From is where the branch came from, as storygit.Opened.From says it.
	From  string        `json:"from"`
	Pack  *ctxpack.Pack `json:"pack"`
	Inbox *inbox.Inbox  `json:"inbox"`
}

// Moved is the started story after its move, with the move's warnings.
type Moved struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Status   string   `json:"status"`
	Warnings []string `json:"warnings"`
}

// Refused is Start's refusal, made before anything changed: the story is not
// ready, or the board holds it.
type Refused struct {
	Story  string
	Status string
	// Hold is the board's hold on the story, nil when it was refused for its
	// state.
	Hold   *workitem.Hold
	Reason string // why, and what clears it
}

func (e *Refused) Error() string { return e.Reason }

// StepError is a step after the move that failed: the story stays in
// progress, and Finish is the command that does what the step did not.
type StepError struct {
	Story  string
	Step   string
	Finish string
	Err    error
}

func (e *StepError) Error() string {
	return fmt.Sprintf("%s is in progress, but its %s step failed: %v; finish it with %s", e.Story, e.Step, e.Err, e.Finish)
}

func (e *StepError) Unwrap() error { return e.Err }

// Start refuses a story that is not ready or that the board holds, and
// otherwise moves it to in-progress as flai move does, its epic following,
// then opens its narrative and its branch in a worktree as flai stream open
// does, builds its prime pack within the budget, and reads the agent's
// inbox. A refusal is a *Refused; a step that fails after the move is a
// *StepError, returned with what the steps before it did.
func Start(ctx context.Context, o Options) (Result, error) {
	if o.Repo == nil || o.Runner == nil {
		return Result{}, errors.New("storystart needs a project and a runner")
	}
	if o.Agent == "" {
		o.Agent = "agent"
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Host == "" {
		o.Host = workitem.ThisHost()
	}
	if o.Log == nil {
		o.Log = slog.New(slog.DiscardHandler)
	}
	story, err := o.Repo.Get(o.Story)
	if err != nil {
		return Result{}, fmt.Errorf("start %s: %w; give the ID of a ready story", o.Story, err)
	}
	if err := refuse(o.Repo, story); err != nil {
		return Result{}, err
	}
	now := o.Now()
	moved, err := o.Repo.TransitionAll(story, workitem.InProgress, o.Agent, "", now, false)
	if err != nil {
		return Result{}, fmt.Errorf("move %s to in-progress: %w", story.ID, err)
	}
	for _, w := range moved.Warnings {
		o.Log.Warn("workflow policy warning", "component", "workitem", "item", story.ID, "detail", w)
	}
	warnings := moved.Warnings
	if warnings == nil {
		warnings = []string{}
	}
	res := Result{Story: Moved{ID: story.ID, Title: story.Title, Status: story.Status, Warnings: warnings}, Followed: moved.Followed}
	o.Log.Info("story started", "component", "storystart", "story", story.ID, "agent", o.Agent)

	if err := openStream(o, story, now, &res); err != nil {
		return res, &StepError{Story: story.ID, Step: StepStream, Finish: "flai stream open " + story.ID + ", then flai prime --story " + story.ID, Err: err}
	}
	if res.Pack, err = ctxpack.ForStory(o.Repo, story.ID, o.Budget); err != nil {
		return res, &StepError{Story: story.ID, Step: StepPrime, Finish: "flai prime --story " + story.ID, Err: err}
	}
	in, err := inbox.Read(ctx, inbox.Options{Repo: o.Repo, Agent: o.Agent, Own: story.ID, Now: o.Now, Runner: o.Runner, Version: o.Version})
	if err != nil {
		return res, &StepError{Story: story.ID, Step: StepInbox, Finish: "the MCP tool inbox", Err: fmt.Errorf("read %s's inbox: %w", o.Agent, err)}
	}
	res.Inbox = &in
	return res, nil
}

// refuse says why story may not start, or nil: it is not a ready story, or
// the board holds it, as inbox and wait_for_work judge it (S-0128).
func refuse(repo *workitem.Repo, story *workitem.Item) error {
	r := &Refused{Story: story.ID, Status: story.Status}
	switch {
	case story.Type != workitem.Story:
		r.Reason = fmt.Sprintf("%s is not a story but of type %s: flai story start starts a ready story", story.ID, story.Type)
	case story.Status == workitem.InProgress:
		r.Reason = fmt.Sprintf("%s is in progress already, so it was not moved: take it up with flai stream open %s and flai prime --story %s", story.ID, story.ID, story.ID)
	case story.Status != workitem.Ready:
		r.Reason = fmt.Sprintf("%s is %s, not ready: only a ready story starts; pull the first ready story the inbox lists instead", story.ID, story.Status)
	}
	if r.Reason != "" {
		return r
	}
	items, err := repo.List(true)
	if err != nil {
		return fmt.Errorf("start %s: read the work items to judge its holds: %w", story.ID, err)
	}
	if h := repo.Holds(items).Of(story); h != nil {
		r.Hold = h
		r.Reason = fmt.Sprintf("%s is %s; nothing was changed, so pull the first ready story the inbox lists that is not held", story.ID, h.Reason)
		return r
	}
	return nil
}

// openStream opens story's narrative, or takes up the one it has when it was
// begun before and sent back, writes the stream index again, and checks the
// story's branch out in its worktree, recording where in res.
func openStream(o Options, story *workitem.Item, now time.Time, res *Result) error {
	opt := workitem.StreamOptions{Agent: o.Agent, Session: o.Session, Host: o.Host, Now: now}
	_, err := o.Repo.OpenStream(story, opt)
	var exists *workitem.StreamExistsError
	reopen := errors.As(err, &exists)
	if reopen {
		_, err = o.Repo.ReopenStream(story, opt)
	}
	if err != nil {
		return fmt.Errorf("open %s's narrative: %w", story.ID, err)
	}
	items, err := o.Repo.List(false)
	if err != nil {
		return fmt.Errorf("list the work items for the stream index: %w", err)
	}
	if err := o.Repo.WriteIndex(items, now); err != nil {
		return fmt.Errorf("write the stream index: %w", err)
	}
	opened, err := storygit.Open(storygit.OpenOptions{Runner: o.Runner, Repo: o.Repo, Story: story.ID, FromRemote: reopen, RelativePaths: o.RelativePaths, Log: o.Log})
	if err != nil {
		return fmt.Errorf("narrative opened but the branch was not: %w", err)
	}
	res.Worktree, res.Branch, res.From = opened.Worktree, opened.Branch, opened.From
	return nil
}
