---
id: T-2386
type: task
nature: improvement
title: flai serve journals the preflight hold, opens one thread to the operator naming the files, resumes the queue on the change that clears it, and journals a landing again
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:41:58Z
updated: 2026-10-31T07:41:58Z
transitions: []
stream: S-0549
tags: [flai]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/serve/serve.go]
after: [T-2383, T-2384]
---
# T-2386 flai serve journals the preflight hold, opens one thread to the operator naming the files, resumes the queue on the change that clears it, and journals a landing again

## Work

In `flai/internal/serve/mergequeue.go`, when a look finds the queue held by the preflight:

- Journal it once per hold through `p.journal` (action `merge-queue`), naming each file and finding, as `notRun` logs a reason once.
- Open one thread to the operator with `threads.New`, anchored where the ADR says, naming each file, line, and rule and saying the queue builds no batch until they are fixed and that the queue changes none of them. Reuse it while it is open, record its ID on the hold so the board can link it, and resolve it, with a reason, when a clean look clears the hold.
- Journal a batch the landing task lands again, saying it landed with no new verification and from which commits.

In `serve.go`, the watcher wakes the queue only for paths under `/kanban/` and `/threads/`, and the 30 s tick does the rest. While the preflight holds the queue, wake it on any change in the main checkout, `wip/` included, and on a move of main's head, so the fix the operator makes resumes it within seconds and with nothing else asked of them.

Tests in `mergequeue_test.go`: a held look journals once and opens one thread, a second look opens none, a fix to the file wakes the queue, clears the hold, resolves the thread, and builds the batch. Story criteria 1 (journal and thread) and 2.

It waits for the preflight task, whose hold it reads, and the landing task, whose landing again it journals; it is the last code task. S-0548 changes other files of `flai/internal/serve` (`agents_test.go`, `serve_test.go`), not these.

## Done when

The serve journal names the held files and findings once, one thread to the operator is open while the hold stands and is resolved when it clears, a change that clears the finding resumes the queue without any other action, a landing again is journaled, and `go test -race ./internal/serve` passes.

## Notes

Drafted by the planner of S-0549 on 2026-10-31.
