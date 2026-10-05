package mdlint

import (
	"regexp"
	"strconv"
	"strings"
)

// The block structure of a document, as far as the rules need it: which
// lines are front matter, code, HTML, headings, list items, tables, and
// paragraph text, in whatever list items and blockquotes hold them, and how
// paragraphs and lists group them. It follows CommonMark closely enough for
// what flai and its agents write; a line it cannot place is left to rules
// that do not report on it.

type kind int

const (
	kText kind = iota
	kBlank
	kFront
	kFenceOpen
	kFenceClose
	kCode
	kIndentCode
	kHTML
	kHeading
	kSetext
	kHR
	kTable
)

type line struct {
	raw   string
	kind  kind
	start int // byte offset of the content in raw: after indentation, list markers, and quote markers
	lead  int // byte offset where the containers' prefix ends: a code span keeps the indentation past it
	para  *para
	fence *fence
}

type para struct {
	depth   int
	lines   []int
	heading *heading // set when the paragraph became a setext heading
}

type fence struct {
	char   byte
	length int
	listed bool // inside a list item
	info   string
	open   int
	close  int // -1 while unclosed
}

type heading struct {
	level  int
	text   string
	start  int // first line index
	end    int // last line index (the underline for setext)
	textAt int // line index of the last text line
	setext bool
	closed bool // an ATX heading with a closing sequence
	raw    string
}

type listItem struct {
	line          int
	ordered       bool
	marker        byte // -, *, + or the ordered delimiter . )
	value         int
	contentIndent int
	at            int // byte offset of the marker in the line
	quote         int // byte offset past the last quote marker before it on the line, and the space after that
	list          *list
}

type list struct {
	ordered bool
	marker  byte
	items   []*listItem
	parent  *listItem // the item the list is nested in; nil at the top level or directly in a quote
}

type doc struct {
	lines      []*line
	front      []string // front matter lines, without the fences
	content    int      // index of the first line after front matter
	headings   []*heading
	lists      []*list
	fences     []*fence
	paras      []*para
	tableCells [][3]int // line, start, end of each table cell
	paraSeg    map[*para]*segment
	paraInfo   map[*para]segInfo
}

var (
	fenceRe   = regexp.MustCompile("^(`{3,}|~{3,})(.*)$")
	atxRe     = regexp.MustCompile(`^(#{1,6})(?:[ \t]+|$)`)
	hrRe      = regexp.MustCompile(`^(?:(?:\*[ \t]*){3,}|(?:-[ \t]*){3,}|(?:_[ \t]*){3,})$`)
	bulletRe  = regexp.MustCompile(`^([-*+])([ \t]+|$)`)
	orderedRe = regexp.MustCompile(`^(\d{1,9})([.)])([ \t]+|$)`)
	setextRe  = regexp.MustCompile(`^(=+|-+)[ \t]*$`)
	htmlRe    = regexp.MustCompile(`^</?[A-Za-z][A-Za-z0-9-]*(?:[\s/>]|$)`)
	delimRow  = regexp.MustCompile(`^\|?[ \t]*:?-+:?[ \t]*(?:\|[ \t]*:?-+:?[ \t]*)*\|?[ \t]*$`)
	closeATX  = regexp.MustCompile(`(?:^|[ \t]+)#+[ \t]*$`)
)

// columnsFrom returns the indentation of s in columns, tabs to the next stop
// of four, and the byte offset where it ends. s starts at column from, where
// tab stops fall as they do on the whole line.
func columnsFrom(s string, from int) (int, int) {
	col := from
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ':
			col++
		case '\t':
			col += 4 - col%4
		default:
			return col - from, i
		}
	}
	return col - from, len(s)
}

// skipColumns returns the byte offset after n columns of indentation.
func skipColumns(s string, n int) int {
	col := 0
	for i := 0; i < len(s); i++ {
		if col >= n {
			return i
		}
		switch s[i] {
		case ' ':
			col++
		case '\t':
			col += 4 - col%4
		default:
			return i
		}
	}
	return len(s)
}

func isBlank(s string) bool { return strings.TrimSpace(s) == "" }

// parse splits content into lines, as markdownlint does (a final newline
// leaves an empty last line), and works out their structure.
func parse(content string) *doc {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	raw := strings.Split(content, "\n")
	d := &doc{}
	for _, r := range raw {
		d.lines = append(d.lines, &line{raw: r})
	}
	if len(raw) > 0 && strings.TrimRight(raw[0], " \t") == "---" {
		for i := 1; i < len(raw); i++ {
			if t := strings.TrimRight(raw[i], " \t"); t == "---" || t == "..." {
				for j := 0; j <= i; j++ {
					d.lines[j].kind = kFront
				}
				d.front = raw[1:i]
				d.content = i + 1
				break
			}
		}
	}
	d.blocks()
	return d
}

type parser struct {
	d       *doc
	stack   []*listItem // the open containers, outermost first: list items, and nil for a block quote
	lists   []*list     // the open list at each depth
	para    *para
	fence   *fence
	htmlEnd string
	html    bool
	table   bool
	quote   int // byte offset past the last quote marker taken from the line
}

func (d *doc) blocks() {
	p := &parser{d: d}
	for i := d.content; i < len(d.lines); i++ {
		p.line(i)
	}
	// blank lines between two indented code lines are code
	for i := d.content; i < len(d.lines); i++ {
		if d.lines[i].kind != kIndentCode {
			continue
		}
		j := i + 1
		for j < len(d.lines) && d.lines[j].kind == kBlank {
			j++
		}
		if j < len(d.lines) && j > i+1 && d.lines[j].kind == kIndentCode {
			for k := i + 1; k < j; k++ {
				d.lines[k].kind = kIndentCode
			}
		}
	}
}

func (p *parser) endLists(depth int) {
	if len(p.lists) > depth {
		p.lists = p.lists[:depth]
	}
}

// startsBlock reports whether content, at most three columns in, begins a
// block that interrupts a paragraph, so the line cannot be a lazy
// continuation of one.
func startsBlock(c string) bool {
	if fenceRe.MatchString(c) || atxRe.MatchString(c) || hrRe.MatchString(c) || strings.HasPrefix(c, ">") || strings.HasPrefix(c, "<!--") {
		return true
	}
	if m := bulletRe.FindStringSubmatch(c); m != nil && !isBlank(c[len(m[0]):]) {
		return true
	}
	if m := orderedRe.FindStringSubmatch(c); m != nil && m[1] == "1" && !isBlank(c[len(m[0]):]) {
		return true
	}
	return false
}

// prefix is how much of a line continues the open containers: the first n
// of them, up to byte pos at column col. The content of the last starts at
// column base, and quote is the byte offset past the last quote marker.
type prefix struct{ n, pos, col, base, quote int }

// match takes the open containers' prefixes from the start of raw: a
// quote's marker, at most three columns in, and a list item's indentation,
// which a blank line also gives.
func (p *parser) match(raw string) prefix {
	var m prefix
	for _, it := range p.stack {
		n, k := columnsFrom(raw[m.pos:], m.col)
		if it == nil {
			at := m.pos + k
			if n > 3 || at == len(raw) || raw[at] != '>' {
				break
			}
			m.pos, m.col, m.base = quoteMarker(raw, at, m.col+n)
			m.quote = m.pos
		} else {
			if !isBlank(raw[m.pos:]) && m.col+n < it.contentIndent {
				break
			}
			m.pos, m.col = skipTo(raw, m.pos, m.col, it.contentIndent)
			m.base = it.contentIndent
		}
		m.n++
	}
	return m
}

// quoteMarker takes the > at byte at, column col, and the space or tab after
// it. It returns the byte offset and column past them, and the column the
// quote's content starts at: a tab gives the marker one column and leaves
// the rest as indentation.
func quoteMarker(raw string, at, col int) (pos, next, base int) {
	pos, next = at+1, col+1
	base = next
	if pos < len(raw) && (raw[pos] == ' ' || raw[pos] == '\t') {
		base++
		_, next = skipTo(raw, pos, next, base)
		pos++
	}
	return pos, next, base
}

// skipTo moves from byte pos, at column col, over indentation up to column
// target, taking whole a tab that runs past it.
func skipTo(s string, pos, col, target int) (int, int) {
	for ; pos < len(s) && col < target; pos++ {
		switch s[pos] {
		case ' ':
			col++
		case '\t':
			col += 4 - col%4
		default:
			return pos, col
		}
	}
	return pos, col
}

func (p *parser) line(i int) {
	d := p.d
	ln := d.lines[i]
	raw := ln.raw
	m := p.match(raw)
	all := m.n == len(p.stack)
	p.quote = m.quote
	n, k := columnsFrom(raw[m.pos:], m.col)
	ind, off := m.col+n, m.pos+k
	c := raw[off:]

	if f := p.fence; f != nil {
		if all {
			t := strings.TrimRight(c, " \t")
			if ind-m.base < 4 && len(t) >= f.length && strings.Count(t, string(f.char)) == len(t) {
				ln.kind, ln.fence, f.close, p.fence = kFenceClose, f, i, nil
				return
			}
			ln.kind, ln.fence, ln.start = kCode, f, off
			return
		}
		p.fence = nil // a container closed, and the fence with it
	}
	if p.htmlEnd != "" {
		if all {
			ln.kind = kHTML
			if strings.Contains(raw, p.htmlEnd) {
				p.htmlEnd = ""
			}
			return
		}
		p.htmlEnd = ""
	}
	if isBlank(c) {
		// a blank line ends the quotes it does not mark
		p.stack = p.stack[:m.n]
		p.rest(i, off, ind, m.base)
		return
	}
	if p.html {
		if all {
			ln.kind = kHTML
			return
		}
		p.html = false
	}
	// any list item closes the item, as a sibling or an ancestor's; past
	// three columns the line cannot start a block
	if !all && p.para != nil && (ind-m.base >= 4 || !startsBlock(c) && !orderedRe.MatchString(c)) {
		ln.lead = m.pos // the containers the line still matches take their prefix
		p.text(i, off)
		return
	}
	if !all {
		p.stack = p.stack[:m.n]
		p.table = false
	}
	if p.para != nil && p.para.depth > len(p.stack) {
		p.para = nil
	}
	ln.lead = m.pos
	p.rest(i, off, ind, m.base)
}

// rest places what follows the containers' prefixes on a line: content at
// byte off, ind columns in, in a container whose content starts at column
// base.
func (p *parser) rest(i, off, ind, base int) {
	ln := p.d.lines[i]
	c := ln.raw[off:]
	depth := len(p.stack)
	if isBlank(c) {
		ln.kind = kBlank
		p.para, p.html, p.table = nil, false, false
		return
	}
	if p.table && strings.Contains(c, "|") {
		ln.kind = kTable
		p.cells(i, off)
		return
	}
	p.table = false
	if ind-base >= 4 {
		if p.para != nil {
			p.text(i, off)
			return
		}
		ln.kind, ln.start = kIndentCode, off
		p.endLists(depth)
		return
	}
	p.block(i, off, depth, ind)
}

// block places a line that is not code, a lazy line, or inside a fence:
// content starts at byte off, ind columns in, at the given depth.
func (p *parser) block(i, off, depth, ind int) {
	d := p.d
	ln := d.lines[i]
	c := ln.raw[off:]
	switch {
	case strings.HasPrefix(c, ">"):
		p.para = nil
		p.endLists(depth)
		p.stack = append(p.stack, nil)
		pos, col, base := quoteMarker(ln.raw, off, ind)
		p.quote = pos
		n, k := columnsFrom(ln.raw[pos:], col)
		p.rest(i, pos+k, col+n, base)
	case fenceRe.MatchString(c) && (c[0] != '`' || !strings.Contains(fenceRe.FindStringSubmatch(c)[2], "`")):
		m := fenceRe.FindStringSubmatch(c)
		f := &fence{char: c[0], length: len(m[1]), info: strings.TrimSpace(m[2]), open: i, close: -1}
		ln.start = off
		for _, it := range p.stack {
			f.listed = f.listed || it != nil
		}
		d.fences = append(d.fences, f)
		ln.kind, ln.fence = kFenceOpen, f
		p.fence, p.para = f, nil
		p.endLists(depth)
	case atxRe.MatchString(c):
		ln.kind = kHeading
		p.atx(i, off)
		p.para = nil
		p.endLists(depth)
	case p.para != nil && p.para.depth == depth && setextRe.MatchString(c) && len(p.para.lines) > 0 && d.lines[p.para.lines[0]].kind == kText:
		pr := p.para
		h := &heading{level: 1, start: pr.lines[0], end: i, textAt: pr.lines[len(pr.lines)-1], setext: true}
		if c[0] == '-' {
			h.level = 2
		}
		var parts []string
		for _, li := range pr.lines {
			parts = append(parts, strings.TrimSpace(d.lines[li].raw[d.lines[li].start:]))
		}
		h.text = strings.Join(parts, "\n")
		h.raw = h.text
		pr.heading = h
		d.headings = append(d.headings, h)
		ln.kind = kSetext
		p.para = nil
	case hrRe.MatchString(c):
		ln.kind = kHR
		p.para = nil
		p.endLists(depth)
	case p.item(i, off, depth, ind):
	case strings.HasPrefix(c, "<!--"):
		ln.kind = kHTML
		if !strings.Contains(c[4:], "-->") {
			p.htmlEnd = "-->"
		}
		p.para = nil
		p.endLists(depth)
	case p.para == nil && htmlRe.MatchString(c):
		ln.kind = kHTML
		p.html = true
		p.endLists(depth)
	case strings.Contains(c, "|") && i+1 < len(d.lines) && strings.Contains(d.lines[i+1].raw, "-") && delimRow.MatchString(strings.TrimSpace(d.lines[i+1].raw[p.match(d.lines[i+1].raw).pos:])):
		ln.kind = kTable
		p.table, p.para = true, nil
		p.cells(i, off)
		p.endLists(depth)
	default:
		if p.para == nil || p.para.depth != depth {
			p.endLists(depth)
		}
		p.text(i, off)
	}
}

// item starts a list item when the content is one, with any items nested on
// the same line, and places what follows the markers.
func (p *parser) item(i, off, depth, ind int) bool {
	d := p.d
	c := d.lines[i].raw[off:]
	it := &listItem{line: i, at: off, quote: p.quote}
	var mlen int
	if m := bulletRe.FindStringSubmatch(c); m != nil {
		it.marker, mlen = m[1][0], len(m[1])
	} else if m := orderedRe.FindStringSubmatch(c); m != nil {
		it.ordered, it.marker, mlen = true, m[2][0], len(m[1])+1
		it.value, _ = strconv.Atoi(m[1])
	} else {
		return false
	}
	rest := c[mlen:]
	if p.para != nil && p.para.depth == depth && (isBlank(rest) || (it.ordered && it.value != 1)) {
		return false // cannot interrupt a paragraph
	}
	sp, ws := columnsFrom(rest, ind+mlen)
	if sp >= 5 || isBlank(rest) {
		sp, ws = 1, skipColumns(rest, 1)
	}
	it.contentIndent = ind + mlen + sp
	if len(p.lists) > depth && p.lists[depth] != nil && p.lists[depth].ordered == it.ordered && p.lists[depth].marker == it.marker {
		it.list = p.lists[depth]
		it.list.items = append(it.list.items, it)
	} else {
		l := &list{ordered: it.ordered, marker: it.marker, items: []*listItem{it}}
		if depth > 0 && depth <= len(p.stack) {
			l.parent = p.stack[depth-1]
		}
		it.list = l
		d.lists = append(d.lists, l)
		p.endLists(depth)
		for len(p.lists) <= depth {
			p.lists = append(p.lists, nil)
		}
		p.lists[depth] = l
	}
	p.endLists(depth + 1)
	p.stack = append(p.stack, it)
	p.para = nil
	if isBlank(rest) {
		d.lines[i].kind = kText
		d.lines[i].start = len(d.lines[i].raw)
		return true
	}
	at := off + mlen + ws
	inner := d.lines[i].raw[at:]
	switch {
	case bulletRe.MatchString(inner) || orderedRe.MatchString(inner):
		if !p.item(i, at, depth+1, it.contentIndent) {
			p.text(i, at)
		}
	case fenceRe.MatchString(inner) || atxRe.MatchString(inner) || hrRe.MatchString(inner) || strings.HasPrefix(inner, ">"):
		p.block(i, at, depth+1, it.contentIndent)
	default:
		p.text(i, at)
	}
	return true
}

func (p *parser) text(i, off int) {
	ln := p.d.lines[i]
	ln.kind, ln.start = kText, off
	if p.para == nil {
		p.para = &para{depth: len(p.stack)}
		p.d.paras = append(p.d.paras, p.para)
	}
	p.para.lines = append(p.para.lines, i)
	ln.para = p.para
}

func (p *parser) atx(i, off int) {
	d := p.d
	c := d.lines[i].raw[off:]
	m := atxRe.FindString(c)
	level := strings.Count(m, "#")
	text := c[len(m):]
	closed := false
	if loc := closeATX.FindStringIndex(text); loc != nil {
		text, closed = text[:loc[0]], true
	}
	h := &heading{level: level, text: strings.TrimSpace(text), start: i, end: i, textAt: i, closed: closed, raw: c}
	d.lines[i].start = off + len(m)
	d.headings = append(d.headings, h)
}

// cells records the cells of a table row, for the inline rules.
func (p *parser) cells(i, off int) {
	raw := p.d.lines[i].raw
	if delimRow.MatchString(strings.TrimSpace(raw[off:])) {
		return
	}
	s, e := off, len(strings.TrimRight(raw, " \t"))
	if s < e && raw[s] == '|' {
		s++
	}
	if e > s && raw[e-1] == '|' && (e < 2 || raw[e-2] != '\\') {
		e--
	}
	from := s
	for j := s; j <= e; j++ {
		if j == e || (raw[j] == '|' && raw[j-1] != '\\') {
			p.d.tableCells = append(p.d.tableCells, [3]int{i, from, j})
			from = j + 1
		}
	}
}
