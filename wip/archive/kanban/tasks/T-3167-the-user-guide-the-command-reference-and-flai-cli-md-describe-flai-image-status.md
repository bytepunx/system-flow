---
id: T-3167
type: task
nature: feature
title: The user guide, the command reference, and flai-cli.md describe flai image status and the one place an agent image is decided
status: done
parent: S-0735
owner: alex
created: 2026-11-12T15:58:56Z
updated: 2026-11-13T08:34:41Z
transitions:
  - to: ready
    at: 2026-11-13T08:30:59Z
    by: agent-S-0735
  - to: in-progress
    at: 2026-11-13T08:30:59Z
    by: agent-S-0735
  - to: done
    at: 2026-11-13T08:34:41Z
    by: agent-S-0735
stream: S-0735
tags: [cli]
touches: [docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [T-3165, T-3166]
usage:
  source: log
  seconds: 222
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 7959
      cache_read: 2813758
      cache_write: 37988
      cost: 0.9703
---
# T-3167 The user guide, the command reference, and flai-cli.md describe flai image status and the one place an agent image is decided

## Work

- `docs/users/flai.md` § In a container (line 498 on): the command, its table, its exits, and the error table's rows (lines 539 to 540) naming it.
- `docs/users/flai-reference.md`: run `make flai-reference` to write it again from the command help. Never edit it by hand.
- `design/system/flai-cli.md` § The container harness (S-0508), § The project's layer image: `containerImageHost` now asks `flai/internal/agentimage`, and `flai image status` reports the same. Bump `updated`.

It waits for T-3165 and T-3166, whose command and wiring it describes.

## Done when

- The three files describe the command and the package as built.
- `TestReferenceIsCurrent` passes, and `flai check --strict` and the markdown lint pass on them.

## Notes
