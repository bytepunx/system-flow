package workflow

import (
	"reflect"
	"slices"
	"strings"
)

// The mappings the workflow file's entries decode into, each read strictly:
// a key that is not here refuses the file. Fields lists their keys from
// these types, so that the list is the decoding's.

// transitionFields is a transition's entry as YAML gives it.
type transitionFields struct {
	From    names     `yaml:"from"`
	To      string    `yaml:"to"`
	Types   []string  `yaml:"types"`
	By      []string  `yaml:"by"`
	Rules   []RuleRef `yaml:"rules"`
	Effects []string  `yaml:"effects"`
	Verify  *Verify   `yaml:"verify"`
}

// ruleFields is a rule entry's mapping, {rule, severity}.
type ruleFields struct {
	Rule     string `yaml:"rule"`
	Severity string `yaml:"severity"`
}

// timerFields is a timer's entry as YAML gives it.
type timerFields struct {
	Name   string   `yaml:"name"`
	Title  string   `yaml:"title"`
	From   *Point   `yaml:"from"`
	To     *Point   `yaml:"to"`
	Sum    []string `yaml:"sum"`
	Until  *Point   `yaml:"until"`
	Types  []string `yaml:"types"`
	Counts string   `yaml:"counts"`
}

// edgeKindFields is the mapping a point's edge gives in place of a state,
// {kind: <kind>}.
type edgeKindFields struct {
	Kind string `yaml:"kind"`
}

// verifyFields is a transition's verify as a mapping, {suites, severity}.
type verifyFields struct {
	Suites   *Suites `yaml:"suites"`
	Severity string  `yaml:"severity"`
}

// runsFields is a state's runs, {suites, kluster}.
type runsFields struct {
	Suites  *Suites  `yaml:"suites"`
	Kluster []string `yaml:"kluster"`
}

// verifyAtFields is verify_at, a mapping of its points to suites.
type verifyAtFields struct {
	Batch   *Suites `yaml:"batch"`
	Publish *Suites `yaml:"publish"`
	CI      *Suites `yaml:"ci"`
}

// Fields are the keys this flai reads in the workflow file's front matter
// and in each mapping inside it, sorted, by the mapping's name: workflow for
// the front matter, workflow.state for a state, and so on. A flai that does
// not know a key refuses the file, so front-matter-fields.txt lists them
// (S-0806): a release that changes them raises flai.minimum, and an older
// flai stops on the manifest, naming the release, before it reads the file.
func Fields() map[string][]string {
	return map[string][]string{
		"workflow":                   yamlKeys(rawFile{}),
		"workflow.state":             yamlKeys(State{}),
		"workflow.state.runs":        yamlKeys(runsFields{}),
		"workflow.transition":        yamlKeys(transitionFields{}),
		"workflow.transition.rule":   yamlKeys(ruleFields{}),
		"workflow.transition.verify": yamlKeys(verifyFields{}),
		"workflow.timer":             yamlKeys(timerFields{}),
		"workflow.timer.point":       sorted([]string{EdgeFirst, EdgeLast}),
		"workflow.timer.point.edge":  yamlKeys(edgeKindFields{}),
		"workflow.verify_at":         yamlKeys(verifyAtFields{}),
	}
}

// yamlKeys are the keys v's struct decodes, sorted.
func yamlKeys(v any) []string {
	t := reflect.TypeOf(v)
	var out []string
	for i := range t.NumField() {
		name, _, _ := strings.Cut(t.Field(i).Tag.Get("yaml"), ",")
		if name != "" && name != "-" {
			out = append(out, name)
		}
	}
	return sorted(out)
}

func sorted(s []string) []string {
	slices.Sort(s)
	return s
}
