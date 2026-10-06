---
id: T-0992
type: task
nature: remediation
title: flai move warns when a story goes to review with unticked criteria
status: done
parent: S-0282
owner: alex
created: 2026-10-06T03:47:48Z
updated: 2026-10-06T03:54:09Z
transitions:
  - to: ready
    at: 2026-10-06T03:48:16Z
    by: agent-S-0282
  - to: in-progress
    at: 2026-10-06T03:48:16Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T03:54:09Z
    by: agent-S-0282
stream: S-0282
tags: [cli]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go]
usage:
  source: log
  seconds: 353
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 6405
      cache_read: 1025976
      cache_write: 30197
      cost: 0.5152
---
# T-0992 flai move warns when a story goes to review with unticked criteria

## Work

In `flai/internal/workitem/rules.go`, a story moving to `review` with an unticked criterion gets a warning, not a refusal, naming the unticked criteria and `flai criteria tick`. A criterion that could not be verified stays unticked (work-management.md), so the move is not refused.

Waits for nothing: no path in common with another task.

## Done when

- [ ] The warning is given on a move to review with an unticked criterion, and not when every criterion is ticked, with a test.
- [ ] `go test -race ./internal/workitem/` passes.

## Notes
