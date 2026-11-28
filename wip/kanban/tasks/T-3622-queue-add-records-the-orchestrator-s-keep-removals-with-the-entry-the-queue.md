---
id: T-3622
type: task
nature: improvement
title: queue_add records the orchestrator's keep_removals with the entry, the queue leaves a kept path out of undone_main, and its landing passes the keeps
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:57:59Z
updated: 2026-11-27T23:57:59Z
transitions: []
stream: S-0852
tags: [flai]
touches: [flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mcpserver/queue.go, flai/internal/mcpserver/queue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go]
after: [T-3620]
---
# T-3622 queue_add records the orchestrator's keep_removals with the entry, the queue leaves a kept path out of undone_main, and its landing passes the keeps

## Work

In I-0415, S-0811 was ineligible for the queue on `undone_main`, and the orchestrator could only enqueue it. Carry the orchestrator's keeps through the queue as the ADR of T-3618 says.

- `flai queue add --by orchestrator` takes `--keep-removal <path> --reason "<why>"`, and the MCP tool `queue_add` takes `keep_removals`. Both check the keeps as `flai accept --by orchestrator` does after T-3620: the permission on, a path the record holds, a reason, and the narrative's `## Decisions` naming the path.
- The entry in `flai/internal/mergequeue/state.go` records them beside `verified` and `evidence`.
- `undoesMain` in `flai/internal/mergequeue/eligible.go` leaves out the paths the entry keeps. A path still undone and not kept remains `undone_main`.
- The batch's landing, `flai accept --batch` in `flai/cmd/accept.go`, reads the entry's keeps and writes them under `### Removals kept` with `by: orchestrator`.

Test it in `flai/internal/mergequeue/eligible_test.go`: an entry whose keeps cover the recorded paths is eligible, and one that keeps only some is not. Test in `flai/cmd/accept_queue_test.go` that a batch lands such a story with its kept reasons in its Notes.

It waits for T-3620, whose keep check and `keptNotes` it reuses. It shares no file with T-3621 and runs beside it.

## Done when

- The two tests pass, and a `queue_add` test shows `keep_removals` recorded and refused with the permission off.
- `flai test` passes on the changed paths.

## Notes
