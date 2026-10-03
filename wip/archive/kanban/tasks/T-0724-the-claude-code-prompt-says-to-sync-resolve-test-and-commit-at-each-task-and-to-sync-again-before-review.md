---
id: T-0724
type: task
nature: improvement
title: The claude-code prompt says to sync, resolve, test, and commit at each task, and to sync again before review
status: done
parent: S-0197
owner: arobson
created: 2026-10-03T00:40:09Z
updated: 2026-10-03T00:59:21Z
transitions:
  - to: ready
    at: 2026-10-03T00:40:34Z
    by: agent-S-0197
  - to: in-progress
    at: 2026-10-03T00:52:36Z
    by: agent-S-0197
  - to: done
    at: 2026-10-03T00:59:21Z
    by: agent-S-0197
stream: S-0197
tags: []
touches: [flai/internal/harness]
usage:
  source: log
  seconds: 405
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 97
      output: 37539
      cache_read: 3287781
      cache_write: 130574
      cost: 2.2133
---
# T-0724 The claude-code prompt says to sync, resolve, test, and commit at each task, and to sync again before review

## Work

Reword the story agent's prompt in `flai/internal/harness/harness.go` so the per-task cycle reads as the git convention has it (the order is settled on TH-0072): commit the task, `flai stream sync`, resolve the conflicts it reports, run the tests for what the task changed, commit any fix; and before review, commit everything, sync again, and resolve. Update the prompt's tests in `harness_test.go`. Waits for nothing: it touches only the harness.

## Done when

- [ ] The `claude-code` prompt names the per-task cycle and the sync before review
- [ ] `harness_test.go` asserts it

## Notes
