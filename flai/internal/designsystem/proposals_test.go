package designsystem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const (
	mainTokens   = "design/system/design-system.tokens.json"
	mainGuidance = "design/system/design-system.md"
	plainTokens  = `{"color":{"primary":{"$type":"color","$value":{"colorSpace":"srgb","components":[0,0,1],"hex":"#0000ff"}}}}`
)

// storyFixture is a story the fixture project holds: its state, and
// whether it has a worktree.
type storyFixture struct {
	title, status string
	worktree      bool
}

// designProject is a main checkout with E-0001, the stories given as
// S-0001 onwards, and T-0001 under S-0001, each story with a worktree
// linked to it as git links one, its manifest naming its token file at
// tokens. It answers the repository and its root.
func designProject(t *testing.T, tokens string, stories ...storyFixture) (*workitem.Repo, string) {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	writeFile(t, filepath.Join(root, "system-flow.yaml"), manifest)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "docs"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Now: now}); err != nil {
		t.Fatal(err)
	}
	for _, s := range stories {
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: s.title, Parent: "E-0001", Now: now})
		if err != nil {
			t.Fatal(err)
		}
		it.Status = s.status
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
		if !s.worktree {
			continue
		}
		wt := repo.WorktreePath(it.ID)
		writeFile(t, filepath.Join(wt, ".git"), "gitdir: "+filepath.Join(root, ".git", "worktrees", it.ID)+"\n")
		writeFile(t, filepath.Join(wt, "system-flow.yaml"), manifest+"design:\n  tokens: "+tokens+"\n")
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Task", Parent: "S-0001", Now: now}); err != nil {
		t.Fatal(err)
	}
	return repo, root
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0384: with a story, the design system is read from its worktree, at
// the paths the worktree's manifest gives, not from the main checkout.
func TestShowStoryReadsItsWorktree(t *testing.T) {
	repo, root := designProject(t, "design/proposed.tokens.json", storyFixture{"Proposal", workitem.InProgress, true})
	wt := repo.WorktreePath("S-0001")
	writeFile(t, filepath.Join(wt, "design/proposed.tokens.json"), plainTokens)
	writeFile(t, filepath.Join(wt, mainGuidance), "# Design system\n\n## Colour\n")
	writeFile(t, filepath.Join(root, "design/proposed.tokens.json"), "{")

	shown, id, err := ShowStory(repo, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	if id != "S-0001" {
		t.Errorf("the story's own ID: %q", id)
	}
	if !shown.Exists || shown.Path != "design/proposed.tokens.json" || shown.Guidance != mainGuidance {
		t.Errorf("the worktree's design system: %+v", shown)
	}
	if !slices.Equal(shown.Headings, []string{"Design system", "Colour"}) {
		t.Errorf("the worktree's guidance headings: %q", shown.Headings)
	}
	if len(shown.Themes) != 2 || len(shown.Themes[0].Tokens) != 1 || shown.Themes[0].Tokens[0].Path != "color.primary" {
		t.Errorf("the worktree's tokens: %+v", shown.Themes)
	}
}

// S-0384: what has no design system of its own to propose is refused, with
// the reason and what to do.
func TestShowStoryRefuses(t *testing.T) {
	repo, _ := designProject(t, mainTokens,
		storyFixture{"Open", workitem.InProgress, true},
		storyFixture{"No worktree", workitem.Ready, false},
		storyFixture{"Done", workitem.Done, true},
		storyFixture{"Cancelled", workitem.Cancelled, true},
	)
	cases := []struct {
		name, story string
		want        []string
	}{
		{"an ID no item has", "S-0099", []string{"read the design system of S-0099", "not found", "name a story of this project"}},
		{"what is not an ID", "colour", []string{"is not a work item ID"}},
		{"a task", "T-0001", []string{"T-0001: it is a task; name a story (S-nnnn)"}},
		{"a done story", "S-0003", []string{"S-0003: it is done", "name an open story"}},
		{"a cancelled story", "S-0004", []string{"S-0004: it is cancelled", "name an open story"}},
		{"a story without a worktree", "S-0002", []string{"S-0002: it has no worktree at " + repo.WorktreePath("S-0002"), "flai stream open S-0002"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			shown, id, err := ShowStory(repo, c.story)
			if err == nil {
				t.Fatalf("read %s: %+v", id, shown)
			}
			for _, w := range c.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("the error %q does not say %q", err, w)
				}
			}
		})
	}
}

// fakeGit answers the git calls storygit.StoryChanges makes in a work tree
// whose main branch is main: each branch's changed files, as git diff
// --name-status -z gives them, and each folder's uncommitted ones, as git
// status --porcelain gives them.
type fakeGit struct {
	branches map[string]string // branch name to its name-status output
	dirty    map[string]string // folder to its porcelain status
}

func (f *fakeGit) Run(dir, _ string, args ...string) (string, error) {
	switch {
	case slices.Equal(args, []string{"rev-parse", "--is-inside-work-tree"}):
		return "true", nil
	case len(args) == 4 && args[1] == "--verify":
		if _, ok := f.branches[strings.TrimPrefix(args[3], "refs/heads/")]; ok {
			return "", nil
		}
		return "", errors.New("exit status 1")
	case slices.Equal(args, []string{"rev-parse", "--abbrev-ref", "HEAD"}):
		return "main", nil
	case args[0] == "merge-base":
		return "base", nil
	case len(args) > 1 && args[1] == "--short":
		return "abc1234", nil
	case args[0] == "rev-list":
		return "1", nil
	case len(args) > 1 && args[0] == "diff" && args[1] == "--name-status":
		return f.branches[args[len(args)-1]], nil
	case args[0] == "status":
		return f.dirty[dir], nil
	}
	return "", nil
}

func (f *fakeGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return f.Run(dir, name, args...)
}

func (f *fakeGit) LookPath(name string) (string, error) { return name, nil }

// S-0384: the proposals are the open stories whose branches or worktrees
// change the token file or the guidance, each with the paths it changes.
func TestProposalsListsOpenStoriesChangingTheDesignSystem(t *testing.T) {
	repo, _ := designProject(t, mainTokens,
		storyFixture{"Tokens on the branch", workitem.InProgress, true},
		storyFixture{"Elsewhere", workitem.InProgress, true},
		storyFixture{"Done", workitem.Done, false},
		storyFixture{"Cancelled", workitem.Cancelled, false},
		storyFixture{"Both", workitem.Review, true},
		storyFixture{"Not started", workitem.Backlog, false},
	)
	git := &fakeGit{
		branches: map[string]string{
			"story/S-0001": "M\x00" + mainTokens + "\x00M\x00flai/main.go\x00",
			"story/S-0002": "M\x00flai/main.go\x00A\x00design/system/other.md\x00",
			"story/S-0003": "M\x00" + mainTokens + "\x00",
			"story/S-0004": "M\x00" + mainGuidance + "\x00",
			"story/S-0005": "M\x00" + mainTokens + "\x00",
		},
		dirty: map[string]string{
			repo.WorktreePath("S-0005"): " M " + mainGuidance + "\n",
		},
	}
	got, err := Proposals(git, repo)
	if err != nil {
		t.Fatal(err)
	}
	want := []Proposal{
		{ID: "S-0001", Title: "Tokens on the branch", Status: workitem.InProgress, Paths: []string{mainTokens}},
		{ID: "S-0005", Title: "Both", Status: workitem.Review, Paths: []string{mainTokens, mainGuidance}},
	}
	if len(got) != len(want) {
		t.Fatalf("proposals:\n got %+v\nwant %+v", got, want)
	}
	for i := range want {
		if got[i].ID != want[i].ID || got[i].Title != want[i].Title || got[i].Status != want[i].Status || !slices.Equal(got[i].Paths, want[i].Paths) {
			t.Errorf("proposal %d:\n got %+v\nwant %+v", i, got[i], want[i])
		}
	}
	data, err := json.Marshal(got[0])
	if err != nil {
		t.Fatal(err)
	}
	if s := string(data); s != `{"id":"S-0001","title":"Tokens on the branch","status":"in-progress","paths":["`+mainTokens+`"]}` {
		t.Errorf("a proposal's JSON: %s", s)
	}
}

// S-0384: a guidance moved away is a change to it.
func TestProposalsOfARename(t *testing.T) {
	repo, _ := designProject(t, mainTokens, storyFixture{"Rename", workitem.InProgress, true})
	git := &fakeGit{branches: map[string]string{"story/S-0001": "R100\x00" + mainGuidance + "\x00design/system/look.md\x00"}}
	got, err := Proposals(git, repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !slices.Equal(got[0].Paths, []string{mainGuidance}) {
		t.Errorf("a renamed guidance: %+v", got)
	}
}

// S-0384: no open story changing the design system is an empty list, not
// nil, so it serialises as [].
func TestProposalsNoneIsAnEmptyList(t *testing.T) {
	repo, _ := designProject(t, mainTokens, storyFixture{"Elsewhere", workitem.InProgress, true})
	git := &fakeGit{branches: map[string]string{"story/S-0001": "M\x00flai/main.go\x00"}}
	got, err := Proposals(git, repo)
	if err != nil {
		t.Fatal(err)
	}
	if data, err := json.Marshal(got); err != nil || string(data) != "[]" {
		t.Errorf("no proposals serialise as []: %s %v", data, err)
	}
}

func TestChangedOfAnUntrackedFolderCoversWhatIsUnderIt(t *testing.T) {
	want := []string{mainTokens, mainGuidance}
	for _, c := range []struct {
		changed []string
		got     []string
	}{
		{[]string{"design/system/"}, want},
		{[]string{"design/"}, want},
		{[]string{"design/sys/"}, nil},
		{[]string{"docs/", mainGuidance}, []string{mainGuidance}},
	} {
		if got := changedOf(want, c.changed); !slices.Equal(got, c.got) {
			t.Errorf("changedOf(%v) = %v, want %v", c.changed, got, c.got)
		}
	}
}
