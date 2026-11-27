---
id: T-3614
type: task
nature: remediation
title: flai serve stops a running rebase and resolve job whose story left review and ends it stopped, outcome left review
status: backlog
parent: S-0822
owner: alex
created: 2026-11-27T23:56:21Z
updated: 2026-11-27T23:56:21Z
transitions: []
stream: S-0822
tags: [flai, serve]
touches: [flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go]
after: [T-3520, T-3613]
---
# T-3614 flai serve stops a running rebase and resolve job whose story left review and ends it stopped, outcome left review

## Work

Build the stop T-3613's ADR decides.

In `flai/internal/serve/verifyjob.go`, the verify pass keeps a cancel for each running rebase and resolve job, derived from the job's context in `runJob`. At each `look`, it reads the state of each story with a running resolve job. When the story is no longer in review, it cancels that job's context. That stops the process group, as `runCloseOut` does when its context ends. `runJob` tells this stop apart from a stop with flai serve, which leaves the job for the next serve to end `interrupted`, and from the time limit. For a left-review stop it writes the job's end and journals it. It announces a `verified` event with verdict `stopped`.

In `flai/internal/serve/resolvejob.go`, `settle` gives a job stopped this way verdict `stopped`, `stopped_at` the step it was at (`rebase` or `close-out`), outcome `left review` (T-3520's `closeout.OutcomeLeft`), and the reason T-3520 writes, which names the state, who moved the story, and when. It sends nothing back. `rebase` aborts a rebase the job started when the stop comes during the sync, so the branch is as it was. A rebase it found waiting is left as it was.

In `flai/internal/serve/resolvejob_test.go`, add a test with a close-out that blocks until its context ends. The test moves the story to in-progress while the close-out runs, then reads the job record: verdict `stopped`, `stopped_at` `close-out`, outcome `left review`, and the reason. It also checks that the story stayed in progress and no send-back thread was opened. Add a case that stops during the rebase and finds no rebase left waiting. In `flai/internal/serve/verifyjob_test.go`, check that a close-out job whose story leaves review is not stopped, and that a stop with flai serve still leaves a resolve job to be ended `interrupted`.

It waits for T-3520, which adds the `left review` outcome and its reason and changes the same `resolvejob.go` and its tests. It also waits for T-3613, whose ADR fixes the record's fields.

## Done when

- A rebase and resolve job whose story leaves review while it runs is stopped within a look of the move. It ends `stopped` at the step it was at, outcome `left review`, with nothing sent back.
- A stop during the rebase leaves no rebase the job started waiting in the worktree.
- A close-out job, a stop with flai serve, and the time limit behave as before.
- The new tests pass with `flai test` on the changed paths.

## Notes
