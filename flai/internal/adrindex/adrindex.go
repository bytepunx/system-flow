// Package adrindex merges the ADR index, the table of every ADR in
// <design>/adrs/README.md, so that a rebase that stops on it because both
// sides added a row resolves it itself, as it does an issue file.
package adrindex

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ErrNotMerged is returned, wrapped with the reason, when the index's
// versions do not merge without a person's judgment.
var ErrNotMerged = errors.New("ADR index not merged")

var (
	headerRe    = regexp.MustCompile(`^\| *ADR *\| *Title *\| *Status *\| *$`)
	separatorRe = regexp.MustCompile(`^\|[-:| ]+\| *$`)
	rowRe       = regexp.MustCompile(`^\| *\[(\d{4})\]\(([^)]+)\)`)
)

// Path is the index's path relative to the repository root, slash-separated,
// as git names it: README.md in the ADR folder of the design layout, or of
// design when the manifest names none.
func Path(r *workitem.Repo) string {
	design := strings.TrimSuffix(r.Manifest.Layout["design"], "/")
	if design == "" {
		design = "design"
	}
	return filepath.ToSlash(filepath.Join(design, "adrs", "README.md"))
}

// index is one version of the index split at its table: the text above the
// table, the header and its separator line, the rows by the file each links,
// and the text below the table. The texts end each of their lines with a
// newline; a row is its line without one.
type index struct {
	above, header, below string
	rows                 map[string]string
}

// parse splits text, the side version of the index, at its table.
func parse(side, text string) (index, error) {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	h := -1
	for i, l := range lines {
		if headerRe.MatchString(l) {
			h = i
			break
		}
	}
	if h < 0 {
		return index{}, fmt.Errorf("%w: the %s version has no table headed | ADR | Title | Status |; merge the file by hand", ErrNotMerged, side)
	}
	end := h + 1
	if end < len(lines) && separatorRe.MatchString(lines[end]) {
		end++
	}
	ix := index{above: joinLines(lines[:h]), header: joinLines(lines[h:end]), rows: map[string]string{}}
	for ; end < len(lines); end++ {
		m := rowRe.FindStringSubmatch(lines[end])
		if m == nil {
			break
		}
		if _, dup := ix.rows[m[2]]; dup {
			return index{}, fmt.Errorf("%w: the %s version has two rows for %s; merge the file by hand, keeping one", ErrNotMerged, side, m[2])
		}
		ix.rows[m[2]] = lines[end]
	}
	ix.below = joinLines(lines[end:])
	return ix, nil
}

// joinLines joins lines, each ended with a newline.
func joinLines(lines []string) string {
	var b strings.Builder
	for _, l := range lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	return b.String()
}

// Merge merges three versions of the index's text, base and the two sides
// that changed it, ours and theirs, row by row, and returns the merged text;
// it returns an error wrapping ErrNotMerged when they do not merge.
func Merge(base, ours, theirs string) (string, error) {
	var v [3]index
	for i, side := range [3]struct{ name, text string }{{"base", base}, {"ours", ours}, {"theirs", theirs}} {
		ix, err := parse(side.name, side.text)
		if err != nil {
			return "", err
		}
		v[i] = ix
	}
	b, o, t := v[0], v[1], v[2]

	var conflicts []string
	str := func(name, inBase, inOurs, inTheirs string) string {
		m, ok := pick(inBase, inOurs, inTheirs)
		if !ok {
			conflicts = append(conflicts, name)
		}
		return m
	}
	m := index{
		above:  str("the text above the table", b.above, o.above, t.above),
		header: str("the table's header", b.header, o.header, t.header),
		below:  str("the text below the table", b.below, o.below, t.below),
		rows:   map[string]string{},
	}
	for _, file := range sortedFiles(b.rows, o.rows, t.rows) {
		if line := str("the row for "+file, b.rows[file], o.rows[file], t.rows[file]); line != "" {
			m.rows[file] = line
		}
	}
	if len(conflicts) > 0 {
		return "", fmt.Errorf("%w: the two sides changed %s differently; merge the file by hand", ErrNotMerged, joinAnd(conflicts))
	}
	return m.text(), nil
}

// text is the index written out, its rows by ADR number then file name.
func (ix index) text() string {
	var b strings.Builder
	b.WriteString(ix.above)
	b.WriteString(ix.header)
	for _, f := range sortedFiles(ix.rows) {
		b.WriteString(ix.rows[f])
		b.WriteByte('\n')
	}
	b.WriteString(ix.below)
	return b.String()
}

// sortedFiles is the files the rows of every version link, by the ADR number
// of their row then file name.
func sortedFiles(versions ...map[string]string) []string {
	number := map[string]string{}
	for _, rows := range versions {
		for f, line := range rows {
			number[f] = rowRe.FindStringSubmatch(line)[1]
		}
	}
	out := make([]string, 0, len(number))
	for f := range number {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if number[out[i]] != number[out[j]] {
			return number[out[i]] < number[out[j]]
		}
		return out[i] < out[j]
	})
	return out
}

// pick is the three-way merge of one value: the side that changed it from
// base, or either when both agree; false when both changed it differently.
func pick(base, ours, theirs string) (string, bool) {
	switch {
	case ours == theirs, theirs == base:
		return ours, true
	case ours == base:
		return theirs, true
	}
	return ours, false
}

// joinAnd joins parts with commas and a final "and".
func joinAnd(parts []string) string {
	if len(parts) < 2 {
		return strings.Join(parts, "")
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}
