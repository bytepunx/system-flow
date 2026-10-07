package verify

import (
	"reflect"
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
