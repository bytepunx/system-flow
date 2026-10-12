---
id: T-1397
type: task
nature: improvement
title: The user and design documents describe flai's protected paths and the hold-and-ask
status: done
parent: S-0354
owner: alex
created: 2026-10-08T08:48:21Z
updated: 2026-10-12T21:49:04Z
transitions:
  - to: ready
    at: 2026-10-12T21:46:02Z
    by: agent-S-0354
  - to: in-progress
    at: 2026-10-12T21:46:05Z
    by: agent-S-0354
  - to: done
    at: 2026-10-12T21:49:04Z
    by: agent-S-0354
stream: S-0354
tags: [cli]
touches: [docs/users/flai.md, docs/operators/settings.md, design/system/flai-cli.md]
after: [T-1396]
usage:
  source: log
  seconds: 179
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 55
      output: 19868
      cache_read: 2505325
      cache_write: 121420
      cost: 1.6572
---
# T-1397 The user and design documents describe flai's protected paths and the hold-and-ask

## Work

- `docs/users/flai.md`: rename "Writes to paths Claude Code protects" to "Writes to protected paths", and describe the list as flai's, kept per harness, and the hold-and-ask with `permission_prompt` as Claude Code's form of it.
- `docs/operators/settings.md` links that anchor: point it at the new one.
- `design/system/flai-cli.md`: `flai/internal/ask`, the `ask` verdict, `guard.Hold`, and the list per harness; cite ADR-0130.

It waits for T-1396, so that it describes the verdict as built.

## Done when

- No link in `docs/` or `design/` points at the old anchor.
- `flai test docs/users/flai.md docs/operators/settings.md design/system/flai-cli.md` passes.

## Notes

Layer 3 of S-0354.
