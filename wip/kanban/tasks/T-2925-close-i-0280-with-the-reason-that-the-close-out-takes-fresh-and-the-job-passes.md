---
id: T-2925
type: task
nature: improvement
title: Close I-0280 with the reason that the close-out takes --fresh and the job passes it on
status: backlog
parent: S-0671
owner: alex
created: 2026-11-10T08:48:36Z
updated: 2026-11-10T08:48:36Z
transitions: []
stream: S-0671
tags: [issues, close-out]
touches: [design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md, design/issues/summary.md]
after: [T-2918, T-2920, T-2921, T-2923]
---
# T-2925 Close I-0280 with the reason that the close-out takes --fresh and the job passes it on

## Work

- In the story's worktree, run `flai issue close I-0280 --reason "<reason>"`. The reason names what fixed it: `scripts/close-out.sh` and the template's copy take `--fresh` and `--all` before the story, or `CLOSE_OUT_FRESH` and `CLOSE_OUT_ALL`, and flai serve's close-out job sets those variables from the request, with the test T-2920 added.
- Waits for every other task of the story, since the reason states what they built and the issue closes only once it is fixed.

## Done when

- I-0280's status is closed, its reason names the fix, and `design/issues/summary.md` no longer lists it.
- Criteria 1 and 2 of S-0671 are ticked with `flai criteria tick`.

## Notes
