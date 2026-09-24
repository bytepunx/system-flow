package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/host"
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
	keyFile := filepath.Join(string(serve.DirFor(cfg)), "dashboard.agent-key")
	info, err := os.Stat(keyFile)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("agent credential: %v %v", info, err)
	}
	key, _ := os.ReadFile(keyFile)
	token, _ := os.ReadFile(filepath.Join(string(serve.DirFor(cfg)), "dashboard.token"))
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
	// S-0106: the dashboard boots flai host, which runs flai serve
	if !strings.Contains(out, "host flai: flai host started (pid 4242); it runs flai serve") {
		t.Errorf("output does not say flai host was started:\n%s", out)
	}
	entries, _ := serve.DirFor(cfg).Projects()
	if len(entries) != 1 || entries[0].Key != "harbour" || entries[0].Root != root || entries[0].URL != "http://127.0.0.1:5555" || entries[0].KeyFile != keyFile {
		t.Errorf("registry: %+v", entries)
	}

	// status, as a running flai serve with this project connected would have it
	f.running["flaiover"] = true
	st := serve.Status{PID: os.Getpid(), Version: "test", Started: "2026-09-20T07:00:00Z", Updated: time.Now().UTC().Format(time.RFC3339),
		Connections: map[string]channel.State{root: {URL: "http://127.0.0.1:5555", Connected: true, Since: "2026-09-20T07:00:01Z"}}}
	data, _ := json.Marshal(st)
	_ = os.WriteFile(filepath.Join(string(serve.DirFor(cfg)), "state.json"), data, 0o600)
	// and flai host, which runs it (S-0106)
	hst, _ := json.Marshal(host.Status{PID: os.Getpid(), Version: "test", Updated: time.Now().UTC().Format(time.RFC3339)})
	_ = os.MkdirAll(string(host.DirFor(cfg)), 0o700)
	_ = os.WriteFile(filepath.Join(string(host.DirFor(cfg)), "state.json"), hst, 0o600)
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "host flai: connected since 2026-09-20T07:00:01Z") || !strings.Contains(out, "serving 1 project(s)") {
		t.Errorf("status:\n%s", out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status", "--json")
	var js struct {
		HostFlai hostFlaiStatus `json:"host_flai"`
	}
	if err := json.Unmarshal([]byte(out), &js); err != nil || !js.HostFlai.Connected || !js.HostFlai.Registered || js.HostFlai.PID != os.Getpid() || !js.HostFlai.HostRunning || js.HostFlai.HostPID != os.Getpid() {
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

// S-0080: one dashboard container serves every project the host flai serves,
// with one shared login token and one shared agent credential.
func TestDashboardSharesOneContainerAcrossProjects(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	serveDir := string(serve.DirFor(cfg))

	first := tempProject(t)
	_ = os.WriteFile(filepath.Join(first, "system-flow.yaml"), []byte("version: 1\nname: First\nkey: first\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n"), 0o644)
	second := tempProject(t)
	_ = os.WriteFile(filepath.Join(second, "system-flow.yaml"), []byte("version: 1\nname: Second\nkey: second\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n"), 0o644)

	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, first, f, "dashboard")
	if code != 0 {
		t.Fatalf("first: %s %s", out, errOut)
	}
	if n := strings.Count(strings.Join(f.calls, "\n"), "docker run "); n != 1 {
		t.Fatalf("one container started: %d runs", n)
	}
	out, errOut, code = runWith(t, second, f, "dashboard")
	if code != 0 {
		t.Fatalf("second: %s %s", out, errOut)
	}
	if n := strings.Count(strings.Join(f.calls, "\n"), "docker run "); n != 1 {
		t.Errorf("a second flai dashboard started a second container: %s", strings.Join(f.calls, "\n"))
	}
	if !strings.Contains(out, "already runs") || !strings.Contains(out, second) {
		t.Errorf("says it now also serves the second project: %s", out)
	}
	entries, _ := serve.DirFor(cfg).Projects()
	if len(entries) != 2 {
		t.Fatalf("both registered: %+v", entries)
	}
	// one shared token and one shared agent credential, not one per project
	if _, err := os.Stat(filepath.Join(first, ".flai-cache", "dashboard.token")); !os.IsNotExist(err) {
		t.Error("no per-project token file")
	}
	if _, err := os.Stat(filepath.Join(serveDir, "dashboard.token")); err != nil {
		t.Errorf("the shared token: %v", err)
	}
	if _, err := os.Stat(filepath.Join(serveDir, "dashboard.agent-key")); err != nil {
		t.Errorf("the shared agent credential: %v", err)
	}
	for _, e := range entries {
		if e.KeyFile != filepath.Join(serveDir, "dashboard.agent-key") {
			t.Errorf("%s does not share the credential file: %s", e.Key, e.KeyFile)
		}
	}

	// stopping the first leaves the container running for the second
	out, _, code = runWith(t, first, f, "dashboard", "stop")
	if code != 0 {
		t.Fatal(out)
	}
	if !strings.Contains(out, "unregistered") || !strings.Contains(out, "second") {
		t.Errorf("stop of the first: %s", out)
	}
	if entries, _ := serve.DirFor(cfg).Projects(); len(entries) != 1 || entries[0].Key != "second" {
		t.Errorf("only the second remains: %+v", entries)
	}
	if strings.Contains(strings.Join(f.calls, "\n"), "docker stop") {
		t.Error("the container was stopped although the second project still needs it")
	}

	// flai dashboard status names every project the container serves
	statusOut, _, _ := runWith(t, second, f, "dashboard", "status")
	if !strings.Contains(statusOut, "serves: second") {
		t.Errorf("status names the project served: %s", statusOut)
	}

	// stopping the last registered project stops the container
	out, _, code = runWith(t, second, f, "dashboard", "stop")
	if code != 0 || !strings.Contains(out, "stopped flaiover") {
		t.Errorf("stop of the last: %d %s", code, out)
	}
	if !strings.Contains(strings.Join(f.calls, "\n"), "docker stop flaiover") {
		t.Error("the container should have been stopped")
	}
	if entries, _ := serve.DirFor(cfg).Projects(); len(entries) != 0 {
		t.Errorf("nothing left registered: %+v", entries)
	}
}
