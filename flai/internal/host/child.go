package host

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// child keeps one process running for as long as it is wanted: started
// again after an exit it was not asked for, after a wait that doubles while
// it keeps failing; stopped, gently and then not, when it is not wanted.
type child struct {
	name, root string
	spec       func() (Spec, error)
	h          *host

	mu        sync.Mutex
	want      bool // false once stopped for good
	paused    bool // stopped by the operator, kept to start again
	pid       int
	since     time.Time
	version   string
	external  int
	restarts  int
	lastExit  string
	lastError string

	kick chan struct{} // something changed: look again at once
	done chan struct{} // the loop has ended
}

func (c *child) poke() {
	select {
	case c.kick <- struct{}{}:
	default:
	}
}

// stop ends the child for good and waits for its process to go.
func (c *child) stop() {
	c.mu.Lock()
	c.want = false
	c.mu.Unlock()
	c.poke()
	<-c.done
}

// pause stops the process and keeps the child, to start again on resume.
func (c *child) pause() {
	c.mu.Lock()
	c.paused = true
	c.mu.Unlock()
	c.poke()
}

func (c *child) resume() {
	c.mu.Lock()
	c.paused = false
	c.mu.Unlock()
	c.poke()
}

// restartNow stops the process and starts it again without waiting.
func (c *child) restartNow() {
	c.mu.Lock()
	c.restarts++
	c.mu.Unlock()
	c.poke()
}

func (c *child) state() Child {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := Child{Name: c.name, Root: c.root, Restarts: c.restarts, LastExit: c.lastExit, LastError: c.lastError}
	switch {
	case (c.paused || !c.want) && c.pid != 0:
		// asked to stop, and its process has not ended yet: it is given the
		// grace period, and says stopped only once it is gone (S-0108)
		out.State, out.PID = "stopping", c.pid
	case c.paused || !c.want:
		out.State = "stopped"
	case c.external != 0:
		out.State, out.PID = "external", c.external
	case c.pid != 0:
		out.State, out.PID, out.Version = "running", c.pid, c.version
		out.Since = c.since.UTC().Format(time.RFC3339)
	default:
		out.State = "waiting"
	}
	return out
}

// wanted says whether the process should run now, and whether the child is
// still kept at all.
func (c *child) wanted() (run, kept bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.want && !c.paused, c.want
}

func (c *child) run(ctx context.Context) {
	defer close(c.done)
	o := c.h.o
	backoff := o.Backoff
	for {
		run, kept := c.wanted()
		if !kept || ctx.Err() != nil {
			return
		}
		if !run {
			if !c.wait(ctx, 0) {
				return
			}
			continue
		}
		spec, err := c.spec()
		switch {
		case err != nil:
			c.set(func() { c.lastError, c.external = err.Error(), 0 })
			o.Logger.Warn("child not started", "component", "host", "process", c.name, "root", c.root, "err", err.Error())
			if !c.wait(ctx, backoff) {
				return
			}
			backoff = min(backoff*2, o.MaxBackoff)
			continue
		case spec.External != 0:
			c.set(func() { c.external, c.lastError = spec.External, "" })
			if !c.wait(ctx, o.Look) {
				return
			}
			continue
		}
		c.set(func() { c.external = 0 })
		ran, asked := c.once(ctx, spec)
		if kept := c.isKept(); !kept || ctx.Err() != nil {
			return
		}
		if asked {
			backoff = o.Backoff
			continue
		}
		if ran > time.Minute {
			backoff = o.Backoff
		}
		c.set(func() { c.restarts++ })
		if !c.wait(ctx, backoff) {
			return
		}
		backoff = min(backoff*2, o.MaxBackoff)
	}
}

func (c *child) isKept() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.want
}

func (c *child) set(f func()) {
	c.mu.Lock()
	f()
	c.mu.Unlock()
}

// wait sleeps d, or until poked; zero waits for a poke. False when ctx ended.
func (c *child) wait(ctx context.Context, d time.Duration) bool {
	var after <-chan time.Time
	if d > 0 {
		t := time.NewTimer(d)
		defer t.Stop()
		after = t.C
	}
	select {
	case <-ctx.Done():
		return false
	case <-c.kick:
	case <-after:
	}
	return true
}

// once runs the process until it exits or is asked to stop; it says how long
// it ran and whether the end was asked for.
func (c *child) once(ctx context.Context, spec Spec) (time.Duration, bool) {
	o := c.h.o
	cmd := exec.Command(spec.Path, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = append(append(os.Environ(), spec.Env...), c.h.childEnv()...)
	if spec.Log != "" {
		f, err := os.OpenFile(spec.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			c.set(func() { c.lastError = err.Error() })
			return 0, false
		}
		defer func() { _ = f.Close() }()
		cmd.Stdout, cmd.Stderr = f, f
	}
	if err := cmd.Start(); err != nil {
		c.set(func() { c.lastError = err.Error() })
		o.Logger.Warn("child not started", "component", "host", "process", c.name, "root", c.root, "err", err.Error())
		return 0, false
	}
	start := time.Now()
	c.set(func() { c.pid, c.since, c.version, c.lastError = cmd.Process.Pid, start, spec.Version, "" })
	o.Logger.Info("child started", "component", "host", "process", c.name, "root", c.root, "pid", cmd.Process.Pid)
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	defer c.set(func() { c.pid = 0 })
	select {
	case err := <-exited:
		how := "exited 0"
		if err != nil {
			how = err.Error()
		}
		c.set(func() { c.lastExit = how })
		o.Logger.Warn("child exited", "component", "host", "process", c.name, "root", c.root, "pid", cmd.Process.Pid, "exit", how)
		return time.Since(start), false
	case <-ctx.Done():
	case <-c.kick:
		// a restart is a kick while wanted; a stop or a pause is a kick while not
	}
	c.end(cmd, exited)
	return time.Since(start), true
}

// end asks the process to stop, and kills it after the grace period.
func (c *child) end(cmd *exec.Cmd, exited <-chan error) {
	o := c.h.o
	_ = terminate(cmd.Process)
	select {
	case <-exited:
	case <-time.After(o.Grace):
		o.Logger.Warn("child killed", "component", "host", "process", c.name, "root", c.root, "pid", cmd.Process.Pid, "grace", o.Grace.String())
		_ = cmd.Process.Kill()
		<-exited
	}
	c.set(func() { c.lastExit = fmt.Sprintf("stopped by the host at %s", time.Now().UTC().Format(time.RFC3339)) })
}
