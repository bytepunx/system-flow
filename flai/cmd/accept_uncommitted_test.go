package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// S-0051, ADR-0248: the dry run lists the main checkout's uncommitted files
// in the shared set, which the acceptance commits, as shared_uncommitted,
// and the rest, which it leaves as they are, as uncommitted; the real run
// asks nothing and refuses nothing for them, and --yes includes the rest.
func TestAcceptReportsUncommittedPathsInTheDryRun(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "docs", "design"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	// two stories in review, with no branch: one accepted without --yes, one
	// moved to done with it
	for i, title := range []string{"Stray", "Second"} {
		story, task := "S-000"+string(rune('1'+i)), "T-000"+string(rune('1'+i))
		if _, errOut, code := runIn(t, root, "story", "new", title, "--epic", "E-0001"); code != 0 {
			t.Fatal(errOut)
		}
		storyFile := filepath.Join(root, "wip/kanban/stories/"+story+"-"+strings.ToLower(title)+".md")
		s, _ := os.ReadFile(storyFile)
		_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
		for _, args := range [][]string{
			{"move", story, "ready"}, {"move", story, "in-progress"}, {"stream", "open", story, "--no-branch"},
			{"task", "new", "Do it", "--story", story},
			{"move", task, "ready"}, {"move", task, "in-progress"}, {"move", task, "done"},
			{"move", story, "review"},
		} {
			if _, errOut, code := runIn(t, root, args...); code != 0 {
				t.Fatalf("flai %v: %s", args, errOut)
			}
		}
	}
	type lists struct {
		Shared      []string `json:"shared_uncommitted"`
		Uncommitted []string `json:"uncommitted"`
		Includes    bool     `json:"includes_uncommitted"`
	}
	preview := func(args ...string) lists {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"accept", "S-0001", "--dry-run", "--json"}, args...)...)
		if code != 0 {
			t.Fatalf("a dry run must not refuse (%v): %d %s", args, code, errOut)
		}
		var res lists
		if err := json.Unmarshal([]byte(out), &res); err != nil {
			t.Fatalf("preview json: %v\n%s", err, out)
		}
		return res
	}
	if got := preview(); len(got.Uncommitted) != 0 || len(got.Shared) == 0 || slices.ContainsFunc(got.Shared, func(p string) bool { return !strings.HasPrefix(p, "wip/") }) {
		t.Errorf("clean tree, wip changes only: uncommitted = %v, shared_uncommitted = %v", got.Uncommitted, got.Shared)
	}

	_ = os.WriteFile(filepath.Join(root, "docs", "stray.md"), []byte("x\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "design", "d.md"), []byte("# flai's\n"), 0o644)
	for _, args := range [][]string{nil, {"--yes"}} {
		got := preview(args...)
		if !slices.Equal(got.Uncommitted, []string{"docs/stray.md"}) || !slices.Contains(got.Shared, "design/d.md") || got.Includes != (len(args) > 0) {
			t.Errorf("dry run %v should report the stray file as uncommitted and the design file as shared, got %+v", args, got)
		}
	}
	out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run")
	for _, want := range []string{
		"uncommitted under design and the shared paths: design/d.md; the acceptance commits them with the acceptance commit, as it commits wip (ADR-0248)",
		"uncommitted outside wip, design, and the shared paths: docs/stray.md; the acceptance leaves them as they are", "--yes includes them",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("text preview should say %q:\n%s", want, out)
		}
	}

	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("the real run asks nothing about the uncommitted files: %d %s", code, errOut)
	}
	files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(files, "design/d.md") || strings.Contains(files, "docs/stray.md") {
		t.Errorf("the acceptance commit holds the design file and not the stray one:\n%s", files)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "?? docs/stray.md" {
		t.Errorf("the stray file alone is left uncommitted: %q", st)
	}

	if _, errOut, code := runIn(t, root, "move", "S-0002", "done", "--yes"); code != 0 {
		t.Fatalf("move to done with --yes should accept and include the file: %s", errOut)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); !strings.Contains(files, "docs/stray.md") {
		t.Errorf("the acceptance commit under --yes should hold the stray file:\n%s", files)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("tree should be clean after acceptance with --yes: %s", st)
	}
}

// statusRunner is the real runner that records, in order, each git status
// --porcelain run in the checkout at root, with the git commands run there
// before it.
type statusRunner struct {
	execx.System
	root  string
	mu    sync.Mutex
	reads [][]string
	since []string
}

func (s *statusRunner) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && samePath(dir, s.root) {
		s.mu.Lock()
		if slices.Contains(args, "status") && slices.Contains(args, "--porcelain") {
			s.reads = append(s.reads, s.since)
			s.since = nil
		} else {
			s.since = append(s.since, strings.Join(args, " "))
		}
		s.mu.Unlock()
	}
	return s.System.Run(dir, name, args...)
}

// samePath reports whether a and b name one folder.
func samePath(a, b string) bool {
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	return errA == nil && errB == nil && ra == rb
}

// S-0609 criterion 2: an acceptance reads the main checkout's git status
// once before it changes anything and once more only before its commit,
// after the merge and the archive changed the checkout; the file outside
// the shared set is left as it is (ADR-0248).
func TestAcceptReadsTheMainCheckoutsStatusOncePerState(t *testing.T) {
	t.Setenv("FLAI_AGENT", "tester")
	log := filepath.Join(t.TempDir(), "lint.log")
	root := lintProject(t, countingTier(log))
	wt := openSyncStory(t, root, 1, "docs")
	commitIn(t, wt, "docs/new.md", "fixme\n")
	issueStoryInReview(t, root, "S-0001", "T-0001")
	for rel, content := range map[string]string{"wip/notes/n.md": "fixme\n", "notes/mine.txt": "mine\n"} {
		_ = os.MkdirAll(filepath.Join(root, filepath.Dir(rel)), 0o755)
		if err := os.WriteFile(filepath.Join(root, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r := &statusRunner{root: root}
	if out, errOut, code := runWithApp(t, &app{configPath: configOf(t), cwd: root, runner: r}, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %d %s %s", code, out, errOut)
	}
	// before S-0609 the preflight lint read it again between the two, with
	// nothing changed since the preview's read
	if len(r.reads) != 2 {
		t.Fatalf("git status --porcelain of the main checkout ran %d times, want 2: %q", len(r.reads), r.reads)
	}
	if !slices.ContainsFunc(r.reads[1], func(c string) bool { return strings.HasPrefix(c, "merge --ff-only") }) {
		t.Errorf("the second read follows the merge, which changed the checkout; it followed %q", r.reads[1])
	}
	if got := runsNaming(t, log, "fix", "wip/notes/n.md"); got != 1 {
		t.Errorf("wip/notes/n.md is fixed once, from the preview's read: %d", got)
	}
	if got := gitIn(t, root, "show", "HEAD:wip/notes/n.md"); got != "fixed" {
		t.Errorf("the acceptance committed wip/notes/n.md as %q", got)
	}
	if got := gitIn(t, root, "status", "--porcelain"); got != "?? notes/" {
		t.Errorf("the file outside the shared set is left as it is: %q", got)
	}
}
