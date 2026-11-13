package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0552 criterion 2: while the merge queue batches and is not held, flai move
// <story> done from review enqueues the story as flai accept does, and lands
// nothing: the story stays in review and the batch under way is not dropped.
func TestMoveToDoneEnqueuesWhileTheQueueBatches(t *testing.T) {
	root, _ := orchestratedQueueProject(t, true)
	mergeQueueOn(t, root)
	queueUnderWay(t, root, false)
	mainBefore := gitIn(t, root, "rev-parse", "main")
	out, errOut, code := runIn(t, root, "move", "S-0001", "done", "--by", "alex")
	if want := "enqueued S-0001 in the merge queue at position 2 as alex, and it stays in review"; code != 0 || !strings.Contains(out, want) {
		t.Fatalf("flai move done must answer %q: %d %s\n%s", want, code, errOut, out)
	}
	assertNothingAccepted(t, root, mainBefore)
	st := readQueueState(t, root)
	if len(st.Enqueued) != 2 || st.Enqueued[1].Story != "S-0001" || st.Enqueued[1].By != "alex" {
		t.Errorf("enqueued %+v, want S-0001 after S-0005, by alex", st.Enqueued)
	}
	assertBatchUnderWay(t, st)
}

// I-0347: S-0678 was accepted and archived with a blocked interval open.
// flai unblock reaches the archived item and closes the interval at its
// close, as the item was not blocked after it.
func TestUnblockClosesAnArchivedItemsIntervalAtItsClose(t *testing.T) {
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	dir := filepath.Join(root, "wip", "archive", "kanban", "stories")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `---
id: S-0001
type: story
nature: improvement
title: Accepted while blocked
status: done
owner: alex
created: 2026-09-15T08:00:00Z
updated: 2026-09-15T10:15:07Z
transitions:
  - to: ready
    at: 2026-09-15T09:00:00Z
    by: alex
  - to: in-progress
    at: 2026-09-15T09:18:41Z
    by: agent-S-0001
  - to: review
    at: 2026-09-15T10:13:48Z
    by: agent-S-0001
  - to: done
    at: 2026-09-15T10:15:07Z
    by: orchestrator
blocked:
  - from: 2026-09-15T09:26:16Z
    reason: pause for urgent fix
---

# S-0001 Accepted while blocked

## Goal

g

## Acceptance criteria

- [x] c

## Notes
`
	path := filepath.Join(dir, "S-0001-accepted-while-blocked.md")
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "unblock", "S-0001")
	if want := "S-0001 unblocked at its close, 2026-09-15T10:15:07Z"; code != 0 || !strings.Contains(out, want) {
		t.Fatalf("flai unblock must answer %q: %d %s\n%s", want, code, errOut, out)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if !it.Archived || it.Path != path {
		t.Errorf("the item must stay in the archive at %s, got %s (archived %v)", path, it.Path, it.Archived)
	}
	if it.IsBlocked() || it.Blocked[0].Until != "2026-09-15T10:15:07Z" {
		t.Errorf("blocked = %+v, want the interval closed at the item's close", it.Blocked)
	}
	if _, errOut, code := runIn(t, root, "unblock", "S-0001"); code == 0 || !strings.Contains(errOut, "not blocked") {
		t.Errorf("a second unblock must say the item is not blocked: %d %s", code, errOut)
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
// whose story/S-0001 branch is at head: rev-parse resolves the branch,
// merge-base --is-ancestor walks the history, and diff --name-only lists
// what each commit changed; anything else answers nothing.
type moveKlusterGit struct{ head string }

var (
	moveKlusterParent = map[string]string{moveKlusterC2: moveKlusterC1, moveKlusterC3: moveKlusterC2}
	moveKlusterFiles  = map[string][]string{moveKlusterC2: {"cli/feature.go"}, moveKlusterC3: {"docs/users/feature.md"}}
)

func (g moveKlusterGit) Run(_, name string, args ...string) (string, error) {
	if name != "git" || len(args) == 0 {
		return "", nil
	}
	switch args[0] {
	case "rev-parse":
		if args[len(args)-1] == "refs/heads/story/S-0001^{commit}" {
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

func (g moveKlusterGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (moveKlusterGit) LookPath(name string) (string, error) { return name, nil }

// moveKlusterProject is a project whose manifest has a kluster: block that
// names a deploy, with S-0001 in progress and its one task done; with
// unticked, one of S-0001's two acceptance criteria is not ticked.
func moveKlusterProject(t *testing.T, unticked bool) string {
	t.Helper()
	root := tempProject(t)
	box := "x"
	if unticked {
		box = " "
	}
	writeFiles(t, root, map[string]string{
		"system-flow.yaml": "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nkluster:\n  config: kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n",
		"kluster.yaml":     "name: t-kluster\nprofile: authstar\nprovider: k3d\n",
		"wip/kanban/stories/S-0001-deployed.md": `---
id: S-0001
type: story
nature: feature
title: Deployed
status: in-progress
owner: alex
created: 2026-09-15T08:00:00Z
updated: 2026-09-15T09:18:41Z
transitions:
  - to: ready
    at: 2026-09-15T09:00:00Z
    by: alex
  - to: in-progress
    at: 2026-09-15T09:18:41Z
    by: agent-S-0001
---

# S-0001 Deployed

## Goal

g

## Acceptance criteria

- [x] one
- [` + box + `] two

## Notes
`,
		"wip/kanban/tasks/T-0001-piece.md": `---
id: T-0001
type: task
nature: feature
title: Piece
status: done
parent: S-0001
owner: alex
created: 2026-09-15T08:00:00Z
updated: 2026-09-15T10:00:00Z
transitions:
  - to: ready
    at: 2026-09-15T09:20:00Z
    by: agent-S-0001
  - to: in-progress
    at: 2026-09-15T09:20:00Z
    by: agent-S-0001
  - to: done
    at: 2026-09-15T10:00:00Z
    by: agent-S-0001
---

# T-0001 Piece
`,
	})
	return root
}

// moveKlusterCases are the states of S-0001's namespace that its move to
// review reads: the kluster host action on or off, the record, the branch's
// head, what the namespace's sentence says (empty for none), and whether a
// warning of its own says it with the commands that deploy and check.
func moveKlusterCases() []struct {
	name   string
	on     bool
	record kluster.Record
	head   string
	says   string
	warns  bool
} {
	c1, c2 := moveKlusterC1[:12], moveKlusterC2[:12]
	return []struct {
		name   string
		on     bool
		record kluster.Record
		head   string
		says   string
		warns  bool
	}{
		{"the action off", false, kluster.Record{Commit: moveKlusterC1, Outcome: kluster.OutcomeFailed}, moveKlusterC2, "", false},
		{"deployed and checked at the head", true, kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC2, ChecksOutcome: kluster.OutcomePassed}, moveKlusterC2,
			"S-0001's namespace s-0001 is deployed and checked at the head " + c2, false},
		{"deployed and checked before a docs commit", true, kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC2, ChecksOutcome: kluster.OutcomePassed}, moveKlusterC3,
			"S-0001's namespace s-0001 is deployed and checked at the head " + moveKlusterC3[:12] + " (at " + c2, false},
		{"not deployed", true, kluster.Record{}, moveKlusterC2,
			"S-0001's namespace s-0001 is not deployed at the head " + c2 + ": no deploy of it is recorded", true},
		{"deployed at an older commit", true, kluster.Record{Commit: moveKlusterC1, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC1, ChecksOutcome: kluster.OutcomePassed}, moveKlusterC2,
			"S-0001's namespace s-0001 is not deployed at the head " + c2 + ": its last deploy, at " + c1 + ", passed, and 1 path outside design, docs, and wip changed", true},
		{"deploy failed", true, kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomeFailed, FailedStep: kluster.StepDeploy}, moveKlusterC2,
			"S-0001's namespace s-0001 is not deployed at the head " + c2 + ": its last deploy there failed at the deploy step", true},
		{"checks failed", true, kluster.Record{Commit: moveKlusterC2, Outcome: kluster.OutcomePassed, ChecksCommit: moveKlusterC2, ChecksOutcome: kluster.OutcomeFailed}, moveKlusterC2,
			"S-0001's namespace s-0001 is deployed at the head " + c2 + ", and its checks run there ended failed", true},
	}
}

// moveKlusterCommands are what a warning of a namespace not deployed and
// checked at the head names.
var moveKlusterCommands = []string{"flai kluster deploy S-0001 --job", "kluster_deploy", "flai checks run S-0001 --job", "kluster_checks", "The move is made all the same"}

// S-0662 criteria 1 and 2, ADR-0292 §6: flai move to review reads S-0001's
// namespace's record against its branch's head. It warns, naming the
// commands that deploy and check, when the namespace was not deployed and
// checked there, adds what stands to the warning of unticked criteria, says
// nothing of it while the kluster host action is off or it passed at the
// head, and refuses nothing.
func TestMoveToReviewSaysWhatItsKlusterNamespaceShowsAtItsHead(t *testing.T) {
	t.Parallel()
	for _, c := range moveKlusterCases() {
		for _, unticked := range []bool{false, true} {
			name := c.name
			if unticked {
				name += " with a criterion unticked"
			}
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				root := moveKlusterProject(t, unticked)
				if c.on {
					klusterOn(t, root, true)
				}
				writeS0001Deploy(t, root, c.record)
				out, errOut, code := runWithApp(t, &app{cwd: root, runner: moveKlusterGit{head: c.head}}, "move", "S-0001", "review", "--json")
				var res struct {
					Status   string   `json:"status"`
					Warnings []string `json:"warnings"`
				}
				if code != 0 || json.Unmarshal([]byte(out), &res) != nil || res.Status != workitem.Review {
					t.Fatalf("the move to review must go through: %d %s\n%s", code, errOut, out)
				}
				var namespace, criteria []string
				for _, w := range res.Warnings {
					switch {
					case strings.HasPrefix(w, "S-0001 goes to review with unticked acceptance criteria (1 of 2)"):
						criteria = append(criteria, w)
					case strings.Contains(w, "namespace"):
						namespace = append(namespace, w)
					}
				}
				switch {
				case !c.warns && len(namespace) > 0:
					t.Errorf("no warning of the namespace is wanted, got %q", namespace)
				case c.warns && (len(namespace) != 1 || !strings.Contains(namespace[0], c.says)):
					t.Errorf("one warning saying %q is wanted, got %q", c.says, namespace)
				case c.warns:
					for _, cmd := range moveKlusterCommands {
						if !strings.Contains(namespace[0], cmd) {
							t.Errorf("the warning must name %q: %s", cmd, namespace[0])
						}
					}
				}
				switch {
				case !unticked && len(criteria) > 0:
					t.Errorf("all ticked, no warning of criteria is wanted: %q", criteria)
				case unticked && len(criteria) != 1:
					t.Fatalf("one warning of unticked criteria is wanted: %q", res.Warnings)
				case unticked && c.says == "" && !strings.HasSuffix(criteria[0], "why any other is left unticked"):
					t.Errorf("with the action off the warning of criteria must say nothing of the namespace: %s", criteria[0])
				case unticked && c.says != "" && !strings.Contains(criteria[0], "left unticked; "+c.says):
					t.Errorf("the warning of criteria must add %q: %s", c.says, criteria[0])
				}
			})
		}
	}
}

// ADR-0292 §6: a record that cannot be read is a warning, not a refusal:
// the move to review is made.
func TestMoveToReviewWarnsOfAKlusterRecordItCannotRead(t *testing.T) {
	t.Parallel()
	root := moveKlusterProject(t, false)
	klusterOn(t, root, true)
	writeFiles(t, root, map[string]string{".flai-cache/kluster/deployments/s-0001.json": "not json"})
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: moveKlusterGit{head: moveKlusterC2}}, "move", "S-0001", "review", "--json")
	if want := "S-0001 went to review, but what its kluster namespace shows at its head was not read"; code != 0 || !strings.Contains(out, want) || !strings.Contains(out, `"review"`) {
		t.Fatalf("the move must go through and warn %q: %d %s\n%s", want, code, errOut, out)
	}
}
