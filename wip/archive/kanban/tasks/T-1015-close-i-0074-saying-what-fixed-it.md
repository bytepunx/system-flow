---
id: T-1015
type: task
nature: improvement
title: Close I-0074 saying what fixed it
status: done
parent: S-0278
owner: alex
created: 2026-10-06T11:32:34Z
updated: 2026-10-06T20:03:39Z
transitions:
  - to: ready
    at: 2026-10-06T19:59:55Z
    by: agent-S-0278
  - to: in-progress
    at: 2026-10-06T19:59:55Z
    by: agent-S-0278
  - to: done
    at: 2026-10-06T20:03:39Z
    by: agent-S-0278
stream: S-0278
tags: [issues]
touches: [design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/summary.md]
after: [T-1012, T-1013]
usage:
  source: log
  seconds: 212
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 29
      cache_read: 193138
      cache_write: 3200
      cost: 0.0853
---
# T-1015 Close I-0074 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0074 --reason`. The reason names T-1011's ADR. It also says what changed: sync and acceptance regenerate `design/issues/summary.md` when a rebase stops on it alone, and the trial merge no longer reports it. The command regenerates `summary.md`.

This is the story's second criterion. It waits for T-1012 and T-1013, so the issue is closed only once the fix is in. It shares no path with T-1014, so the two run together.

## Done when

- I-0074 is `closed`, with a reason naming the fix and the ADR.
- `design/issues/summary.md` no longer lists I-0074.
- `flai check --strict` passes.

## Notes
