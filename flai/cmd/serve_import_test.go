package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0098: the folders named for import are kept in the host's configuration,
// and list says what is found in them.
func TestServeImportNamesFoldersAndListsWhatIsInThem(t *testing.T) {
	home := t.TempDir()
	cfg := filepath.Join(home, "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, "widget", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	if out, _, _ := runIn(t, home, "serve", "import", "list"); !strings.Contains(out, "no folders named") {
		t.Errorf("empty: %s", out)
	}
	if _, errOut, code := runIn(t, home, "serve", "import", "add", filepath.Join(folder, "nope")); code == 0 || !strings.Contains(errOut, "not a folder") {
		t.Errorf("a missing folder: %d %s", code, errOut)
	}
	if out, _, code := runIn(t, home, "serve", "import", "add", folder); code != 0 || !strings.Contains(out, "offered for import") {
		t.Fatalf("add: %d %s", code, out)
	}
	runIn(t, home, "serve", "import", "add", folder) // twice is once
	out, _, _ := runIn(t, home, "serve", "import", "list", "--json")
	var got struct {
		Roots      []string `json:"roots"`
		Candidates []struct {
			Key  string `json:"key"`
			Root string `json:"root"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || len(got.Roots) != 1 || len(got.Candidates) != 1 || got.Candidates[0].Key != "import-widget" {
		t.Fatalf("list: %v %s", err, out)
	}
	if text, _, _ := runIn(t, home, "serve", "import", "list"); !strings.Contains(text, "widget  "+filepath.Join(folder, "widget")) {
		t.Errorf("list text: %s", text)
	}
	runIn(t, home, "serve", "import", "remove", folder)
	if out, _, _ := runIn(t, home, "serve", "import", "list"); !strings.Contains(out, "no folders named") {
		t.Errorf("after remove: %s", out)
	}
}
