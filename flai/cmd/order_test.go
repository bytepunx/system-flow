package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// S-0217: an order computed by a policy is printed with its figures, and
// only --apply writes it to board.md.
func TestOrderByCommand(t *testing.T) {
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
	for _, title := range []string{"One", "Two", "Three"} {
		run("story", "new", title, "--epic", "E-0001")
	}
	run("edit", "S-0001", "--cost-of-delay-value", "100", "--forecast-duration", "1h")
	run("edit", "S-0003", "--cost-of-delay-value", "500", "--forecast-duration", "10h")
	for _, id := range []string{"S-0001", "S-0002", "S-0003"} {
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		_ = os.WriteFile(file[0], []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
		run("move", id, "ready")
	}
	boardPath := filepath.Join(root, "wip/kanban/board.md")
	before, _ := os.ReadFile(boardPath)

	out := run("order", "--by", "cod")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 || lines[0] != "ready by cod:" {
		t.Fatalf("by cod:\n%s", out)
	}
	for i, want := range []string{"1 S-0003 500 Three", "2 S-0001 100 One", "3 S-0002 no cost of delay value Two"} {
		if got := strings.Join(strings.Fields(lines[i+1]), " "); got != want {
			t.Errorf("line %d: got %q want %q", i+1, got, want)
		}
	}
	if after, _ := os.ReadFile(boardPath); string(after) != string(before) {
		t.Error("--by without --apply rewrote board.md")
	}

	var j struct {
		Policy  string `json:"policy"`
		Applied bool   `json:"applied"`
		Stories []struct {
			Position int      `json:"position"`
			ID       string   `json:"id"`
			Figure   *float64 `json:"figure"`
			Missing  string   `json:"missing"`
		} `json:"stories"`
	}
	if err := json.Unmarshal([]byte(run("order", "--by", "wsjf", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	if j.Policy != "wsjf" || j.Applied || len(j.Stories) != 3 ||
		j.Stories[0].ID != "S-0001" || j.Stories[0].Figure == nil || *j.Stories[0].Figure != 100 ||
		j.Stories[1].ID != "S-0003" || j.Stories[1].Figure == nil || *j.Stories[1].Figure != 50 ||
		j.Stories[2].ID != "S-0002" || j.Stories[2].Missing != "cost of delay value and forecast duration" {
		t.Errorf("--json: %+v", j)
	}

	run("order", "--by", "cod", "--apply")
	var b struct {
		Columns map[string][]struct {
			ID string `json:"id"`
		} `json:"columns"`
	}
	if err := json.Unmarshal([]byte(run("board", "--json")), &b); err != nil {
		t.Fatal(err)
	}
	var ready []string
	for _, c := range b.Columns["ready"] {
		ready = append(ready, c.ID)
	}
	if want := []string{"S-0003", "S-0001", "S-0002"}; !reflect.DeepEqual(ready, want) {
		t.Errorf("--apply writes the computed order: got %v want %v", ready, want)
	}
	data, _ := os.ReadFile(boardPath)
	if !strings.Contains(string(data), "order:\n  - S-0003\n  - S-0001\n  - S-0002\n") {
		t.Errorf("board.md front matter:\n%s", data)
	}

	// Refusals leave board.md alone and say why.
	before, _ = os.ReadFile(boardPath)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"order", "--by", "cod", "--top"}, "give one or the other"},
		{[]string{"order", "--by", "cod", "--before", "S-0001"}, "give one or the other"},
		{[]string{"order", "S-0001", "--by", "cod"}, "takes no story"},
		{[]string{"order", "S-0001", "--top", "--apply"}, "give --by too"},
		{[]string{"order", "--by", "random", "--apply"}, "cod, wsjf, throughput, fifo"},
		{[]string{"order", "--top"}, "name the story"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, stderr %q, want %q", c.args, code, errOut, c.want)
		}
	}
	if after, _ := os.ReadFile(boardPath); string(after) != string(before) {
		t.Error("a refused order rewrote board.md")
	}
}

// S-0219: flai order records who placed a story and when, flai move out of
// the column drops the record, and --by --apply keeps the place of a story
// placed by hand within --keep-placed and orders the others around it.
func TestOrderKeepsHandPlacements(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	t.Setenv("FLAI_ROLE", "")
	root := tempProject(t)
	now := time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)
	runAt := func(at time.Time, args ...string) string {
		t.Helper()
		out, errOut, code := runInAt(t, root, at, args...)
		if code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
		return out
	}
	run := func(args ...string) string { t.Helper(); return runAt(now, args...) }
	run("epic", "new", "Epic")
	for i, title := range []string{"One", "Two", "Three", "Four"} {
		run("story", "new", title, "--epic", "E-0001")
		id := fmt.Sprintf("S-%04d", i+1)
		run("edit", id, "--cost-of-delay-value", fmt.Sprint(100*(i+1)))
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		_ = os.WriteFile(file[0], []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)), 0o644)
		run("move", id, "ready")
	}
	boardPath := filepath.Join(root, "wip/kanban/board.md")
	placed := func() map[string]workitem.Placed {
		t.Helper()
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		b, err := repo.LoadBoard()
		if err != nil {
			t.Fatal(err)
		}
		return b.Placed
	}
	ready := func() (ids []string) {
		var v struct {
			Columns map[string][]struct {
				ID string `json:"id"`
			} `json:"columns"`
		}
		if err := json.Unmarshal([]byte(run("board", "--json")), &v); err != nil {
			t.Fatal(err)
		}
		for _, c := range v.Columns["ready"] {
			ids = append(ids, c.ID)
		}
		return
	}

	// Two days ago, by hand; an hour ago, through the dashboard; an hour
	// ago, by the orchestrator, whoever it says it is.
	runAt(now.Add(-48*time.Hour), "order", "S-0002", "--top", "--placed-by", "alex")
	var j struct {
		Placed workitem.Placed `json:"placed"`
	}
	t.Setenv("FLAI_AGENT", "flaiover")
	if err := json.Unmarshal([]byte(runAt(now.Add(-time.Hour), "order", "S-0001", "--after", "S-0002", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	if want := (workitem.Placed{By: "flaiover", At: "2026-09-15T20:00:00Z"}); j.Placed != want {
		t.Errorf("--json names the placement: got %+v want %+v", j.Placed, want)
	}
	t.Setenv("FLAI_ROLE", "orchestrate")
	t.Setenv("FLAI_AGENT", "system-flow")
	runAt(now.Add(-time.Hour), "order", "S-0003", "--bottom", "--placed-by", "alex")
	t.Setenv("FLAI_ROLE", "")
	t.Setenv("FLAI_AGENT", "")
	want := map[string]workitem.Placed{
		"S-0001": {By: "flaiover", At: "2026-09-15T20:00:00Z"},
		"S-0002": {By: "alex", At: "2026-09-13T21:00:00Z"},
		"S-0003": {By: "orchestrator", At: "2026-09-15T20:00:00Z"},
	}
	if got := placed(); !reflect.DeepEqual(got, want) {
		t.Errorf("placed: got %+v want %+v", got, want)
	}
	if got, want := ready(), []string{"S-0002", "S-0001", "S-0004", "S-0003"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("placed by hand: got %v want %v", got, want)
	}

	// With no window, cod alone; with the default, S-0001 stays second.
	var none workitem.ReadyOrder
	if err := json.Unmarshal([]byte(run("order", "--by", "cod", "--keep-placed", "0", "--json")), &none); err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, r := range none.Stories {
		ids = append(ids, r.ID)
		if r.Kept != nil {
			t.Errorf("--keep-placed 0 keeps nothing, but %s: %+v", r.ID, r.Kept)
		}
	}
	if want := []string{"S-0004", "S-0003", "S-0002", "S-0001"}; !reflect.DeepEqual(ids, want) || none.KeepPlaced != "" {
		t.Errorf("--keep-placed 0: got %v (keep %q) want %v", ids, none.KeepPlaced, want)
	}
	out := run("order", "--by", "cod", "--apply")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i, line := range []string{
		"ready by cod:",
		"1 S-0004 400 Four",
		"2 S-0001 100 One (kept: placed by flaiover at 2026-09-15T20:00:00Z)",
		"3 S-0003 300 Three",
		"4 S-0002 200 Two",
		"board.md's ready order is now this one",
	} {
		if i >= len(lines) || strings.Join(strings.Fields(lines[i]), " ") != line {
			t.Errorf("line %d: want %q in\n%s", i, line, out)
		}
	}
	if got, want := ready(), []string{"S-0004", "S-0001", "S-0003", "S-0002"}; !reflect.DeepEqual(got, want) {
		t.Errorf("--apply keeps S-0001 and reorders the rest: got %v want %v", got, want)
	}
	if got := placed(); !reflect.DeepEqual(got, want) {
		t.Errorf("--apply records no placement: got %+v", got)
	}
	data, _ := os.ReadFile(boardPath)
	if !strings.Contains(string(data), "placed:\n  S-0001:\n    by: flaiover\n    at: 2026-09-15T20:00:00Z\n") {
		t.Errorf("board.md front matter:\n%s", data)
	}

	// A story leaving ready leaves its placement behind.
	run("move", "S-0001", "in-progress")
	run("move", "S-0002", "backlog")
	if got := placed(); !reflect.DeepEqual(got, map[string]workitem.Placed{"S-0003": want["S-0003"]}) {
		t.Errorf("moved out of ready, S-0001 and S-0002 lose their placements: %+v", got)
	}

	before, _ := os.ReadFile(boardPath)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"order", "--by", "cod", "--placed-by", "alex"}, "drop --placed-by"},
		{[]string{"order", "S-0003", "--top", "--keep-placed", "1h"}, "give --by too"},
		{[]string{"order", "--by", "cod", "--keep-placed", "-1h"}, "0 for none"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, stderr %q, want %q", c.args, code, errOut, c.want)
		}
	}
	if after, _ := os.ReadFile(boardPath); string(after) != string(before) {
		t.Error("a refused order rewrote board.md")
	}
}
