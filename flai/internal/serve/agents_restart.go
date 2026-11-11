package serve

import (
	"context"
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// resume starts again each agent that ended waiting for an answer or a
// reply, once what answers it is there (answered), in the session it had,
// when its story is open: in progress or in review. Such a story is counted
// in the limit already, and it holds others rather than being held
// (ADR-0046), so neither holds it back. A story in ready is started again by
// the look, once nothing holds it and the limit has room, like any other;
// one in backlog waits until it is ready (S-0182). A blocked story's agent is
// started again only for the answer to the question it asked, not for a
// message or a conversation that closed, which its first inbox lists once it
// is (ADR-0154); a blocked story in ready is never started, so the look needs
// no such rule. The interval flai serve opened when the agent ended asking
// (blockAsked) is no such block: it is closed before the agent is started
// again, whatever answered it, and one an agent or the operator wrote is
// left open (ADR-0205). A story whose agent runs elsewhere is left alone,
// answered or not (ADR-0242).
func (l *launcher) resume(ctx context.Context, cfg AgentConfig, repo *workitem.Repo) {
	st := l.dir.AgentStates()[l.entry.Root]
	for id, run := range st.Stories {
		if !asked(run) {
			continue
		}
		answer := answered(repo, run)
		if answer == "" {
			continue
		}
		it, err := repo.Get(id)
		if err != nil || (it.Status != workitem.InProgress && it.Status != workitem.Review) || harness.Elsewhere(repo.AgentOf(it)) {
			continue
		}
		own := run.Thread != "" && askedOn(it) == run.Thread
		if it.IsBlocked() && !own && answer != run.Thread {
			continue
		}
		if own {
			// on a failure it is started all the same: better marked blocked
			// than left without its agent, and the next look releases it
			if ok, err := l.unblockAsked(repo, id, run.Thread); err != nil {
				l.warn("story not unblocked for its agent's answer", "story", id, "thread", run.Thread, "err", err.Error())
			} else if ok {
				l.log("story unblocked for its agent's answer", "story", id, "thread", run.Thread, "answer", answer)
			}
		}
		l.start(ctx, cfg, readyStory{ID: id, Agent: cfg.storyAgent(repo, it), Asked: run, Answer: answer})
	}
}

// startSentBack starts the agent of each story in progress and not blocked
// that flai sent back from review on a failing required check or a requested
// change, while the thread flai serve opened for it is open (ADR-0167): in
// its last run's session, told the thread, or, when that run has no session,
// in a new one, told to take the story up. It starts once per send back: only
// when the story's newest run started before flai sent it back and has ended.
// A story flai serve never ran an agent for is worked from a session of its
// own, whose inbox lists the thread, so nothing is started beside it. Like
// resume, it spends no automatic restart and does not wait for room under
// the in-progress limit, whose place the story kept, and it leaves alone a
// run the operator stopped, which waits for their retry, a story that
// nothing can start, and one whose agent runs elsewhere (ADR-0242).
func (l *launcher) startSentBack(ctx context.Context, cfg AgentConfig, repo *workitem.Repo) {
	st := l.dir.AgentStates()[l.entry.Root]
	// only a story whose pull request flai serve has read can have been sent back
	files, _ := filepath.Glob(forge.StatusPath(repo.CacheDir(), "*"))
	for _, f := range files {
		it, err := repo.Get(strings.TrimSuffix(filepath.Base(f), ".json"))
		if err != nil || it.Type != workitem.Story || it.Status != workitem.InProgress || it.IsBlocked() {
			continue
		}
		run := st.Stories[it.ID]
		if run == nil || run.live() || run.Stopped != "" {
			continue
		}
		b, what, ok := sentBackFor(repo, it.ID)
		if !ok {
			continue
		}
		if at, err := time.Parse(time.RFC3339, b.At); err != nil || !startedBefore(run, at) {
			continue
		}
		agent := cfg.storyAgent(repo, it) // the default under it, ADR-0251, and the built-in, ADR-0261
		if cfg.refuseHarness(it.ID, agent) != nil {
			continue
		}
		s := readyStory{ID: it.ID, Agent: agent, SentBack: &harness.SentBack{Thread: b.Thread, What: what}}
		if run.Session != "" {
			s.Resume, s.SentBack.Resumed = run, true
		}
		if l.start(ctx, cfg, s) {
			l.log("agent started for a send back", "story", it.ID, "thread", b.Thread, "resumed", s.SentBack.Resumed)
		}
	}
	l.startQueueSentBack(ctx, cfg, repo, st)
	l.startAcceptSentBack(ctx, cfg, repo, st)
	l.startKlusterSentBack(ctx, cfg, repo, st)
}

// startKlusterSentBack starts the agent of each story flai serve's kluster
// pass sent back to in-progress because its deploy, or its checks run against
// its namespace, failed (ADR-0254 §5a), as startSentBack does a forge's: one
// the send back kept with what the pass tried still stands where it put it
// (klusterKept), not blocked, with the thread flai serve opened for it still
// open, whose last run started before the send back and has ended, told of
// the thread and what failed, in its last run's session where it has one.
// The send back names its thread, so no thread is matched by its title.
func (l *launcher) startKlusterSentBack(ctx context.Context, cfg AgentConfig, repo *workitem.Repo, st AgentState) {
	root := mainCheckout(repo)
	for _, ns := range slices.Sorted(maps.Keys(klusterTriedNamespaces(root))) {
		t, ok := readKlusterTried(root, ns)
		if !ok || t.SentBack == nil {
			continue
		}
		it, err := repo.Get(t.Story)
		if err != nil || it.Type != workitem.Story || !klusterKept(it, t) || it.IsBlocked() || !openThread(repo, t.SentBack.Thread) {
			continue
		}
		run := st.Stories[it.ID]
		if run == nil || run.live() || run.Stopped != "" || !startedBefore(run, t.SentBack.At) {
			continue
		}
		agent := cfg.storyAgent(repo, it) // the default under it, ADR-0251, and the built-in, ADR-0261
		if cfg.refuseHarness(it.ID, agent) != nil {
			continue
		}
		s := readyStory{ID: it.ID, Agent: agent, SentBack: &harness.SentBack{Thread: t.SentBack.Thread, What: t.SentBack.Reason, Kluster: true}}
		if run.Session != "" {
			s.Resume, s.SentBack.Resumed = run, true
		}
		if l.start(ctx, cfg, s) {
			l.log("agent started for a send back", "story", it.ID, "thread", t.SentBack.Thread, "resumed", s.SentBack.Resumed)
		}
	}
}

// startQueueSentBack starts the agent of each story the merge queue sent back
// to in-progress (ADR-0176), as startSentBack does a forge's: one whose last
// queue outcome is OutcomeFailing, not blocked, whose last run ended before
// the send back, with the thread flai serve opened on it still open, told of
// that thread and what failed.
func (l *launcher) startQueueSentBack(ctx context.Context, cfg AgentConfig, repo *workitem.Repo, st AgentState) {
	qs, err := mergequeue.ReadState(repo.CacheDir())
	if err != nil {
		return
	}
	for _, o := range qs.Outcomes {
		if o.Outcome != mergequeue.OutcomeFailing {
			continue
		}
		it, err := repo.Get(o.Story)
		if err != nil || it.Type != workitem.Story || it.Status != workitem.InProgress || it.IsBlocked() {
			continue
		}
		run := st.Stories[it.ID]
		if run == nil || run.live() || run.Stopped != "" || !startedBefore(run, o.At) {
			continue
		}
		thread, ok := queueSendBackThread(repo, it.ID)
		if !ok {
			continue
		}
		agent := cfg.storyAgent(repo, it) // the default under it, ADR-0251, and the built-in, ADR-0261
		if cfg.refuseHarness(it.ID, agent) != nil {
			continue
		}
		s := readyStory{ID: it.ID, Agent: agent, SentBack: &harness.SentBack{Thread: thread, What: strings.TrimPrefix(o.Detail, "sent back to in-progress: "), Queue: true}}
		if run.Session != "" {
			s.Resume, s.SentBack.Resumed = run, true
		}
		if l.start(ctx, cfg, s) {
			l.log("agent started for a send back", "story", it.ID, "thread", thread, "resumed", s.SentBack.Resumed)
		}
	}
}

// queueSendBackThread is the open thread the merge queue's send back opened
// on story, by flai serve, and false when there is none.
func queueSendBackThread(repo *workitem.Repo, story string) (string, bool) {
	all, err := threads.For(repo, story)
	if err != nil {
		return "", false
	}
	for _, th := range all {
		if th.Open() && th.Opener() == forgeAuthor && strings.HasPrefix(th.Title, "Merge queue batch ") {
			return th.ID, true
		}
	}
	return "", false
}

// startAcceptSentBack starts the agent of each story an acceptance, whoever
// made it, sent back to in-progress over a rebase it could not settle
// (ADR-0248), as startQueueSentBack does the queue's: one not blocked, with
// the thread flai accept opened on it still open, whose last run ended
// before that thread was opened, told of the thread and the paths it
// stopped on, in its last run's session where it has one.
func (l *launcher) startAcceptSentBack(ctx context.Context, cfg AgentConfig, repo *workitem.Repo, st AgentState) {
	for _, th := range acceptSentBackThreads(repo) {
		it, err := repo.Get(th.Anchor.Item)
		if err != nil || it.Type != workitem.Story || it.Status != workitem.InProgress || it.IsBlocked() {
			continue
		}
		at, err := time.Parse(time.RFC3339, th.Created)
		if err != nil {
			continue
		}
		run := st.Stories[it.ID]
		if run == nil || run.live() || run.Stopped != "" || !startedBefore(run, at) {
			continue
		}
		agent := cfg.storyAgent(repo, it) // the default under it, ADR-0251, and the built-in, ADR-0261
		if cfg.refuseHarness(it.ID, agent) != nil {
			continue
		}
		what := fmt.Sprintf("the acceptance of %s found that %s, which it cannot settle, and sent the story back; the thread says how to settle it in its worktree and sync it with flai stream sync %s",
			it.ID, strings.TrimPrefix(th.Title, storygit.SentBackPrefix), it.ID)
		s := readyStory{ID: it.ID, Agent: agent, SentBack: &harness.SentBack{Thread: th.ID, What: what, Acceptance: true}}
		if run.Session != "" {
			s.Resume, s.SentBack.Resumed = run, true
		}
		if l.start(ctx, cfg, s) {
			l.log("agent started for a send back", "story", it.ID, "thread", th.ID, "resumed", s.SentBack.Resumed)
		}
	}
}

// resolveAcceptSentBack resolves, as flai serve, each open thread flai
// accept opened on a story it sent back (ADR-0248) once the story is back in
// review, moved there after the thread was opened, so that the thread holds
// the next acceptance of it no longer. One that cannot be resolved is warned
// of and tried again at the next look.
func (l *launcher) resolveAcceptSentBack(repo *workitem.Repo) {
	for _, th := range acceptSentBackThreads(repo) {
		it, err := repo.Get(th.Anchor.Item)
		if err != nil || it.Status != workitem.Review {
			continue
		}
		at, err := time.Parse(time.RFC3339, th.Created)
		back := it.EnteredAt()
		if err != nil || back.Before(at) {
			continue
		}
		why := fmt.Sprintf("%s is back in review since %s, after the send back; its next acceptance takes it up again", it.ID, back.UTC().Format(time.RFC3339))
		if _, err := threads.Resolve(repo, th.ID, forgeAuthor, why, l.now()); err != nil {
			l.warn("send back thread not resolved", "story", it.ID, "thread", th.ID, "err", err.Error())
			continue
		}
		l.log("send back thread resolved", "story", it.ID, "thread", th.ID)
	}
}

// acceptSentBackThreads are the open threads flai accept opened on the
// stories an acceptance sent back (ADR-0248); none when the threads cannot
// be read, which, as resume does with what it cannot read,
// starts and resolves nothing.
func acceptSentBackThreads(repo *workitem.Repo) []*threads.Thread {
	all, err := threads.List(repo)
	if err != nil {
		return nil
	}
	var out []*threads.Thread
	for _, th := range all {
		if th.Open() && th.Anchor.Item != "" && th.Opener() == storygit.SentBackAuthor && strings.HasPrefix(th.Title, storygit.SentBackPrefix) {
			out = append(out, th)
		}
	}
	return out
}

// sentBackFor is story's latest send back whose thread is open, and what
// failed, as the move's reason names each failure the thread is for; false
// when there is none, or the kept status cannot be read, which, as resume
// does with what it cannot read, starts nothing.
func sentBackFor(repo *workitem.Repo, story string) (forge.SendBack, string, bool) {
	b, ok, err := OpenSendBack(repo, story)
	if err != nil || !ok {
		return forge.SendBack{}, "", false
	}
	st, _, err := forge.ReadStatus(repo.CacheDir(), story)
	if err != nil {
		return forge.SendBack{}, "", false
	}
	var backs []forge.SendBack
	for _, s := range st.SentBack {
		if s.Thread == b.Thread {
			backs = append(backs, s)
		}
	}
	return b, sendBackReason(st.PullRequest, backs), true
}

// restartEnded starts again, as flai serve agent restart does, the agent of
// each story in progress and not blocked whose newest run started and ended
// failed, neither asking nor stopped by the operator: in a new session, told
// how the last one ended and that flai serve started it on its own. It does
// so up to cfg.AutoRestarts times in a row; a run that ends at the limit gets
// one thread to the operator instead (S-0294, ADR-0108). A story that nothing
// can start, or whose agent runs elsewhere (ADR-0242), which flai serve agent
// restart refuses, is left alone, with no thread at the limit. One
// whose run ended on the harness's usage limit (afterLimit) is no such
// failure: it is started again once its provider's budget state no longer
// holds it, which the start's gate says, as the probe when the hold has just
// passed, spending no restart, its count as it was (S-0534). Every start
// here goes through that gate, so a story on a held provider waits where it
// is, and is counted nothing.
func (l *launcher) restartEnded(ctx context.Context, cfg AgentConfig, repo *workitem.Repo) {
	st := l.dir.AgentStates()[l.entry.Root]
	for _, id := range slices.Sorted(maps.Keys(st.Stories)) {
		run := st.Stories[id]
		again := afterLimit(run)
		if !again && !endedEarly(run) {
			continue
		}
		it, err := repo.Get(id)
		if err != nil || it.Status != workitem.InProgress || it.IsBlocked() {
			continue
		}
		agent := cfg.storyAgent(repo, it) // the default under it, ADR-0251, and the built-in, ADR-0261
		if cfg.refuseHarness(it.ID, agent) != nil {
			continue
		}
		if again {
			l.start(ctx, cfg, readyStory{ID: id, Agent: agent, AutoRestarts: run.AutoRestarts, AfterLimit: true, Restart: restartWhy(run)})
			continue
		}
		if run.AutoRestarts >= cfg.AutoRestarts {
			l.stopRestarting(repo, cfg, run)
			continue
		}
		n := run.AutoRestarts + 1
		if l.start(ctx, cfg, readyStory{ID: id, Agent: agent, AutoRestarts: n, Restart: restartWhy(run)}) {
			l.log("agent restarted on its own", "story", id, "restart", n, "limit", cfg.AutoRestarts)
		}
	}
}

// afterLimit says whether run is one whose story's agent flai serve starts
// again once its provider's budget state no longer holds it (S-0534): it
// started, and ended failed on the harness's usage limit, not asking,
// stopped by the operator, or queued for their retry.
func afterLimit(run *AgentRun) bool {
	return run != nil && run.Started != "" && !run.live() && run.Error == "" && run.Outcome == OutcomeFailed && run.Stopped == "" && run.Queued == "" && run.UsageLimited
}

// autoRestart is which automatic restart n is of limit, as the agent's
// prompt says it, such as "1 of 2"; empty for a start that is not one.
func autoRestart(n, limit int) string {
	if n == 0 {
		return ""
	}
	return fmt.Sprintf("%d of %d", n, limit)
}

// endedEarly says whether run is one flai serve restarts on its own: it
// started, ended failed rather than asking or stopped by the operator, not
// on the harness's usage limit, which its provider's budget state holds
// rather than a restart spent on it (S-0534), and has had no thread opened
// for it at the limit.
func endedEarly(run *AgentRun) bool {
	return run != nil && run.Started != "" && !run.live() && run.Error == "" && run.Outcome == OutcomeFailed && run.Stopped == "" && !run.UsageLimited && run.LimitThread == ""
}

// limitAuthor is who the thread flai serve opens at the limit of automatic
// restarts is by.
const limitAuthor = "flai serve"

// stopRestarting opens a thread to the operator on run's story, which ended
// once more than flai serve restarts it on its own, and records it on run so
// that it is opened once. When the thread cannot be opened, it warns once
// and tries again at the next look of a flai serve started since.
func (l *launcher) stopRestarting(repo *workitem.Repo, cfg AgentConfig, run *AgentRun) {
	key := run.Story + " " + run.Started
	if l.untold[key] {
		return
	}
	ends := run.AutoRestarts + 1
	times := fmt.Sprintf("%d times", ends)
	if ends == 1 {
		times = "once"
	}
	limit := "flai serve does not restart a story's agent on its own on this host (agent.auto_restarts is 0)"
	if cfg.AutoRestarts > 0 {
		limit = fmt.Sprintf("flai serve has restarted it on its own %d %s, as many as agent.auto_restarts allows on this host, and restarts it no more", cfg.AutoRestarts, oneOrMany(cfg.AutoRestarts, "time", "times"))
	}
	text := fmt.Sprintf("%[1]s's agent has ended %[2]s in a row with %[1]s in progress, and %[3]s: the story has no agent now. Its last run %[4]s.\n\n"+
		"%[5]sflai serve agent restart %[1]s, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.",
		run.Story, times, limit, restartWhy(run), countedRuns(repo, run.Story, ends))
	at := l.now().UTC()
	entry := hostapi.Entry{At: at.Format(time.RFC3339), Action: hostapi.ActionAgent, Method: "serve.agent", Project: l.entry.Key, Root: l.entry.Root, By: limitAuthor}
	th, err := threads.New(repo, threads.NewOptions{Title: fmt.Sprintf("%s's agent ended %s with the story in progress; restart it when it should go on", run.Story, times), On: run.Story, Author: limitAuthor, Text: text, Now: at})
	if err != nil {
		if l.untold == nil {
			l.untold = map[string]bool{}
		}
		l.untold[key] = true
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("could not open a thread on %s, whose agent ended %s in progress: %v", run.Story, times, err)
		if l.record != nil {
			l.record(entry)
		}
		l.warn("thread on the restart limit not opened", "story", run.Story, "err", err.Error())
		return
	}
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) {
		for _, r := range []*AgentRun{s.Stories[run.Story], s.Last} {
			if r.same(run) {
				r.LimitThread = th.ID
			}
		}
	})
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("opened %s: %s's agent ended %s in progress, and flai serve restarts it no more", th.ID, run.Story, times)
	if l.record != nil {
		l.record(entry)
	}
	l.log("agent not restarted: at the limit", "story", run.Story, "ended", ends, "thread", th.ID)
}

// countedRuns is the paragraph of the thread at the limit of automatic
// restarts that lists the agent events of story's last ends runs, the runs
// the thread counts, oldest first, from the ends-th last started event on,
// one list line each with its time, what the run came to, and why in a code
// span (S-0532): a started event and its end make two lines, so the operator
// reads each restart and how it ended. Fewer started events than ends, as a
// log that lost some has, list what there is; none leave the paragraph out,
// as an empty string, and the thread reads as it did before the event log.
func countedRuns(repo *workitem.Repo, story string, ends int) string {
	events := inbox.AgentEventsOf(repo, story, time.Time{})
	from, seen := -1, 0
	for i := len(events) - 1; i >= 0 && seen < ends; i-- {
		if events[i].To == inbox.AgentStarted {
			from, seen = i, seen+1
		}
	}
	if from < 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Its agent events since the first of those runs, oldest first:\n\n")
	for _, e := range events[from:] {
		b.WriteString("- " + e.At + " " + e.To)
		if why := strings.Join(strings.Fields(e.Fields[inbox.FieldWhy]), " "); why != "" {
			b.WriteString(": " + codeSpan(why))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

// asked says whether run ended asking: a question on a thread, or on a
// conversation with another story's agent (S-0335).
func asked(run *AgentRun) bool {
	return run != nil && run.Outcome == OutcomeAsked && (run.Thread != "" || len(run.Conversations) > 0)
}

// awaits is what a run that ended asking waits for, as the reason its story
// is not started says it: an answer to its question to the operator, the
// reply of the agent of each story it waits on, or both.
func (r *AgentRun) awaits() string {
	var on []string
	if r.Thread != "" {
		on = append(on, "an answer to "+r.Thread)
	}
	if len(r.Conversations) > 0 {
		who := "the other story's agent"
		switch {
		case len(r.WaitsOn) == 1:
			who = r.WaitsOn[0] + "'s agent"
		case len(r.WaitsOn) > 1:
			who = "the agents of " + strings.Join(r.WaitsOn, ", ")
		}
		on = append(on, who+" to reply on "+strings.Join(r.Conversations, ", "))
	}
	return "waiting for " + strings.Join(on, " and for ")
}

func oneOrMany(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
