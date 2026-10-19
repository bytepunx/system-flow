package issues

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0175: a rebase merges by sections each markdown file inside the
// manifest's claims.shared that the summary, the issue files, and the ADR
// index do not already cover; any other path, and a shared path that is not
// markdown, is not merged by sections.
func TestFilesInMergesTheSharedMarkdownFilesBySections(t *testing.T) {
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
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
		"design/adrs/README.md":               false, // the ADR index, merged row by row
		"design/issues/I-0001-a-friction.md":  false, // an issue file, merged by its instances
		"design/issues/summary.md":            false, // generated
		"docs/users/notes.txt":                false, // shared, not markdown
		"design/system/flaiover-dashboard.md": false, // not shared
		"README.md":                           false,
	} {
		if got := files.Shared.Match(p); got != want {
			t.Errorf("%s merged by sections: %v, want %v", p, got, want)
		}
	}
}
