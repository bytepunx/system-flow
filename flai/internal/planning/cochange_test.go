package planning

import (
	"reflect"
	"testing"
)

// commitLog is a fixed commit list: three commits touch cli/, two of them with
// docs/guide.md, one with web/x.go, and one commit touches neither.
var commitLog = [][]string{
	{"cli/a.go", "docs/guide.md", "web/x.go"},
	{"cli/b.go", "docs/guide.md"},
	{"cli/a.go", "cli/a_test.go"},
	{"web/x.go", "web/y.go"},
}

func TestCoChangedCountsFilesInTheSeedCommits(t *testing.T) {
	for _, c := range []struct {
		name     string
		commits  [][]string
		seeds    []string
		minCount int
		limit    int
		want     Suggestions
	}{
		{
			name: "a folder seed holds the files under it", commits: commitLog, seeds: []string{"cli"}, minCount: 1,
			want: Suggestions{Commits: 4, SeedCommits: 3, Suggestions: []Suggestion{
				{Path: "docs/guide.md", Count: 2, Share: 2.0 / 3},
				{Path: "web/x.go", Count: 1, Share: 1.0 / 3},
			}},
		},
		{
			name: "a file seed holds only itself", commits: commitLog, seeds: []string{"cli/a.go"}, minCount: 1,
			want: Suggestions{Commits: 4, SeedCommits: 2, Suggestions: []Suggestion{
				{Path: "cli/a_test.go", Count: 1, Share: 0.5},
				{Path: "docs/guide.md", Count: 1, Share: 0.5},
				{Path: "web/x.go", Count: 1, Share: 0.5},
			}},
		},
		{
			name: "a seed is not a prefix of a longer name", commits: [][]string{{"cli2/a.go", "x.go"}}, seeds: []string{"cli"}, minCount: 1,
			want: Suggestions{Commits: 1, Suggestions: []Suggestion{}},
		},
		{
			name: "min leaves out the rarer files", commits: commitLog, seeds: []string{"cli"}, minCount: 2,
			want: Suggestions{Commits: 4, SeedCommits: 3, Suggestions: []Suggestion{
				{Path: "docs/guide.md", Count: 2, Share: 2.0 / 3},
			}},
		},
		{
			name: "limit cuts after the ordering, ties by path", commits: commitLog, seeds: []string{"cli/a.go"}, minCount: 1, limit: 2,
			want: Suggestions{Commits: 4, SeedCommits: 2, Suggestions: []Suggestion{
				{Path: "cli/a_test.go", Count: 1, Share: 0.5},
				{Path: "docs/guide.md", Count: 1, Share: 0.5},
			}},
		},
		{
			name: "several seeds count a commit once", commits: commitLog, seeds: []string{"cli", "web/y.go"}, minCount: 1,
			want: Suggestions{Commits: 4, SeedCommits: 4, Suggestions: []Suggestion{
				{Path: "docs/guide.md", Count: 2, Share: 0.5},
				{Path: "web/x.go", Count: 2, Share: 0.5},
			}},
		},
		{
			name: "no history", seeds: []string{"cli"}, minCount: 1,
			want: Suggestions{Suggestions: []Suggestion{}},
		},
		{
			name: "no seed commit", commits: commitLog, seeds: []string{"flaiover"}, minCount: 0,
			want: Suggestions{Commits: 4, Suggestions: []Suggestion{}},
		},
	} {
		if got := CoChanged(c.commits, c.seeds, c.minCount, c.limit); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s:\n got %+v\nwant %+v", c.name, got, c.want)
		}
	}
}
