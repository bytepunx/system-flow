package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand/v2"
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

// portTries is how many free ports StartDevServer starts vite on before it
// gives up, when each is taken before vite binds it.
const portTries = 3

// minPort and maxPort bound the ports freePort chooses from: above the
// ports of well-known and commonly run services, up to the last port.
const (
	minPort = 10000
	maxPort = 65535
)

// ephemeralRangeFile is where Linux says which ports it hands out to a
// listen on port 0 and to outgoing connections.
const ephemeralRangeFile = "/proc/sys/net/ipv4/ip_local_port_range"

// ianaEphemeralLo and ianaEphemeralHi are IANA's dynamic port range, which
// freePort avoids where ephemeralRangeFile cannot be read.
const (
	ianaEphemeralLo = 49152
	ianaEphemeralHi = 65535
)

// bindTries is how many ports freePort tries to listen on before it gives
// up, when each is in use.
const bindTries = 20

// portTaken is what vite's last lines say when its port was in use: Vite's
// "Port 5173 is already in use", or a bind's "address already in use".
const portTaken = "already in use"

// stopGrace is how long Stop gives a dev server to end on its own, and then
// to end once killed. A variable so that a test can shorten it.
var stopGrace = 5 * time.Second

// pickPort finds the port StartDevServer starts vite on. A variable so that a
// test can hand it a port that is taken.
var pickPort = freePort

// launched is told the pid of each vite StartDevServer starts. A variable so
// that a test can learn the pid of a vite stopped before StartDevServer
// returned it, which may never have run far enough to say its own.
var launched = func(int) {}

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
// it answers, the server runs until Stop. A server that exits first because
// its port was taken before it bound it is started again on another free
// port, portTries ports in all. A server that exits first otherwise, that
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
	vite, ok, err := findVite(opts.Dir)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w in %s, and flai runs the sub-project's own vite: install the sub-project's dependencies, e.g. npm install or pnpm install in %s, then run this again", ErrViteMissing, opts.Dir, opts.Dir)
	}
	var taken error
	for range portTries {
		port, err := pickPort()
		if err != nil {
			return nil, err
		}
		s, err := launch(vite, port, opts)
		if err != nil {
			return nil, err
		}
		err = s.await(ctx, opts.Timeout)
		if err == nil {
			return s, nil
		}
		if !s.lostPort() {
			return nil, err
		}
		taken = err
	}
	return nil, fmt.Errorf("the Vite dev server in %s found each of the %d free ports flai picked on 127.0.0.1 taken before it bound it; check what else on this host is opening ports, then run this again: %w", opts.Dir, portTries, taken)
}

// launch starts vite on port and returns the server, not yet waited for.
func launch(vite string, port int, opts DevServerOptions) (*DevServer, error) {
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
	launched(cmd.Process.Pid)
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
	return s, nil
}

// lostPort reports whether the server exited on its own because its port
// was in use, as its last lines say.
func (s *DevServer) lostPort() bool {
	if !closed(s.done) {
		return false
	}
	s.mu.Lock()
	stopped := s.stopped
	s.mu.Unlock()
	return !stopped && bytes.Contains(s.out.last(), []byte(portTaken))
}

// closed reports whether c is closed.
func closed(c <-chan struct{}) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
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

// findVite is the vite launcher of the sub-project in dir, in its
// node_modules/.bin, and whether it is there.
func findVite(dir string) (string, bool, error) {
	vite := filepath.Join(dir, "node_modules", ".bin", viteBin)
	if _, err := os.Stat(vite); errors.Is(err, fs.ErrNotExist) {
		return vite, false, nil
	} else if err != nil {
		return vite, false, fmt.Errorf("cannot look for vite at %s: %w; check that the sub-project's folder can be read", vite, err)
	}
	return vite, true, nil
}

// freePort is a loopback port nothing listens on, chosen at random from
// minPort to maxPort outside the host's ephemeral port range and checked by
// listening on it and closing it. The kernel hands out ports in that range
// to every listen on port 0 and every outgoing connection, so a port taken
// from it, as port 0 gives, can go to another process's test server the
// moment vite stops listening (I-0342); outside it, only a process that asks
// for the port by number takes it. One may still do so before vite binds
// it; --strictPort makes vite exit then, rather than take another, and
// StartDevServer picks another.
func freePort() (int, error) {
	lo, hi := ephemeralRange(os.ReadFile)
	var taken error
	for range bindTries {
		port, ok := portOutside(lo, hi, rand.IntN)
		if !ok {
			return 0, fmt.Errorf("cannot find a port for the Vite dev server: the host's ephemeral port range %d-%d covers every port from %d to %d; narrow it in %s, then run this again", lo, hi, minPort, maxPort, ephemeralRangeFile)
		}
		l, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
		if err != nil {
			taken = err
			continue
		}
		if err := l.Close(); err != nil {
			return 0, fmt.Errorf("cannot free port %d on 127.0.0.1 for the Vite dev server: %w", port, err)
		}
		return port, nil
	}
	return 0, fmt.Errorf("cannot find a free port on 127.0.0.1 for the Vite dev server: each of the %d ports flai tried from %d to %d, outside the ephemeral range %d-%d, could not be listened on; check that loopback networking works and what on this host holds so many ports, then run this again: %w", bindTries, minPort, maxPort, lo, hi, taken)
}

// ephemeralRange is the host's ephemeral port range, lo to hi, as read
// reads it from ephemeralRangeFile, or IANA's 49152-65535 where it cannot
// be read or does not hold one, as on hosts other than Linux.
func ephemeralRange(read func(string) ([]byte, error)) (lo, hi int) {
	data, err := read(ephemeralRangeFile)
	if err != nil {
		return ianaEphemeralLo, ianaEphemeralHi
	}
	lo, hi, ok := parsePortRange(string(data))
	if !ok {
		return ianaEphemeralLo, ianaEphemeralHi
	}
	return lo, hi
}

// parsePortRange reads a port range written as two numbers apart, lo then
// hi, as ip_local_port_range holds it, and whether s is one.
func parsePortRange(s string) (lo, hi int, ok bool) {
	f := strings.Fields(s)
	if len(f) != 2 {
		return 0, 0, false
	}
	lo, err := strconv.Atoi(f[0])
	if err != nil {
		return 0, 0, false
	}
	hi, err = strconv.Atoi(f[1])
	if err != nil || lo < 1 || lo > hi || hi > maxPort {
		return 0, 0, false
	}
	return lo, hi, true
}

// portOutside is a port from minPort to maxPort outside lo to hi, chosen
// with intN, which answers a number from 0 to n-1 as rand.IntN does, and
// whether there is one.
func portOutside(lo, hi int, intN func(n int) int) (int, bool) {
	below := max(0, min(lo, maxPort+1)-minPort)
	start := max(hi+1, minPort)
	above := max(0, maxPort-start+1)
	if below+above == 0 {
		return 0, false
	}
	k := intN(below + above)
	if k < below {
		return minPort + k, true
	}
	return start + k - below, true
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
