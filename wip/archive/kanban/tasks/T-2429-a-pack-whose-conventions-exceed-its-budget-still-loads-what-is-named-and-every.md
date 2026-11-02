---
id: T-2429
type: task
nature: feature
title: A pack whose conventions exceed its budget still loads what is named and every brief, ranking nothing
status: done
parent: S-0563
owner: alex
created: 2026-11-02T04:27:41Z
updated: 2026-11-02T04:34:44Z
transitions:
  - to: ready
    at: 2026-11-02T04:27:58Z
    by: agent-S-0563
  - to: in-progress
    at: 2026-11-02T04:28:00Z
    by: agent-S-0563
  - to: done
    at: 2026-11-02T04:34:44Z
    by: agent-S-0563
stream: S-0563
tags: []
touches: [flai/internal/context/context.go, flai/internal/context/brief_test.go, flai/cmd/prime_test.go]
usage:
  source: log
  seconds: 404
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 6057
      cache_read: 1895709
      cache_write: 39567
      cost: 0.8169
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 22601
          cache_write: 731
          cost: 0.0054
---
# T-2429 A pack whose conventions exceed its budget still loads what is named and every brief, ranking nothing

## Work

In `addDesign` (`flai/internal/context/context.go`), a pack whose conventions exceed the budget keeps `exceeded: conventions` but loads what is named and every brief, as TH-0805 (option 1) chose, and ranks nothing; the header says so. It waits for no task.

## Done when

`TestTheHeaderSaysWhatAloneExceedsTheBudget` asserts the briefs and named items are loaded when the conventions alone exceed, and fails on main; `flai test` on the package passes, and `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` passes.

## Notes
