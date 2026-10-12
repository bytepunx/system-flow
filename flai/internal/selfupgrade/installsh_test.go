package selfupgrade_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// installShPath is the repository's install.sh, from this package's folder.
const installShPath = "../../../install.sh"

// install.sh verifies a release with the release key it carries as PEM
// (S-0233); flai verifies with the one it embeds. They are the same key, or a
// release one of them installs the other refuses.
func TestInstallShCarriesTheReleaseKey(t *testing.T) {
	script := readInstallSh(t)
	const begin, end = "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----"
	if n := strings.Count(script, begin); n != 1 {
		t.Fatalf("install.sh holds %d PEM public keys, want 1", n)
	}
	from := strings.Index(script, begin)
	to := strings.Index(script[from:], end)
	if to < 0 {
		t.Fatalf("install.sh's PEM public key has no %s line", end)
	}
	got := script[from : from+to+len(end)]
	want := strings.TrimSpace(buildinfo.ReleaseKeyPEM())
	if got != want {
		t.Errorf("install.sh's release key differs from flai/internal/buildinfo/cosign.pub; copy the file's PEM into install.sh's RELEASE_KEY\ninstall.sh:\n%s\ncosign.pub:\n%s", got, want)
	}
}

// install.sh installs a release whose checksums.txt the key it carries signed,
// and says so; it refuses, saying why, one whose signature is missing, does
// not verify, or was made by a key it does not know, and installs nothing.
func TestInstallShVerifiesTheReleaseSignature(t *testing.T) {
	key := installShKey(t)
	other := installShKey(t)
	script := installShWithKey(t, key)
	cases := []struct {
		name   string
		script string
		sig    func(sums []byte) []byte // nil: the release has no signature
		want   string
		ok     bool
	}{
		{"good", script, func(sums []byte) []byte { return installShSign(t, key, sums) }, "Signature verified", true},
		{"bad", script, func([]byte) []byte { return installShSign(t, key, []byte("other checksums")) }, "does not verify with the release key", false},
		{"missing", script, nil, "has no checksums.txt.sig, so it cannot be verified and is not installed", false},
		{"unknown key", script, func(sums []byte) []byte { return installShSign(t, other, sums) }, "does not verify with the release key", false},
		// The repository's own install.sh knows only the release key.
		{"test key, real install.sh", installShPath, func(sums []byte) []byte { return installShSign(t, key, sums) }, "does not verify with the release key", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			api := installShRelease(t, c.sig)
			dir := filepath.Join(t.TempDir(), "bin")
			out, code := runInstallShWith(t, c.script, os.Getenv("PATH"), api, dir)
			if !strings.Contains(out, c.want) {
				t.Errorf("install.sh does not say %q:\n%s", c.want, out)
			}
			_, err := os.Stat(filepath.Join(dir, "flai"))
			if c.ok {
				if code != 0 || err != nil {
					t.Errorf("install.sh: exit %d, installed: %v; want it installed\n%s", code, err, out)
				}
				if !strings.Contains(out, "Checksum verified") {
					t.Errorf("install.sh does not say the checksum was verified:\n%s", out)
				}
				return
			}
			if code == 0 || !os.IsNotExist(err) {
				t.Errorf("install.sh: exit %d, installed: %v; want a refusal that installs nothing\n%s", code, err, out)
			}
			if strings.Contains(out, "Checksum verified") {
				t.Errorf("install.sh checked the archive's checksum before the signature:\n%s", out)
			}
		})
	}
}

// Without openssl install.sh cannot verify a signature, so it refuses before
// it fetches anything, and says what to install.
func TestInstallShRefusesWithoutOpenssl(t *testing.T) {
	key := installShKey(t)
	api := installShRelease(t, func(sums []byte) []byte { return installShSign(t, key, sums) })
	tools := t.TempDir()
	for _, tool := range []string{"curl", "tar", "uname", "tr", "sed", "grep", "awk", "head", "mktemp", "rm", "mkdir", "install", "sha256sum", "shasum", "cat", "sleep"} {
		if path, err := exec.LookPath(tool); err == nil {
			if err := os.Symlink(path, filepath.Join(tools, tool)); err != nil {
				t.Fatal(err)
			}
		}
	}
	dir := filepath.Join(t.TempDir(), "bin")
	out, code := runInstallShWith(t, installShWithKey(t, key), tools, api, dir)
	if code == 0 || !strings.Contains(out, "openssl is required to verify the release's signature") {
		t.Errorf("install.sh: exit %d, want a refusal naming openssl\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(dir, "flai")); !os.IsNotExist(err) {
		t.Errorf("something was installed: %v", err)
	}
}

func readInstallSh(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile(installShPath)
	if err != nil {
		t.Fatalf("read install.sh: %v", err)
	}
	return string(b)
}

// installShKey makes an ECDSA P-256 key pair, as the release key is.
func installShKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// installShSign signs data as cosign sign-blob does: base64 of the ASN.1 DER
// ECDSA signature over its SHA-256.
func installShSign(t *testing.T, key *ecdsa.PrivateKey, data []byte) []byte {
	t.Helper()
	digest := sha256.Sum256(data)
	sig, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(sig))
}

// installShWithKey writes a copy of install.sh that carries key's public half
// in place of the release key, and returns its path.
func installShWithKey(t *testing.T, key *ecdsa.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pub := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})))
	script := readInstallSh(t)
	release := strings.TrimSpace(buildinfo.ReleaseKeyPEM())
	if strings.Count(script, release) != 1 {
		t.Fatal("install.sh does not carry the release key once; TestInstallShCarriesTheReleaseKey says how it differs")
	}
	path := filepath.Join(t.TempDir(), "install.sh")
	if err := os.WriteFile(path, []byte(strings.Replace(script, release, pub, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// installShRelease answers install.sh as GitHub's releases API does for o/r
// with release flai/v1.2.3: an archive for every platform install.sh knows,
// each a flai that prints its version, their checksums.txt, and, unless sig is
// nil, checksums.txt.sig, the signature sig gives for checksums.txt.
func installShRelease(t *testing.T, sig func(sums []byte) []byte) string {
	t.Helper()
	content := []byte("#!/bin/sh\necho 'flai 1.2.3'\n")
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "flai", Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg})
	_, _ = tw.Write(content)
	_ = tw.Close()
	_ = gz.Close()
	archive := buf.Bytes()
	sum := sha256.Sum256(archive)

	// Asset 0 is checksums.txt, 1 the signature, and 2 on the archives.
	var names []string
	var sums strings.Builder
	for _, goos := range []string{"linux", "darwin"} {
		for _, arch := range []string{"amd64", "arm64"} {
			name := fmt.Sprintf("flai_1.2.3_%s_%s.tar.gz", goos, arch)
			names = append(names, name)
			sums.WriteString(hex.EncodeToString(sum[:]) + "  " + name + "\n")
		}
	}
	sumsFile := []byte(sums.String())
	var sigFile []byte
	if sig != nil {
		sigFile = sig(sumsFile)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		assets := base + "/repos/o/r/releases/assets/"
		switch {
		case r.URL.Path == "/repos/o/r/releases/tags/flai/v1.2.3":
			list := []string{`{"url":"` + assets + `0","name":"checksums.txt"}`}
			if sigFile != nil {
				list = append(list, `{"url":"`+assets+`1","name":"checksums.txt.sig"}`)
			}
			for i, name := range names {
				list = append(list, fmt.Sprintf(`{"url":"%s%d","name":"%s"}`, assets, i+2, name))
			}
			fmt.Fprintf(w, `{"url":"%s/repos/o/r/releases/1","tag_name":"flai/v1.2.3","assets":[%s]}`, base, strings.Join(list, ","))
		case r.URL.Path == "/repos/o/r/releases/assets/0":
			_, _ = w.Write(sumsFile)
		case r.URL.Path == "/repos/o/r/releases/assets/1" && sigFile != nil:
			_, _ = w.Write(sigFile)
		case strings.HasPrefix(r.URL.Path, "/repos/o/r/releases/assets/"):
			_, _ = w.Write(archive)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

// runInstallShWith runs the install.sh at script with PATH path against the
// API at api, pinned to flai 1.2.3, installing into dir, and returns its
// output and exit code.
func runInstallShWith(t *testing.T, script, path, api, dir string) (string, int) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("install.sh does not install on Windows")
	}
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("install.sh needs sh")
	}
	for _, tool := range []string{"curl", "tar", "openssl"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("install.sh needs %s", tool)
		}
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, sh, script)
	// a token of its own, so that install.sh does not ask gh for one
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + path, "FLAI_API=" + api, "FLAI_REPO=o/r", "FLAI_VERSION=1.2.3", "FLAI_INSTALL_DIR=" + dir, "GITHUB_TOKEN=stand-in"}
	out, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState == nil {
		t.Fatalf("run install.sh: %v", err)
	}
	return string(out), cmd.ProcessState.ExitCode()
}
