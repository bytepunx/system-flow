package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
// agent action is enabled for the project; the story entered ready after
// flai serve began to serve it, and no agent has been started for it since;
// the in-progress limit leaves room, counting an agent started for a story
// still in ready as a story in progress; and nobody else is attending the
// project.
//
// The story says which harness works it, with which model and options (its
// agent, S-0103); package harness turns that into a command, with the
// operator's program and arguments for the harness. A story that names no
// harness is started with the operator's command, when one is set.
//
// Attending is judged from what flai can know without asking anyone: an agent
// connected over MCP rewrites its cursor under .flai-cache/mcp at every look,
// and one holding wait_for_events looks again at least every five minutes; an
// agent at work writes its story's narrative under wip/agents. If the newest
// of those files is younger than the operator's attended_minutes (six by
// default), someone is attending, and they will see the ready story in their
// inbox. The generated index there does not count, nor do the cursor and the
// narrative of an agent flai serve started itself.
//
// Each command is run as it stands in the project's directory, never through
// a shell. The story's ID is flai's own reading of the files, checked against
// the form of an ID before it is put anywhere.

// AgentConfig is the operator's say about starting agents, read afresh at
// every look so that enabling, disabling, and a new command need no restart.
type AgentConfig struct {
	Enabled  bool
	Command  []string
	Name     string        // FLAI_AGENT is this and the story's ID; "agent" when empty
	Attended time.Duration // how recent a sign of an agent counts; 6m when zero
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
	// Outcome is set once it has ended; Why says what went wrong.
	Outcome string `json:"outcome,omitempty"`
	Why     string `json:"why,omitempty"`
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
	stale   map[string]bool // ready when flai serve began, and ready since
	first   bool            // the first look only learns what is ready
	waiting map[int]bool    // the PIDs of agents this launcher waits for
	again   chan struct{}   // an agent ended: look again
}

func newLauncher(o Options, e Entry) *launcher {
	return &launcher{dir: o.Dir, entry: e, config: o.Agent, record: o.Host.Record, now: o.Now, first: true, waiting: map[int]bool{}, again: make(chan struct{}, 1),
		log: func(msg string, args ...any) {
			o.Logger.Info(msg, append([]any{"component", "serve", "project", e.Key}, args...)...)
		}}
}

// readyStory is a story in ready, as the launcher needs it.
type readyStory struct {
	ID      string
	Agent   *manifest.Agent
	Entered time.Time
}

// readyStories are the ready stories in pull order, and how many more
// stories the in-progress limit leaves room for, -1 when there is none.
func readyStories(root string) (ready []readyStory, free int, err error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, 0, err
	}
	items, err := repo.List(false)
	if err != nil {
		return nil, 0, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, 0, err
	}
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	// Pull order and the counts never look at done or archived items, so
	// which ones stay visible there (S-0087) makes no difference here.
	view := workitem.NewBoardView(items, board, time.Now(), false, nil)
	for _, c := range view.ReadyInPullOrder() {
		if it := byID[c.ID]; it != nil && c.Type == workitem.Story && !c.Blocked {
			ready = append(ready, readyStory{ID: c.ID, Agent: it.Agent, Entered: it.EnteredAt()})
		}
	}
	free = -1
	if limit, ok := view.WIPLimits[workitem.InProgress]; ok && limit > 0 {
		free = max(0, limit-view.Counts[workitem.InProgress])
	}
	return ready, free, nil
}

// attended reports the newest sign of an agent, if it is recent enough,
// leaving out the files named in own: those of agents flai serve started.
func attended(root string, within time.Duration, now time.Time, own map[string]bool) (bool, string) {
	repo, err := workitem.Open(root)
	if err != nil {
		return false, ""
	}
	newest, which := time.Time{}, ""
	look := func(pattern string, skip string) {
		files, _ := filepath.Glob(pattern)
		for _, f := range files {
			if filepath.Base(f) == skip || own[f] {
				continue
			}
			if info, err := os.Stat(f); err == nil && info.ModTime().After(newest) {
				newest, which = info.ModTime(), f
			}
		}
	}
	look(filepath.Join(repo.CacheDir(), "mcp", "*.json"), "")
	look(filepath.Join(repo.AgentsDir(), "*.md"), "index.md")
	if which == "" || now.Sub(newest) > within {
		return false, ""
	}
	rel, err := filepath.Rel(root, which)
	if err != nil {
		rel = which
	}
	return true, fmt.Sprintf("%s was written %s ago", filepath.ToSlash(rel), now.Sub(newest).Round(time.Second))
}

// ownSigns are the files an agent flai serve started writes, which are no
// sign of anyone else attending: its MCP cursor and its story's narrative,
// while it runs and for a while after.
func ownSigns(root string, st AgentState, within time.Duration, now time.Time) map[string]bool {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil
	}
	out := map[string]bool{}
	for id, run := range st.Stories {
		ended, err := time.Parse(time.RFC3339, run.Ended)
		if run.live() || (err == nil && now.Sub(ended) <= within) {
			out[filepath.Join(repo.CacheDir(), "mcp", run.Agent+".json")] = true
			out[filepath.Join(repo.AgentsDir(), id+".md")] = true
		}
	}
	return out
}

// judge is how an agent that has ended left its story.
func judge(root, story string, exit *int) (outcome, why string) {
	code := "an exit code nobody saw"
	if exit != nil {
		code = fmt.Sprintf("exit %d", *exit)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); the project could not be read: %v", code, err)
	}
	it, err := repo.Get(story)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); %s could not be read: %v", code, story, err)
	}
	switch it.Status {
	case workitem.Review, workitem.Done:
		return OutcomeWorked, ""
	}
	why = fmt.Sprintf("ended (%s) with %s in %s", code, story, it.Status)
	for _, b := range it.Blocked {
		if b.Until == "" {
			why += ", blocked: " + b.Reason
		}
	}
	return OutcomeFailed, why
}

// look is called when the project's work items may have changed, when an
// agent this launcher started has ended, and now and then.
func (l *launcher) look(ctx context.Context, _ bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.config == nil {
		return
	}
	stories, free, err := readyStories(l.entry.Root)
	if err != nil {
		return
	}
	ready := map[string]bool{}
	for _, s := range stories {
		ready[s.ID] = true
	}
	// What was ready when flai serve started is not news: a restart of flai
	// serve must not start a session nobody asked for. A story stays stale
	// until it leaves ready.
	if l.first {
		l.first, l.stale = false, ready
		return
	}
	for id := range l.stale {
		if !ready[id] {
			delete(l.stale, id)
		}
	}
	l.settleOrphans()
	if len(stories) == 0 {
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = "" })
		return
	}
	cfg := l.config(l.entry.Root)
	if !cfg.Enabled {
		return
	}
	st := l.dir.AgentStates()[l.entry.Root]
	commandSet := cfg.host(harness.Command).Program != ""
	var todo, nothing []readyStory
	reserved := 0
	for _, s := range stories {
		run := st.Stories[s.ID]
		switch {
		case run.live():
			reserved++ // started, and not in progress yet
		case l.stale[s.ID]:
		case run != nil && !startedBefore(run, s.Entered):
			// tried since it entered ready: it waits to be moved again
		case (s.Agent == nil || s.Agent.Harness == "") && !commandSet:
			nothing = append(nothing, s) // nothing to start it with, which is no failure
		default:
			todo = append(todo, s)
		}
	}
	if len(todo) == 0 {
		why := ""
		if len(nothing) > 0 {
			why = fmt.Sprintf("%s %s no harness, and no command is set on the host", strings.Join(ids(nothing), ", "), oneOrMany(len(nothing), "names", "name"))
		}
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = why })
		return
	}
	wait := func(why string) {
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = why })
		l.log("agent not started", "story", todo[0].ID, "why", why)
	}
	within := cfg.Attended
	if within <= 0 {
		within = 6 * time.Minute
	}
	if yes, sign := attended(l.entry.Root, within, l.now(), ownSigns(l.entry.Root, st, within, l.now())); yes {
		wait("an agent is attending the project (" + sign + "); it sees the story in its inbox")
		return
	}
	for i, s := range todo {
		if free >= 0 && reserved >= free {
			wait(fmt.Sprintf("the in-progress limit leaves no room for %s", strings.Join(ids(todo[i:]), ", ")))
			return
		}
		if l.start(ctx, cfg, s) {
			reserved++
		}
	}
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = "" })
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
// for it: flai serve was restarted while they ran.
func (l *launcher) settleOrphans() {
	st := l.dir.AgentStates()[l.entry.Root]
	for _, run := range st.Stories {
		if !run.live() || l.waiting[run.PID] || Alive(run.PID) {
			continue
		}
		ended := *run
		ended.Ended = l.now().UTC().Format(time.RFC3339)
		ended.Outcome, ended.Why = judge(l.entry.Root, run.Story, nil)
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(&ended) })
		l.log("agent ended while flai serve was away", "story", run.Story, "pid", run.PID, "outcome", ended.Outcome)
	}
}

// start starts an agent for story, and says whether it did.
func (l *launcher) start(ctx context.Context, cfg AgentConfig, story readyStory) bool {
	now := l.now().UTC()
	name := cfg.Name
	if name == "" {
		name = "agent"
	}
	run := &AgentRun{Story: story.ID, Agent: name + "-" + story.ID, Started: now.Format(time.RFC3339)}
	if story.Agent != nil {
		run.Model = story.Agent.Model
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
	spec, err := adapter.Start(harness.Request{Story: story.ID, Root: l.entry.Root, Project: l.entry.Key, Agent: story.Agent, Name: run.Agent, Flai: cfg.Flai}, cfg.host(name))
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
	l.waiting[run.PID] = true
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.put(run) })
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) for %s as %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, run.PID, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("agent started", "story", story.ID, "harness", run.Harness, "command", run.Command, "pid", run.PID)
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
		ended.Outcome, ended.Why = judge(l.entry.Root, story.ID, &code)
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
}

// Activity is what each story's newest agent is doing. One that runs is
// waiting when it asked a question on its story that nobody has answered
// yet (an open thread whose last entry is its own) or its story is blocked.
func Activity(root string, st AgentState) map[string]StoryActivity {
	out := map[string]StoryActivity{}
	if len(st.Stories) == 0 {
		return out
	}
	repo, err := workitem.Open(root)
	if err != nil {
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
