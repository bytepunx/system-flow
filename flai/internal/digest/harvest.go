package digest

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Harvest lists the candidates of the kinds given, every one of Kinds when
// none is given, learned at or after since, a zero since being every one:
// the operator's answers and overruled recommendations on every thread,
// resolved or not, the decisions of done stories' narratives, and the causes
// and fixes of closed issues. Each carries in CitedBy the entries, in ID
// order, whose sources name its source. They come sorted by when they were
// learned, then source, ref, and kind, so the same files always give the
// same list.
func Harvest(repo *workitem.Repo, since time.Time, kinds []string) ([]Candidate, error) {
	want, err := wantedKinds(kinds)
	if err != nil {
		return nil, err
	}
	lookup := Lookup(repo)
	var all []Candidate
	if want["answer"] || want["recommendation"] {
		list, err := threads.List(repo)
		if err != nil {
			return nil, fmt.Errorf("reading the threads for the operator's answers: %w; fix the thread file it names so its front matter parses", err)
		}
		for _, th := range list {
			all = append(all, Answers(th, repo.Manifest.Owner, lookup)...)
		}
	}
	if want["decision"] {
		decisions, err := Decisions(repo, since, lookup)
		if err != nil {
			return nil, err
		}
		all = append(all, decisions...)
	}
	if want["cause"] {
		causes, err := Causes(repo, since, lookup)
		if err != nil {
			return nil, err
		}
		all = append(all, causes...)
	}

	entries, err := List(repo)
	if err != nil {
		return nil, err
	}
	cited := map[string][]string{}
	for _, e := range entries {
		for _, s := range clean(e.Sources, canonicalSource) {
			cited[s] = append(cited[s], e.ID)
		}
	}
	out := []Candidate{}
	for _, c := range all {
		if !want[c.Kind] || learned(c).Before(since) {
			continue
		}
		c.CitedBy = listOrEmpty(cited[canonicalSource(c.Source)])
		out = append(out, c)
	}
	slices.SortStableFunc(out, func(a, b Candidate) int {
		if n := learned(a).Compare(learned(b)); n != 0 {
			return n
		}
		if n := strings.Compare(a.Source, b.Source); n != 0 {
			return n
		}
		if n := strings.Compare(a.Ref, b.Ref); n != 0 {
			return n
		}
		return strings.Compare(a.Kind, b.Kind)
	})
	return out, nil
}

// wantedKinds is the set of kinds given, every one of Kinds when none is,
// or an error naming the kinds when one is none of them.
func wantedKinds(kinds []string) (map[string]bool, error) {
	if len(kinds) == 0 {
		kinds = Kinds
	}
	want := map[string]bool{}
	for _, k := range kinds {
		k = strings.TrimSpace(k)
		if !slices.Contains(Kinds, k) {
			return nil, fmt.Errorf("%q is not a kind of digest candidate; give one of %s", k, strings.Join(Kinds, ", "))
		}
		want[k] = true
	}
	return want, nil
}

// learned is when a candidate was learned, from its At in either form it
// takes: a time, or a day alone, which is that day's midnight in UTC. An At
// in neither form is the zero time.
func learned(c Candidate) time.Time {
	if t, err := time.Parse(workitem.TimeFormat, c.At); err == nil {
		return t
	}
	t, _ := time.Parse(time.DateOnly, c.At)
	return t
}
