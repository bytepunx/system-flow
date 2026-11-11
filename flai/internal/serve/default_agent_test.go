//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// withDefaultAgent sets the project's default agent in system-flow.yaml,
// claude-code with opus, as the dashboard's settings or a hand edit do,
// after the lab's stories were made: they carry no agent of their own
// (S-0573).
func (lab *agentLab) withDefaultAgent() {
	lab.t.Helper()
	path := filepath.Join(lab.root, manifest.File)
	data, err := os.ReadFile(path)
	if err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, "agent:\n  harness: claude-code\n  model: opus\n"...), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// S-0573, ADR-0251: a story made before the project had a default agent
// carries none, and flai serve starts it with the project's default as
// system-flow.yaml holds it at the start, on the operator's word and on the
// launcher's own, with no command on the host and nothing written to the
// story. A story that names a harness keeps its own agent. S-0587,
// ADR-0261: with no default and no command, it starts on flai's built-in
// default, claude-code with no model; with a command set it starts on the
// command, as before; and a story whose harness is the command, with none
// set, is refused naming the page and the command that set it.
func TestAStoryWithNoAgentStartsWithTheProjectsDefault(t *testing.T) {
	ctx := context.Background()
	entry := func(lab *agentLab) Entry { return Entry{Key: "t", Name: "t", Root: lab.root} }
	// lab is a project whose host sets no command and runs the stub as
	// claude-code's program.
	lab := func(t *testing.T) *agentLab {
		lab := newAgentLab(t)
		lab.hold()
		lab.claudeStub()
		lab.cfg.Command = nil
		return lab
	}
	t.Run("started on the operator's word with the default", func(t *testing.T) {
		lab := lab(t)
		id := lab.ready("Seeded")
		lab.withDefaultAgent()
		run, err := Start(ctx, lab.o, entry(lab), id)
		if err != nil {
			t.Fatalf("start %s: %v, want it started with the project's default", id, err)
		}
		if run.Harness != harness.ClaudeCode || run.Model != "opus" || !run.live() {
			t.Errorf("run = %+v, want claude-code with opus, running", run)
		}
		if it, err := lab.repo.Get(id); err != nil || it.Agent != nil {
			t.Errorf("the story was written: %+v, %v", it.Agent, err)
		}
		lab.release(id)
	})
	t.Run("started by the launcher with the default", func(t *testing.T) {
		lab := lab(t)
		id := lab.ready("Seeded")
		lab.withDefaultAgent()
		lab.l.look(ctx, false)
		run := lab.run(id)
		if run == nil || run.Harness != harness.ClaudeCode || run.Model != "opus" {
			t.Fatalf("run = %+v, want claude-code with opus; log:\n%s", run, lab.logText())
		}
		if strings.Contains(lab.logText(), "names no harness") {
			t.Errorf("the launcher said it names no harness:\n%s", lab.logText())
		}
		lab.release(id)
	})
	t.Run("restarted in progress with the default", func(t *testing.T) {
		lab := lab(t)
		id := lab.ready("Begun")
		lab.move(id, workitem.InProgress)
		lab.withDefaultAgent()
		run, err := Restart(ctx, lab.o, entry(lab), id)
		if err != nil {
			t.Fatalf("restart %s: %v, want it started with the project's default", id, err)
		}
		if run.Harness != harness.ClaudeCode || run.Model != "opus" {
			t.Errorf("run = %+v, want claude-code with opus", run)
		}
		lab.release(id)
	})
	t.Run("a story that names a harness keeps its own model", func(t *testing.T) {
		lab := lab(t)
		id := lab.readyWith("Own", &manifest.Agent{Harness: harness.ClaudeCode, Model: "sonnet"})
		lab.withDefaultAgent()
		run, err := Start(ctx, lab.o, entry(lab), id)
		if err != nil {
			t.Fatal(err)
		}
		if run.Harness != harness.ClaudeCode || run.Model != "sonnet" {
			t.Errorf("run = %+v, want claude-code with its own sonnet", run)
		}
		lab.release(id)
	})
	t.Run("no default and no command starts on the built-in default", func(t *testing.T) {
		lab := lab(t)
		id := lab.ready("Nothing")
		run, err := Start(ctx, lab.o, entry(lab), id)
		if err != nil {
			t.Fatalf("start %s: %v, want it started on the built-in default", id, err)
		}
		if run.Harness != harness.ClaudeCode || run.Model != "" || !run.live() {
			t.Errorf("run = %+v, want claude-code with no model, running", run)
		}
		if it, err := lab.repo.Get(id); err != nil || it.Agent != nil {
			t.Errorf("the story was written: %+v, %v", it.Agent, err)
		}
		lab.release(id)
	})
	t.Run("the launcher starts it on the built-in default", func(t *testing.T) {
		lab := lab(t)
		id := lab.ready("Nothing")
		lab.l.look(ctx, false)
		run := lab.run(id)
		if run == nil || run.Harness != harness.ClaudeCode || run.Model != "" {
			t.Fatalf("run = %+v, want claude-code with no model; log:\n%s", run, lab.logText())
		}
		if st := lab.state(); strings.Contains(st.Waiting, "no harness") {
			t.Errorf("waiting = %q, want no story waiting for a harness", st.Waiting)
		}
		lab.release(id)
	})
	t.Run("a host command still comes before the built-in default", func(t *testing.T) {
		lab := newAgentLab(t)
		lab.hold()
		lab.claudeStub()
		id := lab.ready("Nothing")
		run, err := Start(ctx, lab.o, entry(lab), id)
		if err != nil {
			t.Fatal(err)
		}
		if run.Harness != harness.Command {
			t.Errorf("run = %+v, want the host's command", run)
		}
		lab.release(id)
	})
	t.Run("a harness that is the command, with none set, names the page and the command", func(t *testing.T) {
		lab := lab(t)
		id := lab.readyWith("Command", &manifest.Agent{Harness: harness.Command})
		_, err := Start(ctx, lab.o, entry(lab), id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, id+"'s harness is the host's command, and none is set for this project: set it on the dashboard's Harnesses page, /settings/harnesses, or with flai serve agent set -- <program> [args...]") {
			t.Errorf("start %s: %v, want refused naming /settings/harnesses and the command", id, err)
		}
		if r := lab.run(id); r != nil {
			t.Errorf("a refused start left a run: %+v", r)
		}
	})
	t.Run("a harness flai cannot start names the page and the command", func(t *testing.T) {
		lab := lab(t)
		id := lab.readyWith("Cursor", &manifest.Agent{Harness: "cursor"})
		_, err := Start(ctx, lab.o, entry(lab), id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, `flai cannot start harness "cursor"`) || !strings.Contains(no.Why, "/settings/agents, or with flai agent set --harness <name>") {
			t.Errorf("start %s: %v, want refused naming /settings/agents and the command", id, err)
		}
	})
}

// S-0587, ADR-0261: the dashboard's stand-in for a story this host has had no
// agent for names the harness it will start with: flai's built-in default,
// claude-code, when the story and the project name none and no command is
// set, and none, the host's command's, when one is.
func TestAStandInNamesTheBuiltinHarnessWhenNoCommandIsSet(t *testing.T) {
	lab := newAgentLab(t)
	id := lab.ready("Begun elsewhere")
	lab.move(id, workitem.InProgress)
	a := Activity(context.Background(), lab.root, lab.state(), false)[id]
	if a.Run == nil || a.Run.Harness != harness.ClaudeCode {
		t.Errorf("no command set: %+v, want a stand-in on claude-code", a.Run)
	}
	if a := Activity(context.Background(), lab.root, lab.state(), true)[id]; a.Run == nil || a.Run.Harness != "" {
		t.Errorf("a command set: %+v, want a stand-in naming no harness", a.Run)
	}
}
