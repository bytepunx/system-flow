package hostapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// DocNode is one folder or Markdown file of the document tree: a file's
// title, not its front matter, which doc.get reads (S-0162).
type DocNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"` // relative to the repository, forward slashes
	Kind     string     `json:"kind"` // dir or file
	Title    string     `json:"title,omitempty"`
	Children []*DocNode `json:"children,omitempty"`
}

// Doc is one Markdown file: its front matter, the body under it, and the file as it is.
type Doc struct {
	Path        string         `json:"path"`
	FrontMatter map[string]any `json:"frontMatter"`
	Body        string         `json:"body"`
	Raw         string         `json:"raw"`
}

// Adr is an ADR's front matter, for the list of decisions.
type Adr struct {
	Path         string   `json:"path"`
	File         string   `json:"file"`
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Status       string   `json:"status"`
	Date         string   `json:"date"`
	Supersedes   []string `json:"supersedes"`
	SupersededBy []string `json:"supersededBy"`
	Refines      []string `json:"refines"`
}

// layoutDirs are the manifest's three folders, in the order the dashboard
// shows them. Only what lies under them is a document to the dashboard.
func layoutDirs(root string) ([]string, error) {
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return nil, err
	}
	var out []string
	for _, k := range []string{"design", "docs", "wip"} {
		dir := filepath.ToSlash(filepath.Clean(m.Layout[k]))
		if dir == "" || dir == "." || filepath.IsAbs(dir) || strings.HasPrefix(dir, "..") {
			return nil, fmt.Errorf("system-flow.yaml: layout.%s must be a folder inside the repository", k)
		}
		out = append(out, dir)
	}
	return out, nil
}

// frontMatter parses a document's front matter into plain values: what the
// file says, with dates as the text written there.
func frontMatter(doc string) (map[string]any, string) {
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil || strings.TrimSpace(fm) == "" {
		if err != nil {
			return nil, doc
		}
		return map[string]any{}, body
	}
	var out map[string]any
	if err := yaml.Unmarshal([]byte(fm), &out); err != nil {
		return nil, doc
	}
	return plain(out).(map[string]any), body
}

// plain turns what the YAML decoder produced into JSON-friendly values.
func plain(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			x[k] = plain(e)
		}
		return x
	case map[any]any:
		out := map[string]any{}
		for k, e := range x {
			out[fmt.Sprint(k)] = plain(e)
		}
		return out
	case []any:
		for i, e := range x {
			x[i] = plain(e)
		}
		return x
	case time.Time:
		if x.Hour() == 0 && x.Minute() == 0 && x.Second() == 0 {
			return x.Format("2006-01-02")
		}
		return x.UTC().Format(time.RFC3339)
	case nil:
		return nil
	}
	return v
}

// docPath checks a path the dashboard sent: a Markdown file under one of the
// manifest's folders, by its real location, and nothing else.
func docPath(root, rel string) (string, *channel.Error) {
	clean := filepath.ToSlash(filepath.Clean(rel))
	if rel == "" || filepath.IsAbs(rel) || strings.HasPrefix(clean, "..") || strings.ContainsRune(rel, 0) {
		return "", bad("%q is not a path inside the repository", rel)
	}
	if !strings.HasSuffix(clean, ".md") {
		return "", bad("only Markdown files are served")
	}
	dirs, err := layoutDirs(root)
	if err != nil {
		return "", failed(err)
	}
	inside := false
	for _, d := range dirs {
		if strings.HasPrefix(clean, d+"/") {
			inside = true
		}
	}
	if !inside {
		return "", bad("%s is not under the project's design, docs, or wip folder", clean)
	}
	abs := filepath.Join(root, filepath.FromSlash(clean))
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", &channel.Error{Code: NotFound, Message: "not found: " + clean}
	}
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", failed(err)
	}
	if !strings.HasPrefix(real, realRoot+string(filepath.Separator)) {
		return "", bad("%s leads outside the repository", clean)
	}
	return abs, nil
}

// racyWindow is how long after a file's modification time its title is
// still not trusted, as internal/workitem's store trusts an item: a second
// write within the timestamp granularity could leave the same time and size.
const racyWindow = 2 * time.Second

// titles keeps each document's title by file, so that a long-running flai
// (serve) reads again only the files that changed since it last walked the
// tree (S-0162). A title is kept while the file is the same file with the
// modification time and size it had when it was read, and that read came
// later than racyWindow after it was modified.
type titles struct {
	mu    sync.Mutex
	files map[string]*keptTitle // by absolute path
	reads int                   // files read, for tests
}

type keptTitle struct {
	info    fs.FileInfo
	title   string
	trusted bool
}

var (
	docTitles   = map[string]*titles{} // by repository root
	docTitlesMu sync.Mutex
)

func titlesFor(root string) *titles {
	docTitlesMu.Lock()
	defer docTitlesMu.Unlock()
	t := docTitles[root]
	if t == nil {
		t = &titles{files: map[string]*keptTitle{}}
		docTitles[root] = t
	}
	return t
}

// title is the front matter's title of the file at abs, read again only when
// the file may have changed since it was kept. ok is false when it cannot be
// read. The caller holds t.mu.
func (t *titles) title(abs string, info fs.FileInfo, seen map[string]bool) (string, bool) {
	seen[abs] = true
	if k := t.files[abs]; k != nil && k.trusted && os.SameFile(k.info, info) &&
		k.info.ModTime().Equal(info.ModTime()) && k.info.Size() == info.Size() {
		return k.title, true
	}
	readAt := time.Now()
	data, err := os.ReadFile(abs)
	if err != nil {
		delete(t.files, abs)
		return "", false
	}
	t.reads++
	fm, _ := frontMatter(string(data))
	title, _ := fm["title"].(string)
	t.files[abs] = &keptTitle{info: info, title: title, trusted: info.ModTime().Before(readAt.Add(-racyWindow))}
	return title, true
}

// docsTree walks the manifest's three folders into their trees, each
// Markdown file with its title, reading only the files that changed since
// the last walk of root.
func docsTree(root string, dirs []string) []*DocNode {
	t := titlesFor(root)
	t.mu.Lock()
	defer t.mu.Unlock()
	seen := map[string]bool{}
	roots := []*DocNode{}
	for _, d := range dirs {
		roots = append(roots, t.tree(root, d, seen))
	}
	for abs := range t.files {
		if !seen[abs] {
			delete(t.files, abs)
		}
	}
	return roots
}

func (t *titles) tree(root, rel string, seen map[string]bool) *DocNode {
	node := &DocNode{Name: filepath.Base(rel), Path: rel, Kind: "dir", Children: []*DocNode{}}
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return node
	}
	for _, e := range entries { // os.ReadDir sorts by name
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		child := rel + "/" + e.Name()
		switch {
		case e.IsDir():
			node.Children = append(node.Children, t.tree(root, child, seen))
		case strings.HasSuffix(e.Name(), ".md") && e.Type().IsRegular():
			info, err := e.Info()
			if err != nil {
				continue // removed since the folder was read
			}
			title, ok := t.title(filepath.Join(root, filepath.FromSlash(child)), info, seen)
			if !ok {
				continue
			}
			node.Children = append(node.Children, &DocNode{Name: e.Name(), Path: child, Kind: "file", Title: title})
		}
	}
	return node
}

func strs(v any) []string {
	out := []string{}
	if list, ok := v.([]any); ok {
		for _, e := range list {
			out = append(out, fmt.Sprint(e))
		}
	}
	return out
}

func str(v any, fallback string) string {
	if v == nil {
		return fallback
	}
	return fmt.Sprint(v)
}

var adrFile = regexp.MustCompile(`^\d{4}-.*\.md$`)

func docMethods() map[string]channel.Method {
	return map[string]channel.Method{
		// docs.tree: the manifest's three folders, every Markdown file with its title.
		"docs.tree": func(ctx context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
			defer perf.Track(ctx, "docs.walk")()
			dirs, err := layoutDirs(p.Root)
			if err != nil {
				return nil, failed(err)
			}
			return docsTree(p.Root, dirs), nil
		},

		// doc.get: one Markdown file under those folders.
		"doc.get": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				Path string `json:"path"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			abs, e := docPath(p.Root, in.Path)
			if e != nil {
				return nil, e
			}
			done := perf.Track(ctx, "doc.read")
			data, err := os.ReadFile(abs)
			done()
			if err != nil {
				return nil, &channel.Error{Code: NotFound, Message: "not found: " + in.Path}
			}
			fm, body := frontMatter(string(data))
			return Doc{Path: filepath.ToSlash(filepath.Clean(in.Path)), FrontMatter: fm, Body: body, Raw: string(data)}, nil
		},

		// adrs.list: the decisions, from their front matter; the template (0000) is not one.
		"adrs.list": func(ctx context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
			defer perf.Track(ctx, "adrs.read")()
			dirs, err := layoutDirs(p.Root)
			if err != nil {
				return nil, failed(err)
			}
			dir := dirs[0] + "/adrs"
			entries, _ := os.ReadDir(filepath.Join(p.Root, filepath.FromSlash(dir)))
			out := []Adr{}
			for _, e := range entries {
				name := e.Name()
				if !adrFile.MatchString(name) || strings.HasPrefix(name, "0000-") {
					continue
				}
				data, err := os.ReadFile(filepath.Join(p.Root, filepath.FromSlash(dir), name))
				if err != nil {
					continue
				}
				fm, _ := frontMatter(string(data))
				out = append(out, Adr{Path: dir + "/" + name, File: name, ID: str(fm["id"], name[:4]), Title: str(fm["title"], name),
					Status: str(fm["status"], ""), Date: str(fm["date"], ""), Supersedes: strs(fm["supersedes"]),
					SupersededBy: strs(fm["superseded_by"]), Refines: strs(fm["refines"])})
			}
			sort.Slice(out, func(i, j int) bool { return out[i].File < out[j].File })
			return out, nil
		},
	}
}
