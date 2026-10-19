package storygit

import (
	"reflect"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// S-0402: of the main checkout's uncommitted files, only those the change
// also changes are named: a dirty file it changes, an untracked file it adds,
// and not a dirty or untracked file it leaves alone, as the analyzer's report
// and issues on TH-0462.
func TestOverwrittenNamesOnlyTheDirtyFilesTheChangeChanges(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	commitFiles(t, dir, "base", "cli/a.go", "cli/b.go")
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	write(t, dir, "cli/a.go", "changed\n")
	commitFiles(t, dir, "feat: [S-0001] change a, add c", "cli/c.go")
	git(t, dir, "checkout", "-q", "main")

	write(t, dir, "cli/a.go", "the operator's edit\n")
	write(t, dir, "cli/b.go", "left alone by the change\n")
	writeFile(t, dir, "cli/c.go")
	writeFile(t, dir, "design/analysis/2026-10-19-all.md")
	writeFile(t, dir, "design/issues/I-0150-x.md")

	got, err := Overwritten(r, dir, "main", "story/S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"cli/a.go", "cli/c.go"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Overwritten = %v, want %v", got, want)
	}
}

// A clean checkout, or one whose dirty files the change leaves alone, has
// nothing overwritten; and a branch made before main moved on is read from
// the commit the two share, so main's own later changes are not the
// branch's.
func TestOverwrittenIsEmptyWhenTheChangeLeavesTheDirtyFilesAlone(t *testing.T) {
	r, dir := execx.Runner(execx.System{}), gitRepo(t)
	git(t, dir, "checkout", "-q", "-b", "story/S-0001")
	commitFiles(t, dir, "feat: [S-0001] add x", "cli/x.go")
	git(t, dir, "checkout", "-q", "main")
	commitFiles(t, dir, "main moves on", "cli/m.go")

	if got, err := Overwritten(r, dir, "main", "story/S-0001"); err != nil || len(got) != 0 {
		t.Fatalf("clean checkout: Overwritten = %v, %v; want none", got, err)
	}
	write(t, dir, "cli/m.go", "the operator's edit of main's own file\n")
	writeFile(t, dir, "notes.md")
	if got, err := Overwritten(r, dir, "main", "story/S-0001"); err != nil || len(got) != 0 {
		t.Fatalf("dirty elsewhere: Overwritten = %v, %v; want none", got, err)
	}
}
