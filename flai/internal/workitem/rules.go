package workitem

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// allowed lists the transitions from each state (ADR-0004, workflow.md).
// Tasks may also go straight from in-progress to done: review is the
// human acceptance step and happens at story level. Every state but done
// has a way back a column (ADR-0055): done is acceptance, and final.
var allowed = map[string][]string{
	Backlog:    {Ready, Cancelled},
	Ready:      {InProgress, Backlog, Cancelled},
	InProgress: {Review, Ready, Cancelled},
	Review:     {Done, InProgress},
	Cancelled:  {Backlog},
}

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
}

// Move validates and applies a state transition. Warnings are non-fatal
// policy breaches such as an exceeded WIP limit.
func (r *Repo) Move(it *Item, to string, opt MoveOptions) (warnings []string, err error) {
	if !contains(States, to) {
		return nil, fmt.Errorf("%q is not a state; use one of %s", to, strings.Join(States, ", "))
	}
	if it.Archived {
		return nil, fmt.Errorf("%s is archived and cannot be moved", it.ID)
	}
	from := it.Status
	if from == to {
		return nil, fmt.Errorf("%s is already %s", it.ID, to)
	}
	ok := Follows(it.Type, from, to) || (opt.Cascade && to == Cancelled && from == Review)
	if !ok {
		return nil, fmt.Errorf("rule: %s cannot go from %s to %s (allowed: %s)", it.ID, from, to, strings.Join(allowedFor(it, from), ", "))
	}
	if (to == Cancelled || (from == Review && to == InProgress)) && strings.TrimSpace(opt.Reason) == "" {
		return nil, fmt.Errorf("rule: moving %s to %s needs --reason", it.ID, to)
	}
	// Nothing open lives under a cancelled parent (ADR-0028), so an item comes
	// back from cancelled only after its parent has.
	if from == Cancelled && it.Parent != "" {
		for _, p := range opt.Items {
			if p.ID == it.Parent && p.Status == Cancelled {
				return nil, fmt.Errorf("rule: %s cannot go back to %s while its parent %s is cancelled; move %s back first", it.ID, to, p.ID, p.ID)
			}
		}
	}
	children := Children(opt.Items, it.ID)
	switch to {
	case Ready:
		// Tasks are not part of ready: the agent that pulls the story
		// writes them once it is in progress (ADR-0021).
		if it.Type == Story && !hasCriteria(it.Body) {
			return nil, fmt.Errorf("rule: a story needs an '## Acceptance criteria' section with at least one checkbox before it is ready")
		}
		if it.Type == Story && it.Draft && !opt.Finalize {
			return nil, fmt.Errorf("rule: %s is a draft: finalize it first (flai edit %s --no-draft, or flai move %s ready --yes)", it.ID, it.ID, it.ID)
		}
	case Review:
		if it.Type == Story && len(children) == 0 {
			return nil, fmt.Errorf("rule: a story needs at least one task before it goes to review (flai task new --story %s \"...\")", it.ID)
		}
		if it.Type == Story {
			if n, err := ReadNarrative(r.NarrativePath(it.ID)); err == nil {
				if qs := OpenQuestions(n.Body); len(qs) > 0 {
					return nil, fmt.Errorf("rule: %s has an open question in its narrative (%s): %q; it needs an answer before review (the operator answers it from the dashboard's inbox, or: flai stream answer %s %q \"<answer>\")", it.ID, n.Path, qs[0], it.ID, qs[0])
				}
			}
		}
		// Work left uncommitted in the worktree stops the acceptance later,
		// when its agent has gone (S-0140). A worktree git cannot read here
		// (I-0017) is not a reason to refuse: say so and move.
		if it.Type == Story {
			dirty, err := r.Uncommitted(it.ID)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("%s's worktree was not checked for uncommitted changes: %v", it.ID, err))
			} else if len(dirty) > 0 {
				return nil, fmt.Errorf("rule: %s", r.UncommittedRule(it.ID, dirty, "it goes to review"))
			}
		}
		// The operator reads the ticks as what the agent verified (S-0282), so
		// an unticked box warns here and refuses only at done.
		if it.Type == Story {
			if open := uncheckedCount(it.Body); open > 0 {
				warnings = append(warnings, fmt.Sprintf("%s goes to review with unticked acceptance criteria (%d of %d): tick each one verified with flai criteria tick %s <n>, and say in its notes why any other is left unticked", it.ID, open, CriteriaCount(it.Body), it.ID))
			}
		}
	case Done:
		for _, c := range children {
			if !c.Closed() {
				return nil, fmt.Errorf("rule: %s cannot be done while %s is %s", it.ID, c.ID, c.Status)
			}
		}
		if it.Type == Story && hasUnchecked(it.Body) {
			return nil, fmt.Errorf("rule: %s has unchecked acceptance criteria", it.ID)
		}
	}
	if to == InProgress && it.IsBlocked() {
		warnings = append(warnings, fmt.Sprintf("%s has an open blocked interval", it.ID))
	}
	// A hold stops flai serve and wait_for_work, not a move (S-0128, ADR-0046).
	if to == InProgress && from == Ready && it.Type == Story {
		if h := r.Holds(opt.Items).Of(it); h != nil {
			warnings = append(warnings, fmt.Sprintf("%s is %s", it.ID, h.Reason))
		}
	}
	// A task's after: is the story's agent's plan, so it warns too (S-0176).
	if to == InProgress && it.Type == Task {
		if p := PlanOf(opt.Items, it.Parent); p != nil {
			if why := p.Waits(it.ID); why != "" {
				warnings = append(warnings, fmt.Sprintf("%s is %s", it.ID, why))
			}
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

// Follows reports whether an item of type typ may move from one state to the
// other, as Move allows it before its other rules.
func Follows(typ, from, to string) bool {
	return contains(allowed[from], to) || (typ == Task && from == InProgress && to == Done)
}

func allowedFor(it *Item, from string) []string {
	out := append([]string{}, allowed[from]...)
	if it.Type == Task && from == InProgress {
		out = append(out, Done)
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

// UnblockItem closes the open interval.
func UnblockItem(it *Item, now time.Time) error {
	for i := range it.Blocked {
		if it.Blocked[i].Until == "" {
			ts := now.UTC().Format(TimeFormat)
			it.Blocked[i].Until = ts
			it.Updated = ts
			return nil
		}
	}
	return fmt.Errorf("%s is not blocked", it.ID)
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

func hasUnchecked(body string) bool {
	return unchecked.MatchString(criteriaSection(body))
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
