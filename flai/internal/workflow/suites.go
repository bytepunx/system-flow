package workflow

import (
	"fmt"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// SuitesAll is the word that names every suite, in place of a list
// (ADR-0312 §2); no suite is called it.
const SuitesAll = "all"

// The cluster steps a state's runs may name (ADR-0312 §4).
const (
	KlusterDeploy = "deploy"
	KlusterChecks = "checks"
)

// KlusterSteps are the cluster steps a state's runs may name, in the order
// they run.
var KlusterSteps = []string{KlusterDeploy, KlusterChecks}

// The points outside an item's moves that verify_at gives suites for
// (ADR-0312 §6).
const (
	VerifyAtBatch   = "batch"
	VerifyAtPublish = "publish"
	VerifyAtCI      = "ci"
)

// VerifyAtPoints are the points verify_at takes, in ADR-0312 §6's order.
var VerifyAtPoints = []string{VerifyAtBatch, VerifyAtPublish, VerifyAtCI}

// Suites are the suites a step names: every suite (All) or a list of them
// (Names), and the file line they are given on.
type Suites struct {
	All   bool     `json:"all,omitempty"`
	Names []string `json:"names,omitempty"`
	// Line is the file line the suites' value starts on: all's, or the
	// list's first entry's.
	Line int `json:"line,omitempty"`
}

// Verify is what a move by a transition needs its last flai verify to have
// passed: the suites, and the severity of the rule, empty for refuse
// (ADR-0312 §2).
type Verify struct {
	Suites   Suites `json:"suites"`
	Severity string `json:"severity,omitempty"`
	// Line is the file line of the transition's verify key.
	Line int `json:"line"`
}

// Runs is what a state runs when a story enters it: suites, the cluster's
// steps, or both (ADR-0312 §4).
type Runs struct {
	Suites  *Suites  `json:"suites,omitempty"`
	Kluster []string `json:"kluster,omitempty"`
	// Line is the file line of the state's runs key.
	Line int `json:"line"`
}

// VerifyAt is the suites each point outside an item's moves runs, nil for
// a point the file does not name (ADR-0312 §6).
type VerifyAt struct {
	Batch   *Suites `json:"batch,omitempty"`
	Publish *Suites `json:"publish,omitempty"`
	CI      *Suites `json:"ci,omitempty"`
}

// At is the suites point runs, nil when the file names none for it or point
// is no point of VerifyAtPoints.
func (v VerifyAt) At(point string) *Suites {
	switch point {
	case VerifyAtBatch:
		return v.Batch
	case VerifyAtPublish:
		return v.Publish
	case VerifyAtCI:
		return v.CI
	}
	return nil
}

const suitesShape = "suites are all, or a list of suite names such as [static, behavior]; a single suite is a list of one, such as [static] (ADR-0312 §2)"

// UnmarshalYAML reads all or a list of suite names.
func (s *Suites) UnmarshalYAML(n ast.Node) error {
	line := n.GetToken().Position.Line
	seq, ok := n.(*ast.SequenceNode)
	if !ok {
		var word string
		if pairs(n) != nil || yaml.NodeToValue(n, &word) != nil || word != SuitesAll {
			return shapeError(n, suitesShape)
		}
		*s = Suites{All: true, Line: line}
		return nil
	}
	out := Suites{Names: make([]string, 0, len(seq.Values)), Line: line}
	if len(seq.Values) > 0 {
		out.Line = seq.Values[0].GetToken().Position.Line
	}
	for _, v := range seq.Values {
		var name string
		if !scalar(v) || yaml.NodeToValue(v, &name) != nil || name == "" {
			at := v
			if at.GetToken() == nil {
				at = n
			}
			return shapeError(at, "a suite is named by a word, such as static; "+suitesShape)
		}
		out.Names = append(out.Names, name)
	}
	*s = out
	return nil
}

const verifyShape = "verify is all, a list of suites such as [static, behavior], or {suites: all | [<suite>, ...], severity: refuse | warn} (ADR-0312 §2)"

// UnmarshalYAML reads a transition's verify: all, a list of suites, or
// {suites, severity}.
func (v *Verify) UnmarshalYAML(n ast.Node) error {
	if _, ok := n.(*ast.MappingNode); !ok && pairs(n) == nil {
		var s Suites
		if err := s.UnmarshalYAML(n); err != nil {
			return err
		}
		*v = Verify{Suites: s}
		return nil
	}
	if err := noValue(n, "suites", "verify's suites is all or a list of suites; give them, or leave verify out"); err != nil {
		return err
	}
	var f verifyFields
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	if f.Suites == nil {
		return shapeError(n, "verify's mapping needs suites, all or a list of suites, besides severity; "+verifyShape)
	}
	*v = Verify{Suites: *f.Suites, Severity: f.Severity}
	return nil
}

const runsShape = "runs is a mapping, {suites: all | [<suite>, ...], kluster: [deploy, checks]}, with suites, kluster, or both (ADR-0312 §4)"

// UnmarshalYAML reads a state's runs: {suites, kluster}.
func (r *Runs) UnmarshalYAML(n ast.Node) error {
	if _, ok := n.(*ast.MappingNode); !ok && pairs(n) == nil {
		return shapeError(n, runsShape)
	}
	for _, key := range []string{"suites", "kluster"} {
		if err := noValue(n, key, fmt.Sprintf("runs gives %s no value; give it, or leave it out: %s", key, runsShape)); err != nil {
			return err
		}
	}
	var f runsFields
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	*r = Runs{Suites: f.Suites, Kluster: f.Kluster}
	return nil
}

const verifyAtShape = "verify_at is a mapping of batch, publish, and ci, each all or a list of suites, such as {batch: all, ci: [static]} (ADR-0312 §6)"

// UnmarshalYAML reads verify_at: a mapping of its points to suites.
func (v *VerifyAt) UnmarshalYAML(n ast.Node) error {
	if _, ok := n.(*ast.MappingNode); !ok && pairs(n) == nil {
		return shapeError(n, verifyAtShape)
	}
	for _, point := range VerifyAtPoints {
		if err := noValue(n, point, fmt.Sprintf("verify_at gives %s no value; give it all or a list of suites, or leave it out", point)); err != nil {
			return err
		}
	}
	var f verifyAtFields
	if err := yaml.NodeToValue(n, &f, yaml.Strict()); err != nil {
		return err
	}
	*v = VerifyAt(f)
	return nil
}

// noValue is an error at key in mapping n when the file gives key no value.
func noValue(n ast.Node, key, msg string) error {
	if kv := pairOf(n, key); kv != nil && (kv.Value == nil || kv.Value.Type() == ast.NullType) {
		return shapeError(kv.Key, msg)
	}
	return nil
}

// pairOf is mapping n's pair of key, nil when it has none.
func pairOf(n ast.Node, key string) *ast.MappingValueNode {
	for _, kv := range pairs(n) {
		if kv.Key.GetToken().Value == key {
			return kv
		}
	}
	return nil
}

// pairLine is the file line of key in mapping n, 0 when it has none.
func pairLine(n ast.Node, key string) int {
	if kv := pairOf(n, key); kv != nil {
		return kv.Key.GetToken().Position.Line
	}
	return 0
}

// scalar reports whether n is a single value, no mapping or list.
func scalar(n ast.Node) bool {
	switch n.Type() {
	case ast.MappingType, ast.MappingValueType, ast.SequenceType, ast.NullType:
		return false
	}
	return true
}

// suites is the problem with the suite names s gives at field, the path of
// their list.
func (p *reader) suites(field string, s Suites) *Problem {
	if s.All {
		return nil
	}
	if len(s.Names) == 0 {
		return p.fail(field, "the list of suites is empty; give suites, or all for every suite (ADR-0312 §2)")
	}
	for j, name := range s.Names {
		at := fmt.Sprintf("%s[%d]", field, j)
		switch {
		case name == SuitesAll:
			return p.fail(at, "all names every suite and is no suite's name; write all alone, in place of the list, or name the suites")
		case !kebab.MatchString(name):
			return p.fail(at, "suite %q is not kebab-case; name a suite as the manifest's tests give it, lowercase words joined by dashes", name)
		case slices.Index(s.Names, name) != j:
			return p.fail(at, "suite %s is named twice; name each suite once", name)
		}
	}
	return nil
}

// verify is the problem with a transition's verify v, at at, the
// transition's path.
func (p *reader) verify(at string, v *Verify) *Problem {
	if v == nil {
		return nil
	}
	field := at + ".verify"
	if _, ok := p.lines[field+".suites"]; ok {
		field += ".suites"
	}
	if prob := p.suites(field, v.Suites); prob != nil {
		return prob
	}
	if v.Severity != "" && !slices.Contains(Severities, v.Severity) {
		return p.fail(at+".verify.severity", "verify's severity %q is not one of %s", v.Severity, strings.Join(Severities, ", "))
	}
	return nil
}

// runs is the problem with state s's runs, at at, the state's path.
func (p *reader) runs(at string, s State) *Problem {
	r := s.Runs
	if r == nil {
		return nil
	}
	if r.Suites == nil && r.Kluster == nil {
		return p.fail(at+".runs", "state %s's runs names nothing; give suites, kluster, or both, or leave runs out: %s", s.Name, runsShape)
	}
	if r.Suites != nil {
		if prob := p.suites(at+".runs.suites", *r.Suites); prob != nil {
			return prob
		}
	}
	if r.Kluster != nil && len(r.Kluster) == 0 {
		return p.fail(at+".runs.kluster", "state %s's kluster is empty; give %s, or leave kluster out", s.Name, strings.Join(KlusterSteps, ", "))
	}
	for j, k := range r.Kluster {
		field := fmt.Sprintf("%s.runs.kluster[%d]", at, j)
		if !slices.Contains(KlusterSteps, k) {
			return p.fail(field, "kluster step %q is not one of %s (ADR-0312 §4)", k, strings.Join(KlusterSteps, ", "))
		}
		if slices.Index(r.Kluster, k) != j {
			return p.fail(field, "kluster step %s is named twice; name each step once", k)
		}
	}
	return nil
}

// verifyAt is the problem with the file's verify_at v, nil when the file
// gives none.
func (p *reader) verifyAt(v VerifyAt) *Problem {
	if _, ok := p.lines["verify_at"]; !ok {
		return nil
	}
	named := false
	for _, point := range VerifyAtPoints {
		s := v.At(point)
		if s == nil {
			continue
		}
		named = true
		if prob := p.suites("verify_at."+point, *s); prob != nil {
			return prob
		}
	}
	if !named {
		return p.fail("verify_at", "verify_at names no point; give %s, or leave verify_at out (ADR-0312 §6)", strings.Join(VerifyAtPoints, ", "))
	}
	return nil
}

// Where a suite is named, as a Named gives it.
const (
	WhereTransition = "transition"
	WhereState      = "state"
	WherePoint      = "point"
)

// Named is a suite the workflow file names and where: in a transition's
// verify, a state's runs, or a point of verify_at.
type Named struct {
	Suite string `json:"suite"`
	// Where is WhereTransition, WhereState, or WherePoint.
	Where string `json:"where"`
	// Name is the transition as "in-progress → review", the state's name, or
	// the point's.
	Name string `json:"name"`
	// Line is the file line the suites are given on.
	Line int `json:"line"`
}

// NamedSuites lists each suite name the file gives, once for every place it
// is given, in the order transitions, states, then verify_at's points, and
// whether any place gives all. A place that gives all lists no name.
func (w Workflow) NamedSuites() (named []Named, all bool) {
	add := func(s *Suites, where, name string) {
		if s == nil {
			return
		}
		if s.All {
			all = true
		}
		for _, suite := range s.Names {
			named = append(named, Named{Suite: suite, Where: where, Name: name, Line: s.Line})
		}
	}
	for _, t := range w.Transitions {
		if t.Verify != nil {
			add(&t.Verify.Suites, WhereTransition, strings.Join(t.From, ", ")+" → "+t.To)
		}
	}
	for _, s := range w.States {
		if s.Runs != nil {
			add(s.Runs.Suites, WhereState, s.Name)
		}
	}
	for _, point := range VerifyAtPoints {
		add(w.VerifyAt.At(point), WherePoint, point)
	}
	return named, all
}

// cloneSuites is a copy of s that shares no list with it.
func cloneSuites(s *Suites) *Suites {
	if s == nil {
		return nil
	}
	c := *s
	c.Names = slices.Clone(s.Names)
	return &c
}
