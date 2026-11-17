---
id: T-3258
type: task
nature: feature
title: flaiover-dashboard.md says the dashboard's views read a state's kind from flai
status: done
parent: S-0758
owner: alex
created: 2026-11-12T16:10:07Z
updated: 2026-11-17T08:25:03Z
transitions:
  - to: ready
    at: 2026-11-17T08:24:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:24:53Z
    by: system-flow
  - to: done
    at: 2026-11-17T08:25:03Z
    by: agent-S-0758
stream: S-0758
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md]
after: [T-3252]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1567
      cache_read: 346955
      cache_write: 7573
      cost: 0.148
---
# T-3258 flaiover-dashboard.md says the dashboard's views read a state's kind from flai

## Work

Waits for T-3252: it describes the module that task adds. Say in `design/system/flaiover-dashboard.md` § Internal structure that the views read a state's kind, colour, and moves from `flaiover/src/lib/workflow.ts`, built from `board.get`, and compare no state's name.

## Done when

§ Internal structure names the module and its source.

## Notes
