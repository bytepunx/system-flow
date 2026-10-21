package verify

import (
	"fmt"
	"strconv"
	"strings"
)

// Text is the result as flai test prints it: a line for each tier, its
// state, name, and duration, or for a skipped tier the tier that covers it,
// then each finding as path:line name: message, a message's further lines
// indented under it.
func (r Result) Text() string {
	if len(r.Tiers) == 0 {
		return "no tier selects these paths\n"
	}
	var b strings.Builder
	for _, t := range r.Tiers {
		switch t.State {
		case NotReached:
			fmt.Fprintf(&b, "%s %s\n", t.State, t.Name)
		case Skipped:
			fmt.Fprintf(&b, "%s %s: covered by %s\n", t.State, t.Name, t.CoveredBy)
		default:
			fmt.Fprintf(&b, "%s %s (%s)\n", t.State, t.Name, t.Duration)
		}
		for _, f := range t.Findings {
			b.WriteString(f.Text())
		}
		if t.Omitted > 0 {
			fmt.Fprintf(&b, "… %d more findings left out\n", t.Omitted)
		}
	}
	return b.String()
}

// StepFindings are the step's findings and how many were left out, as flai
// test prints a tier's, with indent before every line.
func StepFindings(s Step, indent string) string {
	var b strings.Builder
	for _, f := range s.Findings {
		for _, line := range strings.SplitAfter(f.Text(), "\n") {
			if line != "" {
				b.WriteString(indent + line)
			}
		}
	}
	if s.Omitted > 0 {
		fmt.Fprintf(&b, "%s… %d more findings left out\n", indent, s.Omitted)
	}
	return b.String()
}

// Text is the finding as a line, path:line name: message, with the
// message's further lines indented under it.
func (f Finding) Text() string {
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
