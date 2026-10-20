package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// dirtyMain leaves in the main checkout at root, uncommitted, what the
// operator and the analyzer left on TH-0462: an analysis report and an issue,
// untracked in new folders, and an edit of a tracked file; it returns the
// paths as git status names them, each folder once.
func dirtyMain(t *testing.T, root string) []string {
	t.Helper()
	for rel, content := range map[string]string{
		"design/analysis/2026-10-19-all.md": "# report\n",
		"design/issues/I-0150-x.md":         "# issue\n",
		"design/.gitkeep":                   "the operator's edit\n",
	} {
		if err := os.MkdirAll(filepath.Join(root, filepath.Dir(rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return []string{"design/.gitkeep", "design/analysis/", "design/issues/"}
}

// porcelainOutsideWip is git status's paths in root outside wip, sorted.
func porcelainOutsideWip(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	for _, p := range workitem.PorcelainPaths(gitIn(t, root, "status", "--porcelain")) {
		if !strings.HasPrefix(p, "wip/") {
			out = append(out, p)
		}
	}
	slices.Sort(out)
	return out
}

// S-0402 criterion 2: the merge queue's landing lands while the main
// checkout holds uncommitted files no story of the batch changes, and leaves
// them as they were, uncommitted and out of the acceptance commit.
func TestALandingLeavesTheMainCheckoutsUncommittedFilesAsTheyAre(t *testing.T) {
	root, head := orchestratedStoryInReview(t, false)
	landingState(t, root, head, "alex", nil)
	dirty := dirtyMain(t, root)

	preview, _, _ := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--batch", landingBatch, "--dry-run")
	if !strings.Contains(preview, "the merge queue's landing leaves them as they are") || strings.Contains(preview, "refuses") {
		t.Errorf("the dry run of a landing says the files are left as they are:\n%s", preview)
	}
	out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--batch", landingBatch)
	if code != 0 {
		t.Fatalf("flai accept --batch on a dirty checkout: %d %s\n%s", code, errOut, out)
	}
	if !acceptedS0001(t, root).Archived {
		t.Error("S-0001 must be accepted and archived")
	}
	if got := gitIn(t, root, "show", "main:cli/feature.go"); got == "" {
		t.Error("the story's change must be on main")
	}
	if got := porcelainOutsideWip(t, root); !slices.Equal(got, dirty) {
		t.Errorf("uncommitted after the landing %v, want the operator's files %v as they were", got, dirty)
	}
	for _, p := range strings.Split(gitIn(t, root, "show", "--name-only", "--format=", "HEAD"), "\n") {
		if isLeft(p, dirty) {
			t.Errorf("the acceptance commit holds the operator's file %s", p)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(root, "design/.gitkeep")); string(b) != "the operator's edit\n" {
		t.Errorf("the operator's edit was changed: %q", b)
	}
	assertLanded(t, root, "alex")
}

// S-0402 criterion 5: an acceptance by hand still refuses the main
// checkout's uncommitted files outside wip, changing nothing, unless --yes.
func TestAnAcceptanceByHandStillRefusesTheUncommittedFiles(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, false)
	dirtyMain(t, root)
	_, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex")
	if code == 0 || !strings.Contains(errOut, "uncommitted changes outside wip") || !strings.Contains(errOut, "--yes") {
		t.Fatalf("flai accept by hand on a dirty checkout must refuse, naming --yes: %d %s", code, errOut)
	}
	if acceptedS0001(t, root).Status != workitem.Review {
		t.Error("a refused acceptance changes nothing: S-0001 stays in review")
	}
}

// A landing whose story's change would overwrite one of the main checkout's
// uncommitted files is stopped by git's own check before anything changes:
// the story stays in review and the operator's file is as they left it.
func TestALandingThatWouldOverwriteAnUncommittedFileIsStopped(t *testing.T) {
	root, head := orchestratedStoryInReview(t, false)
	landingState(t, root, head, "alex", nil)
	if err := os.WriteFile(filepath.Join(root, "cli/feature.go"), []byte("the operator's own\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--batch", landingBatch); code == 0 {
		t.Fatal("a landing that would overwrite an uncommitted file must fail")
	}
	if acceptedS0001(t, root).Status != workitem.Review {
		t.Error("S-0001 stays in review")
	}
	if b, _ := os.ReadFile(filepath.Join(root, "cli/feature.go")); string(b) != "the operator's own\n" {
		t.Errorf("the operator's file was changed: %q", b)
	}
}

func TestAddArgsLeavesOutEachPathLeft(t *testing.T) {
	if got := addArgs(nil); !slices.Equal(got, []string{"add", "-A"}) {
		t.Errorf("addArgs(nil) = %v", got)
	}
	want := []string{"add", "-A", "--", ".", ":(exclude,literal)design/a.md", ":(exclude,literal)design/analysis"}
	if got := addArgs([]string{"design/a.md", "design/analysis/"}); !slices.Equal(got, want) {
		t.Errorf("addArgs = %v, want %v", got, want)
	}
	if got := withoutLeft([]string{"design/a.md", "design/analysis/r.md", "design/analysisx.md", "cli/b.go"}, []string{"design/a.md", "design/analysis/"}); !slices.Equal(got, []string{"design/analysisx.md", "cli/b.go"}) {
		t.Errorf("withoutLeft = %v", got)
	}
}
