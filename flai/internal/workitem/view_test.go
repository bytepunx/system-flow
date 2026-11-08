package workitem

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// A story with three tasks, one of them archived, and a second story whose
// task is not among the first's children.
func viewProject(t *testing.T) *Repo {
	t.Helper()
	r := newProject(t)
	mustCreate(t, r, Story, "First", "")
	mustCreate(t, r, Story, "Second", "")
	for _, title := range []string{"a", "b", "c"} {
		mustCreate(t, r, Task, title, "S-0001")
	}
	mustCreate(t, r, Task, "other", "S-0002")
	// T-0002 is archived: its file moves to the archive's tasks folder
	from, err := filepath.Glob(filepath.Join(r.ItemDir(Task, false), "T-0002-*.md"))
	if err != nil || len(from) != 1 {
		t.Fatalf("T-0002's file: %v %v", from, err)
	}
	dest := filepath.Join(r.ItemDir(Task, true), filepath.Base(from[0]))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(from[0], dest); err != nil {
		t.Fatal(err)
	}
	age(t, r, time.Now().Add(-time.Hour))
	return r
}

func TestChildrenOfIsChildrenOfTheList(t *testing.T) {
	r := viewProject(t)
	for _, archive := range []bool{false, true} {
		all, err := r.List(archive)
		if err != nil {
			t.Fatal(err)
		}
		want := Children(all, "S-0001")
		got, err := r.ChildrenOf("S-0001", archive)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("archive %v: ChildrenOf gave %v, Children of the list %v", archive, itemIDs(got), itemIDs(want))
		}
	}
	got, _ := r.ChildrenOf("S-0001", true)
	if w := []string{"T-0001", "T-0002", "T-0003"}; !reflect.DeepEqual(itemIDs(got), w) {
		t.Errorf("children %v, want %v in ID order", itemIDs(got), w)
	}
	if !got[1].Archived || got[0].Archived {
		t.Errorf("Archived not set as the list sets it: %v %v", got[0].Archived, got[1].Archived)
	}
}

func TestChildrenOfHandsOutCopies(t *testing.T) {
	r := viewProject(t)
	first, _ := r.ChildrenOf("S-0001", true)
	first[0].Title = "changed"
	first[0].Tags = append(first[0].Tags, "changed")
	again, _ := r.ChildrenOf("S-0001", true)
	if again[0].Title != "a" || len(again[0].Tags) != 0 {
		t.Errorf("a caller's change reached the next call: %+v", again[0])
	}
}

func TestVisitSeesEveryItemOnce(t *testing.T) {
	r := viewProject(t)
	for _, archive := range []bool{false, true} {
		all, err := r.List(archive)
		if err != nil {
			t.Fatal(err)
		}
		want := map[string]bool{}
		for _, it := range all {
			want[it.ID] = it.Archived
		}
		got := map[string]bool{}
		err = r.Visit(archive, func(it *Item, archived bool) {
			if _, twice := got[it.ID]; twice {
				t.Errorf("%s visited twice", it.ID)
			}
			got[it.ID] = archived
		})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("archive %v: visited %v, listed %v", archive, got, want)
		}
	}
}

func TestLessIDOrdersAsListDoes(t *testing.T) {
	cases := []struct {
		a, b string
		less bool
	}{
		{"E-0009", "S-0001", true},
		{"S-0001", "T-0001", true},
		{"S-002", "S-0010", true},
		{"S-0010", "S-002", false},
		{"T-0001", "T-0001", false},
	}
	for _, c := range cases {
		if got := LessID(c.a, c.b); got != c.less {
			t.Errorf("LessID(%s, %s) = %v, want %v", c.a, c.b, got, c.less)
		}
	}
}

func itemIDs(items []*Item) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.ID)
	}
	return out
}
