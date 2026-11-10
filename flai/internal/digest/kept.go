package digest

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RulesVersion is the version of what validate decides of an entry file.
// Raise it whenever a check or a message of validate changes in a way the
// rules hash cannot see, the kinds, the statuses, and the secret shapes, so
// that problems kept under the old rules are not reused (S-0677).
const RulesVersion = 1

// keptVersion is the format of the kept problems file; a file of another
// format is a miss.
const keptVersion = 1

// racyWindow is how long after an entry's modification time its problems
// are still not kept: a second write within the file system's timestamp
// granularity could leave the same time and size, as mdlint's kept findings
// say. An entry checked that soon after it changed is checked again on the
// next run.
const racyWindow = 2 * time.Second

// readEntryFile reads an entry file CheckFolder validates; a test replaces it to
// count the files read.
var readEntryFile = os.ReadFile

// KeptPath is where CheckFolder keeps the problems of the digest folder dir
// between processes in the cache folder cacheDir: one file per folder, named
// by a hash of the folder's path, so that a story's worktree and the main
// checkout do not overwrite each other's.
func KeptPath(cacheDir, dir string) string {
	sum := sha256.Sum256([]byte(filepath.Clean(dir)))
	return filepath.Join(cacheDir, fmt.Sprintf("digest-problems-%x.json", sum[:6]))
}

// FileProblems is one entry file's problems as CheckFolder reports them:
// what Validate finds, or the error that kept the file from being read.
type FileProblems struct {
	Path     string
	Problems []Problem
	Err      error
}

// CheckFolder validates every entry file in the digest folder dir, in name
// order, as Validate does with env, and returns those with a problem or a
// read error. It keeps what each file alone decides between processes in the
// cache folder cacheDir, at KeptPath (S-0677): an entry whose size and
// modification time are those it was checked at, under the same rules, takes
// its kept problems unread. What env answers, whether a source ID or a path
// names something, follows other files and is asked every time, so an entry
// whose source was archived or removed since is still reported. The kept
// file is written again when what it holds changed, and only into a cacheDir
// that exists. A kept file that is missing, damaged, or another folder's or
// rules' is a miss, and a failed write fails nothing: either costs the next
// run a read, never a wrong problem. An empty cacheDir keeps nothing. A
// folder that cannot be read is the error.
func CheckFolder(dir, cacheDir string, env Env) ([]FileProblems, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	keep := cacheDir != "" && isDir(cacheDir)
	var kept map[string]keptEntry
	var keptPath string
	if keep {
		keptPath = KeptPath(cacheDir, dir)
		kept = readKept(keptPath, filepath.Clean(dir))
	}
	checkedAt := time.Now()
	files := map[string]keptEntry{}
	var out []FileProblems
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || Skipped(name) {
			continue
		}
		path := filepath.Join(dir, name)
		var info fs.FileInfo
		if keep {
			// An entry whose stat fails is read, which fails too, and is
			// reported as before.
			info, _ = os.Stat(path)
		}
		var problems []Problem
		if k, ok := kept[name]; ok && info != nil && k.matches(info) {
			problems = k.problems()
			files[name] = k
		} else {
			data, err := readEntryFile(path)
			if err != nil {
				out = append(out, FileProblems{Path: path, Err: err})
				continue
			}
			_, problems = validate(name, string(data))
			if info != nil && info.ModTime().Before(checkedAt.Add(-racyWindow)) {
				files[name] = newKeptEntry(info, problems)
			}
		}
		if ps := Resolve(problems, env); len(ps) > 0 {
			out = append(out, FileProblems{Path: path, Problems: ps})
		}
	}
	if keep && !sameKept(kept, files) {
		_ = writeKept(keptPath, keptFile{Version: keptVersion, Root: filepath.Clean(dir), Hash: rulesHash(), Files: files})
	}
	return out, nil
}

// rulesHash identifies the rules validate applies and the kept file's
// format: RulesVersion, keptVersion, the kinds, the statuses, and the secret
// shapes. Problems kept under one hash are reused only under the same hash.
var rulesHash = sync.OnceValue(func() string {
	h := sha256.New()
	fmt.Fprintf(h, "digest rules %d kept %d\n", RulesVersion, keptVersion)
	fmt.Fprintf(h, "kinds %s\nstatuses %s\n", strings.Join(Kinds, ","), strings.Join(Statuses, ","))
	fmt.Fprintf(h, "assigned %s\n", assigned.String())
	for _, s := range shapes {
		fmt.Fprintf(h, "shape %s %s %s %d %t\n", s.kind, s.name, s.re.String(), s.value, s.also != nil)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
})

// keptFile is what KeptPath holds: the folder checked, the rules hash it was
// checked under, and each entry's key and problems, by its file name.
type keptFile struct {
	Version int                  `json:"version"`
	Root    string               `json:"root"`
	Hash    string               `json:"hash"`
	Files   map[string]keptEntry `json:"files"`
}

// keptEntry is one entry file's problems and the stat they were found at.
type keptEntry struct {
	Size     int64         `json:"size"`
	ModTime  int64         `json:"mtime"` // Unix nanoseconds
	Problems []keptProblem `json:"problems"`
}

// keptProblem is a Problem with what Resolve asks of it, which Problem's
// JSON leaves out.
type keptProblem struct {
	Field   string `json:"field,omitempty"`
	Line    int    `json:"line"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
	Unknown string `json:"unknown,omitempty"`
	Missing string `json:"missing,omitempty"`
}

func newKeptEntry(info fs.FileInfo, problems []Problem) keptEntry {
	k := keptEntry{Size: info.Size(), ModTime: info.ModTime().UnixNano(), Problems: make([]keptProblem, len(problems))}
	for i, p := range problems {
		k.Problems[i] = keptProblem(p)
	}
	return k
}

// matches reports whether the entry was checked at info's size and
// modification time.
func (k keptEntry) matches(info fs.FileInfo) bool {
	return k.Size == info.Size() && k.ModTime == info.ModTime().UnixNano()
}

func (k keptEntry) problems() []Problem {
	out := make([]Problem, len(k.Problems))
	for i, p := range k.Problems {
		out[i] = Problem(p)
	}
	return out
}

// sameKept reports whether two sets of kept entries hold the same files at
// the same keys: problems follow from the key, so they are not compared.
func sameKept(a, b map[string]keptEntry) bool {
	if a == nil || len(a) != len(b) {
		return false
	}
	for name, k := range a {
		o, ok := b[name]
		if !ok || o.Size != k.Size || o.ModTime != k.ModTime {
			return false
		}
	}
	return true
}

// readKept is the problems kept at path for the folder root. A file that is
// missing, cannot be read or parsed, is of another format, or was checked in
// another folder or under other rules, holds none.
func readKept(path, root string) map[string]keptEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var k keptFile
	if err := json.Unmarshal(data, &k); err != nil || k.Version != keptVersion || k.Root != root || k.Hash != rulesHash() {
		return nil
	}
	return k.Files
}

// writeKept writes k to path through a file beside it renamed over it, so
// that a reader finds the old problems or the new ones whole. It never makes
// the cache folder: a project without one keeps nothing.
func writeKept(path string, k keptFile) error {
	data, err := json.Marshal(k)
	if err != nil {
		return fmt.Errorf("encode the kept digest problems for %s: %w", path, err)
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".digest-problems-*")
	if err != nil {
		return fmt.Errorf("write the kept digest problems in %s: %w", dir, err)
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return fmt.Errorf("write the kept digest problems to %s: %w", name, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("close the kept digest problems in %s: %w", name, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("keep the digest problems at %s: %w", path, err)
	}
	return nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
