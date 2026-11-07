package hostapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

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
// S-0586: the dashboard writes the kluster: block's commands as JSON
// argument lists, which flai manifest set reads as they are, and [] unsets
// one; anything else is refused naming the key.
func TestSettingsManifestWritesTheKlusterBlocksCommands(t *testing.T) {
	args, _, e := manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"kluster.deploy":["scripts/kluster-deploy.sh","--wait"],"kluster.build":[],"kluster.config":"kluster.yaml"}}`))
	if e != nil {
		t.Fatal(e)
	}
	want := []string{"manifest", "set", "--autocommit", "--trailer=" + Trailer, "--",
		"kluster.build=[]", "kluster.config=kluster.yaml", `kluster.deploy=["scripts/kluster-deploy.sh","--wait"]`}
	if !slices.Equal(args, want) {
		t.Errorf("args %q\nwant %q", args, want)
	}
	for _, bad := range []string{`"scripts/kluster-deploy.sh"`, `[1,2]`, `{"program":"x"}`, `true`} {
		_, _, e := manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"kluster.teardown":`+bad+`}}`))
		if e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, "kluster.teardown") || !strings.Contains(e.Message, "argument list") {
			t.Errorf("%s is not a command: %+v", bad, e)
		}
	}
}

// S-0595: the dashboard writes kluster.components as a JSON list of
// {name, build}, which flai manifest set reads as it is, and [] unsets it;
// anything else is refused naming the key (ADR-0263).
func TestSettingsManifestWritesTheKlusterBlocksComponents(t *testing.T) {
	args, _, e := manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"kluster.components":[{"name":"api","build":["scripts/build.sh","api"]}]}}`))
	if e != nil {
		t.Fatal(e)
	}
	want := []string{"manifest", "set", "--autocommit", "--trailer=" + Trailer, "--",
		`kluster.components=[{"name":"api","build":["scripts/build.sh","api"]}]`}
	if !slices.Equal(args, want) {
		t.Errorf("args %q\nwant %q", args, want)
	}
	args, _, e = manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"kluster.components":[]}}`))
	if e != nil || !slices.Contains(args, "kluster.components=[]") {
		t.Errorf("[] unsets it: %q %+v", args, e)
	}
	for _, bad := range []string{`"api"`, `[1]`, `[{"name":"api","image":"x"}]`, `{"name":"api"}`} {
		_, _, e := manifestEdit(channel.Project{}, json.RawMessage(`{"set":{"kluster.components":`+bad+`}}`))
		if e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, "kluster.components") || !strings.Contains(e.Message, "list of components") {
			t.Errorf("%s is not a list of components: %+v", bad, e)
		}
	}
}

// S-0586: enabling kluster from the Features page commits the block it
// writes with the dashboard's trailer; other actions commit nothing, and
// keep their command as it was.
func TestEnablingKlusterFromTheDashboardCarriesItsTrailer(t *testing.T) {
	p := channel.Project{Root: t.TempDir()}
	args, _, e := toggleAction(p, json.RawMessage(`{"action":"kluster","on":true}`))
	if e != nil || !slices.Equal(args, []string{"serve", "enable", "kluster", "--trailer=" + Trailer}) {
		t.Errorf("enable kluster: %q %+v", args, e)
	}
	args, _, e = toggleAction(p, json.RawMessage(`{"action":"kluster","on":false}`))
	if e != nil || !slices.Equal(args, []string{"serve", "disable", "kluster"}) {
		t.Errorf("disable kluster: %q %+v", args, e)
	}
	args, _, e = toggleAction(p, json.RawMessage(`{"action":"push","on":true}`))
	if e != nil || !slices.Equal(args, []string{"serve", "enable", "push"}) {
		t.Errorf("enable push: %q %+v", args, e)
	}
	if say := sayToggle(json.RawMessage(`{"action":"kluster","on":true}`)); say != "flai serve enable kluster" {
		t.Errorf("the journal names the command without its trailer: %q", say)
	}
}

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

// S-0587, ADR-0261: settings.agent and settings.harness write the project's
// own agent settings, running flai serve agent in its directory with no scope
// flag, and every project's, the host-wide fallback, only when the request
// says everywhere, with --all-projects.
func TestSettingsAgentAndHarnessWriteTheProjectsOwnUnlessEverywhere(t *testing.T) {
	p := withDocs(t)
	rec := &recorder{ran: Ran{Stdout: []byte(`{}`)}}
	m := writeMethods(rec.run, time.Now, Host{Enabled: func(a, _ string) bool { return a == ActionSettings }})
	for i, c := range []struct{ method, params, want string }{
		{"settings.agent", `{"name":"builder","command":["run-agent","{story}"]`, "serve agent set --name=builder --json -- run-agent {story}"},
		{"settings.agent", `{"name":"builder","everywhere":true`, "serve agent set --all-projects --name=builder --json"},
		{"settings.agent", `{"command":null`, "serve agent clear --json"},
		{"settings.agent", `{"command":null,"everywhere":true`, "serve agent clear --all-projects --json"},
		{"settings.harness", `{"name":"claude-code","program":"/opt/claude"`, "serve agent harness claude-code --program=/opt/claude --json"},
		{"settings.harness", `{"name":"claude-code","program":"/opt/claude","everywhere":true`, "serve agent harness claude-code --all-projects --program=/opt/claude --json"},
		{"settings.harness", `{"name":"claude-code","reset":true,"everywhere":false`, "serve agent harness claude-code --reset --json"},
	} {
		n := len(rec.runs)
		if _, e := m[c.method](context.Background(), p, json.RawMessage(fmt.Sprintf(`%s,"request_id":"req-%08d"}`, c.params, 701+i))); e != nil {
			t.Errorf("%s %s: %+v", c.method, c.params, e)
			continue
		}
		if len(rec.runs) != n+1 {
			t.Errorf("%s %s ran nothing", c.method, c.params)
			continue
		}
		if got := strings.Join(rec.runs[n].Args, " "); got != c.want {
			t.Errorf("%s %s ran %q, want %q", c.method, c.params, got, c.want)
		}
		if rec.runs[n].Dir != p.Root {
			t.Errorf("%s %s ran in %q, not the project's directory %q", c.method, c.params, rec.runs[n].Dir, p.Root)
		}
	}
	for i, c := range []struct{ method, params string }{
		{"settings.agent", `{"everywhere":true`},
		{"settings.harness", `{"name":"claude-code","everywhere":true`},
	} {
		n := len(rec.runs)
		if _, e := m[c.method](context.Background(), p, json.RawMessage(fmt.Sprintf(`%s,"request_id":"req-%08d"}`, c.params, 741+i))); e == nil || !strings.Contains(e.Message, "nothing to set") || len(rec.runs) != n {
			t.Errorf("%s %s: everywhere alone sets nothing, so it is refused and nothing runs: %+v", c.method, c.params, e)
		}
	}
}
