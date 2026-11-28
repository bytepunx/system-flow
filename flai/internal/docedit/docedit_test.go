package docedit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var now = time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)

const (
	planPath = "design/system/plan.md"
	mainPlan = "---\ntitle: Plan\nupdated: 2026-09-01\nstatus: active\n---\n\n# Plan\n\nmain text\n"
	branchPl = "---\ntitle: Plan\nupdated: 2026-09-01\nstatus: active\n---\n\n# Plan\n\nbranch text\n"
)

// fakeGit answers the git calls a save makes: a work tree on branch, or no
// repository when branch is empty, and records every call with its folder.
type fakeGit struct {
	branch string
	calls  []string
}

func (f *fakeGit) Run(dir, name string, args ...string) (string, error) {
	call := dir + ": " + name + " " + strings.Join(args, " ")
	f.calls = append(f.calls, call)
	switch {
	case f.branch == "" && len(args) > 0 && args[0] == "rev-parse":
		return "", errors.New("fatal: not a git repository")
	case len(args) > 1 && args[0] == "rev-parse" && args[1] == "--short":
		return "abc1234", nil
	case len(args) > 0 && args[0] == "symbolic-ref":
		if f.branch == "(detached)" {
			return "", errors.New("exit status 1")
		}
		return f.branch, nil
	}
	return "", nil
}

func (f *fakeGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return f.Run(dir, name, args...)
}

func (f *fakeGit) LookPath(name string) (string, error) { return name, nil }

// project is a main checkout with E-0001, S-0001, and S-0002, and S-0001's
// worktree under .flai-cache, linked to it as git links one, holding its
// own copy of the plan.
func project(t *testing.T) (repo *workitem.Repo, root, wt string) {
	t.Helper()
	root = t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	write := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(root, "system-flow.yaml"), manifest)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "docs"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(root, planPath), mainPlan)
	write(filepath.Join(root, "design/system/main-only.md"), mainPlan)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range []workitem.NewOptions{
		{Type: workitem.Epic, Title: "Epic", Now: now},
		{Type: workitem.Story, Title: "Branch", Parent: "E-0001", Now: now},
		{Type: workitem.Story, Title: "No worktree", Parent: "E-0001", Now: now},
	} {
		if _, err := repo.Create(o); err != nil {
			t.Fatal(err)
		}
	}
	wt = repo.WorktreePath("S-0001")
	write(filepath.Join(wt, ".git"), "gitdir: "+filepath.Join(root, ".git", "worktrees", "S-0001")+"\n")
	write(filepath.Join(wt, "system-flow.yaml"), manifest)
	write(filepath.Join(wt, planPath), branchPl)
	if err := os.MkdirAll(filepath.Join(wt, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	return repo, root, wt
}

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0378: with a story, Show reads the file as it is on the story's branch.
func TestShowWithAStoryReadsItsWorktree(t *testing.T) {
	repo, _, _ := project(t)
	doc, err := Show(repo, planPath, "S-1")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Content != branchPl || doc.Hash != Hash(branchPl) || doc.Mode != Full || doc.Path != planPath {
		t.Errorf("the branch's plan: %+v", doc)
	}
	if doc, err := Show(repo, planPath, ""); err != nil || doc.Content != mainPlan {
		t.Errorf("without a story, the main checkout's: %+v %v", doc, err)
	}
}

// S-0378: with a story, Save writes the worktree's file, checks it there,
// and commits that one path on the story's branch in the worktree.
func TestSaveWithAStoryCommitsOnItsBranch(t *testing.T) {
	repo, root, wt := project(t)
	git := &fakeGit{branch: "story/S-0001"}
	edited := strings.Replace(branchPl, "branch text", "better branch text", 1)
	res, err := Save(repo, git, planPath, edited, SaveOptions{Hash: Hash(branchPl), Story: "S-1", Message: "tidy", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Committed || res.Commit != "abc1234" || res.Path != planPath {
		t.Errorf("result: %+v", res)
	}
	if got := read(t, filepath.Join(wt, planPath)); !strings.Contains(got, "better branch text") || !strings.Contains(got, "updated: 2026-10-13") {
		t.Errorf("the worktree's plan carries the edit and today's date:\n%s", got)
	}
	if got := read(t, filepath.Join(root, planPath)); got != mainPlan {
		t.Errorf("the main checkout's plan is left as it was:\n%s", got)
	}
	want := []string{
		wt + ": git add -- " + planPath,
		wt + ": git commit --quiet -m docs: [S-0001] tidy -- " + planPath,
	}
	joined := strings.Join(git.calls, "\n")
	for _, w := range want {
		if !strings.Contains(joined, w) {
			t.Errorf("git should have run %q; ran:\n%s", w, joined)
		}
	}
	if strings.Contains(joined, root+": git add") || strings.Contains(joined, root+": git commit") {
		t.Errorf("nothing is committed in the main checkout:\n%s", joined)
	}
}

// S-0378: a save on the branch whose file changed after it was loaded is a
// conflict carrying the branch's content, and writes nothing.
func TestSaveWithAStoryRefusesAStaleHash(t *testing.T) {
	repo, _, wt := project(t)
	theirs := strings.Replace(branchPl, "branch text", "their text", 1)
	if err := os.WriteFile(filepath.Join(wt, planPath), []byte(theirs), 0o644); err != nil {
		t.Fatal(err)
	}
	git := &fakeGit{branch: "story/S-0001"}
	_, err := Save(repo, git, planPath, strings.Replace(branchPl, "branch text", "my text", 1), SaveOptions{Hash: Hash(branchPl), Story: "S-0001", Now: now})
	c, ok := IsConflict(err)
	if !ok {
		t.Fatalf("a stale hash must be a conflict: %v", err)
	}
	if c.Current != theirs || c.Hash != Hash(theirs) {
		t.Errorf("the conflict carries the branch's file: %+v", c)
	}
	if got := read(t, filepath.Join(wt, planPath)); got != theirs {
		t.Errorf("a conflict must not overwrite the other change:\n%s", got)
	}
	for _, call := range git.calls {
		if strings.Contains(call, "git commit") {
			t.Errorf("nothing is committed: %s", call)
		}
	}
}

// S-0378: what a story's worktree cannot edit is refused with the reason,
// and nothing is written.
func TestSaveWithAStoryRefuses(t *testing.T) {
	repo, root, wt := project(t)
	wipPath := "wip/kanban/stories/S-0001-branch.md"
	if err := os.MkdirAll(filepath.Join(wt, "wip/kanban/stories"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, wipPath), []byte(read(t, filepath.Join(root, wipPath))), 0o644); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, path, story, branch, want string
	}{
		{"a wip path", wipPath, "S-0001", "story/S-0001", "lives on the main checkout, not on story/S-0001 (ADR-0019)"},
		{"a story without a worktree", planPath, "S-0002", "story/S-0002", "S-0002: it has no worktree at " + repo.WorktreePath("S-0002") + "; open it with flai stream open S-0002"},
		{"an epic", planPath, "E-0001", "", "it is a epic; name a story (S-nnnn)"},
		{"a story that is not there", planPath, "S-0099", "", "name a story of this project"},
		{"a file only on main", "design/system/main-only.md", "S-0001", "story/S-0001", "design/system/main-only.md does not exist in S-0001's worktree"},
		{"a worktree on another branch", planPath, "S-0001", "main", "its worktree is on main, not story/S-0001"},
		{"a worktree with a detached HEAD", planPath, "S-0001", "(detached)", "its worktree is on a detached HEAD, not story/S-0001"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			git := &fakeGit{branch: c.branch}
			content := "---\ntitle: x\n---\n\n# x\n"
			_, err := Save(repo, git, c.path, content, SaveOptions{Hash: Hash(branchPl), Story: c.story, Now: now})
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("want %q, got %v", c.want, err)
			}
			if _, ok := IsConflict(err); ok {
				t.Errorf("a refusal is not a conflict: %v", err)
			}
			if got := read(t, filepath.Join(wt, planPath)); got != branchPl {
				t.Errorf("nothing is written:\n%s", got)
			}
			if c.branch != "main" && c.branch != "(detached)" {
				if _, err := Show(repo, c.path, c.story); err == nil || !strings.Contains(err.Error(), c.want) {
					t.Errorf("show: want %q, got %v", c.want, err)
				}
			}
		})
	}
}

// ADR-0223: a rejected ADR is immutable as an accepted one is, but for its
// topics.
func TestRejectedAdrIsImmutableButForItsTopics(t *testing.T) {
	repo, root, _ := project(t)
	const path = "design/adrs/0001-rejected.md"
	adr := "---\nid: ADR-0001\ntitle: Rejected\nstatus: rejected\ndate: 2026-10-01\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-0001: Rejected\n\n## Decision\n\nWe do not do it.\n"
	if err := os.MkdirAll(filepath.Join(root, "design/adrs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, path), []byte(adr), 0o644); err != nil {
		t.Fatal(err)
	}

	if mode, reason := ModeOf(repo, filepath.Join(root, path)); mode != None || !strings.Contains(reason, "immutable") {
		t.Errorf("a rejected ADR is read-only: %q %q", mode, reason)
	}
	body := strings.Replace(adr, "We do not do it.", "We do it after all.", 1)
	var refused *RefusedError
	if _, err := Save(repo, &fakeGit{}, path, body, SaveOptions{Hash: Hash(adr), NoCommit: true, Now: now}); !errors.As(err, &refused) || !strings.Contains(refused.Reason, "immutable") {
		t.Errorf("a change to its body is refused: %v", err)
	}
	if got := read(t, filepath.Join(root, path)); got != adr {
		t.Errorf("the refused save writes nothing:\n%s", got)
	}

	withTopics := strings.Replace(adr, "superseded_by: []\n---\n", "superseded_by: []\ntopics: [code]\n---\n", 1)
	if _, err := Save(repo, &fakeGit{}, path, withTopics, SaveOptions{Hash: Hash(adr), NoCommit: true, Now: now}); err != nil {
		t.Fatalf("a change to its topics alone is saved: %v", err)
	}
	if got := read(t, filepath.Join(root, path)); got != withTopics {
		t.Errorf("the topics are written and nothing else changes:\n%s", got)
	}
}

// workflowFileProject is project with the baseline workflow file at the
// path the manifest's design folder gives it, on the main checkout and in
// S-0001's worktree, and that path.
func workflowFileProject(t *testing.T) (repo *workitem.Repo, root, wt, rel, content string) {
	t.Helper()
	repo, root, wt = project(t)
	data, err := os.ReadFile(filepath.Join("..", "workflow", "baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	rel, content = workflow.Rel(repo.Manifest.Layout["design"]), string(data)
	for _, dir := range []string{root, wt} {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repo, root, wt, rel, content
}

// workflowFileEdit is the workflow file with a sentence added to its prose.
func workflowFileEdit(t *testing.T, content string) string {
	t.Helper()
	edited := strings.Replace(content, "flai reads only the front matter.", "flai reads only the front matter. Edited.", 1)
	if edited == content {
		t.Fatal("the baseline workflow file no longer has the sentence the edit changes")
	}
	return edited
}

// I-0380, ADR-0301 §7: a save that names no story, as the dashboard's on
// main, refuses the workflow file, writes nothing, and runs no git; Show
// gives it read-only with the same reason.
func TestWorkflowFileSaveOnMainIsRefused(t *testing.T) {
	repo, root, wt, rel, content := workflowFileProject(t)
	if rel != "design/workflows/board.md" {
		t.Fatalf("the fixture's workflow file is %s", rel)
	}
	git := &fakeGit{branch: "main"}
	_, err := Save(repo, git, rel, workflowFileEdit(t, content), SaveOptions{Hash: Hash(content), Now: now})
	r, ok := IsRefused(err)
	if !ok {
		t.Fatalf("a save of the workflow file on main must be refused: %v", err)
	}
	for _, want := range []string{rel + " cannot be edited here", "changes only on a story's branch", "ADR-0301 §7", "name a story (--story S-nnnn)"} {
		if !strings.Contains(r.Reason, want) {
			t.Errorf("the refusal should say %q: %s", want, r.Reason)
		}
	}
	if r.Path != rel || len(r.Findings) != 0 {
		t.Errorf("the refusal names the file and carries no findings: %+v", r)
	}
	for _, dir := range []string{root, wt} {
		if got := read(t, filepath.Join(dir, rel)); got != content {
			t.Errorf("nothing is written in %s:\n%s", dir, got)
		}
	}
	if len(git.calls) != 0 {
		t.Errorf("a refused save runs no git: %v", git.calls)
	}

	doc, err := Show(repo, rel, "")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Mode != None || !strings.HasSuffix(r.Reason, doc.Reason) || doc.Content != content {
		t.Errorf("without a story the workflow file is read-only, for the refusal's reason: %q %q", doc.Mode, doc.Reason)
	}
	if doc, err := Show(repo, rel, "S-0001"); err != nil || doc.Mode != Full || doc.Reason != "" {
		t.Errorf("on a story's branch it is edited whole: %+v %v", doc, err)
	}
	if doc, err := Show(repo, planPath, ""); err != nil || doc.Mode != Full {
		t.Errorf("another design document stays editable on main: %+v %v", doc, err)
	}
}

// I-0380, ADR-0145: a save of the workflow file that names a story is
// checked in its worktree and committed on its branch there, and leaves the
// main checkout's file as it was.
func TestWorkflowFileSaveNamingAStoryCommitsOnItsBranch(t *testing.T) {
	repo, root, wt, rel, content := workflowFileProject(t)
	git := &fakeGit{branch: "story/S-0001"}
	res, err := Save(repo, git, rel, workflowFileEdit(t, content), SaveOptions{Hash: Hash(content), Story: "S-0001", Message: "reshape the board", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Committed || res.Commit != "abc1234" || res.Path != rel {
		t.Errorf("result: %+v", res)
	}
	if got := read(t, filepath.Join(wt, rel)); !strings.Contains(got, "front matter. Edited.") {
		t.Errorf("the worktree's workflow file carries the edit:\n%s", got)
	}
	if got := read(t, filepath.Join(root, rel)); got != content {
		t.Errorf("the main checkout's workflow file is left as it was:\n%s", got)
	}
	joined := strings.Join(git.calls, "\n")
	for _, w := range []string{wt + ": git add -- " + rel, wt + ": git commit --quiet -m docs: [S-0001] reshape the board -- " + rel} {
		if !strings.Contains(joined, w) {
			t.Errorf("git should have run %q; ran:\n%s", w, joined)
		}
	}
	if strings.Contains(joined, root+": git add") || strings.Contains(joined, root+": git commit") {
		t.Errorf("nothing is committed in the main checkout:\n%s", joined)
	}
}

// I-0380: file save reaches no markdown file, so it never saves the
// workflow file, with a story or without one, and writes nothing.
func TestWorkflowFileIsOutOfFileSavesReach(t *testing.T) {
	repo, root, wt, rel, content := workflowFileProject(t)
	if _, err := FileReach(repo.Manifest, rel); err == nil || !strings.Contains(err.Error(), rel+" is a markdown file") {
		t.Errorf("file save does not reach the workflow file: %v", err)
	}
	for _, story := range []string{"", "S-0001"} {
		git := &fakeGit{branch: "story/S-0001"}
		_, err := FileSave(repo, git, rel, workflowFileEdit(t, content), SaveOptions{Hash: Hash(content), Story: story, Now: now})
		if err == nil || !strings.Contains(err.Error(), "is a markdown file") {
			t.Errorf("file save with story %q: %v", story, err)
		}
	}
	for _, dir := range []string{root, wt} {
		if got := read(t, filepath.Join(dir, rel)); got != content {
			t.Errorf("nothing is written in %s:\n%s", dir, got)
		}
	}
}
