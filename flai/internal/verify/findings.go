package verify

import (
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// The findings of the lint a fix left, as data (S-0378): flai accept --json
// prints them in its refusal, the dashboard opens each file they name, and
// the merge queue reports them before it builds a batch (S-0549).

// The checkouts a finding of unfixed lint is to be corrected in (S-0378): the
// main checkout's uncommitted files, or the story branch's in its worktree.
const (
	CheckoutMain  = "main"
	CheckoutStory = "story"
)

// FileFinding is one finding of unfixed lint and where to correct it.
type FileFinding struct {
	Tier string `json:"tier"`
	// Checkout is CheckoutMain or CheckoutStory; Story names the story whose
	// branch holds the file, for CheckoutStory.
	Checkout string `json:"checkout"`
	Story    string `json:"story,omitempty"`
	// Path is relative to the repository's root.
	Path   string `json:"path"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
	// Level is the tool's severity when it gives one: error or warning.
	Level string `json:"level,omitempty"`
	// Rule is the rule's ID, MD009, or the check's or linter's name; Alias
	// is the rule's first name, no-trailing-spaces, when the tool gives one.
	Rule    string `json:"rule,omitempty"`
	Alias   string `json:"alias,omitempty"`
	Message string `json:"message"`
}

// FileFindings are the findings of res's failing tiers, run in the checkout
// at root, as data, in their order, each to be corrected in checkout and, for
// CheckoutStory, story's branch; none, not nil, when they give none.
func FileFindings(root, checkout, story string, res Result) []FileFinding {
	got := []FileFinding{}
	for _, t := range res.Tiers {
		if t.State != Failed {
			continue
		}
		base := FileFinding{Tier: t.Name, Checkout: checkout, Story: story}
		for _, f := range t.Findings {
			got = append(got, fileFindings(root, base, t, f)...)
		}
	}
	return got
}

// markdownlintLine is a line of markdownlint-cli2's default output:
// path:line[:column] [error|warning] MDnnn[/alias...] message.
var markdownlintLine = regexp.MustCompile(`^(.+?):(\d+)(?::(\d+))?\s+(?:(error|warning)\s+)?(MD\d{3})((?:/[A-Za-z0-9_-]+)*)\s+(.+)$`)

// fileFindings are the findings f of tier t gives as data, each base with
// where it is: f itself when it names a file, or each markdownlint-cli2 line
// of a plain tier's output, none when it holds no such line.
func fileFindings(root string, base FileFinding, t TierResult, f Finding) []FileFinding {
	if f.Path != "" {
		u := base
		u.Path, u.Line, u.Message = f.Path, f.Line, strings.TrimSpace(f.Message)
		if f.Name != t.Name {
			u.Rule = f.Name
		}
		return []FileFinding{u}
	}
	var got []FileFinding
	for _, line := range strings.Split(f.Message, "\n") {
		m := markdownlintLine.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		u := base
		u.Path = checkoutRel(root, t.Dir, m[1])
		u.Line, _ = strconv.Atoi(m[2])
		u.Column, _ = strconv.Atoi(m[3])
		u.Level, u.Rule, u.Message = m[4], m[5], m[7]
		if aliases := strings.TrimPrefix(m[6], "/"); aliases != "" {
			u.Alias, _, _ = strings.Cut(aliases, "/")
		}
		got = append(got, u)
	}
	return got
}

// checkoutRel is a path a tool run in dir of the checkout at root printed,
// relative to root: an absolute one below root made relative, and a relative
// one joined to dir.
func checkoutRel(root, dir, p string) string {
	if filepath.IsAbs(p) {
		if r, err := filepath.Rel(root, p); err == nil && r != ".." && !strings.HasPrefix(r, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(r)
		}
		return filepath.ToSlash(p)
	}
	return path.Join(filepath.ToSlash(dir), filepath.ToSlash(p))
}

// WithoutLeft is paths without each path left names, or that lies below a
// folder left names, as git status names an untracked folder.
func WithoutLeft(paths, left []string) []string {
	if len(left) == 0 {
		return paths
	}
	var out []string
	for _, p := range paths {
		if !IsLeft(p, left) {
			out = append(out, p)
		}
	}
	return out
}

// IsLeft reports whether p is, or lies below, a path left names.
func IsLeft(p string, left []string) bool {
	for _, l := range left {
		l = strings.TrimSuffix(l, "/")
		if p == l || strings.HasPrefix(p, l+"/") {
			return true
		}
	}
	return false
}
