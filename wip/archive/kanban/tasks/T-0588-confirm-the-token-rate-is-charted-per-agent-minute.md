---
id: T-0588
type: task
nature: research
title: Confirm the token rate is charted per agent minute
status: done
parent: S-0165
owner: alex
created: 2026-09-29T22:41:06Z
updated: 2026-09-29T22:41:14Z
transitions:
  - to: ready
    at: 2026-09-29T22:41:14Z
    by: agent-S-0165
  - to: in-progress
    at: 2026-09-29T22:41:14Z
    by: agent-S-0165
  - to: done
    at: 2026-09-29T22:41:14Z
    by: agent-S-0165
stream: S-0165
tags: []
usage:
  source: log
  seconds: 0
  models: []
---

# T-0588 Confirm the token rate is charted per agent minute

## Work

Find what charts or prints a token rate per agent hour in `flaiover/src`, `flai/cmd`, and the running dashboard's bundle, and ask the designer what showed it.

## Done when

- [x] Every token rate the dashboard and `flai stats` show is per agent minute, or what is not is changed.

## Notes

All per agent minute since S-0163: `tokenRate` in `flaiover/src/lib/viz/charts.ts`, `SpendTable.svelte`, the charts page's summary, `flai stats`, and the container's bundle (revision 5fe2577). The designer's tab held the old bundle; a hard refresh fixed it (TH-0038).
