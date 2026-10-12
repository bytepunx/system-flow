package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai serve keeps its registry beside the config file in use, so a test
// (and a repository with its own config) never touches the operator's home.
func TestServeStatusReadsTheRegistryBesideTheConfig(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	out, errOut, code := runIn(t, home, "serve", "status")
	if code != 0 || !strings.Contains(out, "not running") || !strings.Contains(out, "no projects registered") {
		t.Fatalf("empty: %d %s %s", code, out, errOut)
	}
	dir := serve.DirFor(cfg)
	if err := dir.Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: "/p/harbour", URL: "http://127.0.0.1:4242", KeyFile: "/p/harbour/k"}); err != nil {
		t.Fatal(err)
	}
	out, _, _ = runIn(t, home, "serve", "status")
	// a root that is not there says so rather than waiting to connect (S-0121)
	if !strings.Contains(out, "harbour  http://127.0.0.1:4242  not served: the folder is gone") {
		t.Errorf("status: %s", out)
	}
	out, _, _ = runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(out), &st); err != nil || st.Running || len(st.Projects) != 1 || st.Dir != filepath.Join(home, "serve") {
		t.Errorf("json: %v %+v", err, st)
	}
	if out, _, code := runIn(t, home, "serve", "stop"); code != 0 || !strings.Contains(out, "not running") {
		t.Errorf("stop with nothing running: %d %s", code, out)
	}
}

// S-0096: flai serve status says where each served project's MCP server
// listens, the one flai serve keeps running, so an agent can be pointed at it.
func TestServeStatusNamesEachProjectsMCPServer(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	root := tempProject(t)
	if err := serve.DirFor(cfg).Register(serve.Entry{Key: "harbour", Name: "Harbour", Root: root, URL: "http://127.0.0.1:4242", KeyFile: filepath.Join(root, "k")}); err != nil {
		t.Fatal(err)
	}
	if out, _, _ := runIn(t, home, "serve", "status"); strings.Contains(out, "mcp:") {
		t.Errorf("an MCP server named while none runs: %s", out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(repo.CacheDir(), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := writeMCPState(repo, mcpState{PID: os.Getpid(), Addr: "127.0.0.1:4244", URL: "http://127.0.0.1:4244/mcp", Started: now, Updated: now}); err != nil {
		t.Fatal(err)
	}
	out, _, _ := runIn(t, home, "serve", "status")
	if !strings.Contains(out, "mcp: http://127.0.0.1:4244/mcp (pid ") {
		t.Errorf("status: %s", out)
	}
	js, _, _ := runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(js), &st); err != nil || st.MCP[root].URL != "http://127.0.0.1:4244/mcp" {
		t.Errorf("json: %v %+v", err, st.MCP)
	}
}

// S-0183: flai serve has its MCP servers kept by the host in its environment
// only when that host runs from the config flai serve was given: one run with
// another --config, by an agent of the operator's, keeps none and says why.
func TestServeKeepsMCPServersOnlyWithTheHostItsConfigNames(t *testing.T) {
	home := t.TempDir()
	hostCfg, otherCfg := filepath.Join(home, "host", "config.json"), filepath.Join(home, "scratch", "config.json")
	dir := host.DirFor(hostCfg)
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dir.TokenFile(), []byte("host-token"), 0o600); err != nil {
		t.Fatal(err)
	}
	serveWith := func(cfg string) (serve.MCP, string) {
		var errOut bytes.Buffer
		a := &app{configPath: cfg, errOut: &errOut}
		return a.serveMCP(), errOut.String()
	}
	t.Setenv(host.URLEnv, "http://127.0.0.1:4241")
	t.Setenv(host.TokenEnv, "host-token")
	if mcp, said := serveWith(hostCfg); mcp == nil {
		t.Errorf("the serve the host started keeps no MCP server: %s", said)
	}
	mcp, said := serveWith(otherCfg)
	if mcp != nil || !strings.Contains(said, "not the one this config names") || !strings.Contains(said, otherCfg) {
		t.Errorf("a serve with another config took the host: %v, said %q", mcp, said)
	}
	t.Setenv(host.TokenEnv, "")
	if mcp, said := serveWith(hostCfg); mcp != nil || !strings.Contains(said, "not under flai host") {
		t.Errorf("outside a host: %v, said %q", mcp, said)
	}
}

// S-0235: flai serve logs one warn at start when a versioned build carries no
// valid release stamp, and none for a development build or a signed one.
func TestServeWarnsOnceWithoutAValidReleaseStamp(t *testing.T) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	statement := "flai 1.42.0 0ef6d26"
	sum := sha256.Sum256([]byte(statement))
	der, err := ecdsa.SignASN1(rand.Reader, k, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	sig := base64.StdEncoding.EncodeToString(der)
	keys := []*ecdsa.PublicKey{&k.PublicKey}
	for _, c := range []struct {
		name          string
		stamp         buildinfo.Stamp
		state, reason string
	}{
		{"versioned without a stamp", buildinfo.VerifyStamp("", "", "1.42.0", "0ef6d26", keys), "missing", ""},
		{"versioned, signed by another key", buildinfo.VerifyStamp(statement, sig, "1.42.0", "0ef6d26", nil), "invalid", "the signature does not verify against the release key"},
		{"development build", buildinfo.VerifyStamp("", "", "dev", "", keys), "", ""},
		{"snapshot build", buildinfo.VerifyStamp("", "", "0.0.0-0ef6d26", "0ef6d26", keys), "", ""},
		{"signed", buildinfo.VerifyStamp(statement, sig, "1.42.0", "0ef6d26", keys), "", ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			var logged bytes.Buffer
			warnReleaseStamp(slog.New(slog.NewJSONHandler(&logged, nil)), c.stamp)
			lines := strings.Split(strings.TrimSpace(logged.String()), "\n")
			if c.state == "" {
				if logged.Len() != 0 {
					t.Errorf("warned: %s", logged.String())
				}
				return
			}
			if len(lines) != 1 {
				t.Fatalf("want one event, got %q", logged.String())
			}
			var ev map[string]any
			if err := json.Unmarshal([]byte(lines[0]), &ev); err != nil {
				t.Fatal(err)
			}
			if ev["level"] != "WARN" || ev["component"] != "serve" || ev["state"] != c.state {
				t.Errorf("event: %v", ev)
			}
			if r, _ := ev["reason"].(string); r != c.reason {
				t.Errorf("reason %q, want %q", r, c.reason)
			}
		})
	}
}
