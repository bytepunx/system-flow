package mdlint

import "strings"

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
