package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/serve"
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

// releasesStandIn answers as GitHub's API does for o/r: the releases given,
// as JSON, and the refs of the flaiover/v tags given. asked lists the paths
// requested so far.
func releasesStandIn(t *testing.T, releases string, flaioverTags ...string) (url string, asked func() []string) {
	t.Helper()
	var mu sync.Mutex
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		switch r.URL.Path {
		case "/repos/o/r/releases":
			fmt.Fprint(w, releases)
		case "/repos/o/r/git/matching-refs/tags/flaiover/v":
			refs := make([]string, len(flaioverTags))
			for i, tag := range flaioverTags {
				refs[i] = `{"ref":"refs/tags/` + tag + `"}`
			}
			fmt.Fprint(w, "["+strings.Join(refs, ",")+"]")
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(paths)
	}
}

// listedReleases are four published flai releases, out of order, beside a
// draft, a prerelease, and a dashboard release, which are not listed.
const listedReleases = `[
 {"tag_name":"flai/v1.2.0","published_at":"2026-10-01T09:00:00Z","assets":[{"name":"checksums.txt","url":"/sums"}]},
 {"tag_name":"flai/v1.10.0","published_at":"2026-10-05T09:00:00Z"},
 {"tag_name":"flai/v2.0.0","draft":true},
 {"tag_name":"flai/v1.11.0-rc.1","prerelease":true},
 {"tag_name":"flaiover/v0.4.0","published_at":"2026-10-02T09:00:00Z"},
 {"tag_name":"flai/v1.1.0","published_at":"2026-09-20T09:00:00Z"},
 {"tag_name":"flai/v1.0.0","published_at":"2026-09-01T09:00:00Z"}
]`

// withMinimum writes a project named name whose flai.minimum is minimum.
func withMinimum(t *testing.T, name, minimum string) string {
	t.Helper()
	root := tempProject(t)
	m := "version: 1\nname: " + name + "\nkey: k\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nflai:\n  minimum: " + minimum + "\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// S-0298: --list prints the published releases newest first and installs
// nothing, marking the installed one, the newest, and those below the
// minimum of the project here and of a project flai serve serves, even one
// whose minimum is above the running flai.
func TestSelfUpgradeList(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	t.Setenv("GITHUB_TOKEN", "tok")
	oldVersion, oldExe := buildinfo.Version, executablePath
	t.Cleanup(func() { buildinfo.Version, executablePath = oldVersion, oldExe })
	buildinfo.Version = "1.2.0"
	installed := filepath.Join(t.TempDir(), "flai")
	executablePath = func() (string, error) { return installed, nil }
	here := withMinimum(t, "Here", "1.2.0")
	other := withMinimum(t, "Other", "1.5.0")
	if err := serve.DirFor(cfg).Register(serve.Entry{Key: "o", Name: "Other", Root: other, URL: "http://127.0.0.1:1", KeyFile: "k"}); err != nil {
		t.Fatal(err)
	}
	api, asked := releasesStandIn(t, listedReleases)

	out, errOut, code := runIn(t, here, "self-upgrade", "--list", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("list: %d %s", code, errOut)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var versions []string
	marks := map[string]string{}
	for _, l := range lines[1:] {
		f := strings.Fields(l)
		versions = append(versions, f[0])
		marks[f[0]] = strings.Join(f[1:], " ")
	}
	if !slices.Equal(versions, []string{"1.10.0", "1.2.0", "1.1.0", "1.0.0"}) {
		t.Fatalf("newest first, published only: %v\n%s", versions, out)
	}
	for v, want := range map[string][]string{
		"1.10.0": {"2026-10-05", "latest"},
		"1.2.0":  {"2026-10-01", "installed", "below minimum 1.5.0 of Other"},
		"1.1.0":  {"below minimum 1.2.0 of Here", "below minimum 1.5.0 of Other"},
	} {
		for _, w := range want {
			if !strings.Contains(marks[v], w) {
				t.Errorf("%s: want %q in %q", v, w, marks[v])
			}
		}
	}
	if strings.Contains(marks["1.10.0"], "below") || strings.Contains(marks["1.10.0"], "installed") {
		t.Errorf("1.10.0 is above every minimum and not installed: %q", marks["1.10.0"])
	}

	out, errOut, code = runIn(t, here, "self-upgrade", "--list", "--repo", "o/r", "--api", api, "--json")
	if code != 0 {
		t.Fatalf("list --json: %d %s", code, errOut)
	}
	var listed []struct {
		Version, Tag string
		Published    time.Time
		Installed    bool
		Latest       bool
		BelowMinimum []projectMinimum `json:"below_minimum"`
	}
	if err := json.Unmarshal([]byte(out), &listed); err != nil || len(listed) != 4 {
		t.Fatalf("json: %v %s", err, out)
	}
	if l := listed[0]; l.Version != "1.10.0" || l.Tag != "flai/v1.10.0" || !l.Latest || l.Installed || len(l.BelowMinimum) != 0 || !l.Published.Equal(time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)) {
		t.Errorf("newest: %+v", l)
	}
	if l := listed[1]; !l.Installed || l.Latest || !slices.Equal(l.BelowMinimum, []projectMinimum{{Project: "Other", Minimum: "1.5.0"}}) {
		t.Errorf("installed: %+v", l)
	}
	if l := listed[2]; !slices.Equal(l.BelowMinimum, []projectMinimum{{Project: "Here", Minimum: "1.2.0"}, {Project: "Other", Minimum: "1.5.0"}}) {
		t.Errorf("below both: %+v", l)
	}
	if strings.Count(out, `"below_minimum"`) != 3 {
		t.Errorf("below_minimum is left out when there is none: %s", out)
	}
	if _, err := os.Stat(installed); !os.IsNotExist(err) {
		t.Errorf("--list installs nothing: %v", err)
	}
	for _, p := range asked() {
		if p != "/repos/o/r/releases" {
			t.Errorf("--list asked for %s", p)
		}
	}

	// a dev build is no published release, and says so
	buildinfo.Version = "dev"
	out, _, _ = runIn(t, here, "self-upgrade", "--list", "--repo", "o/r", "--api", api)
	if strings.Contains(out, "installed,") || !strings.Contains(out, "installed at "+installed+": dev, which is not a published release") {
		t.Errorf("dev build: %s", out)
	}

	for _, args := range [][]string{{"--version", "1.0.0"}, {"--check"}} {
		_, errOut, code := runIn(t, here, append([]string{"self-upgrade", "--list", "--repo", "o/r", "--api", api}, args...)...)
		if code == 0 || !strings.Contains(errOut, "list") {
			t.Errorf("--list with %v is refused: %d %s", args, code, errOut)
		}
	}
}

// S-0298: a version that is not published is refused, naming the published
// ones, before anything is downloaded or replaced.
func TestSelfUpgradeRefusesAnUnpublishedVersion(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "tok")
	api, asked := releasesStandIn(t, listedReleases)
	dir := t.TempDir()
	_, errOut, code := runIn(t, dir, "self-upgrade", "--version", "1.3.0", "--repo", "o/r", "--api", api, "--dir", dir)
	if code == 0 || !strings.Contains(errOut, "published are 1.10.0, 1.2.0, 1.1.0, 1.0.0") {
		t.Fatalf("refused: %d %s", code, errOut)
	}
	if got := asked(); !slices.Equal(got, []string{"/repos/o/r/releases"}) {
		t.Errorf("nothing is downloaded: %v", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "flai")); !os.IsNotExist(err) {
		t.Errorf("nothing is installed: %v", err)
	}
}

// S-0298, ADR-0117 §7: a release below the project's flai.minimum is warned
// about and installed.
func TestSelfUpgradeWarnsBelowTheMinimum(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("tar archives are not used on windows")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "tok")
	here := withMinimum(t, "Here", "1.2.0")
	api := fakeRelease(t, "1.1.0")
	dir := t.TempDir()
	out, errOut, code := runIn(t, here, "self-upgrade", "--version", "1.1.0", "--repo", "o/r", "--api", api, "--dir", dir)
	if code != 0 || !strings.Contains(out, "installed flai 1.1.0") {
		t.Fatalf("installs: %d %s %s", code, out, errOut)
	}
	for _, field := range []string{"minimum", "project", "version"} {
		if !strings.Contains(errOut, `"`+field+`":`) && !strings.Contains(errOut, field+"=") {
			t.Errorf("warning carries %s: %s", field, errOut)
		}
	}
	if !strings.Contains(errOut, "1.2.0") || !strings.Contains(errOut, "Here") || !strings.Contains(errOut, "WARN") {
		t.Errorf("warning: %s", errOut)
	}
}
