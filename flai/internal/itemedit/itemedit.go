// Package itemedit changes a work item after it was created (S-0085): its
// title, nature, tags, touches, parent, and the body below its heading, in
// one step that is checked and committed the way a document save is
// (ADR-0023). What is the item's state stays flai's and is not reachable from
// here: ID, type, status, transitions, blocked intervals, owner, and dates.
//
// A title lives in several places. A retitle keeps them in step: the front
// matter, the heading, the file's name, the line in the parent's list, the
// story's narrative, and links to the old file name in other documents.
package itemedit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Change is what to change; a nil field is left as it is.
type Change struct {
	Title   *string
	Nature  *string
	Tags    *[]string
	Touches *[]string
	Parent  *string
	Body    *string // what lies below the heading; the heading is flai's
}

// Options parameterise an edit.
type Options struct {
	Hash     string // of the file as it was read; empty skips the conflict check
	By       string // who edits, for the notification to agents
	Message  string // commit subject after the prefix; default names what changed
	Trailers []string
	NoCommit bool
	Now      time.Time
}

// View is an item as an editor loads it.
type View struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Status   string   `json:"status"`
	Title    string   `json:"title"`
	Nature   string   `json:"nature"`
	Tags     []string `json:"tags"`
	Touches  []string `json:"touches"`
	Parent   string   `json:"parent,omitempty"`
	Body     string   `json:"body"` // below the heading
	Path     string   `json:"path"`
	Hash     string   `json:"hash"`
	Editable bool     `json:"editable"`
	Reason   string   `json:"reason,omitempty"` // why not
	// Natures and Parents are what the fields may be set to.
	Natures []string `json:"natures"`
	Parents []Option `json:"parents"`
}

// Option is one value a field may take.
type Option struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Result is a successful edit.
type Result struct {
	ID          string          `json:"id"`
	Path        string          `json:"path"`
	Hash        string          `json:"hash"`
	Changed     []string        `json:"changed"` // title, nature, tags, touches, parent, goal, criteria, notes, body
	Renamed     string          `json:"renamed_from,omitempty"`
	Files       []string        `json:"files"` // every file written or removed, relative to the checkout
	Unchanged   bool            `json:"unchanged,omitempty"`
	Committed   bool            `json:"committed"`
	Commit      string          `json:"commit,omitempty"`
	CommitError string          `json:"commit_error,omitempty"`
	Warnings    []check.Finding `json:"warnings"`
}

// InvalidError is a change that is not allowed, whatever the repository
// holds: the caller's mistake, not a failure. A command says it as a rule, so
// that a dashboard answers 400 and not 500.
type InvalidError struct{ msg string }

func (e *InvalidError) Error() string { return e.msg }

func invalid(format string, args ...any) error {
	return &InvalidError{msg: fmt.Sprintf(format, args...)}
}

var heading = regexp.MustCompile(`(?m)\A\s*# [^\n]*\n?`)

// below returns the body without its first heading.
func below(body string) string {
	return strings.TrimLeft(heading.ReplaceAllString(body, ""), "\n")
}

func root(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

func rel(repo *workitem.Repo, abs string) string {
	if r, err := filepath.Rel(root(repo), abs); err == nil {
		return filepath.ToSlash(r)
	}
	return abs
}

func editable(it *workitem.Item) (bool, string) {
	switch {
	case it.Archived:
		return false, it.ID + " is archived; the archive is not edited"
	case it.Closed():
		return false, it.ID + " is " + it.Status + "; a closed item is not edited"
	}
	return true, ""
}

// Show loads an item for an editor.
func Show(repo *workitem.Repo, id string) (*View, error) {
	it, err := repo.Get(id)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(it.Path)
	if err != nil {
		return nil, err
	}
	v := &View{ID: it.ID, Type: it.Type, Status: it.Status, Title: it.Title, Nature: it.Nature, Tags: orEmpty(it.Tags), Touches: orEmpty(it.Touches),
		Parent: it.Parent, Body: below(it.Body), Path: rel(repo, it.Path), Hash: docedit.Hash(string(data)), Natures: workitem.Natures, Parents: []Option{}}
	v.Editable, v.Reason = editable(it)
	if want := parentType(it.Type); want != "" {
		items, err := repo.List(false)
		if err != nil {
			return nil, err
		}
		for _, p := range items {
			if p.Type == want && !p.Closed() {
				v.Parents = append(v.Parents, Option{ID: p.ID, Title: p.Title})
			}
		}
	}
	return v, nil
}

func orEmpty(l []string) []string {
	if l == nil {
		return []string{}
	}
	return l
}

func parentType(typ string) string {
	switch typ {
	case workitem.Story:
		return workitem.Epic
	case workitem.Task:
		return workitem.Story
	}
	return ""
}

var listValue = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 _./@+-]*$`)

// cleanList trims, drops empties and repeats, and refuses what cannot be a
// tag or a path: a comma would split it in the front matter, a leading dash
// would read as a flag somewhere downstream.
func cleanList(what string, in []string) ([]string, error) {
	out := []string{}
	seen := map[string]bool{}
	for _, v := range in {
		v = strings.TrimSuffix(strings.TrimSpace(v), "/")
		if v == "" || seen[v] {
			continue
		}
		if !listValue.MatchString(v) {
			return nil, invalid("%s %q: letters, digits, and _ . / @ + - only, not starting with a dash", what, v)
		}
		seen[v] = true
		out = append(out, v)
	}
	return out, nil
}

func same(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// section returns the text under a "## name" heading.
func section(body, name string) string {
	re := regexp.MustCompile(`(?ms)^## ` + regexp.QuoteMeta(name) + `[ \t]*\n(.*?)(?:^## |\z)`)
	if m := re.FindStringSubmatch(body); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// bodyChanges names the parts of the body that differ, by the sections an
// agent cares about, else "body".
func bodyChanges(before, after string) []string {
	var out []string
	for _, s := range []struct{ heading, name string }{{"Goal", "goal"}, {"Acceptance criteria", "criteria"}, {"Notes", "notes"}} {
		if section(before, s.heading) != section(after, s.heading) {
			out = append(out, s.name)
		}
	}
	strip := func(b string) string {
		for _, h := range []string{"Goal", "Acceptance criteria", "Notes"} {
			b = strings.Replace(b, section(b, h), "", 1)
		}
		return strings.Join(strings.Fields(b), " ")
	}
	if strip(before) != strip(after) {
		out = append(out, "body")
	}
	return out
}

// file is one file an edit touches, with what it held before.
type file struct {
	path    string
	before  []byte
	existed bool
}

type undo struct{ files map[string]*file }

func (u *undo) note(path string) {
	if _, ok := u.files[path]; ok {
		return
	}
	data, err := os.ReadFile(path)
	u.files[path] = &file{path: path, before: data, existed: err == nil}
}

func (u *undo) restore() error {
	var first error
	for _, f := range u.files {
		var err error
		if f.existed {
			err = os.WriteFile(f.path, f.before, 0o644)
		} else if rmErr := os.Remove(f.path); rmErr != nil && !os.IsNotExist(rmErr) {
			err = rmErr
		}
		if err != nil && first == nil {
			first = err
		}
	}
	return first
}

func (u *undo) paths(repo *workitem.Repo) []string {
	out := make([]string, 0, len(u.files))
	for p := range u.files {
		out = append(out, rel(repo, p))
	}
	sort.Strings(out)
	return out
}

// Apply makes the change.
func Apply(repo *workitem.Repo, r execx.Runner, id string, ch Change, opt Options) (*Result, error) {
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	it, err := repo.Get(id)
	if err != nil {
		return nil, err
	}
	oldData, err := os.ReadFile(it.Path)
	if err != nil {
		return nil, err
	}
	if ok, reason := editable(it); !ok {
		return nil, &docedit.RefusedError{Path: rel(repo, it.Path), Reason: reason}
	}
	if opt.Hash != "" && docedit.Hash(string(oldData)) != opt.Hash {
		return nil, &docedit.ConflictError{Path: rel(repo, it.Path), Current: string(oldData), Hash: docedit.Hash(string(oldData))}
	}

	// ---- validate and work out what changes ----
	var changed []string
	oldTitle, oldPath := it.Title, it.Path
	if ch.Title != nil {
		t := strings.Join(strings.Fields(*ch.Title), " ")
		if t == "" {
			return nil, invalid("a title is required")
		}
		if len(t) > 200 {
			return nil, invalid("a title is one line of at most 200 characters; this one has %d", len(t))
		}
		if t != it.Title {
			it.Title = t
			changed = append(changed, "title")
		}
	}
	if ch.Nature != nil && *ch.Nature != it.Nature {
		ok := false
		for _, n := range workitem.Natures {
			ok = ok || n == *ch.Nature
		}
		if !ok {
			return nil, invalid("nature %q must be one of %s", *ch.Nature, strings.Join(workitem.Natures, ", "))
		}
		it.Nature = *ch.Nature
		changed = append(changed, "nature")
	}
	if ch.Tags != nil {
		tags, err := cleanList("tag", *ch.Tags)
		if err != nil {
			return nil, err
		}
		if !same(tags, orEmpty(it.Tags)) {
			it.Tags = tags
			changed = append(changed, "tags")
		}
	}
	if ch.Touches != nil {
		if it.Type == workitem.Epic {
			return nil, invalid("%s is an epic; touches belong to stories and tasks", it.ID)
		}
		touches, err := cleanList("touches", *ch.Touches)
		if err != nil {
			return nil, err
		}
		if !same(touches, orEmpty(it.Touches)) {
			it.Touches = touches
			if len(touches) == 0 {
				it.Touches = nil
			}
			changed = append(changed, "touches")
		}
	}
	var newParent, formerParent *workitem.Item
	if ch.Parent != nil && workitem.CanonicalID(*ch.Parent) != workitem.CanonicalID(it.Parent) {
		want := parentType(it.Type)
		if want == "" {
			return nil, invalid("%s is an epic; epics have no parent", it.ID)
		}
		p, err := repo.Get(*ch.Parent)
		if err != nil {
			return nil, err
		}
		if p.Type != want {
			return nil, invalid("%s is a %s; a %s's parent must be a %s", p.ID, p.Type, it.Type, want)
		}
		if p.Archived || p.Closed() {
			return nil, invalid("%s is %s; pick an open parent", p.ID, p.Status)
		}
		newParent = p
		if it.Parent != "" {
			if fp, err := repo.Get(it.Parent); err == nil && !fp.Archived {
				formerParent = fp
			}
		}
		it.Parent = p.ID
		changed = append(changed, "parent")
	}
	body := below(it.Body)
	if ch.Body != nil {
		nb := strings.TrimSpace(strings.ReplaceAll(*ch.Body, "\r\n", "\n"))
		if nb == "" {
			return nil, invalid("write what the item is for: the body is empty")
		}
		if heading.MatchString(nb) && strings.HasPrefix(strings.TrimSpace(nb), "# ") {
			return nil, invalid("the body starts with a heading of its own; the item's heading is its ID and title, and flai writes it")
		}
		if nb != strings.TrimSpace(body) {
			changed = append(changed, bodyChanges(body, nb)...)
			body = nb + "\n"
		}
	}
	res := &Result{ID: it.ID, Path: rel(repo, it.Path), Hash: docedit.Hash(string(oldData)), Changed: []string{}, Files: []string{}, Warnings: []check.Finding{}}
	if len(changed) == 0 {
		res.Unchanged = true
		return res, nil
	}
	res.Changed = changed
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	it.Body = fmt.Sprintf("# %s %s\n\n%s", it.ID, it.Title, strings.TrimLeft(body, "\n"))
	it.Updated = now

	// ---- write, keeping what was there ----
	before, err := check.Run(repo, opt.Now)
	if err != nil {
		return nil, err
	}
	u := &undo{files: map[string]*file{}}
	fail := func(err error) (*Result, error) {
		if rerr := u.restore(); rerr != nil {
			return nil, fmt.Errorf("%w; and putting the files back failed: %w", err, rerr)
		}
		return nil, err
	}
	retitled := it.Title != oldTitle
	if retitled {
		it.Path = filepath.Join(filepath.Dir(oldPath), workitem.FileName(it.ID, it.Title))
		if it.Path != oldPath {
			if _, err := os.Stat(it.Path); err == nil {
				return nil, invalid("%s already exists", rel(repo, it.Path))
			}
			res.Renamed = rel(repo, oldPath)
		}
	}
	u.note(oldPath)
	u.note(it.Path)
	if err := repo.Save(it); err != nil {
		return fail(err)
	}
	if it.Path != oldPath {
		if err := os.Remove(oldPath); err != nil {
			return fail(err)
		}
	}
	// the parents' lists
	if formerParent != nil {
		u.note(formerParent.Path)
		formerParent.Body = dropChild(formerParent.Body, it.ID)
		formerParent.Updated = now
		if err := repo.Save(formerParent); err != nil {
			return fail(err)
		}
	}
	if newParent != nil {
		u.note(newParent.Path)
		workitem.AppendChild(newParent, it)
		newParent.Updated = now
		if err := repo.Save(newParent); err != nil {
			return fail(err)
		}
	} else if retitled && it.Parent != "" {
		if p, err := repo.Get(it.Parent); err == nil && !p.Archived {
			if nb := retitleChild(p.Body, it.ID, it.Title); nb != p.Body {
				u.note(p.Path)
				p.Body, p.Updated = nb, now
				if err := repo.Save(p); err != nil {
					return fail(err)
				}
			}
		}
	}
	if retitled {
		if err := retitleNarrative(repo, u, it, oldTitle); err != nil {
			return fail(err)
		}
		if it.Path != oldPath {
			if err := relink(repo, u, filepath.Base(oldPath), filepath.Base(it.Path)); err != nil {
				return fail(err)
			}
		}
	}
	if items, err := repo.List(false); err == nil {
		index := filepath.Join(repo.AgentsDir(), "index.md")
		u.note(index)
		if err := repo.WriteIndex(items, opt.Now); err != nil {
			return fail(err)
		}
	}

	// ---- the repository's check, with the change in place ----
	after, err := check.Run(repo, opt.Now)
	if err != nil {
		return fail(err)
	}
	known := map[string]bool{}
	for _, f := range before.Findings {
		known[f.Rule+"\x00"+strings.Replace(f.Path, filepath.Base(oldPath), filepath.Base(it.Path), 1)+"\x00"+f.Message] = true
	}
	var blocking []check.Finding
	own := rel(repo, it.Path)
	for _, f := range after.Findings {
		onPath := filepath.ToSlash(f.Path) == own || strings.HasSuffix(filepath.ToSlash(f.Path), "/"+filepath.Base(it.Path))
		isNew := !known[f.Rule+"\x00"+f.Path+"\x00"+f.Message]
		switch {
		case isNew || (onPath && f.Level == check.Error):
			blocking = append(blocking, f)
		case onPath:
			res.Warnings = append(res.Warnings, f)
		}
	}
	if len(blocking) > 0 {
		if err := u.restore(); err != nil {
			return nil, fmt.Errorf("the check refused the edit of %s and putting the files back failed: %w", it.ID, err)
		}
		return nil, &docedit.RefusedError{Path: own, Reason: fmt.Sprintf("flai check has %d finding(s) with this change; nothing was changed", len(blocking)), Findings: blocking}
	}
	data, err := os.ReadFile(it.Path)
	if err != nil {
		return fail(err)
	}
	res.Path, res.Hash, res.Files = own, docedit.Hash(string(data)), u.paths(repo)
	note(repo, Notice{At: now, By: opt.By, ID: it.ID, Title: it.Title, Type: it.Type, Changed: changed})

	if opt.NoCommit || !repo.Manifest.Autocommit() {
		return res, nil
	}
	subject := opt.Message
	if subject == "" {
		subject = "edit " + strings.Join(changed, ", ")
	}
	sha, inRepo, err := commit(r, root(repo), res.Files, fmt.Sprintf("chore: [%s] %s", it.ID, subject), opt.Trailers)
	switch {
	case !inRepo:
	case err != nil:
		res.CommitError = err.Error()
	default:
		res.Committed, res.Commit = true, sha
	}
	return res, nil
}

var childLine = func(id string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^- ` + regexp.QuoteMeta(id) + `\b[^\n]*\n?`)
}

func dropChild(body, id string) string { return childLine(id).ReplaceAllString(body, "") }

func retitleChild(body, id, title string) string {
	return childLine(id).ReplaceAllStringFunc(body, func(line string) string {
		nl := ""
		if strings.HasSuffix(line, "\n") {
			nl = "\n"
		}
		return "- " + id + " " + title + nl
	})
}

var narrativeTitle = regexp.MustCompile(`(?m)^title: .*$`)

// retitleNarrative keeps a story's narrative in step: its title in the front
// matter and its heading.
func retitleNarrative(repo *workitem.Repo, u *undo, it *workitem.Item, oldTitle string) error {
	path := filepath.Join(repo.AgentsDir(), it.ID+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil //nolint:nilerr // no narrative yet, nothing to keep in step
	}
	doc := string(data)
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil {
		return nil //nolint:nilerr // not a narrative flai wrote; left alone
	}
	newFM := narrativeTitle.ReplaceAllLiteralString(fm, "title: "+workitem.Scalar(it.Title))
	newBody := strings.Replace(body, "# "+it.ID+" "+oldTitle, "# "+it.ID+" "+it.Title, 1)
	if newFM == fm && newBody == body {
		return nil
	}
	u.note(path)
	return os.WriteFile(path, []byte("---\n"+newFM+"---\n"+newBody), 0o644)
}

// relink rewrites links to the old file name in the Markdown of the
// manifest's three folders, the archive excepted: it is not edited.
func relink(repo *workitem.Repo, u *undo, oldName, newName string) error {
	archive := repo.ArchiveDir()
	keys := []string{"wip", "design", "docs"}
	if repo.IsWorktree() {
		// design and docs are this worktree's own, on another branch; one
		// commit cannot span two checkouts, so only wip is kept in step here
		keys = keys[:1]
	}
	for _, key := range keys {
		dir := repo.Manifest.Dir(root(repo), key)
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil //nolint:nilerr // a folder that is not there has no links
			}
			if d.IsDir() {
				if path == archive || strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".md") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil || !strings.Contains(string(data), oldName) {
				return nil //nolint:nilerr // unreadable files are the check's to report
			}
			u.note(path)
			return os.WriteFile(path, []byte(strings.ReplaceAll(string(data), oldName, newName)), 0o644)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// commit records the files of one edit, removals included, on the current
// branch of the checkout.
func commit(r execx.Runner, dir string, files []string, msg string, trailers []string) (sha string, inRepo bool, err error) {
	if _, err := r.Run(dir, "git", "rev-parse", "--is-inside-work-tree"); err != nil {
		return "", false, nil //nolint:nilerr // no repository means nothing to commit to, by design
	}
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	if _, err := r.Run(dir, "git", append([]string{"add", "-A", "--"}, files...)...); err != nil {
		return "", true, err
	}
	if _, err := r.Run(dir, "git", append([]string{"commit", "-q", "-m", msg, "--"}, files...)...); err != nil {
		return "", true, err
	}
	out, err := r.Run(dir, "git", "rev-parse", "--short", "HEAD")
	return strings.TrimSpace(out), true, err
}
