package selfupgrade

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
	"testing"
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

// fakeGitHub serves two flai releases and one flaiover tag, with the newest
// flai release first as the API does.
func fakeGitHub(t *testing.T, archive, sums []byte, wantToken string) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantToken != "" && r.Header.Get("Authorization") != "Bearer "+wantToken {
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		assets := func(v string) string {
			return fmt.Sprintf(`"assets":[{"name":"checksums.txt","url":"%s/assets/sums"},{"name":"flai_%s_linux_amd64.tar.gz","url":"%s/assets/archive"},{"name":"flai_%s_windows_amd64.zip","url":"%s/assets/zip"}]`, srv.URL, v, srv.URL, v, srv.URL)
		}
		switch r.URL.Path {
		case "/repos/o/r/releases":
			fmt.Fprintf(w, `[{"tag_name":"flaiover/v9.0.0","assets":[]},{"tag_name":"flai/v1.2.0","draft":true,%s},{"tag_name":"flai/v1.1.0",%s},{"tag_name":"flai/v1.0.0",%s}]`, assets("1.2.0"), assets("1.1.0"), assets("1.0.0"))
		case "/repos/o/r/releases/tags/flai%2Fv1.0.0", "/repos/o/r/releases/tags/flai/v1.0.0":
			fmt.Fprintf(w, `{"tag_name":"flai/v1.0.0",%s}`, assets("1.0.0"))
		case "/assets/archive":
			if r.Header.Get("Accept") != "application/octet-stream" {
				http.Error(w, "wrong accept", http.StatusBadRequest)
				return
			}
			_, _ = w.Write(archive)
		case "/assets/sums":
			_, _ = w.Write(sums)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestResolveDownloadReplace(t *testing.T) {
	content := []byte("#!/bin/sh\necho new\n")
	archive := tarball(t, "flai", content)
	sum := sha256.Sum256(archive)
	sums := []byte(hex.EncodeToString(sum[:]) + "  flai_1.1.0_linux_amd64.tar.gz\n" + "deadbeef  flai_1.0.0_linux_amd64.tar.gz\n")
	srv := fakeGitHub(t, archive, sums, "tok")
	opt := Options{Repo: "o/r", APIBase: srv.URL, Token: "tok", OS: "linux", Arch: "amd64"}

	rel, err := Resolve(context.Background(), opt)
	if err != nil || rel.Version != "1.1.0" || rel.Tag != "flai/v1.1.0" {
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
