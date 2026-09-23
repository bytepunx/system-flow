package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai serve keeps its registry beside the config file in use, so a test
// (and a repository with its own config) never touches the operator's home.
func TestServeStatusReadsTheRegistryBesideTheConfig(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	out, errOut, code := runIn(t, home, "serve", "status")
	if code != 0 || !strings.Contains(out, "not running") || !strings.Contains(out, "no projects registered") {
		t.Fatalf("empty: %d %s %s", code, out, errOut)
	}
	dir := serve.DirFor(cfg)
	if err := dir.Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: "/p/harbour", URL: "http://127.0.0.1:4242", KeyFile: "/p/harbour/k"}); err != nil {
		t.Fatal(err)
	}
	out, _, _ = runIn(t, home, "serve", "status")
	if !strings.Contains(out, "harbour  http://127.0.0.1:4242  not connected") {
		t.Errorf("status: %s", out)
	}
	out, _, _ = runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(out), &st); err != nil || st.Running || len(st.Projects) != 1 || st.Dir != filepath.Join(home, "serve") {
		t.Errorf("json: %v %+v", err, st)
	}
	if out, _, code := runIn(t, home, "serve", "stop"); code != 0 || !strings.Contains(out, "not running") {
		t.Errorf("stop with nothing running: %d %s", code, out)
	}
}

// S-0096: flai serve status says where each served project's MCP server
// listens, the one flai serve keeps running, so an agent can be pointed at it.
func TestServeStatusNamesEachProjectsMCPServer(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	if err := serve.DirFor(cfg).Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: root, URL: "http://127.0.0.1:4242", KeyFile: filepath.Join(root, "k")}); err != nil {
		t.Fatal(err)
	}
	if out, _, _ := runIn(t, home, "serve", "status"); strings.Contains(out, "mcp:") {
		t.Errorf("an MCP server named while none runs: %s", out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := writeMCPState(repo, mcpState{PID: os.Getpid(), Addr: "127.0.0.1:4244", URL: "http://127.0.0.1:4244/mcp", Started: now, Updated: now}); err != nil {
		t.Fatal(err)
	}
	out, _, _ := runIn(t, home, "serve", "status")
	if !strings.Contains(out, "mcp: http://127.0.0.1:4244/mcp (pid ") {
		t.Errorf("status: %s", out)
	}
	js, _, _ := runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(js), &st); err != nil || st.MCP[root].URL != "http://127.0.0.1:4244/mcp" {
		t.Errorf("json: %v %+v", err, st.MCP)
	}
}
