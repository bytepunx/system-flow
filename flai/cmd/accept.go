package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// acceptOptions are shared by `flai accept` and `flai move <story> done`:
// there is one way for a story to become done, and it is acceptance (S-0046).
type acceptOptions struct {
	by       string
	trailers []string
	dryRun   bool
	// verified and evidence are the orchestrator's (ADR-0093): the commit its
	// verifier passed, and the file holding its evidence, "-" for stdin.
	verified string
	evidence string
	stdin    func() io.Reader // the command's standard input, for --evidence -
}

func newAcceptCmd(a *app) *cobra.Command {
	var o acceptOptions
	c := &cobra.Command{
		Use:   "accept <id>",
		Short: "Accept an item: merge its branch, move to done, archive, commit",
		Long: `The operator's acceptance step as one command, per
design/conventions/work-management.md and git.md:

  0. rebase the story branch and fast-forward it into the main branch
  1. move the item to done (its rules apply: children closed, criteria checked);
     a story's epic follows it, to done when it was the epic's last open story
  2. flai archive for the item and its children and narrative, and for an
     epic that followed its story to done, the epic and its cancelled stories
  3. git commit the work item and archive; while another git process holds
     the index lock, git add and git commit are run again, for about nine
     seconds in all, and the lock is never removed
  4. tell every story in progress or in review whose touches cover a path
     the merge changed which paths those are, for its agent's MCP inbox

Acceptance computes no release, creates no tag, and pushes nothing (S-0087):
that is a deliberate step of its own, run when the operator chooses to
publish what has accumulated on main, not tied to any one item. See flai
release --pending.

flai move <story> done from review runs exactly this. An item that is already
done but was never archived (an older flai, a hand edit) is completed from
step 0 without a second transition. An item that is done and archived while
its archived file is not committed is one whose commit failed: the error
kept git's output and named this command. It is completed from step 3,
committing what the main checkout holds under the usual subject, and step 4
tells the open stories the paths the story's commits changed. A done,
archived item whose file is committed is refused as already done.
The orchestrator completes neither: that is the operator's. --dry-run
changes nothing.

A story whose branch changes a path Claude Code protects (a .claude folder,
.mcp.json, and the others of ADR-0106) is accepted by its operator only: the
story's owner or the project's owner, or anyone but the orchestrator when
neither is named. Anyone else is refused, before anything is merged, with the
files named; --dry-run lists them.

--by orchestrator is the orchestrator's acceptance (ADR-0093), refused unless
orchestration.permissions.accept_reviews is on; under FLAI_ROLE=orchestrate no
other acceptance is allowed. It is refused, before anything is merged, unless
--verified names the story branch's head, every acceptance criterion is
ticked, every file the branch changes is under the story's touches and on no
path Claude Code protects, no thread on the story or its tasks is open, and
--evidence, a Verdict: line and one item "- <n>: <files>" per criterion,
names a changed file for each of them.
The evidence is written, with the commit, under ### Accepted by the
orchestrator in the story's Notes. With --dry-run the evidence is optional.`,
		Example: `  flai accept S-031 --by alex
  flai accept S-031 --by alex --dry-run
  flai accept S-031 --by orchestrator --verified 4f1c2a9 --dry-run
  flai accept S-031 --by orchestrator --verified 4f1c2a9 --evidence evidence.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			res, err := a.acceptItem(repo, it, o)
			if err != nil {
				return err
			}
			return a.printAccept(res)
		},
	}
	addAcceptFlags(c, &o)
	c.Flags().BoolVar(&o.dryRun, "dry-run", false, "show what would block acceptance and change nothing")
	return c
}

// addAcceptFlags registers the flags acceptance shares with move.
func addAcceptFlags(c *cobra.Command, o *acceptOptions) {
	if c.Flags().Lookup("by") == nil {
		c.Flags().StringVar(&o.by, "by", "", "who accepted (default: config author)")
	}
	c.Flags().StringArrayVar(&o.trailers, "trailer", nil, "line appended to the commit message (repeatable)")
	c.Flags().StringVar(&o.verified, "verified", "", "with --by orchestrator: the commit its verifier passed, the story branch's head")
	c.Flags().StringVar(&o.evidence, "evidence", "", "with --by orchestrator: file holding its evidence, a Verdict: line and - <n>: <files> per criterion (- reads standard input)")
	o.stdin = c.InOrStdin
}

// acceptItem runs the acceptance flow for a story or an epic.
func (a *app) acceptItem(repo *workitem.Repo, it *workitem.Item, o acceptOptions) (*preview.Acceptance, error) {
	opts, err := a.acceptance(repo, it, o)
	if err != nil {
		return nil, err
	}
	res, err := preview.AcceptWith(a.runner, repo, it, orDefault(o.by, a.author()), a.now(), a.logger(), opts)
	if err != nil {
		return nil, err
	}
	if opts.Orchestrator && res.Resumed {
		left := "is done but was never archived"
		if res.CommitOnly {
			left = "is done and archived but its acceptance was never committed"
		}
		return nil, fmt.Errorf("rule: %s %s; the orchestrator accepts only a story in review, so the operator completes this acceptance with flai accept %s", it.ID, left, it.ID)
	}
	res.DryRun = o.dryRun
	if !o.dryRun {
		res.By = orDefault(o.by, a.author())
	}
	if dirty := res.Uncommitted; len(dirty) > 0 && !a.yes && !o.dryRun {
		return nil, fmt.Errorf("working tree has uncommitted changes outside wip (%s); commit or stash them so the acceptance commit holds only acceptance, or pass --yes to include them", strings.Join(dirty, ", "))
	}
	if len(res.Blockers) > 0 && !o.dryRun {
		return nil, fmt.Errorf("%s cannot be accepted yet: %s", it.ID, strings.Join(res.Blockers, "; "))
	}
	if o.dryRun {
		return res, nil
	}
	if res.CommitOnly {
		return a.finishCommit(repo, it, res, o)
	}
	useGit, hasBranch := a.inGitWorkTree(repo.MainRoot), res.Branch != ""

	// 0. bring the story branch into the main branch (ADR-0019)
	var changed []string // what the merge brought, for the stories still open
	if it.Type == workitem.Story {
		before := ""
		if hasBranch {
			before = a.headOf(repo.MainRoot)
		}
		merged, err := a.mergeStoryBranch(repo, it.ID)
		if err != nil {
			return nil, err
		}
		res.Merged = merged
		if merged {
			a.acceptStep(it, "merged", res.Branch+" rebased and fast-forwarded into the main branch")
			if before != "" {
				if changed, err = a.changedSince(repo.MainRoot, before); err != nil {
					a.logger().Warn("overlap notices not sent", "component", "accept", "item", it.ID, "err", err)
				}
			}
		}
	}

	// 1. done (epics: walk through review if needed)
	items, err := repo.List(false)
	if err != nil {
		return nil, err
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, err
	}
	from := it.Status
	if !res.Resumed {
		for _, st := range preview.StepsToDone(it.Status) {
			if _, err := repo.Move(it, st, workitem.MoveOptions{By: orDefault(o.by, a.author()), Now: a.now(), Items: items, Board: board}); err != nil {
				return nil, err
			}
		}
		// the orchestrator's evidence goes into the acceptance commit with
		// the story (ADR-0093)
		if opts.Orchestrator {
			it.Body = withNotesSection(it.Body, orchestratorNotes(res, a.now()))
		}
		if err := repo.Save(it); err != nil {
			return nil, err
		}
		if _, err := repo.RollUp(it); err != nil {
			return nil, err
		}
		if it.Type == workitem.Story {
			if err := board.Save(a.now().Format("2006-01-02")); err != nil {
				return nil, err
			}
		}
	}
	res.Status = it.Status
	a.acceptStep(it, "done", it.ID+" moved to done")
	// 1a. the story's epic follows it, to done with its last open story
	// (S-0200); the story's roll-up has already summed the epic's usage
	if it.Type == workitem.Story && !res.Resumed {
		if res.Epic, err = a.followAccepted(repo, it, from, orDefault(o.by, a.author())); err != nil {
			return nil, err
		}
	}
	// 2. archive
	items, _ = repo.List(false)
	ap, err := repo.PlanArchive(items, preview.ArchivedWith(items, it.ID, res.Epic))
	if err != nil {
		return nil, err
	}
	if err := repo.Archive(ap); err != nil {
		return nil, err
	}
	res.Archived = len(ap.Items)
	a.acceptStep(it, "archived", fmt.Sprintf("%d item(s) and the narrative archived", res.Archived))
	if err := a.refreshIndex(repo); err != nil {
		return nil, err
	}
	if !useGit {
		return res, nil
	}
	// 3. commit: the item, the archive, and nothing else. No release is
	// computed, no tag created, no push made (S-0087) — that is a publish, a
	// deliberate step of its own over everything accumulated on main, not
	// tied to any one item's acceptance. See flai release --pending.
	if err := a.commitAcceptance(repo, it, res, o.trailers); err != nil {
		return nil, err
	}
	// 4. tell the stories still open what changed under their claim (S-0132)
	a.tellAccepted(repo, it, res, orDefault(o.by, a.author()), changed)
	return res, nil
}

// finishCommit completes the acceptance of it, done and archived, whose
// commit failed (I-0100): steps 3 and 4 only, with no merge, no transition,
// and no archive, since those happened.
func (a *app) finishCommit(repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance, o acceptOptions) (*preview.Acceptance, error) {
	if err := a.commitAcceptance(repo, it, res, o.trailers); err != nil {
		return nil, err
	}
	// The merge is behind this run, so what it changed cannot be read from
	// the head before it as step 0 reads it. The paths are taken from the
	// story's commits instead, outside the wip folder: what its branch
	// brought into the main branch.
	var changed []string
	if it.Type == workitem.Story {
		committed, err := storygit.Committed(a.runner, repo)
		if err != nil {
			a.logger().Warn("overlap notices not sent", "component", "accept", "item", it.ID, "err", err)
		}
		changed = committed[it.ID]
	}
	a.tellAccepted(repo, it, res, orDefault(o.by, a.author()), changed)
	return res, nil
}

// commitAcceptance is step 3: it commits everything in the main checkout
// under the acceptance's subject. git add and git commit are run again while
// another git process holds the index lock (I-0100). A commit that fails even
// so says what the acceptance has done and how to finish it.
func (a *app) commitAcceptance(repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance, trailers []string) error {
	if _, err := storygit.RunPastIndexLock(a.runner, repo.Root, "add", "-A"); err != nil {
		return notCommitted(it.ID, res, "not yet staged", err)
	}
	msg := acceptSubject(it.ID, res.Epic)
	for _, t := range trailers {
		msg += "\n\n" + t
	}
	if _, err := storygit.RunPastIndexLock(a.runner, repo.Root, "commit", "-q", "-m", msg); err != nil {
		return notCommitted(it.ID, res, "staged", err)
	}
	a.acceptStep(it, "committed", firstLine(msg))
	return nil
}

// acceptSubject is the subject of the acceptance commit of id, naming the
// epic that followed it to done.
func acceptSubject(id string, epic *workitem.Followed) string {
	subject := fmt.Sprintf("chore: [%s] accept and archive", id)
	if epic != nil && epic.To == workitem.Done {
		subject += ", with " + epic.ID
	}
	return subject
}

// notCommitted is the error of an acceptance of id whose commit failed with
// err, git's whole output in it: what is done, where its changes are, and
// that flai accept finishes it.
func notCommitted(id string, res *preview.Acceptance, changes string, err error) error {
	state := id + " is done and archived"
	if res.Merged {
		state = res.Branch + " is merged and " + state
	}
	return fmt.Errorf("commit the acceptance of %s: %w\n%s, and its changes are %s in the main checkout but not committed; once git's error above is resolved, run flai accept %s to commit them", id, err, state, changes, id)
}

// tellAccepted is step 4: it tells the stories still open which of changed
// their claims cover (S-0132).
func (a *app) tellAccepted(repo *workitem.Repo, it *workitem.Item, res *preview.Acceptance, by string, changed []string) {
	res.Overlaps = a.tellOverlaps(repo, it, by, changed)
	if len(res.Overlaps) > 0 {
		ids := make([]string, len(res.Overlaps))
		for i, n := range res.Overlaps {
			ids[i] = n.ID
		}
		a.acceptStep(it, "told", "told "+strings.Join(ids, ", ")+" which changed paths their claims cover")
	}
}

// acceptance is what o adds to the acceptance preview of it, after refusing
// what o may not ask: under FLAI_ROLE=orchestrate an acceptance as anyone but
// the orchestrator, --verified or --evidence for anyone else, the
// orchestrator's acceptance while orchestration.permissions.accept_reviews is
// off, and its real run without evidence or with evidence that cannot be read
// (ADR-0093).
func (a *app) acceptance(repo *workitem.Repo, it *workitem.Item, o acceptOptions) (preview.AcceptOptions, error) {
	orchestrator := workitem.IsOrchestrator(o.by)
	switch {
	case os.Getenv("FLAI_ROLE") == guard.RoleOrchestrate && !orchestrator:
		given := "no --by"
		if o.by != "" {
			given = "--by " + o.by
		}
		return preview.AcceptOptions{}, fmt.Errorf("rule: the orchestrator accepts a story only as itself, so give --by %s, not %s (ADR-0093)", workitem.ActivityOrchestrator, given)
	case !orchestrator && (o.verified != "" || o.evidence != ""):
		return preview.AcceptOptions{}, fmt.Errorf("rule: --verified and --evidence are the orchestrator's acceptance: give them with --by %s, or leave them out", workitem.ActivityOrchestrator)
	case !orchestrator:
		return preview.AcceptOptions{}, nil
	}
	if err := repo.OrchestratorPermits(manifest.PermitAcceptReviews, "accepts a story", it.ID); err != nil {
		return preview.AcceptOptions{}, fmt.Errorf("rule: %w", err)
	}
	opts := preview.AcceptOptions{Orchestrator: true, Verified: o.verified}
	if o.evidence == "" {
		if o.dryRun {
			return opts, nil
		}
		return preview.AcceptOptions{}, fmt.Errorf("rule: the orchestrator's acceptance of %s needs its evidence: --evidence <file> or -, a Verdict: line and one - <n>: <files> item per criterion (ADR-0093)", it.ID)
	}
	text, err := a.readEvidence(o)
	if err != nil {
		return preview.AcceptOptions{}, err
	}
	if opts.Evidence, err = preview.ParseEvidence(text); err != nil {
		return preview.AcceptOptions{}, fmt.Errorf("rule: the orchestrator's evidence for %s cannot be read as ADR-0093 asks: %w", it.ID, err)
	}
	return opts, nil
}

// readEvidence is the text of the orchestrator's evidence: --evidence's file,
// or standard input for "-". A heading in it is refused, since the text goes
// under a heading of its own in the story's Notes.
func (a *app) readEvidence(o acceptOptions) (string, error) {
	var data []byte
	var err error
	if o.evidence == "-" {
		data, err = io.ReadAll(o.stdin())
	} else {
		path := o.evidence
		if !filepath.IsAbs(path) && a.cwd != "" {
			path = filepath.Join(a.cwd, path)
		}
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return "", fmt.Errorf("read the orchestrator's evidence from %s: %w; give a readable file, or - for standard input", o.evidence, err)
	}
	text := strings.TrimSpace(string(data))
	for i, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			return "", fmt.Errorf("rule: line %d of the orchestrator's evidence is a heading (%s): write a Verdict: line and - <n>: <files> items without headings, since flai puts the evidence under ### Accepted by the orchestrator in the story's Notes", i+1, line)
		}
	}
	return text, nil
}

// orchestratorNotes is what the orchestrator's acceptance writes under the
// story's Notes (ADR-0093): the commit its verifier passed, when, and its
// evidence as given.
func orchestratorNotes(res *preview.Acceptance, now time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### Accepted by the orchestrator\n\n- Verified: %s\n- At: %s\n", res.Verified, now.UTC().Format(time.RFC3339))
	if res.Evidence != nil {
		fmt.Fprintf(&b, "\n%s\n", res.Evidence.Text)
	}
	return b.String()
}

// withNotesSection is body with section appended to the end of its "## Notes"
// section, which is added at the end of body when it has none.
func withNotesSection(body, section string) string {
	const notes = "## Notes"
	start := -1
	if strings.HasPrefix(body, notes+"\n") {
		start = 0
	} else if i := strings.Index(body, "\n"+notes+"\n"); i >= 0 {
		start = i + 1
	}
	if start < 0 {
		return strings.TrimRight(body, "\n") + "\n\n" + notes + "\n\n" + section
	}
	end := len(body)
	if i := strings.Index(body[start+len(notes):], "\n## "); i >= 0 {
		end = start + len(notes) + i + 1
	}
	head, tail := strings.TrimRight(body[:end], "\n"), body[end:]
	if tail != "" {
		return head + "\n\n" + section + "\n" + tail
	}
	return head + "\n\n" + section
}

// followAccepted moves and saves the epic of story, accepted from from, the
// way its acceptance takes it, and says where it went; nil when it stays.
func (a *app) followAccepted(repo *workitem.Repo, story *workitem.Item, from, by string) (*workitem.Followed, error) {
	all, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	epic, followed, err := repo.Follow(all, story, from, by, a.now(), true)
	if err != nil || followed == nil {
		return nil, err
	}
	if err := repo.Save(epic); err != nil {
		return nil, err
	}
	a.acceptStep(story, "epic", fmt.Sprintf("%s followed %s to %s", epic.ID, story.ID, epic.Status))
	return followed, nil
}

// acceptStep logs one completed step of an acceptance as an info event with
// stable fields, so a client such as the dashboard can show progress while
// the command runs (S-0041). The message is fixed; what varies is in fields.
func (a *app) acceptStep(it *workitem.Item, step, detail string) {
	a.logger().Info("acceptance step", "component", "accept", "item", it.ID, "step", step, "detail", detail)
}

func (a *app) printAccept(res *preview.Acceptance) error {
	if a.jsonOut {
		return a.printJSON(res)
	}
	if res.DryRun {
		for _, b := range res.Blockers {
			fmt.Fprintf(a.out, "blocked: %s\n", b)
		}
		if len(res.Uncommitted) > 0 {
			effect := "the real run refuses until they are committed or stashed, or --yes includes them"
			if a.yes {
				effect = "--yes includes them in the acceptance commit"
			}
			fmt.Fprintf(a.out, "uncommitted outside wip: %s; %s\n", strings.Join(res.Uncommitted, ", "), effect)
		}
		if res.OperatorOnly != "" {
			fmt.Fprintf(a.out, "%s: %s\n", res.OperatorOnly, strings.Join(res.Protected, ", "))
		}
		if res.Branch != "" {
			fmt.Fprintf(a.out, "would merge %s into the main branch and remove its worktree\n", res.Branch)
		}
		if res.CommitOnly {
			fmt.Fprintf(a.out, "would commit the acceptance of %s, done and archived but not committed: %s\n", res.ID, acceptSubject(res.ID, res.Epic))
		} else if e := res.Epic; e != nil {
			fmt.Fprintf(a.out, "would also move %s %s from %s to %s, following %s", e.ID, e.Title, e.From, e.To, e.Story)
			if e.To == workitem.Done {
				fmt.Fprint(a.out, ", and archive it")
			}
			fmt.Fprintln(a.out)
		}
		fmt.Fprintln(a.out, "dry run: nothing changed")
		return nil
	}
	verb := "accepted"
	if res.Resumed {
		verb = "completed the acceptance of"
	}
	if res.CommitOnly {
		fmt.Fprintf(a.out, "%s %s: done and archived already, committed", verb, res.ID)
	} else {
		fmt.Fprintf(a.out, "%s %s: done, %d items archived, committed", verb, res.ID, res.Archived)
	}
	if res.Merged {
		fmt.Fprintf(a.out, ", %s merged and removed", res.Branch)
	}
	fmt.Fprintln(a.out, "; nothing released yet, run flai release --pending to publish")
	if res.Evidence != nil {
		fmt.Fprintf(a.out, "accepted by the orchestrator at %s; its evidence is under ### Accepted by the orchestrator in the Notes of %s\n", short(res.Verified), res.ID)
	}
	if e := res.Epic; e != nil {
		fmt.Fprintf(a.out, "%s followed it from %s to %s", e.ID, e.From, e.To)
		if e.To == workitem.Done {
			fmt.Fprint(a.out, ", archived")
		}
		fmt.Fprintln(a.out)
	}
	for _, n := range res.Overlaps {
		fmt.Fprintf(a.out, "told %s it overlaps: %s\n", n.ID, strings.Join(n.Paths, ", "))
	}
	return nil
}

func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return s
}
