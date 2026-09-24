package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSelfUpgradeCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "tok")
	content := []byte("#!/bin/sh\necho upgraded\n")
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	name := "flai"
	if runtime.GOOS == "windows" {
		t.Skip("tar archives are not used on windows")
	}
	_ = tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg})
	_, _ = tw.Write(content)
	_ = tw.Close()
	_ = gz.Close()
	archive := buf.Bytes()
	sum := sha256.Sum256(archive)
	archiveName := fmt.Sprintf("flai_2.0.0_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	sums := hex.EncodeToString(sum[:]) + "  " + archiveName + "\n"
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			http.Error(w, "nope", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/repos/o/r/releases":
			fmt.Fprintf(w, `[{"tag_name":"flai/v2.0.0","assets":[{"name":"checksums.txt","url":"%s/sums"},{"name":"%s","url":"%s/archive"}]}]`, srv.URL, archiveName, srv.URL)
		case "/archive":
			_, _ = w.Write(archive)
		case "/sums":
			_, _ = w.Write([]byte(sums))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	dir := t.TempDir()

	// S-0107: FLAI_RELEASES_API is where releases come from when --api is not given
	t.Setenv("FLAI_RELEASES_API", srv.URL)
	if out, errOut, code := runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r"); code != 0 || !strings.Contains(out, "latest is 2.0.0") {
		t.Fatalf("from FLAI_RELEASES_API: %d %s %s", code, out, errOut)
	}
	t.Setenv("FLAI_RELEASES_API", "")
	out, errOut, code := runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r", "--api", srv.URL)
	if code != 0 || !strings.Contains(out, "latest is 2.0.0 (upgrade available)") {
		t.Fatalf("check: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, dir, "self-upgrade", "--repo", "o/r", "--api", srv.URL, "--dir", dir)
	if code != 0 || !strings.Contains(out, "installed flai 2.0.0 to "+filepath.Join(dir, "flai")) {
		t.Fatalf("install: %d %s %s", code, out, errOut)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "flai"))
	if !bytes.Equal(got, content) {
		t.Errorf("installed content: %q", got)
	}
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	_, errOut, code = runWith(t, dir, &fakeRunner{missing: true, images: map[string]bool{}, running: map[string]bool{}}, "self-upgrade", "--check", "--repo", "o/r", "--api", srv.URL)
	if code == 0 || !strings.Contains(errOut, "private repository") {
		t.Errorf("without a token the error should hint at authentication: %d %s", code, errOut)
	}
}
