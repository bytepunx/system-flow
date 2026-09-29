package context

import (
	"slices"
	"strings"
)

// Item is one thing the pack prints of the design documents: a whole
// document, a run of its sections that one reason chose, or a brief.
type Item struct {
	Path    string   `json:"path"`
	ID      string   `json:"id,omitempty"` // ADR-nnnn for an ADR
	Title   string   `json:"title"`
	Heading []string `json:"heading,omitempty"` // heading path of the first section; empty for the whole document
	Label   string   `json:"label,omitempty"`   // the heading path as printed, without the document's title
	Step    string   `json:"step"`              // named, briefed, or ranked
	Reason  string   `json:"reason"`
	Also    []string `json:"also,omitempty"`  // the reasons found for it after the first
	Size    int      `json:"size"`            // bytes of Text
	Whole   int      `json:"whole,omitempty"` // bytes of the whole document, for a brief
	Text    string   `json:"text"`            // what is loaded, or the brief
}

// Items is what the selection prints: what is named, then the briefs of
// design and tech files, then the briefs of ADRs, then what is ranked. Loaded
// items are in the order their reasons were found, and within one reason in
// document order; a reason's sections of one document that follow each
// other print as one item, and all of a document under one reason prints
// whole, front matter included. A briefed document that is loaded whole
// prints once, loaded, with the brief's reasons listed on it.
func (s *Selection) Items() []Item {
	named, ranked := s.loaded(StepNamed), s.loaded(StepRanked)
	var designs, adrs []Item
	for _, b := range s.briefs {
		d := b.doc
		if s.Whole(d) {
			for _, list := range [][]Item{named, ranked} {
				for i := range list {
					if list[i].Path == d.Path && list[i].Heading == nil {
						list[i].Also = appendNew(list[i].Also, append([]string{b.reason}, b.also...)...)
					}
				}
			}
			continue
		}
		it := Item{Path: d.Path, ID: d.ID, Title: d.Title, Step: StepBriefed, Reason: b.reason, Also: b.also, Whole: len(d.Raw)}
		if d.Kind == KindADR {
			it.Text = adrBrief(d)
			it.Size = len(it.Text)
			adrs = append(adrs, it)
			continue
		}
		it.Text = s.designBrief(b)
		it.Size = len(it.Text)
		designs = append(designs, it)
	}
	return append([]Item{}, slices.Concat(named, designs, adrs, ranked)...)
}

func appendNew(list []string, add ...string) []string {
	for _, a := range add {
		if !slices.Contains(list, a) {
			list = append(list, a)
		}
	}
	return list
}

// loaded is what the choices of one step load, as Items prints it.
func (s *Selection) loaded(step string) []Item {
	out := []Item{}
	for ci, c := range s.choices {
		if c.step != step {
			continue
		}
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
				it := Item{Path: d.Path, ID: d.ID, Title: d.Title, Step: c.step, Reason: c.reason, Also: c.also}
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

// Catalog lists every document the selection neither loaded whole nor
// briefed, in document order: those not loaded at all, then those loaded
// in part with the outline of their headings, so an agent can ask for the
// rest. A brief carries its own outline.
func (s *Selection) Catalog() (notLoaded, inPart []Entry) {
	notLoaded, inPart = []Entry{}, []Entry{}
	for _, d := range s.Docs {
		if s.Briefed(d) {
			continue
		}
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
