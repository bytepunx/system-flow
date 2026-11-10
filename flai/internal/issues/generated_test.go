package issues

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0175, ADR-0276: a rebase merges by sections each markdown file inside
// the manifest's claims.shared or under its design or docs folder that the
// generated files, the summary and the ADR index, and the issue files do not
// already cover; any other path, a markdown file under wip, in a sub-project,
// or at the root among them, and a shared path that is not markdown, is not
// merged by sections.
func TestFilesInMergesTheSharedMarkdownFilesBySections(t *testing.T) {
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs/\n  wip: wip\n" +
		"claims:\n  shared:\n    - design/system/flai-cli.md\n    - design/adrs\n    - design/issues\n    - docs/users/**\n    - template/CHANGELOG.md\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	files := FilesIn(r, root, execx.System{}, func() time.Time { return t0 })
	if files.Shared.Match == nil || files.Shared.Merge == nil {
		t.Fatal("no shared files are merged")
	}
	for p, want := range map[string]bool{
		"design/system/flai-cli.md":           true,
		"design/adrs/0001-a-decision.md":      true,
		"docs/users/flai.md":                  true,
		"template/CHANGELOG.md":               true,
		"design/adrs/README.md":               false, // the ADR index, generated
		"design/issues/I-0001-a-friction.md":  false, // an issue file, merged by its instances
		"design/issues/summary.md":            false, // generated
		"docs/users/notes.txt":                false, // shared, not markdown
		"design/system/server-performance.md": true,  // not shared, under the design folder (I-0315)
		"design/digest/D-0001-a-fact.md":      true,  // under the design folder
		"docs/operators/index.md":             true,  // not shared, under the docs folder
		"design/system/tokens.json":           false, // under the design folder, not markdown
		"wip/agents/S-0001.md":                false, // wip is flai's
		"wip/kanban/stories/S-0001-a.md":      false,
		"flai/README.md":                      false, // a sub-project's
		"designs/notes.md":                    false, // a folder that only starts with design
		"README.md":                           false,
	} {
		if got := files.Shared.Match(p); got != want {
			t.Errorf("%s merged by sections: %v, want %v", p, got, want)
		}
	}
}

// ADR-0228 §1: the generated files are the summary and the ADR index, each
// written again when a file it is written from is merged at a stop: the
// summary from the issue files, the index from the ADR files beside it.
func TestGeneratedAreTheSummaryAndTheADRIndex(t *testing.T) {
	r := repo(t)
	gen := Generated(r, "", func() time.Time { return t0 })
	var paths []string
	for _, g := range gen {
		paths = append(paths, g.Path)
	}
	if !slices.Equal(paths, []string{"design/issues/summary.md", "design/adrs/README.md"}) {
		t.Fatalf("generated %q", paths)
	}
	if !gen[1].From("design/adrs/0228-generated.md") || gen[1].From("design/adrs/0000-template.md") || gen[1].From("design/issues/I-0001-a.md") || gen[0].From("design/adrs/0228-generated.md") {
		t.Error("each generated file is written from its own folder's files")
	}
	files := FilesIn(r, r.Root, execx.System{}, func() time.Time { return t0 })
	if files.Merged.Match("design/adrs/README.md") {
		t.Error("the ADR index is generated, not merged row by row")
	}
	if got := files.Unresolved([]string{"design/adrs/README.md", "design/issues/summary.md", "docs/guide.md"}); !slices.Equal(got, []string{"docs/guide.md"}) {
		t.Errorf("the trial merge leaves the generated files out of a pair's conflicts: %q", got)
	}
}

// adrDoc is an ADR file numbered n with title and status.
func adrDoc(n, title, status string) string {
	return "---\nid: ADR-" + n + "\ntitle: " + title + "\nstatus: " + status + "\ndate: 2026-10-31\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-" + n + " " + title + "\n\n## Decision\n\nWe decide.\n"
}

// git runs git in dir and returns its output, failing t when it fails.
func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// commitADR writes file, an ADR with text, and the index in the checkout at
// root and commits them.
func commitADR(t *testing.T, root, file, text, index string) {
	t.Helper()
	for name, body := range map[string]string{file: text, "README.md": index} {
		if err := os.WriteFile(filepath.Join(root, "design", "adrs", name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "add "+file)
}

// S-0550, ADR-0228 §1: a rebase that stops on the ADR index because both
// sides added an ADR writes it again from the rebased tree's ADR files, with
// both sides' rows as the files say, a row the branch wrote by hand
// corrected, and the text around the table kept, and continues with no
// person's help. Integration: it runs real git.
func TestARebaseStoppedOnTheADRIndexWritesItFromBothSidesADRFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nclaims:\n  shared:\n    - design/adrs\n    - design/issues\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "design", "adrs"), 0o755); err != nil {
		t.Fatal(err)
	}
	intro := "# Architecture Decision Records\n\nGenerated.\n\n| ADR | Title | Status |\n|-----|-------|--------|\n"
	row1 := "| [0001](0001-one.md) | One | accepted |\n"
	git(t, root, "init", "-q", "-b", "main")
	commitADR(t, root, "0001-one.md", adrDoc("0001", "One", "accepted"), intro+row1+"\nBelow.\n")
	git(t, root, "checkout", "-q", "-b", "story")
	commitADR(t, root, "0003-three.md", adrDoc("0003", "Three", "proposed"), intro+row1+"| [0003](0003-three.md) | Three | accepted, by hand |\n\nBelow, longer.\n")
	git(t, root, "checkout", "-q", "main")
	commitADR(t, root, "0002-two.md", adrDoc("0002", "Two", "accepted"), intro+row1+"| [0002](0002-two.md) | Two | accepted |\n\nBelow.\n")
	git(t, root, "checkout", "-q", "story")
	if out, err := exec.Command("git", "-C", root, "rebase", "-q", "main").CombinedOutput(); err == nil {
		t.Fatalf("the rebase did not stop on the index:\n%s", out)
	}
	if got := storygit.Conflicts(execx.System{}, root); !slices.Equal(got, []string{"design/adrs/README.md"}) {
		t.Fatalf("the rebase stopped on %q", got)
	}

	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	s, err := storygit.ContinueOverResolved(execx.System{}, root, FilesIn(r, root, execx.System{}, func() time.Time { return t0 }), nil, "S-0001")
	if err != nil || s.Waits {
		t.Fatalf("the stop was left for a person: %v %+v", err, s)
	}
	if !slices.Equal(s.Regenerated, []string{"design/adrs/README.md"}) || len(s.Merged) != 0 {
		t.Errorf("regenerated %q, merged %q; want the index written again, nothing merged", s.Regenerated, s.Merged)
	}
	want := intro + row1 + "| [0002](0002-two.md) | Two | accepted |\n| [0003](0003-three.md) | Three | proposed |\n\nBelow, longer.\n"
	if got := git(t, root, "show", "HEAD:design/adrs/README.md"); got != want {
		t.Errorf("the index on the branch:\n%s\nwant:\n%s", got, want)
	}
	if storygit.RebaseInProgress(execx.System{}, root) || strings.TrimSpace(git(t, root, "status", "--porcelain")) != "" {
		t.Errorf("the rebase did not finish clean:\n%s", git(t, root, "status"))
	}
}
