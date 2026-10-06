package mcpserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// criteriaFixture is setup with the story's criteria three boxes, the second
// ticked, and a box in its notes that is not a criterion.
func criteriaFixture(t *testing.T) (*fixture, string) {
	t.Helper()
	f := setup(t)
	data, err := os.ReadFile(f.story.Path)
	if err != nil {
		t.Fatal(err)
	}
	s := strings.Replace(string(data), "- [x] works\n", "- [ ] one\n- [x] two\n- [ ] three\n", 1)
	s = strings.TrimRight(s, "\n") + "\n\n- [ ] a note, not a criterion\n"
	if !strings.Contains(s, "- [ ] three\n") || !strings.Contains(s, "## Notes") {
		t.Fatalf("the fixture's story is not as expected:\n%s", s)
	}
	if err := os.WriteFile(f.story.Path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	return f, s
}

// criteria is the criteria a criteria_tick answer lists, as "n:ticked:text".
func criteria(t *testing.T, out map[string]any) string {
	t.Helper()
	list, ok := out["criteria"].([]any)
	if !ok {
		t.Fatalf("criteria is not a list: %v", out["criteria"])
	}
	var got []string
	for _, c := range list {
		m := c.(map[string]any)
		mark := " "
		if m["ticked"] == true {
			mark = "x"
		}
		got = append(got, fmt.Sprintf("%v:%s:%v", m["n"], mark, m["text"]))
	}
	return strings.Join(got, ",")
}

// S-0282: criteria_tick ticks and unticks a story's criteria by number, and
// changes those boxes and nothing else, the box in the notes left alone;
// nothing is committed, and the answer lists the criteria after the write.
func TestCriteriaTickTicksAndUnticksByNumber(t *testing.T) {
	f, before := criteriaFixture(t)
	got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID})
	hash := got["hash"].(string)

	out, failed := f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "hash": hash, "tick": []int{1, 3}, "untick": []int{2}})
	if failed != "" {
		t.Fatal(failed)
	}
	if c := criteria(t, out); c != "1:x:one,2: :two,3:x:three" {
		t.Errorf("criteria after: %s", c)
	}
	if out["id"] != f.story.ID || out["hash"] == hash || len(out["hash"].(string)) != 64 || out["unchanged"] == true ||
		strings.Join(toStrings(out["changed"]), ",") != "criteria" || out["path"] != "wip/kanban/stories/"+filepath.Base(f.story.Path) {
		t.Errorf("answer: %v", out)
	}
	data, _ := os.ReadFile(f.story.Path)
	after := string(data)
	want := strings.Replace(before, "- [ ] one\n- [x] two\n- [ ] three\n", "- [x] one\n- [ ] two\n- [x] three\n", 1)
	// the edit stamps updated:, which is front matter flai owns; the body
	// below it is byte for byte what was asked
	if i, j := strings.Index(after, "\n# "), strings.Index(want, "\n# "); i < 0 || j < 0 || after[i:] != want[j:] {
		t.Errorf("the file:\n%s\nwant below the front matter:\n%s", after, want)
	}
	if !strings.Contains(after, "- [ ] a note, not a criterion\n") {
		t.Errorf("the box in the notes changed:\n%s", after)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID}); got["hash"] != out["hash"] {
		t.Errorf("the hash returned is not the file's: %v, item_get %v", out["hash"], got["hash"])
	}

	// unticking with the new hash, and ticking what is ticked, changes only
	// the one box
	out, failed = f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "hash": out["hash"], "tick": []int{1}, "untick": []int{3}})
	if failed != "" {
		t.Fatal(failed)
	}
	if c := criteria(t, out); c != "1:x:one,2: :two,3: :three" {
		t.Errorf("criteria after the untick: %s", c)
	}
	// a write that changes nothing says so and lists the criteria as they are
	out, failed = f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "tick": []int{1}})
	if failed != "" || out["unchanged"] != true || len(toStringsOrNone(out["changed"])) != 0 || criteria(t, out) != "1:x:one,2: :two,3: :three" {
		t.Errorf("an unchanged tick: %v %s", out, failed)
	}
}

// S-0282: a number with no box, a stale hash, and a closed or archived item
// are refused, and nothing is written.
func TestCriteriaTickRefusesWithoutWriting(t *testing.T) {
	f, before := criteriaFixture(t)
	got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID})
	hash := got["hash"].(string)
	unchanged := func(what string) {
		t.Helper()
		if data, _ := os.ReadFile(f.story.Path); string(data) != before {
			t.Errorf("%s wrote the file:\n%s", what, data)
		}
	}

	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "hash": hash, "tick": []int{1, 4}}); !strings.Contains(failed, "there is no acceptance criterion 4") {
		t.Errorf("out of range: %q", failed)
	}
	unchanged("a number out of range")
	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "tick": []int{0}}); !strings.Contains(failed, "no acceptance criterion 0") {
		t.Errorf("zero: %q", failed)
	}
	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "tick": []int{1}, "untick": []int{1}}); !strings.Contains(failed, "both to tick and to untick") {
		t.Errorf("both: %q", failed)
	}
	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": f.story.ID}); !strings.Contains(failed, "name the acceptance criteria") {
		t.Errorf("none named: %q", failed)
	}
	unchanged("a refused tick")

	// changed meanwhile: the stale hash is refused, even with numbers that
	// would be out of range in the file now
	edited := strings.Replace(before, "- [ ] three\n", "- [ ] three\n- [ ] four\n", 1)
	if err := os.WriteFile(f.story.Path, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	before = edited
	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": f.story.ID, "hash": hash, "tick": []int{3}}); !strings.Contains(failed, "changed after it was loaded") {
		t.Errorf("a stale hash: %q", failed)
	}
	unchanged("a stale hash")

	// a closed story, and an archived one
	closed, _ := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Closed", Owner: "alex", Now: t0, Body: "## Goal\n\nG.\n\n## Acceptance criteria\n- [ ] it works\n\n## Tasks\n\n## Notes\n"})
	if _, err := f.repo.Transition(closed, workitem.Cancelled, "alex", "not needed", t0); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": closed.ID, "tick": []int{1}}); !strings.Contains(failed, "a closed item is not edited") {
		t.Errorf("a closed story: %q", failed)
	}
	archived := filepath.Join(f.repo.Root, "wip/archive/kanban/stories/S-0900-done.md")
	content := "---\nid: S-0900\ntype: story\nnature: feature\ntitle: Done\nstatus: done\n---\n# S-0900 Done\n\n## Goal\n\nG.\n\n## Acceptance criteria\n- [ ] it works\n"
	_ = os.MkdirAll(filepath.Dir(archived), 0o755)
	if err := os.WriteFile(archived, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "criteria_tick", map[string]any{"id": "S-0900", "tick": []int{1}}); !strings.Contains(failed, "archived") {
		t.Errorf("an archived story: %q", failed)
	}
	if data, _ := os.ReadFile(archived); string(data) != content {
		t.Errorf("the archived story was written:\n%s", data)
	}
}

func toStringsOrNone(v any) []string {
	if v == nil {
		return nil
	}
	return toStrings(v)
}
