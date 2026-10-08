package selfupgrade

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func tarball(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "README.md", Mode: 0o644, Size: 5, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("hello"))
	_ = tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg})
	_, _ = tw.Write(content)
	_ = tw.Close()
	_ = gz.Close()
	return buf.Bytes()
}

// fakeGitHub serves flai releases over two pages of the releases API, newest
// first by publish date as the API does, mixed with what is not published: a
// draft, a prerelease, a suffixed version, a leading zero, and a flaiover
// release. It serves the flaiover tags over two pages of matching refs.
func fakeGitHub(t *testing.T, archive, sums []byte, wantToken string) *httptest.Server {
	t.Helper()
	return fakeGitHubBehind(t, archive, sums, wantToken, func(h http.Handler) http.Handler { return h })
}

// fakeGitHubBehind is fakeGitHub with front standing before its handler.
func fakeGitHubBehind(t *testing.T, archive, sums []byte, wantToken string, front func(http.Handler) http.Handler) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(front(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantToken != "" && r.Header.Get("Authorization") != "Bearer "+wantToken {
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		assets := func(v string) string {
			return fmt.Sprintf(`"assets":[{"name":"checksums.txt","url":"%s/assets/sums"},{"name":"flai_%s_linux_amd64.tar.gz","url":"%s/assets/archive"},{"name":"flai_%s_windows_amd64.zip","url":"%s/assets/zip"}]`, srv.URL, v, srv.URL, v, srv.URL)
		}
		page := r.URL.Query().Get("page")
		switch {
		case r.URL.Path == "/repos/o/r/releases" && page == "":
			if r.URL.Query().Get("per_page") != "100" {
				http.Error(w, "want per_page=100", http.StatusBadRequest)
				return
			}
			w.Header().Set("Link", fmt.Sprintf(`<%s/repos/o/r/releases?per_page=100&page=2>; rel="next", <%s/repos/o/r/releases?per_page=100&page=2>; rel="last"`, srv.URL, srv.URL))
			fmt.Fprintf(w, `[{"tag_name":"flaiover/v9.0.0","assets":[]},{"tag_name":"flai/v1.12.0","draft":true,"published_at":null,%s},{"tag_name":"flai/v1.9.0","published_at":"2026-09-01T10:00:00Z",%s},{"tag_name":"flai/v1.11.0-rc.1","published_at":"2026-08-30T10:00:00Z",%s}]`, assets("1.12.0"), assets("1.9.0"), assets("1.11.0-rc.1"))
		case r.URL.Path == "/repos/o/r/releases" && page == "2":
			w.Header().Set("Link", fmt.Sprintf(`<%s/repos/o/r/releases?per_page=100&page=1>; rel="prev", <%s/repos/o/r/releases?per_page=100&page=1>; rel="first"`, srv.URL, srv.URL))
			fmt.Fprintf(w, `[{"tag_name":"flai/v1.11.0","prerelease":true,"published_at":"2026-08-20T10:00:00Z",%s},{"tag_name":"flai/v1.10.0","published_at":"2026-08-10T10:00:00Z",%s},{"tag_name":"flai/v01.13.0","published_at":"2026-08-05T10:00:00Z","assets":[]},{"tag_name":"flai/v1.0.0","published_at":"2026-01-02T10:00:00Z",%s}]`, assets("1.11.0"), assets("1.10.0"), assets("1.0.0"))
		case r.URL.Path == "/repos/o/r/git/matching-refs/tags/flaiover/v" && page == "":
			w.Header().Set("Link", fmt.Sprintf(`<%s/repos/o/r/git/matching-refs/tags/flaiover/v?page=2>; rel="next"`, srv.URL))
			_, _ = w.Write([]byte(`[{"ref":"refs/tags/flaiover/v0.9.0"},{"ref":"refs/tags/flaiover/v0.39.0"},{"ref":"refs/tags/flaiover/v0.40.0-rc1"}]`))
		case r.URL.Path == "/repos/o/r/git/matching-refs/tags/flaiover/v" && page == "2":
			_, _ = w.Write([]byte(`[{"ref":"refs/tags/flaiover/v0.40.0"},{"ref":"refs/tags/flaiover/v0.10.0"}]`))
		case r.URL.Path == "/assets/archive":
			if r.Header.Get("Accept") != "application/octet-stream" {
				http.Error(w, "wrong accept", http.StatusBadRequest)
				return
			}
			_, _ = w.Write(archive)
		case r.URL.Path == "/assets/sums":
			_, _ = w.Write(sums)
		default:
			http.NotFound(w, r)
		}
	})))
	t.Cleanup(srv.Close)
	return srv
}

// dropping stands in front of fakeGitHub: it cuts the first answers to a path
// off mid-body, as GitHub sometimes does (I-0086), and counts the requests
// for each path.
type dropping struct {
	t     *testing.T
	next  http.Handler
	mu    sync.Mutex
	drops map[string]int // answers still to cut, by path
	asked map[string]int
}

func dropFirst(t *testing.T, drops map[string]int) (*dropping, func(http.Handler) http.Handler) {
	d := &dropping{t: t, drops: drops, asked: map[string]int{}}
	return d, func(h http.Handler) http.Handler { d.next = h; return d }
}

func (d *dropping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	d.asked[r.URL.Path]++
	drop := d.drops[r.URL.Path] > 0
	if drop {
		d.drops[r.URL.Path]--
	}
	d.mu.Unlock()
	if !drop {
		d.next.ServeHTTP(w, r)
		return
	}
	conn, buf, err := w.(http.Hijacker).Hijack()
	if err != nil {
		d.t.Errorf("hijack %s: %v", r.URL.Path, err)
		return
	}
	_, _ = buf.WriteString("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 4096\r\n\r\n[{\"tag_name\":\"flai/v1.")
	_ = buf.Flush()
	_ = conn.Close()
}

func (d *dropping) count(path string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.asked[path]
}

// pauseFor sets the wait between attempts for one test.
func pauseFor(t *testing.T, pause time.Duration) {
	old := retryPause
	retryPause = pause
	t.Cleanup(func() { retryPause = old })
}

func versions(rels []Release) string {
	var out []string
	for _, r := range rels {
		out = append(out, r.Version)
	}
	return strings.Join(out, " ")
}

func TestList(t *testing.T) {
	srv := fakeGitHub(t, nil, nil, "tok")
	opt := Options{Repo: "o/r", APIBase: srv.URL, Token: "tok"}

	flai, err := List(context.Background(), opt, TagPrefix)
	if err != nil {
		t.Fatal(err)
	}
	if got := versions(flai); got != "1.10.0 1.9.0 1.0.0" {
		t.Errorf("flai releases across pages, newest first by semver: %q", got)
	}
	want := time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC)
	if flai[0].Tag != "flai/v1.10.0" || !flai[0].Published.Equal(want) || len(flai[0].Assets) != 3 {
		t.Errorf("newest: %+v", flai[0])
	}

	dashboard, err := ListTags(context.Background(), opt, "flaiover/v")
	if err != nil {
		t.Fatal(err)
	}
	if got := versions(dashboard); got != "0.40.0 0.39.0 0.10.0 0.9.0" {
		t.Errorf("flaiover tags across pages, newest first by semver: %q", got)
	}
	if dashboard[0].Tag != "flaiover/v0.40.0" || !dashboard[0].Published.IsZero() || dashboard[0].Assets != nil {
		t.Errorf("a tag has no publish date or assets: %+v", dashboard[0])
	}

	if _, err := List(context.Background(), Options{Repo: "o/r", APIBase: srv.URL}, TagPrefix); err == nil || !strings.Contains(err.Error(), "list the flai/v releases of o/r") {
		t.Errorf("an unauthorized list names what it listed: %v", err)
	}
}

func TestResolvePinned(t *testing.T) {
	srv := fakeGitHub(t, nil, nil, "tok")
	for _, pin := range []string{"1.9.0", "v1.9.0", "flai/v1.9.0"} {
		rel, err := Resolve(context.Background(), Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", Version: pin})
		if err != nil || rel.Tag != "flai/v1.9.0" || rel.Version != "1.9.0" {
			t.Errorf("pinned %s: %+v %v", pin, rel, err)
		}
	}
	for _, pin := range []string{"1.11.0", "1.12.0", "1.11.0-rc.1", "2.0.0"} {
		_, err := Resolve(context.Background(), Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", Version: pin})
		if err == nil || !strings.Contains(err.Error(), "published are 1.10.0, 1.9.0, 1.0.0:") || !strings.Contains(err.Error(), "flai/v"+pin) {
			t.Errorf("pinned %s is not published and should be refused naming the published ones: %v", pin, err)
		}
	}
}

func TestVersionList(t *testing.T) {
	var rels []Release
	for i := 12; i > 0; i-- {
		rels = append(rels, Release{Version: fmt.Sprintf("1.%d.0", i)})
	}
	if got := VersionList(rels); got != "1.12.0, 1.11.0, 1.10.0, 1.9.0, 1.8.0, 1.7.0, 1.6.0, 1.5.0, 1.4.0, 1.3.0 and 2 older" {
		t.Errorf("twelve: %q", got)
	}
	if got := VersionList(rels[:2]); got != "1.12.0, 1.11.0" {
		t.Errorf("two: %q", got)
	}
}

func TestNextLink(t *testing.T) {
	for header, want := range map[string]string{
		"": "",
		`<https://api/x?page=2>; rel="next", <https://api/x?page=5>; rel="last"`: "https://api/x?page=2",
		`<https://api/x?page=1>; rel="prev", <https://api/x?page=3>; rel="next"`: "https://api/x?page=3",
		`<https://api/x?page=1>; rel="first"`:                                    "",
	} {
		if got := nextLink(header); got != want {
			t.Errorf("nextLink(%q) = %q, want %q", header, got, want)
		}
	}
}

func TestResolveDownloadReplace(t *testing.T) {
	content := []byte("#!/bin/sh\necho new\n")
	archive := tarball(t, "flai", content)
	sum := sha256.Sum256(archive)
	sums := []byte(hex.EncodeToString(sum[:]) + "  flai_1.10.0_linux_amd64.tar.gz\n" + "deadbeef  flai_1.0.0_linux_amd64.tar.gz\n")
	srv := fakeGitHub(t, archive, sums, "tok")
	opt := Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", OS: "linux", Arch: "amd64"}

	rel, err := Resolve(context.Background(), opt)
	if err != nil || rel.Version != "1.10.0" || rel.Tag != "flai/v1.10.0" {
		t.Fatalf("latest: %+v %v", rel, err)
	}
	bin, err := Download(context.Background(), opt, rel)
	if err != nil || !bytes.Equal(bin, content) {
		t.Fatalf("download: %v", err)
	}
	dest := filepath.Join(t.TempDir(), "flai")
	_ = os.WriteFile(dest, []byte("old"), 0o755)
	if err := Replace(dest, bin); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(dest)
	st, _ := os.Stat(dest)
	if !bytes.Equal(got, content) || st.Mode().Perm()&0o111 == 0 {
		t.Errorf("replaced: %q mode %v", got, st.Mode())
	}
	if leftovers, _ := filepath.Glob(filepath.Join(filepath.Dir(dest), ".flai-upgrade-*")); len(leftovers) != 0 {
		t.Errorf("temp files left: %v", leftovers)
	}

	pinned, err := Resolve(context.Background(), Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", Version: "v1.0.0"})
	if err != nil || pinned.Version != "1.0.0" {
		t.Fatalf("pinned: %+v %v", pinned, err)
	}
	if _, err := Download(context.Background(), opt, pinned); err == nil {
		t.Error("checksum mismatch for 1.0.0 should fail")
	}
	if _, err := Download(context.Background(), Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", OS: "darwin", Arch: "arm64"}, rel); err == nil {
		t.Error("missing platform asset should fail")
	}
	if _, err := Resolve(context.Background(), Options{Repo: "o/r", APIBase: srv.URL}); err == nil {
		t.Error("missing token against a private repository should fail")
	}
}

func TestVerifyAndNames(t *testing.T) {
	if ArchiveName("1.2.3", "windows", "arm64") != "flai_1.2.3_windows_arm64.zip" || ArchiveName("1.2.3", "linux", "amd64") != "flai_1.2.3_linux_amd64.tar.gz" {
		t.Error("archive names")
	}
	data := []byte("x")
	sum := sha256.Sum256(data)
	ok := []byte(hex.EncodeToString(sum[:]) + "  a.tar.gz\n")
	if err := Verify(data, "a.tar.gz", ok); err != nil {
		t.Error(err)
	}
	if err := Verify(data, "b.tar.gz", ok); err == nil {
		t.Error("missing entry should fail")
	}
	if err := Verify([]byte("y"), "a.tar.gz", ok); err == nil {
		t.Error("mismatch should fail")
	}
}

func TestDroppedAnswerIsAskedAgain(t *testing.T) {
	pauseFor(t, 0)
	content := []byte("#!/bin/sh\necho new\n")
	archive := tarball(t, "flai", content)
	sum := sha256.Sum256(archive)
	sums := []byte(hex.EncodeToString(sum[:]) + "  flai_1.10.0_linux_amd64.tar.gz\n")
	d, front := dropFirst(t, map[string]int{"/repos/o/r/releases": 2, "/assets/archive": 1, "/assets/sums": 1})
	srv := fakeGitHubBehind(t, archive, sums, "tok", front)
	opt := Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", OS: "linux", Arch: "amd64"}

	rel, err := Resolve(context.Background(), opt)
	if err != nil || rel.Version != "1.10.0" {
		t.Fatalf("the first page dropped twice mid-body still resolves the newest: %+v %v", rel, err)
	}
	if got := d.count("/repos/o/r/releases"); got != 4 {
		t.Errorf("three attempts at the first page and one at the second, got %d requests", got)
	}
	bin, err := Download(context.Background(), opt, rel)
	if err != nil || !bytes.Equal(bin, content) {
		t.Fatalf("an archive and checksums dropped once mid-body still download: %v", err)
	}
	if d.count("/assets/archive") != 2 || d.count("/assets/sums") != 2 {
		t.Errorf("each download asked twice: archive %d, sums %d", d.count("/assets/archive"), d.count("/assets/sums"))
	}
}

func TestDroppedAnswerIsAskedThreeTimes(t *testing.T) {
	pauseFor(t, 0)
	d, front := dropFirst(t, map[string]int{"/repos/o/r/releases": 3})
	srv := fakeGitHubBehind(t, nil, nil, "", front)

	_, err := List(context.Background(), Options{Repo: "o/r", APIBase: srv.URL}, TagPrefix)
	if err == nil || !strings.Contains(err.Error(), "unexpected EOF") || !strings.Contains(err.Error(), "gave up after 3 attempts") {
		t.Errorf("the last attempt's error, saying how many were made: %v", err)
	}
	if got := d.count("/repos/o/r/releases"); got != 3 {
		t.Errorf("three attempts, got %d requests", got)
	}
}

func TestFailedAnswerIsNotAskedAgain(t *testing.T) {
	pauseFor(t, 0)
	d, front := dropFirst(t, nil)
	srv := fakeGitHubBehind(t, nil, nil, "tok", front)

	_, err := List(context.Background(), Options{Repo: "o/r", APIBase: srv.URL}, TagPrefix)
	if err == nil || !strings.Contains(err.Error(), "404 Not Found (private repository?") || strings.Contains(err.Error(), "attempts") {
		t.Errorf("a 404 fails at once with its hint: %v", err)
	}
	if got := d.count("/repos/o/r/releases"); got != 1 {
		t.Errorf("a 404 is asked once, got %d requests", got)
	}
}

func TestCancelStopsTheWaitToAskAgain(t *testing.T) {
	pauseFor(t, time.Hour)
	d, front := dropFirst(t, map[string]int{"/repos/o/r/releases": 1})
	srv := fakeGitHubBehind(t, nil, nil, "", front)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := List(ctx, Options{Repo: "o/r", APIBase: srv.URL}, TagPrefix)
	if !errors.Is(err, context.DeadlineExceeded) || !strings.Contains(err.Error(), "unexpected EOF") {
		t.Errorf("the deadline and the dropped answer: %v", err)
	}
	if got := d.count("/repos/o/r/releases"); got != 1 {
		t.Errorf("no attempt after the deadline, got %d requests", got)
	}
}
