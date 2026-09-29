---
id: T-0584
type: task
nature: feature
title: Every tier passes and the charts are looked at with this repository's numbers
status: done
parent: S-0163
owner: alex
created: 2026-09-29T20:33:48Z
updated: 2026-09-29T21:00:44Z
transitions:
  - to: ready
    at: 2026-09-29T20:34:13Z
    by: agent-S-0163
  - to: in-progress
    at: 2026-09-29T20:57:45Z
    by: agent-S-0163
  - to: done
    at: 2026-09-29T21:00:44Z
    by: agent-S-0163
stream: S-0163
tags: []
usage:
  source: log
  seconds: 179
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 26
      output: 17535
      cache_read: 2972920
      cache_write: 50411
      cost: 2.6285
---
# T-0584 Every tier passes and the charts are looked at with this repository's numbers

## Work

Run the three tiers and every lint. Render each new chart from this repository's `flai stats --json` and look at it, in light and dark, at each bucket. Tick the criteria that were verified and say in the story's notes how each was.

## Done when

- `make flai-test`, `make flaiover-test`, the markdown lint, and `flai check --strict` pass.
- Each acceptance criterion is ticked with what verified it, or left unticked with why.

## Notes
