package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// keyedProject is a project with this key in a folder named after it.
func keyedProject(t *testing.T, parent, key string) string {
	t.Helper()
	root := filepath.Join(parent, key)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	m := "version: 1\nname: " + key + "\nkey: " + key + "\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func serveProjectRun(t *testing.T, cwd string, args ...string) (string, string, int) {
	t.Helper()
	return runWithApp(t, &app{cwd: cwd}, append([]string{"serve", "project"}, args...)...)
}

// S-0118: flai serve project add registers the project in a folder, the
// current one by default, as flai dashboard does, and refuses what cannot be
// served, naming why.
func TestServeProjectAddRegistersAndRefuses(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	harbour := keyedProject(t, home, "harbour")

	out, errOut, code := serveProjectRun(t, harbour, "add")
	if code != 0 {
		t.Fatalf("add: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "harbour ("+harbour+") is registered with flai serve, for the dashboard at http://localhost:4242") || !strings.Contains(out, "flai host started (pid 4242)") {
		t.Errorf("add said:\n%s", out)
	}
	dir := serve.DirFor(cfg)
	got, _ := dir.Projects()
	keyFile := filepath.Join(string(dir), "dashboard.agent-key")
	if len(got) != 1 || got[0] != (serve.Entry{Key: "harbour", Name: "harbour", Root: harbour, URL: "http://127.0.0.1:4242", KeyFile: keyFile}) {
		t.Errorf("registry: %+v", got)
	}
	if info, err := os.Stat(keyFile); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("the agent credential was not made: %v %v", info, err)
	}

	// by folder, relative to where flai runs
	lighthouse := keyedProject(t, home, "lighthouse")
	if out, errOut, code := serveProjectRun(t, home, "add", "lighthouse"); code != 0 {
		t.Fatalf("add lighthouse: %d %s %s", code, out, errOut)
	}
	if got, _ := dir.Projects(); len(got) != 2 || got[1].Root != lighthouse {
		t.Errorf("registry: %+v", got)
	}

	plain := t.TempDir()
	if _, errOut, code := serveProjectRun(t, plain, "add"); code == 0 || !strings.Contains(errOut, "has no system-flow.yaml") {
		t.Errorf("a folder with no manifest: %d %s", code, errOut)
	}
	keyless := filepath.Join(home, "keyless")
	_ = os.MkdirAll(keyless, 0o755)
	_ = os.WriteFile(filepath.Join(keyless, "system-flow.yaml"), []byte("version: 1\nname: keyless\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	if _, errOut, code := serveProjectRun(t, keyless, "add"); code == 0 || !strings.Contains(errOut, "has no key") {
		t.Errorf("a manifest with no key: %d %s", code, errOut)
	}
	twin := keyedProject(t, filepath.Join(home, "elsewhere"), "harbour")
	if _, errOut, code := serveProjectRun(t, twin, "add"); code == 0 || !strings.Contains(errOut, "the key harbour is already served for "+harbour) {
		t.Errorf("a key another project has: %d %s", code, errOut)
	}
	if got, _ := dir.Projects(); len(got) != 2 {
		t.Errorf("a refusal changed the registry: %+v", got)
	}
}

// S-0118: flai serve project remove unregisters by key or by folder and
// touches none of the project's files.
func TestServeProjectRemoveUnregistersAndTouchesNoFile(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	harbour := keyedProject(t, home, "harbour")
	lighthouse := keyedProject(t, home, "lighthouse")
	for _, p := range []string{harbour, lighthouse} {
		if _, errOut, code := serveProjectRun(t, p, "add"); code != 0 {
			t.Fatal(errOut)
		}
	}
	before := snapshot(t, home)

	// by key, from inside another registered project: the key wins
	out, errOut, code := serveProjectRun(t, lighthouse, "remove", "harbour")
	if code != 0 || !strings.Contains(out, "harbour ("+harbour+") is no longer served, and none of its files was touched") {
		t.Fatalf("remove by key: %d %s %s", code, out, errOut)
	}
	dir := serve.DirFor(cfg)
	if got, _ := dir.Projects(); len(got) != 1 || got[0].Key != "lighthouse" {
		t.Errorf("registry: %+v", got)
	}
	// by folder, from a folder inside it
	sub := filepath.Join(lighthouse, "docs")
	_ = os.MkdirAll(sub, 0o755)
	before[filepath.Join("lighthouse", "docs")] = "dir"
	if _, errOut, code := serveProjectRun(t, sub, "remove", "."); code != 0 {
		t.Fatalf("remove by folder: %s", errOut)
	}
	if got, _ := dir.Projects(); len(got) != 0 {
		t.Errorf("registry: %+v", got)
	}
	if after := snapshot(t, home); !equalMaps(before, after) {
		t.Errorf("files changed:\nbefore %v\nafter  %v", before, after)
	}
	if _, errOut, code := serveProjectRun(t, home, "remove", "harbour"); code == 0 || !strings.Contains(errOut, "no project harbour is served") {
		t.Errorf("removing what is not served: %d %s", code, errOut)
	}
}

// snapshot is every file under root but flai's own state, with its content.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	_ = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(root, p)
		if err != nil || rel == "." {
			return err
		}
		if rel == "serve" || rel == "cfg.json" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			out[rel] = "dir"
			return nil
		}
		data, _ := os.ReadFile(p)
		out[rel] = string(data)
		return nil
	})
	return out
}

func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// S-0118: flai serve project list shows each served project and how it is,
// the repositories offered for import, and the projects under the import
// folders that are not served; flai serve status names an unavailable one.
func TestServeProjectListSaysWhyAProjectIsNotShowing(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	harbour := keyedProject(t, home, "harbour")
	quay := keyedProject(t, home, "quay")
	for _, p := range []string{harbour, quay} {
		if _, errOut, code := serveProjectRun(t, p, "add"); code != 0 {
			t.Fatal(errOut)
		}
	}
	// quay moves away after it was registered
	if err := os.RemoveAll(quay); err != nil {
		t.Fatal(err)
	}
	// an import folder with a repository to offer and a project nobody serves
	git := filepath.Join(home, "git")
	_ = os.MkdirAll(filepath.Join(git, "widget", ".git"), 0o755)
	blog := keyedProject(t, git, "blog")
	c, _, err := config.Load(cfg)
	if err != nil {
		t.Fatal(err)
	}
	c.ImportRoots = []string{git}
	if err := config.Save(cfg, c); err != nil {
		t.Fatal(err)
	}
	// flai serve runs, with harbour connected
	now := time.Now().UTC().Format(time.RFC3339)
	st, _ := json.Marshal(serve.Status{PID: os.Getpid(), Version: "test", Started: now, Updated: now,
		Connections: map[string]channel.State{harbour: {URL: "http://127.0.0.1:4242", Connected: true, Since: "2026-09-26T05:00:00Z"}}})
	_ = os.WriteFile(filepath.Join(string(serve.DirFor(cfg)), "state.json"), st, 0o600)

	out, errOut, code := serveProjectRun(t, home, "list")
	if code != 0 {
		t.Fatalf("list: %d %s", code, errOut)
	}
	for _, want := range []string{
		"harbour  http://127.0.0.1:4242  connected since 2026-09-26T05:00:00Z",
		"quay  http://127.0.0.1:4242  not served: the folder is gone",
		"offered for import on the board (flai serve import):\n  widget  " + filepath.Join(git, "widget"),
		"projects under the import folders that are not served:\n  blog  " + blog + "\n    flai serve project add " + blog,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("list lacks %q:\n%s", want, out)
		}
	}

	js, _, _ := serveProjectRun(t, home, "list", "--json")
	var l projectListing
	if err := json.Unmarshal([]byte(js), &l); err != nil {
		t.Fatalf("json: %v\n%s", err, js)
	}
	if !l.Running || len(l.Served) != 2 || l.Served[0].State != "connected" || l.Served[1].State != "unavailable" || l.Served[1].Reason != "the folder is gone" {
		t.Errorf("served: %+v", l.Served)
	}
	if len(l.Candidates) != 1 || l.Candidates[0].Name != "widget" || len(l.Unserved) != 1 || l.Unserved[0].Key != "blog" {
		t.Errorf("candidates %+v, unserved %+v", l.Candidates, l.Unserved)
	}

	out, _, _ = runIn(t, home, "serve", "status")
	if !strings.Contains(out, "quay  http://127.0.0.1:4242  not served: the folder is gone") {
		t.Errorf("status:\n%s", out)
	}
	js, _, _ = runIn(t, home, "serve", "status", "--json")
	var ss serveStatus
	if err := json.Unmarshal([]byte(js), &ss); err != nil || ss.Unavailable[quay] != "the folder is gone" {
		t.Errorf("status json: %v %+v", err, ss.Unavailable)
	}
}
