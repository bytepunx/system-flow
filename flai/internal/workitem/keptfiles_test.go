package workitem

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// note is the value the kept-files tests keep: a file's text and its words.
type note struct {
	Text  string
	Words []string
}

func parseNote(path string) (*note, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(data)
	if strings.HasPrefix(text, "bad") {
		return nil, errors.New("does not parse")
	}
	return &note{Text: text, Words: strings.Fields(text)}, nil
}

func cloneNote(n *note) *note {
	c := *n
	c.Words = slices.Clone(n.Words)
	return &c
}

func newNotes() *KeptFiles[*note] { return NewKeptFiles(parseNote, cloneNote) }

// longAgo is a modification time well past the racy window.
var longAgo = time.Now().Add(-time.Hour)

// writeNote writes text to name in dir, modified at mod.
func writeNote(t *testing.T, dir, name, text string, mod time.Time) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatal(err)
	}
	return path
}

// texts lists dir and returns each file's base name and text, failing on an
// error.
func texts(t *testing.T, k *KeptFiles[*note], dir string, match func(string) bool) map[string]string {
	t.Helper()
	files, err := k.List(dir, match)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, f := range files {
		if f.Err != nil {
			t.Fatalf("%s: %v", f.Path, f.Err)
		}
		out[filepath.Base(f.Path)] = f.Value.Text
	}
	return out
}

func TestKeptFilesSecondListParsesNothing(t *testing.T) {
	dir := t.TempDir()
	writeNote(t, dir, "a.md", "one", longAgo)
	writeNote(t, dir, "b.md", "two", longAgo)
	k := newNotes()
	first, err := k.List(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if k.Parses() != 2 || !first[0].Parsed || !first[1].Parsed {
		t.Fatalf("first list: %d parses, want 2, each parsed", k.Parses())
	}
	second, err := k.List(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if k.Parses() != 2 {
		t.Fatalf("second list parsed again: %d parses, want 2", k.Parses())
	}
	if second[0].Parsed || second[1].Parsed {
		t.Fatal("second list says it parsed a kept file")
	}
	if second[0].Value.Text != "one" || second[1].Value.Text != "two" || filepath.Base(second[1].Path) != "b.md" {
		t.Fatalf("second list = %+v, %+v", *second[0].Value, *second[1].Value)
	}
}

func TestKeptFilesSeesEachChange(t *testing.T) {
	dir := t.TempDir()
	writeNote(t, dir, "a.md", "one", longAgo)
	writeNote(t, dir, "b.md", "two", longAgo)
	writeNote(t, dir, "c.md", "three", longAgo)
	k := newNotes()
	texts(t, k, dir, nil)

	writeNote(t, dir, "a.md", "one more", longAgo)             // another size
	writeNote(t, dir, "b.md", "owt", longAgo.Add(time.Minute)) // same size, new time
	if err := os.Remove(filepath.Join(dir, "c.md")); err != nil {
		t.Fatal(err)
	}
	writeNote(t, dir, "d.md", "four", longAgo)
	got := texts(t, k, dir, nil)
	want := map[string]string{"a.md": "one more", "b.md": "owt", "d.md": "four"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("after the changes = %v, want %v", got, want)
	}
	if k.Parses() != 6 {
		t.Fatalf("%d parses, want 6: three first, then a, b, and d", k.Parses())
	}
}

func TestKeptFilesSeesAFileRenamedOverAnother(t *testing.T) {
	dir := t.TempDir()
	writeNote(t, dir, "a.md", "old", longAgo)
	k := newNotes()
	texts(t, k, dir, func(name string) bool { return name == "a.md" })
	// A new file of the same size and time takes the old one's place, as an
	// atomic write does: only the file's identity tells them apart.
	writeNote(t, dir, "tmp", "new", longAgo)
	if err := os.Rename(filepath.Join(dir, "tmp"), filepath.Join(dir, "a.md")); err != nil {
		t.Fatal(err)
	}
	if got := texts(t, k, dir, func(name string) bool { return name == "a.md" }); got["a.md"] != "new" {
		t.Fatalf("after the rename a.md = %q, want new", got["a.md"])
	}
}

func TestKeptFilesParsesARecentFileAgainUntilItIsOlderThanTheWindow(t *testing.T) {
	dir := t.TempDir()
	path := writeNote(t, dir, "a.md", "one", time.Now())
	k := newNotes()
	texts(t, k, dir, nil)
	texts(t, k, dir, nil)
	if k.Parses() != 2 {
		t.Fatalf("%d parses of a file changed just now, want one each list", k.Parses())
	}
	if err := os.Chtimes(path, longAgo, longAgo); err != nil {
		t.Fatal(err)
	}
	texts(t, k, dir, nil)
	texts(t, k, dir, nil)
	if k.Parses() != 3 {
		t.Fatalf("%d parses, want 3: once more once the file is older than the window, then kept", k.Parses())
	}
}

func TestKeptFilesKeepsNoParseError(t *testing.T) {
	dir := t.TempDir()
	writeNote(t, dir, "a.md", "bad one", longAgo)
	k := newNotes()
	for range 2 {
		files, err := k.List(dir, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 1 || files[0].Err == nil || files[0].Value != nil || !files[0].Parsed {
			t.Fatalf("list = %+v, want a.md with its parse error", files)
		}
	}
	if k.Parses() != 2 {
		t.Fatalf("%d parses, want 2: an error is parsed again", k.Parses())
	}
	writeNote(t, dir, "a.md", "good one", longAgo)
	if got := texts(t, k, dir, nil); got["a.md"] != "good one" {
		t.Fatalf("after the fix a.md = %q", got["a.md"])
	}
}

func TestKeptFilesListsOnlyWhatMatchAccepts(t *testing.T) {
	dir := t.TempDir()
	writeNote(t, dir, "TH-0001-a.md", "one", longAgo)
	writeNote(t, dir, "README.md", "index", longAgo)
	if err := os.Mkdir(filepath.Join(dir, "TH-0002-folder.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	k := newNotes()
	all := texts(t, k, dir, nil)
	if !reflect.DeepEqual(all, map[string]string{"TH-0001-a.md": "one", "README.md": "index"}) {
		t.Fatalf("all = %v, want the two regular files", all)
	}
	threads := func(name string) bool { return strings.HasPrefix(name, "TH-") }
	if got := texts(t, k, dir, threads); !reflect.DeepEqual(got, map[string]string{"TH-0001-a.md": "one"}) {
		t.Fatalf("matched = %v", got)
	}
	// README.md was dropped: listing all again parses it anew.
	before := k.Parses()
	texts(t, k, dir, nil)
	if k.Parses() != before+1 {
		t.Fatalf("%d parses, want %d: the file match left out is dropped", k.Parses(), before+1)
	}
}

func TestKeptFilesMissingFolderHasNone(t *testing.T) {
	k := newNotes()
	files, err := k.List(filepath.Join(t.TempDir(), "gone"), nil)
	if err != nil || files != nil {
		t.Fatalf("List of a missing folder = %v, %v; want nil, nil", files, err)
	}
}

// A folder's path that names a file holds no files, as a missing folder
// does, so a reader answers nothing rather than failing, as filepath.Glob
// did before the store.
func TestKeptFilesPathThatIsAFileHasNone(t *testing.T) {
	path := writeNote(t, t.TempDir(), "messages", "not a folder", longAgo)
	k := newNotes()
	files, err := k.List(path, nil)
	if err != nil || files != nil {
		t.Fatalf("List of a path that is a file = %v, %v; want nil, nil", files, err)
	}
}

func TestKeptFilesReadSeesChangesAndAGoneFile(t *testing.T) {
	dir := t.TempDir()
	path := writeNote(t, dir, "a.md", "one", longAgo)
	k := newNotes()
	for range 2 {
		n, err := k.Read(path)
		if err != nil || n.Text != "one" {
			t.Fatalf("Read = %v, %v", n, err)
		}
	}
	if k.Parses() != 1 {
		t.Fatalf("%d parses, want 1", k.Parses())
	}
	// Read and List share what is kept, whatever form the path takes.
	if _, err := k.List(dir+string(filepath.Separator)+".", nil); err != nil || k.Parses() != 1 {
		t.Fatalf("List after Read: %v, %d parses, want 1", err, k.Parses())
	}
	writeNote(t, dir, "a.md", "one two", longAgo)
	if n, err := k.Read(path); err != nil || n.Text != "one two" {
		t.Fatalf("Read after the change = %v, %v", n, err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := k.Read(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Read of a gone file: %v, want os.ErrNotExist", err)
	}
}

func TestKeptFilesHandsOutCopies(t *testing.T) {
	dir := t.TempDir()
	path := writeNote(t, dir, "a.md", "one two", longAgo)
	k := newNotes()
	files, err := k.List(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	files[0].Value.Text = "changed"
	files[0].Value.Words[0] = "changed"
	n, err := k.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if n.Text != "one two" || n.Words[0] != "one" {
		t.Fatalf("a caller's change reached the kept value: %+v", *n)
	}
}

func TestKeptFilesReadersAtOnceShareOneParse(t *testing.T) {
	dir := t.TempDir()
	writeNote(t, dir, "a.md", "one", longAgo)
	writeNote(t, dir, "b.md", "two", longAgo)
	k := newNotes()
	const readers = 8
	got := make([]map[string]string, readers)
	errs := make([]error, readers)
	var wg sync.WaitGroup
	for i := range readers {
		wg.Go(func() {
			files, err := k.List(dir, nil)
			errs[i] = err
			got[i] = map[string]string{}
			for _, f := range files {
				got[i][filepath.Base(f.Path)] = f.Value.Text
			}
		})
	}
	wg.Wait()
	for i := range readers {
		if errs[i] != nil {
			t.Fatal(errs[i])
		}
		if !reflect.DeepEqual(got[i], got[0]) || len(got[i]) != 2 {
			t.Fatalf("reader %d got %v, reader 0 %v", i, got[i], got[0])
		}
	}
	if k.Parses() != 2 {
		t.Fatalf("%d parses, want 2: readers at once share one parse of each file", k.Parses())
	}
}
