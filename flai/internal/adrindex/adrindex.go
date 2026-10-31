// Package adrindex renders the ADR index, the table of every ADR in
// <design>/adrs/README.md, from the ADR files (ADR-0228), and merges it so
// that a rebase that stops on it because both sides added a row resolves it
// itself, as it does an issue file.
package adrindex

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"

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

// tableHeader is the table's header and separator lines, as the render writes them.
const tableHeader = "| ADR | Title | Status |\n|-----|-------|--------|\n"

// newIndex is what Write puts above the table of an index that does not exist.
const newIndex = "# Architecture Decision Records\n"

// adrFile matches an ADR's file name, its number in the first group.
var adrFile = regexp.MustCompile(`^(\d{4})-.+\.md$`)

// entry is one ADR file as its row shows it: the number and file name, and
// the front matter's title, status, and the ADRs it supersedes, refines, and
// is superseded by.
type entry struct {
	Number       string   `yaml:"-"`
	File         string   `yaml:"-"`
	Title        string   `yaml:"title"`
	Status       string   `yaml:"status"`
	Supersedes   []string `yaml:"supersedes"`
	Refines      []string `yaml:"refines"`
	SupersededBy []string `yaml:"superseded_by"`
}

// read reads every NNNN-*.md file at the root of folder, 0000-template.md
// left out, in number order then file name.
func read(folder fs.FS) ([]entry, error) {
	files, err := fs.ReadDir(folder, ".")
	if err != nil {
		return nil, err
	}
	var out []entry
	for _, f := range files {
		m := adrFile.FindStringSubmatch(f.Name())
		if f.IsDir() || m == nil || m[1] == "0000" {
			continue
		}
		data, err := fs.ReadFile(folder, f.Name())
		if err != nil {
			return nil, err
		}
		fm, _, err := workitem.SplitFrontMatter(string(data))
		if err != nil {
			return nil, fmt.Errorf("%s: %w; the ADR index is rendered from each ADR's front matter", f.Name(), err)
		}
		e := entry{Number: m[1], File: f.Name()}
		if err := yaml.Unmarshal([]byte(fm), &e); err != nil {
			return nil, fmt.Errorf("%s: front matter does not parse: %w; the ADR index is rendered from it", f.Name(), err)
		}
		out = append(out, e)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Number != out[j].Number {
			return out[i].Number < out[j].Number
		}
		return out[i].File < out[j].File
	})
	return out, nil
}

// row is the entry's line in the table, without a newline.
func (e entry) row() string {
	return fmt.Sprintf("| [%s](%s) | %s | %s |", e.Number, e.File, strings.ReplaceAll(e.Title, "|", `\|`), e.cell())
}

// cell is the status cell ADR-0228 §2 decides: the status word, then the ADRs
// the entry supersedes, refines, and is superseded by, each list when not
// empty.
func (e entry) cell() string {
	parts := []string{e.Status}
	for _, p := range []struct {
		label string
		refs  []string
	}{{"supersedes", e.Supersedes}, {"refines", e.Refines}, {"superseded by", e.SupersededBy}} {
		if nums := numbers(p.refs); len(nums) > 0 {
			parts = append(parts, p.label+" "+joinAnd(nums))
		}
	}
	return strings.Join(parts, ", ")
}

// numbers is refs, each ADR-NNNN or NNNN, as four-digit numbers without
// duplicates in ascending order; a ref that is not a number is kept as
// written, after the numbers.
func numbers(refs []string) []string {
	seen := map[string]bool{}
	var nums []int
	var other []string
	for _, r := range refs {
		r = strings.TrimSpace(r)
		s := strings.TrimPrefix(strings.ToUpper(r), "ADR-")
		key := r
		n, err := strconv.Atoi(s)
		if err == nil {
			key = fmt.Sprintf("%04d", n)
		}
		if r == "" || seen[key] {
			continue
		}
		seen[key] = true
		if err == nil {
			nums = append(nums, n)
		} else {
			other = append(other, r)
		}
	}
	sort.Ints(nums)
	out := make([]string, 0, len(nums)+len(other))
	for _, n := range nums {
		out = append(out, fmt.Sprintf("%04d", n))
	}
	return append(out, other...)
}

// table is the header and one row per entry, each line ended with a newline.
func table(entries []entry) string {
	var b strings.Builder
	b.WriteString(tableHeader)
	for _, e := range entries {
		b.WriteString(e.row())
		b.WriteByte('\n')
	}
	return b.String()
}

// splice puts tbl in place of the table in current, keeping the text above
// and below it; the table is its header line and every line after it that
// starts with |. An index with no table gets one at its end.
func splice(current, tbl string) string {
	text := strings.TrimRight(current, "\n")
	if strings.TrimSpace(text) == "" {
		return tbl
	}
	lines := strings.Split(text, "\n")
	h := -1
	for i, l := range lines {
		if headerRe.MatchString(l) {
			h = i
			break
		}
	}
	if h < 0 {
		return text + "\n\n" + tbl
	}
	end := h + 1
	for end < len(lines) && strings.HasPrefix(lines[end], "|") {
		end++
	}
	return joinLines(lines[:h]) + tbl + joinLines(lines[end:])
}

// Render is current, the index's text, with its table rendered again from
// the ADR files at the root of folder (ADR-0228): one row per NNNN-*.md file,
// 0000-template.md left out, in number order, each the file's link, its
// title with | escaped, and its status cell. The text above and below the
// table is kept; an index with no table gets one at its end. It fails when an
// ADR's front matter cannot be read.
func Render(folder fs.FS, current string) (string, error) {
	entries, err := read(folder)
	if err != nil {
		return "", err
	}
	return splice(current, table(entries)), nil
}

// Write renders the index, README.md in dir, the ADR folder, from the ADR
// files in it and writes it when that changes it, or when it does not exist,
// headed "# Architecture Decision Records". It reports whether it wrote.
func Write(dir string) (bool, error) {
	path := filepath.Join(dir, "README.md")
	current, err := os.ReadFile(path)
	exists := err == nil
	switch {
	case errors.Is(err, fs.ErrNotExist):
		current = []byte(newIndex)
	case err != nil:
		return false, err
	}
	text, err := Render(os.DirFS(dir), string(current))
	if err != nil {
		return false, err
	}
	if exists && text == string(current) {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(text), 0o644)
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
