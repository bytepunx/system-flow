package serve

// The component dev servers flai serve supervises (S-0438, ADR-0163): one
// per project, sub-project, and checkout, main's or a story's worktree
// (S-0442), started by components.start with the sub-project's own vite and
// flai's harness, a story's after its dependencies are installed, answered by
// components.status, which marks it used, stopped by components.stop, by
// components.idle without a components.status, a story's once its story
// leaves review or is no longer in progress or in review, and every one when
// flai serve stops. At most components.max_servers run or start per project.
// One that exits on its own is forgotten, and why is kept as its last error.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// DefaultComponentsIdle is components.idle when the operator has not set it.
const DefaultComponentsIdle = 30 * time.Minute

// DefaultComponentsMaxServers is components.max_servers when the operator
// has not set it: the dev servers that may run or start at once for one
// project, main's and every story's together.
const DefaultComponentsMaxServers = 3

// DefaultComponentsStoryEvery is how often Run looks whether a story whose
// worktree has a dev server has left review, when ComponentsOptions gives no
// StoryEvery.
const DefaultComponentsStoryEvery = 5 * time.Second

// ComponentsOptions configure NewComponents.
type ComponentsOptions struct {
	// Idle is the operator's components.idle, asked at every look; nil, or
	// an answer that is not positive, is DefaultComponentsIdle.
	Idle func() time.Duration
	// MaxServers is the operator's components.max_servers, asked at every
	// start; nil, or an answer that is not positive, is
	// DefaultComponentsMaxServers.
	MaxServers func() int
	// Every is how often Run looks for idle dev servers; a tenth of the idle
	// time, at most a minute and at least a second, when zero.
	Every time.Duration
	// StoryEvery is how often Run looks at the stories whose worktrees have
	// dev servers; DefaultComponentsStoryEvery when zero.
	StoryEvery time.Duration
	// Story reads the story id names in the project at root, for the stop
	// on leaving review; from the project's work items when nil. An ID no
	// item has is an error that wraps workitem.ErrNotFound.
	Story func(root, id string) (*workitem.Item, error)
	// Install runs the package manager that installs a story's worktree's
	// dependencies; components.System when nil.
	Install components.Runner
	// Now is the clock idleness is measured by; time.Now when nil.
	Now func() time.Time
	// Logger is flai serve's log; nothing is logged when nil.
	Logger *slog.Logger
	// Timeout is how long a dev server has to answer once started; zero is
	// components.DevServerTimeout.
	Timeout time.Duration
}

// Components supervises the component dev servers of every project flai
// serve serves. It is safe for concurrent use; NewComponents makes one.
type Components struct {
	o      ComponentsOptions
	ctx    context.Context // ends when Close is called, ending a start in flight
	cancel context.CancelFunc
	wg     sync.WaitGroup // each start in flight

	mu      sync.Mutex
	closed  bool
	servers map[devKey]*devServer
	lastErr map[devKey]string
}

// devKey names one sub-project of one project, by the project's root, in
// the main checkout when story is empty and in story's worktree otherwise.
type devKey struct{ root, sub, story string }

// String names the dev server's sub-project and checkout for an error.
func (k devKey) String() string {
	if k.story == "" {
		return "sub-project " + k.sub
	}
	return "sub-project " + k.sub + " in " + k.story + "'s worktree"
}

// owner names whose dev server it is in a list: main's ui, S-0441's ui.
func (k devKey) owner() string {
	if k.story == "" {
		return "main's " + k.sub
	}
	return k.story + "'s " + k.sub
}

// logger is log with the dev server's root, sub-project, and story.
func (k devKey) logger(log *slog.Logger) *slog.Logger {
	log = log.With("component", "components", "root", k.root, "project", k.sub)
	if k.story != "" {
		log = log.With("story", k.story)
	}
	return log
}

// devServer is one dev server, starting or running.
type devServer struct {
	ready   chan struct{} // closed once the start answered, s or err set
	s       *components.DevServer
	err     error
	started time.Time // when it was asked to start, by Components' clock
	used    time.Time // the last start or status; guarded by Components.mu
}

// NewComponents is a supervisor of no dev servers yet, with its options'
// defaults filled in.
func NewComponents(o ComponentsOptions) *Components {
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Logger == nil {
		o.Logger = slog.New(slog.DiscardHandler)
	}
	if o.Story == nil {
		o.Story = readStory
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Components{o: o, ctx: ctx, cancel: cancel, servers: map[devKey]*devServer{}, lastErr: map[devKey]string{}}
}

// readStory is the story id names among the work items of the project at
// root.
func readStory(root, id string) (*workitem.Item, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	return repo.Get(id)
}

// Idle is components.idle as it stands: a dev server without a Status for
// this long is stopped.
func (c *Components) Idle() time.Duration {
	if c.o.Idle != nil {
		if d := c.o.Idle(); d > 0 {
			return d
		}
	}
	return DefaultComponentsIdle
}

// MaxServers is components.max_servers as it stands: the most dev servers
// that run or start at once for one project.
func (c *Components) MaxServers() int {
	if c.o.MaxServers != nil {
		if n := c.o.MaxServers(); n > 0 {
			return n
		}
	}
	return DefaultComponentsMaxServers
}

// serversFull is the refusal of a start beyond components.max_servers. It
// wraps hostapi.ErrComponentServersFull, which the dashboard is answered as
// a rule.
type serversFull struct{ msg string }

func (e serversFull) Error() string { return e.msg }

func (e serversFull) Unwrap() error { return hostapi.ErrComponentServersFull }

// Start starts sub's dev server for the project at root, whose manifest
// names it project, and answers it once it answers: in the main checkout
// when story is empty, and otherwise in dir, story's worktree, after
// installing the sub-project's dependencies there when they are missing.
// While one runs or is starting for the same sub-project and checkout, it
// answers that one, so that two starts at once start one server. A start
// that would run more than components.max_servers for the project is
// refused, naming the ones running. ctx bounds only the wait: a start it
// gives up on goes on, and Status says when it answers.
func (c *Components) Start(ctx context.Context, root, project string, sub manifest.Project, story, dir string) (hostapi.ComponentServer, error) {
	k := devKey{root, sub.Name, story}
	if story == "" {
		dir = root
	}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return hostapi.ComponentServer{}, fmt.Errorf("cannot start the component dev server of %s: flai serve is stopping; start it again once flai serve runs", k)
	}
	d, already := c.servers[k]
	if !already {
		if err := c.full(k); err != nil {
			c.mu.Unlock()
			return hostapi.ComponentServer{}, err
		}
		d = &devServer{ready: make(chan struct{}), started: c.o.Now()}
		c.servers[k] = d
		c.wg.Add(1)
		go c.launch(k, d, project, sub, dir)
	}
	d.used = c.o.Now()
	c.mu.Unlock()
	select {
	case <-d.ready:
	case <-ctx.Done():
		return hostapi.ComponentServer{}, fmt.Errorf("the component dev server of %s is still starting, and the request ended first (%w); components.status says when it answers", k, ctx.Err())
	}
	if d.err != nil {
		return hostapi.ComponentServer{}, d.err
	}
	out := c.Status(root, sub.Name, story)
	out.AlreadyRunning = already && out.Running
	return out, nil
}

// full is the refusal of a new dev server for k when components.max_servers
// run or start for its project already, or nil. c.mu is held.
func (c *Components) full(k devKey) error {
	limit := c.MaxServers()
	var keys []devKey
	for o := range c.servers {
		if o.root == k.root {
			keys = append(keys, o)
		}
	}
	if len(keys) < limit {
		return nil
	}
	slices.SortFunc(keys, func(a, b devKey) int { return c.servers[a].started.Compare(c.servers[b].started) })
	names := make([]string, len(keys))
	for i, o := range keys {
		names[i] = o.owner()
	}
	servers := "servers"
	if len(keys) == 1 {
		servers = "server"
	}
	return serversFull{msg: fmt.Sprintf("cannot start %s component dev server: flai runs %d component dev %s for this project, components.max_servers: %s; stop one first, or raise components.max_servers on the host with flai serve components max-servers", k.owner(), len(keys), servers, strings.Join(names, ", "))}
}

// launch starts d in dir, keeps the outcome, and watches the server it
// started: Serve in the main checkout, ServeWorktree in a story's worktree.
func (c *Components) launch(k devKey, d *devServer, project string, sub manifest.Project, dir string) {
	defer c.wg.Done()
	log := k.logger(c.o.Logger)
	opts := components.ServeOptions{Root: dir, Project: project, Sub: sub, Timeout: c.o.Timeout}
	var s *components.DevServer
	var err error
	if k.story == "" {
		s, err = components.Serve(c.ctx, opts)
	} else {
		s, err = components.ServeWorktree(c.ctx, components.WorktreeOptions{ServeOptions: opts, Runner: c.o.Install, Logger: log})
	}
	c.mu.Lock()
	d.s, d.err = s, err
	if err != nil {
		if c.servers[k] == d {
			delete(c.servers, k)
		}
		c.lastErr[k] = err.Error()
	} else {
		delete(c.lastErr, k)
		go c.watch(k, d, log)
	}
	close(d.ready)
	c.mu.Unlock()
	if err != nil {
		log.Warn("component dev server not started", "err", err.Error())
		return
	}
	log.Info("component dev server started", "pid", s.PID(), "port", s.Port(), "url", s.URL())
}

// watch forgets d once its server has exited on its own, keeping why. Close
// does not wait for it: a server that could not be stopped never exits.
func (c *Components) watch(k devKey, d *devServer, log *slog.Logger) {
	<-d.s.Done()
	err := d.s.Err()
	if err == nil {
		return // Stop ended it
	}
	c.mu.Lock()
	if c.servers[k] == d {
		delete(c.servers, k)
	}
	c.lastErr[k] = err.Error()
	c.mu.Unlock()
	log.Warn("component dev server exited", "pid", d.s.PID(), "err", err.Error())
}

// Status is how the dev server of the sub-project named sub of the project
// at root stands, in the main checkout when story is empty and in story's
// worktree otherwise, and marks it used.
func (c *Components) Status(root, sub, story string) hostapi.ComponentServer {
	k := devKey{root, sub, story}
	c.mu.Lock()
	defer c.mu.Unlock()
	out := hostapi.ComponentServer{Project: sub, Story: story, LastError: c.lastErr[k]}
	d, ok := c.servers[k]
	if !ok {
		return out
	}
	d.used = c.o.Now()
	select {
	case <-d.ready:
	default:
		out.Starting = true
		return out
	}
	if d.s == nil {
		return out
	}
	out.Running, out.URL, out.Port, out.PID = true, d.s.URL(), d.s.Port(), d.s.PID()
	out.Since = d.s.Started().UTC().Format(time.RFC3339)
	return out
}

// Stop stops the dev server of the sub-project named sub of the project at
// root, in the main checkout when story is empty and in story's worktree
// otherwise, waiting for one that is starting, and says whether one ran.
func (c *Components) Stop(root, sub, story string) (hostapi.ComponentStopped, error) {
	k := devKey{root, sub, story}
	out := hostapi.ComponentStopped{Project: sub, Story: story}
	c.mu.Lock()
	d, ok := c.servers[k]
	c.mu.Unlock()
	if !ok {
		return out, nil
	}
	<-d.ready
	c.mu.Lock()
	if c.servers[k] == d {
		delete(c.servers, k)
	}
	c.mu.Unlock()
	if d.s == nil {
		return out, nil
	}
	out.Stopped, out.PID = true, d.s.PID()
	if err := d.s.Stop(); err != nil {
		return out, fmt.Errorf("cannot stop the component dev server of %s: %w", k, err)
	}
	k.logger(c.o.Logger).Info("component dev server stopped", "pid", out.PID, "reason", "asked")
	return out, nil
}

// Run looks for idle dev servers, and for the stories whose worktrees' dev
// servers are to stop, every so often, stopping each, until ctx ends; then
// it stops every one, as Close does, and returns what Close returns.
func (c *Components) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return c.Close()
		case <-time.After(min(c.every(), c.storyEvery())):
		}
		c.stopIdle()
		c.stopLeft()
	}
}

// every is how long Run waits between looks for idle dev servers.
func (c *Components) every() time.Duration {
	if c.o.Every > 0 {
		return c.o.Every
	}
	return min(max(c.Idle()/10, time.Second), time.Minute)
}

// storyEvery is how long Run waits between looks at the stories.
func (c *Components) storyEvery() time.Duration {
	if c.o.StoryEvery > 0 {
		return c.o.StoryEvery
	}
	return DefaultComponentsStoryEvery
}

// stopIdle stops each running dev server unused for components.idle.
func (c *Components) stopIdle() {
	idle, now := c.Idle(), c.o.Now()
	var stale []*devServer
	var keys []devKey
	c.mu.Lock()
	for k, d := range c.servers {
		select {
		case <-d.ready:
		default:
			continue // a server starting is not idle
		}
		if d.s != nil && now.Sub(d.used) >= idle {
			delete(c.servers, k)
			stale = append(stale, d)
			keys = append(keys, k)
		}
	}
	c.mu.Unlock()
	for i, d := range stale {
		log := keys[i].logger(c.o.Logger).With("pid", d.s.PID())
		if err := d.s.Stop(); err != nil {
			log.Warn("idle component dev server not stopped", "err", err.Error())
			continue
		}
		log.Info("component dev server stopped", "reason", "idle", "idle", idle.String())
	}
}

// storyServer is a running dev server of a story's worktree, as stopLeft
// looks at it.
type storyServer struct {
	k devKey
	d *devServer
}

// stopLeft stops each running dev server of a story's worktree whose story
// StoryLeft says has left. A story that cannot be read, for a reason other
// than its not being there, keeps its servers until the next look.
func (c *Components) stopLeft() {
	type story struct{ root, id string }
	byStory := map[story][]storyServer{}
	c.mu.Lock()
	for k, d := range c.servers {
		if k.story == "" {
			continue
		}
		select {
		case <-d.ready:
		default:
			continue // a server starting is looked at once it answers
		}
		if d.s != nil {
			s := story{k.root, k.story}
			byStory[s] = append(byStory[s], storyServer{k, d})
		}
	}
	c.mu.Unlock()
	for s, servers := range byStory {
		it, err := c.o.Story(s.root, s.id)
		if err != nil && !errors.Is(err, workitem.ErrNotFound) {
			c.o.Logger.Warn("story of a component dev server not read", "component", "components", "root", s.root, "story", s.id, "err", err.Error())
			continue
		}
		if err != nil {
			it = nil
		}
		status := "gone"
		if it != nil {
			status = it.Status
		}
		for _, sv := range servers {
			if !StoryLeft(it, sv.d.started) {
				continue
			}
			c.mu.Lock()
			mine := c.servers[sv.k] == sv.d
			if mine {
				delete(c.servers, sv.k)
			}
			c.mu.Unlock()
			if !mine {
				continue // stopped or replaced since
			}
			log := sv.k.logger(c.o.Logger).With("pid", sv.d.s.PID(), "status", status)
			if err := sv.d.s.Stop(); err != nil {
				log.Warn("component dev server of a story that left review not stopped", "err", err.Error())
				continue
			}
			log.Info("component dev server stopped", "reason", "story left review")
		}
	}
}

// StoryLeft says whether the dev server of the worktree of the story it,
// started at since, is to stop: it is, when the story is no longer on the
// board (it nil, or archived), is neither in progress nor in review, or moved
// out of review, to done, cancelled, or back to in progress, at or after
// since; a story that moved from in progress to review keeps it. A
// transition's time has whole seconds, so since is taken to its second.
func StoryLeft(it *workitem.Item, since time.Time) bool {
	if it == nil || it.Archived || it.Status != workitem.InProgress && it.Status != workitem.Review {
		return true
	}
	from := since.Truncate(time.Second)
	prev := ""
	for _, t := range it.Transitions {
		if prev == workitem.Review && t.To != workitem.Review {
			if at, err := time.Parse(time.RFC3339, t.At); err == nil && !at.Before(from) {
				return true
			}
		}
		prev = t.To
	}
	return false
}

// Close stops every dev server, ends any start in flight, refuses starts
// from then on, and waits until each has ended. It is safe to call more than
// once.
func (c *Components) Close() error {
	c.mu.Lock()
	c.closed = true
	all := make([]*devServer, 0, len(c.servers))
	for k, d := range c.servers {
		all = append(all, d)
		delete(c.servers, k)
	}
	c.mu.Unlock()
	c.cancel()
	var errs []error
	stopped := 0
	for _, d := range all {
		<-d.ready
		if d.s == nil {
			continue
		}
		if err := d.s.Stop(); err != nil {
			errs = append(errs, err)
			continue
		}
		stopped++
	}
	c.wg.Wait()
	if stopped > 0 {
		c.o.Logger.Info("component dev servers stopped", "component", "components", "stopped", stopped, "reason", "flai serve stopping")
	}
	return errors.Join(errs...)
}
