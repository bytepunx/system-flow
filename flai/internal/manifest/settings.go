package manifest

import (
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// The strategic agents' settings that flai manifest set writes (S-0229,
// ADR-0039): what the orchestrator may do and how it orders and releases,
// and the planner's and the analyzer's agents and schedules; and the
// project's test tiers (S-0273). The catalog says of each its kind, its
// values, its default, and what it does; whether a value is allowed is
// always the manifest's own validation (Planning.Errors,
// Orchestration.Errors, Analysis.Errors, TestErrors) on the manifest as it
// would be, so flai manifest set and flai check cannot disagree.

// Kind is what a setting's value is, and so how it is written and edited.
type Kind string

// The kinds of setting.
const (
	// KindBoolean is true or false.
	KindBoolean Kind = "boolean"
	// KindChoice is one of the setting's Values.
	KindChoice Kind = "choice"
	// KindNumber is a number of zero or more; a count is a whole one.
	KindNumber Kind = "number"
	// KindDuration is a Go duration such as 168h or 90m.
	KindDuration Kind = "duration"
	// KindCron is a five-field cron expression in UTC, or daily.
	KindCron Kind = "cron"
	// KindAgent is an agent, the shape of the project's agent, given as JSON
	// such as {"model":"claude-sonnet-5","config":{"effort":"medium"}}.
	KindAgent Kind = "agent"
	// KindText is one line of text, such as an epic ID or a tag.
	KindText Kind = "text"
	// KindTests is a list of test tiers, given as JSON such as
	// [{"name":"test","command":["scripts/test.sh"],"paths":["**"]}].
	KindTests Kind = "tests"
)

// Setting is one key flai manifest set may write.
type Setting struct {
	// Key is the setting's dotted path in system-flow.yaml, such as
	// orchestration.release.policy.
	Key string `json:"key"`
	// Kind is what its value is.
	Kind Kind `json:"kind"`
	// Values are the values a choice or a boolean takes, in the order they
	// are listed to the operator.
	Values []string `json:"values,omitempty"`
	// Default is the value in effect when the key is unset, as it would be
	// written; empty when unset means none, or, for an agent, the project's
	// agent.
	Default string `json:"default,omitempty"`
	// Meaning is a sentence on what the setting does.
	Meaning string `json:"meaning"`
	// Risk is, for a permission of the orchestrator's, a sentence on what
	// can go wrong when it is on.
	Risk string `json:"risk,omitempty"`
	// Whole is set for a number that is a count, a whole number.
	Whole bool `json:"whole,omitempty"`

	// get is the setting's value in a manifest, and whether it is set.
	get func(Manifest) (any, bool)
}

// SettingValue is a setting with its value in a manifest.
type SettingValue struct {
	Setting
	// Value is the setting's value as the manifest has it: a bool, a string,
	// a float64, an int, an *Agent, or a []TestTier; nil when it is unset.
	Value any `json:"value,omitempty"`
	// Set reports whether the manifest sets the key to something other than
	// its zero: a permission written false reads as unset, as it acts.
	Set bool `json:"set"`
}

// withheld are keys of the strategic blocks that flai manifest set refuses
// for a reason of their own, with the reason.
var withheld = map[string]string{
	"planning.currency": "is not a setting flai writes: changing it re-denominates every amount on the items, which flai does not convert; edit it by hand in system-flow.yaml",
}

// permissionText is what each permission lets the orchestrator do, and its
// risk, by its name in PermissionNames.
var permissionText = map[string][2]string{
	PermitPlanBacklogEpics: {
		"Lets the orchestrator ask the planner to draft stories for an epic in the backlog.",
		"The planner runs, and spends, on an epic you may not mean to start yet, and its drafts fill the backlog.",
	},
	PermitFinalizeDrafts: {
		"Lets the orchestrator finalize a draft story, as flai edit --no-draft does.",
		"A story the planner drafted becomes one that can be promoted without your having read it.",
	},
	PermitPromoteToReady: {
		"Lets the orchestrator move a story to ready.",
		"Agents may pull, work, and spend on a story you have not chosen to start.",
	},
	PermitOrderReady: {
		"Lets the orchestrator write the order of the ready column by orchestration.policy.",
		"It replaces the order you set on the board, so the next story pulled is the policy's choice, not yours.",
	},
	PermitAnswerThreads: {
		"How the orchestrator answers threads: off leaves them to you, recommend replies with a recommendation for you, and autonomous answers them itself.",
		"With autonomous, agents act on answers you did not give; with recommend, the decision stays yours.",
	},
	PermitAcceptReviews: {
		"Lets the orchestrator accept a story in review, as flai accept does.",
		"Work is merged to the main branch without your review, and a story accepted wrongly has to be undone by hand.",
	},
	PermitPublish: {
		"Lets the orchestrator release and push accepted work, as flai release --pending and flai push do.",
		"A release reaches everyone who pulls or installs the project, and a published release cannot be taken back.",
	},
}

// catalog is every setting in the order it is listed to the operator.
var catalog = buildCatalog()

func buildCatalog() []Setting {
	var out []Setting
	for _, name := range PermissionNames {
		s := Setting{Key: "orchestration.permissions." + name, Kind: KindBoolean, Values: []string{"true", "false"}, Default: "false",
			Meaning: permissionText[name][0], Risk: permissionText[name][1]}
		s.get = func(m Manifest) (any, bool) { on := m.Orchestration.Permissions.Allows(name); return on, on }
		if name == PermitAnswerThreads {
			s.Kind, s.Values, s.Default = KindChoice, AnswerModes, AnswerOff
			s.get = func(m Manifest) (any, bool) { return text(m.Orchestration.Permissions.AnswerThreads) }
		}
		out = append(out, s)
	}
	rel := func(m Manifest) Release { return m.Orchestration.Release }
	return append(out,
		Setting{Key: "orchestration.policy", Kind: KindChoice, Values: OrderPolicies, Default: OrderFIFO,
			Meaning: "How the ready column is ordered: cod by cost of delay, wsjf by cost of delay over forecast duration, throughput shortest first, and fifo leaves your order alone.",
			get:     func(m Manifest) (any, bool) { return text(m.Orchestration.Policy) }},
		Setting{Key: "orchestration.release.policy", Kind: KindChoice, Values: ReleasePolicies, Default: ReleaseJudgement,
			Meaning: "When accepted stories not yet released are due a release: judgement leaves it to you, threshold at a value or a count of unreleased work, and theme when every story of an epic or a tag is accepted.",
			get:     func(m Manifest) (any, bool) { return text(rel(m).Policy) }},
		Setting{Key: "orchestration.release.value", Kind: KindNumber,
			Meaning: "Under threshold, the unreleased cost of delay per week, in the project's currency, at which a release is due.",
			get:     func(m Manifest) (any, bool) { return number(rel(m).Value) }},
		Setting{Key: "orchestration.release.count", Kind: KindNumber, Whole: true,
			Meaning: "Under threshold, the number of accepted stories not yet released at which a release is due.",
			get: func(m Manifest) (any, bool) {
				if c := rel(m).Count; c != nil {
					return *c, true
				}
				return nil, false
			}},
		Setting{Key: "orchestration.release.epic", Kind: KindText,
			Meaning: "Under theme, the epic, such as E-0001, whose stories, every one accepted, make a release due.",
			get:     func(m Manifest) (any, bool) { return text(rel(m).Epic) }},
		Setting{Key: "orchestration.release.tag", Kind: KindText,
			Meaning: "Under theme, the tag whose stories, every one accepted, make a release due.",
			get:     func(m Manifest) (any, bool) { return text(rel(m).Tag) }},
		Setting{Key: "orchestration.release.whole_epics", Kind: KindBoolean, Values: []string{"true", "false"}, Default: "false",
			Meaning: "Under every policy, holds a release back while a story accepted and not yet released belongs to an epic in neither review nor done.",
			get:     func(m Manifest) (any, bool) { on := rel(m).WholeEpics; return on, on }},
		Setting{Key: "planning.agent", Kind: KindAgent,
			Meaning: "The planner's agent, over the project's: what it sets wins, and what it leaves out is the project's agent's.",
			get:     func(m Manifest) (any, bool) { return agent(m.Planning.Agent) }},
		Setting{Key: "planning.replan", Kind: KindChoice, Values: []string{ReplanNever, ReplanDeterministic, ReplanAgent}, Default: ReplanDeterministic,
			Meaning: "What flai serve does when a story is accepted or cancelled or the pull order changes: never nothing, deterministic plays the board out again and moves forecast deliveries, and agent does that and queues the planner for each story whose delivery moved.",
			get:     func(m Manifest) (any, bool) { return text(m.Planning.Replan) }},
		Setting{Key: "planning.schedule", Kind: KindCron,
			Meaning: "When flai serve runs the planner over the ready column, a five-field cron expression in UTC or daily; unset, there is no schedule.",
			get:     func(m Manifest) (any, bool) { return text(m.Planning.Schedule) }},
		Setting{Key: "planning.hour_rate", Kind: KindNumber,
			Meaning: "What an hour of work costs, in the project's currency, which prices a cost of delay's time lost; unset means unknown, not free.",
			get:     func(m Manifest) (any, bool) { return number(m.Planning.HourRate) }},
		Setting{Key: "planning.cycle", Kind: KindDuration, Default: durationText(DefaultCycle),
			Meaning: "The period a cost of delay's time lost is counted over.",
			get:     func(m Manifest) (any, bool) { return text(m.Planning.Cycle) }},
		Setting{Key: "planning.default_duration", Kind: KindDuration, Default: durationText(DefaultDuration),
			Meaning: "The work a story is forecast to take when there is too little history to forecast it from.",
			get:     func(m Manifest) (any, bool) { return text(m.Planning.DefaultDuration) }},
		Setting{Key: "analysis.agent", Kind: KindAgent,
			Meaning: "The analyzer's agent, over the project's: what it sets wins, and what it leaves out is the project's agent's.",
			get:     func(m Manifest) (any, bool) { return agent(m.Analysis.Agent) }},
		Setting{Key: "analysis.schedule", Kind: KindCron,
			Meaning: "When flai serve runs the analyzer, a five-field cron expression in UTC or daily; unset, it runs only when you ask.",
			get:     func(m Manifest) (any, bool) { return text(m.Analysis.Schedule) }},
		Setting{Key: "tests", Kind: KindTests,
			Meaning: "The project's test tiers, cheapest first, that flai test runs for the paths each selects: each a name, a command, the paths that select it, and the format of its output; unset, one plain tier runs scripts/test.sh when it exists, and [] means none.",
			get: func(m Manifest) (any, bool) {
				if m.Tests == nil {
					return nil, false
				}
				return m.Tests, true
			}},
	)
}

func text(s string) (any, bool) {
	if s == "" {
		return nil, false
	}
	return s, true
}

func number(f *float64) (any, bool) {
	if f == nil {
		return nil, false
	}
	return *f, true
}

func agent(a *Agent) (any, bool) {
	if a.IsZero() {
		return nil, false
	}
	return a, true
}

// durationText is a duration as a manifest writes it: 168h, not 168h0m0s.
func durationText(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}

// Settings are the keys flai manifest set writes, in the order they are
// listed to the operator.
func Settings() []Setting {
	return slices.Clone(catalog)
}

// SettingFor is the setting key names, and whether flai manifest set writes
// it.
func SettingFor(key string) (Setting, bool) {
	i := slices.IndexFunc(catalog, func(s Setting) bool { return s.Key == key })
	if i < 0 {
		return Setting{}, false
	}
	return catalog[i], true
}

// SettingValues are the settings with their values in m, in catalog order.
func (m Manifest) SettingValues() []SettingValue {
	out := make([]SettingValue, len(catalog))
	for i, s := range catalog {
		v, set := s.get(m)
		if !set {
			v = nil
		}
		out[i] = SettingValue{Setting: s, Value: v, Set: set}
	}
	return out
}

// Assignment is a value to write to a setting, as given on the command
// line: a boolean, a choice, a number, a duration, or a cron expression as
// written, or an agent or a list of test tiers as JSON.
type Assignment struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Problem is a refused setting: the field and the reason.
type Problem struct {
	// Field is the dotted key the problem is with, such as
	// orchestration.release.value.
	Field string `json:"field"`
	// Reason says what is wrong and what to write instead.
	Reason string `json:"reason"`
}

// String is the problem as a refusal prints it: field: reason.
func (p Problem) String() string {
	if p.Field == "" {
		return p.Reason
	}
	return p.Field + ": " + p.Reason
}

// RefusedError is an edit of the settings refused, with each problem; the
// manifest is left as it was.
type RefusedError struct {
	Problems []Problem
}

// Error lists the problems, field: reason, separated by semicolons.
func (e *RefusedError) Error() string {
	parts := make([]string, len(e.Problems))
	for i, p := range e.Problems {
		parts[i] = p.String()
	}
	return strings.Join(parts, "; ")
}

// WriteSettings writes the assignments and removes the unset keys in the
// manifest at path, as EditSettings does; on a refusal it writes nothing.
func WriteSettings(path string, set []Assignment, unset []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	out, err := EditSettings(data, set, unset)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

// EditSettings is the manifest data with each assignment written into its
// block and each unset key removed, so that its default applies; every other
// line, the order of the keys, and the comments are kept. It refuses with a
// *RefusedError, naming each field and why, a key outside the catalog, a
// value not of the key's kind, and any problem the manifest as it would be
// has, which flai check would report.
func EditSettings(data []byte, set []Assignment, unset []string) ([]byte, error) {
	if len(set) == 0 && len(unset) == 0 {
		return nil, errors.New("nothing to change: give key=value to set or a key to unset, such as planning.replan=agent")
	}
	writes, problems := parseEdits(set, unset)
	if len(problems) > 0 {
		return nil, &RefusedError{Problems: problems}
	}
	lines := strings.Split(string(data), "\n")
	for _, w := range writes {
		var p *Problem
		if w.unset {
			lines, p = removeKey(lines, strings.Split(w.setting.Key, "."))
		} else {
			lines, p = setKey(lines, strings.Split(w.setting.Key, "."), w.render)
		}
		if p != nil {
			problems = append(problems, *p)
		}
	}
	if len(problems) > 0 {
		return nil, &RefusedError{Problems: problems}
	}
	out := []byte(strings.Join(lines, "\n"))
	var m Manifest
	if err := yaml.Unmarshal(out, &m); err != nil {
		return nil, fmt.Errorf("cannot read system-flow.yaml as it would be after the change: %w; fix its YAML by hand, then run this again", err)
	}
	if problems := m.settingsProblems(); len(problems) > 0 {
		return nil, &RefusedError{Problems: problems}
	}
	if err := readsBack(out, m, writes); err != nil {
		return nil, err
	}
	return out, nil
}

// write is one key to set or remove.
type write struct {
	setting Setting
	unset   bool
	// want is the value the key reads back as once written.
	want any
	// render is the key's lines at an indent, given the line it replaces.
	render func(indent, old string) []string
}

// parseEdits checks each key is in the catalog, given once, and each value
// of its key's kind.
func parseEdits(set []Assignment, unset []string) ([]write, []Problem) {
	var writes []write
	var problems []Problem
	seen := map[string]bool{}
	lookup := func(key string) (Setting, bool) {
		key = strings.TrimSpace(key)
		s, ok := SettingFor(key)
		switch {
		case !ok:
			reason, held := withheld[key]
			if !held {
				reason = "is not a setting flai manifest set writes; write one of " + strings.Join(settingKeys(), ", ")
			}
			problems = append(problems, Problem{Field: key, Reason: reason})
		case seen[key]:
			problems = append(problems, Problem{Field: key, Reason: "is given more than once; set or unset it once"})
			ok = false
		}
		seen[key] = true
		return s, ok
	}
	for _, a := range set {
		s, ok := lookup(a.Key)
		if !ok {
			continue
		}
		w, reason := s.parse(a.Value)
		if reason != "" {
			problems = append(problems, Problem{Field: s.Key, Reason: reason})
			continue
		}
		writes = append(writes, w)
	}
	for _, k := range unset {
		if s, ok := lookup(k); ok {
			writes = append(writes, write{setting: s, unset: true})
		}
	}
	return writes, problems
}

func settingKeys() []string {
	keys := make([]string, len(catalog))
	for i, s := range catalog {
		keys[i] = s.Key
	}
	return keys
}

// parse reads a value given for the setting, or says why it is not one of
// its kind. Whether the value is allowed is the manifest's validation's to
// say, after the write.
func (s Setting) parse(value string) (write, string) {
	v := strings.TrimSpace(value)
	w := write{setting: s}
	if v == "" {
		return w, "is empty; give a value, or unset it for its default"
	}
	if strings.ContainsAny(v, "\n\r\x00") {
		return w, "spans lines; a value is one line"
	}
	scalar := func(text string, want any) (write, string) {
		w.want = want
		w.render = func(indent, old string) []string { return []string{keyLine(indent, s.leaf(), text, old)} }
		return w, ""
	}
	switch s.Kind {
	case KindBoolean:
		if v != "true" && v != "false" {
			return w, fmt.Sprintf("%q is not true or false", v)
		}
		return scalar(v, v == "true")
	case KindNumber:
		if s.Whole {
			n, err := strconv.Atoi(v)
			if err != nil {
				return w, fmt.Sprintf("%q is not a whole number; write one such as 5", v)
			}
			return scalar(strconv.Itoa(n), n)
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
			return w, fmt.Sprintf("%q is not a number; write one such as 500 or 95.5", v)
		}
		return scalar(strconv.FormatFloat(f, 'f', -1, 64), f)
	case KindAgent:
		var a Agent
		if err := yaml.UnmarshalWithOptions([]byte(v), &a, yaml.Strict()); err != nil {
			return w, fmt.Sprintf("is not an agent: %s; write one as JSON, such as {\"harness\":\"claude-code\",\"model\":\"claude-sonnet-5\",\"config\":{\"effort\":\"medium\"}}", firstLine(err.Error()))
		}
		if a.IsZero() {
			w.unset = true
			return w, ""
		}
		w.want = &a
		w.render = func(indent, old string) []string {
			block := strings.Split(strings.TrimSuffix(AgentBlock(&a, indent), "\n"), "\n")
			if old != "" {
				if rest, ok := keyRest(old, s.leaf()); ok {
					if value, _ := splitComment(rest); value == "" {
						block[0] = old // keep the key's line and its comment
					}
				}
			}
			return block
		}
		return w, ""
	case KindTests:
		const example = `such as [{"name":"test","command":["scripts/test.sh"],"paths":["**"]}], or [] for none`
		if !strings.HasPrefix(v, "[") {
			return w, "is not a list of test tiers; write one as JSON, " + example
		}
		tiers := []TestTier{}
		if err := yaml.UnmarshalWithOptions([]byte(v), &tiers, yaml.Strict()); err != nil {
			return w, fmt.Sprintf("is not a list of test tiers: %s; write one as JSON, %s", firstLine(err.Error()), example)
		}
		if tiers == nil {
			tiers = []TestTier{}
		}
		w.want = tiers
		w.render = func(indent, old string) []string { return testsLines(tiers, indent, s.leaf(), old) }
		return w, ""
	}
	return scalar(settingScalar(v), v)
}

// testsLines are the lines of a list of test tiers under key at indent, one
// tier an item and each of its lists on one line, keeping the comment of the
// key's line they replace.
func testsLines(tiers []TestTier, indent, key, old string) []string {
	if len(tiers) == 0 {
		return []string{keyLine(indent, key, "[]", old)}
	}
	head := indent + key + ":"
	if rest, ok := keyRest(old, key); ok {
		_, comment := splitComment(rest)
		head = withComment(head, comment)
	}
	out := []string{head}
	for _, t := range tiers {
		prefix := indent + "  - "
		field := func(k, v string) {
			out = append(out, prefix+k+": "+v)
			prefix = indent + "    "
		}
		field("name", yamlScalar(t.Name))
		field("command", flowList(t.Command))
		if t.Dir != "" {
			field("dir", yamlScalar(t.Dir))
		}
		if len(t.Paths) > 0 {
			field("paths", flowList(t.Paths))
		}
		if t.Format != "" {
			field("format", yamlScalar(t.Format))
		}
		if t.AllOnly {
			field("all_only", "true")
		}
		if len(t.AllCommand) > 0 {
			field("all_command", flowList(t.AllCommand))
		}
	}
	return out
}

// flowList is a list of strings on one line, each plain or double-quoted as
// yamlScalar says.
func flowList(items []string) string {
	out := make([]string, len(items))
	for i, s := range items {
		out[i] = yamlScalar(s)
	}
	return "[" + strings.Join(out, ", ") + "]"
}

// sameTiers reports whether two lists of test tiers say the same, a list
// absent and one empty alike.
func sameTiers(a, b []TestTier) bool {
	return slices.EqualFunc(a, b, func(x, y TestTier) bool {
		return x.Name == y.Name && x.Dir == y.Dir && x.Format == y.Format && x.AllOnly == y.AllOnly &&
			slices.Equal(x.Command, y.Command) && slices.Equal(x.Paths, y.Paths) && slices.Equal(x.AllCommand, y.AllCommand)
	})
}

// leaf is the last part of the setting's key, as its block writes it.
func (s Setting) leaf() string {
	return s.Key[strings.LastIndex(s.Key, ".")+1:]
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	return line
}

// settingScalar is a string written plain when it reads back as the same
// string, such as 168h or daily, and double-quoted otherwise, as a word
// another YAML reader takes for a boolean, such as off, is.
func settingScalar(s string) string {
	q := yamlScalar(s)
	if !strings.HasPrefix(q, `"`) || plainValue.MatchString(s) {
		return q
	}
	var doc map[string]any
	if !strings.ContainsAny(s, `#:"'{}[],&*!|>%@`+"`") && yaml.Unmarshal([]byte("v: "+s), &doc) == nil {
		if got, ok := doc["v"].(string); ok && got == s {
			return s
		}
	}
	return strconv.Quote(s)
}

// keyLine is key: value at indent, keeping the comment of the line it
// replaces in the column it was in when there is room.
func keyLine(indent, key, value, old string) string {
	line := indent + key + ": " + value
	if old == "" {
		return line
	}
	rest, _ := keyRest(old, key)
	_, comment := splitComment(rest)
	if comment == "" {
		return line
	}
	old = strings.TrimRight(old, " \t")
	if col := len(old) - len(comment); col > len(line) {
		return line + strings.Repeat(" ", col-len(line)) + comment
	}
	return line + " " + comment
}

// settingsProblems are what flai check reports of the manifest, as fields
// and reasons.
func (m Manifest) settingsProblems() []Problem {
	var msgs []string
	msgs = append(msgs, m.Agent.Problems()...)
	if _, err := m.Issues.StoryAfterDuration(); err != nil {
		msgs = append(msgs, err.Error())
	}
	msgs = append(msgs, m.Planning.Errors()...)
	msgs = append(msgs, m.Orchestration.Errors()...)
	msgs = append(msgs, m.Analysis.Errors()...)
	msgs = append(msgs, m.TestErrors()...)
	for _, e := range m.Claims.Errors() {
		msgs = append(msgs, e.Error())
	}
	out := make([]Problem, len(msgs))
	for i, msg := range msgs {
		out[i] = splitProblem(msg)
	}
	return out
}

// splitProblem is a sentence of the manifest's validation as a field and a
// reason: each begins with the key it is about, such as planning.cycle "0s"
// is not ..., or tests[1].format "junit" is not ....
func splitProblem(msg string) Problem {
	field, reason, ok := strings.Cut(msg, " ")
	if !ok || strings.Trim(field, "abcdefghijklmnopqrstuvwxyz_.[]0123456789") != "" {
		return Problem{Reason: msg}
	}
	return Problem{Field: field, Reason: reason}
}

// readsBack checks each key written reads back as what was meant, so that a
// manifest laid out in a way the edit does not understand is refused rather
// than written wrong.
func readsBack(data []byte, m Manifest, writes []write) error {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot read system-flow.yaml as it would be after the change: %w; fix its YAML by hand, then run this again", err)
	}
	for _, w := range writes {
		_, present := lookupPath(raw, strings.Split(w.setting.Key, "."))
		got, _ := w.setting.get(m)
		ok := present != w.unset
		switch want := w.want.(type) {
		case nil:
		case *Agent:
			a, _ := got.(*Agent)
			ok = ok && a.Same(want)
		case []TestTier:
			tiers, _ := got.([]TestTier)
			ok = ok && sameTiers(tiers, want)
		default:
			ok = ok && got == want
		}
		if !ok {
			return fmt.Errorf("cannot write %s so that system-flow.yaml reads it back as given; the file is left as it was: edit the key by hand", w.setting.Key)
		}
	}
	return nil
}

// lookupPath is the value at the dotted path in a YAML document read
// generically, and whether it is there.
func lookupPath(doc map[string]any, path []string) (any, bool) {
	var cur any = doc
	for _, k := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		if cur, ok = m[k]; !ok {
			return nil, false
		}
	}
	return cur, true
}

// span is where a key is among a manifest's lines: its key line, and the
// line after the last line of its value.
type span struct {
	at, end int
}

// findKey is where key is among the keys of the mapping in lines[from:to],
// its at -1 when it is not there, and the indent of the mapping's keys: that
// of its first content line, or none at the top of the file.
func findKey(lines []string, from, to int, key string) (span, string) {
	child := ""
	for i := from; i < to; i++ {
		if !isContent(lines[i]) {
			continue
		}
		in := indent(lines[i])
		if child == "" && from > 0 {
			child = in
		}
		if in != child {
			continue
		}
		if _, ok := keyRest(lines[i], key); ok {
			return span{at: i, end: valueEnd(lines, i, to)}, child
		}
	}
	return span{at: -1}, child
}

// valueEnd is the line after the last content line of the value of the key
// at line at: the lines below it indented deeper, and, when nothing follows
// the key on its line, the list items at its own indent, as a list under a
// top-level key is often written; up to to. Comments after that last line
// belong to what follows.
func valueEnd(lines []string, at, to int) int {
	end := at + 1
	keyIn := len(indent(lines[at]))
	t := strings.TrimSpace(lines[at])
	_, rest, _ := strings.Cut(t, ":")
	value, _ := splitComment(rest)
	items := value == "" && !strings.HasPrefix(t, "-")
	for i := at + 1; i < to; i++ {
		if !isContent(lines[i]) {
			continue
		}
		in, l := len(indent(lines[i])), strings.TrimSpace(lines[i])
		item := l == "-" || strings.HasPrefix(l, "- ")
		if in < keyIn || in == keyIn && !(items && item) {
			break
		}
		end = i + 1
	}
	return end
}

// contentEnd is the line after the last content line in lines[from:to], or
// from when there is none.
func contentEnd(lines []string, from, to int) int {
	end := from
	for i := from; i < to; i++ {
		if isContent(lines[i]) {
			end = i + 1
		}
	}
	return end
}

// setKey is lines with the key at path written by render: its lines
// replaced where it is, or added at the end of its block, adding the
// blocks above it that are missing.
func setKey(lines []string, path []string, render func(indent, old string) []string) ([]string, *Problem) {
	from, to, parentIndent := 0, len(lines), ""
	for depth, k := range path {
		sp, child := findKey(lines, from, to, k)
		if depth > 0 && child == "" {
			child = parentIndent + "  "
		}
		if sp.at < 0 {
			at := contentEnd(lines, from, to)
			var add []string
			for _, p := range path[depth : len(path)-1] {
				add = append(add, child+p+":")
				child += "  "
			}
			add = append(add, render(child, "")...)
			return slices.Concat(lines[:at], add, lines[at:]), nil
		}
		if depth == len(path)-1 {
			return slices.Concat(lines[:sp.at], render(child, lines[sp.at]), lines[sp.end:]), nil
		}
		rest, _ := keyRest(lines[sp.at], k)
		if value, _ := splitComment(rest); value != "" {
			field := strings.Join(path[:depth+1], ".")
			return nil, &Problem{Field: field, Reason: fmt.Sprintf("is written on one line (%s: %s); write it as a block, a key a line under it, then run this again", k, value)}
		}
		from, to, parentIndent = sp.at+1, sp.end, child
	}
	return lines, nil
}

// removeKey is lines without the key at path and its value; lines as they
// are when it is not there.
func removeKey(lines []string, path []string) ([]string, *Problem) {
	from, to := 0, len(lines)
	for depth, k := range path {
		sp, _ := findKey(lines, from, to, k)
		if sp.at < 0 {
			return lines, nil
		}
		if depth == len(path)-1 {
			return slices.Delete(slices.Clone(lines), sp.at, sp.end), nil
		}
		rest, _ := keyRest(lines[sp.at], k)
		if value, _ := splitComment(rest); value != "" {
			field := strings.Join(path[:depth+1], ".")
			return nil, &Problem{Field: field, Reason: fmt.Sprintf("is written on one line (%s: %s); write it as a block, a key a line under it, then run this again", k, value)}
		}
		from, to = sp.at+1, sp.end
	}
	return lines, nil
}
