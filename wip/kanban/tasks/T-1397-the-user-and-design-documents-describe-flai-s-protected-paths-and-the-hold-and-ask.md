---
id: T-1397
type: task
nature: improvement
title: The user and design documents describe flai's protected paths and the hold-and-ask
status: backlog
parent: S-0354
owner: alex
created: 2026-10-08T08:48:21Z
updated: 2026-10-08T08:48:30Z
transitions: []
stream: S-0354
tags: [cli]
touches: [docs/users/flai.md, docs/operators/settings.md, design/system/flai-cli.md]
after: [T-1396]
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
