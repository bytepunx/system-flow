package mcpserver

import (
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// I-0234: the server takes the session from Options, not from the process,
// so a test that gives it a story, a role, reasons, a session, or who
// started it runs in parallel with the package's other tests, the
// environment holding none of them.
func TestTheServerTakesItsSessionFromOptions(t *testing.T) {
	t.Parallel()
	t.Run("a story's agent", func(t *testing.T) {
		t.Parallel()
		d := newDeadlines()
		// FLAI_STORY untrimmed, as flai serve sets it
		f := setupWith(t, func(o *Options) { o.After = d.after; o.Story = storysAgentStory })
		f.ask(t, f.story.ID)
		if out := f.endsAtOnce(t, d); out["why"] == "" {
			t.Errorf("the story's agent with only its question pending is told why it ends: %v", out)
		}
	})
	t.Run("the orchestrator for reasons", func(t *testing.T) {
		t.Parallel()
		d := newDeadlines()
		f := setupWith(t, func(o *Options) { o.After = d.after; o.Role = guard.RoleOrchestrate; o.Reasons = "ready-empty" })
		if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
			t.Fatal(failed)
		}
		if out := f.endsAtOnce(t, d); out["why"] != "reasons done" {
			t.Errorf("the orchestrator whose reasons are done is told so: %v", out["why"])
		}
	})
	t.Run("the fields", func(t *testing.T) {
		t.Parallel()
		f := setup(t)
		s := newServer(Options{Repo: f.repo, Story: " S-0007\n", Role: "plan", Session: "s1", Reasons: " \n", StartedBy: guard.StartedByServe}, f.repo)
		if s.story != "S-0007" || s.recordingStory("") != "S-0007" || s.role != "plan" || s.session != "s1" || s.reasons || !s.served {
			t.Errorf("the session from Options: story %q role %q session %q reasons %v served %v", s.story, s.role, s.session, s.reasons, s.served)
		}
		if s := newServer(Options{Repo: f.repo, StartedBy: guard.StartedByServe}, f.repo); !s.served || s.story != "" || s.role != "" {
			t.Errorf("flai serve's mark alone makes a served session: %+v", s)
		}
		if s := newServer(Options{Repo: f.repo, Session: "s2"}, f.repo); s.served || s.reasons {
			t.Errorf("a session flai serve did not start: served %v reasons %v", s.served, s.reasons)
		}
	})
}

// SessionFromEnv reads each variable flai serve sets in an agent's session,
// FLAI_STORY trimmed, as flai mcp passes them to the server. It sets the
// process environment, so it does not run in parallel.
func TestSessionFromEnvReadsEachVariable(t *testing.T) {
	t.Setenv("FLAI_STORY", " S-0003\n")
	t.Setenv("FLAI_ROLE", guard.RoleOrchestrate)
	t.Setenv("FLAI_SESSION", "20261110T085006")
	t.Setenv(harness.ReasonsEnv, "ready-empty")
	t.Setenv(guard.StartedByEnv, guard.StartedByServe)
	want := SessionEnv{Story: "S-0003", Role: guard.RoleOrchestrate, Session: "20261110T085006", Reasons: "ready-empty", StartedBy: guard.StartedByServe}
	if got := SessionFromEnv(); got != want {
		t.Errorf("SessionFromEnv is %+v, want %+v", got, want)
	}
	opt := Options{Agent: "claude"}.WithSession(want)
	if opt.Story != want.Story || opt.Role != want.Role || opt.Session != want.Session || opt.Reasons != want.Reasons || opt.StartedBy != want.StartedBy || opt.Agent != "claude" {
		t.Errorf("WithSession sets the session and keeps the rest: %+v", opt)
	}
}
