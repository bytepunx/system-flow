---
id: T-0612
type: task
nature: feature
title: The design and the user guide say where each inbox entry leads
status: done
parent: S-0173
owner: alex
created: 2026-09-30T01:26:59Z
updated: 2026-10-01T07:46:34Z
transitions:
  - to: ready
    at: 2026-09-30T01:27:10Z
    by: agent-S-0173
  - to: in-progress
    at: 2026-10-01T07:40:58Z
    by: agent-S-0173
  - to: done
    at: 2026-10-01T07:46:34Z
    by: agent-S-0173
stream: S-0173
tags: []
touches: [design/system, docs]
usage:
  source: log
  seconds: 336
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 104
      output: 27295
      cache_read: 6567381
      cache_write: 94058
      cost: 2.6122
---

# T-0612 The design and the user guide say where each inbox entry leads

## Work

design/system/flaiover-dashboard.md and the user guide's inbox section say where each entry kind leads, and that a story's page shows and answers its open questions.

## Done when

The design and docs match the behaviour; flai check --strict passes.

## Notes
