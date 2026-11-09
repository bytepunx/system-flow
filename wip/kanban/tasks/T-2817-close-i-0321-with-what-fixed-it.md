---
id: T-2817
type: task
nature: remediation
title: Close I-0321 with what fixed it
status: done
parent: S-0645
owner: agent
created: 2026-11-09T22:10:39Z
updated: 2026-11-09T22:54:56Z
transitions:
  - to: ready
    at: 2026-11-09T22:54:41Z
    by: agent-S-0645
  - to: in-progress
    at: 2026-11-09T22:54:41Z
    by: agent-S-0645
  - to: done
    at: 2026-11-09T22:54:56Z
    by: agent-S-0645
stream: S-0645
tags: [flai, testing]
touches: [design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md, design/issues/summary.md]
after: [T-2816]
usage:
  source: log
  seconds: 15
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 525
      cache_read: 528527
      cache_write: 8533
      cost: 0.1845
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
