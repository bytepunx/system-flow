package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What wait_for_work answers with (S-0097).
const (
	WorkResume = "resume" // this agent's own story is still in progress
	WorkThread = "thread" // a thread awaiting this agent was written to while it waited
	WorkPull   = "pull"   // a story is ready, the in-progress limit leaves room for it, and review is under its limit

	WaitingForRoom   = "room"   // a story is ready, and the in-progress limit is full
	WaitingForReview = "review" // a story is ready, and review is full: it waits on acceptance (S-0243)
	WaitingForReady  = "ready"  // no story is ready
	WaitingForHeld   = "held"   // every ready story is held by an open story's claim or its after (S-0128, S-0130)
)

// WorkIn bounds one wait.
type WorkIn struct {
	TimeoutSeconds int `json:"timeout_seconds,omitempty" jsonschema:"how long to wait in seconds, default 300 (5 minutes), at most 1800 (30 minutes); call it again when it times out"`
}

// WorkOut is what an idle agent should do next, or why there is nothing yet.
type WorkOut struct {
	Reason     string               `json:"reason" jsonschema:"resume: your story is still in progress, go back to it; thread: a thread awaiting you was written to; pull: pull story now; empty when it timed out"`
	Story      *workitem.BoardCard  `json:"story,omitempty" jsonschema:"the story to resume or to pull"`
	Threads    []ThreadSummary      `json:"threads" jsonschema:"threads awaiting you that were written to since wait_for_work last answered"`
	Ready      []workitem.BoardCard `json:"ready" jsonschema:"stories ready to pull, in pull order; one an open story's claim holds, or that waits for a story it names in after, carries held with the reason and what clears it, and is not offered"`
	CanPull    bool                 `json:"can_pull" jsonschema:"whether the in-progress limit leaves room and review is under its limit, so that one may be pulled"`
	PullHold   string               `json:"pull_hold,omitempty" jsonschema:"why no story may be pulled now: the in-progress limit is full, or review is full and waits on acceptance"`
	WaitingFor string               `json:"waiting_for,omitempty" jsonschema:"when it timed out: room, when a story is ready but the in-progress limit is full; review, when a story is ready but review is full and waits on the operator's acceptance; held, when every ready story is held (ready says why each is); ready, when no story is ready"`
	TimedOut   bool                 `json:"timed_out"`
	// review is set when review being full, not the in-progress limit, is
	// what holds the pull (S-0243).
	review bool
}

// work says what an idle agent should do now, if anything: resume its own
// story, answer a thread written to since `since`, or pull the first ready
// story that is not held when the in-progress limit leaves room and review is
// under its limit (a held one keeps its place and is offered once clear,
// S-0128; a full review holds every one, S-0243). A thread that was already
// awaiting the agent before is not work here, so one the agent cannot answer
// does not wake it again and again; inbox lists every thread awaiting it.
func (s *server) work(ctx context.Context, since time.Time) (WorkOut, error) {
	view, err := s.boardView(ctx, false)
	if err != nil {
		return WorkOut{}, err
	}
	out := WorkOut{Ready: view.ReadyInPullOrder(), CanPull: view.CanPull(), PullHold: view.PullHold(), Threads: []ThreadSummary{}}
	// the hold is review's only when the in-progress limit has room
	out.review = out.PullHold != "" && out.PullHold == view.ReviewHold()
	if out.Ready == nil {
		out.Ready = []workitem.BoardCard{}
	}
	for _, c := range view.Columns[workitem.InProgress] {
		if c.Type == workitem.Story && s.mine(c.ID) {
			card := c
			out.Reason, out.Story = WorkResume, &card
			return out, nil
		}
	}
	done := perf.Track(ctx, "threads.read")
	all, err := threads.List(s.repo)
	done()
	if err != nil {
		return WorkOut{}, err
	}
	for _, th := range all {
		if !th.Open() {
			continue
		}
		sum := s.summary(th)
		updated, err := time.Parse(workitem.TimeFormat, th.Updated)
		if sum.Awaiting == "you" && err == nil && !updated.Before(since.Truncate(time.Second)) {
			out.Threads = append(out.Threads, sum)
		}
	}
	switch {
	case len(out.Threads) > 0:
		out.Reason = WorkThread
	case out.CanPull && workitem.FirstClear(out.Ready) != nil:
		out.Reason, out.Story = WorkPull, workitem.FirstClear(out.Ready)
	case out.CanPull && len(out.Ready) > 0:
		out.WaitingFor = WaitingForHeld
	case len(out.Ready) > 0 && out.review:
		out.WaitingFor = WaitingForReview
	case len(out.Ready) > 0:
		out.WaitingFor = WaitingForRoom
	default:
		out.WaitingFor = WaitingForReady
	}
	return out, nil
}

// mine says whether a story's narrative was last worked on under this
// agent's name, which is how a story in progress is known to be its own.
func (s *server) mine(id string) bool {
	n, err := workitem.ReadNarrative(s.repo.NarrativePath(id))
	return err == nil && n.Agent != "" && n.Agent == s.agent
}

func (s *server) waitForWork(ctx context.Context, req *mcp.CallToolRequest, in WorkIn) (*mcp.CallToolResult, WorkOut, error) {
	timeout := holdFor(in.TimeoutSeconds, workWait, s.maxWait)
	s.workMu.Lock()
	since := s.workSince
	s.workMu.Unlock()
	if since.IsZero() {
		since = s.now()
	}
	answer := func(out WorkOut) WorkOut {
		s.workMu.Lock()
		s.workSince = s.now()
		s.workMu.Unlock()
		return out
	}
	out, err := s.work(ctx, since)
	if err != nil {
		return nil, WorkOut{}, err
	}
	if out.Reason != "" {
		return nil, answer(out), nil
	}
	deadline := s.after(timeout)
	defer keepAlive(ctx, req, s.beat)()
	tick := time.NewTicker(s.poll)
	defer tick.Stop()
	before := s.snapshot()
	for {
		select {
		case <-ctx.Done():
			return nil, out, ctx.Err()
		case <-deadline:
			out.TimedOut = true
			return nil, answer(out), nil
		case <-s.closing:
			out.TimedOut = true
			return nil, answer(out), nil
		case <-tick.C:
			now := s.snapshot()
			if len(diff(before, now)) == 0 {
				continue
			}
			before = now
			if out, err = s.work(ctx, since); err != nil {
				return nil, WorkOut{}, err
			}
			if out.Reason != "" {
				return nil, answer(out), nil
			}
		}
	}
}
