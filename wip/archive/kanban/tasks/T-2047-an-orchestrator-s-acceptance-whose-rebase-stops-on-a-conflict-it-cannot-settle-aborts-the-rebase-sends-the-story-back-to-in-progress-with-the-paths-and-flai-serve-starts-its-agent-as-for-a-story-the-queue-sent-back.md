---
id: T-2047
type: task
nature: improvement
title: An orchestrator's acceptance whose rebase stops on a conflict it cannot settle aborts the rebase, sends the story back to in-progress with the paths, and flai serve starts its agent as for a story the queue sent back
status: done
parent: S-0467
owner: alex
created: 2026-10-22T21:00:19Z
updated: 2026-10-23T05:51:47Z
transitions:
  - to: ready
    at: 2026-10-23T04:51:16Z
    by: agent-S-0467
  - to: in-progress
    at: 2026-10-23T04:51:17Z
    by: agent-S-0467
  - to: done
    at: 2026-10-23T05:20:25Z
    by: agent-S-0467
stream: S-0467
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2037, T-2039]
usage:
  source: log
  seconds: 1748
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 152
      output: 58764
      cache_read: 10412809
      cache_write: 512090
      cost: 6.0861
---
# T-2047 An orchestrator's acceptance whose rebase stops on a conflict it cannot settle aborts the rebase, sends the story back to in-progress with the paths, and flai serve starts its agent as for a story the queue sent back

## Work

Waits for the ADR task, which settles the decision, and for the dirty-main task, which changes `flai/cmd/accept.go` too.

Today `mergeStoryBranch` in `flai/cmd/branch.go` calls `storygit.Rebase`, which runs the merge steps on a stop (`ContinueOverResolved`: the issue files, the ADR index, `summary.md`, and a shared markdown file's sections, ADR-0175) and returns `StopConflicts` when a path is still unsettled. `acceptItem` logs the conflict with `logAcceptConflict` and returns the error, leaving the worktree mid-rebase, and `flai guard` refuses the orchestrator `git rebase --abort`, so the story waited for the operator (TH-0397). The merge queue's batch aborts the same stop itself (`rebase` in `flai/internal/mergequeue/batch.go`).

- Move the batch's abort into `flai/internal/storygit/sync.go` as one function the batch and `flai accept` share: it runs `git rebase --abort` in the worktree, checks the branch out at its head, and reports the conflict's paths; keep `batch.go` calling it so the batch's behaviour does not change.
- In `acceptItem`, for an acceptance by the orchestrator (`opts.Orchestrator`) whose `mergeStoryBranch` stops on conflicts the merge steps could not settle: log the conflict as today, abort with that function, move the story from review back to in-progress by `flai` with the reason naming the paths (`sync onto main stopped on conflicts in <paths>` or the like, through the repository's move so the transition is recorded), and open one thread on the story naming the paths and saying the agent runs `flai stream sync` and resolves them, then resumes the close-out and moves the story to review again. The command's output and `--json` say the story was sent back, with the paths. An acceptance by hand, or `--batch`, keeps today's behaviour.
- In `flai/internal/serve/agents.go`, start the story's agent for this send-back as `startQueueSentBack` does for the queue's: a story in progress, not blocked, whose last run ended before the send-back, with the thread `flai accept` opened still open, told of the thread and what stopped. Reuse `startSentBack`'s shape; the detection may read the story's last transition (`by: flai`, the reason) rather than a queue outcome.

Tests: in `flai/cmd/accept_conflict_test.go` an orchestrator's acceptance that stops on an unsettled conflict leaves the worktree clean at the branch head, moves the story to in-progress with the paths in the reason, opens the thread, and the operator's acceptance still refuses and leaves the stop for them as `TestAcceptLogsTheConflictItsRebaseStopsOn` asserts; in `flai/internal/storygit/sync_test.go` the shared abort restores the worktree; in `flai/internal/serve/agents_test.go` the agent is started for the send-back and not for a story blocked or whose run is live.

## Done when

- An orchestrator's `flai accept` whose rebase stops on a conflict the merge steps cannot settle leaves the worktree clean at the branch head, sends the story back to in-progress by `flai` with the paths, and opens one thread naming them.
- `flai serve` starts the story's agent for that send-back, told of the thread and the paths, with a test.
- The operator's acceptance by hand and the queue's batch behave as before, and `flai test` on the changed paths passes.

## Notes

Layer 3: after the dirty-main task on `flai/cmd/accept.go`. The queue's send-back (ADR-0175 § What the batch does with a stop it cannot settle, ADR-0176) and `startQueueSentBack` are the model. The orchestrator's prompt is changed in the prompt task to expect the send-back rather than a thread to the operator.
