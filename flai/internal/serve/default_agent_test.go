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
// story. A story that names a harness keeps its own agent, and with no
// default and no command the start is refused as before.
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
	t.Run("no default and no command is still refused", func(t *testing.T) {
		lab := lab(t)
		id := lab.ready("Nothing")
		_, err := Start(ctx, lab.o, entry(lab), id)
		var no *Refused
		if !errors.As(err, &no) || !strings.Contains(no.Why, id+" names no harness, and no command is set on the host") {
			t.Errorf("start %s: %v, want refused for naming no harness", id, err)
		}
		if r := lab.run(id); r != nil {
			t.Errorf("a refused start left a run: %+v", r)
		}
	})
}
