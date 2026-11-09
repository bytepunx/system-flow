---
id: T-2817
type: task
nature: remediation
title: Close I-0321 with what fixed it
status: backlog
parent: S-0645
owner: agent
created: 2026-11-09T22:10:39Z
updated: 2026-11-09T22:10:39Z
transitions: []
stream: S-0645
tags: [flai, testing]
touches: [design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md, design/issues/summary.md]
after: [T-2816]
---
# T-2817 Close I-0321 with what fixed it

## Work

Close I-0321 with `flai issue close I-0321 --reason "<reason>"`. The reason names the cause T-2816 confirmed, or the likely one if the failure did not come back: under WSL2 mirrored networking, a loopback port that was free in Linux is sometimes answered by a relay. It also names the fix: the closed-port cases of `TestTCPAnswersAListenerAndNotAClosedOne` dial through a dialer the test controls, and only a listener the test holds open is dialled for real. flai regenerates `design/issues/summary.md`.

This task waits for T-2816 because the reason quotes what T-2816 found and changed.

## Done when

- I-0321's status is closed, and its reason names the cause and the fix.
- `design/issues/summary.md` no longer lists I-0321 as open.
- `flai check --strict` passes on the story.

## Notes
