// Package conflictmark finds the lines git writes into a file it could not
// merge: the markers that open, divide, and close a conflict (S-0253,
// I-0066).
package conflictmark

import "strings"

// markers are the line openings git writes around a conflict: ours, the
// merge base (diff3), and theirs. The divider, seven equals signs, is left
// out: alone on a line it is a markdown setext heading underline, and every
// conflict has an opening and a closing marker besides.
var markers = []string{"<<<<<<<", "|||||||", ">>>>>>>"}

// Lines returns the 1-based numbers of the lines of text that are conflict
// markers, in order.
func Lines(text string) []int {
	var out []int
	for i, line := range strings.Split(text, "\n") {
		if IsMarker(strings.TrimSuffix(line, "\r")) {
			out = append(out, i+1)
		}
	}
	return out
}

// IsMarker reports whether line, without its line ending, opens or closes a
// conflict or opens its merge base.
func IsMarker(line string) bool {
	for _, m := range markers {
		if rest, ok := strings.CutPrefix(line, m); ok && (rest == "" || rest[0] == ' ') {
			return true
		}
	}
	return false
}
