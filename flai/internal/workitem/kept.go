package workitem

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// keptVersion is the format of a kept items file. Raise it whenever Item, or
// a type it holds, changes in a way that the JSON of an older flai would not
// restore exactly, or the file's lines change; a file of another version is
// a miss. Version 1 was one JSON document, items-<folder>.json (S-0600);
// version 2 is a line per item, items-<folder>.jsonl (S-0676).
const keptVersion = 2

// keptHeader is the first line of a kept items file: its format and the
// absolute folder it describes.
type keptHeader struct {
	Version int    `json:"version"`
	Dir     string `json:"dir"`
}

// keptRecord is one item file of a kept items file, one line of it (S-0676):
// its name as a JSON string, its size and modification time in Unix
// nanoseconds, its unknown fields as JSON or null, and the item as
// ParseItem left it, as JSON, each apart from the next by a tab. JSON escapes
// every tab and newline inside a value, so neither can break a line. Size and
// ModTime are the file's stat when it was parsed; an entry is written only
// once that parse was trusted, racyWindow after the file's modification.
// Reading the file splits its lines and reads only the key; the item and its
// unknown fields are decoded only for a file whose key matches.
type keptRecord struct {
	Name    string
	Size    int64
	ModTime int64           // Unix nanoseconds
	Item    json.RawMessage // never changed: it may share the file's bytes
	Unknown json.RawMessage // the item's unknown fields, which Item's JSON leaves out; nil for none
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
	return filepath.Join(cache, "items-"+name+".jsonl")
}

// readKeptItems is the kept items file at path, by file name, when it is of
// this format and describes dir. A file that is missing or cannot be read,
// whose first line is not this format's header for dir, holds none; a line
// that does not split into a record is a miss for its file alone. The items
// are left undecoded, each a slice of the file's bytes.
func readKeptItems(path, dir string) map[string]keptRecord {
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	head, rest, _ := bytes.Cut(data, []byte{'\n'})
	var h keptHeader
	if err := json.Unmarshal(head, &h); err != nil || h.Version != keptVersion || h.Dir != dir {
		return nil
	}
	out := make(map[string]keptRecord, bytes.Count(rest, []byte{'\n'}))
	for len(rest) > 0 {
		var line []byte
		line, rest, _ = bytes.Cut(rest, []byte{'\n'})
		if rec, ok := parseKeptLine(line); ok {
			out[rec.Name] = rec
		}
	}
	return out
}

// parseKeptLine splits one record line of a kept items file; ok is false for
// a line that does not split into one.
func parseKeptLine(line []byte) (keptRecord, bool) {
	var parts [5][]byte
	for i := range 4 {
		var found bool
		parts[i], line, found = bytes.Cut(line, []byte{'\t'})
		if !found {
			return keptRecord{}, false
		}
	}
	parts[4] = line
	name, ok := keptName(parts[0])
	if !ok {
		return keptRecord{}, false
	}
	size, err1 := strconv.ParseInt(string(parts[1]), 10, 64)
	mod, err2 := strconv.ParseInt(string(parts[2]), 10, 64)
	if err1 != nil || err2 != nil || len(parts[4]) == 0 {
		return keptRecord{}, false
	}
	rec := keptRecord{Name: name, Size: size, ModTime: mod, Item: json.RawMessage(parts[4][:len(parts[4]):len(parts[4])])}
	if u := parts[3]; string(u) != "null" {
		rec.Unknown = json.RawMessage(u[:len(u):len(u)])
	}
	return rec, true
}

// keptName decodes a record's name, a JSON string; one with no escape is
// taken as it is, without a decoder.
func keptName(raw []byte) (string, bool) {
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return "", false
	}
	if inner := raw[1 : len(raw)-1]; !bytes.ContainsAny(inner, "\\\"") {
		return string(inner), len(inner) > 0
	}
	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return "", false
	}
	return name, name != ""
}

// restoreItem rebuilds the item of rec as ReadItem gives it for the file
// path: Path from where the file is, Archived left to List, the unknown
// fields put back, and readAs made again from the result.
func restoreItem(path string, rec keptRecord) (*Item, error) {
	var it Item
	if err := json.Unmarshal(rec.Item, &it); err != nil {
		return nil, err
	}
	if len(rec.Unknown) > 0 {
		if err := json.Unmarshal(rec.Unknown, &it.Unknown); err != nil {
			return nil, err
		}
	}
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

// encodeUnknown is the JSON a kept items file holds for an item's unknown
// fields: nil, written as null, for none.
func encodeUnknown(fields []Field) (json.RawMessage, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	return json.Marshal(fields)
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
	recs := make([]keptRecord, 0, len(now))
	for _, name := range slices.Sorted(maps.Keys(now)) {
		k := f.files[name]
		if k.raw == nil {
			raw, err := encodeKept(k.item)
			if err != nil {
				return
			}
			k.raw = raw
		}
		unknown, err := encodeUnknown(k.item.Unknown)
		if err != nil {
			return
		}
		key := now[name]
		recs = append(recs, keptRecord{Name: name, Size: key.size, ModTime: key.mod, Item: k.raw, Unknown: unknown})
	}
	if writeKeptItems(path, dir, recs) == nil {
		f.written = now
	}
}

// encodeKeptItems is the kept items file of dir holding recs.
func encodeKeptItems(dir string, recs []keptRecord) ([]byte, error) {
	head, err := json.Marshal(keptHeader{Version: keptVersion, Dir: dir})
	if err != nil {
		return nil, err
	}
	size := len(head) + 1
	for _, rec := range recs {
		size += len(rec.Name) + len(rec.Item) + len(rec.Unknown) + 64
	}
	buf := bytes.NewBuffer(make([]byte, 0, size))
	buf.Write(head)
	buf.WriteByte('\n')
	for _, rec := range recs {
		name, err := json.Marshal(rec.Name)
		if err != nil {
			return nil, err
		}
		buf.Write(name)
		buf.WriteByte('\t')
		buf.WriteString(strconv.FormatInt(rec.Size, 10))
		buf.WriteByte('\t')
		buf.WriteString(strconv.FormatInt(rec.ModTime, 10))
		buf.WriteByte('\t')
		if len(rec.Unknown) > 0 {
			buf.Write(rec.Unknown)
		} else {
			buf.WriteString("null")
		}
		buf.WriteByte('\t')
		buf.Write(rec.Item)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}

// writeKeptItems writes the kept items file of dir holding recs to path,
// through a temporary file in the same folder renamed over it. It leaves a
// version 1 file, items-<folder>.json, where it is: an older flai on the same
// host, such as an installed MCP server, still reads and writes it, and would
// parse every item again in each process were it removed.
func writeKeptItems(path, dir string, recs []keptRecord) error {
	data, err := encodeKeptItems(dir, recs)
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
