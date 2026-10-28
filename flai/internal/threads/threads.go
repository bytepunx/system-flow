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

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
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

	// Unknown is the front matter this flai does not know, kept for writing
	// back (S-0181).
	Unknown []workitem.Field `yaml:"-" json:"-"`

	Path string `yaml:"-" json:"path"`
	Body string `yaml:"-" json:"-"`
}

// Entry is one dated contribution parsed from the body. Text is the entry as
// written, its source line included. Recommendation and Source are the marks
// of ADR-0090: a recommendation awaits the operator's confirmation, and the
// source is what the entry cites. OperatorOnly is ADR-0206's mark: the
// recommendation is the operator's alone to confirm, so flai serve never
// confirms it after orchestration.permissions.confirm_after.
type Entry struct {
	At             string  `json:"at"`
	Author         string  `json:"author"`
	Text           string  `json:"text"`
	Recommendation bool    `json:"recommendation"`
	OperatorOnly   bool    `json:"operator_only,omitempty"`
	Source         *Source `json:"source,omitempty"`
}

// Source is what an entry cites: a repository path, such as an ADR, a design
// document, or a convention, and optionally a heading in it.
type Source struct {
	Path    string `json:"path"`
	Heading string `json:"heading,omitempty"`
}

// ParseSource reads a reply's source as flai thread reply --source and the
// MCP tool thread_reply take it: a repository path or item ID, optionally
// followed by # and a heading in it. An empty value cites nothing. It checks
// the form only; ReplyWith checks that the path and heading exist.
func ParseSource(v string) (Source, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return Source{}, nil
	}
	path, heading, _ := strings.Cut(v, "#")
	if path = strings.TrimSpace(path); path == "" {
		return Source{}, fmt.Errorf("source %q names no file: give <path> or <path>#<heading>, such as design/system/overview.md#Delivery sequence", v)
	}
	return Source{Path: path, Heading: strings.TrimSpace(heading)}, nil
}

// Marks are the optional marks of a reply: that it is a recommendation, that
// the recommendation is the operator's alone (ADR-0206), and the source it
// cites. The zero value is a plain reply.
type Marks struct {
	Recommendation bool
	// OperatorOnly marks a recommendation on money or a host change as the
	// operator's alone to confirm; it needs Recommendation.
	OperatorOnly bool
	// Source is a repository path or item ID, validated as an anchor is; an
	// empty Path cites nothing.
	Source Source
}

// TimeoutAuthor is the author of the confirmation ConfirmByTimeout writes
// (ADR-0206).
const TimeoutAuthor = "flai serve"

const (
	// recommendationMark ends the heading of a recommendation entry.
	recommendationMark = " (recommendation)"
	// sourcePrefix starts an entry's last paragraph when it cites a source.
	sourcePrefix = "Source: "
	// sourceHeading separates a source's path from its heading.
	sourceHeading = " § "
	// confirmPrefix starts the entry Confirm and ConfirmByTimeout write.
	confirmPrefix = "Confirmed the recommendation of "
	// operatorOnlyPrefix starts the paragraph that marks a recommendation for
	// the operator alone (ADR-0206).
	operatorOnlyPrefix = "For the operator alone: "
	// operatorOnlyLine is that paragraph, written above the source line.
	operatorOnlyLine = operatorOnlyPrefix + "flai serve does not confirm this after orchestration.permissions.confirm_after."
)

var (
	idPattern    = regexp.MustCompile(`^TH-\d{4,}$`)
	storyAgent   = regexp.MustCompile(`^agent-S-\d+$`)
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

// Parse decodes a thread document, keeping front-matter fields this flai
// does not know in Unknown (S-0181).
func Parse(doc string) (*Thread, error) {
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var th Thread
	if err := yaml.Unmarshal([]byte(fm), &th); err != nil {
		return nil, err
	}
	th.Unknown = workitem.UnknownFields(fm, th)
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
	workitem.WriteFields(&b, th.Unknown)
	b.WriteString("---\n")
	b.WriteString(th.Body)
	return b.String()
}

// Save writes the thread to its path.
func (th *Thread) Save() error {
	if err := os.MkdirAll(filepath.Dir(th.Path), 0o755); err != nil {
		return err
	}
	return atomicfile.WriteFile(th.Path, []byte(th.Marshal()), 0o644)
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
		e := Entry{At: th.Body[loc[2]:loc[3]], Text: strings.TrimSpace(th.Body[loc[1]:end])}
		e.Author, e.Recommendation = strings.CutSuffix(th.Body[loc[4]:loc[5]], recommendationMark)
		e.OperatorOnly = e.Recommendation && hasLine(e.Text, operatorOnlyPrefix)
		e.Source = sourceOf(e.Text)
		out = append(out, e)
	}
	return out
}

// hasLine reports whether a line of text starts with prefix.
func hasLine(text, prefix string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

// sourceOf reads the source an entry's text cites in its last paragraph,
// "Source: <path>" or "Source: <path> § <heading>", or nil.
func sourceOf(text string) *Source {
	i := strings.LastIndex(text, "\n\n")
	if i < 0 {
		return nil
	}
	last := text[i+2:]
	if !strings.HasPrefix(last, sourcePrefix) || strings.Contains(last, "\n") {
		return nil
	}
	path, heading, _ := strings.Cut(strings.TrimPrefix(last, sourcePrefix), sourceHeading)
	if path = strings.TrimSpace(path); path == "" {
		return nil
	}
	return &Source{Path: path, Heading: strings.TrimSpace(heading)}
}

// sourceLine renders a source as an entry's last paragraph.
func sourceLine(s Source) string {
	if s.Heading == "" {
		return sourcePrefix + s.Path
	}
	return sourcePrefix + s.Path + sourceHeading + s.Heading
}

// Confirms reports whether the entry is a confirmation Confirm or
// ConfirmByTimeout wrote, alone or joined with another entry of its author's
// in the same second.
func (e Entry) Confirms() bool { return hasLine(e.Text, confirmPrefix) }

// ConfirmsByTimeout reports whether the entry is the confirmation flai serve
// wrote with ConfirmByTimeout in the operator's absence (ADR-0206).
func (e Entry) ConfirmsByTimeout() bool { return e.Author == TimeoutAuthor && e.Confirms() }

// WroteSince reports whether an entry after the given one, in the thread's
// order, is by an author that by accepts, such as OperatorName for the
// thread's story. It is false when the thread has no entry of after's time
// and author.
func (th *Thread) WroteSince(after Entry, by func(author string) bool) bool {
	seen := false
	for _, e := range th.Entries() {
		switch {
		case seen && by(e.Author):
			return true
		case e.At == after.At && e.Author == after.Author:
			seen = true
		}
	}
	return false
}

// OpenedByAgent reports whether the thread's opener is a story's agent
// (agent-S-nnnn) or a planner (planner-…), the openers whose threads flai
// serve may confirm a recommendation on (ADR-0206).
func (th *Thread) OpenedByAgent() bool {
	opener := th.Opener()
	return storyAgent.MatchString(opener) || workitem.IsPlanner("", opener)
}

// PendingRecommendation is the recommendation awaiting the operator's
// confirmation, or nil: the newest entry by someone other than the opener,
// when it is a recommendation, no entry after it confirms it, and the thread
// is not resolved.
func (th *Thread) PendingRecommendation() *Entry {
	if !th.Open() {
		return nil
	}
	es := th.Entries()
	if len(es) == 0 {
		return nil
	}
	opener := es[0].Author
	for i := len(es) - 1; i > 0; i-- {
		e := es[i]
		if e.Confirms() {
			return nil
		}
		if e.Author == opener {
			continue
		}
		if e.Recommendation {
			return &e
		}
		return nil
	}
	return nil
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

// OrchestratorsOwn reports whether the thread is the orchestrator's alone:
// the orchestrator opened it and wrote every entry on it. One entry by
// anyone else, the operator or an agent, makes it a thread others take part
// in. An open thread of the orchestrator's own still holds a story; the
// orchestrator's acceptance names it apart, as the orchestrator's to resolve
// (ADR-0189).
func (th *Thread) OrchestratorsOwn() bool {
	if !workitem.IsOrchestrator(th.Opener()) {
		return false
	}
	for _, e := range th.Entries() {
		if !workitem.IsOrchestrator(e.Author) {
			return false
		}
	}
	return true
}

// List loads every thread, sorted by ID.
func List(r *workitem.Repo) ([]*Thread, error) {
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "TH-*.md"))
	var out []*Thread
	for _, m := range matches {
		th, err := Read(m)
		if err != nil {
			return nil, err
		}
		workitem.WarnUnknown(m, th.Unknown)
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
	if opt.Title = workitem.CleanTitle(opt.Title); opt.Title == "" {
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
		ID: NextID(r), Title: opt.Title, Anchor: anchor, Status: "open",
		Participants: []string{opt.Author}, Created: now, Updated: now,
	}
	th.Path = filepath.Join(Dir(r), th.ID+"-"+orSlug(th.Title)+".md")
	where := anchor.Path
	if anchor.Heading != "" {
		where += " § " + anchor.Heading
	}
	th.Body = fmt.Sprintf("\n# %s %s\n\nOn %s.\n\n## Entries\n\n### %s %s\n%s\n", th.ID, th.Title, where, now, opt.Author, strings.TrimSpace(opt.Text))
	if err := r.LintGuard(th.Path, "", th.Marshal()); err != nil {
		return nil, err
	}
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

// Reply appends a plain entry. The status becomes answered when someone
// other than the opener replies and open when the opener follows up, so a
// resolved thread reopens on any reply.
func Reply(r *workitem.Repo, id, author, text string, now time.Time) (*Thread, error) {
	return ReplyWith(r, id, author, text, now, Marks{})
}

// ReplyWith appends an entry with marks (ADR-0090). A plain reply, with or
// without a source, sets the status as Reply does. A recommendation leaves
// the status as it was, so the thread still awaits the operator, who makes it
// the answer with Confirm; it is refused on a resolved thread and from the
// thread's opener. The operator-alone mark, which needs a recommendation, is
// written as its own paragraph above the source line (ADR-0206). A source
// must exist in the repository, and its heading in it, as an anchor must.
func ReplyWith(r *workitem.Repo, id, author, text string, now time.Time, marks Marks) (*Thread, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("a reply needs text")
	}
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("an author is required")
	}
	if marks.OperatorOnly && !marks.Recommendation {
		return nil, fmt.Errorf("only a recommendation is marked for the operator alone: mark the reply a recommendation too, or drop the operator-alone mark")
	}
	th, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(text)
	if marks.OperatorOnly {
		text += "\n\n" + operatorOnlyLine
	}
	if marks.Source.Path != "" {
		a, err := ResolveAnchor(r, marks.Source.Path, marks.Source.Heading)
		if err != nil {
			return nil, fmt.Errorf("the source of a reply on %s: %w", th.ID, err)
		}
		marks.Source = Source{Path: a.Path, Heading: a.Heading}
		text += "\n\n" + sourceLine(marks.Source)
	}
	heading := author
	if marks.Recommendation {
		if !th.Open() {
			return nil, fmt.Errorf("%s is resolved: a recommendation needs a thread that awaits an answer", th.ID)
		}
		if author == th.Opener() {
			return nil, fmt.Errorf("%s opened %s: a recommendation answers someone else's thread", author, th.ID)
		}
		heading += recommendationMark
	}
	was := th.Marshal()
	stamp := now.UTC().Format(workitem.TimeFormat)
	th.Body = appendEntry(th.Body, stamp, heading, text)
	switch {
	case marks.Recommendation:
		// the thread still awaits the operator: the status stays as it was
	case author == th.Opener():
		th.Status = "open"
	default:
		th.Status = "answered"
	}
	if !contains(th.Participants, author) {
		th.Participants = append(th.Participants, author)
	}
	th.Updated = stamp
	return th, save(r, th, was)
}

// Confirm makes the pending recommendation the answer (ADR-0090): it appends
// an entry by author naming the recommendation's time and author, citing its
// source, and sets the status to answered. It is refused when no
// recommendation is pending and to the recommendation's own author.
func Confirm(r *workitem.Repo, id, author string, now time.Time) (*Thread, error) {
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("an author is required")
	}
	return confirm(r, id, author, now, func(th *Thread, p *Entry) (string, error) {
		if p.Author == author {
			return "", fmt.Errorf("%s made the recommendation on %s and cannot confirm it: the operator confirms it", author, th.ID)
		}
		return "", nil
	})
}

// ConfirmByTimeout is flai serve's confirmation of the pending
// recommendation in the operator's absence (ADR-0206): an entry by
// TimeoutAuthor, written as Confirm writes the operator's and saying it came
// after the given confirm_after with no answer from the operator, with the
// status answered. It is refused when no recommendation is pending, when it
// is marked for the operator alone, when it was written less than after
// before now, and when the operator, the manifest's owner or the story's
// owner, wrote on the thread since it. Which threads it is asked for is the
// caller's.
func ConfirmByTimeout(r *workitem.Repo, id string, after time.Duration, now time.Time) (*Thread, error) {
	if after <= 0 {
		return nil, fmt.Errorf("confirming by timeout needs orchestration.permissions.confirm_after longer than zero, not %s: unset, nothing is confirmed", after)
	}
	return confirm(r, id, TimeoutAuthor, now, func(th *Thread, p *Entry) (string, error) {
		if p.OperatorOnly {
			return "", fmt.Errorf("the recommendation of %s %s on %s is marked for the operator alone: only the operator confirms it", p.At, p.Author, th.ID)
		}
		at, err := time.Parse(workitem.TimeFormat, p.At)
		if err != nil {
			return "", fmt.Errorf("the recommendation on %s has no readable time %q: %w", th.ID, p.At, err)
		}
		if waited := now.Sub(at); waited < after {
			return "", fmt.Errorf("the recommendation of %s %s on %s has waited %s, less than confirm_after %s: it is not confirmed yet", p.At, p.Author, th.ID, readableDuration(waited.Truncate(time.Second)), readableDuration(after))
		}
		story := StoryOf(r, th)
		if th.WroteSince(*p, func(name string) bool { return OperatorName(r, story, name) }) {
			return "", fmt.Errorf("the operator wrote on %s after the recommendation of %s %s: it is theirs to confirm", th.ID, p.At, p.Author)
		}
		return " after " + readableDuration(after) + " with no answer from the operator (orchestration.permissions.confirm_after)", nil
	})
}

// readableDuration writes d as Go does without its zero trailing units:
// 48h0m0s is 48h, 1h30m0s is 1h30m.
func readableDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// confirm appends a confirmation of the pending recommendation by author:
// check refuses it with an error or returns what the confirmation says after
// the recommendation's time and author. The entry cites the recommendation's
// source, and the status becomes answered.
func confirm(r *workitem.Repo, id, author string, now time.Time, check func(th *Thread, p *Entry) (string, error)) (*Thread, error) {
	th, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	p := th.PendingRecommendation()
	if p == nil {
		return nil, fmt.Errorf("%s has no recommendation awaiting confirmation: confirm one that is the newest reply by someone other than the opener", th.ID)
	}
	more, err := check(th, p)
	if err != nil {
		return nil, err
	}
	text := fmt.Sprintf("%s%s %s%s.", confirmPrefix, p.At, p.Author, more)
	if p.Source != nil {
		text += "\n\n" + sourceLine(*p.Source)
	}
	was := th.Marshal()
	stamp := now.UTC().Format(workitem.TimeFormat)
	th.Body = appendEntry(th.Body, stamp, author, text)
	th.Status = "answered"
	if !contains(th.Participants, author) {
		th.Participants = append(th.Participants, author)
	}
	th.Updated = stamp
	return th, save(r, th, was)
}

// save writes the thread unless the change brings markdown the project's
// lint rejects (S-0179).
func save(r *workitem.Repo, th *Thread, was string) error {
	if err := r.LintGuard(th.Path, was, th.Marshal()); err != nil {
		return err
	}
	return th.Save()
}

// Resolve closes a thread, recording who did it and why.
func Resolve(r *workitem.Repo, id, author, reason string, now time.Time) (*Thread, error) {
	th, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	was := th.Marshal()
	stamp := now.UTC().Format(workitem.TimeFormat)
	note := "Resolved."
	if strings.TrimSpace(reason) != "" {
		note = "Resolved: " + strings.TrimSpace(reason)
	}
	th.Body = appendEntry(th.Body, stamp, orDefault(author, "unknown"), note)
	th.Status = "resolved"
	if author != "" && !contains(th.Participants, author) {
		th.Participants = append(th.Participants, author)
	}
	th.Updated = stamp
	return th, save(r, th, was)
}

// appendEntry adds a dated entry to the body. Entries are headed by their
// second and author, and a recommendation's mark; one with the same heading
// as the last joins that entry, as narrative log entries and issue instances
// do (I-0043): a second identical heading fails the duplicate-heading rule
// (MD024).
func appendEntry(body, stamp, author, text string) string {
	body = strings.TrimRight(body, "\n")
	heading := fmt.Sprintf("### %s %s", stamp, author)
	if locs := entryHeading.FindAllStringIndex(body, -1); len(locs) > 0 {
		last := locs[len(locs)-1]
		if body[last[0]:last[1]] == heading {
			return body + "\n\n" + strings.TrimSpace(text) + "\n"
		}
	}
	return body + "\n\n" + heading + "\n" + strings.TrimSpace(text) + "\n"
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

// OperatorName says whether name is one of the operator's names for a thread
// on story: the manifest's owner and the story's owner, those ADR-0097 takes
// an answer from, or "designer" when neither is named, as the dashboard's
// inbox.designer reads it (ADR-0181).
func OperatorName(r *workitem.Repo, story, name string) bool {
	owner := r.Manifest.Owner
	if owner != "" && name == owner {
		return true
	}
	storyOwner := ""
	if story != "" {
		if it, err := r.Get(story); err == nil {
			storyOwner = it.Owner
		}
	}
	if storyOwner != "" {
		return name == storyOwner
	}
	return owner == "" && name == "designer"
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
		before, after := doc[:i], strings.TrimPrefix(doc[j+len(mirrorEnd):], "\n")
		if block == "" && strings.HasSuffix(before, "\n\n") {
			// the block went with a blank line above it: keep one, not two (MD012)
			after = strings.TrimLeft(after, "\n")
		}
		doc = before + block + after
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
	return atomicfile.WriteFile(path, []byte(doc), 0o644)
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
// matter, the path relative to the repository, the entries, the story the
// thread belongs to, and the recommendation awaiting the operator's
// confirmation, null when there is none.
func View(r *workitem.Repo, th *Thread) map[string]any {
	path := th.Path
	if rel, err := filepath.Rel(r.MainRoot, th.Path); err == nil {
		path = filepath.ToSlash(rel)
	}
	return map[string]any{
		"id": th.ID, "title": th.Title, "anchor": th.Anchor, "status": th.Status,
		"participants": th.Participants, "created": th.Created, "updated": th.Updated,
		"path": path, "entries": th.Entries(), "story": StoryOf(r, th),
		"pending_recommendation": th.PendingRecommendation(),
	}
}
