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
// human acceptance step and happens at story level.
var allowed = map[string][]string{
	Backlog:    {Ready, Cancelled},
	Ready:      {InProgress, Cancelled},
	InProgress: {Review, Cancelled},
	Review:     {Done, InProgress},
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
	ok := contains(allowed[from], to) || (it.Type == Task && from == InProgress && to == Done)
	if !ok {
		return nil, fmt.Errorf("rule: %s cannot go from %s to %s (allowed: %s)", it.ID, from, to, strings.Join(allowedFor(it, from), ", "))
	}
	if (to == Cancelled || (from == Review && to == InProgress)) && strings.TrimSpace(opt.Reason) == "" {
		return nil, fmt.Errorf("rule: moving %s to %s needs --reason", it.ID, to)
	}
	children := Children(opt.Items, it.ID)
	switch to {
	case Ready:
		// Tasks are not part of ready: the agent that pulls the story
		// writes them once it is in progress (ADR-0021).
		if it.Type == Story && !hasCriteria(it.Body) {
			return nil, fmt.Errorf("rule: a story needs an '## Acceptance criteria' section with at least one checkbox before it is ready")
		}
	case Review:
		if it.Type == Story && len(children) == 0 {
			return nil, fmt.Errorf("rule: a story needs at least one task before it goes to review (flai task new --story %s \"...\")", it.ID)
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
