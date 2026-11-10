---
id: T-2814
type: task
nature: improvement
title: "Reproduce I-0315 in real git: two stories add a section at one spot of a design document outside claims.shared, and neither sync nor acceptance stops"
status: done
parent: S-0637
owner: alex
created: 2026-11-09T22:09:37Z
updated: 2026-11-10T08:12:57Z
transitions:
  - to: ready
    at: 2026-11-10T08:06:39Z
    by: agent-S-0637
  - to: in-progress
    at: 2026-11-10T08:06:40Z
    by: agent-S-0637
  - to: done
    at: 2026-11-10T08:12:57Z
    by: agent-S-0637
stream: S-0637
tags: [flai, go, test]
touches: [flai/cmd/accept_conflict_test.go, flai/cmd/stream_sync_test.go]
after: [T-2812]
usage:
  source: log
  seconds: 377
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 59
      output: 8648
      cache_read: 8849489
      cache_write: 70710
      cost: 2.4706
---
# T-2814 Reproduce I-0315 in real git: two stories add a section at one spot of a design document outside claims.shared, and neither sync nor acceptance stops

## Work

Add two real-git tests, on the fixtures the files beside them use, that reproduce I-0315's instances:

- In `flai/cmd/accept_conflict_test.go`: two stories each add a section after the same section of `design/system/server-performance.md`, a file not in `claims.shared`. The first is accepted. The second's acceptance lands, with both sections in the file, main's first, and is not sent back to `in-progress`. Beside it, assert that a stop where both stories change the same lines of one section in that file is still sent back, as `TestAnAcceptanceByHandSendsBackAStoryWhoseRebaseItCannotSettle` does today.
- In `flai/cmd/stream_sync_test.go`: the same two sections, the first story accepted, and `flai stream sync` on the second ends clean, with both sections and no conflict left in its worktree.

Name each test for what it shows and cite I-0315 and the new ADR in its comment. Both skip under `-short`, as the real-git tests do.

Waits for T-2812: these tests fail until its matcher takes the file.

## Done when

- Both tests fail with T-2812's change reverted and pass with it.
- `flai test flai/cmd/accept_conflict_test.go flai/cmd/stream_sync_test.go` passes, and the integration tier runs the two tests.

## Notes
