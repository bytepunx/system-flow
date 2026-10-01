package mdlint

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Inline content: emphasis resolved as CommonMark resolves it (the
// delimiter run algorithm), with code spans, HTML, autolinks, link
// destinations, and bare URLs taken out first. What is left of a delimiter
// run once emphasis is resolved is literal text, which is what the
// space-in-emphasis rule pairs up.

// segment is the inline content of one paragraph, heading, or table cell:
// its lines' content joined by newlines, and where each piece came from.
type segment struct {
	text   string
	pieces []piece
	pair   bool // leftover markers are paired for MD037 (not in link text)
}

type piece struct{ at, line, col int }

// where maps an offset in the segment to a line index and a byte column in
// that line's raw text.
func (s *segment) where(pos int) (int, int) {
	k := sort.Search(len(s.pieces), func(i int) bool { return s.pieces[i].at > pos }) - 1
	if k < 0 {
		k = 0
	}
	p := s.pieces[k]
	return p.line, p.col + pos - p.at
}

func (s *segment) add(line int, text string, col int) {
	if len(s.pieces) > 0 {
		s.text += "\n"
	}
	s.pieces = append(s.pieces, piece{at: len(s.text), line: line, col: col})
	s.text += text
}

// emph is one resolved emphasis (one marker) or strong emphasis (two).
type emph struct {
	char         byte
	strong       bool
	openLine     int
	openCol      int
	closeLine    int
	closeCol     int
	from, to     int // the opening marker's start and the closing marker's end in the segment
	inner, inEnd int // the content between the markers
}

// leftover is a delimiter run, or what is left of one, that is literal text.
type leftover struct {
	text      string
	line, col int
	group     int
}

type delim struct {
	ch                byte
	lo, hi, orig      int
	canOpen, canClose bool
	dead              bool
}

func (d *delim) count() int { return d.hi - d.lo }

// inlineOut collects what the inline rules read from every segment.
type inlineOut struct {
	emphs     []emph
	leftovers []leftover
	urls      [][2]int // line, col of each bare URL
	groups    int
}

// segInfo is what parsing one segment found, for no-emphasis-as-heading.
type segInfo struct {
	emphs     []emph
	leftovers int
	other     int // code spans, HTML, links, URLs, escapes
}

var (
	autolinkRe = regexp.MustCompile(`^<[A-Za-z][A-Za-z0-9+.-]{1,31}:[^\s<>]*>`)
	emailRe    = regexp.MustCompile(`^<[A-Za-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*>`)
	openTagRe  = regexp.MustCompile(`^<[A-Za-z][A-Za-z0-9-]*(?:\s+[A-Za-z_:][\w.:-]*(?:\s*=\s*(?:[^\s"'=<>` + "`" + `]+|'[^']*'|"[^"]*"))?)*\s*/?>`)
	closeTagRe = regexp.MustCompile(`^</[A-Za-z][A-Za-z0-9-]*\s*>`)
	domainRe   = regexp.MustCompile(`^[A-Za-z0-9_.-]+`)
)

func isASCIIPunct(c byte) bool {
	return c < 128 && strings.IndexByte("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", c) >= 0
}

func isPunct(r rune) bool {
	if r < 128 {
		return isASCIIPunct(byte(r))
	}
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

func before(t string, i int) rune {
	if i <= 0 {
		return '\n'
	}
	r, _ := utf8.DecodeLastRuneInString(t[:i])
	return r
}

func after(t string, i int) rune {
	if i >= len(t) {
		return '\n'
	}
	r, _ := utf8.DecodeRuneInString(t[i:])
	return r
}

func runLen(t string, i int, c byte) int {
	n := 0
	for i+n < len(t) && t[i+n] == c {
		n++
	}
	return n
}

// codeClose finds the start of the closing backtick run of exactly n.
func codeClose(t string, from, n int) int {
	for i := from; i < len(t); {
		if t[i] != '`' {
			i++
			continue
		}
		m := runLen(t, i, '`')
		if m == n {
			return i
		}
		i += m
	}
	return -1
}

// matchBracket finds the ] that closes the [ at i, past escapes, code, and
// nested brackets.
func matchBracket(t string, i, to int) int {
	depth := 0
	for j := i; j < to; j++ {
		switch t[j] {
		case '\\':
			j++
		case '`':
			n := runLen(t, j, '`')
			if k := codeClose(t[:to], j+n, n); k >= 0 {
				j = k + n - 1
			} else {
				j += n - 1
			}
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return j
			}
		}
	}
	return -1
}

func matchParen(t string, i, to int) int {
	depth := 0
	for j := i; j < to; j++ {
		switch t[j] {
		case '\\':
			j++
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return j
			}
		case '\n':
			if depth > 0 && j+1 < to && t[j+1] == '\n' {
				return -1
			}
		}
	}
	return -1
}

// bareURL returns the length of a GFM literal http(s) autolink at i, or 0.
func bareURL(t string, i int) int {
	s := t[i:]
	var n int
	switch {
	case strings.HasPrefix(s, "https://"):
		n = 8
	case strings.HasPrefix(s, "http://"):
		n = 7
	default:
		return 0
	}
	if i > 0 && (t[i-1] >= 'a' && t[i-1] <= 'z' || t[i-1] >= 'A' && t[i-1] <= 'Z') {
		return 0
	}
	dom := domainRe.FindString(s[n:])
	segs := strings.Split(strings.TrimRight(dom, "."), ".")
	if len(segs) < 2 {
		return 0
	}
	for _, sg := range segs[len(segs)-2:] {
		if sg == "" || strings.Contains(sg, "_") {
			return 0
		}
	}
	end := n
	for end < len(s) && !unicode.IsSpace(rune(s[end])) && s[end] != '<' {
		end++
	}
	return trimURL(s[:end])
}

// trimURL drops the trailing punctuation and unbalanced parentheses GFM
// leaves out of an autolink.
func trimURL(u string) int {
	for len(u) > 0 {
		c := u[len(u)-1]
		switch {
		case strings.IndexByte("?!.,:*_~'\"", c) >= 0:
			u = u[:len(u)-1]
		case c == ')' && strings.Count(u, ")") > strings.Count(u, "("):
			u = u[:len(u)-1]
		default:
			return len(u)
		}
	}
	return 0
}

// parseInline resolves one segment's emphasis and records what the rules
// read into out.
func parseInline(seg *segment, out *inlineOut) segInfo {
	var info segInfo
	parseRange(seg, 0, len(seg.text), seg.pair, out, &info)
	return info
}

func parseRange(seg *segment, from, to int, pair bool, out *inlineOut, info *segInfo) {
	t := seg.text
	var ds []*delim
	for i := from; i < to; {
		c := t[i]
		switch {
		case c == '\\' && i+1 < to && isASCIIPunct(t[i+1]):
			info.other++
			i += 2
		case c == '`':
			n := runLen(t[:to], i, '`')
			if k := codeClose(t[:to], i+n, n); k >= 0 {
				info.other++
				i = k + n
			} else {
				i += n
			}
		case c == '<':
			n := 0
			if strings.HasPrefix(t[i:to], "<!--") {
				if k := strings.Index(t[i+4:to], "-->"); k >= 0 {
					n = 4 + k + 3
				}
			}
			for _, re := range []*regexp.Regexp{autolinkRe, emailRe, openTagRe, closeTagRe} {
				if n == 0 {
					n = len(re.FindString(t[i:to]))
				}
			}
			if n > 0 {
				info.other++
				i += n
			} else {
				i++
			}
		case c == '[' || (c == '!' && i+1 < to && t[i+1] == '['):
			open := i
			if c == '!' {
				open++
			}
			end := matchBracket(t, open, to)
			if end < 0 {
				i = open + 1
				continue
			}
			info.other++
			parseRange(seg, open+1, end, false, out, info)
			j := end + 1
			if j < to && t[j] == '(' {
				if k := matchParen(t, j, to); k >= 0 {
					j = k + 1
				}
			} else if j < to && t[j] == '[' {
				if k := strings.IndexByte(t[j:to], ']'); k >= 0 {
					j += k + 1
				}
			}
			i = j
		case c == '*' || c == '_':
			n := runLen(t[:to], i, c)
			prev, next := before(t[:to], i), after(t[:to], i+n)
			if i == from {
				prev = '\n'
			}
			left := !unicode.IsSpace(next) && (!isPunct(next) || unicode.IsSpace(prev) || isPunct(prev))
			right := !unicode.IsSpace(prev) && (!isPunct(prev) || unicode.IsSpace(next) || isPunct(next))
			d := &delim{ch: c, lo: i, hi: i + n, orig: n}
			if c == '*' {
				d.canOpen, d.canClose = left, right
			} else {
				d.canOpen = left && (!right || isPunct(prev))
				d.canClose = right && (!left || isPunct(next))
			}
			ds = append(ds, d)
			i += n
		case c == 'h' && bareURL(t[:to], i) > 0:
			n := bareURL(t[:to], i)
			l, col := seg.where(i)
			out.urls = append(out.urls, [2]int{l, col})
			info.other++
			i += n
		case c == 'w' && strings.HasPrefix(t[i:to], "www.") && (i == from || strings.ContainsRune(" \t\n(*_~[]", before(t, i))):
			n := 4
			for i+n < to && !unicode.IsSpace(rune(t[i+n])) && t[i+n] != '<' {
				n++
			}
			if k := trimURL(t[i : i+n]); k > 0 {
				n = k
			}
			info.other++
			i += n
		default:
			_, size := utf8.DecodeRuneInString(t[i:])
			i += size
		}
	}

	// resolve emphasis
	var matched []emph
	for ci, cl := range ds {
		if !cl.canClose {
			continue
		}
		for cl.count() > 0 {
			found := -1
			for oi := ci - 1; oi >= 0; oi-- {
				o := ds[oi]
				if o.dead || o.ch != cl.ch || !o.canOpen || o.count() == 0 {
					continue
				}
				if (o.canClose || cl.canOpen) && (o.orig+cl.orig)%3 == 0 && (o.orig%3 != 0 || cl.orig%3 != 0) {
					continue
				}
				found = oi
				break
			}
			if found < 0 {
				break
			}
			o := ds[found]
			use := 1
			if o.count() >= 2 && cl.count() >= 2 {
				use = 2
			}
			e := emph{char: cl.ch, strong: use == 2, from: o.hi - use, to: cl.lo + use, inner: o.hi, inEnd: cl.lo}
			e.openLine, e.openCol = seg.where(o.hi - use)
			e.closeLine, e.closeCol = seg.where(cl.lo)
			matched = append(matched, e)
			o.hi -= use
			cl.lo += use
			for k := found + 1; k < ci; k++ {
				ds[k].dead = true
			}
		}
	}
	out.emphs = append(out.emphs, matched...)
	info.emphs = append(info.emphs, matched...)

	// what is left is literal; group it by the innermost emphasis around it
	base := out.groups
	out.groups += len(matched) + 1
	for _, d := range ds {
		if d.count() == 0 {
			continue
		}
		info.leftovers++
		if !pair {
			continue
		}
		g, span := base+len(matched), -1
		for k, e := range matched {
			if d.lo >= e.inner && d.hi <= e.inEnd && (span < 0 || e.inEnd-e.inner < span) {
				g, span = base+k, e.inEnd-e.inner
			}
		}
		l, col := seg.where(d.lo)
		out.leftovers = append(out.leftovers, leftover{text: t[d.lo:d.hi], line: l, col: col, group: g})
	}
}
