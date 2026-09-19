package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// runStdin is runIn with content on standard input and the exit code kept.
func runStdin(t *testing.T, dir, stdin string, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: dir, clock: func() time.Time { return time.Date(2026, 9, 19, 2, 0, 0, 0, time.UTC) }}
	root := newRootCmdWith(a)
	root.SetArgs(args)
	root.SetIn(strings.NewReader(stdin))
	code := 0
	if err := root.Execute(); err != nil {
		var ee *exitError
		if errors.As(err, &ee) {
			code = ee.code
			if ee.msg != "" {
				a.fail(err)
			}
		} else {
			a.fail(err)
			code = 1
		}
	}
	return out.String(), errOut.String(), code
}

const designDoc = "---\ntitle: Plan\nupdated: 2026-09-01\nstatus: active\n---\n\n# Plan\n\n## Shape\ntext\n"

func docProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GIT_AUTHOR_NAME", "Dana Designer")
	t.Setenv("GIT_AUTHOR_EMAIL", "dana@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "Dana Designer")
	t.Setenv("GIT_COMMITTER_EMAIL", "dana@example.com")
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design/system", "docs"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte(designDoc), 0o644)
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Editable", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	return root
}

func showDoc(t *testing.T, root, path string) (content, hash, mode string) {
	t.Helper()
	out, errOut, code := runIn(t, root, "doc", "show", path, "--json")
	if code != 0 {
		t.Fatalf("doc show %s: %s", path, errOut)
	}
	var d struct{ Content, Hash, Mode, Reason string }
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	return d.Content, d.Hash, d.Mode
}

func TestDocSaveCommitsOnePathAsTheDesigner(t *testing.T) {
	root := docProject(t)
	// an agent's uncommitted wip change must not be swept into the commit
	story := filepath.Join(root, "wip/kanban/stories/S-0001-editable.md")
	s, _ := os.ReadFile(story)
	_ = os.WriteFile(story, append(s, []byte("agent note\n")...), 0o644)

	content, hash, mode := showDoc(t, root, "design/system/plan.md")
	if mode != "full" {
		t.Fatalf("design docs are fully editable, got %s", mode)
	}
	edited := strings.Replace(content, "text\n", "better text\n", 1)
	out, errOut, code := runStdin(t, root, edited, "doc", "save", "design/system/plan.md", "--hash", hash, "--message", "clarify the shape", "--trailer", "Co-Authored-By: flaiover <flaiover@localhost>", "--json")
	if code != 0 {
		t.Fatalf("save: %d %s", code, errOut)
	}
	var res struct {
		Hash      string
		Committed bool
		Commit    string
	}
	_ = json.Unmarshal([]byte(out), &res)
	if !res.Committed || res.Commit == "" || res.Hash == hash {
		t.Errorf("result: %s", out)
	}
	saved, _ := os.ReadFile(filepath.Join(root, "design/system/plan.md"))
	if !strings.Contains(string(saved), "better text") || !strings.Contains(string(saved), "updated: 2026-09-19") {
		t.Errorf("saved content should carry the edit and today's date:\n%s", saved)
	}
	log := gitIn(t, root, "log", "-1", "--format=%an <%ae>%n%B")
	for _, want := range []string{"Dana Designer <dana@example.com>", "docs: clarify the shape", "Co-Authored-By: flaiover <flaiover@localhost>"} {
		if !strings.Contains(log, want) {
			t.Errorf("commit should have %q:\n%s", want, log)
		}
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); files != "design/system/plan.md" {
		t.Errorf("the commit must hold that path only, got %q", files)
	}
	if st := gitIn(t, root, "status", "--porcelain"); !strings.Contains(st, "S-0001-editable.md") {
		t.Errorf("the agent's wip change should still be uncommitted: %q", st)
	}

	// the same content again is not a new commit
	_, hash2, _ := showDoc(t, root, "design/system/plan.md")
	out, _, code = runStdin(t, root, string(saved), "doc", "save", "design/system/plan.md", "--hash", hash2, "--json")
	if code != 0 || !strings.Contains(out, `"unchanged": true`) {
		t.Errorf("unchanged save: %d %s", code, out)
	}
}

func TestDocSaveKeepsTheDesignersOwnDate(t *testing.T) {
	root := docProject(t)
	content, hash, _ := showDoc(t, root, "design/system/plan.md")
	edited := strings.Replace(content, "updated: 2026-09-01", "updated: 2026-09-10", 1)
	if _, errOut, code := runStdin(t, root, edited, "doc", "save", "design/system/plan.md", "--hash", hash); code != 0 {
		t.Fatal(errOut)
	}
	saved, _ := os.ReadFile(filepath.Join(root, "design/system/plan.md"))
	if !strings.Contains(string(saved), "updated: 2026-09-10") {
		t.Errorf("a date the designer set is kept:\n%s", saved)
	}
}

func TestDocSaveWithoutCommit(t *testing.T) {
	root := docProject(t)
	head := gitIn(t, root, "rev-parse", "HEAD")
	content, hash, _ := showDoc(t, root, "design/system/plan.md")
	if out, errOut, code := runStdin(t, root, content+"more\n", "doc", "save", "design/system/plan.md", "--hash", hash, "--no-commit", "--json"); code != 0 || !strings.Contains(out, `"committed": false`) {
		t.Fatalf("--no-commit: %d %s %s", code, out, errOut)
	}
	// and by project setting
	m := filepath.Join(root, "system-flow.yaml")
	data, _ := os.ReadFile(m)
	_ = os.WriteFile(m, append(data, []byte("dashboard:\n  autocommit: false\n")...), 0o644)
	content, hash, _ = showDoc(t, root, "design/system/plan.md")
	if out, errOut, code := runStdin(t, root, content+"and more\n", "doc", "save", "design/system/plan.md", "--hash", hash, "--json"); code != 0 || !strings.Contains(out, `"committed": false`) {
		t.Fatalf("autocommit false: %d %s %s", code, out, errOut)
	}
	if gitIn(t, root, "rev-parse", "HEAD") != head {
		t.Error("nothing should have been committed")
	}
}

func TestDocSaveBodyModeProtectsFrontMatter(t *testing.T) {
	root := docProject(t)
	path := "wip/kanban/stories/S-0001-editable.md"
	content, hash, mode := showDoc(t, root, path)
	if mode != "body" {
		t.Fatalf("work items are body mode, got %s", mode)
	}
	out, errOut, code := runStdin(t, root, strings.Replace(content, "status: backlog", "status: done", 1), "doc", "save", path, "--hash", hash, "--json")
	if code != exitDocRefused || !strings.Contains(errOut, "front matter") || !strings.Contains(out, `"refused"`) {
		t.Errorf("front matter change must be refused: %d %s %s", code, out, errOut)
	}
	if now, _ := os.ReadFile(filepath.Join(root, path)); string(now) != content {
		t.Error("a refused save must leave the file as it was")
	}
	edited := strings.Replace(content, "## Goal\n", "## Goal\nEdit bodies in the dashboard.\n", 1)
	if _, errOut, code := runStdin(t, root, edited, "doc", "save", path, "--hash", hash); code != 0 {
		t.Fatalf("body edit: %s", errOut)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "chore: [S-0001] edit "+path {
		t.Errorf("subject: %q", subject)
	}
}

func TestDocSaveRefusesWhatCannotBeEdited(t *testing.T) {
	root := docProject(t)
	if _, errOut, code := runIn(t, root, "stream", "open", "S-0001", "--no-branch"); code != 0 {
		t.Fatal(errOut)
	}
	_, hash, mode := showDoc(t, root, "wip/agents/index.md")
	if mode != "none" {
		t.Fatalf("index.md is generated, got %s", mode)
	}
	if _, errOut, code := runStdin(t, root, "x\n", "doc", "save", "wip/agents/index.md", "--hash", hash); code != exitDocRefused || !strings.Contains(errOut, "generated") {
		t.Errorf("generated file: %d %s", code, errOut)
	}
	for _, p := range []string{"../outside.md", "system-flow.yaml", "README.md", "design/system/missing.md"} {
		if _, _, code := runIn(t, root, "doc", "show", p); code == 0 {
			t.Errorf("%s must not be served", p)
		}
	}
}

func TestDocSaveConflict(t *testing.T) {
	root := docProject(t)
	content, hash, _ := showDoc(t, root, "design/system/plan.md")
	// someone else saves first
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte(strings.Replace(content, "text\n", "their text\n", 1)), 0o644)
	out, errOut, code := runStdin(t, root, strings.Replace(content, "text\n", "my text\n", 1), "doc", "save", "design/system/plan.md", "--hash", hash, "--json")
	if code != exitDocConflict || !strings.Contains(errOut, "conflict:") {
		t.Fatalf("conflict: %d %s", code, errOut)
	}
	var res struct {
		Conflict struct{ Current, Hash, Diff string }
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	c := res.Conflict
	if !strings.Contains(c.Current, "their text") || c.Hash == hash || !strings.Contains(c.Diff, "-their text") || !strings.Contains(c.Diff, "+my text") {
		t.Errorf("conflict payload: %+v", c)
	}
	if !strings.Contains(c.Diff, "--- design/system/plan.md (current)") || !strings.Contains(c.Diff, "+++ design/system/plan.md (yours)") || strings.Contains(c.Diff, root) || strings.Contains(c.Diff, os.TempDir()) {
		t.Errorf("the diff names the two sides and keeps host paths out:\n%s", c.Diff)
	}
	if now, _ := os.ReadFile(filepath.Join(root, "design/system/plan.md")); !strings.Contains(string(now), "their text") {
		t.Error("a conflict must not overwrite the other change")
	}
	// saving over it is the explicit second save with the current hash
	if _, errOut, code := runStdin(t, root, strings.Replace(content, "text\n", "my text\n", 1), "doc", "save", "design/system/plan.md", "--hash", c.Hash); code != 0 {
		t.Errorf("save over: %s", errOut)
	}
}

func TestDocSaveRefusedByTheCheckRestoresTheFile(t *testing.T) {
	root := docProject(t)
	content, hash, _ := showDoc(t, root, "design/system/plan.md")
	broken := strings.Replace(content, "title: Plan\n", "", 1)
	out, errOut, code := runStdin(t, root, broken, "doc", "save", "design/system/plan.md", "--hash", hash, "--json")
	if code != exitDocRefused || !strings.Contains(errOut, "refused:") {
		t.Fatalf("a document without a title fails the check: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "doc.title") {
		t.Errorf("the findings should name the rule:\n%s", out)
	}
	if now, _ := os.ReadFile(filepath.Join(root, "design/system/plan.md")); string(now) != content {
		t.Error("a refused save must restore the file")
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("nothing should be left behind: %q", st)
	}
}
