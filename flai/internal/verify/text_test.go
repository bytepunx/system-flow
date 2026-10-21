package verify

import "testing"

func TestTextIsALineForEachTierAndFinding(t *testing.T) {
	r := Result{Tiers: []TierResult{
		{Name: "vet", State: Passed, Duration: "1.5s"},
		{Name: "short", State: Skipped, CoveredBy: "test"},
		{Name: "test", State: Failed, Duration: "42ms", Omitted: 3, Findings: []Finding{
			{Name: "TestAdd", Path: "fx/calc/calc_test.go", Line: 7, Message: "Add(2, 3) = -1, want 5\nthe sum is wrong"},
			{Name: "gofmt", Path: "fx/a.go", Message: "not gofmt-formatted"},
			{Name: "smoke", Message: "exit status 2"},
			{Path: "fx/b.go"},
		}},
		{Name: "lint", State: NotReached},
	}}
	want := `passed vet (1.5s)
skipped short: covered by test
failed test (42ms)
fx/calc/calc_test.go:7 TestAdd: Add(2, 3) = -1, want 5
    the sum is wrong
fx/a.go gofmt: not gofmt-formatted
smoke: exit status 2
fx/b.go
… 3 more findings left out
not-reached lint
`
	if got := r.Text(); got != want {
		t.Errorf("text\n%s\nwant\n%s", got, want)
	}
	if got := (Result{Passed: true}).Text(); got != "no tier selects these paths\n" {
		t.Errorf("with no tier: %q", got)
	}
}

func TestStepFindingsIndentsEveryLine(t *testing.T) {
	s := Step{Name: "test", Omitted: 2, Findings: []Finding{
		{Name: "go test ./...", Message: "--- FAIL: TestAdd (0.00s)\n    calc_test.go:7: wrong sum\nFAIL\nexit status 1"},
		{Name: "gofmt", Path: "fx/a.go", Message: "not gofmt-formatted"},
	}}
	want := `    go test ./...: --- FAIL: TestAdd (0.00s)
            calc_test.go:7: wrong sum
        FAIL
        exit status 1
    fx/a.go gofmt: not gofmt-formatted
    … 2 more findings left out
`
	if got := StepFindings(s, "    "); got != want {
		t.Errorf("indented\n%s\nwant\n%s", got, want)
	}
	unindented := "fx/a.go gofmt: not gofmt-formatted\n… 2 more findings left out\n"
	if got := StepFindings(Step{Omitted: 2, Findings: s.Findings[1:]}, ""); got != unindented {
		t.Errorf("with no indent %q, want %q", got, unindented)
	}
	if got := StepFindings(Step{Name: "vet"}, "    "); got != "" {
		t.Errorf("with no findings %q", got)
	}
}
