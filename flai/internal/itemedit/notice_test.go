package itemedit

import (
	"fmt"
	"os"
	"path/filepath"
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
