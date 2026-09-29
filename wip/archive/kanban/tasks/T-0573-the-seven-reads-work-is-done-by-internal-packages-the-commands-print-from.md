---
id: T-0573
type: task
nature: feature
title: The seven reads' work is done by internal packages the commands print from
status: done
parent: S-0159
owner: alex
created: 2026-09-29T19:58:35Z
updated: 2026-09-29T20:03:39Z
transitions:
  - to: ready
    at: 2026-09-29T19:58:49Z
    by: agent-S-0159
  - to: in-progress
    at: 2026-09-29T19:58:49Z
    by: agent-S-0159
  - to: done
    at: 2026-09-29T20:03:39Z
    by: agent-S-0159
stream: S-0159
tags: []
touches: [flai/cmd, flai/internal]
usage:
  source: log
  seconds: 290
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 64
      output: 28828
      cache_read: 6251651
      cache_write: 91306
      cost: 2.5576
---

# T-0573 The seven reads' work is done by internal packages the commands print from

## Work

Move what the seven reads compute out of `flai/cmd` so that flai serve can call it: the story-branch git helpers and `flai stream diff`'s diff into a package of their own, the dry runs of a cancellation, an acceptance, a push, and a publish into one package the commands' real runs share their preflight with, and `--since` parsing into `metrics`. `flai edit --show` already prints from `itemedit.Show`. No command's output changes.

## Done when

- [x] `flai stream diff`, `flai move <id> cancelled --dry-run`, `flai accept --dry-run`, `flai push --pending --dry-run`, `flai release --pending --dry-run`, and `flai stats` print from functions outside `cmd`.
- [x] `make test` and lint pass with the existing command tests unchanged in what they assert.

## Notes
