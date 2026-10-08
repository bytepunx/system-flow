---
id: T-1429
type: task
nature: improvement
title: I-0092 is closed with flai issue close, naming S-0326's merge and the tests that reproduce its instances
status: done
parent: S-0297
owner: alex
created: 2026-10-08T09:44:11Z
updated: 2026-10-08T09:50:02Z
transitions:
  - to: ready
    at: 2026-10-08T09:49:42Z
    by: agent-S-0297
  - to: in-progress
    at: 2026-10-08T09:49:42Z
    by: agent-S-0297
  - to: done
    at: 2026-10-08T09:50:02Z
    by: agent-S-0297
stream: S-0297
tags: []
touches: [design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md, design/issues/summary.md]
after: [T-1428]
usage:
  source: log
  seconds: 20
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 723
      cache_read: 160096
      cache_write: 6972
      cost: 0.0966
---

# T-1429 I-0092 is closed with flai issue close, naming S-0326's merge and the tests that reproduce its instances

## Work

Close I-0092 with `flai issue close I-0092 --reason ... --commit` in the story's worktree, naming S-0326 (ADR-0126) and the tests in `flai/cmd/stream_sync_test.go` that reproduce each instance. Waits for T-1428, whose tests the reason names.

## Done when

- I-0092 is closed on `story/S-0297` and `design/issues/summary.md` no longer lists it as open.

## Notes
