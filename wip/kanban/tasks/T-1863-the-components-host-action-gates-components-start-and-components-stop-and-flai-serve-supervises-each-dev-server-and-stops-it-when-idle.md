---
id: T-1863
type: task
nature: feature
title: The components host action gates components.start and components.stop, and flai serve supervises each dev server and stops it when idle
status: backlog
parent: S-0438
owner: alex
created: 2026-10-16T21:34:14Z
updated: 2026-10-16T21:34:14Z
transitions: []
stream: S-0438
tags: [cli, dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/hostapi.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/components.go, flai/internal/serve/components_test.go]
after: [T-1860]
---
# T-1863 The components host action gates components.start and components.stop, and flai serve supervises each dev server and stops it when idle

## Work

- `writes.go`: `ActionComponents` in `Actions` with its description. Add the specs `components.start` and `components.stop` `{project?}`, gated by it and journalled with the sub-project and the outcome.
- `hostapi/components.go`: the reads `components.status` (running, URL, port, since, last error) and `components.list` (S-0437's catalogue), merged into the table in `hostapi.go`.
- `serve/components.go`: one dev server per project and sub-project, started with T-1860. A `components.status` call marks it used. It stops after `components.idle` without one, and when `flai serve` stops.
- `serve_actions.go`: `flai serve enable components` and `flai serve components idle <duration>`.

Waits for T-1860, whose `Start` and `Stop` it supervises; not for T-1861's command, which it does not call.

## Done when

Tests with a fake `vite` cover the gate refusing while off, a start and a status, a second start answering the running server, the idle stop, and the stop with `flai serve`. `flai test flai/internal/hostapi flai/internal/serve flai/cmd` passes.

## Notes

The dashboard's `REQUIRED_METHODS` gains these in S-0439, with the page that calls them.
