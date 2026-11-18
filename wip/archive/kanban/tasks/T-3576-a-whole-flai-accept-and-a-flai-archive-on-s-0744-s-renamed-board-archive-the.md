---
id: T-3576
type: task
nature: remediation
title: A whole flai accept and a flai archive on S-0744's renamed board archive the shipped story
status: done
parent: S-0842
owner: alex
created: 2026-11-17T12:29:45Z
updated: 2026-11-18T08:05:16Z
transitions:
  - to: ready
    at: 2026-11-18T07:58:23Z
    by: agent-S-0842
  - to: in-progress
    at: 2026-11-18T07:58:23Z
    by: agent-S-0842
  - to: done
    at: 2026-11-18T08:05:16Z
    by: agent-S-0842
stream: S-0842
tags: [flai, workflow, test]
touches: [flai/cmd/accept_workflow_test.go, flai/cmd/archive_test.go]
after: [T-3572]
usage:
  source: log
  seconds: 413
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 6500
      cache_read: 2643878
      cache_write: 79315
      cost: 1.173
---
# T-3576 A whole flai accept and a flai archive on S-0744's renamed board archive the shipped story

## Work

Reproduce I-0377 at the command line. In `flai/cmd/accept_workflow_test.go`, on the board `renamedProject` (`flai/cmd/move_workflow_test.go`) builds, whose accepted state is `shipped`, run a whole `flai accept` of a story in review: it walks the story to `shipped`, archives it and its tasks with the narrative, and commits. Today's test checks only the refusal and the dry run. In `flai/cmd/archive_test.go`, `flai archive` with an ID and without archives a shipped story there. Each test fails at the base of the story and passes after T-3572. Waits for T-3572, whose `Closed` and `PlanArchive` it tests.

## Done when

- The whole acceptance on the renamed board archives the story; run against the code before T-3572, it fails with the refusal I-0377 quotes.
- `flai archive` on the renamed board archives a shipped story, with an ID and without.
- `flai test` on the changed paths passes.

## Notes
