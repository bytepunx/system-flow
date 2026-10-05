package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/planning"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0210: flai touches suggest is a subcommand, while flai touches <id>
// <path> still sets touches; it refuses what is not a story, and a story with
// nothing to start from, before it reads git.
func TestTouchesSuggestRoutesAndRefuses(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, step := range [][]string{{"epic", "new", "E"}, {"story", "new", "S", "--epic", "E-0001"}, {"task", "new", "T", "--story", "S-0001"}} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	if out, errOut, code := runIn(t, root, "touches", "S-0001", "a/b/"); code != 0 || out != "S-0001 touches a/b\n" {
		t.Errorf("touches still sets: %d %q %s", code, out, errOut)
	}
	for _, id := range []string{"E-0001", "T-0001"} {
		if _, errOut, code := runIn(t, root, "touches", "suggest", id); code == 0 || !strings.Contains(errOut, id+" is not a story") || !strings.Contains(errOut, "flai touches suggest S-nnnn <path>") {
			t.Errorf("suggest %s: %d %s", id, code, errOut)
		}
	}
	if _, _, code := runIn(t, root, "touches", "S-0001", "--clear"); code != 0 {
		t.Fatal("clear")
	}
	if _, errOut, code := runIn(t, root, "touches", "suggest", "S-1"); code == 0 || !strings.Contains(errOut, "S-0001 touches nothing and no path was given") || !strings.Contains(errOut, "flai touches suggest S-0001 <path>") {
		t.Errorf("suggest with no seeds: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "touches", "suggest"); code == 0 || !strings.Contains(errOut, "requires at least 1 arg") {
		t.Errorf("suggest with no ID: %d %s", code, errOut)
	}
}

// S-0211: setting a story's touches leaves an edit notice, as flai edit does,
// naming touches and who set them, so that flai serve plans the story again;
// showing them, or setting what is already there, leaves none.
func TestTouchesLeavesAnEditNotice(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "alex")
	root := tempProject(t)
	for _, step := range [][]string{{"epic", "new", "E"}, {"story", "new", "S", "--epic", "E-0001"}, {"touches", "S-0001", "a/b"}, {"touches", "S-0001"}, {"touches", "S-0001", "a/b"}} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	var touched []itemedit.Notice
	for _, n := range itemedit.Notices(repo) {
		if slices.Contains(n.Changed, "touches") {
			touched = append(touched, n)
		}
	}
	if len(touched) != 1 || touched[0].ID != "S-0001" || touched[0].By != "alex" || touched[0].Type != workitem.Story || touched[0].At == "" {
		t.Errorf("one touches notice by alex on S-0001, got %+v", touched)
	}
}

// I-0071: flai touches takes a path that starts with a dot, as flai edit
// does, and refuses one that leaves the repository, writing nothing.
func TestTouchesTakesADotPathAndRefusesOneOutside(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, step := range [][]string{{"epic", "new", "E"}, {"story", "new", "S", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	if out, errOut, code := runIn(t, root, "touches", "S-0001", ".github/workflows/", "flai", ".github/workflows"); code != 0 || out != "S-0001 touches .github/workflows, flai\n" {
		t.Errorf("a dot path: %d %q %s", code, out, errOut)
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-s.md")
	before := read(t, file)
	for _, touch := range []string{"../x", "/etc", "a,b"} {
		if _, errOut, code := runIn(t, root, "touches", "S-0001", touch); code == 0 || !strings.Contains(errOut, "give a repository path or component name") {
			t.Errorf("touch %q: %d %s", touch, code, errOut)
		}
	}
	if read(t, file) != before {
		t.Error("no refusal wrote anything")
	}
}

// S-0210: the seeds are the story's touches, its own and its tasks' not
// cancelled, a component named by a tag read as its path, and the paths
// given; the files changed with them on the main branch are counted.
func TestTouchesSuggestCountsFilesChangedWithTheSeeds(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\n    tags: [command]\n"), 0o644)
	r := execx.System{}
	gitRun := func(args ...string) {
		t.Helper()
		if out, err := r.Run(root, "git", args...); err != nil {
			t.Fatalf("git %v: %s %v", args, out, err)
		}
	}
	commit := func(subject string, files ...string) {
		t.Helper()
		for _, f := range files {
			_ = os.MkdirAll(filepath.Join(root, filepath.Dir(f)), 0o755)
			_ = os.WriteFile(filepath.Join(root, f), []byte(subject+"\n"), 0o644)
		}
		gitRun("add", "-A")
		gitRun("commit", "-q", "-m", subject)
	}
	gitRun("init", "-q", "-b", "main")
	gitRun("config", "user.email", "t@t")
	gitRun("config", "user.name", "t")
	// flai's cache, the edit notices among it, is ignored as a project's is
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	commit("init")
	for _, step := range [][]string{
		{"story", "new", "S"}, {"task", "new", "T", "--story", "S-0001"}, {"task", "new", "U", "--story", "S-0001"},
		{"touches", "S-0001", "command"}, {"touches", "T-0001", "docs/guide.md"}, {"touches", "T-0002", "web"},
		{"move", "T-0002", "cancelled", "--reason", "dropped", "--yes"},
	} {
		if _, errOut, code := runIn(t, root, step...); code != 0 {
			t.Fatalf("%v: %s", step, errOut)
		}
	}
	commit("feat: [S-0001] a", "cli/a.go", "lib/x.go", "lib/y.go")
	commit("feat: b", "cli/b.go", "lib/x.go")
	commit("fix: c", "docs/guide.md", "lib/z.go")
	commit("feat: d", "web/w.go", "lib/z.go")
	commit("chore: publish flai 1.0.0", "cli/CHANGELOG.md", "CHANGELOG.md")

	out, errOut, code := runIn(t, root, "touches", "suggest", "S-0001", "--min", "1", "--json")
	if code != 0 {
		t.Fatalf("suggest --json: %d %s", code, errOut)
	}
	var got touchesSuggestion
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	want := touchesSuggestion{ID: "S-0001", Seeds: []string{"cli", "docs/guide.md"}, Suggestions: planning.Suggestions{
		Commits: 5, SeedCommits: 3, Suggestions: []planning.Suggestion{
			{Path: "lib/x.go", Count: 2, Share: 2.0 / 3},
			{Path: "lib/y.go", Count: 1, Share: 1.0 / 3},
			{Path: "lib/z.go", Count: 1, Share: 1.0 / 3},
		},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("suggest --json:\n got %+v\nwant %+v", got, want)
	}
	for _, key := range []string{`"id"`, `"seeds"`, `"commits"`, `"seed_commits"`, `"suggestions"`, `"path"`, `"count"`, `"share"`} {
		if !strings.Contains(out, key) {
			t.Errorf("suggest --json has no %s: %s", key, out)
		}
	}

	out, errOut, code = runIn(t, root, "touches", "suggest", "S-1")
	if wantOut := "S-0001 from cli, docs/guide.md: 3 of 5 commits changed them\nlib/x.go  2   67%\n"; code != 0 || out != wantOut {
		t.Errorf("suggest: %d %q %s, want %q", code, out, errOut, wantOut)
	}
	out, errOut, code = runIn(t, root, "touches", "suggest", "S-0001", "web/", "--min", "1", "--limit", "1")
	if wantOut := "S-0001 from cli, docs/guide.md, web: 4 of 5 commits changed them\nlib/x.go  2   50%\n"; code != 0 || out != wantOut {
		t.Errorf("suggest with a path and a limit: %d %q %s, want %q", code, out, errOut, wantOut)
	}
	out, errOut, code = runIn(t, root, "touches", "suggest", "S-0001", "--min", "3")
	if code != 0 || !strings.Contains(out, "no other file was changed with them often enough; lower --min") {
		t.Errorf("suggest with none: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "touches", "suggest", "S-0001", "flaiover")
	if code != 0 || !strings.Contains(out, "S-0001 from cli, docs/guide.md, flaiover: 3 of 5") {
		t.Errorf("suggest with a path that changed nothing: %d %q %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Other"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code = runIn(t, root, "touches", "suggest", "S-0002", "flaiover")
	if code != 0 || !strings.Contains(out, "S-0002 from flaiover: 0 of 5 commits changed them\nno commit on the main branch changed them") {
		t.Errorf("suggest with no seed commit: %d %q %s", code, out, errOut)
	}
}
