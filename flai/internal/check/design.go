package check

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
)

// designTokens validates the design token file the manifest names
// (ADR-0148): each problem reading and resolving it gives is an error at its
// token path and line. A project with no token file passes; contrast is
// reported by flai design show, never here.
func (c *checker) designTokens() {
	name := c.repo.Manifest.DesignTokens()
	path := filepath.Join(c.repo.Root, filepath.FromSlash(name))
	sys, err := designsystem.Load(os.DirFS(c.repo.Root), name)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return // optional until flai design init writes one
	case err != nil:
		c.add(Error, "design.tokens", path, 1, "%v", err)
		return
	}
	for _, p := range sys.Problems {
		if p.Path == "" {
			c.add(Error, "design.tokens", path, p.Line, "%s", p.Message)
			continue
		}
		c.add(Error, "design.tokens", path, p.Line, "%s: %s", p.Path, p.Message)
	}
}
