// Package selfupgrade installs a flai release over the running binary: it
// resolves the newest flai/v* GitHub release (or a pinned version), downloads
// the archive for the platform and checksums.txt through the release asset
// API (which works for private repositories with a token), verifies the
// SHA-256, and replaces the executable atomically. install.sh at the
// repository root does the same from a shell.
package selfupgrade

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// TagPrefix is the monorepo tag prefix for flai releases.
const TagPrefix = "flai/v"

// Options configure a resolution and install.
type Options struct {
	Repo    string       // owner/name, default bytepunx/system-flow
	APIBase string       // default https://api.github.com
	Token   string       // optional bearer token
	Version string       // "" for latest, else X.Y.Z or vX.Y.Z
	OS      string       // default runtime.GOOS
	Arch    string       // default runtime.GOARCH
	Client  *http.Client // default http.DefaultClient
}

// Release is a resolved flai release.
type Release struct {
	Tag     string  `json:"tag"`
	Version string  `json:"version"`
	Assets  []Asset `json:"assets"`
}

// Asset is one downloadable file of a release.
type Asset struct {
	Name string `json:"name"`
	URL  string `json:"url"` // API URL; served as octet-stream on request
}

type ghRelease struct {
	TagName    string    `json:"tag_name"`
	Draft      bool      `json:"draft"`
	Prerelease bool      `json:"prerelease"`
	Assets     []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (o *Options) defaults() {
	if o.Repo == "" {
		o.Repo = "bytepunx/system-flow"
	}
	if o.APIBase == "" {
		o.APIBase = "https://api.github.com"
	}
	if o.OS == "" {
		o.OS = runtime.GOOS
	}
	if o.Arch == "" {
		o.Arch = runtime.GOARCH
	}
	if o.Client == nil {
		o.Client = http.DefaultClient
	}
}

// ArchiveName is the GoReleaser archive for a version and platform.
func ArchiveName(version, goos, goarch string) string {
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	return fmt.Sprintf("flai_%s_%s_%s%s", version, goos, goarch, ext)
}

// Resolve finds the newest non-draft flai release, or the pinned version.
func Resolve(ctx context.Context, opt Options) (*Release, error) {
	opt.defaults()
	if opt.Version != "" {
		v := strings.TrimPrefix(strings.TrimPrefix(opt.Version, TagPrefix), "v")
		tag := TagPrefix + v
		var rel ghRelease
		if err := opt.getJSON(ctx, fmt.Sprintf("%s/repos/%s/releases/tags/%s", opt.APIBase, opt.Repo, strings.ReplaceAll(tag, "/", "%2F")), &rel); err != nil {
			return nil, fmt.Errorf("release %s: %w", tag, err)
		}
		return toRelease(rel), nil
	}
	var list []ghRelease
	if err := opt.getJSON(ctx, fmt.Sprintf("%s/repos/%s/releases?per_page=50", opt.APIBase, opt.Repo), &list); err != nil {
		return nil, fmt.Errorf("list releases of %s: %w", opt.Repo, err)
	}
	for _, rel := range list {
		if rel.Draft || rel.Prerelease || !strings.HasPrefix(rel.TagName, TagPrefix) {
			continue
		}
		return toRelease(rel), nil
	}
	return nil, fmt.Errorf("no flai release found in %s", opt.Repo)
}

func toRelease(rel ghRelease) *Release {
	out := &Release{Tag: rel.TagName, Version: strings.TrimPrefix(rel.TagName, TagPrefix)}
	for _, a := range rel.Assets {
		out.Assets = append(out.Assets, Asset(a))
	}
	return out
}

// Asset returns the named asset of the release.
func (r *Release) Asset(name string) (Asset, bool) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a, true
		}
	}
	return Asset{}, false
}

// Download fetches the platform archive and checksums.txt, verifies the
// archive, and returns the extracted flai binary.
func Download(ctx context.Context, opt Options, rel *Release) ([]byte, error) {
	opt.defaults()
	name := ArchiveName(rel.Version, opt.OS, opt.Arch)
	archive, ok := rel.Asset(name)
	if !ok {
		return nil, fmt.Errorf("release %s has no asset %s", rel.Tag, name)
	}
	sums, ok := rel.Asset("checksums.txt")
	if !ok {
		return nil, fmt.Errorf("release %s has no checksums.txt", rel.Tag)
	}
	data, err := opt.getBytes(ctx, archive.URL)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", name, err)
	}
	sumData, err := opt.getBytes(ctx, sums.URL)
	if err != nil {
		return nil, fmt.Errorf("download checksums.txt: %w", err)
	}
	if err := Verify(data, name, sumData); err != nil {
		return nil, err
	}
	bin, err := extract(data, name, opt.OS)
	if err != nil {
		return nil, fmt.Errorf("extract %s: %w", name, err)
	}
	return bin, nil
}

// Verify checks data against the checksums.txt entry for name.
func Verify(data []byte, name string, checksums []byte) error {
	sum := sha256.Sum256(data)
	actual := hex.EncodeToString(sum[:])
	for _, line := range strings.Split(string(checksums), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == name {
			if fields[0] != actual {
				return fmt.Errorf("checksum mismatch for %s", name)
			}
			return nil
		}
	}
	return fmt.Errorf("no checksum entry for %s", name)
}

func extract(data []byte, name, goos string) ([]byte, error) {
	want := "flai"
	if goos == "windows" {
		want = "flai.exe"
	}
	if strings.HasSuffix(name, ".zip") {
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		for _, f := range zr.File {
			if filepath.Base(f.Name) == want {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer func() { _ = rc.Close() }()
				return io.ReadAll(rc)
			}
		}
		return nil, fmt.Errorf("archive does not contain %s", want)
	}
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("archive does not contain %s", want)
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag == tar.TypeReg && filepath.Base(h.Name) == want {
			return io.ReadAll(tr)
		}
	}
}

// Replace writes bin over dest atomically: a temp file beside dest, then a
// rename. On Windows the running executable is moved aside first.
func Replace(dest string, bin []byte) error {
	dir := filepath.Dir(dest)
	tmp, err := os.CreateTemp(dir, ".flai-upgrade-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }
	if _, err := tmp.Write(bin); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		cleanup()
		return err
	}
	if runtime.GOOS == "windows" {
		old := dest + ".old"
		_ = os.Remove(old)
		if err := os.Rename(dest, old); err != nil && !os.IsNotExist(err) {
			cleanup()
			return err
		}
	}
	if err := os.Rename(tmpName, dest); err != nil {
		cleanup()
		return err
	}
	return nil
}

func (o Options) request(ctx context.Context, url, accept string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", "flai-self-upgrade")
	if o.Token != "" {
		req.Header.Set("Authorization", "Bearer "+o.Token)
	}
	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		_ = resp.Body.Close()
		hint := ""
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusUnauthorized {
			hint = " (private repository? set GITHUB_TOKEN or run gh auth login)"
		}
		return nil, fmt.Errorf("%s: %s%s: %s", url, resp.Status, hint, strings.TrimSpace(string(body)))
	}
	return resp, nil
}

func (o Options) getJSON(ctx context.Context, url string, v any) error {
	resp, err := o.request(ctx, url, "application/vnd.github+json")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (o Options) getBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := o.request(ctx, url, "application/octet-stream")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return io.ReadAll(resp.Body)
}
