// Package mdlint checks markdown against the rules of the project's
// markdownlint configuration that what flai writes can break (S-0179).
// flai writes work items, threads, and narratives in the main
// checkout, where no story's lint runs (ADR-0019), so it lints them itself:
// flai check reports what it finds, and the commands that write refuse what
// would bring a finding.
//
// It implements a subset of markdownlint's rules, with markdownlint's
// options and inline markdownlint-disable comments, and matches what
// markdownlint-cli2 reports for them on the fixtures in testdata. Where it
// cannot judge a case it reports nothing: the project's own lint stays the
// authority.
package mdlint

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Finding is one rule breach, as markdownlint reports it.
type Finding struct {
	Line        int    `json:"line"`
	Rule        string `json:"rule"`  // MD024
	Alias       string `json:"alias"` // no-duplicate-heading
	Description string `json:"description"`
	Detail      string `json:"detail,omitempty"`
	Context     string `json:"context,omitempty"`
}

// String is markdownlint's own form, without the file and line.
func (f Finding) String() string {
	s := f.Rule + "/" + f.Alias + " " + f.Description
	if f.Detail != "" {
		s += " [" + f.Detail + "]"
	}
	if f.Context != "" {
		s += fmt.Sprintf(" [Context: %q]", f.Context)
	}
	return s
}

type rule struct {
	id, desc string
	aliases  []string
	tags     []string
	check    func(c *Config, d *doc, in *inlineOut, add adder)
}

type adder func(line int, detail, context string)

// rules are the markdownlint rules flai implements, in markdownlint's order.
var rules []rule

func init() {
	rules = []rule{
		{"MD001", "Heading levels should only increment by one level at a time", []string{"heading-increment"}, []string{"headings"}, md001},
		{"MD004", "Unordered list style", []string{"ul-style"}, []string{"bullet", "ul"}, md004},
		{"MD009", "Trailing spaces", []string{"no-trailing-spaces"}, []string{"whitespace"}, md009},
		{"MD010", "Hard tabs", []string{"no-hard-tabs"}, []string{"whitespace", "hard_tab"}, md010},
		{"MD012", "Multiple consecutive blank lines", []string{"no-multiple-blanks"}, []string{"whitespace", "blank_lines"}, md012},
		{"MD018", "No space after hash on atx style heading", []string{"no-missing-space-atx"}, []string{"headings", "atx", "spaces"}, md018},
		{"MD019", "Multiple spaces after hash on atx style heading", []string{"no-multiple-space-atx"}, []string{"headings", "atx", "spaces"}, md019},
		{"MD022", "Headings should be surrounded by blank lines", []string{"blanks-around-headings"}, []string{"headings", "blank_lines"}, md022},
		{"MD024", "Multiple headings with the same content", []string{"no-duplicate-heading"}, []string{"headings"}, md024},
		{"MD025", "Multiple top-level headings in the same document", []string{"single-title", "single-h1"}, []string{"headings"}, md025},
		{"MD026", "Trailing punctuation in heading", []string{"no-trailing-punctuation"}, []string{"headings"}, md026},
		{"MD029", "Ordered list item prefix", []string{"ol-prefix"}, []string{"ol"}, md029},
		{"MD031", "Fenced code blocks should be surrounded by blank lines", []string{"blanks-around-fences"}, []string{"code", "blank_lines"}, md031},
		{"MD034", "Bare URL used", []string{"no-bare-urls"}, []string{"links", "url"}, md034},
		{"MD036", "Emphasis used instead of a heading", []string{"no-emphasis-as-heading"}, []string{"headings", "emphasis"}, md036},
		{"MD037", "Spaces inside emphasis markers", []string{"no-space-in-emphasis"}, []string{"whitespace", "emphasis"}, md037},
		{"MD040", "Fenced code blocks should have a language specified", []string{"fenced-code-language"}, []string{"code", "language"}, md040},
		{"MD047", "Files should end with a single newline character", []string{"single-trailing-newline"}, []string{"blank_lines"}, md047},
		{"MD049", "Emphasis style", []string{"emphasis-style"}, []string{"emphasis"}, md049},
		{"MD050", "Strong style", []string{"strong-style"}, []string{"emphasis"}, md050},
	}
}

// Rules lists the IDs of the rules flai implements.
func Rules() []string {
	var out []string
	for _, r := range rules {
		out = append(out, r.id)
	}
	return out
}

// Lint checks content, a whole markdown document, against the rules the
// configuration turns on. A nil configuration finds nothing.
func (c *Config) Lint(content string) []Finding {
	if c == nil {
		return nil
	}
	d := parse(content)
	in := d.inline()
	off := disabled(d)
	var out []Finding
	for _, r := range rules {
		if !c.Enabled(r.id) {
			continue
		}
		r.check(c, d, in, func(line int, detail, context string) {
			if off(line, r) {
				return
			}
			out = append(out, Finding{Line: line + 1, Rule: r.id, Alias: r.aliases[0], Description: r.desc, Detail: detail, Context: context})
		})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Line < out[j].Line })
	return out
}

// Introduced returns the findings in after that before does not have: the
// ones a change brings. A finding is known by its rule and the text of its
// line, so findings that only moved are not new.
func (c *Config) Introduced(before, after string) []Finding {
	if c == nil {
		return nil
	}
	key := func(lines []string, f Finding) string {
		text := ""
		if f.Line-1 < len(lines) {
			text = lines[f.Line-1]
		}
		return f.Rule + "\x00" + f.Detail + "\x00" + text
	}
	bl := strings.Split(strings.ReplaceAll(before, "\r\n", "\n"), "\n")
	al := strings.Split(strings.ReplaceAll(after, "\r\n", "\n"), "\n")
	known := map[string]int{}
	for _, f := range c.Lint(before) {
		known[key(bl, f)]++
	}
	var out []Finding
	for _, f := range c.Lint(after) {
		k := key(al, f)
		if known[k] > 0 {
			known[k]--
			continue
		}
		out = append(out, f)
	}
	return out
}

// Error refuses a write whose markdown the project's lint rejects.
type Error struct {
	Path     string
	Findings []Finding
}

func (e *Error) Error() string {
	var parts []string
	for _, f := range e.Findings {
		parts = append(parts, fmt.Sprintf("line %d: %s", f.Line, f))
	}
	return fmt.Sprintf("%s would fail the project's markdown lint, so nothing was written: %s", e.Path, strings.Join(parts, "; "))
}

// Guard returns an *Error when after brings findings before did not have,
// under the configuration in root; nil when it brings none or the project
// has no markdownlint configuration.
func Guard(root, path, before, after string) error {
	c, err := Load(root)
	if err != nil {
		return err
	}
	if f := c.Introduced(before, after); len(f) > 0 {
		return &Error{Path: path, Findings: f}
	}
	return nil
}

var inlineConfigRe = regexp.MustCompile(`(?i)<!--\s*markdownlint-(disable-file|enable-file|disable-next-line|disable-line|disable|enable|capture|restore)((?:\s+[a-z0-9_-]+)*)\s*-->`)

// disabled works out the inline markdownlint comments: which rules are off on
// which line. It reports whether a rule is off on a line index.
func disabled(d *doc) func(int, rule) bool {
	all := func() map[string]bool {
		m := map[string]bool{}
		for _, r := range rules {
			m[r.id] = true
		}
		return m
	}
	apply := func(state map[string]bool, on bool, params string) map[string]bool {
		next := map[string]bool{}
		for k, v := range state {
			next[k] = v
		}
		names := strings.Fields(params)
		if len(names) == 0 {
			for _, r := range rules {
				next[r.id] = on
			}
		}
		for _, n := range names {
			for _, id := range resolve(n) {
				next[id] = on
			}
		}
		return next
	}
	state := all()
	n := len(d.lines)
	for i := d.content; i < n; i++ {
		for _, m := range inlineConfigRe.FindAllStringSubmatch(d.lines[i].raw, -1) {
			switch strings.ToLower(m[1]) {
			case "disable-file":
				state = apply(state, false, m[2])
			case "enable-file":
				state = apply(state, true, m[2])
			}
		}
	}
	per := make([]map[string]bool, n)
	captured := state
	for i := 0; i < n; i++ {
		if i >= d.content {
			for _, m := range inlineConfigRe.FindAllStringSubmatch(d.lines[i].raw, -1) {
				switch strings.ToLower(m[1]) {
				case "disable":
					state = apply(state, false, m[2])
				case "enable":
					state = apply(state, true, m[2])
				case "capture":
					captured = state
				case "restore":
					state = captured
				}
			}
		}
		per[i] = state
	}
	for i := d.content; i < n; i++ {
		for _, m := range inlineConfigRe.FindAllStringSubmatch(d.lines[i].raw, -1) {
			switch strings.ToLower(m[1]) {
			case "disable-line":
				per[i] = apply(per[i], false, m[2])
			case "disable-next-line":
				if i+1 < n {
					per[i+1] = apply(per[i+1], false, m[2])
				}
			}
		}
	}
	return func(line int, r rule) bool {
		return line >= 0 && line < n && !per[line][r.id]
	}
}
