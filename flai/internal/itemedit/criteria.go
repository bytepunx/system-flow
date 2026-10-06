package itemedit

import (
	"regexp"
	"strings"
)

// Criterion is one checkbox of a body's acceptance criteria.
type Criterion struct {
	N      int    `json:"n"`    // 1-based, in the order the boxes appear
	Text   string `json:"text"` // the words after the box
	Ticked bool   `json:"ticked"`
}

var (
	criteriaHeading = regexp.MustCompile(`(?m)^## Acceptance criteria[ \t]*$`)
	nextSection     = regexp.MustCompile(`(?m)^## `)
	// box is flai's checkbox, the one check and the rules count, on one line;
	// the mark is its second group
	box = regexp.MustCompile(`^(\s*- \[)([ xX])\] (\S.*)$`)
)

// mark is one box of the acceptance criteria and the offset of its mark in
// the body.
type mark struct {
	Criterion
	at int
}

// marks finds the boxes under "## Acceptance criteria", up to the next "## "
// heading or the end, nested ones included, in document order.
func marks(body string) []mark {
	loc := criteriaHeading.FindStringIndex(body)
	if loc == nil {
		return nil
	}
	start, end := loc[1], len(body)
	if next := nextSection.FindStringIndex(body[start:]); next != nil {
		end = start + next[0]
	}
	var out []mark
	for off := start; off < end; {
		line := body[off:end]
		if i := strings.IndexByte(line, '\n'); i >= 0 {
			line = line[:i]
		}
		if m := box.FindStringSubmatchIndex(line); m != nil {
			out = append(out, mark{
				Criterion: Criterion{N: len(out) + 1, Text: strings.TrimSpace(line[m[6]:m[7]]), Ticked: line[m[4]] != ' '},
				at:        off + m[4],
			})
		}
		off += len(line) + 1
	}
	return out
}

// Criteria lists a body's acceptance criteria, numbered from 1; none is an
// empty list.
func Criteria(body string) []Criterion {
	out := []Criterion{}
	for _, m := range marks(body) {
		out = append(out, m.Criterion)
	}
	return out
}

// Tick returns body with the acceptance criteria numbered in tick ticked and
// those in untick unticked, and every other byte as it was. A box already
// ticked as [X] and ticked again keeps its capital.
func Tick(body string, tick, untick []int) (string, error) {
	ms := marks(body)
	if len(ms) == 0 {
		return "", invalid("there are no acceptance criteria with a checkbox to tick: write them under ## Acceptance criteria as lines \"- [ ] what is true when it is done\"")
	}
	if len(tick) == 0 && len(untick) == 0 {
		return "", invalid("name the acceptance criteria to tick or untick by their numbers, which flai criteria list shows")
	}
	set := map[int]bool{}
	for _, l := range []struct {
		ns []int
		to bool
	}{{tick, true}, {untick, false}} {
		for _, n := range l.ns {
			if n < 1 || n > len(ms) {
				return "", invalid("there is no acceptance criterion %d: there are %d, numbered 1 to %d as flai criteria list shows them", n, len(ms), len(ms))
			}
			if to, ok := set[n]; ok && to != l.to {
				return "", invalid("acceptance criterion %d is both to tick and to untick: name it once", n)
			}
			set[n] = l.to
		}
	}
	b := []byte(body)
	for _, m := range ms {
		to, ok := set[m.N]
		switch {
		case !ok || to == m.Ticked:
		case to:
			b[m.at] = 'x'
		default:
			b[m.at] = ' '
		}
	}
	return string(b), nil
}
