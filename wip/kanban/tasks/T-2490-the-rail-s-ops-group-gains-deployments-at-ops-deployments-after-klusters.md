---
id: T-2490
type: task
nature: feature
title: The rail's Ops group gains Deployments at /ops/deployments, after Klusters
status: backlog
parent: S-0581
owner: alex
created: 2026-11-06T21:46:00Z
updated: 2026-11-06T21:46:00Z
transitions: []
stream: S-0581
tags: [kluster, flaiover]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts]
---
# T-2490 The rail's Ops group gains Deployments at /ops/deployments, after Klusters

## Work

In `flaiover/src/lib/sitemenu.ts`, add Deployments, at `/ops/deployments`, to the Ops group after Klusters, with its own lucide icon. In `sitemenu.test.ts`, pin it, and test that `locate` makes it current for `/ops/deployments`.

This is the first layer: it waits for no task, and shares no path with the host API task.

## Done when

- The rail lists Deployments under Ops after Klusters.
- The order and `locate` are tested, and `flai test` on the two paths passes.

## Notes
