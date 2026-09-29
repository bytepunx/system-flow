---
id: T-0579
type: task
nature: improvement
title: Record cause 7's remeasurement and the changed contracts
status: done
parent: S-0162
owner: alex
created: 2026-09-29T20:15:55Z
updated: 2026-09-29T20:27:28Z
transitions:
  - to: ready
    at: 2026-09-29T20:16:04Z
    by: agent-S-0162
  - to: in-progress
    at: 2026-09-29T20:23:57Z
    by: agent-S-0162
  - to: done
    at: 2026-09-29T20:27:28Z
    by: agent-S-0162
stream: S-0162
tags: []
touches: [design/system, docs]
usage:
  source: log
  seconds: 211
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 12562
      cache_read: 3536589
      cache_write: 47096
      cost: 1.3355
---
# T-0579 Record cause 7's remeasurement and the changed contracts

## Work

- `design/system/server-performance.md` gets a section for cause 7 under After the stories, with what was measured.
- `design/system/flaiover-dashboard.md` and `design/system/flai-cli.md` say what `items()`, `docs.tree`, `items.count`, and `/_ready` ask now; `docs/` where the user-facing description changes.

## Done when

- `flai check --strict` is clean and the story's criteria are checked.

## Notes
