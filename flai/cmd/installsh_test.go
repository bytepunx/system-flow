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
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// I-0086: GitHub drops the connection partway through the releases listing
// now and then. install.sh tries the listing again and installs; it lists
// releases a small page at a time, asking for the next page only while no
// flai release has turned up.
func TestInstallShRetriesAListingTheNetworkDrops(t *testing.T) {
	api, listed := installStandIn(t, func(w http.ResponseWriter, n int, page string) {
		switch {
		case n == 1:
			cutOff(t, w, `[{"tag_name":"flaiover/v0.9.0"},{"tag_name":"flaiover/v0.8.0"}]`)
		case page == "1":
			fmt.Fprint(w, `[{"tag_name":"flaiover/v0.9.0"},{"tag_name":"flaiover/v0.8.0"}]`)
		case page == "2":
			fmt.Fprint(w, `[{"tag_name":"flai/v1.2.3"},{"tag_name":"flai/v1.2.2"}]`)
		default:
			fmt.Fprint(w, `[]`)
		}
	})
	dir := filepath.Join(t.TempDir(), "bin")
	out, code := runInstallSh(t, api, dir)
	if code != 0 {
		t.Fatalf("install.sh: exit %d\n%s", code, out)
	}
	if !strings.Contains(out, "(curl exit ") || !strings.Contains(out, "attempt 2 of 3") {
		t.Errorf("the retry is not reported:\n%s", out)
	}
	if !strings.Contains(out, "Installed "+filepath.Join(dir, "flai")+" (flai 1.2.3)") {
		t.Errorf("flai 1.2.3 is not reported installed:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, "flai")); err != nil {
		t.Errorf("flai is not in FLAI_INSTALL_DIR: %v", err)
	}
	want := []string{"per_page=10&page=1", "per_page=10&page=1", "per_page=10&page=2"}
	if got := listed(); !slices.Equal(got, want) {
		t.Errorf("listings asked for:\n got %q\nwant %q", got, want)
	}
}

// An HTTP error is no dropped connection: a 404 for the listing, as GitHub
// answers for a private repository without a token, fails at once with the
// hint, and the listing is asked for once.
func TestInstallShFailsAtOnceOnAnHTTPError(t *testing.T) {
	api, listed := installStandIn(t, func(w http.ResponseWriter, _ int, _ string) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"Not Found"}`)
	})
	dir := filepath.Join(t.TempDir(), "bin")
	out, code := runInstallSh(t, api, dir)
	if code == 0 || !strings.Contains(out, "could not list releases of o/r (private repository?") {
		t.Errorf("install.sh: exit %d, want a failure with the hint\n%s", code, out)
	}
	if strings.Contains(out, "trying again") {
		t.Errorf("an HTTP error is retried:\n%s", out)
	}
	if got := listed(); len(got) != 1 {
		t.Errorf("listings asked for: %q, want one", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "flai")); !os.IsNotExist(err) {
		t.Errorf("something was installed: %v", err)
	}
}

// installStandIn answers install.sh as GitHub's releases API does for o/r:
// listing answers each releases listing, given its count from one and the
// page asked for, and release flai/v1.2.3 holds an archive for every platform
// install.sh knows, each a flai that prints its version, and their
// checksums.txt. listed returns the query of each listing asked for so far.
func installStandIn(t *testing.T, listing func(w http.ResponseWriter, n int, page string)) (url string, listed func() []string) {
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
	var names, sums []string
	for _, goos := range []string{"linux", "darwin"} {
		for _, arch := range []string{"amd64", "arm64"} {
			name := fmt.Sprintf("flai_1.2.3_%s_%s.tar.gz", goos, arch)
			names = append(names, name)
			sums = append(sums, hex.EncodeToString(sum[:])+"  "+name+"\n")
		}
	}

	var mu sync.Mutex
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		assets := base + "/repos/o/r/releases/assets/"
		switch r.URL.Path {
		case "/repos/o/r/releases":
			mu.Lock()
			queries = append(queries, r.URL.RawQuery)
			n := len(queries)
			mu.Unlock()
			listing(w, n, r.URL.Query().Get("page"))
		case "/repos/o/r/releases/tags/flai/v1.2.3":
			list := []string{`{"url":"` + assets + `0","name":"checksums.txt"}`}
			for i, name := range names {
				list = append(list, fmt.Sprintf(`{"url":"%s%d","name":"%s"}`, assets, i+1, name))
			}
			fmt.Fprintf(w, `{"url":"%s/repos/o/r/releases/1","tag_name":"flai/v1.2.3","assets":[%s]}`, base, strings.Join(list, ","))
		case "/repos/o/r/releases/assets/0":
			fmt.Fprint(w, strings.Join(sums, ""))
		case "/repos/o/r/releases/assets/1", "/repos/o/r/releases/assets/2", "/repos/o/r/releases/assets/3", "/repos/o/r/releases/assets/4":
			_, _ = w.Write(archive)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(queries)
	}
}

// cutOff answers body as a dropped connection does: the headers promise all
// of it, half of it is sent, and the connection is closed.
func cutOff(t *testing.T, w http.ResponseWriter, body string) {
	conn, rw, err := w.(http.Hijacker).Hijack()
	if err != nil {
		t.Errorf("hijack: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()
	fmt.Fprintf(rw, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", len(body), body[:len(body)/2])
	_ = rw.Flush()
}

// runInstallSh runs the repository's install.sh against the API at api,
// installing into dir, and returns its output and exit code.
func runInstallSh(t *testing.T, api, dir string) (string, int) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("install.sh does not install on Windows")
	}
	for _, tool := range []string{"sh", "curl", "tar"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("install.sh needs %s", tool)
		}
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "../../install.sh")
	// a token of its own, so that install.sh does not ask gh for one
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH"), "FLAI_API=" + api, "FLAI_REPO=o/r", "FLAI_INSTALL_DIR=" + dir, "GITHUB_TOKEN=stand-in"}
	out, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState == nil {
		t.Fatalf("run install.sh: %v", err)
	}
	return string(out), cmd.ProcessState.ExitCode()
}
