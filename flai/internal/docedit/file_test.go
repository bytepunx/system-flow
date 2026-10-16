package docedit

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const (
	tokensPath  = "web/tokens.json"
	sheetPath   = "web/src/theme.css"
	palettePath = "design/system/palette.json"
	mainPalette = "{\"accent\": \"main\"}\n"
	wtPalette   = "{\"accent\": \"branch\"}\n"
	// fileManifest names a token file and a stylesheet outside the design
	// and docs folders, under the web sub-project.
	fileManifest = "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
		"projects:\n  - name: web\n    path: web\n    kind: sveltekit\n" +
		"design:\n  tokens: " + tokensPath + "\n  stylesheet: " + sheetPath + "\n"
)

// starter is the starter token file, and edited is it with another primary
// colour, so that its stylesheet differs.
var (
	starter = string(designsystem.StarterTokens)
	edited  = strings.Replace(starter, `"components": [0.145, 0.388, 0.922], "hex": "#2563eb"`, `"components": [0.114, 0.306, 0.847], "hex": "#1d4ed8"`, 1)
)

func put(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func sheetFor(t *testing.T, tokens string) string {
	t.Helper()
	css, err := StylesheetFor(tokensPath, []byte(tokens))
	if err != nil {
		t.Fatal(err)
	}
	return css
}

// fileProject is project with fileManifest, and extra after it, in the main
// checkout and S-0001's worktree, each holding the starter token file, its
// stylesheet, and a palette; the main checkout also has files out of reach.
func fileProject(t *testing.T, extra string) (repo *workitem.Repo, root, wt string) {
	t.Helper()
	_, root, wt = project(t)
	for dir, palette := range map[string]string{root: mainPalette, wt: wtPalette} {
		put(t, filepath.Join(dir, "system-flow.yaml"), fileManifest+extra)
		put(t, filepath.Join(dir, tokensPath), starter)
		put(t, filepath.Join(dir, sheetPath), sheetFor(t, starter))
		put(t, filepath.Join(dir, palettePath), palette)
	}
	put(t, filepath.Join(root, "docs/notes.txt"), "notes\n")
	put(t, filepath.Join(root, "src/main.go"), "package main\n")
	put(t, filepath.Join(root, "wip/notes.txt"), "notes\n")
	put(t, filepath.Join(root, "design/system/blob.bin"), "GIF89a\x00\x01")
	put(t, filepath.Join(root, "design/system/big.txt"), strings.Repeat("x", MaxFileSize+1))
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo, root, wt
}

func committed(git *fakeGit) bool {
	for _, c := range git.calls {
		if strings.Contains(c, "git add") || strings.Contains(c, "git commit") {
			return true
		}
	}
	return false
}

// S-0436: file show and file save reach a file under the design or docs
// folder and the token file and stylesheet the manifest names, and refuse
// every other path with why.
func TestFileReach(t *testing.T) {
	repo, _, _ := fileProject(t, "")
	cases := []struct {
		name, path, clean, want string
	}{
		{"the design folder", "./" + palettePath, palettePath, ""},
		{"the docs folder", "docs/notes.txt", "docs/notes.txt", ""},
		{"the token file", tokensPath, tokensPath, ""},
		{"the stylesheet", sheetPath, sheetPath, ""},
		{"the wip folder", "wip/notes.txt", "", "wip/notes.txt is under the wip folder wip/"},
		{"a markdown file", planPath, "", "design/system/plan.md is a markdown file; read and save it with flai doc show and flai doc save"},
		{"another path", "src/main.go", "", "src/main.go is outside the design and docs folders and is not the design token file web/tokens.json or the design stylesheet web/src/theme.css"},
		{"a path out of the repository", "../x.json", "", "is not a path inside the repository"},
		{"an absolute path", "/etc/hostname", "", "is not a path inside the repository"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clean, err := FileReach(repo.Manifest, c.path)
			if c.want == "" {
				if err != nil || clean != c.clean {
					t.Fatalf("want %s reached, got %q %v", c.clean, clean, err)
				}
				f, err := FileShow(repo, c.path, "")
				if err != nil || f.Path != c.clean || f.Hash != Hash(f.Content) {
					t.Errorf("show: %+v %v", f, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("want %q, got %q %v", c.want, clean, err)
			}
			if _, err := FileShow(repo, c.path, ""); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("show: want %q, got %v", c.want, err)
			}
			git := &fakeGit{branch: "main"}
			if _, err := FileSave(repo, git, c.path, "x\n", SaveOptions{Hash: Hash("notes\n"), Now: now}); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("save: want %q, got %v", c.want, err)
			}
			if committed(git) {
				t.Errorf("nothing is committed: %v", git.calls)
			}
		})
	}
}

// S-0436: a binary file and one over the size limit are refused, as is
// content of either kind, and nothing is written.
func TestFileRefusesBinaryAndOversize(t *testing.T) {
	repo, root, _ := fileProject(t, "")
	for path, want := range map[string]string{
		"design/system/blob.bin": "design/system/blob.bin is binary, with a NUL byte in its first 8192 bytes",
		"design/system/big.txt":  "design/system/big.txt is 1048577 bytes, over the 1048576-byte limit",
	} {
		if _, err := FileShow(repo, path, ""); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("show %s: want %q, got %v", path, want, err)
		}
		if _, err := FileSave(repo, &fakeGit{}, path, "text\n", SaveOptions{Hash: "x", Now: now}); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("save %s: want %q, got %v", path, want, err)
		}
	}
	for name, content := range map[string]string{
		"binary content":   "a\x00b",
		"oversize content": strings.Repeat("x", MaxFileSize+1),
	} {
		_, err := FileSave(repo, &fakeGit{}, palettePath, content, SaveOptions{Hash: Hash(mainPalette), Now: now})
		if _, ok := IsRefused(err); !ok {
			t.Errorf("%s is refused: %v", name, err)
		}
	}
	if got := read(t, filepath.Join(root, palettePath)); got != mainPalette {
		t.Errorf("nothing is written:\n%s", got)
	}
}

// S-0436: saving the token file writes the stylesheet from the new tokens,
// and commits both paths in one commit on the main checkout's branch.
func TestFileSaveTokensWritesTheStylesheet(t *testing.T) {
	repo, root, _ := fileProject(t, "")
	git := &fakeGit{branch: "main"}
	res, err := FileSave(repo, git, tokensPath, edited, SaveOptions{Hash: Hash(starter), Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Committed || res.Commit != "abc1234" || res.Path != tokensPath || res.Stylesheet != sheetPath || res.Hash != Hash(edited) {
		t.Errorf("result: %+v", res)
	}
	if got := read(t, filepath.Join(root, tokensPath)); got != edited {
		t.Errorf("the token file is saved as it is:\n%s", got)
	}
	want := sheetFor(t, edited)
	if want == sheetFor(t, starter) {
		t.Fatal("the edit must change the stylesheet")
	}
	if got := read(t, filepath.Join(root, sheetPath)); got != want {
		t.Errorf("the stylesheet is written from the new tokens:\n%s", got)
	}
	joined := strings.Join(git.calls, "\n")
	for _, w := range []string{
		root + ": git add -- " + tokensPath + " " + sheetPath,
		root + ": git commit --quiet -m docs: edit " + tokensPath + " and write " + sheetPath + " from it -- " + tokensPath + " " + sheetPath,
	} {
		if !strings.Contains(joined, w) {
			t.Errorf("git should have run %q; ran:\n%s", w, joined)
		}
	}
}

// S-0436: a token save writes the stylesheet the manifest names when it is
// not there yet, and a save of the stylesheet itself regenerates nothing.
func TestFileSaveStylesheet(t *testing.T) {
	repo, root, _ := fileProject(t, "")
	if err := os.RemoveAll(filepath.Join(root, "web/src")); err != nil {
		t.Fatal(err)
	}
	if _, err := FileSave(repo, &fakeGit{}, tokensPath, edited, SaveOptions{Hash: Hash(starter), Now: now}); err != nil {
		t.Fatal(err)
	}
	if got := read(t, filepath.Join(root, sheetPath)); got != sheetFor(t, edited) {
		t.Errorf("the missing stylesheet is written:\n%s", got)
	}

	// A stylesheet already behind its tokens may be saved by hand: design.stale
	// is no new finding, and the tokens are left alone.
	put(t, filepath.Join(root, sheetPath), "/* stale */\n")
	git := &fakeGit{branch: "main"}
	res, err := FileSave(repo, git, sheetPath, "/* by hand */\n", SaveOptions{Hash: Hash("/* stale */\n"), Message: "hand", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if res.Stylesheet != "" || read(t, filepath.Join(root, sheetPath)) != "/* by hand */\n" || read(t, filepath.Join(root, tokensPath)) != edited {
		t.Errorf("only the stylesheet is saved: %+v", res)
	}
	if len(res.Warnings) == 0 || res.Warnings[0].Rule != "design.stale" {
		t.Errorf("the stylesheet's design.stale warning is reported: %+v", res.Warnings)
	}
	if joined := strings.Join(git.calls, "\n"); !strings.Contains(joined, "git commit --quiet -m docs: hand -- "+sheetPath+"\n") {
		t.Errorf("one path is committed:\n%s", joined)
	}
}

// S-0436: a save whose file changed after it was loaded is a conflict with
// the current content and a diff, and writes nothing.
func TestFileSaveRefusesAStaleHash(t *testing.T) {
	repo, root, _ := fileProject(t, "")
	theirs := "{\"accent\": \"theirs\"}\n"
	put(t, filepath.Join(root, palettePath), theirs)
	git := &fakeGit{branch: "main"}
	_, err := FileSave(repo, git, palettePath, "{\"accent\": \"mine\"}\n", SaveOptions{Hash: Hash(mainPalette), Now: now})
	c, ok := IsConflict(err)
	if !ok {
		t.Fatalf("a stale hash must be a conflict: %v", err)
	}
	if c.Current != theirs || c.Hash != Hash(theirs) || c.Path != palettePath {
		t.Errorf("the conflict carries the current file: %+v", c)
	}
	if !strings.Contains(strings.Join(git.calls, "\n"), "git diff --no-index") {
		t.Errorf("the conflict is diffed: %v", git.calls)
	}
	if got := read(t, filepath.Join(root, palettePath)); got != theirs {
		t.Errorf("a conflict must not overwrite the other change:\n%s", got)
	}
	if committed(git) {
		t.Errorf("nothing is committed: %v", git.calls)
	}
}

// S-0436: a finding the save introduces puts back both the token file and
// the stylesheet, or removes a stylesheet that was not there, and commits
// nothing.
func TestFileSaveFindingRestoresBoth(t *testing.T) {
	introduced := check.Finding{Level: check.Warning, Rule: "test.introduced", Path: tokensPath, Line: 1, Message: "introduced by the save"}
	calls := 0
	runCheck = func(repo *workitem.Repo, now time.Time) (*check.Result, error) {
		res, err := check.Run(repo, now)
		if calls++; err == nil && calls%2 == 0 {
			res.Findings = append(res.Findings, introduced)
		}
		return res, err
	}
	t.Cleanup(func() { runCheck = check.Run })
	for _, missing := range []bool{false, true} {
		repo, root, _ := fileProject(t, "")
		if missing {
			if err := os.Remove(filepath.Join(root, sheetPath)); err != nil {
				t.Fatal(err)
			}
		}
		git := &fakeGit{branch: "main"}
		_, err := FileSave(repo, git, tokensPath, edited, SaveOptions{Hash: Hash(starter), Now: now})
		r, ok := IsRefused(err)
		if !ok || len(r.Findings) != 1 || r.Findings[0].Rule != introduced.Rule {
			t.Fatalf("missing %v: the finding refuses the save: %v", missing, err)
		}
		if got := read(t, filepath.Join(root, tokensPath)); got != starter {
			t.Errorf("missing %v: the token file is put back:\n%.200s", missing, got)
		}
		_, statErr := os.Stat(filepath.Join(root, sheetPath))
		switch {
		case missing && !errors.Is(statErr, fs.ErrNotExist):
			t.Errorf("the stylesheet that was not there is removed: %v", statErr)
		case !missing && read(t, filepath.Join(root, sheetPath)) != sheetFor(t, starter):
			t.Error("the stylesheet is put back")
		}
		if committed(git) {
			t.Errorf("nothing is committed: %v", git.calls)
		}
	}
}

// S-0436: token content flai cannot render is refused with its problems,
// and nothing is written.
func TestFileSaveRefusesTokensThatDoNotParse(t *testing.T) {
	repo, root, _ := fileProject(t, "")
	for name, content := range map[string]string{
		"not JSON":       "{\"color\": \n",
		"not all tokens": "{\"color\": {\"$type\": \"color\", \"primary\": {\"$value\": 3}}}\n",
	} {
		_, err := FileSave(repo, &fakeGit{}, tokensPath, content, SaveOptions{Hash: Hash(starter), Now: now})
		r, ok := IsRefused(err)
		if !ok || len(r.Findings) == 0 || r.Findings[0].Rule != "design.tokens" || r.Findings[0].Path != tokensPath {
			t.Errorf("%s is refused with its problems: %v", name, err)
		}
	}
	if read(t, filepath.Join(root, tokensPath)) != starter || read(t, filepath.Join(root, sheetPath)) != sheetFor(t, starter) {
		t.Error("nothing is written")
	}
}

// S-0436: with dashboard.autocommit false, a save writes and commits
// nothing.
func TestFileSaveWithAutocommitOff(t *testing.T) {
	repo, root, _ := fileProject(t, "dashboard:\n  autocommit: false\n")
	git := &fakeGit{branch: "main"}
	res, err := FileSave(repo, git, palettePath, "{\"accent\": \"new\"}\n", SaveOptions{Hash: Hash(mainPalette), Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if res.Committed || committed(git) {
		t.Errorf("nothing is committed: %+v %v", res, git.calls)
	}
	if got := read(t, filepath.Join(root, palettePath)); got != "{\"accent\": \"new\"}\n" {
		t.Errorf("the file is saved:\n%s", got)
	}
}

// S-0436, ADR-0145: with a story, a file is read, saved, checked, and
// committed with its stylesheet in the story's worktree, on its branch, and
// a worktree off its branch is refused.
func TestFileSaveWithAStory(t *testing.T) {
	repo, root, wt := fileProject(t, "")
	f, err := FileShow(repo, palettePath, "S-1")
	if err != nil || f.Content != wtPalette {
		t.Fatalf("show reads the worktree: %+v %v", f, err)
	}
	if _, err := FileSave(repo, &fakeGit{branch: "main"}, tokensPath, edited, SaveOptions{Hash: Hash(starter), Story: "S-0001", Now: now}); err == nil || !strings.Contains(err.Error(), "its worktree is on main, not story/S-0001") {
		t.Errorf("a worktree off its branch is refused: %v", err)
	}
	git := &fakeGit{branch: "story/S-0001"}
	res, err := FileSave(repo, git, tokensPath, edited, SaveOptions{Hash: Hash(starter), Story: "S-1", Message: "brand", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Committed || res.Stylesheet != sheetPath {
		t.Errorf("result: %+v", res)
	}
	if read(t, filepath.Join(wt, tokensPath)) != edited || read(t, filepath.Join(wt, sheetPath)) != sheetFor(t, edited) {
		t.Error("the worktree's token file and stylesheet are saved")
	}
	if read(t, filepath.Join(root, tokensPath)) != starter || read(t, filepath.Join(root, sheetPath)) != sheetFor(t, starter) {
		t.Error("the main checkout's are left as they were")
	}
	joined := strings.Join(git.calls, "\n")
	if w := wt + ": git commit --quiet -m docs: [S-0001] brand -- " + tokensPath + " " + sheetPath; !strings.Contains(joined, w) {
		t.Errorf("git should have run %q; ran:\n%s", w, joined)
	}
	if strings.Contains(joined, root+": git add") || strings.Contains(joined, root+": git commit") {
		t.Errorf("nothing is committed in the main checkout:\n%s", joined)
	}
}
