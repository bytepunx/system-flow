package manifest

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// NamedCommand is one command by name: an argument list, run as it stands,
// never through a shell. In an argument, {story} is replaced by the story's
// ID and {root} by the directory it runs in; nothing else is interpreted.
type NamedCommand struct {
	Name    string   `yaml:"name" json:"name"`
	Command []string `yaml:"command" json:"command"`
}

// TestTier is one tier of the project's tests: a command, the paths that
// select it, and how its output becomes findings (S-0273).
type TestTier struct {
	// Name names the tier, unique among the project's tiers.
	Name string `yaml:"name" json:"name"`
	// Command is an argument list, run as it stands, never through a shell.
	// An argument that is PlaceholderPackages, PlaceholderFiles, or
	// PlaceholderImporters alone stands for what the tier's paths selected;
	// nothing else is interpreted.
	Command []string `yaml:"command" json:"command"`
	// Dir is the folder the command runs in, relative to the repository
	// root; empty means the root.
	Dir string `yaml:"dir,omitempty" json:"dir,omitempty"`
	// Paths are glob patterns relative to the repository root, as
	// claims.shared's are, that select the tier for a changed or named path;
	// one beginning with ! takes paths out again. Required unless AllOnly or
	// Manual.
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	// Format is how the command's output becomes findings: one of
	// TestFormats; empty means FormatPlain.
	Format string `yaml:"format,omitempty" json:"format,omitempty"`
	// AllOnly runs the tier only when every tier is asked for, as flai test
	// --all does, never for paths.
	AllOnly bool `yaml:"all_only,omitempty" json:"all_only,omitempty"`
	// AllCommand, when set, runs instead of Command when every tier is asked
	// for. It takes no placeholders.
	AllCommand []string `yaml:"all_command,omitempty" json:"all_command,omitempty"`
	// FixCommand, when set, is an argument list run before flai commits the
	// files the tier's paths select, with placeholders as Command takes
	// them. It rewrites those files in place; its exit status is ignored.
	FixCommand []string `yaml:"fix_command,omitempty" json:"fix_command,omitempty"`
	// Covers are the names of the tiers this one runs everything of, so that
	// when it is selected with them only this one runs (S-0342).
	Covers []string `yaml:"covers,omitempty" json:"covers,omitempty"`
	// Queue leaves the tier to the merge queue: while the queue runs, flai
	// verify leaves it out and the queue's batch runs it (S-0405).
	Queue bool `yaml:"queue,omitempty" json:"queue,omitempty"`
	// Manual runs the tier only when it is named, as flai test <tier> does:
	// never for paths, --all, a close-out, or the merge queue's batch
	// (S-0548). It takes neither Queue nor FixCommand.
	Manual bool `yaml:"manual,omitempty" json:"manual,omitempty"`
	// When are glob patterns, in the language of Paths, that, while the
	// merge queue can batch, have a story's close-out run a Queue tier too
	// when one matches the story's touches or changed paths, rather than
	// leave it to the batch alone (S-0547). Allowed only beside Queue.
	When []string `yaml:"when,omitempty" json:"when,omitempty"`
	// Parallel lets the tier start beside the tiers before it that are marked
	// Parallel too: it waits only for the unmarked tiers before it in the
	// list, while an unmarked tier waits for every tier before it (S-0610).
	// The verdict is a serial run's. It takes no Manual.
	Parallel bool `yaml:"parallel,omitempty" json:"parallel,omitempty"`
}

// The placeholders a tier's command may hold, each as an argument of its own.
const (
	// PlaceholderPackages stands for the packages of the paths the tier
	// selected.
	PlaceholderPackages = "{packages}"
	// PlaceholderFiles stands for the files the tier selected.
	PlaceholderFiles = "{files}"
	// PlaceholderImporters stands for the packages PlaceholderPackages names
	// and every package of the same module that imports one of them,
	// directly or through other packages, in its own files or its test
	// files.
	PlaceholderImporters = "{importers}"
)

// tierPlaceholders are the placeholders a tier's command may hold, in the
// order they are listed to the operator.
var tierPlaceholders = []string{PlaceholderPackages, PlaceholderFiles, PlaceholderImporters}

// tierPlaceholderList names tierPlaceholders in a sentence.
const tierPlaceholderList = PlaceholderPackages + ", " + PlaceholderFiles + ", or " + PlaceholderImporters

// The values of a tier's format.
const (
	// FormatGoTestJSON is go test -json's events.
	FormatGoTestJSON = "go-test-json"
	// FormatVitestJSON is vitest's JSON reporter.
	FormatVitestJSON = "vitest-json"
	// FormatGolangciJSON is golangci-lint's JSON output.
	FormatGolangciJSON = "golangci-json"
	// FormatGofmtList is gofmt -l's list of files not formatted.
	FormatGofmtList = "gofmt-list"
	// FormatPlain is any output; the exit status alone says whether the tier
	// passed. It is the default.
	FormatPlain = "plain"
)

// TestFormats are the values of a tier's format, in the order they are
// listed to the operator.
var TestFormats = []string{FormatGoTestJSON, FormatVitestJSON, FormatGolangciJSON, FormatGofmtList, FormatPlain}

// DefaultTestScript is the script the default tier runs when the manifest
// has no tests key.
const DefaultTestScript = "scripts/test.sh"

// placeholder is an argument that looks like a placeholder: a word in
// braces.
var placeholder = regexp.MustCompile(`^\{[A-Za-z_][A-Za-z0-9_]*\}$`)

// FormatOrDefault is the tier's format, or FormatPlain when it is empty. It
// does not say whether the format is one of TestFormats; TestErrors does.
func (t TestTier) FormatOrDefault() string {
	if s := strings.TrimSpace(t.Format); s != "" {
		return s
	}
	return FormatPlain
}

// TestTiers are the project's test tiers, cheapest first: the tests key's
// when it is there, and otherwise one plain tier named test that runs
// DefaultTestScript for every path, when fsys, the repository root, holds
// that file; none when it does not. It refuses tiers TestErrors finds wrong,
// with every problem.
func (m Manifest) TestTiers(fsys fs.FS) ([]TestTier, error) {
	if m.Tests != nil {
		if errs := m.TestErrors(); len(errs) > 0 {
			return nil, errors.New(strings.Join(errs, "; "))
		}
		return m.Tests, nil
	}
	if !isFile(fsys, DefaultTestScript) {
		return nil, nil
	}
	return []TestTier{{Name: "test", Command: []string{DefaultTestScript}, Paths: []string{"**"}, Format: FormatPlain}}, nil
}

// isFile reports whether name is a file in fsys; anything it cannot stat is
// not.
func isFile(fsys fs.FS, name string) bool {
	st, err := fs.Stat(fsys, name)
	return err == nil && !st.IsDir()
}

// TestErrors are what is wrong with the tests key, one sentence each, each
// beginning with the tier's index and field, such as tests[1].format; none
// when every tier is valid.
func (m Manifest) TestErrors() []string {
	var errs []string
	tiers := map[string]int{}
	for i, t := range m.Tests {
		if name := strings.TrimSpace(t.Name); name != "" {
			if _, ok := tiers[name]; !ok {
				tiers[name] = i
			}
		}
	}
	seen := map[string]int{}
	for i, t := range m.Tests {
		at := func(field string) string { return testField(i, t, field) }
		name := strings.TrimSpace(t.Name)
		if name == "" {
			errs = append(errs, at("name")+" is empty; give the tier a name, such as test or unit")
		} else if j, ok := seen[name]; ok {
			errs = append(errs, fmt.Sprintf("%s %q is tests[%d]'s name too; give each tier a name of its own", at("name"), t.Name, j))
		} else {
			seen[name] = i
		}
		if len(t.Command) == 0 {
			errs = append(errs, at("command")+" is empty; write the program and its arguments as a list, such as [scripts/test.sh] or [go, test, \"{packages}\"]")
		} else {
			errs = append(errs, commandProblems(at("command"), t.Command, true)...)
		}
		if len(t.AllCommand) > 0 {
			errs = append(errs, commandProblems(at("all_command"), t.AllCommand, false)...)
		}
		if len(t.FixCommand) > 0 {
			errs = append(errs, commandProblems(at("fix_command"), t.FixCommand, true)...)
		}
		if r := dirProblem(t.Dir); r != "" {
			errs = append(errs, fmt.Sprintf("%s %q %s", at("dir"), t.Dir, r))
		}
		selects := 0
		for _, p := range t.Paths {
			glob, out := strings.CutPrefix(p, "!")
			if r := patternProblem(glob); r != "" {
				errs = append(errs, fmt.Sprintf("%s pattern %q %s", at("paths"), p, r))
			} else if !out {
				selects++
			}
		}
		switch {
		case t.AllOnly || t.Manual || selects > 0:
		case len(t.Paths) == 0:
			errs = append(errs, at("paths")+" is empty; write the globs that select the tier, such as flai/**, or set all_only for a tier that runs only when every tier is asked for, or manual for one that runs only by its name")
		default:
			errs = append(errs, at("paths")+" has no pattern that selects a path, only ones beginning with ! or ones not valid; add one, such as flai/**")
		}
		if t.Manual && t.Queue {
			errs = append(errs, at("manual")+" is true with queue, but the merge queue's batch never runs a manual tier; remove queue, or manual for a tier the batch runs")
		}
		if t.Manual && len(t.FixCommand) > 0 {
			errs = append(errs, at("manual")+" is true with fix_command, but a fix runs for the paths a commit changes, which never select a manual tier; remove fix_command, or manual")
		}
		if t.Manual && t.Parallel {
			errs = append(errs, at("manual")+" is true with parallel, but a manual tier runs alone, by its name, with no tier to run beside; remove parallel, or manual")
		}
		if len(t.When) > 0 && !t.Queue {
			errs = append(errs, at("when")+" is set on a tier without queue: true; when only chooses when a queue tier also runs in the close-out, so set queue: true or remove when")
		}
		for _, p := range t.When {
			if r := patternProblem(strings.TrimPrefix(p, "!")); r != "" {
				errs = append(errs, fmt.Sprintf("%s pattern %q %s", at("when"), p, r))
			}
		}
		if f := t.FormatOrDefault(); !slices.Contains(TestFormats, f) {
			errs = append(errs, fmt.Sprintf("%s %q is not a format flai reads; write %s, or remove it for plain", at("format"), t.Format, strings.Join(TestFormats, ", ")))
		}
		named := map[string]bool{}
		for _, c := range t.Covers {
			covered := strings.TrimSpace(c)
			switch _, ok := tiers[covered]; {
			case named[covered]:
				errs = append(errs, fmt.Sprintf("%s names %q twice; name each tier it covers once", at("covers"), c))
			case !ok:
				errs = append(errs, fmt.Sprintf("%s %q names no tier; write the name of another tier, or remove it", at("covers"), c))
			case covered == name:
				errs = append(errs, fmt.Sprintf("%s names the tier itself; name only other tiers, or remove %q", at("covers"), c))
			}
			named[covered] = true
		}
	}
	for _, cycle := range coverCycles(m.Tests, tiers) {
		names := make([]string, len(cycle))
		for k, i := range cycle {
			names[k] = strconv.Quote(strings.TrimSpace(m.Tests[i].Name))
		}
		first := cycle[0]
		errs = append(errs, fmt.Sprintf("%s makes a cycle through the tiers %s, which cover one another; remove a name from one of their covers, so that no tier covers one that covers it", testField(first, m.Tests[first], "covers"), strings.Join(names, ", ")))
	}
	return errs
}

// testField is a tier's field as TestErrors names it, such as tests[1].format
// (tier "unit"), with the tier's name when it has one.
func testField(i int, t TestTier, field string) string {
	s := fmt.Sprintf("tests[%d].%s", i, field)
	if name := strings.TrimSpace(t.Name); name != "" && field != "name" {
		s += fmt.Sprintf(" (tier %q)", name)
	}
	return s
}

// TestWarnings are what a tests key may hold that runs but fails badly, one
// sentence each, named as TestErrors names a field: each command or
// all_command that runs go test with no timeout of its own, since go's
// 10-minute default panics naming only the package (I-0227). A command whose
// program is not go, such as a script, is not inspected.
func (m Manifest) TestWarnings() []string {
	var warns []string
	for i, t := range m.Tests {
		for _, c := range []struct {
			field string
			args  []string
		}{{"command", t.Command}, {"all_command", t.AllCommand}} {
			if runsGoTest(c.args) && !slices.ContainsFunc(c.args[2:], isTimeoutFlag) {
				warns = append(warns, testField(i, t, c.field)+` runs go test with no -timeout, so go's 10-minute default panics naming only the package when a slow host reaches it; add "-timeout" and a duration above the slowest package's, such as "20m"`)
			}
		}
	}
	return warns
}

// runsGoTest reports whether args run go test: the program is go, by name or
// by a path ending in it, and its first argument is test.
func runsGoTest(args []string) bool {
	return len(args) >= 2 && filepath.Base(strings.TrimSpace(args[0])) == "go" && args[1] == "test"
}

// isTimeoutFlag reports whether a is go test's timeout flag, with one dash or
// two and with its value after = or not: -timeout, --timeout=20m, or
// -test.timeout as the test binary takes it.
func isTimeoutFlag(a string) bool {
	if !strings.HasPrefix(a, "-") {
		return false
	}
	name, _, _ := strings.Cut(strings.TrimPrefix(a[1:], "-"), "=")
	return name == "timeout" || name == "test.timeout"
}

// coverCycles are the groups of tiers that cover one another through covers,
// each the tiers' indexes in list order and each group once; tiers maps a
// name to the index of the first tier with it. A tier naming itself is not a
// cycle here, as TestErrors refuses it on its own.
func coverCycles(list []TestTier, tiers map[string]int) [][]int {
	covers := make([][]int, len(list))
	for i, t := range list {
		for _, c := range t.Covers {
			if j, ok := tiers[strings.TrimSpace(c)]; ok && j != i {
				covers[i] = append(covers[i], j)
			}
		}
	}
	reaches := make([][]bool, len(list))
	for i := range list {
		reaches[i] = make([]bool, len(list))
		next := slices.Clone(covers[i])
		for len(next) > 0 {
			j := next[len(next)-1]
			next = next[:len(next)-1]
			if !reaches[i][j] {
				reaches[i][j] = true
				next = append(next, covers[j]...)
			}
		}
	}
	var cycles [][]int
	reported := make([]bool, len(list))
	for i := range list {
		if reported[i] || !reaches[i][i] {
			continue
		}
		var cycle []int
		for j := range list {
			if reaches[i][j] && reaches[j][i] {
				cycle = append(cycle, j)
				reported[j] = true
			}
		}
		cycles = append(cycles, cycle)
	}
	return cycles
}

// commandProblems are what is wrong with a tier's argument list, each
// sentence beginning with field. A placeholder is allowed only where
// placeholders is set, and only as an argument of its own.
func commandProblems(field string, args []string, placeholders bool) []string {
	var errs []string
	if strings.TrimSpace(args[0]) == "" {
		errs = append(errs, field+" has no program first; write the program, then its arguments, such as [scripts/test.sh]")
	} else if placeholder.MatchString(args[0]) {
		errs = append(errs, fmt.Sprintf("%s begins with the placeholder %s; write the program first, then its arguments", field, args[0]))
	}
	for _, a := range args[1:] {
		switch {
		case slices.Contains(tierPlaceholders, a) && !placeholders:
			errs = append(errs, fmt.Sprintf("%s argument %s is a placeholder, and all_command runs for no paths; write it without %s", field, a, tierPlaceholderList))
		case slices.Contains(tierPlaceholders, a):
		case placeholder.MatchString(a):
			errs = append(errs, fmt.Sprintf("%s argument %s is not a placeholder flai fills; write %s, or remove it", field, a, tierPlaceholderList))
		case slices.ContainsFunc(tierPlaceholders, func(p string) bool { return strings.Contains(a, p) }):
			errs = append(errs, fmt.Sprintf("%s argument %q holds a placeholder inside it, which flai does not fill; write %s as an argument of its own", field, a, tierPlaceholderList))
		}
	}
	return errs
}

// dirProblem says what is wrong with a tier's dir, or "" when it is valid: a
// clean path inside the repository, relative to its root, or empty.
func dirProblem(d string) string {
	return repoPathProblem(d, "a folder", "flai", "flai/web")
}

// repoPathProblem says what is wrong with p, a path the manifest names, or ""
// when it is valid: a clean path inside the repository, relative to its root,
// or empty. what is what p names, such as a folder, and example and clean
// are paths the reason gives as ones to write.
func repoPathProblem(p, what, example, clean string) string {
	return pathProblem(p, "the repository", "the repository root", what, example, clean)
}

// pathProblem says what is wrong with p, a path the manifest names inside
// within, such as the repository, whose root is root, such as the repository
// root, or "" when it is valid: a clean path inside it, relative to its root,
// or empty. what, example, and clean are as repoPathProblem's.
func pathProblem(p, within, root, what, example, clean string) string {
	switch {
	case p == "":
		return ""
	case isAbsolute(p):
		return "is an absolute path; write it relative to " + root + ", such as " + example
	case path.Clean(p) == ".." || strings.HasPrefix(path.Clean(p), "../"):
		return "leaves " + within + "; write " + what + " inside it, relative to its root, such as " + example
	case path.Clean(p) != p:
		return "is not a clean path; write it with single slashes, no . or .. segments, and none at the end, such as " + clean
	}
	return ""
}

// isAbsolute reports whether p is an absolute path, on any system: one
// beginning with a slash or a backslash, or with a drive letter.
func isAbsolute(p string) bool {
	return strings.HasPrefix(p, "/") || strings.HasPrefix(p, `\`) || len(p) > 1 && p[1] == ':'
}
