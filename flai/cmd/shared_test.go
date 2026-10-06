package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// sharedProject is a scratch project with the component flai (tag cli), a
// story S-0001 touching a file, a component, and a folder, with a task
// touching a file outside them, and a story S-0002 touching nothing.
func sharedProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, _ := os.ReadFile(file)
	_ = os.WriteFile(file, append(data, "# the project's components\nprojects:\n  - name: flai\n    path: flai\n    kind: go\n    tags: [cli]\n"...), 0o644)
	for _, args := range [][]string{
		{"story", "new", "Claims", "--touches", "docs/users/flai.md,cli,design/adrs"},
		{"story", "new", "Bare"},
		{"task", "new", "--story", "S-0001", "Piece", "--touches", "template/CHANGELOG.md"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	return root
}

// sharedRun runs flai in root and fails the test when it exits non-zero.
func sharedRun(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, errOut, code := runIn(t, root, args...)
	if code != 0 {
		t.Fatalf("flai %v exited %d: %s", args, code, errOut)
	}
	return out
}

func readManifestShared(t *testing.T, root string) []string {
	t.Helper()
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	return m.Claims.Shared
}

// S-0295: flai shared list prints the patterns, and [] with --json when
// there are none.
func TestSharedList(t *testing.T) {
	root := sharedProject(t)
	if out := sharedRun(t, root, "shared", "list"); !strings.Contains(out, "no shared paths") {
		t.Errorf("an empty list says so, got %q", out)
	}
	if out := sharedRun(t, root, "shared", "list", "--json"); out != "[]\n" {
		t.Errorf("an empty list is [] with --json, got %q", out)
	}
	sharedRun(t, root, "shared", "add", "design/adrs", "docs/users/*.md")
	if out := sharedRun(t, root, "shared", "list"); out != "design/adrs\ndocs/users/*.md\n" {
		t.Errorf("list prints one pattern a line, got %q", out)
	}
	var got []string
	if err := json.Unmarshal([]byte(sharedRun(t, root, "shared", "list", "--json")), &got); err != nil {
		t.Fatal(err)
	}
	if want := []string{"design/adrs", "docs/users/*.md"}; !reflect.DeepEqual(got, want) {
		t.Errorf("list --json = %v, want %v", got, want)
	}
}

// S-0295: a pattern that is not valid is listed, and warned about on stderr.
func TestSharedListWarnsOfAnInvalidPattern(t *testing.T) {
	root := sharedProject(t)
	file := filepath.Join(root, manifest.File)
	data, _ := os.ReadFile(file)
	_ = os.WriteFile(file, append(data, "claims:\n  shared:\n    - /abs\n"...), 0o644)
	out, errOut, code := runIn(t, root, "shared", "list")
	if code != 0 || out != "/abs\n" {
		t.Fatalf("list: %d %q %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "not valid") || !strings.Contains(errOut, "absolute") {
		t.Errorf("an invalid pattern is warned about, got %q", errOut)
	}
}

// S-0295: flai shared add and remove change claims.shared, keep the rest of
// the manifest, say what changed, and refuse what they cannot do with the
// reason, writing nothing.
func TestSharedAddRemove(t *testing.T) {
	root := sharedProject(t)
	if out := sharedRun(t, root, "shared", "add", "design/adrs", "docs/users/*.md"); out != "added design/adrs\nadded docs/users/*.md\n" {
		t.Errorf("add says what it added, got %q", out)
	}
	var ch manifest.Change
	if err := json.Unmarshal([]byte(sharedRun(t, root, "shared", "add", "template/**", "--json")), &ch); err != nil {
		t.Fatal(err)
	}
	want := manifest.Change{Added: []string{"template/**"}, Shared: []string{"design/adrs", "docs/users/*.md", "template/**"}}
	if !reflect.DeepEqual(ch, want) {
		t.Errorf("add --json = %+v, want %+v", ch, want)
	}
	data, _ := os.ReadFile(filepath.Join(root, manifest.File))
	if !strings.Contains(string(data), "# the project's components\nprojects:\n") {
		t.Errorf("add kept no other key or comment:\n%s", data)
	}

	for _, tt := range []struct {
		args []string
		why  string
	}{
		{[]string{"add", "/design"}, "absolute"},
		{[]string{"add", "design/../x"}, ".."},
		{[]string{"add", "design/["}, "not a valid glob"},
		{[]string{"add", "flai", "design/adrs"}, "in the list already"},
		{[]string{"remove", "docs"}, "not in the list"},
	} {
		before, _ := os.ReadFile(filepath.Join(root, manifest.File))
		out, errOut, code := runIn(t, root, append([]string{"shared"}, tt.args...)...)
		if code == 0 {
			t.Errorf("flai shared %v was not refused: %s", tt.args, out)
		}
		if !strings.Contains(errOut, tt.why) {
			t.Errorf("flai shared %v: the refusal does not say %q: %s", tt.args, tt.why, errOut)
		}
		if after, _ := os.ReadFile(filepath.Join(root, manifest.File)); string(after) != string(before) {
			t.Errorf("flai shared %v wrote the manifest though refused", tt.args)
		}
	}

	if out := sharedRun(t, root, "shared", "remove", "template/**"); out != "removed template/**\n" {
		t.Errorf("remove says what it removed, got %q", out)
	}
	ch = manifest.Change{}
	if err := json.Unmarshal([]byte(sharedRun(t, root, "shared", "remove", "design/adrs", "--json")), &ch); err != nil {
		t.Fatal(err)
	}
	want = manifest.Change{Removed: []string{"design/adrs"}, Shared: []string{"docs/users/*.md"}}
	if !reflect.DeepEqual(ch, want) {
		t.Errorf("remove --json = %+v, want %+v", ch, want)
	}
	if got := readManifestShared(t, root); !reflect.DeepEqual(got, []string{"docs/users/*.md"}) {
		t.Errorf("the manifest holds %v after the edits", got)
	}
	if _, _, code := runIn(t, root, "shared", "add"); code == 0 {
		t.Error("add without a pattern was not refused")
	}
}

// S-0295: flai shared check says of each path, touches entry, or story's
// claim entry whether it lies inside a pattern, and which.
func TestSharedCheck(t *testing.T) {
	root := sharedProject(t)
	sharedRun(t, root, "shared", "add", "design/adrs", "docs/users/*.md")

	out := sharedRun(t, root, "shared", "check", "docs/users/flai.md", "design/adrs/", "docs/users", "flai/cmd", "cli")
	for _, line := range []string{
		"docs/users/flai.md: shared, inside docs/users/*.md\n",
		"design/adrs/ (design/adrs): shared, inside design/adrs\n",
		"docs/users: not shared\n",
		"flai/cmd: not shared\n",
		"cli (flai): not shared\n",
	} {
		if !strings.Contains(out, line) {
			t.Errorf("check does not print %q:\n%s", line, out)
		}
	}

	out = sharedRun(t, root, "shared", "check", "s-1", "S-02")
	for _, line := range []string{
		"S-0001 docs/users/flai.md: shared, inside docs/users/*.md\n",
		"S-0001 flai: not shared\n",
		"S-0001 design/adrs: shared, inside design/adrs\n",
		"S-0001 template/CHANGELOG.md: not shared\n",
		"S-0002 claims nothing",
	} {
		if !strings.Contains(out, line) {
			t.Errorf("check of stories does not print %q:\n%s", line, out)
		}
	}

	var got []sharedCheck
	if err := json.Unmarshal([]byte(sharedRun(t, root, "shared", "check", "--json", "docs/users/flai.md", "docs/users", "cli", "S-0001")), &got); err != nil {
		t.Fatal(err)
	}
	want := []sharedCheck{
		{Entry: "docs/users/flai.md", Path: "docs/users/flai.md", Shared: true, Pattern: "docs/users/*.md"},
		{Entry: "docs/users", Path: "docs/users"},
		{Entry: "cli", Path: "flai"},
		{Entry: "docs/users/flai.md", Path: "docs/users/flai.md", Shared: true, Pattern: "docs/users/*.md", Story: "S-0001"},
		{Entry: "flai", Path: "flai", Story: "S-0001"},
		{Entry: "design/adrs", Path: "design/adrs", Shared: true, Pattern: "design/adrs", Story: "S-0001"},
		{Entry: "template/CHANGELOG.md", Path: "template/CHANGELOG.md", Story: "S-0001"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("check --json =\n%+v\nwant\n%+v", got, want)
	}
	if out := sharedRun(t, root, "shared", "check", "S-0002", "--json"); out != "[]\n" {
		t.Errorf("a story that claims nothing gives [] with --json, got %q", out)
	}

	for _, tt := range []struct {
		arg, why string
	}{
		{"T-0001", "not a story"},
		{"S-0099", "S-0099"},
		{" ", "empty"},
	} {
		if _, errOut, code := runIn(t, root, "shared", "check", tt.arg); code == 0 || !strings.Contains(errOut, tt.why) {
			t.Errorf("check %q: exit %d, want a refusal saying %q: %s", tt.arg, code, tt.why, errOut)
		}
	}
}

// S-0295: the help of flai shared and of each subcommand gives the glob
// dialect.
func TestSharedHelpGivesTheDialect(t *testing.T) {
	root := sharedProject(t)
	for _, args := range [][]string{{"shared"}, {"shared", "list"}, {"shared", "add"}, {"shared", "remove"}, {"shared", "check"}} {
		out := sharedRun(t, root, append(args, "--help")...)
		for _, want := range []string{"* is any", "** zero or more whole segments", "? one character", "a plain path covers itself and everything below it"} {
			if !strings.Contains(out, want) {
				t.Errorf("flai %v --help does not say %q", args, want)
			}
		}
	}
}

// S-0295: with --autocommit, an edit of the shared paths is committed,
// system-flow.yaml alone, with its trailers; without it, nothing is.
func TestSharedAddAutocommitsTheManifestAlone(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := sharedProject(t)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "start")
	_ = os.WriteFile(filepath.Join(root, "notes.txt"), []byte("not mine\n"), 0o644)

	sharedRun(t, root, "shared", "add", "design/adrs")
	if st := gitIn(t, root, "status", "--porcelain", "--", manifest.File); !strings.Contains(st, manifest.File) {
		t.Fatalf("without --autocommit the manifest stays uncommitted: %q", st)
	}
	sharedRun(t, root, "shared", "add", "--autocommit", "--trailer", "Co-Authored-By: someone <s@s>", "design/issues")
	if st := gitIn(t, root, "status", "--porcelain", "--", manifest.File); st != "" {
		t.Errorf("the manifest is committed: %q", st)
	}
	msg := gitIn(t, root, "log", "-1", "--format=%B")
	if !strings.Contains(msg, "chore: add to the shared paths: design/issues") || !strings.Contains(msg, "Co-Authored-By: someone <s@s>") {
		t.Errorf("commit message: %q", msg)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); strings.TrimSpace(files) != manifest.File {
		t.Errorf("the commit holds the manifest alone: %q", files)
	}
}
