---
id: T-1102
type: task
nature: improvement
title: flai task done T-nnnn -m closes a task from the story's worktree and answers the result as text and --json
status: done
parent: S-0269
owner: alex
created: 2026-10-06T22:53:20Z
updated: 2026-10-07T02:58:33Z
transitions:
  - to: ready
    at: 2026-10-07T02:14:02Z
    by: agent-S-0269
  - to: in-progress
    at: 2026-10-07T02:14:02Z
    by: agent-S-0269
  - to: done
    at: 2026-10-07T02:58:33Z
    by: agent-S-0269
stream: S-0269
tags: [cli]
touches: [flai/cmd/items.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go]
after: [T-1098]
usage:
  source: log
  seconds: 633
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 79
      output: 27248
      cache_read: 4192923
      cache_write: 141740
      cost: 2.1842
    - model: claude-sonnet-5-5
      input: 2
      output: 16
      cache_read: 5573
      cache_write: 6018
      cost: 0.0133
---
# T-1102 flai task done T-nnnn -m closes a task from the story's worktree and answers the result as text and --json

## Work

Add the command on the CLI.

- In `flai/cmd/task_done.go`, write `newTaskDoneCmd`: `flai task done <task> -m "<message>" [--log "<entry>"] [--json]`. It calls `taskdone.Run` (T-1098) with the agent and session from `FLAI_AGENT` and `FLAI_SESSION`.
- In `flai/cmd/items.go`, register it under the `task` command only. A story or an epic is refused, with a reason that names the command for each.
- Text output prints a line for each step done (the commit, the sync, the move, the log, the touches widened, and the check), and then the inbox. A stop prints the step and its findings.
- `--json` prints the result struct.
- The exit code is non-zero when a step stopped the run, with a distinct code for a refused sync and for a failed check. The host channel maps those codes.
- The help says what each step does, that the move keeps its rules, and that the sync still refuses while a rebase is unfinished.
- `flai/cmd/task_done_test.go` covers, through the command:
  - the happy path in text and `--json`;
  - a refused sync;
  - a failed check;
  - `-m` missing;
  - a story ID given in place of a task.

This task waits for T-1098, whose `Run` it calls. It is layer 3, beside the MCP tool, with which it shares no file.

## Done when

- [ ] `flai task done T-nnnn -m` closes a task in one call and answers as text and `--json`, stopping at the first failure with its findings.
- [ ] Tests cover the happy path, a refused sync, and a failed check through the command.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner.
