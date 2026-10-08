---
id: T-1293
type: task
nature: remediation
title: "A test reproduces I-0089: an issue recorded and one closed on main itself while a story branch records its own, and the sync and the acceptance regenerate design/issues/summary.md"
status: done
parent: S-0315
owner: alex
created: 2026-10-07T23:43:55Z
updated: 2026-10-07T23:57:18Z
transitions:
  - to: ready
    at: 2026-10-07T23:54:13Z
    by: agent-S-0315
  - to: in-progress
    at: 2026-10-07T23:54:14Z
    by: agent-S-0315
  - to: done
    at: 2026-10-07T23:57:18Z
    by: agent-S-0315
stream: S-0315
tags: [cli]
touches: [flai/cmd/stream_sync_test.go]
usage:
  source: log
  seconds: 184
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 3172
      cache_read: 784664
      cache_write: 42849
      cost: 0.5632
---
# T-1293 A test reproduces I-0089: an issue recorded and one closed on main itself while a story branch records its own, and the sync and the acceptance regenerate design/issues/summary.md

## Work

I-0089 is the conflict in `design/issues/summary.md` between a story branch that records an issue and an issue recorded on main meanwhile, not through another story's acceptance: S-0220 met it twice on 2026-10-06, before S-0278 was accepted. S-0278 built the fix ([ADR-0098](../../../design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)): when a rebase that `flai stream sync` runs stops on the summary alone, flai regenerates it and continues, and `flai accept` syncs the same way. Its tests in `flai/cmd/stream_sync_test.go` reach main only through the acceptance of another story (`issueStories`).

Add a test there, beside `TestSyncRegeneratesTheIssueSummaryWhenTheRebaseStopsOnItAlone`, that reproduces I-0089 as it happened, reusing `syncProject`, `openSyncStory`, and `recordIssueIn`:

- open one story whose branch records an issue with `flai issue new` and commits it;
- in the main checkout, record a second issue with `flai issue new` and commit it on main, and close a third issue there with `flai issue close`, so that main's summary gains a row and loses one;
- `flai stream sync` the story: it ends clean, with no rebase left in progress, and the story branch's `summary.md` holds the story's issue and main's open issue, without the closed one and without conflict markers;
- move the story to review and `flai accept` it: main's `summary.md` holds both open issues.

If the test fails against the code as it is, the cause is not what ADR-0098 covers: stop, record what it shows in the narrative, and add a task that fixes it in `flai/internal/storygit/sync.go`, widening the story's touches.

It waits for no task.

## Done when

- The test is in `flai/cmd/stream_sync_test.go`, names I-0089 in its comment, and passes with `flai test flai/cmd/stream_sync_test.go`.
- The test fails with the regeneration of ADR-0098 left out, which shows it reproduces the conflict; checked once by hand and said in the narrative.

## Notes

Drafted by the planner for S-0315.
