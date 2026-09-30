---
id: T-0601
type: task
nature: remediation
title: Record the chart changes in an ADR and the design and user docs
status: done
parent: S-0169
owner: alex
created: 2026-09-30T00:21:39Z
updated: 2026-09-30T00:38:24Z
transitions:
  - to: ready
    at: 2026-09-30T00:21:49Z
    by: agent-S-0169
  - to: in-progress
    at: 2026-09-30T00:36:16Z
    by: agent-S-0169
  - to: done
    at: 2026-09-30T00:38:24Z
    by: agent-S-0169
stream: S-0169
tags: []
touches: [design/system, design/adrs, docs/users]
usage:
  source: log
  seconds: 128
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 13704
      cache_read: 3484085
      cache_write: 48021
      cost: 1.3553
---
# T-0601 Record the chart changes in an ADR and the design and user docs

## Work

Record the removed and replaced charts in an ADR, and update design/system/metrics.md, design/system/flaiover-dashboard.md, and docs/users/flaiover.md to the new charts and titles.

## Done when

- The ADR is accepted and linked from metrics.md.
- No design or user document names a removed chart or an old title.
- flai check --strict passes.

## Notes
