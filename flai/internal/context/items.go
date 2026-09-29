package context

import (
	"strings"
)

// Item is one thing the pack prints of the design documents: a whole
// document, or a run of its sections that one reason chose.
type Item struct {
	Path    string   `json:"path"`
	Title   string   `json:"title"`
	Heading []string `json:"heading,omitempty"` // heading path of the first section; empty for the whole document
	Label   string   `json:"label,omitempty"`   // the heading path as printed, without the document's title
	Step    string   `json:"step"`              // topics, linked, or ranked
	Reason  string   `json:"reason"`
	Also    []string `json:"also,omitempty"` // the reasons found for it after the first
	Size    int      `json:"size"`           // bytes of Text
	Text    string   `json:"text"`
}

// Items is what the selection prints, in the order the reasons were found,
// and within one reason in document order. A reason's sections of one
// document that follow each other print as one item; all of a document
// under one reason prints whole, front matter included.
func (s *Selection) Items() []Item {
	out := []Item{}
	for ci, c := range s.choices {
		for _, d := range s.Docs {
			own := s.owner[d]
			for i := 0; i < len(own); i++ {
				if own[i] != ci {
					continue
				}
				j := i
				for j+1 < len(own) && own[j+1] == ci {
					j++
				}
				it := Item{Path: d.Path, Title: d.Title, Step: c.step, Reason: c.reason, Also: c.also}
				if i == 0 && j == len(own)-1 {
					it.Text = d.Raw
				} else {
					var b strings.Builder
					for k := i; k <= j; k++ {
						b.WriteString(d.Sections[k].Text)
					}
					it.Text = b.String()
					it.Heading, it.Label = d.Sections[i].Path, d.Label(i)
				}
				it.Size = len(it.Text)
				out = append(out, it)
				i = j
			}
		}
	}
	return out
}

// Outline is one heading of a document loaded in part.
type Outline struct {
	Level   int    `json:"level"`
	Heading string `json:"heading"`
	Loaded  bool   `json:"loaded"`
}

// Entry is one line of the catalog: a document not loaded, or loaded in
// part with its outline.
type Entry struct {
	Path         string    `json:"path"`
	Title        string    `json:"title"`
	SupersededBy []string  `json:"superseded_by,omitempty"`
	Outline      []Outline `json:"outline,omitempty"` // only for a document loaded in part
}

// Catalog lists every document the selection did not load whole, in
// document order: those not loaded at all, then those loaded in part with
// the outline of their headings, so an agent can ask for the rest.
func (s *Selection) Catalog() (notLoaded, inPart []Entry) {
	notLoaded, inPart = []Entry{}, []Entry{}
	for _, d := range s.Docs {
		e := Entry{Path: d.Path, Title: d.Title}
		if s.Superseded(d) {
			e.SupersededBy = d.SupersededBy
		}
		switch {
		case !s.Loaded(d):
			notLoaded = append(notLoaded, e)
		case !s.Whole(d):
			for i, sec := range d.Sections {
				if sec.Level > 0 {
					e.Outline = append(e.Outline, Outline{Level: sec.Level, Heading: sec.Heading, Loaded: s.owner[d][i] >= 0})
				}
			}
			inPart = append(inPart, e)
		}
	}
	return notLoaded, inPart
}
