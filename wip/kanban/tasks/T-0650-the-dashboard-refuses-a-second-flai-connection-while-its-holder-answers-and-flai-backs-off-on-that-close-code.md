---
id: T-0650
type: task
nature: feature
title: The dashboard refuses a second flai connection while its holder answers, and flai backs off on that close code
status: ready
parent: S-0184
owner: arobson
created: 2026-10-01T08:55:54Z
updated: 2026-10-01T08:56:22Z
transitions:
  - to: ready
    at: 2026-10-01T08:56:22Z
    by: agent-S-0184
stream: S-0184
tags: []
---

# T-0650 The dashboard refuses a second flai connection while its holder answers, and flai backs off on that close code

## Work

- `flaiover/src/lib/server/agent.ts`: when a flai connects for a project that already has one, ping the holder; if it answers in time, close the newcomer with a close code of its own that says the project is held; if it does not, replace it as today.
- `flai/internal/channel` (and its caller in flai serve): on that close code, wait a long back-off before dialling the project again, and log it once.

## Done when

- Tests in flaiover cover the refused and the replaced connection; Go tests cover the back-off on the close code.

## Notes
