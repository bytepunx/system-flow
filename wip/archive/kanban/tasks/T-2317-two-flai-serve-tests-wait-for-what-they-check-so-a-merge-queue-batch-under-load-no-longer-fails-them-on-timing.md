---
id: T-2317
type: task
nature: feature
title: Two flai serve tests wait for what they check, so a merge queue batch under load no longer fails them on timing
status: done
parent: S-0507
owner: alex
created: 2026-10-28T23:43:37Z
updated: 2026-10-28T23:45:20Z
transitions:
  - to: ready
    at: 2026-10-28T23:43:43Z
    by: agent-S-0507
  - to: in-progress
    at: 2026-10-28T23:43:45Z
    by: agent-S-0507
  - to: done
    at: 2026-10-28T23:45:20Z
    by: agent-S-0507
stream: S-0507
tags: []
touches: [flai/internal/serve/agents_test.go, flai/internal/serve/forgemerge_test.go]
usage:
  source: log
  seconds: 95
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 5729
      cache_read: 1042735
      cache_write: 28916
      cost: 0.5116
---

# T-2317 Two flai serve tests wait for what they check, so a merge queue batch under load no longer fails them on timing

## Work

Merge queue batch 20261028-232929 failed S-0507 at integration (TH-0724) on two `flai/internal/serve` tests that pass alone and failed in no earlier batch. Each checks something before it has happened. `TestAStorySentBackIsStartedAfreshOrBeforeARestart` asks for exactly one journal entry once the stub agent ends, but flai serve rightly journals a second one when the agent ends with its story in progress. `TestRequestMergeIsAnsweredByFlaiServesMergePass` reads the merge's log line after waiting only for its journal entry, and the log line is written later. Make the first check that its one start is the only start, and the second wait for the log line. Neither changes what flai does.

## Done when

Both tests check what they checked before without depending on timing, and `flai test flai/internal/serve` passes.

## Notes
