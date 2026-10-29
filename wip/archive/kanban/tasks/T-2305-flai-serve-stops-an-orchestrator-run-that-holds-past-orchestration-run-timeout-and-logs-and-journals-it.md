---
id: T-2305
type: task
nature: feature
title: flai serve stops an orchestrator run that holds past orchestration.run_timeout and logs and journals it
status: done
parent: S-0536
owner: alex
created: 2026-10-28T21:46:15Z
updated: 2026-10-29T08:01:38Z
transitions:
  - to: ready
    at: 2026-10-29T07:46:28Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T07:46:29Z
    by: agent-S-0536
  - to: done
    at: 2026-10-29T08:01:38Z
    by: agent-S-0536
stream: S-0536
tags: [flai, serve, orchestrator]
touches: [flai/internal/serve/orchestrate_timeout.go, flai/internal/serve/orchestrate_timeout_test.go, flai/internal/serve/orchestrate.go]
after: [T-2302]
usage:
  source: log
  seconds: 909
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 12500
      cache_read: 2827972
      cache_write: 90558
      cost: 1.3699
---
# T-2305 flai serve stops an orchestrator run that holds past orchestration.run_timeout and logs and journals it

## Work

ADR-0209 §4, criterion 3 (the timeout). At each look, `orchestrator.look` in `flai/internal/serve/orchestrate.go` stops a run going for longer than `orchestration.run_timeout` (20m by default), as it stops one when the action is turned off. Record it with a stop reason of its own, such as `stopped: past orchestration.run_timeout (20m)`. Put it in the run's `why`, the `orchestrate` journal entry, and the summary of its activity in `wip/agents/orchestrator.md`. A stop for the timeout is not a failure: it does not wait `orchestrateRetry`, and the reasons still queued start the next run as usual. Put the check and its summary in a new `orchestrate_timeout.go`.

Waits for T-2302, which rewrites `look` around the queue.

## Done when

- A test in `orchestrate_timeout_test.go` runs a fake orchestrator past a short `run_timeout` and finds it stopped, the journal entry and the activity summary naming the timeout, and the next queued reason starting a run.
- A run within the timeout is left alone.
- `flai test` passes on the changed paths.

## Notes
