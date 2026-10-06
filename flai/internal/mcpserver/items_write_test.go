package mcpserver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0103: stories are created and edited over MCP with their agent: the
// project's default filled in, what is given winning, and an edit that
// replaces or removes it; a stale hash is refused, not overwritten.
func TestItemNewAndEditCarryTheAgent(t *testing.T) {
	f := setup(t)
	f.repo.Manifest.Agent = &manifest.Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high"}}

	out, failed := f.call(t, "item_new", map[string]any{"type": "story", "title": "From an agent", "agent": map[string]any{"model": "claude-sonnet-5"}, "body": "## Goal\n\nSomething.\n\n## Acceptance criteria\n- [ ] it works\n\n## Tasks\n\n## Notes\n"})
	if failed != "" {
		t.Fatal(failed)
	}
	id := out["id"].(string)
	agent := out["agent"].(map[string]any)
	if agent["harness"] != "claude-code" || agent["model"] != "claude-sonnet-5" || agent["config"].(map[string]any)["effort"] != "high" {
		t.Errorf("created with: %v", agent)
	}
	if !strings.Contains(out["body"].(string), "Something.") {
		t.Errorf("body: %v", out["body"])
	}

	got, _ := f.call(t, "item_get", map[string]any{"id": id})
	hash := got["hash"].(string)
	if got["default_agent"].(map[string]any)["model"] != "claude-opus-5-5" || len(hash) != 64 {
		t.Errorf("item_get: %v", got)
	}
	ed, failed := f.call(t, "item_edit", map[string]any{"id": id, "hash": hash, "agent": map[string]any{"harness": "codex", "model": "gpt-5.1"}})
	if failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "agent" {
		t.Fatalf("replace: %v %s", ed, failed)
	}
	got, _ = f.call(t, "item_get", map[string]any{"id": id})
	if a := got["agent"].(map[string]any); a["harness"] != "codex" || a["config"] != nil {
		t.Errorf("replaced exactly: %v", a)
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": id, "hash": hash, "title": "Stale"}); !strings.Contains(failed, "changed") {
		t.Errorf("a stale hash: %q", failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": id, "clear_agent": true}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "agent" {
		t.Errorf("clear: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": id}); got["agent"] != nil {
		t.Errorf("cleared: %v", got["agent"])
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": id}); !strings.Contains(failed, "nothing to change") {
		t.Errorf("empty edit: %q", failed)
	}
	if _, failed := f.call(t, "item_new", map[string]any{"type": "epic", "title": "E", "agent": map[string]any{"model": "x"}}); !strings.Contains(failed, "only a story carries an agent") {
		t.Errorf("an epic's agent: %q", failed)
	}
}

func toStrings(v any) []string {
	var out []string
	for _, x := range v.([]any) {
		out = append(out, x.(string))
	}
	return out
}

// S-0130: item_edit sets a story's after:, which item_get shows, refuses a
// story that does not exist, and an empty list removes it.
func TestItemEditSetsAndClearsAfter(t *testing.T) {
	f := setup(t)
	later := f.readyStory(t, "Later", t0)
	ed, failed := f.call(t, "item_edit", map[string]any{"id": later.ID, "after": []string{f.story.ID}})
	if failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "after" {
		t.Fatalf("set: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": later.ID}); strings.Join(toStrings(got["after"]), ",") != f.story.ID {
		t.Errorf("item_get: %v", got["after"])
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": later.ID, "after": []string{"S-0999"}}); !strings.Contains(failed, "S-0999, which does not exist") {
		t.Errorf("a missing story: %q", failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": later.ID, "after": []string{}}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "after" {
		t.Errorf("clear: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": later.ID}); got["after"] != nil {
		t.Errorf("cleared: %v", got["after"])
	}
}

// S-0176: item_new and item_edit set a task's after, the tasks of its own
// story it waits for. A creation or an edit whose entry names no task, a task
// of another story, or forms a cycle is refused with the check's finding, and
// a refused creation leaves nothing; item_new sets a story's after too.
func TestItemNewAndEditSetATasksAfter(t *testing.T) {
	f := setup(t)
	out, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Second", "parent": f.story.ID, "after": []string{"t-1"}})
	if failed != "" || strings.Join(toStrings(out["after"]), ",") != f.task.ID {
		t.Fatalf("new task: %v %s", out, failed)
	}
	second := out["id"].(string)
	if _, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Bad", "parent": f.story.ID, "after": []string{"T-0009"}}); !strings.Contains(failed, "task.after: after names T-0009, which does not exist") {
		t.Errorf("a task that does not exist: %q", failed)
	}
	if _, err := f.repo.Get("T-0003"); err == nil {
		t.Error("a refused creation leaves nothing behind")
	}
	other := f.readyStory(t, "Other", t0)
	out, failed = f.call(t, "item_new", map[string]any{"type": "task", "title": "Elsewhere", "parent": other.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	elsewhere := out["id"].(string)
	if _, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Bad", "parent": f.story.ID, "after": []string{elsewhere}}); !strings.Contains(failed, "a task of "+other.ID) {
		t.Errorf("a task of another story: %q", failed)
	}
	if _, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Bad", "parent": f.story.ID, "after": []string{other.ID}}); !strings.Contains(failed, "name tasks of the same story") {
		t.Errorf("a story in a task's after: %q", failed)
	}

	if _, failed := f.call(t, "item_edit", map[string]any{"id": f.task.ID, "after": []string{second}}); !strings.Contains(failed, "after forms a cycle, "+f.task.ID+" waits for "+second+" waits for "+f.task.ID) {
		t.Errorf("a cycle: %q", failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": second, "after": []string{}}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "after" {
		t.Errorf("clear: %v %s", ed, failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": f.task.ID, "after": []string{second}}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "after" {
		t.Errorf("set: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": f.task.ID}); strings.Join(toStrings(got["after"]), ",") != second || got["plan"] != nil {
		t.Errorf("item_get: %v %v", got["after"], got["plan"])
	}
	// the story's plan: the first task waits for the second, which waits for none
	got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID})
	plan, _ := got["plan"].(map[string]any)
	if tasks, _ := plan["tasks"].([]any); len(tasks) != 2 || tasks[0].(map[string]any)["state"] != "waiting" || strings.Join(toStrings(tasks[0].(map[string]any)["waiting_for"]), ",") != second {
		t.Errorf("item_get's plan: %v", plan)
	}
	if layers, _ := plan["layers"].([]any); len(layers) != 2 || strings.Join(toStrings(layers[0]), ",") != second || strings.Join(toStrings(layers[1]), ",") != f.task.ID {
		t.Errorf("item_get's layers: %v", plan["layers"])
	}

	out, failed = f.call(t, "item_new", map[string]any{"type": "story", "title": "Later", "after": []string{other.ID}})
	if failed != "" || strings.Join(toStrings(out["after"]), ",") != other.ID {
		t.Errorf("new story: %v %s", out, failed)
	}
}

// S-0135: item_new and item_edit set a story's or epic's topics, item_get
// shows them, a task is refused them, and an empty list removes them.
func TestItemNewAndEditCarryTopics(t *testing.T) {
	f := setup(t)
	out, failed := f.call(t, "item_new", map[string]any{"type": "epic", "title": "Logs", "topics": []string{"logging"}})
	if failed != "" || strings.Join(toStrings(out["topics"]), ",") != "logging" {
		t.Fatalf("new epic: %v %s", out, failed)
	}
	ed, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "topics": []string{"release", "logging"}})
	if failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "topics" {
		t.Fatalf("set: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID}); strings.Join(toStrings(got["topics"]), ",") != "release,logging" {
		t.Errorf("item_get: %v", got["topics"])
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "topics": []string{"not one"}}); !strings.Contains(failed, "one word") {
		t.Errorf("two words: %q", failed)
	}
	if _, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "T", "parent": f.story.ID, "topics": []string{"logging"}}); !strings.Contains(failed, "stories and epics") {
		t.Errorf("a task's topics: %q", failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "topics": []string{}}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "topics" {
		t.Errorf("clear: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID}); got["topics"] != nil {
		t.Errorf("cleared: %v", got["topics"])
	}
}

// I-0071: item_edit and item_new take a touch that starts with a dot, which
// item_get shows, and refuse one that leaves the repository.
func TestItemEditAndNewTakeATouchStartingWithADot(t *testing.T) {
	f := setup(t)
	ed, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "touches": []string{".claude/agents/planner.md"}})
	if failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "touches" {
		t.Fatalf("set: %v %s", ed, failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID}); strings.Join(toStrings(got["touches"]), ",") != ".claude/agents/planner.md" {
		t.Errorf("item_get: %v", got["touches"])
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "touches": []string{"../x"}}); !strings.Contains(failed, "give a repository path or component name") {
		t.Errorf("a touch outside the repository: %q", failed)
	}
	out, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Dotted", "parent": f.story.ID, "touches": []string{".github/workflows/"}})
	if failed != "" || strings.Join(toStrings(out["touches"]), ",") != ".github/workflows" {
		t.Errorf("new task: %v %s", out, failed)
	}
	if _, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Outside", "parent": f.story.ID, "touches": []string{"a/../../b"}}); !strings.Contains(failed, "climbs out of the repository") {
		t.Errorf("a new task outside the repository: %q", failed)
	}
}

// S-0199: item_new makes a story a draft, which item_get shows, and refuses a
// draft of anything else; item_move refuses a draft story to ready, since
// finalizing is the operator's, and so does item_edit's draft false.
func TestADraftStoryIsMadeButNotFinalizedByAnAgent(t *testing.T) {
	f := setup(t)
	out, failed := f.call(t, "item_new", map[string]any{"type": "story", "title": "Drafted", "draft": true, "body": "## Goal\n\nx\n\n## Acceptance criteria\n- [ ] it works\n\n## Tasks\n\n## Notes\n"})
	if failed != "" || out["draft"] != true {
		t.Fatalf("new draft: %v %s", out, failed)
	}
	id := out["id"].(string)
	for typ, parent := range map[string]string{"epic": "", "task": f.story.ID} {
		if _, failed := f.call(t, "item_new", map[string]any{"type": typ, "title": "Not a story", "parent": parent, "draft": true}); !strings.Contains(failed, "only a story is a draft") {
			t.Errorf("a draft %s: %q", typ, failed)
		}
	}
	if _, failed := f.call(t, "item_move", map[string]any{"id": id, "to": "ready"}); !strings.Contains(failed, id+" is a draft") || !strings.Contains(failed, "operator's") {
		t.Errorf("a draft to ready: %q", failed)
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": id, "draft": false}); !strings.Contains(failed, "operator's") {
		t.Errorf("finalized by an agent: %q", failed)
	}
	if got, _ := f.call(t, "item_get", map[string]any{"id": id}); got["status"] != "backlog" || got["draft"] != true {
		t.Errorf("still a backlog draft: %v %v", got["status"], got["draft"])
	}
	plain, failed := f.call(t, "item_new", map[string]any{"type": "story", "title": "Plain"})
	if failed != "" || plain["draft"] != nil {
		t.Fatalf("new story: %v %s", plain, failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": plain["id"], "draft": true}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "draft" {
		t.Errorf("made a draft: %v %s", ed, failed)
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "draft": true}); !strings.Contains(failed, "only a story in the backlog") {
		t.Errorf("an in-progress story made a draft: %q", failed)
	}
}

// S-0209: item_new with a body is checked as flai story new --body-stdin is,
// so a story the planner writes passes flai check --strict and the markdown
// lint before it is kept. One that leaves out a section, or that the lint
// rejects, is refused with the findings and leaves nothing: no file, and its
// epic as it was. A clean one is kept and listed in its epic.
func TestItemNewWithABodyIsChecked(t *testing.T) {
	f := setup(t)
	f.lintAsThisProjectDoes(t)
	epic, err := f.repo.Get(f.story.Parent)
	if err != nil {
		t.Fatal(err)
	}
	epicWas, _ := os.ReadFile(epic.Path)
	stories := func() string {
		entries, _ := os.ReadDir(filepath.Join(f.repo.Root, "wip/kanban/stories"))
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		return strings.Join(names, " ")
	}
	storiesWere := stories()
	for name, c := range map[string]struct{ body, want string }{
		"a missing section": {"## Goal\n\nx\n\n## Acceptance criteria\n\n- [ ] it works\n\n## Notes\n", `item.heading: body is missing the "## Tasks" section`},
		"a lint finding":    {"## Goal\n\nx\n\n#### Too deep\n\n## Acceptance criteria\n\n- [ ] it works\n\n## Tasks\n\n## Notes\n", "markdown.MD001"},
	} {
		_, failed := f.call(t, "item_new", map[string]any{"type": "story", "title": "Planned", "parent": epic.ID, "draft": true, "body": c.body})
		if !strings.Contains(failed, c.want) || !strings.Contains(failed, "nothing was created") {
			t.Errorf("%s: refused with %q", name, failed)
		}
		if got := stories(); got != storiesWere {
			t.Errorf("%s: a refused story leaves nothing: %s", name, got)
		}
		if got, _ := os.ReadFile(epic.Path); string(got) != string(epicWas) {
			t.Errorf("%s: the epic changed:\n%s", name, got)
		}
	}
	out, failed := f.call(t, "item_new", map[string]any{"type": "story", "title": "Planned", "parent": epic.ID, "draft": true, "body": "## Goal\n\nx\n\n## Acceptance criteria\n\n- [ ] it works\n\n## Tasks\n\n## Notes\n"})
	if failed != "" {
		t.Fatalf("a clean story: %s", failed)
	}
	if got, _ := os.ReadFile(epic.Path); !strings.Contains(string(got), out["id"].(string)) {
		t.Errorf("the epic lists the kept story:\n%s", got)
	}
}

// lintAsThisProjectDoes gives the fixture the lint this project runs, in
// short: a front matter title is not a heading.
func (f *fixture) lintAsThisProjectDoes(t *testing.T) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(f.repo.Root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD025:\n  front_matter_title: \"\"\nMD041: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0255: a task the planner writes over item_new with a body is checked as
// a story's is (S-0209). One the lint rejects, here a list indented one
// space, or that leaves out a section, is refused with the finding and
// leaves nothing: no task file, and its story as it was. One with its work,
// what done means, a nature, tags, touches, and after its sibling is kept as
// written and listed in its story.
func TestItemNewWithABodyChecksAStorysTask(t *testing.T) {
	f := setup(t)
	f.lintAsThisProjectDoes(t)
	storyWas, _ := os.ReadFile(f.story.Path)
	tasks := func() string {
		entries, _ := os.ReadDir(filepath.Join(f.repo.Root, "wip/kanban/tasks"))
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		return strings.Join(names, " ")
	}
	tasksWere := tasks()
	clean := "## Work\n\nTest the guard.\n\n## Done when\n\n- The test fails without the rule.\n\n## Notes\n\nNone.\n"
	task := func(body string) map[string]any {
		return map[string]any{"type": "task", "title": "Test the guard", "parent": f.story.ID, "nature": "improvement",
			"tags": []string{"planner"}, "touches": []string{"flai/internal/guard"}, "body": body}
	}
	// no after, so that it is the body alone that has the task checked
	for name, c := range map[string]struct{ body, want string }{
		"a lint finding":    {strings.ReplaceAll(clean, "\n- ", "\n - "), "markdown.MD007"},
		"a missing section": {"## Work\n\nTest the guard.\n\n## Notes\n\nNone.\n", `item.heading: body is missing the "## Done when" section`},
	} {
		_, failed := f.call(t, "item_new", task(c.body))
		if !strings.Contains(failed, c.want) || !strings.Contains(failed, "nothing was created") {
			t.Errorf("%s: refused with %q", name, failed)
		}
		if got := tasks(); got != tasksWere {
			t.Errorf("%s: a refused task leaves nothing: %s", name, got)
		}
		if got, _ := os.ReadFile(f.story.Path); string(got) != string(storyWas) {
			t.Errorf("%s: the story changed:\n%s", name, got)
		}
	}
	kept := task(clean)
	kept["after"] = []string{f.task.ID}
	out, failed := f.call(t, "item_new", kept)
	if failed != "" {
		t.Fatalf("a clean task: %s", failed)
	}
	id := out["id"].(string)
	if out["parent"] != f.story.ID || out["nature"] != "improvement" || strings.Join(toStrings(out["tags"]), ",") != "planner" ||
		strings.Join(toStrings(out["touches"]), ",") != "flai/internal/guard" || strings.Join(toStrings(out["after"]), ",") != f.task.ID {
		t.Errorf("kept as written: %v", out)
	}
	if !strings.Contains(out["body"].(string), "## Done when\n\n- The test fails without the rule.") {
		t.Errorf("the body is the planner's:\n%s", out["body"])
	}
	if got, _ := os.ReadFile(f.story.Path); !strings.Contains(string(got), id) {
		t.Errorf("the story lists the kept task:\n%s", got)
	}
}

// S-0199: item_edit sets, removes one key of, replaces, and clears a cost of
// delay and a forecast, each stamped with the agent and the time; a cost of
// delay's inputs and value each carry the stamp (ADR-0080).
func TestItemEditSetsAndClearsCostOfDelayAndForecast(t *testing.T) {
	f := setup(t)
	at := f.clock.UTC().Format(workitem.TimeFormat)
	planning := func(key string) map[string]any {
		t.Helper()
		got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID})
		v, _ := got[key].(map[string]any)
		return v
	}
	edit := func(args map[string]any, want string) {
		t.Helper()
		args["id"] = f.story.ID
		ed, failed := f.call(t, "item_edit", args)
		if failed != "" || strings.Join(toStrings(ed["changed"]), ",") != want {
			t.Fatalf("%v: %v %s", args, ed, failed)
		}
	}

	edit(map[string]any{"cost_of_delay": map[string]any{"revenue_per_week": 1200.5, "penalty_per_week": 300, "time_lost_per_cycle": "4h", "value": 2000}}, "cost_of_delay")
	got, _ := f.call(t, "item_get", map[string]any{"id": f.story.ID})
	c := got["cost_of_delay"].(map[string]any)
	in := c["inputs"].(map[string]any)
	if in["revenue_per_week"] != 1200.5 || in["penalty_per_week"] != 300.0 || in["time_lost_per_cycle"] != "4h" || c["value"] != 2000.0 || c["by"] != "claude" || c["at"] != at || in["by"] != "claude" || in["at"] != at || got["currency"] != "USD" {
		t.Errorf("set: %v, currency %v", c, got["currency"])
	}
	edit(map[string]any{"cost_of_delay": map[string]any{"time_lost_per_cycle": ""}}, "cost_of_delay")
	if in := planning("cost_of_delay")["inputs"].(map[string]any); in["time_lost_per_cycle"] != nil || in["revenue_per_week"] != 1200.5 {
		t.Errorf("one key removed: %v", in)
	}
	// an amount is removed by replacing the block with the keys to keep
	edit(map[string]any{"clear_cost_of_delay": true, "cost_of_delay": map[string]any{"revenue_per_week": 1200.5}}, "cost_of_delay")
	// the inputs and the value are stamped apart (ADR-0080): with the value
	// gone, so is its stamp
	if c := planning("cost_of_delay"); c["value"] != nil || c["by"] != nil || c["at"] != nil || len(c["inputs"].(map[string]any)) != 3 || c["inputs"].(map[string]any)["by"] != "claude" {
		t.Errorf("replaced: %v", c)
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": f.story.ID, "cost_of_delay": map[string]any{"time_lost_per_cycle": "soon"}}); !strings.Contains(failed, "not a Go duration") {
		t.Errorf("a bad duration: %q", failed)
	}
	edit(map[string]any{"clear_cost_of_delay": true}, "cost_of_delay")
	if c := planning("cost_of_delay"); c != nil {
		t.Errorf("cleared: %v", c)
	}

	edit(map[string]any{"forecast": map[string]any{"duration": "6h", "delivery": "2026-10-09T17:00:00Z", "basis": "three tasks like the last story's"}}, "forecast")
	if fc := planning("forecast"); fc["duration"] != "6h" || fc["delivery"] != "2026-10-09T17:00:00Z" || fc["basis"] != "three tasks like the last story's" || fc["by"] != "claude" || fc["at"] != at {
		t.Errorf("set: %v", fc)
	}
	edit(map[string]any{"forecast": map[string]any{"basis": ""}}, "forecast")
	if fc := planning("forecast"); fc["basis"] != nil || fc["duration"] != "6h" {
		t.Errorf("one key removed: %v", fc)
	}
	edit(map[string]any{"clear_forecast": true, "forecast": map[string]any{"delivery": "2026-10-10T09:00:00Z"}}, "forecast")
	if fc := planning("forecast"); fc["duration"] != nil || fc["delivery"] != "2026-10-10T09:00:00Z" {
		t.Errorf("replaced: %v", fc)
	}
	edit(map[string]any{"clear_forecast": true}, "forecast")
	if fc := planning("forecast"); fc != nil {
		t.Errorf("cleared: %v", fc)
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": f.story.Parent, "forecast": map[string]any{"duration": "6h"}}); !strings.Contains(failed, "a forecast belongs to stories") {
		t.Errorf("an epic's forecast: %q", failed)
	}
	if ed, failed := f.call(t, "item_edit", map[string]any{"id": f.story.Parent, "cost_of_delay": map[string]any{"value": 500}}); failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "cost_of_delay" {
		t.Errorf("an epic's cost of delay: %v %s", ed, failed)
	}
}

// I-0059: a task whose touches grow its story's claim into another story in
// progress names that story and the paths, through item_new and item_edit,
// and each of the two stories is told once, as an overlapped change caused by
// the other. A write that adds nothing inside another claim tells nobody.
func TestAWriteThatGrowsAClaimIntoAnotherTellsBothStories(t *testing.T) {
	f := setup(t)
	other := f.readyStory(t, "Other", t0, "docs/other")
	if _, err := f.repo.Transition(other, workitem.InProgress, "alex", "", t0); err != nil {
		t.Fatal(err)
	}
	if out, _ := f.call(t, "inbox", map[string]any{}); out == nil {
		t.Fatal("first look")
	}
	overlapped := func() map[string]map[string]any {
		t.Helper()
		out, _ := f.call(t, "inbox", map[string]any{})
		got := map[string]map[string]any{}
		for _, c := range out["changes"].([]any) {
			c := c.(map[string]any)
			if c["kind"] != Overlapped {
				continue
			}
			if got[c["id"].(string)] != nil {
				t.Errorf("told twice: %v", c)
			}
			got[c["id"].(string)] = c
		}
		return got
	}
	toldBoth := func(paths string) {
		t.Helper()
		got := overlapped()
		mine, theirs := got[f.story.ID], got[other.ID]
		if len(got) != 2 || mine == nil || theirs == nil {
			t.Fatalf("one overlapped change for each story: %v", got)
		}
		if mine["cause"] != other.ID || theirs["cause"] != f.story.ID || mine["to"] != paths || theirs["to"] != paths || mine["by"] != "claude" {
			t.Errorf("the changes: %v / %v", mine, theirs)
		}
		want := f.story.ID + "'s claim grew to overlap " + other.ID + "'s on " + strings.ReplaceAll(paths, ",", ", ") + ", written by claude"
		if s := theirs["summary"].(string); !strings.Contains(s, want) || !strings.Contains(s, "coordinate with "+f.story.ID+"'s agent before "+other.ID+" Other changes them") {
			t.Errorf("summary: %s", s)
		}
		if s := mine["summary"].(string); !strings.Contains(s, "coordinate with "+other.ID+"'s agent") || strings.Contains(s, "accepted") {
			t.Errorf("summary: %s", s)
		}
		if again := overlapped(); len(again) != 0 {
			t.Errorf("told once: %v", again)
		}
	}
	reached := func(out map[string]any, paths ...string) {
		t.Helper()
		list, _ := out["overlaps"].([]any)
		if len(list) != 1 {
			t.Fatalf("overlaps: %v", out["overlaps"])
		}
		o := list[0].(map[string]any)
		if o["story"] != other.ID || o["title"] != "Other" || strings.Join(toStrings(o["paths"]), ",") != strings.Join(paths, ",") {
			t.Errorf("overlap: %v", o)
		}
	}

	*f.clock = f.clock.Add(time.Minute)
	out, failed := f.call(t, "item_new", map[string]any{"type": "task", "title": "Reach", "parent": f.story.ID, "touches": []string{"docs/other/a.md", "flai/internal/mcpserver/x.go"}})
	if failed != "" {
		t.Fatal(failed)
	}
	reached(out, "docs/other/a.md")
	toldBoth("docs/other/a.md")

	*f.clock = f.clock.Add(time.Minute)
	ed, failed := f.call(t, "item_edit", map[string]any{"id": f.task.ID, "touches": []string{"docs/other/b.md", "docs/other/a.md"}})
	if failed != "" {
		t.Fatal(failed)
	}
	reached(ed, "docs/other/b.md")
	toldBoth("docs/other/b.md")

	// the story itself, in progress, widened: the write stands
	*f.clock = f.clock.Add(time.Minute)
	ed, failed = f.call(t, "item_edit", map[string]any{"id": f.story.ID, "touches": []string{"flai/internal/mcpserver", "docs/other/c.md"}})
	if failed != "" {
		t.Fatal(failed)
	}
	reached(ed, "docs/other/c.md")
	toldBoth("docs/other/c.md")

	// paths already claimed, or outside every other claim, tell nobody
	*f.clock = f.clock.Add(time.Minute)
	for _, touches := range [][]string{{"docs/other/a.md"}, {"design/system/plan.md"}} {
		ed, failed = f.call(t, "item_edit", map[string]any{"id": f.task.ID, "touches": touches})
		if failed != "" || ed["overlaps"] != nil {
			t.Errorf("%v: %v %s", touches, ed, failed)
		}
	}
	out, failed = f.call(t, "item_new", map[string]any{"type": "task", "title": "Apart", "parent": f.story.ID, "touches": []string{"flai/cmd"}})
	if failed != "" || out["overlaps"] != nil {
		t.Errorf("apart: %v %s", out, failed)
	}
	if got := overlapped(); len(got) != 0 {
		t.Errorf("nothing grew into another claim: %v", got)
	}
}

// An overlap notice written at acceptance by a flai that predates grown
// claims (I-0059) still reports as an acceptance.
func TestAnOldOverlapNoticeReportsAsBefore(t *testing.T) {
	f := setup(t)
	if out, _ := f.call(t, "inbox", map[string]any{}); out == nil {
		t.Fatal("first look")
	}
	*f.clock = f.clock.Add(time.Minute)
	line := `{"at":"` + f.clock.Format(workitem.TimeFormat) + `","by":"alex","id":"` + f.story.ID + `","title":"Open","accepted":"S-0099","paths":["flai/cmd/a.go"]}` + "\n"
	_ = os.MkdirAll(filepath.Dir(itemedit.OverlapsPath(f.repo)), 0o700)
	if err := os.WriteFile(itemedit.OverlapsPath(f.repo), []byte(line), 0o600); err != nil {
		t.Fatal(err)
	}
	out, _ := f.call(t, "inbox", map[string]any{})
	changes := out["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("one overlap: %v", changes)
	}
	c := changes[0].(map[string]any)
	if c["kind"] != Overlapped || c["cause"] != "S-0099" || c["to"] != "flai/cmd/a.go" || c["summary"] != "S-0099 was accepted by alex and changed flai/cmd/a.go, which "+f.story.ID+" Open claims. Run flai stream sync "+f.story.ID+" and the tests before you go on" {
		t.Errorf("the change: %v", c)
	}
}

// orchestrated is a server whose caller is the orchestrator, as flai serve
// starts it, and a maker of stories in its backlog under the fixture's epic:
// each with a goal, a criterion, and the touches given, a draft when draft
// is true, and with a forecast and a cost of delay when planned is true.
func orchestrated(t *testing.T) (*fixture, func(title string, draft, planned bool, touches ...string) string) {
	t.Helper()
	t.Setenv("FLAI_ROLE", "orchestrate")
	f := setupWith(t, func(o *Options) { o.Agent = "orchestrator" })
	story := func(title string, draft, planned bool, touches ...string) string {
		t.Helper()
		args := map[string]any{"type": "story", "title": title, "parent": f.story.Parent, "draft": draft, "body": "## Goal\n\nx\n\n## Acceptance criteria\n- [ ] it works\n\n## Tasks\n\n## Notes\n"}
		if len(touches) > 0 {
			args["touches"] = touches
		}
		out, failed := f.call(t, "item_new", args)
		if failed != "" {
			t.Fatalf("%s: %s", title, failed)
		}
		id := out["id"].(string)
		if planned {
			if _, failed := f.call(t, "item_edit", map[string]any{"id": id, "cost_of_delay": map[string]any{"value": 300}, "forecast": map[string]any{"duration": "2h", "delivery": "2026-09-20T12:00:00Z"}}); failed != "" {
				t.Fatalf("%s planned: %s", title, failed)
			}
		}
		return id
	}
	return f, story
}

// S-0219: the orchestrator finalizes a draft with item_edit's draft false and
// nothing else, and only a complete one, which then names it as who
// finalized it; an incomplete one is refused naming what it lacks.
func TestTheOrchestratorFinalizesACompleteDraft(t *testing.T) {
	f, story := orchestrated(t)
	complete := story("Complete", true, true, "docs")
	incomplete := story("Incomplete", true, false)
	if _, failed := f.call(t, "item_edit", map[string]any{"id": complete, "draft": false, "title": "Renamed"}); !strings.Contains(failed, complete+": the orchestrator finalizes a draft with draft false and nothing else") {
		t.Errorf("with another change: %q", failed)
	}
	if _, failed := f.call(t, "item_edit", map[string]any{"id": incomplete, "draft": false}); !strings.Contains(failed, incomplete+" is not complete, so the orchestrator does not finalize it (flai promote --drafts): no touches; no forecast duration; no forecast delivery; no cost of delay value") {
		t.Errorf("an incomplete draft: %q", failed)
	}
	for _, id := range []string{complete, incomplete} {
		if it, _ := f.repo.Get(id); !it.Draft || it.Finalized != nil {
			t.Errorf("%s finalized by a refused call: %+v", id, it.Finalized)
		}
	}
	ed, failed := f.call(t, "item_edit", map[string]any{"id": complete, "draft": false})
	if failed != "" || strings.Join(toStrings(ed["changed"]), ",") != "draft" {
		t.Fatalf("a complete draft: %v %s", ed, failed)
	}
	it, _ := f.repo.Get(complete)
	if it.Draft || it.Finalized == nil || it.Finalized.By != "orchestrator" || it.Finalized.At != f.clock.UTC().Format(workitem.TimeFormat) {
		t.Errorf("finalized = %+v, want the orchestrator, now", it.Finalized)
	}
}

// S-0219: the orchestrator moves to ready only a story flai promote
// --candidates lists, and only while ready is under its WIP limit; a held
// story, a draft, and any other item are refused with the reason.
func TestTheOrchestratorPromotesACandidateWhileReadyHasRoom(t *testing.T) {
	f, story := orchestrated(t)
	candidate := story("Candidate", false, true, "docs/a")
	next := story("Next", false, true, "docs/b")
	held := story("Held", false, true, "flai/internal/mcpserver/x")
	drafted := story("Drafted", true, true, "docs/c")
	unplanned := story("Unplanned", false, false, "docs/d")
	for id, want := range map[string]string{
		held:      held + " is not a candidate to go to ready (flai promote --candidates): held (overlap)",
		drafted:   drafted + " is not a candidate to go to ready (flai promote --candidates): draft: finalize it first",
		unplanned: unplanned + " is not a candidate to go to ready (flai promote --candidates): no forecast duration; no cost of delay value",
		f.task.ID: f.task.ID + " is a task: the orchestrator moves a story to ready",
	} {
		if _, failed := f.call(t, "item_move", map[string]any{"id": id, "to": "ready"}); !strings.Contains(failed, want) {
			t.Errorf("%s: %q, want %q", id, failed, want)
		}
	}
	if out, failed := f.call(t, "item_move", map[string]any{"id": candidate, "to": "ready"}); failed != "" || out["status"] != "ready" {
		t.Fatalf("a candidate: %v %s", out, failed)
	}
	b, err := f.repo.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if err := b.SetLimit(workitem.Ready, 1); err != nil {
		t.Fatal(err)
	}
	if err := b.Save("2026-09-18"); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "item_move", map[string]any{"id": next, "to": "ready"}); !strings.Contains(failed, "the ready column is at its WIP limit (1 of 1): the orchestrator moves no story to ready until one leaves it") {
		t.Errorf("a candidate with ready full: %q", failed)
	}
	if it, _ := f.repo.Get(next); it.Status != workitem.Backlog {
		t.Errorf("%s moved: %s", next, it.Status)
	}
}
