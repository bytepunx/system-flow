package mcpserver

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// S-0275: adr_new with commit writes the ADR, its index row, and the
// superseded ADR's superseded_by in the story's worktree, commits those and
// only those on the story's branch as <ADR-nnnn> <decision> with the
// story's prefix and the trailers, and widens the story's touches.
func TestAdrNewCommitsOnTheStorysBranch(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	wt := storyGit(t, f)
	writeIn(t, wt.Root, "design/system/plan.md", "---\ntitle: Plan\n---\n\n# Plan\n\nhalf done\n")
	id := f.story.ID

	out, failed := f.call(t, "adr_new", map[string]any{"decision": "Tokens are one per project", "status": "accepted", "supersedes": []string{"ADR-0001"},
		"body": adrBody, "story": id, "commit": true, "trailers": []string{"Co-Authored-By: t <t@t>"}})
	if failed != "" {
		t.Fatal(failed)
	}
	rel, _ := filepath.Rel(f.repo.Root, wt.Root)
	under := filepath.ToSlash(rel) + "/"
	adr := "design/adrs/0002-tokens-are-one-per-project.md"
	if out["id"] != "ADR-0002" || out["title"] != "Tokens are one per project" || out["status"] != "accepted" || out["story"] != id || out["path"] != under+adr {
		t.Errorf("adr_new: %v", out)
	}
	if got, want := strs(out["changed"]), []string{under + adr, under + "design/adrs/README.md", under + "design/adrs/0001-first.md"}; !reflect.DeepEqual(got, want) {
		t.Errorf("changed %v, want %v", got, want)
	}
	committed := []string{"design/adrs/0001-first.md", adr, "design/adrs/README.md"}
	subject := "docs: [" + id + "] ADR-0002 Tokens are one per project"
	if field(out, "commit.subject") != subject || !reflect.DeepEqual(strs(field(out, "commit.paths")), committed) || !reflect.DeepEqual(strs(out["touches_added"]), committed) {
		t.Errorf("commit %v, touches_added %v", out["commit"], out["touches_added"])
	}
	message, files, ahead := lastCommit(t, wt.Root)
	if message != subject+"\n\nCo-Authored-By: t <t@t>" || !reflect.DeepEqual(files, committed) || ahead != "1" {
		t.Errorf("the story's branch should hold one commit of the ADR's files: %q %v %s", message, files, ahead)
	}
	if st := gitIn(t, wt.Root, "status", "--porcelain"); st != "M design/system/plan.md" {
		t.Errorf("the worktree's other change should stay uncommitted, and nothing else: %q", st)
	}
	if data, _ := os.ReadFile(filepath.Join(wt.Root, "design/adrs/0001-first.md")); !strings.Contains(string(data), "superseded_by: [ADR-0002]") {
		t.Errorf("ADR-0001 should be superseded by ADR-0002:\n%s", data)
	}
	story, err := f.repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if want := append([]string{"flai/internal/mcpserver"}, committed...); !reflect.DeepEqual(story.Touches, want) {
		t.Errorf("the story's touches in the project %v, want %v", story.Touches, want)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, adr)); err == nil {
		t.Error("the ADR should not be written in the main checkout")
	}
}

// S-0275: without commit, adr_new writes in the story's worktree when the
// story has one and in the project otherwise, commits nothing, and refuses
// what flai adr new refuses, writing nothing.
func TestAdrNewWritesInTheStorysWorktreeOrTheProject(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	wt := f.repo.WorktreePath(f.story.ID)
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	out, failed := f.call(t, "adr_new", map[string]any{"decision": "Issues stay on the branch", "body": adrBody, "story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	rel, _ := filepath.Rel(f.repo.Root, wt)
	path := filepath.ToSlash(rel) + "/design/adrs/0001-issues-stay-on-the-branch.md"
	if out["id"] != "ADR-0001" || out["status"] != "proposed" || out["story"] != f.story.ID || out["path"] != path || !reflect.DeepEqual(strs(out["changed"]), []string{path}) {
		t.Errorf("adr_new in the worktree: %v", out)
	}
	if out["commit"] != nil || out["touches_added"] != nil {
		t.Errorf("without commit nothing is committed: %v", out)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, path)); err != nil {
		t.Errorf("the ADR should be in the worktree: %v", err)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, "design/adrs")); err == nil {
		t.Error("the main checkout's ADR folder should be left alone while the story has a worktree")
	}

	out, failed = f.call(t, "adr_new", map[string]any{"decision": "Decisions without a story go to the project", "body": adrBody})
	if failed != "" {
		t.Fatal(failed)
	}
	p, _ := out["path"].(string)
	if out["story"] != nil || !strings.HasPrefix(p, "design/adrs/") {
		t.Errorf("adr_new without a story writes in the project: %v", out)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, p)); err != nil {
		t.Errorf("the ADR should be in the project: %v", err)
	}

	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"decision": "Bad number", "supersedes": []string{"seven"}}, `adr_new supersedes "seven" is not an ADR number`},
		{map[string]any{"decision": "Bad status", "status": "rejected"}, `status "rejected" must be proposed or accepted`},
		{map[string]any{"decision": "No such ADR", "refines": []string{"ADR-0042"}}, "there is no ADR-0042 to supersede or refine"},
		{map[string]any{"decision": " "}, "an ADR needs a title"},
	} {
		c.args["story"] = f.story.ID
		if _, failed := f.call(t, "adr_new", c.args); !strings.Contains(failed, c.want) {
			t.Errorf("adr_new %v should be refused with %q: %q", c.args, c.want, failed)
		}
	}
	if entries, _ := os.ReadDir(filepath.Join(wt, "design/adrs")); len(entries) != 1 {
		t.Errorf("a refusal should write nothing: %v", entries)
	}
}
