---
id: T-1560
type: task
nature: feature
title: The manifest, settings, flai design, and user docs and the reference describe design.stylesheet, flai design css --write, and design.stale
status: done
parent: S-0385
owner: alex
created: 2026-10-13T17:48:59Z
updated: 2026-10-15T02:50:30Z
transitions:
  - to: ready
    at: 2026-10-15T02:23:03Z
    by: agent-S-0385
  - to: in-progress
    at: 2026-10-15T02:48:42Z
    by: agent-S-0385
  - to: done
    at: 2026-10-15T02:50:30Z
    by: agent-S-0385
stream: S-0385
tags: [cli]
touches: [design/system/flai-cli.md, docs/users/flai.md, design/system/design-system-tool.md]
after: [T-1558, T-1559]
usage:
  source: log
  seconds: 108
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 10525
      cache_read: 1624126
      cache_write: 61024
      cost: 0.9176
---
# T-1560 The manifest, settings, flai design, and user docs and the reference describe design.stylesheet, flai design css --write, and design.stale

## Work

Describe `flai design css --write`, the stylesheet `flai design init` writes, and `design.stale` in `design/system/flai-cli.md` and `docs/users/flai.md`, and say in `design/system/design-system-tool.md` that S-0385 settled how the stylesheet selects the dark values (ADR-0159). It waits for T-1558 and T-1559, whose behaviour it describes.

## Done when

- The three documents describe it, and the markdown lint passes on them.

## Notes

- Layer 3. The settings row and the reference moved to T-1556 and T-1558, whose tests check them.
