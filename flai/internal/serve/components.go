package serve

// The component dev servers flai serve supervises (S-0438, ADR-0163): one
// per project and sub-project, started by components.start with the
// sub-project's own vite and flai's harness, answered by components.status,
// which marks it used, stopped by components.stop, by components.idle
// without a components.status, and every one when flai serve stops. One
// that exits on its own is forgotten, and why is kept as its last error.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// DefaultComponentsIdle is components.idle when the operator has not set it.
const DefaultComponentsIdle = 30 * time.Minute

// ComponentsOptions configure NewComponents.
type ComponentsOptions struct {
	// Idle is the operator's components.idle, asked at every look; nil, or
	// an answer that is not positive, is DefaultComponentsIdle.
	Idle func() time.Duration
	// Every is how often Run looks for idle dev servers; a tenth of the idle
	// time, at most a minute and at least a second, when zero.
	Every time.Duration
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

// devKey names one sub-project of one project, by the project's root.
type devKey struct{ root, sub string }

// devServer is one dev server, starting or running.
type devServer struct {
	ready chan struct{} // closed once the start answered, s or err set
	s     *components.DevServer
	err   error
	used  time.Time // the last start or status; guarded by Components.mu
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
	ctx, cancel := context.WithCancel(context.Background())
	return &Components{o: o, ctx: ctx, cancel: cancel, servers: map[devKey]*devServer{}, lastErr: map[devKey]string{}}
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

// Start starts sub's dev server for the project at root, whose manifest
// names it project, and answers it once it answers. While one runs or is
// starting for the same sub-project, it answers that one, so that two starts
// at once start one server. ctx bounds only the wait: a start it gives up on
// goes on, and Status says when it answers.
func (c *Components) Start(ctx context.Context, root, project string, sub manifest.Project) (hostapi.ComponentServer, error) {
	k := devKey{root, sub.Name}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return hostapi.ComponentServer{}, fmt.Errorf("cannot start sub-project %s's component dev server: flai serve is stopping; start it again once flai serve runs", sub.Name)
	}
	d, already := c.servers[k]
	if !already {
		d = &devServer{ready: make(chan struct{})}
		c.servers[k] = d
		c.wg.Add(1)
		go c.launch(k, d, project, sub)
	}
	d.used = c.o.Now()
	c.mu.Unlock()
	select {
	case <-d.ready:
	case <-ctx.Done():
		return hostapi.ComponentServer{}, fmt.Errorf("sub-project %s's component dev server is still starting, and the request ended first (%w); components.status says when it answers", sub.Name, ctx.Err())
	}
	if d.err != nil {
		return hostapi.ComponentServer{}, d.err
	}
	out := c.Status(root, sub.Name)
	out.AlreadyRunning = already && out.Running
	return out, nil
}

// launch starts d, keeps the outcome, and watches the server it started.
func (c *Components) launch(k devKey, d *devServer, project string, sub manifest.Project) {
	defer c.wg.Done()
	s, err := components.Serve(c.ctx, components.ServeOptions{Root: k.root, Project: project, Sub: sub, Timeout: c.o.Timeout})
	log := c.o.Logger.With("component", "components", "root", k.root, "project", sub.Name)
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
// at root stands, and marks it used.
func (c *Components) Status(root, sub string) hostapi.ComponentServer {
	k := devKey{root, sub}
	c.mu.Lock()
	defer c.mu.Unlock()
	out := hostapi.ComponentServer{Project: sub, LastError: c.lastErr[k]}
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
// root, waiting for one that is starting, and says whether one ran.
func (c *Components) Stop(root, sub string) (hostapi.ComponentStopped, error) {
	k := devKey{root, sub}
	out := hostapi.ComponentStopped{Project: sub}
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
		return out, fmt.Errorf("cannot stop sub-project %s's component dev server: %w", sub, err)
	}
	c.o.Logger.Info("component dev server stopped", "component", "components", "root", root, "project", sub, "pid", out.PID, "reason", "asked")
	return out, nil
}

// Run looks for idle dev servers every so often, stopping each that had no
// Status for components.idle, until ctx ends; then it stops every one, as
// Close does, and returns what Close returns.
func (c *Components) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return c.Close()
		case <-time.After(c.every()):
		}
		c.stopIdle()
	}
}

// every is how long Run waits between looks.
func (c *Components) every() time.Duration {
	if c.o.Every > 0 {
		return c.o.Every
	}
	return min(max(c.Idle()/10, time.Second), time.Minute)
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
		log := c.o.Logger.With("component", "components", "root", keys[i].root, "project", keys[i].sub, "pid", d.s.PID())
		if err := d.s.Stop(); err != nil {
			log.Warn("idle component dev server not stopped", "err", err.Error())
			continue
		}
		log.Info("component dev server stopped", "reason", "idle", "idle", idle.String())
	}
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
