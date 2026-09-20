// Package mcphttp serves flai's MCP server over Streamable HTTP on the host
// (ADR-0030). The tools, the rules, and the cursor are mcpserver's, the same
// as on stdio; this package decides only who may ask, under what name, and
// how many at once.
//
// Two revisions of the transport are served side by side, because the SDK
// serves them differently: up to 2025-11-25 a client initializes a session,
// which idles out and is capped; from 2026-07-28 there are no sessions and
// every request stands alone. The Mcp-Protocol-Version header says which.
package mcphttp

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	// AgentHeader names the agent behind a session or a request.
	AgentHeader = "X-Flai-Agent"
	// ProjectKeyHeader names the project on every answer (ADR-0024, kept by
	// ADR-0030).
	ProjectKeyHeader = "X-Flai-Project-Key"
	// ProjectNameHeader carries the project's name, URI-encoded.
	ProjectNameHeader = "X-Flai-Project-Name"

	sessionHeader  = "Mcp-Session-Id"
	versionHeader  = "Mcp-Protocol-Version"
	sessionless    = "2026-07-28" // the first revision without sessions
	clientInfoMeta = "io.modelcontextprotocol/clientInfo"
	maxBody        = 4 << 20 // the SDK's own bound on a request body
	maxAgents      = 256     // servers kept by agent name before idle ones are dropped
)

// Options configure the handler.
type Options struct {
	Token       string // the bearer token; required
	ProjectKey  string
	ProjectName string
	// NewServer builds the MCP server for an agent name. It is called once
	// per name and revision; the server is shared by that agent's sessions.
	NewServer   func(agent string) *mcp.Server
	MaxSessions int           // sessions at once, default 16
	MaxInFlight int           // requests at once, default 64
	Idle        time.Duration // a session with no request for this long ends, default 30m
	Logger      *slog.Logger
}

type ctxKey struct{}

type handler struct {
	opt       Options
	stateful  *pool
	stateless *pool
	sessions  http.Handler
	alone     http.Handler
	inFlight  chan struct{}
}

// pool keeps one server per agent name.
type pool struct {
	mu      sync.Mutex
	build   func(string) *mcp.Server
	servers map[string]*mcp.Server
}

func (p *pool) get(agent string) *mcp.Server {
	p.mu.Lock()
	defer p.mu.Unlock()
	if s, ok := p.servers[agent]; ok {
		return s
	}
	if len(p.servers) >= maxAgents {
		for name, s := range p.servers {
			if count(s) == 0 {
				delete(p.servers, name)
			}
		}
	}
	s := p.build(agent)
	p.servers[agent] = s
	return s
}

func (p *pool) live() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := 0
	for _, s := range p.servers {
		n += count(s)
	}
	return n
}

func count(s *mcp.Server) int {
	n := 0
	for range s.Sessions() {
		n++
	}
	return n
}

// Handler serves MCP at whatever path it is mounted on.
func Handler(opt Options) http.Handler {
	if opt.MaxSessions <= 0 {
		opt.MaxSessions = 16
	}
	if opt.MaxInFlight <= 0 {
		opt.MaxInFlight = 64
	}
	if opt.Idle <= 0 {
		opt.Idle = 30 * time.Minute
	}
	if opt.Logger == nil {
		opt.Logger = slog.New(slog.DiscardHandler)
	}
	h := &handler{
		opt:       opt,
		stateful:  &pool{build: opt.NewServer, servers: map[string]*mcp.Server{}},
		stateless: &pool{build: opt.NewServer, servers: map[string]*mcp.Server{}},
		inFlight:  make(chan struct{}, opt.MaxInFlight),
	}
	pick := func(p *pool) func(*http.Request) *mcp.Server {
		return func(r *http.Request) *mcp.Server {
			agent, _ := r.Context().Value(ctxKey{}).(string)
			return p.get(agent)
		}
	}
	// JSON answers, not event streams: nothing is sent while a call runs, and
	// a plain body is what a proxy and curl expect.
	h.sessions = mcp.NewStreamableHTTPHandler(pick(h.stateful), &mcp.StreamableHTTPOptions{JSONResponse: true, SessionTimeout: opt.Idle, Logger: opt.Logger})
	h.alone = mcp.NewStreamableHTTPHandler(pick(h.stateless), &mcp.StreamableHTTPOptions{JSONResponse: true, Stateless: true, Logger: opt.Logger})
	return h
}

func (h *handler) refuse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": nil, "error": map[string]any{"code": -32000, "message": message}})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ProjectKeyHeader, h.opt.ProjectKey)
	w.Header().Set(ProjectNameHeader, url.PathEscape(h.opt.ProjectName))
	// An agent is not a browser. A page in one, on any site, must not reach
	// this server with a token it found or guessed, so an Origin is refused
	// before the token is looked at.
	if r.Header.Get("Origin") != "" {
		h.refuse(w, http.StatusForbidden, "requests from a browser are refused: this endpoint is for agents, which send no Origin header")
		return
	}
	given, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || h.opt.Token == "" || subtle.ConstantTimeCompare([]byte(strings.TrimSpace(given)), []byte(h.opt.Token)) != 1 {
		w.Header().Set("WWW-Authenticate", `Bearer realm="flai mcp"`)
		h.refuse(w, http.StatusUnauthorized, "send the token as Authorization: Bearer <token>; flai mcp token prints it")
		return
	}
	select {
	case h.inFlight <- struct{}{}:
		defer func() { <-h.inFlight }()
	default:
		h.refuse(w, http.StatusServiceUnavailable, "too many requests at once; try again shortly")
		return
	}

	alone := r.Header.Get(versionHeader) >= sessionless
	// The limit bounds what sessions cost; two that open in the same instant
	// may both pass it, which is not worth a lock held across a request.
	opening := !alone && r.Method == http.MethodPost && r.Header.Get(sessionHeader) == ""
	if opening && h.stateful.live() >= h.opt.MaxSessions {
		h.refuse(w, http.StatusServiceUnavailable, "this server already has its limit of MCP sessions; end one with DELETE or wait for an idle one to expire")
		return
	}
	if alone || opening {
		agent := SafeAgent(r.Header.Get(AgentHeader))
		if agent == "" && r.Method == http.MethodPost {
			agent = SafeAgent(clientName(w, r))
		}
		if agent == "" {
			agent = "agent"
		}
		r = r.WithContext(context.WithValue(r.Context(), ctxKey{}, agent))
		if opening {
			h.opt.Logger.Info("mcp session opening", "agent", agent, "sessions", h.stateful.live())
		}
	}
	if alone {
		h.alone.ServeHTTP(w, r)
		return
	}
	h.sessions.ServeHTTP(w, r)
}

var unsafeAgent = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// SafeAgent makes a name fit for a thread entry and a cursor file.
func SafeAgent(s string) string {
	s = strings.Trim(unsafeAgent.ReplaceAllString(s, "-"), "-.")
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

// clientName reads the client's name from the body and puts the body back:
// from initialize where there are sessions, else from the request's _meta.
func clientName(w http.ResponseWriter, r *http.Request) string {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return ""
	}
	type message struct {
		Method string `json:"method"`
		Params struct {
			ClientInfo struct {
				Name string `json:"name"`
			} `json:"clientInfo"`
			Meta map[string]json.RawMessage `json:"_meta"`
		} `json:"params"`
	}
	var many []message
	if json.Unmarshal(body, &many) != nil {
		var one message
		if json.Unmarshal(body, &one) != nil {
			return ""
		}
		many = []message{one}
	}
	for _, m := range many {
		if m.Method == "initialize" && m.Params.ClientInfo.Name != "" {
			return m.Params.ClientInfo.Name
		}
		var info struct {
			Name string `json:"name"`
		}
		if raw, ok := m.Params.Meta[clientInfoMeta]; ok && json.Unmarshal(raw, &info) == nil && info.Name != "" {
			return info.Name
		}
	}
	return ""
}
