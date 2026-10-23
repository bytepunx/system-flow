---
id: T-2090
type: task
nature: feature
title: The digest step's over-budget rule follows TH-0624, and agent-context.md records the packs of three archived stories with and without the step
status: done
parent: S-0430
owner: alex
created: 2026-10-23T19:48:53Z
updated: 2026-10-23T20:40:34Z
transitions:
  - to: ready
    at: 2026-10-23T20:31:36Z
    by: agent-S-0430
  - to: in-progress
    at: 2026-10-23T20:31:37Z
    by: agent-S-0430
  - to: done
    at: 2026-10-23T20:40:34Z
    by: agent-S-0430
stream: S-0430
tags: []
touches: [design/system/agent-context.md, design/adrs/0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md, flai/internal/context/context.go, flai/internal/context/digest.go, flai/internal/context/digest_test.go, flai/internal/context/strategic_test.go, flai/cmd/prime.go, docs/users/flai-reference.md, flai/internal/mcpserver/folder.go, docs/users/flai.md, design/system/flai-cli.md]
after: [T-1822]
usage:
  source: log
  seconds: 537
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 21245
      cache_read: 3984031
      cache_write: 107246
      cost: 1.8753
---
# T-2090 The digest step's over-budget rule follows TH-0624, and agent-context.md records the packs of three archived stories with and without the step

## Work

Every pack on this repository is over its budget before the digest step (about 145 KB of conventions against 80 KB), so under ADR-0194 as first written the step loads nothing here. TH-0624 asks the operator whether the step should load within its share on a pack already over budget. Apply the answer: with A, amend ADR-0194 on the branch, change `addDesign` and `ForStrategic`'s step to run whatever `Exceeded` says, within the share, test both cases, and change the over-budget sentence in `prime.go`, `folder.go`, `flai.md`, `flai-cli.md`, and `agent-context.md`; with B, keep the code. Then measure three archived stories' packs with and without the step and record them under `agent-context.md` § The digest › Measured, with the budget each was measured at.

It waits for T-1822, which wrote the section it fills.

## Done when

- [ ] ADR-0194 and the code follow TH-0624's answer, with a test of the over-budget case.
- [ ] `agent-context.md` records the pack sizes of three archived stories with and without the step, and that no pack goes over its budget by the step alone.

## Notes

Split from T-1822 on 2026-10-23, when the measurement showed every pack here over budget before the step.
