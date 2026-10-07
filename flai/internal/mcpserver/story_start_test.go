package mcpserver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// startClock is when story_start starts a story, an hour after the
// projects' stories were made.
var startClock = t0.Add(time.Hour)

// startConvention is a convention whose rules are rules.
func startConvention(rules string) string {
	return "---\ntitle: Git\nupdated: 2026-08-01\naudience: agent\norder: 10\nstatus: active\n---\n\n# Git\n\n## Rules\n" + rules + "\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n"
}

// startProject is a project with key at dir, with one epic, which it
// returns, and convention as its one convention, or no conventions folder,
// so that no pack can be built, when convention is "".
func startProject(t *testing.T, dir, key, convention string) (*workitem.Repo, *workitem.Item) {
	t.Helper()
	for _, d := range []string{"design/system", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(dir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeIn(t, dir, "system-flow.yaml", "version: 1\nname: "+key+"\nkey: "+key+"\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	writeIn(t, dir, "design/system/plan.md", "---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\ntext\n")
	if convention != "" {
		writeIn(t, dir, "design/conventions/README.md", "# Agent conventions\n\n| Order | File | Governs |\n|-------|------|---------|\n| 10 | [git.md](git.md) | Commits |\n")
		writeIn(t, dir, "design/conventions/git.md", convention)
	}
	repo, err := workitem.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic of " + key, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	return repo, epic
}

// storyOfEpic makes a story of epic with a criterion and touches, moved by
// the operator through states, and returns it as it is then.
func storyOfEpic(t *testing.T, repo *workitem.Repo, epic *workitem.Item, title string, touches []string, states ...string) *workitem.Item {
	t.Helper()
	s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Owner: "alex", Touches: touches, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	writeIn(t, filepath.Dir(s.Path), filepath.Base(s.Path), strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1))
	for _, st := range states {
		if s, err = repo.Get(s.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.Transition(s, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	if s, err = repo.Get(s.ID); err != nil {
		t.Fatal(err)
	}
	return s
}

// starter is an MCP server on repo as agent-x, at startClock, running git
// with runner.
func starter(t *testing.T, repo *workitem.Repo, runner execx.Runner, with func(*Options)) *fixture {
	t.Helper()
	opt := Options{Repo: repo, Runner: runner, Agent: "agent-x", Version: "test", Now: func() time.Time { return startClock }}
	if with != nil {
		with(&opt)
	}
	return &fixture{repo: repo, cs: connectServer(t, opt)}
}

// S-0274: story_start starts a ready story in one call as the server's
// agent: in progress, its epic following, its narrative and its branch
// opened in its worktree, linked as the server's RelativePaths says, and
// the answer holds every key flai story start --json gives, the pack whole
// when it fits as part 1 of 1. Integration: it runs real git.
func TestStoryStartStartsAReadyStory(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	t.Setenv("FLAI_SESSION", "s1")
	dir := t.TempDir()
	repo, epic := startProject(t, dir, "t", startConvention("- Commit at landing."))
	story := storyOfEpic(t, repo, epic, "Start me", []string{"docs/start"}, workitem.Ready)
	writeIn(t, dir, ".gitignore", ".flai-cache/\n")
	gitIn(t, dir, "init", "-q", "-b", "main")
	gitIn(t, dir, "add", "-A")
	gitIn(t, dir, "commit", "-q", "-m", "project")
	var asked atomic.Bool
	f := starter(t, repo, execx.System{}, func(o *Options) {
		o.RelativePaths = func() bool { asked.Store(true); return false }
	})

	out, failed := f.call(t, "story_start", map[string]any{"story": story.ID, "budget": "40KB"})
	if failed != "" {
		t.Fatalf("story_start: %s", failed)
	}
	for _, key := range []string{"story", "followed", "worktree", "branch", "from", "pack", "inbox"} {
		if _, ok := out[key]; !ok {
			t.Errorf("the answer has no %s: %v", key, out)
		}
	}
	if s := out["story"].(map[string]any); s["id"] != story.ID || s["status"] != workitem.InProgress || s["title"] != "Start me" {
		t.Errorf("story: %v", s)
	}
	if fl, _ := out["followed"].(map[string]any); fl["id"] != epic.ID || fl["from"] != workitem.Ready || fl["to"] != workitem.InProgress {
		t.Errorf("followed: %v", out["followed"])
	}
	if out["worktree"] != repo.WorktreePath(story.ID) || out["branch"] != storygit.Branch(story.ID) || out["from"] != storygit.FromMain {
		t.Errorf("worktree %v, branch %v, from %v", out["worktree"], out["branch"], out["from"])
	}
	if !asked.Load() {
		t.Error("the new worktree was linked without asking RelativePaths")
	}
	if head := gitIn(t, repo.WorktreePath(story.ID), "rev-parse", "--abbrev-ref", "HEAD"); head != storygit.Branch(story.ID) {
		t.Errorf("the worktree is on %q", head)
	}
	want, _ := ctxpack.ParseSize("40KB")
	pack := out["pack"].(map[string]any)
	if pack["story"] != story.ID || pack["part"] != 1.0 || pack["parts"] != 1.0 || pack["budget"] != float64(want) || len(pack["conventions"].([]any)) == 0 {
		t.Errorf("pack: %v", pack)
	}
	if in := out["inbox"].(map[string]any); in["agent"] != "agent-x" {
		t.Errorf("inbox: %v", in)
	}
	got, err := repo.Get(story.ID)
	if err != nil {
		t.Fatal(err)
	}
	if last := got.Transitions[len(got.Transitions)-1]; got.Status != workitem.InProgress || last.By != "agent-x" || last.At != startClock.Format(workitem.TimeFormat) {
		t.Errorf("story %s, last transition %+v", got.Status, last)
	}
	n, err := workitem.ReadNarrative(repo.NarrativePath(story.ID))
	if err != nil {
		t.Fatalf("no narrative: %v", err)
	}
	if n.Agent != "agent-x" || n.Session != "s1" {
		t.Errorf("narrative agent %q, session %q", n.Agent, n.Session)
	}
}

// S-0274: a held story is refused with the hold's reason and what clears
// it, and nothing is changed; a story in progress already is refused with
// what an agent wait_for_work sent does then.
func TestStoryStartRefusesAHeldStory(t *testing.T) {
	repo, epic := startProject(t, t.TempDir(), "t", startConvention("- Commit at landing."))
	theirs := storyOfEpic(t, repo, epic, "Theirs", []string{"docs/shared"}, workitem.Ready, workitem.InProgress)
	mine := storyOfEpic(t, repo, epic, "Mine", []string{"docs/shared/page.md"}, workitem.Ready)
	f := starter(t, repo, noGitRunner{}, nil)

	_, failed := f.call(t, "story_start", map[string]any{"story": mine.ID})
	if !strings.Contains(failed, mine.ID) || !strings.Contains(failed, theirs.ID+" moves to review") || !strings.Contains(failed, "nothing was changed") {
		t.Errorf("refusal: %q", failed)
	}
	after, err := repo.Get(mine.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != workitem.Ready || len(after.Transitions) != len(mine.Transitions) {
		t.Errorf("%s changed: %s, %d transitions, was %d", mine.ID, after.Status, len(after.Transitions), len(mine.Transitions))
	}
	if _, err := os.Stat(repo.NarrativePath(mine.ID)); err == nil {
		t.Errorf("%s has a narrative", mine.ID)
	}

	_, failed = f.call(t, "story_start", map[string]any{"story": theirs.ID})
	if !strings.Contains(failed, "in progress already") || !strings.Contains(failed, "call wait_for_work again") {
		t.Errorf("refusal of a story in progress: %q", failed)
	}
}

// S-0274: a step that fails after the move is the tool's error, naming the
// step and the command that finishes it; the story stays in progress.
func TestStoryStartSaysWhichStepFailed(t *testing.T) {
	repo, epic := startProject(t, t.TempDir(), "t", "")
	story := storyOfEpic(t, repo, epic, "No pack", []string{"docs/none"}, workitem.Ready)
	f := starter(t, repo, noGitRunner{}, nil)

	_, failed := f.call(t, "story_start", map[string]any{"story": story.ID})
	if !strings.Contains(failed, story.ID+" is in progress, but its prime step failed") || !strings.Contains(failed, "flai prime --story "+story.ID) {
		t.Errorf("step error: %q", failed)
	}
	if got, err := repo.Get(story.ID); err != nil || got.Status != workitem.InProgress {
		t.Errorf("the story is %v: %v", got, err)
	}
}

// ADR-0104: an answer that would be over ctxpack.PartLimit bytes with part 1
// of the pack holds the pack's header alone, with parts and no part, and
// prime answers part 1. Thirty conventions of 2.5 KB fill part 1, and ten
// threads awaiting the agent make the inbox 20 KB.
func TestStoryStartAnswersThePacksHeaderWhenPartOneDoesNotFit(t *testing.T) {
	dir := t.TempDir()
	rule := "- Commit at landing, with the story's ID in the subject line and the reason in the body.\n"
	for i := range 30 {
		writeIn(t, dir, fmt.Sprintf("design/conventions/c%02d.md", i), strings.Replace(startConvention(strings.Repeat(rule, 28)), "order: 10", fmt.Sprintf("order: %d", 20+i), 1))
	}
	repo, epic := startProject(t, dir, "t", startConvention("- Commit at landing."))
	story := storyOfEpic(t, repo, epic, "Large", []string{"docs/large"}, workitem.Ready)
	for i := range 10 {
		if _, err := threads.New(repo, threads.NewOptions{Title: fmt.Sprintf("Question %d", i), On: story.ID, Author: "alex", Text: strings.Repeat("Which library do we take here? ", 64), Now: t0}); err != nil {
			t.Fatal(err)
		}
	}
	f := starter(t, repo, noGitRunner{}, nil)

	out, failed := f.call(t, "story_start", map[string]any{"story": story.ID})
	if failed != "" {
		t.Fatalf("story_start: %s", failed)
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > ctxpack.PartLimit {
		t.Errorf("the answer is %d bytes, over %d", len(data), ctxpack.PartLimit)
	}
	if in := out["inbox"].(map[string]any); in["awaiting_you"] != 10.0 {
		t.Errorf("the inbox does not hold the threads: %v", in["awaiting_you"])
	}
	pack := out["pack"].(map[string]any)
	parts, _ := pack["parts"].(float64)
	if _, ok := pack["part"]; ok || parts < 2 || pack["story"] != story.ID || len(pack["conventions"].([]any)) != 0 || pack["readme"] != nil {
		t.Errorf("pack: %v", pack)
	}
	first, failed := f.call(t, "prime", map[string]any{"story": story.ID, "part": 1})
	if failed != "" || first["part"] != 1.0 || first["parts"] != parts || len(first["conventions"].([]any)) == 0 {
		t.Errorf("prime part 1: %v %s", first, failed)
	}
}

// ADR-0104: an answer over ctxpack.PartLimit bytes with part 1 of the pack
// and the whole inbox keeps part 1 and leaves out the inbox's oldest
// changes, counted in changes_omitted. Ten conventions of 2.5 KB make part
// 1, and forty stories the operator moved to ready and back make the inbox's changes.
func TestStoryStartLeavesOutOldChangesToKeepPartOne(t *testing.T) {
	dir := t.TempDir()
	rule := "- Commit at landing, with the story's ID in the subject line and the reason in the body.\n"
	for i := range 10 {
		writeIn(t, dir, fmt.Sprintf("design/conventions/c%02d.md", i), strings.Replace(startConvention(strings.Repeat(rule, 28)), "order: 10", fmt.Sprintf("order: %d", 20+i), 1))
	}
	repo, epic := startProject(t, dir, "t", startConvention("- Commit at landing."))
	story := storyOfEpic(t, repo, epic, "Started", []string{"docs/started"}, workitem.Ready)
	for i := range 40 {
		storyOfEpic(t, repo, epic, fmt.Sprintf("Returned story %02d, whose title is long enough to make its change in the inbox weigh several hundred bytes", i), []string{fmt.Sprintf("docs/c%02d", i)}, workitem.Ready, workitem.Backlog)
	}
	f := starter(t, repo, noGitRunner{}, nil)

	out, failed := f.call(t, "story_start", map[string]any{"story": story.ID})
	if failed != "" {
		t.Fatalf("story_start: %s", failed)
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > ctxpack.PartLimit {
		t.Errorf("the answer is %d bytes, over %d", len(data), ctxpack.PartLimit)
	}
	pack := out["pack"].(map[string]any)
	if pack["part"] != 1.0 || len(pack["conventions"].([]any)) == 0 {
		t.Errorf("pack is not part 1: part %v, %d conventions", pack["part"], len(pack["conventions"].([]any)))
	}
	in := out["inbox"].(map[string]any)
	changes, _ := in["changes"].([]any)
	omitted, _ := in["changes_omitted"].(float64)
	if omitted == 0 || len(changes) == 0 || int(omitted)+len(changes) < 80 {
		t.Errorf("inbox: %d changes, %v omitted", len(changes), omitted)
	}
	if len(changes) > 0 && !strings.Contains(fmt.Sprint(changes[len(changes)-1]), "Returned story 39") {
		t.Errorf("the newest change was not kept: %v", changes[len(changes)-1])
	}
}

// S-0274, S-0101: in a folder of projects, story_start takes project and
// starts the story there, as the folder server's agent.
func TestStoryStartTakesTheProjectInAFolder(t *testing.T) {
	root := t.TempDir()
	startProject(t, filepath.Join(root, "alpha"), "alpha", startConvention("- Commit at landing."))
	beta, epic := startProject(t, filepath.Join(root, "org", "beta"), "beta", startConvention("- Commit at landing."))
	story := storyOfEpic(t, beta, epic, "Beta work", []string{"docs/beta"}, workitem.Ready)
	f := folderSetupWith(t, root, func(o *Options) { o.Runner = noGitRunner{} })

	if _, failed := f.call(t, "story_start", map[string]any{"story": story.ID}); !strings.Contains(failed, "name one with project: alpha (alpha), beta (org/beta)") {
		t.Errorf("no project named: %q", failed)
	}
	out, failed := f.call(t, "story_start", map[string]any{"project": "beta", "story": story.ID})
	if failed != "" {
		t.Fatalf("story_start: %s", failed)
	}
	if s := out["story"].(map[string]any); s["id"] != story.ID || s["status"] != workitem.InProgress {
		t.Errorf("story: %v", s)
	}
	if out["pack"].(map[string]any)["story"] != story.ID || out["inbox"].(map[string]any)["agent"] != "claude" {
		t.Errorf("pack or inbox: %v", out)
	}
	got, err := beta.Get(story.ID)
	if err != nil {
		t.Fatal(err)
	}
	if last := got.Transitions[len(got.Transitions)-1]; got.Status != workitem.InProgress || last.By != "claude" {
		t.Errorf("beta's story %s, last transition %+v", got.Status, last)
	}
}

// The server tells an idle agent to pull with story_start, not with
// item_move and flai stream open, in its instructions and in wait_for_work.
func TestTheServersPullWithStoryStart(t *testing.T) {
	root := t.TempDir()
	startProject(t, filepath.Join(root, "alpha"), "alpha", "")
	for name, fx := range map[string]*folderFixture{"project": {cs: setup(t).cs}, "folder": folderSetup(t, root)} {
		res, err := fx.cs.ListTools(t.Context(), nil)
		if err != nil {
			t.Fatal(err)
		}
		texts := map[string]string{"instructions": fx.cs.InitializeResult().Instructions}
		for _, tool := range res.Tools {
			if tool.Name == "wait_for_work" {
				texts["wait_for_work"] = tool.Description
			}
		}
		for what, text := range texts {
			if !strings.Contains(text, "story_start") || !strings.Contains(text, "call wait_for_work again") || strings.Contains(text, "flai stream open") || strings.Contains(text, "item_move it to in-progress") {
				t.Errorf("%s's %s does not pull with story_start: %s", name, what, text)
			}
		}
	}
}
