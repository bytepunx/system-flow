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
	t.Parallel()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}, probeAddr: "0.0.0.0:4242", health: map[string]string{"flaiover": "healthy"}}

	var probed []string
	answers := true
	status := func(args ...string) string {
		t.Helper()
		a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, healthProbe: func(url string) bool { probed = append(probed, url); return answers }}
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
	t.Parallel()
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
	t.Parallel()
	root := dashboardProject(t)
	dir := string((&app{configPath: configOf(t)}).serveDir())
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard", "--tag", "0.2.0"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	r, ok := readDashboardRecord(dir)
	if !ok || r.Name != "flaiover" || r.Ref != flaioverRef("0.2.0") || r.Publish != "0.0.0.0:5555:3000" {
		t.Fatalf("record after start: %+v %v", r, ok)
	}

	// dashboard.tag latest resolves to the newest release, run by digest
	if _, errOut, code := runWithApp(t, &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f}, "dashboard", "upgrade"); code != 0 {
		t.Fatalf("upgrade: %s", errOut)
	}
	if r, _ = readDashboardRecord(dir); r.Ref != flaioverRef("0.10.0") || r.Image != f.idFor(flaioverRef("0.10.0")) {
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
	t.Parallel()
	root := dashboardProject(t)
	answers := true
	f := &fakeRunner{images: map[string]bool{"sha256:v020": true}, running: map[string]bool{}, imageIDs: map[string]string{flaioverRef("0.2.0"): "sha256:v020"}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, healthProbe: func(string) bool { return answers }}
	if _, _, code := runWithApp(t, a, "dashboard", "--tag", "0.2.0"); code != 0 {
		t.Fatal("start")
	}
	w := a.dashboardWatch()
	ctx := context.Background()
	if got, _ := w.Look(ctx); got != host.DashboardRunning {
		t.Errorf("running: %s", got)
	}
	answers = false
	if got, _ := w.Look(ctx); got != host.DashboardNotAnswering {
		t.Errorf("not answering: %s", got)
	}
	delete(f.running, "flaiover")
	if got, _ := w.Look(ctx); got != host.DashboardGone {
		t.Errorf("gone: %s", got)
	}

	before := len(f.calls)
	if err := w.Restart(ctx); err != nil {
		t.Fatalf("restart: %v", err)
	}
	calls := strings.Join(f.calls[before:], "\n")
	if !strings.Contains(calls, "docker rm -f flaiover") || !strings.Contains(calls, "--publish 0.0.0.0:5555:3000") || !strings.HasSuffix(calls, "--label "+dashboardRefLabel+"="+flaioverRef("0.2.0")+" sha256:v020") || strings.Contains(calls, "docker pull") {
		t.Errorf("restart calls:\n%s", calls)
	}
	if !f.running["flaiover"] {
		t.Error("restarted container runs")
	}

	if _, errOut, code := runWithApp(t, &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f}, "config", "set", "dashboard.no_restart", "true"); code != 0 {
		t.Fatalf("config set: %s", errOut)
	}
	if got, _ := w.Look(ctx); got != host.DashboardOff {
		t.Errorf("dashboard.no_restart: %s", got)
	}

	forgetDashboard(string(a.serveDir()))
	if err := w.Restart(ctx); err == nil {
		t.Error("nothing recorded, nothing to restart")
	}
}

// I-0243: a container Docker no longer lists keeps its docker-proxy on the
// dashboard's port, so docker publishes nothing for every container flai
// starts after it, while the port answers, from the orphan. The look says
// port held and names the proxy, rather than not answering, which the host
// restarted every five minutes; with the port not answering it is not
// answering as before (ADR-0267).
func TestLookDashboardFindsThePortHeldByAnOrphanedContainersProxy(t *testing.T) {
	t.Parallel()
	root := dashboardProject(t)
	proc := t.TempDir()
	writeProc(t, proc, "188220", []string{"/usr/bin/docker-proxy", "-proto", "tcp", "-host-ip", "0.0.0.0", "-host-port", "5555", "-container-ip", "172.17.0.2", "-container-port", "3000", "-use-listen-fd"})
	writeProc(t, proc, "188237", []string{"node", "server.js"})
	held := true
	var probed []string
	asked := 0
	f := &fakeRunner{images: map[string]bool{"sha256:v020": true}, running: map[string]bool{}, imageIDs: map[string]string{flaioverRef("0.2.0"): "sha256:v020"}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, procRoot: proc,
		healthProbe: func(string) bool { asked++; return true },
		portProbe:   func(publish string) bool { probed = append(probed, publish); return held }}
	if _, errOut, code := runWithApp(t, a, "dashboard", "--tag", "0.2.0"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	f.noPort, asked = true, 0 // docker inspect shows an empty port map: the orphan's proxy has the port
	w := a.dashboardWatch()
	ctx := context.Background()
	state, detail := w.Look(ctx)
	if state != host.DashboardPortHeld || detail != "docker-proxy (pid 188220) for 172.17.0.2:3000" {
		t.Errorf("look: %s %q, want port-held by the orphan's docker-proxy", state, detail)
	}
	if len(probed) != 1 || probed[0] != "0.0.0.0:5555:3000" {
		t.Errorf("the recorded publish is probed: %v", probed)
	}

	before := len(f.calls)
	for range 40 { // ten minutes of looks: the host restarted twice in that time
		if s, _ := w.Look(ctx); s != host.DashboardPortHeld {
			t.Fatalf("look: %s", s)
		}
	}
	for _, c := range f.calls[before:] {
		if strings.HasPrefix(c, "docker rm") || strings.HasPrefix(c, "docker run") {
			t.Errorf("a look changes nothing: %s", c)
		}
	}

	a.procRoot = t.TempDir() // no docker-proxy names the port
	if state, detail = w.Look(ctx); state != host.DashboardPortHeld || detail != "a process other than its container" {
		t.Errorf("look with no holder found: %s %q", state, detail)
	}

	if asked != 0 {
		t.Errorf("nothing is published, so /_health was asked at no address, yet asked %d times", asked)
	}

	held = false // the holder is gone: nothing answers, and the host restarts it as before
	if state, detail = w.Look(ctx); state != host.DashboardNotAnswering || detail != "" {
		t.Errorf("look with the port free: %s %q", state, detail)
	}
}

// I-0243, ADR-0267: flai dashboard status says a dashboard whose port another
// process holds is port held, names the docker-proxy that holds it and what
// to do, and gives the holder under port_holder in --json; with no holder
// found, it names ss to find the process.
func TestDashboardStatusSaysWhatHoldsThePort(t *testing.T) {
	t.Parallel()
	root := dashboardProject(t)
	proc := t.TempDir()
	writeProc(t, proc, "188220", []string{"/usr/bin/docker-proxy", "-proto", "tcp", "-host-ip", "0.0.0.0", "-host-port", "5555", "-container-ip", "172.17.0.2", "-container-port", "3000"})
	f := &fakeRunner{images: map[string]bool{"sha256:v020": true}, running: map[string]bool{}, imageIDs: map[string]string{flaioverRef("0.2.0"): "sha256:v020"}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, procRoot: proc,
		portProbe: func(string) bool { return true }}
	if _, errOut, code := runWithApp(t, a, "dashboard", "--tag", "0.2.0"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	f.noPort = true

	out, errOut, code := runWithApp(t, a, "dashboard", "status")
	if code != 0 {
		t.Fatalf("status: %d %s", code, errOut)
	}
	for _, want := range []string{
		"flaiover port held at ",
		"port 5555 is held by docker-proxy (pid 188220) for 172.17.0.2:3000, not by its container",
		"so flai host does not restart it",
		"Stop that process (sudo kill 188220) or restart the Docker daemon, and flai host starts the dashboard again on its own",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("status lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "not answering") {
		t.Errorf("a held port is not said as not answering:\n%s", out)
	}

	js, errOut, code := runWithApp(t, a, "dashboard", "status", "--json")
	var st struct {
		State  string      `json:"state"`
		Holder *portHolder `json:"port_holder"`
	}
	if code != 0 || json.Unmarshal([]byte(js), &st) != nil {
		t.Fatalf("status --json: %d %s %s", code, js, errOut)
	}
	want := portHolder{PID: 188220, HostIP: "0.0.0.0", HostPort: 5555, ContainerIP: "172.17.0.2", ContainerPort: "3000"}
	if st.State != host.DashboardPortHeld || st.Holder == nil || *st.Holder != want {
		t.Errorf("--json: state %q, holder %+v, want port-held by %+v", st.State, st.Holder, want)
	}

	a.procRoot = t.TempDir() // no docker-proxy names the port
	out, _, _ = runWithApp(t, a, "dashboard", "status")
	if want := "port 5555 is held by a process other than its container, so flai host does not restart it. Find that process with sudo ss -ltnp 'sport = :5555'"; !strings.Contains(out, want) {
		t.Errorf("status with no holder lacks %q:\n%s", want, out)
	}
	js, _, _ = runWithApp(t, a, "dashboard", "status", "--json")
	if strings.Contains(js, "port_holder") || !strings.Contains(js, `"state": "port-held"`) {
		t.Errorf("--json with no holder: %s", js)
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
	t.Parallel()
	root := mirrorProject(t)
	const latest = mirrorLatest
	f := &fakeRunner{images: map[string]bool{"sha256:old": true}, running: map[string]bool{}, imageIDs: map[string]string{latest: "sha256:old"}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f}
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
	t.Parallel()
	root := dashboardProject(t)
	useSignedDashboardReleases(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	// the restart reads the configuration, which it creates, saying so
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, log: slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))}
	dir := string(a.serveDir())
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	older := `{"name": "flaiover", "ref": "registry.example/flaiover:0.1.0", "publish": "0.0.0.0:5555:3000"}`
	if err := os.WriteFile(filepath.Join(dir, dashboardRecordFile), []byte(older), 0o600); err != nil {
		t.Fatal(err)
	}
	if r, ok := readDashboardRecord(dir); !ok || r.Ref != "registry.example/flaiover:0.1.0" || r.Image != "" {
		t.Fatalf("an older record is read: %+v %v", r, ok)
	}
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("restart: %v", err)
	}
	if run := lastRun(f.calls); !strings.HasSuffix(run, " registry.example/flaiover:0.1.0") || !strings.Contains(run, "--publish 0.0.0.0:5555:3000") {
		t.Errorf("an older record of another repository restarts from its ref: %q", run)
	}
}

// S-0236, ADR-0070: a record an older flai wrote names the signed repository
// by tag; a restart, by flai host's watch or flai dashboard restart, runs the
// release that tag resolves to by digest and records it, and one whose tag
// is refused starts nothing.
func TestARestartPinsAnOlderRecordOfTheSignedRepositoryByDigest(t *testing.T) {
	t.Parallel()
	root := dashboardProject(t)
	useSignedDashboardReleases(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, log: slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))}
	dir := string(a.serveDir())
	record := func(ref string) {
		t.Helper()
		if err := writeDashboardRecord(dir, dashboardRecord{Name: "flaiover", Ref: ref, Image: "sha256:by-tag", Publish: "0.0.0.0:5555:3000"}); err != nil {
			t.Fatal(err)
		}
		f.images["sha256:by-tag"] = true
	}
	record(signedDashboardImage + ":0.3.0")
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("watch restart: %v", err)
	}
	want := flaioverRef("0.3.0")
	if run := lastRun(f.calls); !strings.HasSuffix(run, " "+want) || !strings.Contains(run, "--label "+dashboardRefLabel+"="+want+" ") {
		t.Errorf("the watch runs the tag's release by digest, not the image pulled by tag: %q", run)
	}
	if r, _ := readDashboardRecord(dir); r.Ref != want {
		t.Errorf("the record keeps the digest: %+v", r)
	}

	record(signedDashboardImage + ":0.4.0")
	delete(f.running, "flaiover")
	out, errOut, code := runWithApp(t, a, "dashboard", "restart", "--json")
	if want := flaioverRef("0.4.0"); code != 0 || !strings.Contains(out, `"ref": "`+want+`"`) || !strings.HasSuffix(lastRun(f.calls), " "+want) {
		t.Errorf("flai dashboard restart runs the tag's release by digest: %d %s %s", code, out, errOut)
	}

	record(signedDashboardImage + ":edge")
	delete(f.running, "flaiover")
	f.calls = nil
	if err := a.dashboardWatch().Restart(context.Background()); err == nil || !strings.Contains(err.Error(), "is not a flaiover release version") {
		t.Errorf("a refused tag is not restarted: %v", err)
	}
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker run ") || strings.HasPrefix(c, "docker rm ") {
			t.Errorf("nothing is removed or started for a refused tag: %s", c)
		}
	}
}

// I-0116: an image ID that is no longer a local image, as after docker image
// prune, cannot be started: the watch says so at warn and starts the
// reference, and records the ID that now runs.
func TestDashboardWatchStartsTheRefWhenTheRecordedImageIsGone(t *testing.T) {
	t.Parallel()
	root := dashboardProject(t)
	ref := flaioverRef("0.2.0")
	f := &fakeRunner{images: map[string]bool{ref: true}, running: map[string]bool{}, imageIDs: map[string]string{ref: "sha256:v020"}}
	var logged bytes.Buffer
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, log: slog.New(slog.NewTextHandler(&logged, nil))}
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
	t.Parallel()
	root := mirrorProject(t)
	const latest = mirrorLatest
	f := &fakeRunner{images: map[string]bool{"sha256:old": true}, running: map[string]bool{}, imageIDs: map[string]string{latest: "sha256:old"}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f}
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

// S-0489: flai host's watch starts the recorded container again from its
// image with the allowed hosts the configuration names now, as flai dashboard
// restart does, and records them; when the configuration cannot be read it
// keeps the hosts the record names.
func TestDashboardWatchRestartsWithTheConfiguredAllowedHosts(t *testing.T) {
	t.Parallel()
	root := dashboardProject(t)
	ref := flaioverRef("0.2.0")
	f := &fakeRunner{images: map[string]bool{"sha256:v020": true}, running: map[string]bool{}, imageIDs: map[string]string{ref: "sha256:v020"}}
	a := &app{configPath: configOf(t), releases: releasesOf(t), cwd: root, runner: f, log: slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))}
	dir := string(a.serveDir())
	if err := writeDashboardRecord(dir, dashboardRecord{Name: "flaiover", Ref: ref, Image: "sha256:v020", Publish: "0.0.0.0:5555:3000", AllowedHosts: "old.example"}); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runWithApp(t, a, "config", "set", "dashboard.allowed_hosts", "New.Example, tunnel.example:8443"); code != 0 {
		t.Fatalf("config set: %s", errOut)
	}
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("restart: %v", err)
	}
	const configured = "new.example,tunnel.example:8443"
	if run := lastRun(f.calls); !strings.Contains(run, "--env "+allowedHostsEnv+"="+configured+" ") || strings.Contains(run, "old.example") || !strings.HasSuffix(run, " sha256:v020") {
		t.Errorf("the watch starts the recorded image with the configuration's hosts: %q", run)
	}
	if r, _ := readDashboardRecord(dir); r.AllowedHosts != configured || r.Image != "sha256:v020" || r.Ref != ref {
		t.Errorf("the record names the hosts it runs with now: %+v", r)
	}

	if err := os.WriteFile(configOf(t), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	delete(f.running, "flaiover")
	if err := a.dashboardWatch().Restart(context.Background()); err != nil {
		t.Fatalf("restart with no configuration to read: %v", err)
	}
	if run := lastRun(f.calls); !strings.Contains(run, "--env "+allowedHostsEnv+"="+configured+" ") {
		t.Errorf("with no configuration to read, the watch keeps the recorded hosts: %q", run)
	}
}
