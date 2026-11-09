---
id: T-2821
type: task
nature: improvement
title: Close I-0320 with what fixed it
status: backlog
parent: S-0643
owner: agent
created: 2026-11-09T22:13:30Z
updated: 2026-11-09T22:13:30Z
transitions: []
stream: S-0643
tags: [issues]
touches: [design/issues/I-0320-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md, design/issues/summary.md]
after: [T-2819, T-2820]
---
# T-2821 Close I-0320 with what fixed it

## Work

Close I-0320 with `flai issue close I-0320 --reason`. The reason names the `--require-pending` flag, the publish step's `stop: [2]`, and the test that reproduces the race. The command regenerates `design/issues/summary.md`.

This task waits for T-2819 and T-2820, because the reason names what they deliver.

## Done when

- I-0320 is closed, with a reason naming the flag, the formula's stop, and the test.
- `design/issues/summary.md` no longer lists I-0320.

## Notes
