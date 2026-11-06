---
id: T-2486
type: task
nature: feature
title: The rail's site menu gains an Ops group, between Design and Documentation, whose first page is Klusters at /ops/klusters
status: backlog
parent: S-0580
owner: alex
created: 2026-11-06T21:45:06Z
updated: 2026-11-06T21:45:24Z
transitions: []
stream: S-0580
tags: [kluster, flaiover]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/navrail.svelte.ts]
---
# T-2486 The rail's site menu gains an Ops group, between Design and Documentation, whose first page is Klusters at /ops/klusters

## Work

In `flaiover/src/lib/sitemenu.ts`, add an Ops group to `SITE_MENU`, between Design and Documentation, with a lucide icon imported as the others are (such as `server` or `container`). Its one page is Klusters, at `/ops/klusters`, with its own icon. `locate` then marks `/ops` and anything under it as Ops, with Klusters current for `/ops/klusters`. Add the group's key, `ops`, to the `keys` list in `flaiover/src/lib/navrail.svelte.ts`, so the rail remembers whether it is expanded.

In `sitemenu.test.ts`, pin the new group order and the Klusters page. Test that `locate` places `/ops/klusters` and an unknown `/ops/...` page in Ops.

This is the first layer: it waits for no task, and shares no path with the host API task.

## Done when

- The rail lists Ops between Design and Documentation, with Klusters under it, and remembers its expanded state.
- `locate` and the order are tested, and `flai test` on the three paths passes.

## Notes
