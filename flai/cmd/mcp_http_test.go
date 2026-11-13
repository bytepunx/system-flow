package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/ask"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/mcphttp"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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
func testMCPOverHTTP(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	inEnv(t, "FLAI_AGENT", "tester")
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
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: root}
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

// S-0583, ADR-0258: the sample flai mcp status prints names the per-server
// timeout flai's own MCP file names, so that a hand-written client copies it.
func TestTheClientSampleNamesThePerServerTimeout(t *testing.T) {
	t.Parallel()
	var sample struct {
		MCPServers map[string]struct {
			Type    string `json:"type"`
			URL     string `json:"url"`
			Timeout int64  `json:"timeout"`
		} `json:"mcpServers"`
	}
	text := mcpClientConfig("http://127.0.0.1:4242/mcp")
	if err := json.Unmarshal([]byte(text), &sample); err != nil {
		t.Fatalf("the sample is not JSON: %v: %s", err, text)
	}
	flai := sample.MCPServers["flai"]
	if flai.Type != "http" || flai.URL != "http://127.0.0.1:4242/mcp" || flai.Timeout != 1800000 || flai.Timeout != mcphttp.ClientTimeout.Milliseconds() {
		t.Errorf("the sample's server is %+v, want type http, the URL, and timeout 1800000: %s", flai, text)
	}
}

func TestMCPStateGoesStale(t *testing.T) {
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
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

// S-0105: a server flai serve started is started again with a new token as
// flai serve's still, and one started by hand as it was.
func TestARotatedServerKeepsWhatItWentWith(t *testing.T) {
	t.Parallel()
	st := mcpState{Addr: "127.0.0.1:4247", ExitWith: os.Getpid()}
	if got := strings.Join(restartArgs(st, "/c/config.json"), " "); got != "mcp http --addr 127.0.0.1:4247 --exit-with "+strconv.Itoa(os.Getpid())+" --config /c/config.json" {
		t.Errorf("serve's: %s", got)
	}
	st.ExitWith = 0
	if got := strings.Join(restartArgs(st, ""), " "); got != "mcp http --addr 127.0.0.1:4247" {
		t.Errorf("by hand: %s", got)
	}
	gone := exec.Command("true")
	_ = gone.Run()
	st.ExitWith = gone.Process.Pid
	if got := strings.Join(restartArgs(st, ""), " "); strings.Contains(got, "--exit-with") {
		t.Errorf("a flai serve that has gone is not waited on: %s", got)
	}
}

// claudeCodeHTTPTimeout is how long Claude Code waits for an MCP call over
// HTTP before it drops the call, as S-0519 measured it (ADR-0215).
const claudeCodeHTTPTimeout = 60 * time.Second

// I-0217: over HTTP, permission_prompt held an unanswered write four minutes,
// so Claude Code dropped the call after 60 seconds and the agent saw the
// write time out while its thread stayed open with nobody waiting on it. The
// server an agent reaches over HTTP now holds a write less than Claude Code
// waits: with both scaled down alike, a client that gives up when Claude Code
// does gets the refusal naming the open thread.
func testMCPOverHTTPRefusesAnUnansweredWriteBeforeClaudeCodeGivesUp(t *testing.T) {
	root := heldProject(t) // S-0001 in progress
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: root}
	opt := a.mcpHTTPOptions(repo, "claude", nil, "")
	if opt.PermissionHold <= 0 || opt.PermissionHold >= claudeCodeHTTPTimeout {
		t.Errorf("over HTTP permission_prompt holds a write for %s, not less than the %s Claude Code waits", opt.PermissionHold, claudeCodeHTTPTimeout)
	}
	const scale = 20
	opt.PermissionHold /= scale
	srv := httptest.NewServer(mcphttp.Handler(mcphttp.Options{Token: "token", NewServer: func(string) *mcp.Server { return mcpserver.New(opt) }}))
	defer srv.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: srv.URL, HTTPClient: &http.Client{Transport: mcpBearer{"token", "claude"}}, DisableStandaloneSSE: true}, &mcp.ClientSessionOptions{ProtocolVersion: "2025-11-25"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range tools.Tools {
		if want := "waits at most " + ask.Spoken(opt.PermissionHold); tool.Name == "permission_prompt" && !strings.Contains(tool.Description, want) {
			t.Errorf("permission_prompt's description states the bound in force, %q: %s", want, tool.Description)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), claudeCodeHTTPTimeout/scale)
	defer cancel()
	path := filepath.Join(repo.WorktreePath("S-0001"), ".claude", "settings.json")
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "permission_prompt", Arguments: map[string]any{"tool_name": "Write", "input": map[string]any{"file_path": path, "content": "{}"}}})
	if err != nil {
		t.Fatalf("the call outlives the client's %s, as it outlived Claude Code's %s: %v", claudeCodeHTTPTimeout/scale, claudeCodeHTTPTimeout, err)
	}
	all, err := threads.List(repo)
	if err != nil || len(all) != 1 || !all[0].Open() {
		t.Fatalf("the write is asked on one thread, left open: %v %+v", err, all)
	}
	text, _ := res.Content[0].(*mcp.TextContent)
	if text == nil || !strings.Contains(text.Text, `"behavior":"deny"`) || !strings.Contains(text.Text, "No answer on "+all[0].ID+" yet") {
		t.Errorf("the unanswered write is refused naming the open thread: %+v", res.Content)
	}
}

// I-0228: over HTTP, wait_for_events and wait_for_work held a call for the
// timeout_seconds asked, so Claude Code dropped it after 60 seconds and the
// agent saw an error instead of timed_out. The server an agent reaches over
// HTTP now holds one at most HTTPLongestWait (ADR-0244), which is checked to
// be under Claude Code's 60 seconds. With that bound scaled down, both asked
// for 600 seconds answer timed_out before a client whose deadline is far
// above the bound and far below the 600 seconds asked, so that only a wait
// held for what was asked fails, never a bound met late on a loaded host
// (I-0309).
func TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp(t *testing.T) {
	t.Parallel()
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: root}
	opt := a.mcpHTTPOptions(repo, "idle", nil, "")
	if opt.MaxWait != mcpserver.HTTPLongestWait || opt.MaxWait <= 0 || opt.MaxWait >= claudeCodeHTTPTimeout {
		t.Errorf("over HTTP a wait is held at most %s, want mcpserver.HTTPLongestWait (%s), less than the %s Claude Code waits", opt.MaxWait, mcpserver.HTTPLongestWait, claudeCodeHTTPTimeout)
	}
	const scale = 20
	opt.MaxWait /= scale
	srv := httptest.NewServer(mcphttp.Handler(mcphttp.Options{Token: "token", NewServer: func(string) *mcp.Server { return mcpserver.New(opt) }}))
	defer srv.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: srv.URL, HTTPClient: &http.Client{Transport: mcpBearer{"token", "idle"}}, DisableStandaloneSSE: true}, &mcp.ClientSessionOptions{ProtocolVersion: "2025-11-25"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	// clientDeadline is over five times the scaled bound and a fiftieth of
	// the 600 seconds asked: a wait held for what was asked outlives it, and
	// one held to the bound answers within it on a loaded host (I-0309).
	const clientDeadline = 12 * time.Second
	if bound := opt.MaxWait; clientDeadline < 4*bound {
		t.Fatalf("the client's deadline %s is not far above the scaled bound %s", clientDeadline, bound)
	}
	ctx, cancel := context.WithTimeout(context.Background(), clientDeadline)
	defer cancel()
	waits := []string{"wait_for_events", "wait_for_work"}
	answers := make(chan string, len(waits))
	for _, tool := range waits {
		go func() {
			res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: map[string]any{"timeout_seconds": 600}})
			var got struct {
				TimedOut bool `json:"timed_out"`
			}
			switch {
			case err != nil:
				answers <- fmt.Sprintf("%s outlives the client's %s, held past its bound of %s, as it outlived Claude Code's %s: %v", tool, clientDeadline, opt.MaxWait, claudeCodeHTTPTimeout, err)
			case res.IsError || len(res.Content) == 0:
				answers <- fmt.Sprintf("%s answers an error: %+v", tool, res.Content)
			default:
				text, _ := res.Content[0].(*mcp.TextContent)
				if text == nil || json.Unmarshal([]byte(text.Text), &got) != nil || !got.TimedOut {
					answers <- fmt.Sprintf("%s with nothing to answer answers timed_out within the bound: %+v", tool, res.Content)
					return
				}
				answers <- ""
			}
		}()
	}
	for range waits {
		if failed := <-answers; failed != "" {
			t.Error(failed)
		}
	}
}

// mcpHelperEnv, set to 1, makes TestMCPHelperProcess run flai with the
// arguments after "--", on this process's own standard streams: flai mcp on
// stdio reads and writes those, so a test reaches the server that command
// builds only from a process of its own.
const mcpHelperEnv = "FLAI_TEST_MCP_HELPER"

func TestMCPHelperProcess(t *testing.T) {
	if os.Getenv(mcpHelperEnv) != "1" {
		t.Skip("flaiProcess runs it as flai; it does nothing on its own")
	}
	if code := run(flag.Args(), os.Stdout, os.Stderr); code != 0 {
		t.Fatalf("flai %s exited %d", strings.Join(flag.Args(), " "), code)
	}
}

// flaiProcess is flai with args, run in dir by this test binary through
// TestMCPHelperProcess with t's config file, its standard error kept in a
// file that stderr reads once the process has ended.
func flaiProcess(t *testing.T, dir string, args ...string) (cmd *exec.Cmd, stderr func() string) {
	t.Helper()
	logPath := filepath.Join(t.TempDir(), "stderr.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = logFile.Close() })
	cmd = exec.Command(os.Args[0], append([]string{"-test.run=^TestMCPHelperProcess$", "--", "--config", configOf(t)}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), mcpHelperEnv+"=1")
	cmd.Stderr = logFile
	return cmd, func() string {
		data, _ := os.ReadFile(logPath)
		return string(data)
	}
}

// mcpStdioSession is a session with flai mcp on stdio, serving agent in root.
func mcpStdioSession(ctx context.Context, t *testing.T, root, agent string) *mcp.ClientSession {
	t.Helper()
	cmd, stderr := flaiProcess(t, root, "mcp", "--agent", agent)
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect to flai mcp on stdio: %v\n%s", err, stderr())
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

// mcpHTTPSession is a session with flai mcp http, serving repo, as agent
// names itself by X-Flai-Agent with token.
func mcpHTTPSession(ctx context.Context, t *testing.T, repo *workitem.Repo, token, agent string) *mcp.ClientSession {
	t.Helper()
	cmd, stderr := flaiProcess(t, repo.Root, "mcp", "http", "--addr", "127.0.0.1:0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	exited := make(chan struct{})
	go func() { _ = cmd.Wait(); close(exited) }()
	t.Cleanup(func() {
		_ = cmd.Process.Signal(os.Interrupt)
		select {
		case <-exited:
		case <-time.After(10 * time.Second):
			_ = cmd.Process.Kill()
			<-exited
		}
	})
	var st mcpState
	for alive := false; !alive; {
		select {
		case <-exited:
			t.Fatalf("flai mcp http stopped at once:\n%s", stderr())
		case <-ctx.Done():
			t.Fatal("flai mcp http never said it was up")
		case <-time.After(20 * time.Millisecond):
		}
		st, alive = readMCPState(repo, time.Now())
		alive = alive && st.PID == cmd.Process.Pid
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: st.URL, HTTPClient: &http.Client{Transport: mcpBearer{token, agent}}, DisableStandaloneSSE: true}, &mcp.ClientSessionOptions{ProtocolVersion: "2025-11-25"})
	if err != nil {
		t.Fatalf("connect to flai mcp http at %s: %v", st.URL, err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

// mcpToolDescriptions are the tools a session offers, each with its
// description.
func mcpToolDescriptions(ctx context.Context, t *testing.T, cs *mcp.ClientSession) map[string]string {
	t.Helper()
	tools := map[string]string{}
	for tool, err := range cs.Tools(ctx, nil) {
		if err != nil {
			t.Fatal(err)
		}
		tools[tool.Name] = tool.Description
	}
	return tools
}

// differsByTransport names each tool that the stdio and HTTP servers offer
// with different descriptions, and why.
var differsByTransport = map[string]string{
	"permission_prompt": "its description states how long it holds a write for the operator's answer: StdioPermissionHold over stdio, and HTTPPermissionHold over HTTP, below the 60 seconds after which Claude Code drops a call there (ADR-0215)",
	"wait_for_events":   "its description states the most it holds one call: LongestWait over stdio, and HTTPLongestWait over HTTP, the same bound as permission_prompt's there (ADR-0244)",
	"wait_for_work":     "its description states the most it holds one call: LongestWait over stdio, and HTTPLongestWait over HTTP, the same bound as permission_prompt's there (ADR-0244)",
}

// S-0506: an agent in a container reaches flai only over HTTP, so the server
// flai mcp http builds offers every tool the one flai mcp builds on stdio
// offers, for the same agent, and the tools S-0505 named answer alike on
// both. Each server is built by its own command in a flai process of its own.
func TestMCPOverHTTPOffersEveryToolStdioDoes(t *testing.T) {
	t.Parallel()
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := mcpToken(repo, false)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	const agent = "parity"
	stdio := mcpStdioSession(ctx, t, root, agent)
	overHTTP := mcpHTTPSession(ctx, t, repo, token, agent)

	onStdio, onHTTP := mcpToolDescriptions(ctx, t, stdio), mcpToolDescriptions(ctx, t, overHTTP)
	if len(onStdio) == 0 {
		t.Fatal("flai mcp on stdio offers no tools")
	}
	for name, desc := range onStdio {
		other, ok := onHTTP[name]
		why, differs := differsByTransport[name]
		switch {
		case !ok:
			t.Errorf("%s is offered on stdio and not over HTTP", name)
		case other != desc && !differs:
			t.Errorf("%s is described differently over HTTP than on stdio, and differsByTransport gives no reason:\nstdio: %s\nHTTP:  %s", name, desc, other)
		case other == desc && differs:
			t.Errorf("%s is described alike on both now, so differsByTransport's reason no longer holds: %s", name, why)
		}
	}
	for name := range onHTTP {
		if _, ok := onStdio[name]; !ok {
			t.Errorf("%s is offered over HTTP and not on stdio", name)
		}
	}

	// A refusal is an answer: the tool ran on the host and said why.
	outside := filepath.Join(root, "README.md")
	for _, c := range []struct {
		call mcp.CallToolParams
		says string
		ok   func(string) bool
	}{
		{mcp.CallToolParams{Name: "task_done", Arguments: map[string]any{"task": "T-9999"}}, "refuses a task that does not exist, naming it",
			func(got string) bool { return strings.HasPrefix(got, "error: ") && strings.Contains(got, "T-9999") }},
		{mcp.CallToolParams{Name: "permission_prompt", Arguments: map[string]any{"tool_name": "Write", "input": map[string]any{"file_path": outside, "content": "x"}}}, "refuses a write outside every worktree, naming it",
			func(got string) bool {
				return strings.Contains(got, `"behavior":"deny"`) && strings.Contains(got, outside)
			}},
	} {
		want, got := mcpAnswer(ctx, t, stdio, c.call), mcpAnswer(ctx, t, overHTTP, c.call)
		if got != want {
			t.Errorf("%s answers over HTTP otherwise than on stdio:\nstdio: %s\nHTTP:  %s", c.call.Name, want, got)
		}
		if !c.ok(got) {
			t.Errorf("%s over HTTP %s: %s", c.call.Name, c.says, got)
		}
	}
	for _, cs := range []*mcp.ClientSession{stdio, overHTTP} {
		var out struct {
			TimedOut bool `json:"timed_out"`
		}
		got := mcpAnswer(ctx, t, cs, mcp.CallToolParams{Name: "wait_for_events", Arguments: map[string]any{"timeout_seconds": 1}})
		if err := json.Unmarshal([]byte(got), &out); err != nil || !out.TimedOut {
			t.Errorf("wait_for_events with nothing to report runs out its second: %s", got)
		}
	}
}

// mcpAnswer is the text a tool call answers, prefixed "error: " when the
// tool reports an error; a protocol error fails t.
func mcpAnswer(ctx context.Context, t *testing.T, cs *mcp.ClientSession, call mcp.CallToolParams) string {
	t.Helper()
	res, err := cs.CallTool(ctx, &call)
	if err != nil {
		t.Fatalf("%s: %v", call.Name, err)
	}
	var text []string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			text = append(text, tc.Text)
		}
	}
	answer := strings.Join(text, "\n")
	if res.IsError {
		answer = "error: " + answer
	}
	return answer
}

// ADR-0278 §3: the server over HTTP, a child of flai host with no FLAI_ROLE
// of its callers, gives a caller the orchestrator's role only when its name
// is the orchestrator's run flai serve records for the project and still
// runs; any other name, a run that ended, or another project's gets none.
func TestMCPOverHTTPKnowsTheOrchestratorByFlaiServesRun(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("the run is a process sleep stands in for")
	}
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	a := &app{configPath: cfg, out: &bytes.Buffer{}, errOut: &bytes.Buffer{}}
	if a.mcpHTTPRole("/p", "orchestrator") != "" {
		t.Error("with no state flai serve wrote, no caller is the orchestrator")
	}
	run := exec.Command("sleep", "30")
	serve.Detach(run) // the leader of its own session, as flai serve starts an agent
	if err := run.Start(); err != nil {
		t.Fatal(err)
	}
	waited := make(chan struct{})
	go func() { _ = run.Wait(); close(waited) }()
	t.Cleanup(func() { _ = run.Process.Kill(); <-waited })
	now := time.Now().UTC().Format(time.RFC3339)
	write := func(state map[string]serve.AgentState) {
		t.Helper()
		dir := serve.DirFor(cfg)
		if err := os.MkdirAll(string(dir), 0o700); err != nil {
			t.Fatal(err)
		}
		data, _ := json.Marshal(state)
		if err := os.WriteFile(filepath.Join(string(dir), "agents.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write(map[string]serve.AgentState{
		"/p": {Orchestrator: &serve.AgentRun{Agent: "orchestrator", PID: run.Process.Pid, Started: now},
			Stories: map[string]*serve.AgentRun{"S-0001": {Story: "S-0001", Agent: "agent-S-0001", PID: run.Process.Pid, Started: now}}},
		"/q": {Orchestrator: &serve.AgentRun{Agent: "orchestrator", PID: run.Process.Pid, Started: now, Ended: now}},
	})
	for _, c := range []struct {
		why         string
		root, agent string
		want        string
	}{
		{"the orchestrator's run", "/p", "orchestrator", guard.RoleOrchestrate},
		{"a story's agent", "/p", "agent-S-0001", ""},
		{"a name any holder of the token can claim", "/p", "claude", ""},
		{"no name", "/p", "", ""},
		{"an orchestrator's run that ended", "/q", "orchestrator", ""},
		{"another project", "/r", "orchestrator", ""},
	} {
		if got := a.mcpHTTPRole(c.root, c.agent); got != c.want {
			t.Errorf("%s: role %q, want %q", c.why, got, c.want)
		}
	}
	_ = run.Process.Kill()
	<-waited
	if a.mcpHTTPRole("/p", "orchestrator") != "" {
		t.Error("once the process flai serve started has gone, its name is no longer the orchestrator")
	}
}

// needSecondLoopback skips t where 127.0.0.2, loopback on Linux, cannot be
// listened on, as on macOS without an alias.
func needSecondLoopback(t *testing.T) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.2:0")
	if err != nil {
		t.Skipf("127.0.0.2 cannot be listened on here: %v", err)
	}
	_ = ln.Close()
}

// gatewayDocker answers docker network ls and network inspect as an engine
// whose networks have gateways, or fails every call; it records each call.
type gatewayDocker struct {
	networks, gateways string
	fail               bool
	calls              []string
}

func (d *gatewayDocker) Run(_, name string, args ...string) (string, error) {
	d.calls = append(d.calls, name+" "+strings.Join(args, " "))
	if d.fail {
		return "Cannot connect to the Docker daemon", errors.New("exit status 1")
	}
	if len(args) > 1 && args[0] == "network" && args[1] == "ls" {
		return d.networks, nil
	}
	if len(args) > 1 && args[0] == "network" && args[1] == "inspect" {
		return d.gateways, nil
	}
	return "", fmt.Errorf("gatewayDocker does not answer %s %s", name, strings.Join(args, " "))
}

func (d *gatewayDocker) RunInput(dir, name, _ string, args ...string) (string, error) {
	return d.Run(dir, name, args...)
}

func (d *gatewayDocker) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

// S-0732: flai mcp http listens on every --addr, on one port, so that a
// client here reaches it on loopback and an agent in a container on its
// network's gateway. Every address answers, the state records them all as
// bound, and addr and url name the loopback one, as older readers expect.
func TestMCPOverHTTPListensOnEveryAddress(t *testing.T) {
	t.Parallel()
	needSecondLoopback(t)
	root := tempProject(t)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := mcpToken(repo, false)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	docker := &gatewayDocker{}
	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: root, runner: docker}
	server := newRootCmdWith(a)
	server.SetArgs([]string{"mcp", "http", "--addr", "127.0.0.1:0", "--addr", "127.0.0.2:0"})
	served := make(chan error, 1)
	go func() { served <- server.ExecuteContext(ctx) }()
	var st mcpState
	for alive := false; !alive; {
		select {
		case err := <-served:
			t.Fatalf("the server stopped: %v %s", err, errOut.String())
		case <-ctx.Done():
			t.Fatal("the server never said it was up")
		case <-time.After(20 * time.Millisecond):
		}
		st, alive = readMCPState(repo, time.Now())
	}

	if len(st.Addrs) != 2 {
		t.Fatalf("the state records both addresses as bound: %+v", st)
	}
	host0, port0, _ := net.SplitHostPort(st.Addrs[0])
	host1, port1, _ := net.SplitHostPort(st.Addrs[1])
	if host0 != "127.0.0.1" || host1 != "127.0.0.2" || port0 != port1 || port0 == "0" {
		t.Errorf("both addresses on one port the first listener picked: %v", st.Addrs)
	}
	if st.Addr != st.Addrs[0] || st.URL != "http://"+st.Addrs[0]+"/mcp" {
		t.Errorf("addr and url name the first loopback listener: %+v", st)
	}
	if data, _ := os.ReadFile(filepath.Join(root, ".flai-cache", "mcp-http.addr")); strings.TrimSpace(string(data)) != "127.0.0.1:0" {
		t.Errorf("the remembered address is the one addr names: %q", data)
	}
	if len(docker.calls) != 0 {
		t.Errorf("docker is not asked for gateways when every address is loopback: %v", docker.calls)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	for _, addr := range st.Addrs {
		url := "http://" + addr + "/mcp"
		res, err := http.Get(url)
		if err != nil {
			t.Fatalf("%s: %v", addr, err)
		}
		_ = res.Body.Close()
		if res.StatusCode != http.StatusUnauthorized || res.Header.Get("X-Flai-Project-Key") != "t" {
			t.Errorf("%s answers a request without the token with 401: %d %v", addr, res.StatusCode, res.Header)
		}
		cs, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: url, HTTPClient: &http.Client{Transport: mcpBearer{token, "both"}}, DisableStandaloneSSE: true}, &mcp.ClientSessionOptions{ProtocolVersion: "2025-11-25"})
		if err != nil {
			t.Fatalf("connect at %s: %v", addr, err)
		}
		if tools := mcpToolDescriptions(ctx, t, cs); len(tools) == 0 {
			t.Errorf("%s serves no tools", addr)
		}
		_ = cs.Close()
	}

	cancel()
	select {
	case err := <-served:
		if err != nil {
			t.Errorf("a stop is not a failure: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the server did not stop")
	}
	for _, addr := range st.Addrs {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			t.Errorf("a stop closes every listener, but %s is still held: %v", addr, err)
			continue
		}
		_ = ln.Close()
	}
}

// S-0732: a listen failure on any address closes those already opened and
// names the address that failed. It listens on real addresses, so it keeps
// what listenMCP does with the operating system covered. The test holds
// 127.0.0.1:<port> and nothing at 127.0.0.2:<port>, which another test's
// connection may take as its source before listenMCP opens it or after it
// closes it (I-0372): such a try is lost, not failed, and the test tries
// again on a new port. A listener listenMCP did not close fails every try.
func TestMCPListenFailureClosesWhatItOpened(t *testing.T) {
	t.Parallel()
	needSecondLoopback(t)
	const tries = 20
	var lost []string
	for try := 0; try < tries; try++ {
		why, ok := listenFailureClosesWhatItOpened(t)
		if ok {
			return
		}
		lost = append(lost, why)
	}
	t.Fatalf("every one of %d tries lost 127.0.0.2 at its port, or found it still held after the failure:\n%s", tries, strings.Join(lost, "\n"))
}

// listenFailureClosesWhatItOpened is one try of
// TestMCPListenFailureClosesWhatItOpened: ok when it showed the failure
// named and the address opened before it closed again, else why the try was
// lost. A failure that names the wrong thing fails the test.
func listenFailureClosesWhatItOpened(t *testing.T) (string, bool) {
	t.Helper()
	held, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = held.Close() }()
	_, port, _ := net.SplitHostPort(held.Addr().String())
	first, taken := net.JoinHostPort("127.0.0.2", port), net.JoinHostPort("127.0.0.1", port)
	lns, err := listenMCP([]string{first, taken})
	if err == nil {
		closeListeners(lns)
		t.Fatal("a held address is a failure")
	}
	if strings.Contains(err.Error(), first) {
		return fmt.Sprintf("%s was taken before listenMCP opened it: %v", first, err), false
	}
	if !strings.Contains(err.Error(), taken) || !strings.Contains(err.Error(), "--addr") {
		t.Fatalf("the failure names the address and what to do: %v", err)
	}
	ln, err := net.Listen("tcp", first)
	if err != nil {
		return fmt.Sprintf("%s was held after the failure: %v", first, err), false
	}
	_ = ln.Close()
	return "", true
}

// S-0732: the server listens on one port, so addresses whose ports differ
// are refused before anything listens or is remembered.
func TestMCPOverHTTPRefusesAddressesOnDifferentPorts(t *testing.T) {
	t.Parallel()
	root := tempProject(t)
	_, errOut, code := runIn(t, root, "mcp", "http", "--addr", "127.0.0.1:4243", "--addr", "127.0.0.2:4244")
	if code == 0 || !strings.Contains(errOut, "different ports") || !strings.Contains(errOut, "same port") {
		t.Errorf("addresses on different ports are refused, saying what to do: %d %s", code, errOut)
	}
	for _, name := range []string{"mcp-http.addr", "mcp-http.json"} {
		if _, err := os.Stat(filepath.Join(root, ".flai-cache", name)); !os.IsNotExist(err) {
			t.Errorf("a refused start leaves no %s: %v", name, err)
		}
	}
	for _, c := range []struct {
		addrs []string
		ok    bool
	}{
		{[]string{"127.0.0.1:4243"}, true},
		{[]string{"127.0.0.1:4243", "172.17.0.1:4243", "[::1]:4243"}, true},
		{[]string{"127.0.0.1:0", "127.0.0.2:0"}, true},
		{[]string{"127.0.0.1:0", "127.0.0.2:4243"}, false},
		{[]string{"127.0.0.1"}, false},
	} {
		if err := samePort(c.addrs); (err == nil) != c.ok {
			t.Errorf("samePort(%v) = %v, want ok %v", c.addrs, err, c.ok)
		}
	}
}

// reachLog is what logReach logs for addrs with docker, each line's level
// and addr, at debug and above.
func reachLog(t *testing.T, docker *gatewayDocker, addrs []string) (warned, gateways []string, debug int) {
	t.Helper()
	var buf bytes.Buffer
	a := &app{configPath: configOf(t), out: &bytes.Buffer{}, errOut: &bytes.Buffer{}, runner: docker, log: slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))}
	log := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", "mcp")
	a.logReach(log, t.TempDir(), addrs)
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line == "" {
			continue
		}
		var e struct {
			Level     string `json:"level"`
			Component string `json:"component"`
			Addr      string `json:"addr"`
			Err       string `json:"err"`
		}
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("not a JSON line: %s", line)
		}
		if e.Component != "mcp" {
			t.Errorf("the line is not component mcp's: %s", line)
		}
		switch e.Level {
		case "WARN":
			warned = append(warned, e.Addr)
		case "INFO":
			gateways = append(gateways, e.Addr)
		case "DEBUG":
			if e.Err != "" {
				debug++
			}
		}
	}
	return warned, gateways, debug
}

// S-0732: the warning for listening beyond this machine is logged for each
// address that is neither loopback nor a Docker network's gateway, and a
// gateway gets an info line instead; docker is asked only when an address
// could be a gateway, and when it cannot answer, every such address is
// warned about and why is logged at debug.
func TestMCPWarnsOnlyBeyondLoopbackAndTheContainerGateways(t *testing.T) {
	t.Parallel()
	answers := func() *gatewayDocker {
		return &gatewayDocker{networks: "bridge\nflai\nhost\n", gateways: "172.17.0.1 \n172.18.0.1 \n\n"}
	}
	docker := answers()
	warned, gateways, _ := reachLog(t, docker, []string{"127.0.0.1:4243", "172.17.0.1:4243", "203.0.113.7:4243", "[::]:4243", "172.18.0.1:4243"})
	if !slices.Equal(warned, []string{"203.0.113.7:4243", "[::]:4243"}) {
		t.Errorf("warned of %v, want the address beyond this machine and the unspecified one", warned)
	}
	if !slices.Equal(gateways, []string{"172.17.0.1:4243", "172.18.0.1:4243"}) {
		t.Errorf("the gateways get an info line: %v", gateways)
	}
	if len(docker.calls) == 0 || !strings.HasPrefix(docker.calls[0], "docker network ls") {
		t.Errorf("the container harness's docker client is asked for the networks: %v", docker.calls)
	}

	for _, addrs := range [][]string{{"127.0.0.1:4243", "[::1]:4243", "localhost:4243"}, {"127.0.0.1:4243", "0.0.0.0:4243"}} {
		docker := answers()
		warned, _, _ := reachLog(t, docker, addrs)
		if len(docker.calls) != 0 {
			t.Errorf("for %v no address could be a gateway, yet docker was asked: %v", addrs, docker.calls)
		}
		if want := slices.DeleteFunc(slices.Clone(addrs), loopback); !slices.Equal(warned, want) {
			t.Errorf("for %v warned of %v, want %v", addrs, warned, want)
		}
	}

	down := &gatewayDocker{fail: true}
	warned, gateways, debug := reachLog(t, down, []string{"127.0.0.1:4243", "172.17.0.1:4243"})
	if !slices.Equal(warned, []string{"172.17.0.1:4243"}) || len(gateways) != 0 {
		t.Errorf("when docker cannot answer, the address is warned about as before: warned %v, gateways %v", warned, gateways)
	}
	if debug == 0 {
		t.Error("why docker could not answer is logged at debug")
	}
}

// S-0732: a rotated server starts again on every address it listened on,
// and one whose state predates Addrs on its one Addr; the address the state
// names is the loopback one.
func TestARotatedServerListensOnEveryAddressAgain(t *testing.T) {
	t.Parallel()
	st := mcpState{Addr: "127.0.0.1:4247", Addrs: []string{"172.17.0.1:4247", "127.0.0.1:4247"}}
	if got := strings.Join(restartArgs(st, ""), " "); got != "mcp http --addr 172.17.0.1:4247 --addr 127.0.0.1:4247" {
		t.Errorf("several: %s", got)
	}
	st.Addrs = nil
	if got := strings.Join(restartArgs(st, ""), " "); got != "mcp http --addr 127.0.0.1:4247" {
		t.Errorf("an older state's one address: %s", got)
	}
	if got := primaryAddr([]string{"172.17.0.1:4247", "127.0.0.1:4247"}); got != 1 {
		t.Errorf("the state names the loopback address, not the first: %d", got)
	}
	if got := primaryAddr([]string{"172.17.0.1:4247", "203.0.113.7:4247"}); got != 0 {
		t.Errorf("with no loopback address, the first: %d", got)
	}
}
