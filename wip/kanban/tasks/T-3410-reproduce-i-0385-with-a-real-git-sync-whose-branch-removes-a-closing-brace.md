---
id: T-3410
type: task
nature: remediation
title: Reproduce I-0385 with a real-git sync whose branch removes a closing brace elsewhere in a Go file main added one to
status: done
parent: S-0800
owner: alex
created: 2026-11-16T01:45:03Z
updated: 2026-11-16T02:19:26Z
transitions:
  - to: ready
    at: 2026-11-16T01:47:32Z
    by: agent-S-0800
  - to: in-progress
    at: 2026-11-16T02:10:11Z
    by: agent-S-0800
  - to: done
    at: 2026-11-16T02:19:26Z
    by: agent-S-0800
stream: S-0800
tags: [storygit, go, test]
touches: [flai/internal/storygit/sync_test.go]
after: [T-3407]
usage:
  source: log
  seconds: 555
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 5141
      cache_read: 1435471
      cache_write: 53088
      cost: 0.7333
---
# T-3410 Reproduce I-0385 with a real-git sync whose branch removes a closing brace elsewhere in a Go file main added one to

## Work

Add a real-git test to `flai/internal/storygit/sync_test.go`, beside S-0767's I-0360 reproduction, that builds I-0385's shape. A story branch replaces a loop in a Go file, which removes the loop's closing brace. Main then adds a function to the same file, ending in a tab and a closing brace. `flai stream sync`'s rebase is clean. The test checks that the sync reports nothing under `UndoesMain`, that the story's undone record holds no path, and that a record written by the old rule for that path drops out when `RecordedUndone` reads it, as S-0739's will.

Add the I-0124 shape alongside, or check that an existing test covers it: a branch that changed a file and reverted it, rebased over main's identical change, is still reported and recorded.

It waits for T-3407: the test passes only with the new count in place.

## Done when

- `flai test flai/internal/storygit/sync_test.go` passes with the I-0385 test, and that test fails against the multiset count of `undone.go` before T-3407.
- The I-0124 revert is still reported by a real-git test.

## Notes
