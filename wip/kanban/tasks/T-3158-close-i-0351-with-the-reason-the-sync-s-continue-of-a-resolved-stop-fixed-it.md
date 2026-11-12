---
id: T-3158
type: task
nature: remediation
title: Close I-0351 with the reason the sync's continue of a resolved stop fixed it
status: backlog
parent: S-0726
owner: alex
created: 2026-11-12T15:57:01Z
updated: 2026-11-12T15:57:01Z
transitions: []
stream: S-0726
tags: [flai]
touches: [design/issues/I-0351-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md, design/issues/summary.md]
after: [T-3156, T-3157]
---
# T-3158 Close I-0351 with the reason the sync's continue of a resolved stop fixed it

## Work

Run `flai issue close I-0351 --reason "<reason>"` in the story's worktree. The reason names the ADR T-3154 wrote. It says that `flai stream sync`, run again after the agent resolves and stages a stop, now continues the rebase itself and settles later stops on issue files and `design/issues/summary.md`. It also names the test T-3155 added. The command regenerates `design/issues/summary.md`.

Waits for T-3156 and T-3157, so that the fix, its advice, and its docs are all in before the issue is closed. This task is the last layer.

## Done when

- I-0351 is closed with that reason.
- `design/issues/summary.md` no longer lists it.
- Criterion 2 can be ticked.

## Notes
