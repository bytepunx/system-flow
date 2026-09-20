// Package channel is the host end of the connection between flai and the
// dashboard (ADR-0029). flai opens a WebSocket to the dashboard's agent
// endpoint and the two speak JSON-RPC 2.0 over it; the dashboard never
// connects to the host. Everything the dashboard can ask for is a named
// method in a table given to the client: there is no method that takes a
// command line, and none is added by the other side.
package channel

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// Protocol is the version both sides name in hello.
const Protocol = 1

// MaxMessage caps one message in either direction, where the dashboard's
// buffer for a flai command's output was (16 MiB).
const MaxMessage = 16 << 20

// Path is the dashboard's agent endpoint.
const Path = "/agent"

// JSON-RPC error codes used here.
const (
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternal       = -32603
	CodeUnknownProject = -32001
)

// Error is a JSON-RPC error object.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

// Project is what flai serves over one connection.
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Root string `json:"-"`
}

// Method answers one named request for a project. Params are data to be
// validated, never text to be executed.
type Method func(ctx context.Context, p Project, params json.RawMessage) (any, *Error)

type message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// State is what flai serve reports about one connection.
type State struct {
	URL       string `json:"url"`
	Connected bool   `json:"connected"`
	Since     string `json:"since,omitempty"`      // when the current connection was proven
	LastError string `json:"last_error,omitempty"` // why the last attempt ended
	Attempts  int    `json:"attempts"`
}

// Client keeps one connection to one dashboard alive.
type Client struct {
	URL     string // http(s)://host:port of the dashboard
	Key     []byte // the agent credential; never sent, only proven
	Project Project
	Methods map[string]Method
	Version string
	Logger  *slog.Logger

	// Tests shorten these.
	MaxBytes   int // cap on one message, MaxMessage when zero
	PingEvery  time.Duration
	MinBackoff time.Duration
	MaxBackoff time.Duration
	Now        func() time.Time

	mu    sync.Mutex
	state State
}

func (c *Client) defaults() {
	if c.MaxBytes <= 0 {
		c.MaxBytes = MaxMessage
	}
	if c.PingEvery <= 0 {
		c.PingEvery = 5 * time.Second
	}
	if c.MinBackoff <= 0 {
		c.MinBackoff = 250 * time.Millisecond
	}
	if c.MaxBackoff <= 0 {
		c.MaxBackoff = 4 * time.Second
	}
	if c.Now == nil {
		c.Now = time.Now
	}
	if c.Logger == nil {
		c.Logger = slog.New(slog.DiscardHandler)
	}
}

// State reports the connection as it is now.
func (c *Client) State() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

func (c *Client) setState(f func(*State)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	f(&c.state)
}

// Run connects and serves until ctx ends, reconnecting with backoff and
// jitter. A connection that held for a while resets the backoff.
func (c *Client) Run(ctx context.Context) {
	c.defaults()
	c.setState(func(s *State) { s.URL = c.URL })
	backoff := c.MinBackoff
	for ctx.Err() == nil {
		started := c.Now()
		err := c.serveOnce(ctx)
		c.setState(func(s *State) {
			s.Connected, s.Since, s.Attempts = false, "", s.Attempts+1
			if err != nil {
				s.LastError = err.Error()
			}
		})
		if ctx.Err() != nil {
			return
		}
		if c.Now().Sub(started) > 2*c.MaxBackoff {
			backoff = c.MinBackoff
		}
		c.Logger.Info("dashboard connection ended", "component", "channel", "url", c.URL, "project", c.Project.Key, "err", errText(err), "retry_in", backoff.String())
		select {
		case <-ctx.Done():
			return
		case <-time.After(jitter(backoff)):
		}
		if backoff *= 2; backoff > c.MaxBackoff {
			backoff = c.MaxBackoff
		}
	}
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// jitter spreads d over [d/2, d) so that several clients do not return together.
func jitter(d time.Duration) time.Duration {
	half := int64(d / 2)
	if half <= 0 {
		return d
	}
	n, err := rand.Int(rand.Reader, big.NewInt(half))
	if err != nil {
		return d
	}
	return time.Duration(half + n.Int64())
}

func wsURL(base string) (string, error) {
	switch {
	case strings.HasPrefix(base, "http://"):
		return "ws://" + strings.TrimSuffix(strings.TrimPrefix(base, "http://"), "/") + Path, nil
	case strings.HasPrefix(base, "https://"):
		return "wss://" + strings.TrimSuffix(strings.TrimPrefix(base, "https://"), "/") + Path, nil
	}
	return "", fmt.Errorf("dashboard address %q is not http or https", base)
}

// Proof is HMAC-SHA256(key, role|first|second) in hex. The role keeps one
// side's proof from being replayed as the other's.
func Proof(key []byte, role, first, second string) string {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(role + "|" + first + "|" + second))
	return hex.EncodeToString(m.Sum(nil))
}

func nonce() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type helloParams struct {
	Protocol int     `json:"protocol"`
	Nonce    string  `json:"nonce"`
	Flai     string  `json:"flai"`
	Project  Project `json:"project"`
}

type helloResult struct {
	Nonce     string `json:"nonce"`
	Proof     string `json:"proof"`
	Dashboard string `json:"dashboard,omitempty"`
}

// serveOnce dials, proves, and serves requests until the connection ends.
func (c *Client) serveOnce(ctx context.Context) error {
	u, err := wsURL(c.URL)
	if err != nil {
		return err
	}
	dialCtx, cancelDial := context.WithTimeout(ctx, 10*time.Second)
	conn, resp, err := websocket.Dial(dialCtx, u, &websocket.DialOptions{
		HTTPHeader: http.Header{"X-Flai-Agent": {"flai/" + c.Version}},
	})
	cancelDial()
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		return err
	}
	defer func() { _ = conn.CloseNow() }()
	conn.SetReadLimit(int64(c.MaxBytes))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wmu sync.Mutex
	send := func(m message) error {
		m.JSONRPC = "2.0"
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}
		if len(b) > c.MaxBytes {
			return fmt.Errorf("an answer of %d bytes is over the %d byte cap", len(b), c.MaxBytes)
		}
		wmu.Lock()
		defer wmu.Unlock()
		wctx, done := context.WithTimeout(ctx, 30*time.Second)
		defer done()
		return conn.Write(wctx, websocket.MessageText, b)
	}

	// hello: flai's nonce out, the dashboard's nonce and proof back.
	mine := nonce()
	params, _ := json.Marshal(helloParams{Protocol: Protocol, Nonce: mine, Flai: c.Version, Project: c.Project})
	if err := send(message{ID: json.RawMessage(`"hello"`), Method: "hello", Params: params}); err != nil {
		return err
	}
	hctx, hdone := context.WithTimeout(ctx, 10*time.Second)
	_, data, err := conn.Read(hctx)
	hdone()
	if err != nil {
		return fmt.Errorf("no answer to hello: %w", err)
	}
	var ans message
	if err := json.Unmarshal(data, &ans); err != nil || string(ans.ID) != `"hello"` {
		return errors.New("the first message back was not the answer to hello")
	}
	if ans.Error != nil {
		return fmt.Errorf("the dashboard refused hello: %s", ans.Error.Message)
	}
	var hr helloResult
	if err := json.Unmarshal(ans.Result, &hr); err != nil || hr.Nonce == "" {
		return errors.New("the answer to hello has no nonce")
	}
	if !hmac.Equal([]byte(hr.Proof), []byte(Proof(c.Key, "dashboard", mine, hr.Nonce))) {
		_ = conn.Close(websocket.StatusPolicyViolation, "not my dashboard")
		return errors.New("the dashboard did not prove it holds the agent credential; is something else listening on this port?")
	}
	proof, _ := json.Marshal(map[string]string{"proof": Proof(c.Key, "flai", hr.Nonce, mine)})
	if err := send(message{Method: "hello.prove", Params: proof}); err != nil {
		return err
	}
	c.setState(func(s *State) {
		s.Connected, s.Since, s.LastError = true, c.Now().UTC().Format(time.RFC3339), ""
	})
	c.Logger.Info("connected to the dashboard", "component", "channel", "url", c.URL, "project", c.Project.Key, "dashboard", hr.Dashboard)

	// A dashboard that stops answering pings is gone, however open the socket looks.
	go func() {
		t := time.NewTicker(c.PingEvery)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				pctx, done := context.WithTimeout(ctx, c.PingEvery)
				err := conn.Ping(pctx)
				done()
				if err != nil {
					cancel()
					return
				}
			}
		}
	}()

	var inflight sync.Map // request id -> cancel
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return err
		}
		var m message
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		switch {
		case m.Method == "$/cancel":
			var p struct {
				ID json.RawMessage `json:"id"`
			}
			if json.Unmarshal(m.Params, &p) == nil {
				if stop, ok := inflight.Load(string(p.ID)); ok {
					stop.(context.CancelFunc)()
				}
			}
		case m.Method != "" && m.ID != nil:
			rctx, stop := context.WithCancel(ctx)
			inflight.Store(string(m.ID), stop)
			go func(m message) {
				defer func() { stop(); inflight.Delete(string(m.ID)) }()
				out := message{ID: m.ID}
				res, rerr := c.call(rctx, m)
				if rerr != nil {
					out.Error = rerr
				} else if b, err := json.Marshal(res); err != nil {
					out.Error = &Error{Code: CodeInternal, Message: err.Error()}
				} else {
					out.Result = b
				}
				if err := send(out); err != nil && out.Error == nil {
					_ = send(message{ID: m.ID, Error: &Error{Code: CodeInternal, Message: err.Error()}})
				}
			}(m)
		}
	}
}

// call finds the method and checks the project before anything runs.
func (c *Client) call(ctx context.Context, m message) (any, *Error) {
	method, ok := c.Methods[m.Method]
	if !ok {
		return nil, &Error{Code: CodeMethodNotFound, Message: "flai offers no method " + m.Method}
	}
	var p struct {
		Project string `json:"project"`
	}
	if len(m.Params) > 0 {
		if err := json.Unmarshal(m.Params, &p); err != nil {
			return nil, &Error{Code: CodeInvalidParams, Message: "params must be an object"}
		}
	}
	if p.Project != c.Project.Key {
		return nil, &Error{Code: CodeUnknownProject, Message: fmt.Sprintf("this flai does not serve project %q on this connection", p.Project)}
	}
	return method(ctx, c.Project, m.Params)
}
