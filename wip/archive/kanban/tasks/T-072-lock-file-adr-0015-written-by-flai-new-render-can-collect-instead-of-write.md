---
id: T-072
type: task
nature: feature
title: "Lock file (ADR-0015) written by flai new; render can collect instead of write"
status: done
parent: S-020
owner: alex
created: 2026-09-17T03:19:34Z
updated: 2026-09-17T03:25:21Z
transitions:
  - to: ready
    at: 2026-09-17T03:25:20Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:25:20Z
    by: agent
  - to: done
    at: 2026-09-17T03:25:21Z
    by: agent
stream: S-020
tags: [cli, template]
---

# T-072 Lock file (ADR-0015) written by flai new; render can collect instead of write

## Work
ADR-0015: system-flow.lock.yaml at the root records the template source and a sha256 per rendered path; internal/lock loads, saves, and computes hashes; template.Render gains a Collect option that returns rendered contents without writing; flai new writes the lock after rendering.

## Done when
New projects carry the lock; tests cover round-trip and hashing.

## Notes
