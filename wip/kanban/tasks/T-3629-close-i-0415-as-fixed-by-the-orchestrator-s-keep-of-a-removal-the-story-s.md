---
id: T-3629
type: task
nature: remediation
title: Close I-0415 as fixed by the orchestrator's keep of a removal the story's narrative records as a rewrite
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:58:24Z
updated: 2026-11-27T23:58:24Z
transitions: []
stream: S-0852
tags: [flai]
touches: [design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md, design/issues/summary.md]
after: [T-3624, T-3626]
---
# T-3629 Close I-0415 as fixed by the orchestrator's keep of a removal the story's narrative records as a rewrite

## Work

Run `flai issue close I-0415 --reason "<what fixed it>"`. Name the new ADR of T-3618, the permission, and the reproduction test of T-3620, and say that the orchestrator's `accept` and `queue_add` now take `keep_removals`. It rewrites `design/issues/summary.md`. Tick the story's two criteria with `flai criteria tick` once each is verified.

It waits for T-3624 and T-3626: the issue is closed only once the role and the docs say what fixed it.

## Done when

- I-0415's status is closed, with that reason, and `design/issues/summary.md` lists it so.
- `flai check --strict` passes.

## Notes
