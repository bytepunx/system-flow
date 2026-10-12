package dashboard

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// cacheFolder is the cache's folder under flai serve's state directory.
const cacheFolder = "dashboard-digests"

// listExt ends a cached list's file name; its signature's adds ".sig".
const listExt = ".digests"

// Cache keeps each verified digest list with its signature, by release
// version, and verifies each again whenever it is read.
type Cache struct {
	Dir  string             // the folder dashboard-digests under flai serve's state directory
	Keys []*ecdsa.PublicKey // the keys a cached list must verify against
}

// NewCache is the cache under serveDir, flai serve's state directory, trusting keys.
func NewCache(serveDir string, keys []*ecdsa.PublicKey) Cache {
	return Cache{Dir: filepath.Join(serveDir, cacheFolder), Keys: keys}
}

// Put verifies version's list and its signature and keeps both, refusing a
// list that does not verify or is for another version.
func (c Cache) Put(version string, data, sig []byte) error {
	l, err := Verify(data, sig, c.Keys)
	if err != nil {
		return fmt.Errorf("caching the flaiover %s digest list refused: %w", version, err)
	}
	if l.Version != version {
		return fmt.Errorf("caching the digest list as flaiover %s refused: it lists version %s; fetch the list from the flaiover v%s release", version, l.Version, version)
	}
	if err := os.MkdirAll(c.Dir, 0o700); err != nil {
		return fmt.Errorf("caching the flaiover %s digest list failed: %w; check that %s can be created and written", version, err, c.Dir)
	}
	listPath, sigPath := c.paths(version)
	for _, f := range []struct {
		path string
		data []byte
	}{{listPath, data}, {sigPath, sig}} {
		if err := atomicfile.WriteFile(f.path, f.data, 0o600); err != nil {
			return fmt.Errorf("caching the flaiover %s digest list failed writing %s: %w; check that %s is writable", version, f.path, err, c.Dir)
		}
	}
	return nil
}

// Get is version's cached list, verified again; false when none is cached,
// an error naming the file when the one cached does not verify.
func (c Cache) Get(version string) (List, bool, error) {
	if !buildinfo.Bare(version) {
		return List{}, false, fmt.Errorf("reading the cached flaiover digest list for %q refused: it is not a release version X.Y.Z without its v", version)
	}
	listPath, sigPath := c.paths(version)
	data, err := os.ReadFile(listPath)
	if errors.Is(err, fs.ErrNotExist) {
		return List{}, false, nil
	}
	if err != nil {
		return List{}, false, unreadable(listPath, err)
	}
	sig, err := os.ReadFile(sigPath)
	if err != nil {
		return List{}, false, unreadable(listPath, err)
	}
	l, err := Verify(data, sig, c.Keys)
	if err != nil {
		return List{}, false, unreadable(listPath, err)
	}
	if l.Version != version {
		return List{}, false, unreadable(listPath, fmt.Errorf("it lists version %s", l.Version))
	}
	return l, true, nil
}

// All is every cached list that verifies, the newest version first, and an
// error naming each cached list skipped because it does not.
func (c Cache) All() ([]List, []error) {
	entries, err := os.ReadDir(c.Dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, []error{fmt.Errorf("reading the cached flaiover digest lists in %s failed: %w; check that it is readable", c.Dir, err)}
	}
	var lists []List
	var skipped []error
	for _, e := range entries {
		name := e.Name()
		version, ok := strings.CutSuffix(name, listExt)
		if !ok || strings.HasPrefix(name, ".") || e.IsDir() {
			continue
		}
		l, found, err := c.Get(version)
		switch {
		case err != nil:
			skipped = append(skipped, err)
		case found:
			lists = append(lists, l)
		}
	}
	sort.Slice(lists, func(i, j int) bool { return buildinfo.Below(lists[j].Version, lists[i].Version) })
	return lists, skipped
}

// Lookup is the newest cached list that verifies and names digest; when
// none does, the error names each cached list skipped because it does not
// verify, nil when none was.
func (c Cache) Lookup(digest string) (List, bool, error) {
	lists, skipped := c.All()
	for _, l := range lists {
		if l.Has(digest) {
			return l, true, nil
		}
	}
	return List{}, false, errors.Join(skipped...)
}

// paths are version's cached list and its signature.
func (c Cache) paths(version string) (list, sig string) {
	list = filepath.Join(c.Dir, version+listExt)
	return list, list + ".sig"
}

// unreadable refuses the cached list at path, which is not trusted.
func unreadable(path string, err error) error {
	return fmt.Errorf("reading the cached flaiover digest list %s failed: %w; it is not trusted, so delete it and its .sig for the list to be fetched again from its release", path, err)
}
