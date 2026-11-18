// Package verifyrecord reads and merges what flai verify records of a
// story's suites: for each suite, the tiers a run selected of it, how each
// ended, and the head and base they passed at (ADR-0312 §7). It reads the
// reports verify stores, a story's last report and its verdict at each
// commit, in the shape verify writes them, and imports nothing of flai's, so
// that what moves a story can read the record without importing verify.
package verifyrecord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// State is how a tier ended in a run, in the words verify's report writes.
type State string

// The states a tier ends a run in, as verify names them.
const (
	Passed     State = "passed"
	Failed     State = "failed"
	NotReached State = "not-reached"
	Reused     State = "reused"
	Skipped    State = "skipped"
	Queued     State = "queued"
)

// Tier is how one tier of a suite ended in the run that recorded the suite.
type Tier struct {
	Name  string `json:"name"`
	State State  `json:"state"`
	// CoveredBy names the tier that ran in place of a Skipped one.
	CoveredBy string `json:"covered_by,omitempty"`
	// ReusedFrom is when the run that ran a Reused tier began.
	ReusedFrom *time.Time `json:"reused_from,omitempty"`
}

// Suite is what a run recorded of one suite: the tiers it selected of it, in
// list order, how each ended, and the head and base the run was at.
type Suite struct {
	Name string `json:"name"`
	// Passed says no tier of the suite failed or was left unreached: each
	// passed, was reused, was left Queued to the merge queue, or was Skipped
	// for a tier that one of those did.
	Passed bool `json:"passed"`
	// Commit is the head the run verified, Base the main branch, BaseCommit
	// the commit Base had, and MergeBase where the branch left it.
	Commit     string `json:"commit"`
	Base       string `json:"base,omitempty"`
	BaseCommit string `json:"base_commit,omitempty"`
	MergeBase  string `json:"merge_base,omitempty"`
	// Patch is the stable patch id of the branch's own change, "" when the
	// worktree had anything uncommitted (ADR-0156).
	Patch string `json:"patch,omitempty"`
	// Fingerprint hashes the project's tiers and the suite's selected ones'
	// commands; "" for a suite read from a report written before suites were
	// recorded, which no later run reuses.
	Fingerprint string `json:"fingerprint,omitempty"`
	// Clean says the verdict is of exactly Commit: the run began with
	// nothing uncommitted, or the close-out committed what it verified.
	Clean bool      `json:"clean,omitempty"`
	RanAt time.Time `json:"ran_at"`
	Tiers []Tier    `json:"tiers"`
}

// Tier is the suite's tier of that name and true, or false when the run did
// not select it.
func (s Suite) Tier(name string) (Tier, bool) {
	for _, t := range s.Tiers {
		if t.Name == name {
			return t, true
		}
	}
	return Tier{}, false
}

// TierNames are the names of the suite's tiers, in its order.
func (s Suite) TierNames() []string {
	names := make([]string, 0, len(s.Tiers))
	for _, t := range s.Tiers {
		names = append(names, t.Name)
	}
	return names
}

// ranNone says the run ran none of the suite's tiers: each was left to the
// merge queue or not reached.
func (s Suite) ranNone() bool {
	return !slices.ContainsFunc(s.Tiers, func(t Tier) bool { return t.State != Queued && t.State != NotReached })
}

// Step is what the record reads of one step of a stored report.
type Step struct {
	Name       string     `json:"name"`
	Tier       bool       `json:"tier,omitempty"`
	State      State      `json:"state"`
	CoveredBy  string     `json:"covered_by,omitempty"`
	ReusedFrom *time.Time `json:"reused_from,omitempty"`
}

// Record is what the record reads of a stored verify report: a story's last
// report, or the verdict kept at one of its commits.
type Record struct {
	Story      string    `json:"story"`
	Commit     string    `json:"commit"`
	Base       string    `json:"base"`
	BaseCommit string    `json:"base_commit,omitempty"`
	MergeBase  string    `json:"merge_base,omitempty"`
	Patch      string    `json:"patch,omitempty"`
	RanAt      time.Time `json:"ran_at"`
	Passed     bool      `json:"passed"`
	Clean      bool      `json:"clean,omitempty"`
	Steps      []Step    `json:"steps"`
	Suites     []Suite   `json:"suites,omitempty"`
}

// SuitesOf is the record's suites: those it holds, or, for a report written
// before suites were recorded, its tier steps grouped as Group groups them.
func (r Record) SuitesOf(suiteOf func(tier string) string) []Suite {
	if len(r.Suites) > 0 {
		return r.Suites
	}
	return Group(r, suiteOf)
}

// Group groups the record's tier steps into suites by suiteOf, the suite of
// each tier (nil, or "" for a tier, is the tier's own name), in the order of
// each suite's first step, each at the record's head and base and with no
// Fingerprint. A Skipped tier counts toward Passed as the tier that covered it
// ended.
func Group(r Record, suiteOf func(tier string) string) []Suite {
	ended := map[string]Step{}
	for _, s := range r.Steps {
		if _, seen := ended[s.Name]; s.Tier && !seen {
			ended[s.Name] = s
		}
	}
	var out []Suite
	at := map[string]int{}
	for _, s := range r.Steps {
		if !s.Tier {
			continue
		}
		name := ""
		if suiteOf != nil {
			name = suiteOf(s.Name)
		}
		if name == "" {
			name = s.Name
		}
		i, ok := at[name]
		if !ok {
			i, at[name] = len(out), len(out)
			out = append(out, Suite{Name: name, Passed: true, Commit: r.Commit, Base: r.Base, BaseCommit: r.BaseCommit,
				MergeBase: r.MergeBase, Patch: r.Patch, Clean: r.Clean, RanAt: r.RanAt, Tiers: []Tier{}})
		}
		if _, dup := out[i].Tier(s.Name); dup {
			continue
		}
		out[i].Tiers = append(out[i].Tiers, Tier{Name: s.Name, State: s.State, CoveredBy: s.CoveredBy, ReusedFrom: s.ReusedFrom})
		if !passing(ended, s) {
			out[i].Passed = false
		}
	}
	return out
}

// passing says the step does not fail its suite: it passed, was reused, or
// was left to the merge queue, or it was skipped for a tier that did one of
// those; a chain of covers is followed at most as far as there are steps.
func passing(ended map[string]Step, s Step) bool {
	for range len(ended) + 1 {
		switch s.State {
		case Passed, Reused, Queued:
			return true
		case Skipped:
			next, ok := ended[s.CoveredBy]
			if !ok {
				return false
			}
			s = next
		default:
			return false
		}
	}
	return false
}

// Merge is old's suites with add's in their place by name, so that a run of
// fewer suites erases none that another run recorded; but an entry of add
// that ran none of its tiers, each left to the merge queue or not reached,
// leaves old's in place. old's order is kept and add's other suites follow in
// theirs. kept names the suites whose entry is old's.
func Merge(old, add []Suite) (merged []Suite, kept []string) {
	byName := map[string]Suite{}
	for _, s := range add {
		if _, seen := byName[s.Name]; !seen {
			byName[s.Name] = s
		}
	}
	placed := map[string]bool{}
	for _, o := range old {
		if placed[o.Name] {
			continue
		}
		placed[o.Name] = true
		if s, ok := byName[o.Name]; ok && !s.ranNone() {
			merged = append(merged, s)
			continue
		}
		merged, kept = append(merged, o), append(kept, o.Name)
	}
	for _, s := range add {
		if !placed[s.Name] {
			placed[s.Name] = true
			merged = append(merged, s)
		}
	}
	return merged, kept
}

// AllPassed says whether every suite of names passed at commit among suites:
// each has an entry at that commit that passed. No names is true.
func AllPassed(suites []Suite, commit string, names []string) bool {
	for _, name := range names {
		i := slices.IndexFunc(suites, func(s Suite) bool { return s.Name == name })
		if i < 0 || !suites[i].Passed || suites[i].Commit != commit {
			return false
		}
	}
	return true
}

// Verdict is what the record of a story at a head says of one of its suites.
type Verdict struct {
	Suite string
	// Passed says the suite passed at the head, in a run that verified
	// exactly it, with every tier of it held: none failed, not reached, or
	// left to the merge queue.
	Passed bool
	// Tiers are the tiers the run selected of the suite, in its order, and
	// Held those of them that passed or were reused, or were skipped for a
	// tier that was.
	Tiers, Held []string
}

// AtHead answers a Verdict for each suite of the verdict the record keeps of
// story at head in the project's cache folder cacheDir, in the record's
// order; ok is false, with no error, when it keeps none there. suiteOf gives
// a tier's suite for a record written before suites were recorded, as Group
// takes it.
func AtHead(cacheDir, story, head string, suiteOf func(tier string) string) (verdicts []Verdict, ok bool, err error) {
	rec, ok, err := Read(CommitPath(cacheDir, story, head))
	if err != nil || !ok {
		return nil, false, err
	}
	suites := rec.SuitesOf(suiteOf)
	tiers := map[string]Tier{}
	for _, s := range suites {
		for _, t := range s.Tiers {
			if _, seen := tiers[t.Name]; !seen {
				tiers[t.Name] = t
			}
		}
	}
	for _, s := range suites {
		v := Verdict{Suite: s.Name, Tiers: s.TierNames()}
		for _, t := range s.Tiers {
			if held(tiers, t) {
				v.Held = append(v.Held, t.Name)
			}
		}
		v.Passed = s.Passed && s.Commit == head && s.Clean && len(v.Held) == len(v.Tiers)
		verdicts = append(verdicts, v)
	}
	return verdicts, true, nil
}

// held says the tier passed or was reused, or was skipped for a tier that
// was, among the record's tiers by name.
func held(tiers map[string]Tier, t Tier) bool {
	for range len(tiers) + 1 {
		switch t.State {
		case Passed, Reused:
			return true
		case Skipped:
			next, ok := tiers[t.CoveredBy]
			if !ok {
				return false
			}
			t = next
		default:
			return false
		}
	}
	return false
}

// ReportPath is where a story's last report is stored in the project's cache
// folder cacheDir: verify/<story>.json.
func ReportPath(cacheDir, story string) string {
	return filepath.Join(cacheDir, "verify", story+".json")
}

// CommitPath is where the record keeps a story's verdict at a commit in the
// project's cache folder cacheDir: verify/<story>/<commit>.json (ADR-0229).
func CommitPath(cacheDir, story, commit string) string {
	return filepath.Join(cacheDir, "verify", story, commit+".json")
}

// Read reads the stored report at path; ok is false, with no error, when
// there is none.
func Read(path string) (rec Record, ok bool, err error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Record{}, false, nil
	}
	if err != nil {
		return Record{}, false, fmt.Errorf("read the verify record %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, false, fmt.Errorf("read the verify record %s: %w; delete it and verify again", path, err)
	}
	return rec, true, nil
}
