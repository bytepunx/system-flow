---
id: T-3638
type: task
nature: improvement
title: The merge queue's preflight holds the queue on unfixable lint in a shared-set file only when an enqueued story changes it
status: backlog
parent: S-0853
owner: alex
created: 2026-11-27T23:58:46Z
updated: 2026-11-27T23:58:46Z
transitions: []
stream: S-0853
tags: [flai]
touches: [flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go]
after: [T-3631]
---
# T-3638 The merge queue's preflight holds the queue on unfixable lint in a shared-set file only when an enqueued story changes it

## Work

Build rule 2 of T-3631's ADR. ADR-0234's preflight, `preflightLint` and `confirmLint` in `flai/internal/mergequeue/preflight.go`, lints the main checkout's uncommitted files but `d.Rest` before it builds a batch, and holds the queue with one thread to the operator on any finding. An uncommitted digest entry the curator wrote holds the queue the same way it refused the acceptances in I-0411.

Hold only on a finding in a file an enqueued story's change touches. Log a finding in any other file at warn, naming it, and build the batch: the landing, after T-3636, leaves that file uncommitted. Keep the preflight's other checks, flai check's error findings on main, as they are.

It waits for T-3631: it builds the rule that ADR records. It shares no path with T-3636 and runs beside it.

## Done when

- A test in `flai/internal/mergequeue/preflight_test.go` shows an uncommitted digest entry with an MD028 finding that no enqueued story changes holds nothing, and the report names it.
- A test shows a finding in a file an enqueued story changes still holds the queue.
- `flai test` passes on the changed paths.

## Notes
