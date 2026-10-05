package check

import (
	"strings"
	"testing"
)

// S-0217: each thing wrong with the manifest's orchestration block is an
// error on its key; a valid block, and no block, are not reported.
func TestOrchestrationSettingsAreChecked(t *testing.T) {
	for _, extra := range []string{
		"",
		"orchestration:\n  policy: wsjf\n  release:\n    policy: threshold\n    value: 500\n    count: 5\n",
		"orchestration:\n  release:\n    policy: theme\n    tag: cli\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.orchestration"); len(got) != 0 {
			t.Errorf("a valid orchestration block %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, "orchestration:\n  policy: lifo\n  release:\n    policy: threshold\n    count: -1\n")
	got := findings(t, repo, "manifest.orchestration")
	if len(got) != 2 {
		t.Fatalf("want two errors, policy and release.count: %+v", got)
	}
	for i, want := range []string{`orchestration.policy "lifo"`, "orchestration.release.count -1"} {
		f := got[i]
		if f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, want) {
			t.Errorf("finding %d: %+v, want an error on line 8 saying %q", i, f, want)
		}
	}
	repo = planningProject(t, "orchestration:\n  release:\n    policy: theme\n    epic: E-0001\n    tag: cli\n")
	if got := findings(t, repo, "manifest.orchestration"); len(got) != 1 || !strings.Contains(got[0].Message, "both epic and tag") {
		t.Errorf("want one error, a theme with both epic and tag: %+v", got)
	}
}

// S-0218: a permission key that names none, an answer_threads outside its
// values, and an orchestration.agent that is not valid are errors on the
// orchestration key; valid permissions and agent are not reported.
func TestOrchestrationPermissionsAndAgentAreChecked(t *testing.T) {
	valid := "orchestration:\n  permissions:\n    promote_to_ready: true\n    answer_threads: recommend\n  agent:\n    model: claude-sonnet-5\n"
	if got := findings(t, planningProject(t, valid), "manifest.orchestration"); len(got) != 0 {
		t.Errorf("valid permissions and agent: %+v", got)
	}
	repo := planningProject(t, "orchestration:\n  permissions:\n    promote: true\n    answer_threads: always\n  agent:\n    harness: Claude Code\n")
	got := findings(t, repo, "manifest.orchestration")
	if len(got) != 3 {
		t.Fatalf("want three errors, the unknown key, answer_threads, and the agent: %+v", got)
	}
	for i, want := range []string{`orchestration.permissions has no permission "promote"`, `orchestration.permissions.answer_threads "always"`, `orchestration.agent harness "Claude Code"`} {
		f := got[i]
		if f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, want) {
			t.Errorf("finding %d: %+v, want an error on line 8 saying %q", i, f, want)
		}
	}
}
