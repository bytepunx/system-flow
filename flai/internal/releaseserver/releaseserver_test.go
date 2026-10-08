package releaseserver

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

const repo = "local/flai"

// archive is a tar.gz holding a file named flai with content.
func archive(t *testing.T, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "flai", Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// dist writes a GoReleaser-style dist folder for version on linux/amd64,
// with a checksums.txt holding the archive's real sha256 and the files
// GoReleaser writes beside them, and answers its path.
func dist(t *testing.T, version, content string) string {
	t.Helper()
	dir := t.TempDir()
	name := selfupgrade.ArchiveName(version, "linux", "amd64")
	data := archive(t, content)
	sum := sha256.Sum256(data)
	write := func(n string, b []byte) {
		if err := os.WriteFile(filepath.Join(dir, n), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(name, data)
	write(ChecksumsName, fmt.Appendf(nil, "%s  %s\n", hex.EncodeToString(sum[:]), name))
	write("artifacts.json", []byte("[]"))
	write("flai_snapshot_linux_amd64.tar.gz", []byte("not a plain version"))
	if err := os.Mkdir(filepath.Join(dir, "flai_linux_amd64_v1"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// serve scans dirs and serves them under repo.
func serve(t *testing.T, dirs ...string) *httptest.Server {
	t.Helper()
	releases, err := Scan(dirs...)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(New(repo, releases))
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url, accept string) (int, http.Header, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, resp.Header, body
}

func TestResolveAndDownloadAgainstTheServer(t *testing.T) {
	srv := serve(t, dist(t, "9.9.9", "the new flai"), dist(t, "9.9.8", "the old flai"))
	opt := selfupgrade.Options{Repo: repo, APIBase: srv.URL, OS: "linux", Arch: "amd64"}
	ctx := context.Background()

	rel, err := selfupgrade.Resolve(ctx, opt)
	if err != nil {
		t.Fatalf("resolve the newest: %v", err)
	}
	if rel.Tag != "flai/v9.9.9" || rel.Version != "9.9.9" {
		t.Fatalf("resolved %s (%s), want flai/v9.9.9", rel.Tag, rel.Version)
	}
	if rel.Published.IsZero() {
		t.Error("resolved release has no publish time")
	}
	bin, err := selfupgrade.Download(ctx, opt, rel)
	if err != nil {
		t.Fatalf("download the newest: %v", err)
	}
	if string(bin) != "the new flai" {
		t.Errorf("downloaded %q, want the new flai", bin)
	}

	opt.Version = "9.9.8"
	old, err := selfupgrade.Resolve(ctx, opt)
	if err != nil {
		t.Fatalf("resolve the pinned 9.9.8: %v", err)
	}
	bin, err = selfupgrade.Download(ctx, opt, old)
	if err != nil {
		t.Fatalf("download 9.9.8: %v", err)
	}
	if string(bin) != "the old flai" {
		t.Errorf("downloaded %q, want the old flai", bin)
	}

	listed, err := selfupgrade.List(ctx, opt, selfupgrade.TagPrefix)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 2 {
		t.Fatalf("listed %d releases, want 2: %+v", len(listed), listed)
	}
}

func TestListingIsNewestFirstWithNoLinkHeader(t *testing.T) {
	srv := serve(t, dist(t, "1.9.0", "a"), dist(t, "1.10.0", "b"))
	code, h, body := get(t, srv.URL+"/repos/"+repo+"/releases?per_page=50", "application/vnd.github+json")
	if code != http.StatusOK {
		t.Fatalf("listing answered %d: %s", code, body)
	}
	if link := h.Get("Link"); link != "" {
		t.Errorf("listing sent Link %q, want none", link)
	}
	// install.sh takes the first tag_name in the listing as the newest.
	first := regexp.MustCompile(`"tag_name": *"flai/v[^"]*"`).Find(body)
	if string(first) != `"tag_name":"flai/v1.10.0"` {
		t.Errorf("first tag in the listing is %s, want flai/v1.10.0", first)
	}
}

func TestTagRouteTakesAnEscapedSlash(t *testing.T) {
	srv := serve(t, dist(t, "9.9.9", "x"))
	for _, tag := range []string{"flai%2Fv9.9.9", "flai%2fv9.9.9", "flai/v9.9.9"} {
		code, _, body := get(t, srv.URL+"/repos/"+repo+"/releases/tags/"+tag, "application/vnd.github+json")
		if code != http.StatusOK {
			t.Fatalf("tag %s answered %d: %s", tag, code, body)
		}
		var rel struct {
			TagName string `json:"tag_name"`
			Assets  []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"assets"`
		}
		if err := json.Unmarshal(body, &rel); err != nil {
			t.Fatal(err)
		}
		if rel.TagName != "flai/v9.9.9" {
			t.Errorf("tag %s answered %s", tag, rel.TagName)
		}
		var names []string
		for _, a := range rel.Assets {
			names = append(names, a.Name)
			if !regexp.MustCompile(`/releases/assets/[0-9]+$`).MatchString(a.URL) {
				t.Errorf("asset %s has url %s, which install.sh's asset_url does not match", a.Name, a.URL)
			}
		}
		if got := strings.Join(names, ","); got != "checksums.txt,flai_9.9.9_linux_amd64.tar.gz" {
			t.Errorf("assets are %s, want checksums.txt and the archive only", got)
		}
	}
}

// install.sh pairs each asset's url with the name after it, so url must come
// first in each asset's JSON.
func TestAssetURLComesBeforeItsName(t *testing.T) {
	srv := serve(t, dist(t, "9.9.9", "x"))
	_, _, body := get(t, srv.URL+"/repos/"+repo+"/releases/tags/flai%2Fv9.9.9", "application/vnd.github+json")
	pairs := regexp.MustCompile(`"(url|name)": *"[^"]*"`).FindAllString(string(body), -1)
	var keys []string
	for _, p := range pairs {
		keys = append(keys, p[1:strings.Index(p[1:], `"`)+1])
	}
	// the release's name, then url and name for each of two assets
	if got := strings.Join(keys, ","); got != "name,url,name,url,name" {
		t.Errorf("url and name keys come in the order %s, want name,url,name,url,name", got)
	}
}

func TestAssetAnswersBytesOnlyForOctetStream(t *testing.T) {
	dir := dist(t, "9.9.9", "x")
	srv := serve(t, dir)
	want, err := os.ReadFile(filepath.Join(dir, ChecksumsName))
	if err != nil {
		t.Fatal(err)
	}
	url := srv.URL + "/repos/" + repo + "/releases/assets/1"
	code, h, body := get(t, url, "application/octet-stream")
	if code != http.StatusOK || !bytes.Equal(body, want) {
		t.Fatalf("asset 1 answered %d %q, want 200 with checksums.txt", code, body)
	}
	if ct := h.Get("Content-Type"); ct != "application/octet-stream" {
		t.Errorf("asset bytes answered Content-Type %q", ct)
	}
	code, _, body = get(t, url, "application/vnd.github+json")
	var desc struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}
	if code != http.StatusOK || json.Unmarshal(body, &desc) != nil || desc.Name != ChecksumsName || desc.Size != int64(len(want)) {
		t.Errorf("asset 1 without octet-stream answered %d %s, want its description", code, body)
	}
}

func TestUnknownAnswers404(t *testing.T) {
	srv := serve(t, dist(t, "9.9.9", "x"))
	for _, path := range []string{
		"/repos/" + repo + "/releases/tags/flai%2Fv1.0.0",
		"/repos/" + repo + "/releases/assets/99",
		"/repos/" + repo + "/releases/assets/x",
		"/repos/other/repo/releases",
		"/repos/" + repo + "/git/matching-refs/tags/flaiover/v",
		"/",
	} {
		if code, _, body := get(t, srv.URL+path, "application/octet-stream"); code != http.StatusNotFound {
			t.Errorf("%s answered %d %s, want 404", path, code, body)
		}
	}
	resp, err := http.Post(srv.URL+"/repos/"+repo+"/releases", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST to the listing answered %d, want 404", resp.StatusCode)
	}
	rel, err := selfupgrade.Resolve(context.Background(), selfupgrade.Options{Repo: repo, APIBase: srv.URL, Version: "1.0.0"})
	if err == nil {
		t.Errorf("resolve an unserved version answered %+v, want an error", rel)
	}
}

func TestScanRefusesWhatItCannotServe(t *testing.T) {
	empty := t.TempDir()
	noSums := t.TempDir()
	if err := os.WriteFile(filepath.Join(noSums, "flai_1.0.0_linux_amd64.tar.gz"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	twice := dist(t, "2.0.0", "a")
	cases := []struct {
		name string
		dirs []string
	}{
		{"no folder", nil},
		{"a missing folder", []string{filepath.Join(empty, "missing")}},
		{"no archive", []string{empty}},
		{"no checksums", []string{noSums}},
		{"a version twice", []string{twice, dist(t, "2.0.0", "b")}},
	}
	for _, c := range cases {
		if _, err := Scan(c.dirs...); err == nil {
			t.Errorf("Scan with %s answered no error", c.name)
		}
	}
}

func TestPublishedIsTheArchiveModTime(t *testing.T) {
	dir := dist(t, "9.9.9", "x")
	at := time.Date(2026, 10, 8, 9, 30, 0, 0, time.UTC)
	if err := os.Chtimes(filepath.Join(dir, "flai_9.9.9_linux_amd64.tar.gz"), at, at); err != nil {
		t.Fatal(err)
	}
	releases, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 || !releases[0].Published.Equal(at) {
		t.Errorf("published %+v, want one release at %s", releases, at)
	}
}
