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
