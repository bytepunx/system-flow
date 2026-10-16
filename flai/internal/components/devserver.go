package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DevServerTimeout is how long StartDevServer waits for a dev server to
// answer when DevServerOptions gives no Timeout.
const DevServerTimeout = time.Minute

// tailLines is how many of a dev server's last lines of output it keeps for
// its errors.
const tailLines = 20

// maxLine is the longest line of output kept whole; a longer one is kept cut
// at this many bytes.
const maxLine = 4096

// pollEvery is how often StartDevServer asks a starting dev server whether
// it answers.
const pollEvery = 150 * time.Millisecond

// probeTimeout bounds one of those asks.
const probeTimeout = 2 * time.Second

// waitDelay is how long a dev server that has exited may leave its output
// open, held by something it started, before it is no longer read.
const waitDelay = 5 * time.Second

// stopGrace is how long Stop gives a dev server to end on its own, and then
// to end once killed. A variable so that a test can shorten it.
var stopGrace = 5 * time.Second

// ErrViteMissing is wrapped by the error StartDevServer returns when the
// sub-project has no vite in its node_modules.
var ErrViteMissing = errors.New("vite is not installed")

// DevServerOptions says how to start a sub-project's own Vite dev server.
type DevServerOptions struct {
	// Dir is the sub-project's folder, absolute; vite runs there, from its
	// node_modules/.bin.
	Dir string
	// Config is the Vite config to pass with --config, absolute; empty runs
	// vite with the sub-project's own config.
	Config string
	// Path is the path to wait for, such as "/__flai/components/"; it is the
	// path of URL. Empty is "/".
	Path string
	// Timeout is how long to wait for Path to answer; zero is
	// DevServerTimeout.
	Timeout time.Duration
	// Output, when set, also gets Vite's standard output and standard error,
	// as they come.
	Output io.Writer
}

// DevServer is a Vite dev server flai started, in a process group of its
// own. It runs until Stop ends it, or it exits on its own.
type DevServer struct {
	dir     string
	url     string
	port    int
	pid     int
	started time.Time
	cmd     *exec.Cmd
	out     *tail
	done    chan struct{}

	mu      sync.Mutex
	exited  bool // vite has exited
	stopped bool // Stop asked it to end while it ran

	stopOnce sync.Once
	stopErr  error
}

// StartDevServer starts the Vite dev server of the sub-project in
// opts.Dir, on 127.0.0.1 and a free port, and returns once an HTTP GET of its
// URL answers 200, as the harness page does. ctx bounds only that wait: once
// it answers, the server runs until Stop. A server that exits first, that
// does not answer 200 within the timeout, or whose wait ctx ends, is stopped,
// with nothing it started left behind, and the error says why, with its last
// lines of output. A sub-project without vite installed is refused with
// ErrViteMissing.
func StartDevServer(ctx context.Context, opts DevServerOptions) (*DevServer, error) {
	if !filepath.IsAbs(opts.Dir) {
		return nil, fmt.Errorf("cannot start a Vite dev server in %q: flai needs the sub-project's folder as an absolute path", opts.Dir)
	}
	if opts.Timeout <= 0 {
		opts.Timeout = DevServerTimeout
	}
	if !strings.HasPrefix(opts.Path, "/") {
		opts.Path = "/" + opts.Path
	}
	vite := filepath.Join(opts.Dir, "node_modules", ".bin", viteBin)
	if _, err := os.Stat(vite); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w in %s, and flai runs the sub-project's own vite: install the sub-project's dependencies, e.g. npm install or pnpm install in %s, then run this again", ErrViteMissing, opts.Dir, opts.Dir)
	} else if err != nil {
		return nil, fmt.Errorf("cannot look for vite at %s: %w; check that the sub-project's folder can be read", vite, err)
	}
	port, err := freePort()
	if err != nil {
		return nil, err
	}
	var args []string
	if opts.Config != "" {
		args = append(args, "--config", opts.Config)
	}
	args = append(args, "--host", "127.0.0.1", "--port", strconv.Itoa(port), "--strictPort")
	cmd := exec.Command(vite, args...)
	cmd.Dir = opts.Dir
	out := &tail{out: opts.Output}
	cmd.Stdout, cmd.Stderr = out, out
	cmd.WaitDelay = waitDelay
	ownGroup(cmd)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cannot start %s in %s: %w; check that it is executable, or install the sub-project's dependencies again", vite, opts.Dir, err)
	}
	s := &DevServer{
		dir:     opts.Dir,
		url:     "http://" + net.JoinHostPort("127.0.0.1", strconv.Itoa(port)) + opts.Path,
		port:    port,
		pid:     cmd.Process.Pid,
		started: time.Now(),
		cmd:     cmd,
		out:     out,
		done:    make(chan struct{}),
	}
	go s.wait()
	if err := s.await(ctx, opts.Timeout); err != nil {
		return nil, err
	}
	return s, nil
}

// URL is where the server answers: http://127.0.0.1:<port><Path>.
func (s *DevServer) URL() string { return s.url }

// Port is the loopback port the server listens on.
func (s *DevServer) Port() int { return s.port }

// PID is the process ID of vite, which leads the server's process group.
func (s *DevServer) PID() int { return s.pid }

// Started is when vite was started.
func (s *DevServer) Started() time.Time { return s.started }

// Done is closed once vite has exited and what it started has been ended.
func (s *DevServer) Done() <-chan struct{} { return s.done }

// Err says why the server exited when it exited on its own, with its last
// lines of output; it is nil while it runs and when Stop ended it.
func (s *DevServer) Err() error {
	select {
	case <-s.done:
	default:
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return nil
	}
	return fmt.Errorf("the Vite dev server in %s (pid %d) exited: %s; fix what its last lines say, then start it again%s", s.dir, s.pid, s.exitStatus(), detail(s.out.last()))
}

// Stop ends the server's whole process group, asking first and killing what
// has not ended within a grace period, and waits for vite to exit. It is
// safe to call more than once, and from more than one goroutine; each call
// returns the first's answer.
func (s *DevServer) Stop() error {
	s.stopOnce.Do(func() { s.stopErr = s.stop() })
	return s.stopErr
}

// stop is Stop, once.
func (s *DevServer) stop() error {
	s.mu.Lock()
	ran := !s.exited
	s.stopped = ran
	s.mu.Unlock()
	if !ran {
		<-s.done
		return nil
	}
	_ = terminateGroup(s.pid)
	timer := time.NewTimer(stopGrace)
	defer timer.Stop()
	select {
	case <-s.done:
		return nil
	case <-timer.C:
	}
	_ = killGroup(s.pid)
	timer.Reset(stopGrace)
	select {
	case <-s.done:
		return nil
	case <-timer.C:
	}
	return fmt.Errorf("cannot stop the Vite dev server in %s: pid %d had not exited %s after it was killed; end it by its PID", s.dir, s.pid, stopGrace)
}

// wait waits for vite to exit, ends what it started that outlives it, and
// closes done.
func (s *DevServer) wait() {
	_ = s.cmd.Wait()
	s.mu.Lock()
	s.exited = true
	s.mu.Unlock()
	_ = killGroup(s.pid)
	close(s.done)
}

// exitStatus says how vite exited, such as "exit status 1"; done is closed.
func (s *DevServer) exitStatus() string {
	if s.cmd.ProcessState == nil {
		return "it could not be waited for"
	}
	return s.cmd.ProcessState.String()
}

// await waits for the server's URL to answer 200, polling every pollEvery,
// and stops the server when it does not within timeout or ctx ends.
func (s *DevServer) await(ctx context.Context, timeout time.Duration) error {
	client := &http.Client{
		Transport: &http.Transport{Proxy: nil, DisableKeepAlives: true},
		Timeout:   probeTimeout,
	}
	defer client.CloseIdleConnections()
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	tick := time.NewTicker(pollEvery)
	defer tick.Stop()
	last := 0
	for {
		if status := probe(ctx, client, s.url); status == http.StatusOK {
			return nil
		} else if status != 0 {
			last = status
		}
		select {
		case <-s.done:
			return fmt.Errorf("the Vite dev server in %s exited before it answered at %s: %s; fix what its last lines say, then run this again%s", s.dir, s.url, s.exitStatus(), detail(s.out.last()))
		case <-ctx.Done():
			return errors.Join(fmt.Errorf("stopped the Vite dev server in %s before it answered at %s: %w", s.dir, s.url, ctx.Err()), s.Stop())
		case <-deadline.C:
			answered := ""
			if last != 0 {
				answered = fmt.Sprintf(" (it answered %d, and the harness answers 200)", last)
			}
			return errors.Join(fmt.Errorf("the Vite dev server in %s did not answer at %s within %s%s, so flai stopped it; fix what its last lines say, or allow it longer%s", s.dir, s.url, timeout, answered, detail(s.out.last())), s.Stop())
		case <-tick.C:
		}
	}
}

// probe is the status an HTTP GET of url answers, or 0 when it does not.
func probe(ctx context.Context, client *http.Client, url string) int {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return resp.StatusCode
}

// freePort is a loopback port nothing listens on, found by listening on
// port 0 and closing it. Another process may take it before vite does;
// --strictPort makes vite exit then, rather than take another.
func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("cannot find a free port on 127.0.0.1 for the Vite dev server: %w; check that loopback networking works", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, fmt.Errorf("cannot free port %d on 127.0.0.1 for the Vite dev server: %w", port, err)
	}
	return port, nil
}

// tail keeps the last tailLines lines written to it that are not blank, and
// copies everything written to out when it is set.
type tail struct {
	mu    sync.Mutex
	out   io.Writer
	lines []string
	part  []byte // the line being written, not ended yet
}

// Write keeps p's lines and copies p to out. It never fails, so that vite's
// output is always read; a failing out loses what it was given.
func (t *tail) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.out != nil {
		_, _ = t.out.Write(p)
	}
	t.part = append(t.part, p...)
	for {
		i := bytes.IndexByte(t.part, '\n')
		if i < 0 {
			break
		}
		t.keep(t.part[:i])
		t.part = t.part[i+1:]
	}
	if len(t.part) > maxLine {
		t.keep(t.part)
		t.part = nil
	}
	return len(p), nil
}

// keep adds line, cut at maxLine, unless it is blank, dropping the oldest
// past tailLines.
func (t *tail) keep(line []byte) {
	line = bytes.TrimRight(line[:min(len(line), maxLine)], "\r")
	if len(bytes.TrimSpace(line)) == 0 {
		return
	}
	t.lines = append(t.lines, string(line))
	if len(t.lines) > tailLines {
		t.lines = slices.Delete(t.lines, 0, len(t.lines)-tailLines)
	}
}

// last is the lines kept, and the one not ended yet, one per line.
func (t *tail) last() []byte {
	t.mu.Lock()
	defer t.mu.Unlock()
	lines := slices.Clone(t.lines)
	if len(bytes.TrimSpace(t.part)) > 0 {
		lines = append(lines, string(t.part))
	}
	return []byte(strings.Join(lines, "\n"))
}
