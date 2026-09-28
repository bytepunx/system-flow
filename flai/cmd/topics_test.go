package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0135: flai story new and flai epic new take --topics, flai edit sets,
// shows, and clears them, and a task has no such flag.
func TestTopicsOnStoriesAndEpics(t *testing.T) {
	root := heldProject(t)
	if _, errOut, code := runIn(t, root, "epic", "new", "Logs", "--topics", "logging"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Ship", "--epic", "E-0002", "--topics", "release,logging"); code != 0 {
		t.Fatal(errOut)
	}
	story, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", "S-0003-*.md"))
	if item, _ := os.ReadFile(story[0]); !strings.Contains(string(item), "\ntopics: [release, logging]\n") {
		t.Errorf("the story:\n%s", item)
	}
	epic, _ := filepath.Glob(filepath.Join(root, "wip/kanban/epics", "E-0002-*.md"))
	if item, _ := os.ReadFile(epic[0]); !strings.Contains(string(item), "\ntopics: [logging]\n") {
		t.Errorf("the epic:\n%s", item)
	}
	if _, errOut, code := runIn(t, root, "task", "new", "T", "--story", "S-0003", "--topics", "logging"); code == 0 || !strings.Contains(errOut, "unknown flag: --topics") {
		t.Errorf("a task takes no topics: %d %s", code, errOut)
	}

	out, errOut, code := runIn(t, root, "edit", "S-0003", "--topics", "release")
	if code != 0 || !strings.Contains(out, "S-0003: changed topics") {
		t.Fatalf("edit: %d %s %s", code, out, errOut)
	}
	if out, _, _ := runIn(t, root, "edit", "S-0003", "--show"); !strings.Contains(out, "  topics: release\n") {
		t.Errorf("show:\n%s", out)
	}
	for _, c := range []struct {
		args []string
		says string
	}{
		{[]string{"edit", "S-0003", "--topics", "a/b"}, "not one word"},
		{[]string{"edit", "S-0003", "--topics", "x", "--clear-topics"}, "contradict"},
		{[]string{"story", "new", "Bad", "--topics", "two words"}, "not one word"},
	} {
		if out, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(out+errOut, c.says) {
			t.Errorf("%v: %d\n%s%s", c.args, code, out, errOut)
		}
	}
	if out, errOut, code := runIn(t, root, "edit", "E-0002", "--clear-topics"); code != 0 || !strings.Contains(out, "changed topics") {
		t.Fatalf("clear: %d %s %s", code, out, errOut)
	}
	if item, _ := os.ReadFile(epic[0]); strings.Contains(string(item), "topics:") {
		t.Errorf("topics stay:\n%s", item)
	}
}
