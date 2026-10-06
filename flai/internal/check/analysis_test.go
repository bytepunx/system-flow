package check

import (
	"strings"
	"testing"
)

// S-0223: an analysis.schedule that does not parse and an analysis.agent that
// is not valid are errors on the analysis key, each saying what to write; a
// valid block, and no block, are not reported.
func TestAnalysisSettingsAreChecked(t *testing.T) {
	for _, extra := range []string{
		"",
		"analysis:\n  schedule: daily\n",
		"analysis:\n  schedule: \"0 6 * * 1\"\n  agent:\n    model: claude-sonnet-5\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.analysis"); len(got) != 0 {
			t.Errorf("a valid analysis block %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, "analysis:\n  schedule: weekly\n  agent:\n    harness: Claude Code\n")
	got := findings(t, repo, "manifest.analysis")
	if len(got) != 2 {
		t.Fatalf("want two errors, the schedule and the agent: %+v", got)
	}
	for i, want := range []string{`analysis.schedule "weekly" has 1 fields, not five (minute, hour, day of month, month, day of week); write one such as`, `analysis.agent harness "Claude Code" is not a name such as claude-code`} {
		f := got[i]
		if f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, want) {
			t.Errorf("finding %d: %+v, want an error on line 8 saying %q", i, f, want)
		}
	}
	repo = planningProject(t, "analysis:\n  schedule: \"0 0 31 2 *\"\n")
	if got := findings(t, repo, "manifest.analysis"); len(got) != 1 || !strings.Contains(got[0].Message, "never comes round") {
		t.Errorf("want one error, a schedule that never comes round: %+v", got)
	}
}
