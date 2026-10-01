---
id: T-0630
type: task
nature: feature
title: Diagnose and fix the stop test's darwin failure
status: done
parent: S-0180
owner: arobson
created: 2026-10-01T08:22:54Z
updated: 2026-10-01T08:23:04Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
  - to: in-progress
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
  - to: done
    at: 2026-10-01T08:23:04Z
    by: agent-S-0180
stream: S-0180
tags: []
usage:
  source: log
  seconds: 0
  models: []
---

# T-0630 Diagnose and fix the stop test's darwin failure

## Work

Find why `TestTheOperatorStopsAStorysAgent` fails on darwin and fix the cause; record it in the narrative.

## Done when

The stop test passes on macOS repeatedly, and the cause is in the narrative's Decisions.

## Notes
