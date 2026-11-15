package verify

import (
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

const overlapMod = "github.com/bytepunx/system-flow/flai"

// goSource is a Go file of package name that imports each of pkgs, given
// relative to the flai module's root.
func goSource(name string, pkgs ...string) *fstest.MapFile {
	var b strings.Builder
	b.WriteString("package " + name + "\n")
	for _, p := range pkgs {
		b.WriteString("\nimport _ \"" + overlapMod + "/" + p + "\"\n")
	}
	return &fstest.MapFile{Data: []byte(b.String())}
}

// textFile is a file that is not Go.
func textFile() *fstest.MapFile {
	return &fstest.MapFile{Data: []byte("# text\n")}
}

// overlapRepo is a repository whose flai/ holds the flai module, with files added.
func overlapRepo(files fstest.MapFS) fstest.MapFS {
	fsys := fstest.MapFS{
		"flai/go.mod":  {Data: []byte("module " + overlapMod + "\n\ngo 1.26\n")},
		"flai/main.go": goSource("main"),
	}
	for p, f := range files {
		fsys[p] = f
	}
	return fsys
}

func checkOverlap(t *testing.T, mainFS, storyFS fstest.MapFS, changed, mainPaths, want []string) {
	t.Helper()
	if got := OverlappingPaths(mainFS, storyFS, changed, mainPaths); !reflect.DeepEqual(got, want) {
		t.Errorf("OverlappingPaths(%q, %q) = %q, want %q", changed, mainPaths, got, want)
	}
}

// I-0229: the story's change imports a package main changed, so main's change
// can break it.
func TestOverlapMainChangesAPackageTheStoryImports(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/internal/guard/roles.go":      goSource("guard", "internal/roles"),
		"flai/internal/roles/conditions.go": goSource("roles"),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/internal/guard/roles.go"},
		[]string{"flai/internal/roles/conditions.go"},
		[]string{"flai/internal/roles/conditions.go"})
}

func TestOverlapMainImportsAPackageTheStoryChanges(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/internal/verify/x.go":       goSource("verify"),
		"flai/internal/preview/accept.go": goSource("preview", "internal/verify"),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/internal/verify/x.go"},
		[]string{"flai/internal/preview/accept.go"},
		[]string{"flai/internal/preview/accept.go"})
}

func TestOverlapSamePackage(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/cmd/a.go": goSource("cmd"),
		"flai/cmd/b.go": goSource("cmd"),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/cmd/a.go"},
		[]string{"flai/cmd/b.go"},
		[]string{"flai/cmd/b.go"})
}

// An embedded file belongs to the package that embeds it.
func TestOverlapEmbeddedFileBelongsToItsPackage(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/internal/roles/roles.go":                 goSource("roles"),
		"flai/internal/roles/baseline/orchestrator.md": textFile(),
		"flai/internal/guard/roles.go":                 goSource("guard", "internal/roles"),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/internal/guard/roles.go"},
		[]string{"flai/internal/roles/baseline/orchestrator.md"},
		[]string{"flai/internal/roles/baseline/orchestrator.md"})
}

// A testdata file belongs to the package whose tests read it, though
// testdata holds Go files.
func TestOverlapTestdataBelongsToItsPackage(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/internal/verify/x.go":                    goSource("verify"),
		"flai/internal/verify/testdata/proj/main.go":   goSource("main"),
		"flai/internal/verify/testdata/proj/README.md": textFile(),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/internal/verify/x.go"},
		[]string{"flai/internal/verify/testdata/proj/main.go", "flai/internal/verify/testdata/proj/README.md"},
		[]string{"flai/internal/verify/testdata/proj/main.go", "flai/internal/verify/testdata/proj/README.md"})
}

func TestOverlapUnrelatedPathsDoNotOverlap(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/internal/guard/roles.go": goSource("guard", "internal/roles"),
		"flai/internal/roles/roles.go": goSource("roles"),
		"flai/internal/metrics/m.go":   goSource("metrics", "internal/roles"),
		"design/digest/D-0627-x.md":    textFile(),
		"docs/users/flai.md":           textFile(),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/internal/guard/roles.go"},
		[]string{"flai/internal/metrics/m.go", "design/digest/D-0627-x.md", "docs/users/flai.md"},
		nil)
}

func TestOverlapSameNonGoFile(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{"docs/users/flai.md": textFile()})
	checkOverlap(t, fsys, fsys,
		[]string{"docs/users/flai.md"},
		[]string{"docs/users/flai.md"},
		[]string{"docs/users/flai.md"})
}

// A path outside every module has no package: a neighbour does not overlap.
func TestOverlapOutsideAModuleIsTheSamePathAlone(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"design/x.md": textFile(),
		"design/y.md": textFile(),
		"design/z.go": goSource("design"),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"design/x.md", "design/z.go"},
		[]string{"design/y.md"},
		nil)
}

// Only direct imports count: the story's A imports B, B imports C, and main
// changed C.
func TestOverlapTransitiveImportsDoNotCount(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/internal/a/a.go": goSource("a", "internal/b"),
		"flai/internal/b/b.go": goSource("b", "internal/c"),
		"flai/internal/c/c.go": goSource("c"),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"flai/internal/a/a.go"},
		[]string{"flai/internal/c/c.go"},
		nil)
}

// A file main deleted is read in the story's tree, so it overlaps a package
// the story changes, whether or not main kept the package.
func TestOverlapFileMainDeletedIsReadInTheStory(t *testing.T) {
	mainFS := overlapRepo(fstest.MapFS{"flai/internal/x/keep.go": goSource("x")})
	storyFS := overlapRepo(fstest.MapFS{
		"flai/internal/x/keep.go": goSource("x"),
		"flai/internal/x/old.go":  goSource("x"),
		"flai/internal/gone/g.go": goSource("gone"),
		"flai/internal/gone/h.go": goSource("gone"),
	})
	checkOverlap(t, mainFS, storyFS,
		[]string{"flai/internal/x/keep.go", "flai/internal/gone/g.go"},
		[]string{"flai/internal/x/old.go", "flai/internal/gone/h.go"},
		[]string{"flai/internal/x/old.go", "flai/internal/gone/h.go"})
}

func TestOverlapKeepsMainsOrderWithoutDuplicates(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{
		"flai/cmd/a.go": goSource("cmd"),
		"flai/cmd/b.go": goSource("cmd"),
		"docs/a.md":     textFile(),
		"docs/b.md":     textFile(),
		"docs/c.md":     textFile(),
	})
	checkOverlap(t, fsys, fsys,
		[]string{"docs/a.md", "docs/b.md", "flai/cmd/a.go"},
		[]string{"docs/b.md", "docs/c.md", "flai/cmd/b.go", "docs/a.md", "flai/cmd/b.go", "docs/b.md"},
		[]string{"docs/b.md", "flai/cmd/b.go", "docs/a.md"})
}

func TestOverlapEmptyInputs(t *testing.T) {
	fsys := overlapRepo(fstest.MapFS{"docs/a.md": textFile()})
	checkOverlap(t, fsys, fsys, nil, []string{"docs/a.md"}, nil)
	checkOverlap(t, fsys, fsys, []string{"docs/a.md"}, nil, nil)
}
