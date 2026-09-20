---
id: T-0251
type: task
nature: improvement
title: "flai check: an open story under a cancelled epic, and an open task under a cancelled story"
status: done
parent: S-0070
owner: alex
created: 2026-09-20T06:16:43Z
updated: 2026-09-20T06:30:42Z
transitions:
  - to: ready
    at: 2026-09-20T06:23:29Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T06:23:29Z
    by: system-flow
  - to: done
    at: 2026-09-20T06:30:42Z
    by: system-flow
stream: S-0070
tags: []
---
# T-0251 flai check: an open story under a cancelled epic, and an open task under a cancelled story

## Work
A rule in `flai/internal/check` that reports an open child under a cancelled parent, with a message that says how to fix it (cancel the parent again with this flai, or move the child). Listed in the design's table of rules.

## Done when
- Tests for both levels and for the clean case
- This repository passes `flai check --strict`

## Notes
