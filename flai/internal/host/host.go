// Package host is flai host: one process per machine that runs flai serve
// and each served project's HTTP MCP server as its children, keeps them
// running, and stops them when it stops (S-0106). flai serve says which
// projects need an MCP server; the host starts and stops them. A control API
// on a fixed loopback address, behind a bearer token, is how flai serve and
// the flai host command reach it, and binding that address is what makes
// the host one per machine.
package host

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// DefaultAddr is where the host listens unless AddrEnv says otherwise.
const DefaultAddr = "127.0.0.1:4241"

// The environment a child is started with: where its host is and how to prove
// it is one of its children, and the variable that moves the host's address.
const (
	URLEnv   = "FLAI_HOST_URL"
	TokenEnv = "FLAI_HOST_TOKEN"
	AddrEnv  = "FLAI_HOST_ADDR"
)

// Addr is the host's address: AddrEnv when set, else DefaultAddr.
func Addr() string {
	if v := os.Getenv(AddrEnv); v != "" {
		return v
	}
	return DefaultAddr
}

// The processes the host runs, by the names the control API takes.
const (
	Serve = "serve"
	MCP   = "mcp"
	All   = "all"
)

// ErrRestart is what Run returns when an upgrade installed a new flai: every
// child is stopped and the address is free, and the caller starts the new
// binary in its place.
var ErrRestart = errors.New("flai host restarts on the upgraded binary")

// Dir is the folder beside flai's config file that holds the host's state,
// token, and log.
type Dir string

// DirFor is the host's folder for a config file.
func DirFor(configPath string) Dir {
	return Dir(filepath.Join(filepath.Dir(configPath), "host"))
}

func (d Dir) state() string { return filepath.Join(string(d), "state.json") }

// TokenFile is where the running host keeps its token, mode 0600.
func (d Dir) TokenFile() string { return filepath.Join(string(d), "token") }

// Log is where a detached host writes.
func (d Dir) Log() string { return filepath.Join(string(d), "host.log") }

func (d Dir) write(path string, data []byte) error {
	if err := os.MkdirAll(string(d), 0o700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Status is what a running host last wrote about itself and its children.
type Status struct {
	PID      int     `json:"pid"`
	Version  string  `json:"version"`
	Started  string  `json:"started"`
	Updated  string  `json:"updated"`
	Addr     string  `json:"addr"`
	Config   string  `json:"config,omitempty"`
	Children []Child `json:"children"`
}

// Child is one process the host keeps, as it stands.
type Child struct {
	// Name is serve, or mcp for a project's MCP server, with Root.
	Name string `json:"name"`
	Root string `json:"root,omitempty"`
	// State is running; waiting, between an exit and the next start;
	// external, when a process the host did not start already does the job;
	// or stopped, when the operator stopped it.
	State string `json:"state"`
	PID   int    `json:"pid,omitempty"`
	Since string `json:"since,omitempty"`
	// Version is the flai it was started from.
	Version   string `json:"version,omitempty"`
	Restarts  int    `json:"restarts"`
	LastExit  string `json:"last_exit,omitempty"`
	LastError string `json:"last_error,omitempty"`
}

// StaleAfter is how old a status may be before its host counts as gone.
const StaleAfter = 10 * time.Second

// ReadStatus returns the last status and whether its host still runs.
func (d Dir) ReadStatus(now time.Time) (Status, bool) {
	var st Status
	data, err := os.ReadFile(d.state())
	if err != nil || json.Unmarshal(data, &st) != nil {
		return Status{}, false
	}
	updated, err := time.Parse(time.RFC3339, st.Updated)
	if err != nil || now.Sub(updated) > StaleAfter {
		return st, false
	}
	return st, alive(st.PID)
}

// Client is a client of the running host this folder records: its address,
// and the token it wrote.
func (d Dir) Client(now time.Time) (*Client, Status, error) {
	st, ok := d.ReadStatus(now)
	if !ok {
		return nil, st, errors.New("flai host is not running; start it with flai host start")
	}
	token, err := os.ReadFile(d.TokenFile())
	if err != nil {
		return nil, st, fmt.Errorf("flai host runs (pid %d) but its token is unreadable: %w", st.PID, err)
	}
	return &Client{URL: "http://" + st.Addr, Token: string(token)}, st, nil
}

// Spec is how to start one child. External, when not zero, is the PID of a
// process the host did not start that already does the job: nothing is
// started, and the host looks again later.
type Spec struct {
	Path     string
	Args     []string
	Dir      string
	Env      []string // added to the host's environment
	Log      string   // standard output and error are appended here
	Version  string
	External int
}

// Options configure Run.
type Options struct {
	Dir     Dir
	Addr    string // DefaultAddr when empty
	Version string
	Config  string // the config file the host was started with, for /_health
	Logger  *slog.Logger
	Now     func() time.Time
	// Serve is how to start flai serve; MCP how to start a project's MCP server.
	Serve func() (Spec, error)
	MCP   func(root string) (Spec, error)
	// Check says whether a newer flai is published; Upgrade installs it and
	// says whether it did (installed), so that the host restarts on it.
	Check   func(ctx context.Context) (any, error)
	Upgrade func(ctx context.Context) (res any, installed bool, err error)
	// Every is how often the state is written; Backoff the first wait after a
	// child ends unasked, doubled up to MaxBackoff; Grace how long a child is
	// given to stop before it is killed; Look how often an external process
	// is looked at again.
	Every, Backoff, MaxBackoff, Grace, Look time.Duration
}

func (o *Options) defaults() {
	if o.Addr == "" {
		o.Addr = DefaultAddr
	}
	if o.Logger == nil {
		o.Logger = slog.New(slog.DiscardHandler)
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Every <= 0 {
		o.Every = time.Second
	}
	if o.Backoff <= 0 {
		o.Backoff = time.Second
	}
	if o.MaxBackoff <= 0 {
		o.MaxBackoff = 30 * time.Second
	}
	if o.Grace <= 0 {
		o.Grace = 10 * time.Second
	}
	if o.Look <= 0 {
		o.Look = 15 * time.Second
	}
}

// Health is what /_health answers anyone on the machine, without the token:
// enough to say which host holds the address.
type Health struct {
	PID     int    `json:"pid"`
	Version string `json:"version"`
	Config  string `json:"config,omitempty"`
	Started string `json:"started"`
}

// AlreadyRunningError is Run's answer when the address is held by a host.
type AlreadyRunningError struct {
	Addr   string
	Health Health
}

func (e *AlreadyRunningError) Error() string {
	return fmt.Sprintf("flai host already runs on this machine at %s (pid %d, flai %s, config %s); one host serves the machine", e.Addr, e.Health.PID, e.Health.Version, e.Health.Config)
}

// Probe asks the address who holds it; false when no host answers there.
func Probe(addr string) (Health, bool) {
	c := &http.Client{Timeout: time.Second}
	resp, err := c.Get("http://" + addr + "/_health")
	if err != nil {
		return Health{}, false
	}
	defer func() { _ = resp.Body.Close() }()
	var h Health
	if resp.StatusCode != http.StatusOK || json.NewDecoder(resp.Body).Decode(&h) != nil || h.PID == 0 {
		return Health{}, false
	}
	return h, true
}

// host is a running host.
type host struct {
	o       Options
	token   string
	addr    string
	url     string
	started string
	ctx     context.Context

	mu       sync.Mutex
	serve    *child
	mcp      map[string]*child // by root
	serveOff bool              // the operator stopped serve
	mcpOff   bool              // the operator stopped the MCP servers
	restart  chan struct{}     // an upgrade installed a new flai
}

// Run is the host: it binds its address, starts flai serve, answers the
// control API, and when ctx ends, or an upgrade asks for a restart, stops
// every child before it returns.
func Run(ctx context.Context, o Options) error {
	o.defaults()
	ln, err := net.Listen("tcp", o.Addr)
	if err != nil {
		if h, ok := Probe(o.Addr); ok {
			return &AlreadyRunningError{Addr: o.Addr, Health: h}
		}
		return fmt.Errorf("flai host cannot listen on %s: %w; another program holds it, so set %s to another loopback address", o.Addr, err, AddrEnv)
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		_ = ln.Close()
		return err
	}
	h := &host{o: o, token: base64.RawURLEncoding.EncodeToString(buf), addr: ln.Addr().String(), url: "http://" + ln.Addr().String(),
		started: o.Now().UTC().Format(time.RFC3339), mcp: map[string]*child{}, restart: make(chan struct{}, 1)}
	if err := o.Dir.write(o.Dir.TokenFile(), []byte(h.token)); err != nil {
		_ = ln.Close()
		return err
	}
	cctx, stopChildren := context.WithCancel(context.WithoutCancel(ctx))
	h.ctx = cctx
	srv := &http.Server{Handler: h.api(), ReadHeaderTimeout: 5 * time.Second}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	o.Logger.Info("host started", "component", "host", "addr", ln.Addr().String(), "pid", os.Getpid(), "version", o.Version)

	h.mu.Lock()
	h.serve = h.newChild(Serve, "", o.Serve, false)
	h.mu.Unlock()
	h.writeStatus()

	tick := time.NewTicker(o.Every)
	defer tick.Stop()
	restart := false
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-h.restart:
			restart = true
			break loop
		case err := <-served:
			o.Logger.Error("control api stopped", "component", "host", "err", err.Error())
			break loop
		case <-tick.C:
			h.writeStatus()
		}
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = srv.Shutdown(shutdown)
	cancel()
	stopChildren()
	h.stopAll()
	_ = os.Remove(o.Dir.state())
	_ = os.Remove(o.Dir.TokenFile())
	o.Logger.Info("host stopped", "component", "host", "restart", restart)
	if restart {
		return ErrRestart
	}
	return nil
}

// childEnv is what every child is given: the way back to its host.
func (h *host) childEnv() []string {
	return []string{URLEnv + "=" + h.url, TokenEnv + "=" + h.token}
}

// newChild keeps a new child, paused when the operator stopped its kind.
func (h *host) newChild(name, root string, spec func() (Spec, error), paused bool) *child {
	c := &child{name: name, root: root, spec: spec, h: h, want: true, paused: paused, kick: make(chan struct{}, 1), done: make(chan struct{})}
	go c.run(h.ctx)
	return c
}

// children is every child, serve first, then the MCP servers by root.
func (h *host) children() []*child {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := []*child{}
	if h.serve != nil {
		out = append(out, h.serve)
	}
	roots := make([]string, 0, len(h.mcp))
	for r := range h.mcp {
		roots = append(roots, r)
	}
	sort.Strings(roots)
	for _, r := range roots {
		out = append(out, h.mcp[r])
	}
	return out
}

func (h *host) status() Status {
	st := Status{PID: os.Getpid(), Version: h.o.Version, Started: h.started, Updated: h.o.Now().UTC().Format(time.RFC3339), Addr: h.addr, Config: h.o.Config, Children: []Child{}}
	for _, c := range h.children() {
		st.Children = append(st.Children, c.state())
	}
	h.mu.Lock()
	if h.serveOff && h.serve == nil {
		st.Children = append([]Child{{Name: Serve, State: "stopped"}}, st.Children...)
	}
	h.mu.Unlock()
	return st
}

func (h *host) writeStatus() {
	data, err := json.MarshalIndent(h.status(), "", "  ")
	if err == nil {
		err = h.o.Dir.write(h.o.Dir.state(), append(data, '\n'))
	}
	if err != nil {
		h.o.Logger.Warn("status not written", "component", "host", "err", err.Error())
	}
}

// stopAll stops every child and waits for each.
func (h *host) stopAll() {
	for _, c := range h.children() {
		c.stop()
	}
}

// act starts, stops, or restarts serve, the MCP servers, or both.
func (h *host) act(process, action string) error {
	if process != Serve && process != MCP && process != All {
		return fmt.Errorf("%q is not a process of the host: serve, mcp, or all", process)
	}
	if action != "start" && action != "stop" && action != "restart" {
		return fmt.Errorf("%q is not an action: start, stop, or restart", action)
	}
	if process == Serve || process == All {
		h.actServe(action)
	}
	if process == MCP || process == All {
		h.actMCP(action)
	}
	h.o.Logger.Info("process "+action, "component", "host", "process", process)
	return nil
}

// actServe acts on serve; a restart of a stopped serve starts it.
func (h *host) actServe(action string) {
	h.mu.Lock()
	c := h.serve
	switch {
	case action == "stop":
		h.serveOff, h.serve = true, nil
	case c == nil:
		h.serveOff = false
		h.serve = h.newChild(Serve, "", h.o.Serve, false)
	}
	h.mu.Unlock()
	switch {
	case c == nil:
	case action == "stop":
		c.stop()
	case action == "restart":
		c.restartNow()
	}
}

// actMCP acts on every MCP server; a restart of stopped ones starts them.
func (h *host) actMCP(action string) {
	h.mu.Lock()
	var cs []*child
	for _, c := range h.mcp {
		cs = append(cs, c)
	}
	wasOff := h.mcpOff
	h.mcpOff = action == "stop"
	h.mu.Unlock()
	for _, c := range cs {
		switch {
		case action == "stop":
			c.pause()
		case action == "start" || wasOff:
			c.resume()
		default:
			c.restartNow()
		}
	}
}

// wantMCP makes the MCP servers the host keeps exactly those of roots: it
// starts what is missing and stops what is no longer asked for.
func (h *host) wantMCP(roots []string) {
	want := map[string]bool{}
	for _, r := range roots {
		if r != "" {
			want[r] = true
		}
	}
	h.mu.Lock()
	var gone []*child
	for r, c := range h.mcp {
		if !want[r] {
			gone = append(gone, c)
			delete(h.mcp, r)
		}
	}
	for r := range want {
		if _, ok := h.mcp[r]; ok {
			continue
		}
		root := r
		c := h.newChild(MCP, root, func() (Spec, error) { return h.o.MCP(root) }, h.mcpOff)
		h.mcp[root] = c
		h.o.Logger.Info("mcp server kept", "component", "host", "root", root)
	}
	h.mu.Unlock()
	for _, c := range gone {
		c.stop()
		h.o.Logger.Info("mcp server released", "component", "host", "root", c.root)
	}
}

// upgrade installs a newer flai and, when one was installed, restarts the
// host on it once the answer has gone.
func (h *host) upgrade(ctx context.Context) (any, error) {
	if h.o.Upgrade == nil {
		return nil, errors.New("this host cannot upgrade flai")
	}
	res, installed, err := h.o.Upgrade(ctx)
	if err != nil {
		return nil, err
	}
	if installed {
		go func() {
			time.Sleep(200 * time.Millisecond) // the answer leaves first
			select {
			case h.restart <- struct{}{}:
			default:
			}
		}()
	}
	return map[string]any{"upgrade": res, "restarting": installed}, nil
}
