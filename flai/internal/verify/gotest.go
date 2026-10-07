package verify

import (
	"cmp"
	"encoding/json"
	"io/fs"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// goEvent is the part of a go test -json event read here.
type goEvent struct {
	ImportPath  string
	Action      string
	Package     string
	Test        string
	Output      string
	FailedBuild string
}

var (
	// compileLine is a compiler's or vet's file:line[:col]: message.
	compileLine = regexp.MustCompile(`^(\S+\.go):(\d+)(?::\d+)?: (.+)$`)
	// testLine is a test's failure, file.go:line: message, indented.
	testLine = regexp.MustCompile(`^(\s*)([^\s:]+\.go):(\d+): ?(.*)$`)
	// frameLine is a stack frame's /path/file.go:line +0x...
	frameLine = regexp.MustCompile(`^(\S+\.go):(\d+)(?: \+0x[0-9a-f]+)?$`)
)

// goTest gathers go test -json events into findings, in the order the
// failures happened.
type goTest struct {
	e       env
	module  string
	got     []Finding
	built   map[string]bool     // packages with a build finding
	output  map[string][]string // a test's or a package's output lines
	failed  map[string][]string // each package's failed tests
	current string              // the package whose build output is being read
	last    int                 // the build finding a continuation line adds to, or -1
}

// parseGoTest is a finding for each test go test -json reports failed, by
// the file and line of its first failure, for each package that failed to
// build, by its compiler errors, and for each package that failed with no
// test failing.
func parseGoTest(stdout, stderr string, e env) []Finding {
	g := &goTest{
		e: e, module: modulePath(e), built: map[string]bool{},
		output: map[string][]string{}, failed: map[string][]string{}, last: -1,
	}
	if !strings.Contains(stdout, `"Action":"build-output"`) {
		// a go before build output came as events printed it on stderr
		for _, line := range strings.Split(stderr, "\n") {
			g.buildOutput("", line)
		}
	}
	for _, line := range strings.Split(stdout, "\n") {
		var ev goEvent
		if !strings.HasPrefix(line, "{") || json.Unmarshal([]byte(line), &ev) != nil {
			continue
		}
		g.event(ev)
	}
	return g.got
}

// event takes one event in.
func (g *goTest) event(ev goEvent) {
	key := ev.Package + "\x00" + ev.Test
	switch ev.Action {
	case "build-output":
		g.buildOutput(ev.ImportPath, ev.Output)
	case "output":
		g.output[key] = append(g.output[key], strings.TrimRight(ev.Output, "\n"))
	case "fail":
		switch {
		case ev.Test != "":
			g.testFailed(ev.Package, ev.Test, g.output[key])
		case ev.FailedBuild != "" || buildFailed(g.output[key]):
			if pkg := cmp.Or(ev.FailedBuild, ev.Package); !g.built[bare(pkg)] {
				dir, _ := g.packageDir(pkg)
				g.got = append(g.got, Finding{Name: "build", Path: dir, Message: packageMessage(g.output[key])})
			}
		case len(g.failed[ev.Package]) == 0:
			dir, _ := g.packageDir(ev.Package)
			g.got = append(g.got, Finding{Name: ev.Package, Path: dir, Message: packageMessage(g.output[key])})
		}
	}
}

// buildOutput takes a line the build printed in: a file:line: message is a
// finding, and an indented line after one adds to its message.
func (g *goTest) buildOutput(importPath, line string) {
	line = strings.TrimRight(line, "\n")
	if pkg, ok := strings.CutPrefix(line, "# "); ok {
		g.current = bare(pkg)
		return
	}
	if importPath == "" {
		importPath = g.current
	}
	if strings.HasPrefix(line, "\t") && g.last >= 0 {
		f := &g.got[g.last]
		if strings.Count(f.Message, "\n")+1 < messageLines {
			f.Message += "\n" + strings.TrimSpace(line)
		}
		return
	}
	m := compileLine.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		g.last = -1
		return
	}
	n, _ := strconv.Atoi(m[2])
	g.got = append(g.got, Finding{Name: "build", Path: g.e.rel(m[1]), Line: n, Message: m[3]})
	g.last = len(g.got) - 1
	if importPath != "" {
		g.built[bare(importPath)] = true
	}
}

// buildFailed reports whether a package's output says its build failed.
func buildFailed(lines []string) bool {
	for _, l := range lines {
		if strings.Contains(l, "[build failed]") {
			return true
		}
	}
	return false
}

// testFailed takes a failed test in, unless one of its subtests failed and
// so already speaks for it.
func (g *goTest) testFailed(pkg, test string, lines []string) {
	g.last = -1
	for _, t := range g.failed[pkg] {
		if strings.HasPrefix(t, test+"/") {
			g.failed[pkg] = append(g.failed[pkg], test)
			return
		}
	}
	g.failed[pkg] = append(g.failed[pkg], test)
	dir, known := g.packageDir(pkg)
	f := Finding{Name: test, Path: dir}
	if i, m := firstMatch(testLine, lines); m != nil {
		f.Line, _ = strconv.Atoi(m[3])
		switch file := path.Join(dir, m[2]); {
		case !known:
			f.Path = m[2]
		case g.exists(file):
			f.Path = file
		default:
			// not the package's own file, as testing.go is when the race
			// detector fails a test: the stack, if any, says where
			f.Path, f.Line = g.frame(lines)
			if f.Path == "" {
				f.Path = dir
			}
		}
		f.Message = clip(strings.Join(append([]string{m[4]}, continuation(lines[i+1:], len(m[1]))...), "\n"))
	} else if p := panicLine(lines); p != "" {
		f.Message = p
		if loc, line := g.frame(lines); loc != "" {
			f.Path, f.Line = loc, line
		}
	} else {
		f.Message = packageMessage(lines)
	}
	g.got = append(g.got, f)
}

// exists reports whether the checkout holds the root-relative file.
func (g *goTest) exists(file string) bool {
	_, err := fs.Stat(g.e.fsys, file)
	return err == nil
}

// frame is the first stack frame in lines that lies below the root, and its
// line.
func (g *goTest) frame(lines []string) (string, int) {
	for _, l := range lines {
		m := frameLine.FindStringSubmatch(strings.TrimSpace(l))
		if m == nil {
			continue
		}
		if p := g.e.rel(m[1]); !path.IsAbs(p) {
			n, _ := strconv.Atoi(m[2])
			return p, n
		}
	}
	return "", 0
}

// packageDir is the directory, relative to the root, of the package with
// the import path, and whether the module's go.mod says it.
func (g *goTest) packageDir(importPath string) (string, bool) {
	p := bare(importPath)
	switch {
	case g.module == "":
		return "", false
	case p == g.module:
		return g.e.dir, true
	}
	if rest, ok := strings.CutPrefix(p, g.module+"/"); ok {
		return path.Join(g.e.dir, rest), true
	}
	return "", false
}

// modulePath is the module the go.mod in the tier's Dir declares, "" when
// there is none.
func modulePath(e env) string {
	data, err := fs.ReadFile(e.fsys, path.Join(e.dir, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if f := strings.Fields(line); len(f) >= 2 && f[0] == "module" {
			return strings.Trim(f[1], `"`)
		}
	}
	return ""
}

// bare is an import path without the test variant go names after it, as
// in "x/y [x/y.test]".
func bare(importPath string) string {
	p, _, _ := strings.Cut(importPath, " ")
	return p
}

// firstMatch is the index and submatches of the first line re matches.
func firstMatch(re *regexp.Regexp, lines []string) (int, []string) {
	for i, l := range lines {
		if m := re.FindStringSubmatch(l); m != nil {
			return i, m
		}
	}
	return -1, nil
}

// continuation are the lines after a failure's first that go on with its
// message: those indented further than it, up to the next of go test's own.
func continuation(lines []string, indent int) []string {
	var out []string
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t == "" || len(l)-len(strings.TrimLeft(l, " \t")) <= indent || framing(t) {
			break
		}
		out = append(out, t)
	}
	return out
}

// panicLine is the panic a test's output reports, "" when it reports none.
func panicLine(lines []string) string {
	for _, l := range lines {
		if t := strings.TrimSpace(l); strings.HasPrefix(t, "panic: ") {
			return t
		}
	}
	return ""
}

// packageMessage is a test's or a package's output without go test's own
// lines, its last few.
func packageMessage(lines []string) string {
	var keep []string
	for _, l := range lines {
		if t := strings.TrimSpace(l); t != "" && !framing(t) && t != "FAIL" && !strings.HasPrefix(t, "FAIL\t") {
			keep = append(keep, t)
		}
	}
	return lastLines(strings.Join(keep, "\n"), messageLines)
}

// framing reports whether a trimmed line is go test's own, not the test's.
func framing(t string) bool {
	for _, p := range []string{"=== RUN", "=== PAUSE", "=== CONT", "=== NAME", "--- FAIL", "--- PASS", "--- SKIP"} {
		if strings.HasPrefix(t, p) {
			return true
		}
	}
	return false
}
