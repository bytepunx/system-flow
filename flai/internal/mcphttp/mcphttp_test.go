package mcphttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const token = "test-token"

// whoami answers with the name the server was built for, which is all these
// tests need of a server: the tools themselves are mcpserver's.
func newServer(agent string) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "flai", Version: "test"}, nil)
	mcp.AddTool(s, &mcp.Tool{Name: "whoami", Description: "the agent's name"}, func(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: agent}}}, nil, nil
	})
	return s
}

func lab(t *testing.T, o Options) *httptest.Server {
	t.Helper()
	o.Token, o.ProjectKey, o.ProjectName, o.NewServer = token, "lab", "the lab", newServer
	srv := httptest.NewServer(Handler(o))
	t.Cleanup(srv.Close)
	return srv
}

type reply struct {
	status  int
	header  http.Header
	body    string
	message struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
}

func send(t *testing.T, method, url, body string, headers map[string]string) reply {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("Authorization", "Bearer "+token)
	for k, v := range headers {
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	data, _ := io.ReadAll(res.Body)
	r := reply{status: res.StatusCode, header: res.Header, body: string(data)}
	_ = json.Unmarshal(data, &r.message)
	return r
}

func initialize(name string) string {
	return `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"` + name + `","version":"1"}}}`
}

const (
	initialized = `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	whoami      = `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"whoami","arguments":{}}}`
)

// open starts a session and returns its ID.
func open(t *testing.T, url, client string, headers map[string]string) string {
	t.Helper()
	r := send(t, "POST", url, initialize(client), headers)
	sid := r.header.Get(sessionHeader)
	if r.status != 200 || sid == "" {
		t.Fatalf("initialize: %d %s", r.status, r.body)
	}
	if r := send(t, "POST", url, initialized, map[string]string{sessionHeader: sid}); r.status != 202 {
		t.Fatalf("initialized: %d %s", r.status, r.body)
	}
	return sid
}

func TestSessions(t *testing.T) {
	srv := lab(t, Options{})
	sid := open(t, srv.URL, "Claude Code (test)", nil)
	r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: sid})
	if r.status != 200 || !strings.Contains(r.body, `"Claude-Code-test"`) {
		t.Errorf("the client's name, made safe, names the agent: %d %s", r.status, r.body)
	}
	if !strings.HasPrefix(r.header.Get("Content-Type"), "application/json") {
		t.Errorf("answers are JSON, not an event stream: %s", r.header.Get("Content-Type"))
	}
	if r.header.Get(ProjectKeyHeader) != "lab" || r.header.Get(ProjectNameHeader) != "the%20lab" {
		t.Errorf("every answer names its project: %v", r.header)
	}

	named := open(t, srv.URL, "whatever", map[string]string{AgentHeader: "builder/1"})
	if r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: named}); !strings.Contains(r.body, `"builder-1"`) {
		t.Errorf("X-Flai-Agent comes first: %s", r.body)
	}
	// the first session is still its own agent
	if r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: sid}); !strings.Contains(r.body, `"Claude-Code-test"`) {
		t.Errorf("sessions keep their names apart: %s", r.body)
	}

	if r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: "nope"}); r.status != 404 {
		t.Errorf("an unknown session is 404, and the client starts again: %d %s", r.status, r.body)
	}
	if r := send(t, "DELETE", srv.URL, "", map[string]string{sessionHeader: sid}); r.status != 204 {
		t.Errorf("DELETE ends a session: %d %s", r.status, r.body)
	}
	if r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: sid}); r.status != 404 {
		t.Errorf("an ended session is gone: %d", r.status)
	}
}

func TestRefusals(t *testing.T) {
	srv := lab(t, Options{})
	for name, c := range map[string]struct {
		headers map[string]string
		status  int
		says    string
	}{
		"no token":    {map[string]string{"Authorization": ""}, 401, "flai mcp token"},
		"wrong token": {map[string]string{"Authorization": "Bearer nope"}, 401, "Bearer"},
		"not bearer":  {map[string]string{"Authorization": "Basic " + token}, 401, "Bearer"},
		"a browser":   {map[string]string{"Origin": "http://127.0.0.1"}, 403, "agents"},
	} {
		r := send(t, "POST", srv.URL, initialize("x"), c.headers)
		if r.status != c.status || r.message.Error == nil || !strings.Contains(r.message.Error.Message, c.says) {
			t.Errorf("%s: %d %s", name, r.status, r.body)
		}
		if r.header.Get(ProjectKeyHeader) != "lab" {
			t.Errorf("%s: a refusal names the project too", name)
		}
	}
	if r := send(t, "POST", srv.URL, initialize("x"), map[string]string{"Authorization": ""}); r.header.Get("WWW-Authenticate") == "" {
		t.Error("401 says how to authenticate")
	}
}

func TestSessionCapAndIdle(t *testing.T) {
	srv := lab(t, Options{MaxSessions: 2, Idle: 300 * time.Millisecond})
	a := open(t, srv.URL, "a", nil)
	open(t, srv.URL, "b", nil)
	r := send(t, "POST", srv.URL, initialize("c"), nil)
	if r.status != 503 || !strings.Contains(r.body, "limit of MCP sessions") {
		t.Fatalf("a third session is refused: %d %s", r.status, r.body)
	}
	if r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: a}); r.status != 200 {
		t.Errorf("the sessions there are still work: %d", r.status)
	}
	if r := send(t, "DELETE", srv.URL, "", map[string]string{sessionHeader: a}); r.status != 204 {
		t.Fatalf("delete: %d", r.status)
	}
	open(t, srv.URL, "c", nil) // room again
	// both remaining sessions idle out, and there is room for two more
	deadline := time.Now().Add(5 * time.Second)
	for {
		if r := send(t, "POST", srv.URL, initialize("d"), nil); r.status == 200 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("idle sessions never expired")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestInFlightCap(t *testing.T) {
	release := make(chan struct{})
	started := make(chan struct{}, 1)
	o := Options{MaxInFlight: 1}
	o.Token = token
	o.NewServer = func(string) *mcp.Server {
		s := mcp.NewServer(&mcp.Implementation{Name: "flai", Version: "test"}, nil)
		mcp.AddTool(s, &mcp.Tool{Name: "whoami"}, func(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, any, error) {
			started <- struct{}{}
			<-release
			return &mcp.CallToolResult{}, nil, nil
		})
		return s
	}
	srv := httptest.NewServer(Handler(o))
	defer srv.Close()
	sid := open(t, srv.URL, "a", nil)
	done := make(chan reply, 1)
	go func() { done <- send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: sid}) }()
	<-started
	if r := send(t, "POST", srv.URL, whoami, map[string]string{sessionHeader: sid}); r.status != 503 {
		t.Errorf("a request beyond the limit is refused, not queued: %d %s", r.status, r.body)
	}
	close(release)
	if r := <-done; r.status != 200 {
		t.Errorf("the held request finishes: %d %s", r.status, r.body)
	}
}

// bearer adds what an agent's configuration adds: the token and its name.
type bearer struct{ agent string }

func (b bearer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+token)
	if b.agent != "" {
		r.Header.Set(AgentHeader, b.agent)
	}
	return http.DefaultTransport.RoundTrip(r)
}

// The SDK's client speaks the newest revision it knows unless told otherwise:
// 2026-07-28, with no sessions. Both it and the revision before are served.
func TestBothRevisionsWithTheSDKClient(t *testing.T) {
	srv := lab(t, Options{})
	for name, c := range map[string]struct {
		version string
		agent   string
		want    string
		session bool
	}{
		"2026-07-28, named by header":   {"", "builder", "builder", false},
		"2026-07-28, named by client":   {"", "", "sdk-client", false},
		"2025-11-25, a session, header": {"2025-11-25", "builder", "builder", true},
		"2025-11-25, a session, client": {"2025-11-25", "", "sdk-client", true},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			client := mcp.NewClient(&mcp.Implementation{Name: "sdk client", Version: "1"}, nil)
			transport := &mcp.StreamableClientTransport{Endpoint: srv.URL, HTTPClient: &http.Client{Transport: bearer{c.agent}}, DisableStandaloneSSE: true}
			cs, err := client.Connect(ctx, transport, &mcp.ClientSessionOptions{ProtocolVersion: c.version})
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = cs.Close() }()
			if got := cs.InitializeResult().ProtocolVersion; (c.version != "" && got != c.version) || (c.version == "" && got < sessionless) {
				t.Errorf("negotiated %s", got)
			}
			if (cs.ID() != "") != c.session {
				t.Errorf("session ID %q", cs.ID())
			}
			res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "whoami"})
			if err != nil {
				t.Fatal(err)
			}
			if text, _ := res.Content[0].(*mcp.TextContent); text == nil || text.Text != c.want {
				t.Errorf("agent: %+v", res.Content[0])
			}
		})
	}
}
