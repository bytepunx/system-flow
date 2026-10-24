package digest

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// citationID is an entry ID written anywhere in a text, in any zero padding.
var citationID = regexp.MustCompile(`\bD-\d+\b`)

// storyID is a narrative's stream when it names a story.
var storyID = regexp.MustCompile(`^S-\d+$`)

// Citation is who cited an entry and when last (ADR-0199): the work items,
// the stories whose narratives, and the threads whose files name its ID, in
// ID order, and the latest of those files' updated times, zero when none.
type Citation struct {
	Citers []string
	Latest time.Time
}

// citing is a file that can cite an entry: the file, and the citer it counts
// as, empty for a narrative, whose stream names its story.
type citing struct {
	file  string
	citer string
}

// Citations reads every work item, kanban and archive, every narrative in
// the agents folder and its archive, and every thread for the IDs of the
// entries given, and gives each cited entry's ID its Citation. The digest
// folder is not read: entries citing one another do not count. A citing file
// whose updated time cannot be read is an error naming it.
func Citations(repo *workitem.Repo, entries []Entry) (map[string]Citation, error) {
	known := map[string]bool{}
	for _, e := range entries {
		known[e.ID] = true
	}
	files, err := citingFiles(repo)
	if err != nil {
		return nil, err
	}
	out := map[string]Citation{}
	for _, f := range files {
		data, err := os.ReadFile(f.file)
		if err != nil {
			return nil, fmt.Errorf("reading %s for the digest entries it cites: %w", repoPath(repo, f.file), err)
		}
		var ids []string
		for _, id := range citationID.FindAllString(string(data), -1) {
			if id = CanonicalID(id); known[id] && !slices.Contains(ids, id) {
				ids = append(ids, id)
			}
		}
		if len(ids) == 0 {
			continue
		}
		citer, at, err := citedBy(f, string(data))
		if err != nil {
			return nil, fmt.Errorf("%s cites %s: %w", repoPath(repo, f.file), strings.Join(ids, ", "), err)
		}
		for _, id := range ids {
			c := out[id]
			if !slices.Contains(c.Citers, citer) {
				c.Citers = append(c.Citers, citer)
			}
			if at.After(c.Latest) {
				c.Latest = at
			}
			out[id] = c
		}
	}
	for id, c := range out {
		slices.SortFunc(c.Citers, compareIDs)
		out[id] = c
	}
	return out, nil
}

// citingFiles are the work items' files, kanban and archive, the narratives'
// in the agents folder and its archive, and the threads', each with the
// citer its name gives. A missing folder holds none.
func citingFiles(repo *workitem.Repo) ([]citing, error) {
	type folder struct {
		dir       string
		narrative bool
	}
	var folders []folder
	for _, typ := range workitem.Types {
		folders = append(folders, folder{dir: repo.ItemDir(typ, false)}, folder{dir: repo.ItemDir(typ, true)})
	}
	folders = append(folders, folder{dir: repo.AgentsDir(), narrative: true}, folder{dir: filepath.Join(repo.ArchiveDir(), "agents"), narrative: true}, folder{dir: threads.Dir(repo)})
	var out []citing
	for _, fo := range folders {
		files, err := os.ReadDir(fo.dir)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("listing %s for digest citations: %w", repoPath(repo, fo.dir), err)
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			file := filepath.Join(fo.dir, f.Name())
			if fo.narrative {
				if narrativeFile.MatchString(f.Name()) {
					out = append(out, citing{file: file})
				}
				continue
			}
			if m := sourceFile.FindStringSubmatch(f.Name()); m != nil {
				if id := canonicalSource(m[1]); citedID.MatchString(id) {
					out = append(out, citing{file: file, citer: id})
				}
			}
		}
	}
	return out, nil
}

// citedBy is the citer a file counts as, a narrative's being the story its
// stream names, else the one it is named for, and the file's updated time,
// RFC 3339 or a day, which is that day's midnight in UTC.
func citedBy(f citing, content string) (string, time.Time, error) {
	fm, _, err := workitem.SplitFrontMatter(content)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("it has no front matter, so its citation has no date (%w); fix its front matter", err)
	}
	var head struct {
		Stream  string `yaml:"stream"`
		Updated string `yaml:"updated"`
	}
	if err := yaml.Unmarshal([]byte(fm), &head); err != nil {
		return "", time.Time{}, fmt.Errorf("its front matter does not parse, so its citation has no date: %s; fix it", firstLine(err.Error()))
	}
	citer := f.citer
	if citer == "" {
		citer = "S-" + narrativeFile.FindStringSubmatch(filepath.Base(f.file))[1]
		if s := strings.TrimSpace(head.Stream); storyID.MatchString(s) {
			citer = s
		}
		citer = workitem.CanonicalID(citer)
	}
	at, err := updatedTime(head.Updated)
	if err != nil {
		return "", time.Time{}, err
	}
	return citer, at, nil
}

// updatedTime reads an updated front matter value: an RFC 3339 time, or a
// day, YYYY-MM-DD, which is that day's midnight in UTC.
func updatedTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("its updated is %q, not an RFC 3339 time or a day, so its citation has no date; give it the time it was last changed", s)
}

// compareIDs orders IDs by their prefix, then their number: E, S, T, TH.
func compareIDs(a, b string) int {
	pa, na, _ := strings.Cut(a, "-")
	pb, nb, _ := strings.Cut(b, "-")
	if n := strings.Compare(pa, pb); n != 0 {
		return n
	}
	x, _ := strconv.Atoi(na)
	y, _ := strconv.Atoi(nb)
	return cmp.Compare(x, y)
}

// CitedEntry is an entry with its citers and its latest citation, as Cited
// lists it. Changed says its cited front matter was written, or in a dry
// run of Compact would be.
type CitedEntry struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Status string   `json:"status"`
	Date   string   `json:"date"`
	Path   string   `json:"path"`
	Cited  []string `json:"cited"`
	// Latest is the latest citation, RFC 3339 in UTC, empty when none.
	Latest  string `json:"latest"`
	Changed bool   `json:"changed"`

	latest time.Time
}

// Cited lists every entry with its citers and its latest citation, and
// writes each entry's cited as its citers, in ID order, where that changed.
// Each write is checked as Edit's is; the first the check refuses stops it,
// with the entries before it written.
func Cited(repo *workitem.Repo, env Env) ([]CitedEntry, error) {
	return cite(repo, env, true)
}

// cite is Cited, writing only when write is true.
func cite(repo *workitem.Repo, env Env, write bool) ([]CitedEntry, error) {
	entries, err := List(repo)
	if err != nil {
		return nil, err
	}
	citations, err := Citations(repo, entries)
	if err != nil {
		return nil, err
	}
	out := make([]CitedEntry, 0, len(entries))
	for _, e := range entries {
		c := citations[e.ID]
		citers := c.Citers
		if citers == nil {
			citers = []string{}
		}
		ce := CitedEntry{ID: e.ID, Title: e.Title, Status: e.Status, Date: e.Date, Path: e.Path, Cited: citers, Changed: !slices.Equal(e.Cited, citers), latest: c.Latest}
		if !c.Latest.IsZero() {
			ce.Latest = c.Latest.UTC().Format(workitem.TimeFormat)
		}
		if ce.Changed && write {
			if _, _, err := rewrite(repo, e.Path, env, func(fm, body string) (string, string) {
				return setKey(fm, "cited", flowList(citers)), body
			}); err != nil {
				return out, err
			}
		}
		out = append(out, ce)
	}
	return out, nil
}

// AnchorFunc is the anchor of a heading in a document, given the document's
// path from the project root and its text, or an error naming the headings
// it has when the heading is none of them.
type AnchorFunc func(path, content, heading string) (string, error)

// AbsorbOptions say what document now records an entry's fact.
type AbsorbOptions struct {
	By      string     // the document's path from the project root
	Heading string     // a heading of it, or empty for the whole document
	Anchor  AnchorFunc // finds Heading's anchor; needed only with a Heading
	Env     Env        // what Validate checks the entry against
}

// Absorb marks an entry absorbed by the document opts names (ADR-0199):
// status absorbed and absorbed_by the path, or <path>#<anchor> with a
// heading, keeping the body. An entry absorbed again names the newer
// document. A path that is not local or names no file, or a heading the
// document does not have, is a *RefusedError and writes nothing.
func Absorb(repo *workitem.Repo, id string, opts AbsorbOptions) (Result, error) {
	e, err := Get(repo, id)
	if err != nil {
		return Result{}, err
	}
	refuse := func(format string, args ...any) (Result, error) {
		return Result{}, &RefusedError{Path: e.Path, Problems: []Problem{{Field: "absorbed_by", Line: fieldLine(repo, e.Path, "absorbed_by"), Rule: RuleEntry, Message: fmt.Sprintf(format, args...)}}}
	}
	by := strings.TrimSpace(filepath.ToSlash(opts.By))
	if by == "" {
		return refuse("no document is given; give the path, from the project root, of the document that now records the fact")
	}
	by = path.Clean(by)
	if !filepath.IsLocal(filepath.FromSlash(by)) {
		return refuse("%q is not a path from the project root; give one such as design/system/overview.md", opts.By)
	}
	file := filepath.Join(repo.Root, filepath.FromSlash(by))
	if info, err := os.Stat(file); err != nil || info.IsDir() {
		return refuse("%s names no file; give the path, from the project root, of the document that now records the fact", by)
	}
	value := by
	if heading := strings.TrimSpace(opts.Heading); heading != "" {
		if opts.Anchor == nil {
			return Result{}, errors.New("absorbing by a heading needs a way to find its anchor; this is a flai bug: report it")
		}
		content, err := os.ReadFile(file)
		if err != nil {
			return Result{}, fmt.Errorf("reading %s for its headings: %w", by, err)
		}
		anchor, err := opts.Anchor(by, string(content), heading)
		if err != nil {
			return refuse("%v", err)
		}
		value = by + "#" + anchor
	}
	got, changed, err := rewrite(repo, e.Path, opts.Env, func(fm, body string) (string, string) {
		return setKey(setKey(fm, "status", "absorbed"), "absorbed_by", workitem.Scalar(value)), body
	})
	if err != nil {
		return Result{}, err
	}
	if !changed {
		got = e
	}
	return Result{ID: got.ID, Path: e.Path, Title: got.Title, Entry: got}, nil
}

// CompactOptions say how Compact judges and whether it writes.
type CompactOptions struct {
	Window time.Duration // digest.compact_after; zero turns compaction off
	Now    time.Time
	DryRun bool // list what would be compacted and write nothing
	Env    Env
}

// Compacted is an entry Compact compacted, or in a dry run would, with what
// it was judged by: its latest citation, empty when none, and its date.
type Compacted struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Path   string `json:"path"`
	Latest string `json:"latest"`
	Date   string `json:"date"`
}

// CompactResult is what Compact did: the window, whether compaction is off,
// the time an active entry last cited or dated before is compacted, every
// entry as Cited listed it, and the entries compacted.
type CompactResult struct {
	Window    string       `json:"window"`
	Off       bool         `json:"off"`
	DryRun    bool         `json:"dry_run"`
	Before    string       `json:"before"`
	Cited     []CitedEntry `json:"cited"`
	Compacted []Compacted  `json:"compacted"`
}

// Compact runs Cited, or in a dry run works out the citations without
// writing them, then compacts each active entry whose latest citation, or
// its own date when that is later or it has none, is before opts.Now less
// opts.Window (ADR-0199): status compacted, and its body emptied below its
// "# D-nnnn <title>" heading. Absorbed and compacted entries are left as
// they are, so a second run compacts nothing more. A zero window turns
// compaction off: it does nothing and says so in Off.
func Compact(repo *workitem.Repo, opts CompactOptions) (CompactResult, error) {
	res := CompactResult{Window: opts.Window.String(), DryRun: opts.DryRun, Cited: []CitedEntry{}, Compacted: []Compacted{}}
	if opts.Window < 0 {
		return res, fmt.Errorf("the compaction window %s is negative; give a duration such as 720h, or 0 to turn compaction off", opts.Window)
	}
	if opts.Window == 0 {
		res.Off = true
		return res, nil
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}
	before := opts.Now.Add(-opts.Window).UTC()
	res.Before = before.Format(workitem.TimeFormat)
	cited, err := cite(repo, opts.Env, !opts.DryRun)
	res.Cited = cited
	if err != nil {
		return res, err
	}
	for _, c := range cited {
		if c.Status != "active" {
			continue
		}
		last, err := time.Parse(time.DateOnly, c.Date)
		if err != nil {
			return res, fmt.Errorf("the digest entry %s has date %q, not a day, YYYY-MM-DD, so it cannot be judged for compaction; fix it as flai check says and run it again", c.Path, c.Date)
		}
		if c.latest.After(last) {
			last = c.latest
		}
		if !last.Before(before) {
			continue
		}
		if !opts.DryRun {
			id, title := c.ID, c.Title
			if _, _, err := rewrite(repo, c.Path, opts.Env, func(fm, body string) (string, string) {
				head, _ := splitHeading(body)
				if head == "" {
					head = "\n# " + id + " " + title + "\n"
				}
				if !strings.HasSuffix(head, "\n") {
					head += "\n"
				}
				return setKey(fm, "status", "compacted"), head
			}); err != nil {
				return res, err
			}
		}
		res.Compacted = append(res.Compacted, Compacted{ID: c.ID, Title: c.Title, Path: c.Path, Latest: c.Latest, Date: c.Date})
	}
	return res, nil
}

// rewrite changes the entry file at rel, a path from the project root, with
// change, given the file's front matter and body, and checks the result as
// Edit does: when the check reports a problem the change introduced, it
// puts the file back as it was and returns a *RefusedError. It returns the
// entry as written and whether the file changed; a change that leaves the
// file as it was writes nothing.
func rewrite(repo *workitem.Repo, rel string, env Env, change func(fm, body string) (string, string)) (Entry, bool, error) {
	name := path.Base(rel)
	file := filepath.Join(repo.Root, filepath.FromSlash(rel))
	data, err := os.ReadFile(file)
	if err != nil {
		return Entry{}, false, fmt.Errorf("reading the digest entry %s: %w", rel, err)
	}
	fm, body, err := workitem.SplitFrontMatter(string(data))
	if err != nil {
		return Entry{}, false, fmt.Errorf("the digest entry %s has no front matter to change (%w); fix it by hand first", rel, err)
	}
	fm, body = change(fm, body)
	text := "---\n" + fm + "---\n" + body
	if text == string(data) {
		return Entry{}, false, nil
	}
	_, before := Validate(name, string(data), env)
	if err := atomicfile.WriteFile(file, []byte(text), 0o644); err != nil {
		return Entry{}, false, fmt.Errorf("writing the digest entry %s: %w", rel, err)
	}
	got, after := Validate(name, text, env)
	if added := introduced(before, after); len(added) > 0 {
		if err := atomicfile.WriteFile(file, data, 0o644); err != nil {
			return Entry{}, false, errors.Join(&RefusedError{Path: rel, Problems: added}, fmt.Errorf("putting %s back failed; restore it with git: %w", rel, err))
		}
		return Entry{}, false, &RefusedError{Path: rel, Problems: added}
	}
	got.Path = rel
	return got, true, nil
}

// fieldLine is the file line of a front matter key in the entry at rel, or
// 1 when the file or the key cannot be read.
func fieldLine(repo *workitem.Repo, rel, key string) int {
	data, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(rel)))
	if err != nil {
		return 1
	}
	fm, _, err := workitem.SplitFrontMatter(string(data))
	if err != nil {
		return 1
	}
	return keyLine(fm, key)
}
