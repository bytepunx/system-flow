package cmd

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// treeOf reads every file under root, by its path.
func treeOf(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(p)
		out[p] = string(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// S-0217: flai promote --candidates lists the backlog stories that could go
// to ready by the manifest's policy, caps them with --limit, says why each
// other one cannot, and writes nothing.
func TestPromoteCandidatesCommand(t *testing.T) {
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
	fill := func(id string) {
		t.Helper()
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		body := strings.Replace(string(s), "## Goal\n", "## Goal\n\nDo it.\n", 1)
		body = strings.Replace(body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)
		if err := os.WriteFile(file[0], []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("epic", "new", "Epic")
	run("story", "new", "Open", "--epic", "E-0001", "--touches", "flai")
	run("story", "new", "Cheap", "--epic", "E-0001", "--touches", "docs")
	run("story", "new", "Dear", "--epic", "E-0001", "--touches", "template")
	run("story", "new", "Quick", "--epic", "E-0001", "--touches", "scripts")
	run("story", "new", "Held", "--epic", "E-0001", "--touches", "flai/cmd")
	run("story", "new", "Unvalued", "--epic", "E-0001", "--touches", "design")
	run("story", "new", "Drafted", "--epic", "E-0001", "--touches", "wip", "--draft")
	run("story", "new", "Unwritten", "--epic", "E-0001", "--touches", "Makefile")
	for _, id := range []string{"S-0001", "S-0002", "S-0003", "S-0004", "S-0005", "S-0006", "S-0007"} {
		fill(id)
	}
	run("edit", "S-0002", "--cost-of-delay-value", "100", "--forecast-duration", "1h")
	run("edit", "S-0003", "--cost-of-delay-value", "900", "--forecast-duration", "30h")
	run("edit", "S-0004", "--cost-of-delay-value", "300", "--forecast-duration", "30m")
	run("edit", "S-0005", "--cost-of-delay-value", "300", "--forecast-duration", "1h")
	run("edit", "S-0006", "--forecast-duration", "1h")
	run("edit", "S-0007", "--cost-of-delay-value", "300", "--forecast-duration", "1h")
	run("move", "S-0001", "ready")
	run("move", "S-0001", "in-progress")

	before := treeOf(t, root)
	out := run("promote", "--candidates")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i, want := range []string{"candidates by fifo:", "1 S-0002 2026-09-15T21:00:00Z Cheap"} {
		if i >= len(lines) || strings.Join(strings.Fields(lines[i]), " ") != want {
			t.Fatalf("with no policy set, fifo; line %d want %q:\n%s", i, want, out)
		}
	}
	for _, want := range []string{
		"not candidates:\n  S-0005  Held\n    - held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches;",
		"  S-0006  Unvalued\n    - no cost of delay value\n",
		"  S-0007  Drafted\n    - draft: finalize it first\n",
		"  S-0008  Unwritten\n    - no goal\n    - no acceptance criteria with a checkbox\n    - no forecast duration\n    - no cost of delay value\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
	if !reflect.DeepEqual(treeOf(t, root), before) {
		t.Error("flai promote --candidates changed a file")
	}

	manifest := filepath.Join(root, "system-flow.yaml")
	m, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, append(m, []byte("orchestration:\n  policy: wsjf\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	var j struct {
		Policy     string `json:"policy"`
		Candidates []struct {
			Position int    `json:"position"`
			ID       string `json:"id"`
			Text     string `json:"text"`
		} `json:"candidates"`
		Others []struct {
			ID      string   `json:"id"`
			Reasons []string `json:"reasons"`
			Held    *struct {
				Code string `json:"code"`
			} `json:"held"`
		} `json:"others"`
	}
	if err := json.Unmarshal([]byte(run("promote", "--candidates", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, c := range j.Candidates {
		ids = append(ids, c.ID)
	}
	if j.Policy != "wsjf" || !reflect.DeepEqual(ids, []string{"S-0004", "S-0002", "S-0003"}) || j.Candidates[0].Text != "600/h" {
		t.Errorf("by the manifest's wsjf: %+v", j)
	}
	if len(j.Others) != 4 || j.Others[0].ID != "S-0005" || j.Others[0].Held == nil || j.Others[0].Held.Code != "overlap" {
		t.Errorf("others: %+v", j.Others)
	}

	j.Candidates, j.Others = nil, nil
	if err := json.Unmarshal([]byte(run("promote", "--candidates", "--limit", "2", "--json")), &j); err != nil {
		t.Fatal(err)
	}
	if len(j.Candidates) != 2 || j.Candidates[1].ID != "S-0002" || len(j.Others) != 4 {
		t.Errorf("--limit 2 caps the candidates and leaves the others as they were: %+v", j)
	}
	out = run("promote", "--candidates", "--limit", "1")
	if !strings.Contains(out, "candidates by wsjf:\n") || strings.Contains(out, "S-0002  2026") || strings.Count(out, "/h") != 1 {
		t.Errorf("--limit 1:\n%s", out)
	}

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"promote"}, "give --candidates"},
		{[]string{"promote", "--candidates", "--limit", "-1"}, "--limit"},
		{[]string{"promote", "--candidates", "S-0002"}, "unknown command"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, stderr %q, want %q", c.args, code, errOut, c.want)
		}
	}
}
