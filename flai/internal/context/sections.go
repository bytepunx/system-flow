package context

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/search"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// LoadSections reads every markdown file under dirs, cut as the pack cuts
// design (Doc.Cuts), for doc_search and flai doc search (ADR-0049). Folder
// READMEs, which are indexes, and the ADR template are left out, and so are
// hidden folders. A dir that does not exist has no sections. root names
// paths.
func LoadSections(root string, dirs ...string) ([]search.Section, error) {
	var out []search.Section
	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) && p == dir {
					return filepath.SkipDir
				}
				return err
			}
			name := e.Name()
			if e.IsDir() {
				if p != dir && strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(name, ".md") || name == "README.md" || (filepath.Base(filepath.Dir(p)) == "adrs" && strings.HasPrefix(name, "0000-")) {
				return nil
			}
			raw, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return err
			}
			d, err := ParseDoc(filepath.ToSlash(rel), "", string(raw))
			if err != nil {
				return fmt.Errorf("%s: %w", rel, err)
			}
			out = append(out, d.Cuts()...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// SearchSections ranks the sections of the project's design and docs
// folders, the conventions among them, against a query: at most limit, and
// never more than search.MaxSectionHits.
func SearchSections(repo *workitem.Repo, query string, limit int) ([]search.SectionHit, error) {
	if len(search.Terms(query)) == 0 {
		return nil, fmt.Errorf("search for what? %q has no words to rank by (single letters and words such as \"the\" are left out)", query)
	}
	var dirs []string
	for _, key := range []string{"design", "docs"} {
		if repo.Manifest.Layout[key] != "" {
			dirs = append(dirs, repo.Manifest.Dir(repo.Root, key))
		}
	}
	secs, err := LoadSections(repo.Root, dirs...)
	if err != nil {
		return nil, fmt.Errorf("reading the design and docs folders: %w", err)
	}
	return search.IndexSections(secs).Search(query, limit), nil
}
