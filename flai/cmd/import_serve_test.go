package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// fakeRunning writes a state file in dir, as a running flai host or flai
// serve would, with this process's PID so that it reads as alive.
func fakeRunning(t *testing.T, dir string) {
	t.Helper()
	now := time.Now().UTC().Format(time.RFC3339)
	data, _ := json.Marshal(map[string]any{"pid": os.Getpid(), "started": now, "updated": now})
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "state.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func importForServe(t *testing.T, args ...string) (root, out, errOut string, code int) {
	t.Helper()
	if _, err := os.Stat(filepath.Join("..", "..", "template", "template.yaml")); err != nil {
		t.Skip("prototype template not present")
	}
	root = legacyRepo(t)
	out, errOut, code = runIn(t, ".", append([]string{"import", root, "--template", "../../template", "--yes", "--var", "project_key=legacy"}, args...)...)
	return root, out, errOut, code
}

// S-0117: with no flai host running, nothing is registered and the import
// names the one command that serves the project.
func TestImportWithNoHostSaysWhatServesIt(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root, out, errOut, code := importForServe(t)
	if code != 0 {
		t.Fatalf("import: %d %s", code, errOut)
	}
	if want := "dashboard: not served yet (flai host is not running); flai dashboard in " + root + " serves it"; !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	if projects, _ := serve.DirFor(cfg).Projects(); len(projects) != 0 {
		t.Errorf("registered with no host: %+v", projects)
	}
}

// S-0117: with a flai host running, the import registers the project with
// flai serve as flai dashboard would, and says where the dashboard shows it,
// with or without --commit, and committed or not: the legacy repository's
// go module has no Go files, so go test fails and --commit exits 5.
func TestImportWithAHostRegistersTheProject(t *testing.T) {
	for args, exit := range map[string]int{"": 0, "--commit": exitImportNotCommitted} {
		t.Run("import "+args, func(t *testing.T) {
			cfg := filepath.Join(t.TempDir(), "cfg.json")
			t.Setenv("FLAI_CONFIG", cfg)
			fakeRunning(t, string(host.DirFor(cfg)))
			dir := serve.DirFor(cfg)
			fakeRunning(t, string(dir))

			root, out, errOut, code := importForServe(t, strings.Fields(args)...)
			if code != exit {
				t.Fatalf("import: %d %s\n%s", code, errOut, out)
			}
			if want := "dashboard: served by the host flai at http://localhost:4242"; !strings.Contains(out, want) {
				t.Errorf("missing %q in:\n%s", want, out)
			}
			projects, err := dir.Projects()
			if err != nil || len(projects) != 1 {
				t.Fatalf("registry: %+v %v", projects, err)
			}
			p := projects[0]
			want := serve.Entry{Key: "legacy", Name: filepath.Base(root), Root: root, URL: "http://127.0.0.1:4242", KeyFile: filepath.Join(string(dir), "dashboard.agent-key")}
			if p != want {
				t.Errorf("entry:\n got %+v\nwant %+v", p, want)
			}
			if b, err := os.ReadFile(p.KeyFile); err != nil || strings.TrimSpace(string(b)) == "" {
				t.Errorf("credential not written: %v", err)
			}
		})
	}
}

// A host whose flai serve was stopped: registered, and told how to start it.
func TestImportWithServeStoppedSaysToStartIt(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	fakeRunning(t, string(host.DirFor(cfg)))

	_, out, errOut, code := importForServe(t, "--json")
	if code != 0 {
		t.Fatalf("import: %d %s", code, errOut)
	}
	var v struct {
		Serve importServed `json:"serve"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if !v.Serve.Served || v.Serve.Next != "flai serve start" || v.Serve.Dashboard != "http://localhost:4242" {
		t.Errorf("serve: %+v", v.Serve)
	}
	if projects, _ := serve.DirFor(cfg).Projects(); len(projects) != 1 {
		t.Errorf("not registered: %+v", projects)
	}
}
