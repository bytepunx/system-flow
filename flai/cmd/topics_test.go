package cmd

import (
	"encoding/json"
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

// S-0135: flai show prints a story's topics with their sources, and --json
// returns every source; an epic shows its own.
func TestShowPrintsAStorysTopics(t *testing.T) {
	root := heldProject(t)
	if _, errOut, code := runIn(t, root, "edit", "E-0001", "--topics", "release"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "edit", "S-0002", "--topics", "logging"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code := runIn(t, root, "show", "S-0002")
	if code != 0 {
		t.Fatal(errOut)
	}
	want := "  topics:\n    logging  own\n    release  epic E-0001\n    flai     flai by touches cli (S-0002)\n    cli      flai by touches cli (S-0002)\n    go       flai by touches cli (S-0002)\n    code     flai is code\n    all      every story\n"
	if !strings.Contains(out, want) {
		t.Errorf("show:\n%s\nwant:\n%s", out, want)
	}
	out, _, _ = runIn(t, root, "show", "S-0002", "--json")
	var view struct {
		Item   struct{ Topics []string }
		Topics []struct {
			Topic   string
			Sources []struct{ Kind, Item, Project, Via string }
		}
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatal(err)
	}
	if strings.Join(view.Item.Topics, ",") != "logging" || len(view.Topics) != 7 || view.Topics[2].Topic != "flai" ||
		view.Topics[2].Sources[0] != (struct{ Kind, Item, Project, Via string }{"claim", "S-0002", "flai", "cli"}) {
		t.Errorf("json: %+v", view)
	}
	if out, _, _ := runIn(t, root, "show", "E-0001"); !strings.Contains(out, "  topics: release\n") {
		t.Errorf("epic:\n%s", out)
	}
}
