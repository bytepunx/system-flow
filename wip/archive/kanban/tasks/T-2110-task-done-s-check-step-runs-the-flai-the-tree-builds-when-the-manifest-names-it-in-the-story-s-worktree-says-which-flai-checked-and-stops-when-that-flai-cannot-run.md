---
id: T-2110
type: task
nature: improvement
title: task done's check step runs the flai the tree builds when the manifest names it, in the story's worktree, says which flai checked, and stops when that flai cannot run
status: done
parent: S-0495
owner: alex
created: 2026-10-25T20:33:10Z
updated: 2026-10-25T20:55:55Z
transitions:
  - to: ready
    at: 2026-10-25T20:36:40Z
    by: agent-S-0495
  - to: in-progress
    at: 2026-10-25T20:45:08Z
    by: agent-S-0495
  - to: done
    at: 2026-10-25T20:55:55Z
    by: agent-S-0495
stream: S-0495
tags: [cli, mcp]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/buildinfo/task-done-rules.txt]
after: [T-2109]
usage:
  source: log
  seconds: 647
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 39
      output: 14329
      cache_read: 2090856
      cache_write: 112070
      cost: 1.4079
---
# T-2110 task done's check step runs the flai the tree builds when the manifest names it, in the story's worktree, says which flai checked, and stops when that flai cannot run

## Work

In `flai/internal/taskdone/taskdone.go`, `run.check` calls `check.RunScoped(r.wtRepo, …)` in-process, so the binary running the close checks the worktree by its own rules: the cause of I-0194. When the worktree's manifest sets the key T-2109 adds, the step instead runs `<worktree>/<command> check --strict --story <S-nnnn> --json` in the worktree, through `r.o.Runner` so that tests can fake it, with `FLAI_AGENT` and `FLAI_SESSION` as the run's, decodes the JSON into the `Check` the result carries (`Passed`, `Findings`, `Notes`, as `newCheck` builds them, `MaxNotes` kept), and records on the result which flai checked, `checked_by`, naming the command and the version its `version --json` prints, as `closed_by` names the installed flai (ADR-0160). Without the key the step is unchanged.

A command that cannot run, exits without a result, or prints output that is not a result stops the step as a failed check does, with the command's output in the error, so a branch whose flai does not build closes no task; the ADR of T-2108 settles this rule and whether `flai/internal/buildinfo/task-done-rules.txt` is raised by one, since ADR-0180 says a change to what task done checks raises it; raise it when the ADR says so, in this task.

The MCP `task_done` and the host channel's `task.done` go through `taskdone.Run`, so they need no change; `TestTaskDoneAnswersAFailedCheck` in `flai/internal/mcpserver/task_test.go` must still pass.

Reproduce I-0194 in `taskdone_test.go`, with a fake tree flai as `newFakeFlai` in `flai/internal/mcpserver/installed_test.go` fakes the installed one: a worktree carrying a finding the in-process check reports, such as an examples file with a key `components.examples` does not know, and a manifest naming a fake that prints a passing result: the close passes and `checked_by` names the fake. A fake that prints findings stops the step at the check, after the commit, sync, and move, as `TestRunStopsAtAFailedCheckAfterTheMoveAndTheLog` shows for the in-process check. A fake that exits without a result stops the step with its output.

Waits for T-2109, which adds the key it reads, and through it for T-2108. The docs task and the issue close wait for this one.

## Done when

`go test ./internal/taskdone/... ./internal/mcpserver/...` passes with the three new tests; `flai test` on the changed paths is clean; and, run from the story's worktree, `scripts/flai.sh task done T-<this task>` closes this task with `checked_by` naming `scripts/flai.sh`, since the worktree's manifest names it and the tree's flai carries this change.

## Notes

The MCP server's own `task_done` runs the installed flai, which gains this change only when the release that carries it is installed; until then the tree's `scripts/flai.sh task done` is where the change can be seen.
