---
id: T-3562
type: task
nature: remediation
title: Close I-0355 with what fixed it
status: done
parent: S-0829
owner: alex
created: 2026-11-17T12:29:00Z
updated: 2026-11-17T13:12:16Z
transitions:
  - to: ready
    at: 2026-11-17T13:12:10Z
    by: agent-S-0829
  - to: in-progress
    at: 2026-11-17T13:12:10Z
    by: agent-S-0829
  - to: done
    at: 2026-11-17T13:12:16Z
    by: agent-S-0829
stream: S-0829
tags: [issues]
touches: [design/issues/I-0355-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md, design/issues/summary.md]
after: [T-3559]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 411
      cache_read: 248522
      cache_write: 9142
      cost: 0.1311
---
# T-3562 Close I-0355 with what fixed it

## Work

Run `flai issue close I-0355 --reason "<reason>"` in the story's worktree, the reason saying that `cliReads` in `flai/internal/guard/guard_tables.go` now lists `kluster status` and `deployments`, so every sub-agent and strategic role runs them, and naming the test T-3559 added. It writes the issue and `design/issues/summary.md`.

It waits for T-3559 because the reason names that task's fix and test, which must exist and pass first.

## Done when

- I-0355's status is closed, its reason names the fix and the test.
- `design/issues/summary.md` no longer lists I-0355 as open, and `flai check --strict` is clean.

## Notes
