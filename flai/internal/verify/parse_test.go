package verify

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

// The fixtures under testdata are the tools' own output for a small failing
// module, recorded with its temporary folder renamed /home/dev/proj/fx:
// go test -json ./..., golangci-lint run --output.json.path stdout, and
// gofmt -l ., each run in fx. vitest.txt follows vitest 4's JSON reporter
// for a project at /home/dev/proj/flaiover.
const fixtureRoot = "/home/dev/proj"

func fixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func fixtureFS() fstest.MapFS {
	return fstest.MapFS{
		"fx/go.mod":            {Data: []byte("module example.com/fx\n\ngo 1.26\n")},
		"fx/calc/calc_test.go": {},
		"fx/boom/boom_test.go": {},
	}
}

func TestEachFormatTurnsAFailingOutputIntoFindings(t *testing.T) {
	cases := []struct {
		format, dir, fixture string
		want                 []Finding
	}{
		{FormatGoTestJSON, "fx", "go-test.json", []Finding{
			{Name: "build", Path: "fx/broken/broken.go", Line: 4, Message: "undefined: undefinedThing"},
			{Name: "TestAdd", Path: "fx/calc/calc_test.go", Line: 7, Message: "Add(2, 3) = -1, want 5\nthe sum is wrong\nfor small numbers"},
			{Name: "TestTable/case", Path: "fx/calc/calc_test.go", Line: 15, Message: "Add(1, 1) = 0, want 2"},
			{Name: "TestPanics", Path: "fx/boom/boom_test.go", Line: 7, Message: "panic: assignment to entry in nil map [recovered, repanicked]"},
		}},
		{FormatGolangciJSON, "fx", "golangci.json", []Finding{
			{Name: "errcheck", Path: "fx/lint/lint.go", Line: 7, Message: "Error return value of `os.Remove` is not checked"},
			{Name: "errcheck", Path: "fx/lint/more.go", Line: 7, Message: "Error return value of `os.WriteFile` is not checked"},
			{Name: "ineffassign", Path: "fx/lint/lint.go", Line: 12, Message: "ineffectual assignment to x"},
			{Name: "ineffassign", Path: "fx/lint/more.go", Line: 9, Message: "ineffectual assignment to s"},
		}},
		{FormatGofmtList, "fx", "gofmt.txt", []Finding{
			{Name: "gofmt", Path: "fx/lint/lint.go", Message: "not gofmt-formatted"},
			{Name: "gofmt", Path: "fx/pass/ugly.go", Message: "not gofmt-formatted"},
		}},
		{FormatVitestJSON, "flaiover", "vitest.txt", []Finding{
			{Name: "sum adds small numbers", Path: "flaiover/src/lib/sum.test.ts", Line: 5, Message: "AssertionError: expected 5 to be 4 // Object.is equality"},
			{Name: "sum negative subtracts", Path: "flaiover/src/lib/sum.test.ts", Line: 13, Message: "AssertionError: expected -1 to deeply equal 1"},
			{Name: "suite", Path: "flaiover/src/lib/broken.test.ts", Message: `Failed to resolve import "./missing" from "src/lib/broken.test.ts". Does the file exist?`},
		}},
	}
	for _, c := range cases {
		t.Run(c.format, func(t *testing.T) {
			out := fixture(t, c.fixture)
			tier := Tier{Name: c.format, Dir: c.dir, Format: c.format}
			got := parse(tier, output{stdout: out, combined: out, exit: 1}, env{root: fixtureRoot, dir: c.dir, fsys: fixtureFS()})
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("findings\n%s\nwant\n%s", show(got), show(c.want))
			}

			// through a run, capped at two with the rest counted
			proc := &fakeProc{steps: map[string]step{"tool": {stdout: out, exit: 1}}}
			res, err := Run(context.Background(), fixtureRoot, []Selected{{Tier: tier, Argv: []string{"tool"}}},
				RunOptions{Max: 2, FS: fixtureFS(), Proc: proc})
			if err != nil || res.Passed || len(res.Tiers) != 1 {
				t.Fatalf("run %+v, %v", res, err)
			}
			tr := res.Tiers[0]
			if !reflect.DeepEqual(tr.Findings, c.want[:2]) || tr.Omitted != len(c.want)-2 {
				t.Errorf("capped findings %s, omitted %d; want the first two and %d", show(tr.Findings), tr.Omitted, len(c.want)-2)
			}
		})
	}
}

func TestPlainIsTheLastLinesAndTheExitStatus(t *testing.T) {
	var lines []string
	for i := 1; i <= 30; i++ {
		lines = append(lines, fmt.Sprintf("line %d", i))
	}
	got := parse(Tier{Name: "smoke"}, output{combined: strings.Join(lines, "\n") + "\n\n", exit: 2}, env{})
	want := []Finding{{Name: "smoke", Message: strings.Join(lines[10:], "\n") + "\nexit status 2"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s", show(got))
	}

	got = parse(Tier{Name: "gofmt", Dir: "fx", Format: FormatGofmtList},
		output{stdout: "a.go\n", stderr: "b.go:3:1: expected declaration, found x\n", exit: 2}, env{dir: "fx"})
	want = []Finding{
		{Name: "gofmt", Path: "fx/b.go", Line: 3, Message: "expected declaration, found x"},
		{Name: "gofmt", Path: "fx/a.go", Message: "not gofmt-formatted"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("gofmt that cannot parse a file: %s", show(got))
	}

	got = parse(Tier{Name: "lint", Format: FormatGolangciJSON}, output{err: errors.New(`exec: "golangci-lint": executable file not found in $PATH`)}, env{})
	if len(got) != 1 || got[0].Message != `exec: "golangci-lint": executable file not found in $PATH` {
		t.Errorf("a tier that cannot start: %s", show(got))
	}
}

func TestAFormatThatFindsNothingFallsBackToPlain(t *testing.T) {
	for _, format := range []string{FormatGoTestJSON, FormatVitestJSON, FormatGolangciJSON, FormatGofmtList} {
		stdout := "{not json\n"
		if format == FormatGofmtList {
			stdout = ""
		}
		out := output{stdout: stdout, stderr: "go: cannot find main module\n", combined: "go: cannot find main module\n", exit: 1}
		got := parse(Tier{Name: "t", Format: format}, out, env{fsys: fstest.MapFS{}})
		want := []Finding{{Name: "t", Message: "go: cannot find main module\nexit status 1"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %s", format, show(got))
		}
	}
}

func TestGoTestReportsAPackageThatFailsWithNoTestFailing(t *testing.T) {
	out := strings.Join([]string{
		`{"Action":"start","Package":"example.com/fx/calc"}`,
		`{"Action":"output","Package":"example.com/fx/calc","Output":"TestMain exited early\n"}`,
		`{"Action":"output","Package":"example.com/fx/calc","Output":"FAIL\texample.com/fx/calc\t0.01s\n"}`,
		`{"Action":"fail","Package":"example.com/fx/calc","Elapsed":0.01}`,
	}, "\n")
	got := parseGoTest(out, "", env{dir: "fx", fsys: fixtureFS()})
	want := []Finding{{Name: "example.com/fx/calc", Path: "fx/calc", Message: "TestMain exited early"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s", show(got))
	}
}

func TestGoTestReadsBuildErrorsFromStderrWhenAnOlderGoPrintsThemThere(t *testing.T) {
	stdout := `{"Action":"output","Package":"example.com/fx/broken","Output":"FAIL\texample.com/fx/broken [build failed]\n"}
{"Action":"fail","Package":"example.com/fx/broken"}`
	stderr := "# example.com/fx/broken\nbroken/broken.go:4:27: cannot use x (variable of type int) as string value\n\thave int\n"
	got := parseGoTest(stdout, stderr, env{dir: "fx", fsys: fixtureFS()})
	want := []Finding{{Name: "build", Path: "fx/broken/broken.go", Line: 4, Message: "cannot use x (variable of type int) as string value\nhave int"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s", show(got))
	}
}

func TestGoTestPlacesARaceByItsStackNotByTesting(t *testing.T) {
	out := strings.Join([]string{
		`{"Action":"output","Package":"example.com/fx/calc","Test":"TestAdd","Output":"==================\n"}`,
		`{"Action":"output","Package":"example.com/fx/calc","Test":"TestAdd","Output":"WARNING: DATA RACE\n"}`,
		`{"Action":"output","Package":"example.com/fx/calc","Test":"TestAdd","Output":"  example.com/fx/calc.TestAdd.func1()\n"}`,
		`{"Action":"output","Package":"example.com/fx/calc","Test":"TestAdd","Output":"      /home/dev/proj/fx/calc/calc_test.go:12 +0x44\n"}`,
		`{"Action":"output","Package":"example.com/fx/calc","Test":"TestAdd","Output":"    testing.go:1617: race detected during execution of test\n"}`,
		`{"Action":"fail","Package":"example.com/fx/calc","Test":"TestAdd"}`,
	}, "\n")
	got := parseGoTest(out, "", env{root: fixtureRoot, dir: "fx", fsys: fixtureFS()})
	want := []Finding{{Name: "TestAdd", Path: "fx/calc/calc_test.go", Line: 12, Message: "race detected during execution of test"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s", show(got))
	}
}

func show(fs []Finding) string {
	var b strings.Builder
	for _, f := range fs {
		fmt.Fprintf(&b, "  %#v\n", f)
	}
	return b.String()
}
