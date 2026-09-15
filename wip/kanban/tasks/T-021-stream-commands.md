---
id: T-021
type: task
nature: feature
title: stream open and stream log, index maintenance
status: done
parent: S-007
owner: agent
created: 2026-09-15T17:42:41Z
updated: 2026-09-15T17:52:36Z
transitions:
  - to: ready
    at: 2026-09-15T17:42:41Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:52:36Z
    by: agent
  - to: done
    at: 2026-09-15T17:52:36Z
    by: agent
stream: S-007
tags: [cli, workitems]
---

# T-021 stream open and stream log, index maintenance

## Work
stream open renders the narrative template for a story; stream log appends a timestamped entry and bumps updated; index.md is regenerated from active narratives on every open, log, move, and archive.

## Done when
Tests cover open, log, and index regeneration.

## Notes
