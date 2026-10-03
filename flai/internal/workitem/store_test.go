package workitem

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// parses is how many files the store has parsed in dir.
func parses(dir string) int {
	f := parsed.folder(dir)
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.parses
}

// totalParses is parses across the project's item folders, archive included.
func totalParses(r *Repo) int {
	n := 0
	for _, archived := range []bool{false, true} {
		for _, typ := range Types {
			n += parses(r.ItemDir(typ, archived))
		}
	}
	return n
}

// age sets every item file's modification time to at, as a file written
// long before it is read.
func age(t *testing.T, r *Repo, at time.Time) {
	t.Helper()
	for _, archived := range []bool{false, true} {
		for _, typ := range Types {
			files, _ := filepath.Glob(filepath.Join(r.ItemDir(typ, archived), "*.md"))
			for _, f := range files {
				if err := os.Chtimes(f, at, at); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

func listed(t *testing.T, r *Repo, archive bool) map[string]*Item {
	t.Helper()
	items, err := r.List(archive)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]*Item{}
	for _, it := range items {
		out[it.ID] = it
	}
	return out
}

// sameSizeEdit rewrites the file with one byte of its title changed, so that
// only its content and modification time tell the edit apart.
func sameSizeEdit(t *testing.T, path, from, to string) {
	t.Helper()
	if len(from) != len(to) {
		t.Fatal("edit must keep the size")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	edited := strings.Replace(string(data), from, to, 1)
	if edited == string(data) {
		t.Fatalf("%q not in %s", from, path)
	}
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteAt([]byte(edited), 0); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestListReadsAgainOnlyTheFilesThatChanged(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "Epic", "")
	mustCreate(t, r, Story, "Story one", "E-0001")
	mustCreate(t, r, Task, "Task one", "S-0001")
	s, err := r.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	s.Body = strings.Replace(s.Body, "- [ ]\n", "- [x] ok\n", 1)
	if err := r.Save(s); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-time.Hour)
	age(t, r, past)
	listed(t, r, true)
	before := totalParses(r)
	listed(t, r, true)
	if got := totalParses(r) - before; got != 0 {
		t.Fatalf("a second list of unchanged files parsed %d", got)
	}

	// created by flai
	mustCreate(t, r, Story, "Story two", "E-0001")
	if got := listed(t, r, false); got["S-0002"] == nil || got["S-0002"].Title != "Story two" {
		t.Fatalf("created story not listed: %v", got)
	}
	age(t, r, past.Add(time.Minute))
	listed(t, r, true)
	before = totalParses(r)
	listed(t, r, true)
	if got := totalParses(r) - before; got != 0 {
		t.Fatalf("unchanged files parsed again: %d", got)
	}

	// moved by flai
	s, err = r.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	mustMove(t, r, s, Ready, "")
	if got := listed(t, r, false)["S-0001"].Status; got != Ready {
		t.Errorf("moved story lists as %s", got)
	}
	if got, _ := r.Get("S-0001"); got.Status != Ready {
		t.Errorf("moved story gets as %s", got.Status)
	}

	// edited by hand, same size, a later time
	age(t, r, past.Add(2*time.Minute))
	listed(t, r, true)
	path := s.Path
	sameSizeEdit(t, path, "Story one", "Story One")
	later := past.Add(3 * time.Minute)
	if err := os.Chtimes(path, later, later); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, false)["S-0001"].Title; got != "Story One" {
		t.Errorf("hand-edited title lists as %q", got)
	}
	if got, _ := r.Get("S-0001"); got.Title != "Story One" {
		t.Errorf("hand-edited title gets as %q", got.Title)
	}

	// created by hand
	data, _ := os.ReadFile(path)
	other := strings.Replace(strings.Replace(string(data), "id: S-0001", "id: S-0009", 1), "Story One", "Story nine", 1)
	handmade := filepath.Join(r.ItemDir(Story, false), "S-0009-story-nine.md")
	if err := os.WriteFile(handmade, []byte(other), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, false)["S-0009"]; got == nil || got.Title != "Story nine" {
		t.Errorf("hand-made story not listed: %v", got)
	}

	// deleted by hand
	if err := os.Remove(handmade); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, false)["S-0009"]; got != nil {
		t.Error("deleted story still listed")
	}
	if _, err := r.Get("S-0009"); err == nil {
		t.Error("deleted story still found")
	}

	// archived by flai, and archived by hand
	task, _ := r.Get("T-0001")
	for _, to := range []string{Ready, InProgress, Done} {
		mustMove(t, r, task, to, "")
	}
	items, _ := r.List(false)
	plan, err := r.PlanArchive(items, []string{"T-0001"})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Archive(plan); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, false)["T-0001"]; got != nil {
		t.Error("archived task still listed as open")
	}
	if got := listed(t, r, true)["T-0001"]; got == nil || !got.Archived {
		t.Errorf("archived task not listed as archived: %v", got)
	}
	epic, _ := r.Get("E-0001")
	if err := os.MkdirAll(r.ItemDir(Epic, true), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(epic.Path, filepath.Join(r.ItemDir(Epic, true), filepath.Base(epic.Path))); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, true)["E-0001"]; got == nil || !got.Archived {
		t.Errorf("epic archived by hand not listed as archived: %v", got)
	}
	if got, _ := r.Get("E-0001"); got == nil || !got.Archived {
		t.Errorf("epic archived by hand not found archived: %v", got)
	}
}

// A write in the same timestamp tick as the read before it, keeping the size,
// leaves the time and size as they were: the store reads such a file again
// until its time is well before the read.
func TestListReadsAgainAFileWrittenTooRecentlyToTell(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "Epic", "")
	e, _ := r.Get("E-0001")
	tick := time.Now().Add(-time.Second / 2)
	if err := os.Chtimes(e.Path, tick, tick); err != nil {
		t.Fatal(err)
	}
	listed(t, r, false)
	sameSizeEdit(t, e.Path, "Epic", "Eric")
	if err := os.Chtimes(e.Path, tick, tick); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, false)["E-0001"].Title; got != "Eric" {
		t.Errorf("a same-size write in the same tick lists as %q", got)
	}
}

// A file renamed over another, as git and editors do, is another file even
// with the same time and size.
func TestListReadsAgainAFileRenamedOverAnother(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "Epic", "")
	e, _ := r.Get("E-0001")
	past := time.Now().Add(-time.Hour)
	age(t, r, past)
	listed(t, r, false)
	data, _ := os.ReadFile(e.Path)
	tmp := filepath.Join(r.Root, "epic.tmp")
	if err := os.WriteFile(tmp, []byte(strings.Replace(string(data), "title: Epic", "title: Eric", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(tmp, past, past); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp, e.Path); err != nil {
		t.Fatal(err)
	}
	if got := listed(t, r, false)["E-0001"].Title; got != "Eric" {
		t.Errorf("a file renamed over lists as %q", got)
	}
}

func TestListsAtOnceShareOneRead(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "Epic", "")
	for range 20 {
		mustCreate(t, r, Story, "Story", "E-0001")
	}
	age(t, r, time.Now().Add(-time.Hour)) // every file changed since it was last read
	before := totalParses(r)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			if items, err := r.List(true); err != nil || len(items) != 21 {
				t.Errorf("list: %d items, %v", len(items), err)
			}
		})
	}
	wg.Wait()
	if got := totalParses(r) - before; got != 21 {
		t.Errorf("8 lists at once parsed %d files, want each of 21 once", got)
	}
}

func TestAListedItemIsTheCallersOwn(t *testing.T) {
	r := newProject(t)
	e := mustCreate(t, r, Epic, "Epic", "")
	e.Transitions = []Transition{{To: Backlog, At: "2026-09-15T20:00:00Z", By: "test"}}
	if err := r.Save(e); err != nil {
		t.Fatal(err)
	}
	age(t, r, time.Now().Add(-time.Hour))
	first := listed(t, r, false)["E-0001"]
	first.Title = "changed"
	first.Transitions[0].To = "changed"
	first.Tags = append(first.Tags, "changed")
	again := listed(t, r, false)["E-0001"]
	if again.Title != "Epic" || again.Transitions[0].To == "changed" || len(again.Tags) != 0 {
		t.Errorf("a caller's change reached the next list: %+v", again)
	}
}

// clone shares nothing a caller could change: every slice, map, and pointer
// of a copy is its own, whatever fields Item gains.
func TestCloneSharesNothing(t *testing.T) {
	it := &Item{
		ID: "S-0001", Transitions: []Transition{{To: Ready}}, Blocked: []Block{{Reason: "r"}},
		Tags: []string{"t"}, Touches: []string{"p"}, Topics: []string{"x"}, After: []string{"S-0002"},
		Agent:   &manifest.Agent{Model: "m", Config: map[string]string{"k": "v"}, Roles: map[string]manifest.Role{"verify": {Model: "s", Config: map[string]string{"k": "v"}}}},
		Usage:   &usage.Usage{Models: []usage.Model{{Model: "m"}}},
		Unknown: []Field{{Name: "hold", Raw: "hold: x\n"}},
		CostOfDelay: &CostOfDelay{
			Inputs: &CostInputs{RevenuePerWeek: new(float64), PenaltyPerWeek: new(float64)},
			Value:  new(float64),
		},
		Forecast:  &Forecast{Duration: "4h"},
		Finalized: &Finalized{By: "alex"},
	}
	c := it.clone()
	if !reflect.DeepEqual(it, c) {
		t.Fatalf("clone differs: %+v", c)
	}
	var shared func(path string, a, b reflect.Value)
	shared = func(path string, a, b reflect.Value) {
		switch a.Kind() {
		case reflect.Pointer:
			if a.IsNil() {
				t.Errorf("%s: fill it in this test", path)
				return
			}
			if a.Pointer() == b.Pointer() {
				t.Errorf("%s is shared", path)
			}
			shared(path, a.Elem(), b.Elem())
		case reflect.Slice, reflect.Map:
			if a.Len() == 0 {
				t.Errorf("%s: fill it in this test", path)
				return
			}
			if a.Pointer() == b.Pointer() {
				t.Errorf("%s is shared", path)
			}
		case reflect.Struct:
			for i := range a.NumField() {
				shared(path+"."+a.Type().Field(i).Name, a.Field(i), b.Field(i))
			}
		}
	}
	shared("Item", reflect.ValueOf(it), reflect.ValueOf(c))
}
