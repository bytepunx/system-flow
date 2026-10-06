---
id: T-0906
type: task
nature: feature
title: flai stats counts the thread waits the orchestrator ended apart from the operator's, and a recommendation ends no wait
status: done
parent: S-0220
owner: alex
created: 2026-10-05T04:47:46Z
updated: 2026-10-06T06:19:41Z
transitions:
  - to: ready
    at: 2026-10-06T06:06:49Z
    by: agent-S-0220
  - to: in-progress
    at: 2026-10-06T06:06:50Z
    by: agent-S-0220
  - to: done
    at: 2026-10-06T06:19:41Z
    by: agent-S-0220
stream: S-0220
tags: [flai]
touches: [flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, design/system/metrics.md, design/adrs]
after: [T-0894]
usage:
  source: log
  seconds: 771
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 102
      output: 39265
      cache_read: 4862921
      cache_write: 145053
      cost: 2.5947
---
# T-0906 flai stats counts the thread waits the orchestrator ended apart from the operator's, and a recommendation ends no wait

## Work

`threadWait` (`flai/internal/metrics/waiting.go`, line 77) ends a thread's wait at the first later entry by another author. With T-0894's marks:

- A recommendation ends no wait. The wait runs on to the operator's confirmation, or to the next answer.
- Each wait records who ended it: the orchestrator, when its answer ended it; the operator, when the operator's confirmation of a recommendation ended it, counted as confirmed; anyone else as today. The orchestrator is the author name S-0218 gives its run.
- `items[]` gains `wait_threads_orchestrator_seconds`, the part of `wait_threads_seconds` whose waits the orchestrator ended. `waiting.weeks[].threads` gains `orchestrator` and `confirmed`, each with `count` and `total_seconds`. The existing fields keep their meaning, so the dashboard's charts read them as before.

`design/system/metrics.md` is the contract between `flai stats` and the dashboard, and CLAUDE.md allows changing it only with an ADR. Record the change as an ADR (`flai adr new`) and update `### Waiting` under `## Planning, waiting, and claims (S-0205)`.

It waits for T-0894, whose marks tell a recommendation from an answer. It runs with T-0902, whose paths it does not share.

## Done when

- a test with a fixture thread finds a recommendation ending no wait, an orchestrator answer counted under `orchestrator`, and a confirmed recommendation counted under `confirmed`
- a thread with no orchestrator entry gives the same figures as before
- the ADR is proposed, and `metrics.md` defines the new fields
- `go test ./internal/metrics/` passes

## Notes
