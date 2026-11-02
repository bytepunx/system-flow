package hostapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// scheduledReviewer is reviewerRole on no item, started on its own schedule
// too, and declaring the publish permission, its own, and one of the
// orchestrator's, which is not its to name.
var scheduledReviewer = strings.NewReplacer(
	"  - asked\n", "  - asked\n  - {schedule: roles.reviewer.schedule}\n",
	"item: [story]", "item: none",
	"permissions: []", "permissions:\n  - setting: roles.reviewer.permissions.publish\n    tools: [release_publish]\n    commands: [{call: flai release, when: own-thread}]\n  - setting: orchestration.permissions.accept_reviews\n    tools: [verify]",
).Replace(reviewerRole)

// S-0420: a defined role's settings are its agent, the schedule its trigger
// names, and the permissions its file declares as its own, each with its
// value in the manifest and, for a permission, the calls its file lists
// under it; the built-in roles have none here, since their keys are their
// own blocks'.
func TestDefinedRoleSettingsListARolesAgentScheduleAndPermissions(t *testing.T) {
	p := harbour(t)
	dir := filepath.Join(p.Root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "reviewer.md"), []byte(scheduledReviewer), 0o644); err != nil {
		t.Fatal(err)
	}
	m := manifest.Manifest{Roles: map[string]manifest.RoleSettings{"reviewer": {Schedule: "daily", Permissions: manifest.Permissions{Publish: true}}}}
	got := DefinedRoleSettings(p.Root, m)
	if len(got) != 1 {
		t.Fatalf("want the reviewer alone, got %+v", got)
	}
	r := got[0]
	if r.Name != "reviewer" || r.Title != "Reviewer" || r.Action != "review" || r.Path != "design/workflows/roles/reviewer.md" || !strings.HasPrefix(r.Description, "Reviews a story") {
		t.Errorf("the role: %+v", r)
	}
	var keys []string
	for _, s := range r.Settings {
		keys = append(keys, s.Key)
	}
	if want := []string{"roles.reviewer.agent", "roles.reviewer.schedule", "roles.reviewer.permissions.publish"}; !slices.Equal(keys, want) {
		t.Fatalf("keys %v, want %v", keys, want)
	}
	if s := r.Settings[1]; s.Kind != manifest.KindCron || s.Value != "daily" || !s.Set {
		t.Errorf("the schedule: %+v", s)
	}
	pub := r.Settings[2]
	if pub.Kind != manifest.KindBoolean || pub.Value != true || !pub.Set || pub.Default != "false" || pub.Meaning == "" ||
		!slices.Equal(pub.Allows, []string{"release_publish", "flai release when own-thread"}) {
		t.Errorf("the permission: %+v", pub)
	}
	if r.Settings[0].Allows != nil {
		t.Errorf("the agent allows nothing: %+v", r.Settings[0])
	}
	out, err := json.Marshal(pub)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"key":"roles.reviewer.permissions.publish"`, `"kind":"boolean"`, `"value":true`, `"set":true`, `"allows":["release_publish","flai release when own-thread"]`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("JSON %s lacks %s", out, want)
		}
	}
	if got := DefinedRoleSettings(harbour(t).Root, manifest.Manifest{}); got == nil || len(got) != 0 {
		t.Errorf("with only the built-in roles: %#v, want none", got)
	}
}

// S-0420: settings.manifest writes a defined role's agent, schedule, and
// permissions with flai manifest set, as it writes the strategic agents'
// keys, and refuses a key flai manifest set does not write, naming it.
func TestSettingsManifestWritesARolesKeys(t *testing.T) {
	args, _, e := manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"roles.reviewer.agent":{"model":"claude-sonnet-5"},"roles.reviewer.schedule":"daily","roles.reviewer.permissions.publish":true},"unset":["roles.reviewer.permissions.answer_threads"]}`))
	if e != nil {
		t.Fatal(e)
	}
	want := []string{"manifest", "set", "--autocommit", "--trailer=" + Trailer, "--unset=roles.reviewer.permissions.answer_threads", "--",
		`roles.reviewer.agent={"model":"claude-sonnet-5"}`, "roles.reviewer.permissions.publish=true", "roles.reviewer.schedule=daily"}
	if !slices.Equal(args, want) {
		t.Errorf("args %q\nwant %q", args, want)
	}
	_, _, e = manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"roles.reviewer.permissions.fly":true}}`))
	if e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, "roles.reviewer.permissions.fly") {
		t.Errorf("an unknown permission: %+v", e)
	}
}
