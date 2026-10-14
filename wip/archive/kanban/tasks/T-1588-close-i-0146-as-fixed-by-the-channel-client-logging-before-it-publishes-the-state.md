---
id: T-1588
type: task
nature: remediation
title: Close I-0146 as fixed by the channel client logging before it publishes the state
status: done
parent: S-0389
owner: alex
created: 2026-10-14T23:11:05Z
updated: 2026-10-14T23:14:59Z
transitions:
  - to: ready
    at: 2026-10-14T23:14:50Z
    by: agent-S-0389
  - to: in-progress
    at: 2026-10-14T23:14:51Z
    by: agent-S-0389
  - to: done
    at: 2026-10-14T23:14:59Z
    by: agent-S-0389
stream: S-0389
tags: [issues]
touches: [design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md, design/issues/summary.md]
after: [T-1587]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 715
      cache_read: 275478
      cache_write: 5147
      cost: 0.1106
---
# T-1588 Close I-0146 as fixed by the channel client logging before it publishes the state

## Work

Close I-0146 with `flai issue close I-0146 --reason "<reason>"` from the story's worktree. The reason names the reorder in `Client.Run` and the test T-1587 added that makes the race certain.

Waits for T-1587: the reason names the fix and the test it lands.

## Done when

- I-0146 is closed with that reason, and `design/issues/summary.md` no longer lists it.

## Notes
