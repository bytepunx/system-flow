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
