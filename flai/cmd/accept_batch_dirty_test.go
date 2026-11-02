package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// sharedDirty are the paths of the shared set dirtyMain leaves: what the
// analyzer left on TH-0462, an analysis report and an issue, untracked in
// new folders under design, which every acceptance commits (ADR-0248).
var sharedDirty = []string{"design/analysis/2026-10-19-all.md", "design/issues/I-0150-x.md"}

// dirtyMain leaves in the main checkout at root, uncommitted, sharedDirty
// and the operator's own: an edit of a tracked file and an untracked note
// outside the shared set; it returns the operator's paths as git status
// names them, the rest every acceptance leaves as they are.
func dirtyMain(t *testing.T, root string) []string {
	t.Helper()
	files := map[string]string{"cli/.gitkeep": "the operator's edit\n", "notes.txt": "the operator's note\n"}
	for _, rel := range sharedDirty {
		files[rel] = "# " + filepath.Base(rel) + "\n"
	}
	for rel, content := range files {
		if err := os.MkdirAll(filepath.Join(root, filepath.Dir(rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return []string{"cli/.gitkeep", "notes.txt"}
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

// subjects are the subjects of root's main branch, newest first.
func subjects(t *testing.T, root string) []string {
	t.Helper()
	return strings.Split(gitIn(t, root, "log", "--format=%s", "main"), "\n")
}

// assertSharedCommittedRestLeft fails unless, after an acceptance in root,
// the commit at HEAD holds sharedDirty and none of the operator's files
// dirty, which are still uncommitted and as the operator left them, and no
// commit of the shared set was made before the merge.
func assertSharedCommittedRestLeft(t *testing.T, root string, dirty []string) {
	t.Helper()
	if got := porcelainOutsideWip(t, root); !slices.Equal(got, dirty) {
		t.Errorf("uncommitted after the acceptance %v, want the operator's files %v as they were", got, dirty)
	}
	committed := strings.Split(gitIn(t, root, "show", "--name-only", "--format=", "HEAD"), "\n")
	for _, p := range committed {
		if verify.IsLeft(p, dirty) {
			t.Errorf("the acceptance commit holds the operator's file %s", p)
		}
	}
	for _, p := range sharedDirty {
		if !slices.Contains(committed, p) {
			t.Errorf("the acceptance commit must hold %s of the shared set:\n%s", p, strings.Join(committed, "\n"))
		}
	}
	if b, _ := os.ReadFile(filepath.Join(root, "cli/.gitkeep")); string(b) != "the operator's edit\n" {
		t.Errorf("the operator's edit was changed: %q", b)
	}
	if slices.Contains(subjects(t, root), sharedSubject("S-0001")) {
		t.Error("the shared set the story does not change goes with the acceptance commit, not one of its own")
	}
}

// S-0402 criterion 2, ADR-0248: the merge queue's landing lands while the
// main checkout holds uncommitted files no story of the batch changes: it
// commits those of the shared set with the acceptance and leaves the
// operator's as they were, uncommitted and out of the acceptance commit.
func TestALandingLeavesTheMainCheckoutsUncommittedFilesAsTheyAre(t *testing.T) {
	t.Parallel()
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
	assertSharedCommittedRestLeft(t, root, dirty)
	assertLanded(t, root, "alex")
}

// ADR-0248: an acceptance by hand on a dirty main checkout asks nothing: it
// commits the shared set with the acceptance, leaves the operator's files
// as they are, and its dry run says so, and that --yes includes them.
func TestAnAcceptanceByHandCommitsTheSharedSetAndLeavesTheRest(t *testing.T) {
	t.Parallel()
	root, _ := orchestratedStoryInReview(t, false)
	dirty := dirtyMain(t, root)
	dry, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--dry-run")
	if code != 0 || strings.Contains(dry, "blocked:") {
		t.Fatalf("the dry run of an acceptance by hand on a dirty checkout blocks nothing: %d %s\n%s", code, errOut, dry)
	}
	for _, want := range []string{
		"uncommitted under design and the shared paths: design/analysis/2026-10-19-all.md, design/issues/I-0150-x.md; the acceptance commits them with the acceptance commit, as it commits wip, or before the merge when story/S-0001 would overwrite one (ADR-0248)",
		"uncommitted outside wip, design, and the shared paths: cli/.gitkeep, notes.txt; the acceptance leaves them as they are; it is refused only while the story's change would overwrite one, and --yes includes them",
	} {
		if !strings.Contains(dry, want) {
			t.Errorf("the dry run must say %q:\n%s", want, dry)
		}
	}
	out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex")
	if code != 0 {
		t.Fatalf("flai accept by hand on a dirty checkout: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done {
		t.Errorf("S-0001 must be accepted and archived: %s, archived %v", story.Status, story.Archived)
	}
	assertSharedCommittedRestLeft(t, root, dirty)
}

// ADR-0248: --yes, by hand from the shell, includes the operator's files in
// the acceptance commit too; its dry run says so.
func TestAnAcceptanceByHandWithYesIncludesTheRest(t *testing.T) {
	t.Parallel()
	root, _ := orchestratedStoryInReview(t, false)
	dirty := dirtyMain(t, root)
	dry, _, _ := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--dry-run", "--yes")
	if !strings.Contains(dry, "uncommitted outside wip, design, and the shared paths: cli/.gitkeep, notes.txt; --yes includes them in the acceptance commit") {
		t.Errorf("the dry run under --yes says the operator's files are included:\n%s", dry)
	}
	if out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--yes"); code != 0 {
		t.Fatalf("flai accept --yes on a dirty checkout: %d %s\n%s", code, errOut, out)
	}
	committed := strings.Split(gitIn(t, root, "show", "--name-only", "--format=", "HEAD"), "\n")
	for _, p := range slices.Concat(dirty, sharedDirty) {
		if !slices.Contains(committed, p) {
			t.Errorf("the acceptance commit under --yes must hold %s:\n%s", p, strings.Join(committed, "\n"))
		}
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("tree should be clean after acceptance with --yes: %s", st)
	}
}

// ADR-0248: an acceptance by hand of a story whose change would overwrite
// one of the operator's files is refused before anything changes, with or
// without --yes, and its dry run lists the blocker dirty_main, naming the
// file and the command that sets it aside, among the blockers and not among
// the orchestrator's.
func TestAnAcceptanceByHandIsRefusedAStoryWhoseChangeWouldOverwriteTheOperatorsFile(t *testing.T) {
	t.Parallel()
	root, _ := orchestratedStoryInReview(t, false)
	dirtyMain(t, root)
	if err := os.WriteFile(filepath.Join(root, "cli/feature.go"), []byte("the operator's own\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mainBefore := gitIn(t, root, "rev-parse", "main")
	want := preview.DirtyMainMessage("S-0001", []string{"cli/feature.go"})
	if !strings.Contains(want, "git stash push -- cli/feature.go") || !strings.Contains(want, "(ADR-0248)") {
		t.Fatalf("the blocker names the command that sets the file aside: %q", want)
	}

	out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run reports blockers rather than failing: %d %s", code, errOut)
	}
	var res preview.Acceptance
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if !slices.Equal(res.Blockers, []string{want}) || len(res.OrchestratorBlockers) != 0 {
		t.Errorf("the dry run must list %s alone, %q: %+v, orchestrator's %+v", preview.BlockDirtyMain, want, res.Blockers, res.OrchestratorBlockers)
	}
	// the shared set holds the story's own wip files beside what dirtyMain
	// left under design
	if !slices.Equal(res.Uncommitted, []string{"cli/.gitkeep", "cli/feature.go", "notes.txt"}) || !slices.Equal(outsideFolder(res.SharedUncommitted, "wip/"), sharedDirty) || len(res.SharedUncommitted) == len(sharedDirty) {
		t.Errorf("uncommitted %v and shared_uncommitted %v, want the operator's files and the shared set", res.Uncommitted, res.SharedUncommitted)
	}

	for _, args := range [][]string{
		{"accept", "S-0001", "--by", "alex"},
		{"accept", "S-0001", "--by", "alex", "--yes"},
	} {
		if _, errOut, code := runInAt(t, root, landedAt, args...); code == 0 || !strings.Contains(errOut, "S-0001 cannot be accepted yet: "+want) {
			t.Errorf("flai %v must be refused naming the file: %d\n%s", args, code, errOut)
		}
	}
	assertNothingAccepted(t, root, mainBefore)
	if b, _ := os.ReadFile(filepath.Join(root, "cli/feature.go")); string(b) != "the operator's own\n" {
		t.Errorf("the operator's file was changed: %q", b)
	}
}

// ADR-0248: when the story's change would overwrite a file of the shared
// set, the acceptance commits the shared set on the main branch first, under
// its own subject, so that the rebase settles the file by git's merge in
// place of refusing to overwrite it; the operator's file outside the set is
// left as it is, out of both commits.
func TestAnAcceptanceCommitsTheSharedSetBeforeTheMergeWhenTheStoryWouldOverwriteIt(t *testing.T) {
	t.Parallel()
	root, _ := orchestratedStoryInReview(t, false)
	lines := make([]string, 12)
	for i := range lines {
		lines[i] = "line " + string(rune('a'+i))
	}
	logFile := filepath.Join(root, "design", "log.md")
	if err := os.WriteFile(logFile, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "-A", "--", "design/log.md")
	gitIn(t, root, "commit", "-q", "-m", "chore: a log under design")
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	gitIn(t, wt, "rebase", "-q", "main")
	branchLines := slices.Clone(lines)
	branchLines[11] = "the branch's last line"
	if err := os.WriteFile(filepath.Join(wt, "design", "log.md"), []byte(strings.Join(branchLines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] write the log's last line")
	mainLines := slices.Clone(lines)
	mainLines[0] = "flai's first line"
	for rel, content := range map[string]string{
		"design/log.md":    strings.Join(mainLines, "\n") + "\n",
		"design/digest.md": "# a digest entry flai wrote\n",
		"cli/.gitkeep":     "the operator's edit\n",
		"notes.txt":        "the operator's note\n",
	} {
		if err := os.WriteFile(filepath.Join(root, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	dirty := []string{"cli/.gitkeep", "notes.txt"}
	mainBefore := gitIn(t, root, "rev-parse", "main")

	out, errOut, code := runInAt(t, root, landedAt, "accept", "S-0001", "--by", "alex")
	if code != 0 {
		t.Fatalf("flai accept with a shared file the story changes dirty: %d %s\n%s", code, errOut, out)
	}
	if story := acceptedS0001(t, root); !story.Archived || story.Status != workitem.Done {
		t.Errorf("S-0001 must be accepted and archived: %s, archived %v", story.Status, story.Archived)
	}
	subject := sharedSubject("S-0001")
	if subject != "chore: [S-0001] commit the main checkout's wip and design before accepting" {
		t.Fatalf("the pre-commit's subject: %q", subject)
	}
	log := subjects(t, root)
	i := slices.Index(log, subject)
	if i < 0 || log[0] != "chore: [S-0001] accept and archive, with E-0001" || !slices.Contains(log[:i], "feat: [S-0001] write the log's last line") {
		t.Fatalf("the shared set is committed before the branch's commits and the acceptance commit:\n%s", strings.Join(log, "\n"))
	}
	pre := gitIn(t, root, "rev-parse", "main~"+strconv.Itoa(i))
	if parent := gitIn(t, root, "rev-parse", pre+"^"); parent != mainBefore {
		t.Errorf("the pre-commit %s follows main as it was, %s, not %s", pre, mainBefore, parent)
	}
	committed := strings.Split(gitIn(t, root, "show", "--name-only", "--format=", pre), "\n")
	for _, p := range []string{"design/log.md", "design/digest.md"} {
		if !slices.Contains(committed, p) {
			t.Errorf("the pre-commit must hold %s of the shared set:\n%s", p, strings.Join(committed, "\n"))
		}
	}
	for _, p := range committed {
		if verify.IsLeft(p, dirty) {
			t.Errorf("the pre-commit holds the operator's file %s", p)
		}
	}
	if got := gitIn(t, root, "show", "main:design/log.md"); !strings.HasPrefix(got, "flai's first line\n") || !strings.HasSuffix(got, "\nthe branch's last line") {
		t.Errorf("main's log must hold flai's first line and the branch's last:\n%s", got)
	}
	if got := porcelainOutsideWip(t, root); !slices.Equal(got, dirty) {
		t.Errorf("uncommitted after the acceptance %v, want the operator's files %v as they were", got, dirty)
	}
}

func TestAddArgsLeavesOutEachPathLeft(t *testing.T) {
	t.Parallel()
	if got := addArgs(nil); !slices.Equal(got, []string{"add", "-A"}) {
		t.Errorf("addArgs(nil) = %v", got)
	}
	want := []string{"add", "-A", "--", ".", ":(exclude,literal)design/a.md", ":(exclude,literal)design/analysis"}
	if got := addArgs([]string{"design/a.md", "design/analysis/"}); !slices.Equal(got, want) {
		t.Errorf("addArgs = %v, want %v", got, want)
	}
}

// ADR-0248: leftUncommitted is the rest for every acceptance, none when
// --yes includes it.
func TestLeftUncommittedIsTheRestUnlessIncluded(t *testing.T) {
	t.Parallel()
	res := &preview.Acceptance{SharedUncommitted: []string{"design/a.md"}, Uncommitted: []string{"notes.txt"}}
	for _, tc := range []struct {
		name                string
		orchestrator, batch bool
		includes            bool
		want                []string
	}{
		{name: "by hand", want: []string{"notes.txt"}},
		{name: "the orchestrator's", orchestrator: true, want: []string{"notes.txt"}},
		{name: "the queue's landing", batch: true, want: []string{"notes.txt"}},
		{name: "by hand with --yes", includes: true},
	} {
		r := *res
		r.Orchestrator, r.IncludesUncommitted = tc.orchestrator, tc.includes
		if tc.batch {
			r.Batch = landingBatch
		}
		if got := leftUncommitted(&r); !slices.Equal(got, tc.want) {
			t.Errorf("%s: leftUncommitted = %v, want %v", tc.name, got, tc.want)
		}
	}
	if got := leftUncommitted(nil); got != nil {
		t.Errorf("leftUncommitted(nil) = %v", got)
	}
}
