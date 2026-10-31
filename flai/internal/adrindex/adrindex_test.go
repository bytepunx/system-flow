package adrindex

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

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

// adrText is an ADR file with the front matter given after its id.
func adrText(id, frontMatter string) string {
	return "---\nid: " + id + "\n" + frontMatter + "date: 2026-10-31\n---\n\n# " + id + " A title\n\n## Decision\n\nWe decide.\n"
}

// ADR-0228 §2: the status word, then supersedes, refines, and superseded by,
// each list in ascending order, the last two joined by "and", never "in part".
func TestRenderWritesTheStatusCell(t *testing.T) {
	cases := []struct{ name, frontMatter, cell string }{
		{"the status word alone", "title: A\nstatus: accepted\nsupersedes: []\nsuperseded_by: []\n", "accepted"},
		{"refines two", "title: A\nstatus: proposed\nsupersedes: []\nsuperseded_by: []\nrefines: [ADR-0131, ADR-0130]\n", "proposed, refines 0130 and 0131"},
		{"supersedes and refines", "title: A\nstatus: accepted\nsupersedes: [ADR-0024]\nsuperseded_by: []\nrefines: [ADR-0029, ADR-0020]\n", "accepted, supersedes 0024, refines 0020 and 0029"},
		{"refines and superseded by", "title: A\nstatus: accepted\nsupersedes: []\nsuperseded_by: [ADR-0041]\nrefines: [ADR-0029, ADR-0037]\n", "accepted, refines 0029 and 0037, superseded by 0041"},
		{"superseded whole", "title: A\nstatus: superseded\nsupersedes: []\nsuperseded_by: [ADR-0031]\n", "superseded, superseded by 0031"},
		{"all three, three numbers with no comma before and", "title: A\nstatus: accepted\nsupersedes: [ADR-0133, ADR-0085, ADR-0122]\nsuperseded_by: [ADR-0300, ADR-0299]\nrefines: [ADR-0007]\n", "accepted, supersedes 0085, 0122 and 0133, refines 0007, superseded by 0299 and 0300"},
		{"a block list and a duplicate written with fewer digits", "title: A\nstatus: accepted\nsupersedes:\n  - ADR-0016\n  - ADR-16\nsuperseded_by: []\n", "accepted, supersedes 0016"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			folder := fstest.MapFS{"0042-a.md": {Data: []byte(adrText("ADR-0042", c.frontMatter))}}
			got, err := Render(folder, header)
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			if want := header + "| [0042](0042-a.md) | A | " + c.cell + " |\n"; got != want {
				t.Errorf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

// ADR-0228 §1: one row per NNNN-*.md file in number order, the template and
// any other file left out, titles with | escaped, and the text above and
// below the table kept while the table's old rows, notes and all, go.
func TestRenderRewritesTheTableOnly(t *testing.T) {
	folder := fstest.MapFS{
		"0000-template.md": {Data: []byte(adrText("ADR-0000", "title: Template\nstatus: template\n"))},
		"0010-tenth.md":    {Data: []byte(adrText("ADR-0010", "title: \"Pipes | in a title\"\nstatus: accepted\nsupersedes: [ADR-0002]\n"))},
		"0002-second.md":   {Data: []byte(adrText("ADR-0002", "title: Second\nstatus: accepted\nsuperseded_by: [ADR-0010]\n"))},
		"README.md":        {Data: []byte("not an ADR")},
		"notes.md":         {Data: []byte("not an ADR either")},
		"0003-dir.md":      {Mode: os.ModeDir},
	}
	current := intro + header +
		"| [0002](0002-second.md) | Second | accepted, superseded in part by 0010 |\n" +
		"| [0009](0009-gone.md) | Gone | accepted |\n" +
		"| a row by hand that links nothing |\n" +
		"\nText below the table stays.\n\n| ADR | not | the table |\n"
	got, err := Render(folder, current)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	want := intro + header +
		"| [0002](0002-second.md) | Second | accepted, superseded by 0010 |\n" +
		"| [0010](0010-tenth.md) | Pipes \\| in a title | accepted, supersedes 0002 |\n" +
		"\nText below the table stays.\n\n| ADR | not | the table |\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	again, err := Render(folder, got)
	if err != nil || again != got {
		t.Errorf("a second render changes nothing: %v\n%s", err, again)
	}
}

func TestRenderOfAnIndexWithNoTableOrNoADRs(t *testing.T) {
	one := fstest.MapFS{"0001-one.md": {Data: []byte(adrText("ADR-0001", "title: One\nstatus: proposed\n"))}}
	cases := []struct {
		name    string
		folder  fstest.MapFS
		current string
		want    string
	}{
		{"no ADRs leaves the header alone", fstest.MapFS{"0000-template.md": {Data: []byte(adrText("ADR-0000", "title: T\nstatus: template\n"))}}, intro + header, intro + header},
		{"an index with no table gets one at its end", one, intro, intro + header + "| [0001](0001-one.md) | One | proposed |\n"},
		{"an empty index is the table", one, "", header + "| [0001](0001-one.md) | One | proposed |\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Render(c.folder, c.current)
			if err != nil || got != c.want {
				t.Errorf("Render: %v\ngot:\n%s\nwant:\n%s", err, got, c.want)
			}
		})
	}
}

func TestRenderRefusesFrontMatterItCannotRead(t *testing.T) {
	for _, c := range []struct{ name, text string }{
		{"no front matter", "# ADR-0005 Nothing above\n"},
		{"front matter that is not YAML", "---\ntitle: [unclosed\nstatus: accepted\n---\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := Render(fstest.MapFS{"0005-bad.md": {Data: []byte(c.text)}}, intro+header)
			if err == nil || !strings.Contains(err.Error(), "0005-bad.md") {
				t.Errorf("want an error naming the file, got %v", err)
			}
		})
	}
}

// Write writes the index only when the render changes it, and creates one
// that is not there.
func TestWrite(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "0001-one.md"), []byte(adrText("ADR-0001", "title: One\nstatus: accepted\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	index := filepath.Join(dir, "README.md")
	wrote, err := Write(dir)
	data, _ := os.ReadFile(index)
	if err != nil || !wrote || string(data) != "# Architecture Decision Records\n\n"+header+"| [0001](0001-one.md) | One | accepted |\n" {
		t.Fatalf("an absent index is created: %v %v\n%s", wrote, err, data)
	}
	if wrote, err := Write(dir); err != nil || wrote {
		t.Errorf("a current index is not written: %v %v", wrote, err)
	}
	if err := os.WriteFile(index, []byte(intro+header+"| [0001](0001-one.md) | One | proposed |\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if wrote, err := Write(dir); err != nil || !wrote {
		t.Errorf("a stale index is written: %v %v", wrote, err)
	}
	if data, _ := os.ReadFile(index); string(data) != intro+header+"| [0001](0001-one.md) | One | accepted |\n" {
		t.Errorf("the intro is kept and the row follows the file:\n%s", data)
	}
}
