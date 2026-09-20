package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// editProject is a git repository with two epics, a story in progress with a
// narrative, and a design document that links to the story's file.
func editProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "design", "system"), 0o755)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	for _, args := range [][]string{{"epic", "new", "First epic"}, {"epic", "new", "Second epic"}, {"story", "new", "A plain story", "--epic", "E-0001", "--tag", "cli"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] it works\n", 1)), 0o644)
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte("---\ntitle: Plan\nupdated: 2026-09-01\n---\n\n# Plan\n\nSee [the story](../../wip/kanban/stories/S-0001-a-plain-story.md).\n"), 0o644)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: in progress")
	return root
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// S-0085: a title lives in six places, and a retitle keeps them in step, in
// one commit.
func TestEditRetitlesEverywhere(t *testing.T) {
	root := editProject(t)
	out, errOut, code := runIn(t, root, "edit", "S-0001", "--title", "  A better   name ", "--autocommit", "--trailer", "Co-Authored-By: someone <s@s>")
	if code != 0 || !strings.Contains(out, "changed title") || !strings.Contains(out, "renamed from wip/kanban/stories/S-0001-a-plain-story.md") {
		t.Fatalf("edit: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md")); !os.IsNotExist(err) {
		t.Error("the old file is gone")
	}
	item := read(t, filepath.Join(root, "wip/kanban/stories/S-0001-a-better-name.md"))
	for _, want := range []string{"title: A better name\n", "\n# S-0001 A better name\n", "- [ ] it works", "status: in-progress"} {
		if !strings.Contains(item, want) {
			t.Errorf("the item lacks %q:\n%s", want, item)
		}
	}
	if epic := read(t, filepath.Join(root, "wip/kanban/epics/E-0001-first-epic.md")); !strings.Contains(epic, "- S-0001 A better name\n") || strings.Contains(epic, "A plain story") {
		t.Errorf("the parent's list:\n%s", epic)
	}
	if n := read(t, filepath.Join(root, "wip/agents/S-0001.md")); !strings.Contains(n, "title: A better name\n") || !strings.Contains(n, "# S-0001 A better name\n") || strings.Contains(n, "A plain story") {
		t.Errorf("the narrative:\n%s", n)
	}
	if plan := read(t, filepath.Join(root, "design/system/plan.md")); !strings.Contains(plan, "stories/S-0001-a-better-name.md") {
		t.Errorf("a link to the old file name is kept working:\n%s", plan)
	}
	if idx := read(t, filepath.Join(root, "wip/agents/index.md")); strings.Contains(idx, "A plain story") {
		t.Errorf("the index:\n%s", idx)
	}
	if st := gitIn(t, root, "status", "--porcelain"); strings.TrimSpace(st) != "" {
		t.Errorf("one commit holds every file, the removal included:\n%s", st)
	}
	msg := gitIn(t, root, "log", "-1", "--format=%B")
	if !strings.HasPrefix(msg, "chore: [S-0001] edit title") || !strings.Contains(msg, "Co-Authored-By: someone") {
		t.Errorf("commit message: %s", msg)
	}
	// the fixture is not a whole project, so the check has things to say about its layout;
	// none of them may be about the story, its parent, its narrative, or the link
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "S-0001") || strings.Contains(out, "E-0001") || strings.Contains(out, "plan.md") {
		t.Errorf("the check has nothing to say about what the retitle touched:\n%s", out)
	}
}

func TestEditFieldsParentAndBody(t *testing.T) {
	root := editProject(t)
	out, errOut, code := runIn(t, root, "edit", "S-0001", "--nature", "remediation", "--tag", "dashboard,cli", "--touches", "flai/cmd/", "--touches", "flaiover/src", "--parent", "E-2", "--json")
	var res struct {
		Changed []string `json:"changed"`
		Files   []string `json:"files"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil || strings.Join(res.Changed, ",") != "nature,tags,touches,parent" {
		t.Fatalf("edit: %d %s %s", code, out, errOut)
	}
	item := read(t, filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md"))
	for _, want := range []string{"nature: remediation\n", "tags: [dashboard, cli]\n", "touches: [flai/cmd, flaiover/src]\n", "parent: E-0002\n"} {
		if !strings.Contains(item, want) {
			t.Errorf("the item lacks %q:\n%s", want, item)
		}
	}
	if was := read(t, filepath.Join(root, "wip/kanban/epics/E-0001-first-epic.md")); strings.Contains(was, "S-0001") {
		t.Errorf("it left the old parent's list:\n%s", was)
	}
	if now := read(t, filepath.Join(root, "wip/kanban/epics/E-0002-second-epic.md")); !strings.Contains(now, "- S-0001 A plain story") {
		t.Errorf("and joined the new one's:\n%s", now)
	}
	if out, _, _ := runIn(t, root, "edit", "S-0001", "--clear-tags", "--clear-touches"); !strings.Contains(out, "changed tags, touches") {
		t.Errorf("clearing: %s", out)
	}
	if out, _, _ := runIn(t, root, "edit", "S-0001", "--nature", "remediation"); !strings.Contains(out, "unchanged") {
		t.Errorf("the same again changes nothing: %s", out)
	}

	// the body: what lies below the heading; the heading stays flai's
	show, _, _ := runIn(t, root, "edit", "S-0001", "--show", "--json")
	var v struct {
		Body, Hash string
		Editable   bool
		Parents    []struct{ ID string }
	}
	if err := json.Unmarshal([]byte(show), &v); err != nil || !v.Editable || strings.HasPrefix(v.Body, "# ") || !strings.Contains(v.Body, "## Goal") || len(v.Parents) != 2 {
		t.Fatalf("show: %v %s", err, show)
	}
	body := strings.Replace(v.Body, "- [ ] it works", "- [x] it works\n- [ ] and it is fast", 1)
	out, errOut, code = runStdin(t, root, body, "edit", "S-0001", "--body-stdin", "--hash", v.Hash)
	if code != 0 || !strings.Contains(out, "changed criteria") {
		t.Fatalf("body: %d %s %s", code, out, errOut)
	}
	item = read(t, filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md"))
	if !strings.Contains(item, "\n# S-0001 A plain story\n\n## Goal") || !strings.Contains(item, "- [ ] and it is fast") {
		t.Errorf("the body:\n%s", item)
	}
	// the hash is of the file as it was read: it no longer matches
	if _, errOut, code := runStdin(t, root, body+"\nmore\n", "edit", "S-0001", "--body-stdin", "--hash", v.Hash); code != exitDocConflict || !strings.Contains(errOut, "changed after it was loaded") {
		t.Errorf("a change made meanwhile is a conflict: %d %s", code, errOut)
	}
}

func TestEditRefusals(t *testing.T) {
	root := editProject(t)
	file := filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md")
	before := read(t, file)
	epicBefore := read(t, filepath.Join(root, "wip/kanban/epics/E-0001-first-epic.md"))

	// the check refuses a story in progress without criteria, and everything is put back,
	// the retitle that came with it included
	out, errOut, code := runStdin(t, root, "## Goal\nNo criteria any more.\n", "edit", "S-0001", "--title", "Renamed too", "--body-stdin")
	if code != exitDocRefused || !strings.Contains(out+errOut, "story.criteria") {
		t.Fatalf("refused by the check: %d %s %s", code, out, errOut)
	}
	if read(t, file) != before || read(t, filepath.Join(root, "wip/kanban/epics/E-0001-first-epic.md")) != epicBefore {
		t.Error("every file is as it was")
	}
	if _, err := os.Stat(filepath.Join(root, "wip/kanban/stories/S-0001-renamed-too.md")); !os.IsNotExist(err) {
		t.Error("and no renamed file is left behind")
	}
	if st := gitIn(t, root, "status", "--porcelain"); strings.TrimSpace(st) != "" {
		t.Errorf("git sees no change:\n%s", st)
	}

	for name, c := range map[string]struct {
		args []string
		says string
	}{
		"no title":              {[]string{"edit", "S-0001", "--title", "  "}, "a title is required"},
		"a nature there is not": {[]string{"edit", "S-0001", "--nature", "urgent"}, "must be one of"},
		"a tag that is a flag":  {[]string{"edit", "S-0001", "--tag=--owner=eve"}, "not starting with a dash"},
		"a story as parent":     {[]string{"edit", "S-0001", "--parent", "S-0001"}, "must be a epic"},
		"a parent there is not": {[]string{"edit", "S-0001", "--parent", "E-0009"}, "not found"},
		"an epic's parent":      {[]string{"edit", "E-0001", "--parent", "E-0002"}, "epics have no parent"},
		"an epic's touches":     {[]string{"edit", "E-0001", "--touches", "flai"}, "touches belong to stories and tasks"},
		"nothing":               {[]string{"edit", "S-0001"}, "nothing to change"},
		"no such item":          {[]string{"edit", "S-0042", "--title", "x"}, "not found"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.says) {
			t.Errorf("%s: %d %s", name, code, errOut)
		}
		// what the caller got wrong is said as a rule, so that a dashboard answers 400 and
		// not 500; an item that is not there is said as that, which it answers with 404
		if wantRule := !strings.Contains(c.says, "not found") && name != "nothing"; wantRule != strings.Contains(errOut, "rule: ") {
			t.Errorf("%s: a rule or not: %s", name, errOut)
		}
	}
	if _, errOut, code := runStdin(t, root, "# My own heading\n\n## Goal\nx\n", "edit", "S-0001", "--body-stdin"); code == 0 || !strings.Contains(errOut, "heading of its own") {
		t.Errorf("a body with its own heading: %d %s", code, errOut)
	}
	if read(t, file) != before {
		t.Error("no refusal wrote anything")
	}

	// a closed item is not edited
	if _, errOut, code := runIn(t, root, "move", "S-0001", "cancelled", "--reason", "test"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "edit", "S-0001", "--title", "Too late"); code != exitDocRefused || !strings.Contains(errOut, "cancelled") {
		t.Errorf("a closed item: %d %s", code, errOut)
	}
}
