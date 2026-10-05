---
id: S-0268
type: story
nature: improvement
title: The story agent's prompt says to batch independent edits and commands in one turn, to take a new task through ready in one command, and to prove a test fails without the shared stash
status: ready
owner: alex
created: 2026-10-05T00:06:36Z
updated: 2026-10-05T00:34:46Z
transitions:
  - to: ready
    at: 2026-10-05T00:14:14Z
    by: alex
tags: [flai, template]
topics: [cli, conventions]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/git.md, template/root/design/conventions/git.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/system/flai-cli.md, docs/operators/index.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 25m
  delivery: 2026-10-05T01:01:00Z
  basis: "Its own forecast of 25m; 2nd in the pull order with an in-progress limit of 3, behind S-0249 and S-0266."
  by: flai
  at: 2026-10-05T00:34:46Z
---
# S-0268 The story agent's prompt says to batch independent edits and commands in one turn, to take a new task through ready in one command, and to prove a test fails without the shared stash

## Goal

In S-0248's run the main agent made 64 model calls, median 2.8 seconds, 3.5 minutes in all, for three edits and a test: 32 of its 48 API messages carried one tool call, three consecutive edits to one file went out as three turns, and a `flai move T-nnnn in-progress` straight from backlog was refused and retried through ready. To prove its new test fails without the fix, it pushed and dropped a `git stash` in the stash stack every worktree on the host shares, which another session's stash could have been caught by. The prompt flai serve builds for a story agent (`flai/internal/harness/harness.go`) and the conventions it points at (`design/conventions/delegation.md`, `git.md`, and the template's copies) say none of this. Each is a sentence; together they save a minute or two per story and remove one hazard.

## Acceptance criteria
- [ ] The story agent's prompt says to make independent edits and commands in one turn and to move a task it has just written to ready and in-progress in one command
- [ ] `design/conventions/git.md` and the template's copy say how to show a test fails without the change under test (check the files out from the main branch into a scratch copy, or build the old binary) and that `git stash` is not used in a worktree, since the stash stack is shared across the host's worktrees
- [ ] The prompt's tests in `flai/internal/harness` cover the new sentences, and the design (`design/system/flai-cli.md` or where the prompt is described) records them

## Tasks
- T-0853 The story agent's prompt says to make independent edits and commands in one turn and to move a new task to ready and in-progress in one command
- T-0854 git.md says how to show a test fails without the change and never to use git stash in a worktree, and delegation.md says to batch independent calls
- T-0855 The design and the operator docs record the prompt's new sentences and git.md's rule against git stash

## Notes

Found by the operator's review of the S-0248 agent log on 2026-10-04 (`~/.flai/serve/agents/sf-S-0248-20261004T231346Z.log`, 23:14:54Z and 23:15:45Z to 23:16:14Z).

### Planning

The story declared no touches, so `flai touches suggest` was started from the paths that the goal and the criteria name. Where each touch came from:

- **Design (the goal and the criteria name them):**
  - `flai/internal/harness/harness.go`
  - `design/conventions/git.md` and `template/root/design/conventions/git.md`
  - `design/conventions/delegation.md` and `template/root/design/conventions/delegation.md`
  - `design/system/flai-cli.md`, whose `flai serve agent` row describes the prompt
- **Co-change, from `flai touches suggest`:**
  - `flai/internal/harness/harness_test.go`, changed with the seeds in 12% of their commits. The criteria ask for the prompt's tests there.
  - `docs/operators/index.md` (33%). Its Harnesses bullet says what the prompt tells the agent.
  - The suggestion also listed `template/CHANGELOG.md` and `template/template.yaml`. They are left out because the release tooling writes them, as `git.md` says.
  - `docs/users/flai.md` (66%) is left out, because it does not describe the story agent's prompt.

**Tags and topics:** tags `[flai, template]`, because the story delivers to both components. Topics `[cli, conventions]`, which are the topics of `design/system/flai-cli.md` and `design/system/agent-context.md`, the documents that describe the prompt.

**Forecast:** 25m, adjusted from `flai forecast`'s 16m. That figure is the median 86 s per unit of size over 7 done medium improvement stories, at size 11 (3 criteria, 8 touches). It was raised by 9m, because the change to `harness.go` makes the close-out run the flai tier, about two minutes a run, and the verifier runs it twice.

**Delivery:** 01:47Z, moved from `flai forecast`'s 01:05Z. That figure plays out the pull order but not the hold. This story shares `harness.go`, `harness_test.go`, `delegation.md`, and `flai-cli.md` with S-0266, so it starts only when S-0266 is accepted, forecast at 01:22Z, and then takes its 25m.

**Cost of delay:** not set yet. The story has no inputs and no epic, so `flai cod` refuses. TH-0117 asks the operator for `time_lost_per_cycle`, recommending 1h (150 USD/week).

**Overlap:** S-0266's T-0851 also edits `delegation()` in `harness.go`. Whichever story is pulled second syncs over the first.
