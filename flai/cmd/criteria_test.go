package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// criteriaProject is editProject with three criteria on S-0001, the second
// ticked.
func criteriaProject(t *testing.T) (root, file string) {
	t.Helper()
	root = editProject(t)
	file = filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md")
	s := read(t, file)
	_ = os.WriteFile(file, []byte(strings.Replace(s, "- [ ] it works\n", "- [ ] it works\n- [x] it is fast\n  - [ ] it is small\n", 1)), 0o644)
	gitIn(t, root, "commit", "-qam", "criteria")
	return root, file
}

// criteriaJSON is what flai criteria list --json and tick --json print.
type criteriaJSON struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Hash     string `json:"hash"`
	Changed  []string
	Criteria []struct {
		N      int    `json:"n"`
		Text   string `json:"text"`
		Ticked bool   `json:"ticked"`
	} `json:"criteria"`
}

// S-0282: flai criteria list numbers the boxes as tick takes them.
func TestCriteriaList(t *testing.T) {
	root, _ := criteriaProject(t)
	out, errOut, code := runIn(t, root, "criteria", "list", "S-0001")
	want := "S-0001: 1 of 3 ticked\n  1 [ ] it works\n  2 [x] it is fast\n  3 [ ] it is small\n"
	if code != 0 || out != want {
		t.Fatalf("list: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "criteria", "list", "S-0001", "--json")
	var v criteriaJSON
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil {
		t.Fatalf("list --json: %d %s %s", code, out, errOut)
	}
	if v.ID != "S-0001" || v.Path != "wip/kanban/stories/S-0001-a-plain-story.md" || v.Hash == "" || len(v.Criteria) != 3 || v.Criteria[1].N != 2 || !v.Criteria[1].Ticked || v.Criteria[2].Text != "it is small" {
		t.Errorf("list --json: %+v", v)
	}
	// none is an empty list, not null
	if out, _, _ := runIn(t, root, "criteria", "list", "E-0001", "--json"); !strings.Contains(out, `"criteria": []`) {
		t.Errorf("an epic without boxes: %s", out)
	}
}

// S-0282: tick and untick change the one box, report it, and commit the item
// alone with --autocommit; ticking what is ticked changes nothing.
func TestCriteriaTickAndUntick(t *testing.T) {
	root, file := criteriaProject(t)
	before := read(t, file)
	out, errOut, code := runIn(t, root, "criteria", "tick", "S-0001", "1,3", "--autocommit", "--trailer", "Co-Authored-By: someone <s@s>")
	if code != 0 || !strings.HasPrefix(out, "S-0001: ticked 1, 3\n  committed ") || !strings.HasSuffix(out, "S-0001: 3 of 3 ticked\n  1 [x] it works\n  2 [x] it is fast\n  3 [x] it is small\n") {
		t.Fatalf("tick: %d %q %s", code, out, errOut)
	}
	after := read(t, file)
	want := strings.Replace(strings.Replace(before, "- [ ] it works", "- [x] it works", 1), "- [ ] it is small", "- [x] it is small", 1)
	if below(after) != below(want) {
		t.Errorf("only the boxes change:\n%s", after)
	}
	if msg := gitIn(t, root, "log", "-1", "--format=%B"); !strings.HasPrefix(msg, "chore: [S-0001] tick criteria 1, 3") || !strings.Contains(msg, "Co-Authored-By: someone") {
		t.Errorf("commit message: %s", msg)
	}
	if out, _, code := runIn(t, root, "criteria", "tick", "S-0001", "2"); code != 0 || !strings.HasPrefix(out, "S-0001: unchanged\n") {
		t.Errorf("ticking a ticked box: %d %s", code, out)
	}
	out, errOut, code = runIn(t, root, "criteria", "untick", "S-0001", "2", "--json")
	var v criteriaJSON
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil || strings.Join(v.Changed, ",") != "criteria" || len(v.Criteria) != 3 || v.Criteria[1].Ticked || !v.Criteria[0].Ticked {
		t.Fatalf("untick --json: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(read(t, file), "- [ ] it is fast\n") {
		t.Errorf("untick: %s", read(t, file))
	}
}

// below is a file from its heading on, which a write leaves as it was but
// for what it changes.
func below(file string) string {
	return file[strings.Index(file, "\n# "):]
}

// S-0282: what the caller gets wrong is a rule, and nothing is written; a
// stale hash is a conflict; a closed or archived item is refused, and listed.
func TestCriteriaRefusals(t *testing.T) {
	root, file := criteriaProject(t)
	before := read(t, file)
	for name, c := range map[string]struct {
		args []string
		says string
	}{
		"out of range": {[]string{"criteria", "tick", "S-0001", "4"}, "there is no acceptance criterion 4"},
		"zero":         {[]string{"criteria", "tick", "S-0001", "1,0"}, `0\" is not the number`},
		"not a number": {[]string{"criteria", "untick", "S-0001", "1", "two"}, `two\" is not the number`},
		"no number":    {[]string{"criteria", "tick", "S-0001"}, "name the acceptance criteria"},
		"no boxes":     {[]string{"criteria", "tick", "E-0001", "1"}, "there are no acceptance criteria"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, "rule: ") || !strings.Contains(errOut, c.says) {
			t.Errorf("%s: %d %s", name, code, errOut)
		}
	}
	if read(t, file) != before {
		t.Error("no refusal wrote anything")
	}

	var v criteriaJSON
	out, _, _ := runIn(t, root, "criteria", "list", "S-0001", "--json")
	_ = json.Unmarshal([]byte(out), &v)
	if _, errOut, code := runIn(t, root, "criteria", "tick", "S-0001", "1", "--hash", v.Hash); code != 0 {
		t.Fatalf("tick with the current hash: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "criteria", "untick", "S-0001", "1", "--hash", v.Hash); code != exitDocConflict || !strings.Contains(errOut, "changed after it was loaded") {
		t.Errorf("a stale hash: %d %s", code, errOut)
	}
	if !strings.Contains(read(t, file), "- [x] it works") {
		t.Error("the conflict wrote nothing")
	}

	if _, errOut, code := runIn(t, root, "move", "S-0001", "cancelled", "--reason", "test"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "criteria", "untick", "S-0001", "1"); code != exitDocRefused || !strings.Contains(errOut, "cancelled") {
		t.Errorf("a closed item: %d %s", code, errOut)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-qm", "cancelled")
	if _, errOut, code := runIn(t, root, "archive"); code != 0 {
		t.Fatalf("archive: %s", errOut)
	}
	if out, errOut, code := runIn(t, root, "criteria", "list", "S-0001"); code != 0 || !strings.HasPrefix(out, "S-0001: 2 of 3 ticked\n") {
		t.Errorf("an archived item is listed: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "criteria", "tick", "S-0001", "3"); code != exitDocRefused || !strings.Contains(errOut, "archived") {
		t.Errorf("an archived item: %d %s", code, errOut)
	}
}
