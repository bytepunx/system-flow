---
id: T-2202
type: task
nature: remediation
title: flai serve raises on a thread a curator retry that ends on the usage limit again, or a limit whose reset time it cannot read
status: backlog
parent: S-0514
owner: alex
created: 2026-10-28T12:54:11Z
updated: 2026-10-28T12:54:11Z
transitions: []
stream: S-0514
tags: [flai, serve, curator]
touches: [flai/internal/serve/curate.go, flai/internal/serve/curate_test.go]
after: [T-2201]
---
# T-2202 flai serve raises on a thread a curator retry that ends on the usage limit again, or a limit whose reset time it cannot read

## Work

Waits for T-2201, which marks a retry on the curator's run.

In `flai/internal/serve/curate.go`, when `curateEnded` records a run that ended on the usage limit and the run is a retry (T-2201's mark), or the limit names no reset time, open one thread to the operator, as `stopRestarting` in `agents.go` and `claudeChecker.tell` in `claudecheck.go` do: on `manifest.File` (`system-flow.yaml`), since the curator has no story, by `flai serve` (`limitAuthor`), titled such as `The curator ended on the harness's usage limit twice; start it when the limit allows`, with a text that names the run (its start, its trigger, the harness's text, the reset time or that none was read), what it means (no digest entry is written until the curator runs), and what starts the curator again: `flai curate`, the next acceptance, or `digest.schedule`. Record the thread's ID in the run's `LimitThread` so it is opened once; journal the opening with a `curate` action entry as the start is journalled; when the thread cannot be opened, warn and record it in the journal as `failed`, as `stopRestarting` does. The scheduler schedules no further retry for a retry that ended on the limit (T-2201 already leaves it alone).

Tests in `curate_test.go`: a retry that ends on the limit again opens exactly one thread on `system-flow.yaml`, by `flai serve`, naming the reset time, recorded on the run; a first run on the limit, with a reset time, opens none; a run on the limit with no reset time opens one; a thread that cannot be opened is warned of and journalled.

## Done when

- A curator retry that ends on the usage limit again, and a run whose limit names no reset time, open one thread to the operator on `system-flow.yaml`, recorded on the run, journalled, and opened once.
- A first run on the limit with a reset time opens no thread.
- `flai test` on the task's paths passes.

## Notes

Layer 4 of S-0514's plan. `curate.go` and `curate_test.go` are shared with T-2200 and T-2201 in the layers before.
