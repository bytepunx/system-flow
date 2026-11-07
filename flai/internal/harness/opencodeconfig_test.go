package harness

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/protected"
)

// envOf is the value of the variable name in env, and whether env sets it.
func envOf(env []string, name string) (string, bool) {
	for _, kv := range env {
		if k, v, ok := strings.Cut(kv, "="); ok && k == name {
			return v, true
		}
	}
	return "", false
}

// objectInOrder is a JSON object's keys in the order written, and their
// values.
func objectInOrder(t *testing.T, data string) ([]string, map[string]json.RawMessage) {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(data))
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		t.Fatalf("%s: not a JSON object: %v %v", data, tok, err)
	}
	var keys []string
	values := map[string]json.RawMessage{}
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("%s: %v", data, err)
		}
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			t.Fatalf("%s: %v", data, err)
		}
		keys = append(keys, tok.(string))
		values[tok.(string)] = v
	}
	return keys, values
}

// editRules are the edit rules of an OPENCODE_PERMISSION, in order, each
// pattern and its action, and the permission's top-level keys.
func editRules(t *testing.T, permission string) (top []string, rules [][2]string) {
	t.Helper()
	top, values := objectInOrder(t, permission)
	patterns, actions := objectInOrder(t, string(values["edit"]))
	for _, p := range patterns {
		var a string
		if err := json.Unmarshal(actions[p], &a); err != nil {
			t.Fatalf("edit %q: %s is not an action", p, actions[p])
		}
		rules = append(rules, [2]string{p, a})
	}
	return top, rules
}

// S-0406, ADR-0260 §1: a start gives OpenCode flai's MCP server in
// OPENCODE_CONFIG_CONTENT and nothing else: a local server, flai mcp with
// the session's agent, on, with mcphttp.ClientTimeout as its timeout; with
// no agent's name it passes no --agent, as claude-code's does.
func TestOpenCodeGivesFlaisMCPServerInItsConfigContent(t *testing.T) {
	r := req(openCodeAgent("anthropic/claude-opus-5-5"))
	s, err := (openCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	content, ok := envOf(s.Env, "OPENCODE_CONFIG_CONTENT")
	want := `{"mcp":{"flai":{"type":"local","command":["/usr/local/bin/flai","mcp","--agent","agent-S-0104"],"enabled":true,"timeout":1800000}}}`
	if !ok || content != want {
		t.Errorf("OPENCODE_CONFIG_CONTENT %q (set %v),\nwant %q", content, ok, want)
	}

	r.Name = ""
	content, err = openCodeConfigContent(r)
	want = `{"mcp":{"flai":{"type":"local","command":["/usr/local/bin/flai","mcp"],"enabled":true,"timeout":1800000}}}`
	if err != nil || content != want {
		t.Errorf("no agent's name: %q, %v;\nwant %q", content, err, want)
	}
}

// S-0406, ADR-0260 §2: under deny_protected a start sets OPENCODE_PERMISSION
// to an edit object with a deny rule for every path of OpenCode's protected
// list, a folder at the root and below it, a file at the root and below it,
// in the list's order, and writes no file into the root; without it, the
// variable is not set.
func TestOpenCodeDeniesTheProtectedPathsUnderDenyProtected(t *testing.T) {
	r := req(openCodeAgent("anthropic/claude-opus-5-5"))
	r.Root = t.TempDir()
	s, err := (openCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := envOf(s.Env, "OPENCODE_PERMISSION"); ok {
		t.Errorf("no deny_protected: OPENCODE_PERMISSION=%s, want it unset", v)
	}

	s, err = (openCode{}).Start(r, Host{DenyProtected: true})
	if err != nil {
		t.Fatal(err)
	}
	permission, ok := envOf(s.Env, "OPENCODE_PERMISSION")
	if !ok {
		t.Fatalf("deny_protected: env %q, want OPENCODE_PERMISSION", s.Env)
	}
	top, rules := editRules(t, permission)
	if !slices.Equal(top, []string{"edit"}) {
		t.Errorf("top-level keys %q, want edit alone", top)
	}
	l := protected.For(protected.OpenCode)
	if len(rules) != 2*(len(l.Folders)+len(l.Files)) {
		t.Errorf("%d rules for %d folders and %d files, want two each: %q", len(rules), len(l.Folders), len(l.Files), rules)
	}
	var patterns []string
	for _, rule := range rules {
		if rule[1] != "deny" {
			t.Errorf("rule %q, want deny", rule)
		}
		patterns = append(patterns, rule[0])
	}
	for _, p := range []string{".git/*", "*/.git/*", ".claude/*", "*/.claude/*", ".config/git/*", "*/.config/git/*", ".opencode/*", "*/.opencode/*",
		"opencode.json", "*/opencode.json", "opencode.jsonc", "*/opencode.jsonc", "AGENTS.md", "*/AGENTS.md", ".mcp.json", "*/.mcp.json"} {
		if !slices.Contains(patterns, p) {
			t.Errorf("no deny rule for %q in %q", p, patterns)
		}
	}
	if patterns[0] != ".git/*" || patterns[1] != "*/.git/*" {
		t.Errorf("rules begin %q, want the list's order, .git first", patterns[:2])
	}
	again, err := (openCode{}).Start(r, Host{DenyProtected: true})
	if p, _ := envOf(again.Env, "OPENCODE_PERMISSION"); err != nil || p != permission {
		t.Errorf("a second start: %q, %v; want the same permission, %q", p, err, permission)
	}
	if entries, err := os.ReadDir(r.Root); err != nil || len(entries) != 0 {
		t.Errorf("the root holds %v, %v after the starts; want no file from flai", entries, err)
	}

	r.Root = ""
	if _, err := (openCode{}).Start(r, Host{DenyProtected: true}); err == nil || !strings.Contains(err.Error(), "the start names no project root, where flai reads the project's opencode.json and opencode.jsonc") {
		t.Errorf("no root: %v", err)
	}
}

// ADR-0260 §2: a root folder is denied as <path>/* and a root file as
// itself, below the other paths.
func TestOpenCodeDeniesARootPathAsItself(t *testing.T) {
	got := openCodeDenies(protected.List{Folders: []string{".git"}, Files: []string{"AGENTS.md"}, RootFolders: []string{"flai/images"}, RootFiles: []string{"scripts/agent-image.sh"}})
	want := []string{".git/*", "*/.git/*", "AGENTS.md", "*/AGENTS.md", "flai/images/*", "scripts/agent-image.sh"}
	if !slices.Equal(got, want) {
		t.Errorf("denies %q,\nwant %q", got, want)
	}
}

// S-0406, ADR-0260 §3: a project whose opencode.json or opencode.jsonc gives
// its permission, or its edit, as one action keeps it as the first edit
// rule, *, and a whole permission as the top-level * too, as OpenCode merges
// the files: opencode.jsonc over opencode.json, a string replacing what came
// before whole, an object replacing a string and keeping an earlier edit it
// does not give.
func TestOpenCodeKeepsTheProjectsPermissionGivenAsOneAction(t *testing.T) {
	for _, c := range []struct {
		what        string
		json, jsonc string // the files' content, absent when empty
		empty       bool   // opencode.json is there, and empty
		top, edit   string // the actions kept: top-level and edit's *, empty for none
	}{
		{what: "no files"},
		{what: "empty file", empty: true},
		{what: "no permission", json: `{"model":"anthropic/claude-opus-5-5"}`},
		{what: "permission one action", json: `{"permission":"ask"}`, top: "ask", edit: "ask"},
		{what: "edit one action", json: `{"permission":{"bash":"allow","edit":"ask"}}`, edit: "ask"},
		{what: "edit an object", json: `{"permission":{"edit":{"*.md":"allow"}}}`},
		{what: "jsonc over json", json: `{"permission":"deny"}`, jsonc: `{"permission":"allow"}`, top: "allow", edit: "allow"},
		{what: "jsonc's edit over json's", json: `{"permission":{"edit":"allow"}}`, jsonc: `{"permission":{"edit":"ask"}}`, edit: "ask"},
		{what: "an object replaces an action", json: `{"permission":"ask"}`, jsonc: `{"permission":{"bash":"allow"}}`},
		{what: "an object keeps an earlier edit", json: `{"permission":{"edit":"ask"}}`, jsonc: `{"permission":{"bash":"deny"}}`, edit: "ask"},
		{what: "an action replaces an earlier edit", json: `{"permission":{"edit":"ask"}}`, jsonc: `{"permission":"allow"}`, top: "allow", edit: "allow"},
	} {
		fsys := fstest.MapFS{}
		if c.empty {
			fsys["opencode.json"] = &fstest.MapFile{}
		}
		if c.json != "" {
			fsys["opencode.json"] = &fstest.MapFile{Data: []byte(c.json)}
		}
		if c.jsonc != "" {
			fsys["opencode.jsonc"] = &fstest.MapFile{Data: []byte(c.jsonc)}
		}
		permission, err := openCodePermission("/p/flow", fsys)
		if err != nil {
			t.Errorf("%s: %v", c.what, err)
			continue
		}
		top, rules := editRules(t, permission)
		wantTop := []string{"edit"}
		if c.top != "" {
			wantTop = []string{"*", "edit"}
			var action string
			_, values := objectInOrder(t, permission)
			if json.Unmarshal(values["*"], &action) != nil || action != c.top {
				t.Errorf("%s: top-level * is %s, want %q", c.what, values["*"], c.top)
			}
		}
		if !slices.Equal(top, wantTop) {
			t.Errorf("%s: top-level keys %q, want %q", c.what, top, wantTop)
		}
		first := rules[0]
		switch {
		case c.edit != "" && first != [2]string{"*", c.edit}:
			t.Errorf("%s: edit begins %q, want * %s", c.what, first, c.edit)
		case c.edit != "" && rules[1] != [2]string{".git/*", "deny"}:
			t.Errorf("%s: the denies begin %q, want them after *", c.what, rules[1])
		case c.edit == "" && first != [2]string{".git/*", "deny"}:
			t.Errorf("%s: edit begins %q, want the denies alone", c.what, first)
		}
	}
}

// ADR-0260 §3: the project's opencode.json at the root a start runs in is
// read through the start, kept, and left as it was.
func TestOpenCodeReadsTheProjectsConfigAtItsRoot(t *testing.T) {
	r := req(openCodeAgent("anthropic/claude-opus-5-5"))
	r.Root = t.TempDir()
	config := []byte(`{"permission":"ask","mcp":{"flai":{"type":"local","command":["flai","mcp"]}}}`)
	path := filepath.Join(r.Root, "opencode.json")
	if err := os.WriteFile(path, config, 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := (openCode{}).Start(r, Host{DenyProtected: true})
	if err != nil {
		t.Fatal(err)
	}
	permission, _ := envOf(s.Env, "OPENCODE_PERMISSION")
	if !strings.HasPrefix(permission, `{"*":"ask","edit":{"*":"ask",".git/*":"deny",`) {
		t.Errorf("OPENCODE_PERMISSION %s, want the project's ask kept first", permission)
	}
	if read, err := os.ReadFile(path); err != nil || !bytes.Equal(read, config) {
		t.Errorf("the project's opencode.json after the start: %s, %v; want it as it was", read, err)
	}
}

// ADR-0260 §3: JSONC as OpenCode reads it, in either file, gives up its
// comments and trailing commas, and keeps what looks like them in a string.
func TestOpenCodeReadsJSONCWithCommentsAndTrailingCommas(t *testing.T) {
	jsonc := `{
  // the project's own
  "$schema": "https://opencode.ai/config.json", /* a URL holds // */
  "permission": {
    "edit": "ask", // ask before each write
    "bash": { "git *": "allow", },
  },
  "instructions": ["a /* not a comment */ b", "c\"//d", ", }",],
}`
	var got struct {
		Schema       string   `json:"$schema"`
		Instructions []string `json:"instructions"`
	}
	if err := json.Unmarshal(stripJSONC([]byte(jsonc)), &got); err != nil {
		t.Fatalf("stripJSONC: %v\n%s", err, stripJSONC([]byte(jsonc)))
	}
	if want := []string{"a /* not a comment */ b", `c"//d`, ", }"}; got.Schema != "https://opencode.ai/config.json" || !slices.Equal(got.Instructions, want) {
		t.Errorf("stripped %+v, want the schema's URL and instructions %q whole", got, want)
	}
	for _, name := range []string{"opencode.json", "opencode.jsonc"} {
		kept, err := openCodeProjectPermission("/p/flow", fstest.MapFS{name: {Data: []byte(jsonc)}})
		if err != nil || kept != (openCodeKept{edit: "ask"}) {
			t.Errorf("%s: %+v, %v; want edit ask kept", name, kept, err)
		}
	}
}

// ADR-0260 §3: a project file that is not JSONC refuses a start under
// deny_protected, naming it and saying to fix it; without deny_protected
// flai reads none.
func TestOpenCodeRefusesAProjectConfigThatDoesNotParse(t *testing.T) {
	for _, bad := range []string{`{"permission": "ask"`, `{"permission": "ask"} /* never closed`, `["permission"]`, " \n"} {
		_, err := openCodePermission("/p/flow", fstest.MapFS{"opencode.json": {Data: []byte(`{}`)}, "opencode.jsonc": {Data: []byte(bad)}})
		if err == nil || !strings.Contains(err.Error(), "the project's /p/flow/opencode.jsonc is not a JSONC object, which OpenCode would refuse too") || !strings.HasSuffix(err.Error(), "; fix it") {
			t.Errorf("%q: %v", bad, err)
		}
	}

	r := req(openCodeAgent("anthropic/claude-opus-5-5"))
	r.Root = t.TempDir()
	if err := os.WriteFile(filepath.Join(r.Root, "opencode.json"), []byte(`{"permission":`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := (openCode{}).Start(r, Host{DenyProtected: true}); err == nil || !strings.Contains(err.Error(), filepath.Join(r.Root, "opencode.json")) {
		t.Errorf("deny_protected: %v; want the start refused, naming the file", err)
	}
	if _, err := (openCode{}).Start(r, Host{}); err != nil {
		t.Errorf("no deny_protected: %v; want the file unread", err)
	}
}

// S-0406: an agent's effort is OpenCode's model variant, passed as
// --variant after the model; an effort OpenCode names no variant for is
// refused, saying the ones it does.
func TestOpenCodeTakesEffortAsItsVariant(t *testing.T) {
	r := req(&manifest.Agent{Harness: OpenCode, Model: "anthropic/claude-opus-5-5", Config: map[string]string{"effort": "high"}})
	r.Session, r.Answered = "ses_0b1c2d3e4f5aAbCdEfGhIjKlMn", "TH-0001"
	s, err := (openCode{}).Start(r, Host{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"opencode", "run", "--format", "json", "--model", "anthropic/claude-opus-5-5", "--variant", "high", "--session", r.Session, mustPrompt(t, r)}
	if !slices.Equal(s.Argv, want) {
		t.Errorf("argv %q,\nwant %q", s.Argv, want)
	}
	for _, effort := range []string{"none", "minimal", "low", "medium", "xhigh", "max"} {
		r.Agent.Config["effort"] = effort
		if s, err := (openCode{}).Start(r, Host{}); err != nil || after(s.Argv, "--variant") != effort {
			t.Errorf("effort %s: %q, %v", effort, s.Argv, err)
		}
	}
	r.Agent.Config["effort"] = "high --auto"
	if _, err := (openCode{}).Start(r, Host{}); err == nil || !strings.Contains(err.Error(), `opencode config effort="high --auto": it must be none, minimal, low, medium, high, xhigh, or max`) {
		t.Errorf("an effort with a flag: %v", err)
	}
}

// S-0406: opencode's adapter reads the session a run made from its log, the
// sessionID of its first event, past lines that are no event, and only an
// ID OpenCode makes; claude-code, which names its own sessions, does not.
func TestOpenCodeReadsTheSessionFromItsLog(t *testing.T) {
	reader, ok := Adapters[OpenCode].(SessionReader)
	if !ok {
		t.Fatal("opencode's adapter is no SessionReader")
	}
	if _, ok := Adapters[ClaudeCode].(SessionReader); ok {
		t.Error("claude-code's adapter is a SessionReader; it names its own sessions")
	}
	const stepStart = `{"type":"step_start","timestamp":1794045600000,"sessionID":"ses_0b1c2d3e4f5aAbCdEfGhIjKlMn","part":{"id":"prt_9a8b7c6d5e4fAbCdEfGhIjKlM1","sessionID":"ses_0b1c2d3e4f5aAbCdEfGhIjKlMn","messageID":"msg_9a8b7c6d5e40AbCdEfGhIjKlMn","type":"step-start","snapshot":"4b825dc642cb6eb9a060e54bf8d69288fbee4904"}}`
	const errorEvent = `{"type":"error","timestamp":1794045605000,"sessionID":"ses_1c2d3e4f5a6bBcDeFgHiJkLmNo","error":{"name":"APIError","data":{"message":"Unauthorized","statusCode":401,"isRetryable":false}}}`
	for _, c := range []struct {
		what, log, want string
	}{
		{"a step first", stepStart + "\n" + errorEvent + "\n", "ses_0b1c2d3e4f5aAbCdEfGhIjKlMn"},
		{"an error first", errorEvent + "\n" + stepStart, "ses_1c2d3e4f5a6bBcDeFgHiJkLmNo"},
		{"past stderr and a JSON line with no type", "permission requested: edit (.git/config); auto-rejecting\n{\"level\":\"warn\"}\n\n" + stepStart + "\n", "ses_0b1c2d3e4f5aAbCdEfGhIjKlMn"},
		{"empty", "", ""},
		{"no event", "opencode: command not found\n", ""},
		{"not OpenCode's ID", `{"type":"step_start","sessionID":"6f1c2d3e-4f5a-4b6c-8d7e-9f0a1b2c3d4e"}` + "\n" + stepStart, ""},
		{"no ID", `{"type":"error","error":{"name":"UnknownError"}}`, ""},
	} {
		got, ok := reader.Session([]byte(c.log))
		if got != c.want || ok != (c.want != "") {
			t.Errorf("%s: %q, %v; want %q", c.what, got, ok, c.want)
		}
	}
}
