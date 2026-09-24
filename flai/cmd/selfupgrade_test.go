package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
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

// fakeRelease serves release v of a flai that is a shell script answering
// "version --json" with v, as GitHub's releases API does.
func fakeRelease(t *testing.T, v string) string {
	t.Helper()
	content := []byte("#!/bin/sh\necho '{\"version\":\"" + v + "\"}'\n")
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "flai", Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg})
	_, _ = tw.Write(content)
	_ = tw.Close()
	_ = gz.Close()
	archive := buf.Bytes()
	sum := sha256.Sum256(archive)
	archiveName := fmt.Sprintf("flai_%s_%s_%s.tar.gz", v, runtime.GOOS, runtime.GOARCH)
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/o/r/releases":
			fmt.Fprintf(w, `[{"tag_name":"flai/v%s","assets":[{"name":"checksums.txt","url":"%s/sums"},{"name":"%s","url":"%s/archive"}]}]`, v, srv.URL, archiveName, srv.URL)
		case "/archive":
			_, _ = w.Write(archive)
		case "/sums":
			_, _ = w.Write([]byte(hex.EncodeToString(sum[:]) + "  " + archiveName + "\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

// S-0111: a flai run from inside a project (a checkout's bin/flai) never
// installs a release there: it goes under the home folder, and whether it is
// up to date is asked of the flai there.
func TestSelfUpgradeNeverInstallsIntoAProject(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a shell script stands in for flai")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("FLAI_INSTALL_DIR", "")
	api := fakeRelease(t, "2.0.0")
	project := tempProject(t)
	inCheckout := filepath.Join(project, "bin", "flai")
	old := executablePath
	t.Cleanup(func() { executablePath = old })
	executablePath = func() (string, error) { return inCheckout, nil }

	want := filepath.Join(home, ".flai", "bin", "flai")
	out, errOut, code := runIn(t, project, "self-upgrade", "--repo", "o/r", "--api", api)
	if code != 0 || !strings.Contains(out, "installed flai 2.0.0 to "+want+" (was not installed there)") || !strings.Contains(out, "put "+filepath.Dir(want)+" on your PATH ahead of "+filepath.Dir(inCheckout)) {
		t.Fatalf("install: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(inCheckout); !os.IsNotExist(err) {
		t.Errorf("nothing is written into the checkout: %v", err)
	}
	// the flai under home is the one asked: it is the latest now
	if out, _, _ := runIn(t, project, "self-upgrade", "--repo", "o/r", "--api", api); !strings.Contains(out, "already the latest") {
		t.Errorf("up to date against the home install: %s", out)
	}
	js, _, _ := runIn(t, project, "self-upgrade", "--check", "--repo", "o/r", "--api", api, "--json")
	if !strings.Contains(js, `"path": "`+want+`"`) || !strings.Contains(js, `"up_to_date": true`) {
		t.Errorf("check: %s", js)
	}
	// FLAI_INSTALL_DIR is honoured, as install.sh honours it
	other := t.TempDir()
	t.Setenv("FLAI_INSTALL_DIR", other)
	if out, _, code := runIn(t, project, "self-upgrade", "--repo", "o/r", "--api", api, "--json"); code != 0 || !strings.Contains(out, `"path": "`+filepath.Join(other, "flai")+`"`) || !strings.Contains(out, `"running": "`+inCheckout+`"`) {
		t.Errorf("FLAI_INSTALL_DIR: %s", out)
	}

	// an installed flai outside any project is replaced where it is, as before
	outside := filepath.Join(t.TempDir(), "flai")
	executablePath = func() (string, error) { return outside, nil }
	if out, errOut, code := runIn(t, project, "self-upgrade", "--repo", "o/r", "--api", api); code != 0 || !strings.Contains(out, "installed flai 2.0.0 to "+outside) || strings.Contains(out, "PATH") {
		t.Errorf("in place: %d %s %s", code, out, errOut)
	}
}

// S-0111: flai host restarts on the binary the upgrade installed, which is
// not its own when it ran from a checkout.
func TestTheHostRestartsOnWhatTheUpgradeInstalled(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a shell script stands in for flai")
	}
	exe := filepath.Join(t.TempDir(), "flai")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\necho '{\"previous\":\"1.0.0\",\"installed\":\"2.0.0\",\"path\":\"/home/me/.flai/bin/flai\"}'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	l := &hostLauncher{a: &app{}, exe: exe, config: filepath.Join(t.TempDir(), "cfg.json"), taken: map[string]string{}}
	if l.installedTo() != "" {
		t.Fatal("nothing installed yet")
	}
	if _, installed, err := l.upgrade(context.Background()); err != nil || !installed {
		t.Fatalf("upgrade: %v %v", installed, err)
	}
	if got := l.installedTo(); got != "/home/me/.flai/bin/flai" {
		t.Errorf("restarts on %q", got)
	}
}
