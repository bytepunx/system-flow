// Package topics reads which stories a document is for (ADR-0047): the
// topics on its front matter and on its headings, and the sections its body
// splits into, each with the topics that apply to it.
package topics

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// All is the topic that every story has.
const All = "all"

// Section is one heading and the text down to the next heading of any level,
// or the text before the first heading (Level 0, no Heading).
type Section struct {
	Heading string   `json:"heading,omitempty"` // text, without the #s and the topics comment
	Level   int      `json:"level"`             // 1 to 6; 0 for the text before the first heading
	Path    []string `json:"path"`              // headings from the outermost down to this one
	Line    int      `json:"line"`              // 1-based line of the heading in the whole file
	Own     []string `json:"own,omitempty"`     // topics on the heading's own line
	Topics  []string `json:"topics"`            // own, else the parent's, else the file's
	Text    string   `json:"text"`              // the heading line and what follows it
}

// Doc is a document read for its topics.
type Doc struct {
	Topics   []string  `json:"topics"`   // the file's: front matter, else the default given
	Declared bool      `json:"declared"` // the front matter carries topics
	Sections []Section `json:"sections"`
}

var (
	headingRe = regexp.MustCompile(`^(#{1,6})[ \t]+(.*?)[ \t]*$`)
	commentRe = regexp.MustCompile(`[ \t]*<!--[ \t]*topics:(.*?)-->[ \t]*$`)
	fenceRe   = regexp.MustCompile("^ {0,3}(```+|~~~+)")
	closing   = regexp.MustCompile(`(^|[ \t]+)#+$`)
)

// Parse reads a whole markdown file, front matter included. A file whose
// front matter has no topics takes def: [all] for a convention, nil for a
// document that is selected by topic only when it declares them.
func Parse(content string, def []string) (*Doc, error) {
	doc := &Doc{Topics: def}
	body, offset := content, 0
	if fm, b, err := workitem.SplitFrontMatter(content); err == nil {
		t, err := FromFrontMatter(fm)
		if err != nil {
			return nil, err
		}
		if len(t) > 0 {
			doc.Topics, doc.Declared = t, true
		}
		body = b
		offset = strings.Count(content[:len(content)-len(b)], "\n")
	}
	doc.Sections = split(body, offset, doc.Topics)
	return doc, nil
}

// FromFrontMatter reads the topics of a front matter block alone.
func FromFrontMatter(fm string) ([]string, error) {
	var f struct {
		Topics []string `yaml:"topics"`
	}
	if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
		return nil, fmt.Errorf("the front matter does not parse, or topics is not a list such as topics: [cli, go]: %w", err)
	}
	return clean(f.Topics), nil
}

// Heading reads one line: whether it is a heading, its level, its text, and
// the topics its comment carries.
func Heading(line string) (level int, text string, own []string, ok bool) {
	m := headingRe.FindStringSubmatch(line)
	if m == nil {
		return 0, "", nil, false
	}
	text = m[2]
	if c := commentRe.FindStringSubmatchIndex(text); c != nil {
		own = Split(text[c[2]:c[3]])
		text = strings.TrimRight(text[:c[0]], " \t")
	}
	return len(m[1]), closing.ReplaceAllString(text, ""), own, true
}

// Split reads a comma or space separated list of topics.
func Split(s string) []string {
	return clean(strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' || r == '\t' }))
}

// clean trims, drops empties, and keeps the first of each topic.
func clean(in []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, t := range in {
		t = strings.TrimSpace(t)
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

// split cuts the body at every heading outside fenced code. offset is the
// number of lines before the body, so Line counts in the whole file.
func split(body string, offset int, file []string) []Section {
	type open struct {
		level  int
		text   string
		topics []string
	}
	var (
		out   []Section
		stack []open
		cur   = Section{Topics: file, Path: []string{}, Line: offset + 1}
		text  strings.Builder
		fence string
	)
	flush := func() {
		cur.Text = text.String()
		if cur.Level > 0 || strings.TrimSpace(cur.Text) != "" {
			out = append(out, cur)
		}
		text.Reset()
	}
	lines := strings.SplitAfter(body, "\n")
	for i, raw := range lines {
		if raw == "" {
			continue
		}
		line := strings.TrimRight(raw, "\r\n")
		if m := fenceRe.FindStringSubmatch(line); m != nil {
			switch {
			case fence == "":
				fence = m[1]
			case strings.HasPrefix(m[1], fence[:1]) && len(m[1]) >= len(fence) && strings.TrimSpace(line[len(m[0]):]) == "":
				fence = ""
			}
		}
		if level, heading, own, ok := Heading(line); ok && fence == "" {
			flush()
			for len(stack) > 0 && stack[len(stack)-1].level >= level {
				stack = stack[:len(stack)-1]
			}
			effective := own
			if len(effective) == 0 {
				effective = file
				if len(stack) > 0 {
					effective = stack[len(stack)-1].topics
				}
			}
			stack = append(stack, open{level, heading, effective})
			path := make([]string, len(stack))
			for j, s := range stack {
				path[j] = s.text
			}
			cur = Section{Heading: heading, Level: level, Path: path, Line: offset + i + 1, Own: own, Topics: effective}
		}
		text.WriteString(raw)
	}
	flush()
	return out
}

var (
	wordRe   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
	keyRe    = regexp.MustCompile(`^topics:`)
	nestedRe = regexp.MustCompile(`^([ \t]+|-[ \t]|-$)`)
)

// Valid reports whether a topic is one word: letters, digits, dot, dash, or
// underscore, starting with a letter or digit.
func Valid(t string) bool { return wordRe.MatchString(t) }

// SetInFrontMatter writes topics into a front matter block, in place of the
// key when it is there and at the end otherwise, and changes nothing else.
// No topics removes the key.
func SetInFrontMatter(fm string, list []string) string {
	line := ""
	if len(list) > 0 {
		line = "topics: [" + strings.Join(list, ", ") + "]"
	}
	lines := strings.Split(fm, "\n")
	var out []string
	done := false
	for i := 0; i < len(lines); i++ {
		if !keyRe.MatchString(lines[i]) {
			out = append(out, lines[i])
			continue
		}
		for i+1 < len(lines) && nestedRe.MatchString(lines[i+1]) {
			i++ // a block list's items belong to the key
		}
		if line != "" && !done {
			out = append(out, line)
		}
		done = true
	}
	if !done && line != "" {
		if n := len(out); n > 0 && out[n-1] == "" {
			out = append(out[:n-1], line, "")
		} else {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// WithoutTopics is a front matter block with its topics key taken out, to
// compare two blocks on everything but their topics.
func WithoutTopics(fm string) string { return SetInFrontMatter(fm, nil) }
