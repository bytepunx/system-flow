---
id: T-2422
type: task
nature: improvement
title: metrics.md's merge queue event log names the preflight and landing-again events, recorded by an ADR
status: done
parent: S-0549
owner: alex
created: 2026-11-02T00:38:31Z
updated: 2026-11-02T00:44:27Z
transitions:
  - to: ready
    at: 2026-11-02T00:39:13Z
    by: agent-S-0549
  - to: in-progress
    at: 2026-11-02T00:39:18Z
    by: agent-S-0549
  - to: done
    at: 2026-11-02T00:44:26Z
    by: agent-S-0549
stream: S-0549
tags: []
touches: [design/system/metrics.md, design/adrs/0237-the-merge-queue-s-event-log-gains-preflight-held-preflight-cleared-and-landing.md, design/adrs/README.md]
usage:
  source: log
  seconds: 308
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 8818
      cache_read: 1754627
      cache_write: 37987
      cost: 0.7575
---
# T-2422 metrics.md's merge queue event log names the preflight and landing-again events, recorded by an ADR

## Work

`design/system/metrics.md` lists every kind of the merge queue event log and the fields each carries; it is the contract between `flai stats` and the dashboard and changes only with an ADR. S-0549 adds `preflight_held` (files, rules), `preflight_cleared` (held_seconds), and `landing_resumed` (batch, stories, waited_seconds), which `flai stats` counts in no metric. Record that with a short ADR refining ADR-0177 and ADR-0234, and add the kinds and fields to the event table, saying they are not counted. It waits for no task.

## Done when

The ADR is written and the event table names the three kinds and their fields. Story criterion 4.

## Notes

Found by T-2387: metrics.md changes only with an ADR.
