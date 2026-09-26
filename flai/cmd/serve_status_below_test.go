package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// belowFolder is a folder with two git repositories with a system-flow.yaml:
// blog with a key, and nokey without one.
func belowFolder(t *testing.T) string {
	t.Helper()
	folder := t.TempDir()
	for name, key := range map[string]string{"blog": "key: blog\n", "nokey": ""} {
		root := filepath.Join(folder, name)
		if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
		m := "version: 1\nname: " + name + "\n" + key + "layout:\n  design: design\n  docs: docs\n  wip: wip\n"
		if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return folder
}

// S-0117: with flai serve not running, status still lists the projects below
// the folders named for import, and why each is not served.
func TestServeStatusListsTheProjectsBelowTheImportRootsNotServed(t *testing.T) {
	home := t.TempDir()
	t.Setenv("FLAI_CONFIG", filepath.Join(home, "cfg.json"))
	folder := belowFolder(t)
	if _, errOut, code := runIn(t, home, "serve", "import", "add", folder); code != 0 {
		t.Fatalf("add: %s", errOut)
	}
	out, _, _ := runIn(t, home, "serve", "status")
	if want := filepath.Join(folder, "blog") + "\n    no dashboard is known to serve it for (flai dashboard starts one)\n"; !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	// a dashboard known: blog would be served, but flai serve does not run
	if err := serve.DirFor(filepath.Join(home, "cfg.json")).Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: "/p/harbour", URL: "http://127.0.0.1:4242", KeyFile: "/p/harbour/k"}); err != nil {
		t.Fatal(err)
	}
	out, _, _ = runIn(t, home, "serve", "status")
	for _, want := range []string{
		"not served:\n",
		filepath.Join(folder, "blog") + "\n    flai serve is not running\n",
		filepath.Join(folder, "nokey") + "\n    its system-flow.yaml has no key (flai check says how to add one)\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	js, _, _ := runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(js), &st); err != nil || len(st.Unserved) != 2 {
		t.Errorf("json: %v %s", err, js)
	}
}

// With flai serve running, status shows what it serves from the folders
// named for import and what it does not, from its state.
func TestServeStatusShowsWhatServeServesBelowTheImportRoots(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	now := time.Now().UTC().Format(time.RFC3339)
	st := serve.Status{PID: os.Getpid(), Version: "test", Started: now, Updated: now,
		ImportProjects: []serve.Entry{{Key: "blog", Root: "/g/blog", URL: "http://127.0.0.1:4242"}},
		Unserved:       []serve.Found{{Entry: serve.Entry{Root: "/g/twin"}, Reason: "its key blog is served already, for /g/blog"}}}
	data, _ := json.Marshal(st)
	dir := serve.DirFor(cfg)
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(string(dir), "state.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	out, _, _ := runIn(t, home, "serve", "status")
	for _, want := range []string{
		"served from the folders named for import (flai serve import list):\n  blog  http://127.0.0.1:4242  not connected\n    /g/blog\n",
		"not served:\n  /g/twin\n    its key blog is served already, for /g/blog\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}
