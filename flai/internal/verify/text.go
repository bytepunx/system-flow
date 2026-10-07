package verify

import (
	"fmt"
	"strconv"
	"strings"
)

// Text is the result as flai test prints it: a line for each tier, its
// state, name, and duration, then each finding as path:line name: message,
// a message's further lines indented under it.
func (r Result) Text() string {
	if len(r.Tiers) == 0 {
		return "no tier selects these paths\n"
	}
	var b strings.Builder
	for _, t := range r.Tiers {
		if t.State == NotReached {
			fmt.Fprintf(&b, "%s %s\n", t.State, t.Name)
		} else {
			fmt.Fprintf(&b, "%s %s (%s)\n", t.State, t.Name, t.Duration)
		}
		for _, f := range t.Findings {
			b.WriteString(f.text())
		}
		if t.Omitted > 0 {
			fmt.Fprintf(&b, "… %d more findings left out\n", t.Omitted)
		}
	}
	return b.String()
}

// text is the finding as a line, path:line name: message, with the
// message's further lines indented under it.
func (f Finding) text() string {
	loc := f.Path
	if f.Line > 0 && loc != "" {
		loc += ":" + strconv.Itoa(f.Line)
	}
	head := strings.TrimSpace(loc + " " + f.Name)
	lines := strings.Split(f.Message, "\n")
	switch {
	case f.Message == "":
	case head == "":
		head = lines[0]
	default:
		head += ": " + lines[0]
	}
	var b strings.Builder
	b.WriteString(head + "\n")
	if f.Message != "" {
		for _, l := range lines[1:] {
			b.WriteString("    " + l + "\n")
		}
	}
	return b.String()
}
