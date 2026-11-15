package check

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The codes flai check reports a workflow under besides the parser's,
// workflow.file, workflow.rule, and workflow.timer (ADR-0301 §2, §3, §6).
const (
	codeState       = "workflow.state"
	codeUnreachable = "workflow.unreachable"
	codeDeadEnd     = "workflow.dead-end"
	codeNoPull      = "workflow.no-pull"
	codeNoAccept    = "workflow.no-accept"
	codeItem        = "workflow.item"
)

// workflowRel is the path, relative to the root, of the project's workflow
// file under the design folder, whether or not the project has it.
func (c *checker) workflowRel() string {
	return path.Join(c.repo.Manifest.Layout["design"], workflow.File)
}

// workflow validates the workflow the project runs (ADR-0301): the design
// folder's workflow file when the project has one, else c.repo.Flow(), the
// baseline flai carries or the one a test gives the repository, reported on
// the path the project's file would have, on line 1. A file that does not
// parse is one error under its problem's code, workflow.file, workflow.rule,
// or workflow.timer, and nothing else is checked on it; a workflow that
// parses is held to what flai needs to run it (flowCheck).
func (c *checker) workflow() {
	rel := c.workflowRel()
	abs := filepath.Join(c.repo.Root, filepath.FromSlash(rel))
	data, err := os.ReadFile(abs)
	switch {
	case errors.Is(err, os.ErrNotExist):
		flowCheck{c: c, w: c.repo.Flow(), path: abs}.run()
		return
	case err != nil:
		c.add(Error, workflow.CodeFile, abs, 1, "%s cannot be read: %v; make it readable, or remove it to run the baseline workflow flai carries", rel, err)
		return
	}
	w, prob := workflow.Parse(rel, string(data))
	if prob != nil {
		msg := prob.Message
		if prob.Field != "" {
			msg = prob.Field + ": " + msg
		}
		c.add(Error, prob.Code, filepath.Join(c.repo.Root, filepath.FromSlash(prob.Path)), prob.Line, "%s", msg)
		return
	}
	flowCheck{c: c, w: w, path: abs, file: true}.run()
}

// flowCheck holds one workflow that parses to what flai needs to run it,
// reporting on path: the workflow file's, on the lines of its entries, or
// with file false the path a project's file would have, on line 1.
type flowCheck struct {
	c    *checker
	w    workflow.Workflow
	path string
	file bool
}

// run reports every finding of the workflow: its states' counts and the
// accepted state left (counts), the states no item enters or leaves for good
// (reach), the pull and the acceptance (effect), and the rules and timers a
// workflow built in code may name without a parser to refuse them
// (vocabulary).
func (f flowCheck) run() {
	f.counts()
	f.reach()
	f.effect(effectRule{effect: workflow.EffectPull, code: codeNoPull, name: "pull", does: "start a story's work",
		from: workflow.KindQueue, to: workflow.KindActive, fromWords: "a queue state", toWords: "an active state"})
	f.effect(effectRule{effect: workflow.EffectAccept, code: codeNoAccept, name: "acceptance", does: "accept a story",
		from: workflow.KindWaiting, to: workflow.KindAccepted, fromWords: "a waiting state", toWords: "the accepted state", byAcceptance: true})
	f.vocabulary()
}

// line is an entry's line in the file, or 1 for a workflow with no file.
func (f flowCheck) line(n int) int {
	if f.file && n > 0 {
		return n
	}
	return 1
}

// keyLine is a front matter key's line in the file, or 1 for a workflow with
// no file.
func (f flowCheck) keyLine(key string) int {
	if f.file {
		return keyLine(f.path, key)
	}
	return 1
}

// add reports a finding on the workflow's path.
func (f flowCheck) add(level, code string, line int, format string, args ...any) {
	f.c.add(level, code, f.path, line, format, args...)
}

// move names the transition at i as the parser does: transitions[i] (from
// to to).
func move(i int, t workflow.Transition) string {
	return fmt.Sprintf("transitions[%d] (%s to %s)", i, strings.Join(t.From, ", "), t.To)
}

// kindOf is how a message names the state name: "a queue state", "an active
// state", or "no state of the workflow".
func (f flowCheck) kindOf(name string) string {
	switch k := f.w.Kind(name); k {
	case "":
		return "no state of the workflow"
	case workflow.KindActive, workflow.KindAccepted:
		return "an " + k + " state"
	default:
		return "a " + k + " state"
	}
}

// counts reports workflow.state, an error, for a workflow without exactly one
// initial state, one accepted state, and one cancelled state, and for each
// transition out of the accepted state (ADR-0301 §2).
func (f flowCheck) counts() {
	var initial []workflow.State
	for _, s := range f.w.States {
		if s.Initial {
			initial = append(initial, s)
		}
	}
	for _, want := range []struct {
		what   string
		states []workflow.State
		fix    string
	}{
		{"initial", initial, "mark exactly one queue state initial: true, the state items are created in"},
		{workflow.KindAccepted, f.w.OfKind(workflow.KindAccepted), "give exactly one state kind: accepted, where items close delivered"},
		{workflow.KindCancelled, f.w.OfKind(workflow.KindCancelled), "give exactly one state kind: cancelled, where items close undelivered"},
	} {
		switch len(want.states) {
		case 1:
		case 0:
			f.add(Error, codeState, f.keyLine("states"), "the workflow has no %s state; %s (ADR-0301 §2)", want.what, want.fix)
		default:
			f.add(Error, codeState, f.line(want.states[1].Line), "the workflow has %d %s states, %s; %s (ADR-0301 §2)", len(want.states), want.what, strings.Join(names(want.states), ", "), want.fix)
		}
	}
	for i, t := range f.w.Transitions {
		for _, from := range t.From {
			if f.w.Kind(from) == workflow.KindAccepted {
				f.add(Error, codeState, f.line(t.Line), "%s leaves %s, the accepted state, which is final: no transition leaves it (ADR-0301 §2); remove %s from the transition's from", move(i, t), from, from)
			}
		}
	}
}

// reach reports, for each state, workflow.unreachable, a warning, when no
// chain of transitions from an initial state enters it, and workflow.dead-end,
// an error, when no chain of transitions out of it enters an accepted or a
// cancelled state, which are themselves no dead end. Transitions of every
// item type count. A workflow with no initial state, or with neither an
// accepted nor a cancelled state, has its workflow.state error instead.
func (f flowCheck) reach() {
	reached := map[string]bool{}
	var initial []string
	for _, s := range f.w.States {
		if s.Initial {
			reached[s.Name] = true
			initial = append(initial, s.Name)
		}
	}
	queue := slices.Clone(initial)
	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]
		for _, t := range f.w.Transitions {
			if slices.Contains(t.From, from) && !reached[t.To] {
				reached[t.To] = true
				queue = append(queue, t.To)
			}
		}
	}
	closes := map[string]bool{}
	for _, s := range f.w.States {
		if s.Kind == workflow.KindAccepted || s.Kind == workflow.KindCancelled {
			closes[s.Name] = true
		}
	}
	for grew := true; grew; {
		grew = false
		for _, t := range f.w.Transitions {
			if !closes[t.To] {
				continue
			}
			for _, from := range t.From {
				if !closes[from] {
					closes[from], grew = true, true
				}
			}
		}
	}
	ends := strings.Join(append(names(f.w.OfKind(workflow.KindAccepted)), names(f.w.OfKind(workflow.KindCancelled))...), " or ")
	for _, s := range f.w.States {
		if len(initial) > 0 && !reached[s.Name] {
			f.add(Warning, codeUnreachable, f.line(s.Line), "state %s is entered by no chain of transitions from the initial state %s, so no item can reach it; add a transition into it, or remove it and list its name in the aliases of the state its items go to (ADR-0301 §2, §6)", s.Name, strings.Join(initial, ", "))
		}
		if ends != "" && !closes[s.Name] {
			f.add(Error, codeDeadEnd, f.line(s.Line), "state %s leads by no chain of transitions to %s, so an item in it can never close; add a transition out of it that leads there (ADR-0301 §2)", s.Name, ends)
		}
	}
}

// names are the states' names.
func names(states []workflow.State) []string {
	out := make([]string, len(states))
	for i, s := range states {
		out[i] = s.Name
	}
	return out
}

// effectRule is what ADR-0301 §3 asks of the one transition that gives
// stories an effect: its code, the effect's name in a message, what flai does
// with it, the kinds it leaves and enters with how a message says them, and
// whether only an acceptance makes it.
type effectRule struct {
	effect, code, name, does string
	from, to                 string
	fromWords, toWords       string
	byAcceptance             bool
}

// effect reports r.code, an error, unless exactly one transition gives
// stories r.effect, from a state of kind r.from to one of kind r.to, and with
// r.byAcceptance by [acceptance] alone: one finding for no such transition,
// one for each beyond the first, and one for each state of the wrong kind and
// for the wrong by.
func (f flowCheck) effect(r effectRule) {
	var found []int
	for i, t := range f.w.Transitions {
		if t.Takes(workflow.TypeStory) && t.HasEffect(r.effect) {
			found = append(found, i)
		}
	}
	if len(found) == 0 {
		f.add(Error, r.code, f.keyLine("transitions"), "no transition gives stories the %s effect, so flai cannot %s; give the move from %s to %s effects: [%s] (ADR-0301 §3)", r.effect, r.does, r.fromWords, r.toWords, r.effect)
		return
	}
	first := move(found[0], f.w.Transitions[found[0]])
	for _, i := range found {
		t := f.w.Transitions[i]
		if i != found[0] {
			f.add(Error, r.code, f.line(t.Line), "%s gives stories the %s effect, as %s does; exactly one transition does for stories: take %s out of the effects of all but one (ADR-0301 §3)", move(i, t), r.effect, first, r.effect)
		}
		for _, s := range t.From {
			if f.w.Kind(s) != r.from {
				f.add(Error, r.code, f.line(t.Line), "%s gives stories the %s effect from %s, %s; the %s leaves %s: give the effect to the move from %s to %s instead (ADR-0301 §3)", move(i, t), r.effect, s, f.kindOf(s), r.name, r.fromWords, r.fromWords, r.toWords)
			}
		}
		if f.w.Kind(t.To) != r.to {
			f.add(Error, r.code, f.line(t.Line), "%s gives stories the %s effect into %s, %s; the %s enters %s: give the effect to the move from %s to %s instead (ADR-0301 §3)", move(i, t), r.effect, t.To, f.kindOf(t.To), r.name, r.toWords, r.fromWords, r.toWords)
		}
		if r.byAcceptance && !slices.Equal(t.By, []string{workflow.ByAcceptance}) {
			f.add(Error, r.code, f.line(t.Line), "%s accepts stories by [%s]; only an acceptance accepts a story: give it by: [acceptance] (ADR-0301 §3)", move(i, t), strings.Join(t.By, ", "))
		}
	}
}

// vocabulary reports workflow.rule, an error, for each rule a transition
// names that flai does not implement, and workflow.timer, an error, for each
// state a timer names that the workflow lacks. The parser refuses both in a
// file, so they report only a workflow built in code, such as a test's.
func (f flowCheck) vocabulary() {
	for i, t := range f.w.Transitions {
		for _, r := range t.Rules {
			if _, ok := workflow.LookupRule(r.Rule); !ok {
				f.add(Error, workflow.CodeRule, f.line(t.Line), "%s names rule %s, which flai does not implement; name one of %s (ADR-0301 §4), or ask for a story that adds it", move(i, t), r.Rule, strings.Join(workflow.RuleNames(), ", "))
			}
		}
	}
	states := f.w.Names()
	for _, tm := range f.w.Timers {
		parts := map[string][]string{"sum": tm.Sum}
		for part, pt := range map[string]*workflow.Point{"from": tm.From, "to": tm.To, "until": tm.Until} {
			if pt != nil && pt.State != "" {
				parts[part] = []string{pt.State}
			}
		}
		for _, part := range []string{"from", "to", "sum", "until"} {
			for _, s := range parts[part] {
				if !slices.Contains(states, s) {
					f.add(Error, workflow.CodeTimer, f.line(tm.Line), "timer %s's %s names %q, which is no state of the workflow; name one of %s, or add the state (ADR-0301 §5)", tm.Name, part, s, strings.Join(states, ", "))
				}
			}
		}
	}
}

// itemStates reports workflow.item, an error, once for each name item's
// status or a transition of its history gives that the workflow the
// repository runs neither has as a state nor lists as an alias (ADR-0301
// §6): flai cannot tell which state the item is in, or was. It is on the
// status line when the status names it, else on the transitions line.
func (c *checker) itemStates(it *workitem.Item) {
	flow := c.repo.Flow()
	seen := map[string]bool{}
	report := func(name, key, where string) {
		if _, ok := flow.State(name); ok || name == "" || seen[name] {
			return
		}
		seen[name] = true
		c.add(Error, codeItem, it.Path, keyLine(it.Path, key), "%s's %s names state %s, which the workflow neither has nor lists as an alias; add the state to %s, or add %s to the aliases of the state its items go to (ADR-0301 §6)", it.ID, where, name, c.workflowRel(), name)
	}
	report(it.Status, "status", "status")
	for i, tr := range it.Transitions {
		report(tr.To, "transitions", fmt.Sprintf("transitions[%d]", i))
	}
}
