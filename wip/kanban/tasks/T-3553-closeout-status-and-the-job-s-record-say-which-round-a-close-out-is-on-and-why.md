---
id: T-3553
type: task
nature: improvement
title: closeout.Status and the job's record say which round a close-out is on and why it started another, and resolve.status answers them
status: backlog
parent: S-0824
owner: alex
created: 2026-11-17T08:33:41Z
updated: 2026-11-17T08:33:41Z
transitions: []
stream: S-0824
tags: [flai, serve]
touches: [flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/hostapi/resolve_status.go, flai/internal/hostapi/resolve_status_test.go]
after: [T-3551]
---
# T-3553 closeout.Status and the job's record say which round a close-out is on and why it started another, and resolve.status answers them

## Work

Give the job's progress the round fields T-3551's ADR names. It waits for T-3551 for their names and the round line's words.

Today the round is only text in the log. `sync_round` in the close-out scripts prints `close-out: <why>; flai stream sync <story> (round N of M)`, and the last line ends `in round N of M`. `JobStatus`, `Progress`, and `Job` in `flai/internal/closeout/closeout.go` carry no round.

- `readLog` and `parseLine` read each round line into the round the job is on, the rounds it may run, and why the round started. Keep the why for each round, so that a reader sees why every round after the first ran.
- The ended job's status keeps the round it ended on.
- Fix `stepAfter`, which matches only `main moved during the run; flai stream sync`. A sync for `the branch is behind main` then reads as `commit`. Read every round line's sync as `StepSync`.
- `resolve.status` in `flai/internal/hostapi/resolve_status.go` answers the new fields, for a close-out job as for a resolve job.

Tests:

- in `flai/internal/closeout/closeout_test.go`, extend the `round 2 of 3` fixture and `TestRunningAndStatusReadAResolveJobsProgressFromItsLog`, so that a log of three rounds, each with its why, gives round 3 of 3 and both whys, and a behind-main sync reads as `StepSync`;
- in `flai/internal/hostapi/resolve_status_test.go`, `resolve.status` carries them.

## Done when

- `closeout.Status` and `resolve.status` give the round, the rounds, and each later round's why, for a running and an ended job.
- A sync for a branch behind main reads as the sync step.
- The tests pass with `flai test` on the changed paths.

## Notes
