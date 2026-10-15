---
id: T-1564
type: task
nature: feature
title: The template ships the design-system convention as baseline and lists it in its conventions README
status: done
parent: S-0386
owner: alex
created: 2026-10-13T17:49:24Z
updated: 2026-10-15T04:40:35Z
transitions:
  - to: ready
    at: 2026-10-15T04:21:14Z
    by: agent-S-0386
  - to: in-progress
    at: 2026-10-15T04:21:15Z
    by: agent-S-0386
  - to: done
    at: 2026-10-15T04:40:35Z
    by: agent-S-0386
stream: S-0386
tags: [template]
touches: [template/root/design/conventions/design-system.md, template/root/design/conventions/README.md]
after: [T-1563]
usage:
  source: log
  seconds: 1160
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 4037
      cache_read: 577996
      cache_write: 22203
      cost: 0.3262
---
# T-1564 The template ships the design-system convention as baseline and lists it in its conventions README

## Work

Copy the baseline of `design/conventions/design-system.md` to `template/root/design/conventions/design-system.md`, with an empty Project additions below the marker, and list it in `template/root/design/conventions/README.md`. Check that a project rendered from the template passes `flai check --strict`. It waits for T-1563, whose convention it ships, as `CLAUDE.md` asks: `design/system` first, then `template/`.

## Done when

- The template's convention matches this repository's above the marker, its README lists it, and a rendered project passes `flai check --strict`.

## Notes

- Layer 2.
