package docedit

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// MaxFileSize is the largest text file, in bytes, that FileShow reads and
// FileSave writes (S-0436): 1 MiB, far above a token file or a stylesheet.
const MaxFileSize = 1 << 20

// binarySniff is how much of a file's start is searched for a NUL byte,
// the mark of a binary file, as git sniffs one.
const binarySniff = 8 << 10

// File is a text file that is not markdown as an editor loads it (ADR-0147).
type File struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Hash    string `json:"hash"`
}

// FileResult is a successful FileSave: the save's result and, when it saved
// the design token file, the stylesheet it wrote from them (ADR-0150).
type FileResult struct {
	Result
	Stylesheet string `json:"stylesheet,omitempty"`
}

// runCheck is the check FileSave runs before and after it writes: a
// variable so that a test can give the finding no valid token file gives.
var runCheck = check.Run

// FileReach maps a repository path to a file FileShow and FileSave reach
// under the manifest m, and returns it cleaned: a file under the design or
// docs folder, or the design token file or stylesheet m names. The wip
// folder, a markdown file, and any other path are refused, each with why.
func FileReach(m manifest.Manifest, rel string) (string, error) {
	clean := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(rel, "./")))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || filepath.IsAbs(clean) || strings.HasPrefix(clean, "/") {
		return "", fmt.Errorf("%q is not a path inside the repository; name a file relative to its root", rel)
	}
	under := func(key string) (string, bool) {
		dir := strings.TrimSuffix(m.Layout[key], "/")
		return dir, dir != "" && (clean == dir || strings.HasPrefix(clean, dir+"/"))
	}
	if dir, ok := under("wip"); ok {
		return "", fmt.Errorf("%s is under the wip folder %s/, which flai's own commands write; file show and file save do not reach it", clean, dir)
	}
	if strings.EqualFold(path.Ext(clean), ".md") {
		return "", fmt.Errorf("%s is a markdown file; read and save it with flai doc show and flai doc save", clean)
	}
	sheet := m.DesignStylesheet()
	if clean == path.Clean(m.DesignTokens()) || (sheet != "" && clean == path.Clean(sheet)) {
		return clean, nil
	}
	if _, ok := under("design"); ok {
		return clean, nil
	}
	if _, ok := under("docs"); ok {
		return clean, nil
	}
	named := "the design token file " + m.DesignTokens()
	if sheet != "" {
		named += " or the design stylesheet " + sheet
	}
	return "", fmt.Errorf("%s is outside the design and docs folders and is not %s; file show and file save reach only those", clean, named)
}

// FileShow loads a text file that is not markdown for editing, from story's
// worktree when story is not empty (ADR-0145, ADR-0147).
func FileShow(repo *workitem.Repo, rel, story string) (*File, error) {
	at, err := locateFile(repo, rel, story, false)
	if err != nil {
		return nil, err
	}
	content, err := readText(at.abs, at.clean)
	if err != nil {
		return nil, err
	}
	return &File{Path: at.clean, Content: content, Hash: Hash(content)}, nil
}

// FileSave writes content, as it is, to a text file that is not markdown by
// doc.save's rules (ADR-0147): a stale opt.Hash is a *ConflictError, a
// finding of the check the content introduces puts the file back and is a
// *RefusedError, and the save is committed as docs unless autocommit is
// off; with opt.Story, in that story's worktree and on its branch
// (ADR-0145). An empty opt.Hash creates a file that is absent, with its
// folders, which a refusal removes again, and is a *ConflictError on a file
// that exists (ADR-0164). When the file is the design token file and the
// manifest names a stylesheet, the stylesheet is written from the new
// tokens, checked with them, put back with them, and committed with them in
// one commit (ADR-0150).
func FileSave(repo *workitem.Repo, r execx.Runner, rel, content string, opt SaveOptions) (*FileResult, error) {
	at, err := locateFile(repo, rel, opt.Story, opt.Hash == "")
	if err != nil {
		return nil, err
	}
	abs, clean := at.abs, at.clean
	repo, opt.Story = at.repo, at.story
	if opt.Story != "" {
		if err := onStoryBranch(r, repo.Root, opt.Story, clean); err != nil {
			return nil, err
		}
	}
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	if reason := notText([]byte(content)); reason != "" {
		return nil, &RefusedError{Path: clean, Reason: "the content you would save " + reason + "; nothing was saved"}
	}
	var old string
	if !at.absent {
		if old, err = readText(abs, clean); err != nil {
			return nil, err
		}
		// No content hashes to "", so an empty hash on a file that exists is
		// a conflict, as a stale one is.
		if Hash(old) != opt.Hash {
			return nil, conflict(r, abs, clean, content)
		}
	}
	var sheet, css string
	if s := repo.Manifest.DesignStylesheet(); s != "" && clean == path.Clean(repo.Manifest.DesignTokens()) && clean != path.Clean(s) {
		if css, err = StylesheetFor(clean, []byte(content)); err != nil {
			return nil, err
		}
		sheet = path.Clean(s)
	}
	if !at.absent && content == old {
		return &FileResult{Result: Result{Path: clean, Hash: Hash(old), Unchanged: true, Warnings: []check.Finding{}}}, nil
	}

	before, err := runCheck(repo, opt.Now)
	if err != nil {
		return nil, err
	}
	var restore func() error
	if at.absent {
		restore, err = createFile(abs, clean, content)
		if errors.Is(err, fs.ErrExist) {
			// Another writer created it after it was found absent.
			return nil, conflict(r, abs, clean, content)
		}
	} else {
		err = writeAtomic(abs, content)
		restore = func() error { return writeAtomic(abs, old) }
	}
	if err != nil {
		return nil, err
	}
	saved := []string{abs}
	if sheet != "" {
		sheetAbs := filepath.Join(at.root, filepath.FromSlash(sheet))
		oldSheet, readErr := os.ReadFile(sheetAbs)
		putBack := restore
		restore = func() error {
			if err := putBack(); err != nil {
				return err
			}
			if readErr != nil {
				return os.Remove(sheetAbs)
			}
			return writeAtomic(sheetAbs, string(oldSheet))
		}
		if err := writeStylesheet(at.root, sheet, css); err != nil {
			_ = restore()
			return nil, err
		}
		saved = append(saved, sheetAbs)
	}
	after, err := runCheck(repo, opt.Now)
	if err != nil {
		_ = restore()
		return nil, err
	}
	res := &FileResult{Result: Result{Path: clean, Hash: Hash(content)}, Stylesheet: sheet}
	warnings, blocking := judge(repo, before, after, saved...)
	res.Warnings = warnings
	if len(blocking) > 0 {
		if err := restore(); err != nil {
			return nil, fmt.Errorf("the check refused %s and restoring it failed: %w", clean, err)
		}
		return nil, &RefusedError{Path: clean, Reason: fmt.Sprintf("flai check has %d finding(s) with this content; nothing was saved", len(blocking)), Findings: blocking}
	}

	if opt.NoCommit || !repo.Manifest.Autocommit() {
		return res, nil
	}
	paths := make([]string, 0, len(saved))
	for _, p := range saved {
		inRoot, _ := filepath.Rel(at.root, p)
		paths = append(paths, inRoot)
	}
	if sheet != "" && strings.TrimSpace(opt.Message) == "" {
		opt.Message = "edit " + clean + " and write " + sheet + " from it"
	}
	// A failed commit is not a failed save, as for Save.
	sha, inRepo, err := commit(r, at.root, commitMessage(repo, clean, opt), paths...)
	switch {
	case !inRepo:
	case err != nil:
		res.CommitError = err.Error()
	default:
		res.Committed, res.Commit = true, sha
	}
	return res, nil
}

// StylesheetFor is the design stylesheet generated from content, the token
// file tokens as it would be saved (ADR-0150). Content that is not a token
// file flai can render is a *RefusedError with each problem as a
// design.tokens finding.
func StylesheetFor(tokens string, content []byte) (string, error) {
	var sys *designsystem.System
	f, err := designsystem.Parse(content)
	var pe *designsystem.ParseError
	switch {
	case errors.As(err, &pe):
		sys = &designsystem.System{Problems: []designsystem.Problem{pe.Problem()}}
	case err != nil:
		return "", fmt.Errorf("read %s as a token file: %w; fix it and save again", tokens, err)
	default:
		sys = designsystem.Resolve(f)
	}
	if len(sys.Problems) > 0 {
		findings := make([]check.Finding, 0, len(sys.Problems))
		for _, p := range sys.Problems {
			msg := p.Message
			if p.Path != "" {
				msg = p.Path + ": " + msg
			}
			findings = append(findings, check.Finding{Level: check.Error, Rule: "design.tokens", Path: tokens, Line: p.Line, Message: msg})
		}
		return "", &RefusedError{Path: tokens, Reason: fmt.Sprintf("%s has %d problem(s) as a token file, so no stylesheet can be written from it; fix them; nothing was saved", tokens, len(findings)), Findings: findings}
	}
	css := designsystem.Stylesheet(sys, tokens)
	if css == "" {
		return "", &RefusedError{Path: tokens, Reason: tokens + " gives no stylesheet: it resolves to no light and dark theme; give it both; nothing was saved"}
	}
	return css, nil
}

// WriteStylesheet writes the design stylesheet, the repository path sheet
// under root, generated from content, the token file tokens (ADR-0150),
// creating its folders. Content StylesheetFor refuses writes nothing.
func WriteStylesheet(root, sheet, tokens string, content []byte) error {
	css, err := StylesheetFor(tokens, content)
	if err != nil {
		return err
	}
	return writeStylesheet(root, sheet, css)
}

// writeStylesheet writes css to the stylesheet sheet under root, replacing
// it in one step when it is there and creating it and its folders when not.
func writeStylesheet(root, sheet, css string) error {
	abs := filepath.Join(root, filepath.FromSlash(sheet))
	if _, err := os.Stat(abs); err == nil {
		if err := writeAtomic(abs, css); err != nil {
			return fmt.Errorf("write the design stylesheet %s: %w", sheet, err)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return fmt.Errorf("write the design stylesheet %s: %w", sheet, err)
	}
	if err := os.WriteFile(abs, []byte(css), 0o644); err != nil {
		return fmt.Errorf("write the design stylesheet %s: %w", sheet, err)
	}
	return nil
}

// filePlace is where a file edit is made, and the root of the tree that
// holds the file, in which it is committed; absent when the file is not
// there yet and the edit creates it.
type filePlace struct {
	place
	root   string
	absent bool
}

// locateFile resolves rel for a file edit as locate does for a document: in
// the tree flai runs in, falling back to the main checkout, or with a story,
// in that story's worktree only, under its manifest. A file none of them
// holds is refused, unless create, when it is placed, absent, in the first.
func locateFile(repo *workitem.Repo, rel, story string, create bool) (filePlace, error) {
	in, id, roots, where := repo, "", []string{repo.Root, repo.MainRoot}, ""
	if story != "" {
		var err error
		if id, in, err = storyRepo(repo, story); err != nil {
			return filePlace{}, err
		}
		roots, where = []string{in.Root}, fmt.Sprintf(" in %s's worktree; name a file on its branch", id)
	}
	clean, err := FileReach(in.Manifest, rel)
	if err != nil {
		return filePlace{}, err
	}
	for _, root := range roots {
		if root == "" {
			continue
		}
		p := filepath.Join(root, filepath.FromSlash(clean))
		if st, err := os.Stat(p); err == nil && st.Mode().IsRegular() {
			return filePlace{place: place{repo: in, abs: p, clean: clean, story: id}, root: root}, nil
		}
	}
	if create {
		p := filepath.Join(roots[0], filepath.FromSlash(clean))
		return filePlace{place: place{repo: in, abs: p, clean: clean, story: id}, root: roots[0], absent: true}, nil
	}
	return filePlace{}, fmt.Errorf("%s does not exist%s", clean, where)
}

// createFile writes content to the file abs, clean in the repository, which
// must not exist, creating the folders it lacks, and returns what removes
// the file and those folders again. A file that is there by then is an
// error matching fs.ErrExist, and nothing is left behind.
func createFile(abs, clean, content string) (func() error, error) {
	var made []string // deepest first
	for dir := filepath.Dir(abs); ; dir = filepath.Dir(dir) {
		_, err := os.Lstat(dir)
		if err == nil {
			break
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("create %s: %w", clean, err)
		}
		made = append(made, dir)
	}
	undo := func(file bool) error {
		if file {
			if err := os.Remove(abs); err != nil {
				return err
			}
		}
		for _, dir := range made {
			if err := os.Remove(dir); err != nil {
				return err
			}
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		_ = undo(false)
		return nil, fmt.Errorf("create the folders of %s: %w; check that no file stands in their way", clean, err)
	}
	f, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		_ = undo(false)
		return nil, fmt.Errorf("create %s: %w", clean, err)
	}
	_, err = f.WriteString(content)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = undo(true)
		return nil, fmt.Errorf("write %s: %w", clean, err)
	}
	return func() error { return undo(true) }, nil
}

// conflict is the *ConflictError of a save of content over the file at
// abs, clean in the repository, as it is now.
func conflict(r execx.Runner, abs, clean, content string) error {
	current, err := readText(abs, clean)
	if err != nil {
		return err
	}
	return &ConflictError{Path: clean, Current: current, Hash: Hash(current), Diff: diff(r, abs, clean, content)}
}

// readText reads the file at abs, clean in the repository, refusing a file
// over MaxFileSize, before it is read, and a binary one.
func readText(abs, clean string) (string, error) {
	st, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", clean, err)
	}
	if st.Size() > MaxFileSize {
		return "", fmt.Errorf("%s is %d bytes, over the %d-byte limit of file show and file save; edit it in the repository", clean, st.Size(), MaxFileSize)
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", clean, err)
	}
	if reason := notText(data); reason != "" {
		return "", fmt.Errorf("%s %s; edit it in the repository", clean, reason)
	}
	return string(data), nil
}

// notText says why data is not text that file show and file save handle:
// over MaxFileSize, or binary, with a NUL byte near its start; "" when it
// is text.
func notText(data []byte) string {
	switch {
	case len(data) > MaxFileSize:
		return fmt.Sprintf("is over the %d-byte limit of file show and file save", MaxFileSize)
	case bytes.IndexByte(data[:min(len(data), binarySniff)], 0) >= 0:
		return fmt.Sprintf("is binary, with a NUL byte in its first %d bytes, and file show and file save handle text only", binarySniff)
	}
	return ""
}
