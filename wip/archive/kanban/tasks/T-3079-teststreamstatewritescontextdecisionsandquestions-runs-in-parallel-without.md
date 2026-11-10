---
id: T-3079
type: task
nature: feature
title: TestStreamStateWritesContextDecisionsAndQuestions runs in parallel without setting the process, as TestNoTestSetsTheProcess requires
status: done
parent: S-0665
owner: alex
created: 2026-11-10T11:03:52Z
updated: 2026-11-10T11:08:19Z
transitions:
  - to: ready
    at: 2026-11-10T11:04:25Z
    by: agent-S-0665
  - to: in-progress
    at: 2026-11-10T11:04:26Z
    by: agent-S-0665
  - to: done
    at: 2026-11-10T11:08:19Z
    by: agent-S-0665
stream: S-0665
tags: []
touches: [flai/internal/mcpserver/stream_sections_test.go]
usage:
  source: log
  seconds: 233
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 872
      cache_read: 684979
      cache_write: 17789
      cost: 0.2627
---
# T-3079 TestStreamStateWritesContextDecisionsAndQuestions runs in parallel without setting the process, as TestNoTestSetsTheProcess requires

## Work

Main is red, and so is this branch once synced: S-0668 added `TestStreamStateWritesContextDecisionsAndQuestions` in `flai/internal/mcpserver/stream_sections_test.go` with `t.Setenv("FLAI_SESSION", "s9")` and no `t.Parallel()`, and S-0684 added `TestNoTestSetsTheProcess`, which fails on such a test (I-0234). The two were accepted apart. The go-test tier runs mcpserver's tests for any change to `flai/internal/serve`, which mcpserver imports, so S-0665's close-out cannot pass while it stands.

`streamFixture` already gives the server the session `s9` through `Options`, so the `t.Setenv` is not needed: replace it with `t.Parallel()`.

It waits for nothing and shares no path with S-0665's other tasks.

## Done when

- The test calls `t.Parallel()` and sets no environment variable.
- `flai test` on `flai/internal/mcpserver` passes, `TestNoTestSetsTheProcess` and the changed test included.

## Notes
