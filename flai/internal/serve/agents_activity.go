package serve

import (
	"context"
	"crypto/rand"
	"fmt"
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Activity states of a story's agent (S-0104), as the dashboard shows them.
const (
	ActivityWorking = "working" // running
	ActivityWaiting = "waiting" // waiting for the designer or another story's agent, running or not
	ActivityFailed  = "failed"  // ended, or never started, without its story in review
	ActivityWorked  = "worked"  // ended with its story in review or done
)

// StoryActivity is what a story's agent is doing.
type StoryActivity struct {
	State string    `json:"state"`
	Why   string    `json:"why,omitempty"`
	Run   *AgentRun `json:"run"`
	// Thread is the question to the operator a waiting agent asked, when it
	// waits on one.
	Thread string `json:"thread,omitempty"`
	// WaitsOn are the stories whose agents' replies a waiting agent waits on,
	// rather than or as well as the operator's answer, and Conversations the
	// conversations it waits on them in, or, once it has ended, those with a
	// message to its story it is started again to read (S-0335).
	WaitsOn       []string `json:"waits_on,omitempty"`
	Conversations []string `json:"conversations,omitempty"`
	// Hold is why a story in ready waits for another's claim (S-0128); Why
	// says the same.
	Hold *workitem.Hold `json:"hold,omitempty"`
	// Elsewhere is where a story in progress that this host has had no agent
	// for was begun (S-0177, ADR-0064); Why says the same.
	Elsewhere *Elsewhere `json:"elsewhere,omitempty"`
}

// Elsewhere is who began a story, when, and on which host, as its files say.
type Elsewhere struct {
	By    string `json:"by"`              // who last moved it to in-progress
	At    string `json:"at"`              // when
	Agent string `json:"agent,omitempty"` // the agent its narrative names
	Host  string `json:"host,omitempty"`  // the host its narrative was opened on
	Here  bool   `json:"here,omitempty"`  // that host is this one: begun outside flai serve
}

// Said is the line the dashboard shows for it.
func (e *Elsewhere) Said() string {
	who := e.Agent
	if who == "" {
		who = orSomeone(e.By)
	}
	where := "on another host"
	switch {
	case e.Here:
		where = "on this host, outside flai serve"
	case e.Host != "":
		where = "on " + e.Host
	}
	when := ""
	if e.At != "" {
		when = " at " + e.At
	}
	return "begun by " + who + " " + where + when + "; no agent here"
}

func orSomeone(who string) string {
	if who == "" {
		return "someone"
	}
	return who
}

// Activity is what each story's newest agent is doing. One that runs is
// waiting when it asked a question on its story that nobody has answered
// yet (an open thread whose last entry is its own, or one it opened with a
// recommendation awaiting the operator's confirmation), when a conversation
// of its story awaits another story's agent's reply, which WaitsOn names
// (S-0335), or when its story is blocked; the operator's question, when
// there is one, is its Thread either way.
// One that ended is waiting when it asked, or when the operator queued
// another for its story in ready (S-0118). A story in ready that a claim
// holds is waiting with the hold's reason, whether or not it has had an
// agent (S-0128): one that never had one is given a run that names only the
// story and the harness and model it asks for. A story in progress that
// this host has had no agent for is waiting, with where it was begun
// (S-0177), and the same kind of run.
//
// commandSet is whether a command is set on the host for the project: a
// stand-in's harness is then the story's and the project's, else flai's
// built-in default under them when they name none (ADR-0261).
//
// It times its parts as phases of the request ctx carries a timer for
// (I-0313): activity.conversations, the conversations read once; activity.runs,
// the runs judged; activity.holds, the holds built;
// activity.trial, each call of the trial judge among them; and
// activity.elsewhere, the stories begun elsewhere.
func Activity(ctx context.Context, root string, st AgentState, commandSet bool) map[string]StoryActivity {
	repo, err := workitem.Open(root)
	if err != nil {
		return map[string]StoryActivity{}
	}
	withTrial(repo, nil)
	return activityWith(ctx, repo, repo.Trial, st, commandSet)
}

// activityWith is Activity read from repo, whose overlaps judge, when not
// nil, judges by a trial merge.
func activityWith(ctx context.Context, repo *workitem.Repo, judge workitem.TrialJudge, st AgentState, commandSet bool) map[string]StoryActivity {
	repo.Trial = timedTrial(ctx, judge)
	done := perf.Track(ctx, "activity.conversations")
	convs := readConversations(repo)
	done()
	done = perf.Track(ctx, "activity.runs")
	out := runActivity(repo, st, convs)
	done()
	done = perf.Track(ctx, "activity.holds")
	holds := held(repo, st, convs)
	done()
	for id, h := range holds {
		out[id] = StoryActivity{State: ActivityWaiting, Why: h.Reason, Run: standIn(repo, st, id, "", commandSet), Hold: h}
	}
	done = perf.Track(ctx, "activity.elsewhere")
	begun := elsewhere(repo, st)
	done()
	for id, e := range begun {
		out[id] = StoryActivity{State: ActivityWaiting, Why: e.Said(), Run: standIn(repo, st, id, e.Agent, commandSet), Elsewhere: e}
	}
	return out
}

// timedTrial is judge with each call timed as the phase activity.trial of
// the request ctx carries a timer for, so that the phase sums them; nil when
// judge is, so that an overlap holds as it does without one.
func timedTrial(ctx context.Context, judge workitem.TrialJudge) workitem.TrialJudge {
	if judge == nil {
		return nil
	}
	return func(held, holder string) ([]string, error) {
		defer perf.Track(ctx, "activity.trial")()
		return judge(held, holder)
	}
}

// Unfinished is what agent.status answers of stories and plans without
// finished (S-0611): each story's activity whose run is live or whose state
// is working or waiting, and the newest planner run of each item that is
// live, with the newest run of each other story or planned item kept only
// while that item is open: listed outside the archive, neither done nor
// cancelled. So the board's failed dot and a story's agent panel keep their
// run, and the runs of finished work, which only the pages that list past
// runs show, are left out. A project whose items cannot be read keeps only
// the live and waiting ones.
func Unfinished(root string, stories map[string]StoryActivity, plans map[string]*AgentRun) (map[string]StoryActivity, map[string]*AgentRun) {
	open := map[string]bool{}
	if repo, err := workitem.Open(root); err == nil {
		if items, err := repo.List(false); err == nil {
			for _, it := range items {
				if it.Status != workitem.Done && it.Status != workitem.Cancelled {
					open[it.ID] = true
				}
			}
		}
	}
	keptStories := make(map[string]StoryActivity, len(stories))
	for id, a := range stories {
		if open[id] || a.Run.live() || a.State == ActivityWorking || a.State == ActivityWaiting {
			keptStories[id] = a
		}
	}
	keptPlans := make(map[string]*AgentRun, len(plans))
	for id, run := range plans {
		if open[id] || run.live() {
			keptPlans[id] = run
		}
	}
	return keptStories, keptPlans
}

// standIn is story's newest run, or, when it has none, one that names only
// the story, the agent given, and the harness and model the story will start
// with, the project's default's when its own names no harness (ADR-0251),
// and flai's built-in default's when neither does and no command is set
// (ADR-0261).
func standIn(repo *workitem.Repo, st AgentState, story, agent string, commandSet bool) *AgentRun {
	if run := st.Stories[story]; run != nil {
		return run
	}
	run := &AgentRun{Story: story, Agent: agent}
	if it, err := repo.Get(story); err == nil {
		if a := repo.StartAgentOf(it, commandSet); a != nil {
			run.Harness, run.Model = a.Harness, a.Model
		}
	}
	return run
}

// elsewhere are the stories in progress that this host has had no agent
// for, none started or only one queued (ADR-0064), with where each was
// begun.
func elsewhere(repo *workitem.Repo, st AgentState) map[string]*Elsewhere {
	items, err := repo.List(false)
	if err != nil {
		return nil
	}
	here := workitem.ThisHost()
	out := map[string]*Elsewhere{}
	for _, it := range items {
		if it.Type != workitem.Story || it.Status != workitem.InProgress {
			continue
		}
		if run := st.Stories[it.ID]; run != nil && run.Started != "" {
			continue
		}
		b := begunBy(repo, it)
		out[it.ID] = &Elsewhere{By: b.By, At: b.At, Agent: b.Agent, Host: b.Host, Here: b.Host != "" && b.Host == here}
	}
	return out
}

// runActivity is what each story's newest run is doing, judged from convs,
// the conversations read once for the call.
func runActivity(repo *workitem.Repo, st AgentState, convs *conversations) map[string]StoryActivity {
	out := map[string]StoryActivity{}
	if len(st.Stories) == 0 {
		return out
	}
	asked := map[string]*threads.Thread{} // story → the question its agent waits on
	if all, err := threads.List(repo); err == nil {
		for _, th := range all {
			story := threads.StoryOf(repo, th)
			run := st.Stories[story]
			if !th.Open() || !run.live() {
				continue
			}
			if awaitsAnswer(th, run.Agent) {
				asked[story] = th
			}
		}
	}
	for id, run := range st.Stories {
		a := StoryActivity{Run: run}
		switch {
		case run.live():
			a.State = ActivityWorking
			th, waits := asked[id], convs.awaitingOther(repo, id)
			if th != nil || len(waits) > 0 {
				a.State, a.Why = ActivityWaiting, waitWhy(th, nil, waits, nil, id)
				if th != nil {
					a.Thread = th.ID
				}
				a.WaitsOn, a.Conversations = others(waits, id), conversationIDs(waits)
			} else if it, err := repo.Get(id); err == nil {
				for _, b := range it.Blocked {
					if b.Until == "" {
						a.State, a.Why = ActivityWaiting, "blocked: "+b.Reason
					}
				}
			}
		case run.imageHeld():
			// its start waits for the image job that makes its agent images
			// (S-0751): waiting, not failed
			a.State, a.Why = ActivityWaiting, run.Why
		case run.Outcome == OutcomeWorked:
			a.State = ActivityWorked
		case run.Outcome == OutcomeAsked:
			a.State, a.Thread, a.Why = ActivityWaiting, run.Thread, run.Why
			a.WaitsOn, a.Conversations = run.WaitsOn, run.Conversations
		case run.Queued != "" && inReady(repo, id):
			a.State, a.Why = ActivityWaiting, "queued: flai serve starts another agent when there is room"
		default:
			a.State, a.Why = ActivityFailed, run.Why
			if a.Why == "" {
				a.Why = run.Error
			}
		}
		out[id] = a
	}
	return out
}

// held are the stories in ready with no agent running that a claim holds,
// counting the claim of each story in ready whose agent runs, and judging an
// overlap by a trial merge when repo has a judge, as the launcher does. One
// whose agent ended asking is left out until the question or its
// conversation is answered: it waits for the answer, and then for its hold
// (S-0182, S-0335). The shares, the asks, and the answers come from convs,
// the conversations read once for the call (S-0603).
func held(repo *workitem.Repo, st AgentState, convs *conversations) map[string]*workitem.Hold {
	items, err := repo.List(false)
	if err != nil {
		return nil
	}
	holds := repo.HoldsFrom(items, convs.shares, convs.asks)
	var ready []*workitem.Item
	for _, it := range items {
		if it.Type != workitem.Story || it.Status != workitem.Ready {
			continue
		}
		switch run := st.Stories[it.ID]; {
		case run.live():
			holds.Open(it, agentStarted)
		case !asked(run) || answeredFrom(repo, run, convs) != "":
			ready = append(ready, it)
		}
	}
	out := map[string]*workitem.Hold{}
	for _, it := range ready {
		if h := holds.Of(it); h != nil {
			out[it.ID] = h
		}
	}
	return out
}

// inReady says whether story is in ready: a retry queued for it waits only
// while it is.
func inReady(repo *workitem.Repo, story string) bool {
	it, err := repo.Get(story)
	return err == nil && it.Status == workitem.Ready
}

// asking is the open thread on story that awaits an answer to agent's
// question: one it asked that nobody has answered yet.
func asking(repo *workitem.Repo, story, agent string) *threads.Thread {
	all, err := threads.List(repo)
	if err != nil {
		return nil
	}
	for _, th := range all {
		if awaitsAnswer(th, agent) && threads.StoryOf(repo, th) == story {
			return th
		}
	}
	return nil
}

// awaitsAnswer says whether th is open and awaits an answer to agent's
// question: its last entry is agent's, or agent opened it and a
// recommendation on it awaits the operator's confirmation, which is no
// answer yet (ADR-0090).
func awaitsAnswer(th *threads.Thread, agent string) bool {
	if !th.Open() {
		return false
	}
	if e := th.Entries(); len(e) > 0 && e[len(e)-1].Author == agent {
		return true
	}
	return th.PendingRecommendation() != nil && th.Opener() == agent
}

// askedAnswered is the thread on run's story on which its agent wrote after
// the run started and whose question has an answer now, as threadAnswered
// reads one, or nil: an answer that came before the end was judged (I-0095).
// An entry from before the second the run started, or in it, is no such
// question: it can be the last run's, which was started again for its answer.
func askedAnswered(repo *workitem.Repo, run *AgentRun) *threads.Thread {
	started, err := time.Parse(time.RFC3339, run.Started)
	if err != nil {
		return nil
	}
	all, err := threads.List(repo)
	if err != nil {
		return nil
	}
	wrote := func(e threads.Entry) bool {
		at, err := time.Parse(time.RFC3339, e.At)
		return e.Author == run.Agent && err == nil && at.After(started)
	}
	for _, th := range all {
		e := th.Entries()
		if len(e) == 0 || e[len(e)-1].Author == run.Agent || th.PendingRecommendation() != nil {
			continue
		}
		if slices.ContainsFunc(e, wrote) && threads.StoryOf(repo, th) == run.Story {
			return th
		}
	}
	return nil
}

// answered is what answers a run that ended asking, or "" while nothing does:
// its thread, once the question has an answer (threadAnswered); a
// conversation it ended on, once that no longer awaits the other story: it
// reads as closed, so that no reply will come (S-0335), or the other story's
// agent has written last on it, a settled reply among what it writes, which
// leaves it awaiting neither story (ADR-0154); or a conversation on which a
// message to its story came after it ended (messaged), even while its
// question to the operator is unanswered. While that question is unanswered,
// a conversation the other story wrote only task-done notices on since the
// run ended answers nothing: the answer starts the agent, and it answers the
// notices then (ADR-0252).
func answered(repo *workitem.Repo, run *AgentRun) string {
	return answeredFrom(repo, run, nil)
}

// answeredFrom is answered judged from convs, the conversations read once
// already, or, when convs is nil, from each conversation read as it is
// needed.
func answeredFrom(repo *workitem.Repo, run *AgentRun, convs *conversations) string {
	asking := false
	if run.Thread != "" {
		if threadAnswered(repo, run) {
			return run.Thread
		}
		asking = true
	}
	for _, id := range run.Conversations {
		c, err := convs.get(repo, id)
		if err != nil {
			continue
		}
		if closed, _ := c.Closed(repo); closed || c.Awaiting() == "" {
			return c.ID
		}
		if wakes(c, run, asking) {
			return c.ID
		}
	}
	return messaged(repo, run, asking, convs)
}

// messaged is the first conversation of run's story, which ended asking,
// that awaits its reply to a message the other story's agent wrote after the
// run ended (sinceEnd), or "". It never looks before the run started: the
// message it was started again for is not new to it. A settled reply awaits
// no one, so it is passed over (ADR-0154); so is a conversation that gained
// only task-done notices since the end while the run is asking, its own
// question unanswered (ADR-0252).
func messaged(repo *workitem.Repo, run *AgentRun, asking bool, convs *conversations) string {
	since, ok := sinceEnd(run)
	if !ok {
		return ""
	}
	cs, err := convs.toSince(repo, run.Story, since)
	if err != nil {
		return ""
	}
	for _, c := range cs {
		if wakes(c, run, asking) {
			return c.ID
		}
	}
	return ""
}

// wakes reports whether c, which may await run's story, starts its ended run
// again: it awaits the story's reply, and, while the run is asking, its own
// question unanswered, the other story wrote more than task-done notices on
// it since the end (ADR-0252).
func wakes(c *messages.Conversation, run *AgentRun, asking bool) bool {
	return owesReply(c, run.Story) && (!asking || !onlyNoticesSinceEnd(c, run))
}

// sinceEnd is the time after which a message to run's story is new to its
// ended run: a second before the end, as times are kept to the second, so a
// message written in the second the run ended, after judgeRun read the
// conversations, is new to it too, while one its agent answered in that
// second no longer awaits its reply; but never before the run started. ok is
// false when the end cannot be read.
func sinceEnd(run *AgentRun) (time.Time, bool) {
	ended, err := time.Parse(time.RFC3339, run.Ended)
	if err != nil {
		return time.Time{}, false
	}
	since := ended.Add(-time.Second)
	if started, err := time.Parse(time.RFC3339, run.Started); err == nil && started.After(since) {
		since = started
	}
	return since, true
}

// onlyNoticesSinceEnd reports whether the other story's entries on c after
// sinceEnd are at least one and every one of them a task-done notice
// (ADR-0252), so that they need not start a run whose own question to the
// operator is unanswered. One entry that is not a notice says no.
func onlyNoticesSinceEnd(c *messages.Conversation, run *AgentRun) bool {
	since, ok := sinceEnd(run)
	if !ok {
		return false
	}
	other := workitem.CanonicalID(c.Other(run.Story))
	seen := false
	for _, e := range c.Entries() {
		at, err := time.Parse(workitem.TimeFormat, e.At)
		if err != nil || !at.After(since) || workitem.CanonicalID(e.Story) != other {
			continue
		}
		if !e.Notice {
			return false
		}
		seen = true
	}
	return seen
}

// threadAnswered says whether the question on the thread a run ended
// waiting on has an answer: an entry by someone else after the agent's that
// is not a recommendation awaiting the operator's confirmation, the
// confirmation of one, or the thread resolved.
func threadAnswered(repo *workitem.Repo, run *AgentRun) bool {
	th, err := threads.Get(repo, run.Thread)
	if err != nil {
		return false
	}
	if !th.Open() {
		return true
	}
	if th.PendingRecommendation() != nil {
		return false
	}
	e := th.Entries()
	return len(e) > 0 && e[len(e)-1].Author != run.Agent
}

// newSession is a random UUID, the form Claude Code's sessions take.
func newSession() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
