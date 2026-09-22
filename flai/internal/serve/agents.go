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
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Starting an agent when a story becomes ready (S-0079, ADR-0029).
//
// flai serve sees the project's files change. When a story has entered
// ready, it starts the operator's command, once, if all of this holds: the
// agent action is enabled for the project and a command is set; the
// in-progress limit leaves room for a pull; no agent it started for the
// project is still running; and nobody is attending the project.
//
// Attending is judged from what flai can know without asking anyone: an agent
// connected over MCP rewrites its cursor under .flai-cache/mcp at every look,
// and one holding wait_for_events looks again at least every five minutes; an
// agent at work writes its story's narrative under wip/agents. If the newest
// of those files is younger than the operator's attended_minutes (six by
// default), someone is attending, and they will see the ready story in their
// inbox. The generated index there does not count: flai rewrites it itself.
//
// The command is the operator's argument list, run as it stands in the
// project's directory, never through a shell. {story} and {root} in an
// argument are replaced; the story's ID is flai's own reading of the files,
// checked against the form of an ID before it is put anywhere.

// AgentConfig is the operator's say about starting agents, read afresh at
// every look so that enabling, disabling, and a new command need no restart.
type AgentConfig struct {
	Enabled  bool
	Command  []string
	Name     string        // FLAI_AGENT; "agent" when empty
	Attended time.Duration // how recent a sign of an agent counts; 6m when zero
}

// AgentRun is one agent flai serve started, or failed to.
type AgentRun struct {
	Story   string `json:"story"`
	Command string `json:"command"` // the program's name, not its arguments
	Agent   string `json:"agent"`
	PID     int    `json:"pid,omitempty"`
	Started string `json:"started"`
	Ended   string `json:"ended,omitempty"`
	Exit    *int   `json:"exit,omitempty"`
	Error   string `json:"error,omitempty"` // why it could not be started
	Log     string `json:"log,omitempty"`
}

// AgentState is what is known of agents for one project.
type AgentState struct {
	Running *AgentRun `json:"running,omitempty"`
	Last    *AgentRun `json:"last,omitempty"`
	// Waiting says why a ready story has not had an agent started for it.
	Waiting string `json:"waiting,omitempty"`
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

var storyID = regexp.MustCompile(`^S-\d{3,}$`)

// launcher starts agents for one project.
type launcher struct {
	dir    Dir
	entry  Entry
	config func(root string) AgentConfig
	record func(hostapi.Entry)
	now    func() time.Time
	log    func(msg string, args ...any)

	mu    sync.Mutex
	known map[string]bool // stories that were ready at the last look
	first bool            // the first look only learns what is ready
	again chan struct{}   // an agent ended: look again
}

func newLauncher(o Options, e Entry) *launcher {
	return &launcher{dir: o.Dir, entry: e, config: o.Agent, record: o.Host.Record, now: o.Now, first: true, again: make(chan struct{}, 1),
		log: func(msg string, args ...any) {
			o.Logger.Info(msg, append([]any{"component", "serve", "project", e.Key}, args...)...)
		}}
}

// readyStories are the ready stories in pull order, and whether the
// in-progress limit leaves room for one more.
func readyStories(root string) (ids []string, room bool, err error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, false, err
	}
	items, err := repo.List(false)
	if err != nil {
		return nil, false, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, false, err
	}
	// Pull order and CanPull never look at done or archived items, so which
	// ones stay visible there (S-0087) makes no difference here.
	view := workitem.NewBoardView(items, board, time.Now(), false, nil)
	for _, c := range view.ReadyInPullOrder() {
		if c.Type == workitem.Story && !c.Blocked {
			ids = append(ids, c.ID)
		}
	}
	return ids, view.CanPull(), nil
}

// attended reports the newest sign of an agent, if it is recent enough.
func attended(root string, within time.Duration, now time.Time) (bool, string) {
	repo, err := workitem.Open(root)
	if err != nil {
		return false, ""
	}
	newest, which := time.Time{}, ""
	look := func(pattern string, skip string) {
		files, _ := filepath.Glob(pattern)
		for _, f := range files {
			if filepath.Base(f) == skip {
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

// look is called when the project's work items may have changed, and when an
// agent this launcher started has ended.
func (l *launcher) look(ctx context.Context, ended bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.config == nil {
		return
	}
	ids, room, err := readyStories(l.entry.Root)
	if err != nil {
		return
	}
	entered := false
	ready := map[string]bool{}
	for _, id := range ids {
		ready[id] = true
		entered = entered || !l.known[id]
	}
	first := l.first
	l.known, l.first = ready, false
	// What was ready when flai serve started is not news: a restart of flai
	// serve must not start a session nobody asked for.
	if first || (!entered && !ended) || len(ids) == 0 {
		if len(ids) == 0 {
			l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = "" })
		}
		return
	}
	cfg := l.config(l.entry.Root)
	if !cfg.Enabled || len(cfg.Command) == 0 {
		return
	}
	wait := func(why string) {
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Waiting = why })
		l.log("agent not started", "story", ids[0], "why", why)
	}
	if st := l.dir.AgentStates()[l.entry.Root]; st.Running != nil && Alive(st.Running.PID) {
		wait(fmt.Sprintf("the agent started for %s is still running (pid %d); the next is started when it ends", st.Running.Story, st.Running.PID))
		return
	}
	if !room {
		wait("the in-progress limit leaves no room to pull a story")
		return
	}
	within := cfg.Attended
	if within <= 0 {
		within = 6 * time.Minute
	}
	if yes, sign := attended(l.entry.Root, within, l.now()); yes {
		wait("an agent is attending the project (" + sign + "); it sees the story in its inbox")
		return
	}
	l.start(ctx, cfg, ids[0])
}

func (l *launcher) start(ctx context.Context, cfg AgentConfig, story string) {
	now := l.now().UTC()
	run := &AgentRun{Story: story, Command: filepath.Base(cfg.Command[0]), Agent: cfg.Name, Started: now.Format(time.RFC3339)}
	if run.Agent == "" {
		run.Agent = "agent"
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionAgent, Method: "serve.agent", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) {
		run.Error = err.Error()
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Running, s.Last, s.Waiting = nil, run, "" })
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s for %s could not be started: %s", run.Command, story, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("agent could not be started", "story", story, "command", run.Command, "err", run.Error)
	}
	if !storyID.MatchString(story) {
		fail(fmt.Errorf("%q is not a story's ID", story))
		return
	}
	argv := Substitute(cfg.Command, story, l.entry.Root)
	logDir := filepath.Join(string(l.dir), "agents")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		fail(err)
		return
	}
	run.Log = filepath.Join(logDir, fmt.Sprintf("%s-%s-%s.log", l.entry.Key, story, now.Format("20060102T150405Z")))
	out, err := os.OpenFile(run.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		fail(err)
		return
	}
	// It outlives the look that started it and flai serve itself: a session
	// at work is not ended because its starter is restarted.
	cmd := exec.CommandContext(context.WithoutCancel(ctx), argv[0], argv[1:]...)
	cmd.Dir = l.entry.Root
	cmd.Stdout, cmd.Stderr = out, out
	cmd.Env = append(os.Environ(), "FLAI_AGENT="+run.Agent, "FLAI_STORY="+story, "FLAI_SESSION="+now.Format("20060102T150405"), "FLAI_STARTED_BY=flai-serve")
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		_ = out.Close()
		var ee *exec.Error
		if errors.As(err, &ee) {
			err = fmt.Errorf("%s: %w", ee.Name, ee.Err)
		}
		fail(err)
		return
	}
	run.PID = cmd.Process.Pid
	l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Running, s.Waiting = run, "" })
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s for %s as %s (pid %d); log %s", run.Command, story, run.Agent, run.PID, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("agent started", "story", story, "command", run.Command, "pid", run.PID)
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
		l.dir.updateAgent(l.entry.Root, func(s *AgentState) { s.Running, s.Last = nil, &ended })
		l.log("agent ended", "story", story, "pid", run.PID, "exit", code)
		select {
		case l.again <- struct{}{}:
		default:
		}
	}()
}
