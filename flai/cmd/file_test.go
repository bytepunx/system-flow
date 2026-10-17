package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The files flai file show and flai file save reach in a fileCmdProject: a
// token file and a stylesheet the manifest names under the web sub-project,
// and a file under the design folder.
const (
	fileTokens  = "web/tokens.json"
	fileSheet   = "web/src/theme.css"
	filePalette = "design/system/palette.json"
	palette     = "{\"accent\": \"main\"}\n"
)

// starterTokens is the starter token file, and editedTokens is it with
// another primary colour, so that its stylesheet differs.
var (
	starterTokens = string(designsystem.StarterTokens)
	editedTokens  = strings.Replace(starterTokens, `"components": [0.145, 0.388, 0.922], "hex": "#2563eb"`, `"components": [0.114, 0.306, 0.847], "hex": "#1d4ed8"`, 1)
)

// fileCmdProject is a docProject whose manifest names fileTokens and
// fileSheet under the web sub-project, holding the starter token file, its
// stylesheet, a palette, and a file out of reach, all committed.
func fileCmdProject(t *testing.T) string {
	t.Helper()
	root := docProject(t)
	writeProjectFile(t, root, manifest.File, readProjectFile(t, root, manifest.File)+
		"projects:\n  - name: web\n    path: web\n    kind: sveltekit\n"+
		"design:\n  tokens: "+fileTokens+"\n  stylesheet: "+fileSheet+"\n")
	writeProjectFile(t, root, fileTokens, starterTokens)
	writeProjectFile(t, root, fileSheet, sheetOf(t, fileTokens, starterTokens))
	writeProjectFile(t, root, filePalette, palette)
	writeProjectFile(t, root, "src/main.go", "package main\n")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "design system")
	return root
}

// showFile is flai file show --json's content and hash for path.
func showFile(t *testing.T, root, path string, flags ...string) (content, hash string) {
	t.Helper()
	out, errOut, code := runIn(t, root, append([]string{"file", "show", path, "--json"}, flags...)...)
	if code != 0 {
		t.Fatalf("file show %s: %s", path, errOut)
	}
	var f docedit.File
	if err := json.Unmarshal([]byte(out), &f); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if f.Path != path || f.Hash != docedit.Hash(f.Content) {
		t.Fatalf("file show --json: %s", out)
	}
	return f.Content, f.Hash
}

// S-0436: flai file show prints the raw content on standard output and the
// path and hash on standard error, or {path, content, hash} with --json.
func TestFileShow(t *testing.T) {
	root := fileCmdProject(t)
	out, errOut, code := runIn(t, root, "file", "show", "./"+filePalette)
	if code != 0 || out != palette {
		t.Fatalf("text: %d %q %s", code, out, errOut)
	}
	if want := filePalette + "\n  hash: " + docedit.Hash(palette) + "\n"; errOut != want {
		t.Errorf("standard error %q, want %q", errOut, want)
	}
	if content, _ := showFile(t, root, fileTokens); content != starterTokens {
		t.Error("the token file is shown as it is")
	}
	if content, _ := showFile(t, root, fileSheet); content != sheetOf(t, fileTokens, starterTokens) {
		t.Error("the stylesheet is shown as it is")
	}
}

// S-0436: a path out of reach, the wip folder, a markdown file, a binary
// file, and a file over the size limit are refused by show and save with
// exit status 1 and the reason, and nothing is written.
func TestFileRefusesWhatItDoesNotReach(t *testing.T) {
	root := fileCmdProject(t)
	writeProjectFile(t, root, "wip/notes.txt", "notes\n")
	writeProjectFile(t, root, "design/system/blob.bin", "GIF89a\x00\x01")
	writeProjectFile(t, root, "design/system/big.txt", strings.Repeat("x", docedit.MaxFileSize+1))
	for path, want := range map[string]string{
		"src/main.go":            "src/main.go is outside the design and docs folders and is not the design token file " + fileTokens + " or the design stylesheet " + fileSheet,
		"../outside.json":        "is not a path inside the repository",
		"wip/notes.txt":          "wip/notes.txt is under the wip folder wip/",
		"design/system/plan.md":  "is a markdown file; read and save it with flai doc show and flai doc save",
		"design/system/blob.bin": "is binary",
		"design/system/big.txt":  "over the 1048576-byte limit",
		"design/system/none.txt": "design/system/none.txt does not exist",
	} {
		if out, errOut, code := runIn(t, root, "file", "show", path); code != 1 || out != "" || !strings.Contains(errOut, want) {
			t.Errorf("show %s: %d %q %s", path, code, out, errOut)
		}
		if _, errOut, code := runStdin(t, root, "x\n", "file", "save", path, "--hash", docedit.Hash("notes\n")); code != 1 || !strings.Contains(errOut, want) {
			t.Errorf("save %s: %d %s", path, code, errOut)
		}
	}
	if got := readProjectFile(t, root, "wip/notes.txt"); got != "notes\n" {
		t.Errorf("nothing is written: %q", got)
	}
	if _, errOut, code := runStdin(t, root, "x\n", "file", "save", filePalette); code != exitDocConflict {
		t.Errorf("save without --hash on a file that exists is a conflict: %d %s", code, errOut)
	}
	if got := readProjectFile(t, root, filePalette); got != palette {
		t.Errorf("a save without --hash leaves a file that exists as it was: %q", got)
	}
}

// S-0443, ADR-0164: a save without --hash creates a file that is absent,
// with its folders, and commits it.
func TestFileSaveWithoutHashCreates(t *testing.T) {
	root := fileCmdProject(t)
	path := "design/system/drafts/new.json"
	out, errOut, code := runStdin(t, root, "{}\n", "file", "save", path)
	if code != 0 {
		t.Fatalf("create: %d %s", code, errOut)
	}
	sha := gitIn(t, root, "rev-parse", "--short", "HEAD")
	if want := "saved " + path + ", committed " + sha + "\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	if got := readProjectFile(t, root, path); got != "{}\n" {
		t.Errorf("the created content: %q", got)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); files != path {
		t.Errorf("the commit must hold the created path only, got %q", files)
	}
}

// S-0436: a save writes the content as it is and commits the one path as
// docs, as text and as JSON; the same content again is unchanged.
func TestFileSaveCommitsOnePath(t *testing.T) {
	root := fileCmdProject(t)
	_, hash := showFile(t, root, filePalette)
	mine := "{\"accent\": \"mine\"}\n"
	out, errOut, code := runStdin(t, root, mine, "file", "save", filePalette, "--hash", hash, "--message", "warmer accent", "--trailer", "Co-Authored-By: flaiover <flaiover@localhost>")
	if code != 0 {
		t.Fatalf("save: %d %s", code, errOut)
	}
	sha := gitIn(t, root, "rev-parse", "--short", "HEAD")
	if want := "saved " + filePalette + ", committed " + sha + "\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	if got := readProjectFile(t, root, filePalette); got != mine {
		t.Errorf("the content is saved as it is: %q", got)
	}
	if log := gitIn(t, root, "log", "-1", "--format=%B"); !strings.HasPrefix(log, "docs: warmer accent\n") || !strings.Contains(log, "Co-Authored-By: flaiover <flaiover@localhost>") {
		t.Errorf("commit message:\n%s", log)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); files != filePalette {
		t.Errorf("the commit must hold that path only, got %q", files)
	}

	_, hash = showFile(t, root, filePalette)
	out, errOut, code = runStdin(t, root, mine, "file", "save", filePalette, "--hash", hash, "--json")
	var res docedit.FileResult
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil || !res.Unchanged || res.Committed || res.Stylesheet != "" || res.Hash != hash {
		t.Errorf("unchanged save: %d %s %s", code, out, errOut)
	}
}

// S-0436, ADR-0150: a save of the token file writes the stylesheet from the
// new tokens and commits both paths in one commit, and flai check finds the
// stylesheet in step.
func TestFileSaveTokensCommitsTheStylesheet(t *testing.T) {
	if editedTokens == starterTokens || sheetOf(t, fileTokens, editedTokens) == sheetOf(t, fileTokens, starterTokens) {
		t.Fatal("the edit must change the stylesheet")
	}
	root := fileCmdProject(t)
	_, hash := showFile(t, root, fileTokens)
	out, errOut, code := runStdin(t, root, editedTokens, "file", "save", fileTokens, "--hash", hash, "--message", "darker primary")
	if code != 0 {
		t.Fatalf("save: %d %s", code, errOut)
	}
	sha := gitIn(t, root, "rev-parse", "--short", "HEAD")
	if want := "saved " + fileTokens + " and wrote " + fileSheet + " from it, committed " + sha + "\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	if readProjectFile(t, root, fileTokens) != editedTokens || readProjectFile(t, root, fileSheet) != sheetOf(t, fileTokens, editedTokens) {
		t.Error("the token file and the stylesheet from it are saved")
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); files != fileSheet+"\n"+fileTokens {
		t.Errorf("the commit must hold both paths, got %q", files)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "docs: darker primary" {
		t.Errorf("subject: %q", subject)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("nothing is left uncommitted: %q", st)
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "design.stale") || strings.Contains(out, "design.tokens") {
		t.Errorf("flai check has something to say about the design system:\n%s", out)
	}

	// --json gives the result with the stylesheet written; the default
	// message names both paths.
	out, errOut, code = runStdin(t, root, starterTokens, "file", "save", fileTokens, "--hash", docedit.Hash(editedTokens), "--json")
	var res docedit.FileResult
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil || !res.Committed || res.Stylesheet != fileSheet || res.Path != fileTokens || res.Hash != docedit.Hash(starterTokens) {
		t.Fatalf("save --json: %d %s %s", code, out, errOut)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "docs: edit "+fileTokens+" and write "+fileSheet+" from it" {
		t.Errorf("subject: %q", subject)
	}
	if readProjectFile(t, root, fileSheet) != sheetOf(t, fileTokens, starterTokens) {
		t.Error("the stylesheet follows the tokens back")
	}
}

// S-0436: a stale hash is a conflict, exit 3, with the current content and
// a diff, as JSON and as text, and the other change is kept.
func TestFileSaveConflict(t *testing.T) {
	root := fileCmdProject(t)
	_, hash := showFile(t, root, filePalette)
	theirs := "{\"accent\": \"theirs\"}\n"
	writeProjectFile(t, root, filePalette, theirs)
	mine := "{\"accent\": \"mine\"}\n"
	out, errOut, code := runStdin(t, root, mine, "file", "save", filePalette, "--hash", hash, "--json")
	if code != exitDocConflict || !strings.Contains(errOut, "conflict:") {
		t.Fatalf("conflict: %d %s", code, errOut)
	}
	var res struct{ Conflict docedit.ConflictError }
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	c := res.Conflict
	if c.Path != filePalette || c.Current != theirs || c.Hash != docedit.Hash(theirs) || !strings.Contains(c.Diff, "-{\"accent\": \"theirs\"}") || !strings.Contains(c.Diff, "+{\"accent\": \"mine\"}") {
		t.Errorf("conflict payload: %+v", c)
	}
	out, errOut, code = runStdin(t, root, mine, "file", "save", filePalette, "--hash", hash)
	if code != exitDocConflict || !strings.Contains(out, "+{\"accent\": \"mine\"}") || !strings.Contains(errOut, "conflict: "+filePalette+" changed after it was loaded") {
		t.Errorf("conflict as text: %d %s %s", code, out, errOut)
	}
	if got := readProjectFile(t, root, filePalette); got != theirs {
		t.Errorf("a conflict must not overwrite the other change: %q", got)
	}
}

// S-0436: content flai will not save is refused, exit 4, with its findings,
// and the files are left as they were: a token file it cannot render, a
// stylesheet the check finds out of step with the tokens, and binary
// content.
func TestFileSaveRefused(t *testing.T) {
	root := fileCmdProject(t)
	head := gitIn(t, root, "rev-parse", "HEAD")
	_, tokensHash := showFile(t, root, fileTokens)
	out, errOut, code := runStdin(t, root, "{\"color\": \n", "file", "save", fileTokens, "--hash", tokensHash, "--json")
	if code != exitDocRefused || !strings.Contains(errOut, "refused:") {
		t.Fatalf("tokens that do not parse: %d %s %s", code, out, errOut)
	}
	var res struct{ Refused docedit.RefusedError }
	if err := json.Unmarshal([]byte(out), &res); err != nil || res.Refused.Path != fileTokens || len(res.Refused.Findings) == 0 || res.Refused.Findings[0].Rule != "design.tokens" {
		t.Errorf("refusal payload: %v %s", err, out)
	}

	_, sheetHash := showFile(t, root, fileSheet)
	out, errOut, code = runStdin(t, root, "/* by hand */\n", "file", "save", fileSheet, "--hash", sheetHash)
	if code != exitDocRefused || !strings.Contains(out, ": warning: design.stale: ") || !strings.Contains(errOut, "refused: flai check has") {
		t.Errorf("a stylesheet out of step with the tokens: %d %s %s", code, out, errOut)
	}

	_, paletteHash := showFile(t, root, filePalette)
	if _, errOut, code := runStdin(t, root, "a\x00b", "file", "save", filePalette, "--hash", paletteHash); code != exitDocRefused || !strings.Contains(errOut, "is binary") {
		t.Errorf("binary content: %d %s", code, errOut)
	}

	if readProjectFile(t, root, fileTokens) != starterTokens || readProjectFile(t, root, fileSheet) != sheetOf(t, fileTokens, starterTokens) || readProjectFile(t, root, filePalette) != palette {
		t.Error("a refused save must leave the files as they were")
	}
	if gitIn(t, root, "rev-parse", "HEAD") != head || gitIn(t, root, "status", "--porcelain") != "" {
		t.Error("nothing is committed or left behind")
	}
}

// S-0436, ADR-0145: with --story, show reads and save writes, checks, and
// commits the token file and its stylesheet in the story's worktree, on its
// branch, and the main checkout is left as it was.
func TestFileSaveWithAStory(t *testing.T) {
	root := fileCmdProject(t)
	if _, errOut, code := runIn(t, root, "file", "show", filePalette, "--story", "S-0001"); code != 1 || !strings.Contains(errOut, "S-0001: it has no worktree") {
		t.Errorf("a story without a worktree: %d %s", code, errOut)
	}
	wt := storyWorktreeIn(t, root)
	writeProjectFile(t, wt, filePalette, "{\"accent\": \"branch\"}\n")
	gitIn(t, wt, "commit", "-q", "-am", "docs: [S-0001] branch palette")
	mainHead := gitIn(t, root, "rev-parse", "HEAD")

	if content, _ := showFile(t, root, filePalette, "--story", "S-0001"); content != "{\"accent\": \"branch\"}\n" {
		t.Errorf("show --story reads the branch's file: %q", content)
	}
	_, hash := showFile(t, root, fileTokens, "--story", "S-0001")
	out, errOut, code := runStdin(t, root, editedTokens, "file", "save", fileTokens, "--story", "S-0001", "--hash", hash, "--message", "brand", "--json")
	var res docedit.FileResult
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil || !res.Committed || res.Stylesheet != fileSheet {
		t.Fatalf("save --story: %d %s %s", code, out, errOut)
	}
	if got := gitIn(t, root, "rev-parse", "--short", "story/S-0001"); got != res.Commit {
		t.Errorf("the commit %s should be story/S-0001's head, which is %s", res.Commit, got)
	}
	if subject := gitIn(t, wt, "log", "-1", "--format=%s"); subject != "docs: [S-0001] brand" {
		t.Errorf("subject: %q", subject)
	}
	if files := gitIn(t, wt, "show", "--name-only", "--format=", "HEAD"); files != fileSheet+"\n"+fileTokens {
		t.Errorf("the commit must hold both paths, got %q", files)
	}
	if readProjectFile(t, wt, fileTokens) != editedTokens || readProjectFile(t, wt, fileSheet) != sheetOf(t, fileTokens, editedTokens) {
		t.Error("the worktree's token file and stylesheet are saved")
	}
	if readProjectFile(t, root, fileTokens) != starterTokens || readProjectFile(t, root, fileSheet) != sheetOf(t, fileTokens, starterTokens) {
		t.Error("the main checkout's are left as they were")
	}
	if gitIn(t, root, "rev-parse", "HEAD") != mainHead || gitIn(t, root, "status", "--porcelain") != "" {
		t.Error("the main checkout is left as it was")
	}
}
