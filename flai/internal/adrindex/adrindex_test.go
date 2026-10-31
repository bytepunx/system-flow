package adrindex

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const (
	intro  = "# Architecture Decision Records\n\nOne decision per file.\n\n"
	header = "| ADR | Title | Status |\n|-----|-------|--------|\n"
	row1   = "| [0001](0001-one.md) | One | accepted |\n"
	row2   = "| [0002](0002-two.md) | Two | accepted |\n"
	row3   = "| [0003](0003-three.md) | Three | proposed |\n"
)

// both is an ADR folder as a rebase that stopped on the index sees it: ADR
// 0001 from the base, 0002 the one side added, 0003 the other.
func both() fstest.MapFS {
	return fstest.MapFS{
		"0001-one.md":   {Data: []byte(adrText("ADR-0001", "title: One\nstatus: accepted\n"))},
		"0002-two.md":   {Data: []byte(adrText("ADR-0002", "title: Two\nstatus: accepted\n"))},
		"0003-three.md": {Data: []byte(adrText("ADR-0003", "title: Three\nstatus: proposed\n"))},
	}
}

// ADR-0228 §1: a stop on the index writes it from the rebased tree's ADR
// files, whatever rows each side's version holds, keeping the text around the
// table that one side changed.
func TestRenderStoppedWritesTheTableFromTheFolder(t *testing.T) {
	cases := []struct {
		name, base, ours, theirs, want string
	}{
		{"each side added a different row (I-0099)",
			intro + header + row1, intro + header + row1 + row2, intro + header + row1 + row3,
			intro + header + row1 + row2 + row3},
		{"both sides changed one row differently",
			intro + header + row1, intro + header + "| [0001](0001-one.md) | One | rejected |\n", intro + header + "| [0001](0001-one.md) | One | accepted, in part |\n",
			intro + header + row1 + row2 + row3},
		{"one side changed the text above the table",
			intro + header + row1, "# ADRs\n\n" + header + row1 + row2, intro + header + row1 + row3,
			"# ADRs\n\n" + header + row1 + row2 + row3},
		{"one side changed the text below it",
			intro + header + row1 + "\nA note.\n", intro + header + row1 + row2 + "\nA note.\n", intro + header + row1 + row3 + "\nA longer note.\n",
			intro + header + row1 + row2 + row3 + "\nA longer note.\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := RenderStopped(both(), c.base, c.ours, c.theirs)
			if err != nil || got != c.want {
				t.Errorf("RenderStopped: %v\ngot:\n%s\nwant:\n%s", err, got, c.want)
			}
		})
	}
}

// The text around the table both sides changed differently is a person's to
// merge, and so is an ADR whose front matter does not read.
func TestRenderStoppedRefuses(t *testing.T) {
	if _, err := RenderStopped(both(), intro+header+row1, "# ADRs\n\n"+header+row1, "# Decisions\n\n"+header+row1); err == nil || !strings.Contains(err.Error(), "text around the table") {
		t.Errorf("text changed differently on both sides: %v", err)
	}
	broken := both()
	broken["0004-broken.md"] = &fstest.MapFile{Data: []byte("# ADR-0004 No front matter\n")}
	if _, err := RenderStopped(broken, intro+header+row1, intro+header+row1+row2, intro+header+row1+row3); err == nil || !strings.Contains(err.Error(), "0004-broken.md") {
		t.Errorf("an unreadable ADR: %v", err)
	}
}

// Diff names the rows that differ, are missing, or link no file, by number,
// and the first line that differs.
func TestDiff(t *testing.T) {
	rendered := intro + header + row1 + row2 + row3 + "\nAfter.\n"
	cases := []struct {
		name, current string
		want          Difference
	}{
		{"the same", rendered, Difference{}},
		{"a row changed, one missing, one extra, one twice",
			intro + header + row1 + row1 + "| [0002](0002-two.md) | Two | proposed |\n" + "| [0009](0009-gone.md) | Gone | accepted |\n" + "\nAfter.\n",
			Difference{Line: 8, Changed: []string{"0001", "0002"}, Missing: []string{"0003"}, Extra: []string{"0009"}}},
		{"the text alone", "# ADRs\n\n" + header + row1 + row2 + row3 + "\nAfter.\n", Difference{Line: 1}},
		{"rows missing at the end, at the last line", intro + header + row1 + row2, Difference{Line: 8, Missing: []string{"0003"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Diff(c.current, rendered); !reflect.DeepEqual(got, c.want) {
				t.Errorf("Diff = %+v, want %+v", got, c.want)
			}
		})
	}
	if got := Diff("# adrs\n", "# adrs\n\n"+header); !reflect.DeepEqual(got, Difference{}) {
		t.Errorf("an index with no table and no ADRs to list does not differ: %+v", got)
	}
	if got := Diff("# adrs\n", "# adrs\n\n"+header+row1); !reflect.DeepEqual(got, Difference{Line: 1, Missing: []string{"0001"}}) {
		t.Errorf("an index with no table and an ADR to list: %+v", got)
	}
}

func TestSource(t *testing.T) {
	r := &workitem.Repo{Manifest: manifest.Manifest{Layout: map[string]string{"design": "design"}}}
	for p, want := range map[string]bool{
		"design/adrs/0228-generated.md":    true,
		"design/adrs/0000-template.md":     false,
		"design/adrs/README.md":            false,
		"design/adrs/old/0001-moved.md":    false,
		"design/system/0001-not-an-adr.md": false,
		"design/adrs/0001-notes.txt":       false,
	} {
		if got := Source(r, p); got != want {
			t.Errorf("Source(%s) = %v, want %v", p, got, want)
		}
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
