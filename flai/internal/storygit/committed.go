package storygit

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// storyID is a story's ID as a commit subject names it in brackets, in any
// zero padding.
var storyID = regexp.MustCompile(`^S-\d+$`)

// committedRead is what Committed last read in one repository root, and the
// key it read it at.
type committedRead struct {
	key       string // the main branch, the wip folder, and git show-ref --heads
	committed map[string][]string
}

// keptCommitted holds Committed's last read per repository root. Committed
// holds its lock throughout, so two calls never read the log at once.
var keptCommitted = struct {
	sync.Mutex
	roots map[string]committedRead
}{roots: map[string]committedRead{}}

// Committed is, per story by its canonical ID, the files changed by the
// commits on the main branch and the story branches whose subject names it in
// brackets ([S-nnnn]), merges left out, and the files under the wip folder
// left out, since flai writes them (S-0205). Paths are relative to the
// project's root, sorted and without repeats; a story none of whose commits
// changed a file outside wip has no entry. It reads git log, in the main
// checkout, once per change of a branch tip, keeping the result per
// repository root in memory for the life of the process (I-0184) and in
// CommittedPath of the main checkout's cache folder between processes
// (I-0209), when that folder exists: a call while no tip moved costs two git
// processes, rev-parse and show-ref, and no log, in a new process as in one
// that read it before. Each call returns a copy of its own. It fails when git or the repository cannot
// be read.
func Committed(r execx.Runner, repo *workitem.Repo) (map[string][]string, error) {
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	if _, err := r.LookPath("git"); err != nil {
		return nil, fmt.Errorf("cannot read the stories' commits: git is not on PATH (%w); install git to compare touches with commits", err)
	}
	main, err := MainBranch(r, root)
	if err != nil {
		return nil, fmt.Errorf("cannot read the stories' commits in %s: %w; run flai in a git checkout with a branch and a commit on it", root, err)
	}
	wip := strings.TrimSuffix(repo.Manifest.Layout["wip"], "/") + "/"
	keptCommitted.Lock()
	defer keptCommitted.Unlock()
	// Every local branch's tip, story or not: a move of any of them reads the
	// log again, which is never wrong and cheap to tell.
	tips, err := r.Run(root, "git", "show-ref", "--heads")
	if err != nil {
		return nil, fmt.Errorf("cannot read the branch tips in %s to tell whether the stories' commits changed: %w; run flai in a git checkout with a branch and a commit on it", root, err)
	}
	key := main + "\n" + wip + "\n" + tips
	if prev, ok := keptCommitted.roots[root]; ok && prev.key == key {
		return copyCommitted(prev.committed), nil
	}
	path := CommittedPath(repo.CacheDir())
	if committed, ok := readKeptCommitted(path, root, key); ok {
		keptCommitted.roots[root] = committedRead{key: key, committed: committed}
		return copyCommitted(committed), nil
	}
	committed, err := readCommitted(r, root, main, wip)
	if err != nil {
		return nil, err
	}
	keptCommitted.roots[root] = committedRead{key: key, committed: committed}
	// A file not written costs the next process a read of the log, never a
	// wrong answer, so the call answers all the same.
	_ = writeKeptCommitted(path, keptCommittedFile{Root: root, Key: key, Committed: committed})
	return copyCommitted(committed), nil
}

// CommittedPath is where Committed keeps its last read between processes in
// the cache folder cacheDir.
func CommittedPath(cacheDir string) string {
	return filepath.Join(cacheDir, "stories-commits.json")
}

// keptCommittedFile is what CommittedPath holds: Committed's last read, the
// repository root it read it in, and the key it read it at.
type keptCommittedFile struct {
	Root      string              `json:"root"`
	Key       string              `json:"key"` // the main branch, the wip folder, and git show-ref --heads
	Committed map[string][]string `json:"committed"`
}

// readKeptCommitted is the answer kept at path when it was read in root at
// key. A file that is missing, cannot be read or parsed, or was read in
// another root or at another key, holds none.
func readKeptCommitted(path, root, key string) (map[string][]string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var kept keptCommittedFile
	if err := json.Unmarshal(data, &kept); err != nil || kept.Root != root || kept.Key != key || kept.Committed == nil {
		return nil, false
	}
	return kept.Committed, true
}

// writeKeptCommitted writes kept to path through a file beside it renamed
// over it, so that a reader finds the old answer or the new one whole. It
// never makes the cache folder: a project without one, such as one with no
// .gitignore to ignore it, keeps nothing rather than gain a folder that
// flai check would report (layout.gitignore).
func writeKeptCommitted(path string, kept keptCommittedFile) error {
	data, err := json.Marshal(kept)
	if err != nil {
		return fmt.Errorf("encode the stories' commits for %s: %w", path, err)
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".stories-commits-*")
	if err != nil {
		return fmt.Errorf("write the stories' commits in %s: %w", dir, err)
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return fmt.Errorf("write the stories' commits to %s: %w", name, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("close the stories' commits in %s: %w", name, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("keep the stories' commits at %s: %w", path, err)
	}
	return nil
}

// readCommitted reads Committed's answer from git log in root, the main
// branch and the story branches, leaving out the files under wip.
func readCommitted(r execx.Runner, root, main, wip string) (map[string][]string, error) {
	// --relative and the pathspec keep the files under the project's root and
	// name them from it; --full-history keeps every commit that touched them,
	// whichever side of a merge it is on.
	out, err := r.Run(root, "git", "-c", "core.quotePath=false", "log", "--no-merges", "--no-renames", "--full-history", "--relative",
		"--name-only", "--format=%x1e%s", main, "--branches="+Prefix+"*", "--", ".")
	if err != nil {
		return nil, fmt.Errorf("cannot read the commits of %s and the story branches in %s: %w", main, root, err)
	}
	files := map[string]map[string]bool{}
	for _, rec := range strings.Split(out, "\x1e") {
		lines := strings.Split(rec, "\n")
		ids := subjectStories(lines[0])
		if len(ids) == 0 {
			continue
		}
		for _, f := range lines[1:] {
			if f == "" || strings.HasPrefix(f, wip) {
				continue
			}
			for _, id := range ids {
				if files[id] == nil {
					files[id] = map[string]bool{}
				}
				files[id][f] = true
			}
		}
	}
	committed := make(map[string][]string, len(files))
	for id, set := range files {
		committed[id] = slices.Sorted(maps.Keys(set))
	}
	return committed, nil
}

// copyCommitted is a copy of committed that shares no slice with it, so that
// a caller who changes its answer cannot change what is kept.
func copyCommitted(committed map[string][]string) map[string][]string {
	c := make(map[string][]string, len(committed))
	for id, files := range committed {
		c[id] = slices.Clone(files)
	}
	return c
}

// subjectStories is every story a commit subject names in brackets, by its
// canonical ID, read as release reads the IDs of a message.
func subjectStories(subject string) []string {
	var ids []string
	for i := 0; i < len(subject); i++ {
		if subject[i] != '[' {
			continue
		}
		end := strings.IndexByte(subject[i+1:], ']')
		if end < 0 {
			break
		}
		if tok := subject[i+1 : i+1+end]; storyID.MatchString(tok) {
			if id := workitem.CanonicalID(tok); !slices.Contains(ids, id) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}
