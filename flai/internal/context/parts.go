package context

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/topics"
)

// Parts cuts the pack into parts that each encode to at most limit bytes of
// JSON, as json.Marshal encodes them and the MCP prime tool returns them
// (ADR-0104). The pack's content and order are unchanged; only its delivery
// is paged. Part 1 carries the header and then the conventions and items in
// the pack's order while they fit; each later part carries the story, the
// item, the role, and the title, to say whose pack it is, and the next
// conventions and items. The catalog goes in the last part, or starts a new
// one when it does not fit what the last has left. A part ends before the
// convention or item that would take it over limit; one too large for any
// part is split at line boundaries into chunks, each starting a part. Every
// part carries its number and the number of parts. A pack that fits one part
// is one part. The pack itself is not changed.
func (p *Pack) Parts(limit int) []*Pack {
	whole := *p
	whole.PartNumber, whole.PartCount = 1, 1
	if encoded(&whole) <= limit {
		return []*Pack{&whole}
	}
	// The parts are measured with a count of parts as long in digits as the
	// one they come to, so that writing the count in takes none over limit.
	for count := 1; ; {
		parts := p.cut(limit, count)
		if len(strconv.Itoa(len(parts))) <= len(strconv.Itoa(count)) {
			for _, q := range parts {
				q.PartCount = len(parts)
			}
			return parts
		}
		count = len(parts)
	}
}

// Part is part n of the pack cut into parts of at most limit bytes, the
// first being 1. A number outside the parts is refused with how many there
// are.
func (p *Pack) Part(n, limit int) (*Pack, error) {
	parts := p.Parts(limit)
	switch {
	case len(parts) == 1 && n != 1:
		return nil, fmt.Errorf("prime: part %d of a pack in 1 part; give part 1", n)
	case n < 1 || n > len(parts):
		return nil, fmt.Errorf("prime: part %d of a pack in %d parts; give part 1 to %d", n, len(parts), len(parts))
	}
	return parts[n-1], nil
}

// cutter cuts one pack into parts, measuring each as it fills.
type cutter struct {
	pack  *Pack
	limit int
	count int // the count of parts each part is measured with
	parts []*Pack
	used  int // bytes the last part encodes to
}

// cut is the parts of the pack, each written with count as its count of
// parts.
func (p *Pack) cut(limit, count int) []*Pack {
	c := &cutter{pack: p, limit: limit, count: count}
	c.open()
	place(c, p.Conventions, func(q *Pack) *[]Convention { return &q.Conventions })
	place(c, p.Items, func(q *Pack) *[]Item { return &q.Items })
	c.catalog()
	return c.parts
}

// open starts the next part: the pack's header for the first, the names
// that say whose pack it is for the rest, with no conventions, items, or
// catalog yet.
func (c *cutter) open() {
	p := c.pack
	q := &Pack{Story: p.Story, Item: p.Item, Title: p.Title, Role: p.Role, Topics: []topics.StoryTopic{}, Omitted: []string{}}
	if len(c.parts) == 0 {
		h := *p
		q = &h
	}
	q.Conventions, q.Items, q.Catalog = []Convention{}, []Item{}, Catalog{NotLoaded: []Entry{}, InPart: []Entry{}}
	q.PartNumber, q.PartCount = len(c.parts)+1, c.count
	c.parts = append(c.parts, q)
	c.used = encoded(q)
}

// last is the part being filled.
func (c *cutter) last() *Pack { return c.parts[len(c.parts)-1] }

// add appends e to the last part's list when it fits what the part has
// left, or always when force is set, and says whether it did.
func add[T any](c *cutter, e T, list func(*Pack) *[]T, force bool) bool {
	l := list(c.last())
	n := encoded(e)
	if len(*l) > 0 {
		n++ // the comma before it
	}
	if !force && c.used+n > c.limit {
		return false
	}
	*l = append(*l, e)
	c.used += n
	return true
}

// chunkable is a convention or an item: what Parts splits when it is too
// large for any part.
type chunkable[T any] interface {
	text() string
	chunk(text string, k, n int) T
}

// place puts each element of all into the parts in order: in the last part
// while it fits, else at the start of a new part, else split into chunks.
func place[T chunkable[T]](c *cutter, all []T, list func(*Pack) *[]T) {
	for _, e := range all {
		if add(c, e, list, false) {
			continue
		}
		c.open()
		if add(c, e, list, false) {
			continue
		}
		split(c, e, list)
	}
}

// split puts an element too large for an empty part into consecutive parts
// as chunks of its text, cut at line boundaries, the first in the part just
// opened. Each chunk takes what its part has room for, measured with the
// element's whole length as its count of chunks, which no count exceeds.
func split[T chunkable[T]](c *cutter, e T, list func(*Pack) *[]T) {
	rest := e.text()
	if rest == "" {
		add(c, e, list, true)
		return
	}
	most := len(rest)
	var texts []string
	var parts []*Pack
	for rest != "" {
		if len(texts) > 0 {
			c.open()
		}
		k := len(texts) + 1
		n := cut(rest, c.limit-c.used-encoded(e.chunk("", k, most)))
		texts, parts = append(texts, rest[:n]), append(parts, c.last())
		rest = rest[n:]
		add(c, e.chunk(texts[k-1], k, most), list, true)
	}
	for i, q := range parts {
		(*list(q))[0] = e.chunk(texts[i], i+1, len(texts))
	}
	c.used = encoded(c.last())
}

// cut is how many bytes from the start of s take at most room bytes inside
// a JSON string: whole lines while they fit, or, when not even the first
// line fits, its runes while they fit, and always at least one rune. The
// encoder escapes rune by rune, so the pieces' encoded lengths add up.
func cut(s string, room int) int {
	n, used := 0, 0
	for n < len(s) {
		end := len(s)
		if i := strings.IndexByte(s[n:], '\n'); i >= 0 {
			end = n + i + 1
		}
		w := escaped(s[n:end])
		if used+w > room {
			break
		}
		n, used = end, used+w
	}
	if n > 0 {
		return n
	}
	for n < len(s) {
		_, size := utf8.DecodeRuneInString(s[n:])
		w := escaped(s[n : n+size])
		if n > 0 && used+w > room {
			break
		}
		n, used = n+size, used+w
	}
	return n
}

// catalog puts the pack's catalog in the last part when it fits what that
// part has left, else in a new part, else its entries across new parts.
func (c *cutter) catalog() {
	cat := c.pack.Catalog
	if len(cat.NotLoaded)+len(cat.InPart) == 0 {
		return
	}
	n := encoded(cat) - encoded(Catalog{NotLoaded: []Entry{}, InPart: []Entry{}})
	if c.used+n > c.limit {
		c.open()
	}
	if c.used+n <= c.limit {
		c.last().Catalog = cat
		c.used += n
		return
	}
	for _, e := range cat.NotLoaded {
		entry(c, e, func(q *Pack) *[]Entry { return &q.Catalog.NotLoaded })
	}
	for _, e := range cat.InPart {
		entry(c, e, func(q *Pack) *[]Entry { return &q.Catalog.InPart })
	}
}

// entry puts one catalog entry in the last part, or in a new one when it
// does not fit; an entry too large for any part has one of its own.
func entry(c *cutter, e Entry, list func(*Pack) *[]Entry) {
	if !add(c, e, list, false) {
		c.open()
		add(c, e, list, true)
	}
}

func (c Convention) text() string { return c.Text }

// chunk is the convention holding text as chunk k of n. Its kept and
// left-out sections are on the first chunk only.
func (c Convention) chunk(text string, k, n int) Convention {
	c.Text, c.Chunk, c.Chunks = text, k, n
	if k > 1 {
		c.Kept, c.LeftOut = []Piece{}, []Piece{}
	}
	return c
}

func (it Item) text() string { return it.Text }

// chunk is the item holding text as chunk k of n. The reasons found for it
// after the first are on the first chunk only.
func (it Item) chunk(text string, k, n int) Item {
	it.Text, it.Chunk, it.Chunks = text, k, n
	if k > 1 {
		it.Also = nil
	}
	return it
}

// encoded is how many bytes v encodes to as JSON, with the standard
// encoder's escaping, as the MCP go-sdk encodes a tool's result.
func encoded(v any) int {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err) // a pack holds strings, numbers, and lists, which always encode
	}
	return len(b)
}

// escaped is how many bytes s takes inside a JSON string, its quotes left
// out.
func escaped(s string) int { return encoded(s) - 2 }
