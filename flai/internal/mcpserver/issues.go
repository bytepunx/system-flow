package mcpserver

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const issueStoryDescription = "Make a backlog story that remediates an open issue (S-0198), as flai issue story makes one: it takes the issue's title; its nature is remediation for a defect or a blocker, improvement otherwise; its goal links the issue and carries the issue's recommended solution. That link is what ties the issue to the story. The story is a draft (S-0203): the operator finalizes it before it can be ready. When the issue gives them, it carries the issue's cost of delay inputs, set by flai: time_lost_per_cycle from the issue's cost and count over the planning cycles since it was first reported, and the figures its Impact section gives, which take precedence; its Notes say how each was set. What strategic agents spent on the issue, its usage's strategic entries, is charged to the story and its epic (S-0227), and its Notes say so; the issue keeps its own. The issue's Remediation section then names the story. flai check runs with the story in place and refuses it, leaving nothing, if it reports anything the story introduces. A closed issue, or one an open story already links (named), is refused and nothing changes. epic is the story's epic (optional); story reads the issue from that story's worktree, where the issues it recorded are until it is accepted, and the story is named in the issue there, for that story to commit; the new story is made in wip/ as always. Nothing is committed."

// IssueStoryIn makes a story from an issue (S-0198).
type IssueStoryIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"the issue, such as I-0007 (any zero padding)"`
	Epic    string `json:"epic,omitempty" jsonschema:"the story's epic (optional: a story need not belong to one)"`
	Story   string `json:"story,omitempty" jsonschema:"read the issue from this story's worktree when it has one"`
}

func (in IssueStoryIn) project() string { return in.Project }

// IssueStoryOut is the story made and the issue it remediates.
type IssueStoryOut struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Nature string `json:"nature"`
	Path   string `json:"path"`
	Issue  string `json:"issue"`
	Draft  bool   `json:"draft"`
}

func (s *server) issueStory(ctx context.Context, _ *mcp.CallToolRequest, in IssueStoryIn) (*mcp.CallToolResult, IssueStoryOut, error) {
	cycle, err := s.repo.Manifest.Planning.CycleDuration()
	if err != nil {
		return nil, IssueStoryOut{}, fmt.Errorf("making a story for %s: %w", in.ID, err)
	}
	is, err := issues.Get(issues.RepoFor(s.repo, in.Story), in.ID)
	if err != nil {
		return nil, IssueStoryOut{}, err
	}
	if is.Status != "open" {
		return nil, IssueStoryOut{}, fmt.Errorf("%s is %s, so no story was made for it; reopen it by editing its status, or record a new issue", is.ID, is.Status)
	}
	items, err := s.listItems(ctx)
	if err != nil {
		return nil, IssueStoryOut{}, err
	}
	if id := issues.LinkedBy(is, items); id != "" {
		return nil, IssueStoryOut{}, fmt.Errorf("%s is already linked by open story %s, so no story was made for it; work it in %s, or cancel %s first", is.ID, id, id, id)
	}
	owner := s.repo.Manifest.Owner
	if owner == "" {
		owner = s.agent
	}
	now := s.now()
	draft := issues.ForStory(is, now, cycle)
	res, err := itemnew.Create(s.repo, s.runner, itemnew.Options{New: workitem.NewOptions{
		Type: workitem.Story, Title: draft.Title, Nature: draft.Nature, Parent: in.Epic,
		Owner: owner, Body: draft.Body, Now: now, Draft: draft.Draft, CostOfDelay: draft.CostOfDelay,
	}, Then: func(it *workitem.Item) ([]string, error) {
		if _, err := issues.CarryStrategic(s.repo, it.ID, draft.Strategic); err != nil {
			return nil, fmt.Errorf("carrying %s's strategic usage to it: %w", is.ID, err)
		}
		if err := issues.LinkStory(is, it.ID, now); err != nil {
			return nil, fmt.Errorf("naming it in %s's Remediation section: %w", is.ID, err)
		}
		return nil, nil
	}})
	if err != nil {
		return nil, IssueStoryOut{}, fmt.Errorf("making a story for %s: %w", is.ID, refusal(err))
	}
	return nil, IssueStoryOut{ID: res.Item.ID, Title: res.Item.Title, Nature: res.Item.Nature, Path: s.rel(res.Item.Path), Issue: is.ID, Draft: res.Item.Draft}, nil
}

// commitWhat says what commit does for the tools that write in a story's
// worktree (S-0275).
const commitWhat = " With commit, the write is one call (S-0275): the files it wrote, and only those, are committed on the story's branch story/S-nnnn in its worktree as docs: [S-nnnn] <subject>, with trailers (such as Co-Authored-By lines) ending the message, leaving everything else in the worktree as it was, and the story's touches are widened to the files committed, as flai touches --add would; commit gives the commit's hash, subject, and paths, and touches_added the paths added to the touches. commit is refused, and nothing is written, when no story resolves, when the story has no worktree (open it with flai stream open), when its worktree is not on the story's branch, or when a rebase is unfinished there. Without commit, nothing is committed."

// issueWhere says where issue_new and issue_bump record, and what they leave.
const issueWhere = " story is the story the occurrence belongs to, named in its instance; by default FLAI_STORY, else the story in this agent's name of the form agent-S-nnnn, else the story branch checked out where this server runs, else none. The issue is written in that story's worktree when it has one, where the issues a story records stay until it is accepted, for that story to commit, and in the project otherwise. An amount is a number of zero or more in planning.currency, time_lost_per_cycle a Go duration longer than zero such as 4h, and report the markdown file under design/analysis that found it, from the project root; anything else, like a class not in the list, is refused as flai refuses it and nothing is written. design/issues/summary.md is regenerated. Returns the issue's id, title, count, path (from the project root), the story recorded, and outcome." + commitWhat

const issueNewDescription = "Record an issue as flai issue new does (S-0224): a title, a class (defect, blocker, efficiency, or impression), the cost of this occurrence (a duration such as 20m), and a note, recorded as its first instance. The analyzer files each actionable finding of its report this way, with its impact: revenue_per_week, penalty_per_week, or time_lost_per_cycle, and the evidence for them, which become the issue's Impact section, carried over by issue_story as the cost of delay inputs of the story it makes; and report, the report that found it, which the instance names and the issue's Remediation section links, so that the report can link the issue this returns. With report, an open issue of the same title is bumped, with the report, impact, and note, rather than a second one opened (outcome bumped), and left as it is when an instance already names that report (outcome already recorded), so that filing a report's findings again counts nothing twice; otherwise outcome is opened. Without report a new issue is always opened. With commit the subject is record <I-nnnn> <title>, or bump <I-nnnn> <title> when it bumped, and an occurrence already recorded commits nothing and leaves summary.md as it was." + issueWhere

const issueBumpDescription = "Record another occurrence of an open issue as flai issue bump does (S-0224): its count, last reported, and average cost (cost is this occurrence's), and an instance with the note. revenue_per_week, penalty_per_week, and time_lost_per_cycle replace those figures in its Impact section, added when it has none, and keep the others; evidence adds the words behind them. report links the analysis report that found it, as issue_new does: the analyzer bumps this way a finding it judges the same as an open issue under another title, rather than filing a second. A closed issue is refused. outcome is bumped. With commit the subject is bump <I-nnnn> <title>." + issueWhere

const issueCloseDescription = "Close an open issue as flai issue close does (S-0275): its status becomes closed and reason (a story ID, a fix, or why it no longer applies) is written under its Remediation section; design/issues/summary.md is regenerated without it. story is the story closing it; by default FLAI_STORY, else the story in this agent's name of the form agent-S-nnnn, else the story branch checked out where this server runs, else none. The issue is closed in that story's worktree when it has one, for that story to commit, and in the project otherwise; an issue recorded on main since the story's branch was last synced is not in its worktree until flai stream sync brings it. An issue already closed, or one not found, is refused and nothing is written. Returns the issue's id, title, count, path (from the project root), the story whose worktree it was closed in, and outcome closed. With commit the subject is close <I-nnnn> <title>." + commitWhat

// IssueImpactIn is a finding's impact (S-0224): the cost of delay inputs an
// issue's Impact section gives, as item_edit takes them, and the evidence.
type IssueImpactIn struct {
	RevenuePerWeek   *float64 `json:"revenue_per_week,omitempty" jsonschema:"Impact: the revenue lost each week it stays open, an amount of zero or more in the project's currency"`
	PenaltyPerWeek   *float64 `json:"penalty_per_week,omitempty" jsonschema:"Impact: the penalty paid each week it stays open, an amount of zero or more in the project's currency"`
	TimeLostPerCycle string   `json:"time_lost_per_cycle,omitempty" jsonschema:"Impact: the time it loses each planning cycle, a Go duration longer than zero such as 4h"`
	Evidence         string   `json:"evidence,omitempty" jsonschema:"Impact: the evidence for the figures, in any words"`
	Report           string   `json:"report,omitempty" jsonschema:"the analysis report under design/analysis that found it, from the project root, such as design/analysis/2026-10-06-risk.md: the instance names it and the Remediation section links it"`
}

// impact is the issues.Impact the inputs give, amounts written as the CLI's
// flags take them.
func (in IssueImpactIn) impact() issues.Impact {
	text := func(v *float64) string {
		if s := amount(v); s != nil {
			return *s
		}
		return ""
	}
	return issues.Impact{RevenuePerWeek: text(in.RevenuePerWeek), PenaltyPerWeek: text(in.PenaltyPerWeek), TimeLostPerCycle: in.TimeLostPerCycle, Evidence: in.Evidence}
}

// IssueNewIn records an issue, or an analysis report's finding (S-0224).
type IssueNewIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Title   string `json:"title" jsonschema:"the issue's title; with report, an open issue of exactly this title is bumped instead of a second opened"`
	Class   string `json:"class" jsonschema:"defect, blocker, efficiency, or impression"`
	Cost    string `json:"cost,omitempty" jsonschema:"the wall-clock cost of this occurrence, a duration such as 20m"`
	Note    string `json:"note,omitempty" jsonschema:"what happened, recorded as the instance"`
	Story   string `json:"story,omitempty" jsonschema:"the story this occurrence belongs to (default: FLAI_STORY, else the story in the agent's name agent-S-nnnn, else the story branch checked out)"`
	IssueImpactIn
	CommitIn
}

func (in IssueNewIn) project() string { return in.Project }

// IssueBumpIn records another occurrence of an issue (S-0224).
type IssueBumpIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"the issue, such as I-0007 (any zero padding)"`
	Cost    string `json:"cost,omitempty" jsonschema:"the wall-clock cost of this occurrence, a duration such as 20m; the average is updated"`
	Note    string `json:"note,omitempty" jsonschema:"what happened this time"`
	Story   string `json:"story,omitempty" jsonschema:"the story this occurrence belongs to (default: FLAI_STORY, else the story in the agent's name agent-S-nnnn, else the story branch checked out)"`
	IssueImpactIn
	CommitIn
}

func (in IssueBumpIn) project() string { return in.Project }

// IssueCloseIn closes an issue (S-0275).
type IssueCloseIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"the issue, such as I-0007 (any zero padding)"`
	Reason  string `json:"reason" jsonschema:"what closed it: a story ID, a fix, or why it no longer applies"`
	Story   string `json:"story,omitempty" jsonschema:"the story closing it, in whose worktree it is closed when it has one (default: FLAI_STORY, else the story in the agent's name agent-S-nnnn, else the story branch checked out)"`
	CommitIn
}

func (in IssueCloseIn) project() string { return in.Project }

// CommitIn asks a tool that writes in a story's worktree to commit what it
// wrote on the story's branch and widen the story's touches (S-0275).
type CommitIn struct {
	Commit   bool     `json:"commit,omitempty" jsonschema:"commit the files written, and only those, on the story's branch in its worktree with the story's prefix, and widen the story's touches to them; refused, writing nothing, without a story whose worktree has its branch checked out"`
	Trailers []string `json:"trailers,omitempty" jsonschema:"with commit, lines that end the commit's message after a blank line, such as Co-Authored-By: Name <address>"`
}

// CommittedOut is what commit did: the commit, and what the story's touches
// gained.
type CommittedOut struct {
	Commit       *storygit.Commit `json:"commit,omitempty" jsonschema:"with commit, the commit made on the story's branch: its hash, subject, and paths (relative to the worktree, as git names them); absent when nothing was committed"`
	TouchesAdded []string         `json:"touches_added,omitempty" jsonschema:"with commit, the paths added to the story's touches; absent when its touches already covered every file committed"`
}

// IssueFiledOut is the issue issue_new or issue_bump recorded, and what it did.
type IssueFiledOut struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Count   int    `json:"count"`
	Path    string `json:"path" jsonschema:"from the project root; under the story's worktree when the issue was written there"`
	Story   string `json:"story,omitempty" jsonschema:"the story the occurrence was recorded for; for issue_close, the story in whose worktree it was closed"`
	Outcome string `json:"outcome" jsonschema:"opened, bumped, already recorded (with report, nothing changed), or closed"`
	CommittedOut
}

var (
	agentStoryRe  = regexp.MustCompile(`^agent-(S-\d+)$`)
	branchStoryRe = regexp.MustCompile(`^` + regexp.QuoteMeta(storygit.Prefix) + `(S-\d+)$`)
)

// recordingStory is the story an occurrence is recorded for, as flai issue
// new and bump resolve it: the one given, else FLAI_STORY, else the S-nnnn
// in this agent's name of the form agent-S-nnnn, else the story whose
// branch is checked out where the server runs. Empty is none.
func (s *server) recordingStory(given string) string {
	if v := strings.TrimSpace(given); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv("FLAI_STORY")); v != "" {
		return v
	}
	if m := agentStoryRe.FindStringSubmatch(strings.TrimSpace(s.agent)); m != nil {
		return m[1]
	}
	out, err := s.gitRunner().Run(s.repo.Root, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}
	if m := branchStoryRe.FindStringSubmatch(strings.TrimSpace(out)); m != nil {
		return m[1]
	}
	return ""
}

// gitRunner is the server's runner, execx.System when it has none.
func (s *server) gitRunner() execx.Runner {
	if s.runner == nil {
		return execx.System{}
	}
	return s.runner
}

// storyCheckout is the checkout a tool named what writes in for story, as
// issues.RepoFor finds it: the story's worktree when it has one, else the
// project. With commit it is refused, before anything is written, when no
// story resolves, when the story has no worktree, when its worktree does not
// have the story's branch checked out, and when a rebase is unfinished there.
func (s *server) storyCheckout(story string, commit bool, what string) (*workitem.Repo, error) {
	r := issues.RepoFor(s.repo, story)
	if !commit {
		return r, nil
	}
	if story == "" {
		return nil, fmt.Errorf("%s with commit commits on a story's branch, and no story was given or resolved: give story, or set FLAI_STORY, or leave commit out to write uncommitted in the project; nothing was written", what)
	}
	id := workitem.CanonicalID(story)
	if wt := s.repo.WorktreePath(id); r.Root != wt {
		return nil, fmt.Errorf("%s with commit commits on %s's branch in its worktree, and %s has none at %s: open it with flai stream open %s on the host, or leave commit out; nothing was written", what, id, id, s.rel(wt), id)
	}
	// With no paths, CommitPaths only checks that the worktree has the
	// story's branch checked out and no rebase unfinished, so that a refusal
	// comes before the write.
	if _, err := storygit.CommitPaths(storygit.CommitOptions{Runner: s.gitRunner(), Dir: r.Root, Story: id, Subject: what}); err != nil {
		return nil, fmt.Errorf("%s with commit: %w; nothing was written", what, err)
	}
	return r, nil
}

// commitWritten commits paths, written in the story's worktree r and
// relative to its root, on the story's branch as subject with trailers, and
// widens the story's touches in the project to the paths committed.
func (s *server) commitWritten(r *workitem.Repo, story string, paths []string, subject string, trailers []string) (CommittedOut, error) {
	id := workitem.CanonicalID(story)
	c, err := storygit.CommitPaths(storygit.CommitOptions{Runner: s.gitRunner(), Dir: r.Root, Story: id, Paths: paths, Subject: subject, Trailers: trailers})
	if err != nil {
		return CommittedOut{}, fmt.Errorf("%s was written in %s but not committed: %w", strings.Join(paths, ", "), s.rel(r.Root), err)
	}
	out := CommittedOut{Commit: c}
	if c == nil {
		return out, nil
	}
	w, err := itemedit.WidenStory(itemedit.WidenOptions{Repo: s.repo, Story: id, Paths: c.Paths, By: s.agent, Now: s.now})
	if err != nil {
		return out, fmt.Errorf("committed %s as %s on %s, but %s's touches were not widened: %w; add the paths with flai touches %s --add", strings.Join(c.Paths, ", "), c.Hash, storygit.Branch(id), id, err, id)
	}
	if w.Untold != nil && s.logger != nil {
		s.logger.Warn("overlap notices not sent", "component", "mcp", "agent", s.agent, "item", id, "err", w.Untold.Error())
	}
	if len(w.Story) > 0 {
		out.TouchesAdded = w.Story
	}
	return out, nil
}

// filed regenerates the summary in the checkout the issue is in, commits the
// issue and the summary as c asks, the subject beginning with verb, and says
// what was recorded. With commit, a call that wrote nothing, as for an
// occurrence already recorded, leaves the summary as it was and commits
// nothing.
func (s *server) filed(r *workitem.Repo, is *issues.Issue, story string, outcome issues.Outcome, c CommitIn, verb string) (IssueFiledOut, error) {
	written, err := is.Changed(r)
	if err != nil {
		return IssueFiledOut{}, err
	}
	if !c.Commit || len(written) > 0 {
		if _, err := issues.WriteSummary(r, s.now()); err != nil {
			return IssueFiledOut{}, fmt.Errorf("%s was recorded, but %s was not regenerated: %w", is.ID, issues.SummaryPath(r), err)
		}
	}
	if story != "" {
		story = workitem.CanonicalID(story)
	}
	out := IssueFiledOut{ID: is.ID, Title: is.Title, Count: is.Count, Path: s.rel(is.Path), Story: story, Outcome: string(outcome)}
	if c.Commit && len(written) > 0 {
		out.CommittedOut, err = s.commitWritten(r, story, append(written, issues.SummaryPath(r)), verb+" "+is.ID+" "+is.Title, c.Trailers)
	}
	return out, err
}

func (s *server) issueNew(_ context.Context, _ *mcp.CallToolRequest, in IssueNewIn) (*mcp.CallToolResult, IssueFiledOut, error) {
	story := s.recordingStory(in.Story)
	r, err := s.storyCheckout(story, in.Commit, "issue_new")
	if err != nil {
		return nil, IssueFiledOut{}, err
	}
	is, outcome, err := issues.NewOrBump(r, issues.NewOptions{Title: in.Title, Class: in.Class, Cost: in.Cost, Note: in.Note, Story: story,
		Impact: in.impact(), Report: in.Report, Now: s.now(), Runner: s.runner})
	if err != nil {
		return nil, IssueFiledOut{}, err
	}
	verb := "record"
	if outcome == issues.Bumped {
		verb = "bump"
	}
	out, err := s.filed(r, is, story, outcome, in.CommitIn, verb)
	return nil, out, err
}

func (s *server) issueBump(_ context.Context, _ *mcp.CallToolRequest, in IssueBumpIn) (*mcp.CallToolResult, IssueFiledOut, error) {
	story := s.recordingStory(in.Story)
	r, err := s.storyCheckout(story, in.Commit, "issue_bump")
	if err != nil {
		return nil, IssueFiledOut{}, err
	}
	is, err := issues.Get(r, in.ID)
	if err != nil {
		return nil, IssueFiledOut{}, err
	}
	if err := issues.BumpWith(r, is, issues.BumpOptions{Story: story, Cost: in.Cost, Note: in.Note, Impact: in.impact(), Report: in.Report, Now: s.now()}); err != nil {
		return nil, IssueFiledOut{}, err
	}
	out, err := s.filed(r, is, story, issues.Bumped, in.CommitIn, "bump")
	return nil, out, err
}

// issueClosed is issue_close's outcome.
const issueClosed issues.Outcome = "closed"

func (s *server) issueClose(_ context.Context, _ *mcp.CallToolRequest, in IssueCloseIn) (*mcp.CallToolResult, IssueFiledOut, error) {
	if strings.TrimSpace(in.Reason) == "" {
		return nil, IssueFiledOut{}, fmt.Errorf("issue_close needs reason, what closed %s: a story ID, a fix, or why it no longer applies; nothing was written", in.ID)
	}
	story := s.recordingStory(in.Story)
	r, err := s.storyCheckout(story, in.Commit, "issue_close")
	if err != nil {
		return nil, IssueFiledOut{}, err
	}
	if r == s.repo {
		// closed in the project, not in a story's worktree
		story = ""
	}
	is, err := issues.Get(r, in.ID)
	if err != nil {
		if story != "" {
			id := workitem.CanonicalID(story)
			return nil, IssueFiledOut{}, fmt.Errorf("%w in %s's worktree %s: if it was recorded on main since the story's branch was last synced, run flai stream sync %s first", err, id, s.rel(r.Root), id)
		}
		return nil, IssueFiledOut{}, err
	}
	if err := issues.Close(is, in.Reason, s.now()); err != nil {
		return nil, IssueFiledOut{}, err
	}
	out, err := s.filed(r, is, story, issueClosed, in.CommitIn, "close")
	return nil, out, err
}
