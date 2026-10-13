package cmd

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
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// selfUpgradeConfig points FLAI_CONFIG at a new config whose cache_dir is a
// new temporary folder, so that no test reads or writes ~/.flai/cache, and
// answers the config's path and that folder.
func selfUpgradeConfig(t *testing.T) (path, cacheDir string) {
	t.Helper()
	path, cacheDir = filepath.Join(t.TempDir(), "cfg.json"), t.TempDir()
	cfg := config.Default()
	cfg.CacheDir = cacheDir
	if err := config.Save(path, cfg); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FLAI_CONFIG", path)
	return path, cacheDir
}

func TestSelfUpgradeCommand(t *testing.T) {
	selfUpgradeConfig(t)
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
	sums := []byte(hex.EncodeToString(sum[:]) + "  " + archiveName + "\n")
	sig := signBlob(t, trustTestKey(t), sums)
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			http.Error(w, "nope", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/repos/o/r/releases":
			fmt.Fprintf(w, `[{"tag_name":"flai/v2.0.0","assets":[{"name":"checksums.txt","url":"%s/sums"},{"name":"checksums.txt.sig","url":"%s/sig"},{"name":"%s","url":"%s/archive"}]}]`, srv.URL, srv.URL, archiveName, srv.URL)
		case "/archive":
			_, _ = w.Write(archive)
		case "/sums":
			_, _ = w.Write(sums)
		case "/sig":
			_, _ = w.Write(sig)
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
// "version --json" with v, as GitHub's releases API does, signed with a key
// the test makes the release key.
func fakeRelease(t *testing.T, v string) string {
	t.Helper()
	url, _ := serveFlaiReleases(t, flaiRelease(t, v, trustTestKey(t)))
	return url
}

// trustTestKey makes a new key the release key self-upgrade verifies with for
// the rest of t, and answers it to sign releases with.
func trustTestKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	old := releaseKeys
	t.Cleanup(func() { releaseKeys = old })
	releaseKeys = []*ecdsa.PublicKey{&key.PublicKey}
	return key
}

// signBlob is the checksums.txt.sig cosign sign-blob writes for data with
// key: base64 of an ASN.1 ECDSA signature over its SHA-256.
func signBlob(t *testing.T, key *ecdsa.PrivateKey, data []byte) []byte {
	t.Helper()
	sum := sha256.Sum256(data)
	sig, err := ecdsa.SignASN1(rand.Reader, key, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(sig) + "\n")
}

// publicPEM is key's public half as cosign.pub holds it.
func publicPEM(t *testing.T, key *ecdsa.PrivateKey) []byte {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

// fakeFlai is one release serveFlaiReleases serves: its version, and its
// assets by name.
type fakeFlai struct {
	version string
	assets  map[string][]byte
}

// flaiRelease is release v whose archive for this platform holds a shell
// script answering "version --json" with v, with its checksums.txt and, when
// key is not nil, the checksums.txt.sig key signs.
func flaiRelease(t *testing.T, v string, key *ecdsa.PrivateKey) fakeFlai {
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
	sums := []byte(hex.EncodeToString(sum[:]) + "  " + archiveName + "\n")
	rel := fakeFlai{version: v, assets: map[string][]byte{archiveName: archive, "checksums.txt": sums}}
	if key != nil {
		rel.assets["checksums.txt.sig"] = signBlob(t, key, sums)
	}
	return rel
}

// serveFlaiReleases answers as GitHub's API does for o/r with rels, newest
// first as given, each asset at /assets/<version>/<name>; asked counts the
// requests for each path so far.
func serveFlaiReleases(t *testing.T, rels ...fakeFlai) (url string, asked func(path string) int) {
	t.Helper()
	var mu sync.Mutex
	count := map[string]int{}
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count[r.URL.Path]++
		mu.Unlock()
		if r.URL.Path == "/repos/o/r/releases" {
			list := make([]string, len(rels))
			for i, rel := range rels {
				var assets []string
				for name := range rel.assets {
					assets = append(assets, fmt.Sprintf(`{"name":%q,"url":"%s/assets/%s/%s"}`, name, srv.URL, rel.version, name))
				}
				list[i] = fmt.Sprintf(`{"tag_name":"flai/v%s","published_at":"2026-10-0%dT09:00:00Z","assets":[%s]}`, rel.version, 9-i, strings.Join(assets, ","))
			}
			fmt.Fprint(w, "["+strings.Join(list, ",")+"]")
			return
		}
		for _, rel := range rels {
			for name, data := range rel.assets {
				if r.URL.Path == "/assets/"+rel.version+"/"+name {
					_, _ = w.Write(data)
					return
				}
			}
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv.URL, func(path string) int {
		mu.Lock()
		defer mu.Unlock()
		return count[path]
	}
}

// S-0111: a flai run from inside a project (a checkout's bin/flai) never
// installs a release there: it goes under the home folder, and whether it is
// up to date is asked of the flai there.
func TestSelfUpgradeNeverInstallsIntoAProject(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a shell script stands in for flai")
	}
	selfUpgradeConfig(t)
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
	if _, installed, err := l.upgradeTo(context.Background(), ""); err != nil || !installed {
		t.Fatalf("upgrade: %v %v", installed, err)
	}
	if got := l.installedTo(); got != "/home/me/.flai/bin/flai" {
		t.Errorf("restarts on %q", got)
	}
}

// releasesStandIn answers as GitHub's API does for o/r: the releases given,
// as JSON. asked lists the paths requested so far.
func releasesStandIn(t *testing.T, releases string) (url string, asked func() []string) {
	t.Helper()
	var mu sync.Mutex
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		if r.URL.Path != "/repos/o/r/releases" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, releases)
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
	cfg, _ := selfUpgradeConfig(t)
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
	selfUpgradeConfig(t)
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
	selfUpgradeConfig(t)
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

// unverifiableReleases serves four releases to a flai whose release key is a
// new test key: 1.3.0 signed by another key, which it publishes as
// cosign.pub; 1.2.0 signed by the release key; 1.1.0 signed by another key it
// does not publish; and 1.0.0, from before signing, with no signature.
func unverifiableReleases(t *testing.T) (api string, asked func(path string) int) {
	t.Helper()
	trusted := trustTestKey(t)
	other, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	unknown := flaiRelease(t, "1.3.0", other)
	unknown.assets["cosign.pub"] = publicPEM(t, other)
	return serveFlaiReleases(t, unknown, flaiRelease(t, "1.2.0", trusted), flaiRelease(t, "1.1.0", other), flaiRelease(t, "1.0.0", nil))
}

// noArchiveAsked fails t when an archive of the versions was downloaded.
func noArchiveAsked(t *testing.T, asked func(string) int, versions ...string) {
	t.Helper()
	for _, v := range versions {
		if n := asked(fmt.Sprintf("/assets/%s/flai_%s_%s_%s.tar.gz", v, v, runtime.GOOS, runtime.GOARCH)); n != 0 {
			t.Errorf("the archive of %s was downloaded %d times", v, n)
		}
	}
}

// S-0233: --list marks each release this flai cannot verify, and why, in text
// and with --json, where a release it can verify carries no mark; a release
// with no signature costs no download, and no archive is downloaded.
func TestSelfUpgradeListMarksUnverifiableReleases(t *testing.T) {
	selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	oldExe := executablePath
	t.Cleanup(func() { executablePath = oldExe })
	executablePath = func() (string, error) { return filepath.Join(t.TempDir(), "flai"), nil }
	api, asked := unverifiableReleases(t)
	dir := t.TempDir()

	out, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("list: %d %s", code, errOut)
	}
	lines := map[string]string{}
	for _, l := range strings.Split(out, "\n") {
		if f := strings.Fields(l); strings.HasPrefix(l, "  ") && len(f) > 0 {
			lines[f[0]] = l
		}
	}
	for v, want := range map[string]string{
		"1.3.0": "latest, not installable: it is signed by a key this flai does not trust",
		"1.1.0": "not installable: its signature does not verify, so it may have been altered",
		"1.0.0": "not installable: it publishes no checksums.txt.sig",
	} {
		if !strings.HasSuffix(lines[v], want) {
			t.Errorf("%s: %q, want it to end %q", v, lines[v], want)
		}
	}
	if l := lines["1.2.0"]; l == "" || strings.Contains(l, "not installable") {
		t.Errorf("1.2.0 is verifiable: %q\n%s", l, out)
	}
	if n := asked("/assets/1.0.0/checksums.txt"); n != 0 {
		t.Errorf("a release with no signature is marked without a download: %d", n)
	}
	noArchiveAsked(t, asked, "1.3.0", "1.2.0", "1.1.0", "1.0.0")

	out, errOut, code = runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api, "--json")
	if code != 0 {
		t.Fatalf("list --json: %d %s", code, errOut)
	}
	var listed []map[string]any
	if err := json.Unmarshal([]byte(out), &listed); err != nil || len(listed) != 4 {
		t.Fatalf("json: %v %s", err, out)
	}
	for i, want := range []map[string]any{
		{"version": "1.3.0", "unverifiable": true, "signature": "unknown-key", "reason": "it is signed by a key this flai does not trust"},
		{"version": "1.2.0"},
		{"version": "1.1.0", "unverifiable": true, "signature": "bad", "reason": "its signature does not verify, so it may have been altered"},
		{"version": "1.0.0", "unverifiable": true, "signature": "missing", "reason": "it publishes no checksums.txt.sig"},
	} {
		for _, k := range []string{"version", "unverifiable", "signature", "reason"} {
			if got, ok := listed[i][k]; got != want[k] || ok != (want[k] != nil) {
				t.Errorf("%v %s: %v, want %v", want["version"], k, got, want[k])
			}
		}
	}
}

// S-0233: a release whose signature could not be checked is not marked as
// one that cannot be verified: the list fails, naming it.
func TestSelfUpgradeListFailsWhenASignatureCannotBeFetched(t *testing.T) {
	selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	trustTestKey(t)
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	t.Cleanup(broken.Close)
	api, _ := releasesStandIn(t, `[{"tag_name":"flai/v1.2.0","assets":[{"name":"checksums.txt","url":"`+broken.URL+`/sums"},{"name":"checksums.txt.sig","url":"`+broken.URL+`/sig"}]},{"tag_name":"flai/v1.0.0"}]`)
	for _, args := range [][]string{nil, {"--json"}} {
		out, errOut, code := runIn(t, t.TempDir(), append([]string{"self-upgrade", "--list", "--repo", "o/r", "--api", api}, args...)...)
		if code == 0 || !strings.Contains(errOut, "check whether this flai can verify release flai/v1.2.0: download checksums.txt of release flai/v1.2.0") || strings.Contains(errOut, "is not installed") {
			t.Errorf("%v: %d %s %s", args, code, out, errOut)
		}
	}
}

// S-0233: self-upgrade refuses a release this flai cannot verify, naming it
// and why, and for an unknown key the newest release it can verify; it exits
// non-zero, downloads no archive, and replaces nothing.
func TestSelfUpgradeRefusesAnUnverifiableRelease(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("tar archives are not used on windows")
	}
	selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	api, asked := unverifiableReleases(t)
	dir := t.TempDir()
	for _, tc := range []struct {
		args []string
		want []string
	}{
		{nil, []string{"release flai/v1.3.0 is not installed: it is signed by the key it publishes as cosign.pub, which this flai does not trust", "install flai 1.2.0 first", "flai self-upgrade --version 1.2.0"}},
		{[]string{"--version", "1.1.0"}, []string{"release flai/v1.1.0 is not installed: the signature over its checksums.txt does not verify", "may have been altered"}},
		{[]string{"--version", "1.0.0"}, []string{"release flai/v1.0.0 is not installed: it publishes no checksums.txt.sig, so its signature cannot be verified"}},
	} {
		out, errOut, code := runIn(t, dir, append([]string{"self-upgrade", "--repo", "o/r", "--api", api, "--dir", dir}, tc.args...)...)
		if code == 0 {
			t.Errorf("%v is installed: %s", tc.args, out)
		}
		for _, w := range tc.want {
			if !strings.Contains(errOut, w) {
				t.Errorf("%v: want %q in %s", tc.args, w, errOut)
			}
		}
	}
	noArchiveAsked(t, asked, "1.3.0", "1.1.0", "1.0.0")
	if _, err := os.Stat(filepath.Join(dir, "flai")); !os.IsNotExist(err) {
		t.Errorf("nothing is installed: %v", err)
	}
	if out, errOut, code := runIn(t, dir, "self-upgrade", "--repo", "o/r", "--api", api, "--dir", dir, "--version", "1.2.0"); code != 0 || !strings.Contains(out, "installed flai 1.2.0") {
		t.Errorf("the verifiable release installs: %d %s %s", code, out, errOut)
	}
}

// S-0233: --check says when the upgrade it reports cannot be verified, and
// checks nothing when there is none to report.
func TestSelfUpgradeCheckSaysTheUpgradeCannotBeVerified(t *testing.T) {
	selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	oldVersion := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = oldVersion })
	buildinfo.Version = "1.1.0"
	oldExe := executablePath
	t.Cleanup(func() { executablePath = oldExe })
	executablePath = func() (string, error) { return filepath.Join(t.TempDir(), "flai"), nil }
	api, asked := unverifiableReleases(t)
	dir := t.TempDir()

	out, errOut, code := runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r", "--api", api)
	if code != 0 || !strings.Contains(out, "latest is 1.3.0 (upgrade available, but not installable: it is signed by a key this flai does not trust; flai self-upgrade --list marks the releases this flai can verify)") {
		t.Errorf("check: %d %s %s", code, out, errOut)
	}
	out, _, _ = runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r", "--api", api, "--json")
	var said map[string]any
	if err := json.Unmarshal([]byte(out), &said); err != nil || said["unverifiable"] != true || said["signature"] != "unknown-key" || said["reason"] != "it is signed by a key this flai does not trust" || said["latest"] != "1.3.0" {
		t.Errorf("check --json: %v %s", err, out)
	}
	out, _, _ = runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r", "--api", api, "--version", "1.2.0", "--json")
	if strings.Contains(out, "unverifiable") || !strings.Contains(out, `"latest": "1.2.0"`) {
		t.Errorf("a verifiable upgrade carries no mark: %s", out)
	}

	buildinfo.Version = "1.3.0"
	before := asked("/assets/1.3.0/checksums.txt")
	out, _, _ = runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r", "--api", api)
	if !strings.Contains(out, "latest is 1.3.0 (up to date)") || asked("/assets/1.3.0/checksums.txt") != before {
		t.Errorf("up to date, nothing is checked: %s", out)
	}
}

// signedAssetsAsked is how many times the checksums.txt, checksums.txt.sig,
// and cosign.pub of the versions were asked for.
func signedAssetsAsked(asked func(string) int, versions ...string) int {
	n := 0
	for _, v := range versions {
		for _, name := range []string{"checksums.txt", "checksums.txt.sig", "cosign.pub"} {
			n += asked("/assets/" + v + "/" + name)
		}
	}
	return n
}

// I-0142: a second --list marks the same releases as the first without
// downloading any signed release's checksums.txt or signature again, and
// --check answers from the same verdicts.
func TestSelfUpgradeListKeepsEachVerdict(t *testing.T) {
	_, cache := selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	oldVersion, oldExe := buildinfo.Version, executablePath
	t.Cleanup(func() { buildinfo.Version, executablePath = oldVersion, oldExe })
	buildinfo.Version = "1.1.0"
	installed := filepath.Join(t.TempDir(), "flai")
	executablePath = func() (string, error) { return installed, nil }
	api, asked := unverifiableReleases(t)
	dir := t.TempDir()
	signed := []string{"1.3.0", "1.2.0", "1.1.0"}

	first, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("first list: %d %s", code, errOut)
	}
	downloaded := signedAssetsAsked(asked, signed...)
	if downloaded == 0 {
		t.Fatal("the first list downloads the signed releases' checksums.txt and signature")
	}
	if _, err := os.Stat(filepath.Join(cache, "releases", "verdicts.json")); err != nil {
		t.Errorf("the verdicts are kept under cache_dir: %v", err)
	}
	second, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("second list: %d %s", code, errOut)
	}
	if n := signedAssetsAsked(asked, signed...); n != downloaded {
		t.Errorf("the second list downloaded %d signed assets again", n-downloaded)
	}
	if second != first {
		t.Errorf("the second list marks the releases as the first did:\n%s\nwant\n%s", second, first)
	}
	out, errOut, code := runIn(t, dir, "self-upgrade", "--check", "--repo", "o/r", "--api", api)
	if code != 0 || !strings.Contains(out, "latest is 1.3.0 (upgrade available, but not installable: it is signed by a key this flai does not trust") {
		t.Errorf("check: %d %s %s", code, out, errOut)
	}
	if n := signedAssetsAsked(asked, signed...); n != downloaded {
		t.Errorf("--check downloaded %d signed assets the list had checked", n-downloaded)
	}
	noArchiveAsked(t, asked, "1.3.0", "1.2.0", "1.1.0", "1.0.0")
}

// I-0142: a release whose signature asset is replaced, and so has a new URL,
// is downloaded and verified again, and marked by the new verdict; a release
// whose assets are unchanged is not downloaded again.
func TestSelfUpgradeListChecksAReplacedSignatureAgain(t *testing.T) {
	selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	trusted := trustTestKey(t)
	other, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	// 1.2.0 was published with a signature by a key this flai does not
	// trust, which is replaced by one the release key makes
	replaced := flaiRelease(t, "1.2.0", other)
	replaced.assets["replaced.sig"] = signBlob(t, trusted, replaced.assets["checksums.txt"])
	assets, asked := serveFlaiReleases(t, replaced, flaiRelease(t, "1.1.0", trusted))
	releases := func(sig120 string) string {
		at := func(v, name string) string { return fmt.Sprintf(`"%s/assets/%s/%s"`, assets, v, name) }
		return `[{"tag_name":"flai/v1.2.0","published_at":"2026-10-09T09:00:00Z","assets":[{"name":"checksums.txt","url":` + at("1.2.0", "checksums.txt") + `},{"name":"checksums.txt.sig","url":` + at("1.2.0", sig120) + `}]},` +
			`{"tag_name":"flai/v1.1.0","published_at":"2026-10-08T09:00:00Z","assets":[{"name":"checksums.txt","url":` + at("1.1.0", "checksums.txt") + `},{"name":"checksums.txt.sig","url":` + at("1.1.0", "checksums.txt.sig") + `}]}]`
	}
	before, _ := releasesStandIn(t, releases("checksums.txt.sig"))
	after, _ := releasesStandIn(t, releases("replaced.sig"))
	dir := t.TempDir()
	list := func(api string) map[string]any {
		t.Helper()
		out, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api, "--json")
		var listed []map[string]any
		if err := json.Unmarshal([]byte(out), &listed); code != 0 || err != nil || len(listed) != 2 {
			t.Fatalf("list: %d %v %s %s", code, err, out, errOut)
		}
		return listed[0]
	}

	if got := list(before); got["signature"] != "bad" {
		t.Fatalf("1.2.0 before its signature is replaced: %v", got)
	}
	unchanged := signedAssetsAsked(asked, "1.1.0")
	if got := list(after); got["unverifiable"] != nil {
		t.Errorf("1.2.0 verifies once its signature is replaced: %v", got)
	}
	if n := asked("/assets/1.2.0/replaced.sig"); n != 1 {
		t.Errorf("the replaced signature was downloaded %d times, want 1", n)
	}
	if n := signedAssetsAsked(asked, "1.1.0"); n != unchanged {
		t.Errorf("1.1.0, unchanged, was downloaded again: %d more", n-unchanged)
	}
}

// I-0142: a listing whose signature download fails keeps no verdict, so the
// next listing downloads it again, and keeps the verdict it finds then.
func TestSelfUpgradeListKeepsNoVerdictFromAFailedDownload(t *testing.T) {
	_, cache := selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	rel := flaiRelease(t, "1.2.0", trustTestKey(t))
	var down atomic.Bool
	down.Store(true)
	var mu sync.Mutex
	count := map[string]int{}
	flaky := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		mu.Lock()
		count[name]++
		mu.Unlock()
		if down.Load() {
			http.Error(w, "down", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(rel.assets[name])
	}))
	t.Cleanup(flaky.Close)
	asked := func() int {
		mu.Lock()
		defer mu.Unlock()
		return count["checksums.txt"] + count["checksums.txt.sig"]
	}
	api, _ := releasesStandIn(t, `[{"tag_name":"flai/v1.2.0","assets":[{"name":"checksums.txt","url":"`+flaky.URL+`/checksums.txt"},{"name":"checksums.txt.sig","url":"`+flaky.URL+`/checksums.txt.sig"}]}]`)
	dir := t.TempDir()
	verdicts := filepath.Join(cache, "releases", "verdicts.json")

	if out, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api); code == 0 {
		t.Fatalf("the list fails while the download does: %s %s", out, errOut)
	}
	if _, err := os.Stat(verdicts); !os.IsNotExist(err) {
		t.Errorf("a failed download keeps no verdict: %v", err)
	}
	failed := asked()
	down.Store(false)
	if out, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api); code != 0 || strings.Contains(out, "not installable") {
		t.Fatalf("the next list checks 1.2.0: %d %s %s", code, out, errOut)
	}
	checked := asked()
	if checked == failed {
		t.Error("the next list downloads the signature again")
	}
	if _, err := os.Stat(verdicts); err != nil {
		t.Errorf("the verdict found then is kept: %v", err)
	}
	if _, _, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api); code != 0 || asked() != checked {
		t.Errorf("a third list downloads nothing: %d, %d more", code, asked()-checked)
	}
}

// I-0142: installing a release a listing has checked still downloads its
// checksums.txt and signature and verifies them, rather than reading the
// listing's verdict.
func TestSelfUpgradeInstallVerifiesAfterACachedListing(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("tar archives are not used on windows")
	}
	selfUpgradeConfig(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	api, asked := unverifiableReleases(t)
	dir := t.TempDir()
	if _, errOut, code := runIn(t, dir, "self-upgrade", "--list", "--repo", "o/r", "--api", api); code != 0 {
		t.Fatalf("list: %d %s", code, errOut)
	}
	sums, sig := asked("/assets/1.2.0/checksums.txt"), asked("/assets/1.2.0/checksums.txt.sig")
	if out, errOut, code := runIn(t, dir, "self-upgrade", "--repo", "o/r", "--api", api, "--dir", dir, "--version", "1.2.0"); code != 0 || !strings.Contains(out, "installed flai 1.2.0") {
		t.Fatalf("install: %d %s %s", code, out, errOut)
	}
	if asked("/assets/1.2.0/checksums.txt") == sums || asked("/assets/1.2.0/checksums.txt.sig") != sig+1 {
		t.Errorf("the install downloads checksums.txt and its signature again: %d to %d, %d to %d", sums, asked("/assets/1.2.0/checksums.txt"), sig, asked("/assets/1.2.0/checksums.txt.sig"))
	}
}
