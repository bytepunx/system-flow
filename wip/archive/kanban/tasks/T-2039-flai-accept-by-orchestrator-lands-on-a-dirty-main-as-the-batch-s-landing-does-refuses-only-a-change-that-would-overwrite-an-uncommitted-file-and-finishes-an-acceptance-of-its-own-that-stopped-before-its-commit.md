---
id: T-2039
type: task
nature: improvement
title: flai accept --by orchestrator lands on a dirty main as the batch's landing does, refuses only a change that would overwrite an uncommitted file, and finishes an acceptance of its own that stopped before its commit
status: done
parent: S-0467
owner: alex
created: 2026-10-22T20:59:30Z
updated: 2026-10-23T05:51:47Z
transitions:
  - to: ready
    at: 2026-10-23T04:26:16Z
    by: agent-S-0467
  - to: in-progress
    at: 2026-10-23T04:26:17Z
    by: agent-S-0467
  - to: done
    at: 2026-10-23T04:50:47Z
    by: agent-S-0467
stream: S-0467
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/internal/preview/accept.go, flai/cmd/accept_orchestrator_dirty_test.go]
after: [T-2037]
usage:
  source: log
  seconds: 1470
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 72
      output: 28037
      cache_read: 4968032
      cache_write: 244322
      cost: 2.9037
---
# T-2039 flai accept --by orchestrator lands on a dirty main as the batch's landing does, refuses only a change that would overwrite an uncommitted file, and finishes an acceptance of its own that stopped before its commit

## Work

Waits for the ADR task, which settles the decision this builds.

Today `flai accept` refuses the real run when `storygit.DirtyOutsideWip` lists anything, unless `--yes` includes the files in the commit or `--batch` leaves them (`flai/cmd/accept.go`, the `res.Uncommitted` check; `leftUncommitted` and `addArgs` in `flai/cmd/accept_batch_dirty.go`). The dry run only prints the files as a note, so the orchestrator reads "no blocker" and then fails on the real run (TH-0394). Change it so that an acceptance by the orchestrator (`opts.Orchestrator`, from `--by orchestrator`) takes the batch's path:

- `leftUncommitted` leaves the main checkout's uncommitted files outside `wip/` as they are for an orchestrator's acceptance as it does under `--batch`: out of `fixUncommitted` and `refuseUnfixedUncommitted`, and out of `addArgs`, so the acceptance commit holds only the acceptance.
- In `preview.AcceptWith`, when the acceptance is the orchestrator's and the story's change would overwrite one of those files (`storygit.Overwritten`, as `mergequeue.Land` uses it), add a blocker naming the files and the command that frees them, `git stash push -- <files>`, or a commit of them. The dry run reports it under the blockers, and the real run refuses on it before anything is merged. Files the change leaves alone are no blocker and no refusal.
- Narrow the rule at the top of `acceptItem` that refuses the orchestrator a resumed acceptance: keep it for an acceptance someone else made, and let the orchestrator finish one whose done transition is `by: orchestrator` (`res.Resumed`, with or without `CommitOnly`), through `finishCommit` as the operator's `flai accept` does today. The refusal's text names who made the acceptance.

Tests: in `flai/cmd/accept_orchestrator_test.go` or `accept_batch_dirty_test.go`, an orchestrator's acceptance with an unrelated uncommitted file outside `wip/` lands and leaves the file uncommitted and out of the commit; one whose change would overwrite the file is refused, and its dry run names the file and the command under the blockers; the orchestrator finishes an acceptance of its own that stopped before its commit, and is still refused one the operator made (the case `accept.go` had no test for). Keep `TestAnAcceptanceByHandStillRefusesTheUncommittedFiles` true: the operator's acceptance by hand keeps today's refusal and `--yes`.

## Done when

- An orchestrator's `flai accept` without the queue lands while the main checkout has uncommitted files outside `wip/` that its change leaves alone, leaves them out of the commit, and refuses one its change would overwrite, as a blocker the dry run lists too.
- The orchestrator finishes a resumed acceptance it made itself, and is refused one someone else made, with a test for each.
- `flai test` on the changed paths passes, and the operator's `flai accept` behaves as before.

## Notes

Layer 2, beside the own-thread task; the rebase task waits for this one because both change `flai/cmd/accept.go`. The batch's landing is the model: `flai/cmd/accept_batch_dirty.go`, `flai/internal/mergequeue/land.go` (`overwrittenBy`, `holdBack`), and ADR-0178.
