package serve

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting agents when stories become ready (S-0079, S-0104, ADR-0029).
//
// flai serve sees the project's files change. It starts an agent for each
// story in ready that has none, in pull order, if all of this holds: the
// agent action is enabled for the project; no agent has been started for the
// story since it entered ready, whether that was before or after flai serve
// began to serve it (S-0112: serve/agents.json remembers each story's runs
// across a restart), or its agent has been changed since the last one was
// started (S-0116); and the in-progress limit leaves room, counting an agent
// started for a story still in ready as a story in progress. Each ready
// story it does not start, and has no agent running, is named in the state's
// waiting with the reason, and logged when its reason changes.
//
// The story says which harness works it, with which model and options (its
// agent, S-0103); package harness turns that into a command, with the
// operator's program and arguments for the harness. A story that names no
// harness is started with the operator's command, when one is set.
//
// Nobody attending holds a ready story back any more (S-0116, ADR-0043): the
// in-progress limit does, and so does a claim (S-0128, ADR-0046). A ready
// story whose claim overlaps the claim of a story in progress, in review, or
// in ready with an agent started for it is held: it is skipped, keeps its
// place, and is started first once it is clear. An agent holding wait_for_work and the
// launcher may both go for a story; the second move to in-progress is
// refused, and the loser pulls the next one.
//
// Each command is run as it stands in the project's directory, never through
// a shell. The story's ID is flai's own reading of the files, checked against
// the form of an ID before it is put anywhere.

// AgentConfig is the operator's say about starting agents, read afresh at
// every look so that enabling, disabling, and a new command need no restart.
type AgentConfig struct {
	Enabled bool
	Command []string
	Name    string // FLAI_AGENT is this and the story's ID; "agent" when empty
	// Harnesses are the program and arguments each harness runs with on this
	// host, the command among them when one is set (S-0104).
	Harnesses map[string]harness.Host
	// Flai is this flai's executable: the agents reach flai's MCP server
	// through it.
	Flai string
}

// host is what harness name runs with: the operator's, or for the command
// the command itself.
func (c AgentConfig) host(name string) harness.Host {
	if h, ok := c.Harnesses[name]; ok {
		return h
	}
	if name == harness.Command && len(c.Command) > 0 {
		return harness.Host{Program: c.Command[0], Args: c.Command[1:]}
	}
	return harness.Host{}
}

// Outcomes of an agent that has ended.
const (
	OutcomeWorked = "worked" // it left its story in review or done
	OutcomeAsked  = "asked"  // it ended waiting for an answer, and is started again when there is one
	OutcomeFailed = "failed" // it could not be started, or left its story anywhere else
)

// AgentRun is one agent flai serve started, or failed to.
type AgentRun struct {
	Story   string `json:"story"`
	Harness string `json:"harness,omitempty"`
	Model   string `json:"model,omitempty"`
	Command string `json:"command"` // the program's name, not its arguments
	Agent   string `json:"agent"`
	PID     int    `json:"pid,omitempty"`
	Started string `json:"started"`
	Ended   string `json:"ended,omitempty"`
	Exit    *int   `json:"exit,omitempty"`
	Error   string `json:"error,omitempty"` // why it could not be started
	Log     string `json:"log,omitempty"`
	// Session is the harness's session, which a start after an answer resumes.
	Session string `json:"session,omitempty"`
	// Answered is the question this run was started again for.
	Answered string `json:"answered,omitempty"`
	// Outcome is set once it has ended; Why says what went wrong, and Thread
	// is the question it ended waiting on.
	Outcome string `json:"outcome,omitempty"`
	Why     string `json:"why,omitempty"`
	Thread  string `json:"thread,omitempty"`
	// StoryAgent is the story's agent as it was when this run was started,
	// empty when it named none, and nil on runs recorded before S-0116: a
	// story whose agent has changed since is started again.
	StoryAgent *manifest.Agent `json:"story_agent,omitempty"`
	// Queued is when the operator asked for another agent after this one
	// ended, for its story in ready while the in-progress limit was full:
	// flai serve starts it as soon as there is room (S-0118).
	Queued string `json:"queued,omitempty"`
}

// agentChanged says whether the story's agent says something else than it
// did when run was started. A run that did not record it counts as unchanged.
func agentChanged(run *AgentRun, now *manifest.Agent) bool {
	return run.StoryAgent != nil && !run.StoryAgent.Same(now)
}

// live is a run that has not ended.
func (r *AgentRun) live() bool { return r != nil && r.Ended == "" && r.Error == "" && r.PID > 0 }

// AgentState is what is known of agents for one project.
type AgentState struct {
	// Running is the newest run that has not ended, and Last the newest that
	// has: what dashboards before S-0104 show.
	Running *AgentRun `json:"running,omitempty"`
	Last    *AgentRun `json:"last,omitempty"`
	// Waiting says why a ready story has not had an agent started for it.
	Waiting string `json:"waiting,omitempty"`
	// Stories are each story's newest run, by ID (S-0104).
	Stories map[string]*AgentRun `json:"stories,omitempty"`
}

// put records a run as its story's newest, and keeps Running and Last true.
func (s *AgentState) put(r *AgentRun) {
	if s.Stories == nil {
		s.Stories = map[string]*AgentRun{}
	}
	s.Stories[r.Story] = r
	s.Running = nil
	for _, run := range s.Stories {
		if run.live() && (s.Running == nil || run.Started > s.Running.Started) {
			s.Running = run
		}
	}
	if !r.live() {
		s.Last = r
	}
}

func (d Dir) agents() string { return filepath.Join(string(d), "agents.json") }

var agentsMu sync.Mutex

// AgentStates reads the state of every project's agents.
func (d Dir) AgentStates() map[string]AgentState {
	agentsMu.Lock()
	defer agentsMu.Unlock()
	return d.agentStates()
}

func (d Dir) agentStates() map[string]AgentState {
	out := map[string]AgentState{}
	if data, err := os.ReadFile(d.agents()); err == nil {
		_ = json.Unmarshal(data, &out)
	}
	return out
}

func (d Dir) updateAgent(root string, change func(*AgentState)) {
	agentsMu.Lock()
	defer agentsMu.Unlock()
	// flai serve agent restart writes here from a process of its own (S-0116).
	defer lockFile(d.agents() + ".lock")()
	all := d.agentStates()
	st := all[root]
	change(&st)
	all[root] = st
	_ = d.write(d.agents(), all)
}

// StoryID matches a story ID: S-, then three or more digits.
var StoryID = regexp.MustCompile(`^S-\d{3,}$`)

// launcher starts agents for one project.
type launcher struct {
	dir    Dir
	entry  Entry
	config func(root string) AgentConfig
	record func(hostapi.Entry)
	now    func() time.Time
	log    func(msg string, args ...any)

	mu      sync.Mutex
	said    map[string]string // why each skipped story waits, as last logged
	waiting map[int]bool      // the PIDs of agents this launcher waits for
	// handOver leaves each agent it starts to the serving flai, which settles
	// the run when its process is gone, instead of waiting for it: set for a
	// start on the operator's word, from a process that exits once it has
	// started the agent (S-0115).
	handOver bool
	again    chan struct{} // an agent ended: look again
}

func newLauncher(o Options, e Entry) *launcher {
	return &launcher{dir: o.Dir, entry: e, config: o.Agent, record: o.Host.Record, now: o.Now, waiting: map[int]bool{}, again: make(chan struct{}, 1),
		log: func(msg string, args ...any) {
			o.Logger.Info(msg, append([]any{"component", "serve", "project", e.Key}, args...)...)
		}}
}

// readyStory is a story in ready, as the launcher needs it.
type readyStory struct {
	ID      string
	Agent   *manifest.Agent
	Entered time.Time
	// Restart says how its last agent ended, when the operator has it
	// started again (S-0116); empty when it entered ready.
	Restart string
	item    *workitem.Item
}

// readyStories are the ready stories in pull order, the claims of the open
// stories they are held by, and how many more stories the in-progress limit
// leaves room for, -1 when there is none.
func readyStories(root string) (ready []readyStory, holds *workitem.Holds, free int, err error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, nil, 0, err
	}
	items, err := repo.List(false)
	if err != nil {
		return nil, nil, 0, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, nil, 0, err
	}
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	// Pull order and the counts never look at done or archived items, so
	// which ones stay visible there (S-0087) makes no difference here.
	view := workitem.NewBoardView(items, board, time.Now(), false, nil, repo.Manifest.Projects)
	for _, c := range view.ReadyInPullOrder() {
		if it := byID[c.ID]; it != nil && c.Type == workitem.Story && !c.Blocked {
			ready = append(ready, readyStory{ID: c.ID, Agent: it.Agent, Entered: it.EnteredAt(), item: it})
		}
	}
	free = -1
	if limit, ok := view.WIPLimits[workitem.InProgress]; ok && limit > 0 {
		free = max(0, limit-view.Counts[workitem.InProgress])
	}
	return ready, repo.Holds(items), free, nil
}

// agentStarted is how a hold names a story in ready whose agent has been
// started: it claims its paths from then, not once its agent pulls it.
const agentStarted = "its agent started"

// judge is how an agent that has ended left its story: worked, asked (a
// question of its own on the story is open and unanswered), or failed.
func judge(root, story, agent string, exit *int) (outcome, why, thread string) {
	code := "an exit code nobody saw"
	if exit != nil {
		code = fmt.Sprintf("exit %d", *exit)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); the project could not be read: %v", code, err), ""
	}
	it, err := repo.Get(story)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); %s could not be read: %v", code, story, err), ""
	}
	switch it.Status {
	case workitem.Review, workitem.Done:
		return OutcomeWorked, "", ""
	}
	if th := asking(repo, story, agent); th != nil {
		return OutcomeAsked, "waiting for an answer to " + th.ID + ": " + th.Title, th.ID
	}
	why = fmt.Sprintf("ended (%s) with %s in %s", code, story, it.Status)
	for _, b := range it.Blocked {
		if b.Until == "" {
			why += ", blocked: " + b.Reason
		}
	}
	return OutcomeFailed, why, ""
}

// look is called when the project's work items may have changed, when an
// agent this launcher started has ended, and now and then.
func (l *launcher) look(ctx context.Context, _ bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.config == nil {
		return
	}
	stories, holds, free, err := readyStories(l.entry.Root)
	if err != nil {
		return
	}
	// The first look after flai serve starts is like any other: what keeps a
	// restart from starting a story twice is its run in serve/agents.json,
	// not what happened to be ready when flai serve began (S-0112).
	l.settleOrphans()
	cfg := l.config(l.entry.Root)
	if cfg.Enabled {
		l.resume(ctx, cfg)
	}
	var sk skips
	defer func() { l.say(sk) }()
	if len(stories) == 0 {
		return
	}
	if !cfg.Enabled {
		sk.add("the agent host action is off for this project", stories...)
		return
	}
	st := l.dir.AgentStates()[l.entry.Root]
	commandSet := cfg.host(harness.Command).Program != ""
	var todo, nothing, noRoom []readyStory
	reserved := 0
	for _, s := range stories {
		run := st.Stories[s.ID]
		switch {
		case run.live():
			reserved++ // started, and not in progress yet: it has its agent
			holds.Open(s.item, agentStarted)
		case run != nil && run.Queued == "" && !startedBefore(run, s.Entered) && !agentChanged(run, s.Agent):
			sk.add(tried(s.ID, run), s)
		case (s.Agent == nil || s.Agent.Harness == "") && !commandSet:
			nothing = append(nothing, s) // nothing to start it with, which is no failure
		default:
			if run != nil && run.Queued != "" {
				s.Restart = restartWhy(run) // the operator's retry, waiting for room
			}
			todo = append(todo, s)
		}
	}
	if len(nothing) > 0 {
		sk.add(fmt.Sprintf("%s %s no harness, and no command is set on the host", strings.Join(ids(nothing), ", "), oneOrMany(len(nothing), "names", "name")), nothing...)
	}
	for _, s := range todo {
		switch h := holds.Of(s.item); {
		case h != nil:
			sk.add(s.ID+" "+h.Reason, s)
		case free >= 0 && reserved >= free:
			noRoom = append(noRoom, s)
		case l.start(ctx, cfg, s):
			reserved++
			holds.Open(s.item, agentStarted)
		}
	}
	if len(noRoom) > 0 {
		sk.add("the in-progress limit leaves no room for "+strings.Join(ids(noRoom), ", "), noRoom...)
	}
}

// skips are the ready stories a look did not start, and why, in the order
// said.
type skips struct {
	whys  []string
	story map[string]string
}

func (k *skips) add(why string, stories ...readyStory) {
	if k.story == nil {
		k.story = map[string]string{}
	}
	k.whys = append(k.whys, why)
	for _, s := range stories {
		k.story[s.ID] = why
	}
}

// say puts why each skipped story waits in the state's waiting, and logs a
// story's reason when it is not the kind last logged, so that a look every
// minute does not repeat it.
func (l *launcher) say(k skips) {
	why := strings.Join(k.whys, "; ")
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = why })
	for _, id := range slices.Sorted(maps.Keys(k.story)) {
		if _, said := l.said[id]; !said || kind(l.said[id]) != kind(k.story[id]) {
			l.log("agent not started", "story", id, "why", k.story[id])
		}
	}
	l.said = k.story
}

// kind is a reason without the detail in its parentheses: when a run
// started, or how it ended, is not a new reason.
func kind(why string) string {
	if i := strings.Index(why, " ("); i >= 0 {
		return why[:i]
	}
	return why
}

// tried says of a story that has had its agent since it entered ready what
// became of it, and what starts another.
func tried(id string, run *AgentRun) string {
	what := "started " + run.Started
	switch {
	case run.Error != "":
		what += ", could not start: " + run.Error
	case run.Why != "":
		what += ", " + run.Outcome + ": " + run.Why
	case run.Outcome != "":
		what += ", " + run.Outcome
	}
	return id + " has had its agent since it entered ready (" + what + "); it gets another when it enters ready again, its agent is changed, or it is restarted (flai serve agent restart)"
}

func ids(stories []readyStory) []string {
	out := make([]string, len(stories))
	for i, s := range stories {
		out[i] = s.ID
	}
	return out
}

// startedBefore says whether run began before the story entered ready.
func startedBefore(run *AgentRun, entered time.Time) bool {
	started, err := time.Parse(time.RFC3339, run.Started)
	return err == nil && started.Before(entered.Truncate(time.Second))
}

// settleOrphans ends the runs whose process is gone while no launcher waits
// for it: flai serve was restarted while they ran, or a command started them
// on the operator's word and handed them over (S-0115, S-0116).
func (l *launcher) settleOrphans() {
	st := l.dir.AgentStates()[l.entry.Root]
	for _, run := range st.Stories {
		if !run.live() || l.waiting[run.PID] || Alive(run.PID) {
			continue
		}
		ended := *run
		ended.Ended = l.now().UTC().Format(time.RFC3339)
		ended.Outcome, ended.Why, ended.Thread = judge(l.entry.Root, run.Story, run.Agent, nil)
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(&ended) })
		l.log("agent ended, seen at a look", "story", run.Story, "pid", run.PID, "outcome", ended.Outcome)
	}
}

// resume starts again each agent that ended waiting for an answer, once the
// answer is there, in the session it had. Its story is its own and already
// in progress, so neither the limit nor anyone attending holds it back.
func (l *launcher) resume(ctx context.Context, cfg AgentConfig) {
	st := l.dir.AgentStates()[l.entry.Root]
	repo, err := workitem.Open(l.entry.Root)
	if err != nil {
		return
	}
	for id, run := range st.Stories {
		if run.Outcome != OutcomeAsked || run.Thread == "" || !answered(repo, run) {
			continue
		}
		it, err := repo.Get(id)
		if err != nil || it.Closed() {
			continue
		}
		l.start(ctx, cfg, readyStory{ID: id, Agent: it.Agent}, run)
	}
}

// start starts an agent for story, or starts again the one that ended asking
// in after, and says whether it did.
func (l *launcher) start(ctx context.Context, cfg AgentConfig, story readyStory, after ...*AgentRun) bool {
	now := l.now().UTC()
	name := cfg.Name
	if name == "" {
		name = "agent"
	}
	run := &AgentRun{Story: story.ID, Agent: name + "-" + story.ID, Started: now.Format(time.RFC3339), Session: newSession()}
	if len(after) > 0 && after[0] != nil {
		run.Agent, run.Answered = after[0].Agent, after[0].Thread
		if after[0].Session != "" {
			run.Session = after[0].Session
		}
	}
	run.StoryAgent = &manifest.Agent{}
	if story.Agent != nil {
		run.Model = story.Agent.Model
		run.StoryAgent = story.Agent
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionAgent, Method: "serve.agent", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) bool {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(run) })
		what := run.Command
		if what == "" {
			what = "an agent"
		}
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s for %s could not be started: %s", what, story.ID, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("agent could not be started", "story", story.ID, "harness", run.Harness, "err", run.Error)
		return false
	}
	if !StoryID.MatchString(story.ID) {
		return fail(fmt.Errorf("%q is not a story's ID", story.ID))
	}
	name, adapter, err := harness.For(story.Agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	spec, err := adapter.Start(harness.Request{Story: story.ID, Root: l.entry.Root, Project: l.entry.Key, Agent: story.Agent, Name: run.Agent, Flai: cfg.Flai,
		Session: run.Session, Answered: run.Answered, Restart: story.Restart}, cfg.host(name))
	if err != nil {
		return fail(err)
	}
	argv := spec.Argv
	run.Command = filepath.Base(argv[0])
	logDir := filepath.Join(string(l.dir), "agents")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		return fail(err)
	}
	run.Log = filepath.Join(logDir, fmt.Sprintf("%s-%s-%s.log", l.entry.Key, story.ID, now.Format("20060102T150405Z")))
	out, err := os.OpenFile(run.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fail(err)
	}
	// It outlives the look that started it and flai serve itself: a session
	// at work is not ended because its starter is restarted.
	cmd := exec.CommandContext(context.WithoutCancel(ctx), argv[0], argv[1:]...)
	cmd.Dir = l.entry.Root
	cmd.Stdout, cmd.Stderr = out, out
	cmd.Env = append(os.Environ(), "FLAI_AGENT="+run.Agent, "FLAI_STORY="+story.ID, "FLAI_SESSION="+now.Format("20060102T150405"), "FLAI_STARTED_BY=flai-serve")
	cmd.Env = append(cmd.Env, spec.Env...)
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		_ = out.Close()
		var ee *exec.Error
		if errors.As(err, &ee) {
			err = fmt.Errorf("%s: %w", ee.Name, ee.Err)
		}
		return fail(err)
	}
	run.PID = cmd.Process.Pid
	if !l.handOver {
		l.waiting[run.PID] = true
	}
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(run) })
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) for %s as %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, run.PID, run.Log)
	if run.Answered != "" {
		entry.Detail = fmt.Sprintf("started %s (%s) again for %s as %s, %s answered (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, run.Answered, run.PID, run.Log)
	}
	if l.record != nil {
		l.record(entry)
	}
	l.log("agent started", "story", story.ID, "harness", run.Harness, "command", run.Command, "pid", run.PID)
	if l.handOver {
		_ = out.Close()
		_ = cmd.Process.Release()
		return true
	}
	go func() {
		err := cmd.Wait()
		_ = out.Close()
		code := 0
		var xe *exec.ExitError
		if errors.As(err, &xe) {
			code = xe.ExitCode()
		} else if err != nil {
			code = -1
		}
		ended := *run
		ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), &code
		ended.Outcome, ended.Why, ended.Thread = judge(l.entry.Root, story.ID, run.Agent, &code)
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(&ended) })
		l.mu.Lock()
		delete(l.waiting, run.PID)
		l.mu.Unlock()
		l.log("agent ended", "story", story.ID, "pid", run.PID, "exit", code, "outcome", ended.Outcome)
		select {
		case l.again <- struct{}{}:
		default:
		}
	}()
	return true
}

func oneOrMany(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// Activity states of a story's agent (S-0104), as the dashboard shows them.
const (
	ActivityWorking = "working" // running
	ActivityWaiting = "waiting" // running, and waiting for the designer
	ActivityFailed  = "failed"  // ended, or never started, without its story in review
	ActivityWorked  = "worked"  // ended with its story in review or done
)

// StoryActivity is what a story's agent is doing.
type StoryActivity struct {
	State string    `json:"state"`
	Why   string    `json:"why,omitempty"`
	Run   *AgentRun `json:"run"`
	// Thread is the question a waiting agent asked, when it waits on one.
	Thread string `json:"thread,omitempty"`
	// Hold is why a story in ready waits for another's claim (S-0128); Why
	// says the same.
	Hold *workitem.Hold `json:"hold,omitempty"`
}

// Activity is what each story's newest agent is doing. One that runs is
// waiting when it asked a question on its story that nobody has answered
// yet (an open thread whose last entry is its own) or its story is blocked.
// One that ended is waiting when it asked, or when the operator queued
// another for its story in ready (S-0118). A story in ready that a claim
// holds is waiting with the hold's reason, whether or not it has had an
// agent (S-0128): one that never had one is given a run that names only the
// story and the harness and model it asks for.
func Activity(root string, st AgentState) map[string]StoryActivity {
	repo, err := workitem.Open(root)
	if err != nil {
		return map[string]StoryActivity{}
	}
	out := runActivity(repo, st)
	for id, h := range held(repo, st) {
		run := st.Stories[id]
		if run == nil {
			run = &AgentRun{Story: id}
			if it, err := repo.Get(id); err == nil && it.Agent != nil {
				run.Harness, run.Model = it.Agent.Harness, it.Agent.Model
			}
		}
		out[id] = StoryActivity{State: ActivityWaiting, Why: h.Reason, Run: run, Hold: h}
	}
	return out
}

// runActivity is what each story's newest run is doing.
func runActivity(repo *workitem.Repo, st AgentState) map[string]StoryActivity {
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
			if e := th.Entries(); len(e) > 0 && e[len(e)-1].Author == run.Agent {
				asked[story] = th
			}
		}
	}
	for id, run := range st.Stories {
		a := StoryActivity{Run: run}
		switch {
		case run.live():
			a.State = ActivityWorking
			if th := asked[id]; th != nil {
				a.State, a.Thread, a.Why = ActivityWaiting, th.ID, "waiting for an answer to "+th.ID+": "+th.Title
			} else if it, err := repo.Get(id); err == nil {
				for _, b := range it.Blocked {
					if b.Until == "" {
						a.State, a.Why = ActivityWaiting, "blocked: "+b.Reason
					}
				}
			}
		case run.Outcome == OutcomeWorked:
			a.State = ActivityWorked
		case run.Outcome == OutcomeAsked:
			a.State, a.Thread, a.Why = ActivityWaiting, run.Thread, run.Why
		case run.Queued != "" && inReady(repo, id):
			a.State, a.Why = ActivityWaiting, "queued: flai serve starts another agent when the in-progress limit has room"
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
// counting the claim of each story in ready whose agent runs, as the
// launcher does. One whose agent ended asking is left out: it is started
// again when it is answered, held or not.
func held(repo *workitem.Repo, st AgentState) map[string]*workitem.Hold {
	items, err := repo.List(false)
	if err != nil {
		return nil
	}
	holds := repo.Holds(items)
	var ready []*workitem.Item
	for _, it := range items {
		if it.Type != workitem.Story || it.Status != workitem.Ready {
			continue
		}
		switch run := st.Stories[it.ID]; {
		case run.live():
			holds.Open(it, agentStarted)
		case run == nil || run.Outcome != OutcomeAsked:
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

// asking is the open thread on story whose last entry is agent's: a question
// it asked that nobody has answered yet.
func asking(repo *workitem.Repo, story, agent string) *threads.Thread {
	all, err := threads.List(repo)
	if err != nil {
		return nil
	}
	for _, th := range all {
		if !th.Open() || threads.StoryOf(repo, th) != story {
			continue
		}
		if e := th.Entries(); len(e) > 0 && e[len(e)-1].Author == agent {
			return th
		}
	}
	return nil
}

// answered says whether the question a run ended waiting on has an answer:
// an entry by someone else after the agent's, or the thread resolved.
func answered(repo *workitem.Repo, run *AgentRun) bool {
	th, err := threads.Get(repo, run.Thread)
	if err != nil {
		return false
	}
	if !th.Open() {
		return true
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
