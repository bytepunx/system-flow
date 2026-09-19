package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// S-0057: flai order places a story in the pull order, flai board reads the
// columns in it, and a story that becomes ready joins the ready section.
func TestOrderCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	run := func(args ...string) string {
		t.Helper()
		out, errOut, code := runIn(t, root, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		return out
	}
	run("epic", "new", "Epic")
	for _, title := range []string{"One", "Two", "Three", "Four"} {
		run("story", "new", title, "--epic", "E-0001")
	}
	order := func() []string {
		var v struct {
			Order []string `json:"order"`
		}
		if err := json.Unmarshal([]byte(run("board", "--json")), &v); err != nil {
			t.Fatal(err)
		}
		return v.Order
	}
	column := func(state string) (ids []string) {
		var v struct {
			Columns map[string][]struct {
				ID string `json:"id"`
			} `json:"columns"`
		}
		if err := json.Unmarshal([]byte(run("board", "--json")), &v); err != nil {
			t.Fatal(err)
		}
		for _, c := range v.Columns[state] {
			ids = append(ids, c.ID)
		}
		return
	}

	if got, want := column("backlog"), []string{"S-0001", "S-0002", "S-0003", "S-0004"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("an unplaced backlog reads by ID: got %v want %v", got, want)
	}
	out := run("order", "S-0004", "--before", "S-0002")
	if !strings.Contains(out, "backlog: S-0001, S-0004, S-0002, S-0003") {
		t.Errorf("the command prints the column's new sequence, got %q", out)
	}
	if got, want := column("backlog"), []string{"S-0001", "S-0004", "S-0002", "S-0003"}; !reflect.DeepEqual(got, want) {
		t.Errorf("flai board agrees: got %v want %v", got, want)
	}
	if got, want := order(), []string{"S-0001", "S-0004"}; !reflect.DeepEqual(got, want) {
		t.Errorf("board.md names only as far as it must: got %v want %v", got, want)
	}
	data, _ := os.ReadFile(filepath.Join(root, "wip/kanban/board.md"))
	if !strings.Contains(string(data), "order:\n  - S-0001\n  - S-0004\n") {
		t.Errorf("board.md front matter:\n%s", data)
	}

	// Becoming ready: after the ready stories, before the backlog names.
	for _, id := range []string{"S-0003", "S-0004"} {
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		_ = os.WriteFile(file[0], []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
		run("move", id, "ready")
	}
	if got, want := order(), []string{"S-0003", "S-0004", "S-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ready names lead the list: got %v want %v", got, want)
	}
	run("order", "S-0004", "--top")
	if got, want := column("ready"), []string{"S-0004", "S-0003"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ready reordered: got %v want %v", got, want)
	}

	// Refusals leave board.md alone and say why.
	before, _ := os.ReadFile(filepath.Join(root, "wip/kanban/board.md"))
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"order", "S-0004", "--before", "S-0001"}, "within one column"},
		{[]string{"order", "E-0001", "--top"}, "only stories"},
		{[]string{"order", "S-0004"}, "exactly one of"},
		{[]string{"order", "S-0004", "--top", "--bottom"}, "exactly one of"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("flai %v: code %d, stderr %q, want %q", c.args, code, errOut, c.want)
		}
	}
	after, _ := os.ReadFile(filepath.Join(root, "wip/kanban/board.md"))
	if string(before) != string(after) {
		t.Error("a refused placement rewrote board.md")
	}

	var j struct {
		ID       string   `json:"id"`
		Status   string   `json:"status"`
		Sequence []string `json:"sequence"`
		Order    []string `json:"order"`
	}
	if err := json.Unmarshal([]byte(run("order", "S-0003", "--top", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	if j.ID != "S-0003" || j.Status != "ready" || !reflect.DeepEqual(j.Sequence, []string{"S-0003", "S-0004"}) || !reflect.DeepEqual(j.Order, []string{"S-0003", "S-0004", "S-0001"}) {
		t.Errorf("--json: %+v", j)
	}
}
