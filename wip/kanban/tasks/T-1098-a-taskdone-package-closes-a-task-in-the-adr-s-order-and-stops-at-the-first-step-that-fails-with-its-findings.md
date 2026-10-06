---
id: T-1098
type: task
nature: improvement
title: A taskdone package closes a task in the ADR's order and stops at the first step that fails, with its findings
status: backlog
parent: S-0269
owner: alex
created: 2026-10-06T22:53:09Z
updated: 2026-10-06T22:53:36Z
transitions: []
stream: S-0269
tags: [cli, git]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/move.go, flai/cmd/touches.go, flai/cmd/check.go]
after: [T-1080, T-1085, T-1092]
---
# T-1098 A taskdone package closes a task in the ADR's order and stops at the first step that fails, with its findings

## Work

Compose the steps once, in a new package `flai/internal/taskdone`, so that the CLI, MCP, and the host channel share one implementation.

- `taskdone.go` has `Run`. It takes the repository, the runner, the task ID, the message, the log entry, the agent and session, and the clock, and runs the ADR's steps in its order (T-1080):
  1. commit in the story's worktree, where nothing to commit is not a failure;
  2. the storygit sync (T-1085);
  3. `workitem.TransitionAll` to done;
  4. `repo.LogStream`;
  5. widen the touches of the task and its story with the paths its commit changed that they do not cover, through `itemedit`'s claim, as `flai touches` records it;
  6. `check.Run`, scoped to the story;
  7. the inbox (T-1092).
- It answers one result struct with the fields the ADR names, JSON-tagged, and the step it stopped at. A failing step ends the run with that step's findings, and the steps after it are left undone.
- Where `flai move`, `flai touches`, or `flai check` keeps a rule or a scoping in its RunE that the task close needs, move that rule into a function both call. Leave each command's behaviour as it is. That is why `flai/cmd/move.go`, `touches.go`, and `check.go` are touched; leave untouched those it does not need.
- `taskdone_test.go` uses the `gittest` fixtures and covers three cases:
  - the happy path: commit, sync, move, log, touches widened, check clean, and the inbox answered;
  - a refused sync, a rebase left unfinished or uncommitted changes the commit cannot take, that stops before the move;
  - a failed check that answers its findings after the move and the log.

This task waits for T-1080, for the contract, and for T-1085 and T-1092, for the sync and the inbox it calls. It is layer 2.

## Done when

- [ ] `taskdone.Run` performs the seven steps in order and stops at the first failure, with its findings.
- [ ] Tests cover the happy path, a refused sync, and a failed check.
- [ ] `flai move`, `flai touches`, and `flai check` behave as before, and `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner. This task meets criterion 4, the tests, at the package level; T-1102 repeats the three cases through the command.
