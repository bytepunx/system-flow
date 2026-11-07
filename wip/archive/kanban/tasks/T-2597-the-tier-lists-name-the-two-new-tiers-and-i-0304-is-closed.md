---
id: T-2597
type: task
nature: improvement
title: The tier lists name the two new tiers, and I-0304 is closed
status: done
parent: S-0592
owner: alex
created: 2026-11-07T22:03:33Z
updated: 2026-11-07T22:24:38Z
transitions:
  - to: ready
    at: 2026-11-07T22:23:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T22:23:52Z
    by: system-flow
  - to: done
    at: 2026-11-07T22:24:38Z
    by: agent-S-0592
stream: S-0592
tags: [docs, issues]
touches: [design/conventions/tooling.md, design/system/devex.md, design/issues/I-0304-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md, design/issues/summary.md]
after: [T-2595, T-2596]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 19330
          cache_write: 573
          cost: 0.005
---
# T-2597 The tier lists name the two new tiers, and I-0304 is closed

## Work

Say in the two places that list what `flai test` runs that a change under `flai/internal/hostapi/` or to `flai/cmd/hostapi.go` runs `flaiover-flai`, flaiover's tests of the built flai, and a change to `docs/operators/settings.md` runs `settings-index`, naming I-0304: the `flai test` bullet under `## Project additions` in `design/conventions/tooling.md`, below the marker, and the close-out row in `design/system/devex.md`. Then close the issue with `flai issue close I-0304 --reason` naming both tiers and the `-short` key checks, which updates `design/issues/summary.md`.

It waits for T-2595 and T-2596, whose tiers it describes as they were built.

## Done when

- Both documents name `flaiover-flai` and `settings-index` and the paths that select each.
- I-0304 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai test design/conventions/tooling.md design/system/devex.md` passes.

## Notes
