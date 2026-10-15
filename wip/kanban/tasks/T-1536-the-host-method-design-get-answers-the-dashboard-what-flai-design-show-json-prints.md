---
id: T-1536
type: task
nature: feature
title: The host method design.get answers the dashboard what flai design show --json prints
status: in-progress
parent: S-0381
owner: alex
created: 2026-10-13T17:46:45Z
updated: 2026-10-15T00:30:29Z
transitions:
  - to: ready
    at: 2026-10-15T00:30:28Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:30:29Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/hostapi/hostapi.go]
after: [T-1533, T-1590]
---
# T-1536 The host method design.get answers the dashboard what flai design show --json prints

## Work

Add `designMethods()` in `flai/internal/hostapi/design.go` with the read `design.get`, answering `designsystem.Show` for the project, and merge it into the method table in `hostapi.go`. Test it with a document, without one, and with a token that does not parse. It waits for T-1533, whose `Show` it answers.

## Done when

- `go test ./internal/hostapi/` passes, the contract tests included.

## Notes

- Layer 3, beside T-1535. The dashboard's side, `REQUIRED_METHODS` and the endpoint, is S-0382's.
