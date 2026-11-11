---
id: T-3102
type: task
nature: improvement
title: "Reproduce I-0346 in real git: two stories change different sentences of one paragraph of a design document, and neither sync nor acceptance stops"
status: done
parent: S-0719
owner: alex
created: 2026-11-11T08:54:58Z
updated: 2026-11-11T14:50:41Z
transitions:
  - to: ready
    at: 2026-11-11T14:40:26Z
    by: agent-S-0719
  - to: in-progress
    at: 2026-11-11T14:40:26Z
    by: agent-S-0719
  - to: done
    at: 2026-11-11T14:50:41Z
    by: agent-S-0719
stream: S-0719
tags: [flai, go, test]
touches: [flai/cmd/accept_conflict_test.go, flai/cmd/stream_sync_test.go]
after: [T-3097]
usage:
  source: log
  seconds: 615
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 10605
      cache_read: 2161730
      cache_write: 88668
      cost: 1.2407
---
# T-3102 Reproduce I-0346 in real git: two stories change different sentences of one paragraph of a design document, and neither sync nor acceptance stops

## Work

Add two real-git tests, on the fixtures S-0637's I-0315 tests beside them use, that reproduce I-0346's instance:

- In `flai/cmd/accept_conflict_test.go`: two stories each change a different sentence of one one-line paragraph of a design document under the design folder, as S-0666 and S-0688 did in `design/system/agent-adapters.md`. The first is accepted. The second's acceptance lands, with main's paragraph and the second story's sentence applied, and the story is not sent back to `in-progress`. Beside it, assert that both stories changing the same sentence of that paragraph is still sent back, as `TestAnAcceptanceByHandSendsBackAStoryWhoseRebaseItCannotSettle` does.
- In `flai/cmd/stream_sync_test.go`: the same two changes, the first story accepted, and `flai stream sync` on the second ends clean, with both sentences in the paragraph and no conflict left in its worktree.

Name each test for what it shows and cite I-0346 and the new ADR in its comment. Both skip under `-short`, as the real-git tests do.

Waits for T-3097: these tests fail until its sentence merge settles the paragraph.

## Done when

- Both tests pass with T-3097's change, and fail without it.
- `go test ./cmd/ -run` on the two tests passes, and `scripts/flai-test.sh` lint is clean.

## Notes
