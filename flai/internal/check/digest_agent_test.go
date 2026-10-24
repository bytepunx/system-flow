package check

import (
	"strings"
	"testing"
)

// ADR-0196: a digest.agent that is not valid is an error on the digest key,
// saying what to write; a valid block, and no block, are not reported.
func TestDigestAgentIsChecked(t *testing.T) {
	for _, extra := range []string{
		"",
		"digest:\n  agent:\n    model: claude-sonnet-5\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.digest"); len(got) != 0 {
			t.Errorf("a valid digest block %q: %+v", extra, got)
		}
	}
	got := findings(t, planningProject(t, "digest:\n  agent:\n    harness: Claude Code\n"), "manifest.digest")
	if len(got) != 1 {
		t.Fatalf("want one error, the agent: %+v", got)
	}
	f := got[0]
	if want := `digest.agent harness "Claude Code" is not a name such as claude-code`; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, want) {
		t.Errorf("finding: %+v, want an error on line 8 saying %q", f, want)
	}
}
