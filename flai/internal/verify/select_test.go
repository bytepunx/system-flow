package verify

import (
	"os"
	"path/filepath"
	"reflect"
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
	run, queued := SelectStoryLeavingQueue(fsys, tiers, paths)
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
	run, _ = SelectStoryLeavingQueue(fsys, alone, paths)
	if got, want := coverage(run), []string{"go-test", "vet"}; !reflect.DeepEqual(got, want) {
		t.Errorf("run without lint %q, want %q", got, want)
	}
	if got, want := coverage(SelectStory(fsys, alone, paths)), []string{"go-test covered by integration", "integration", "vet covered by integration", "smoke"}; !reflect.DeepEqual(got, want) {
		t.Errorf("SelectStory %q, want %q", got, want)
	}
}

// I-0105: this repository's manifest has flai test check a changed flaiover
// file's format and lint with flaiover-lint, outside --all and before vitest,
// rather than leaving prettier and eslint to the close-out's flaiover tier.
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
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("selected %q, want %q", got, want)
	}
}

// I-0201: this repository's go-test tier runs the changed packages' tests
// whole, without -short, so a real-git or monorepo test a story adds runs in
// flai test and the close-out rather than only in the merge queue's batch;
// it skips the paid gateway tests, as scripts/integration.sh does. Not gated
// under -short, so that make test catches a manifest that brings -short back.
func TestRepositoryManifestRunsChangedGoPackagesWithoutShort(t *testing.T) {
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
	want := []string{"go", "test", "-json", "-race", "-count=1", "-skip", "^TestGateway", "./internal/verify"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("go-test runs %q, want %q", got, want)
	}
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
