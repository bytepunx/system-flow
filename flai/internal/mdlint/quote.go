package mdlint

import (
	"strings"
	"unicode"
)

// QuoteBareURLs returns line, one line of inline markdown, with each bare
// URL, www. literal, and email address that MD034 finds in it wrapped in a
// code span, so that text flai takes from elsewhere can be written lint
// clean. Everything else is left as it is.
func QuoteBareURLs(line string) string {
	var seg segment
	seg.add(0, line, 0)
	var out inlineOut
	parseInline(&seg, &out)
	if len(out.urls) == 0 {
		return line
	}
	var b strings.Builder
	at := 0
	for _, u := range out.urls {
		b.WriteString(line[at:u.col])
		b.WriteString(codeSpanOf(line[u.col : u.col+u.n]))
		at = u.col + u.n
	}
	b.WriteString(line[at:])
	return b.String()
}

// codeSpanOf is s in a code span whose fence is one backtick longer than the
// longest run of backticks in s, padded with a space on each side when s
// starts or ends with a backtick.
func codeSpanOf(s string) string {
	longest := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '`' {
			n := runLen(s, i, '`')
			longest = max(longest, n)
			i += n - 1
		}
	}
	fence := strings.Repeat("`", longest+1)
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		s = " " + s + " "
	}
	return fence + s + fence
}

// Fenced returns text, the output of a tool, in a fenced code block of
// language text that flai's own lint passes: its line ends made \n, its tabs
// four spaces (MD010), its ANSI escape sequences and other control
// characters dropped, its trailing spaces trimmed (MD009), its runs of blank
// lines folded to one and those at its start and end dropped (MD012), and
// its fence one backtick longer than its longest run of backticks and at
// least three. A thread that quoted go test's summary line with its tabs was
// refused by the lint, so the merge queue could not send its story back
// (I-0242). An empty text gives an empty block.
func Fenced(text string) string {
	text = strings.NewReplacer("\r\n", "\n", "\r", "\n", "\t", "    ").Replace(text)
	text = dropControls(text)
	var lines []string
	for l := range strings.SplitSeq(text, "\n") {
		l = strings.TrimRight(l, " ")
		if l == "" && (len(lines) == 0 || lines[len(lines)-1] == "") {
			continue
		}
		lines = append(lines, l)
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	longest := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '`' {
			n := runLen(text, i, '`')
			longest = max(longest, n)
			i += n - 1
		}
	}
	fence := strings.Repeat("`", max(3, longest+1))
	if len(lines) == 0 {
		return fence + "text\n" + fence
	}
	return fence + "text\n" + strings.Join(lines, "\n") + "\n" + fence
}

// dropControls is s without its ANSI escape sequences, CSI (ESC [ ... final
// byte), OSC (ESC ] ... ended by BEL or ESC \), and the others (ESC, its
// intermediate bytes, and a final byte, as ESC ( B), and without every other
// control character but \n.
func dropControls(s string) string {
	var b strings.Builder
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		c := r[i]
		switch {
		case c == '\x1b' && i+1 < len(r) && r[i+1] == '[':
			i += 2
			for i < len(r) && (r[i] < 0x40 || r[i] > 0x7e) {
				i++
			}
		case c == '\x1b' && i+1 < len(r) && r[i+1] == ']':
			i += 2
			for i < len(r) && r[i] != '\a' && (r[i] != '\x1b' || i+1 >= len(r) || r[i+1] != '\\') {
				i++
			}
			if i < len(r) && r[i] == '\x1b' {
				i++
			}
		case c == '\x1b':
			for i+1 < len(r) && r[i+1] >= 0x20 && r[i+1] <= 0x2f {
				i++
			}
			if i+1 < len(r) && r[i+1] >= 0x30 && r[i+1] <= 0x7e {
				i++
			}
		case c == '\n' || !unicode.IsControl(c):
			b.WriteRune(c)
		}
	}
	return b.String()
}
