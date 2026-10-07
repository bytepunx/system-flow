//go:build !windows

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/host"
)

// hostReleases is the list the stand-in's self-upgrade --list --json prints.
const hostReleases = `[{"version":"1.2.0","tag":"flai/v1.2.0","published":"2026-10-01T09:00:00Z","installed":true,"latest":true},` +
	`{"version":"1.1.0","tag":"flai/v1.1.0","published":"2026-09-20T09:00:00Z","installed":false,"latest":false,"below_minimum":[{"project":"Here","minimum":"1.2.0"}]},` +
	`{"version":"1.0.0","tag":"flai/v1.0.0","installed":false,"latest":false}]`

// selfUpgradeStandIn writes a flai that answers self-upgrade as the real one
// does with --json for the published releases 1.2.0, 1.1.0, and 1.0.0, and
// appends each command line it is given to the file it answers.
func selfUpgradeStandIn(t *testing.T, installedPath string) (exe, calls string) {
	t.Helper()
	dir := t.TempDir()
	exe, calls = filepath.Join(dir, "flai"), filepath.Join(dir, "calls")
	refused := `{"level":"FATAL","msg":"command failed","err":"resolve flai 1.3.0: o/r has no published release tagged flai/v1.3.0; published are 1.2.0, 1.1.0, 1.0.0: choose one of them, or give no version for the newest"}`
	script := "#!/bin/sh\necho \"$*\" >> '" + calls + "'\ncase \"$*\" in\n" +
		"*--list*) echo '" + hostReleases + "' ;;\n" +
		"*'--version 1.3.0'*) echo '" + refused + "' >&2; exit 1 ;;\n" +
		"*'--version 1.1.0'*) echo '{\"previous\":\"1.2.0\",\"installed\":\"1.1.0\",\"tag\":\"flai/v1.1.0\",\"path\":\"" + installedPath + "\"}' ;;\n" +
		"*) echo '{\"current\":\"1.2.0\",\"latest\":\"1.2.0\",\"up_to_date\":true,\"path\":\"" + installedPath + "\"}' ;;\nesac\n"
	if err := os.WriteFile(exe, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return exe, calls
}

// inProcessHost is a host run in this test for the config cfg, whose
// launcher's flai is exe; done gets what Run returned.
type inProcessHost struct {
	l    *hostLauncher
	done chan error
}

func runHostFor(t *testing.T, cfg, exe string) *inProcessHost {
	t.Helper()
	l := &hostLauncher{a: &app{}, exe: exe, config: cfg, taken: map[string]string{}}
	ctx, cancel := context.WithCancel(context.Background())
	h := &inProcessHost{l: l, done: make(chan error, 1)}
	ended := make(chan struct{})
	go func() {
		err := host.Run(ctx, host.Options{
			Dir: host.DirFor(cfg), Addr: "127.0.0.1:0", Version: "1.2.0", Config: cfg, Every: 50 * time.Millisecond,
			Serve:    func() (host.Spec, error) { return host.Spec{External: os.Getpid()}, nil },
			Versions: l.versions, Upgrade: l.upgradeTo,
		})
		close(ended)
		h.done <- err
	}()
	t.Cleanup(func() { cancel(); <-ended })
	waitUntil(t, "the host's state", func() bool { _, _, err := host.DirFor(cfg).Client(time.Now()); return err == nil })
	return h
}

func (h *inProcessHost) restarts(t *testing.T) {
	t.Helper()
	select {
	case err := <-h.done:
		if !errors.Is(err, host.ErrRestart) {
			t.Fatalf("the host ended with %v, not to restart", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the host did not end to restart on what was installed")
	}
}

func (h *inProcessHost) stillRuns(t *testing.T) {
	t.Helper()
	select {
	case err := <-h.done:
		t.Fatalf("the host ended: %v", err)
	case <-time.After(400 * time.Millisecond): // longer than an installed upgrade waits to restart
	}
}

func lastCall(t *testing.T, calls string) string {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(readFile(calls)), "\n")
	return lines[len(lines)-1]
}

// S-0298: flai host versions prints the host's list of published releases,
// as text and as it is with --json.
func TestHostVersions(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	exe, calls := selfUpgradeStandIn(t, "/home/me/.flai/bin/flai")
	dir := t.TempDir()

	if _, errOut, code := runIn(t, dir, "host", "versions"); code == 0 || !strings.Contains(errOut, "flai host is not running") {
		t.Errorf("no host: %d %s", code, errOut)
	}

	runHostFor(t, cfg, exe)
	out, errOut, code := runIn(t, dir, "host", "versions")
	if code != 0 {
		t.Fatalf("versions: %d %s", code, errOut)
	}
	if got := lastCall(t, calls); got != "self-upgrade --list --json --config "+cfg {
		t.Errorf("the host ran %q", got)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 || lines[0] != "published flai releases, newest first (the host runs flai 1.2.0):" {
		t.Fatalf("versions:\n%s", out)
	}
	for i, want := range [][]string{
		{"1.2.0", "2026-10-01", "installed, latest"},
		{"1.1.0", "2026-09-20", "below minimum 1.2.0 of Here"},
		{"1.0.0", "-"},
	} {
		f := strings.Fields(lines[i+1])
		if len(f) < 2 || f[0] != want[0] || f[1] != want[1] || (len(want) == 3 && !strings.HasSuffix(lines[i+1], want[2])) {
			t.Errorf("line %d: %q, want %q", i+1, lines[i+1], want)
		}
	}
	if strings.Contains(lines[3], "installed") || strings.Contains(lines[3], "latest") {
		t.Errorf("1.0.0 is neither installed nor the newest: %q", lines[3])
	}

	out, errOut, code = runIn(t, dir, "host", "versions", "--json")
	if code != 0 {
		t.Fatalf("versions --json: %d %s", code, errOut)
	}
	var got, want any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	_ = json.Unmarshal([]byte(hostReleases), &want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("--json prints the list as the host answers it:\n%s", out)
	}
}

// S-0298: flai host upgrade --version installs that published release
// through the host, which restarts on what was installed as for the newest;
// a malformed or unpublished version is refused before anything is replaced.
func TestHostUpgradeToAVersion(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	installed := "/home/me/.flai/bin/flai"
	exe, calls := selfUpgradeStandIn(t, installed)
	dir := t.TempDir()

	if _, errOut, code := runIn(t, dir, "host", "upgrade", "--version", "1.1.0"); code == 0 || !strings.Contains(errOut, "flai host is not running") {
		t.Errorf("no host: %d %s", code, errOut)
	}

	h := runHostFor(t, cfg, exe)
	for _, v := range []string{"v1.1.0", "1.1", "1.1.0-rc.1"} {
		if _, errOut, code := runIn(t, dir, "host", "upgrade", "--version", v); code != 1 || !strings.Contains(errOut, "is not a release version") {
			t.Errorf("%s: %d %s", v, code, errOut)
		}
	}
	if c := readFile(calls); c != "" {
		t.Errorf("a malformed version ran flai: %s", c)
	}
	_, errOut, code := runIn(t, dir, "host", "upgrade", "--version", "1.3.0")
	if code != 1 || !strings.Contains(errOut, "published are 1.2.0, 1.1.0, 1.0.0") {
		t.Errorf("unpublished: %d %s", code, errOut)
	}
	h.stillRuns(t)
	if to := h.l.installedTo(); to != "" {
		t.Errorf("a refused version is restarted on: %q", to)
	}

	out, errOut, code := runIn(t, dir, "host", "upgrade", "--version", "1.1.0")
	if code != 0 || out != "installed flai 1.1.0 (was 1.2.0); the host restarts on it with its processes\n" {
		t.Fatalf("upgrade --version: %d %q %s", code, out, errOut)
	}
	if got := lastCall(t, calls); got != "self-upgrade --version 1.1.0 --json --config "+cfg {
		t.Errorf("the host ran %q", got)
	}
	h.restarts(t)
	if to := h.l.installedTo(); to != installed {
		t.Errorf("the host restarts on %q, want %q", to, installed)
	}

	// --json is the host's answer as it is: what was installed, over what,
	// and that the host restarts on it
	h = runHostFor(t, cfg, exe)
	out, errOut, code = runIn(t, dir, "host", "upgrade", "--version", "1.1.0", "--json")
	var said struct {
		Restarting bool              `json:"restarting"`
		Upgrade    map[string]string `json:"upgrade"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &said) != nil || !said.Restarting ||
		said.Upgrade["installed"] != "1.1.0" || said.Upgrade["previous"] != "1.2.0" || said.Upgrade["path"] != installed {
		t.Fatalf("upgrade --version --json: %d %s %s", code, out, errOut)
	}
	h.restarts(t)

	// with no version it is the newest, as before
	h = runHostFor(t, cfg, exe)
	out, errOut, code = runIn(t, dir, "host", "upgrade")
	if code != 0 || out != "flai 1.2.0 is already the latest release\n" || lastCall(t, calls) != "self-upgrade --json --config "+cfg {
		t.Errorf("upgrade: %d %q %s, ran %q", code, out, errOut, lastCall(t, calls))
	}
	h.stillRuns(t)
}
