---
id: T-1401
type: task
nature: feature
title: The user and design documents say what a harness without the guard or asking loses and what lifts each refusal
status: done
parent: S-0355
owner: alex
created: 2026-10-08T08:49:25Z
updated: 2026-10-13T03:55:09Z
transitions:
  - to: ready
    at: 2026-10-13T03:51:58Z
    by: agent-S-0355
  - to: in-progress
    at: 2026-10-13T03:51:59Z
    by: agent-S-0355
  - to: done
    at: 2026-10-13T03:55:09Z
    by: agent-S-0355
stream: S-0355
tags: [cli]
touches: [docs/users/flai.md, design/system/flai-cli.md]
after: [T-1399, T-1400]
usage:
  source: log
  seconds: 190
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 54
      output: 19263
      cache_read: 2785727
      cache_write: 108042
      cost: 1.5663
---
# T-1401 The user and design documents say what a harness without the guard or asking loses and what lifts each refusal

## Work

- `docs/users/flai.md`, "Sub-agents" and "Starting an agent": what a harness without the guard loses (its roles), what a harness that cannot ask loses (its start, unless auto-approve or `deny_protected`), `guard: none` and what it risks, and that `command` is such a harness today.
- `design/system/flai-cli.md`: the capability checks before a start and their refusal texts; cite ADR-0131.

It waits for T-1399 and T-1400, so that it describes the refusals and the settings as built.

## Done when

- The two documents match the code.
- `flai test docs/users/flai.md design/system/flai-cli.md` passes.

## Notes

Layer 3 of S-0355.
