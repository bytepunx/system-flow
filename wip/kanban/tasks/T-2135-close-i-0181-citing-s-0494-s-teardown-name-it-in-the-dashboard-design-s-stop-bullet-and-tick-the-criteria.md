---
id: T-2135
type: task
nature: remediation
title: Close I-0181 citing S-0494's teardown, name it in the dashboard design's stop bullet, and tick the criteria
status: backlog
parent: S-0499
owner: alex
created: 2026-10-26T04:21:10Z
updated: 2026-10-26T04:21:10Z
transitions: []
stream: S-0499
tags: [dashboard, issues]
touches: [design/issues/I-0181-events-socket-test-ts-closes-every-events-socket-with-1001-when-the-server-stops-fails-now-and-then-a-change-listener-is-still-attached-when-the-server-has-closed.md, design/issues/summary.md, design/system/flaiover-dashboard.md]
after: [T-2134]
---
# T-2135 Close I-0181 citing S-0494's teardown, name it in the dashboard design's stop bullet, and tick the criteria

## Work

Waits for T-2134, which shows which fix removed the race and whether `events-socket.ts` needed another change: the close reason says so.

- In `design/system/flaiover-dashboard.md`, the SIGTERM bullet under the agent endpoint section says each live-events socket is taken off the project's `Repo` before its 1001 is sent "(I-0205, S-0494)": name I-0181 there too, as the first instance of the same defect, and bump `updated`.
- Close the issue from the story's worktree: `flai issue close I-0181 --reason "..."`, the reason saying that I-0181's instance was the race I-0205 recorded and S-0494 fixed in `d207daaf3` (`EventsSocket.close()` tears each socket down before it sends 1001), that the server-stop test in `flaiover/src/lib/server/events-socket.test.ts` asserts no listener right after `close()` and fails against the pre-fix module, and how many repeated runs T-2134 passed; name any further change T-2134 made. The command rewrites the issue file and `design/issues/summary.md`.
- Tick the story's two criteria with `flai criteria tick S-0499 <n>` once each is verified.

## Done when

- I-0181 reads `status: closed` with a reason that names S-0494 and T-2134's evidence.
- `design/issues/summary.md` no longer lists I-0181 as open.
- `design/system/flaiover-dashboard.md` names I-0181 beside I-0205 and S-0494 in the SIGTERM bullet.
- Both acceptance criteria of S-0499 are ticked.
- `flai check --strict` is clean for the story.

## Notes

The second layer: after T-2134.
