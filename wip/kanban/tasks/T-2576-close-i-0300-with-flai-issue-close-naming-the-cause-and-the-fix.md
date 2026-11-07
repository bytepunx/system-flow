---
id: T-2576
type: task
nature: remediation
title: Close I-0300 with flai issue close, naming the cause and the fix
status: backlog
parent: S-0589
owner: alex
created: 2026-11-07T10:06:49Z
updated: 2026-11-07T10:06:49Z
transitions: []
stream: S-0589
tags: [issues]
touches: [design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md, design/issues/summary.md]
after: [T-2574]
---
# T-2576 Close I-0300 with flai issue close, naming the cause and the fix

## Work

Close I-0300 with `flai issue close I-0300 --reason "<reason>"`, which writes the issue file and regenerates `design/issues/summary.md`. It waits for T-2574 because the reason names the fix that task made and the test that now passes. The reason says the cause T-2573 recorded, the file and function T-2574 changed, and the test that reproduces I-0300.

## Done when

- I-0300's status is closed, with a reason naming the cause, the fix, and the test.
- `design/issues/summary.md` no longer lists I-0300 as open.

## Notes
