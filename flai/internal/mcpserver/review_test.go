package mcpserver

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// statusRunner answers git status in a worktree with what it holds.
type statusRunner struct {
	execx.System
	status string
}

func (r *statusRunner) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && len(args) > 0 && args[0] == "status" {
		return r.status, nil
	}
	return "", nil
}

// S-0140: item_move refuses a story whose worktree holds uncommitted work,
// naming it, and moves it once the worktree is clean. The server hands its
// runner to the repo, so the rule asks git as flai move does.
func TestItemMoveRefusesReviewWithUncommittedWork(t *testing.T) {
	t.Parallel()
	f := setup(t)
	git := &statusRunner{status: "?? docs/left.md\n M flai/cmd/move.go"}
	New(Options{Repo: f.repo, Runner: git})
	if f.repo.Git != execx.Runner(git) {
		t.Fatal("the server must give the repo its runner")
	}
	if err := os.MkdirAll(f.repo.WorktreePath(f.story.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{"ready", "in-progress", "done"} {
		if _, failed := f.call(t, "item_move", map[string]any{"id": f.task.ID, "to": to}); failed != "" {
			t.Fatal(failed)
		}
	}
	_, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "review"})
	if !strings.Contains(failed, "has uncommitted changes (docs/left.md, flai/cmd/move.go)") || !strings.Contains(failed, "commit them on story/"+f.story.ID) {
		t.Fatalf("review with uncommitted work: %q", failed)
	}
	if it, _ := f.repo.Get(f.story.ID); it.Status != workitem.InProgress {
		t.Errorf("a refused move changes nothing: %s", it.Status)
	}
	git.status = ""
	if out, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "review"}); failed != "" || out["status"] != "review" {
		t.Errorf("a clean worktree moves: %v %s", out, failed)
	}
}

// The commits of moveKlusterGit's history, oldest first: the first two
// change code, the third only the docs folder.
var (
	moveKlusterC1 = strings.Repeat("a", 40)
	moveKlusterC2 = strings.Repeat("b", 40)
	moveKlusterC3 = strings.Repeat("c", 40)
)

// moveKlusterGit is an execx.Runner that answers git from a linear history
// whose branch story/<story> is at head: rev-parse resolves the branch,
// merge-base --is-ancestor walks the history, diff --name-only lists what
// each commit changed, and anything else, git status among it, answers
// nothing.
type moveKlusterGit struct {
	execx.System
	story, head string
}

var (
	moveKlusterParent = map[string]string{moveKlusterC2: moveKlusterC1, moveKlusterC3: moveKlusterC2}
	moveKlusterFiles  = map[string][]string{moveKlusterC2: {"flai/cmd/move.go"}, moveKlusterC3: {"docs/users/move.md"}}
)

func (g moveKlusterGit) Run(_, name string, args ...string) (string, error) {
	if name != "git" || len(args) == 0 {
		return "", nil
	}
	switch args[0] {
	case "rev-parse":
		if args[len(args)-1] == "refs/heads/story/"+g.story+"^{commit}" {
			return g.head + "\n", nil
		}
		return "", errors.New("exit status 1")
	case "merge-base":
		for c := args[len(args)-1]; c != ""; c = moveKlusterParent[c] {
			if c == args[len(args)-2] {
				return "", nil
			}
		}
		return "", errors.New("exit status 1")
	case "diff":
		from, to := args[len(args)-2], args[len(args)-1]
		var out []string
		for c := to; c != "" && c != from; c = moveKlusterParent[c] {
			out = append(out, moveKlusterFiles[c]...)
		}
		return strings.Join(out, "\x00"), nil
	}
	return "", nil
}

// moveKlusterRecord writes r as story's namespace's record under root, or,
// with r's Commit empty, writes none.
func moveKlusterRecord(t *testing.T, root, story string, r kluster.Record) {
	t.Helper()
	if r.Commit == "" {
		return
	}
	ns, err := kluster.StoryNamespace(story)
	if err != nil {
		t.Fatal(err)
	}
	r.Story, r.Namespace = story, ns
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	path := kluster.RecordPath(root, ns)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// moveKlusterActions is a HostAction with the kluster and checks host
// actions on while on is, and every other off.
func moveKlusterActions(on bool) HostAction {
	return func(action, _ string) (bool, error) {
		return on && (action == hostapi.ActionKluster || action == hostapi.ActionChecks), nil
	}
}

// S-0662 criteria 1 and 2, ADR-0292 §6: item_move to review reads the
// story's namespace's record against its branch's head, in the words flai
// move says: a warning naming the commands that deploy and check when it
// was not deployed and checked there, what stands added to the warning of
// unticked criteria, nothing while the kluster host action is off or it
// passed at the head, and never a refusal.
func TestItemMoveToReviewSaysWhatTheKlusterNamespaceShowsAtItsHead(t *testing.T) {
	t.Parallel()
	c1, c2 := moveKlusterC1[:12], moveKlusterC2[:12]
	passed := kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC2, ChecksOutcome: kluster.OutcomePassed}
	for _, c := range []struct {
		name   string
		on     bool
		record kluster.Record
		head   string
		says   string // what the namespace's sentence says, empty for none
		warns  bool   // a warning of its own says it
	}{
		{"the action off", false, kluster.Record{Commit: moveKlusterC1, Outcome: kluster.OutcomeFailed}, moveKlusterC2, "", false},
		{"deployed and checked at the head", true, passed, moveKlusterC2, "'s namespace %s is deployed and checked at the head " + c2, false},
		{"deployed and checked before a docs commit", true, passed, moveKlusterC3, "'s namespace %s is deployed and checked at the head " + moveKlusterC3[:12] + " (at " + c2, false},
		{"not deployed", true, kluster.Record{}, moveKlusterC2, "'s namespace %s is not deployed at the head " + c2 + ": no deploy of it is recorded", true},
		{"deployed at an older commit", true, kluster.Record{Commit: moveKlusterC1, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC1, ChecksOutcome: kluster.OutcomePassed}, moveKlusterC2,
			"'s namespace %s is not deployed at the head " + c2 + ": its last deploy, at " + c1 + ", passed, and 1 path outside design, docs, and wip changed", true},
		{"deploy failed", true, kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomeFailed, FailedStep: kluster.StepDeploy}, moveKlusterC2,
			"'s namespace %s is not deployed at the head " + c2 + ": its last deploy there failed at the deploy step", true},
		{"checks failed", true, kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC2, ChecksOutcome: kluster.OutcomeFailed}, moveKlusterC2,
			"'s namespace %s is deployed at the head " + c2 + ", and its checks run there ended failed", true},
	} {
		for _, unticked := range []bool{false, true} {
			name := c.name
			if unticked {
				name += " with a criterion unticked"
			}
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				git := &moveKlusterGit{head: c.head}
				f := setupWith(t, func(o *Options) { o.Runner, o.HostAction = git, moveKlusterActions(c.on) })
				id := f.story.ID
				git.story = id
				ns, _ := kluster.StoryNamespace(id)
				says := ""
				if c.says != "" {
					says = id + strings.Replace(c.says, "%s", ns, 1)
				}
				f.repo.Manifest.Kluster = &manifest.Kluster{Deploy: []string{"scripts/kluster-deploy.sh"}}
				moveKlusterRecord(t, projectRoot(f.repo), id, c.record)
				if unticked {
					data, err := os.ReadFile(f.story.Path)
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(f.story.Path, []byte(strings.Replace(string(data), "- [x] works\n", "- [x] works\n- [ ] more\n", 1)), 0o644); err != nil {
						t.Fatal(err)
					}
				}
				for _, to := range []string{"ready", "in-progress", "done"} {
					if _, failed := f.call(t, "item_move", map[string]any{"id": f.task.ID, "to": to}); failed != "" {
						t.Fatal(failed)
					}
				}
				out, failed := f.call(t, "item_move", map[string]any{"id": id, "to": "review"})
				if failed != "" || out["status"] != "review" {
					t.Fatalf("the move to review must go through: %v %s", out, failed)
				}
				var namespace, criteria []string
				warnings := moveKlusterStrings(out["warnings"])
				for _, w := range warnings {
					switch {
					case strings.HasPrefix(w, id+" goes to review with unticked acceptance criteria (1 of 2)"):
						criteria = append(criteria, w)
					case strings.Contains(w, "namespace"):
						namespace = append(namespace, w)
					}
				}
				switch {
				case !c.warns && len(namespace) > 0:
					t.Errorf("no warning of the namespace is wanted, got %q", namespace)
				case c.warns && (len(namespace) != 1 || !strings.Contains(namespace[0], says)):
					t.Errorf("one warning saying %q is wanted, got %q", says, namespace)
				case c.warns:
					for _, cmd := range []string{"flai kluster deploy " + id + " --job", "kluster_deploy", "flai checks run " + id + " --job", "kluster_checks", "The move is made all the same"} {
						if !strings.Contains(namespace[0], cmd) {
							t.Errorf("the warning must name %q: %s", cmd, namespace[0])
						}
					}
				}
				switch {
				case !unticked && len(criteria) > 0:
					t.Errorf("all ticked, no warning of criteria is wanted: %q", criteria)
				case unticked && len(criteria) != 1:
					t.Fatalf("one warning of unticked criteria is wanted: %q", warnings)
				case unticked && says == "" && !strings.HasSuffix(criteria[0], "why any other is left unticked"):
					t.Errorf("with the action off the warning of criteria must say nothing of the namespace: %s", criteria[0])
				case unticked && says != "" && !strings.Contains(criteria[0], "left unticked; "+says):
					t.Errorf("the warning of criteria must add %q: %s", says, criteria[0])
				}
			})
		}
	}
}

// ADR-0292 §6: host actions that cannot be read are a warning, not a
// refusal: item_move makes the move to review.
func TestItemMoveToReviewWarnsOfHostActionsItCannotRead(t *testing.T) {
	t.Parallel()
	git := &moveKlusterGit{head: moveKlusterC2}
	f := setupWith(t, func(o *Options) {
		o.Runner = git
		o.HostAction = func(string, string) (bool, error) { return false, errors.New("parse the configuration cfg.json: bad") }
	})
	git.story = f.story.ID
	f.repo.Manifest.Kluster = &manifest.Kluster{Deploy: []string{"scripts/kluster-deploy.sh"}}
	for _, to := range []string{"ready", "in-progress", "done"} {
		if _, failed := f.call(t, "item_move", map[string]any{"id": f.task.ID, "to": to}); failed != "" {
			t.Fatal(failed)
		}
	}
	out, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": "review"})
	want := f.story.ID + " went to review, but what its kluster namespace shows at its head was not read: cannot tell whether the kluster and checks host actions are on: parse the configuration cfg.json: bad"
	if failed != "" || out["status"] != "review" || !strings.Contains(strings.Join(moveKlusterStrings(out["warnings"]), "\n"), want) {
		t.Fatalf("the move must go through and warn %q: %v %s", want, out, failed)
	}
}

// moveKlusterStrings is v, a JSON array of strings, as strings.
func moveKlusterStrings(v any) []string {
	a, _ := v.([]any)
	out := make([]string, 0, len(a))
	for _, e := range a {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
