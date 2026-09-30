---
id: T-0597
type: task
nature: feature
title: Time in state is drawn on a time axis that spans the window
status: in-progress
parent: S-0168
owner: alex
created: 2026-09-29T23:58:07Z
updated: 2026-09-30T00:14:07Z
transitions:
  - to: ready
    at: 2026-09-29T23:58:27Z
    by: agent-S-0168
  - to: in-progress
    at: 2026-09-30T00:14:07Z
    by: agent-S-0168
stream: S-0168
tags: []
touches: [flaiover/src, design, docs]
---

# T-0597 Time in state is drawn on a time axis that spans the window

## Work

- Draw time in state on a time axis from the window's start to the report's now, in the layout the designer chooses in TH-0039.
- Record the change to the chart's definition in an ADR refining ADR-0054, and update `design/system/metrics.md` and `docs/users/flaiover.md`.
- Unit and page tests.

## Done when

- The time in state chart's x axis runs from the window's start to now for every window, and changes when the window does, in tests and in a browser against this repository's stats.
- `scripts/flaiover-test.sh` passes; `flai check --strict` is clean.

## Notes
