// Package conflictmark finds the lines git writes into a file it could not
// merge: the markers that open, divide, and close a conflict (S-0253,
// I-0066).
package conflictmark

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

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

// Branch lists as path:line, by path and then line, the conflict markers in
// the text files branch adds or changes since it left base, read in the
// repository at root. Acceptance refuses by it and its preview reports by it,
// so the two cannot disagree (S-0276). It reads the blobs git names, so a
// file's own leading blank lines and odd names are read as they are; binary
// files are skipped.
func Branch(r execx.Runner, root, base, branch string) ([]string, error) {
	raw, err := r.Run(root, "git", "diff", "--raw", "-z", "--no-abbrev", "--no-renames", "--diff-filter=AM", base+"..."+branch)
	if err != nil {
		return nil, err
	}
	// records are ":oldmode newmode oldsha newsha status\0path\0"
	var paths, blobs []string
	fields := strings.Split(strings.TrimRight(raw, "\x00"), "\x00")
	for i := 0; i+1 < len(fields); i += 2 {
		meta := strings.Fields(fields[i])
		if len(meta) < 5 || meta[1] == "160000" { // a submodule has no content here
			continue
		}
		paths, blobs = append(paths, fields[i+1]), append(blobs, meta[3])
	}
	if len(blobs) == 0 {
		return nil, nil
	}
	out, err := r.RunInput(root, "git", strings.Join(blobs, "\n")+"\n", "cat-file", "--batch")
	if err != nil {
		return nil, err
	}
	// each blob is "sha type size\n" then size bytes and a newline
	marks := map[string][]int{}
	for _, p := range paths {
		head, rest, ok := strings.Cut(out, "\n")
		if !ok {
			break
		}
		meta := strings.Fields(head)
		if len(meta) != 3 {
			out = rest
			continue
		}
		size, err := strconv.Atoi(meta[2])
		if err != nil {
			return nil, fmt.Errorf("read %s from %s: %q", p, branch, head)
		}
		size = min(size, len(rest))
		text := rest[:size]
		out = strings.TrimPrefix(rest[size:], "\n")
		if !strings.Contains(text, "\x00") {
			if lines := Lines(text); len(lines) > 0 {
				marks[p] = lines
			}
		}
	}
	var found []string
	for _, p := range slices.Sorted(maps.Keys(marks)) {
		for _, n := range marks[p] {
			found = append(found, fmt.Sprintf("%s:%d", p, n))
		}
	}
	return found, nil
}
