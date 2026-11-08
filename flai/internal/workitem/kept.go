package workitem

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// keptVersion is the format of a kept items file. Raise it whenever Item, or
// a type it holds, changes in a way that the JSON of an older flai would not
// restore exactly; a file of another version is a miss.
const keptVersion = 1

// keptItemsFile is what a folder's kept items file holds between processes
// (S-0600): the absolute folder it describes, its format, and per item file
// the key it was parsed at and the item as ParseItem left it.
type keptItemsFile struct {
	Version int          `json:"version"`
	Dir     string       `json:"dir"`
	Files   []keptRecord `json:"files"`
}

// keptRecord is one item file of a kept items file. Size and ModTime are its
// stat when it was parsed; an entry is written only once that parse was
// trusted, racyWindow after the file's modification.
type keptRecord struct {
	Name    string          `json:"name"`
	Size    int64           `json:"size"`
	ModTime int64           `json:"mod_time"` // Unix nanoseconds
	Item    json.RawMessage `json:"item"`
	Unknown []Field         `json:"unknown,omitempty"` // Item's JSON leaves them out
}

// fileKey is what tells a file's kept parse still holds: its size and
// modification time.
type fileKey struct {
	size int64
	mod  int64
}

func keyOf(info os.FileInfo) fileKey {
	return fileKey{size: info.Size(), mod: info.ModTime().UnixNano()}
}

// keptItemsPath is the kept items file of the item folder dir in the cache
// folder cache, named after dir's path below the project's root, or after a
// hash of dir when it is not below it. With no cache, there is none.
func keptItemsPath(cache, dir string) string {
	if cache == "" {
		return ""
	}
	name := ""
	if rel, err := filepath.Rel(filepath.Dir(cache), dir); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
		name = strings.ReplaceAll(filepath.ToSlash(rel), "/", "-")
	} else {
		sum := sha256.Sum256([]byte(dir))
		name = "dir-" + hex.EncodeToString(sum[:8])
	}
	return filepath.Join(cache, "items-"+name+".json")
}

// readKeptItems is the kept items file at path, by file name, when it is of
// this format and describes dir. A file that is missing, cannot be read or
// decoded, is of another version, or describes another folder holds none.
func readKeptItems(path, dir string) map[string]keptRecord {
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var kf keptItemsFile
	if err := json.Unmarshal(data, &kf); err != nil || kf.Version != keptVersion || kf.Dir != dir {
		return nil
	}
	out := make(map[string]keptRecord, len(kf.Files))
	for _, rec := range kf.Files {
		if rec.Name != "" && len(rec.Item) > 0 {
			out[rec.Name] = rec
		}
	}
	return out
}

// restoreItem rebuilds the item of rec as ReadItem gives it for the file
// path: Path from where the file is, Archived left to List, the unknown
// fields put back, and readAs made again from the result.
func restoreItem(path string, rec keptRecord) (*Item, error) {
	var it Item
	if err := json.Unmarshal(rec.Item, &it); err != nil {
		return nil, err
	}
	it.Unknown = rec.Unknown
	it.Path = path
	it.Archived = false
	it.readAs = it.Marshal()
	return &it, nil
}

// encodeKept is the JSON a kept items file holds for it, leaving out the
// fields that come from where the file is.
func encodeKept(it *Item) (json.RawMessage, error) {
	c := *it
	c.Path, c.Archived = "", false
	return json.Marshal(&c)
}

// writeKept writes the folder's trusted parses to its kept items file, when
// they differ from what the file was last known to hold, through a file
// beside it renamed over it, so that a reader finds the old file or the new
// one whole. It never makes the cache folder: a project without one keeps
// nothing (D-0402). A file not written costs the next process its parses,
// never a wrong answer, so a failure is dropped. The caller holds f.mu.
func (f *folder) writeKept(path, dir string) {
	if path == "" {
		return
	}
	now := map[string]fileKey{}
	for name, k := range f.files {
		if k.trusted {
			now[name] = keyOf(k.info)
		}
	}
	if f.written == nil && len(now) == 0 || f.written != nil && maps.Equal(now, f.written) {
		return
	}
	kf := keptItemsFile{Version: keptVersion, Dir: dir, Files: make([]keptRecord, 0, len(now))}
	for _, name := range slices.Sorted(maps.Keys(now)) {
		k := f.files[name]
		if k.raw == nil {
			raw, err := encodeKept(k.item)
			if err != nil {
				return
			}
			k.raw = raw
		}
		key := now[name]
		kf.Files = append(kf.Files, keptRecord{Name: name, Size: key.size, ModTime: key.mod, Item: k.raw, Unknown: k.item.Unknown})
	}
	if writeKeptItems(path, kf) == nil {
		f.written = now
	}
}

// writeKeptItems writes kf to path through a temporary file in the same
// folder renamed over it.
func writeKeptItems(path string, kf keptItemsFile) error {
	data, err := json.Marshal(kf)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".items-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return err
	}
	return nil
}
