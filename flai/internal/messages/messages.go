// Package messages manages wip/messages: conversations between the agents of
// two open stories, one file per conversation, kept apart from the
// operator's threads and closed when either story leaves (ADR-0120).
package messages

import (
	"errors"
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
const Folder = "messages"

// Statuses is the closed list a conversation's status takes: open takes
// replies, closed takes none and is never reopened.
var Statuses = []string{StatusOpen, StatusClosed}

// The statuses a conversation is stored with.
const (
	StatusOpen   = "open"
	StatusClosed = "closed"
)

// titleLimit is the most runes a conversation's title keeps of the first
// line of its first message.
const titleLimit = 72

// closedPrefix starts the line of the entry that closes a conversation.
const closedPrefix = "Closed: "

// Conversation is one file: the messages between two stories.
type Conversation struct {
	ID           string   `yaml:"id" json:"id"`
	Title        string   `yaml:"title" json:"title"`
	From         string   `yaml:"from" json:"from"`
	To           string   `yaml:"to" json:"to"`
	About        []string `yaml:"about,omitempty" json:"about"`
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

// Entry is one dated contribution parsed from the body: its time, its
// author, the story it was written for, empty on an entry that closes the
// conversation, and its text.
type Entry struct {
	At     string `json:"at"`
	Author string `json:"author"`
	Story  string `json:"story"`
	Text   string `json:"text"`
}

var (
	idPattern    = regexp.MustCompile(`^MS-\d{4,}$`)
	looseID      = regexp.MustCompile(`^(?:(?i:ms)-?)?0*(\d+)$`)
	storyPattern = regexp.MustCompile(`^S-\d{3,}$`)
	entryHeading = regexp.MustCompile(`(?m)^### (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z) (.+?)(?: (S-\d{3,}))?$`)
)

// Dir is <layout.wip>/messages, always in the main checkout.
func Dir(r *workitem.Repo) string { return filepath.Join(r.WipDir(), Folder) }

// CanonicalID pads a conversation ID typed in any form (ms-7, MS-007, 7) to
// MS-0007. Anything else is returned unchanged.
func CanonicalID(id string) string {
	m := looseID.FindStringSubmatch(strings.TrimSpace(id))
	if m == nil {
		return id
	}
	return fmt.Sprintf("MS-%0*s", workitem.IDWidth, m[1])
}

// Parse decodes a conversation document, keeping front-matter fields this
// flai does not know in Unknown (S-0181).
func Parse(doc string) (*Conversation, error) {
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var c Conversation
	if err := yaml.Unmarshal([]byte(fm), &c); err != nil {
		return nil, err
	}
	c.Unknown = workitem.UnknownFields(fm, c)
	c.Body = body
	return &c, nil
}

// Read loads one file.
func Read(path string) (*Conversation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	c.Path = path
	return c, nil
}

// Validate checks the schema of one conversation's front matter.
func (c *Conversation) Validate() error {
	var errs []string
	if !idPattern.MatchString(c.ID) {
		errs = append(errs, fmt.Sprintf("id %q must look like MS-0001", c.ID))
	}
	if strings.TrimSpace(c.Title) == "" {
		errs = append(errs, "title is required")
	}
	for name, v := range map[string]string{"from": c.From, "to": c.To} {
		if !storyPattern.MatchString(v) {
			errs = append(errs, fmt.Sprintf("%s %q must be a story ID such as S-0001", name, v))
		}
	}
	if c.From != "" && workitem.CanonicalID(c.From) == workitem.CanonicalID(c.To) {
		errs = append(errs, fmt.Sprintf("from and to are both %s: a conversation is between two stories", c.From))
	}
	for _, p := range c.About {
		if _, err := cleanPath(p); err != nil {
			errs = append(errs, "about: "+err.Error())
		}
	}
	if !contains(Statuses, c.Status) {
		errs = append(errs, fmt.Sprintf("status %q must be one of %s", c.Status, strings.Join(Statuses, ", ")))
	}
	if len(c.Participants) == 0 {
		errs = append(errs, "participants must name at least the sender")
	}
	for name, v := range map[string]string{"created": c.Created, "updated": c.Updated} {
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
func (c *Conversation) Marshal() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "id: %s\n", c.ID)
	fmt.Fprintf(&b, "title: %s\n", workitem.Scalar(c.Title))
	fmt.Fprintf(&b, "from: %s\n", c.From)
	fmt.Fprintf(&b, "to: %s\n", c.To)
	if len(c.About) > 0 {
		fmt.Fprintf(&b, "about: %s\n", workitem.FlowList(c.About))
	}
	fmt.Fprintf(&b, "status: %s\n", c.Status)
	fmt.Fprintf(&b, "participants: %s\n", workitem.FlowList(c.Participants))
	fmt.Fprintf(&b, "created: %s\n", c.Created)
	fmt.Fprintf(&b, "updated: %s\n", c.Updated)
	workitem.WriteFields(&b, c.Unknown)
	b.WriteString("---\n")
	b.WriteString(c.Body)
	return b.String()
}

// Save writes the conversation to its path, through a temporary file renamed
// into place.
func (c *Conversation) Save() error {
	if err := os.MkdirAll(filepath.Dir(c.Path), 0o755); err != nil {
		return err
	}
	return atomicfile.WriteFile(c.Path, []byte(c.Marshal()), 0o644)
}

// Entries parses the dated entries from the body, in order. An entry headed
// with no story is one that closed the conversation.
func (c *Conversation) Entries() []Entry {
	locs := entryHeading.FindAllStringSubmatchIndex(c.Body, -1)
	var out []Entry
	for i, loc := range locs {
		end := len(c.Body)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		e := Entry{At: c.Body[loc[2]:loc[3]], Author: c.Body[loc[4]:loc[5]], Text: strings.TrimSpace(c.Body[loc[1]:end])}
		if loc[6] >= 0 {
			e.Story = c.Body[loc[6]:loc[7]]
		}
		out = append(out, e)
	}
	return out
}

// Awaiting is the story of the two that did not write the last entry written
// for a story: the addressee until it replies, then the sender, and so on.
func (c *Conversation) Awaiting() string {
	es := c.Entries()
	for i := len(es) - 1; i >= 0; i-- {
		switch workitem.CanonicalID(es[i].Story) {
		case "":
			continue
		case workitem.CanonicalID(c.From):
			return c.To
		default:
			return c.From
		}
	}
	return c.To
}

// Names reports whether the story, in any padding, is one of the two.
func (c *Conversation) Names(story string) bool {
	s := workitem.CanonicalID(story)
	return s == workitem.CanonicalID(c.From) || s == workitem.CanonicalID(c.To)
}

// Closed reports whether the conversation reads as closed, and why: its
// status is closed, or either story is done, cancelled, archived, or gone. A
// story whose file cannot be read does not close it; Reply reports the error.
func (c *Conversation) Closed(r *workitem.Repo) (closed bool, reason string) {
	if c.Status == StatusClosed {
		return true, c.closedReason()
	}
	for _, s := range []string{c.From, c.To} {
		if why := gone(r, s); why != "" {
			return true, why
		}
	}
	return false, ""
}

// closedReason is the reason the newest closing entry gives, or that the
// status says closed when no entry gives one.
func (c *Conversation) closedReason() string {
	es := c.Entries()
	for i := len(es) - 1; i >= 0; i-- {
		lines := strings.Split(es[i].Text, "\n")
		for j := len(lines) - 1; j >= 0; j-- {
			if why, ok := strings.CutPrefix(lines[j], closedPrefix); ok && strings.TrimSpace(why) != "" {
				return strings.TrimSpace(why)
			}
		}
	}
	return "its status is closed"
}

// gone says why a story has left, so that its conversations read as closed,
// or "" while it is open or cannot be read.
func gone(r *workitem.Repo, story string) string {
	it, err := r.Get(story)
	switch {
	case errors.Is(err, workitem.ErrNotFound):
		return story + " does not exist"
	case err != nil:
		return ""
	case it.Status == workitem.Done:
		return it.ID + " was accepted"
	case it.Status == workitem.Cancelled:
		return it.ID + " was cancelled"
	case it.Archived:
		return it.ID + " is archived"
	}
	return ""
}

// List loads every conversation, sorted by ID. A missing folder has none.
func List(r *workitem.Repo) ([]*Conversation, error) {
	matches, err := filepath.Glob(filepath.Join(Dir(r), "MS-*.md"))
	if err != nil {
		return nil, err
	}
	var out []*Conversation
	for _, m := range matches {
		c, err := Read(m)
		if err != nil {
			return nil, err
		}
		workitem.WarnUnknown(m, c.Unknown)
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return num(out[i].ID) < num(out[j].ID) })
	return out, nil
}

// Get finds one conversation by ID in any padding.
func Get(r *workitem.Repo, id string) (*Conversation, error) {
	canon := CanonicalID(id)
	for _, cand := range []string{id, canon} {
		matches, _ := filepath.Glob(filepath.Join(Dir(r), cand+"-*.md"))
		if len(matches) > 0 {
			return Read(matches[0])
		}
	}
	return nil, fmt.Errorf("conversation %s not found in %s: flai message list --all lists them", canon, Dir(r))
}

// For returns the conversations of a story, sent or received, the story
// named in any padding, in ID order.
func For(r *workitem.Repo, story string) ([]*Conversation, error) {
	all, err := List(r)
	if err != nil {
		return nil, err
	}
	var out []*Conversation
	for _, c := range all {
		if c.Names(story) {
			out = append(out, c)
		}
	}
	return out, nil
}

// Other is the story of the two that is not story, named in any padding.
func (c *Conversation) Other(story string) string {
	if workitem.CanonicalID(story) == workitem.CanonicalID(c.From) {
		return c.To
	}
	return c.From
}

// AwaitingOther returns story's conversations that read as open and await
// the other story's reply: those story sent or answered last, in ID order.
func AwaitingOther(r *workitem.Repo, story string) ([]*Conversation, error) {
	all, err := For(r, story)
	if err != nil {
		return nil, err
	}
	var out []*Conversation
	for _, c := range all {
		if closed, _ := c.Closed(r); closed {
			continue
		}
		if workitem.CanonicalID(c.Awaiting()) != workitem.CanonicalID(story) {
			out = append(out, c)
		}
	}
	return out, nil
}

// ToSince returns story's conversations with an entry the other story wrote
// after t, in ID order: those a message to story arrived in since then.
func ToSince(r *workitem.Repo, story string, t time.Time) ([]*Conversation, error) {
	all, err := For(r, story)
	if err != nil {
		return nil, err
	}
	var out []*Conversation
	for _, c := range all {
		other := workitem.CanonicalID(c.Other(story))
		for _, e := range c.Entries() {
			at, err := time.Parse(workitem.TimeFormat, e.At)
			if err == nil && at.After(t) && workitem.CanonicalID(e.Story) == other {
				out = append(out, c)
				break
			}
		}
	}
	return out, nil
}

// NextID allocates the next MS-nnnn: one more than the highest in the folder.
func NextID(r *workitem.Repo) string {
	highest := 0
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "MS-*.md"))
	for _, m := range matches {
		n, _, _ := strings.Cut(strings.TrimPrefix(filepath.Base(m), "MS-"), "-")
		if v, err := strconv.Atoi(strings.TrimSuffix(n, ".md")); err == nil && v > highest {
			highest = v
		}
	}
	return fmt.Sprintf("MS-%0*d", workitem.IDWidth, highest+1)
}

// SendOptions describe the first message of a conversation.
type SendOptions struct {
	From   string   // the sender's story
	To     string   // the addressee's story
	Author string   // the agent writing it
	Text   string   // the message
	About  []string // repository paths the conversation is about
	Now    time.Time
}

// Send starts a conversation between two stories in progress or in review
// with its first message, from opt.From to opt.To. Each about path must exist
// in the repository. The file is refused, and nothing written, when the
// project's lint rejects it.
func Send(r *workitem.Repo, opt SendOptions) (*Conversation, error) {
	m, err := prepare(r, opt)
	if err != nil {
		return nil, err
	}
	if m.about, err = resolveAbout(r, opt.About); err != nil {
		return nil, err
	}
	return begin(r, m)
}

// Notify sends a message from opt.From to opt.To on the conversation still
// open between the two stories, whichever of them started it, adding to its
// about the paths it does not name yet; with none open, it starts one as Send
// does. Unlike Send, it takes the about paths as given, cleaned, without
// asking that each exist: flai sends it for what git reports, and a path a
// commit deleted, or one only a story's worktree holds, is still one the
// message is about. The file is refused, and nothing written, when the
// project's lint rejects it.
func Notify(r *workitem.Repo, opt SendOptions) (*Conversation, error) {
	m, err := prepare(r, opt)
	if err != nil {
		return nil, err
	}
	if m.about, err = cleanAbout(opt.About); err != nil {
		return nil, err
	}
	c, err := openBetween(r, m.from.ID, m.to.ID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return begin(r, m)
	}
	was := c.Marshal()
	before := betweenLine(c.From, c.To, c.About)
	for _, p := range m.about {
		if !contains(c.About, p) {
			c.About = append(c.About, p)
		}
	}
	c.Body = strings.Replace(c.Body, "\n"+before+"\n", "\n"+betweenLine(c.From, c.To, c.About)+"\n", 1)
	c.add(m.now, m.author, storyOf(c, m.from.ID), m.text)
	return c, save(r, c, was)
}

// message is a first message, or a notice, checked and ready to write.
type message struct {
	text, author string
	from, to     *workitem.Item
	about        []string
	now          string
}

// prepare checks what Send and Notify share: the text, the author, and two
// different stories, each in progress or in review.
func prepare(r *workitem.Repo, opt SendOptions) (message, error) {
	m := message{text: strings.TrimSpace(opt.Text), author: cleanAuthor(opt.Author), now: opt.Now.UTC().Format(workitem.TimeFormat)}
	if m.text == "" {
		return m, fmt.Errorf("a message needs text")
	}
	if m.author == "" {
		return m, fmt.Errorf("an author is required (--by, FLAI_AGENT, or config author)")
	}
	if s := workitem.CanonicalID(strings.TrimSpace(opt.From)); s != "" && s == workitem.CanonicalID(strings.TrimSpace(opt.To)) {
		return m, fmt.Errorf("%s cannot message itself: name another story in progress or in review", s)
	}
	var err error
	if m.from, err = sendable(r, opt.From, "sender"); err != nil {
		return m, err
	}
	if m.to, err = sendable(r, opt.To, "addressee"); err != nil {
		return m, err
	}
	return m, nil
}

// begin writes a new conversation holding the message as its first entry.
func begin(r *workitem.Repo, m message) (*Conversation, error) {
	c := &Conversation{
		ID: NextID(r), Title: titleOf(m.text, m.from.ID), From: m.from.ID, To: m.to.ID, About: m.about,
		Status: StatusOpen, Participants: []string{m.author}, Created: m.now, Updated: m.now,
	}
	c.Path = filepath.Join(Dir(r), c.ID+"-"+orSlug(c.Title)+".md")
	c.Body = fmt.Sprintf("\n# %s %s\n\n%s\n\n## Entries\n\n### %s %s %s\n%s\n", c.ID, c.Title, betweenLine(c.From, c.To, c.About), m.now, m.author, c.From, m.text)
	if err := r.LintGuard(c.Path, "", c.Marshal()); err != nil {
		return nil, err
	}
	if err := c.Save(); err != nil {
		return nil, err
	}
	return c, nil
}

// betweenLine is the sentence under a conversation's heading naming its two
// stories and the paths it is about.
func betweenLine(from, to string, about []string) string {
	line := fmt.Sprintf("Between %s and %s", from, to)
	if len(about) > 0 {
		quoted := make([]string, len(about))
		for i, p := range about {
			quoted[i] = "`" + p + "`"
		}
		line += ", about " + strings.Join(quoted, ", ")
	}
	return line + "."
}

// openBetween is the newest conversation between the two stories, in either
// direction, that is stored open and reads as open, or nil when there is none.
func openBetween(r *workitem.Repo, a, b string) (*Conversation, error) {
	all, err := For(r, a)
	if err != nil {
		return nil, err
	}
	for i := len(all) - 1; i >= 0; i-- {
		c := all[i]
		if !c.Names(b) || c.Status != StatusOpen {
			continue
		}
		if closed, _ := c.Closed(r); !closed {
			return c, nil
		}
	}
	return nil, nil
}

// Reply appends a message from story, one of the conversation's two, which
// must still be in progress or in review. A conversation that reads as closed
// takes no reply.
func Reply(r *workitem.Repo, id, story, author, text string, now time.Time) (*Conversation, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("a reply needs text")
	}
	if author = cleanAuthor(author); author == "" {
		return nil, fmt.Errorf("an author is required (--by, FLAI_AGENT, or config author)")
	}
	c, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	if !c.Names(story) {
		return nil, fmt.Errorf("%s is not in %s, which is between %s and %s: reply from one of its two stories, or start a conversation of its own with flai message send", workitem.CanonicalID(story), c.ID, c.From, c.To)
	}
	if closed, why := c.Closed(r); closed {
		return nil, fmt.Errorf("%s is closed (%s) and takes no reply: send a new message with flai message send to start another", c.ID, why)
	}
	it, err := sendable(r, story, "replying story")
	if err != nil {
		return nil, err
	}
	was := c.Marshal()
	c.add(now.UTC().Format(workitem.TimeFormat), author, storyOf(c, it.ID), text)
	return c, save(r, c, was)
}

// add appends an entry by author for story at stamp, names author among the
// participants, and marks the conversation updated then.
func (c *Conversation) add(stamp, author, story, text string) {
	c.Body = appendEntry(c.Body, stamp+" "+author+" "+story, text)
	if !contains(c.Participants, author) {
		c.Participants = append(c.Participants, author)
	}
	c.Updated = stamp
}

// Close closes a conversation that is still stored open, with an entry by
// author, who may be the operator rather than a story's agent, giving the
// reason. A closed conversation is never reopened.
func Close(r *workitem.Repo, id, author, reason string, now time.Time) (*Conversation, error) {
	if author = cleanAuthor(author); author == "" {
		return nil, fmt.Errorf("an author is required (--by, FLAI_AGENT, or config author)")
	}
	c, err := Get(r, id)
	if err != nil {
		return nil, err
	}
	if c.Status == StatusClosed {
		return nil, fmt.Errorf("%s is closed already (%s): a closed conversation stays closed", c.ID, c.closedReason())
	}
	note := "Closed."
	if reason = strings.TrimSpace(reason); reason != "" {
		note = closedPrefix + reason
	}
	was := c.Marshal()
	stamp := now.UTC().Format(workitem.TimeFormat)
	c.Body = appendEntry(c.Body, stamp+" "+author, note)
	c.Status = StatusClosed
	if !contains(c.Participants, author) {
		c.Participants = append(c.Participants, author)
	}
	c.Updated = stamp
	return c, save(r, c, was)
}

// CloseOn closes every conversation still stored open that names one of the
// stories, in any padding, each with an entry by author reading "Closed:
// <story> was <verb>" (accepted, archived), and returns the IDs it closed in
// ID order. flai accept and flai archive call it for the stories they archive.
func CloseOn(r *workitem.Repo, stories []string, author, verb string, now time.Time) ([]string, error) {
	want := map[string]bool{}
	for _, s := range stories {
		want[workitem.CanonicalID(s)] = true
	}
	all, err := List(r)
	if err != nil {
		return nil, err
	}
	var closed []string
	for _, c := range all {
		if c.Status == StatusClosed {
			continue
		}
		var leaving []string
		for _, s := range []string{c.From, c.To} {
			if want[workitem.CanonicalID(s)] {
				leaving = append(leaving, s)
			}
		}
		if len(leaving) == 0 {
			continue
		}
		reason := fmt.Sprintf("%s was %s", leaving[0], verb)
		if len(leaving) == 2 {
			reason = fmt.Sprintf("%s and %s were %s", leaving[0], leaving[1], verb)
		}
		if _, err := Close(r, c.ID, author, reason, now); err != nil {
			return closed, fmt.Errorf("close %s, between %s and %s: %w; the conversations listed before it are closed, and it and those after it are still open", c.ID, c.From, c.To, err)
		}
		closed = append(closed, c.ID)
	}
	return closed, nil
}

// View is a conversation as flai prints it with --json and the dashboard
// reads it: the front matter, whether it reads as closed and why, the story
// it awaits while open, the path relative to the repository, and the entries.
func View(r *workitem.Repo, c *Conversation) map[string]any {
	path := c.Path
	if rel, err := filepath.Rel(r.MainRoot, c.Path); err == nil {
		path = filepath.ToSlash(rel)
	}
	closed, why := c.Closed(r)
	awaiting := ""
	if !closed {
		awaiting = c.Awaiting()
	}
	about, entries := c.About, c.Entries()
	if about == nil {
		about = []string{}
	}
	if entries == nil {
		entries = []Entry{}
	}
	return map[string]any{
		"id": c.ID, "title": c.Title, "from": c.From, "to": c.To, "about": about,
		"status": c.Status, "closed": closed, "closed_reason": why, "awaiting": awaiting,
		"participants": c.Participants, "created": c.Created, "updated": c.Updated,
		"path": path, "entries": entries,
	}
}

// sendable finds a story a message may go to or from: a story ID, of a story
// that exists, is not archived, and is in progress or in review. role names
// the side in the error.
func sendable(r *workitem.Repo, id, role string) (*workitem.Item, error) {
	canon := workitem.CanonicalID(strings.TrimSpace(id))
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("the %s's story is required: a message goes from one story to another, such as S-0001", role)
	}
	if !storyPattern.MatchString(canon) {
		return nil, fmt.Errorf("%q, the %s, is not a story ID: a message goes from one story to another, such as S-0001", id, role)
	}
	it, err := r.Get(canon)
	if errors.Is(err, workitem.ErrNotFound) {
		return nil, fmt.Errorf("%s, the %s, does not exist: a message goes only between stories in progress or in review", canon, role)
	}
	if err != nil {
		return nil, fmt.Errorf("read %s, the %s: %w", canon, role, err)
	}
	if it.Archived {
		return nil, fmt.Errorf("%s is archived; a message goes only between stories in progress or in review", it.ID)
	}
	if it.Status != workitem.InProgress && it.Status != workitem.Review {
		return nil, fmt.Errorf("%s is %s; a message goes only between stories in progress or in review", it.ID, stateWords(it.Status))
	}
	return it, nil
}

// stateWords reads a status in a sentence: "S-0012 is in the backlog".
func stateWords(status string) string {
	switch status {
	case workitem.Backlog:
		return "in the backlog"
	case workitem.InProgress:
		return "in progress"
	case workitem.Review:
		return "in review"
	}
	return status
}

// storyOf is the conversation's own spelling of a story it names, so that
// entries match its front matter.
func storyOf(c *Conversation, story string) string {
	if workitem.CanonicalID(story) == workitem.CanonicalID(c.From) {
		return c.From
	}
	return c.To
}

// resolveAbout checks that each path exists in the repository, in the
// working checkout or the main one, and returns them cleaned, with forward
// slashes, without repeats.
func resolveAbout(r *workitem.Repo, paths []string) ([]string, error) {
	out, err := cleanAbout(paths)
	if err != nil {
		return nil, err
	}
	for _, rel := range out {
		if _, err := os.Stat(filepath.Join(r.Root, filepath.FromSlash(rel))); err != nil {
			if _, err2 := os.Stat(filepath.Join(r.MainRoot, filepath.FromSlash(rel))); err2 != nil {
				return nil, fmt.Errorf("--about %s does not exist in the repository: name a file or folder relative to its root", rel)
			}
		}
	}
	return out, nil
}

// cleanAbout returns the paths cleaned, with forward slashes, without
// repeats, whether they exist or not.
func cleanAbout(paths []string) ([]string, error) {
	var out []string
	for _, p := range paths {
		rel, err := cleanPath(p)
		if err != nil {
			return nil, fmt.Errorf("--about: %w", err)
		}
		if !contains(out, rel) {
			out = append(out, rel)
		}
	}
	return out, nil
}

// cleanPath is a repository path as a conversation keeps it: relative to the
// root, cleaned, with forward slashes.
func cleanPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", fmt.Errorf("a path is empty: name a file or folder relative to the repository root")
	}
	rel := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(filepath.FromSlash(p))), "./")
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") || rel == ".." || strings.HasPrefix(rel, "../") || rel == "." {
		return "", fmt.Errorf("%s is not a path inside the repository: name a file or folder relative to its root", p)
	}
	return rel, nil
}

// titleOf is the conversation's title: the first non-empty line of the first
// message, cleaned as an item title is and cut on a word to titleLimit runes.
// A line that leaves nothing names the sender.
func titleOf(text, from string) string {
	line := ""
	for _, l := range strings.Split(text, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			line = l
			break
		}
	}
	title := workitem.CleanTitle(line)
	if runes := []rune(title); len(runes) > titleLimit {
		cut := string(runes[:titleLimit])
		if i := strings.LastIndex(cut, " "); i > 0 {
			cut = cut[:i]
		}
		title = workitem.CleanTitle(cut)
	}
	if title == "" {
		return "From " + from
	}
	return title
}

// cleanAuthor is an author name on one line with single spaces.
func cleanAuthor(s string) string { return strings.Join(strings.Fields(s), " ") }

// appendEntry adds a dated entry under its heading, the time, author, and
// story. One with the same heading as the last entry joins it, as threads'
// entries do (I-0043): a second identical heading fails the duplicate-heading
// rule (MD024).
func appendEntry(body, heading, text string) string {
	body = strings.TrimRight(body, "\n")
	heading = "### " + heading
	if locs := entryHeading.FindAllStringIndex(body, -1); len(locs) > 0 {
		last := locs[len(locs)-1]
		if body[last[0]:last[1]] == heading {
			return body + "\n\n" + text + "\n"
		}
	}
	return body + "\n\n" + heading + "\n" + text + "\n"
}

// save writes the conversation unless the change brings markdown the
// project's lint rejects (S-0179).
func save(r *workitem.Repo, c *Conversation, was string) error {
	if err := r.LintGuard(c.Path, was, c.Marshal()); err != nil {
		return err
	}
	return c.Save()
}

func num(id string) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, "MS-"))
	return n
}

func orSlug(title string) string {
	if s := template.Slug(title); s != "" {
		return s
	}
	return "message"
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
