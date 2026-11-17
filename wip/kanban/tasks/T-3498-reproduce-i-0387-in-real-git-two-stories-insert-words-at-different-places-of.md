---
id: T-3498
type: task
nature: improvement
title: "Reproduce I-0387 in real git: two stories insert words at different places of one table row and one sentence, and neither sync nor acceptance stops"
status: backlog
parent: S-0817
owner: alex
created: 2026-11-17T02:31:55Z
updated: 2026-11-17T02:31:55Z
transitions: []
stream: S-0817
tags: [flai, merge, test]
touches: [flai/cmd/accept_conflict_test.go, flai/cmd/stream_sync_test.go]
after: [T-3496]
---
# T-3498 Reproduce I-0387 in real git: two stories insert words at different places of one table row and one sentence, and neither sync nor acceptance stops

## Work

Beside S-0719's reproductions of I-0346 (`flai/cmd/stream_sync_test.go` near its I-0346 cases, `flai/cmd/accept_conflict_test.go` near its ADR-0287 cases), add real-git tests of S-0797's stop:

- A design document with a commands table and a one-line paragraph. Story S-0001 inserts words in one place of a table row and of a sentence; main, through another accepted story, inserts words in other places of the same row and sentence.
- `flai stream sync` on S-0001 ends clean. The file holds both sides' words, and the sync records no undone line (ADR-0297, ADR-0309).
- `flai accept` of S-0001, after the other is accepted, lands with no send-back and no thread.
- The stop stays where both stories insert at the same point of the row: sync stops naming the section, and acceptance sends the story back as today.

Cite I-0387 and the ADR in the test comments.

Waits for T-3496: the tests pass only with its merge.

## Done when

- The new tests fail on main and pass on the story's branch.
- `go test ./cmd/ -run` with their names passes, and `flai test` on the two files is clean.

## Notes
