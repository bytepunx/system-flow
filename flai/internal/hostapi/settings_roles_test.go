package hostapi

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// S-0748, ADR-0304: a permission's calls name a to condition with the state
// it names, as the role file writes it, by the old name or the new.
func TestADefinedRolesPermissionNamesTheStateItsToConditionNames(t *testing.T) {
	p := harbour(t)
	dir := filepath.Join(p.Root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	file := strings.Replace(reviewerRole, "permissions: []", "permissions:\n  - setting: roles.reviewer.permissions.publish\n    tools: [{call: item_move, when: {to: backlog}}]\n    commands: [{call: flai move, when: to-backlog}]", 1)
	if err := os.WriteFile(filepath.Join(dir, "reviewer.md"), []byte(file), 0o644); err != nil {
		t.Fatal(err)
	}
	got := DefinedRoleSettings(p.Root, manifest.Manifest{})
	if len(got) != 1 || len(got[0].Settings) == 0 {
		t.Fatalf("want the reviewer and its settings, got %+v", got)
	}
	s := got[0].Settings[len(got[0].Settings)-1]
	if want := []string{"item_move when {to: backlog}", "flai move when {to: backlog}"}; s.Key != "roles.reviewer.permissions.publish" || !slices.Equal(s.Allows, want) {
		t.Errorf("the permission: %+v, want it to allow %v", s, want)
	}
}
