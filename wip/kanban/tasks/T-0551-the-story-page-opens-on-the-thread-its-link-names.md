---
id: T-0551
type: task
nature: feature
title: The story page opens on the thread its link names
status: in-progress
parent: S-0155
owner: alex
created: 2026-09-29T07:02:04Z
updated: 2026-09-29T07:03:01Z
transitions:
  - to: ready
    at: 2026-09-29T07:02:22Z
    by: agent-S-0155
  - to: in-progress
    at: 2026-09-29T07:03:01Z
    by: agent-S-0155
stream: S-0155
tags: []
---

# T-0551 The story page opens on the thread its link names

## Work

`Threads.svelte` takes the thread to show first (`select`). The item page passes `?thread=` from its URL. When the thread is in the list the pager opens on it and the section scrolls into view; otherwise the page behaves as before.

## Done when

A component test shows the named thread, not the first, with the pager reading its position; the whole flaiover suite and lint pass.

## Notes
