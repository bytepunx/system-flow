package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/gittest"
)

// The token file and stylesheet a designProject's manifest names, under its
// web sub-project, as the Design page's token save finds them.
const (
	designTokens = "web/tokens.json"
	designSheet  = "web/src/theme.css"
)

// designProject is a project in git whose manifest names designTokens and
// designSheet, holding the starter token file and its stylesheet, committed.
func designProject(t *testing.T) channel.Project {
	t.Helper()
	gittest.Identity(t, "Olive Owner", "olive@example.com")
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "config.json"))
	root := t.TempDir()
	sheet, err := docedit.StylesheetFor(designTokens, designsystem.StarterTokens)
	if err != nil {
		t.Fatal(err)
	}
	for rel, content := range map[string]string{
		"system-flow.yaml": "version: 1\nname: Harbour\nkey: harbour\nowner: olive\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
			"projects:\n  - name: web\n    path: web\n    kind: sveltekit\ndesign:\n  tokens: " + designTokens + "\n  stylesheet: " + designSheet + "\n",
		".gitignore":                  ".flai-cache/\n",
		"wip/kanban/epics/.gitkeep":   "",
		"wip/kanban/stories/.gitkeep": "",
		"wip/kanban/tasks/.gitkeep":   "",
		"wip/agents/.gitkeep":         "",
		"wip/archive/.gitkeep":        "",
		"design/system/palette.json":  "{\"accent\": \"main\"}\n",
		designTokens:                  string(designsystem.StarterTokens),
		designSheet:                   sheet,
		"docs/users/.gitkeep":         "",
	} {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "design system")
	return channel.Project{Key: "harbour", Name: "Harbour", Root: root}
}

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// S-0436, end to end with a real flai built from this tree: the Design
// page's token save at /api/design/tokens is file.get then file.save of the
// token file, and on the host it saves the tokens, writes the stylesheet
// from them, and commits both in one commit with the dashboard's trailer
// (ADR-0147, ADR-0150). A stale hash is then a Conflict with the current
// content, and a token file flai cannot render is Refused with its
// findings, leaving the commit and the files as they were.
func TestTheDesignPagesTokenSaveCommitsTheTokensAndTheStylesheet(t *testing.T) {
	if testing.Short() {
		t.Skip("builds flai and runs it on a project in git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	exe := filepath.Join(t.TempDir(), "flai")
	if out, err := exec.Command("go", "build", "-o", exe, "github.com/bytepunx/system-flow/flai").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	p := designProject(t)
	methods := writeMethods(func(ctx context.Context, r Run) (Ran, error) { return runFlai(ctx, exe, r) }, time.Now, Host{})
	call := func(method, params string, into any) *channel.Error {
		t.Helper()
		res, e := methods[method](context.Background(), p, json.RawMessage(params))
		if e != nil {
			return e
		}
		if err := json.Unmarshal(res.(Written).Data, into); err != nil {
			t.Fatalf("%s answered %s: %v", method, res.(Written).Data, err)
		}
		return nil
	}
	params := func(v map[string]any) string {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	var loaded docedit.File
	if e := call("file.get", `{"path":"`+designTokens+`"}`, &loaded); e != nil {
		t.Fatalf("file.get: %+v", e)
	}
	if loaded.Path != designTokens || loaded.Content != string(designsystem.StarterTokens) || loaded.Hash != docedit.Hash(loaded.Content) {
		t.Fatalf("file.get answered %+v", loaded)
	}

	edited := strings.Replace(loaded.Content, `"components": [0.145, 0.388, 0.922], "hex": "#2563eb"`, `"components": [0.114, 0.306, 0.847], "hex": "#1d4ed8"`, 1)
	if edited == loaded.Content {
		t.Fatal("the starter token file no longer has the primary colour this test edits")
	}
	head := strings.TrimSpace(gitIn(t, p.Root, "rev-parse", "HEAD"))
	var saved docedit.FileResult
	if e := call("file.save", params(map[string]any{"path": designTokens, "content": edited, "hash": loaded.Hash, "message": "darker primary", "request_id": "req-00000001"}), &saved); e != nil {
		t.Fatalf("file.save: %+v", e)
	}
	if !saved.Committed || saved.Stylesheet != designSheet || saved.Hash != docedit.Hash(edited) {
		t.Errorf("file.save answered %+v", saved)
	}
	if parent := strings.TrimSpace(gitIn(t, p.Root, "rev-parse", "HEAD~1")); parent != head {
		t.Errorf("the save is one commit on %s, HEAD~1 is %s", head, parent)
	}
	if files := strings.Fields(gitIn(t, p.Root, "show", "--name-only", "--format=", "HEAD")); strings.Join(files, " ") != designSheet+" "+designTokens {
		t.Errorf("the commit holds %v, want the token file and the stylesheet", files)
	}
	if msg := gitIn(t, p.Root, "log", "-1", "--format=%B"); !strings.HasPrefix(msg, "docs: darker primary") || !strings.Contains(msg, Trailer) {
		t.Errorf("commit message %q", msg)
	}
	want, err := docedit.StylesheetFor(designTokens, []byte(edited))
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(p.Root, designSheet)); string(got) != want {
		t.Error("the stylesheet is written from the saved tokens")
	}
	if status := gitIn(t, p.Root, "status", "--porcelain"); status != "" {
		t.Errorf("the save left changes uncommitted:\n%s", status)
	}

	head = strings.TrimSpace(gitIn(t, p.Root, "rev-parse", "HEAD"))
	e := call("file.save", params(map[string]any{"path": designTokens, "content": loaded.Content, "hash": loaded.Hash, "request_id": "req-00000002"}), &saved)
	if e == nil || e.Code != Conflict || !strings.Contains(e.Message, designTokens+" changed after it was loaded") {
		t.Fatalf("a stale hash: %+v", e)
	}
	if d := e.Data.(map[string]any); d["current"] != edited || d["hash"] != docedit.Hash(edited) {
		t.Errorf("a conflict carries the current content and its hash: %+v", d)
	}

	e = call("file.save", params(map[string]any{"path": designTokens, "content": "{\"color\": \n", "hash": docedit.Hash(edited), "request_id": "req-00000003"}), &saved)
	if e == nil || e.Code != Refused {
		t.Fatalf("a token file flai cannot render: %+v", e)
	}
	findings, _ := e.Data.(map[string]any)["findings"].([]any)
	if len(findings) == 0 || findings[0].(map[string]any)["rule"] != "design.tokens" {
		t.Errorf("a refusal carries its findings: %+v", e.Data)
	}
	if now := strings.TrimSpace(gitIn(t, p.Root, "rev-parse", "HEAD")); now != head {
		t.Error("a conflict and a refusal commit nothing")
	}
	if got, _ := os.ReadFile(filepath.Join(p.Root, designTokens)); string(got) != edited {
		t.Error("a refusal leaves the token file as it was")
	}
}
