package workitem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// Narrative is the front matter of wip/agents/<story>.md.
type Narrative struct {
	Stream  string `yaml:"stream" json:"stream"`
	Title   string `yaml:"title" json:"title"`
	Updated string `yaml:"updated" json:"updated"`
	Agent   string `yaml:"agent" json:"agent"`
	Session string `yaml:"session" json:"session"`
	// Host names the host whose flai stream open last opened the stream
	// (ADR-0064); empty in a narrative opened before it was recorded.
	Host string `yaml:"host,omitempty" json:"host,omitempty"`

	Path string `yaml:"-" json:"path"`
	Body string `yaml:"-" json:"body"`
}

// NarrativePath is wip/agents/<id>.md.
func (r *Repo) NarrativePath(storyID string) string {
	return filepath.Join(r.AgentsDir(), storyID+".md")
}

// ReadNarrative loads a narrative.
func ReadNarrative(path string) (*Narrative, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fm, body, err := SplitFrontMatter(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	var n Narrative
	if err := yaml.Unmarshal([]byte(fm), &n); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	n.Path, n.Body = path, body
	return &n, nil
}

// Marshal renders the narrative document.
func (n *Narrative) Marshal() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "stream: %s\n", n.Stream)
	fmt.Fprintf(&b, "title: %s\n", Scalar(n.Title))
	fmt.Fprintf(&b, "updated: %s\n", n.Updated)
	fmt.Fprintf(&b, "agent: %s\n", Scalar(n.Agent))
	fmt.Fprintf(&b, "session: %s\n", Scalar(n.Session))
	if n.Host != "" {
		fmt.Fprintf(&b, "host: %s\n", Scalar(n.Host))
	}
	b.WriteString("---\n")
	b.WriteString(n.Body)
	return b.String()
}

// Save writes the narrative.
func (n *Narrative) Save() error {
	if err := os.MkdirAll(filepath.Dir(n.Path), 0o755); err != nil {
		return err
	}
	return atomicfile.WriteFile(n.Path, []byte(n.Marshal()), 0o644)
}

// StreamOptions identify the agent opening or logging a stream.
type StreamOptions struct {
	Agent   string
	Session string
	// Host is the name of the host opening the stream, recorded by
	// OpenStream and ReopenStream.
	Host string
	Now  time.Time
}

// ThisHost is the name this host goes by, as a narrative records it, or ""
// when the system does not say.
func ThisHost() string {
	h, err := os.Hostname()
	if err != nil {
		return ""
	}
	return h
}

// StreamExistsError is OpenStream's refusal of a story whose narrative
// exists already.
type StreamExistsError struct{ ID, Path string }

func (e *StreamExistsError) Error() string {
	return fmt.Sprintf("stream %s already exists at %s", e.ID, e.Path)
}

// OpenStream creates the narrative for a story from the template. It fails
// with a *StreamExistsError if one already exists.
func (r *Repo) OpenStream(story *Item, opt StreamOptions) (*Narrative, error) {
	if story.Type != Story {
		return nil, fmt.Errorf("%s is a %s; streams belong to stories", story.ID, story.Type)
	}
	path := r.NarrativePath(story.ID)
	if _, err := os.Stat(path); err == nil {
		return nil, &StreamExistsError{ID: story.ID, Path: path}
	}
	now := opt.Now.UTC().Format(TimeFormat)
	doc, err := r.renderItemTemplate("narrative", map[string]any{
		"id": story.ID, "title": story.Title, "now": now, "today": now[:10],
		"agent": orDefault(opt.Agent, "agent"), "session": opt.Session,
	})
	if err != nil {
		return nil, err
	}
	fm, body, err := SplitFrontMatter(doc)
	if err != nil {
		return nil, fmt.Errorf("narrative template: %w", err)
	}
	var n Narrative
	if err := yaml.Unmarshal([]byte(fm), &n); err != nil {
		return nil, fmt.Errorf("narrative template: %w", err)
	}
	n.Path, n.Body, n.Host = path, body, opt.Host
	if err := n.Save(); err != nil {
		return nil, err
	}
	return &n, nil
}

// ReopenStream takes up a story's narrative that exists already, as on a
// host other than the one that opened it (ADR-0064): it keeps the body, and
// records the agent, session, and host opening it now where they are given
// and differ. It saves only what changed.
func (r *Repo) ReopenStream(story *Item, opt StreamOptions) (*Narrative, error) {
	n, err := ReadNarrative(r.NarrativePath(story.ID))
	if err != nil {
		return nil, err
	}
	was := n.Marshal()
	if opt.Agent != "" {
		n.Agent = opt.Agent
	}
	if opt.Session != "" {
		n.Session = opt.Session
	}
	if opt.Host != "" {
		n.Host = opt.Host
	}
	if n.Marshal() == was {
		return n, nil
	}
	n.Updated = opt.Now.UTC().Format(TimeFormat)
	if err := r.LintGuard(n.Path, was, n.Marshal()); err != nil {
		return nil, err
	}
	if err := n.Save(); err != nil {
		return nil, err
	}
	return n, nil
}

// LogStream appends a timestamped entry to the narrative's Log section.
func (r *Repo) LogStream(storyID, entry string, opt StreamOptions) (*Narrative, error) {
	if strings.TrimSpace(entry) == "" {
		return nil, fmt.Errorf("an entry is required")
	}
	var n *Narrative
	var err error
	for _, cand := range idCandidates(storyID) {
		n, err = ReadNarrative(r.NarrativePath(cand))
		if !os.IsNotExist(err) {
			break
		}
	}
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no stream for %s; open one with `flai stream open %s`", storyID, CanonicalID(storyID))
	}
	if err != nil {
		return nil, err
	}
	was := n.Marshal()
	now := opt.Now.UTC().Format(TimeFormat)
	if !strings.Contains(n.Body, "\n## Log") && !strings.HasPrefix(n.Body, "## Log") {
		n.Body = strings.TrimRight(n.Body, "\n") + "\n\n## Log\n"
	}
	// Headings are second-resolution timestamps; two entries in the same
	// second share one heading so the narrative never carries duplicate
	// headings (markdownlint MD024).
	body := strings.TrimRight(n.Body, "\n")
	if lastHeading(body) == "### "+now {
		n.Body = body + "\n\n" + strings.TrimSpace(entry) + "\n"
	} else {
		n.Body = body + "\n\n### " + now + "\n" + strings.TrimSpace(entry) + "\n"
	}
	n.Updated = now
	if opt.Agent != "" {
		n.Agent = opt.Agent
	}
	if opt.Session != "" {
		n.Session = opt.Session
	}
	if err := r.LintGuard(n.Path, was, n.Marshal()); err != nil {
		return nil, err
	}
	if err := n.Save(); err != nil {
		return nil, err
	}
	return n, nil
}

// The sections SetStreamState writes.
const (
	CurrentState = "Current state"
	NextSteps    = "Next steps"
)

// NarrativeSections are a narrative's second-level headings, without "## ",
// in the order the narrative template gives them.
var NarrativeSections = []string{"Context", CurrentState, NextSteps, "Decisions", "Open questions", "Log"}

// StreamRefusedError is SetStreamState's refusal of a story whose narrative
// state it may not write: one not in progress or in review, or one with no
// narrative.
type StreamRefusedError struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (e *StreamRefusedError) Error() string { return e.Reason }

// StreamState is what SetStreamState did to a story's narrative.
type StreamState struct {
	Stream string `json:"stream"`
	Path   string `json:"path"`
	// Updated is the narrative's updated stamp after the write, or as it was
	// when nothing changed.
	Updated string `json:"updated"`
	// Written are the sections given text, CurrentState and NextSteps, in
	// that order.
	Written []string `json:"written"`
	// Changed is false when the narrative already held the text given, and
	// nothing was written.
	Changed bool `json:"changed"`
}

// stateHeading is a line that would end a section written into a narrative:
// a first- or second-level heading.
var stateHeading = regexp.MustCompile(`^ {0,3}#{1,2}(\s|$)`)

// SetStreamState replaces what is under a story's narrative's ## Current
// state with current and under ## Next steps with next, leaves every other
// section as it was, and writes wip/agents/index.md again. An empty text
// leaves its section alone; a missing section is inserted where the template
// puts it. It stamps updated, and records the agent and session given, as
// LogStream does, and appends nothing to the log. It refuses, with a
// *StreamRefusedError, a story not in progress or in review and one with no
// narrative, and, with an *mdlint.Error, text the project's lint rejects.
func (r *Repo) SetStreamState(storyID, current, next string, opt StreamOptions) (*StreamState, error) {
	current, next = strings.TrimSpace(current), strings.TrimSpace(next)
	if current == "" && next == "" {
		return nil, fmt.Errorf("nothing to write to %s's narrative: give the current state, the next steps, or both", storyID)
	}
	texts := []struct{ name, text string }{{CurrentState, current}, {NextSteps, next}}
	for _, t := range texts {
		for _, l := range strings.Split(t.text, "\n") {
			if stateHeading.MatchString(l) {
				return nil, fmt.Errorf("the %s text for %s holds the heading %q, which would end the section: use ### or deeper for a heading inside it", strings.ToLower(t.name), storyID, l)
			}
		}
	}
	story, err := r.Get(storyID)
	if err != nil {
		return nil, err
	}
	if story.Type != Story {
		return nil, fmt.Errorf("%s is a %s; narratives belong to stories, so give its story's ID", story.ID, story.Type)
	}
	if story.Status != InProgress && story.Status != Review {
		reason := fmt.Sprintf("%s is %s, and only a story in progress or in review has its narrative's state written; move it to in-progress first with `flai move %s in-progress`", story.ID, story.Status, story.ID)
		if story.Status == Done || story.Status == Cancelled {
			reason = fmt.Sprintf("%s is %s, so its narrative is finished and keeps the state it closed with; write new work's state in the narrative of the story that takes it up", story.ID, story.Status)
		}
		return nil, &StreamRefusedError{ID: story.ID, Status: story.Status, Reason: reason}
	}
	n, err := ReadNarrative(r.NarrativePath(story.ID))
	if os.IsNotExist(err) {
		return nil, &StreamRefusedError{ID: story.ID, Status: story.Status,
			Reason: fmt.Sprintf("%s is %s but has no narrative at %s; open one with `flai stream open %s`", story.ID, story.Status, r.NarrativePath(story.ID), story.ID)}
	}
	if err != nil {
		return nil, err
	}
	res := &StreamState{Stream: n.Stream, Path: n.Path, Written: []string{}}
	body := n.Body
	for _, t := range texts {
		if t.text != "" {
			body = replaceSection(body, t.name, t.text)
			res.Written = append(res.Written, t.name)
		}
	}
	if body == n.Body {
		res.Updated = n.Updated
		return res, nil
	}
	was := n.Marshal()
	n.Body = body
	n.Updated = opt.Now.UTC().Format(TimeFormat)
	if opt.Agent != "" {
		n.Agent = opt.Agent
	}
	if opt.Session != "" {
		n.Session = opt.Session
	}
	if err := r.LintGuard(n.Path, was, n.Marshal()); err != nil {
		return nil, err
	}
	if err := n.Save(); err != nil {
		return nil, err
	}
	res.Updated, res.Changed = n.Updated, true
	items, err := r.List(false)
	if err == nil {
		err = r.WriteIndex(items, opt.Now)
	}
	if err != nil {
		return nil, fmt.Errorf("write the narratives' index after writing %s's state, which was written; any flai move writes the index again: %w", story.ID, err)
	}
	return res, nil
}

// Section is one second-level section of a narrative.
type Section struct {
	// Name is the heading without "## ".
	Name string
	// Line is the heading's line, 1-based in the text it was found in.
	Line int
	// Text is what is under the heading up to the next second-level
	// heading, without the blank lines at either end.
	Text string
}

// NarrativeSection finds the section headed "## <name>" in text, a
// narrative's body or its whole file; ok is false when text has none.
func NarrativeSection(text, name string) (s Section, ok bool) {
	lines := strings.Split(text, "\n")
	start, end := sectionSpan(lines, name)
	if start < 0 {
		return Section{}, false
	}
	under := lines[start+1 : end]
	for len(under) > 0 && strings.TrimSpace(under[0]) == "" {
		under = under[1:]
	}
	for len(under) > 0 && strings.TrimSpace(under[len(under)-1]) == "" {
		under = under[:len(under)-1]
	}
	return Section{Name: name, Line: start + 1, Text: strings.Join(under, "\n")}, true
}

// listMarker is the indent and list marker that open a line, which
// scripts/close-out.sh strips before it asks whether a line says anything.
var listMarker = regexp.MustCompile(`^\s*([-*]|[0-9]+\.)?\s*`)

// Written reports whether the section says something: a line with more than
// a list marker, as the close-out requires of Current state and Next steps.
// The template's lone "1." under Next steps is not written.
func (s Section) Written() bool {
	for _, l := range strings.Split(s.Text, "\n") {
		if listMarker.ReplaceAllString(l, "") != "" {
			return true
		}
	}
	return false
}

// sectionSpan is the index in lines of the heading "## <name>" and of the
// next second-level heading after it, or len(lines) when none follows; start
// is -1 when lines have no such heading.
func sectionSpan(lines []string, name string) (start, end int) {
	start = -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "## "+name {
			start = i
			break
		}
	}
	if start < 0 {
		return -1, -1
	}
	for i := start + 1; i < len(lines); i++ {
		if narrativeSectionHeading.MatchString(lines[i]) {
			return start, i
		}
	}
	return start, len(lines)
}

// replaceSection puts text under body's "## <name>", in place of what was
// there, with one blank line after the heading and one before the next.
// Without such a heading it inserts one before the first section that
// follows it in NarrativeSections, or at the end when none does.
func replaceSection(body, name, text string) string {
	lines := strings.Split(body, "\n")
	start, end := sectionSpan(lines, name)
	if start < 0 {
		at := len(lines)
		for _, later := range NarrativeSections[slices.Index(NarrativeSections, name)+1:] {
			if s, _ := sectionSpan(lines, later); s >= 0 {
				at = s
				break
			}
		}
		if at == len(lines) {
			lines = append(strings.Split(strings.TrimRight(body, "\n"), "\n"), "", "## "+name)
			start, end = len(lines)-1, len(lines)
		} else {
			heading := []string{"## " + name}
			if at > 0 && strings.TrimSpace(lines[at-1]) != "" {
				heading = []string{"", "## " + name}
			}
			lines = slices.Concat(lines[:at], heading, lines[at:])
			start = at + len(heading) - 1
			end = start + 1
		}
	}
	out := slices.Concat(lines[:start+1], []string{""}, strings.Split(text, "\n"), []string{""}, lines[end:])
	return strings.Join(out, "\n")
}

var narrativeSectionHeading = regexp.MustCompile(`^##\s`)
var narrativeBullet = regexp.MustCompile(`^\s*[-*]\s+(.*\S)\s*$`)
var narrativeContinuation = regexp.MustCompile(`^\s{2,}\S`)

// OpenQuestions are the bullets under a narrative's ## Open questions,
// outside the block flai generates to mirror threads (a mirrored thread is
// tracked, and closes, as a thread; a hand-written bullet stays until it is
// answered with AnswerOpenQuestion or removed by hand).
func OpenQuestions(body string) []string {
	lines := strings.Split(body, "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "## Open questions" {
			start = i
			break
		}
	}
	if start < 0 {
		return nil
	}
	var out []string
	generated := false
	for _, line := range lines[start+1:] {
		if narrativeSectionHeading.MatchString(line) {
			break
		}
		switch {
		case strings.Contains(line, "<!-- threads:start -->"):
			generated = true
		case strings.Contains(line, "<!-- threads:end -->"):
			generated = false
		case !generated:
			if m := narrativeBullet.FindStringSubmatch(line); m != nil {
				out = append(out, m[1])
			} else if len(out) > 0 && narrativeContinuation.MatchString(line) {
				out[len(out)-1] += " " + strings.TrimSpace(line)
			}
		}
	}
	return out
}

// AnswerOpenQuestion removes the bullet under a story's narrative's ## Open
// questions matching question (exactly as OpenQuestions returns it, so the
// caller passes back a title an inbox entry already showed) and records the
// answer under ## Decisions (design/system/agent-narrative.md: "Answered
// questions move to Decisions"). It refuses if none matches, so an answer
// can never silently land nowhere.
func (r *Repo) AnswerOpenQuestion(storyID, question, answer, by string, now time.Time) (*Narrative, error) {
	var n *Narrative
	var err error
	for _, cand := range idCandidates(storyID) {
		n, err = ReadNarrative(r.NarrativePath(cand))
		if !os.IsNotExist(err) {
			break
		}
	}
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no stream for %s; open one with `flai stream open %s`", storyID, CanonicalID(storyID))
	}
	if err != nil {
		return nil, err
	}

	lines := strings.Split(n.Body, "\n")
	secStart := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "## Open questions" {
			secStart = i
			break
		}
	}
	if secStart < 0 {
		return nil, fmt.Errorf("%s's narrative has no ## Open questions section", n.Stream)
	}
	secEnd := len(lines)
	for i := secStart + 1; i < len(lines); i++ {
		if narrativeSectionHeading.MatchString(lines[i]) {
			secEnd = i
			break
		}
	}

	bStart, bEnd, generated := -1, -1, false
	for i := secStart + 1; i < secEnd; {
		switch {
		case strings.Contains(lines[i], "<!-- threads:start -->"):
			generated, i = true, i+1
		case strings.Contains(lines[i], "<!-- threads:end -->"):
			generated, i = false, i+1
		case generated:
			i++
		default:
			m := narrativeBullet.FindStringSubmatch(lines[i])
			if m == nil {
				i++
				continue
			}
			text, j := m[1], i+1
			for j < secEnd && narrativeContinuation.MatchString(lines[j]) {
				text += " " + strings.TrimSpace(lines[j])
				j++
			}
			if text == question {
				bStart, bEnd = i, j
				i = secEnd // found it; stop
			} else {
				i = j
			}
		}
	}
	if bStart < 0 {
		return nil, fmt.Errorf("%s's narrative has no open question matching %q", n.Stream, question)
	}
	lines = append(lines[:bStart], lines[bEnd:]...)

	entry := fmt.Sprintf("- %s: %s — %s (answered by %s)", now.UTC().Format("2006-01-02"), question, strings.TrimSpace(answer), orDefault(by, "designer"))
	dStart := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "## Decisions" {
			dStart = i
			break
		}
	}
	if dStart < 0 {
		n.Body = strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n\n## Decisions\n" + entry + "\n"
	} else {
		dEnd := len(lines)
		for i := dStart + 1; i < len(lines); i++ {
			if narrativeSectionHeading.MatchString(lines[i]) {
				dEnd = i
				break
			}
		}
		for dEnd > dStart+1 && strings.TrimSpace(lines[dEnd-1]) == "" {
			dEnd--
		}
		out := append([]string{}, lines[:dEnd]...)
		out = append(out, entry)
		out = append(out, lines[dEnd:]...)
		n.Body = strings.Join(out, "\n")
	}
	n.Updated = now.UTC().Format(TimeFormat)
	if err := n.Save(); err != nil {
		return nil, err
	}
	return n, nil
}

// ActiveNarratives lists narratives under wip/agents: the S-*.md files, so
// not the strategic agents' activity documents.
func (r *Repo) ActiveNarratives() ([]*Narrative, error) {
	matches, _ := filepath.Glob(filepath.Join(r.AgentsDir(), "S-*.md"))
	var out []*Narrative
	for _, m := range matches {
		n, err := ReadNarrative(m)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return lessID(out[i].Stream, out[j].Stream) })
	return out, nil
}

// WriteIndex regenerates wip/agents/index.md from the active narratives and,
// under a heading of their own, the strategic agents' activity documents.
func (r *Repo) WriteIndex(items []*Item, now time.Time) error {
	narratives, err := r.ActiveNarratives()
	if err != nil {
		return err
	}
	status := map[string]string{}
	for _, it := range items {
		status[it.ID] = it.Status
	}
	var b strings.Builder
	fmt.Fprintf(&b, "---\ntitle: Active streams\nupdated: %s\n---\n\n# Active streams\n\n", now.UTC().Format(TimeFormat))
	b.WriteString("| Stream | Title | Status | Last agent | Updated |\n|--------|-------|--------|------------|---------|\n")
	for _, n := range narratives {
		fmt.Fprintf(&b, "| [%s](%s.md) | %s | %s | %s | %s |\n", n.Stream, n.Stream, n.Title, orDefault(status[n.Stream], "unknown"), n.Agent, n.Updated)
	}
	// An activity document that does not parse is listed as such rather than
	// failing the index, which every move rewrites; flai check names the fault.
	var rows []string
	for _, kind := range ActivityKinds {
		a, err := ReadActivity(r.ActivityPath(kind))
		switch {
		case errors.Is(err, os.ErrNotExist):
		case err != nil:
			rows = append(rows, fmt.Sprintf("| [%s](%s.md) | unreadable, see flai check | | | |\n", kind, kind))
		default:
			rows = append(rows, fmt.Sprintf("| [%s](%s.md) | %d | %s USD | %d | %s |\n", a.Kind, a.Kind, a.TasksCompleted, formatCost(a.AccruedCost), a.AccruedSeconds, orDefault(a.LastRun, "none")))
		}
	}
	if len(rows) > 0 {
		b.WriteString("\n## Strategic agents\n\n| Agent | Activities | Cost | Seconds | Last run |\n|-------|------------|------|---------|----------|\n")
		for _, row := range rows {
			b.WriteString(row)
		}
	}
	if err := os.MkdirAll(r.AgentsDir(), 0o755); err != nil {
		return err
	}
	return atomicfile.WriteFile(filepath.Join(r.AgentsDir(), "index.md"), []byte(b.String()), 0o644)
}

// lastHeading returns the last markdown heading line in body, or "".
func lastHeading(body string) string {
	lines := strings.Split(body, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.HasPrefix(lines[i], "#") {
			return strings.TrimSpace(lines[i])
		}
	}
	return ""
}
