---
id: T-2903
type: task
nature: remediation
title: Close I-0318 with what fixed it
status: backlog
parent: S-0680
owner: alex
created: 2026-11-10T08:46:54Z
updated: 2026-11-10T08:46:54Z
transitions: []
stream: S-0680
tags: [issues]
touches: [design/issues/I-0318-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md, design/issues/summary.md]
after: [T-2901]
---
# T-2903 Close I-0318 with what fixed it

## Work

Run `flai issue close I-0318 --reason "<reason>"` in the story's worktree, the reason naming `readsServe` in `flai/internal/guard/guard_tables.go`, the read forms it passes, and the tests that pin it, so that `design/issues/summary.md` drops the issue.

Waits for T-2901: the reason names what that task changed.

## Done when

- I-0318 is closed with that reason, `design/issues/summary.md` no longer lists it, and `flai check --strict` passes.

## Notes
