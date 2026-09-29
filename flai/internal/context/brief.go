package context

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// brief is a document the pack describes instead of loading (ADR-0049): a
// design or tech file its topics select, or an ADR reached one step away.
type brief struct {
	doc    *Doc
	reason string
	also   []string
	secs   map[int]bool // the sections of a design file its topics selected
}

// addBrief briefs d for a reason, or lists the reason on d's brief when it
// has one. secs are the sections that selected it, if any.
func (s *Selection) addBrief(d *Doc, reason string, secs []int) {
	b := s.briefOf[d]
	if b == nil {
		b = &brief{doc: d, reason: reason, secs: map[int]bool{}}
		s.briefOf[d] = b
		s.briefs = append(s.briefs, b)
	} else if b.reason != reason && !slices.Contains(b.also, reason) {
		b.also = append(b.also, reason)
	}
	for _, i := range secs {
		b.secs[i] = true
	}
}

// briefADR briefs an ADR, or what supersedes it when it is superseded, each
// replacement with the reason "supersedes ADR-nnnn". An ADR the pack loads
// whole is not briefed; the reason is listed on what loads it.
func (s *Selection) briefADR(d *Doc, reason string) {
	for _, r := range s.Replace(d) {
		why := reason
		if r.by != nil {
			why = "supersedes " + r.by.ID
		}
		if s.Whole(r.doc) {
			s.choose(r.doc, r.doc.all(), "", why)
			continue
		}
		s.addBrief(r.doc, why, nil)
	}
}

// Briefed says whether d is briefed.
func (s *Selection) Briefed(d *Doc) bool { return s.briefOf[d] != nil }

// FirstParagraph is the first paragraph of a document's prose: the text
// under its title before the first subheading, or else under the first
// heading that has any, with its lines joined.
func (d *Doc) FirstParagraph() string {
	for _, sec := range d.Sections {
		body := sec.Text
		if sec.Level > 0 {
			_, body, _ = strings.Cut(body, "\n")
		} else if _, rest, err := workitem.SplitFrontMatter(body); err == nil {
			body = rest
		}
		if p := paragraph(body); p != "" {
			return p
		}
	}
	return ""
}

// paragraph is the first run of non-blank lines in text, joined with
// spaces; empty when text has only blank lines.
func paragraph(text string) string {
	var lines []string
	for _, l := range strings.Split(text, "\n") {
		t := strings.TrimSpace(l)
		if t == "" {
			if len(lines) > 0 {
				break
			}
			continue
		}
		lines = append(lines, t)
	}
	return strings.Join(lines, " ")
}

// DecisionSentence is an ADR's brief: the first sentence under its
// "## Decision" heading, with bold markers taken out, and the one after it
// too when the first is a lead-in of under five words ("MCP over HTTP.").
// ok is false when the section is missing or does not open with a sentence:
// it opens with a list, a table, a quote, a code block, or a heading, or
// its first paragraph does not end a sentence (ADR-0049).
func DecisionSentence(raw string) (sentence string, ok bool) {
	d, err := ParseDoc("", KindADR, raw)
	if err != nil {
		return "", false
	}
	return d.Decision()
}

// Decision is the ADR's decision sentence, as DecisionSentence gives it.
func (d *Doc) Decision() (sentence string, ok bool) {
	for _, sec := range d.Sections {
		if sec.Level != 2 || sec.Heading != "Decision" {
			continue
		}
		_, body, _ := strings.Cut(sec.Text, "\n")
		first := strings.TrimSpace(firstLine(body))
		if first == "" || opensBlock(first) {
			return "", false
		}
		p := strings.NewReplacer("**", "", "__", "").Replace(paragraph(body))
		one, rest, ended := cutSentence(p)
		if !ended {
			return p, false
		}
		if len(strings.Fields(one)) < 5 && rest != "" {
			if two, _, ok := cutSentence(rest); ok {
				one += " " + two
			}
		}
		return one, true
	}
	return "", false
}

func firstLine(text string) string {
	for _, l := range strings.Split(text, "\n") {
		if strings.TrimSpace(l) != "" {
			return l
		}
	}
	return ""
}

// opensBlock says whether a line starts something other than a paragraph.
func opensBlock(line string) bool {
	for _, p := range []string{"- ", "* ", "+ ", "|", ">", "#", "```", "~~~", "<"} {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	digits := strings.TrimLeft(line, "0123456789")
	return len(digits) < len(line) && (strings.HasPrefix(digits, ". ") || strings.HasPrefix(digits, ") "))
}

// cutSentence cuts p after its first sentence: the first ., !, or ?
// outside backticks that ends p or is followed by a space. ended is false
// when there is none.
func cutSentence(p string) (sentence, rest string, ended bool) {
	code := false
	for i, r := range p {
		switch {
		case r == '`':
			code = !code
		case !code && (r == '.' || r == '!' || r == '?') && (i+1 == len(p) || p[i+1] == ' '):
			return p[:i+1], strings.TrimSpace(p[i+1:]), true
		}
	}
	return p, "", false
}

// designBrief is the text of a design or tech file's brief: its first
// paragraph and the outline of its headings below the title, each heading
// its topics selected marked when they did not select all of it, and each
// loaded marked.
func (s *Selection) designBrief(b *brief) string {
	d := b.doc
	var out strings.Builder
	if p := d.FirstParagraph(); p != "" {
		out.WriteString(p + "\n\n")
	}
	partial := false
	for i, sec := range d.Sections {
		partial = partial || (sec.Level > 0 && !b.secs[i])
	}
	for i, sec := range d.Sections {
		if sec.Level < 2 {
			continue
		}
		var marks []string
		if partial && b.secs[i] {
			marks = append(marks, "selected")
		}
		if s.owner[d][i] >= 0 {
			marks = append(marks, "loaded")
		}
		fmt.Fprintf(&out, "%s- %s", strings.Repeat("  ", sec.Level-2), sec.Heading)
		if len(marks) > 0 {
			fmt.Fprintf(&out, " (%s)", strings.Join(marks, ", "))
		}
		out.WriteString("\n")
	}
	return out.String()
}

// adrBrief is the text of an ADR's brief: its decision sentence, or its
// first paragraph when the decision does not open with one.
func adrBrief(d *Doc) string {
	if sen, ok := d.Decision(); ok || sen != "" {
		return sen
	}
	return d.FirstParagraph()
}
