package mcpserver

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// S-0103: stories are created and edited over MCP with their agent: the
// project's default filled in, what is given winning, and an edit that
// replaces or removes it; a stale hash is refused, not overwritten.
func TestItemNewAndEditCarryTheAgent(t *testing.T) {
	f := setup(t)
	f.repo.Manifest.Agent = &manifest.Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high"}}

	out, failed := f.call(t, "item_new", map[string]any{"type": "story", "title": "From an agent", "agent": map[string]any{"model": "claude-sonnet-5"}, "body": "## Goal\n\nSomething.\n\n## Acceptance criteria\n- [ ] it works\n"})
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
