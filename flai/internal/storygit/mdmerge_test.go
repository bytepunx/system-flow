package storygit

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

// md joins lines into a file ending in a newline.
func md(lines ...string) string { return strings.Join(lines, "\n") + "\n" }

func TestMergeMarkdownSettles(t *testing.T) {
	doc := md("# T", "", "## A", "", "a1", "", "## B", "", "b1")
	fm := func(updated, status string, body ...string) string {
		return md(append([]string{"---", "title: X", "updated: " + updated, "status: " + status, "---", ""}, body...)...)
	}
	for _, tc := range []struct {
		name, base, ours, theirs, want string
	}{
		{
			name:   "a section only ours changed takes ours",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1 ours", "", "## B", "", "b1"),
			theirs: doc,
			want:   md("# T", "", "## A", "", "a1 ours", "", "## B", "", "b1"),
		},
		{
			name:   "a section only theirs changed takes theirs",
			base:   doc,
			ours:   doc,
			theirs: md("# T", "", "## A", "", "a1", "", "## B", "", "b1 theirs"),
			want:   md("# T", "", "## A", "", "a1", "", "## B", "", "b1 theirs"),
		},
		{
			name:   "each side changes its own section",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1 ours", "", "## B", "", "b1"),
			theirs: md("# T", "", "## A", "", "a1", "", "## B", "", "b1 theirs"),
			want:   md("# T", "", "## A", "", "a1 ours", "", "## B", "", "b1 theirs"),
		},
		{
			name:   "identical changes appear once",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1 both", "", "## B", "", "b1"),
			theirs: md("# T", "", "## A", "", "a1 both", "", "## B", "", "b1"),
			want:   md("# T", "", "## A", "", "a1 both", "", "## B", "", "b1"),
		},
		{
			name:   "both add a section at the end, ours first, a blank line between",
			base:   md("# T", "", "## A", "", "a1"),
			ours:   md("# T", "", "## A", "", "a1", "", "## O", "", "o1"),
			theirs: md("# T", "", "## A", "", "a1", "", "## H", "", "h1"),
			want:   md("# T", "", "## A", "", "a1", "", "## O", "", "o1", "", "## H", "", "h1"),
		},
		{
			name:   "both add the same section once",
			base:   md("# T", "", "## A", "", "a1"),
			ours:   md("# T", "", "## A", "", "a1", "", "## N", "", "n1"),
			theirs: md("# T", "", "## A", "", "a1", "", "## N", "", "n1"),
			want:   md("# T", "", "## A", "", "a1", "", "## N", "", "n1"),
		},
		{
			name:   "both add a section of one heading with different text, each once, ours first",
			base:   md("# T", "", "## A", "", "a1"),
			ours:   md("# T", "", "## A", "", "a1", "", "## N", "", "n ours"),
			theirs: md("# T", "", "## A", "", "a1", "", "## N", "", "n theirs"),
			want:   md("# T", "", "## A", "", "a1", "", "## N", "", "n ours", "", "## N", "", "n theirs"),
		},
		{
			name:   "both add a section in the middle at one place, ours first",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1", "", "## O", "", "o1", "", "## B", "", "b1"),
			theirs: md("# T", "", "## A", "", "a1", "", "## H", "", "h1", "", "## B", "", "b1"),
			want:   md("# T", "", "## A", "", "a1", "", "## O", "", "o1", "", "## H", "", "h1", "", "## B", "", "b1"),
		},
		{
			name:   "theirs adds a section where it put it while ours changes another",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1", "", "## B", "", "b1 ours"),
			theirs: md("# T", "", "## H", "", "h1", "", "## A", "", "a1", "", "## B", "", "b1"),
			want:   md("# T", "", "## H", "", "h1", "", "## A", "", "a1", "", "## B", "", "b1 ours"),
		},
		{
			name:   "both add a bullet at the end of the same list, ours first",
			base:   md("# T", "", "## L", "", "- a", "- b", "", "## M", "", "m"),
			ours:   md("# T", "", "## L", "", "- a", "- b", "- c", "", "## M", "", "m"),
			theirs: md("# T", "", "## L", "", "- a", "- b", "- d", "", "## M", "", "m"),
			want:   md("# T", "", "## L", "", "- a", "- b", "- c", "- d", "", "## M", "", "m"),
		},
		{
			name:   "adjacent line edits in one section are both kept",
			base:   md("# T", "", "## A", "", "one", "two", "three"),
			ours:   md("# T", "", "## A", "", "one ours", "two", "three"),
			theirs: md("# T", "", "## A", "", "one", "two theirs", "three"),
			want:   md("# T", "", "## A", "", "one ours", "two theirs", "three"),
		},
		{
			name:   "a section one side removes and the other leaves is removed",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1"),
			theirs: md("# T", "", "## A", "", "a1 theirs", "", "## B", "", "b1"),
			want:   md("# T", "", "## A", "", "a1 theirs"),
		},
		{
			name:   "a section both remove is removed",
			base:   doc,
			ours:   md("# T", "", "## A", "", "a1"),
			theirs: md("# T", "", "## A", "", "a1"),
			want:   md("# T", "", "## A", "", "a1"),
		},
		{
			name:   "sections of one heading under different parents are kept apart",
			base:   md("# T", "## A", "### X", "x", "## B", "### X", "x"),
			ours:   md("# T", "## A", "### X", "x ours", "## B", "### X", "x"),
			theirs: md("# T", "## A", "### X", "x", "## B", "### X", "x theirs"),
			want:   md("# T", "## A", "### X", "x ours", "## B", "### X", "x theirs"),
		},
		{
			name:   "the text before the first heading merges by lines",
			base:   md("intro one", "intro two", "", "# T"),
			ours:   md("intro one ours", "intro two", "", "# T"),
			theirs: md("intro one", "intro two theirs", "", "# T"),
			want:   md("intro one ours", "intro two theirs", "", "# T"),
		},
		{
			name:   "a heading inside a fenced code block does not cut a section",
			base:   md("# T", "", "## A", "", "```", "## B", "```", "", "## B", "", "b1"),
			ours:   md("# T", "", "## A", "", "```", "## B", "```", "", "## B", "", "b1 ours"),
			theirs: md("# T", "", "## A", "", "## B", "", "b1"),
			want:   md("# T", "", "## A", "", "## B", "", "b1 ours"),
		},
		{
			name:   "a heading inside a tilde fence does not cut a section",
			base:   md("# T", "", "~~~~", "# X", "~~~", "~~~~", "x"),
			ours:   md("# T", "", "~~~~", "# X", "~~~", "~~~~", "x ours"),
			theirs: md("# T", "", "~~~~", "# Y", "~~~", "~~~~", "x"),
			want:   md("# T", "", "~~~~", "# Y", "~~~", "~~~~", "x ours"),
		},
		{
			name:   "both bump updated, theirs later",
			base:   fm("2026-10-01", "active", "# T"),
			ours:   fm("2026-10-18", "active", "# T"),
			theirs: fm("2026-10-19", "active", "# T"),
			want:   fm("2026-10-19", "active", "# T"),
		},
		{
			name:   "both bump updated, ours later",
			base:   fm("2026-10-01", "active", "# T"),
			ours:   fm("2026-10-19", "active", "# T"),
			theirs: fm("2026-10-18", "active", "# T"),
			want:   fm("2026-10-19", "active", "# T"),
		},
		{
			name:   "both bump a quoted RFC 3339 updated, the later",
			base:   fm(`"2026-10-01T00:00:00Z"`, "active", "# T"),
			ours:   fm(`"2026-10-19T07:00:00Z"`, "active", "# T"),
			theirs: fm(`"2026-10-19T08:00:00Z"`, "active", "# T"),
			want:   fm(`"2026-10-19T08:00:00Z"`, "active", "# T"),
		},
		{
			name:   "front matter keys each side changes, and updated, with the body",
			base:   fm("2026-10-01", "active", "# T", "", "t"),
			ours:   fm("2026-10-18", "draft", "# T", "", "t ours"),
			theirs: fm("2026-10-19", "active", "# T", "", "t", "", "## H", "", "h"),
			want:   fm("2026-10-19", "draft", "# T", "", "t ours", "", "## H", "", "h"),
		},
		{
			name:   "keys theirs adds follow ours' keys",
			base:   md("---", "title: X", "updated: 2026-10-01", "---", "# T"),
			ours:   md("---", "title: X", "topics: [cli]", "updated: 2026-10-01", "---", "# T"),
			theirs: md("---", "title: X", "updated: 2026-10-01", "owner: me", "---", "# T"),
			want:   md("---", "title: X", "topics: [cli]", "updated: 2026-10-01", "owner: me", "---", "# T"),
		},
		{
			name:   "a multi-line value is one key",
			base:   md("---", "title: X", "tags:", "  - a", "status: active", "---", "# T"),
			ours:   md("---", "title: X", "tags:", "  - a", "  - b", "status: active", "---", "# T"),
			theirs: md("---", "title: Y", "tags:", "  - a", "status: active", "---", "# T"),
			want:   md("---", "title: Y", "tags:", "  - a", "  - b", "status: active", "---", "# T"),
		},
		{
			name:   "a file without a final newline stays without one",
			base:   "# T\n\nt",
			ours:   "# T\n\nt ours",
			theirs: "# T\n\nt",
			want:   "# T\n\nt ours",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MergeMarkdown(tc.base, tc.ours, tc.theirs)
			if err != nil {
				t.Fatalf("MergeMarkdown: %v", err)
			}
			if got != tc.want {
				t.Errorf("merged %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMergeMarkdownLeavesWhatItCannotSettle(t *testing.T) {
	doc := md("---", "title: X", "status: active", "---", "", "# T", "", "## A", "", "one", "two", "three", "", "## B", "", "b1")
	for _, tc := range []struct {
		name, ours, theirs, where string
	}{
		{
			name:   "the same line changed differently",
			ours:   strings.Replace(doc, "two", "two ours", 1),
			theirs: strings.Replace(doc, "two", "two theirs", 1),
			where:  `section "T > A": both sides change the same lines`,
		},
		{
			name:   "lines inserted inside lines the other changes",
			ours:   strings.Replace(doc, "one\ntwo", "one ours\ntwo ours", 1),
			theirs: strings.Replace(doc, "one\ntwo", "one\nnew\ntwo", 1),
			where:  `section "T > A": both sides change the same lines`,
		},
		{
			name:   "a section ours removes and theirs changes",
			ours:   strings.Replace(doc, "three\n\n## B\n\nb1\n", "three\n", 1),
			theirs: strings.Replace(doc, "b1", "b1 theirs", 1),
			where:  `section "T > B": one side removes it and the other changes it`,
		},
		{
			name:   "a section theirs removes and ours changes",
			ours:   strings.Replace(doc, "b1", "b1 ours", 1),
			theirs: strings.Replace(doc, "three\n\n## B\n\nb1\n", "three\n", 1),
			where:  `section "T > B": one side removes it and the other changes it`,
		},
		{
			name:   "a front matter key other than updated changed differently",
			ours:   strings.Replace(doc, "status: active", "status: draft", 1),
			theirs: strings.Replace(doc, "status: active", "status: superseded", 1),
			where:  `front matter key "status": both sides change it differently`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MergeMarkdown(doc, tc.ours, tc.theirs)
			if !errors.Is(err, ErrNotSettled) {
				t.Fatalf("MergeMarkdown = %q, %v; want an error wrapping ErrNotSettled", got, err)
			}
			if !strings.Contains(err.Error(), tc.where) {
				t.Errorf("error %q does not name %q", err, tc.where)
			}
		})
	}
}

func TestCutSectionsKeysByHeadingPathOutsideFences(t *testing.T) {
	doc := parseMarkdown(md("---", "title: X", "---", "intro", "# T", "```go", "# not a heading", "```", "## A", "#hashtag", "    # indented code", "## A", "### X ###", "x"))
	var where []string
	for _, s := range doc.sections {
		where = append(where, s.where)
	}
	want := []string{
		"the text before the first heading",
		`section "T"`,
		`section "T > A"`,
		`section "T > A" (occurrence 2)`,
		`section "T > A > X"`,
	}
	if !slices.Equal(where, want) {
		t.Errorf("sections %q, want %q", where, want)
	}
	if doc.front == nil || len(doc.front.entries) != 1 || doc.front.entries[0].where != `front matter key "title"` {
		t.Errorf("front matter %+v, want the one key title", doc.front)
	}
}
