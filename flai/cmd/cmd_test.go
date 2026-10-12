package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

func runCLI(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	var out, errOut bytes.Buffer
	code = run(args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestVersionPlainAndJSON(t *testing.T) {
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
