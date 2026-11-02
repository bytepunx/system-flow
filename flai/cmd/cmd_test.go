package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// sessionEnv are the variables flai reads that the shell running the tests
// may have set: who runs flai, in what role and session, where its host and
// releases are, how it logs, and the forge's tokens, which an agent's session,
// flai serve, the host, or the operator sets.
var sessionEnv = []string{
	"FLAI_AGENT", "FLAI_ROLE", "FLAI_STORY", "FLAI_SESSION", "FLAI_ITEM", "FLAI_FOCUS", "FLAI_SINCE",
	guard.StartedByEnv, harness.ReasonsEnv, host.URLEnv, host.TokenEnv, serve.MCPURLEnv, serve.MCPTokenEnv, pprofEnv,
	"FLAI_RELEASES_API", "FLAI_INSTALL_DIR", "LOG_LEVEL", "LOG_FORMAT", "GITHUB_TOKEN", "GH_TOKEN",
}

// noConfig is FLAI_CONFIG while the tests run: a path no config can be
// created at, so that a command that reads the configuration fails unless
// its test named a file of its own with useConfig or FLAI_CONFIG.
var noConfig string

// TestMain clears sessionEnv, so that no test inherits the session of the
// agent that runs the tests and a test that needs none of them set need not
// clear them with t.Setenv, which would keep it from running in parallel
// (S-0548). For the same reason it points FLAI_CACHE_DIR, the cache_dir a new
// configuration is written with, at a folder of the run's own, which every
// test shares, and FLAI_CONFIG at noConfig.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	dir, err := os.MkdirTemp("", "flai-cmd-test-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = os.RemoveAll(dir) }()
	file := filepath.Join(dir, "not-a-folder")
	noConfig = filepath.Join(file, "config.json")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for _, v := range sessionEnv {
		if err := os.Unsetenv(v); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	for k, v := range map[string]string{config.EnvVar: noConfig, config.CacheEnvVar: filepath.Join(dir, "cache")} {
		if err := os.Setenv(k, v); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	return m.Run()
}

// together runs tests as parallel subtests of t, a test that does not run
// in parallel and has set what they all need of the process, an environment
// variable or a package variable: a test that sets either cannot run in
// parallel with the package's other tests, but tests that need the same can
// run with each other (S-0548). Each subtest is named for its function,
// testX running as X.
func together(t *testing.T, tests ...func(*testing.T)) {
	t.Helper()
	for _, test := range tests {
		name := runtime.FuncForPC(reflect.ValueOf(test).Pointer()).Name()
		t.Run(strings.TrimPrefix(name[strings.LastIndex(name, ".")+1:], "test"), func(t *testing.T) {
			t.Parallel()
			test(t)
		})
	}
}

// inEnv fails t unless the environment has key set to want, as the test
// that runs t with together sets it.
func inEnv(t *testing.T, key, want string) {
	t.Helper()
	if got := os.Getenv(key); got != want {
		t.Fatalf("%s is %q, not %q: run this test with together from the test that sets it", key, got, want)
	}
}

// perTest is what a test, and the subtests it runs, set for the commands
// they run, in place of an environment or package variable that would keep
// them from running in parallel (S-0548). It is keyed by the test's name, so
// that a subtest finds what its parent set.
type perTest[V any] struct{ m sync.Map }

// set makes v the value for t and its subtests until t ends.
func (p *perTest[V]) set(t *testing.T, v V) {
	t.Helper()
	name := t.Name()
	p.m.Store(name, v)
	t.Cleanup(func() { p.m.Delete(name) })
}

// get is the value t or the nearest test that runs it set.
func (p *perTest[V]) get(t *testing.T) (V, bool) {
	for name := t.Name(); ; {
		if v, ok := p.m.Load(name); ok {
			return v.(V), true
		}
		i := strings.LastIndex(name, "/")
		if i < 0 {
			var zero V
			return zero, false
		}
		name = name[:i]
	}
}

// testConfigs are the config files tests named with useConfig.
var testConfigs perTest[string]

// useConfig makes path the config file of the commands t and its subtests
// run, as FLAI_CONFIG would, without setting it.
func useConfig(t *testing.T, path string) {
	t.Helper()
	testConfigs.set(t, path)
}

// configOf is the config file of t's commands: the one t or its parent named
// with useConfig, else a file of t's own, which the first command that reads
// it creates. A test that sets FLAI_CONFIG gets what it set.
func configOf(t *testing.T) string {
	t.Helper()
	if env := os.Getenv(config.EnvVar); env != noConfig {
		return env
	}
	if path, ok := testConfigs.get(t); ok {
		return path
	}
	path := filepath.Join(t.TempDir(), "config.json")
	useConfig(t, path)
	return path
}

// testRoot is the command tree of a, with configOf(t) as its config file
// unless a names one.
func testRoot(t *testing.T, a *app) *cobra.Command {
	t.Helper()
	if a.configPath == "" {
		a.configPath = configOf(t)
	}
	return newRootCmdWith(a)
}

// runCLI runs the CLI as Execute does, with configOf(t) as its config file
// unless args name one; FLAI_CONFIG set to "" leaves flai's default.
func runCLI(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	if path := configOf(t); path != "" {
		args = append([]string{"--config", path}, args...)
	}
	var out, errOut bytes.Buffer
	code = run(args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestVersionPlainAndJSON(t *testing.T) {
	t.Parallel()
	out, _, code := runCLI(t, "version")
	if code != 0 || !strings.HasPrefix(out, "flai dev") || !strings.Contains(out, "commit:") || !strings.Contains(out, "built:") {
		t.Fatalf("plain version: code=%d out=%q", code, out)
	}
	out, _, code = runCLI(t, "version", "--json")
	if code != 0 {
		t.Fatalf("json version exit %d", code)
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	for _, k := range []string{"version", "commit", "date", "go"} {
		if s, _ := v[k].(string); s == "" {
			t.Errorf("missing %s in %v", k, v)
		}
	}
	// a test build, like scripts/flai.sh's, carries no stamp (S-0235)
	var stamp struct {
		Release map[string]any `json:"release"`
	}
	if err := json.Unmarshal([]byte(out), &stamp); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"statement": "", "signature": "", "state": "missing", "development": true}
	if !reflect.DeepEqual(stamp.Release, want) {
		t.Errorf("release = %v, want %v", stamp.Release, want)
	}
	plain, _, _ := runCLI(t, "version")
	if !strings.HasSuffix(plain, "\nrelease: development build, no stamp\n") {
		t.Errorf("plain version names no release line: %q", plain)
	}
}

// A stamped build shows its statement, its signature, and signed; an invalid
// stamp says why (S-0235).
func TestVersionShowsTheReleaseStamp(t *testing.T) {
	t.Parallel()
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
	info := buildinfo.Info{Version: "1.42.0", Commit: "0ef6d26", Date: "2026-10-12T00:00:00Z", Go: "go1.26"}

	var out bytes.Buffer
	a := &app{out: &out, jsonOut: true}
	if err := a.printVersion(versionReport{Info: info, Release: buildinfo.VerifyStamp(statement, sig, "1.42.0", "0ef6d26", keys)}); err != nil {
		t.Fatal(err)
	}
	var v struct {
		Version string         `json:"version"`
		Release map[string]any `json:"release"`
	}
	if err := json.Unmarshal(out.Bytes(), &v); err != nil {
		t.Fatalf("not json: %v\n%s", err, out.String())
	}
	want := map[string]any{"statement": statement, "signature": sig, "state": "signed", "development": false}
	if v.Version != "1.42.0" || !reflect.DeepEqual(v.Release, want) {
		t.Errorf("json version = %s", out.String())
	}

	out.Reset()
	a.jsonOut = false
	if err := a.printVersion(versionReport{Info: info, Release: buildinfo.VerifyStamp(statement, sig, "1.42.0", "0ef6d26", keys)}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(out.String(), "\nrelease: signed\n") {
		t.Errorf("plain signed: %q", out.String())
	}

	out.Reset()
	other, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.printVersion(versionReport{Info: info, Release: buildinfo.VerifyStamp(statement, sig, "1.42.0", "0ef6d26", []*ecdsa.PublicKey{&other.PublicKey})}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "\nrelease: invalid stamp: the signature does not verify") {
		t.Errorf("plain signed by another key: %q", out.String())
	}

	out.Reset()
	if err := a.printVersion(versionReport{Info: info, Release: buildinfo.VerifyStamp("", "", "1.42.0", "0ef6d26", keys)}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(out.String(), "\nrelease: missing stamp\n") {
		t.Errorf("plain versioned without a stamp: %q", out.String())
	}
}

func TestConfigFirstRunGetSet(t *testing.T) {
	t.Setenv("FLAI_CONFIG", "")
	path := filepath.Join(t.TempDir(), "cfg.json")

	out, errOut, code := runCLI(t, "config", "path", "--config", path)
	if code != 0 || strings.TrimSpace(out) != path {
		t.Fatalf("config path: %q %q %d", out, errOut, code)
	}
	if _, err := os.Stat(path); err == nil {
		t.Fatal("config path must not create the file")
	}

	out, errOut, code = runCLI(t, "config", "get", "--config", path)
	if code != 0 || !strings.Contains(errOut, "created") || !strings.Contains(out, "\"template\"") {
		t.Fatalf("first run get: %q %q %d", out, errOut, code)
	}

	out, _, code = runCLI(t, "config", "set", "dashboard.port", "9999", "--config", path)
	if code != 0 || strings.TrimSpace(out) != "dashboard.port = 9999" {
		t.Fatalf("set: %q %d", out, code)
	}
	out, _, code = runCLI(t, "config", "get", "dashboard.port", "--config", path)
	if code != 0 || strings.TrimSpace(out) != "9999" {
		t.Fatalf("get after set: %q %d", out, code)
	}
	out, _, code = runCLI(t, "--config", path, "config", "get", "template")
	if code != 0 || !strings.Contains(out, "\"repo\"") {
		t.Fatalf("get object: %q %d", out, code)
	}

	t.Setenv("FLAI_CONFIG", path)
	out, _, code = runCLI(t, "config", "get", "dashboard.port")
	if code != 0 || strings.TrimSpace(out) != "9999" {
		t.Fatalf("env config: %q %d", out, code)
	}

	_, errOut, code = runCLI(t, "config", "set", "dashboard.port", "abc", "--config", path)
	if code == 0 || !strings.Contains(errOut, "integer") {
		t.Fatalf("bad set should fail: %q %d", errOut, code)
	}
	_, errOut, code = runCLI(t, "config", "get", "nope", "--config", path)
	if code == 0 || !strings.Contains(errOut, "unknown key") {
		t.Fatalf("unknown key should fail: %q %d", errOut, code)
	}
}

func TestHelpListsCommands(t *testing.T) {
	t.Parallel()
	out, _, code := runCLI(t, "--help")
	if code != 0 {
		t.Fatalf("help exit %d", code)
	}
	for _, c := range []string{"version", "config"} {
		if !strings.Contains(out, c) {
			t.Errorf("help missing %s:\n%s", c, out)
		}
	}
}
