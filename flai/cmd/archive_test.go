package cmd

import (
	"encoding/json"
	"strings"
	"testing"
)

// I-0073: flai archive resolves the threads still open or answered on the
// items it archives (a story and its tasks), naming each as archived, and
// leaves a thread on an item it does not archive open. The dry run names the
// threads and resolves none.
func TestArchiveResolvesTheThreadsOnWhatItArchives(t *testing.T) {
	t.Setenv("FLAI_CONFIG", t.TempDir()+"/cfg.json")
	t.Setenv("FLAI_AGENT", "alex")
	root := tempProject(t)
	// the author flai archive resolves as is the config's, which a fresh config fills with this
	// host's user: pin it, so the test does not depend on who runs it (I-0013)
	for _, args := range [][]string{
		{"config", "set", "author", "alex"},
		{"story", "new", "Slice"},
		{"task", "new", "--story", "S-0001", "Piece"},
		{"story", "new", "Other"},
		{"thread", "new", "--on", "S-0001", "--by", "alex", "Is the flag right?", "Say which."},
		{"thread", "new", "--on", "T-0001", "--by", "claude", "Is the task done?", "Say so."},
		{"thread", "new", "--on", "S-0002", "--by", "claude", "When is it next?", "Say when."},
		{"thread", "reply", "TH-0001", "--by", "claude", "Use the flag."},
		{"move", "S-0001", "cancelled", "--reason", "not needed"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}

	out, errOut, code := runIn(t, root, "archive", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the dry run: %d %s", code, errOut)
	}
	var res struct {
		Archived []string `json:"archived"`
		Resolved []string `json:"resolved_threads"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if got := strings.Join(res.Resolved, " "); got != "TH-0001 TH-0002" {
		t.Errorf("the dry run names the threads it would resolve: %q", got)
	}
	if out, _, _ := runIn(t, root, "archive", "S-0001", "--dry-run"); !strings.Contains(out, "would resolve TH-0001, TH-0002, open or answered on what it archives\n") {
		t.Errorf("the dry run's text names them:\n%s", out)
	}
	for _, id := range []string{"TH-0001", "TH-0002"} {
		if th := thread(t, root, id); !th.Open() {
			t.Errorf("the dry run resolves nothing: %s is %s", id, th.Status)
		}
	}

	out, errOut, code = runIn(t, root, "archive", "S-0001")
	if code != 0 {
		t.Fatalf("the archive: %d %s\n%s", code, errOut, out)
	}
	if !strings.Contains(out, "resolved TH-0001, TH-0002, open or answered on what it archives\n") {
		t.Errorf("the output names the threads resolved:\n%s", out)
	}
	for id, item := range map[string]string{"TH-0001": "S-0001", "TH-0002": "T-0001"} {
		th := thread(t, root, id)
		es := th.Entries()
		last := es[len(es)-1]
		if th.Status != "resolved" || last.Author != "alex" || last.Text != "Resolved: "+item+" was archived" {
			t.Errorf("%s is resolved as %s archived: status %s, last entry %+v", id, item, th.Status, last)
		}
	}
	if th := thread(t, root, "TH-0003"); th.Status != "open" {
		t.Errorf("a thread on another item is left open: %s", th.Status)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "threads.archived") {
		t.Errorf("flai check finds no thread open on an archived item:\n%s", out)
	}
}

// ADR-0120: flai archive closes the conversations still open of the stories
// it archives, each with the entry "Closed: S-0001 was archived" by the
// author, and leaves a conversation between two other stories open. The dry
// run names the conversations and closes none.
func TestArchiveClosesTheConversationsOfWhatItArchives(t *testing.T) {
	t.Setenv("FLAI_CONFIG", t.TempDir()+"/cfg.json")
	t.Setenv("FLAI_AGENT", "alex")
	root := tempProject(t)
	for _, args := range [][]string{
		{"config", "set", "author", "alex"},
		{"story", "new", "Slice"},
		{"story", "new", "Other"},
		{"story", "new", "Third"},
		{"move", "S-0001", "cancelled", "--reason", "not needed"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	writeConversation(t, root, "MS-0001", "S-0002", "S-0001", "open")
	writeConversation(t, root, "MS-0002", "S-0002", "S-0003", "open")

	out, errOut, code := runIn(t, root, "archive", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the dry run: %d %s", code, errOut)
	}
	var res struct {
		Closed []string `json:"closed_messages"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if got := strings.Join(res.Closed, " "); got != "MS-0001" {
		t.Errorf("the dry run names the conversations it would close: %q", got)
	}
	if out, _, _ := runIn(t, root, "archive", "S-0001", "--dry-run"); !strings.Contains(out, "would close MS-0001, open with S-0002\n") {
		t.Errorf("the dry run's text names it with the story it is open with:\n%s", out)
	}
	if c := conversation(t, root, "MS-0001"); c.Status != "open" {
		t.Errorf("the dry run closes nothing: MS-0001 is %s", c.Status)
	}

	out, errOut, code = runIn(t, root, "archive", "S-0001", "--json")
	if code != 0 {
		t.Fatalf("the archive: %d %s\n%s", code, errOut, out)
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || strings.Join(res.Closed, " ") != "MS-0001" {
		t.Errorf("--json lists the conversations closed in closed_messages: %v %s", err, out)
	}
	assertClosedAs(t, root, "Closed: S-0001 was archived", "MS-0001")
	if c := conversation(t, root, "MS-0002"); c.Status != "open" || len(c.Entries()) != 1 {
		t.Errorf("a conversation between two other stories is left open: %s %+v", c.Status, c.Entries())
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "messages.") {
		t.Errorf("flai check finds nothing in the conversations:\n%s", out)
	}
}

// The text of flai archive names each conversation it closes with the story
// it was open with, or both stories when it archives both.
func TestArchiveNamesTheConversationsItCloses(t *testing.T) {
	t.Setenv("FLAI_CONFIG", t.TempDir()+"/cfg.json")
	t.Setenv("FLAI_AGENT", "alex")
	root := tempProject(t)
	for _, args := range [][]string{
		{"config", "set", "author", "alex"},
		{"story", "new", "Slice"},
		{"story", "new", "Other"},
		{"story", "new", "Third"},
		{"move", "S-0001", "cancelled", "--reason", "not needed"},
		{"move", "S-0002", "cancelled", "--reason", "not needed"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	writeConversation(t, root, "MS-0001", "S-0001", "S-0002", "open")
	writeConversation(t, root, "MS-0002", "S-0003", "S-0002", "open")
	out, errOut, code := runIn(t, root, "archive", "S-0001", "S-0002")
	if code != 0 || !strings.Contains(out, "closed MS-0001, open between S-0001 and S-0002\nclosed MS-0002, open with S-0003\n") {
		t.Fatalf("the archive names the conversations it closed: %d %s\n%s", code, errOut, out)
	}
	assertClosedAs(t, root, "Closed: S-0001 and S-0002 were archived", "MS-0001")
	assertClosedAs(t, root, "Closed: S-0002 was archived", "MS-0002")
}
