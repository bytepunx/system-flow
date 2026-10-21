---
id: T-1913
type: task
nature: feature
title: The CLI design and the user docs say a new item keeps a wip.overlap as a warning, and I-0169 is closed
status: done
parent: S-0451
owner: alex
created: 2026-10-20T19:04:11Z
updated: 2026-10-20T19:06:34Z
transitions:
  - to: ready
    at: 2026-10-20T19:06:06Z
    by: agent-S-0451
  - to: in-progress
    at: 2026-10-20T19:06:06Z
    by: agent-S-0451
  - to: done
    at: 2026-10-20T19:06:34Z
    by: agent-S-0451
stream: S-0451
tags: []
touches: [design/system/flai-cli.md, docs/users/flai.md, design/issues]
after: [T-1912]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2539
      cache_read: 767711
      cache_write: 30409
      cost: 0.4392
---
# T-1913 The CLI design and the user docs say a new item keeps a wip.overlap as a warning, and I-0169 is closed

## Work

Say in `design/system/flai-cli.md` (the item creation row), `docs/users/flai.md` (creating items), and `design/system/agent-coordination.md` where it fits, that a creation keeps a `wip.overlap` it introduces as a warning rather than refusing it. Close I-0169 with `flai issue close I-0169 --reason` saying what fixed it. Waits for the code task, whose behaviour it describes.

## Done when

- The docs describe the behaviour the code task built.
- I-0169 is closed and `design/issues/summary.md` is current.

## Notes
