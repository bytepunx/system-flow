package guard

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/ask"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0353: Decide holds every caller's rules on a Call that names no harness:
// its kind, its bare tool, whether a sub-agent makes it, and the role and
// story it carries, rather than the guard's, decide.
func TestDecideHoldsEachCallerOnANeutralCall(t *testing.T) {
	running := Guard{Commands: g.Commands, Running: func() ([]SubAgent, error) {
		return []SubAgent{{ID: "a1", Type: "explorer"}}, nil
	}}
	off := Guard{Commands: g.Commands, Permissions: except(manifest.PermitPromoteToReady)}
	sub := func(c Call) Call { c.Subagent, c.Agent = true, "verifier"; return c }
	for _, tc := range []struct {
		name  string
		guard Guard
		call  Call
		want  Outcome
		why   string // what the reason says, for a Deny
		needs string
	}{
		{"a sub-agent's write through flai", g, sub(Call{Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: "review"}}), Deny, "a sub-agent (verifier) cannot call item_move: it reads", ""},
		{"a sub-agent's read through flai", g, sub(Call{Kind: FlaiTool, Tool: "item_get", Args: Args{ID: "S-0001"}}), Allow, "", ""},
		{"a sub-agent's ADR without commit", g, sub(Call{Kind: FlaiTool, Tool: ADRNew}), Allow, "", ""},
		{"a sub-agent's ADR with commit", g, sub(Call{Kind: FlaiTool, Tool: ADRNew, Args: Args{Commit: true}}), Deny, "cannot call adr_new with commit true", ""},
		{"a sub-agent's commit", g, sub(Call{Kind: Shell, Command: "git add -A && git commit -m x"}), Deny, `a sub-agent (verifier) cannot run "git add -A"`, ""},
		{"a sub-agent's move through flai's command line", g, sub(Call{Kind: Shell, Command: "flai move S-0001 review"}), Deny, "flai commands that change work items", ""},
		{"a sub-agent's test run", g, sub(Call{Kind: Shell, Command: "scripts/flai.sh test flai/internal/guard"}), Allow, "", ""},
		{"a sub-agent's protected write", g, sub(Call{Kind: FileWrite, Name: "write_file", Path: "/w/.claude/settings.json"}), Deny, "a sub-agent (verifier) cannot use write_file on /w/.claude/settings.json", ""},
		{"a sub-agent's protected write with auto-approve on", Guard{Commands: g.Commands, AutoApprove: true}, sub(Call{Kind: FileWrite, Name: "write_file", Path: "/w/.claude/settings.json"}), Allow, "", ""},
		{"a sub-agent's write of code", g, sub(Call{Kind: FileWrite, Name: "write_file", Path: "/w/flai/x.go"}), Allow, "", ""},
		{"the story's agent's wait while a sub-agent runs", running, Call{Story: "S-0001", Kind: FlaiTool, Tool: "wait_for_events"}, Deny, "the story's agent cannot call wait_for_events while its sub-agents run (explorer a1) and no thread on S-0001", ""},
		{"a wait outside a story", running, Call{Kind: FlaiTool, Tool: "wait_for_events"}, Allow, "", ""},
		{"the planner's file write", g, Call{Role: RolePlan, Kind: FileWrite, Name: "write_file", Path: "flai/x.go"}, Deny, "the planner cannot use write_file: it plans through flai", ""},
		{"the planner's move to backlog", g, Call{Role: RolePlan, Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: Backlog}}, Allow, "", ""},
		{"the planner's story that is not a draft", g, Call{Role: RolePlan, Kind: FlaiTool, Tool: "item_new", Args: Args{Type: "story"}}, Deny, "the planner cannot create a story that is not a draft", ""},
		{"the orchestrator's promotion, its permission off", off, Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: Ready}}, Deny, "the orchestrator cannot move S-0001 to ready: it needs orchestration.permissions." + manifest.PermitPromoteToReady + ", which is off", "orchestration.permissions." + manifest.PermitPromoteToReady},
		{"the orchestrator's promotion, its permission on", Guard{Commands: g.Commands, Permissions: allOn}, Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: Ready}}, Allow, "", ""},
		{"the orchestrator's file write", off, Call{Role: RoleOrchestrate, Kind: FileWrite, Name: "write_file", Path: "x.md"}, Deny, "the orchestrator cannot use write_file: the orchestrator never does it", ""},
		{"the analyzer's story", g, Call{Role: RoleAnalyze, Kind: FlaiTool, Tool: "item_new", Args: Args{Type: "story"}}, Deny, "the analyzer cannot call item_new: the analyzer authors no stories", ""},
		{"the analyzer's commit", g, Call{Role: RoleAnalyze, Kind: Shell, Command: "git commit -m x"}, Deny, `the analyzer cannot run "git commit -m x": it runs only git's reads`, ""},
		{"the analyzer's issue", g, Call{Role: RoleAnalyze, Kind: FlaiTool, Tool: "issue_new"}, Allow, "", ""},
		{"a call of no kind the guard reads", g, sub(Call{Kind: Other, Name: "read_file"}), Allow, "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := Decide(tc.guard, tc.call)
			if v.Outcome != tc.want || !strings.Contains(v.Reason, tc.why) || v.Needs != tc.needs {
				t.Errorf("%+v: got %+v, want outcome %d with %q, needs %q", tc.call, v, tc.want, tc.why, tc.needs)
			}
			if v.Outcome == Allow && v.Reason != "" {
				t.Errorf("allowed with a reason: %q", v.Reason)
			}
		})
	}
}

// S-0353: a served agent writes only as itself whatever the guard's own role
// and story say: the Call's decide whose names it writes as (ADR-0143).
func TestDecideTakesTheRoleAndStoryFromTheCall(t *testing.T) {
	agent := Guard{Commands: g.Commands, Agent: "system-flow"}
	run := Call{Story: "S-0001", Kind: Shell, Command: "flai move S-0001 review --by agent-S-0001"}
	if v := Decide(agent, run); v.Outcome != Allow {
		t.Errorf("the story's agent as its run's name: %+v", v)
	}
	run.Command = "flai move S-0001 review --by alex"
	if v := Decide(agent, run); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the story's agent cannot run ") || !strings.Contains(v.Reason, "ADR-0143") {
		t.Errorf("the story's agent as another: %+v", v)
	}
	// the operator's own session, with neither role nor story, is not held
	run.Story = ""
	if v := Decide(agent, run); v.Outcome != Allow {
		t.Errorf("the operator's session held: %+v", v)
	}
}

// S-0354, ADR-0130: a story's own agent's write to a path flai protects in
// its worktree is put to the story's owner while auto-approve is off; a
// sub-agent's is refused as before (ADR-0102), and an unprotected path, a
// path in .git, and a protected path outside the worktree are not asked.
func TestDecideAsksForAStoryAgentsProtectedWriteInItsWorktree(t *testing.T) {
	wt := "/home/op/proj/.flai-cache/worktrees/S-0001"
	story := Guard{Commands: g.Commands, Served: true}
	write := func(path string) Call {
		return Call{Story: "S-0001", Kind: FileWrite, Name: "write_file", Path: path}
	}
	sub := write(wt + "/.claude/settings.json")
	sub.Subagent, sub.Agent = true, "general-purpose"
	operator := write(wt + "/.claude/settings.json")
	operator.Story = ""
	for _, tc := range []struct {
		name  string
		guard Guard
		call  Call
		want  Outcome
	}{
		{"the story's agent's protected write", story, write(wt + "/.claude/settings.json"), Ask},
		{"a protected file nested in the worktree", story, write(wt + "/template/root/.mcp.json"), Ask},
		{"auto-approve on", Guard{Commands: g.Commands, Served: true, AutoApprove: true}, write(wt + "/.claude/settings.json"), Allow},
		{"a sub-agent's protected write", story, sub, Deny},
		{"an unprotected path", story, write(wt + "/flai/cmd/guard.go"), Allow},
		{"a path in .git", story, write(wt + "/.git/config"), Allow},
		{"the main checkout's protected path", story, write("/home/op/proj/.claude/settings.json"), Allow},
		{"another story's worktree", story, write("/home/op/proj/.flai-cache/worktrees/S-0002/.claude/settings.json"), Allow},
		{"a path that leaves the worktree", story, write(wt + "/../S-0002/.claude/settings.json"), Allow},
		{"the operator's own session", story, operator, Allow},
		{"the planner's write", story, Call{Role: RolePlan, Story: "S-0001", Kind: FileWrite, Name: "write_file", Path: wt + "/.claude/settings.json"}, Deny},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := Decide(tc.guard, tc.call)
			if v.Outcome != tc.want {
				t.Fatalf("%+v: got %+v, want outcome %d", tc.call, v, tc.want)
			}
			if tc.want == Ask {
				rel := strings.TrimPrefix(tc.call.Path, wt+"/")
				for _, want := range []string{"write_file of " + rel + " waits on the story's owner", "a path flai protects is written only with the owner's approval on a thread on S-0001", "ADR-0086, ADR-0130"} {
					if !strings.Contains(v.Reason, want) {
						t.Errorf("reason %q lacks %q", v.Reason, want)
					}
				}
			}
			if tc.want == Deny && tc.call.Subagent && !strings.Contains(v.Reason, "a path flai protects") {
				t.Errorf("the sub-agent's refusal: %q", v.Reason)
			}
		})
	}
}

// holder is a reader for a harness that calls no permission tool of flai's:
// every input is a write of path, and it answers with exit 1 and the reason
// for anything but an allow.
type holder struct{ path string }

func (r holder) Read(_ []byte, g Guard) (Hook, error) {
	return Hook{Call: Call{Session: "s1", Role: g.Role, Story: g.Story, Kind: FileWrite, Name: "write_file", Path: r.path}}, nil
}

func (holder) Answer(v Verdict, stderr io.Writer) int {
	if v.Outcome == Allow {
		return 0
	}
	fmt.Fprintln(stderr, v.Reason)
	return 1
}

// holdFixture is a project with an in-progress story owned by alex, and an
// Asker on it that reads its threads every few milliseconds.
func holdFixture(t *testing.T) (*workitem.Repo, *workitem.Item, ask.Asker) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"design/system", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)
	story, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story", Owner: "alex", Now: at, Body: "## Goal\n\nWork.\n\n## Acceptance criteria\n\n- [ ] works\n"})
	if err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress} {
		if _, err := repo.Transition(story, to, "alex", "", at); err != nil {
			t.Fatal(err)
		}
	}
	return repo, story, ask.Asker{Repo: repo, Poll: 5 * time.Millisecond}
}

// S-0354, ADR-0130: for a reader that holds the call itself, Hold asks the
// story's owner on a thread within the harness's bound: the owner's allow
// lets the write through, and no answer within the bound refuses it naming
// the thread, which stays open.
func TestHoldAsksTheOwnerForAReaderThatHolds(t *testing.T) {
	repo, story, asker := holdFixture(t)
	gd := Guard{Commands: g.Commands, Story: story.ID, Served: true}
	r := holder{path: filepath.Join(repo.WorktreePath(story.ID), ".claude", "settings.json")}
	h, err := r.Read(nil, gd)
	if err != nil {
		t.Fatal(err)
	}

	var stderr bytes.Buffer
	if code := r.Answer(Hold(context.Background(), gd, asker, h.Call, 50*time.Millisecond), &stderr); code != 1 {
		t.Fatalf("unanswered: exit %d", code)
	}
	all, err := threads.List(repo)
	if err != nil || len(all) != 1 {
		t.Fatalf("threads %v, %v", all, err)
	}
	th := all[0]
	if !strings.Contains(stderr.String(), "No answer on "+th.ID) || !th.Open() {
		t.Errorf("unanswered: stderr %q, thread %s", stderr.String(), th.Status)
	}
	if first := th.Entries()[0]; first.Author != "agent-"+story.ID || !strings.Contains(first.Text, "write_file `.claude/settings.json`") || !strings.Contains(first.Text, "flai protects the path") {
		t.Errorf("the question: %+v", first)
	}

	if _, err := threads.Reply(repo, th.ID, "alex", "allow", time.Now()); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := r.Answer(Hold(context.Background(), gd, asker, h.Call, time.Second), &stderr); code != 0 || stderr.Len() != 0 {
		t.Errorf("the owner's allow: exit %d, stderr %q", code, stderr.String())
	}
	if back, err := threads.Get(repo, th.ID); err != nil || back.Open() {
		t.Errorf("the thread is settled: %v, %v", back, err)
	}

	// a call Decide does not ask about comes back as Decide gives it
	h.Call.Path = filepath.Join(repo.WorktreePath(story.ID), "flai", "x.go")
	if v := Hold(context.Background(), gd, asker, h.Call, time.Second); v.Outcome != Allow {
		t.Errorf("an unprotected write: %+v", v)
	}
	if all, _ := threads.List(repo); len(all) != 1 {
		t.Errorf("an unprotected write asked: %d threads", len(all))
	}
}
