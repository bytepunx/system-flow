package mcpserver

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/gittest"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0198: issue_story makes a backlog story from an open issue as flai issue
// story does, and refuses a closed issue or one an open story already links,
// naming it. S-0203: the story is a draft carrying the issue's cost of delay
// inputs, and the issue's Remediation section names it.
func TestIssueStoryMakesAStoryFromAnOpenIssue(t *testing.T) {
	f := setup(t)
	defect, err := issues.New(f.repo, issues.NewOptions{Title: "Fixture was ignored", Class: "defect", Cost: "20m", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	defect.Body += "\n## Impact\n- revenue_per_week: 1200\n"
	if err := defect.Save(); err != nil {
		t.Fatal(err)
	}

	out, failed := f.call(t, "issue_story", map[string]any{"id": "I-1"})
	if failed != "" {
		t.Fatal(failed)
	}
	id, _ := out["id"].(string)
	want := map[string]any{"title": "Fixture was ignored", "nature": "remediation", "issue": "I-0001", "path": "wip/kanban/stories/" + id + "-fixture-was-ignored.md", "draft": true}
	for k, v := range want {
		if out[k] != v {
			t.Errorf("issue_story %s = %v, want %v", k, out[k], v)
		}
	}
	story, err := f.repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if story.Type != workitem.Story || story.Status != workitem.Backlog || story.Nature != "remediation" || story.Owner != "claude" || !issues.Links(story.Body, defect) {
		t.Errorf("the story should be a backlog remediation story the agent owns that links I-0001: %+v\n%s", story, story.Body)
	}
	if !story.Draft {
		t.Error("the story should be a draft")
	}
	cod := story.CostOfDelay
	if cod == nil || cod.Inputs == nil || cod.Inputs.TimeLostPerCycle != "20m" || cod.Inputs.RevenuePerWeek == nil || *cod.Inputs.RevenuePerWeek != 1200 || cod.Inputs.By != "flai" || cod.By != "" {
		t.Errorf("the story should carry the issue's time lost per cycle and its Impact's revenue, set by flai: %+v", cod)
	}
	for _, want := range []string{"time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h", "revenue_per_week 1200 carried over from I-0001's Impact section."} {
		if !strings.Contains(story.Body, want) {
			t.Errorf("the story's Notes should say %q:\n%s", want, story.Body)
		}
	}
	if is, err := issues.Get(f.repo, "I-0001"); err != nil || !strings.Contains(is.Body, "## Remediation\n\nStory "+id+" remediates this issue, created from it at ") {
		t.Fatalf("the issue's Remediation section should name %s: %v %+v", id, err, is)
	}

	if _, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001"}); !strings.Contains(failed, "I-0001 is already linked by open story "+id) {
		t.Errorf("an issue an open story links should be refused, naming it: %q", failed)
	}

	closed, err := issues.New(f.repo, issues.NewOptions{Title: "Go not on PATH", Class: "blocker", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if err := issues.Close(closed, "fixed", t0); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "issue_story", map[string]any{"id": closed.ID}); !strings.Contains(failed, "I-0002 is closed, so no story was made for it") {
		t.Errorf("a closed issue should be refused: %q", failed)
	}
	if list, _ := f.repo.List(false); countStories(list) != 2 {
		t.Errorf("a refusal should make no story: %d stories", countStories(list))
	}
}

// S-0227: issue_story charges what an analyzer spent on the issue to the
// story and its epic, as the CLI does, and the issue keeps it; a story made
// from an issue with no usage has none.
func TestIssueStoryCarriesTheIssuesStrategicUsage(t *testing.T) {
	f := setup(t)
	epic := f.story.Parent
	for _, title := range []string{"Review waits a day", "Fixture was ignored"} {
		if _, err := issues.New(f.repo, issues.NewOptions{Title: title, Class: "efficiency", Now: t0}); err != nil {
			t.Fatal(err)
		}
	}
	models := []usage.Model{{Model: "claude-opus-5-5", Input: 1, Output: 200, CacheRead: 5000, CacheWrite: 300, Cost: 0.42}}
	if _, err := issues.ChargeStrategic(f.repo, "I-0001", workitem.ActivityAnalyzer, &usage.Usage{Seconds: 723, Models: models}); err != nil {
		t.Fatal(err)
	}
	out, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001", "epic": epic})
	if failed != "" {
		t.Fatal(failed)
	}
	id, _ := out["id"].(string)
	for _, of := range []string{id, epic} {
		it, err := f.repo.Get(of)
		if err != nil {
			t.Fatal(err)
		}
		if u := it.Usage; u == nil || len(u.Strategic) != 1 || u.Strategic[0].Kind != workitem.ActivityAnalyzer || u.Strategic[0].Seconds != 723 || !u.Strategic[0].Estimated ||
			len(u.Strategic[0].Models) != 1 || u.Strategic[0].Models[0] != models[0] {
			t.Errorf("%s should carry the issue's analyzer entry: %+v", of, it.Usage)
		}
	}
	is, err := issues.Get(f.repo, "I-0001")
	if err != nil {
		t.Fatal(err)
	}
	if u := is.Usage; u == nil || len(u.Strategic) != 1 || u.Strategic[0].Seconds != 723 || u.Strategic[0].Cost() != 0.42 {
		t.Errorf("the issue should keep its entry: %+v", is.Usage)
	}

	out, failed = f.call(t, "issue_story", map[string]any{"id": "I-0002"})
	if failed != "" {
		t.Fatal(failed)
	}
	if story, err := f.repo.Get(out["id"].(string)); err != nil || story.Usage != nil {
		t.Errorf("a story made from an issue with no usage has none: %v %+v", err, story)
	}
}

// S-0198: story reads the issue from that story's worktree, where the issues
// it recorded are until it is accepted, and names the story there (S-0203);
// the story is made in the main wip.
func TestIssueStoryReadsAStorysWorktree(t *testing.T) {
	f := setup(t)
	wt := *f.repo
	wt.Root = f.repo.WorktreePath(f.story.ID)
	if err := os.MkdirAll(wt.Root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := issues.New(&wt, issues.NewOptions{Title: "Only on the branch", Class: "efficiency", Story: f.story.ID, Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001"}); !strings.Contains(failed, "I-0001 not found") {
		t.Errorf("without story the issue is read from the project, which has none: %q", failed)
	}
	out, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001", "story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["nature"] != "improvement" || !strings.HasPrefix(out["path"].(string), "wip/kanban/stories/") {
		t.Errorf("issue_story with story: %v", out)
	}
	if _, err := os.Stat(filepath.Join(wt.Root, "wip")); err == nil {
		t.Error("the story should be made in the main checkout's wip, not the worktree's")
	}
	if is, err := issues.Get(&wt, "I-0001"); err != nil || !strings.Contains(is.Body, "Story "+out["id"].(string)+" remediates this issue") {
		t.Errorf("the issue in the worktree should name the story: %v %+v", err, is)
	}
}

// S-0224: issue_new files an analysis report's finding with its impact and
// links the report; filed again from another report under the same title it
// bumps the open issue rather than open a second, and from the same report it
// changes nothing. issue_story then carries the impact over.
func TestIssueNewFilesAFindingAndBumpsAnOpenIssueOfTheSameTitle(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	finding := map[string]any{"title": "Review waits a day for the operator", "class": "efficiency", "cost": "30m",
		"time_lost_per_cycle": "6h", "revenue_per_week": 1200, "evidence": "12 stories waited 18h on average in review",
		"report": "design/analysis/2026-10-06-bottlenecks.md"}
	out, failed := f.call(t, "issue_new", finding)
	if failed != "" {
		t.Fatal(failed)
	}
	want := map[string]any{"id": "I-0001", "title": "Review waits a day for the operator", "count": float64(1), "outcome": "opened",
		"path": "design/issues/I-0001-review-waits-a-day-for-the-operator.md"}
	for k, v := range want {
		if out[k] != v {
			t.Errorf("issue_new %s = %v, want %v", k, out[k], v)
		}
	}
	is, err := issues.Get(f.repo, "I-0001")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Impact\n12 stories waited 18h on average in review\n\n- revenue_per_week: 1200\n- time_lost_per_cycle: 6h\n", "Report: design/analysis/2026-10-06-bottlenecks.md."} {
		if !strings.Contains(is.Body, want) {
			t.Errorf("the issue should say %q:\n%s", want, is.Body)
		}
	}
	if links := issues.ReportLinks(is); len(links) != 1 || links[0] != "../analysis/2026-10-06-bottlenecks.md" {
		t.Errorf("the Remediation section should link the report: %v\n%s", links, is.Body)
	}
	if data, err := os.ReadFile(filepath.Join(f.repo.Root, "design/issues/summary.md")); err != nil || !strings.Contains(string(data), "I-0001") {
		t.Errorf("summary.md should be regenerated with I-0001: %v\n%s", err, data)
	}

	finding["report"], finding["penalty_per_week"] = "design/analysis/2026-10-13-bottlenecks.md", 300
	delete(finding, "revenue_per_week")
	out, failed = f.call(t, "issue_new", finding)
	if failed != "" {
		t.Fatal(failed)
	}
	if out["id"] != "I-0001" || out["outcome"] != "bumped" || out["count"] != float64(2) {
		t.Errorf("a finding of an open issue's title should bump it: %v", out)
	}
	list, _ := issues.List(f.repo)
	if len(list) != 1 {
		t.Fatalf("no second issue should be opened: %d issues", len(list))
	}
	is = list[0]
	for _, want := range []string{"- revenue_per_week: 1200\n- penalty_per_week: 300\n- time_lost_per_cycle: 6h\n", "Report: design/analysis/2026-10-13-bottlenecks.md."} {
		if !strings.Contains(is.Body, want) {
			t.Errorf("the bumped issue should say %q:\n%s", want, is.Body)
		}
	}
	if links := issues.ReportLinks(is); len(links) != 2 {
		t.Errorf("the Remediation section should link both reports: %v", links)
	}

	out, failed = f.call(t, "issue_new", finding)
	if failed != "" {
		t.Fatal(failed)
	}
	if out["outcome"] != "already recorded" || out["count"] != float64(2) {
		t.Errorf("a finding filed again from the same report should change nothing: %v", out)
	}

	made, failed := f.call(t, "issue_story", map[string]any{"id": "I-0001"})
	if failed != "" {
		t.Fatal(failed)
	}
	story, err := f.repo.Get(made["id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if cod := story.CostOfDelay; cod == nil || cod.Inputs == nil || cod.Inputs.RevenuePerWeek == nil || *cod.Inputs.RevenuePerWeek != 1200 ||
		cod.Inputs.PenaltyPerWeek == nil || *cod.Inputs.PenaltyPerWeek != 300 || cod.Inputs.TimeLostPerCycle != "6h" {
		t.Errorf("the story should carry the finding's impact as its cost of delay inputs: %+v", cod)
	}
}

// S-0224: issue_bump records a finding the analyzer judged the same as an
// issue under another title: it links the report and replaces the figures it
// gives in the Impact section, keeping the others.
func TestIssueBumpLinksAReportAndUpdatesTheImpact(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	if _, err := issues.New(f.repo, issues.NewOptions{Title: "Releases slip", Class: "defect", Cost: "20m",
		Impact: issues.Impact{RevenuePerWeek: "100", TimeLostPerCycle: "2h"}, Now: t0}); err != nil {
		t.Fatal(err)
	}
	out, failed := f.call(t, "issue_bump", map[string]any{"id": "I-1", "cost": "40m", "note": "two releases slipped this cycle",
		"penalty_per_week": 300, "revenue_per_week": 250.5, "evidence": "v1.4 and v1.5 shipped a week late", "report": "design/analysis/2026-10-06-risk.md"})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["id"] != "I-0001" || out["outcome"] != "bumped" || out["count"] != float64(2) || out["path"] != "design/issues/I-0001-releases-slip.md" {
		t.Errorf("issue_bump: %v", out)
	}
	is, err := issues.Get(f.repo, "I-0001")
	if err != nil {
		t.Fatal(err)
	}
	if is.Count != 2 || is.Cost != "30m" {
		t.Errorf("count and average cost: %d %s", is.Count, is.Cost)
	}
	for _, want := range []string{"v1.4 and v1.5 shipped a week late\n\n- revenue_per_week: 250.5\n- penalty_per_week: 300\n- time_lost_per_cycle: 2h\n", "Report: design/analysis/2026-10-06-risk.md.\ntwo releases slipped this cycle"} {
		if !strings.Contains(is.Body, want) {
			t.Errorf("the issue should say %q:\n%s", want, is.Body)
		}
	}
	if links := issues.ReportLinks(is); len(links) != 1 || links[0] != "../analysis/2026-10-06-risk.md" {
		t.Errorf("the Remediation section should link the report: %v", links)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, "design/issues/summary.md")); err != nil {
		t.Errorf("summary.md should be regenerated: %v", err)
	}
}

// S-0224: issue_new and issue_bump refuse what flai issue new and bump
// refuse, with the same errors, and write nothing.
func TestIssueNewAndBumpRefuseAndWriteNothing(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"title": "Bad class", "class": "annoyance"}, "--class must be one of defect, blocker, efficiency, impression"},
		{map[string]any{"title": "Bad amount", "class": "defect", "revenue_per_week": -5}, `impact revenue_per_week "-5" is not an amount of zero or more`},
		{map[string]any{"title": "Bad duration", "class": "defect", "time_lost_per_cycle": "soon"}, `impact time_lost_per_cycle "soon" is not a duration longer than zero`},
		{map[string]any{"title": "Bad report", "class": "defect", "report": "notes/risk.md"}, `report "notes/risk.md" is not an analysis report`},
		{map[string]any{"title": "Bad story", "class": "defect", "story": "T-0001"}, `story "T-0001" is not a story ID`},
	} {
		if _, failed := f.call(t, "issue_new", c.args); !strings.Contains(failed, c.want) {
			t.Errorf("issue_new %v should be refused with %q: %q", c.args, c.want, failed)
		}
	}
	if entries, err := os.ReadDir(filepath.Join(f.repo.Root, "design/issues")); err == nil && len(entries) > 0 {
		t.Errorf("a refusal should write nothing: %v", entries)
	}

	is, err := issues.New(f.repo, issues.NewOptions{Title: "Kept as it was", Class: "efficiency", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(is.Path)
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"id": is.ID, "penalty_per_week": -1}, `impact penalty_per_week "-1" is not an amount of zero or more`},
		{map[string]any{"id": is.ID, "report": "design/analysis/README.md"}, "is not an analysis report"},
		{map[string]any{"id": is.ID, "cost": "a while"}, "--cost must be a duration like 20m"},
		{map[string]any{"id": "I-0009"}, "I-0009 not found"},
	} {
		if _, failed := f.call(t, "issue_bump", c.args); !strings.Contains(failed, c.want) {
			t.Errorf("issue_bump %v should be refused with %q: %q", c.args, c.want, failed)
		}
	}
	if after, _ := os.ReadFile(is.Path); string(after) != string(before) {
		t.Errorf("a refused bump should leave the issue as it was:\n%s", after)
	}
	if err := issues.Close(is, "fixed", t0); err != nil {
		t.Fatal(err)
	}
	closed, _ := os.ReadFile(is.Path)
	if _, failed := f.call(t, "issue_bump", map[string]any{"id": is.ID}); !strings.Contains(failed, "I-0001 is closed") {
		t.Errorf("a closed issue should be refused: %q", failed)
	}
	if after, _ := os.ReadFile(is.Path); string(after) != string(closed) {
		t.Errorf("a refused bump of a closed issue should leave it as it was:\n%s", after)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, "design/issues/summary.md")); err == nil {
		t.Error("a refusal should not regenerate summary.md")
	}
}

// S-0224: the story an occurrence is recorded for defaults to FLAI_STORY, as
// flai issue new's does, and the issue is written in that story's worktree.
func TestIssueNewRecordsForTheSessionsStoryInItsWorktree(t *testing.T) {
	f := setup(t)
	t.Setenv("FLAI_STORY", f.story.ID)
	wt := *f.repo
	wt.Root = f.repo.WorktreePath(f.story.ID)
	if err := os.MkdirAll(wt.Root, 0o755); err != nil {
		t.Fatal(err)
	}
	out, failed := f.call(t, "issue_new", map[string]any{"title": "Fixture was ignored", "class": "defect"})
	if failed != "" {
		t.Fatal(failed)
	}
	rel, _ := filepath.Rel(f.repo.Root, wt.Root)
	if out["story"] != f.story.ID || out["path"] != filepath.ToSlash(rel)+"/design/issues/I-0001-fixture-was-ignored.md" {
		t.Errorf("issue_new should record for FLAI_STORY in its worktree: %v", out)
	}
	is, err := issues.Get(&wt, "I-0001")
	if err != nil {
		t.Fatal(err)
	}
	if got := issues.Stories(is); len(got) != 1 || got[0] != f.story.ID {
		t.Errorf("the instance should name %s: %v", f.story.ID, got)
	}
	if _, err := issues.Get(f.repo, "I-0001"); err == nil {
		t.Error("the issue should not be written in the main checkout")
	}
}

// storyGit makes f's project a git repository whose main branch holds it,
// with ADR 0001 and an ADR index, and checks the story's branch out in its
// worktree, which it returns as the project there. Integration: it runs real
// git.
func storyGit(t *testing.T, f *fixture) *workitem.Repo {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	gittest.Identity(t, "t", "t@t")
	root := f.repo.Root
	writeIn(t, root, ".gitignore", ".flai-cache/\n")
	writeIn(t, root, "design/adrs/0001-first.md", adrDoc("0001", "First"))
	writeIn(t, root, "design/adrs/README.md", adrIndex)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "project")
	wt := *f.repo
	wt.Root = f.repo.WorktreePath(f.story.ID)
	gitIn(t, root, "worktree", "add", "-q", "-b", storygit.Branch(f.story.ID), wt.Root, "main")
	return &wt
}

// adrIndex is design/adrs/README.md with ADR 0001's row.
const adrIndex = "---\ntitle: ADRs\nupdated: 2026-09-01\nstatus: active\n---\n\n# Architecture Decision Records\n\n| ADR | Title | Status |\n|-----|-------|--------|\n| [0001](0001-first.md) | First | accepted |\n"

// adrDoc is an accepted ADR numbered n.
func adrDoc(n, title string) string {
	return "---\nid: ADR-" + n + "\ntitle: " + title + "\nstatus: accepted\ndate: 2026-09-01\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-" + n + " " + title + "\n\n" + adrBody
}

// adrBody is an ADR's sections below its heading.
const adrBody = "## Context\n\nc\n\n## Decision\n\nd.\n\n## Consequences\n\ne\n\n## Alternatives considered\n\nf\n"

// lastCommit is the subject, the message, and the files of the commit at
// the head of dir's branch, and how many commits the branch has beyond main.
func lastCommit(t *testing.T, dir string) (message string, files []string, ahead string) {
	t.Helper()
	message = gitIn(t, dir, "log", "-1", "--format=%B")
	files = strings.Fields(gitIn(t, dir, "show", "--name-only", "--format=", "HEAD"))
	return message, files, gitIn(t, dir, "rev-list", "--count", "main..HEAD")
}

// strs is a decoded JSON list of strings.
func strs(v any) []string {
	out := []string{}
	list, _ := v.([]any)
	for _, e := range list {
		s, _ := e.(string)
		out = append(out, s)
	}
	return out
}

// S-0275: with commit, issue_new, issue_bump, and issue_close each write in
// the story's worktree, commit the issue and summary.md, and only those, on
// the story's branch with the story's prefix and the trailers, and widen the
// story's touches; an occurrence already recorded commits nothing.
func TestIssueToolsCommitOnTheStorysBranch(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	wt := storyGit(t, f)
	writeIn(t, wt.Root, "design/system/plan.md", "---\ntitle: Plan\n---\n\n# Plan\n\nhalf done\n")
	issue, summary := "design/issues/I-0001-fixture-was-ignored.md", "design/issues/summary.md"
	id := f.story.ID

	out, failed := f.call(t, "issue_new", map[string]any{"title": "Fixture was ignored", "class": "defect", "story": id,
		"report": "design/analysis/2026-10-06-risk.md", "commit": true, "trailers": []string{"Co-Authored-By: t <t@t>"}})
	if failed != "" {
		t.Fatal(failed)
	}
	rel, _ := filepath.Rel(f.repo.Root, wt.Root)
	if out["id"] != "I-0001" || out["outcome"] != "opened" || out["story"] != id || out["path"] != filepath.ToSlash(rel)+"/"+issue {
		t.Errorf("issue_new: %v", out)
	}
	subject := "docs: [" + id + "] record I-0001 Fixture was ignored"
	if field(out, "commit.subject") != subject || !reflect.DeepEqual(strs(field(out, "commit.paths")), []string{issue, summary}) {
		t.Errorf("commit: %v", out["commit"])
	}
	if got := strs(out["touches_added"]); !reflect.DeepEqual(got, []string{issue, summary}) {
		t.Errorf("touches_added %v", got)
	}
	message, files, ahead := lastCommit(t, wt.Root)
	if message != subject+"\n\nCo-Authored-By: t <t@t>" || !reflect.DeepEqual(files, []string{issue, summary}) || ahead != "1" {
		t.Errorf("the story's branch should hold one commit of the issue and the summary: %q %v %s", message, files, ahead)
	}
	if hash := gitIn(t, wt.Root, "rev-parse", "HEAD"); field(out, "commit.hash") != hash {
		t.Errorf("commit.hash %v, head %s", field(out, "commit.hash"), hash)
	}
	if st := gitIn(t, wt.Root, "status", "--porcelain"); st != "M design/system/plan.md" {
		t.Errorf("the worktree's other change should stay uncommitted, and nothing else: %q", st)
	}
	story, err := f.repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"flai/internal/mcpserver", issue, summary}; !reflect.DeepEqual(story.Touches, want) {
		t.Errorf("the story's touches in the project %v, want %v", story.Touches, want)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, "design/issues")); err == nil {
		t.Error("nothing should be written in the main checkout's design/issues")
	}

	out, failed = f.call(t, "issue_new", map[string]any{"title": "Fixture was ignored", "class": "defect", "story": id,
		"report": "design/analysis/2026-10-06-risk.md", "commit": true})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["outcome"] != "already recorded" || out["commit"] != nil || out["touches_added"] != nil {
		t.Errorf("an occurrence already recorded should commit nothing: %v", out)
	}
	if _, _, ahead := lastCommit(t, wt.Root); ahead != "1" {
		t.Errorf("no commit should be made: %s ahead", ahead)
	}

	out, failed = f.call(t, "issue_bump", map[string]any{"id": "I-1", "story": id, "note": "again", "commit": true})
	if failed != "" {
		t.Fatal(failed)
	}
	subject = "docs: [" + id + "] bump I-0001 Fixture was ignored"
	if out["outcome"] != "bumped" || out["count"] != float64(2) || field(out, "commit.subject") != subject || out["touches_added"] != nil {
		t.Errorf("issue_bump: %v", out)
	}
	if message, files, ahead := lastCommit(t, wt.Root); message != subject || !reflect.DeepEqual(files, []string{issue, summary}) || ahead != "2" {
		t.Errorf("the bump's commit: %q %v %s", message, files, ahead)
	}

	out, failed = f.call(t, "issue_close", map[string]any{"id": "I-0001", "reason": "fixed by " + id, "story": id, "commit": true})
	if failed != "" {
		t.Fatal(failed)
	}
	subject = "docs: [" + id + "] close I-0001 Fixture was ignored"
	if out["outcome"] != "closed" || out["story"] != id || field(out, "commit.subject") != subject {
		t.Errorf("issue_close: %v", out)
	}
	if message, files, ahead := lastCommit(t, wt.Root); message != subject || !reflect.DeepEqual(files, []string{issue, summary}) || ahead != "3" {
		t.Errorf("the close's commit: %q %v %s", message, files, ahead)
	}
	if is, err := issues.Get(wt, "I-0001"); err != nil || is.Status != "closed" || !strings.Contains(is.Body, "fixed by "+id) {
		t.Errorf("the issue in the worktree should be closed with the reason: %v %+v", err, is)
	}
	if st := gitIn(t, wt.Root, "status", "--porcelain"); st != "M design/system/plan.md" {
		t.Errorf("the worktree should hold only its other change: %q", st)
	}
}

// S-0275: commit is refused, and nothing is written, when the story's
// worktree does not have the story's branch checked out.
func TestCommitIsRefusedOffTheStorysBranch(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	wt := storyGit(t, f)
	gitIn(t, wt.Root, "checkout", "-q", "-b", "elsewhere")
	for _, c := range []struct {
		tool string
		args map[string]any
	}{
		{"issue_new", map[string]any{"title": "Off the branch", "class": "defect"}},
		{"adr_new", map[string]any{"decision": "Off the branch", "body": adrBody}},
	} {
		c.args["story"], c.args["commit"] = f.story.ID, true
		if _, failed := f.call(t, c.tool, c.args); !strings.Contains(failed, "has elsewhere checked out, not "+storygit.Branch(f.story.ID)) || !strings.Contains(failed, "nothing was written") {
			t.Errorf("%s off the story's branch should be refused: %q", c.tool, failed)
		}
	}
	if st := gitIn(t, wt.Root, "status", "--porcelain"); st != "" {
		t.Errorf("a refusal should write nothing: %q", st)
	}
}

// S-0275: commit is refused before anything is written when no story
// resolves, and when the story has no worktree, by every tool that takes it.
func TestCommitIsRefusedWithoutAStoryOrItsWorktree(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	if _, err := issues.New(f.repo, issues.NewOptions{Title: "Kept as it was", Class: "efficiency", Now: t0}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadDir(filepath.Join(f.repo.Root, "design/issues"))
	if err != nil {
		t.Fatal(err)
	}
	calls := []struct {
		tool string
		args map[string]any
	}{
		{"issue_new", map[string]any{"title": "Not recorded", "class": "defect"}},
		{"issue_bump", map[string]any{"id": "I-0001"}},
		{"issue_close", map[string]any{"id": "I-0001", "reason": "fixed"}},
		{"adr_new", map[string]any{"decision": "Not recorded", "body": adrBody}},
	}
	for _, c := range calls {
		c.args["commit"] = true
		if _, failed := f.call(t, c.tool, c.args); !strings.Contains(failed, c.tool+" with commit commits on a story's branch, and no story was given or resolved") {
			t.Errorf("%s with commit and no story should be refused: %q", c.tool, failed)
		}
		c.args["story"] = f.story.ID
		if _, failed := f.call(t, c.tool, c.args); !strings.Contains(failed, f.story.ID+" has none at .flai-cache/worktrees/"+f.story.ID+": open it with flai stream open") {
			t.Errorf("%s with commit and a story with no worktree should be refused: %q", c.tool, failed)
		}
	}
	after, _ := os.ReadDir(filepath.Join(f.repo.Root, "design/issues"))
	if len(after) != len(before) {
		t.Errorf("a refusal should write nothing: %v", after)
	}
	if is, err := issues.Get(f.repo, "I-0001"); err != nil || is.Count != 1 || is.Status != "open" {
		t.Errorf("the issue should be as it was: %v %+v", err, is)
	}
	if _, err := os.Stat(filepath.Join(f.repo.Root, "design/adrs")); err == nil {
		t.Error("a refused adr_new should write no ADR")
	}
}

// S-0275: without commit, issue_close closes the issue in the story's
// worktree when the story has one, regenerating its summary there, and in the
// project otherwise, as flai issue close does; nothing is committed.
func TestIssueCloseClosesInTheStorysWorktreeOrTheProject(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	wt := *f.repo
	wt.Root = f.repo.WorktreePath(f.story.ID)
	is, err := issues.New(f.repo, issues.NewOptions{Title: "Flaky test", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	// the worktree holds the issue as the story's branch had it from main
	data, err := os.ReadFile(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	writeIn(t, wt.Root, "design/issues/"+filepath.Base(is.Path), string(data))
	out, failed := f.call(t, "issue_close", map[string]any{"id": "I-1", "reason": "fixed", "story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	rel, _ := filepath.Rel(f.repo.Root, wt.Root)
	if out["outcome"] != "closed" || out["story"] != f.story.ID || out["path"] != filepath.ToSlash(rel)+"/design/issues/I-0001-flaky-test.md" || out["commit"] != nil {
		t.Errorf("issue_close in the worktree: %v", out)
	}
	if is, _ := issues.Get(&wt, "I-0001"); is == nil || is.Status != "closed" {
		t.Errorf("the worktree's issue should be closed: %+v", is)
	}
	if data, err := os.ReadFile(filepath.Join(wt.Root, "design/issues/summary.md")); err != nil || strings.Contains(string(data), "I-0001") {
		t.Errorf("the worktree's summary should leave the closed issue out: %v\n%s", err, data)
	}
	if is, _ := issues.Get(f.repo, "I-0001"); is == nil || is.Status != "open" {
		t.Errorf("the project's issue should be left open: %+v", is)
	}
	if _, failed := f.call(t, "issue_close", map[string]any{"id": "I-1", "reason": "fixed", "story": f.story.ID}); !strings.Contains(failed, "I-0001 is already closed") {
		t.Errorf("an issue already closed should be refused: %q", failed)
	}
	if _, failed := f.call(t, "issue_close", map[string]any{"id": "I-9", "reason": "fixed", "story": f.story.ID}); !strings.Contains(failed, "run flai stream sync "+f.story.ID+" first") {
		t.Errorf("an issue not in the worktree should be refused, saying to sync: %q", failed)
	}
	if _, failed := f.call(t, "issue_close", map[string]any{"id": "I-1", "reason": " "}); !strings.Contains(failed, "issue_close needs reason") {
		t.Errorf("a close without a reason should be refused: %q", failed)
	}

	out, failed = f.call(t, "issue_close", map[string]any{"id": "I-0001", "reason": "no longer applies"})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["story"] != nil || out["path"] != "design/issues/I-0001-flaky-test.md" {
		t.Errorf("issue_close without a story closes in the project: %v", out)
	}
	if is, _ := issues.Get(f.repo, "I-0001"); is == nil || is.Status != "closed" || !strings.Contains(is.Body, "no longer applies") {
		t.Errorf("the project's issue should be closed with the reason: %+v", is)
	}
}

func countStories(items []*workitem.Item) int {
	n := 0
	for _, it := range items {
		if it.Type == workitem.Story {
			n++
		}
	}
	return n
}
