---
id: T-0650
type: task
nature: feature
title: The dashboard refuses a second flai connection while its holder answers, and flai backs off on that close code
status: done
parent: S-0184
owner: arobson
created: 2026-10-01T08:55:54Z
updated: 2026-10-01T09:06:58Z
transitions:
  - to: ready
    at: 2026-10-01T08:56:22Z
    by: agent-S-0184
  - to: in-progress
    at: 2026-10-01T09:02:51Z
    by: agent-S-0184
  - to: done
    at: 2026-10-01T09:06:58Z
    by: agent-S-0184
stream: S-0184
tags: []
usage:
  source: log
  seconds: 247
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 20288
      cache_read: 4576406
      cache_write: 81854
      cost: 1.8817
---

# T-0650 The dashboard refuses a second flai connection while its holder answers, and flai backs off on that close code

## Work

- `flaiover/src/lib/server/agent.ts`: when a flai connects for a project that already has one, ping the holder; if it answers in time, close the newcomer with a close code of its own that says the project is held; if it does not, replace it as today.
- `flai/internal/channel` (and its caller in flai serve): on that close code, wait a long back-off before dialling the project again, and log it once.

## Done when

- Tests in flaiover cover the refused and the replaced connection; Go tests cover the back-off on the close code.

## Notes
