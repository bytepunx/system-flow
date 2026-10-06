package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// S-0219: flai promote --drafts lists each draft story in the backlog by the
// manifest's policy, complete or with what it lacks, and writes nothing.
func TestPromoteDraftsCommand(t *testing.T) {
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
	fill := func(id, criterion string) {
		t.Helper()
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		body := strings.Replace(string(s), "## Goal\n", "## Goal\n\nDo it.\n", 1)
		body = strings.Replace(body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ]"+criterion+"\n", 1)
		if err := os.WriteFile(file[0], []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if out := run("promote", "--drafts"); out != "no draft stories in the backlog\n" {
		t.Fatalf("no drafts:\n%s", out)
	}
	run("story", "new", "Empty", "--draft")

	run("epic", "new", "Epic")
	run("story", "new", "Complete", "--epic", "E-0001", "--touches", "flai", "--draft")
	run("story", "new", "Uncriteried", "--epic", "E-0001", "--touches", "docs", "--draft")
	run("story", "new", "Untouched", "--epic", "E-0001", "--draft")
	run("story", "new", "Unforecast", "--epic", "E-0001", "--touches", "scripts", "--draft")
	run("story", "new", "Unvalued", "--epic", "E-0001", "--touches", "design", "--draft")
	run("story", "new", "Final", "--epic", "E-0001", "--touches", "wip")
	for _, id := range []string{"S-0002", "S-0004", "S-0005", "S-0006", "S-0007"} {
		fill(id, " it works")
	}
	fill("S-0003", "")
	forecast := []string{"--forecast-duration", "2h", "--forecast-delivery", "2026-09-20T12:00:00Z"}
	run(append([]string{"edit", "S-0002", "--cost-of-delay-value", "100"}, forecast...)...)
	run(append([]string{"edit", "S-0003", "--cost-of-delay-value", "900"}, forecast...)...)
	run(append([]string{"edit", "S-0004", "--cost-of-delay-value", "400"}, forecast...)...)
	run("edit", "S-0005", "--cost-of-delay-value", "300", "--forecast-duration", "1h")
	run(append([]string{"edit", "S-0006"}, forecast...)...)

	before := treeOf(t, root)
	out := run("promote", "--drafts")
	for _, want := range []string{
		"drafts by fifo:\n  1  S-0001  incomplete  Empty\n    - no goal\n    - no acceptance criteria with a checkbox\n    - no touches\n" +
			"    - no forecast duration\n    - no forecast delivery\n    - no cost of delay value\n",
		"  2  S-0002  complete  Complete\n  3  S-0003  incomplete  Uncriteried\n    - no acceptance criteria with a checkbox\n",
		"  4  S-0004  incomplete  Untouched\n    - no touches\n",
		"  5  S-0005  incomplete  Unforecast\n    - no forecast delivery\n",
		"  6  S-0006  incomplete  Unvalued\n    - no cost of delay value\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "S-0007") {
		t.Errorf("a story that is not a draft is listed:\n%s", out)
	}
	if !reflect.DeepEqual(treeOf(t, root), before) {
		t.Error("flai promote --drafts changed a file")
	}

	manifest := filepath.Join(root, "system-flow.yaml")
	m, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, append(m, []byte("orchestration:\n  policy: wsjf\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	var j struct {
		Policy string `json:"policy"`
		Drafts []struct {
			ID       string   `json:"id"`
			Complete bool     `json:"complete"`
			Lacks    []string `json:"lacks"`
		} `json:"drafts"`
	}
	if err := json.Unmarshal([]byte(run("promote", "--drafts", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, d := range j.Drafts {
		ids = append(ids, d.ID)
	}
	if j.Policy != "wsjf" || !reflect.DeepEqual(ids, []string{"S-0003", "S-0005", "S-0004", "S-0002", "S-0001", "S-0006"}) {
		t.Fatalf("by the manifest's wsjf: %+v", j)
	}
	if d := j.Drafts[3]; !d.Complete || d.Lacks == nil || len(d.Lacks) != 0 {
		t.Errorf("a complete draft lacks nothing, as an empty list: %+v", d)
	}
	if d := j.Drafts[0]; d.Complete || !reflect.DeepEqual(d.Lacks, []string{"no acceptance criteria with a checkbox"}) {
		t.Errorf("an incomplete draft says what it lacks: %+v", d)
	}

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"promote", "--candidates", "--drafts"}, "give --candidates or --drafts"},
		{[]string{"promote", "--drafts", "--limit", "2"}, "--limit caps the candidates"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, stderr %q, want %q", c.args, code, errOut, c.want)
		}
	}
}
