---
id: T-1483
type: task
nature: improvement
title: flai task done leaves the paths sync regenerates out of its changed-paths notice
status: done
parent: S-0372
owner: alex
created: 2026-10-12T17:22:53Z
updated: 2026-10-14T21:13:23Z
transitions:
  - to: ready
    at: 2026-10-14T20:57:34Z
    by: agent-S-0372
  - to: in-progress
    at: 2026-10-14T20:57:35Z
    by: agent-S-0372
  - to: done
    at: 2026-10-14T21:13:23Z
    by: agent-S-0372
stream: S-0372
tags: [messages, cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
usage:
  source: log
  seconds: 948
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 72
      output: 398
      cache_read: 2373865
      cache_write: 78534
      cost: 1.1291
---
# T-1483 flai task done leaves the paths sync regenerates out of its changed-paths notice

## Work

In `flai/internal/taskdone/taskdone.go`, `(*run).tell` passes every committed path to `itemedit.Covering` and then to `messages.Notify`. Before `Covering`, drop the paths that sync regenerates, which `issues.Generated` gives as `storygit.GeneratedFile` entries. `sync()` already reads them through `issues.SyncFiles`. Drop them first because `Covering` tells a story with an empty claim about every path. A commit that changes only such paths tells no one.

This task waits for nothing. It shares no path with the other tasks, and the ADR task records the rule without changing what this code needs.

## Done when

- A test in `flai/internal/taskdone/taskdone_test.go` reproduces I-0125's third instance and fails without the change. In it, a task's commit changes `design/issues/summary.md` and a file another open story claims. The other story is told about the file alone. A commit that changes only `design/issues/summary.md` writes no notice and leaves `Result.Told` empty.
- `TestRunTellsEachOpenStoryWhoseClaimCoversTheCommit` and the other tests of the package still pass.
- `flai test` passes on the two files.

## Notes

S-0367, in progress, also touches `flai/internal/taskdone/taskdone.go`. Its overlap holds S-0372 in ready until S-0367 is accepted, so this change starts from S-0367's.
