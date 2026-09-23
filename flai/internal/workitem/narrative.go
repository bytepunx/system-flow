package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// Narrative is the front matter of wip/agents/<story>.md.
type Narrative struct {
	Stream  string `yaml:"stream" json:"stream"`
	Title   string `yaml:"title" json:"title"`
	Updated string `yaml:"updated" json:"updated"`
	Agent   string `yaml:"agent" json:"agent"`
	Session string `yaml:"session" json:"session"`

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
	b.WriteString("---\n")
	b.WriteString(n.Body)
	return b.String()
}

// Save writes the narrative.
func (n *Narrative) Save() error {
	if err := os.MkdirAll(filepath.Dir(n.Path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(n.Path, []byte(n.Marshal()), 0o644)
}

// StreamOptions identify the agent opening or logging a stream.
type StreamOptions struct {
	Agent   string
	Session string
	Now     time.Time
}

// OpenStream creates the narrative for a story from the template. It fails
// if one already exists.
func (r *Repo) OpenStream(story *Item, opt StreamOptions) (*Narrative, error) {
	if story.Type != Story {
		return nil, fmt.Errorf("%s is a %s; streams belong to stories", story.ID, story.Type)
	}
	path := r.NarrativePath(story.ID)
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("stream %s already exists at %s", story.ID, path)
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
	n.Path, n.Body = path, body
	if err := n.Save(); err != nil {
		return nil, err
	}
	return &n, nil
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
	if err := n.Save(); err != nil {
		return nil, err
	}
	return n, nil
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

// ActiveNarratives lists narratives under wip/agents.
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

// WriteIndex regenerates wip/agents/index.md from the active narratives.
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
	if err := os.MkdirAll(r.AgentsDir(), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(r.AgentsDir(), "index.md"), []byte(b.String()), 0o644)
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
