---
id: T-0624
type: task
nature: research
title: Measure delegation on two stories against runs without it
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T07:53:02Z
updated: 2026-10-01T08:19:34Z
transitions:
  - to: ready
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T08:12:57Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T08:19:34Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [design/system/agent-context.md]
usage:
  source: log
  seconds: 397
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 61
      output: 23820
      cache_read: 5680623
      cache_write: 72151
      cost: 2.1097
---
# T-0624 Measure delegation on two stories against runs without it

## Work

- From flai serve's logs, with ADR-0051's reading (sub-agent usage included): cost, cache reads, and turns of at least two stories whose agents delegated, against comparable runs that did not.
- Record the table in `design/system/agent-context.md`.

## Done when

- The table is in `agent-context.md`, or the criterion is left unchecked with why in the story's notes.

## Notes
