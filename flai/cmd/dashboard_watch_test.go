package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
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
	older := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}, containerRef: map[string]string{"flaiover": "ghcr.io/bytepunx/flaiover:0.1.0"}, containerImage: map[string]string{"flaiover": "sha256:v010"}}
	if _, errOut, code := runWith(t, root, older, "dashboard"); code != 0 {
		t.Fatalf("already running: %s", errOut)
	}
	if r, _ = readDashboardRecord(dir); r.Ref != "ghcr.io/bytepunx/flaiover:0.1.0" || r.Image != "sha256:v010" {
		t.Errorf("record of a running container: %+v", r)
	}
}

// S-0184: the host's look is off without a record or with
// dashboard.no_restart, and a restart starts the recorded image again, by its
// image ID (I-0116).
func TestDashboardWatchLooksAndRestartsAsRecorded(t *testing.T) {
	root := dashboardProject(t)
	answers := true
	f := &fakeRunner{images: map[string]bool{"sha256:v020": true}, running: map[string]bool{}, imageIDs: map[string]string{"ghcr.io/bytepunx/flaiover:0.2.0": "sha256:v020"}}
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
	if !strings.Contains(calls, "docker rm -f flaiover") || !strings.Contains(calls, "--publish 0.0.0.0:5555:3000") || !strings.HasSuffix(calls, "--label "+dashboardRefLabel+"=ghcr.io/bytepunx/flaiover:0.2.0 sha256:v020") || strings.Contains(calls, "docker pull") {
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

// lastRun is the last docker run among calls, "" when there is none.
func lastRun(calls []string) string {
	for i := len(calls) - 1; i >= 0; i-- {
		if strings.HasPrefix(calls[i], "docker run ") {
			return calls[i]
		}
	}
	return ""
}

// I-0116: the watch starts the image ID the dashboard ran, labelled with its
// tag, never what the tag names now, so a newer image flai dashboard check
// pulled under a floating tag is not started without an upgrade.
func TestDashboardWatchRestartsTheImageIDNotWhatItsTagNamesNow(t *testing.T) {
	root := dashboardProject(t)
	const latest = "ghcr.io/bytepunx/flaiover:latest"
	f := &fakeRunner{images: map[string]bool{"sha256:old": true}, running: map[string]bool{}, imageIDs: map[string]string{latest: "sha256:old"}}
	a := &app{cwd: root, runner: f}
	if _, errOut, code := runWithApp(t, a, "dashboard"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	if r, _ := readDashboardRecord(string(a.serveDir())); r.Ref != latest || r.Image != "sha256:old" {
		t.Errorf("record after start: %+v", r)
	}

	f.imageIDs[latest], f.images["sha256:new"] = "sha256:new", true // a check pulled a newer latest
	delete(f.running, "flaiover")
	before := len(f.calls)
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("restart: %v", err)
	}
	run := lastRun(f.calls[before:])
	if !strings.HasSuffix(run, " sha256:old") || !strings.Contains(run, "--label "+dashboardRefLabel+"="+latest+" ") || strings.Contains(strings.Join(f.calls[before:], "\n"), "docker pull") {
		t.Errorf("the watch restarts the image ID it recorded, labelled with its tag: %q", run)
	}
}

// I-0116: a record an older flai wrote has no image ID, and the watch starts
// its reference as before.
func TestDashboardWatchRestartsAnOlderRecordFromItsRef(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	a := &app{cwd: root, runner: f}
	dir := string(a.serveDir())
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	older := `{"name": "flaiover", "ref": "ghcr.io/bytepunx/flaiover:0.1.0", "publish": "0.0.0.0:5555:3000"}`
	if err := os.WriteFile(filepath.Join(dir, dashboardRecordFile), []byte(older), 0o600); err != nil {
		t.Fatal(err)
	}
	if r, ok := readDashboardRecord(dir); !ok || r.Ref != "ghcr.io/bytepunx/flaiover:0.1.0" || r.Image != "" {
		t.Fatalf("an older record is read: %+v %v", r, ok)
	}
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("restart: %v", err)
	}
	if run := lastRun(f.calls); !strings.HasSuffix(run, " ghcr.io/bytepunx/flaiover:0.1.0") || !strings.Contains(run, "--publish 0.0.0.0:5555:3000") {
		t.Errorf("an older record restarts from its ref: %q", run)
	}
}

// I-0116: an image ID that is no longer a local image, as after docker image
// prune, cannot be started: the watch says so at warn and starts the
// reference, and records the ID that now runs.
func TestDashboardWatchStartsTheRefWhenTheRecordedImageIsGone(t *testing.T) {
	root := dashboardProject(t)
	const ref = "ghcr.io/bytepunx/flaiover:0.2.0"
	f := &fakeRunner{images: map[string]bool{ref: true}, running: map[string]bool{}, imageIDs: map[string]string{ref: "sha256:v020"}}
	var logged bytes.Buffer
	a := &app{cwd: root, runner: f, log: slog.New(slog.NewTextHandler(&logged, nil))}
	dir := string(a.serveDir())
	if err := writeDashboardRecord(dir, dashboardRecord{Name: "flaiover", Ref: ref, Image: "sha256:pruned", Publish: "0.0.0.0:5555:3000"}); err != nil {
		t.Fatal(err)
	}
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("restart: %v", err)
	}
	if run := lastRun(f.calls); !strings.HasSuffix(run, " "+ref) || !strings.Contains(run, "--label "+dashboardRefLabel+"="+ref+" ") {
		t.Errorf("a gone image restarts from its ref: %q", run)
	}
	if log := logged.String(); !strings.Contains(log, "level=WARN") || !strings.Contains(log, "recorded dashboard image is gone") || !strings.Contains(log, "image=sha256:pruned") {
		t.Errorf("a gone image is logged at warn: %s", log)
	}
	if r, _ := readDashboardRecord(dir); r.Ref != ref || r.Image != "sha256:v020" {
		t.Errorf("the record names the image that runs now: %+v", r)
	}
}

// I-0116: flai dashboard restart starts the image ID and records the tag it
// shows beside it, so the watch's restart after it is labelled with the tag,
// not the ID.
func TestDashboardRestartRecordsTheTagAndTheImageID(t *testing.T) {
	root := dashboardProject(t)
	const latest = "ghcr.io/bytepunx/flaiover:latest"
	f := &fakeRunner{images: map[string]bool{"sha256:old": true}, running: map[string]bool{}, imageIDs: map[string]string{latest: "sha256:old"}}
	a := &app{cwd: root, runner: f}
	if _, errOut, code := runWithApp(t, a, "dashboard"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	f.imageIDs[latest], f.images["sha256:new"] = "sha256:new", true
	if _, errOut, code := runWithApp(t, a, "dashboard", "restart"); code != 0 {
		t.Fatalf("restart: %s", errOut)
	}
	if r, _ := readDashboardRecord(string(a.serveDir())); r.Ref != latest || r.Image != "sha256:old" {
		t.Errorf("record after flai dashboard restart: %+v", r)
	}

	delete(f.running, "flaiover")
	before := len(f.calls)
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("watch restart: %v", err)
	}
	if run := lastRun(f.calls[before:]); !strings.HasSuffix(run, " sha256:old") || !strings.Contains(run, "--label "+dashboardRefLabel+"="+latest+" ") {
		t.Errorf("the watch after a restart: %q", run)
	}
}
