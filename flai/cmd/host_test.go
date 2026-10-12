package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

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
