package components

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// The kept read (S-0601). Starting node costs 0.6 to 4.2 s a sub-project
// though nothing changed since the last read, so ReadProps keeps each
// component's answer and asks node only for the components whose answer may
// have changed. An answer stays true while its key does:
//
//   - the component file's size and modification time;
//   - the script's hash, scriptHash, so that a flai with another script reads
//     again;
//   - the sub-project's svelte version, from the svelte/package.json node
//     would resolve, and the hash of its package.json, which says whether it
//     is SvelteKit and what it depends on;
//   - for a component whose skeletons typescript resolved, the stamps of what
//     typescript read, which the script names in its answer (stampOf).
//
// The kept read is a cache: it lives under .flai-cache and is rebuilt from
// the components whenever it is missing or does not parse.

// keptVersion is the kept file's format; a file of another is a miss.
const keptVersion = 1

// KeptPropsPath is the file that keeps the props read for the sub-project at
// project, relative to root: .flai-cache/component-props-kept/<project>/props.json,
// outside PropsCacheDir, whose files the script prunes.
func KeptPropsPath(root, project string) string {
	return filepath.Join(root, ".flai-cache", "component-props-kept", filepath.FromSlash(project), "props.json")
}

// keptMu keeps one process's reads and writes of a kept file apart; the file
// is written through a rename, so another process reads one whole file or
// the other.
var keptMu sync.Mutex

// keptFile is a kept file.
type keptFile struct {
	Version int    `json:"version"`
	Script  string `json:"script"`
	Svelte  string `json:"svelte"`
	Package string `json:"package"`
	// Sets are the stamps of each node run whose answer is kept, each
	// stamp [kind, path, stamp] as the script gives it.
	Sets       [][][3]string            `json:"sets"`
	Components map[string]keptComponent `json:"components"`
}

// keptComponent is one component's kept answer, under its path.
type keptComponent struct {
	Size  int64  `json:"size"`
	Mod   int64  `json:"mod"`
	Props []Prop `json:"props"`
	Error string `json:"error,omitempty"`
	// Set is the index in Sets of the stamps its skeletons depend on, or -1.
	Set int `json:"set"`
}

// fileStat is a component file's size and modification time in nanoseconds,
// and whether it could be read.
type fileStat struct {
	size, mod int64
	ok        bool
}

// keptRead is the kept read of one sub-project for one ReadProps; nil when
// nothing is kept for it.
type keptRead struct {
	path, dir string
	head      keptFile // the key's parts every component shares, as they are now
	file      keptFile // the kept file as it was read
	valid     map[int]bool
	pre       map[string]fileStat // each component asked for, as it was before node read it
}

// openKept is the kept read of the sub-project at project, or nil when root
// has no .flai-cache folder, which it never makes (D-0402), or when the
// sub-project's svelte cannot be found, as in a fixture whose node is a fake.
func openKept(root, project string) *keptRead {
	if st, err := os.Stat(filepath.Join(root, ".flai-cache")); err != nil || !st.IsDir() {
		return nil
	}
	dir := filepath.Join(root, filepath.FromSlash(project))
	svelte := svelteVersion(dir)
	if svelte == "" {
		return nil
	}
	pkg := "-"
	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		sum := sha256.Sum256(data)
		pkg = hex.EncodeToString(sum[:])
	}
	k := &keptRead{
		path:  KeptPropsPath(root, project),
		dir:   dir,
		head:  keptFile{Version: keptVersion, Script: scriptHash, Svelte: svelte, Package: pkg},
		valid: map[int]bool{},
		pre:   map[string]fileStat{},
	}
	keptMu.Lock()
	k.file = k.read()
	keptMu.Unlock()
	return k
}

// svelteVersion is the version in the svelte/package.json that require
// resolves from dir: in node_modules in dir or a folder above it. It is
// empty when there is none.
func svelteVersion(dir string) string {
	for d := dir; ; d = filepath.Dir(d) {
		if data, err := os.ReadFile(filepath.Join(d, "node_modules", "svelte", "package.json")); err == nil {
			var pkg struct {
				Version string `json:"version"`
			}
			if json.Unmarshal(data, &pkg) == nil && pkg.Version != "" {
				return pkg.Version
			}
			return ""
		}
		if filepath.Dir(d) == d {
			return ""
		}
	}
}

// read is the kept file when it parses and its key's shared parts are as
// they are now; otherwise an empty one with them. The caller holds keptMu.
func (k *keptRead) read() keptFile {
	empty := k.head
	empty.Components = map[string]keptComponent{}
	data, err := os.ReadFile(k.path)
	if err != nil {
		return empty
	}
	var f keptFile
	if json.Unmarshal(data, &f) != nil || f.Version != k.head.Version || f.Script != k.head.Script || f.Svelte != k.head.Svelte || f.Package != k.head.Package || f.Components == nil {
		return empty
	}
	return f
}

// answer sets on out each component at idx whose kept answer is still true,
// and gives the indexes of the rest, which node is to read. With no kept
// read, that is all of idx.
func (k *keptRead) answer(out []Component, idx []int) []int {
	if k == nil {
		return idx
	}
	var miss []int
	for _, n := range idx {
		p := out[n].Path
		st := statFile(filepath.Join(k.dir, filepath.FromSlash(p)))
		k.pre[p] = st
		c, ok := k.file.Components[p]
		if !ok || !st.ok || c.Size != st.size || c.Mod != st.mod || !k.setValid(c.Set) {
			miss = append(miss, n)
			continue
		}
		out[n].Props, out[n].Error = c.Props, c.Error
		if out[n].Props == nil {
			out[n].Props = []Prop{}
		}
	}
	return miss
}

// setValid is whether set is no set, or one whose every stamp is as it is
// now.
func (k *keptRead) setValid(set int) bool {
	if set < 0 {
		return true
	}
	if set >= len(k.file.Sets) {
		return false
	}
	if v, ok := k.valid[set]; ok {
		return v
	}
	v := true
	for _, s := range k.file.Sets[set] {
		if now, ok := stampOf(s[0], s[1]); !ok || now != s[2] {
			v = false
			break
		}
	}
	k.valid[set] = v
	return v
}

// keep keeps the answers node gave for out's components at idx, each under
// its key as it was before node read it, merged into the kept file as it is
// now. A component node gave no answer for, or whose file could not be read,
// is not kept, nor is a stamped one when the answer's stamps do not parse.
// Keeping is a cache: a file that cannot be written leaves the next read to
// node, which is all a missing kept read costs.
func (k *keptRead) keep(out []Component, idx []int, a *answer) {
	if k == nil {
		return
	}
	got := make(map[string]answered, len(a.Components))
	for _, c := range a.Components {
		got[c.Path] = c
	}
	stamps, stampsOK := parseStamps(a.Stamps)

	keptMu.Lock()
	defer keptMu.Unlock()
	f := k.read()
	newSet := -1
	for _, n := range idx {
		p := out[n].Path
		c, answered := got[p]
		st := k.pre[p]
		if !answered || !st.ok || (c.Stamped && !stampsOK) {
			delete(f.Components, p)
			continue
		}
		set := -1
		if c.Stamped {
			if newSet < 0 {
				f.Sets = append(f.Sets, stamps)
				newSet = len(f.Sets) - 1
			}
			set = newSet
		}
		f.Components[p] = keptComponent{Size: st.size, Mod: st.mod, Props: out[n].Props, Error: out[n].Error, Set: set}
	}
	compactSets(&f)
	_ = writeKept(k.path, f)
}

// compactSets drops the sets no component refers to and renumbers the rest.
func compactSets(f *keptFile) {
	used := map[int]int{}
	var sets [][][3]string
	for p, c := range f.Components {
		if c.Set < 0 {
			continue
		}
		if c.Set >= len(f.Sets) {
			delete(f.Components, p)
			continue
		}
		to, ok := used[c.Set]
		if !ok {
			to = len(sets)
			sets = append(sets, f.Sets[c.Set])
			used[c.Set] = to
		}
		c.Set = to
		f.Components[p] = c
	}
	f.Sets = sets
}

// writeKept writes f to path through a file of its own in the same folder,
// renamed over path.
func writeKept(path string, f keptFile) error {
	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "props-*.tmp")
	if err != nil {
		return err
	}
	_, werr := tmp.Write(data)
	cerr := tmp.Close()
	if err := errors.Join(werr, cerr); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	return nil
}

// parseStamps is the answer's stamps, each [[kind, path], stamp], as
// [kind, path, stamp], and false when one is not of that shape.
func parseStamps(raw []json.RawMessage) ([][3]string, bool) {
	out := make([][3]string, 0, len(raw))
	for _, r := range raw {
		var e [2]json.RawMessage
		var kp [2]string
		var s string
		if json.Unmarshal(r, &e) != nil || json.Unmarshal(e[0], &kp) != nil || json.Unmarshal(e[1], &s) != nil {
			return nil, false
		}
		out = append(out, [3]string{kp[0], kp[1], s})
	}
	return out, true
}

// statFile is a file's size and modification time.
func statFile(path string) fileStat {
	st, err := os.Stat(path)
	if err != nil {
		return fileStat{}
	}
	return fileStat{size: st.Size(), mod: st.ModTime().UnixNano(), ok: true}
}

// stampOf is what path is now as kind sees it, the string props.mjs's stamp
// gives: for real, what it resolves to; otherwise "-" when it is missing,
// else "d" or "f", then for list its modification time in nanoseconds, and
// for file that and its size. ok is false when path cannot be looked at for
// another reason, which makes the stamp stale.
func stampOf(kind, path string) (string, bool) {
	if kind == "real" {
		r, err := filepath.EvalSymlinks(path)
		if err != nil {
			return "-", true
		}
		return r, true
	}
	st, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "-", true
	}
	if err != nil {
		return "", false
	}
	typ := "f"
	if st.IsDir() {
		typ = "d"
	}
	switch kind {
	case "dir":
		return typ, true
	case "list":
		return typ + ":" + strconv.FormatInt(st.ModTime().UnixNano(), 10), true
	case "file":
		return typ + ":" + strconv.FormatInt(st.ModTime().UnixNano(), 10) + ":" + strconv.FormatInt(st.Size(), 10), true
	default:
		return "", false
	}
}
