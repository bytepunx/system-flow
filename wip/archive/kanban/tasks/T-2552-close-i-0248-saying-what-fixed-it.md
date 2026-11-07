---
id: T-2552
type: task
nature: remediation
title: Close I-0248 saying what fixed it
status: done
parent: S-0561
owner: alex
created: 2026-11-07T08:19:03Z
updated: 2026-11-07T08:29:14Z
transitions:
  - to: ready
    at: 2026-11-07T08:29:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:29:03Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:29:14Z
    by: agent-S-0561
stream: S-0561
tags: [cli]
touches: [design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md, design/issues/summary.md]
after: [T-2550]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 337
      cache_read: 495025
      cache_write: 6653
      cost: 0.159
---
# T-2552 Close I-0248 saying what fixed it

## Work

Run `flai issue close I-0248 --commit --reason "<reason>"` in the story's worktree, the reason naming the wait T-2550 added and why it removes the failure: the case now waits for the ended event, which `endRun` appends just after it records the outcome, before it stops the agent.

It waits for T-2550 because the reason says what that task fixed, and only once its runs have shown it.

## Done when

- I-0248's file has status closed, and its body ends with the reason.
- `design/issues/summary.md` no longer lists I-0248.

## Notes
