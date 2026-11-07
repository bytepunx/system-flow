package release

import "strings"

// SettleChangelog is a changelog as main has it with the lines a branch
// added to it under main's top "## Unreleased" (ADR-0257). main is the base
// branch's file now, base the file at the branch's old merge base ("" when it
// had none), and branch the branch's file before the rebase. The branch's
// lines are those branch holds beyond base, compared as a multiset of line
// text, in branch's order, less blank lines, headings (lines starting with
// "#"), and lines main's top "## Unreleased" already holds. They go at the end
// of that section's body; a main with no "## Unreleased" gets one above its
// first "## " heading, as flai release writes it. With no such lines, main is
// returned as it is.
func SettleChangelog(main, base, branch string) string {
	own := addedLines(base, branch)
	lines := strings.SplitAfter(main, "\n")
	start, end := unreleasedSection(lines)
	if start >= 0 {
		own = without(own, lines[start+1:end])
	}
	if len(own) == 0 {
		return main
	}
	block := strings.Join(own, "\n") + "\n"
	var out string
	if start >= 0 {
		head := strings.Join(lines[:start+1], "")
		if !strings.HasSuffix(head, "\n") {
			head += "\n"
		}
		body := strings.TrimLeft(strings.TrimRight(strings.Join(lines[start+1:end], ""), " \t\r\n"), "\n")
		if body != "" {
			body += "\n"
		}
		out = head + "\n" + body + block
		if rest := strings.Join(lines[end:], ""); rest != "" {
			out += "\n" + rest
		}
	} else {
		section := unreleasedHeading + "\n\n" + block
		switch i := firstSection(main); {
		case i >= 0:
			out = main[:i] + section + "\n" + main[i:]
		case strings.TrimSpace(main) == "":
			out = section
		default:
			out = strings.TrimRight(main, "\n") + "\n\n" + section
		}
	}
	if main != "" && !strings.HasSuffix(main, "\n") {
		out = strings.TrimRight(out, "\n")
	}
	return out
}

// unreleasedSection is the index of the first "## Unreleased" heading among
// lines and the index of the next "## " heading after it, or len(lines);
// start is -1 when there is none.
func unreleasedSection(lines []string) (start, end int) {
	start = -1
	for i, l := range lines {
		if strings.TrimRight(l, " \t\r\n") == unreleasedHeading {
			start = i
			break
		}
	}
	if start < 0 {
		return -1, len(lines)
	}
	end = len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			end = i
			break
		}
	}
	return start, end
}

// addedLines is the entry lines branch holds beyond base, as a multiset of
// line text with trailing white space trimmed, in branch's order: blank lines
// and headings are left out.
func addedLines(base, branch string) []string {
	seen := map[string]int{}
	for _, l := range strings.Split(base, "\n") {
		seen[strings.TrimRight(l, " \t\r")]++
	}
	var out []string
	for _, l := range strings.Split(branch, "\n") {
		l = strings.TrimRight(l, " \t\r")
		if strings.TrimSpace(l) == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if seen[l] > 0 {
			seen[l]--
			continue
		}
		out = append(out, l)
	}
	return out
}

// without is lines less those body holds, as a multiset of line text with
// trailing white space trimmed.
func without(lines, body []string) []string {
	held := map[string]int{}
	for _, l := range body {
		held[strings.TrimRight(l, " \t\r\n")]++
	}
	var out []string
	for _, l := range lines {
		if held[l] > 0 {
			held[l]--
			continue
		}
		out = append(out, l)
	}
	return out
}
