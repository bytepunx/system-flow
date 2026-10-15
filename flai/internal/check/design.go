package check

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// designTokens validates the design token file the manifest names
// (ADR-0148): each problem reading and resolving it gives is an error at its
// token path and line. A project with no token file passes; contrast is
// reported by flai design show, never here. It then checks the generated
// stylesheet the manifest names, if any, against the tokens (designStylesheet,
// designStale).
func (c *checker) designTokens() {
	inProject := c.designStylesheet()
	name := c.repo.Manifest.DesignTokens()
	file := filepath.Join(c.repo.Root, filepath.FromSlash(name))
	sys, err := designsystem.Load(os.DirFS(c.repo.Root), name)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return // optional until flai design init writes one
	case err != nil:
		c.add(Error, "design.tokens", file, 1, "%v", err)
		return
	}
	for _, p := range sys.Problems {
		if p.Path == "" {
			c.add(Error, "design.tokens", file, p.Line, "%s", p.Message)
			continue
		}
		c.add(Error, "design.tokens", file, p.Line, "%s: %s", p.Path, p.Message)
	}
	if inProject && len(sys.Problems) == 0 {
		c.designStale(sys, name)
	}
}

// designStylesheet reports design.stylesheet when the manifest names a
// stylesheet that lies outside every sub-project (ADR-0152): an error on the
// manifest at the key's line. It is true when the manifest names a
// stylesheet under a sub-project, the one case designStale reads it in.
func (c *checker) designStylesheet() bool {
	m := c.repo.Manifest
	sheet := m.DesignStylesheet()
	if sheet == "" {
		return false
	}
	if _, ok := m.DesignStylesheetProject(); ok {
		return true
	}
	mf := filepath.Join(c.repo.Root, manifest.File)
	c.add(Error, "design.stylesheet", mf, designKeyLine(mf, "stylesheet"), "design.stylesheet %q lies outside every sub-project; name a path under the path of one of the manifest's projects", sheet)
	return false
}

// designStale reports design.stale when the stylesheet the manifest names is
// missing or differs from what flai design css --write writes for sys, read
// from the token file tokens (ADR-0152, ADR-0159): a warning on the
// stylesheet and one on the token file, so that a branch that changed either
// stops at its close-out. A stylesheet flai cannot read is an error on it.
func (c *checker) designStale(sys *designsystem.System, tokens string) {
	sheet := path.Clean(c.repo.Manifest.DesignStylesheet())
	sheetPath := filepath.Join(c.repo.Root, filepath.FromSlash(sheet))
	tokensPath := filepath.Join(c.repo.Root, filepath.FromSlash(tokens))
	content, err := fs.ReadFile(os.DirFS(c.repo.Root), sheet)
	var msg string
	switch {
	case errors.Is(err, fs.ErrNotExist):
		msg = "the stylesheet " + sheet + " for the tokens in " + tokens + " is missing; run flai design css --write to write it"
	case err != nil:
		c.add(Error, "design.stale", sheetPath, 1, "%v", err)
		return
	case designsystem.Stale(sys, tokens, content):
		msg = "the stylesheet " + sheet + " does not match the tokens in " + tokens + "; run flai design css --write to bring it in step"
	default:
		return
	}
	c.add(Warning, "design.stale", sheetPath, 1, "%s", msg)
	c.add(Warning, "design.stale", tokensPath, 1, "%s", msg)
}

// designKeyLine is the line of key in the manifest's design block, else the
// line of design itself, else 1.
func designKeyLine(mf, key string) int {
	at := keyLine(mf, "design")
	f, err := os.Open(mf)
	if err != nil {
		return at
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	n, in := 0, false
	for sc.Scan() {
		n++
		l := sc.Text()
		switch {
		case strings.HasPrefix(l, "design:"):
			in = true
		case in && strings.TrimSpace(l) != "" && !strings.HasPrefix(l, " ") && !strings.HasPrefix(l, "\t"):
			in = false
		case in && strings.HasPrefix(strings.TrimSpace(l), key+":"):
			return n
		}
	}
	return at
}
