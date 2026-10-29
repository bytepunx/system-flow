---
id: T-2272
type: task
nature: improvement
title: flai serve appends a started, ended, failed, or stopped agent event for each story's agent run and for any run the operator stops
status: done
parent: S-0532
owner: alex
created: 2026-10-28T21:08:56Z
updated: 2026-10-29T02:36:42Z
transitions:
  - to: ready
    at: 2026-10-29T02:29:30Z
    by: agent-S-0532
  - to: in-progress
    at: 2026-10-29T02:29:31Z
    by: agent-S-0532
  - to: done
    at: 2026-10-29T02:36:42Z
    by: agent-S-0532
stream: S-0532
tags: [flai, serve, events]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/stop.go, flai/internal/serve/stop_test.go]
after: [T-2274]
usage:
  source: log
  seconds: 431
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 21605
      cache_read: 3360009
      cache_write: 119902
      cost: 1.8126
---
# T-2272 flai serve appends a started, ended, failed, or stopped agent event for each story's agent run and for any run the operator stops

## Work

Call the serve helper from the story's agent's run path in `flai/internal/serve/agents.go`:

- `started` when `start()` has spawned the agent, with the start's reason as `why` (the operator's start, ready, a restart on an answer, or an automatic restart with its count);
- `failed` when `start()` could not start it (`fail()`), with that error;
- on `endRun()`, after `judgeRun()`, `ended` for a run judged `worked` or `asked` and `failed` for one judged `failed`, with the run's `Why`;
- on `settleOrphans()`, the same for a run whose process was gone when flai serve looked.

In `stop.go`, `Stop()` appends `stopped`, with the operator's name, for a story's agent and for a role's run alike: the stop path is shared, so the roles' task leaves it to this one.

Each event is written once: a run judged twice, or settled after it ended, writes nothing more. This task waits for T-2271, which gives the event, its log, and the helper.

## Done when

- [ ] A test per `to` drives a story's agent run through start, end, a failed start, and an operator stop and reads exactly one event of each back from the log, with its item, `to`, and `why`
- [ ] An orphaned run settled at a look writes its end once
- [ ] `flai test` passes on the changed paths

## Notes

Planned by planner-S-0532.
