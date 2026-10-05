package itemedit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The notices are a log beside the cursors, and stay bounded.
func TestNoticesStayBounded(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if Notices(repo) != nil {
		t.Fatal("none to begin with")
	}
	for i := 0; i < 2*keep+5; i++ {
		Record(repo, Notice{At: "2026-09-20T15:00:00Z", By: "alex", ID: fmt.Sprintf("S-%04d", i), Changed: []string{"title"}})
	}
	all := Notices(repo)
	if len(all) > 2*keep || len(all) < keep || all[len(all)-1].ID != fmt.Sprintf("S-%04d", 2*keep+4) {
		t.Errorf("kept %d, the newest last: %+v", len(all), all[len(all)-1])
	}
	if st, err := os.Stat(NoticesPath(repo)); err != nil || st.Mode().Perm() != 0o600 {
		t.Errorf("the file is the owner's: %v %v", st, err)
	}
}

// Overlap notices are a log of their own: an older flai reads every line of
// the edit notices as an edit (S-0132).
func TestOverlapsAreALogOfTheirOwn(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	RecordOverlap(repo, Overlap{At: "2026-09-26T18:00:00Z", By: "alex", ID: "S-0002", Title: "Open", Accepted: "S-0001", Paths: []string{"flai/cmd/a.go"}})
	got := Overlaps(repo)
	if len(got) != 1 || got[0].ID != "S-0002" || got[0].Accepted != "S-0001" || len(got[0].Paths) != 1 {
		t.Errorf("overlaps: %+v", got)
	}
	if Notices(repo) != nil {
		t.Errorf("an overlap is not an edit notice: %+v", Notices(repo))
	}
	if OverlapsPath(repo) == NoticesPath(repo) {
		t.Error("two logs, two files")
	}
}

// A notice of claims that grew to overlap names both stories and no
// acceptance (I-0059); a notice written at acceptance before it reads as it
// did, and one that names only one of the two stories is not read.
func TestOverlapsReadOldAndGrownNotices(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.MkdirAll(filepath.Dir(OverlapsPath(repo)), 0o700)
	old := `{"at":"2026-09-26T18:00:00Z","by":"alex","id":"S-0002","title":"Open","accepted":"S-0001","paths":["flai/cmd/a.go"]}` + "\n" +
		`{"at":"2026-09-26T18:00:01Z","by":"claude","id":"S-0003","title":"Half","grew":"S-0003","paths":["docs"]}` + "\n"
	if err := os.WriteFile(OverlapsPath(repo), []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	RecordOverlap(repo, Overlap{At: "2026-09-26T18:01:00Z", By: "claude", ID: "S-0004", Title: "Grew", Grew: "S-0004", Reached: "S-0005", Paths: []string{"docs/a.md"}})
	RecordOverlap(repo, Overlap{At: "2026-09-26T18:01:00Z", By: "claude", ID: "S-0005", Title: "Reached", Grew: "S-0004", Reached: "S-0005", Paths: []string{"docs/a.md"}})
	got := Overlaps(repo)
	if len(got) != 3 {
		t.Fatalf("read: %+v", got)
	}
	if o := got[0]; o.ID != "S-0002" || o.Accepted != "S-0001" || o.Cause() != "S-0001" || o.Grew != "" {
		t.Errorf("the old notice: %+v", o)
	}
	if got[1].Cause() != "S-0005" || got[2].Cause() != "S-0004" || got[1].Accepted != "" {
		t.Errorf("each story's cause is the other: %+v", got[1:])
	}
	data, _ := os.ReadFile(OverlapsPath(repo))
	if lines := strings.Split(strings.TrimSpace(string(data)), "\n"); strings.Contains(lines[len(lines)-1], `"accepted"`) {
		t.Errorf("a grown notice names no acceptance, so an older flai passes over it: %s", lines[len(lines)-1])
	}
}
