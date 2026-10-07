package verify

import "testing"

func TestTextIsALineForEachTierAndFinding(t *testing.T) {
	r := Result{Tiers: []TierResult{
		{Name: "vet", State: Passed, Duration: "1.5s"},
		{Name: "test", State: Failed, Duration: "42ms", Omitted: 3, Findings: []Finding{
			{Name: "TestAdd", Path: "fx/calc/calc_test.go", Line: 7, Message: "Add(2, 3) = -1, want 5\nthe sum is wrong"},
			{Name: "gofmt", Path: "fx/a.go", Message: "not gofmt-formatted"},
			{Name: "smoke", Message: "exit status 2"},
			{Path: "fx/b.go"},
		}},
		{Name: "lint", State: NotReached},
	}}
	want := `passed vet (1.5s)
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
