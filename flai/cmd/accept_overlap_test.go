package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Which open stories acceptance tells, and of which paths (S-0132): a claim
// covers a changed path by the prefix rule, a component name counts as its
// path, and an empty claim covers everything.
func TestOverlapsOf(t *testing.T) {
	t.Parallel()
	story := func(id, status string, touches ...string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: status, Title: id, Touches: touches}
	}
	task := &workitem.Item{ID: "T-0001", Type: workitem.Task, Status: workitem.InProgress, Parent: "S-0005", Touches: []string{"design/system"}}
	items := []*workitem.Item{
		story("S-0001", workitem.Done, "flai/cmd"),                 // the accepted one
		story("S-0002", workitem.InProgress, "flai/cmd/accept.go"), // a file the merge changed
		story("S-0003", workitem.Review, "cli"),                    // the component, by tag
		story("S-0004", workitem.InProgress, "flaiover"),           // elsewhere
		story("S-0005", workitem.InProgress, "flaiover"),           // its task claims design/system
		story("S-0006", workitem.InProgress),                       // no touches: may change anything
		story("S-0007", workitem.Ready, "flai/cmd"),                // not open
		task,
	}
	projects := []manifest.Project{{Name: "flai", Path: "flai", Tags: []string{"cli"}}, {Name: "flaiover", Path: "flaiover"}}
	changed := []string{"flai/cmd/accept.go", "flai/cmd/move.go", "design/system/workflow.md"}
	got := map[string][]string{}
	for _, o := range overlapsOf(items, projects, manifest.Claims{}, "S-0001", changed) {
		if o.Accepted != "S-0001" || o.Title != o.ID {
			t.Errorf("notice: %+v", o)
		}
		got[o.ID] = o.Paths
	}
	want := map[string][]string{
		"S-0002": {"flai/cmd/accept.go"},
		"S-0003": {"flai/cmd/accept.go", "flai/cmd/move.go"},
		"S-0005": {"design/system/workflow.md"},
		"S-0006": changed,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("told:\n got %v\nwant %v", got, want)
	}
	if overlapsOf(items, projects, manifest.Claims{}, "S-0001", nil) != nil {
		t.Error("a merge that changed nothing tells nobody")
	}
	// S-0550: a changed path inside claims.shared is nobody's notice
	shared := manifest.Claims{Shared: []string{"design/system"}}
	got = map[string][]string{}
	for _, o := range overlapsOf(items, projects, shared, "S-0001", changed) {
		got[o.ID] = o.Paths
	}
	delete(want, "S-0005")
	want["S-0006"] = changed[:2]
	if !reflect.DeepEqual(got, want) {
		t.Errorf("told beside a shared path:\n got %v\nwant %v", got, want)
	}
	if o := overlapsOf(items, projects, shared, "S-0001", []string{"design/system/workflow.md"}); o != nil {
		t.Errorf("a merge wholly inside the shared paths told %+v", o)
	}
}

// Accepting a story records a notice for an open story whose claim covers
// what it changed, and none for one whose claim does not (S-0132).
func TestAcceptTellsOverlappingOpenStories(t *testing.T) {
	told, notices := acceptBesideOpenStories(t, "", [3]string{"docs/guide.md", "docs", "design"}, "docs/guide.md")
	if len(told) != 1 || told[0].ID != "S-0002" || !reflect.DeepEqual(told[0].Paths, []string{"docs/guide.md"}) {
		t.Errorf("told: %+v", told)
	}
	if len(notices) != 1 {
		t.Fatalf("notices recorded: %+v", notices)
	}
	if n := notices[0]; n.ID != "S-0002" || n.Accepted != "S-0001" || n.By != "alex" || n.At == "" || !reflect.DeepEqual(n.Paths, []string{"docs/guide.md"}) {
		t.Errorf("notice: %+v", n)
	}
}

// S-0550: accepting a story tells no open story of a changed path inside
// claims.shared. A change beside it is told as before; a change wholly
// inside it tells nobody, though both open stories' claims cover it.
func TestAcceptTellsNoOpenStoryOfASharedPath(t *testing.T) {
	const shared = "claims:\n  shared:\n    - docs/users/flai.md\n"
	touches := [3]string{"docs", "docs", "docs/users/flai.md"}
	t.Run("a shared path and one beside it", func(t *testing.T) {
		told, notices := acceptBesideOpenStories(t, shared, touches, "docs/users/flai.md", "docs/guide.md")
		if len(told) != 1 || told[0].ID != "S-0002" || !reflect.DeepEqual(told[0].Paths, []string{"docs/guide.md"}) {
			t.Errorf("told: %+v", told)
		}
		if len(notices) != 1 || notices[0].ID != "S-0002" || !reflect.DeepEqual(notices[0].Paths, []string{"docs/guide.md"}) {
			t.Errorf("notices recorded: %+v", notices)
		}
	})
	t.Run("only a shared path", func(t *testing.T) {
		told, notices := acceptBesideOpenStories(t, shared, touches, "docs/users/flai.md")
		if len(told) != 0 || len(notices) != 0 {
			t.Errorf("told %+v, recorded %+v", told, notices)
		}
	})
}

// acceptBesideOpenStories makes a project whose system-flow.yaml ends with
// manifestTail, and stories S-0001, S-0002, and S-0003 touching touches, in
// progress. S-0001 commits files, moves to review, and is accepted by alex;
// it returns what the acceptance says it told and the notices recorded.
func acceptBesideOpenStories(t *testing.T, manifestTail string, touches [3]string, files ...string) (told, notices []itemedit.Overlap) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	inEnv(t, "FLAI_AGENT", "tester")
	root := tempProject(t)
	if manifestTail != "" {
		m, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
		_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), append(m, manifestTail...), 0o644)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design", "docs"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{
		{"epic", "new", "Epic"},
		{"story", "new", "Accepted", "--epic", "E-0001", "--touches", touches[0]},
		{"story", "new", "Overlaps", "--epic", "E-0001", "--touches", touches[1]},
		{"story", "new", "Elsewhere", "--epic", "E-0001", "--touches", touches[2]},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	for _, f := range []string{"S-0001-accepted.md", "S-0002-overlaps.md", "S-0003-elsewhere.md"} {
		p := filepath.Join(root, "wip/kanban/stories", f)
		s, _ := os.ReadFile(p)
		_ = os.WriteFile(p, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	}
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"move", "S-0002", "ready"}, {"move", "S-0002", "in-progress"},
		{"move", "S-0003", "ready"}, {"move", "S-0003", "in-progress"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	for _, f := range files {
		_ = os.MkdirAll(filepath.Dir(filepath.Join(wt, f)), 0o755)
		_ = os.WriteFile(filepath.Join(wt, f), []byte(f+"\n"), 0o644)
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] change")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("accept: %d %s", code, errOut)
	}
	var res struct {
		Overlaps []itemedit.Overlap `json:"overlaps"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("accept json: %v\n%s", err, out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return res.Overlaps, itemedit.Overlaps(repo)
}
