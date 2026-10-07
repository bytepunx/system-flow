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
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte("---\ntitle: Plan\nupdated: 2026-09-01\ntopics: [all]\n---\n\n# Plan\n\nSee [the story](../../wip/kanban/stories/S-0001-a-plain-story.md).\n"), 0o644)
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
	// none of them may be about the story, its parent, its narrative, or the link,
	// save narrative.state on the narrative the fixture opened and never wrote
	checked, _, _ := runIn(t, root, "check")
	var kept []string
	for _, l := range strings.Split(checked, "\n") {
		if !strings.Contains(l, ": narrative.state: ") {
			kept = append(kept, l)
		}
	}
	if checked = strings.Join(kept, "\n"); strings.Contains(checked, "S-0001") || strings.Contains(checked, "E-0001") || strings.Contains(checked, "plan.md") {
		t.Errorf("the check has nothing to say about what the retitle touched:\n%s", checked)
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
	if was := read(t, filepath.Join(root, "wip/kanban/epics/E-0001-first-epic.md")); strings.Contains(was, "- S-0001 ") {
		t.Errorf("it left the old parent's list:\n%s", was)
	}
	// S-0001 is in progress: E-0002, in backlog, follows it there (S-0200)
	if now := read(t, filepath.Join(root, "wip/kanban/epics/E-0002-second-epic.md")); !strings.Contains(now, "status: in-progress\n") || !strings.Contains(now, "follows S-0001, which joined it from E-0001") {
		t.Errorf("the epic it joined follows it:\n%s", now)
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

// I-0071: a touch that starts with a dot is a repository path like any other,
// and flai edit records it; one that leaves the repository is refused as a
// rule, and nothing is written.
func TestEditTouchesStartingWithADot(t *testing.T) {
	root := editProject(t)
	file := filepath.Join(root, "wip/kanban/stories/S-0001-a-plain-story.md")
	if out, errOut, code := runIn(t, root, "edit", "S-0001", "--touches", ".claude/agents/planner.md", "--touches", ".github/workflows/"); code != 0 || !strings.Contains(out, "changed touches") {
		t.Fatalf("edit: %d %s %s", code, out, errOut)
	}
	if out, _, _ := runIn(t, root, "touches", "S-0001"); out != "S-0001 touches .claude/agents/planner.md, .github/workflows\n" {
		t.Errorf("the story's touches: %q", out)
	}
	before := read(t, file)
	for _, touch := range []string{"../x", "a/../b", "/etc"} {
		if _, errOut, code := runIn(t, root, "edit", "S-0001", "--touches", touch); code == 0 || !strings.Contains(errOut, "rule: ") || !strings.Contains(errOut, "give a repository path or component name") {
			t.Errorf("touch %q: %d %s", touch, code, errOut)
		}
	}
	if read(t, file) != before {
		t.Error("no refusal wrote anything")
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

// S-0176: flai task new --after and flai edit --after set the tasks of the
// same story a task waits for; the creation and the edit are checked, so an
// entry that names no task, a task of another story, or a cycle is refused
// and nothing is written; --clear-after removes them. flai story new --after
// sets a story's.
func TestTaskAfterIsSetByTaskNewAndEdit(t *testing.T) {
	root := heldProject(t)
	for _, args := range [][]string{
		{"task", "new", "First", "--story", "S-0001"},
		{"task", "new", "Second", "--story", "S-0001", "--after", "t-1"},
		{"task", "new", "Elsewhere", "--story", "S-0002"},
	} {
		if out, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s %s", args, code, out, errOut)
		}
	}
	second := filepath.Join(root, "wip/kanban/tasks/T-0002-second.md")
	if item := read(t, second); !strings.Contains(item, "\nafter: [T-0001]\n") {
		t.Errorf("the task:\n%s", item)
	}
	if out, _, _ := runIn(t, root, "edit", "T-0002", "--show"); !strings.Contains(out, "  after: T-0001\n") {
		t.Errorf("show:\n%s", out)
	}

	for _, c := range []struct {
		args []string
		code int
		says string
	}{
		{[]string{"task", "new", "Bad", "--story", "S-0001", "--after", "T-0009"}, exitDocRefused, "task.after: after names T-0009, which does not exist"},
		{[]string{"task", "new", "Bad", "--story", "S-0001", "--after", "T-0003"}, exitDocRefused, "after names T-0003, a task of S-0002; a task waits only for tasks of its own story, S-0001"},
		{[]string{"task", "new", "Bad", "--story", "S-0001", "--after", "S-0002"}, 1, "name tasks of the same story, such as T-0001"},
		{[]string{"edit", "T-0001", "--after", "T-0002"}, exitDocRefused, "after forms a cycle, T-0001 waits for T-0002 waits for T-0001"},
		{[]string{"edit", "T-0001", "--after", "T-0003"}, exitDocRefused, "a task of S-0002"},
		{[]string{"edit", "T-0001", "--after", "T-1"}, 1, "a task cannot wait for itself"},
		{[]string{"edit", "E-0001", "--after", "S-0001"}, 1, "only a story or a task waits for others in after, not an epic"},
	} {
		out, errOut, code := runIn(t, root, c.args...)
		if code != c.code || !strings.Contains(out+errOut, c.says) {
			t.Errorf("%v: %d\n%s%s", c.args, code, out, errOut)
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/tasks", "T-0004-*.md")); len(matches) != 0 {
		t.Errorf("a refused creation leaves nothing behind: %v", matches)
	}
	if story := read(t, filepath.Join(root, "wip/kanban/stories/S-0001-open.md")); strings.Contains(story, "Bad") {
		t.Errorf("nor a line in the story's list:\n%s", story)
	}

	if out, errOut, code := runIn(t, root, "edit", "T-0001", "--after", "T-0002", "--clear-after"); code == 0 || !strings.Contains(errOut, "contradict") {
		t.Errorf("both: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, root, "edit", "T-0002", "--clear-after"); code != 0 || !strings.Contains(out, "T-0002: changed after") {
		t.Fatalf("clear: %d %s %s", code, out, errOut)
	}
	if item := read(t, second); strings.Contains(item, "after:") {
		t.Errorf("after: stays:\n%s", item)
	}
	if out, errOut, code := runIn(t, root, "edit", "T-0001", "--after", "T-0002"); code != 0 || !strings.Contains(out, "T-0001: changed after") {
		t.Errorf("set: %d %s %s", code, out, errOut)
	}

	if out, errOut, code := runIn(t, root, "story", "new", "Later", "--epic", "E-0001", "--after", "S-1"); code != 0 {
		t.Fatalf("story new --after: %d %s %s", code, out, errOut)
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", "S-0003-*.md")); len(matches) != 1 || !strings.Contains(read(t, matches[0]), "\nafter: [S-0001]\n") {
		t.Errorf("a story's after at creation: %v", matches)
	}
}

// S-0219: in the orchestrator's shell flai edit --no-draft finalizes only a
// draft that flai promote --drafts finds complete, with nothing else
// changed, and names the orchestrator as who finalized it; an incomplete one
// is refused naming what it lacks. Outside it --no-draft is as it was.
func TestTheOrchestratorFinalizesOnlyACompleteDraft(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "orchestrator")
	t.Setenv("FLAI_ROLE", "")
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
	run("story", "new", "Complete", "--epic", "E-0001", "--touches", "flai", "--draft")
	run("story", "new", "Incomplete", "--epic", "E-0001", "--draft")
	run("story", "new", "Also incomplete", "--epic", "E-0001", "--draft")
	stories, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", "S-000*.md"))
	complete := stories[0]
	s := read(t, complete)
	s = strings.Replace(s, "## Goal\n", "## Goal\n\nDo it.\n", 1)
	s = strings.Replace(s, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] it works\n", 1)
	if err := os.WriteFile(complete, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	run("edit", "S-0001", "--cost-of-delay-value", "300", "--forecast-duration", "2h", "--forecast-delivery", "2026-09-20T12:00:00Z")

	t.Setenv("FLAI_ROLE", "orchestrate")
	if _, errOut, code := runIn(t, root, "edit", "S-0001", "--no-draft"); code == 0 || !strings.Contains(errOut, "rule: the orchestrator finalizes a draft only with orchestration.permissions.finalize_drafts, which is off") {
		t.Errorf("a complete draft without finalize_drafts: exit %d %s", code, errOut)
	}
	permitOrchestrator(t, root, "finalize_drafts")
	if _, errOut, code := runIn(t, root, "edit", "S-0002", "--no-draft"); code == 0 || !strings.Contains(errOut, "rule: S-0002 is not complete, so the orchestrator does not finalize it (flai promote --drafts): no goal; no acceptance criteria with a checkbox; no touches; no forecast duration; no forecast delivery; no cost of delay value") {
		t.Errorf("an incomplete draft: exit %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "edit", "S-0001", "--no-draft", "--title", "Renamed"); code == 0 || !strings.Contains(errOut, "rule: S-0001: the orchestrator finalizes a draft with --no-draft and nothing else") {
		t.Errorf("with another change: exit %d %s", code, errOut)
	}
	if item := read(t, complete); !strings.Contains(item, "\ndraft: true\n") || strings.Contains(item, "finalized:") {
		t.Fatalf("a refused edit finalized the story:\n%s", item)
	}
	if out := run("edit", "S-0001", "--no-draft"); !strings.Contains(out, "S-0001: changed draft") {
		t.Errorf("a complete draft: %s", out)
	}
	if item := read(t, complete); strings.Contains(item, "\ndraft: true\n") || !strings.Contains(item, "finalized:\n  by: orchestrator\n") {
		t.Errorf("finalized by the orchestrator:\n%s", item)
	}

	t.Setenv("FLAI_ROLE", "")
	if out := run("edit", "S-0003", "--no-draft"); !strings.Contains(out, "S-0003: changed draft") {
		t.Errorf("an incomplete draft outside the orchestrator's shell, as before: %s", out)
	}
}
