package check

import (
	"strings"
	"testing"
)

// S-0273: a test tier that is not valid is an error on the tests key, naming
// its index, its field, and what to write; valid tiers, and no tests key, are
// not reported.
func TestTestTiersAreChecked(t *testing.T) {
	for _, extra := range []string{
		"",
		"tests: []\n",
		"tests:\n  - name: unit\n    command: [go, test, \"{packages}\"]\n    dir: flai\n    paths: [\"flai/**\"]\n    format: go-test-json\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.tests"); len(got) != 0 {
			t.Errorf("valid tests %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, "tests:\n  - name: unit\n    command: [go, test]\n    paths: [\"flai/**\"]\n    format: junit\n")
	got := findings(t, repo, "manifest.tests")
	if len(got) != 1 {
		t.Fatalf("want one error, the format: %+v", got)
	}
	if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, `tests[0].format (tier "unit") "junit" is not a format flai reads`) {
		t.Errorf("finding: %+v, want an error on line 8 naming tests[0].format", f)
	}
}

// S-0548, I-0227: a tier that runs go test with no -timeout is a warning on
// the tests key, which --strict fails on, naming the tier and saying to add
// one; a tier that names its timeout, and one that runs a script, are not.
func TestGoTestTiersWithoutATimeoutAreWarnedOf(t *testing.T) {
	for _, extra := range []string{
		"tests:\n  - name: go-test\n    command: [go, test, \"-timeout\", \"20m\", \"{packages}\"]\n    paths: [\"flai/**\"]\n",
		"tests:\n  - name: integration\n    command: [scripts/integration.sh]\n    all_only: true\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.tests-timeout"); len(got) != 0 {
			t.Errorf("tests %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, "tests:\n  - name: go-test\n    command: [go, test, \"-race\", \"{packages}\"]\n    paths: [\"flai/**\"]\n")
	got := findings(t, repo, "manifest.tests-timeout")
	if len(got) != 1 {
		t.Fatalf("want one warning: %+v", got)
	}
	if f := got[0]; f.Level != Warning || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, `tests[0].command (tier "go-test") runs go test with no -timeout`) || !strings.Contains(f.Message, `add "-timeout"`) {
		t.Errorf("finding: %+v, want a warning on line 8 naming tests[0].command and saying to add -timeout", f)
	}
	if got := findings(t, repo, "manifest.tests"); len(got) != 0 {
		t.Errorf("a tier without a timeout is valid, not an error: %+v", got)
	}
	if got[0].advisory {
		t.Errorf("the warning is advisory, so --strict passes over it: %+v", got[0])
	}
}
