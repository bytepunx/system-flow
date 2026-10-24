// Package digest knows the digest (ADR-0191): one markdown file per fact a
// story learned, D-nnnn-<slug>.md in the design folder's digest folder, whose
// front matter flai writes and checks. Nothing in it is ever deleted: an
// entry is absorbed or compacted by a change of status.
package digest

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the digest folder's name inside the design folder.
const Folder = "digest"

// Index is the folder's README.md, which says what the folder is for and is
// not an entry.
const Index = "README.md"

// Dir is the digest folder, relative to the project root, slash separated.
func Dir(m manifest.Manifest) string {
	return path.Join(m.Layout["design"], Folder)
}

// Kinds are what an entry records: a decision, the operator's answer on a
// thread, a recommendation the operator confirmed or overruled, or an
// issue's cause and fix.
var Kinds = []string{"decision", "answer", "recommendation", "cause"}

// Statuses are an entry's states: active is primed, absorbed is recorded by
// a document now and no longer primed, compacted is primed by its title only.
var Statuses = []string{"active", "absorbed", "compacted"}

// Rules are the check rules a problem is reported under.
const (
	RuleEntry  = "digest.entry"
	RuleSecret = "digest.secret"
)

// Entry is one digest entry: its front matter, where it is, and its body.
type Entry struct {
	ID         string   `yaml:"id" json:"id"`
	Title      string   `yaml:"title" json:"title"`
	Kind       string   `yaml:"kind" json:"kind"`
	Sources    []string `yaml:"sources" json:"sources"`
	Topics     []string `yaml:"topics" json:"topics"`
	Paths      []string `yaml:"paths" json:"paths"`
	Date       string   `yaml:"date" json:"date"`
	Status     string   `yaml:"status" json:"status"`
	AbsorbedBy string   `yaml:"absorbed_by" json:"absorbed_by"`
	AskedAgain []string `yaml:"asked_again" json:"asked_again"`
	Cited      []string `yaml:"cited" json:"cited"`

	// Path is the file's path from the project root, slash separated.
	Path string `yaml:"-" json:"path"`
	// Body is the text below the entry's "# D-nnnn <title>" heading.
	Body string `yaml:"-" json:"body"`
}

// Problem is one way an entry falls short. Field is the front matter key it
// is about, or empty for the file; Line is the 1-based line in the file of
// that key, or of the text it is about, else 1.
type Problem struct {
	Field   string `json:"field,omitempty"`
	Line    int    `json:"line"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// String is the problem as "line N: rule: message".
func (p Problem) String() string {
	return fmt.Sprintf("line %d: %s: %s", p.Line, p.Rule, p.Message)
}

// Env answers what Validate cannot see in an entry's text. A nil func skips
// the checks that need it.
type Env struct {
	// Known reports whether a work item (E-, S-, T-), thread (TH-), or issue
	// (I-) with the ID exists in the project or its archive.
	Known func(id string) bool
	// Exists reports whether a path from the project root exists.
	Exists func(path string) bool
}

var (
	fileName  = regexp.MustCompile(`^D-(\d{4,})-([a-z0-9]+(?:-[a-z0-9]+)*)\.md$`)
	entryFile = regexp.MustCompile(`^D-(\d+)-.*\.md$`)
	idShape   = regexp.MustCompile(`^D-\d{4,}$`)
	looseID   = regexp.MustCompile(`^(?i:d)-?0*(\d+)$`)
	sourceID  = regexp.MustCompile(`^(?:[ESTI]|TH)-\d{3,}$`)
	citedID   = regexp.MustCompile(`^(?:[EST]|TH)-\d{3,}$`)
	threadID  = regexp.MustCompile(`^TH-\d{3,}$`)
	errAt     = regexp.MustCompile(`^\[(\d+):\d+\]`)
	// sourceFile is the ID a work item, thread, or issue file is named for.
	sourceFile = regexp.MustCompile(`^((?:[ESTI]|TH)-\d+)(?:-.*)?\.md$`)
)

// Skipped reports whether a file in the digest folder is not an entry: the
// folder's index.
func Skipped(name string) bool { return name == Index }

// CanonicalID writes an entry ID typed in any form (d-7, D-007) as D-0007;
// anything else is returned unchanged.
func CanonicalID(id string) string {
	m := looseID.FindStringSubmatch(strings.TrimSpace(id))
	if m == nil {
		return id
	}
	return fmt.Sprintf("D-%0*s", workitem.IDWidth, m[1])
}

// number is an entry ID's or file name's number, 0 when it has none.
func number(s string) int {
	m := entryFile.FindStringSubmatch(s)
	if m == nil {
		m = looseID.FindStringSubmatch(s)
	}
	if m == nil {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

// List reads every entry in the checkout's digest folder, sorted by ID. A
// missing folder is no entries.
func List(repo *workitem.Repo) ([]Entry, error) {
	dir := filepath.Join(repo.Root, filepath.FromSlash(Dir(repo.Manifest)))
	files, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading the digest folder %s: %w", Dir(repo.Manifest), err)
	}
	var out []Entry
	for _, f := range files {
		if f.IsDir() || !entryFile.MatchString(f.Name()) {
			continue
		}
		e, err := read(repo, f.Name())
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	slices.SortFunc(out, func(a, b Entry) int { return number(path.Base(a.Path)) - number(path.Base(b.Path)) })
	return out, nil
}

// Get reads the entry with the ID, given in any zero padding (D-7 is D-0007).
func Get(repo *workitem.Repo, id string) (Entry, error) {
	n := number(CanonicalID(id))
	if n == 0 {
		return Entry{}, fmt.Errorf("%q is not a digest entry ID; an entry's ID is D-nnnn, such as D-0007", id)
	}
	dir := filepath.Join(repo.Root, filepath.FromSlash(Dir(repo.Manifest)))
	files, err := os.ReadDir(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Entry{}, fmt.Errorf("reading the digest folder %s: %w", Dir(repo.Manifest), err)
	}
	for _, f := range files {
		if !f.IsDir() && entryFile.MatchString(f.Name()) && number(f.Name()) == n {
			return read(repo, f.Name())
		}
	}
	return Entry{}, fmt.Errorf("%s is not in %s; flai digest list lists the entries", CanonicalID(id), Dir(repo.Manifest))
}

// read parses the entry file named name in the checkout's digest folder.
func read(repo *workitem.Repo, name string) (Entry, error) {
	rel := path.Join(Dir(repo.Manifest), name)
	data, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(rel)))
	if err != nil {
		return Entry{}, fmt.Errorf("reading digest entry %s: %w", rel, err)
	}
	e, err := parse(string(data))
	if err != nil {
		return Entry{}, fmt.Errorf("digest entry %s: %w; flai check names what to fix", rel, err)
	}
	e.Path = rel
	return e, nil
}

// parse decodes an entry's front matter and body, with every list empty
// rather than nil.
func parse(content string) (Entry, error) {
	var e Entry
	fm, body, err := workitem.SplitFrontMatter(content)
	if err != nil {
		return e, err
	}
	if err := yaml.Unmarshal([]byte(fm), &e); err != nil {
		return e, fmt.Errorf("front matter does not parse: %s", firstLine(err.Error()))
	}
	_, e.Body = splitHeading(body)
	e.fill()
	return e, nil
}

// fill makes every nil list empty, so an entry serialises its lists as [].
func (e *Entry) fill() {
	for _, l := range []*[]string{&e.Sources, &e.Topics, &e.Paths, &e.AskedAgain, &e.Cited} {
		if *l == nil {
			*l = []string{}
		}
	}
}

// splitHeading cuts a body after its first "# " heading line: head is the
// text up to and including that line, rest what follows it without leading
// blank lines or trailing space. A body with no heading is all rest.
func splitHeading(body string) (head, rest string) {
	at := 0
	for _, l := range strings.SplitAfter(body, "\n") {
		at += len(l)
		if strings.HasPrefix(l, "# ") {
			return body[:at], strings.TrimSpace(body[at:])
		}
	}
	return "", strings.TrimSpace(body)
}

func firstLine(s string) string {
	l, _, _ := strings.Cut(s, "\n")
	return strings.TrimSpace(l)
}

// keyLine is the file line of a front matter key, the front matter starting
// on the file's second line, or 1 when the key is not there.
func keyLine(fm, key string) int {
	for i, l := range strings.Split(fm, "\n") {
		if strings.HasPrefix(l, key+":") {
			return i + 2
		}
	}
	return 1
}

// Validate checks an entry file named name with content: its name, its
// front matter against ADR-0191, what env can see of its sources and
// absorbed_by, and that no text in it has the shape of a secret.
func Validate(name, content string, env Env) (Entry, []Problem) {
	var out []Problem
	entry := func(field string, line int, format string, args ...any) {
		out = append(out, Problem{Field: field, Line: line, Rule: RuleEntry, Message: fmt.Sprintf(format, args...)})
	}
	var secrets []Problem
	for _, s := range Secrets(content) {
		secrets = append(secrets, Problem{Line: s.Line, Rule: RuleSecret, Message: s.Description + "; an entry says where a secret lives, never its value: remove the value and name the file or variable that holds it"})
	}
	m := fileName.FindStringSubmatch(name)
	if m == nil {
		entry("", 1, "the file is named %s; an entry is named D-nnnn-<slug>.md, its ID and its title's slug in lowercase words joined by dashes, as flai digest new names it", name)
	}
	fm, body, err := workitem.SplitFrontMatter(content)
	if err != nil {
		entry("", 1, "no front matter; an entry opens with --- and its id, title, kind, sources, date, and status, as flai digest new writes it")
		return Entry{}, append(out, secrets...)
	}
	var e Entry
	if err := yaml.Unmarshal([]byte(fm), &e); err != nil {
		line := 1
		if at := errAt.FindStringSubmatch(err.Error()); at != nil {
			n, _ := strconv.Atoi(at[1])
			line = n + 1
		}
		entry("", line, "front matter does not parse: %s; fix the YAML so each key holds the value ADR-0191 gives it", firstLine(err.Error()))
		return e, append(out, secrets...)
	}
	_, e.Body = splitHeading(body)
	e.fill()
	at := func(key string) int { return keyLine(fm, key) }

	for _, k := range []struct{ name, value string }{{"id", e.ID}, {"title", e.Title}, {"kind", e.Kind}, {"date", e.Date}, {"status", e.Status}} {
		if strings.TrimSpace(k.value) == "" {
			entry(k.name, at(k.name), "front matter needs %s; add it", k.name)
		}
	}
	switch {
	case e.ID == "":
	case !idShape.MatchString(e.ID):
		entry("id", at("id"), "id is %q, not an entry ID, D-nnnn; give it the number the file is named for", e.ID)
	case m != nil && e.ID != "D-"+m[1]:
		entry("id", at("id"), "id is %s but the file is named for D-%s; make them the same", e.ID, m[1])
	}
	if e.Kind != "" && !slices.Contains(Kinds, e.Kind) {
		entry("kind", at("kind"), "kind is %q, which is none of %s", e.Kind, strings.Join(Kinds, ", "))
	}
	if e.Status != "" && !slices.Contains(Statuses, e.Status) {
		entry("status", at("status"), "status is %q, which is none of %s", e.Status, strings.Join(Statuses, ", "))
	}
	if e.Date != "" {
		if _, err := time.Parse(time.DateOnly, e.Date); err != nil {
			entry("date", at("date"), "date is %q, not a date, YYYY-MM-DD; give the day the fact was learned", e.Date)
		}
	}
	if len(e.Sources) == 0 {
		entry("sources", at("sources"), "sources names nothing; give at least one work item (E-, S-, T-), thread (TH-), or issue (I-) the fact comes from")
	}
	for _, s := range e.Sources {
		switch {
		case !sourceID.MatchString(s):
			entry("sources", at("sources"), "source %q is not the ID of a work item (E-, S-, T-), thread (TH-), or issue (I-)", s)
		case env.Known != nil && !env.Known(s):
			entry("sources", at("sources"), "source %s names nothing in the project or its archive; give an ID that exists", s)
		}
	}
	for _, t := range e.Topics {
		if !workitem.ValidTopic(t) {
			entry("topics", at("topics"), "topic %q is not a word; a topic is letters, digits, dot, dash, or underscore, such as cli or go", t)
		}
	}
	// absorbed_by names a document, or a heading in it as <path>#<anchor>;
	// the path is what must be local and exist.
	absorbedFile, _, _ := strings.Cut(e.AbsorbedBy, "#")
	switch {
	case e.AbsorbedBy != "" && e.Status != "absorbed":
		entry("absorbed_by", at("absorbed_by"), "absorbed_by is set but status is %s; set status: absorbed, or empty absorbed_by", e.Status)
	case e.AbsorbedBy == "" && e.Status == "absorbed":
		entry("absorbed_by", at("absorbed_by"), "status is absorbed but absorbed_by is empty; give the path, from the project root, of the document that now records the fact")
	case e.AbsorbedBy != "" && !filepath.IsLocal(filepath.FromSlash(absorbedFile)):
		entry("absorbed_by", at("absorbed_by"), "absorbed_by is %q, not a path from the project root; give one such as design/system/overview.md", e.AbsorbedBy)
	case e.AbsorbedBy != "" && env.Exists != nil && !env.Exists(absorbedFile):
		entry("absorbed_by", at("absorbed_by"), "absorbed_by is %s, which names no file; give the path of the document that records the fact", e.AbsorbedBy)
	}
	for _, id := range e.AskedAgain {
		if !threadID.MatchString(id) {
			entry("asked_again", at("asked_again"), "asked_again holds %q, not a thread ID, TH-nnnn", id)
		}
	}
	for _, id := range e.Cited {
		if !citedID.MatchString(id) {
			entry("cited", at("cited"), "cited holds %q, not the ID of a work item (E-, S-, T-) or thread (TH-)", id)
		}
	}
	return e, append(out, secrets...)
}

// EnvFor is the Env of a project: Known reads the file names of its work
// items, kanban and archive, its threads, and its issues, once, on first
// use; Exists looks under the checkout's root.
func EnvFor(repo *workitem.Repo) Env {
	var once sync.Once
	ids := map[string]bool{}
	load := func() {
		var dirs []string
		for _, typ := range workitem.Types {
			dirs = append(dirs, repo.ItemDir(typ, false), repo.ItemDir(typ, true))
		}
		dirs = append(dirs, threads.Dir(repo), issues.Dir(repo))
		for _, dir := range dirs {
			files, _ := os.ReadDir(dir)
			for _, f := range files {
				if m := sourceFile.FindStringSubmatch(f.Name()); m != nil && !f.IsDir() {
					ids[canonicalSource(m[1])] = true
				}
			}
		}
	}
	return Env{
		Known: func(id string) bool {
			once.Do(load)
			return ids[canonicalSource(id)]
		},
		Exists: func(p string) bool {
			local := filepath.FromSlash(p)
			if !filepath.IsLocal(local) {
				return false
			}
			_, err := os.Stat(filepath.Join(repo.Root, local))
			return err == nil
		},
	}
}

// canonicalSource pads a work item, issue, or thread ID typed in any form.
func canonicalSource(id string) string {
	if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(id)), "TH") {
		return threads.CanonicalID(id)
	}
	return workitem.CanonicalID(id)
}
