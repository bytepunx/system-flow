// Package selfupgrade lists the published flai releases and dashboard tags,
// and installs a flai release over the running binary: it resolves the newest
// published flai/v* GitHub release (or a pinned published one), downloads
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
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
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

// Release is a published release: a flai GitHub release with its assets, or
// a dashboard tag, which has neither assets nor a publish date.
type Release struct {
	Tag       string    `json:"tag"`
	Version   string    `json:"version"`
	Published time.Time `json:"published,omitzero"`
	Assets    []Asset   `json:"assets,omitempty"`
}

// Asset is one downloadable file of a release.
type Asset struct {
	Name string `json:"name"`
	URL  string `json:"url"` // API URL; served as octet-stream on request
}

type ghRelease struct {
	TagName     string    `json:"tag_name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []ghAsset `json:"assets"`
}

type ghRef struct {
	Ref string `json:"ref"`
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

// shownVersions is how many published versions a refused pin names.
const shownVersions = 10

// Resolve finds the newest published flai release, or the pinned version when
// it is published.
func Resolve(ctx context.Context, opt Options) (*Release, error) {
	opt.defaults()
	published, err := List(ctx, opt, TagPrefix)
	if err != nil {
		return nil, err
	}
	if len(published) == 0 {
		return nil, fmt.Errorf("no published flai release in %s: no release tagged %sX.Y.Z that is neither a draft nor a prerelease; check the repository the releases come from", opt.Repo, TagPrefix)
	}
	if opt.Version == "" {
		return &published[0], nil
	}
	v := strings.TrimPrefix(strings.TrimPrefix(opt.Version, TagPrefix), "v")
	for i := range published {
		if published[i].Version == v {
			return &published[i], nil
		}
	}
	return nil, fmt.Errorf("resolve flai %s: %s has no published release tagged %s%s (drafts and prereleases are not published); published are %s: choose one of them, or give no version for the newest", v, opt.Repo, TagPrefix, v, versionList(published))
}

// versionList names the newest published versions, and how many older ones
// it leaves out.
func versionList(published []Release) string {
	names := make([]string, 0, shownVersions)
	for i := 0; i < len(published) && i < shownVersions; i++ {
		names = append(names, published[i].Version)
	}
	out := strings.Join(names, ", ")
	if more := len(published) - len(names); more > 0 {
		out += fmt.Sprintf(" and %d older", more)
	}
	return out
}

// List answers the published GitHub releases of opt.Repo whose tag is prefix
// followed by a plain X.Y.Z, without drafts or prereleases, newest first.
func List(ctx context.Context, opt Options, prefix string) ([]Release, error) {
	opt.defaults()
	var out []Release
	for url := fmt.Sprintf("%s/repos/%s/releases?per_page=100", opt.APIBase, opt.Repo); url != ""; {
		var page []ghRelease
		next, err := opt.getPage(ctx, url, &page)
		if err != nil {
			return nil, fmt.Errorf("list the %s releases of %s: %w", prefix, opt.Repo, err)
		}
		for _, rel := range page {
			v, ok := plainVersion(rel.TagName, prefix)
			if !ok || rel.Draft || rel.Prerelease {
				continue
			}
			r := Release{Tag: rel.TagName, Version: v, Published: rel.PublishedAt}
			for _, a := range rel.Assets {
				r.Assets = append(r.Assets, Asset(a))
			}
			out = append(out, r)
		}
		url = next
	}
	newestFirst(out)
	return out, nil
}

// ListTags answers the git tags of opt.Repo that are prefix followed by a
// plain X.Y.Z, newest first, without a publish date: dashboard releases are
// tags with no GitHub release.
func ListTags(ctx context.Context, opt Options, prefix string) ([]Release, error) {
	opt.defaults()
	var out []Release
	for url := fmt.Sprintf("%s/repos/%s/git/matching-refs/tags/%s", opt.APIBase, opt.Repo, prefix); url != ""; {
		var page []ghRef
		next, err := opt.getPage(ctx, url, &page)
		if err != nil {
			return nil, fmt.Errorf("list the %s tags of %s: %w", prefix, opt.Repo, err)
		}
		for _, ref := range page {
			tag := strings.TrimPrefix(ref.Ref, "refs/tags/")
			if v, ok := plainVersion(tag, prefix); ok {
				out = append(out, Release{Tag: tag, Version: v})
			}
		}
		url = next
	}
	newestFirst(out)
	return out, nil
}

// plainVersion is the X.Y.Z after prefix in tag, refusing a suffix, a second
// v, or leading zeros.
func plainVersion(tag, prefix string) (string, bool) {
	rest, ok := strings.CutPrefix(tag, prefix)
	if !ok {
		return "", false
	}
	n, ok := buildinfo.Semver(rest)
	if !ok || fmt.Sprintf("%d.%d.%d", n[0], n[1], n[2]) != rest {
		return "", false
	}
	return rest, true
}

// newestFirst orders releases by semver, the newest first.
func newestFirst(rels []Release) {
	slices.SortStableFunc(rels, func(a, b Release) int {
		switch {
		case buildinfo.Below(b.Version, a.Version):
			return -1
		case buildinfo.Below(a.Version, b.Version):
			return 1
		}
		return 0
	})
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

// getPage decodes one page of a list into v and answers the URL of the next
// page from the Link header, "" on the last.
func (o Options) getPage(ctx context.Context, url string, v any) (string, error) {
	resp, err := o.request(ctx, url, "application/vnd.github+json")
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return "", fmt.Errorf("read %s: %w", url, err)
	}
	return nextLink(resp.Header.Get("Link")), nil
}

// nextLink is the rel="next" URL of a Link header, "" when there is none.
func nextLink(header string) string {
	for _, part := range strings.Split(header, ",") {
		target, params, ok := strings.Cut(part, ";")
		if !ok {
			continue
		}
		for _, p := range strings.Split(params, ";") {
			if strings.TrimSpace(p) == `rel="next"` {
				return strings.Trim(strings.TrimSpace(target), "<>")
			}
		}
	}
	return ""
}

func (o Options) getBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := o.request(ctx, url, "application/octet-stream")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return io.ReadAll(resp.Body)
}
