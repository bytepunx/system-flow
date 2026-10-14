package adrindex

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const (
	intro  = "# Architecture Decision Records\n\nOne decision per file.\n\n"
	header = "| ADR | Title | Status |\n|-----|-------|--------|\n"
	r0154  = "| [0154](0154-a-reply-may-settle.md) | A reply may settle a conversation | accepted |\n"
	r0154s = "| [0154](0154-a-reply-may-settle.md) | A reply may settle a conversation | accepted, superseded by 0157 |\n"
	r0154r = "| [0154](0154-a-reply-may-settle.md) | A reply may settle a conversation | accepted, refines 0019 |\n"
	r0155  = "| [0155](0155-adr-index-merges.md) | The ADR index merges on a rebase | accepted |\n"
	r0156  = "| [0156](0156-a-close-out-syncs.md) | A close-out that main outran syncs again | accepted |\n"
	r0157  = "| [0157](0157-supersedes-0154.md) | A reply settles nothing | accepted, supersedes 0154 |\n"
)

// doc is an index with the intro above its table and rows in the order given.
func doc(rows ...string) string { return intro + header + strings.Join(rows, "") }

func TestMergeMerges(t *testing.T) {
	cases := []struct {
		name, base, ours, theirs, want string
	}{
		{"each side appends a different row (I-0099)",
			doc(r0154), doc(r0154, r0157), doc(r0154, r0155),
			doc(r0154, r0155, r0157)},
		{"one side notes a row superseded while the other appends a row",
			doc(r0154, r0156), doc(r0154s, r0156, r0157), doc(r0154, r0155, r0156),
			doc(r0154s, r0155, r0156, r0157)},
		{"both sides append the same row",
			doc(r0154), doc(r0154, r0155), doc(r0154, r0155),
			doc(r0154, r0155)},
		{"one side removes a row the other leaves alone",
			doc(r0154, r0155), doc(r0154), doc(r0154, r0155, r0156),
			doc(r0154, r0156)},
		{"one side changes the text above the table",
			doc(r0154), "# ADRs\n\n" + header + r0154, doc(r0154, r0155),
			"# ADRs\n\n" + header + r0154 + r0155},
		{"text below the table is kept and changed by one side",
			doc(r0154) + "\nA note.\n", doc(r0154) + "\nA longer note.\n", doc(r0154, r0155) + "\nA note.\n",
			doc(r0154, r0155) + "\nA longer note.\n"},
		{"the trailing newlines are one",
			doc(r0154), doc(r0154, r0155) + "\n\n", doc(r0154),
			doc(r0154, r0155)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Merge(c.base, c.ours, c.theirs)
			if err != nil {
				t.Fatalf("Merge: %v", err)
			}
			if got != c.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, c.want)
			}
		})
	}
}

func TestMergeRefuses(t *testing.T) {
	cases := []struct {
		name, base, ours, theirs, names string
	}{
		{"both sides change one row differently",
			doc(r0154), doc(r0154s), doc(r0154r),
			"the row for 0154-a-reply-may-settle.md"},
		{"both sides change the text above the table differently",
			doc(r0154), "# ADRs\n\n" + header + r0154, "# Decisions\n\n" + header + r0154,
			"the text above the table"},
		{"one side removes a row the other changes",
			doc(r0154, r0155), doc(r0155), doc(r0154s, r0155),
			"the row for 0154-a-reply-may-settle.md"},
		{"a version has no table",
			doc(r0154), intro, doc(r0154, r0155),
			"the ours version has no table"},
		{"a version has two rows for one file",
			doc(r0154), doc(r0154, r0154s), doc(r0154),
			"two rows for 0154-a-reply-may-settle.md"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := Merge(c.base, c.ours, c.theirs)
			if !errors.Is(err, ErrNotMerged) || !strings.Contains(err.Error(), c.names) {
				t.Errorf("want ErrNotMerged naming %s, got %v:\n%s", c.names, err, out)
			}
		})
	}
}

func TestPath(t *testing.T) {
	cases := []struct {
		name, design, want string
	}{
		{"no design layout", "", "design/adrs/README.md"},
		{"the design layout", "design", "design/adrs/README.md"},
		{"a design layout of its own, with a trailing slash", "notes/arch/", "notes/arch/adrs/README.md"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := &workitem.Repo{Manifest: manifest.Manifest{Layout: map[string]string{"design": c.design}}}
			if got := Path(r); got != c.want {
				t.Errorf("Path = %q, want %q", got, c.want)
			}
		})
	}
}
