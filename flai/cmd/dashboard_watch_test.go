package cmd

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/host"
)

// S-0184: status probes the published port and says running, not answering,
// or gone, with docker's HEALTHCHECK verdict beside it.
func TestDashboardStatusSaysRunningNotAnsweringOrGone(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}, probeAddr: "0.0.0.0:4242", health: map[string]string{"flaiover": "healthy"}}

	var probed []string
	answers := true
	status := func(args ...string) string {
		t.Helper()
		a := &app{cwd: root, runner: f, healthProbe: func(url string) bool { probed = append(probed, url); return answers }}
		out, errOut, code := runWithApp(t, a, append([]string{"dashboard", "status"}, args...)...)
		if code != 0 {
			t.Fatalf("status: %d %s", code, errOut)
		}
		return out
	}

	out := status()
	if !strings.Contains(out, "flaiover running at") || !strings.Contains(out, "docker: healthy") {
		t.Errorf("running: %s", out)
	}
	if len(probed) != 1 || probed[0] != "http://127.0.0.1:4242/_health" {
		t.Errorf("a port published on every interface is probed on loopback: %v", probed)
	}
	var js struct {
		Running bool   `json:"running"`
		State   string `json:"state"`
		Health  string `json:"health"`
	}
	if err := json.Unmarshal([]byte(status("--json")), &js); err != nil || !js.Running || js.State != "running" || js.Health != "healthy" {
		t.Errorf("running --json: %+v %v", js, err)
	}

	answers, f.health["flaiover"] = false, "unhealthy"
	out = status()
	if !strings.Contains(out, "flaiover not answering at") || !strings.Contains(out, "docker: unhealthy") || !strings.Contains(out, "flai dashboard restart") {
		t.Errorf("not answering: %s", out)
	}
	js = struct {
		Running bool   `json:"running"`
		State   string `json:"state"`
		Health  string `json:"health"`
	}{}
	if err := json.Unmarshal([]byte(status("--json")), &js); err != nil || !js.Running || js.State != "not-answering" {
		t.Errorf("not answering --json: %+v %v", js, err)
	}

	delete(f.running, "flaiover")
	probed = nil
	out = status()
	if !strings.Contains(out, "flaiover gone") || !strings.Contains(out, "flai dashboard") {
		t.Errorf("gone: %s", out)
	}
	if len(probed) != 0 {
		t.Errorf("no container, nothing to probe: %v", probed)
	}
	js.Health = ""
	if err := json.Unmarshal([]byte(status("--json")), &js); err != nil || js.Running || js.State != "gone" || js.Health != "" {
		t.Errorf("gone --json: %+v %v", js, err)
	}
}

func TestProbeHostProbesAWildcardOnThisHost(t *testing.T) {
	for in, want := range map[string]string{
		"0.0.0.0:4242":     "127.0.0.1:4242",
		"[::]:4242":        "127.0.0.1:4242",
		"127.0.0.1:4242":   "127.0.0.1:4242",
		"192.168.1.5:4242": "192.168.1.5:4242",
		"nonsense":         "nonsense",
	} {
		if got := probeHost(in); got != want {
			t.Errorf("probeHost(%q) = %q, want %q", in, got, want)
		}
	}
}

// S-0184: flai dashboard records the container it starts, from the image it
// ran, an upgrade records the new image, and a stop that leaves no project
// forgets it, so that flai host does not bring back what the operator stopped.
func TestDashboardRecordsWhatFlaiHostRestarts(t *testing.T) {
	root := dashboardProject(t)
	dir := string((&app{}).serveDir())
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard", "--tag", "0.2.0"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	r, ok := readDashboardRecord(dir)
	if !ok || r.Name != "flaiover" || r.Ref != "ghcr.io/bytepunx/flaiover:0.2.0" || r.Publish != "0.0.0.0:5555:3000" {
		t.Fatalf("record after start: %+v %v", r, ok)
	}

	f.imageIDs = map[string]string{"ghcr.io/bytepunx/flaiover:latest": "sha256:newer"}
	if _, errOut, code := runWithApp(t, &app{cwd: root, runner: f}, "dashboard", "upgrade"); code != 0 {
		t.Fatalf("upgrade: %s", errOut)
	}
	if r, _ = readDashboardRecord(dir); r.Ref != "ghcr.io/bytepunx/flaiover:latest" {
		t.Errorf("record after upgrade: %+v", r)
	}

	if _, errOut, code := runWith(t, root, f, "dashboard", "stop"); code != 0 {
		t.Fatalf("stop: %s", errOut)
	}
	if _, ok := readDashboardRecord(dir); ok {
		t.Error("a dashboard stopped with no project left is forgotten")
	}

	// one an older flai started, already running, is recorded as it runs
	older := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}, containerRef: map[string]string{"flaiover": "ghcr.io/bytepunx/flaiover:0.1.0"}}
	if _, errOut, code := runWith(t, root, older, "dashboard"); code != 0 {
		t.Fatalf("already running: %s", errOut)
	}
	if r, _ = readDashboardRecord(dir); r.Ref != "ghcr.io/bytepunx/flaiover:0.1.0" {
		t.Errorf("record of a running container: %+v", r)
	}
}

// S-0184: the host's look is off without a record or with
// dashboard.no_restart, and a restart starts the recorded image again.
func TestDashboardWatchLooksAndRestartsAsRecorded(t *testing.T) {
	root := dashboardProject(t)
	answers := true
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	a := &app{cwd: root, runner: f, healthProbe: func(string) bool { return answers }}
	if _, _, code := runWithApp(t, a, "dashboard", "--tag", "0.2.0"); code != 0 {
		t.Fatal("start")
	}
	w := a.dashboardWatch()
	ctx := context.Background()
	if got := w.Look(ctx); got != host.DashboardRunning {
		t.Errorf("running: %s", got)
	}
	answers = false
	if got := w.Look(ctx); got != host.DashboardNotAnswering {
		t.Errorf("not answering: %s", got)
	}
	delete(f.running, "flaiover")
	if got := w.Look(ctx); got != host.DashboardGone {
		t.Errorf("gone: %s", got)
	}

	before := len(f.calls)
	if err := w.Restart(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	calls := strings.Join(f.calls[before:], "\n")
	if !strings.Contains(calls, "docker rm -f flaiover") || !strings.Contains(calls, "--publish 0.0.0.0:5555:3000") || !strings.HasSuffix(calls, "ghcr.io/bytepunx/flaiover:0.2.0") || strings.Contains(calls, "docker pull") {
		t.Errorf("restart calls:\n%s", calls)
	}
	if !f.running["flaiover"] {
		t.Error("restarted container runs")
	}

	if _, errOut, code := runWithApp(t, &app{cwd: root, runner: f}, "config", "set", "dashboard.no_restart", "true"); code != 0 {
		t.Fatalf("config set: %s", errOut)
	}
	if got := w.Look(ctx); got != host.DashboardOff {
		t.Errorf("dashboard.no_restart: %s", got)
	}

	forgetDashboard(string(a.serveDir()))
	if err := w.Restart(ctx); err == nil {
		t.Error("nothing recorded, nothing to restart")
	}
}
