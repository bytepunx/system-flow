---
id: T-0970
type: task
nature: feature
title: flaiover-dashboard.md and the user guide describe the cost of delay charts
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:47:28Z
updated: 2026-10-07T08:04:26Z
transitions:
  - to: ready
    at: 2026-10-07T08:00:08Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T08:00:08Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T08:04:26Z
    by: agent-S-0213
stream: S-0213
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0967]
usage:
  source: log
  seconds: 258
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 15428
      cache_read: 2612904
      cache_write: 65101
      cost: 1.2219
---
# T-0970 flaiover-dashboard.md and the user guide describe the cost of delay charts

## Work

Describe the three charts where the others are described:

- `design/system/flaiover-dashboard.md`:
  - a row per chart in the Views table, after the usage charts, each naming what it draws, the `flai stats --json` keys it reads with a link to `metrics.md`, and the ADR T-0953 wrote
  - the paragraph on the chart groups names the planning group and that `cod-order`'s axis is the projection's, not the window's
- `docs/users/flaiover.md`: a Planning table under `## Charts`, beside Flow and Usage, saying what each chart shows, how to read the saving, and what the count without a value means

It waits for T-0967, whose kinds and behaviour it describes. It runs with the page task, whose paths it does not share.

## Done when

- both documents describe the three charts and the planning group
- `flai check --strict` and the markdown lint pass

## Notes
