package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/formulas"
)

// S-0811, ADR-0312: flai release --pending exits exitReleaseSuiteFailed when
// a suite verify_at.publish lists fails at the main branch's head. That exit
// is a failure of the publish formula's publish step, not a clean stop
// (ADR-0271, ADR-0239 §3): the step's stop still lists 2 alone, so its
// instance's agent blocks the task and asks, and the formula's body says so.
func TestPublishFormulaFailsOnAFailedSuite(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile(publishFormulaFile)
	if err != nil {
		t.Fatal(err)
	}
	f, problem := formulas.Parse("design/workflows/formulas/publish.md", string(content))
	if problem != nil {
		t.Fatalf("the template's publish formula does not parse: %s", problem)
	}
	var publish *formulas.Step
	for i := range f.Steps {
		if f.Steps[i].ID == "publish" {
			publish = &f.Steps[i]
		}
	}
	if publish == nil {
		t.Fatalf("the publish formula has no publish step: %+v", f.Steps)
	}
	if !strings.Contains(publish.Run.Flai, "release --pending") {
		t.Errorf("the publish step runs flai release --pending: %+v", publish.Run)
	}
	if publish.Stops(exitReleaseSuiteFailed) {
		t.Errorf("the publish step's stop %v lists %d, so a failed suite would stop the run cleanly", publish.Stop, exitReleaseSuiteFailed)
	}
	if !publish.Stops(2) {
		t.Errorf("the publish step's stop %v still lists 2, nothing pending", publish.Stop)
	}
	for _, want := range []string{"exits 5", "`verify_at.publish`", "block its task", "the suite, the tier, and the log path"} {
		if !strings.Contains(f.Body, want) {
			t.Errorf("the publish formula's body says %q:\n%s", want, f.Body)
		}
	}
}
