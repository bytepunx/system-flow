---
id: T-2089
type: task
nature: feature
title: TestSyncListsPathsChangedOutsideTheClaim's fixture names the file its task touches, not the folder ADR-0188 refuses, so the merge queue's integration tier passes
status: done
parent: S-0429
owner: alex
created: 2026-10-23T18:41:10Z
updated: 2026-10-23T18:41:50Z
transitions:
  - to: ready
    at: 2026-10-23T18:41:19Z
    by: agent-S-0429
  - to: in-progress
    at: 2026-10-23T18:41:20Z
    by: agent-S-0429
  - to: done
    at: 2026-10-23T18:41:50Z
    by: agent-S-0429
stream: S-0429
tags: []
touches: [flai/cmd/stream_sync_test.go]
usage:
  source: log
  seconds: 30
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 5067
      cache_read: 830563
      cache_write: 26898
      cost: 0.4291
---

# T-2089 TestSyncListsPathsChangedOutsideTheClaim's fixture names the file its task touches, not the folder ADR-0188 refuses, so the merge queue's integration tier passes

## Work

`TestSyncListsPathsChangedOutsideTheClaim` in `flai/cmd/stream_sync_test.go` gives its task the folder touch `design`, which `flai task new` refuses since S-0468 (ADR-0188), so the test fails on main (I-0200) and failed S-0429's merge queue batch at integration (TH-0623). Name the file the task changes, `design/note.md`, instead: the test still shows that a path inside a task's touches is not reported as outside the claim. Close I-0200 with the fix.

Added after review, outside the plan's layers: it waits for nothing and no other task touches the file.

## Done when

The test passes in the integration tier, and I-0200 is closed naming the fix.

## Notes
