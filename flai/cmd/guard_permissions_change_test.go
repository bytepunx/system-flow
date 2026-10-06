package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0229: flai guard reads orchestration.permissions from the manifest at
// each call, so a permission the operator turns on or off while the
// orchestrator runs holds its next call: off to on lets the next call
// through, and on to off refuses it.
func TestGuardHoldsARunningOrchestratorToChangedPermissions(t *testing.T) {
	t.Setenv("FLAI_ROLE", "orchestrate")
	root := tempProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	base, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	setPromote := func(on string) {
		t.Helper()
		data := append(append([]byte{}, base...), "orchestration:\n  permissions:\n    promote_to_ready: "+on+"\n"...)
		if err := os.WriteFile(manifest, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	const promote = `{"tool_name":"mcp__flai__item_move","tool_input":{"id":"S-1","to":"ready"}}`
	for _, step := range []struct {
		name, on string
		code     int
		err      string
	}{
		{"off", "false", 2, "it needs orchestration.permissions.promote_to_ready, which is off"},
		{"turned on", "true", 0, ""},
		{"turned off", "false", 2, "it needs orchestration.permissions.promote_to_ready, which is off"},
	} {
		setPromote(step.on)
		_, errOut, code := runStdin(t, root, promote, "guard")
		if code != step.code || (step.err == "" && errOut != "") || !strings.Contains(errOut, step.err) {
			t.Errorf("promote_to_ready %s: code %d, stderr %q", step.name, code, errOut)
		}
	}
}
