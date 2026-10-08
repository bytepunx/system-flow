---
id: T-1383
type: task
nature: remediation
title: Close I-0119 with the reason that names the sync step's wip rule and the close-out's flai verify --sync-only
status: backlog
parent: S-0347
owner: alex
created: 2026-10-08T08:44:35Z
updated: 2026-10-08T08:44:35Z
transitions: []
stream: S-0347
tags: [docs]
touches: [design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/summary.md]
after: [T-1382]
---
# T-1383 Close I-0119 with the reason that names the sync step's wip rule and the close-out's flai verify --sync-only

## Work

Close I-0119 in the story's worktree with `flai issue close I-0119 --reason "<reason>"`. The reason names the ADR, the `sync` step of `flai verify` that passes over commits on main changing only `wip/` paths the branch does not change, and the close-out's last check, now `flai verify --sync-only`. Say that a commit outside `wip/`, such as an acceptance, still needs a sync and a new run, as the issue's third instance did.

Waits for T-1382: the issue closes when the close-out it describes no longer stops on a `wip`-only commit.

## Done when

- I-0119's status is closed with that reason, and `design/issues/summary.md` no longer lists it.
- Criterion 2 of S-0347 is ticked with `flai criteria tick`.

## Notes

Drafted by the planner for S-0347.
