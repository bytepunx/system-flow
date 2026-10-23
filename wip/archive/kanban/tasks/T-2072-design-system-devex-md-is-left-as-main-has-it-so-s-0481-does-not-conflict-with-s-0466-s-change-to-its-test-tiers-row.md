---
id: T-2072
type: task
nature: feature
title: design/system/devex.md is left as main has it, so S-0481 does not conflict with S-0466's change to its test tiers row
status: done
parent: S-0481
owner: alex
created: 2026-10-23T07:19:38Z
updated: 2026-10-23T07:19:45Z
transitions:
  - to: ready
    at: 2026-10-23T07:19:38Z
    by: agent-S-0481
  - to: in-progress
    at: 2026-10-23T07:19:39Z
    by: agent-S-0481
  - to: done
    at: 2026-10-23T07:19:45Z
    by: agent-S-0481
stream: S-0481
tags: []
touches: [design/system/devex.md]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 687
      cache_read: 152988
      cache_write: 4794
      cost: 0.0827
---

# T-2072 design/system/devex.md is left as main has it, so S-0481 does not conflict with S-0466's change to its test tiers row

## Work

Waits for no task. A trial merge at T-2070's sync found its sentence in `design/system/devex.md`'s test tiers row conflicts with S-0466, in review, which changes the same one-line row (MS-0178). `scripts/README.md` and `design/tech/node.md` already say what the fix does, so the sentence goes and the row stays as main has it.

## Done when

- The story's branch does not change `design/system/devex.md`, and the next sync finds S-0466 and S-0481 merging cleanly.

## Notes
