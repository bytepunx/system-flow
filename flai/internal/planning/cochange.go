package planning

import (
	"sort"
	"strings"
)

// Suggestions is what CoChanged found in a project's commits.
type Suggestions struct {
	Commits     int          `json:"commits"`
	SeedCommits int          `json:"seed_commits"`
	Suggestions []Suggestion `json:"suggestions"`
}

// Suggestion is a file changed together with the seeds: in how many of the
// seed commits, and that count's share of them.
type Suggestion struct {
	Path  string  `json:"path"`
	Count int     `json:"count"`
	Share float64 `json:"share"`
}

// CoChanged lists the files most often changed in the commits that changed a
// file under a seed, the files under a seed left out. A seed holds a file
// when it is the file or a folder that holds it, as a claim does. Each
// commit is the files it changed, without repeats. A file is kept when it is
// in at least minCount seed commits; files are ordered by that count, most
// first, then by path, and at most limit are kept, every one when limit is
// zero or less.
func CoChanged(commits [][]string, seeds []string, minCount, limit int) Suggestions {
	out := Suggestions{Commits: len(commits), Suggestions: []Suggestion{}}
	counts := map[string]int{}
	for _, files := range commits {
		var others []string
		seeded := false
		for _, f := range files {
			if under(f, seeds) {
				seeded = true
			} else {
				others = append(others, f)
			}
		}
		if !seeded {
			continue
		}
		out.SeedCommits++
		for _, f := range others {
			counts[f]++
		}
	}
	for f, n := range counts {
		if n >= minCount {
			out.Suggestions = append(out.Suggestions, Suggestion{Path: f, Count: n, Share: float64(n) / float64(out.SeedCommits)})
		}
	}
	sort.Slice(out.Suggestions, func(i, j int) bool {
		a, b := out.Suggestions[i], out.Suggestions[j]
		if a.Count != b.Count {
			return a.Count > b.Count
		}
		return a.Path < b.Path
	})
	if limit > 0 && len(out.Suggestions) > limit {
		out.Suggestions = out.Suggestions[:limit]
	}
	return out
}

// under reports whether a seed is the file or a folder that holds it.
func under(file string, seeds []string) bool {
	for _, s := range seeds {
		if file == s || strings.HasPrefix(file, s+"/") {
			return true
		}
	}
	return false
}
