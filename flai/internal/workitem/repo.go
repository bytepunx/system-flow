package workitem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/execx"
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
	// MainAbsent is the main checkout a linked worktree's .git file names
	// when that checkout's system-flow.yaml is not here, as in an agent
	// container that mounts the worktree alone (I-0219); MainRoot is then
	// Root. Empty otherwise.
	MainAbsent string
	// Git runs git for the rules that ask it: a story goes to review only
	// with its worktree committed (S-0140). Nil skips them.
	Git execx.Runner
	// Trial judges by a trial merge whether an overlap holds a ready story
	// (S-0551), as Repo.Holds applies it. Nil leaves the overlap hold as it
	// was.
	Trial TrialJudge
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
	main, absent := mainCheckout(root)
	return &Repo{Root: root, Manifest: m, MainRoot: main, MainAbsent: absent}, nil
}

// mainCheckout resolves a linked worktree's main checkout from its .git
// file ("gitdir: <main>/.git/worktrees/<name>"); a normal checkout returns
// itself. When the main checkout it names has no system-flow.yaml, it
// returns root and, as absent, the main checkout's path.
func mainCheckout(root string) (main, absent string) {
	data, err := os.ReadFile(filepath.Join(root, ".git"))
	if err != nil {
		return root, ""
	}
	line := strings.TrimSpace(string(data))
	gitdir, ok := strings.CutPrefix(line, "gitdir: ")
	if !ok {
		return root, ""
	}
	gitdir = filepath.Clean(gitdir)
	marker := string(filepath.Separator) + ".git" + string(filepath.Separator) + "worktrees" + string(filepath.Separator)
	i := strings.Index(gitdir, marker)
	if i < 0 {
		return root, ""
	}
	main = gitdir[:i]
	if !filepath.IsAbs(main) {
		main = filepath.Join(root, main)
	}
	if _, err := os.Stat(filepath.Join(main, "system-flow.yaml")); err != nil {
		return root, main
	}
	return main, ""
}

// IsWorktree reports whether Root is a linked worktree of MainRoot.
func (r *Repo) IsWorktree() bool { return r.MainRoot != "" && r.MainRoot != r.Root }

// RequireMain returns an error naming the command what when Root is a story
// worktree whose main checkout is not here, and nil otherwise.
func (r *Repo) RequireMain(what string) error {
	if r.MainAbsent == "" {
		return nil
	}
	return fmt.Errorf("%s needs the main checkout %s, which is not here: %s is a story worktree mounted without it, as in an agent container; run it on the host, or call the MCP tool verify, task_done, or stream_sync, which run there", what, r.MainAbsent, r.Root)
}

// CacheDir is .flai-cache in the main checkout: token, worktrees, tool caches.
func (r *Repo) CacheDir() string { return filepath.Join(r.mainRoot(), ".flai-cache") }

// WorktreePath is where a story's branch is checked out: Root itself when
// Root is that story's worktree and its main checkout is not here.
func (r *Repo) WorktreePath(storyID string) string {
	if r.MainAbsent != "" && filepath.Base(r.Root) == storyID {
		return r.Root
	}
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
// by ID. Each is the caller's own copy; only files that changed since the
// process last read them are read again (S-0156), and a process's first list
// of a folder parses only the files that changed since a process last kept
// it in the cache folder (S-0600). The folders are listed at once, each under
// its own lock, sharing one bound of parses at once (S-0676); an error is the
// first folder's, in the order current before archived and Types.
func (r *Repo) List(includeArchive bool) ([]*Item, error) {
	type folderList struct {
		dir      string
		archived bool
		items    []*Item
		err      error
	}
	var lists []*folderList
	scopes := []bool{false}
	if includeArchive {
		scopes = append(scopes, true)
	}
	for _, archived := range scopes {
		for _, typ := range Types {
			lists = append(lists, &folderList{dir: r.ItemDir(typ, archived), archived: archived})
		}
	}
	cache, slots := r.CacheDir(), parseSlots()
	var wg sync.WaitGroup
	for _, l := range lists {
		wg.Go(func() { l.items, l.err = parsed.listOn(l.dir, cache, slots) })
	}
	wg.Wait()
	var items []*Item
	for _, l := range lists {
		if l.err != nil {
			return nil, l.err
		}
		for _, it := range l.items {
			it.Archived = l.archived
		}
		items = append(items, l.items...)
	}
	sort.Slice(items, func(i, j int) bool { return lessID(items[i].ID, items[j].ID) })
	return items, nil
}

// Generation counts, in this process, the changes to the items List reads,
// the archive's among them (S-0608). It moves when a list or get finds an
// item file added, changed, or removed, when the store is told of a change in
// an item folder, and when this process saves an item; it stays the same
// while none of that happens, and never goes back. Read it before and after a
// List: when the two agree, the items listed are those of that generation, so
// what is computed from them may be kept under it.
func (r *Repo) Generation() uint64 {
	var n uint64
	for _, archived := range []bool{false, true} {
		for _, typ := range Types {
			n += parsed.generation(r.ItemDir(typ, archived))
		}
	}
	return n
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

// ErrNotFound means no item has the ID, in kanban or archive.
var ErrNotFound = errors.New("not found")

// Get finds one item by ID in kanban or archive. IDs may be given in short
// form (S-32) or with any zero padding; the item's own ID is returned. An ID
// no item has is ErrNotFound.
func (r *Repo) Get(id string) (*Item, error) {
	typ := TypeOfID(CanonicalID(id))
	if typ == "" {
		return nil, fmt.Errorf("%q is not a work item ID (E-0001, S-0001, T-0001)", id)
	}
	for _, cand := range idCandidates(id) {
		for _, archived := range []bool{false, true} {
			dir := r.ItemDir(typ, archived)
			var it *Item
			if parsed.snapshotting(dir) {
				// the folder's snapshot answers, with no glob (S-0606)
				found, err := parsed.find(dir, cand+"-", r.CacheDir())
				if err != nil {
					return nil, err
				}
				if found == nil {
					continue
				}
				it = found
			} else {
				matches, _ := filepath.Glob(filepath.Join(dir, cand+"-*.md"))
				if len(matches) == 0 {
					continue
				}
				found, err := parsed.get(matches[0])
				if err != nil {
					return nil, err
				}
				it = found
			}
			if it.ID != cand {
				return nil, fmt.Errorf("%s: front matter id is %s", it.Path, it.ID)
			}
			it.Archived = archived
			return it, nil
		}
	}
	return nil, fmt.Errorf("%s %w", id, ErrNotFound)
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

// FileName is <ID>-<slug>.md, its slug cut at template.SlugMax bytes on a word.
func FileName(id, title string) string {
	slug := template.CutSlug(template.Slug(title), template.SlugMax)
	if slug == "" {
		return id + ".md"
	}
	return id + "-" + slug + ".md"
}

// Save writes the item to its Path, creating the folder if needed. What
// another writer saved to the file since the item was read is kept, and a
// change both made to the same key refuses the save with a *ConflictError,
// writing nothing (I-0175). An updated before created or the last transition
// is raised to the later of the two (I-0183). On success the item is what was
// written.
func (r *Repo) Save(it *Item) error {
	unlock, err := r.lockItems()
	if err != nil {
		return fmt.Errorf("save %s: %w", it.ID, err)
	}
	defer unlock()
	return r.saveLocked(it)
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
