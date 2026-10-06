package context

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/topics"
)

// largePack is a pack over the fixture documents whose conventions are
// files of the given texts and whose briefs take it over its budget, as the
// packs I-0068 records are.
func largePack(t *testing.T, texts ...string) *Pack {
	t.Helper()
	set := &conventions.Set{Dir: "/r/design/conventions"}
	for i, text := range texts {
		name := fmt.Sprintf("c%d.md", i+1)
		set.Files = append(set.Files, conventions.File{Path: "design/conventions/" + name, Name: name, Title: name, Order: 10 * (i + 1), Raw: text})
	}
	names := []string{"cli", "go", "dashboard"}
	var st []topics.StoryTopic
	for _, n := range names {
		st = append(st, topics.StoryTopic{Topic: n, Sources: []topics.Source{{Kind: topics.FromOwn, Item: "S-0001"}}})
	}
	build := func(budget int, query string) *Pack {
		p, err := Build("/r", "S-0001", "A story", st, set, "", budget)
		if err != nil {
			t.Fatal(err)
		}
		p.AddDesign(Design(fixture(t).Docs, names, nil, BriefOver(budget)), query)
		return p
	}
	base := build(DefaultBudget, "")
	return build(base.Size.Bytes-100, "herons")
}

// rules is a convention of about n bytes of rules, with what the encoder
// escapes on every line.
func rules(title string, n int) string {
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	for i := 0; b.Len() < n; i++ {
		fmt.Fprintf(&b, "- Rule %d: a <section> & its \"heading\" stay whole.\n", i)
	}
	return b.String()
}

func size(t *testing.T, v any) int {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return len(b)
}

// rejoin is the conventions and items of the parts in order, each chunked
// one put back together from its chunks.
func rejoin(t *testing.T, parts []*Pack) ([]Convention, []Item) {
	t.Helper()
	convs, items := []Convention{}, []Item{}
	for _, q := range parts {
		for _, c := range q.Conventions {
			if c.Chunk > 1 {
				last := &convs[len(convs)-1]
				if last.Path != c.Path || c.Chunk != last.Chunk+1 {
					t.Fatalf("chunk %d of %s follows chunk %d of %s", c.Chunk, c.Path, last.Chunk, last.Path)
				}
				last.Text += c.Text
				last.Chunk = c.Chunk
				continue
			}
			convs = append(convs, c)
		}
		for _, it := range q.Items {
			if it.Chunk > 1 {
				last := &items[len(items)-1]
				last.Text += it.Text
				last.Chunk = it.Chunk
				continue
			}
			items = append(items, it)
		}
	}
	for i := range convs {
		if convs[i].Chunk != convs[i].Chunks {
			t.Errorf("%s ends at chunk %d of %d", convs[i].Path, convs[i].Chunk, convs[i].Chunks)
		}
		convs[i].Chunk, convs[i].Chunks = 0, 0
	}
	for i := range items {
		items[i].Chunk, items[i].Chunks = 0, 0
	}
	return convs, items
}

// checkParts says what every set of parts must be: each within limit,
// numbered in order with the count, the header in part 1 alone, the catalog
// in the last, and the conventions and items of the pack in its order.
func checkParts(t *testing.T, p *Pack, parts []*Pack, limit int) {
	t.Helper()
	for i, q := range parts {
		if n := size(t, q); n > limit {
			t.Errorf("part %d encodes to %d bytes, over %d", i+1, n, limit)
		}
		if q.PartNumber != i+1 || q.PartCount != len(parts) {
			t.Errorf("part %d says part %d of %d, of %d", i+1, q.PartNumber, q.PartCount, len(parts))
		}
		if q.Story != p.Story || q.Title != p.Title || q.Conventions == nil || q.Items == nil {
			t.Errorf("part %d: story %q, title %q, conventions %v, items %v", i+1, q.Story, q.Title, q.Conventions, q.Items)
		}
		if i > 0 && (q.Budget != 0 || q.Exceeded != "" || len(q.Omitted) != 0 || q.README != nil) {
			t.Errorf("part %d carries the header: %+v", i+1, q)
		}
	}
	h := *parts[0]
	h.PartNumber, h.PartCount = 0, 0
	h.Conventions, h.Items, h.Catalog = p.Conventions, p.Items, p.Catalog
	if !reflect.DeepEqual(&h, p) {
		t.Errorf("part 1 does not carry the pack's header")
	}
	var catalog Catalog
	for _, q := range parts {
		catalog.NotLoaded = append(catalog.NotLoaded, q.Catalog.NotLoaded...)
		catalog.InPart = append(catalog.InPart, q.Catalog.InPart...)
	}
	if !reflect.DeepEqual(catalog.NotLoaded, p.Catalog.NotLoaded) || len(catalog.InPart) != len(p.Catalog.InPart) {
		t.Errorf("the parts' catalog %+v is not the pack's %+v", catalog, p.Catalog)
	}
	if last := parts[len(parts)-1]; len(p.Catalog.NotLoaded) > 0 && len(last.Catalog.NotLoaded) == 0 {
		t.Errorf("the catalog is not in the last part")
	}
	convs, items := rejoin(t, parts)
	if !reflect.DeepEqual(convs, p.Conventions) {
		t.Errorf("the parts' conventions are not the pack's")
	}
	if !reflect.DeepEqual(items, p.Items) {
		t.Errorf("the parts' items are not the pack's")
	}
}

// I-0068: a story's pack, its briefs over the budget, encodes to more than
// the 50,000 characters over which Claude Code saves a tool result to a
// file. Its parts each fit PartLimit and together are the pack.
func TestAPackOverTheToolResultLimitIsCutIntoPartsThatFit(t *testing.T) {
	p := largePack(t, rules("One", 15000), rules("Two", 15000), rules("Three", 15000), rules("Four", 15000))
	if p.Exceeded != ExceededBriefs || len(p.Catalog.NotLoaded) == 0 {
		t.Fatalf("exceeded %q, catalog %+v", p.Exceeded, p.Catalog)
	}
	before, _ := json.Marshal(p)
	if len(before) <= 50000 {
		t.Fatalf("the pack encodes to %d bytes, not over 50,000", len(before))
	}
	parts := p.Parts(PartLimit)
	if len(parts) < 2 {
		t.Fatalf("%d parts", len(parts))
	}
	checkParts(t, p, parts, PartLimit)
	for _, q := range parts {
		for _, c := range q.Conventions {
			if c.Chunk != 0 {
				t.Errorf("%s fits a part of its own but is chunked", c.Path)
			}
		}
	}
	if after, _ := json.Marshal(p); string(after) != string(before) {
		t.Error("Parts changed the pack")
	}
	if q, err := p.Part(len(parts), PartLimit); err != nil || !reflect.DeepEqual(q, parts[len(parts)-1]) {
		t.Errorf("Part(%d): %v", len(parts), err)
	}
	if _, err := p.Part(len(parts)+1, PartLimit); err == nil || !strings.Contains(err.Error(), fmt.Sprintf("in %d parts; give part 1 to %d", len(parts), len(parts))) {
		t.Errorf("part past the last: %v", err)
	}
	if _, err := p.Part(0, PartLimit); err == nil {
		t.Error("part 0 taken")
	}
}

func TestAConventionTooLargeForAPartIsSplitIntoChunks(t *testing.T) {
	long := "# Long\n\n" + strings.Repeat("héron <à> &   marée ", 600) + "\n" + rules("Tail", 3000)
	p := largePack(t, rules("Small", 2000), rules("Large", 12000), long, rules("After", 1000))
	const limit = 5000
	before, _ := json.Marshal(p)
	parts := p.Parts(limit)
	checkParts(t, p, parts, limit)
	chunks := map[string]int{}
	for _, q := range parts {
		for i, c := range q.Conventions {
			if c.Chunk == 0 {
				continue
			}
			chunks[c.Path]++
			if i != 0 {
				t.Errorf("chunk %d of %s does not start its part", c.Chunk, c.Path)
			}
			if !utf8.ValidString(c.Text) {
				t.Errorf("chunk %d of %s is cut inside a rune", c.Chunk, c.Path)
			}
			if c.Chunk > 1 && len(c.Kept) != 0 {
				t.Errorf("chunk %d of %s carries its kept sections", c.Chunk, c.Path)
			}
		}
	}
	if chunks["design/conventions/c2.md"] < 3 || chunks["design/conventions/c3.md"] < 3 || chunks["design/conventions/c1.md"] != 0 || chunks["design/conventions/c4.md"] != 0 {
		t.Errorf("chunks %v", chunks)
	}
	if after, _ := json.Marshal(p); string(after) != string(before) {
		t.Error("Parts changed the pack")
	}
}

func TestACatalogTooLargeForAPartIsSplitAcrossParts(t *testing.T) {
	// The conventions alone exceed a budget of one byte, so every document
	// is in the catalog.
	p, _ := budgetPack(t, 1, []string{"cli"}, nil, "")
	empty := size(t, &Pack{Story: p.Story, Title: p.Title, PartNumber: 9, PartCount: 9, Topics: []topics.StoryTopic{}, Conventions: []Convention{}, Omitted: []string{}, Items: []Item{}, Catalog: Catalog{NotLoaded: []Entry{}, InPart: []Entry{}}})
	limit := empty + size(t, p.Catalog)/2
	if size(t, p.Catalog)+empty <= limit {
		t.Fatalf("the catalog fits a part of %d", limit)
	}
	parts := p.Parts(limit)
	checkParts(t, p, parts, limit)
	n := 0
	for _, q := range parts {
		if len(q.Catalog.NotLoaded) > 0 {
			n++
		}
	}
	if n < 2 || len(parts[len(parts)-1].Catalog.NotLoaded) == 0 {
		t.Errorf("the catalog is in %d parts", n)
	}
}

func TestASmallPackIsOnePart(t *testing.T) {
	p, _ := budgetPack(t, DefaultBudget, []string{"cli"}, nil, "")
	parts := p.Parts(PartLimit)
	if len(parts) != 1 || parts[0].PartNumber != 1 || parts[0].PartCount != 1 {
		t.Fatalf("%d parts, the first %d of %d", len(parts), parts[0].PartNumber, parts[0].PartCount)
	}
	whole := *parts[0]
	whole.PartNumber, whole.PartCount = 0, 0
	if !reflect.DeepEqual(&whole, p) || p.PartNumber != 0 {
		t.Error("the one part is not the whole pack")
	}
	b, _ := json.Marshal(parts[0])
	if !strings.Contains(string(b), `"part":1,"parts":1,`) {
		t.Errorf("no part and parts in %s", b)
	}
	if b, _ := json.Marshal(p); strings.Contains(string(b), `"part"`) || strings.Contains(string(b), `"chunk"`) {
		t.Errorf("the whole pack names a part: %s", b)
	}
	if q, err := p.Part(1, PartLimit); err != nil || q.PartCount != 1 {
		t.Errorf("part 1: %v", err)
	}
	for _, n := range []int{0, 2, -1} {
		if _, err := p.Part(n, PartLimit); err == nil || err.Error() != fmt.Sprintf("prime: part %d of a pack in 1 part; give part 1", n) {
			t.Errorf("part %d: %v", n, err)
		}
	}
}

func TestCutKeepsLinesWholeAndRunesWhole(t *testing.T) {
	for _, c := range []struct {
		s    string
		room int
		want int
	}{
		{"ab\ncd\n", 6, 3},     // "ab\n" is 4 escaped, "cd\n" would take it to 8
		{"ab\ncd\n", 8, 6},     // both lines
		{"abcdef\ng\n", 3, 3},  // no whole line fits: runes
		{"<<\n", 6, 1},         // < escapes to six bytes
		{"<<\n", 2, 1},         // always one rune
		{"éé", 2, 2},           // é is one rune of two bytes
		{"a b", 4, 1},          // U+2028 escapes to six bytes
		{"\xffa\n", 6, 1},      // an invalid byte encodes as �
		{"one line\n", 100, 9}, // all of it
	} {
		if got := cut(c.s, c.room); got != c.want {
			t.Errorf("cut(%q, %d) = %d, want %d", c.s, c.room, got, c.want)
		}
	}
}
