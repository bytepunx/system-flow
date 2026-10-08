package workitem

import (
	"os"
	"path/filepath"
	"regexp"
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
// The board names, on a held card, the conversation in which flai asked about
// the hold, and on a card started on a share, the share (S-0338).

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
	// Conversation is the ID of the conversation that keeps the share, which
	// Repo.Conversations fills from its id for the board (S-0338); it is never
	// written to the front matter.
	Conversation string `yaml:"-" json:"conversation,omitempty"`
}

// Ask is one message by flai asking the agent of a story in progress about
// its hold on a ready story, as an open conversation between the two keeps
// it (S-0338).
type Ask struct {
	Conversation string // the conversation's ID
	Held         string // the ready story the message was written for
	Holder       string // the conversation's other story, whose agent was asked
	At           string // when, in TimeFormat
}

// askAuthor is the author of the message that asks about a hold; package
// messages names it too, as AskAuthor.
const askAuthor = "flai"

// askHeading is the heading of an entry by flai for a story, as package
// messages writes an entry's: "### <time> <author> <story>".
var askHeading = regexp.MustCompile(`(?m)^### (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z) ` + askAuthor + ` (S-\d{3,})$`)

// conversationFrontMatter is what Repo.Conversations reads of a
// conversation's front matter.
type conversationFrontMatter struct {
	ID     string  `yaml:"id"`
	From   string  `yaml:"from"`
	To     string  `yaml:"to"`
	Status string  `yaml:"status"`
	Shares []Share `yaml:"shares"`
}

// Shares is every share kept in a conversation stored open under
// <layout.wip>/messages, in file order, as Conversations reads them.
func (r *Repo) Shares() []Share {
	shares, _ := r.Conversations()
	return shares
}

// Conversations reads the conversations stored open under
// <layout.wip>/messages once, in file order, for the shares each keeps, each
// naming the conversation, and for each message by flai asking about a hold
// (S-0338): one written for one of the conversation's two stories, the held,
// whose holder is the other. A conversation that cannot be read or parsed is
// skipped, so that holds never fail on one.
func (r *Repo) Conversations() ([]Share, []Ask) {
	files, err := filepath.Glob(filepath.Join(r.WipDir(), messagesFolder, "MS-*.md"))
	if err != nil {
		return nil, nil
	}
	sort.Strings(files)
	var shares []Share
	var asks []Ask
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		fm, body, err := SplitFrontMatter(string(data))
		if err != nil {
			continue
		}
		var c conversationFrontMatter
		if err := yaml.Unmarshal([]byte(fm), &c); err != nil || c.Status != "open" {
			continue
		}
		for _, s := range c.Shares {
			s.Conversation = c.ID
			shares = append(shares, s)
		}
		from, to := CanonicalID(c.From), CanonicalID(c.To)
		for _, m := range askHeading.FindAllStringSubmatch(body, -1) {
			held, holder := CanonicalID(m[2]), ""
			switch held {
			case from:
				holder = to
			case to:
				holder = from
			}
			if holder != "" && holder != held {
				asks = append(asks, Ask{Conversation: c.ID, Held: held, Holder: holder, At: m[1]})
			}
		}
	}
	return shares, asks
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
	return h.overlapsBy(story, h.of(story))
}

// overlapsBy is OverlapsBy with hold, story's as of judges it.
func (h *Holds) overlapsBy(story *Item, hold *Hold) []Overlap {
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

// SharesWith is the shares in force whose held story is id, in the order they
// were given: what the board shows on the card of a story started on a share
// (S-0338).
func (h *Holds) SharesWith(id string) []Share {
	id = CanonicalID(id)
	var out []Share
	for _, s := range h.shares {
		if s.Held == id {
			out = append(out, s)
		}
	}
	return out
}

// WithAsks gives h the messages by flai that asked about a hold, which a hold
// on overlap names (S-0338), and returns h.
func (h *Holds) WithAsks(asks []Ask) *Holds {
	h.asks = nil
	for _, a := range asks {
		a.Held, a.Holder = CanonicalID(a.Held), CanonicalID(a.Holder)
		h.asks = append(h.asks, a)
	}
	return h
}

// asked is, for each story in overlaps that holds story, in their order, the
// newest message by flai that asked its agent about the hold since story
// last entered ready, as package messages judges an ask; a holder not asked
// is left out, and nil when none was.
func (h *Holds) asked(story *Item, overlaps []Overlap) []HoldAsk {
	if len(h.asks) == 0 {
		return nil
	}
	since := readySince(story)
	var out []HoldAsk
	for _, o := range overlaps {
		var newest *Ask
		var newestAt time.Time
		for i, a := range h.asks {
			if a.Held != story.ID || a.Holder != o.By {
				continue
			}
			at, err := time.Parse(TimeFormat, a.At)
			if err != nil || at.Before(since) || (newest != nil && at.Before(newestAt)) {
				continue
			}
			newest, newestAt = &h.asks[i], at
		}
		if newest != nil {
			out = append(out, HoldAsk{By: o.By, Conversation: newest.Conversation, At: newest.At})
		}
	}
	return out
}

// readySince is when story last entered ready, else when it was created; a
// time that does not parse reads as the zero time.
func readySince(story *Item) time.Time {
	for i := len(story.Transitions) - 1; i >= 0; i-- {
		if story.Transitions[i].To == Ready {
			t, _ := time.Parse(TimeFormat, story.Transitions[i].At)
			return t
		}
	}
	t, _ := time.Parse(TimeFormat, story.Created)
	return t
}
