---
id: T-0551
type: task
nature: feature
title: The story page opens on the thread its link names
status: done
parent: S-0155
owner: alex
created: 2026-09-29T07:02:04Z
updated: 2026-09-29T07:05:23Z
transitions:
  - to: ready
    at: 2026-09-29T07:02:22Z
    by: agent-S-0155
  - to: in-progress
    at: 2026-09-29T07:03:01Z
    by: agent-S-0155
  - to: done
    at: 2026-09-29T07:05:23Z
    by: agent-S-0155
stream: S-0155
tags: []
usage:
  source: log
  seconds: 142
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 26
      output: 204
      cache_read: 1556830
      cache_write: 25421
      cost: 0.6242
---

# T-0551 The story page opens on the thread its link names

## Work

`Threads.svelte` takes the thread to show first (`select`). The item page passes `?thread=` from its URL. When the thread is in the list the pager opens on it and the section scrolls into view; otherwise the page behaves as before.

## Done when

A component test shows the named thread, not the first, with the pager reading its position; the whole flaiover suite and lint pass.

## Notes
