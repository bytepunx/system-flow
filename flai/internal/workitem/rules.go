package workitem

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/verifyrecord"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// MoveOptions parameterise a transition.
type MoveOptions struct {
	By     string
	Reason string
	Now    time.Time
	// Items is the full item list used for parent and child rules and WIP
	// counts. Load it once with Repo.List(false).
	Items []*Item
	Board *Board
	// Cascade marks a cancellation that follows a parent's. It alone may
	// cancel an item in review (ADR-0028): cancelling the parent must close
	// everything under it, and the item's branch is left as it is.
	Cascade bool
	// Finalize lets a draft story go to ready, and makes it no longer a
	// draft (S-0199), recording By and Now as who finalized it (S-0201).
	Finalize bool
	// Kluster answers the kluster-verified rule on a transition that names
	// it (ADR-0312); without one the rule passes saying nothing.
	Kluster KlusterCheck
}

// KlusterVerdict is what a KlusterCheck found of a story's kluster namespace
// at its branch's head, in the words a move says of it.
type KlusterVerdict struct {
	// Applies is false where the rule has nothing to judge: no kluster:
	// block or no deploy in it, the kluster host action off, no agent
	// namespace, or no git. The rule then passes saying nothing (ADR-0312 §5).
	Applies bool
	// Passed says the namespace was deployed at the head and, while the
	// checks host action is on, its checks passed there.
	Passed bool
	// Refusal is what a refusal of the move says after "rule: ", and Warning
	// what a warning of it says, of a namespace that did not pass.
	Refusal, Warning string
}

// KlusterCheck answers the kluster-verified rule for a story. The namespace's
// record is the host's, not the repository's, so flai move and item_move pass
// one with the move; an error is a check that could not be made, which warns
// and lets the move through, whatever the rule's severity.
type KlusterCheck func(story string) (KlusterVerdict, error)

// Move validates and applies a state transition. Warnings are non-fatal
// policy breaches such as an exceeded WIP limit.
func (r *Repo) Move(it *Item, to string, opt MoveOptions) (warnings []string, err error) {
	wf := r.Flow()
	st, known := wf.State(to)
	if !known {
		return nil, fmt.Errorf("%q is not a state; use one of %s", to, strings.Join(wf.Names(), ", "))
	}
	// A move names the state's current name, never an alias (ADR-0301 §6).
	to = st.Name
	if it.Archived {
		return nil, fmt.Errorf("%s is archived and cannot be moved", it.ID)
	}
	from := it.Status
	if from == to {
		return nil, fmt.Errorf("%s is already %s", it.ID, to)
	}
	ok := follows(wf, it.Type, from, to)
	if !ok && opt.Cascade {
		t, found := wf.Find(it.Type, from, to)
		ok = found && t.Allows(workflow.ByCascade)
	}
	if !ok {
		return nil, fmt.Errorf("rule: %s cannot go from %s to %s (allowed: %s)", it.ID, from, to, strings.Join(allowedFor(wf, it, from), ", "))
	}
	// The rules the transition lists run in its order: the first refusal
	// stops the move, and warnings come in the order of the rules that gave
	// them (ADR-0301 §4).
	t, _ := wf.Find(it.Type, from, to)
	c := &moveCheck{r: r, wf: wf, t: t, it: it, to: to, opt: opt}
	for _, ref := range wf.MoveRules(it.Type, t) {
		if rule, _ := workflow.LookupRule(ref.Rule); !rule.Applies(it.Type) {
			continue
		}
		f := ruleChecks[ref.Rule](c)
		switch {
		case f == nil:
		case f.unchecked:
			warnings = append(warnings, f.warning)
		case ref.SeverityOf() == workflow.SeverityRefuse:
			return nil, fmt.Errorf("rule: %s", f.refusal)
		default:
			warnings = append(warnings, f.warning)
		}
	}
	if it.Type == Story && opt.Board != nil {
		if limit, ok := opt.Board.WIPLimits[to]; ok && limit > 0 {
			count := 1
			for _, x := range opt.Items {
				if x.Type == Story && x.Status == to && x.ID != it.ID {
					count++
				}
			}
			if count > limit {
				warnings = append(warnings, fmt.Sprintf("WIP limit for %s is %d, this makes %d", to, limit, count))
			}
		}
	}

	now := opt.Now.UTC().Format(TimeFormat)
	// A closed item is never blocked (I-0347): a block left open at done or
	// cancelled is reported by every check after it and counted to now by
	// flai stats, so the closing move ends it, and says so.
	if to == Done || to == Cancelled {
		for i := range it.Blocked {
			if it.Blocked[i].Until == "" {
				it.Blocked[i].Until = now
				warnings = append(warnings, fmt.Sprintf("%s was blocked (%s); the move to %s closed its blocked interval", it.ID, it.Blocked[i].Reason, to))
			}
		}
	}
	it.Transitions = append(it.Transitions, Transition{To: to, At: now, By: orDefault(opt.By, "agent")})
	it.Status = to
	it.Updated = now
	if to == Ready && opt.Finalize {
		it.Finalize(opt.By, opt.Now)
	}
	if opt.Reason != "" {
		it.Body = appendNote(it.Body, fmt.Sprintf("- %s: moved to %s: %s", now, to, opt.Reason))
	}
	if opt.Board != nil && it.Type == Story {
		if to == Ready {
			opt.Board.PlaceReadyLast(it.ID, opt.Items)
		} else {
			opt.Board.RemoveFromOrder(it.ID)
		}
	}
	return warnings, nil
}

// moveCheck is the move a rule's check reads: the transition, the item, the
// state it goes to, and the move's options, in the repository's workflow.
type moveCheck struct {
	r     *Repo
	wf    workflow.Workflow
	t     workflow.Transition
	it    *Item
	to    string
	opt   MoveOptions
	holds *Holds
}

// heldBy is why the story is held, or nil, reading the holds once a move.
func (c *moveCheck) heldBy() *Hold {
	if c.holds == nil {
		c.holds = c.r.Holds(c.opt.Items)
	}
	return c.holds.Of(c.it)
}

// finding is what a rule's check found against a move: refusal is what a
// refusal says after "rule: ", warning what a warning says.
type finding struct {
	refusal, warning string
	// unchecked marks a check that could not be made: it warns and lets the
	// move through, whatever the rule's severity.
	unchecked bool
}

// found is a finding that a refusal and a warning say alike.
func found(format string, args ...any) *finding {
	text := fmt.Sprintf(format, args...)
	return &finding{refusal: text, warning: text}
}

// ruleChecks is the check of each rule of workflow.Vocabulary, by its name:
// what the rule finds against a move, or nil when the move keeps it. Move
// runs only the rules that apply to the item's type.
var ruleChecks = map[string]func(c *moveCheck) *finding{
	// Tasks are not part of ready: the agent that pulls the story writes them
	// once it is in progress (ADR-0021).
	workflow.RuleCriteriaPresent: func(c *moveCheck) *finding {
		if hasCriteria(c.it.Body) {
			return nil
		}
		return found("a story needs an '## Acceptance criteria' section with at least one checkbox before it is ready")
	},
	workflow.RuleNotDraft: func(c *moveCheck) *finding {
		if !c.it.Draft || c.opt.Finalize {
			return nil
		}
		return found("%s is a draft: finalize it first (flai edit %s --no-draft, or flai move %s ready --yes)", c.it.ID, c.it.ID, c.it.ID)
	},
	workflow.RuleHasTask: func(c *moveCheck) *finding {
		if len(Children(c.opt.Items, c.it.ID)) > 0 {
			return nil
		}
		return found("a story needs at least one task before it goes to review (flai task new --story %s \"...\")", c.it.ID)
	},
	workflow.RuleNoOpenQuestion: func(c *moveCheck) *finding {
		n, err := ReadNarrative(c.r.NarrativePath(c.it.ID))
		if err != nil {
			return nil
		}
		qs := OpenQuestions(n.Body)
		if len(qs) == 0 {
			return nil
		}
		return found("%s has an open question in its narrative (%s): %q; it needs an answer before review (the operator answers it from the dashboard's inbox, or: flai stream answer %s %q \"<answer>\")", c.it.ID, n.Path, qs[0], c.it.ID, qs[0])
	},
	// Work left uncommitted in the worktree stops the acceptance later, when
	// its agent has gone (S-0140). A worktree git cannot read here (I-0017)
	// is not a reason to refuse: say so and move.
	workflow.RuleCleanWorktree: func(c *moveCheck) *finding {
		dirty, err := c.r.Uncommitted(c.it.ID)
		if err != nil {
			return &finding{warning: fmt.Sprintf("%s's worktree was not checked for uncommitted changes: %v", c.it.ID, err), unchecked: true}
		}
		if len(dirty) == 0 {
			return nil
		}
		return found("%s", c.r.UncommittedRule(c.it.ID, dirty, "it goes to review"))
	},
	// The operator reads the ticks as what the agent verified (S-0282), so
	// the baseline warns of an unticked box at review and refuses it at done.
	workflow.RuleCriteriaTicked: func(c *moveCheck) *finding {
		open := uncheckedCount(c.it.Body)
		if open == 0 {
			return nil
		}
		return &finding{
			refusal: fmt.Sprintf("%s has unchecked acceptance criteria", c.it.ID),
			warning: fmt.Sprintf("%s goes to review with unticked acceptance criteria (%d of %d): tick each one verified with flai criteria tick %s <n>, and say in its notes why any other is left unticked", c.it.ID, open, CriteriaCount(c.it.Body), c.it.ID),
		}
	},
	// The namespace's record is the host's, not the repository's: the move's
	// KlusterCheck reads it, and a cluster that is not there never holds a
	// story (ADR-0312 §5).
	workflow.RuleKlusterVerified: func(c *moveCheck) *finding {
		if c.opt.Kluster == nil {
			return nil
		}
		v, err := c.opt.Kluster(c.it.ID)
		switch {
		case err != nil:
			return &finding{warning: fmt.Sprintf("%s went to %s, but what its kluster namespace shows at its head was not read: %v", c.it.ID, c.to, err), unchecked: true}
		case !v.Applies || v.Passed:
			return nil
		}
		return &finding{refusal: v.Refusal, warning: v.Warning}
	},
	workflow.RuleChildrenClosed: func(c *moveCheck) *finding {
		for _, ch := range Children(c.opt.Items, c.it.ID) {
			if k := c.wf.Kind(ch.Status); k != workflow.KindAccepted && k != workflow.KindCancelled {
				return found("%s cannot be done while %s is %s", c.it.ID, ch.ID, ch.Status)
			}
		}
		return nil
	},
	workflow.RuleReasonGiven: func(c *moveCheck) *finding {
		if strings.TrimSpace(c.opt.Reason) != "" {
			return nil
		}
		return found("moving %s to %s needs --reason", c.it.ID, c.to)
	},
	// Nothing open lives under a cancelled parent (ADR-0028), so an item
	// comes back from cancelled only after its parent has.
	workflow.RuleParentOpen: func(c *moveCheck) *finding {
		if c.it.Parent == "" {
			return nil
		}
		for _, p := range c.opt.Items {
			if p.ID == c.it.Parent && c.wf.Kind(p.Status) == workflow.KindCancelled {
				return found("%s cannot go back to %s while its parent %s is cancelled; move %s back first", c.it.ID, c.to, p.ID, p.ID)
			}
		}
		return nil
	},
	workflow.RuleNotBlocked: func(c *moveCheck) *finding {
		if !c.it.IsBlocked() {
			return nil
		}
		return found("%s has an open blocked interval", c.it.ID)
	},
	// A hold stops flai serve and wait_for_work, not a move (S-0128,
	// ADR-0046).
	workflow.RuleNotHeld: func(c *moveCheck) *finding {
		h := c.heldBy()
		if h == nil {
			return nil
		}
		return found("%s is %s", c.it.ID, h.Reason)
	},
	// A story moved to ready while a story it names in after: is not done
	// waits there until that one is accepted; the move is the operator's, so
	// it warns rather than refuses (ADR-0279). An overlap alone does not
	// warn: a trial merge may clear it at the next look.
	workflow.RuleAfterDone: func(c *moveCheck) *finding {
		h := c.heldBy()
		if h == nil || h.Code != HoldAfter {
			return nil
		}
		return found("%s goes to ready %s; flai serve starts no agent for it until then", c.it.ID, h.Reason)
	},
	// A task's after: is the story's agent's plan, so it warns too (S-0176).
	workflow.RuleTaskAfter: func(c *moveCheck) *finding {
		p := PlanOf(c.opt.Items, c.it.Parent)
		if p == nil {
			return nil
		}
		why := p.Waits(c.it.ID)
		if why == "" {
			return nil
		}
		return found("%s is %s", c.it.ID, why)
	},
	// A move by a transition with verify needs the suites it lists passed at
	// the story's branch's head (ADR-0312 §2); a repository with no git to
	// read the head skips it, as clean-worktree's check does.
	workflow.RuleVerified: func(c *moveCheck) *finding {
		if c.r.Git == nil {
			return nil
		}
		return c.verified()
	},
}

// verified is what the verified rule finds against the move: nil when the
// verify record at the head of the story's worktree holds, for every suite
// the transition's verify needs, a verdict that passed there, all meaning
// every suite the record holds. A suite it needs that the record holds no
// entry of passes when the record's run passed, clean, at the head: flai
// verify records every suite whose tiers the story's paths select, so a
// passing run without the suite selected none of its tiers. A suite whose
// tiers each passed, were reused, or were left to the merge queue passes, as
// the close-out does (ADR-0174). No worktree, a
// head git cannot read, and no record at the head fail the rule.
func (c *moveCheck) verified() *finding {
	id := c.it.ID
	suites := workflow.Suites{All: true}
	if v := c.wf.VerifyOf(c.it.Type, c.t); v != nil {
		suites = v.Suites
	}
	fail := func(format string, args ...any) *finding {
		return found("moving %s to %s needs its last flai verify to have passed %s at its branch's head, but %s; run flai verify %s, then move it again",
			id, c.to, suitesText(suites), fmt.Sprintf(format, args...), id)
	}
	path := c.r.WorktreePath(id)
	if !exists(path) {
		return fail("it has no worktree at %s", c.r.RelToMain(path))
	}
	out, err := c.r.Git.Run(path, "git", "rev-parse", "HEAD")
	head := strings.TrimSpace(out)
	if err == nil && head == "" {
		err = errors.New("git named no commit")
	}
	if err != nil {
		return fail("its head was not read in %s: %v", c.r.RelToMain(path), err)
	}
	cache := c.r.CacheDir()
	verdicts, ok, err := verifyrecord.AtHead(cache, id, head, c.r.suiteOf)
	if err != nil {
		return fail("%v", err)
	}
	if !ok {
		return fail("flai verify recorded no run of it at its head %s", shortCommit(head))
	}
	rec, _, err := verifyrecord.Read(verifyrecord.CommitPath(cache, id, head))
	if err != nil {
		return fail("%v", err)
	}
	ran := rec.Passed && rec.Clean && rec.Commit == head
	need := suites.Names
	if suites.All {
		need = nil
		for _, v := range verdicts {
			need = append(need, v.Suite)
		}
		if len(need) == 0 && !ran {
			return fail("its run at %s did not pass", shortCommit(head))
		}
	}
	var short []string
	for _, name := range need {
		i := slices.IndexFunc(verdicts, func(v verifyrecord.Verdict) bool { return v.Suite == name })
		if i < 0 {
			if !ran {
				short = append(short, name+" (not in the record, whose run did not pass)")
			}
			continue
		}
		switch v := verdicts[i]; {
		case v.Passed:
		case v.AtHead && len(v.Held)+len(v.Queued) == len(v.Tiers):
			// a tier left to the merge queue runs in its batch, and the
			// close-out passes without it (ADR-0174); what the acceptance
			// needs of it is the acceptance's gates' (S-0808)
		case len(v.Held) == len(v.Tiers):
			short = append(short, name+" (its tiers passed, but not in a clean run at the head)")
		default:
			short = append(short, fmt.Sprintf("%s (%d of %d tiers passed)", name, len(v.Held), len(v.Tiers)))
		}
	}
	if len(short) == 0 {
		return nil
	}
	return fail("at %s it did not pass %s", shortCommit(head), strings.Join(short, ", "))
}

// suitesText names suites for a message: "every suite", or "suite static" and
// "suites static, behavior".
func suitesText(s workflow.Suites) string {
	switch {
	case s.All:
		return "every suite"
	case len(s.Names) == 1:
		return "suite " + s.Names[0]
	}
	return "suites " + strings.Join(s.Names, ", ")
}

// shortCommit is commit's first seven characters, for messages.
func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// suiteOf is the suite of the manifest's tier of that name, "" for a name no
// tier has, for a verify record written before suites were recorded.
func (r *Repo) suiteOf(tier string) string {
	for _, t := range r.Manifest.Tests {
		if strings.TrimSpace(t.Name) == tier {
			return t.SuiteOrDefault()
		}
	}
	return ""
}

// Follows reports whether an item of type typ may move from one state to the
// other in the baseline workflow, as Move allows it before its other rules.
func Follows(typ, from, to string) bool {
	return follows(baselineFlow, typ, from, to)
}

// follows reports whether wf has a transition for typ from one state to the
// other that is not made only as a parent's cancellation takes the item
// along: that move needs MoveOptions.Cascade (ADR-0028).
func follows(wf workflow.Workflow, typ, from, to string) bool {
	t, ok := wf.Find(typ, from, to)
	return ok && !cascadeOnly(t)
}

// cascadeOnly reports whether only a parent's cancellation makes t.
func cascadeOnly(t workflow.Transition) bool {
	return len(t.By) == 1 && t.By[0] == workflow.ByCascade
}

// allowedFor are the states an item of its type may move to from from in wf,
// in the file's order, leaving out the moves only a cascade makes.
func allowedFor(wf workflow.Workflow, it *Item, from string) []string {
	var out []string
	for _, t := range wf.From(it.Type, from) {
		if !cascadeOnly(t) && !slices.Contains(out, t.To) {
			out = append(out, t.To)
		}
	}
	if len(out) == 0 {
		return []string{"none, " + from + " is terminal"}
	}
	return out
}

// BlockItem opens a blocked interval.
func BlockItem(it *Item, reason string, now time.Time) error {
	if strings.TrimSpace(reason) == "" {
		return errors.New("a reason is required")
	}
	if it.IsBlocked() {
		return fmt.Errorf("%s is already blocked; unblock it first", it.ID)
	}
	if it.Closed() {
		return fmt.Errorf("%s is %s and cannot be blocked", it.ID, it.Status)
	}
	ts := now.UTC().Format(TimeFormat)
	it.Blocked = append(it.Blocked, Block{From: ts, Reason: reason})
	it.Updated = ts
	return nil
}

// UnblockItem closes the open interval: at now, or, on a done or cancelled
// item, at the item's close, since it was not blocked after it (I-0347).
func UnblockItem(it *Item, now time.Time) error {
	for i := range it.Blocked {
		if it.Blocked[i].Until == "" {
			it.Blocked[i].Until = unblockAt(it, it.Blocked[i], now)
			it.Updated = now.UTC().Format(TimeFormat)
			return nil
		}
	}
	return fmt.Errorf("%s is not blocked", it.ID)
}

// unblockAt is when UnblockItem closes the open interval b of it: now, or
// the close of a done or cancelled item, never before the interval began.
func unblockAt(it *Item, b Block, now time.Time) string {
	at := now.UTC()
	if closed := it.CompletedAt(); !closed.IsZero() {
		at = closed.UTC()
		if from, err := time.Parse(TimeFormat, b.From); err == nil && from.After(at) {
			at = from.UTC()
		}
	}
	return at.Format(TimeFormat)
}

var (
	criteriaHeading = regexp.MustCompile(`(?m)^## Acceptance criteria\s*$`)
	checkbox        = regexp.MustCompile(`(?m)^\s*- \[[ xX]\] \S`)
	unchecked       = regexp.MustCompile(`(?m)^\s*- \[ \] \S`)
)

func criteriaSection(body string) string {
	loc := criteriaHeading.FindStringIndex(body)
	if loc == nil {
		return ""
	}
	rest := body[loc[1]:]
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

func hasCriteria(body string) bool {
	return checkbox.MatchString(criteriaSection(body))
}

// uncheckedCount is the number of unticked checkbox lines in a body's
// "## Acceptance criteria" section; CriteriaCount counts them all.
func uncheckedCount(body string) int {
	return len(unchecked.FindAllString(criteriaSection(body), -1))
}

// appendNote adds a line under "## Notes", creating the section if missing.
func appendNote(body, line string) string {
	body = strings.TrimRight(body, "\n") + "\n"
	idx := strings.LastIndex(body, "\n## Notes")
	if idx < 0 && !strings.HasPrefix(body, "## Notes") {
		return body + "\n## Notes\n" + line + "\n"
	}
	return body + line + "\n"
}

// AppendChild adds "- ID Title" under the parent's Stories or Tasks section.
func AppendChild(parent *Item, child *Item) {
	heading := "## Tasks"
	if parent.Type == Epic {
		heading = "## Stories"
	}
	line := fmt.Sprintf("- %s %s", child.ID, child.Title)
	body := parent.Body
	idx := strings.Index(body, heading+"\n")
	if idx < 0 {
		parent.Body = strings.TrimRight(body, "\n") + "\n\n" + heading + "\n" + line + "\n"
		return
	}
	start := idx + len(heading) + 1
	rest := body[start:]
	end := strings.Index(rest, "\n## ")
	var section, tail string
	if end < 0 {
		section, tail = rest, ""
	} else {
		section, tail = rest[:end], rest[end:]
	}
	section = strings.TrimRight(section, "\n")
	if strings.Contains(section, child.ID+" ") {
		return
	}
	if section == "" {
		section = line
	} else {
		section += "\n" + line
	}
	parent.Body = body[:start] + section + "\n" + tail
	if tail == "" {
		parent.Body = strings.TrimRight(parent.Body, "\n") + "\n"
	} else if !strings.HasPrefix(tail, "\n\n") {
		parent.Body = body[:start] + section + "\n" + tail
	}
}
