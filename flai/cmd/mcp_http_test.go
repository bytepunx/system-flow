package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type mcpBearer struct{ token, agent string }

func (b mcpBearer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+b.token)
	r.Header.Set("X-Flai-Agent", b.agent)
	return http.DefaultTransport.RoundTrip(r)
}

// S-0076: flai mcp http is the server flai mcp is, over HTTP, with its token,
// state, and address kept under .flai-cache.
func TestMCPOverHTTP(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	runIn(t, root, "epic", "new", "Big thing")
	runIn(t, root, "story", "new", "Slice", "--epic", "E-0001")

	if out, _, _ := runIn(t, root, "mcp", "status"); !strings.Contains(out, "not running over HTTP") || !strings.Contains(out, ".mcp.json") {
		t.Errorf("status before a start: %s", out)
	}
	token, _, code := runIn(t, root, "mcp", "token")
	token = strings.TrimSpace(token)
	if code != 0 || len(token) < 40 {
		t.Fatalf("token: %d %q", code, token)
	}
	if info, err := os.Stat(filepath.Join(root, ".flai-cache", "mcp.token")); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("the token file is the owner's only: %v %v", info, err)
	}
	if again, _, _ := runIn(t, root, "mcp", "token"); strings.TrimSpace(again) != token {
		t.Error("the token is kept, not made again")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root}
	served := make(chan error, 1)
	// built here, not in the goroutine: building a command tree writes to a
	// table cobra shares, and the status calls below build one each
	server := newRootCmdWith(a)
	server.SetArgs([]string{"mcp", "http", "--addr", "127.0.0.1:0"})
	go func() { served <- server.ExecuteContext(ctx) }()
	var state struct {
		Running bool `json:"running"`
		State   struct {
			URL string `json:"url"`
			PID int    `json:"pid"`
		} `json:"state"`
	}
	deadline := time.Now().Add(5 * time.Second)
	for !state.Running {
		select {
		case err := <-served:
			t.Fatalf("the server stopped: %v %s", err, errOut.String())
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("the server never said it was up")
		}
		js, _, _ := runIn(t, root, "mcp", "status", "--json")
		_ = json.Unmarshal([]byte(js), &state)
		time.Sleep(20 * time.Millisecond)
	}
	if state.State.PID != os.Getpid() || !strings.HasPrefix(state.State.URL, "http://127.0.0.1:") || !strings.HasSuffix(state.State.URL, "/mcp") {
		t.Fatalf("state: %+v", state)
	}
	if text, _, _ := runIn(t, root, "mcp", "status"); !strings.Contains(text, state.State.URL) || !strings.Contains(text, "X-Flai-Agent") || strings.Contains(text, token) {
		t.Errorf("status says where and how to connect, and never prints the token: %s", text)
	}
	if _, second, code := runIn(t, root, "mcp", "http", "--addr", "127.0.0.1:0"); code == 0 || !strings.Contains(second, "already serving") {
		t.Errorf("one server per project: %d %s", code, second)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	transport := &mcp.StreamableClientTransport{Endpoint: state.State.URL, HTTPClient: &http.Client{Transport: mcpBearer{token, "remote one"}}, DisableStandaloneSSE: true}
	cs, err := client.Connect(ctx, transport, &mcp.ClientSessionOptions{ProtocolVersion: "2025-11-25"})
	if err != nil {
		t.Fatal(err)
	}
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "inbox"})
	if err != nil || res.IsError {
		t.Fatalf("inbox: %v %+v", err, res)
	}
	text, _ := res.Content[0].(*mcp.TextContent)
	if text == nil || !strings.Contains(text.Text, `"agent":"remote-one"`) {
		t.Errorf("the agent is the header's name, made safe: %+v", res.Content[0])
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "mcp", "remote-one.json")); err != nil {
		t.Errorf("the cursor is kept by agent name, as on stdio: %v", err)
	}
	moved, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "item_move", Arguments: map[string]any{"id": "S-1", "to": "done"}})
	if err != nil || !moved.IsError {
		t.Errorf("the rules are the same: an agent does not accept a story: %v %+v", err, moved)
	}
	_ = cs.Close()

	// An idle agent holds wait_for_events. A stop ends that call as if its
	// time had passed, instead of cutting the agent off after a wait.
	holder, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: state.State.URL, HTTPClient: &http.Client{Transport: mcpBearer{token, "holder"}}, DisableStandaloneSSE: true}, &mcp.ClientSessionOptions{ProtocolVersion: "2025-11-25"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := holder.CallTool(ctx, &mcp.CallToolParams{Name: "inbox"}); err != nil {
		t.Fatal(err)
	}
	held := make(chan string, 1)
	go func() {
		res, err := holder.CallTool(context.Background(), &mcp.CallToolParams{Name: "wait_for_events", Arguments: map[string]any{"timeout_seconds": 120}})
		if err != nil || len(res.Content) == 0 {
			held <- "failed: " + fmt.Sprint(err)
			return
		}
		text, _ := res.Content[0].(*mcp.TextContent)
		held <- text.Text
	}()
	time.Sleep(300 * time.Millisecond) // let the call arrive and start waiting

	res2, err := http.Get(state.State.URL)
	if err != nil {
		t.Fatal(err)
	}
	_ = res2.Body.Close()
	if res2.StatusCode != 401 || res2.Header.Get("X-Flai-Project-Key") != "t" {
		t.Errorf("no token: %d %v", res2.StatusCode, res2.Header)
	}

	cancel()
	select {
	case got := <-held:
		if !strings.Contains(got, `"timed_out":true`) {
			t.Errorf("the held wait ends as a wait that ran out: %s", got)
		}
	case <-time.After(3 * time.Second):
		t.Error("the held wait was not ended by the stop")
	}
	select {
	case err := <-served:
		if err != nil {
			t.Errorf("a stop is not a failure: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the server did not stop")
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "mcp-http.json")); !os.IsNotExist(err) {
		t.Errorf("the state file goes with the server: %v", err)
	}
	if out, _, _ := runIn(t, root, "mcp", "stop"); !strings.Contains(out, "not running") {
		t.Errorf("stop with nothing to stop: %s", out)
	}
	rotated, _, _ := runIn(t, root, "mcp", "token", "--rotate")
	if strings.TrimSpace(rotated) == token || len(strings.TrimSpace(rotated)) < 40 {
		t.Error("--rotate replaces the token")
	}
}

func TestMCPStateGoesStale(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, ".flai-cache"), 0o700)
	old := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	_ = os.WriteFile(filepath.Join(root, ".flai-cache", "mcp-http.json"), []byte(`{"pid":`+"1"+`,"addr":"127.0.0.1:1","url":"http://127.0.0.1:1/mcp","updated":"`+old+`"}`), 0o600)
	if out, _, _ := runIn(t, root, "mcp", "status"); !strings.Contains(out, "not running") {
		t.Errorf("a state file nobody refreshes is a dead server's: %s", out)
	}
	if out, _, _ := runIn(t, root, "mcp", "stop"); !strings.Contains(out, "not running") {
		t.Errorf("stop signals nobody: %s", out)
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "mcp-http.json")); !os.IsNotExist(err) {
		t.Error("and clears what was left")
	}
}
