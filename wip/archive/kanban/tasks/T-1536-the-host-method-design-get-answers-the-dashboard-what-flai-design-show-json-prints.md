---
id: T-1536
type: task
nature: feature
title: The host method design.get answers the dashboard what flai design show --json prints
status: done
parent: S-0381
owner: alex
created: 2026-10-13T17:46:45Z
updated: 2026-10-15T00:54:59Z
transitions:
  - to: ready
    at: 2026-10-15T00:30:28Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:30:29Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:54:59Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/hostapi/hostapi.go, flaiover/src/lib/server/agent.ts]
after: [T-1533, T-1590]
usage:
  source: log
  seconds: 1470
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 15785
      cache_read: 1482739
      cache_write: 71683
      cost: 1.0061
---
# T-1536 The host method design.get answers the dashboard what flai design show --json prints

## Work

Add `designMethods()` in `flai/internal/hostapi/design.go` with the read `design.get`, answering `designsystem.Show` for the project, and merge it into the method table in `hostapi.go`. Test it with a document, without one, and with a token that does not parse. It waits for T-1533, whose `Show` it answers.

## Done when

- `go test ./internal/hostapi/` passes, the contract tests included.

## Notes

- Layer 3, beside T-1535. The dashboard's side, `REQUIRED_METHODS` and the endpoint, is S-0382's.
