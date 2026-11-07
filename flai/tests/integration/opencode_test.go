package integration

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// openCodeVersion is the OpenCode release the flai-guard plugin was read
// against and is tested on (agent-adapters.md § Read from the source at
// v1.19.5).
const openCodeVersion = "1.19.5"

// guardPlugin is the template's flai-guard plugin, from this package's
// folder.
const guardPlugin = "../../../template/root/.opencode/plugins/flai-guard.js"

// subAgentRefused is the start of flai guard's refusal of a sub-agent's call.
const subAgentRefused = "a sub-agent ("

// buildFlai builds flai from this module into a temporary folder and returns
// that folder, to put first on PATH.
func buildFlai(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	build := exec.Command("go", "build", "-o", filepath.Join(dir, "flai"), "../..")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building flai: %v\n%s", err, out)
	}
	return dir
}

// TestFlaiGuardPluginUnderNode loads the template's flai-guard plugin with
// node, a fake OpenCode client, and the flai this module builds, and shows
// that a sub-agent's call to flai's item_move throws the guard's refusal
// while the story's own agent's passes (S-0407). It needs node on PATH and
// skips without it:
//
//	go test ./tests/integration/ -run TestFlaiGuardPluginUnderNode -v
func TestFlaiGuardPluginUnderNode(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not on PATH, so the flai-guard plugin cannot be loaded outside OpenCode")
	}
	plugin, err := filepath.Abs(guardPlugin)
	if err != nil {
		t.Fatal(err)
	}
	bin := buildFlai(t)
	script := filepath.Join(t.TempDir(), "drive.mjs")
	const drive = `const { FlaiGuard } = await import(process.argv[2]);
const client = { session: { get: async ({ path }) => ({ data: path.id === "ses_child" ? { id: "ses_child", parentID: "ses_root", agent: "explore" } : { id: path.id } }) } };
const hooks = await FlaiGuard({ client, directory: process.cwd() });
const out = {};
for (const id of ["ses_child", "ses_root"]) {
  try {
    await hooks["tool.execute.before"]({ tool: "flai_item_move", sessionID: id, callID: "call_1" }, { args: { id: "S-0001", to: "review" } });
    out[id] = "";
  } catch (e) {
    out[id] = e.message;
  }
}
console.log(JSON.stringify(out));
`
	if err := os.WriteFile(script, []byte(drive), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, script, plugin)
	cmd.Dir = t.TempDir()
	cmd.Env = append(cleanFlaiEnv(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "FLAI_STORY=S-0001")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("node %s: %v", script, err)
	}
	var got map[string]string
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("reading %q: %v", out, err)
	}
	if !strings.Contains(got["ses_child"], subAgentRefused) || !strings.Contains(got["ses_child"], "item_move") {
		t.Errorf("the sub-agent's item_move was not refused by the guard: %q", got["ses_child"])
	}
	if got["ses_root"] != "" {
		t.Errorf("the story's own agent's item_move was refused: %q", got["ses_root"])
	}
}

// cleanFlaiEnv is the test's environment without the FLAI_ variables of the
// session that runs it, so that the guard sees only what the test sets, and
// without PWD, which names this package's folder inside a system-flow
// checkout: OpenCode takes its directory from PWD, so a run would otherwise
// work in this repository rather than in the test's project.
func cleanFlaiEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "FLAI_") || strings.HasPrefix(kv, "PWD=") {
			continue
		}
		env = append(env, kv)
	}
	return env
}

// TestFlaiGuardPluginOnOpenCode runs opencode run, at the release the plugin
// was read against, in a project holding the template's flai-guard plugin,
// against a fake provider that scripts the session (S-0407). The story's own
// agent hands a task to the built-in general sub-agent, which runs flai move
// through bash; then the agent runs the same command itself. The test shows
// from the tool results OpenCode sends back to the provider that the plugin
// refused the sub-agent's call with flai guard's reason and let the agent's
// own through. It needs opencode 1.19.5 on PATH and skips without it, for
// instance:
//
//	npm install --prefix .flai-cache/scratch/opencode opencode-ai@1.19.5
//	PATH="$PWD/.flai-cache/scratch/opencode/node_modules/.bin:$PATH" \
//	go test ./flai/tests/integration/ -run TestFlaiGuardPluginOnOpenCode -v
//
// (the go test from the flai folder, with the PATH entry absolute).
func TestFlaiGuardPluginOnOpenCode(t *testing.T) {
	opencode, err := exec.LookPath("opencode")
	if err != nil {
		t.Skipf("opencode is not on PATH: install opencode-ai@%s to run the flai-guard plugin on OpenCode", openCodeVersion)
	}
	version, err := exec.Command(opencode, "--version").Output()
	if err != nil {
		t.Skipf("opencode --version failed (%v): the flai-guard plugin is tested on OpenCode %s", err, openCodeVersion)
	}
	if v := strings.TrimSpace(string(version)); v != openCodeVersion {
		t.Skipf("opencode on PATH is %s, not %s, the release the flai-guard plugin was read against", v, openCodeVersion)
	}
	plugin, err := os.ReadFile(guardPlugin)
	if err != nil {
		t.Fatal(err)
	}
	bin := buildFlai(t)

	provider := &scriptedProvider{t: t}
	srv := httptest.NewServer(provider)
	defer srv.Close()

	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".opencode", "plugins"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".opencode", "plugins", "flai-guard.js"), plugin, 0o600); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", project, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	config, err := json.Marshal(map[string]any{
		"$schema":     "https://opencode.ai/config.json",
		"autoupdate":  false,
		"share":       "disabled",
		"model":       "fake/scripted",
		"small_model": "fake/scripted",
		"provider": map[string]any{
			"fake": map[string]any{
				"npm":     "@ai-sdk/openai-compatible",
				"name":    "Fake",
				"options": map[string]any{"baseURL": srv.URL + "/v1", "apiKey": "test"},
				"models":  map[string]any{"scripted": map[string]any{"name": "scripted"}},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, opencode, "run", "--format", "json", "--model", "fake/scripted", "Run flai move S-0001 review through a sub-agent, then yourself.")
	cmd.Dir = project
	cmd.Env = append(cleanFlaiEnv(),
		"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"PWD="+project,
		"HOME="+home,
		"XDG_CONFIG_HOME="+filepath.Join(home, ".config"),
		"XDG_DATA_HOME="+filepath.Join(home, ".local", "share"),
		"XDG_STATE_HOME="+filepath.Join(home, ".local", "state"),
		"XDG_CACHE_HOME="+filepath.Join(home, ".cache"),
		"OPENCODE_CONFIG_CONTENT="+string(config),
		"OPENCODE_DISABLE_AUTOUPDATE=1",
		"OPENCODE_DISABLE_CLAUDE_CODE=1",
		"OPENCODE_DISABLE_LSP_DOWNLOAD=true",
		"FLAI_STORY=S-0001",
	)
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("opencode run did not end within its time: %v\n%s", ctx.Err(), out)
	}
	if err != nil {
		t.Logf("opencode run: %v\n%s", err, out)
	}
	child, root := provider.results()
	if len(child) == 0 {
		t.Fatalf("the sub-agent sent back no tool result; opencode printed:\n%s", out)
	}
	if !strings.Contains(child[0], subAgentRefused) || !strings.Contains(child[0], flaiWrite) {
		t.Errorf("the sub-agent's flai move was not refused by flai guard; its tool result was %q", child[0])
	}
	if len(root) < 2 {
		t.Fatalf("the story's own agent sent back %d tool results, want the task's and its own flai move's; opencode printed:\n%s", len(root), out)
	}
	// the guard let the agent's own call through when flai itself ran and
	// answered, here that the test's project is no system-flow project
	if own := root[len(root)-1]; strings.Contains(own, subAgentRefused) || !strings.Contains(own, "no system-flow.yaml found") {
		t.Errorf("the story's own agent's flai move did not reach flai in the test's project: %q", own)
	}
}

// flaiWrite is the flai write both the sub-agent and the agent run.
const flaiWrite = "flai move S-0001 review"

// scriptedProvider is an OpenAI-compatible chat completions endpoint that
// plays a fixed session: to the root session it answers a task for the
// general sub-agent, then a bash flai write of its own, then text; to the
// sub-agent's session a bash flai write, then text; to a request with no
// tools, such as the title, text. It keeps the tool results each session
// sends back.
type scriptedProvider struct {
	t     *testing.T
	mu    sync.Mutex
	child []string
	root  []string
	calls int
}

// results are the tool results the sub-agent and the root session sent
// back, in order.
func (p *scriptedProvider) results() (child, root []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.child...), append([]string(nil), p.root...)
}

// chatRequest is the part of a chat completions request the provider reads.
type chatRequest struct {
	Tools    []json.RawMessage `json:"tools"`
	Messages []struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	} `json:"messages"`
}

func (p *scriptedProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
		http.NotFound(w, r)
		return
	}
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var results []string
	for _, m := range req.Messages {
		if m.Role != "tool" {
			continue
		}
		var s string
		if json.Unmarshal(m.Content, &s) != nil {
			s = string(m.Content)
		}
		results = append(results, s)
	}
	sub := r.Header.Get("x-opencode-parent-session-id") != "" || r.Header.Get("x-parent-session-id") != ""
	p.mu.Lock()
	p.calls++
	n := p.calls
	switch {
	case len(req.Tools) == 0:
	case sub:
		p.child = results
	default:
		p.root = results
	}
	p.mu.Unlock()
	w.Header().Set("Content-Type", "text/event-stream")
	switch {
	case len(req.Tools) == 0:
		streamText(w, "Scripted session")
	case sub && len(results) == 0:
		streamTool(w, n, "bash", map[string]any{"command": flaiWrite, "description": "Move the story to review"})
	case sub:
		streamText(w, "The guard answered.")
	case len(results) == 0:
		streamTool(w, n, "task", map[string]any{"subagent_type": "general", "description": "Move the story", "prompt": "Run " + flaiWrite + " with bash."})
	case len(results) == 1:
		streamTool(w, n, "bash", map[string]any{"command": flaiWrite, "description": "Move the story to review"})
	default:
		streamText(w, "Done.")
	}
}

// streamChunk writes one chat completion chunk as a server-sent event.
func streamChunk(w http.ResponseWriter, delta map[string]any, finish any) {
	chunk := map[string]any{
		"id": "chatcmpl-scripted", "object": "chat.completion.chunk", "created": 0, "model": "scripted",
		"choices": []any{map[string]any{"index": 0, "delta": delta, "finish_reason": finish}},
	}
	if finish != nil {
		chunk["usage"] = map[string]any{"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
	}
	b, _ := json.Marshal(chunk)
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "data: %s\n\n", b)
	_ = bw.Flush()
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// streamText answers with text and ends the turn.
func streamText(w http.ResponseWriter, text string) {
	streamChunk(w, map[string]any{"role": "assistant", "content": text}, nil)
	streamChunk(w, map[string]any{}, "stop")
	fmt.Fprint(w, "data: [DONE]\n\n")
}

// streamTool answers with one call of tool with args.
func streamTool(w http.ResponseWriter, n int, tool string, args map[string]any) {
	b, _ := json.Marshal(args)
	streamChunk(w, map[string]any{"role": "assistant", "tool_calls": []any{map[string]any{
		"index": 0, "id": fmt.Sprintf("call_%d", n), "type": "function",
		"function": map[string]any{"name": tool, "arguments": string(b)},
	}}}, nil)
	streamChunk(w, map[string]any{}, "tool_calls")
	fmt.Fprint(w, "data: [DONE]\n\n")
}
