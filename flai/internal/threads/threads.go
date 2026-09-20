// Package threads manages wip/threads: conversations between the designer
// and agents anchored to a document, a heading, or a work item. Files are
// the durable channel (ADR-0020); the dashboard renders them and flai mcp
// serves them.
package threads

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the fixed subfolder under layout.wip.
const Folder = "threads"

// Statuses is the closed list: open awaits the other party, answered means
// the addressed party replied, resolved is closed.
var Statuses = []string{"open", "answered", "resolved"}

// Anchor is what a thread is about: a repository path, optionally a heading
// in it, and the work item when the path is an item file.
type Anchor struct {
	Path    string `yaml:"path" json:"path"`
	Heading string `yaml:"heading,omitempty" json:"heading,omitempty"`
	Item    string `yaml:"item,omitempty" json:"item,omitempty"`
}

// Thread is one file.
type Thread struct {
	ID           string   `yaml:"id" json:"id"`
	Title        string   `yaml:"title" json:"title"`
	Anchor       Anchor   `yaml:"anchor" json:"anchor"`
	Status       string   `yaml:"status" json:"status"`
	Participants []string `yaml:"participants" json:"participants"`
	Created      string   `yaml:"created" json:"created"`
	Updated      string   `yaml:"updated" json:"updated"`

	Path string `yaml:"-" json:"path"`
	Body string `yaml:"-" json:"-"`
}

// Entry is one dated contribution parsed from the body.
type Entry struct {
	At     string `json:"at"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

var (
	idPattern    = regexp.MustCompile(`^TH-\d{4,}$`)
	looseID      = regexp.MustCompile(`^(?i:th)-?0*(\d+)$`)
	entryHeading = regexp.MustCompile(`(?m)^### (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z) (.+)$`)
)

// Dir is <layout.wip>/threads, always in the main checkout.
func Dir(r *workitem.Repo) string { return filepath.Join(r.WipDir(), Folder) }

// CanonicalID pads a thread ID typed in any form (th-7, TH-007) to TH-0007.
func CanonicalID(id string) string {
	m := looseID.FindStringSubmatch(strings.TrimSpace(id))
	if m == nil {
		return id
	}
	return fmt.Sprintf("TH-%0*s", workitem.IDWidth, m[1])
}

// Parse decodes a thread document.
func Parse(doc string) (*Thread, error) {
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var th Thread
	if err := yaml.UnmarshalWithOptions([]byte(fm), &th, yaml.Strict()); err != nil {
		return nil, err
	}
	th.Body = body
	return &th, nil
}

// Read loads one file.
func Read(path string) (*Thread, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	th, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	th.Path = path
	return th, nil
}

// Validate checks the schema of one thread.
func (th *Thread) Validate() error {
	var errs []string
	if !idPattern.MatchString(th.ID) {
		errs = append(errs, fmt.Sprintf("id %q must look like TH-0001", th.ID))
	}
	if strings.TrimSpace(th.Title) == "" {
		errs = append(errs, "title is required")
	}
	if th.Anchor.Path == "" {
		errs = append(errs, "anchor.path is required")
	}
	if !contains(Statuses, th.Status) {
		errs = append(errs, fmt.Sprintf("status %q must be one of %s", th.Status, strings.Join(Statuses, ", ")))
	}
	if len(th.Participants) == 0 {
		errs = append(errs, "participants must name at least the opener")
	}
	for name, v := range map[string]string{"created": th.Created, "updated": th.Updated} {
		if _, err := time.Parse(workitem.TimeFormat, v); err != nil {
			errs = append(errs, fmt.Sprintf("%s %q is not a UTC timestamp", name, v))
		}
	}
	if len(errs) > 0 {
		sort.Strings(errs)
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

// Marshal renders the document in the same hand-written style as items.
func (th *Thread) Marshal() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "id: %s\n", th.ID)
	fmt.Fprintf(&b, "title: %s\n", workitem.Scalar(th.Title))
	b.WriteString("anchor:\n")
	fmt.Fprintf(&b, "  path: %s\n", workitem.Scalar(th.Anchor.Path))
	if th.Anchor.Heading != "" {
		fmt.Fprintf(&b, "  heading: %s\n", workitem.Scalar(th.Anchor.Heading))
	}
	if th.Anchor.Item != "" {
		fmt.Fprintf(&b, "  item: %s\n", th.Anchor.Item)
	}
	fmt.Fprintf(&b, "status: %s\n", th.Status)
	fmt.Fprintf(&b, "participants: %s\n", workitem.FlowList(th.Participants))
	fmt.Fprintf(&b, "created: %s\n", th.Created)
	fmt.Fprintf(&b, "updated: %s\n", th.Updated)
	b.WriteString("---\n")
	b.WriteString(th.Body)
	return b.String()
}

// Save writes the thread to its path.
func (th *Thread) Save() error {
	if err := os.MkdirAll(filepath.Dir(th.Path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(th.Path, []byte(th.Marshal()), 0o644)
}

// Entries parses the dated entries from the body, in order.
func (th *Thread) Entries() []Entry {
	locs := entryHeading.FindAllStringSubmatchIndex(th.Body, -1)
	var out []Entry
	for i, loc := range locs {
		end := len(th.Body)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		out = append(out, Entry{
			At:     th.Body[loc[2]:loc[3]],
			Author: th.Body[loc[4]:loc[5]],
			Text:   strings.TrimSpace(th.Body[loc[1]:end]),
		})
	}
	return out
}

// Opener is the author of the first entry, or the first participant.
func (th *Thread) Opener() string {
	if e := th.Entries(); len(e) > 0 {
		return e[0].Author
	}
	if len(th.Participants) > 0 {
		return th.Participants[0]
	}
	return ""
}

// Open reports whether the thread still needs someone.
func (th *Thread) Open() bool { return th.Status != "resolved" }

// List loads every thread, sorted by ID.
func List(r *workitem.Repo) ([]*Thread, error) {
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "TH-*.md"))
	var out []*Thread
	for _, m := range matches {
		th, err := Read(m)
		if err != nil {
			return nil, err
		}
		out = append(out, th)
	}
	sort.Slice(out, func(i, j int) bool { return num(out[i].ID) < num(out[j].ID) })
	return out, nil
}

// Get finds one thread by ID in any padding.
func Get(r *workitem.Repo, id string) (*Thread, error) {
	canon := CanonicalID(id)
	for _, cand := range []string{id, canon} {
		matches, _ := filepath.Glob(filepath.Join(Dir(r), cand+"-*.md"))
		if len(matches) > 0 {
			return Read(matches[0])
		}
	}
	return nil, fmt.Errorf("%s not found", id)
}

// For returns the threads anchored to a repository path or an item ID.
func For(r *workitem.Repo, pathOrItem string) ([]*Thread, error) {
	all, err := List(r)
	if err != nil {
		return nil, err
	}
	want := strings.TrimSuffix(filepath.ToSlash(pathOrItem), "/")
	// An item is matched in any padding, the anchor's included: S-4, S-004,
	// and S-0004 name one story, and older projects anchored with three digits.
	canon := workitem.CanonicalID(want)
	var out []*Thread
	for _, th := range all {
		if th.Anchor.Path == want || (th.Anchor.Item != "" && workitem.CanonicalID(th.Anchor.Item) == canon) {
			out = append(out, th)
		}
	}
	return out, nil
}

// NextID allocates the next TH-nnnn.
func NextID(r *workitem.Repo) string {
	max := 0
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "TH-*.md"))
	for _, m := range matches {
		if n := num(strings.SplitN(filepath.Base(m), "-", 3)[0] + "-" + strings.SplitN(filepath.Base(m), "-", 3)[1]); n > max {
			max = n
		}
	}
	return fmt.Sprintf("TH-%0*d", workitem.IDWidth, max+1)
}

// NewOptions describe a thread to open.
type NewOptions struct {
	Title   string
	On      string // repository path or item ID
	Heading string
	Author  string
	Text    string
	Now     time.Time
}

// New opens a thread. The anchor is validated: an item must exist, a path
// must exist under the repository, and a heading must be present in the file.
func New(r *workitem.Repo, opt NewOptions) (*Thread, error) {
	if strings.TrimSpace(opt.Title) == "" {
		return nil, fmt.Errorf("a title is required")
	}
	if strings.TrimSpace(opt.Author) == "" {
		return nil, fmt.Errorf("an author is required (--by, FLAI_AGENT, or config author)")
	}
	anchor, err := ResolveAnchor(r, opt.On, opt.Heading)
	if err != nil {
		return nil, err
	}
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	th := &Thread{
		ID: NextID(r), Title: strings.TrimSpace(opt.Title), Anchor: anchor, Status: "open",
		Participants: []string{opt.Author}, Created: now, Updated: now,
	}
	th.Path = filepath.Join(Dir(r), th.ID+"-"+orSlug(th.Title)+".md")
	where := anchor.Path
	if anchor.Heading != "" {
		where += " § " + anchor.Heading
	}
	th.Body = fmt.Sprintf("\n# %s %s\n\nOn %s.\n\n## Entries\n\n### %s %s\n%s\n", th.ID, th.Title, where, now, opt.Author, strings.TrimSpace(opt.Text))
	if err := th.Save(); err != nil {
		return nil, err
	}
	return th, nil
}

// ResolveAnchor turns a path or item ID (plus optional heading) into a
// validated anchor.
func ResolveAnchor(r *workitem.Repo, on, heading string) (Anchor, error) {
	on = strings.TrimSpace(on)
	if on == "" {
		return Anchor{}, fmt.Errorf("an anchor is required: --on <path> or --on <item ID>")
	}
	var a Anchor
	if workitem.TypeOfID(workitem.CanonicalID(on)) != "" && !strings.Contains(on, "/") && !strings.HasSuffix(on, ".md") {
		it, err := r.Get(on)
		if err != nil {
			return Anchor{}, err
		}
		a.Item = it.ID
		rel, err := filepath.Rel(r.MainRoot, it.Path)
		if err != nil {
			return Anchor{}, err
		}
		a.Path = filepath.ToSlash(rel)
	} else {
		rel := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(on)), "./")
		abs := filepath.Join(r.Root, filepath.FromSlash(rel))
		if _, err := os.Stat(abs); err != nil {
			if _, err2 := os.Stat(filepath.Join(r.MainRoot, filepath.FromSlash(rel))); err2 != nil {
				return Anchor{}, fmt.Errorf("anchor %s does not exist in the repository", rel)
			}
		}
		a.Path = rel
	}
	if heading != "" {
		abs := filepath.Join(r.Root, filepath.FromSlash(a.Path))
		data, err := os.ReadFile(abs)
		if err != nil {
			data, err = os.ReadFile(filepath.Join(r.MainRoot, filepath.FromSlash(a.Path)))
		}
		if err != nil {
			return Anchor{}, err
		}
		if !HasHeading(string(data), heading) {
			return Anchor{}, fmt.Errorf("heading %q is not in %s", heading, a.Path)
		}
		a.Heading = strings.TrimSpace(heading)
	}
	return a, nil
}

// HasHeading reports whether a markdown document has the heading (any level,
// case-insensitive, ignoring surrounding whitespace).
func HasHeading(doc, heading string) bool {
	want := strings.ToLower(strings.TrimSpace(heading))
	for _, line := range strings.Split(doc, "\n") {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		text := strings.ToLower(strings.TrimSpace(strings.TrimLeft(line, "#")))
		if text == want {
			return true
		}
	}
	return false
}

// Reply appends an entry. The status becomes answered when someone other
// than the opener replies and open when the opener follows up, so a
// resolved thread reopens on any reply.
func Reply(r *workitem.Repo, id, author, text string, now time.Time) (*Thread, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("a reply needs text")
	}
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("an author is required")
	}
	th, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	stamp := now.UTC().Format(workitem.TimeFormat)
	th.Body = strings.TrimRight(th.Body, "\n") + fmt.Sprintf("\n\n### %s %s\n%s\n", stamp, author, strings.TrimSpace(text))
	if author == th.Opener() {
		th.Status = "open"
	} else {
		th.Status = "answered"
	}
	if !contains(th.Participants, author) {
		th.Participants = append(th.Participants, author)
	}
	th.Updated = stamp
	return th, th.Save()
}

// Resolve closes a thread, recording who did it and why.
func Resolve(r *workitem.Repo, id, author, reason string, now time.Time) (*Thread, error) {
	th, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	stamp := now.UTC().Format(workitem.TimeFormat)
	note := "Resolved."
	if strings.TrimSpace(reason) != "" {
		note = "Resolved: " + strings.TrimSpace(reason)
	}
	th.Body = strings.TrimRight(th.Body, "\n") + fmt.Sprintf("\n\n### %s %s\n%s\n", stamp, orDefault(author, "unknown"), note)
	th.Status = "resolved"
	if author != "" && !contains(th.Participants, author) {
		th.Participants = append(th.Participants, author)
	}
	th.Updated = stamp
	return th, th.Save()
}

// StoryOf returns the story a thread belongs to for the narrative mirror:
// the anchored story, or the parent of an anchored task.
func StoryOf(r *workitem.Repo, th *Thread) string {
	if th.Anchor.Item == "" {
		return ""
	}
	switch workitem.TypeOfID(th.Anchor.Item) {
	case workitem.Story:
		return th.Anchor.Item
	case workitem.Task:
		if it, err := r.Get(th.Anchor.Item); err == nil {
			return it.Parent
		}
	}
	return ""
}

const (
	mirrorStart = "<!-- threads:start -->"
	mirrorEnd   = "<!-- threads:end -->"
)

// MirrorNarrative rewrites the generated block under "## Open questions" in
// the story's narrative to list its unresolved threads, so a fresh agent
// sees them without the dashboard. Missing narrative: nothing to do.
func MirrorNarrative(r *workitem.Repo, storyID string) error {
	path := r.NarrativePath(storyID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil // no narrative yet, nothing to mirror into
	}
	if err != nil {
		return err
	}
	all, err := List(r)
	if err != nil {
		return err
	}
	var lines []string
	for _, th := range all {
		if th.Open() && StoryOf(r, th) == storyID {
			lines = append(lines, fmt.Sprintf("- %s (%s, %s, %s): %s", th.ID, th.Status, th.Opener(), th.Updated[:10], th.Title))
		}
	}
	doc := string(data)
	block := ""
	if len(lines) > 0 {
		block = mirrorStart + "\nThreads on this story (flai thread show <id>):\n" + strings.Join(lines, "\n") + "\n" + mirrorEnd + "\n"
	}
	if i := strings.Index(doc, mirrorStart); i >= 0 {
		j := strings.Index(doc, mirrorEnd)
		if j < i {
			return fmt.Errorf("%s: malformed threads block", path)
		}
		doc = doc[:i] + block + strings.TrimPrefix(doc[j+len(mirrorEnd):], "\n")
	} else if block != "" {
		h := "## Open questions"
		i := strings.Index(doc, "\n"+h)
		if i < 0 {
			return nil
		}
		eol := strings.Index(doc[i+1:], "\n")
		if eol < 0 {
			doc += "\n" + block
		} else {
			at := i + 1 + eol + 1
			doc = doc[:at] + block + doc[at:]
		}
	}
	if doc == string(data) {
		return nil
	}
	return os.WriteFile(path, []byte(doc), 0o644)
}

func num(id string) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, "TH-"))
	return n
}

func orSlug(title string) string {
	if s := template.Slug(title); s != "" {
		return s
	}
	return "thread"
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// View is a thread as flai prints it and the dashboard reads it: the front
// matter, the path relative to the repository, the entries, and the story
// the thread belongs to.
func View(r *workitem.Repo, th *Thread) map[string]any {
	path := th.Path
	if rel, err := filepath.Rel(r.MainRoot, th.Path); err == nil {
		path = filepath.ToSlash(rel)
	}
	return map[string]any{
		"id": th.ID, "title": th.Title, "anchor": th.Anchor, "status": th.Status,
		"participants": th.Participants, "created": th.Created, "updated": th.Updated,
		"path": path, "entries": th.Entries(), "story": StoryOf(r, th),
	}
}
