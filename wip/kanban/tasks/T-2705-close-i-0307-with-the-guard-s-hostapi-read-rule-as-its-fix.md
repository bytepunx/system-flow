---
id: T-2705
type: task
nature: remediation
title: Close I-0307 with the guard's hostapi read rule as its fix
status: done
parent: S-0612
owner: alex
created: 2026-11-09T07:58:24Z
updated: 2026-11-09T09:06:47Z
transitions:
  - to: ready
    at: 2026-11-09T09:06:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:06:37Z
    by: system-flow
  - to: done
    at: 2026-11-09T09:06:47Z
    by: agent-S-0612
stream: S-0612
tags: [issues]
touches: [design/issues/I-0307-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md, design/issues/summary.md]
after: [T-2703, T-2704]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 32
      cache_read: 584370
      cache_write: 1836
      cost: 0.2131
---
# T-2705 Close I-0307 with the guard's hostapi read rule as its fix

## Work

Waits for T-2703 and T-2704: the issue is closed once the fix, its test, and its documentation are in.

Run `flai issue close I-0307 --reason "..."` in the story's worktree, the reason naming the fix: `flai guard` lets every agent run `flai hostapi` alone and with the read methods on its list, `--timing` included, refuses the write methods, and a test in `flai/internal/hostapi` keeps the list in step with the method table (S-0612). It rewrites `design/issues/summary.md`.

## Done when

- I-0307's status is closed with that reason, `design/issues/summary.md` no longer lists it, and `flai check --strict` is clean.

## Notes
