---
id: T-0339
type: task
nature: feature
title: "The dashboard refuses a newcomer while a healthy connection answers; a refused flai backs off instead of redialling"
status: cancelled
parent: S-0084
owner: alex
created: 2026-09-23T01:10:04Z
updated: 2026-09-23T01:18:00Z
transitions:
  - to: ready
    at: 2026-09-23T01:15:06Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T01:15:09Z
    by: system-flow
  - to: cancelled
    at: 2026-09-23T01:18:00Z
    by: alex
stream: S-0084
tags: []
---

# T-0339 The dashboard refuses a newcomer while a healthy connection answers; a refused flai backs off instead of redialling

## Work

## Done when

## Notes
- 2026-09-23T01:18:00Z: moved to cancelled: S-0084 cancelled: this defect no longer aligns correctly with a more global approach to having flai/flaiover capable of managing multiple repositories per host
