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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/channel"
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

// writeServeState writes flai serve's state beside cfg as a serve that runs
// now, in this process, would.
func writeServeState(t *testing.T, cfg string, st serve.Status) {
	t.Helper()
	now := time.Now().UTC().Format(time.RFC3339)
	st.PID, st.Started, st.Updated = os.Getpid(), now, now
	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	dir := serve.DirFor(cfg)
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(string(dir), "state.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// refusedConnections are a connection the dashboard refused as unsigned, one
// this flai refused, and one connected, as channel.Client records them.
var refusedConnections = map[string]channel.State{
	"/p/harbour": {URL: "http://127.0.0.1:4242", Attempts: 2, RefusedBy: channel.RefusedByDashboard, RefusedReason: "missing stamp",
		LastError: "the dashboard refused this flai: unsigned peer: missing stamp; trying again in 1m0s"},
	"/p/quay": {URL: "http://127.0.0.1:4243", Attempts: 1, RefusedBy: channel.RefusedByFlai, RefusedReason: "wrong component flai", PeerVersion: "1.2.3",
		LastError: "flai refused the dashboard: unsigned peer: wrong component flai; trying again in 1m0s"},
	"/p/dock": {URL: "http://127.0.0.1:4244", Connected: true, Since: "2026-10-12T17:00:00Z", PeerVersion: "1.2.4", PeerSigned: true},
}

// S-0237: flai serve status says which side refused the other as unsigned,
// with the reason and the dashboard's version where known, in place of the
// bare last error, for a registered project and one served from the folder;
// a connected project is shown as before, and --json carries the fields.
func TestServeStatusSaysWhichSideRefusedTheOtherAsUnsigned(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	for key, root := range map[string]string{"harbour": tempProject(t), "dock": tempProject(t)} {
		if err := os.WriteFile(filepath.Join(root, "k"), []byte("credential"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := serve.DirFor(cfg).Register(serve.Entry{Key: key, Name: key, Root: root, URL: "http://127.0.0.1:4242", KeyFile: filepath.Join(root, "k")}); err != nil {
			t.Fatal(err)
		}
	}
	projects, err := serve.DirFor(cfg).Projects()
	if err != nil {
		t.Fatal(err)
	}
	conns := map[string]channel.State{"/p/quay": refusedConnections["/p/quay"]}
	for _, p := range projects {
		conns[p.Root] = refusedConnections["/p/"+p.Key]
	}
	writeServeState(t, cfg, serve.Status{Version: "test", Connections: conns, Folder: "/p",
		FolderProjects: []serve.Entry{{Key: "quay", Root: "/p/quay", URL: "http://127.0.0.1:4243"}}})
	out, _, _ := runIn(t, home, "serve", "status")
	for _, want := range []string{
		"harbour  http://127.0.0.1:4242  not connected: the dashboard refused this flai: unsigned (missing stamp)\n",
		"quay  http://127.0.0.1:4243  not connected: flai refused the dashboard: unsigned (wrong component flai, dashboard 1.2.3)\n",
		"dock  http://127.0.0.1:4242  connected since 2026-10-12T17:00:00Z\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "unsigned peer") || strings.Contains(out, "trying again") {
		t.Errorf("the bare last error printed beside the refusal:\n%s", out)
	}
	js, _, _ := runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(js), &st); err != nil || st.Status == nil {
		t.Fatalf("json: %v %s", err, js)
	}
	if got := st.Status.Connections["/p/quay"]; got != refusedConnections["/p/quay"] {
		t.Errorf("flai's refusal: %+v", got)
	}
	for _, p := range projects {
		if got := st.Status.Connections[p.Root]; got != refusedConnections["/p/"+p.Key] {
			t.Errorf("%s: %+v", p.Key, got)
		}
	}
	for _, field := range []string{`"refused_by": "dashboard"`, `"refused_by": "flai"`, `"refused_reason": "missing stamp"`, `"peer_version": "1.2.3"`, `"peer_signed": true`} {
		if !strings.Contains(js, field) {
			t.Errorf("--json lacks %s:\n%s", field, js)
		}
	}
}

// A close that gave no reason and a dashboard that named no version are said
// without the parentheses rather than with an empty or repeated reason.
func TestRefusalLeavesOutWhatIsNotKnown(t *testing.T) {
	for _, c := range []struct {
		state channel.State
		want  string
	}{
		{channel.State{}, ""},
		{channel.State{LastError: "dial tcp: connection refused"}, ""},
		{channel.State{RefusedBy: channel.RefusedByDashboard, RefusedReason: "unsigned"}, "the dashboard refused this flai: unsigned"},
		{channel.State{RefusedBy: channel.RefusedByDashboard, RefusedReason: "invalid signature"}, "the dashboard refused this flai: unsigned (invalid signature)"},
		{channel.State{RefusedBy: channel.RefusedByFlai, RefusedReason: "missing stamp"}, "flai refused the dashboard: unsigned (missing stamp)"},
		{channel.State{RefusedBy: channel.RefusedByFlai, RefusedReason: "invalid signature", PeerVersion: "1.2.3"}, "flai refused the dashboard: unsigned (invalid signature, dashboard 1.2.3)"},
	} {
		if got := refusal(c.state); got != c.want {
			t.Errorf("refusal(%+v) = %q, want %q", c.state, got, c.want)
		}
	}
}

// unsignedImageReason is the reason flai serve's image gate refuses a dial
// of an image with no RepoDigests with, as channel.State records it.
const unsignedImageReason = "the dashboard runs an unsigned image sha256:c0c0 (no RepoDigests): " +
	"its image has no RepoDigests, as a local build has none, so no signed digest list names it"

// imageVerdicts are a listed image, an unlisted one with no RepoDigests, and
// one not measured, as flai serve's status records them (S-0238).
var imageVerdicts = map[string]serve.ImageVerdict{
	"/p/dock": {Kind: serve.ImageListed, Container: "d0d0", Image: "sha256:a1a1", Version: "1.2.4",
		RepoDigests: []string{"ghcr.io/bytepunx/flaiover@sha256:b2b2"}, Digest: "sha256:b2b2", Release: "1.2.4"},
	"/p/quay": {Kind: serve.ImageUnlisted, Container: "d1d1", Image: "sha256:c0c0", Version: "1.2.3",
		Reason: "its image has no RepoDigests, as a local build has none, so no signed digest list names it"},
	"/p/harbour": {Kind: serve.ImageNotMeasured, Reason: "docker is not on PATH where flai serve runs"},
}

// imageConnections are the connections of imageVerdicts' projects: the
// unlisted image refused by the gate, the others connected.
var imageConnections = map[string]channel.State{
	"/p/dock":    {URL: "http://127.0.0.1:4244", Connected: true, Since: "2026-10-12T17:00:00Z", PeerVersion: "1.2.4", PeerSigned: true},
	"/p/quay":    {URL: "http://127.0.0.1:4243", Attempts: 3, RefusedBy: channel.RefusedByFlai, RefusedReason: unsignedImageReason, LastError: "flai refused the dashboard: " + unsignedImageReason + "; trying again in 1m0s"},
	"/p/harbour": {URL: "http://127.0.0.1:4242", Connected: true, Since: "2026-10-12T17:00:01Z", PeerVersion: "1.2.2", PeerSigned: true},
}

// S-0238: each verdict on the dashboard's image is one line, in the words
// flai dashboard status uses; a kind this flai does not know is said by none.
func TestImageVerdictLineSaysEachVerdict(t *testing.T) {
	for _, c := range []struct {
		v    serve.ImageVerdict
		want string
	}{
		{imageVerdicts["/p/dock"], "image: signed, flaiover 1.2.4, sha256:b2b2"},
		{imageVerdicts["/p/quay"], "image: the dashboard runs an unsigned image sha256:c0c0 (no RepoDigests); flai does not dial it"},
		{serve.ImageVerdict{Kind: serve.ImageUnlisted, Image: "sha256:e4e4", RepoDigests: []string{"ghcr.io/x/flaiover@sha256:f5f5", "ghcr.io/x/flaiover@sha256:a6a6"}},
			"image: the dashboard runs an unsigned image sha256:e4e4 (ghcr.io/x/flaiover@sha256:f5f5, ghcr.io/x/flaiover@sha256:a6a6); flai does not dial it"},
		{imageVerdicts["/p/harbour"], "image: not measured: docker is not on PATH where flai serve runs; the release stamp alone is checked"},
		{serve.ImageVerdict{Kind: "weighed"}, ""},
	} {
		if got := imageVerdictLine(c.v); got != c.want {
			t.Errorf("imageVerdictLine(%+v) = %q, want %q", c.v, got, c.want)
		}
	}
}

// S-0238: a dial the image gate refused is said in the gate's words, not
// wrapped as a refused stamp is.
func TestRefusalSaysAnImageRefusalInTheGatesWords(t *testing.T) {
	if got, want := refusal(imageConnections["/p/quay"]), "flai refused the dashboard: "+unsignedImageReason; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// S-0238: flai serve status says under each project, registered or served
// from the folder, whether its dashboard's image is signed, unsigned and not
// dialled, or not measured, and nothing for a project with no verdict yet;
// --json carries each verdict under status.images by root.
func TestServeStatusSaysEachProjectsImageVerdict(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	roots := map[string]string{"harbour": tempProject(t), "dock": tempProject(t), "pier": tempProject(t)}
	for key, root := range roots {
		if err := os.WriteFile(filepath.Join(root, "k"), []byte("credential"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := serve.DirFor(cfg).Register(serve.Entry{Key: key, Name: key, Root: root, URL: "http://127.0.0.1:4242", KeyFile: filepath.Join(root, "k")}); err != nil {
			t.Fatal(err)
		}
	}
	conns := map[string]channel.State{"/p/quay": imageConnections["/p/quay"], roots["pier"]: {URL: "http://127.0.0.1:4242"}}
	images := map[string]serve.ImageVerdict{"/p/quay": imageVerdicts["/p/quay"]}
	for _, key := range []string{"harbour", "dock"} {
		conns[roots[key]], images[roots[key]] = imageConnections["/p/"+key], imageVerdicts["/p/"+key]
	}
	writeServeState(t, cfg, serve.Status{Version: "test", Connections: conns, Images: images, Folder: "/p",
		FolderProjects: []serve.Entry{{Key: "quay", Root: "/p/quay", URL: "http://127.0.0.1:4243"}}})
	out, _, _ := runIn(t, home, "serve", "status")
	for _, want := range []string{
		"  harbour  http://127.0.0.1:4242  connected since 2026-10-12T17:00:01Z\n    " + roots["harbour"] +
			"\n    image: not measured: docker is not on PATH where flai serve runs; the release stamp alone is checked\n",
		"  dock  http://127.0.0.1:4242  connected since 2026-10-12T17:00:00Z\n    " + roots["dock"] +
			"\n    image: signed, flaiover 1.2.4, sha256:b2b2\n",
		"  quay  http://127.0.0.1:4243  not connected: flai refused the dashboard: " + unsignedImageReason + "\n    /p/quay" +
			"\n    image: the dashboard runs an unsigned image sha256:c0c0 (no RepoDigests); flai does not dial it\n",
		"  pier  http://127.0.0.1:4242  not connected\n    " + roots["pier"] + "\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if n := strings.Count(out, "image: "); n != 3 {
		t.Errorf("%d image lines, want 3 (none for pier):\n%s", n, out)
	}
	js, _, _ := runIn(t, home, "serve", "status", "--json")
	var st serveStatus
	if err := json.Unmarshal([]byte(js), &st); err != nil || st.Status == nil {
		t.Fatalf("json: %v %s", err, js)
	}
	if !reflect.DeepEqual(st.Status.Images, images) {
		t.Errorf("status.images = %+v, want %+v", st.Status.Images, images)
	}
	for _, field := range []string{`"kind": "listed"`, `"kind": "unlisted"`, `"kind": "not-measured"`, `"release": "1.2.4"`, `"digest": "sha256:b2b2"`,
		`"repo_digests": [`, `"reason": "docker is not on PATH where flai serve runs"`} {
		if !strings.Contains(js, field) {
			t.Errorf("--json lacks %s:\n%s", field, js)
		}
	}
}
