package mdlint

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Each rule follows markdownlint 0.40's implementation of it: the same
// lines, the same options and defaults, the same detail text.

var (
	entityEndRe = regexp.MustCompile(`&(?:#\d+|#[xX][\da-fA-F]+|[a-zA-Z]{2,31}|blk\d{2}|frac\d{2}|sup\d|there4);$`)
	gemojiEndRe = regexp.MustCompile(`:(?:[abmovx]|[-+]1|100|1234|(?:1st|2nd|3rd)_place_medal|8ball|clock\d{1,4}|e-mail|non-potable_water|o2|t-rex|u5272|u5408|u55b6|u6307|u6708|u6709|u6e80|u7121|u7533|u7981|u7a7a|[a-z]{2,15}2?|[a-z]{1,14}(?:_[a-z\d]{1,16})+):$`)
	noSpaceATX  = regexp.MustCompile(`^#+[^# \t]`)
	endsHash    = regexp.MustCompile(`#\s*$`)
	manySpace   = regexp.MustCompile(`^#+[ \t]{2,}\S`)
	tabsRe      = regexp.MustCompile(`\t+`)
	startSpace  = regexp.MustCompile(`^\s+\S`)
	endSpace    = regexp.MustCompile(`\S\s+$`)
)

// blankLine is markdownlint's isBlankLine: empty once HTML comments and
// blockquote markers are taken out. A line before the start or past the end
// of the document counts as blank.
func blankLine(d *doc, i int) bool {
	if i < d.content || i >= len(d.lines) {
		return true
	}
	s := d.lines[i].raw
	if strings.TrimSpace(s) == "" {
		return true
	}
	for {
		start, end := strings.Index(s, "<!--"), strings.Index(s, "-->")
		switch {
		case end != -1 && (start == -1 || end < start):
			s = s[end+3:]
			continue
		case start != -1 && end != -1:
			s = s[:start] + s[end+3:]
			continue
		case start != -1:
			s = s[:start]
		}
		break
	}
	return strings.TrimSpace(strings.ReplaceAll(s, ">", "")) == ""
}

// frontTitle reports whether the front matter has a title the rule counts as
// the document's top-level heading.
func frontTitle(c *Config, id string, d *doc) bool {
	pattern, set := c.strOpt(id, "front_matter_title", "")
	if set && pattern == "" {
		return false
	}
	if pattern == "" {
		pattern = `^\s*"?title"?\s*[:=]`
	}
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return false
	}
	for _, l := range d.front {
		if re.MatchString(l) {
			return true
		}
	}
	return false
}

// inline parses every paragraph, heading, and table cell.
func (d *doc) inline() *inlineOut {
	out := &inlineOut{}
	d.paraSeg = map[*para]*segment{}
	d.paraInfo = map[*para]segInfo{}
	for _, p := range d.paras {
		s := &segment{pair: true}
		for _, i := range p.lines {
			s.add(i, d.lines[i].raw[d.lines[i].start:], d.lines[i].start)
		}
		d.paraSeg[p] = s
		d.paraInfo[p] = parseInline(s, out)
	}
	for _, h := range d.headings {
		if h.setext {
			continue
		}
		raw := d.lines[h.start].raw
		at := d.lines[h.start].start
		if k := strings.Index(raw[at:], h.text); k >= 0 && h.text != "" {
			s := &segment{pair: true}
			s.add(h.start, h.text, at+k)
			parseInline(s, out)
		}
	}
	for _, cell := range d.tableCells {
		s := &segment{pair: true}
		s.add(cell[0], d.lines[cell[0]].raw[cell[1]:cell[2]], cell[1])
		parseInline(s, out)
	}
	return out
}

func md001(c *Config, d *doc, _ *inlineOut, add adder) {
	prev := 1 << 30
	if frontTitle(c, "MD001", d) {
		prev = 1
	}
	for _, h := range d.headings {
		if h.level > prev+1 {
			add(h.start, fmt.Sprintf("Expected: h%d; Actual: h%d", prev+1, h.level), "")
		}
		prev = h.level
	}
}

func md004(c *Config, d *doc, _ *inlineOut, add adder) {
	style, _ := c.strOpt("MD004", "style", "consistent")
	names := map[byte]string{'*': "asterisk", '-': "dash", '+': "plus"}
	switch style {
	case "consistent", "asterisk", "dash", "plus":
	default:
		return // sublist depends on nesting this package does not judge
	}
	for _, l := range d.lists {
		if l.ordered {
			continue
		}
		for _, it := range l.items {
			got := names[it.marker]
			if style == "consistent" {
				style = got
			}
			if got != style {
				add(it.line, fmt.Sprintf("Expected: %s; Actual: %s", style, got), "")
			}
		}
	}
}

func md007(c *Config, d *doc, _ *inlineOut, add adder) {
	indent := c.intOpt("MD007", "indent", 0)
	if indent == 0 {
		indent = 2
	}
	start := 0
	if c.boolOpt("MD007", "start_indented", false) {
		if start = c.intOpt("MD007", "start_indent", 0); start == 0 {
			start = indent
		}
	}
	// The actual indent is in characters, a tab counting one, as markdownlint
	// counts it, from the last quote marker on the line; nesting counts the
	// lists up to the nearest quote. Footnotes are not parsed, so lists in
	// them are not judged.
	for _, l := range d.lists {
		if l.ordered {
			continue
		}
		nesting := 0
		for p := l.parent; p != nil && nesting >= 0; p = p.list.parent {
			nesting++
			if p.ordered {
				nesting = -1 // markdownlint skips lists under an ordered list
			}
		}
		if nesting < 0 {
			continue
		}
		for _, it := range l.items {
			want, got := start+nesting*indent, it.at-it.quote
			if got != want {
				add(it.line, fmt.Sprintf("Expected: %d; Actual: %d", want, got), "")
			}
		}
	}
}

func md009(c *Config, d *doc, _ *inlineOut, add adder) {
	if c.boolOpt("MD009", "strict", false) {
		return // strict depends on paragraph ends this package does not judge
	}
	br := c.intOpt("MD009", "br_spaces", 2)
	expected := br
	if br < 2 {
		expected = 0
	}
	for i := d.content; i < len(d.lines); i++ {
		ln := d.lines[i]
		if ln.kind == kCode || ln.kind == kIndentCode {
			continue
		}
		n := len(ln.raw) - len(strings.TrimRightFunc(ln.raw, unicode.IsSpace))
		if n == 0 || n == expected {
			continue
		}
		detail := fmt.Sprintf("Expected: 0 or %d; Actual: %d", expected, n)
		if expected == 0 {
			detail = fmt.Sprintf("Expected: 0; Actual: %d", n)
		}
		add(i, detail, "")
	}
}

func md010(c *Config, d *doc, _ *inlineOut, add adder) {
	code := c.boolOpt("MD010", "code_blocks", true)
	ignore := c.listOpt("MD010", "ignore_code_languages")
	for i := d.content; i < len(d.lines); i++ {
		ln := d.lines[i]
		if ln.kind == kCode || ln.kind == kIndentCode {
			if !code {
				continue
			}
			if ln.fence != nil && contains(ignore, strings.ToLower(firstWord(ln.fence.info))) {
				continue
			}
		}
		for _, loc := range tabsRe.FindAllStringIndex(ln.raw, -1) {
			add(i, fmt.Sprintf("Column: %d", loc[0]+1), "")
		}
	}
}

func md012(c *Config, d *doc, _ *inlineOut, add adder) {
	max := c.intOpt("MD012", "maximum", 1)
	count := 0
	for i := d.content; i < len(d.lines); i++ {
		ln := d.lines[i]
		if ln.kind == kCode || ln.kind == kIndentCode || !isBlank(ln.raw) {
			count = 0
			continue
		}
		count++
		if count > max {
			add(i, fmt.Sprintf("Expected: %d; Actual: %d", max, count), "")
		}
	}
}

func md018(_ *Config, d *doc, _ *inlineOut, add adder) {
	for i := d.content; i < len(d.lines); i++ {
		ln := d.lines[i]
		switch ln.kind {
		case kCode, kIndentCode, kHTML, kFenceOpen, kFenceClose:
			continue
		}
		if noSpaceATX.MatchString(ln.raw) && !endsHash.MatchString(ln.raw) && !strings.HasPrefix(ln.raw, "#️⃣") {
			add(i, "", strings.TrimSpace(ln.raw))
		}
	}
}

func md019(_ *Config, d *doc, _ *inlineOut, add adder) {
	for _, h := range d.headings {
		if !h.setext && !h.closed && manySpace.MatchString(h.raw) {
			add(h.start, "", strings.TrimSpace(h.raw))
		}
	}
}

// linesOpt reads lines_above or lines_below: a number, or one per level.
func linesOpt(c *Config, name string, level int) int {
	if l, ok := c.options["MD022"][name].([]any); ok {
		if level-1 < len(l) {
			switch v := l[level-1].(type) {
			case uint64:
				return int(v)
			case int64:
				return int(v)
			}
		}
		return -1
	}
	return c.intOpt("MD022", name, 1)
}

func md022(c *Config, d *doc, _ *inlineOut, add adder) {
	for _, h := range d.headings {
		text := strings.TrimSpace(d.lines[h.start].raw)
		if want := linesOpt(c, "lines_above", h.level); want >= 0 {
			got := 0
			for got < want && blankLine(d, h.start-1-got) {
				got++
			}
			if got != want {
				add(h.start, fmt.Sprintf("Expected: %d; Actual: %d; Above", want, got), text)
			}
		}
		if want := linesOpt(c, "lines_below", h.level); want >= 0 {
			got := 0
			for got < want && blankLine(d, h.end+1+got) {
				got++
			}
			if got != want {
				add(h.start, fmt.Sprintf("Expected: %d; Actual: %d; Below", want, got), text)
			}
		}
	}
}

func md024(c *Config, d *doc, _ *inlineOut, add adder) {
	siblings := c.boolOpt("MD024", "siblings_only", false) || c.boolOpt("MD024", "allow_different_nesting", false)
	known := [][]string{nil, {}}
	last := 1
	current := known[1]
	for _, h := range d.headings {
		text := strings.Join(strings.Fields(h.text), " ")
		if siblings {
			for last < h.level {
				last++
				for len(known) <= last {
					known = append(known, nil)
				}
				known[last] = []string{}
			}
			if last > h.level {
				last = h.level
			}
			for len(known) <= h.level {
				known = append(known, []string{})
			}
			current = known[h.level]
		}
		if contains(current, text) {
			add(h.start, "", text)
			continue
		}
		current = append(current, text)
		if siblings {
			known[h.level] = current
		} else {
			known[1] = current
		}
	}
}

func md025(c *Config, d *doc, _ *inlineOut, add adder) {
	level := c.intOpt("MD025", "level", 1)
	fm := frontTitle(c, "MD025", d)
	top := false
	for _, h := range d.headings {
		if h.level != level {
			continue
		}
		if top || fm {
			add(h.start, "", h.text)
		} else if h.start == d.content {
			top = true
		}
	}
}

func md026(c *Config, d *doc, _ *inlineOut, add adder) {
	punct, _ := c.strOpt("MD026", "punctuation", ".,;:!。，；：！")
	if punct == "" {
		return
	}
	for _, h := range d.headings {
		t := strings.TrimRight(h.text, " \t")
		end := len(t)
		for end > 0 {
			r, size := utf8.DecodeLastRuneInString(t[:end])
			if !strings.ContainsRune(punct, r) {
				break
			}
			end -= size
		}
		if end == len(t) || entityEndRe.MatchString(t) || gemojiEndRe.MatchString(t) {
			continue
		}
		add(h.textAt, fmt.Sprintf("Punctuation: '%s'", t[end:]), "")
	}
}

func md029(c *Config, d *doc, _ *inlineOut, add adder) {
	style, _ := c.strOpt("MD029", "style", "one_or_ordered")
	examples := map[string]string{"one": "1/1/1", "ordered": "1/2/3", "zero": "0/0/0"}
	for _, l := range d.lists {
		if !l.ordered {
			continue
		}
		expected, incrementing := 1, false
		if len(l.items) >= 2 && (l.items[1].value != 1 || l.items[0].value == 0) {
			incrementing = true
			if l.items[0].value == 0 {
				expected = 0
			}
		}
		s := style
		switch s {
		case "one_or_ordered":
			s = "one"
			if incrementing {
				s = "ordered"
			}
		case "zero":
			expected = 0
		case "one":
			expected = 1
		}
		if _, ok := examples[s]; !ok {
			return
		}
		for _, it := range l.items {
			if it.value != expected {
				add(it.line, fmt.Sprintf("Expected: %d; Actual: %d; Style: %s", expected, it.value, examples[s]), "")
			}
			if s == "ordered" {
				expected++
			}
		}
	}
}

func md031(c *Config, d *doc, _ *inlineOut, add adder) {
	items := c.boolOpt("MD031", "list_items", true)
	for _, f := range d.fences {
		if !items && f.listed {
			continue
		}
		if !blankLine(d, f.open-1) {
			add(f.open, "", strings.TrimSpace(d.lines[f.open].raw))
		}
		if f.close >= 0 && !blankLine(d, f.close+1) {
			add(f.close, "", strings.TrimSpace(d.lines[f.close].raw))
		}
	}
}

func md034(_ *Config, _ *doc, in *inlineOut, add adder) {
	for _, u := range in.urls {
		add(u[0], "", "")
	}
}

func md036(c *Config, d *doc, _ *inlineOut, add adder) {
	punct, _ := c.strOpt("MD036", "punctuation", ".,;:!?。，；：！？")
	for _, p := range d.paras {
		if p.depth != 0 || p.heading != nil || len(p.lines) != 1 {
			continue
		}
		info := d.paraInfo[p]
		if len(info.emphs) != 1 || info.leftovers != 0 || info.other != 0 {
			continue
		}
		s, e := d.paraSeg[p].text, info.emphs[0]
		lead := len(s) - len(strings.TrimLeft(s, " \t"))
		if e.from != lead || e.to != len(strings.TrimRight(s, " \t")) {
			continue
		}
		inner := s[e.inner:e.inEnd]
		if inner == "" || strings.ContainsAny(inner, "&\\*_`<[]\t\n") {
			continue
		}
		if r, _ := utf8.DecodeLastRuneInString(inner); strings.ContainsRune(punct, r) {
			continue
		}
		add(p.lines[0], "", inner)
	}
}

func md037(_ *Config, d *doc, in *inlineOut, add adder) {
	type key struct {
		group int
		text  string
	}
	pending := map[key]*leftover{}
	for k := range in.leftovers {
		lo := &in.leftovers[k]
		switch lo.text {
		case "*", "**", "***", "_", "__", "___":
		default:
			continue
		}
		kk := key{lo.group, lo.text}
		open := pending[kk]
		if open == nil {
			pending[kk] = lo
			continue
		}
		delete(pending, kk)
		if raw := d.lines[open.line].raw; open.col+len(open.text) <= len(raw) {
			if m := startSpace.FindString(raw[open.col+len(open.text):]); m != "" {
				add(open.line, "", open.text+m)
			}
		}
		if raw := d.lines[lo.line].raw; lo.col <= len(raw) {
			if m := endSpace.FindString(raw[:lo.col]); m != "" {
				add(lo.line, "", m+lo.text)
			}
		}
	}
}

// jsSpace is JavaScript's \s, which markdownlint's rules match with.
func jsSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xa0, 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff:
		return true
	}
	return r >= 0x2000 && r <= 0x200a
}

// ellipsify shortens a context as markdownlint does: past 30 characters it
// keeps the start, the end, or both, of what the finding is about.
func ellipsify(text string, start, end bool) string {
	r := []rune(text)
	switch {
	case len(r) <= 30:
		return text
	case start && end:
		return string(r[:15]) + "..." + string(r[len(r)-15:])
	case end:
		return "..." + string(r[len(r)-30:])
	}
	return string(r[:30]) + "..."
}

// codeLine is the part of a code span on one line.
type codeLine struct {
	line int
	text string
}

// spanLines splits seg's text from from to to by line. A line after the
// first keeps the indentation past its containers' prefix, as micromark
// keeps it in a code span; with raw, it keeps all of it, as the source has
// it.
func spanLines(d *doc, seg *segment, from, to int, raw bool) []codeLine {
	var out []codeLine
	for k, p := range seg.pieces {
		end := len(seg.text)
		if k+1 < len(seg.pieces) {
			end = seg.pieces[k+1].at - 1
		}
		if p.at > to || end < from {
			continue
		}
		lo, hi := max(p.at, from), min(end, to)
		text := seg.text[lo:hi]
		if p.at > from {
			lead := d.lines[p.line].lead
			if raw {
				lead = 0
			}
			if lead <= p.col {
				text = d.lines[p.line].raw[lead:p.col] + text
			}
		}
		out = append(out, codeLine{p.line, text})
	}
	return out
}

// md038 follows micromark's code text: one space or line ending at each end
// is padding when both ends have one and the code is not all spaces, and
// markdownlint reports what whitespace is left at the start of the first
// line of code and at the end of the last, but for a space before a
// backtick the padding cannot hold.
func md038(_ *Config, d *doc, in *inlineOut, add adder) {
	for _, cs := range in.codes {
		ls := spanLines(d, cs.seg, cs.from, cs.to, false)
		if len(ls) == 0 {
			continue
		}
		first, last := ls[0].text, ls[len(ls)-1].text
		pads := func(s string, at int) bool {
			if len(ls) > 1 && s == "" {
				return true // a line ending
			}
			return s != "" && s[at] == ' '
		}
		data := false
		for _, l := range ls {
			if strings.Trim(l.text, " ") != "" {
				data = true
			}
		}
		padding := data && pads(first, 0) && pads(last, len(last)-1)
		if padding {
			if first == "" {
				ls = ls[1:]
			} else {
				ls[0].text = first[1:]
			}
			n := len(ls) - 1
			if last == "" {
				ls = ls[:n]
			} else {
				ls[n].text = ls[n].text[:len(ls[n].text)-1]
			}
		}
		var datas []codeLine
		for _, l := range ls {
			if l.text != "" {
				datas = append(datas, l)
			}
		}
		if len(datas) == 0 {
			continue
		}
		var parts []string
		for _, l := range spanLines(d, cs.seg, cs.open, cs.close, true) {
			parts = append(parts, l.text)
		}
		context := strings.Join(parts, " ")
		// the whitespace before the first non-space and after the last
		count := func(s string, fromEnd bool) int {
			n := 0
			for s != "" {
				var r rune
				var size int
				if fromEnd {
					r, size = utf8.DecodeLastRuneInString(s)
				} else {
					r, size = utf8.DecodeRuneInString(s)
				}
				if !jsSpace(r) {
					if n > 0 && r == '`' && !padding {
						n--
					}
					return n
				}
				n++
				if fromEnd {
					s = s[:len(s)-size]
				} else {
					s = s[size:]
				}
			}
			return 0 // all whitespace
		}
		if s := datas[0]; count(s.text, false) > 0 {
			add(s.line, "", ellipsify(context, true, false))
		}
		if s := datas[len(datas)-1]; count(s.text, true) > 0 {
			add(s.line, "", ellipsify(context, false, true))
		}
	}
}

func md040(c *Config, d *doc, _ *inlineOut, add adder) {
	allowed := c.listOpt("MD040", "allowed_languages")
	only := c.boolOpt("MD040", "language_only", false)
	for _, f := range d.fences {
		lang := firstWord(f.info)
		ln := d.lines[f.open]
		text := strings.TrimSpace(ln.raw[ln.start:])
		switch {
		case lang == "":
			add(f.open, "", text)
		case len(allowed) > 0 && !contains(allowed, lang):
			add(f.open, fmt.Sprintf("Language %q is not allowed", lang), "")
		case only && f.info != lang:
			add(f.open, fmt.Sprintf("Info string contains more than language: %q", f.info), "")
		}
	}
}

func md047(_ *Config, d *doc, _ *inlineOut, add adder) {
	if last := len(d.lines) - 1; last >= d.content && !blankLine(d, last) {
		add(last, "", "")
	}
}

func emphasisStyle(c *Config, id string, d *doc, in *inlineOut, strong bool, add adder) {
	style, _ := c.strOpt(id, "style", "consistent")
	var list []emph
	for _, e := range in.emphs {
		if e.strong == strong {
			list = append(list, e)
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].openLine != list[j].openLine {
			return list[i].openLine < list[j].openLine
		}
		return list[i].openCol < list[j].openCol
	})
	for _, e := range list {
		got := "asterisk"
		if e.char == '_' {
			got = "underscore"
		}
		if style == "consistent" {
			style = got
		}
		if style == got {
			continue
		}
		if style == "underscore" && intraword(d, e) {
			continue
		}
		detail := fmt.Sprintf("Expected: %s; Actual: %s", style, got)
		add(e.openLine, detail, "")
		add(e.closeLine, detail, "")
	}
}

// intraword reports whether emphasis sits inside a word, where only
// asterisks work.
func intraword(d *doc, e emph) bool {
	word := func(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) }
	open, close := d.lines[e.openLine].raw, d.lines[e.closeLine].raw
	n := 1
	if e.strong {
		n = 2
	}
	if e.openCol > 0 {
		if r, _ := utf8.DecodeLastRuneInString(open[:e.openCol]); word(r) {
			return true
		}
	}
	if e.closeCol+n < len(close) {
		if r, _ := utf8.DecodeRuneInString(close[e.closeCol+n:]); word(r) {
			return true
		}
	}
	return false
}

func md049(c *Config, d *doc, in *inlineOut, add adder) { emphasisStyle(c, "MD049", d, in, false, add) }
func md050(c *Config, d *doc, in *inlineOut, add adder) { emphasisStyle(c, "MD050", d, in, true, add) }

func firstWord(s string) string {
	if f := strings.Fields(s); len(f) > 0 {
		return f[0]
	}
	return ""
}
