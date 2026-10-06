---
id: T-1016
type: task
nature: remediation
title: The orchestrator stub prints its stream before it says what it was given, so a stop cannot come before its usage is in the run's log
status: done
parent: S-0292
owner: alex
created: 2026-10-06T11:34:08Z
updated: 2026-10-06T11:36:56Z
transitions:
  - to: ready
    at: 2026-10-06T11:34:28Z
    by: agent-S-0292
  - to: in-progress
    at: 2026-10-06T11:34:28Z
    by: agent-S-0292
  - to: done
    at: 2026-10-06T11:36:43Z
    by: agent-S-0292
stream: S-0292
tags: []
touches: [design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/I-0090-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md, design/issues/summary.md, flai/internal/serve/orchestrate_test.go]
usage:
  source: log
  seconds: 135
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 7283
      cache_read: 1763226
      cache_write: 49294
      cost: 0.8928
---
# T-1016 The orchestrator stub prints its stream before it says what it was given, so a stop cannot come before its usage is in the run's log

## Work

The orchestrator stub in `flai/internal/serve/orchestrate_test.go` writes what it was given to `orchestrator-<pid>.txt` and only then prints `stream-orchestrator`, the usage the run's activity is measured from. `TestTheOrchestratorIsStoppedWhenTheActionIsTurnedOff` stops the run as soon as `given` sees that file. Under load the stop's TERM can land before `cat` has written the stream, so the run's log holds no usage, `logRunEnd` logs no activity, and the test waits out its 5 s for an entry that never comes (I-0090).

Print the stream first, then write what it was given, so that `given` returning means the usage is in the log. Show the cause with a delay before `cat` in the old order, and that the new order holds with the same delay. Close I-0090. It waits for no other task.

## Done when

- The stub prints its stream before it writes `orchestrator-<pid>.txt`, and its doc comment says why.
- With a delay injected before the stream, the test fails in the old order and passes in the new one, recorded in the narrative.
- The serve package's orchestrator tests pass with `-race`.
- I-0090 is closed with a reason saying what fixed it.

## Notes
