// Package adrindex renders the ADR index, the table of every ADR in
// <design>/adrs/README.md, from the ADR files (ADR-0228), tells how an index
// differs from its render, and renders it at a rebase that stopped on it, so
// that the rebase writes it again from the rebased tree's ADR files.
package adrindex

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var (
	headerRe = regexp.MustCompile(`^\| *ADR *\| *Title *\| *Status *\| *$`)
	rowRe    = regexp.MustCompile(`^\| *\[(\d{4})\]\(([^)]+)\)`)
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

// Source says whether p, relative to the repository root and slash-separated
// as git names it, is an ADR file the index is rendered from: an NNNN-*.md
// file beside the index, 0000-template.md left out.
func Source(r *workitem.Repo, p string) bool {
	m := adrFile.FindStringSubmatch(path.Base(p))
	return m != nil && m[1] != "0000" && path.Dir(p) == path.Dir(Path(r))
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
	file := filepath.Join(dir, "README.md")
	current, err := os.ReadFile(file)
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
	return true, atomicfile.WriteFile(file, []byte(text), 0o644) // a rename, which the watcher sees at once (S-0606)
}

// RenderStopped is the index rendered at a rebase that stopped on it, from
// the ADR files at the root of folder, the rebased tree's at the stop, and
// the stop's three versions of the index, base, ours, and theirs (ADR-0228):
// each version rendered, so that the two sides differ only in the text
// around the table, and the side that changed that text from base taken. It
// fails when both sides changed the text around the table differently, or
// when an ADR's front matter cannot be read.
func RenderStopped(folder fs.FS, base, ours, theirs string) (string, error) {
	var v [3]string
	for i, text := range [3]string{base, ours, theirs} {
		r, err := Render(folder, text)
		if err != nil {
			return "", err
		}
		v[i] = r
	}
	text, ok := pick(v[0], v[1], v[2])
	if !ok {
		return "", errors.New("the two sides changed the text around the table differently; merge that text by hand, then run flai adr index")
	}
	return text, nil
}

// Difference is how an index differs from its render: the line of the index
// where they first differ, counted from 1 and at most its last line, or 0
// when they do not, and, each by ADR number in ascending order, the rows the
// render writes differently (a file linked twice among them), the rows it
// adds, and the rows it drops, which link a file the folder does not have.
type Difference struct {
	Line                    int
	Changed, Missing, Extra []string
}

// Diff is how current, the index's text, differs from rendered, its render.
// An index with no table whose render has no rows does not differ: the
// render would only add an empty table to it.
func Diff(current, rendered string) Difference {
	var d Difference
	have, want := rows(current), rows(rendered)
	if current == rendered || (len(want) == 0 && !slices.ContainsFunc(strings.Split(current, "\n"), headerRe.MatchString)) {
		return d
	}
	cl, rl := strings.Split(current, "\n"), strings.Split(rendered, "\n")
	for d.Line < len(cl) && d.Line < len(rl) && cl[d.Line] == rl[d.Line] {
		d.Line++
	}
	d.Line = min(d.Line+1, max(1, len(strings.Split(strings.TrimSuffix(current, "\n"), "\n"))))
	for file, lines := range want {
		switch got := have[file]; {
		case len(got) == 0:
			d.Missing = append(d.Missing, number(lines[0]))
		case !slices.Equal(got, lines):
			d.Changed = append(d.Changed, number(lines[0]))
		}
	}
	for file, lines := range have {
		if len(want[file]) == 0 {
			d.Extra = append(d.Extra, number(lines[0]))
		}
	}
	for _, l := range []*[]string{&d.Changed, &d.Missing, &d.Extra} {
		slices.Sort(*l)
		*l = slices.Compact(*l)
	}
	return d
}

// rows is every line of text that is an index row, by the file it links.
func rows(text string) map[string][]string {
	out := map[string][]string{}
	for _, l := range strings.Split(text, "\n") {
		if m := rowRe.FindStringSubmatch(l); m != nil {
			out[m[2]] = append(out[m[2]], l)
		}
	}
	return out
}

// number is the ADR number of row, an index row.
func number(row string) string { return rowRe.FindStringSubmatch(row)[1] }

// joinLines joins lines, each ended with a newline.
func joinLines(lines []string) string {
	var b strings.Builder
	for _, l := range lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	return b.String()
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
