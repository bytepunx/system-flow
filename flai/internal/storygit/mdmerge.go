package storygit

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ErrNotSettled marks a markdown file that MergeMarkdown cannot merge by sections.
var ErrNotSettled = errors.New("the markdown file is not settled by sections; resolve the stop by hand")

// MergeMarkdown merges a markdown file from a stopped rebase's three versions by sections (ADR-0175).
//
// ours is the branch being rebased onto and theirs the commit being replayed,
// as Stages names them. A file it cannot settle is an error wrapping
// ErrNotSettled that names each section by its heading path, and each front
// matter key, it could not settle.
func MergeMarkdown(base, ours, theirs string) (string, error) {
	b, o, t := parseMarkdown(base), parseMarkdown(ours), parseMarkdown(theirs)
	front, unsettled := mergeFrontMatter(b.front, o.front, t.front)
	parts, why := mergeSections(b.sections, o.sections, t.sections)
	unsettled = append(unsettled, why...)
	if len(unsettled) > 0 {
		return "", fmt.Errorf("%w: %s", ErrNotSettled, strings.Join(unsettled, "; "))
	}
	return renderMarkdown(front, parts, o.sections, t.sections, o.newline || t.newline), nil
}

// mdDoc is a markdown file cut by sections.
type mdDoc struct {
	front    *frontMatter // nil when the file has none
	sections []mdSection  // the text before the first heading always first
	newline  bool         // the file ends with a newline
}

// frontMatter is a file's YAML front matter, cut by its top-level keys.
type frontMatter struct {
	open, close string // the delimiter lines as written
	entries     []mdPart
}

// mdPart is a front matter key or a section: its key, how an error names
// it, and its lines without their newlines.
type mdPart struct {
	key, where string
	lines      []string
}

// mdSection is a section, or the text before the first heading.
type mdSection = mdPart

// preambleKey keys the text before the first heading, which every version has.
const preambleKey = ""

// parseMarkdown cuts text into its front matter and sections.
func parseMarkdown(text string) mdDoc {
	lines, newline := splitLines(text)
	doc := mdDoc{newline: newline}
	if len(lines) > 0 && isDelimiter(lines[0]) {
		for i := 1; i < len(lines); i++ {
			if isDelimiter(lines[i]) {
				doc.front = &frontMatter{open: lines[0], close: lines[i], entries: frontEntries(lines[1:i])}
				lines = lines[i+1:]
				break
			}
		}
	}
	doc.sections = cutSections(lines)
	return doc
}

// splitLines splits text into lines without their newlines, and reports
// whether it ends with one.
func splitLines(text string) ([]string, bool) {
	if text == "" {
		return nil, false
	}
	newline := strings.HasSuffix(text, "\n")
	return strings.Split(strings.TrimSuffix(text, "\n"), "\n"), newline
}

// isDelimiter reports whether line opens or closes front matter.
func isDelimiter(line string) bool { return strings.TrimRight(line, "\r") == "---" }

// frontEntries cuts front matter lines at its top-level keys; the lines after
// a key up to the next, such as list items and indented lines, belong to it.
func frontEntries(lines []string) []mdPart {
	var out []mdPart
	seen := map[string]int{}
	for _, l := range lines {
		if key, ok := frontKey(l); ok {
			seen[key]++
			out = append(out, mdPart{
				key:   key + "\x00" + strconv.Itoa(seen[key]),
				where: fmt.Sprintf("front matter key %q", key) + occurrence(seen[key]),
				lines: []string{l},
			})
			continue
		}
		if len(out) == 0 {
			out = append(out, mdPart{key: preambleKey, where: "the front matter before its first key"})
		}
		out[len(out)-1].lines = append(out[len(out)-1].lines, l)
	}
	return out
}

// frontKey reads the key of a top-level `key: value` line.
func frontKey(line string) (string, bool) {
	if line == "" || strings.ContainsRune(" \t-#", rune(line[0])) {
		return "", false
	}
	key, _, ok := strings.Cut(line, ":")
	return strings.TrimSpace(key), ok
}

// frontValue is a front matter line's value, without spaces or quotes.
func frontValue(line string) string {
	_, v, _ := strings.Cut(line, ":")
	return strings.Trim(strings.TrimSpace(v), `"'`)
}

// occurrence names the second and later uses of a repeated key or heading path.
func occurrence(n int) string {
	if n < 2 {
		return ""
	}
	return fmt.Sprintf(" (occurrence %d)", n)
}

// cutSections cuts the body of a markdown file at its ATX headings outside
// fenced code blocks, keying each section by its heading path.
func cutSections(lines []string) []mdSection {
	out := []mdSection{{key: preambleKey, where: "the text before the first heading"}}
	type heading struct {
		level int
		text  string
	}
	var stack []heading
	seen := map[string]int{}
	fence := ""
	for _, l := range lines {
		cur := &out[len(out)-1]
		switch {
		case fence != "":
			if closesFence(l, fence) {
				fence = ""
			}
		case fenceMarker(l) != "":
			fence = fenceMarker(l)
		default:
			level, text, ok := atxHeading(l)
			if !ok {
				break
			}
			for len(stack) > 0 && stack[len(stack)-1].level >= level {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, heading{level, text})
			texts := make([]string, len(stack))
			for i, h := range stack {
				texts[i] = h.text
			}
			path := strings.Join(texts, "\x00")
			seen[path]++
			out = append(out, mdSection{
				key:   path + "\x00" + strconv.Itoa(seen[path]),
				where: fmt.Sprintf("section %q", strings.Join(texts, " > ")) + occurrence(seen[path]),
				lines: []string{l},
			})
			continue
		}
		cur.lines = append(cur.lines, l)
	}
	return out
}

// fenceMarker is the run of three or more backticks or tildes that opens a
// fenced code block on line, or empty when line opens none.
func fenceMarker(line string) string {
	s := strings.TrimLeft(line, " ")
	if len(line)-len(s) > 3 || (!strings.HasPrefix(s, "```") && !strings.HasPrefix(s, "~~~")) {
		return ""
	}
	n := len(s) - len(strings.TrimLeft(s, s[:1]))
	if s[0] == '`' && strings.Contains(s[n:], "`") {
		return ""
	}
	return s[:n]
}

// closesFence reports whether line closes the fenced code block open opened.
func closesFence(line, open string) bool {
	m := fenceMarker(line)
	return m != "" && m[0] == open[0] && len(m) >= len(open) &&
		strings.TrimSpace(strings.TrimLeft(line, " ")[len(m):]) == ""
}

// atxHeading reads an ATX heading's level and text.
func atxHeading(line string) (int, string, bool) {
	s := strings.TrimLeft(strings.TrimRight(line, "\r"), " ")
	if len(line)-len(strings.TrimLeft(line, " ")) > 3 {
		return 0, "", false
	}
	level := len(s) - len(strings.TrimLeft(s, "#"))
	rest := s[level:]
	if level == 0 || level > 6 || (rest != "" && rest[0] != ' ' && rest[0] != '\t') {
		return 0, "", false
	}
	text := strings.TrimSpace(rest)
	if closed := strings.TrimRight(text, "#"); closed == "" || strings.HasSuffix(closed, " ") || strings.HasSuffix(closed, "\t") {
		text = strings.TrimSpace(closed)
	}
	return level, text, true
}

// sameText reports whether two versions of a part, nil when absent, are alike.
func sameText(a, b *mdPart) bool {
	if a == nil || b == nil {
		return a == b
	}
	return slices.Equal(a.lines, b.lines)
}

// byKey indexes parts by key.
func byKey(parts []mdPart) map[string]*mdPart {
	m := make(map[string]*mdPart, len(parts))
	for i := range parts {
		m[parts[i].key] = &parts[i]
	}
	return m
}

// mergeFrontMatter merges front matter key by key: ours' keys in ours'
// order, then the keys only theirs has.
func mergeFrontMatter(b, o, t *frontMatter) (*frontMatter, []string) {
	entriesOf := func(f *frontMatter) []mdPart {
		if f == nil {
			return nil
		}
		return f.entries
	}
	oe, te := entriesOf(o), entriesOf(t)
	bm, om, tm := byKey(entriesOf(b)), byKey(oe), byKey(te)
	var entries []mdPart
	var unsettled []string
	keep := func(key, where string) {
		e, ok := mergeFrontKey(key, bm[key], om[key], tm[key])
		switch {
		case ok && e != nil:
			entries = append(entries, *e)
		case !ok && om[key] != nil && tm[key] != nil:
			unsettled = append(unsettled, where+": both sides change it differently")
		case !ok:
			unsettled = append(unsettled, where+": one side removes it and the other changes it")
		}
	}
	for _, e := range oe {
		keep(e.key, e.where)
	}
	for _, e := range te {
		if om[e.key] == nil {
			keep(e.key, e.where)
		}
	}
	present := o != nil // the front matter itself follows the side that adds or removes it
	if (o != nil) == (b != nil) {
		present = t != nil
	}
	if !present && len(entries) == 0 {
		return nil, unsettled
	}
	delims := o
	if delims == nil {
		delims = t
	}
	if delims == nil {
		delims = b
	}
	return &frontMatter{open: delims.open, close: delims.close, entries: entries}, unsettled
}

// mergeFrontKey merges one front matter key; updated, changed differently by
// both, takes the later date. It reports false when the key is not settled.
func mergeFrontKey(key string, b, o, t *mdPart) (*mdPart, bool) {
	switch {
	case sameText(o, t), sameText(t, b):
		return o, true
	case sameText(o, b):
		return t, true
	case strings.HasPrefix(key, "updated\x00") && o != nil && t != nil:
		if frontValue(t.lines[0]) > frontValue(o.lines[0]) {
			return t, true
		}
		return o, true
	}
	return nil, false
}

// mdChunk is a section of the merged file, and whether it is ours' addition
// that theirs lacks, which theirs' additions at the same place follow.
type mdChunk struct {
	mdSection
	oursOnly bool
}

// mergeSections merges the sections in ours' order, placing each that only
// theirs adds after the section preceding it in theirs and after ours'
// additions there.
func mergeSections(b, o, t []mdSection) ([]mdChunk, []string) {
	bm, om, tm := byKey(b), byKey(o), byKey(t)
	var out []mdChunk
	var unsettled []string
	theirsToPlace := map[string]bool{}
	for i := range o {
		s := &o[i]
		bs, ts := bm[s.key], tm[s.key]
		switch {
		case sameText(s, ts), sameText(ts, bs):
			out = append(out, mdChunk{*s, bs == nil && !sameText(s, ts)})
		case sameText(s, bs):
			if ts != nil {
				out = append(out, mdChunk{mdSection: *ts})
			}
		case bs == nil:
			out = append(out, mdChunk{*s, true})
			theirsToPlace[s.key] = true
		case ts == nil:
			unsettled = append(unsettled, s.where+": one side removes it and the other changes it")
		default:
			lines, ok := mergeLines(bs.lines, s.lines, ts.lines)
			if !ok {
				unsettled = append(unsettled, s.where+": both sides change the same lines")
				continue
			}
			out = append(out, mdChunk{mdSection: mdSection{key: s.key, where: s.where, lines: lines}})
		}
	}
	for i := range t {
		s := &t[i]
		bs := bm[s.key]
		switch {
		case om[s.key] != nil && !theirsToPlace[s.key]:
			continue
		case om[s.key] == nil && bs != nil:
			if !sameText(s, bs) {
				unsettled = append(unsettled, s.where+": one side removes it and the other changes it")
			}
			continue
		}
		out = slices.Insert(out, placeAfter(out, t[:i]), mdChunk{mdSection: *s})
	}
	return out, unsettled
}

// placeAfter is where in out a section of theirs goes that follows the
// sections before: after the last of them out holds, and after ours'
// additions that follow it.
func placeAfter(out []mdChunk, before []mdSection) int {
	at := -1
	for i := len(before) - 1; i >= 0 && at < 0; i-- {
		for j := len(out) - 1; j >= 0 && at < 0; j-- {
			if out[j].key == before[i].key {
				at = j + 1
			}
		}
	}
	if at < 0 {
		at = len(out)
	}
	for at < len(out) && out[at].oursOnly {
		at++
	}
	return at
}

// renderMarkdown writes the merged file. A section that ends without a blank
// line gets one before the next heading when no side put the two together,
// as when it ended its side's file.
func renderMarkdown(front *frontMatter, parts []mdChunk, o, t []mdSection, newline bool) string {
	adjacent := map[[2]string]bool{}
	for _, side := range [][]mdSection{o, t} {
		for i := 1; i < len(side); i++ {
			adjacent[[2]string{side[i-1].key, side[i].key}] = true
		}
	}
	var lines []string
	if front != nil {
		lines = append(lines, front.open)
		for _, e := range front.entries {
			lines = append(lines, e.lines...)
		}
		lines = append(lines, front.close)
	}
	for i, p := range parts {
		if i > 0 {
			prev := parts[i-1].lines
			if len(prev) > 0 && strings.TrimSpace(prev[len(prev)-1]) != "" && !adjacent[[2]string{parts[i-1].key, p.key}] {
				lines = append(lines, "")
			}
		}
		lines = append(lines, p.lines...)
	}
	text := strings.Join(lines, "\n")
	if newline && len(lines) > 0 {
		text += "\n"
	}
	return text
}

// lineHunk replaces the base lines [start, end) with lines; start == end inserts.
type lineHunk struct {
	start, end int
	lines      []string
}

// mergeLines merges two sides' changes to base line by line. Changes that
// touch different base lines are both kept, even when adjacent; insertions
// at one place take ours' lines, then theirs'; a change both make alike is
// kept once. It reports false when both change or remove the same base lines.
func mergeLines(base, ours, theirs []string) ([]string, bool) {
	ho, ht := lineHunks(base, ours), lineHunks(base, theirs)
	var out []string
	pos := 0
	for len(ho) > 0 || len(ht) > 0 {
		var h lineHunk
		switch {
		case len(ht) == 0:
			h, ho = ho[0], ho[1:]
		case len(ho) == 0:
			h, ht = ht[0], ht[1:]
		default:
			a, c := ho[0], ht[0]
			switch {
			case a.start == c.start && a.end == c.end && slices.Equal(a.lines, c.lines):
				h, ho, ht = a, ho[1:], ht[1:]
			case a.start == a.end && c.start == c.end && a.start == c.start:
				h, ho, ht = lineHunk{a.start, a.end, slices.Concat(a.lines, c.lines)}, ho[1:], ht[1:]
			case touches(a, c):
				return nil, false
			case a.start < c.start || a.start == a.end:
				h, ho = a, ho[1:]
			default:
				h, ht = c, ht[1:]
			}
		}
		out = append(out, base[pos:h.start]...)
		out = append(out, h.lines...)
		pos = h.end
	}
	return append(out, base[pos:]...), true
}

// touches reports whether two hunks change the same base lines, or one
// inserts lines inside those the other changes.
func touches(a, c lineHunk) bool {
	switch {
	case a.start == a.end:
		return c.start < a.start && a.start < c.end
	case c.start == c.end:
		return a.start < c.start && c.start < a.end
	}
	return a.start < c.end && c.start < a.end
}

// lineHunks diffs other against base by their longest common subsequence,
// after setting aside the lines they share at the start and the end. Time
// and memory grow with the product of the lines left in each, so a section
// of thousands of lines both sides rewrite whole costs tens of megabytes.
func lineHunks(base, other []string) []lineHunk {
	p := 0
	for p < len(base) && p < len(other) && base[p] == other[p] {
		p++
	}
	s := 0
	for s < len(base)-p && s < len(other)-p && base[len(base)-1-s] == other[len(other)-1-s] {
		s++
	}
	a, b := base[p:len(base)-s], other[p:len(other)-s]
	n, m := len(a), len(b)
	lcs := make([]int, (n+1)*(m+1)) // lcs[i*(m+1)+j] is the LCS of a[i:] and b[j:]
	at := func(i, j int) int { return lcs[i*(m+1)+j] }
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i*(m+1)+j] = 1 + at(i+1, j+1)
			} else {
				lcs[i*(m+1)+j] = max(at(i+1, j), at(i, j+1))
			}
		}
	}
	var hunks []lineHunk
	var cur *lineHunk
	i, j := 0, 0
	for i < n || j < m {
		if i < n && j < m && a[i] == b[j] {
			if cur != nil {
				hunks = append(hunks, *cur)
				cur = nil
			}
			i, j = i+1, j+1
			continue
		}
		if cur == nil {
			cur = &lineHunk{start: p + i, end: p + i}
		}
		if j < m && (i == n || at(i, j+1) >= at(i+1, j)) {
			cur.lines = append(cur.lines, b[j])
			j++
		} else {
			i++
			cur.end = p + i
		}
	}
	if cur != nil {
		hunks = append(hunks, *cur)
	}
	return hunks
}
