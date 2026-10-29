package digest

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// order is the front matter's keys in the order New writes them.
var order = []string{"id", "title", "kind", "sources", "topics", "paths", "date", "status", "absorbed_by", "asked_again", "cited"}

// NewOptions describe an entry to write.
type NewOptions struct {
	Title   string // the fact, as one sentence
	Kind    string // one of Kinds
	Sources []string
	Topics  []string
	Paths   []string
	Body    string // below the heading: what was asked or decided, the answer and its reason
	Now     time.Time
	Env     Env // what Validate checks sources and absorbed_by against
}

// EditOptions are the changes Edit makes. A nil field is left as it is; a
// non-nil one replaces the entry's own, and an empty list empties it.
type EditOptions struct {
	Body       *string
	Topics     *[]string
	Paths      *[]string
	Sources    *[]string
	AskedAgain *[]string
	Env        Env
}

// Result is the entry New or Edit wrote.
type Result struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Title string `json:"title"`
	Entry Entry  `json:"entry"`
}

// RefusedError is a write the check refused: each problem the write would
// have introduced. Nothing it refused is left written.
type RefusedError struct {
	Path     string
	Problems []Problem
}

// Error names the entry and lists each problem on a line of its own.
func (e *RefusedError) Error() string {
	lines := []string{fmt.Sprintf("the digest entry %s has %d problem(s), so nothing was written; fix them and run it again:", e.Path, len(e.Problems))}
	for _, p := range e.Problems {
		lines = append(lines, "  "+e.Path+": "+p.String())
	}
	return strings.Join(lines, "\n")
}

// NextNumber is one more than the highest entry number in this checkout, on
// main, in any worktree, and on any story branch, so parallel writers do not
// take the same one. Gaps are not filled. A nil r is execx.System.
func NextNumber(repo *workitem.Repo, r execx.Runner) (int, error) {
	if r == nil {
		r = execx.System{}
	}
	folder := Dir(repo.Manifest)
	files, err := os.ReadDir(filepath.Join(repo.Root, filepath.FromSlash(folder)))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("reading the digest folder %s: %w", folder, err)
	}
	names := storygit.FolderNames(r, repo, folder)
	for _, f := range files {
		names = append(names, f.Name())
	}
	next := 1
	for _, name := range names {
		if entryFile.MatchString(name) {
			next = max(next, number(name)+1)
		}
	}
	return next, nil
}

// slug is a title's slug for a file name, cut at template.SlugMax bytes on a word.
func slug(title string) string {
	s := template.CutSlug(template.Slug(title), template.SlugMax)
	if s == "" {
		return "entry"
	}
	return s
}

// New writes a new entry in the checkout's digest folder, creating the
// folder when it is missing, and checks what it wrote with Validate and
// opts.Env. When the check reports anything, it removes the file and any
// folder it created and returns a *RefusedError.
func New(repo *workitem.Repo, r execx.Runner, opts NewOptions) (Result, error) {
	title := strings.Join(strings.Fields(opts.Title), " ")
	if title == "" {
		return Result{}, errors.New("rule: a digest entry needs a title: the fact, as one sentence")
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}
	n, err := NextNumber(repo, r)
	if err != nil {
		return Result{}, err
	}
	e := Entry{
		ID:         fmt.Sprintf("D-%0*d", workitem.IDWidth, n),
		Title:      title,
		Kind:       strings.TrimSpace(opts.Kind),
		Sources:    clean(opts.Sources, canonicalSource),
		Topics:     clean(opts.Topics, nil),
		Paths:      clean(opts.Paths, filepath.ToSlash),
		Date:       opts.Now.UTC().Format(time.DateOnly),
		Status:     "active",
		AskedAgain: []string{},
		Cited:      []string{},
	}
	name := e.ID + "-" + slug(title) + ".md"
	e.Path = path.Join(Dir(repo.Manifest), name)
	dir := filepath.Join(repo.Root, filepath.FromSlash(Dir(repo.Manifest)))
	file := filepath.Join(dir, name)

	created, err := mkdirs(dir)
	if err != nil {
		return Result{}, fmt.Errorf("creating the digest folder %s: %w", Dir(repo.Manifest), err)
	}
	undo := func() error {
		errs := []error{}
		if err := os.Remove(file); err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, err)
		}
		for _, d := range created {
			if err := os.Remove(d); err != nil && !errors.Is(err, os.ErrNotExist) {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}
	fail := func(err error) (Result, error) {
		if uerr := undo(); uerr != nil {
			return Result{}, errors.Join(err, fmt.Errorf("removing the refused entry %s failed; remove it by hand: %w", e.Path, uerr))
		}
		return Result{}, err
	}
	content := frontMatter(e) + "\n# " + e.ID + " " + title + "\n"
	if body := strings.TrimSpace(opts.Body); body != "" {
		content += "\n" + body + "\n"
	}
	if err := atomicfile.WriteFile(file, []byte(content), 0o644); err != nil {
		return fail(fmt.Errorf("writing the digest entry %s: %w", e.Path, err))
	}
	written, err := os.ReadFile(file)
	if err != nil {
		return fail(fmt.Errorf("reading back the digest entry %s: %w", e.Path, err))
	}
	got, problems := Validate(name, string(written), opts.Env)
	if len(problems) > 0 {
		return fail(&RefusedError{Path: e.Path, Problems: problems})
	}
	got.Path = e.Path
	return Result{ID: e.ID, Path: e.Path, Title: title, Entry: got}, nil
}

// Edit changes an entry's body, topics, paths, sources, or asked_again, and
// keeps every other byte of the file. When the check reports a problem the
// change introduced, it puts the file back as it was and returns a
// *RefusedError. A change that leaves the file as it was writes nothing.
func Edit(repo *workitem.Repo, id string, opts EditOptions) (Result, error) {
	e, err := Get(repo, id)
	if err != nil {
		return Result{}, err
	}
	name := path.Base(e.Path)
	file := filepath.Join(repo.Root, filepath.FromSlash(e.Path))
	data, err := os.ReadFile(file)
	if err != nil {
		return Result{}, fmt.Errorf("reading the digest entry %s: %w", e.Path, err)
	}
	fm, body, err := workitem.SplitFrontMatter(string(data))
	if err != nil {
		return Result{}, fmt.Errorf("the digest entry %s has no front matter to edit (%w); fix it by hand first", e.Path, err)
	}
	for _, k := range []struct {
		key   string
		list  *[]string
		canon func(string) string
	}{{"sources", opts.Sources, canonicalSource}, {"topics", opts.Topics, nil}, {"paths", opts.Paths, filepath.ToSlash}, {"asked_again", opts.AskedAgain, canonicalSource}} {
		if k.list != nil {
			fm = setKey(fm, k.key, flowList(clean(*k.list, k.canon)))
		}
	}
	if opts.Body != nil {
		head, _ := splitHeading(body)
		if head == "" {
			head = "\n# " + e.ID + " " + e.Title + "\n"
		}
		if !strings.HasSuffix(head, "\n") {
			head += "\n"
		}
		body = head
		if b := strings.TrimSpace(*opts.Body); b != "" {
			body += "\n" + b + "\n"
		}
	}
	text := "---\n" + fm + "---\n" + body
	res := Result{ID: e.ID, Path: e.Path, Title: e.Title, Entry: e}
	if text == string(data) {
		return res, nil
	}
	_, before := Validate(name, string(data), opts.Env)
	if err := atomicfile.WriteFile(file, []byte(text), 0o644); err != nil {
		return Result{}, fmt.Errorf("writing the digest entry %s: %w", e.Path, err)
	}
	got, after := Validate(name, text, opts.Env)
	if added := introduced(before, after); len(added) > 0 {
		if err := atomicfile.WriteFile(file, data, 0o644); err != nil {
			return Result{}, errors.Join(&RefusedError{Path: e.Path, Problems: added}, fmt.Errorf("putting %s back failed; restore it with git: %w", e.Path, err))
		}
		return Result{}, &RefusedError{Path: e.Path, Problems: added}
	}
	got.Path = e.Path
	res.Title, res.Entry = got.Title, got
	return res, nil
}

// introduced is the problems in after that before does not have as many of,
// compared by rule, field, and message, so that a problem a change only
// moved to another line is not counted.
func introduced(before, after []Problem) []Problem {
	seen := map[string]int{}
	for _, p := range before {
		seen[p.Rule+"\x00"+p.Field+"\x00"+p.Message]++
	}
	var out []Problem
	for _, p := range after {
		k := p.Rule + "\x00" + p.Field + "\x00" + p.Message
		if seen[k] > 0 {
			seen[k]--
			continue
		}
		out = append(out, p)
	}
	return out
}

// mkdirs creates dir and returns the folders it created, deepest first.
func mkdirs(dir string) ([]string, error) {
	var created []string
	for d := dir; ; d = filepath.Dir(d) {
		if _, err := os.Stat(d); err == nil || filepath.Dir(d) == d {
			break
		}
		created = append(created, d)
	}
	return created, os.MkdirAll(dir, 0o755)
}

// frontMatter renders an entry's front matter in the order New writes it.
func frontMatter(e Entry) string {
	values := map[string]string{
		"id":          e.ID,
		"title":       workitem.Scalar(e.Title),
		"kind":        workitem.Scalar(e.Kind),
		"sources":     flowList(e.Sources),
		"topics":      flowList(e.Topics),
		"paths":       flowList(e.Paths),
		"date":        e.Date,
		"status":      e.Status,
		"absorbed_by": workitem.Scalar(e.AbsorbedBy),
		"asked_again": flowList(e.AskedAgain),
		"cited":       flowList(e.Cited),
	}
	var b strings.Builder
	b.WriteString("---\n")
	for _, k := range order {
		b.WriteString(k + ": " + values[k] + "\n")
	}
	b.WriteString("---\n")
	return b.String()
}

// flowList renders a list as [a, b], quoting an item YAML would misread.
func flowList(items []string) string {
	parts := make([]string, len(items))
	for i, s := range items {
		if strings.ContainsAny(s, ",[]{}") {
			parts[i] = strconv.Quote(s)
		} else {
			parts[i] = workitem.Scalar(s)
		}
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// clean trims each item, drops empty ones and repeats, and writes each in
// canon's form when canon is given.
func clean(in []string, canon func(string) string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, s := range in {
		s = strings.TrimSpace(s)
		if canon != nil {
			s = canon(s)
		}
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// setKey replaces a front matter key's value, with any lines it continues
// on, or adds the key after the last key that comes before it in order.
func setKey(fm, key, value string) string {
	lines := strings.SplitAfter(fm, "\n")
	line := key + ": " + value + "\n"
	if from, to, ok := keyBlock(lines, key); ok {
		return strings.Join(lines[:from], "") + line + strings.Join(lines[to:], "")
	}
	at := 0
	for _, k := range order {
		if k == key {
			break
		}
		if _, to, ok := keyBlock(lines, k); ok {
			at = to
		}
	}
	return strings.Join(lines[:at], "") + line + strings.Join(lines[at:], "")
}

// keyBlock is the lines a top-level key holds: its own and those indented
// or listed under it.
func keyBlock(lines []string, key string) (from, to int, ok bool) {
	for i, l := range lines {
		if !strings.HasPrefix(l, key+":") {
			continue
		}
		j := i + 1
		for j < len(lines) && lines[j] != "" && (strings.HasPrefix(lines[j], " ") || strings.HasPrefix(lines[j], "\t") || strings.HasPrefix(lines[j], "- ")) {
			j++
		}
		return i, j, true
	}
	return 0, 0, false
}
