---
id: T-1794
type: task
nature: feature
title: The user guides and the dashboard's design say that strategic usage is charged, reported, and charted per role
status: done
parent: S-0425
owner: alex
created: 2026-10-16T05:49:40Z
updated: 2026-11-02T19:42:20Z
transitions:
  - to: ready
    at: 2026-11-02T19:34:44Z
    by: agent-S-0425
  - to: in-progress
    at: 2026-11-02T19:34:47Z
    by: agent-S-0425
  - to: done
    at: 2026-11-02T19:42:20Z
    by: agent-S-0425
stream: S-0425
tags: [cli, dashboard]
touches: [docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/continuous-improvement.md, design/tech/charts.md, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-1792, T-1793]
usage:
  source: log
  seconds: 453
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 80
      output: 26311
      cache_read: 4721064
      cache_write: 124462
      cost: 2.2119
---
# T-1794 The user guides and the dashboard's design say that strategic usage is charged, reported, and charted per role

## Work

- In `docs/users/flai.md`, describe `flai stats` `strategic` per role.
- In `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` § Views, describe the strategic charts per role and the role page's cost.

Last layer.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
