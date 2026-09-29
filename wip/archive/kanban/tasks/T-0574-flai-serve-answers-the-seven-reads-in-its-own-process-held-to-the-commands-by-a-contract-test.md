---
id: T-0574
type: task
nature: feature
title: flai serve answers the seven reads in its own process, held to the commands by a contract test
status: done
parent: S-0159
owner: alex
created: 2026-09-29T19:58:35Z
updated: 2026-09-29T20:09:00Z
transitions:
  - to: ready
    at: 2026-09-29T19:58:49Z
    by: agent-S-0159
  - to: in-progress
    at: 2026-09-29T20:03:39Z
    by: agent-S-0159
  - to: done
    at: 2026-09-29T20:09:00Z
    by: agent-S-0159
stream: S-0159
tags: []
touches: [flai/internal/hostapi]
usage:
  source: log
  seconds: 321
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 25850
      cache_read: 5605702
      cache_write: 81872
      cost: 2.2933
---

# T-0574 flai serve answers the seven reads in its own process, held to the commands by a contract test

## Work

Answer `publish.preview`, `push.pending`, `stats.get`, `stream.diff`, `item.show`, `item.move.preview`, and `accept.preview` in `hostapi` from those functions, with the same answer shape (`data`, `warnings`) and the same error codes as when flai ran, and drop their specs from the write table. Name their phases for the work, not `exec.flai`.

## Done when

- [x] None of the seven has an `exec.flai` phase in its `request answered` event.
- [x] A contract test runs each command with `--json` and the method in-process on a fixture repository with git, and compares the answers and errors.

## Notes
