// Package conventions loads design/conventions: the agent norms every
// session primes with. See design/system/conventions.md and ADR-0013.
package conventions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the fixed subfolder name under layout.design.
const Folder = "conventions"

// Marker separates the template baseline from project additions.
const Marker = "<!-- system-flow:end-of-baseline -->"

// MaxLines is the length limit from conventions.md.
const MaxLines = 120

// File is one convention file.
type File struct {
	Path     string `json:"path"` // relative to the repo root
	Name     string `json:"name"` // file name
	Title    string `json:"title"`
	Updated  string `json:"updated"`
	Audience string `json:"audience"`
	Order    int    `json:"order"`
	Status   string `json:"status"`
	Lines    int    `json:"lines"`
	Markers  int    `json:"-"`
	HasAdds  bool   `json:"-"` // "## Project additions" after the marker
	Body     string `json:"-"`
	Raw      string `json:"-"`
}

// Set is the loaded folder.
type Set struct {
	Dir     string // absolute folder path
	Files   []File // sorted by Order, then name
	README  string // README.md content, empty if missing
	Missing bool   // folder does not exist
}

type front struct {
	Title    string `yaml:"title"`
	Updated  string `yaml:"updated"`
	Audience string `yaml:"audience"`
	Order    *int   `yaml:"order"`
	Status   string `yaml:"status"`
}

// Dir returns <layout.design>/conventions for the repo.
func Dir(r *workitem.Repo) string {
	return filepath.Join(r.Manifest.Dir(r.Root, "design"), Folder)
}

// Load reads every convention file. Parse problems are returned inside
// File.Title as empty with the error collected in errs so check can report
// per file; prime ignores errs.
func Load(r *workitem.Repo) (*Set, map[string]error, error) {
	dir := Dir(r)
	set := &Set{Dir: dir}
	errs := map[string]error{}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		set.Missing = true
		return set, errs, nil
	}
	if err != nil {
		return nil, nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		rel, _ := filepath.Rel(r.Root, path)
		if e.Name() == "README.md" {
			set.README = string(data)
			continue
		}
		f := File{Path: rel, Name: e.Name(), Raw: string(data), Lines: strings.Count(string(data), "\n")}
		fm, body, err := workitem.SplitFrontMatter(string(data))
		if err != nil {
			errs[rel] = err
		} else {
			var fr front
			if err := yaml.Unmarshal([]byte(fm), &fr); err != nil {
				errs[rel] = err
			} else {
				f.Title, f.Updated, f.Audience, f.Status = fr.Title, fr.Updated, fr.Audience, fr.Status
				if fr.Order != nil {
					f.Order = *fr.Order
				} else {
					f.Order = -1
				}
			}
			f.Body = body
		}
		f.Markers = strings.Count(string(data), Marker)
		if i := strings.Index(string(data), Marker); i >= 0 {
			f.HasAdds = strings.Contains(string(data)[i:], "## Project additions")
		}
		set.Files = append(set.Files, f)
	}
	sort.Slice(set.Files, func(i, j int) bool {
		if set.Files[i].Order != set.Files[j].Order {
			return set.Files[i].Order < set.Files[j].Order
		}
		return set.Files[i].Name < set.Files[j].Name
	})
	return set, errs, nil
}

var linkPattern = regexp.MustCompile(`\]\(([^)#]+\.md)\)`)

// IndexedNames returns the file names README.md links to, with counts.
func (s *Set) IndexedNames() map[string]int {
	out := map[string]int{}
	for _, m := range linkPattern.FindAllStringSubmatch(s.README, -1) {
		out[filepath.Base(m[1])]++
	}
	return out
}

// Finding is one rule breach: level, rule name, repo-relative path, message.
type Finding struct {
	Level, Rule, Path, Message string
}

// Validate checks the set against conventions.md.
func (s *Set) Validate(errs map[string]error) []Finding {
	var out []Finding
	add := func(level, rule, path, msg string, args ...any) {
		out = append(out, Finding{level, rule, path, fmt.Sprintf(msg, args...)})
	}
	dirRel := filepath.Base(filepath.Dir(s.Dir)) + "/" + Folder
	if s.Missing {
		add("warning", "conventions.missing", dirRel, "no %s folder; render it from the template or run flai upgrade", dirRel)
		return out
	}
	if s.README == "" {
		add("error", "conventions.index", dirRel+"/README.md", "README.md is missing; it lists every convention file in read order")
	}
	indexed := s.IndexedNames()
	orders := map[int]string{}
	for _, f := range s.Files {
		if err, ok := errs[f.Path]; ok {
			add("error", "conventions.front-matter", f.Path, "%v", err)
			continue
		}
		if f.Title == "" || f.Updated == "" || f.Status == "" {
			add("error", "conventions.front-matter", f.Path, "front matter needs title, updated, audience, order, status")
		}
		if f.Audience != "agent" {
			add("error", "conventions.front-matter", f.Path, "audience must be \"agent\", got %q", f.Audience)
		}
		if f.Order < 0 {
			add("error", "conventions.front-matter", f.Path, "order is required and governs read order")
		} else if prev, dup := orders[f.Order]; dup {
			add("error", "conventions.order", f.Path, "order %d is also used by %s", f.Order, prev)
		} else {
			orders[f.Order] = f.Name
		}
		switch f.Markers {
		case 0:
			add("error", "conventions.marker", f.Path, "missing the baseline marker %s", Marker)
		case 1:
			if !f.HasAdds {
				add("warning", "conventions.marker", f.Path, "no \"## Project additions\" section after the marker")
			}
		default:
			add("error", "conventions.marker", f.Path, "marker appears %d times; exactly one", f.Markers)
		}
		if f.Lines > MaxLines {
			add("warning", "conventions.length", f.Path, "%d lines; conventions stay under %d so they are read at every session start", f.Lines, MaxLines)
		}
		if s.README != "" {
			switch indexed[f.Name] {
			case 0:
				add("error", "conventions.index", dirRel+"/README.md", "README.md does not list %s", f.Name)
			case 1:
			default:
				add("warning", "conventions.index", dirRel+"/README.md", "README.md lists %s %d times", f.Name, indexed[f.Name])
			}
		}
	}
	names := map[string]bool{}
	for _, f := range s.Files {
		names[f.Name] = true
	}
	for name := range indexed {
		if !names[name] && name != "README.md" {
			add("error", "conventions.index", dirRel+"/README.md", "README.md links to %s which does not exist", name)
		}
	}
	return out
}
