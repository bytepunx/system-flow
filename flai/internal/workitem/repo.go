package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

// Repo is a conforming project opened for work item operations.
type Repo struct {
	Root     string
	Manifest manifest.Manifest
	// TemplateDir, when set, points at a template checkout whose items/ are
	// used for new item bodies. Otherwise embedded defaults apply.
	TemplateDir string
	// MainRoot is the main checkout when Root is a linked git worktree
	// (a story branch under .flai-cache/worktrees, ADR-0019), else Root.
	// wip/ always lives in the main checkout so the board stays live.
	MainRoot string
}

// Open finds the manifest above start and returns a Repo.
func Open(start string) (*Repo, error) {
	path, err := manifest.Find(start)
	if err != nil {
		return nil, err
	}
	m, err := manifest.Load(path)
	if err != nil {
		return nil, err
	}
	root := filepath.Dir(path)
	return &Repo{Root: root, Manifest: m, MainRoot: mainCheckout(root)}, nil
}

// mainCheckout resolves a linked worktree's main checkout from its .git
// file ("gitdir: <main>/.git/worktrees/<name>"); a normal checkout returns
// itself.
func mainCheckout(root string) string {
	data, err := os.ReadFile(filepath.Join(root, ".git"))
	if err != nil {
		return root
	}
	line := strings.TrimSpace(string(data))
	gitdir, ok := strings.CutPrefix(line, "gitdir: ")
	if !ok {
		return root
	}
	gitdir = filepath.Clean(gitdir)
	marker := string(filepath.Separator) + ".git" + string(filepath.Separator) + "worktrees" + string(filepath.Separator)
	i := strings.Index(gitdir, marker)
	if i < 0 {
		return root
	}
	main := gitdir[:i]
	if !filepath.IsAbs(main) {
		main = filepath.Join(root, main)
	}
	if _, err := os.Stat(filepath.Join(main, "system-flow.yaml")); err != nil {
		return root
	}
	return main
}

// IsWorktree reports whether Root is a linked worktree of MainRoot.
func (r *Repo) IsWorktree() bool { return r.MainRoot != "" && r.MainRoot != r.Root }

// CacheDir is .flai-cache in the main checkout: token, worktrees, tool caches.
func (r *Repo) CacheDir() string { return filepath.Join(r.mainRoot(), ".flai-cache") }

// WorktreePath is where a story's branch is checked out.
func (r *Repo) WorktreePath(storyID string) string {
	return filepath.Join(r.CacheDir(), "worktrees", storyID)
}

func (r *Repo) mainRoot() string {
	if r.MainRoot != "" {
		return r.MainRoot
	}
	return r.Root
}

// WipDir is the work-in-process folder, always in the main checkout.
func (r *Repo) WipDir() string { return r.Manifest.Dir(r.mainRoot(), "wip") }

// KanbanDir is wip/kanban.
func (r *Repo) KanbanDir() string { return filepath.Join(r.WipDir(), "kanban") }

// AgentsDir is wip/agents.
func (r *Repo) AgentsDir() string { return filepath.Join(r.WipDir(), "agents") }

// ArchiveDir is wip/archive.
func (r *Repo) ArchiveDir() string { return filepath.Join(r.WipDir(), "archive") }

// Folder is the plural folder name for a type.
func Folder(typ string) string {
	if typ == Story {
		return "stories"
	}
	return typ + "s"
}

// ItemDir returns the kanban or archive folder for a type.
func (r *Repo) ItemDir(typ string, archived bool) string {
	if archived {
		return filepath.Join(r.ArchiveDir(), "kanban", Folder(typ))
	}
	return filepath.Join(r.KanbanDir(), Folder(typ))
}

// List loads every item, optionally including the archive. Items are sorted
// by ID.
func (r *Repo) List(includeArchive bool) ([]*Item, error) {
	var items []*Item
	scopes := []bool{false}
	if includeArchive {
		scopes = append(scopes, true)
	}
	for _, archived := range scopes {
		for _, typ := range Types {
			dir := r.ItemDir(typ, archived)
			entries, err := os.ReadDir(dir)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return nil, err
			}
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasPrefix(e.Name(), "_") {
					continue
				}
				it, err := ReadItem(filepath.Join(dir, e.Name()))
				if err != nil {
					return nil, err
				}
				it.Archived = archived
				items = append(items, it)
			}
		}
	}
	sort.Slice(items, func(i, j int) bool { return lessID(items[i].ID, items[j].ID) })
	return items, nil
}

// IDWidth is the number of digits flai zero-pads new IDs to. Older items
// with three digits stay valid; see design/adrs/0017-four-digit-ids.md.
const IDWidth = 4

var looseID = regexp.MustCompile(`^([ESTIesti])-?0*(\d+)$`)

// CanonicalID normalises a work item or issue ID as typed (s-32, S-032,
// S-0032, I-7) to the padded form flai allocates (S-0032, I-0007). Anything
// that is not an ID is returned unchanged.
func CanonicalID(id string) string {
	m := looseID.FindStringSubmatch(strings.TrimSpace(id))
	if m == nil {
		return id
	}
	return fmt.Sprintf("%s-%0*s", strings.ToUpper(m[1]), IDWidth, m[2])
}

// idCandidates lists the file-name forms an ID may be stored under: as given,
// canonical, and the legacy three-digit form.
func idCandidates(id string) []string {
	out := []string{id}
	canon := CanonicalID(id)
	if m := looseID.FindStringSubmatch(canon); m != nil {
		legacy := fmt.Sprintf("%s-%03s", m[1], m[2])
		for _, c := range []string{canon, legacy} {
			if !contains(out, c) {
				out = append(out, c)
			}
		}
	}
	return out
}

// Get finds one item by ID in kanban or archive. IDs may be given in short
// form (S-32) or with any zero padding; the item's own ID is returned.
func (r *Repo) Get(id string) (*Item, error) {
	typ := TypeOfID(CanonicalID(id))
	if typ == "" {
		return nil, fmt.Errorf("%q is not a work item ID (E-0001, S-0001, T-0001)", id)
	}
	for _, cand := range idCandidates(id) {
		for _, archived := range []bool{false, true} {
			dir := r.ItemDir(typ, archived)
			matches, _ := filepath.Glob(filepath.Join(dir, cand+"-*.md"))
			if len(matches) == 0 {
				continue
			}
			it, err := ReadItem(matches[0])
			if err != nil {
				return nil, err
			}
			if it.ID != cand {
				return nil, fmt.Errorf("%s: front matter id is %s", matches[0], it.ID)
			}
			it.Archived = archived
			return it, nil
		}
	}
	return nil, fmt.Errorf("%s not found", id)
}

// Children returns items whose parent is id, from the given list.
func Children(items []*Item, id string) []*Item {
	var out []*Item
	for _, it := range items {
		if it.Parent == id {
			out = append(out, it)
		}
	}
	return out
}

// NextID allocates the next ID for a type across kanban and archive.
func (r *Repo) NextID(typ string) (string, error) {
	prefix := strings.ToUpper(typ[:1]) + "-"
	max := 0
	for _, archived := range []bool{false, true} {
		matches, _ := filepath.Glob(filepath.Join(r.ItemDir(typ, archived), prefix+"*.md"))
		for _, m := range matches {
			base := strings.TrimPrefix(filepath.Base(m), prefix)
			num, _, _ := strings.Cut(base, "-")
			if n, err := strconv.Atoi(num); err == nil && n > max {
				max = n
			}
		}
	}
	return fmt.Sprintf("%s%0*d", prefix, IDWidth, max+1), nil
}

// FileName is <ID>-<slug>.md.
func FileName(id, title string) string {
	slug := template.Slug(title)
	if slug == "" {
		return id + ".md"
	}
	return id + "-" + slug + ".md"
}

// Save writes the item to its Path, creating the folder if needed.
func (r *Repo) Save(it *Item) error {
	if it.Path == "" {
		it.Path = filepath.Join(r.ItemDir(it.Type, it.Archived), FileName(it.ID, it.Title))
	}
	if err := os.MkdirAll(filepath.Dir(it.Path), 0o755); err != nil {
		return err
	}
	return atomicfile.WriteFile(it.Path, []byte(it.Marshal()), 0o644)
}

func lessID(a, b string) bool {
	ta, tb := TypeOfID(a), TypeOfID(b)
	if ta != tb {
		return typeRank(ta) < typeRank(tb)
	}
	na, _ := strconv.Atoi(strings.TrimLeft(a[2:], "0"))
	nb, _ := strconv.Atoi(strings.TrimLeft(b[2:], "0"))
	return na < nb
}

func typeRank(t string) int {
	switch t {
	case Epic:
		return 0
	case Story:
		return 1
	default:
		return 2
	}
}
