package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// S-0237: flai host status says under its serve child which projects' flai
// serve was refused as unsigned, and by which side, from flai serve's state;
// a connected project adds nothing, and without a running flai serve no
// refusal is said.
func TestDescribeHostSaysWhichSideRefusedTheOtherAsUnsigned(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	st := host.Status{Children: []host.Child{
		{Name: host.Serve, State: "running", PID: 41, Since: "2026-10-12T17:00:00Z", Version: "1.44.0"},
		{Name: host.MCP, Root: "/p/dock", State: "running", PID: 42, Since: "2026-10-12T17:00:01Z", Version: "1.44.0"},
	}}
	describe := func() string {
		var out bytes.Buffer
		(&app{configPath: cfg, out: &out}).describeHost(st)
		return out.String()
	}
	plain := describe()
	if strings.Contains(plain, "refused") {
		t.Errorf("a refusal said with no flai serve running:\n%s", plain)
	}
	writeServeState(t, cfg, serve.Status{Version: "1.44.0", Connections: refusedConnections})
	want := "  serve: running since 2026-10-12T17:00:00Z (pid 41, flai 1.44.0)\n" +
		"    /p/harbour: the dashboard refused this flai: unsigned (missing stamp)\n" +
		"    /p/quay: flai refused the dashboard: unsigned (wrong component flai, dashboard 1.2.3)\n" +
		"  mcp /p/dock: running since 2026-10-12T17:00:01Z (pid 42, flai 1.44.0)\n"
	if got := describe(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// S-0238: flai host status says under its serve child, beside each refusal,
// whether each project's dashboard image is signed, unsigned and not dialled,
// or not measured, from flai serve's state.
func TestDescribeHostSaysEachProjectsImageVerdict(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	st := host.Status{Children: []host.Child{
		{Name: host.Serve, State: "running", PID: 41, Since: "2026-10-12T17:00:00Z", Version: "1.44.0"},
	}}
	writeServeState(t, cfg, serve.Status{Version: "1.44.0", Connections: imageConnections, Images: imageVerdicts})
	var out bytes.Buffer
	(&app{configPath: cfg, out: &out}).describeHost(st)
	want := "  serve: running since 2026-10-12T17:00:00Z (pid 41, flai 1.44.0)\n" +
		"    /p/dock: image: signed, flaiover 1.2.4, sha256:b2b2\n" +
		"    /p/harbour: image: not measured: docker is not on PATH where flai serve runs; the release stamp alone is checked\n" +
		"    /p/quay: flai refused the dashboard: " + unsignedImageReason + "\n" +
		"    /p/quay: image: the dashboard runs an unsigned image sha256:c0c0 (no RepoDigests); flai does not dial it\n"
	if got := out.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// S-0238: flai host status --json carries flai serve's verdicts by root while
// the host runs, and none when flai serve has no state.
func TestHostStatusJSONCarriesEachProjectsImageVerdict(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	now := time.Now().UTC().Format(time.RFC3339)
	hst, err := json.Marshal(host.Status{PID: os.Getpid(), Version: "1.44.0", Started: now, Updated: now})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(string(host.DirFor(cfg)), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(string(host.DirFor(cfg)), "state.json"), hst, 0o600); err != nil {
		t.Fatal(err)
	}
	status := func() (hostStatus, string) {
		js, _, code := runIn(t, home, "host", "status", "--json")
		var st hostStatus
		if err := json.Unmarshal([]byte(js), &st); err != nil || code != 0 || !st.Running {
			t.Fatalf("json (exit %d): %v %s", code, err, js)
		}
		return st, js
	}
	if st, js := status(); st.Images != nil || strings.Contains(js, `"images"`) {
		t.Errorf("images with no flai serve state:\n%s", js)
	}
	writeServeState(t, cfg, serve.Status{Version: "1.44.0", Connections: imageConnections, Images: imageVerdicts})
	st, js := status()
	if !reflect.DeepEqual(st.Images, imageVerdicts) {
		t.Errorf("images = %+v, want %+v", st.Images, imageVerdicts)
	}
	for _, field := range []string{`"kind": "listed"`, `"kind": "unlisted"`, `"kind": "not-measured"`, `"release": "1.2.4"`, `"digest": "sha256:b2b2"`} {
		if !strings.Contains(js, field) {
			t.Errorf("--json lacks %s:\n%s", field, js)
		}
	}
}
