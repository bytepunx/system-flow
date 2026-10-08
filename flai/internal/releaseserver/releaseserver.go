// Package releaseserver answers the GitHub REST endpoints that install.sh and
// flai self-upgrade read, from GoReleaser dist folders on this host: the
// release listing, a release by tag, and an asset download. The smoke tier
// points FLAI_API and FLAI_RELEASES_API at it, so that the installer and
// self-upgrade run unchanged against a release built from the tree and no
// close-out depends on GitHub (S-0340).
package releaseserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// TagPrefix is the tag prefix of a flai release; a release's tag is it
// followed by the archives' version.
const TagPrefix = "flai/v"

// ChecksumsName is the checksums file GoReleaser writes into a dist folder.
const ChecksumsName = "checksums.txt"

// archiveName matches a GoReleaser archive of flai: flai_<version>_<os>_<arch>
// with .tar.gz or .zip.
var archiveName = regexp.MustCompile(`^flai_([^_]+)_[^_]+_[^_]+\.(tar\.gz|zip)$`)

// Asset is one file of a release, read from the folder that holds it.
type Asset struct {
	Name string
	FS   fs.FS
}

// Release is one release the server answers: its tag, its publish time, and
// its assets.
type Release struct {
	Tag       string
	Published time.Time
	Assets    []Asset
}

// Scan reads GoReleaser dist folders into releases, newest version first: one
// release per plain X.Y.Z version among the folders' flai archives, holding
// those archives and the checksums.txt of their folder, published at the
// newest archive's modification time.
func Scan(dirs ...string) ([]Release, error) {
	if len(dirs) == 0 {
		return nil, errors.New("scan releases: no dist folder given; pass the folder scripts/flai-snapshot.sh --local X.Y.Z builds (flai/dist)")
	}
	byVersion := map[string]*Release{}
	from := map[string]string{}
	for _, dir := range dirs {
		found, err := scanDir(dir, os.DirFS(dir))
		if err != nil {
			return nil, err
		}
		for version, rel := range found {
			if other, ok := from[version]; ok {
				return nil, fmt.Errorf("scan releases: version %s has archives in both %s and %s; give each version from one folder only", version, other, dir)
			}
			from[version] = dir
			byVersion[version] = rel
		}
	}
	out := make([]Release, 0, len(byVersion))
	for _, rel := range byVersion {
		out = append(out, *rel)
	}
	slices.SortFunc(out, func(a, b Release) int {
		va, vb := strings.TrimPrefix(a.Tag, TagPrefix), strings.TrimPrefix(b.Tag, TagPrefix)
		switch {
		case buildinfo.Below(vb, va):
			return -1
		case buildinfo.Below(va, vb):
			return 1
		}
		return 0
	})
	return out, nil
}

// scanDir finds the releases in one dist folder, by version.
func scanDir(dir string, fsys fs.FS) (map[string]*Release, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("scan releases: read dist folder %s: %w", dir, err)
	}
	out := map[string]*Release{}
	hasSums := false
	for _, e := range entries {
		if !e.Type().IsRegular() {
			continue
		}
		name := e.Name()
		if name == ChecksumsName {
			hasSums = true
			continue
		}
		m := archiveName.FindStringSubmatch(name)
		if m == nil || !buildinfo.Bare(m[1]) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("scan releases: stat %s: %w", filepath.Join(dir, name), err)
		}
		rel, ok := out[m[1]]
		if !ok {
			rel = &Release{Tag: TagPrefix + m[1]}
			out[m[1]] = rel
		}
		rel.Assets = append(rel.Assets, Asset{Name: name, FS: fsys})
		if t := info.ModTime().UTC().Truncate(time.Second); t.After(rel.Published) {
			rel.Published = t
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("scan releases: %s holds no archive named flai_X.Y.Z_<os>_<arch>.tar.gz or .zip with a plain X.Y.Z; build one with scripts/flai-snapshot.sh --local X.Y.Z", dir)
	}
	if !hasSums {
		return nil, fmt.Errorf("scan releases: %s holds no %s beside its archives; give the dist folder GoReleaser wrote, which has one", dir, ChecksumsName)
	}
	for _, rel := range out {
		rel.Assets = append(rel.Assets, Asset{Name: ChecksumsName, FS: fsys})
		slices.SortFunc(rel.Assets, func(a, b Asset) int { return strings.Compare(a.Name, b.Name) })
	}
	return out, nil
}

// served is an asset with the id its URL carries.
type served struct {
	Asset
	id int
}

// Server answers GitHub's release endpoints for one repository.
type Server struct {
	repo     string
	releases []Release
	assets   [][]served // by release, in the order of releases
	byID     map[int]served
}

// New answers GitHub's release endpoints for repo, such as local/flai, with
// releases in the order given, newest first as Scan answers them.
func New(repo string, releases []Release) *Server {
	s := &Server{repo: repo, releases: releases, byID: map[int]served{}}
	id := 0
	for _, rel := range releases {
		var list []served
		for _, a := range rel.Assets {
			id++
			sa := served{Asset: a, id: id}
			list = append(list, sa)
			s.byID[id] = sa
		}
		s.assets = append(s.assets, list)
	}
	return s
}

// ghRelease is a release as GitHub's API answers it, with the fields
// install.sh and selfupgrade read.
type ghRelease struct {
	ID          int       `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []ghAsset `json:"assets"`
}

// ghAsset is a release asset as GitHub's API answers it. url comes before
// name: install.sh's asset_url pairs each url with the name after it.
type ghAsset struct {
	URL         string `json:"url"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// ServeHTTP answers GET and HEAD of the release listing, a release by tag,
// and an asset; everything else answers 404 as GitHub does.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		notFound(w)
		return
	}
	rest, ok := strings.CutPrefix(r.URL.EscapedPath(), "/repos/"+s.repo+"/releases")
	if !ok {
		notFound(w)
		return
	}
	switch {
	case rest == "":
		s.list(w, r)
	case strings.HasPrefix(rest, "/tags/"):
		tag, err := url.PathUnescape(strings.TrimPrefix(rest, "/tags/"))
		if err != nil {
			notFound(w)
			return
		}
		s.tag(w, r, tag)
	case strings.HasPrefix(rest, "/assets/"):
		id, err := strconv.Atoi(strings.TrimPrefix(rest, "/assets/"))
		if err != nil {
			notFound(w)
			return
		}
		s.asset(w, r, id)
	default:
		notFound(w)
	}
}

// list answers every release, whatever per_page asks, with no Link header:
// there is never a next page.
func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	out := make([]ghRelease, 0, len(s.releases))
	for i := range s.releases {
		rel, err := s.render(r, i)
		if err != nil {
			serverError(w, err)
			return
		}
		out = append(out, rel)
	}
	writeJSON(w, out)
}

// tag answers the release tagged tag.
func (s *Server) tag(w http.ResponseWriter, r *http.Request, tag string) {
	for i := range s.releases {
		if s.releases[i].Tag != tag {
			continue
		}
		rel, err := s.render(r, i)
		if err != nil {
			serverError(w, err)
			return
		}
		writeJSON(w, rel)
		return
	}
	notFound(w)
}

// asset answers an asset's bytes when the request accepts
// application/octet-stream, and its description otherwise, as GitHub does.
func (s *Server) asset(w http.ResponseWriter, r *http.Request, id int) {
	a, ok := s.byID[id]
	if !ok {
		notFound(w)
		return
	}
	if r.Header.Get("Accept") != "application/octet-stream" {
		desc, err := s.describe(r, a)
		if err != nil {
			serverError(w, err)
			return
		}
		writeJSON(w, desc)
		return
	}
	data, err := fs.ReadFile(a.FS, a.Name)
	if err != nil {
		serverError(w, fmt.Errorf("read asset %s: %w", a.Name, err))
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodGet {
		_, _ = w.Write(data)
	}
}

// render is release i as GitHub answers it, its asset URLs under the
// request's host.
func (s *Server) render(r *http.Request, i int) (ghRelease, error) {
	rel := s.releases[i]
	out := ghRelease{
		ID:          i + 1,
		TagName:     rel.Tag,
		Name:        "flai v" + strings.TrimPrefix(rel.Tag, TagPrefix),
		PublishedAt: rel.Published,
		Assets:      []ghAsset{},
	}
	for _, a := range s.assets[i] {
		desc, err := s.describe(r, a)
		if err != nil {
			return ghRelease{}, err
		}
		out.Assets = append(out.Assets, desc)
	}
	return out, nil
}

// describe is an asset as GitHub answers it.
func (s *Server) describe(r *http.Request, a served) (ghAsset, error) {
	info, err := fs.Stat(a.FS, a.Name)
	if err != nil {
		return ghAsset{}, fmt.Errorf("stat asset %s: %w", a.Name, err)
	}
	return ghAsset{
		URL:         fmt.Sprintf("%s/repos/%s/releases/assets/%d", baseURL(r), s.repo, a.id),
		ID:          a.id,
		Name:        a.Name,
		ContentType: "application/octet-stream",
		Size:        info.Size(),
	}, nil
}

// baseURL is the scheme and host the request reached the server at.
func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`{"message":"Not Found"}` + "\n"))
}

func serverError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
}
