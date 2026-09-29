package context

import (
	"fmt"
	"strings"
)

// Part is one section of a document fetched by its heading, with the
// sections below it: what doc_get and flai doc show --heading return
// (ADR-0049).
type Part struct {
	Path    string `json:"path"`
	Heading string `json:"heading" jsonschema:"the heading path of the section found, below the document's title"`
	Line    int    `json:"line" jsonschema:"1-based line of the heading in the file"`
	Size    int    `json:"size" jsonschema:"bytes in text"`
	Text    string `json:"text" jsonschema:"the section from its heading line down to the next heading at its level or higher"`
}

// Section reads the part of a markdown document a heading names. The
// heading is matched by its text, by a heading path whose parts are joined
// by › or >, which matches at the end of a section's path, or by its anchor
// slug, and by none of these more strictly than Slug compares. When it
// names more than one section, only those whose own heading is the same
// text are kept; still more than one, or none, is an error naming the
// sections there are to choose from.
func Section(path, raw, heading string) (*Part, error) {
	d, err := ParseDoc(path, "", raw)
	if err != nil {
		return nil, err
	}
	sec, err := d.find(heading)
	if err != nil {
		return nil, err
	}
	var text strings.Builder
	for _, i := range d.under(sec) {
		text.WriteString(d.Sections[i].Text)
	}
	return &Part{Path: path, Heading: d.Label(sec), Line: d.Sections[sec].Line, Size: text.Len(), Text: text.String()}, nil
}

func (d *Doc) find(heading string) (int, error) {
	want := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(heading), "#"))
	if want == "" {
		return 0, fmt.Errorf("give a heading of %s: %s", d.Path, d.headings())
	}
	parts := strings.FieldsFunc(want, func(r rune) bool { return r == '›' || r == '>' })
	found := d.matching([]string{want})
	if len(found) == 0 {
		found = d.matching(parts)
	}
	if len(found) > 1 && len(parts) > 0 {
		var exact []int
		for _, i := range found {
			if strings.TrimSpace(parts[len(parts)-1]) == d.Sections[i].Heading {
				exact = append(exact, i)
			}
		}
		if len(exact) > 0 {
			found = exact
		}
	}
	switch len(found) {
	case 1:
		return found[0], nil
	case 0:
		return 0, fmt.Errorf("no heading %q in %s; its headings are: %s", heading, d.Path, d.headings())
	}
	labels := make([]string, len(found))
	for n, i := range found {
		labels[n] = fmt.Sprintf("%q (line %d)", d.Label(i), d.Sections[i].Line)
	}
	return 0, fmt.Errorf("heading %q names %d sections of %s: %s; give more of the heading path, such as %q", heading, len(found), d.Path, strings.Join(labels, ", "), d.Label(found[0]))
}

// matching is the sections whose heading path ends with parts, compared by
// slug.
func (d *Doc) matching(parts []string) []int {
	want := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := Slug(strings.TrimSpace(p)); s != "" {
			want = append(want, s)
		}
	}
	if len(want) == 0 {
		return nil
	}
	var out []int
	for i, sec := range d.Sections {
		if sec.Level == 0 || len(sec.Path) < len(want) {
			continue
		}
		tail, ok := sec.Path[len(sec.Path)-len(want):], true
		for n, w := range want {
			ok = ok && Slug(tail[n]) == w
		}
		if ok {
			out = append(out, i)
		}
	}
	return out
}

// headings lists every heading path of the document, quoted, for an error.
func (d *Doc) headings() string {
	var out []string
	for i, sec := range d.Sections {
		if sec.Level > 0 {
			out = append(out, fmt.Sprintf("%q", d.Label(i)))
		}
	}
	if len(out) == 0 {
		return "none"
	}
	return strings.Join(out, ", ")
}
