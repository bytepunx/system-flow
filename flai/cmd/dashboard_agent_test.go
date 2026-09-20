package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// S-0072: flai dashboard hands the container the agent credential, registers
// the project with flai serve, and reports the connection.
func TestDashboardHandsOverTheAgentCredentialAndRegistersTheProject(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: Harbour\nkey: harbour\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n"), 0o644)
	f := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatalf("%s %s", out, errOut)
	}
	keyFile := filepath.Join(root, ".flai-cache", "dashboard.agent-key")
	info, err := os.Stat(keyFile)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("agent credential: %v %v", info, err)
	}
	key, _ := os.ReadFile(keyFile)
	token, _ := os.ReadFile(filepath.Join(root, ".flai-cache", "dashboard.token"))
	if len(strings.TrimSpace(string(key))) < 32 || string(key) == string(token) {
		t.Error("the agent credential must be its own secret, not the login token")
	}
	run := ""
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker run ") {
			run = c
		}
	}
	for _, want := range []string{
		"--mount type=bind,source=" + keyFile + ",target=/run/secrets/flaiover_agent_key,readonly",
		"--env FLAIOVER_AGENT_KEY_FILE=/run/secrets/flaiover_agent_key",
	} {
		if !strings.Contains(run, want) {
			t.Errorf("docker run lacks %q:\n%s", want, run)
		}
	}
	if strings.Contains(run, strings.TrimSpace(string(key))) || strings.Contains(out, strings.TrimSpace(string(key))) {
		t.Error("the credential itself appears in the command line or the output")
	}
	if !strings.Contains(out, "host flai: flai serve started (pid 4242)") {
		t.Errorf("output does not say flai serve was started:\n%s", out)
	}
	entries, _ := serve.DirFor(cfg).Projects()
	if len(entries) != 1 || entries[0].Key != "harbour" || entries[0].Root != root || entries[0].URL != "http://127.0.0.1:5555" || entries[0].KeyFile != keyFile {
		t.Errorf("registry: %+v", entries)
	}

	// status, as a running flai serve with this project connected would have it
	f.running["flaiover-harbour"] = true
	st := serve.Status{PID: os.Getpid(), Version: "test", Started: "2026-09-20T07:00:00Z", Updated: time.Now().UTC().Format(time.RFC3339),
		Connections: map[string]channel.State{root: {URL: "http://127.0.0.1:5555", Connected: true, Since: "2026-09-20T07:00:01Z"}}}
	data, _ := json.Marshal(st)
	_ = os.WriteFile(filepath.Join(string(serve.DirFor(cfg)), "state.json"), data, 0o600)
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "host flai: connected since 2026-09-20T07:00:01Z") || !strings.Contains(out, "serving 1 project(s)") {
		t.Errorf("status:\n%s", out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status", "--json")
	var js struct {
		HostFlai hostFlaiStatus `json:"host_flai"`
	}
	if err := json.Unmarshal([]byte(out), &js); err != nil || !js.HostFlai.Connected || !js.HostFlai.Registered || js.HostFlai.PID != os.Getpid() {
		t.Errorf("status json: %v %+v\n%s", err, js, out)
	}

	// stop leaves flai serve alone and stops it dialling this project
	out, _, code = runWith(t, root, f, "dashboard", "stop")
	if code != 0 {
		t.Fatal(out)
	}
	if entries, _ := serve.DirFor(cfg).Projects(); len(entries) != 0 {
		t.Errorf("still registered after stop: %+v", entries)
	}
	if _, alive := serve.DirFor(cfg).ReadStatus(time.Now()); !alive {
		t.Error("flai dashboard stop must not stop flai serve")
	}
	for _, c := range f.calls {
		if strings.Contains(c, "serve stop") {
			t.Errorf("stop touched flai serve: %s", c)
		}
	}
}

func TestDashboardNoServeRegistersNothing(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	f := &fakeRunner{images: map[string]bool{"ghcr.io/bytepunx/flaiover:latest": true}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "--no-serve")
	if code != 0 {
		t.Fatalf("%s %s", out, errOut)
	}
	if entries, _ := serve.DirFor(cfg).Projects(); len(entries) != 0 {
		t.Errorf("registered with --no-serve: %+v", entries)
	}
	if !strings.Contains(out, "host flai: not started (--no-serve)") {
		t.Errorf("output:\n%s", out)
	}
}

func TestDialURLFollowsTheBindAddress(t *testing.T) {
	for bind, want := range map[string]string{"0.0.0.0": "http://127.0.0.1:4242", "": "http://127.0.0.1:4242", "127.0.0.1": "http://127.0.0.1:4242", "192.168.1.5": "http://192.168.1.5:4242", "::1": "http://[::1]:4242"} {
		if got := (dashboardSettings{Bind: bind, Port: 4242}).dialURL(); got != want {
			t.Errorf("bind %q: %s, want %s", bind, got, want)
		}
	}
}
