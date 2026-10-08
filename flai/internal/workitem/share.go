package workitem

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// Shares (ADR-0134). The agent of a story in progress that holds a ready story
// on overlap may share the overlapping paths with it, with a split of who
// changes what. The share is kept in the front matter of the two stories'
// conversation under wip/messages, and while it is in force an overlap
// between the two whose narrower path lies inside a shared path does not hold.

// messagesFolder is where conversations are kept under layout.wip; package
// messages, which imports this one, names it too.
const messagesFolder = "messages"

// Share is one sharing of paths by the holding story with the held one, as a
// conversation's front matter keeps it under shares (ADR-0134).
type Share struct {
	Holder string   `yaml:"holder" json:"holder"` // the story in progress that held
	Held   string   `yaml:"held" json:"held"`     // the ready story it held
	Paths  []string `yaml:"paths" json:"paths"`   // what it shares, each inside the two claims' overlap
	Split  string   `yaml:"split" json:"split"`   // who changes what
	By     string   `yaml:"by" json:"by"`         // who shared
	At     string   `yaml:"at" json:"at"`         // when, in TimeFormat
}

// shareFrontMatter is what Shares reads of a conversation's front matter.
type shareFrontMatter struct {
	Status string  `yaml:"status"`
	Shares []Share `yaml:"shares"`
}

// Shares is every share kept in a conversation stored open under
// <layout.wip>/messages, in file order. A conversation that cannot be read
// or parsed is skipped, so that holds never fail on one.
func (r *Repo) Shares() []Share {
	files, err := filepath.Glob(filepath.Join(r.WipDir(), messagesFolder, "MS-*.md"))
	if err != nil {
		return nil
	}
	sort.Strings(files)
	var out []Share
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		fm, _, err := SplitFrontMatter(string(data))
		if err != nil {
			continue
		}
		var c shareFrontMatter
		if err := yaml.Unmarshal([]byte(fm), &c); err != nil || c.Status != "open" {
			continue
		}
		out = append(out, c.Shares...)
	}
	return out
}

// WithShares has h treat an overlap that a share in force clears as no
// overlap (ADR-0134), keeping of shares only those in force, and returns h.
// A share is in force while both its stories exist, are not archived, are
// ready, in progress, or in review, and neither has moved to backlog, done,
// or cancelled at or after it was made; one whose time does not parse is not.
func (h *Holds) WithShares(shares []Share) *Holds {
	h.shares = nil
	for _, s := range shares {
		s.Holder, s.Held = CanonicalID(s.Holder), CanonicalID(s.Held)
		if h.inForce(s) {
			h.shares = append(h.shares, s)
		}
	}
	return h
}

// InForce says whether share s is in force, as WithShares judges it, with
// the stories h knows: flai serve asks it of the shares it names in the
// prompt of a story started on one.
func (h *Holds) InForce(s Share) bool {
	s.Holder, s.Held = CanonicalID(s.Holder), CanonicalID(s.Held)
	return h.inForce(s)
}

// inForce says whether share s still holds, judged with the stories h knows.
func (h *Holds) inForce(s Share) bool {
	at, err := time.Parse(TimeFormat, s.At)
	if err != nil || s.Holder == s.Held {
		return false
	}
	for _, id := range []string{s.Holder, s.Held} {
		it := h.story(id)
		if it == nil || it.Archived || (it.Status != Ready && it.Status != InProgress && it.Status != Review) {
			return false
		}
		for _, tr := range it.Transitions {
			if tr.To != Backlog && tr.To != Done && tr.To != Cancelled {
				continue
			}
			if when, err := time.Parse(TimeFormat, tr.At); err == nil && !when.Before(at) {
				return false
			}
		}
	}
	return true
}

// sharedBy says whether a share in force by holder with held covers path p:
// p is one of its paths or lies below one.
func (h *Holds) sharedBy(holder, held, p string) bool {
	for _, s := range h.shares {
		if s.Holder != holder || s.Held != held {
			continue
		}
		for _, sp := range s.Paths {
			if e := h.path(sp); e != "" && (p == e || strings.HasPrefix(p, e+"/")) {
				return true
			}
		}
	}
	return false
}

// Overlap is one open story that holds a ready story on overlap, and the
// paths of that hold.
type Overlap struct {
	By    string   `json:"by"`
	Paths []string `json:"paths"`
}

// OverlapsBy is, for a story held on overlap alone (held (overlap), with no
// after: and no story without touches holding it), each open story that holds
// it, in ID order, with the paths it is held on that no share clears: for
// each holding pair of entries the narrower, the path both stories change,
// in the ready story's claim order, without duplicates. Otherwise nil.
func (h *Holds) OverlapsBy(story *Item) []Overlap {
	hold := h.Of(story)
	if hold == nil || hold.Code != HoldOverlap {
		return nil
	}
	claim := h.Claim(story)
	var out []Overlap
	for _, o := range h.open {
		if o.id == story.ID {
			continue
		}
		if len(o.paths) == 0 {
			return nil // a story with no touches holds it too: not overlap alone
		}
		var paths []string
		for _, pr := range h.holding(story.ID, claim, o) {
			if p := narrower(pr[0], pr[1]); !contains(paths, p) {
				paths = append(paths, p)
			}
		}
		if len(paths) > 0 {
			out = append(out, Overlap{By: o.id, Paths: paths})
		}
	}
	return out
}
