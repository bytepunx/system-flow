---
id: T-3292
type: task
nature: remediation
title: Reproduce I-0360 with a real-git sync that stops on one long markdown line and resumes with main's words kept
status: done
parent: S-0767
owner: alex
created: 2026-11-13T11:48:23Z
updated: 2026-11-13T12:40:00Z
transitions:
  - to: ready
    at: 2026-11-13T12:14:43Z
    by: agent-S-0767
  - to: in-progress
    at: 2026-11-13T12:14:44Z
    by: agent-S-0767
  - to: done
    at: 2026-11-13T12:40:00Z
    by: agent-S-0767
stream: S-0767
tags: [flai, git, test]
touches: [flai/internal/storygit/sync_test.go]
after: [T-3291]
usage:
  source: log
  seconds: 1516
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 9736
      cache_read: 3171336
      cache_write: 76968
      cost: 1.3588
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 20200
          cache_write: 757
          cost: 0.0061
---
# T-3292 Reproduce I-0360 with a real-git sync that stops on one long markdown line and resumes with main's words kept

## Work

In `flai/internal/storygit/sync_test.go`, beside `TestSyncFindsNothingUndoneOnABranchThatKeepsMainsChanges`, add a real-git test of I-0360's instance, with `syncProject` and `openStory`:

1. A markdown file under the story's touches holds one long table row on a single line.
2. The branch extends the row with a sentence and commits; main changes the same row differently and commits.
3. `Sync` stops on the file. Resolve it as S-0761 did: main's row with the branch's sentence added, so the branch keeps every word of main's. `git add` it and run `Sync` again, which resumes the rebase (ADR-0294).
4. The result's `UndoesMain` is empty and no record is written under `UndoneRecordPath`.

Add the counter-case: the same steps resolved with the branch's row alone, dropping main's words, still report the path and record it. The test runs under the integration tier, as the other real-git sync tests do; skip it under `-short` the way they skip. It waits for T-3291, since it fails until the check keeps an extended line.

## Done when

- [ ] The new test fails on the code before T-3291 and passes after it
- [ ] The counter-case reports the path in `UndoesMain` and in the story's record
- [ ] `flai test flai/internal/storygit/sync_test.go` and the integration tier pass

## Notes
