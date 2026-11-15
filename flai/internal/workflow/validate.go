package workflow

import (
	"bufio"
	"fmt"
	"slices"
	"strings"
)

// The codes a workflow that parses is judged under besides the parser's,
// workflow.rule and workflow.timer (ADR-0301 §2, §3, §6).
const (
	CodeState       = "workflow.state"
	CodeUnreachable = "workflow.unreachable"
	CodeDeadEnd     = "workflow.dead-end"
	CodeNoPull      = "workflow.no-pull"
	CodeNoAccept    = "workflow.no-accept"
)

// The levels of a Finding: an error keeps flai from running the workflow, a
// warning does not (ADR-0307 §2).
const (
	LevelError   = "error"
	LevelWarning = "warning"
)

// Finding is one way a workflow falls short of what flai needs to run it:
// its level, its code, the line of the entry it is on (0 when it is on Key,
// or on an entry built in code), the front matter key it is on when it is on
// no entry, and what is wrong and what to do.
type Finding struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Line    int    `json:"line"`
	Key     string `json:"key,omitempty"`
	Message string `json:"message"`
}

// LineIn is the finding's line in the workflow file content: its key's line
// when it is on a key, else its entry's, else 1.
func (f Finding) LineIn(content string) int {
	if f.Key != "" {
		return keyLine(content, f.Key)
	}
	if f.Line > 0 {
		return f.Line
	}
	return 1
}

// keyLine is the line of a top-level front matter key in content, or 1.
func keyLine(content, key string) int {
	sc := bufio.NewScanner(strings.NewReader(content))
	n := 0
	for sc.Scan() {
		n++
		if n > 1 && sc.Text() == "---" {
			break
		}
		if strings.HasPrefix(sc.Text(), key+":") {
			return n
		}
	}
	return 1
}

// errorsOf are the findings of level error.
func errorsOf(fs []Finding) []Finding {
	var out []Finding
	for _, f := range fs {
		if f.Level == LevelError {
			out = append(out, f)
		}
	}
	return out
}

// Findings judge the workflow by what flai needs to run it: its states'
// counts and the accepted state left (counts), the states no item enters or
// leaves for good (reach), the pull and the acceptance (effect), and the
// rules and timers a workflow built in code may name without a parser to
// refuse them (vocabulary).
func (w Workflow) Findings() []Finding {
	j := &judge{w: w}
	j.counts()
	j.reach()
	j.effect(effectRule{effect: EffectPull, code: CodeNoPull, name: "pull", does: "start a story's work",
		from: KindQueue, to: KindActive, fromWords: "a queue state", toWords: "an active state"})
	j.effect(effectRule{effect: EffectAccept, code: CodeNoAccept, name: "acceptance", does: "accept a story",
		from: KindWaiting, to: KindAccepted, fromWords: "a waiting state", toWords: "the accepted state", byAcceptance: true})
	j.vocabulary()
	return j.out
}

// judge gathers the findings of one workflow.
type judge struct {
	w   Workflow
	out []Finding
}

// add records a finding on the entry at line.
func (j *judge) add(level, code string, line int, format string, args ...any) {
	j.out = append(j.out, Finding{Level: level, Code: code, Line: line, Message: fmt.Sprintf(format, args...)})
}

// addKey records a finding on the front matter key.
func (j *judge) addKey(level, code, key, format string, args ...any) {
	j.out = append(j.out, Finding{Level: level, Code: code, Key: key, Message: fmt.Sprintf(format, args...)})
}

// move names the transition at i as the parser does: transitions[i] (from
// to to).
func move(i int, t Transition) string {
	return fmt.Sprintf("transitions[%d] (%s to %s)", i, strings.Join(t.From, ", "), t.To)
}

// kindOf is how a message names the state name: "a queue state", "an active
// state", or "no state of the workflow".
func (j *judge) kindOf(name string) string {
	switch k := j.w.Kind(name); k {
	case "":
		return "no state of the workflow"
	case KindActive, KindAccepted:
		return "an " + k + " state"
	default:
		return "a " + k + " state"
	}
}

// stateNames are the states' names.
func stateNames(states []State) []string {
	out := make([]string, len(states))
	for i, s := range states {
		out[i] = s.Name
	}
	return out
}

// counts finds workflow.state, an error, for a workflow without exactly one
// initial state, one accepted state, and one cancelled state, and for each
// transition out of the accepted state (ADR-0301 §2).
func (j *judge) counts() {
	var initial []State
	for _, s := range j.w.States {
		if s.Initial {
			initial = append(initial, s)
		}
	}
	for _, want := range []struct {
		what   string
		states []State
		fix    string
	}{
		{"initial", initial, "mark exactly one queue state initial: true, the state items are created in"},
		{KindAccepted, j.w.OfKind(KindAccepted), "give exactly one state kind: accepted, where items close delivered"},
		{KindCancelled, j.w.OfKind(KindCancelled), "give exactly one state kind: cancelled, where items close undelivered"},
	} {
		switch len(want.states) {
		case 1:
		case 0:
			j.addKey(LevelError, CodeState, "states", "the workflow has no %s state; %s (ADR-0301 §2)", want.what, want.fix)
		default:
			j.add(LevelError, CodeState, want.states[1].Line, "the workflow has %d %s states, %s; %s (ADR-0301 §2)", len(want.states), want.what, strings.Join(stateNames(want.states), ", "), want.fix)
		}
	}
	for i, t := range j.w.Transitions {
		for _, from := range t.From {
			if j.w.Kind(from) == KindAccepted {
				j.add(LevelError, CodeState, t.Line, "%s leaves %s, the accepted state, which is final: no transition leaves it (ADR-0301 §2); remove %s from the transition's from", move(i, t), from, from)
			}
		}
	}
}

// reach finds, for each state, workflow.unreachable, a warning, when no
// chain of transitions from an initial state enters it, and workflow.dead-end,
// an error, when no chain of transitions out of it enters an accepted or a
// cancelled state, which are themselves no dead end. Transitions of every
// item type count. A workflow with no initial state, or with neither an
// accepted nor a cancelled state, has its workflow.state error instead.
func (j *judge) reach() {
	reached := map[string]bool{}
	var initial []string
	for _, s := range j.w.States {
		if s.Initial {
			reached[s.Name] = true
			initial = append(initial, s.Name)
		}
	}
	queue := slices.Clone(initial)
	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]
		for _, t := range j.w.Transitions {
			if slices.Contains(t.From, from) && !reached[t.To] {
				reached[t.To] = true
				queue = append(queue, t.To)
			}
		}
	}
	closes := map[string]bool{}
	for _, s := range j.w.States {
		if s.Kind == KindAccepted || s.Kind == KindCancelled {
			closes[s.Name] = true
		}
	}
	for grew := true; grew; {
		grew = false
		for _, t := range j.w.Transitions {
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
	ends := strings.Join(append(stateNames(j.w.OfKind(KindAccepted)), stateNames(j.w.OfKind(KindCancelled))...), " or ")
	for _, s := range j.w.States {
		if len(initial) > 0 && !reached[s.Name] {
			j.add(LevelWarning, CodeUnreachable, s.Line, "state %s is entered by no chain of transitions from the initial state %s, so no item can reach it; add a transition into it, or remove it and list its name in the aliases of the state its items go to (ADR-0301 §2, §6)", s.Name, strings.Join(initial, ", "))
		}
		if ends != "" && !closes[s.Name] {
			j.add(LevelError, CodeDeadEnd, s.Line, "state %s leads by no chain of transitions to %s, so an item in it can never close; add a transition out of it that leads there (ADR-0301 §2)", s.Name, ends)
		}
	}
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

// effect finds r.code, an error, unless exactly one transition gives
// stories r.effect, from a state of kind r.from to one of kind r.to, and with
// r.byAcceptance by [acceptance] alone: one finding for no such transition,
// one for each beyond the first, and one for each state of the wrong kind and
// for the wrong by.
func (j *judge) effect(r effectRule) {
	var found []int
	for i, t := range j.w.Transitions {
		if t.Takes(TypeStory) && t.HasEffect(r.effect) {
			found = append(found, i)
		}
	}
	if len(found) == 0 {
		j.addKey(LevelError, r.code, "transitions", "no transition gives stories the %s effect, so flai cannot %s; give the move from %s to %s effects: [%s] (ADR-0301 §3)", r.effect, r.does, r.fromWords, r.toWords, r.effect)
		return
	}
	first := move(found[0], j.w.Transitions[found[0]])
	for _, i := range found {
		t := j.w.Transitions[i]
		if i != found[0] {
			j.add(LevelError, r.code, t.Line, "%s gives stories the %s effect, as %s does; exactly one transition does for stories: take %s out of the effects of all but one (ADR-0301 §3)", move(i, t), r.effect, first, r.effect)
		}
		for _, s := range t.From {
			if j.w.Kind(s) != r.from {
				j.add(LevelError, r.code, t.Line, "%s gives stories the %s effect from %s, %s; the %s leaves %s: give the effect to the move from %s to %s instead (ADR-0301 §3)", move(i, t), r.effect, s, j.kindOf(s), r.name, r.fromWords, r.fromWords, r.toWords)
			}
		}
		if j.w.Kind(t.To) != r.to {
			j.add(LevelError, r.code, t.Line, "%s gives stories the %s effect into %s, %s; the %s enters %s: give the effect to the move from %s to %s instead (ADR-0301 §3)", move(i, t), r.effect, t.To, j.kindOf(t.To), r.name, r.toWords, r.fromWords, r.toWords)
		}
		if r.byAcceptance && !slices.Equal(t.By, []string{ByAcceptance}) {
			j.add(LevelError, r.code, t.Line, "%s accepts stories by [%s]; only an acceptance accepts a story: give it by: [acceptance] (ADR-0301 §3)", move(i, t), strings.Join(t.By, ", "))
		}
	}
}

// vocabulary finds workflow.rule, an error, for each rule a transition
// names that flai does not implement, and workflow.timer, an error, for each
// state a timer names that the workflow lacks. The parser refuses both in a
// file, so they find only in a workflow built in code, such as a test's.
func (j *judge) vocabulary() {
	for i, t := range j.w.Transitions {
		for _, r := range t.Rules {
			if _, ok := LookupRule(r.Rule); !ok {
				j.add(LevelError, CodeRule, t.Line, "%s names rule %s, which flai does not implement; name one of %s (ADR-0301 §4), or ask for a story that adds it", move(i, t), r.Rule, strings.Join(RuleNames(), ", "))
			}
		}
	}
	states := j.w.Names()
	for _, tm := range j.w.Timers {
		parts := map[string][]string{"sum": tm.Sum}
		for part, pt := range map[string]*Point{"from": tm.From, "to": tm.To, "until": tm.Until} {
			if pt != nil && pt.State != "" {
				parts[part] = []string{pt.State}
			}
		}
		for _, part := range []string{"from", "to", "sum", "until"} {
			for _, s := range parts[part] {
				if !slices.Contains(states, s) {
					j.add(LevelError, CodeTimer, tm.Line, "timer %s's %s names %q, which is no state of the workflow; name one of %s, or add the state (ADR-0301 §5)", tm.Name, part, s, strings.Join(states, ", "))
				}
			}
		}
	}
}
