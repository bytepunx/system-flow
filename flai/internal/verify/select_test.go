package verify

import (
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

// repoTiers are tiers like this repository's: Go's cheapest first, then
// the dashboard's, a markdown check, and a smoke test only for all.
var repoTiers = []Tier{
	{Name: "gofmt", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"gofmt", "-l", "{files}"}, Format: FormatGofmtList},
	{Name: "go-vet", Dir: "flai", Paths: []string{"flai/**/*.go", "flai/go.mod"}, Command: []string{"go", "vet", "{packages}"}},
	{Name: "go-test", Dir: "flai", Paths: []string{"flai", "!flai/**/*.md"}, Command: []string{"go", "test", "-json", "{packages}"}, Format: FormatGoTestJSON},
	{Name: "golangci", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"golangci-lint", "run", "{packages}"},
		AllCommand: []string{"golangci-lint", "run"}, Format: FormatGolangciJSON},
	{Name: "vitest", Dir: "flaiover", Paths: []string{"flaiover/src/**"}, Command: []string{"npx", "vitest", "run", "{files}"}, Format: FormatVitestJSON},
	{Name: "markdown", Paths: []string{"**/*.md"}, Command: []string{"flai", "check", "--strict"}},
	{Name: "scripts", Paths: []string{"scripts/**"}, Command: []string{"go", "vet", "{packages}"}},
	{Name: "smoke", AllOnly: true, Command: []string{"scripts/template-test.sh"}},
}

func repoFS() fstest.MapFS {
	return fstest.MapFS{
		"flai/go.mod":                                {Data: []byte("module example.com/flai\n")},
		"flai/main.go":                               {},
		"flai/README.md":                             {},
		"flai/internal/verify/verify.go":             {},
		"flai/internal/verify/select.go":             {},
		"flai/internal/verify/testdata/go-test.json": {},
		"flai/internal/verify/testdata/proj/a.go":    {},
		"flai/internal/docs/guide.txt":               {},
		"flaiover/package.json":                      {},
		"flaiover/src/lib/sum.ts":                    {},
		"flaiover/src/lib/sum.test.ts":               {},
		"design/system/x.md":                         {},
		"scripts/flai.sh":                            {},
	}
}

func TestSelectPicksTiersAndFillsTheirPlaceholders(t *testing.T) {
	cases := []struct {
		name  string
		paths []string
		all   bool
		want  map[string][]string // tier name → argv, in list order below
		order []string
	}{
		{
			name:  "a Go file selects the Go tiers with its package and file",
			paths: []string{"flai/internal/verify/select.go"},
			order: []string{"gofmt", "go-vet", "go-test", "golangci"},
			want: map[string][]string{
				"gofmt":    {"gofmt", "-l", "internal/verify/select.go"},
				"go-vet":   {"go", "vet", "./internal/verify"},
				"go-test":  {"go", "test", "-json", "./internal/verify"},
				"golangci": {"golangci-lint", "run", "./internal/verify"},
			},
		},
		{
			name:  "Go files of two packages and go.mod give each package once, sorted",
			paths: []string{"flai/internal/verify/verify.go", "flai/go.mod", "flai/internal/verify/select.go", "flai/main.go"},
			order: []string{"gofmt", "go-vet", "go-test", "golangci"},
			want: map[string][]string{
				"gofmt":    {"gofmt", "-l", "internal/verify/select.go", "internal/verify/verify.go", "main.go"},
				"go-vet":   {"go", "vet", ".", "./internal/verify"},
				"go-test":  {"go", "test", "-json", ".", "./internal/verify"},
				"golangci": {"golangci-lint", "run", ".", "./internal/verify"},
			},
		},
		{
			name:  "a fixture under testdata selects the tests of the package above it",
			paths: []string{"flai/internal/verify/testdata/go-test.json", "flai/internal/verify/testdata/proj/a.go"},
			order: []string{"gofmt", "go-vet", "go-test", "golangci"},
			want: map[string][]string{
				"gofmt":    {"gofmt", "-l", "internal/verify/testdata/proj/a.go"},
				"go-vet":   {"go", "vet", "./internal/verify"},
				"go-test":  {"go", "test", "-json", "./internal/verify"},
				"golangci": {"golangci-lint", "run", "./internal/verify"},
			},
		},
		{
			name:  "a file in a folder with no Go file belongs to the package above",
			paths: []string{"flai/internal/docs/guide.txt"},
			order: []string{"go-test"},
			want:  map[string][]string{"go-test": {"go", "test", "-json", "."}},
		},
		{
			name:  "a flaiover source file selects vitest with the file relative to flaiover",
			paths: []string{"flaiover/src/lib/sum.test.ts", "flaiover/package.json"},
			order: []string{"vitest"},
			want:  map[string][]string{"vitest": {"npx", "vitest", "run", "src/lib/sum.test.ts"}},
		},
		{
			name:  "markdown selects only the markdown check, the ! pattern keeping go-test out",
			paths: []string{"design/system/x.md", "flai/README.md"},
			order: []string{"markdown"},
			want:  map[string][]string{"markdown": {"flai", "check", "--strict"}},
		},
		{
			name:  "a placeholder that comes to nothing leaves its tier out",
			paths: []string{"scripts/flai.sh"},
		},
		{
			name: "no paths select nothing",
		},
		{
			name:  "all selects every tier, with ./... and . or its own command",
			paths: []string{"design/system/x.md"},
			all:   true,
			order: []string{"gofmt", "go-vet", "go-test", "golangci", "vitest", "markdown", "scripts", "smoke"},
			want: map[string][]string{
				"gofmt":    {"gofmt", "-l", "."},
				"go-vet":   {"go", "vet", "./..."},
				"go-test":  {"go", "test", "-json", "./..."},
				"golangci": {"golangci-lint", "run"},
				"vitest":   {"npx", "vitest", "run", "."},
				"markdown": {"flai", "check", "--strict"},
				"scripts":  {"go", "vet", "./..."},
				"smoke":    {"scripts/template-test.sh"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Select(repoFS(), repoTiers, c.paths, c.all)
			var names []string
			for _, s := range got {
				names = append(names, s.Tier.Name)
				if !reflect.DeepEqual(s.Argv, c.want[s.Tier.Name]) {
					t.Errorf("%s argv %q, want %q", s.Tier.Name, s.Argv, c.want[s.Tier.Name])
				}
			}
			if !reflect.DeepEqual(names, c.order) {
				t.Errorf("selected %v, want %v", names, c.order)
			}
		})
	}
}

func TestSelectWithTheTierAtTheRoot(t *testing.T) {
	fsys := fstest.MapFS{"go.mod": {}, "main.go": {}, "cmd/x/x.go": {}, "cmd/x/x.txt": {}}
	tier := Tier{Name: "test", Dir: "./", Paths: []string{"**"}, Command: []string{"go", "test", "{packages}", "{files}"}}
	got := Select(fsys, []Tier{tier}, []string{"cmd/x/x.txt", "go.mod"}, false)
	want := []string{"go", "test", ".", "./cmd/x", "cmd/x/x.txt", "go.mod"}
	if len(got) != 1 || !reflect.DeepEqual(got[0].Argv, want) {
		t.Errorf("got %+v, want argv %q", got, want)
	}
}

func TestSelectStoryAddsTheAllOnlyTiersThePathsSelectOrThatHaveNone(t *testing.T) {
	tiers := []Tier{
		{Name: "go-test", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}},
		{Name: "integration", AllOnly: true, Dir: "flai", Paths: []string{"flai/**"}, Command: []string{"go", "test", "{packages}"}},
		{Name: "template", AllOnly: true, Paths: []string{"template/**", "!template/**/*.md"}, Command: []string{"template-test.sh"}},
		{Name: "flaiover", AllOnly: true, Paths: []string{"flaiover/**"}, Command: []string{"npm", "test", "{files}"}, AllCommand: []string{"npm", "run", "check"}},
		{Name: "smoke", AllOnly: true, Command: []string{"smoke.sh"}},
	}
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/main.go": {}}
	for _, tc := range []struct {
		name  string
		paths []string
		want  []string // name and argv, joined
	}{
		{name: "a Go file selects its tests, the full run, and smoke", paths: []string{"flai/main.go"},
			want: []string{"go-test go test .", "integration go test ./...", "smoke smoke.sh"}},
		{name: "a file of a sub-project selects its whole run with its own command", paths: []string{"flaiover/src/a.ts", "template/x.sh"},
			want: []string{"template template-test.sh", "flaiover npm run check", "smoke smoke.sh"}},
		{name: "an excluded path selects nothing", paths: []string{"template/README.md"},
			want: []string{"smoke smoke.sh"}},
		{name: "no paths still run the tiers with none", want: []string{"smoke smoke.sh"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			for _, s := range SelectStory(fsys, tiers, tc.paths) {
				got = append(got, s.Tier.Name+" "+strings.Join(s.Argv, " "))
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("selected %q, want %q", got, tc.want)
			}
			for _, s := range Select(fsys, tiers, tc.paths, false) {
				if s.Tier.AllOnly {
					t.Errorf("Select chose %s, which only all runs", s.Tier.Name)
				}
			}
		})
	}
}

// coverage is each selected tier's name, with " covered by " and the tier
// that covers it when one does.
func coverage(selected []Selected) []string {
	var out []string
	for _, s := range selected {
		if s.CoveredBy != "" {
			out = append(out, s.Tier.Name+" covered by "+s.CoveredBy)
		} else {
			out = append(out, s.Tier.Name)
		}
	}
	return out
}

// S-0342: a tier selected with one that covers it is marked covered by it,
// with all and in a story's close-out; selected alone, as Select for paths
// leaves it without the all_only tier, it runs.
func TestSelectMarksATierCoveredByASelectedTier(t *testing.T) {
	tiers := []Tier{
		{Name: "go-test", Dir: "flai", Paths: []string{"flai", "!flai/**/*.md"}, Command: []string{"go", "test", "-short", "{packages}"}},
		{Name: "integration", AllOnly: true, Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}, Covers: []string{" go-test "}},
	}
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/main.go": {}}
	for _, tc := range []struct {
		name string
		got  []Selected
		want []string
	}{
		{"Select with all", Select(fsys, tiers, nil, true), []string{"go-test covered by integration", "integration"}},
		{"SelectStory with the covering tier", SelectStory(fsys, tiers, []string{"flai/main.go"}), []string{"go-test covered by integration", "integration"}},
		{"SelectStory without the covering tier", SelectStory(fsys, tiers, []string{"flai/go.mod"}), []string{"go-test"}},
		{"Select for paths, which leaves the all_only tier out", Select(fsys, tiers, []string{"flai/main.go"}, false), []string{"go-test"}},
	} {
		if got := coverage(tc.got); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: selected %q, want %q", tc.name, got, tc.want)
		}
	}
}

// S-0342: covers chain, so a covered tier names the nearest covering tier
// that runs, through tiers that were not selected too; a tier only an
// unselected tier covers runs.
func TestSelectResolvesAChainOfCoversToTheTierThatRuns(t *testing.T) {
	tiers := []Tier{
		{Name: "a", Paths: []string{"x/**"}, Command: []string{"a"}, Covers: []string{"b"}},
		{Name: "b", Paths: []string{"y/**"}, Command: []string{"b"}, Covers: []string{"c"}},
		{Name: "c", Paths: []string{"**"}, Command: []string{"c"}},
		{Name: "d", Paths: []string{"w/**"}, Command: []string{"d"}, Covers: []string{"c"}},
	}
	for _, tc := range []struct {
		name  string
		paths []string
		all   bool
		want  []string
	}{
		{name: "all three of a chain resolve to its head", paths: []string{"x/1", "y/1"}, want: []string{"a", "b covered by a", "c covered by a"}},
		{name: "a chain passes through a tier not selected", paths: []string{"x/1"}, want: []string{"a", "c covered by a"}},
		{name: "a covering tier that runs is nearer than one that is covered", all: true, want: []string{"a", "b covered by a", "c covered by d", "d"}},
		{name: "a tier only an unselected tier covers runs", paths: []string{"y/1"}, want: []string{"b", "c covered by b"}},
		{name: "alone a covered tier runs", paths: []string{"z/1"}, want: []string{"c"}},
	} {
		if got := coverage(Select(fstest.MapFS{}, tiers, tc.paths, tc.all)); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: selected %q, want %q", tc.name, got, tc.want)
		}
	}
}

// ADR-0174: SelectStoryLeavingQueue answers the tiers SelectStory selects
// that are marked queue apart, in list order, and a tier that only a queue
// tier covers runs, while one a chain through a queue tier to a tier that
// runs covers stays covered; SelectStory itself leaves none out.
func TestSelectStoryLeavingQueueAnswersTheQueueTiersApart(t *testing.T) {
	tiers := []Tier{
		{Name: "go-test", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}},
		{Name: "lint", Paths: []string{"flai/**/*.go"}, Command: []string{"lint"}, Covers: []string{"integration"}},
		{Name: "integration", AllOnly: true, Queue: true, Dir: "flai", Paths: []string{"flai/**"}, Command: []string{"go", "test", "{packages}"}, Covers: []string{"go-test", "vet"}},
		{Name: "vet", Paths: []string{"flai/**/*.go"}, Command: []string{"vet"}},
		{Name: "smoke", AllOnly: true, Queue: true, Command: []string{"smoke.sh"}},
	}
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/main.go": {}}
	paths := []string{"flai/main.go"}
	run, queued := SelectStoryLeavingQueue(fsys, tiers, paths, nil)
	if got, want := coverage(run), []string{"go-test covered by lint", "lint", "vet covered by lint"}; !reflect.DeepEqual(got, want) {
		t.Errorf("run %q, want %q", got, want)
	}
	var names []string
	for _, s := range queued {
		names = append(names, s.Tier.Name+" "+strings.Join(s.Argv, " "))
	}
	if want := []string{"integration go test ./...", "smoke smoke.sh"}; !reflect.DeepEqual(names, want) {
		t.Errorf("queued %q, want %q", names, want)
	}

	// without lint, nothing that runs covers go-test or vet
	alone := append(append([]Tier{}, tiers[0]), tiers[2:]...)
	run, _ = SelectStoryLeavingQueue(fsys, alone, paths, nil)
	if got, want := coverage(run), []string{"go-test", "vet"}; !reflect.DeepEqual(got, want) {
		t.Errorf("run without lint %q, want %q", got, want)
	}
	if got, want := coverage(SelectStory(fsys, alone, paths)), []string{"go-test covered by integration", "integration", "vet covered by integration", "smoke"}; !reflect.DeepEqual(got, want) {
		t.Errorf("SelectStory %q, want %q", got, want)
	}
}

// manualTiers are a tier for paths, an all_only one, and two manual ones, one
// with paths that would select it and one that covers the tier for paths.
var manualTiers = []Tier{
	{Name: "go-test", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}},
	{Name: "smoke", AllOnly: true, Command: []string{"smoke.sh"}},
	{Name: "gateway", Manual: true, Command: []string{"gateway-smoke.sh"}},
	{Name: "deep", Manual: true, Dir: "flai", Paths: []string{"flai/**"}, Command: []string{"go", "test", "-race", "{packages}"}, Covers: []string{"go-test"}},
}

// S-0548: no run but one that names it selects a manual tier: not paths that
// match its own, not all, not a story's close-out, not the merge queue's
// batch; and a manual tier left out covers nothing, so the tier it covers
// runs.
func TestSelectLeavesTheManualTiersOut(t *testing.T) {
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/main.go": {}}
	paths := []string{"flai/main.go"}
	run, queued := SelectStoryLeavingQueue(fsys, manualTiers, paths, nil)
	for _, tc := range []struct {
		name string
		got  []Selected
		want []string
	}{
		{"Select for paths", Select(fsys, manualTiers, paths, false), []string{"go-test"}},
		{"Select with all", Select(fsys, manualTiers, paths, true), []string{"go-test", "smoke"}},
		{"SelectStory", SelectStory(fsys, manualTiers, paths), []string{"go-test", "smoke"}},
		{"SelectStoryLeavingQueue's run", run, []string{"go-test", "smoke"}},
		{"SelectStoryLeavingQueue's queued", queued, nil},
	} {
		if got := coverage(tc.got); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: selected %q, want %q", tc.name, got, tc.want)
		}
	}
}

// S-0548: SelectTier chooses the tier of the name, trimmed, alone, with the
// argv it runs with all, whether it is manual, all_only, or neither, and
// answers false for a name no tier has.
func TestSelectTierChoosesOneTierByItsName(t *testing.T) {
	tiers := append(slices.Clone(manualTiers), Tier{Name: "golangci", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"golangci-lint", "run", "{packages}"}, AllCommand: []string{"golangci-lint", "run"}})
	for _, tc := range []struct{ name, want string }{
		{"gateway", "gateway-smoke.sh"},
		{" deep ", "go test -race ./..."},
		{"smoke", "smoke.sh"},
		{"go-test", "go test ./..."},
		{"golangci", "golangci-lint run"},
	} {
		s, ok := SelectTier(tiers, tc.name)
		if !ok || s.Tier.Name != strings.TrimSpace(tc.name) || strings.Join(s.Argv, " ") != tc.want || s.CoveredBy != "" {
			t.Errorf("SelectTier(%q) = %+v, %v; want tier %s running %q", tc.name, s, ok, strings.TrimSpace(tc.name), tc.want)
		}
	}
	for _, name := range []string{"nope", "", "flai/main.go"} {
		if s, ok := SelectTier(tiers, name); ok {
			t.Errorf("SelectTier(%q) chose %+v", name, s)
		}
	}
}

// ADR-0233: SelectStoryLeavingQueue runs a queue tier whose when matches an
// entry of the story's claim, or else a changed path, naming the first
// pattern and the claim entry, then path, it matches; it leaves to the queue
// one whose when matches neither and one with no when; and a queue tier that
// runs covers what it covers.
func TestSelectStoryLeavingQueueRunsAQueueTierWhoseWhenMatches(t *testing.T) {
	t.Parallel()
	queueTier := func(name string, when ...string) Tier {
		return Tier{Name: name, AllOnly: true, Queue: true, Command: []string{name + ".sh"}, When: when}
	}
	tiers := []Tier{
		{Name: "go-test", Dir: "flai", Paths: []string{"flai/**/*.go"}, Command: []string{"go", "test", "{packages}"}},
		{Name: "integration", AllOnly: true, Queue: true, Dir: "flai", Paths: []string{"flai/**"}, Command: []string{"go", "test", "{packages}"}, Covers: []string{"go-test"}, When: []string{"docs/**", "flai/**"}},
		queueTier("by-path", "!flai/internal/x/x_test.go", "flai/internal/x/*.go"),
		queueTier("by-neither", "template/**", "!flai/**"),
		queueTier("by-none"),
	}
	fsys := fstest.MapFS{"flai/go.mod": {}, "flai/internal/x/x.go": {}}
	paths := []string{"flai/internal/x/x.go"}
	claim := []string{"README.md", "flai/internal/check"}
	run, queued := SelectStoryLeavingQueue(fsys, tiers, paths, claim)
	var got []string
	for _, s := range run {
		got = append(got, s.Tier.Name+" covered by "+s.CoveredBy+", ran for "+s.RanFor)
	}
	want := []string{
		"go-test covered by integration, ran for ",
		"integration covered by , ran for its when flai/** matches flai/internal/check",
		"by-path covered by , ran for its when flai/internal/x/*.go matches flai/internal/x/x.go",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("run\n%q\nwant\n%q", got, want)
	}
	var names []string
	for _, s := range queued {
		names = append(names, s.Tier.Name+" ran for "+s.RanFor)
	}
	if want := []string{"by-neither ran for ", "by-none ran for "}; !reflect.DeepEqual(names, want) {
		t.Errorf("queued %q, want %q", names, want)
	}

	// with no claim, integration runs for the changed path, and with a when
	// that matches nothing it is left to the queue and covers nothing
	run, _ = SelectStoryLeavingQueue(fsys, tiers, paths, nil)
	if run[1].RanFor != "its when flai/** matches flai/internal/x/x.go" {
		t.Errorf("integration ran for %q with no claim", run[1].RanFor)
	}
	tiers[1].When = []string{"docs/**"}
	run, queued = SelectStoryLeavingQueue(fsys, tiers, paths, nil)
	if got, want := coverage(run), []string{"go-test", "by-path"}; !reflect.DeepEqual(got, want) {
		t.Errorf("run with integration's when matching nothing %q, want %q", got, want)
	}
	if len(queued) != 3 || queued[0].Tier.Name != "integration" {
		t.Errorf("queued %+v, want integration first of three", queued)
	}
}

// ADR-0233: a when pattern such as flai/** matches a claim entry that is a
// folder below the sub-project and one that is the sub-project's own path,
// as Holds.Claim reads a sub-project's name; a narrower pattern does not
// match a wider folder.
func TestFirstMatchingReadsAClaimsFoldersAndSubProjects(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		patterns []string
		entry    string
		want     string
	}{
		{[]string{"flai/**"}, "flai/internal/check", "flai/**"},
		{[]string{"flai/**"}, "flai", "flai/**"},
		{[]string{"flai"}, "flai/internal/check", "flai"},
		{[]string{"docs/**", "flai/**"}, "flai/cmd/verify.go", "flai/**"},
		{[]string{"flai/**", "!flai/internal/**"}, "flai/internal/check", ""},
		{[]string{"flai/internal/verify/**"}, "flai", ""},
	} {
		pattern, entry, ok := firstMatching(tc.patterns, []string{tc.entry})
		if pattern != tc.want || ok != (tc.want != "") || (ok && entry != tc.entry) {
			t.Errorf("firstMatching(%q, %q) = %q, %q, %v, want %q", tc.patterns, tc.entry, pattern, entry, ok, tc.want)
		}
	}
}

// I-0105: this repository's manifest has flai test check a changed flaiover
// file's format and lint with flaiover-lint, outside --all and before vitest,
// rather than leaving prettier and eslint to the close-out's flaiover tier.
// S-0548: this repository's gateway tier runs only by name, never for a
// path, a story, or --all, so no close-out or batch spends on the paid
// gateway checks.
func TestRepositoryManifestRunsTheGatewayChecksOnlyByName(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	s, ok := SelectTier(tiers, "gateway")
	if !ok || !s.Tier.Manual || s.Tier.Queue || strings.Join(s.Argv, " ") != "scripts/gateway-smoke.sh" {
		t.Fatalf("the gateway tier: %+v %v", s, ok)
	}
	paths := []string{"scripts/gateway-smoke.sh", "scripts/smoke.sh"}
	runs := map[string][]Selected{
		"paths": Select(os.DirFS(root), tiers, paths, false),
		"all":   Select(os.DirFS(root), tiers, nil, true),
		"story": SelectStory(os.DirFS(root), tiers, paths),
	}
	for name, run := range runs {
		for _, sel := range run {
			if sel.Tier.Name == "gateway" {
				t.Errorf("a %s run selects the gateway tier", name)
			}
		}
	}
}

func TestRepositoryManifestLintsChangedFlaioverFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range Select(os.DirFS(root), tiers, []string{"flaiover/src/routes/+page.svelte", "flaiover/package.json"}, false) {
		got = append(got, s.Tier.Name+" "+strings.Join(s.Argv, " "))
	}
	want := []string{
		"flaiover-lint ../scripts/flaiover-lint.sh package.json src/routes/+page.svelte",
		"vitest ../scripts/flaiover-unit.sh related --run --reporter=json src/routes/+page.svelte",
		"design-routes ../scripts/flaiover-unit.sh run --reporter=json src/lib/server/design-routes.test.ts",
		"svelte-check ../scripts/flaiover-check.sh",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("selected %q, want %q", got, want)
	}
}

// I-0206: this repository's manifest has flai test type-check flaiover with
// svelte-check, outside --all, when a script, component, or configuration
// changes, and not for a change svelte-check cannot see, such as the README.
// svelte-check takes no file list, so it ran only in the --all flaiover tier,
// and a type error the flaiover-lint and vitest tiers passed over first showed
// in the close-out. Not gated under -short, so that make test catches a
// manifest that drops the tier.
func TestRepositoryManifestTypeChecksChangedFlaioverFiles(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	const changed = "flaiover/src/lib/age.ts"
	if _, err := os.Stat(filepath.Join(root, changed)); err != nil {
		t.Fatalf("%s: %v; name another flaiover/src/lib file", changed, err)
	}
	argv := func(path string) ([]string, bool) {
		for _, s := range Select(os.DirFS(root), tiers, []string{path}, false) {
			if s.Tier.Name == "svelte-check" {
				return s.Argv, true
			}
		}
		return nil, false
	}
	got, ok := argv(changed)
	if want := []string{"../scripts/flaiover-check.sh"}; !ok || !reflect.DeepEqual(got, want) {
		t.Errorf("%s: svelte-check runs %q (selected %v), want %q", changed, got, ok, want)
	}
	if got, ok := argv("flaiover/README.md"); ok {
		t.Errorf("flaiover/README.md selects svelte-check, running %q; want no svelte-check tier", got)
	}
}

// I-0304: this repository's manifest has flai test run flaiover's tests of
// the built flai, every test file that imports src/lib/server/testing.ts, when
// flai's host API changes. The vitest tier picks tests by the TypeScript
// import graph and keeps only paths under flaiover, so S-0587's change to a
// gate in flai/internal/hostapi/writes.go ran none of them, and
// writes.test.ts first failed in the close-out. Under --all the flaiover tier
// covers it, so those tests run once. Not gated under -short, so that make
// test catches a manifest that drops the tier.
func TestRepositoryManifestRunsFlaioverTestsOfTheBuiltFlaiForHostAPIChanges(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	const changed = "flai/internal/hostapi/writes.go"
	if _, err := os.Stat(filepath.Join(root, changed)); err != nil {
		t.Fatalf("%s: %v; name another flai/internal/hostapi file", changed, err)
	}
	var got []string
	ok := false
	for _, s := range Select(os.DirFS(root), tiers, []string{changed}, false) {
		if s.Tier.Name == "flaiover-flai" {
			got, ok = s.Argv, true
		}
	}
	want := []string{"../scripts/flaiover-unit.sh", "related", "--run", "--reporter=json", "src/lib/server/testing.ts"}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Errorf("%s: flaiover-flai runs %q (selected %v), want %q", changed, got, ok, want)
	}
	covered := ""
	for _, s := range Select(os.DirFS(root), tiers, nil, true) {
		if s.Tier.Name == "flaiover-flai" {
			covered = s.CoveredBy
		}
	}
	if covered != "flaiover" {
		t.Errorf("under --all flaiover-flai is covered by %q, want the flaiover tier, which runs every flaiover test", covered)
	}
}

// I-0304: this repository's manifest has flai test check the settings index
// against the configuration and the manifest when the index changes. A
// change to docs/operators/settings.md alone selected only the markdown tier,
// so a row named wrongly first failed in the close-out's integration tier.
// The settings-index tier runs TestSettingsIndexHasEveryKey under -short, and
// under --all the integration tier covers it, running the whole test. Not
// gated under -short, so that make test catches a manifest that drops it.
//
// I-0308: a change to the manifest's settings catalog passed flai test while
// docs/users/flai-reference.md went stale, because TestReferenceIsCurrent and
// TestFlagIndexIsCurrent skipped under -short. They no longer do, so go-test
// runs them for a change to flai/internal/manifest/settings.go, through its
// importer flai/cmd, and settings-index runs them for a hand edit of either
// generated page.
func TestRepositoryManifestChecksTheSettingsIndexWhenItChanges(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"go", "test", "-json", "-short", "-count=1", "-timeout", "5m", "-run", "^(TestSettingsIndex|TestReferenceIsCurrent$|TestFlagIndexIsCurrent$)", "./cmd"}
	for _, changed := range []string{"docs/operators/settings.md", "docs/users/flai-reference.md"} {
		var got []string
		ok := false
		for _, s := range Select(os.DirFS(root), tiers, []string{changed}, false) {
			if s.Tier.Name == "settings-index" {
				got, ok = s.Argv, true
			}
		}
		if !ok || !reflect.DeepEqual(got, want) {
			t.Errorf("%s: settings-index runs %q (selected %v), want %q", changed, got, ok, want)
		}
	}
	const catalog = "flai/internal/manifest/settings.go"
	var goTest []string
	for _, s := range Select(os.DirFS(root), tiers, []string{catalog}, false) {
		if s.Tier.Name == "go-test" {
			goTest = s.Argv
		}
	}
	if !slices.Contains(goTest, "-short") || !slices.Contains(goTest, "./cmd") {
		t.Errorf("%s: go-test runs %q, want flai/cmd's tests under -short, which hold TestReferenceIsCurrent", catalog, goTest)
	}
	covered := ""
	for _, s := range Select(os.DirFS(root), tiers, nil, true) {
		if s.Tier.Name == "settings-index" {
			covered = s.CoveredBy
		}
	}
	if covered != "integration" {
		t.Errorf("under --all settings-index is covered by %q, want the integration tier, which runs the whole test", covered)
	}
}

// I-0322: this repository's manifest has flai test run design-routes.test.ts
// when the dashboard's design or a route changes. The test reads
// design/system/flaiover-dashboard.md and the routes folder from disk and
// imports neither, so the vitest tier, which runs vitest related on the
// changed files, never selected it: a change to the design alone selected only
// the markdown tier, and a row naming a route flaiover lacks first failed
// under --all or in the merge queue's batch. Under --all the flaiover tier
// covers it. Not gated under -short, so that make test catches a manifest that
// drops the tier.
func TestRepositoryManifestChecksTheDashboardDesignsRoutesWhenEitherChanges(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"../scripts/flaiover-unit.sh", "run", "--reporter=json", "src/lib/server/design-routes.test.ts"}
	for _, changed := range []string{"design/system/flaiover-dashboard.md", "flaiover/src/routes/+layout.svelte"} {
		if _, err := os.Stat(filepath.Join(root, changed)); err != nil {
			t.Fatalf("%s: %v; name another file the tier selects", changed, err)
		}
		var got []string
		ok := false
		for _, s := range Select(os.DirFS(root), tiers, []string{changed}, false) {
			if s.Tier.Name == "design-routes" {
				got, ok = s.Argv, true
			}
		}
		if !ok || !reflect.DeepEqual(got, want) {
			t.Errorf("%s: design-routes runs %q (selected %v), want %q", changed, got, ok, want)
		}
	}
	covered := ""
	for _, s := range Select(os.DirFS(root), tiers, nil, true) {
		if s.Tier.Name == "design-routes" {
			covered = s.CoveredBy
		}
	}
	if covered != "flaiover" {
		t.Errorf("under --all design-routes is covered by %q, want the flaiover tier, which runs every flaiover test", covered)
	}
}

// S-0548: this repository's go-test tier runs the changed packages' tests
// with -short, reversing I-0201's choice. The tests that skip under -short,
// those that need real git or the monorepo, are the integration tier's own;
// among them the four that check, lint, and prime the whole monorepo
// (TestMonorepoIsClean, TestRepositoryLintsClean, TestPrimeOnThisRepository,
// TestPrimeComesInPartsOnThisRepository) took minutes under -race in every
// close-out that touched their packages and test the repository, not the
// change. The integration tier runs scripts/integration.sh, without -short,
// so they all still run in the merge queue's batch. go-test skips the paid
// gateway tests, as integration.sh does, and names its timeout. It runs the
// changed packages' importers too (S-0547, I-0221). Not gated
// under -short, so that make test catches a manifest that changes either.
func TestRepositoryManifestRunsChangedGoPackagesShortAndIntegrationWhole(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range Select(os.DirFS(root), tiers, []string{"flai/internal/verify/select.go"}, false) {
		if s.Tier.Name == "go-test" {
			got = s.Argv
		}
	}
	// one package at a time, every go process on one core, and no race
	// detector, so that close-outs at once do not saturate the host (S-0722)
	flags := []string{"env", "GOMAXPROCS=1", "go", "test", "-json", "-p", "1", "-short", "-count=1", "-timeout", "10m", "-skip", "^TestGateway"}
	if len(got) <= len(flags) || !reflect.DeepEqual(got[:len(flags)], flags) {
		t.Fatalf("go-test runs %q, want %q and the packages", got, flags)
	}
	// a change to verify runs cmd's tests too, which call it
	for _, pkg := range []string{"./internal/verify", "./cmd"} {
		if !slices.Contains(got[len(flags):], pkg) {
			t.Errorf("go-test runs %q, want it to hold %s", got, pkg)
		}
	}
	var integration []string
	for _, s := range Select(os.DirFS(root), tiers, []string{"flai/internal/verify/select.go"}, true) {
		if s.Tier.Name == "integration" {
			integration = s.Argv
		}
	}
	if want := []string{"../scripts/integration.sh", "-json"}; !reflect.DeepEqual(integration, want) {
		t.Errorf("integration runs %q, want %q, which runs go test without -short", integration, want)
	}
}

// S-0722: the race detector runs in this repository's race tier alone, the
// full suite under -race on every core, and only by name, flai test race, as
// the gateway tier does: no close-out, --all run, or batch runs it unasked.
// Not gated under -short, so that make test catches a manifest that changes it.
func TestRepositoryManifestRunsTheRaceDetectorOnlyByName(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	s, ok := SelectTier(tiers, "race")
	if !ok || !s.Tier.Manual || s.Tier.Queue || strings.Join(s.Argv, " ") != "../scripts/integration.sh --race -json" {
		t.Fatalf("the race tier: %+v %v", s, ok)
	}
	paths := []string{"scripts/integration.sh", "flai/internal/verify/select.go"}
	runs := map[string][]Selected{
		"paths": Select(os.DirFS(root), tiers, paths, false),
		"all":   Select(os.DirFS(root), tiers, nil, true),
		"story": SelectStory(os.DirFS(root), tiers, paths),
	}
	for name, run := range runs {
		for _, sel := range run {
			if sel.Tier.Name == "race" {
				t.Errorf("a %s run selects the race tier", name)
			}
			if slices.Contains(sel.Argv, "-race") || slices.Contains(sel.Argv, "--race") {
				t.Errorf("a %s run's tier %s runs %q, with the race detector", name, sel.Tier.Name, sel.Argv)
			}
		}
	}
}

// I-0199: a change to a script runs ShellCheck on it, and each script test
// under scripts/ is run, in flai verify and the merge queue's batch, by a tier
// that both selects it and runs it from its command, directly or through the
// scripts that command runs. S-0481 changed scripts/flaiover-install.sh and
// added scripts/flaiover-install-test.sh, and no tier ran the test. Not gated
// under -short, so that make test catches a manifest, or a script test added
// later, that opens the gap again.
func TestRepositoryManifestRunsScriptTestsForScriptChanges(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	tiers, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(root)
	got := map[string][]string{}
	for _, s := range SelectStory(fsys, tiers, []string{"scripts/flaiover-install.sh", "scripts/flaiover-install-test.sh"}) {
		got[s.Tier.Name] = s.Argv
	}
	if want := []string{"scripts/shellcheck.sh", "scripts/flaiover-install-test.sh", "scripts/flaiover-install.sh"}; !reflect.DeepEqual(got["shellcheck"], want) {
		t.Errorf("shellcheck runs %q, want %q", got["shellcheck"], want)
	}
	if _, ok := got["flaiover"]; !ok {
		t.Errorf("S-0481's change selects %v, want the flaiover tier, which runs scripts/flaiover-install-test.sh", slices.Sorted(maps.Keys(got)))
	}

	// The scripts named *-test.sh that no tier runs, each with its reason.
	notRun := map[string]string{
		"scripts/install-published-test.sh": "needs the network and a token; CI's install-published.yml runs it",
		"scripts/flai-test.sh":              "runs every tier, as make flai-test; it tests no script",
	}
	tests, err := filepath.Glob(filepath.Join(root, "scripts", "*-test.sh"))
	if err != nil || len(tests) == 0 {
		t.Fatalf("no scripts/*-test.sh under %s: %v", root, err)
	}
	for _, abs := range tests {
		p := "scripts/" + filepath.Base(abs)
		if _, ok := notRun[p]; ok {
			continue
		}
		var by []string
		for _, s := range SelectStory(fsys, tiers, []string{p}) {
			if s.Tier.Name != "shellcheck" && scriptsRun(t, root, s.Argv)[p] {
				by = append(by, s.Tier.Name)
			}
		}
		if len(by) == 0 {
			t.Errorf("%s: no tier both selects it and runs it; add it to the paths of the tier whose command runs it in system-flow.yaml (I-0199)", p)
		}
	}
}

// scriptName finds the scripts a command or a script names, by base name.
var scriptName = regexp.MustCompile(`[A-Za-z0-9_.-]+\.sh\b`)

// scriptsRun are the scripts/*.sh that argv runs: those it names, and those
// each of them names outside a comment, followed in turn.
func scriptsRun(t *testing.T, root string, argv []string) map[string]bool {
	t.Helper()
	seen := map[string]bool{}
	var queue []string
	add := func(text string) {
		for _, name := range scriptName.FindAllString(text, -1) {
			p := "scripts/" + name
			if _, err := os.Stat(filepath.Join(root, p)); err == nil && !seen[p] {
				seen[p] = true
				queue = append(queue, p)
			}
		}
	}
	add(strings.Join(argv, " "))
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		data, err := os.ReadFile(filepath.Join(root, p))
		if err != nil {
			t.Fatal(err)
		}
		for line := range strings.Lines(string(data)) {
			if !strings.HasPrefix(strings.TrimSpace(line), "#") {
				add(line)
			}
		}
	}
	return seen
}

// I-0128: a tier whose command holds {importers} runs on the changed
// packages and every package of their module that imports them, where
// {packages} names the changed ones only; with all, on ./...; with no
// go.mod, on what {packages} names.
func TestSelectFillsImportersWithTheModulesImporters(t *testing.T) {
	tiers := []Tier{
		{Name: "vet", Dir: "mod", Paths: []string{"mod/**"}, Command: []string{"go", "vet", "{importers}"}},
		{Name: "go-test", Dir: "mod", Paths: []string{"mod/**"}, Command: []string{"go", "test", "{packages}"}},
		{Name: "loose", Dir: "loose", Paths: []string{"loose/**"}, Command: []string{"go", "vet", "{importers}"}},
		{Name: "integration", AllOnly: true, Dir: "mod", Paths: []string{"mod/**"}, Command: []string{"go", "test", "{importers}"}},
	}
	fsys := moduleFS()
	fsys["loose/p/p.go"] = &fstest.MapFile{Data: []byte("package p\n")}
	fsys["loose/q/q.go"] = &fstest.MapFile{Data: []byte("package q\n\nimport \"example.com/loose/p\"\n")}
	fsys["loose/notes/n.txt"] = &fstest.MapFile{}
	for _, tc := range []struct {
		name string
		got  []Selected
		want []string // name and argv, joined
	}{
		{name: "Select for a change to a", got: Select(fsys, tiers, []string{"mod/a/a.go"}, false),
			want: []string{"vet go vet ./a ./b ./c ./d", "go-test go test ./a"}},
		{name: "Select for changes to c and the root", got: Select(fsys, tiers, []string{"mod/c/c.go", "mod/go.mod"}, false),
			want: []string{"vet go vet . ./c ./root", "go-test go test . ./c"}},
		{name: "SelectStory for a change to a", got: SelectStory(fsys, tiers, []string{"mod/a/a.go"}),
			want: []string{"vet go vet ./a ./b ./c ./d", "go-test go test ./a", "integration go test ./..."}},
		{name: "Select with all", got: Select(fsys, tiers, nil, true),
			want: []string{"vet go vet ./...", "go-test go test ./...", "loose go vet ./...", "integration go test ./..."}},
		{name: "no go.mod gives the packages", got: Select(fsys, tiers, []string{"loose/p/p.go"}, false),
			want: []string{"loose go vet ./p"}},
		{name: "no package leaves the tier out", got: Select(fsys, tiers, []string{"loose/notes/n.txt"}, false)},
	} {
		var got []string
		for _, s := range tc.got {
			got = append(got, s.Tier.Name+" "+strings.Join(s.Argv, " "))
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: selected %q, want %q", tc.name, got, tc.want)
		}
	}
}
