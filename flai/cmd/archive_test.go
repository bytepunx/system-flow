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
