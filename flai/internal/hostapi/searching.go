package hostapi

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/search"
)

var (
	anyHeading = regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`)
	topHeading = regexp.MustCompile(`(?m)^#\s+(.+)$`)
	workItemID = regexp.MustCompile(`^[EST]-\d+$`)
)

// indexed is a project's search index and the state of the files it was
// built from. The index is rebuilt when that state changes: a walk that only
// stats is cheap, and it means search needs no word from the watcher.
type indexed struct {
	signature uint64
	index     *search.Index
}

var (
	indexes   = map[string]*indexed{}
	indexesMu sync.Mutex
)

// markdownUnder walks the manifest's folders in the order search scopes them.
func markdownUnder(root string, visit func(scope, rel string, info fs.FileInfo)) error {
	dirs, err := layoutDirs(root)
	if err != nil {
		return err
	}
	for _, s := range []struct{ scope, dir string }{{"design", dirs[0]}, {"wip", dirs[2]}, {"docs", dirs[1]}} {
		base := filepath.Join(root, filepath.FromSlash(s.dir))
		_ = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil //nolint:nilerr // a folder that cannot be read is skipped, not fatal
			}
			if strings.HasPrefix(d.Name(), ".") && path != base {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			info, err := d.Info()
			if err != nil || !info.Mode().IsRegular() {
				return nil //nolint:nilerr // gone between the listing and the stat
			}
			rel, err := filepath.Rel(root, path)
			if err == nil {
				visit(s.scope, filepath.ToSlash(rel), info)
			}
			return nil
		})
	}
	return nil
}

func searchIndex(root string) (*search.Index, error) {
	h := fnv.New64a()
	var files []struct{ scope, rel string }
	err := markdownUnder(root, func(scope, rel string, info fs.FileInfo) {
		files = append(files, struct{ scope, rel string }{scope, rel})
		_, _ = h.Write([]byte(rel))
		var b [16]byte
		m, s := uint64(info.ModTime().UnixNano()), uint64(info.Size())
		for i := 0; i < 8; i++ {
			b[i], b[8+i] = byte(m>>(8*i)), byte(s>>(8*i))
		}
		_, _ = h.Write(b[:])
	})
	if err != nil {
		return nil, err
	}
	sig := h.Sum64()
	indexesMu.Lock()
	defer indexesMu.Unlock()
	if have := indexes[root]; have != nil && have.signature == sig {
		return have.index, nil
	}
	docs := make([]search.Doc, 0, len(files))
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(f.rel)))
		if err != nil {
			continue
		}
		fm, body := frontMatter(string(data))
		id, _ := fm["id"].(string)
		typ, _ := fm["type"].(string)
		doc := search.Doc{Path: f.rel, Kind: "doc", Body: body, Scope: f.scope, Type: typ}
		if workItemID.MatchString(id) && typ != "" {
			doc.Kind, doc.ItemID = "item", id
		}
		doc.Nature, _ = fm["nature"].(string)
		doc.Status, _ = fm["status"].(string)
		if title, ok := fm["title"].(string); ok {
			doc.Title = title
		} else if m := topHeading.FindStringSubmatch(body); m != nil {
			doc.Title = m[1]
		} else {
			doc.Title = filepath.Base(f.rel)
		}
		doc.Tags = strings.Join(strs(fm["tags"]), " ")
		var heads []string
		for _, m := range anyHeading.FindAllStringSubmatch(body, -1) {
			heads = append(heads, m[1])
		}
		doc.Headings = strings.Join(heads, " ")
		docs = append(docs, doc)
	}
	ix := search.Build(docs)
	indexes[root] = &indexed{signature: sig, index: ix}
	return ix, nil
}

// SearchResult is what search.query answers.
type SearchResult struct {
	Query   string       `json:"query"`
	Indexed int          `json:"indexed"`
	Hits    []search.Hit `json:"hits"`
}

func searchMethods() map[string]channel.Method {
	return map[string]channel.Method{
		// search.query: full text over design and wip, docs on request.
		"search.query": func(_ context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				Q     string `json:"q"`
				Docs  bool   `json:"docs"`
				Limit int    `json:"limit"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			if len(in.Q) > 500 {
				return nil, bad("a query is at most 500 characters")
			}
			if in.Limit < 0 || in.Limit > 100 {
				return nil, bad("limit is between 1 and 100")
			}
			ix, err := searchIndex(p.Root)
			if err != nil {
				return nil, failed(err)
			}
			q := strings.TrimSpace(in.Q)
			return SearchResult{Query: q, Indexed: ix.Size(), Hits: ix.Search(q, in.Docs, in.Limit)}, nil
		},
	}
}
