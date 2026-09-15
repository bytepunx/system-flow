package template

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Source is a resolved template location.
type Source struct {
	Repo  string // as configured: a URL or a local path
	Ref   string // branch, tag, or commit; empty means the default branch
	Local bool   // true when Repo is a directory on disk
	Dir   string // where the template files are (Repo itself, or the cache clone)
}

// IsLocal reports whether repo names a directory rather than a git remote.
func IsLocal(repo string) bool {
	if strings.Contains(repo, "://") || strings.HasPrefix(repo, "git@") {
		return false
	}
	st, err := os.Stat(config.ExpandHome(repo))
	return err == nil && st.IsDir()
}

// Resolve turns a repo and ref into a Source. Git sources map to a cache
// directory but are not cloned; call Ensure for that.
func Resolve(repo, ref, cacheDir string) (Source, error) {
	if repo == "" {
		return Source{}, fmt.Errorf("no template source configured; use --template or `flai template use`")
	}
	if IsLocal(repo) {
		dir, err := filepath.Abs(config.ExpandHome(repo))
		if err != nil {
			return Source{}, err
		}
		return Source{Repo: repo, Ref: "", Local: true, Dir: dir}, nil
	}
	sum := sha256.Sum256([]byte(repo + "@" + ref))
	dir := filepath.Join(config.ExpandHome(cacheDir), "templates", hex.EncodeToString(sum[:])[:12])
	return Source{Repo: repo, Ref: ref, Local: false, Dir: dir}, nil
}

// Cached reports whether a git source has already been cloned.
func (s Source) Cached() bool {
	if s.Local {
		return true
	}
	_, err := os.Stat(filepath.Join(s.Dir, ManifestFile))
	return err == nil
}

// Ensure makes s.Dir hold the template, cloning if needed. With refresh it
// discards any existing clone first.
func (s Source) Ensure(r execx.Runner, refresh bool) error {
	if s.Local {
		if _, err := os.Stat(filepath.Join(s.Dir, ManifestFile)); err != nil {
			return fmt.Errorf("%s is not a template: no %s", s.Repo, ManifestFile)
		}
		return nil
	}
	if s.Cached() && !refresh {
		return nil
	}
	if err := execx.Require(r, "git", "Install git to fetch templates, or point --template at a local directory."); err != nil {
		return err
	}
	if err := os.RemoveAll(s.Dir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.Dir), 0o755); err != nil {
		return err
	}
	args := []string{"clone", "--quiet", "--depth", "1"}
	if s.Ref != "" {
		args = append(args, "--branch", s.Ref)
	}
	args = append(args, s.Repo, s.Dir)
	if _, err := r.Run("", "git", args...); err != nil {
		if s.Ref == "" {
			return fmt.Errorf("clone template: %w", err)
		}
		// --branch only accepts branches and tags. Fall back to a full clone
		// and a checkout so commit hashes work too.
		_ = os.RemoveAll(s.Dir)
		if _, err2 := r.Run("", "git", "clone", "--quiet", s.Repo, s.Dir); err2 != nil {
			return fmt.Errorf("clone template: %w", err2)
		}
		if _, err2 := r.Run(s.Dir, "git", "checkout", "--quiet", s.Ref); err2 != nil {
			_ = os.RemoveAll(s.Dir)
			return fmt.Errorf("ref %q not found in %s: %w", s.Ref, s.Repo, err2)
		}
	}
	if _, err := os.Stat(filepath.Join(s.Dir, ManifestFile)); err != nil {
		_ = os.RemoveAll(s.Dir)
		return fmt.Errorf("%s is not a template: no %s at its root", s.Repo, ManifestFile)
	}
	return nil
}
