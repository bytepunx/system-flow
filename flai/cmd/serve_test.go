package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/serve"
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
